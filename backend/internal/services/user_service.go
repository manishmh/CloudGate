package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"cloudgate-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserService handles user-related operations
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateOrUpdateUser creates a new user or updates existing user from Keycloak data
func (s *UserService) CreateOrUpdateUser(keycloakID, email, username, firstName, lastName string) (*models.User, error) {
	var user models.User

	// Try to find existing user by Keycloak ID
	err := s.db.Where("keycloak_id = ?", keycloakID).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		// Create new user
		user = models.User{
			KeycloakID: keycloakID,
			Email:      email,
			Username:   username,
			FirstName:  firstName,
			LastName:   lastName,
			IsActive:   true,
		}

		if err := s.db.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		// Log user creation
		s.LogAudit(user.ID, "user.created", "user", user.ID.String(), "", "", "User account created")

	} else if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	} else {
		// Update existing user
		user.Email = email
		user.Username = username
		user.FirstName = firstName
		user.LastName = lastName
		user.LastLoginAt = &time.Time{}
		*user.LastLoginAt = time.Now()

		if err := s.db.Save(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		// Log user login
		s.LogAudit(user.ID, "user.login", "user", user.ID.String(), "", "", "User logged in")
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.Where("id = ? AND is_active = ?", userID, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByKeycloakID retrieves a user by Keycloak ID
func (s *UserService) GetUserByKeycloakID(keycloakID string) (*models.User, error) {
	var user models.User
	err := s.db.Where("keycloak_id = ? AND is_active = ?", keycloakID, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserProfile updates user profile information
func (s *UserService) UpdateUserProfile(userID uuid.UUID, firstName, lastName, profilePictureURL string) error {
	updates := map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
	}

	if profilePictureURL != "" {
		updates["profile_picture_url"] = profilePictureURL
	}

	err := s.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update user profile: %w", err)
	}

	// Log profile update
	s.LogAudit(userID, "user.profile_updated", "user", userID.String(), "", "", "User profile updated")

	return nil
}

// CreateEmailVerification creates a new email verification token
func (s *UserService) CreateEmailVerification(userID uuid.UUID, email string) (*models.EmailVerification, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Create verification record
	verification := models.EmailVerification{
		UserID:    userID,
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours expiry
	}

	if err := s.db.Create(&verification).Error; err != nil {
		return nil, fmt.Errorf("failed to create email verification: %w", err)
	}

	// Log verification creation
	s.LogAudit(userID, "email.verification_created", "email_verification", verification.ID.String(), "", "", "Email verification token created")

	return &verification, nil
}

// VerifyEmail verifies an email using the verification token
func (s *UserService) VerifyEmail(token string) (*models.User, error) {
	var verification models.EmailVerification

	// Find verification token
	err := s.db.Where("token = ? AND used_at IS NULL", token).First(&verification).Error
	if err != nil {
		return nil, fmt.Errorf("invalid or expired verification token")
	}

	// Check if token is expired
	if verification.IsExpired() {
		return nil, fmt.Errorf("verification token has expired")
	}

	// Mark token as used
	now := time.Now()
	verification.UsedAt = &now
	if err := s.db.Save(&verification).Error; err != nil {
		return nil, fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Update user email verification status
	var user models.User
	err = s.db.Model(&user).Where("id = ?", verification.UserID).Updates(map[string]interface{}{
		"email_verified":    true,
		"email_verified_at": now,
	}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update user email verification: %w", err)
	}

	// Get updated user
	if err := s.db.Where("id = ?", verification.UserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	// Log email verification
	s.LogAudit(user.ID, "email.verified", "user", user.ID.String(), "", "", "Email address verified")

	return &user, nil
}

// LogAudit creates an audit log entry
func (s *UserService) LogAudit(userID uuid.UUID, action, resource, resourceID, ipAddress, userAgent, details string) {
	auditLog := models.AuditLog{
		UserID:     &userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Details:    details,
		Status:     "success",
	}

	// Don't fail the main operation if audit logging fails
	if err := s.db.Create(&auditLog).Error; err != nil {
		fmt.Printf("Failed to create audit log: %v\n", err)
	}
}

// GetUserAuditLogs retrieves audit logs for a user
func (s *UserService) GetUserAuditLogs(userID uuid.UUID, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	query := s.db.Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	return logs, nil
}

// DeactivateUser deactivates a user account
func (s *UserService) DeactivateUser(userID uuid.UUID) error {
	err := s.db.Model(&models.User{}).Where("id = ?", userID).Update("is_active", false).Error
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	// Deactivate all user sessions
	err = s.db.Model(&models.Session{}).Where("user_id = ?", userID).Update("is_active", false).Error
	if err != nil {
		return fmt.Errorf("failed to deactivate user sessions: %w", err)
	}

	// Log user deactivation
	s.LogAudit(userID, "user.deactivated", "user", userID.String(), "", "", "User account deactivated")

	return nil
}

// GetOrCreateDemoUser gets or creates the demo user for development
func (us *UserService) GetOrCreateDemoUser() (*models.User, error) {
	demoUserUUID, _ := uuid.Parse("12345678-1234-1234-1234-123456789012")

	// Try to get existing demo user
	user, err := us.GetUserByID(demoUserUUID)
	if err == nil {
		return user, nil
	}

	// Create demo user if it doesn't exist
	demoUser := &models.User{
		ID:            demoUserUUID,
		KeycloakID:    "demo-keycloak-id",
		Email:         "demo@cloudgate.com",
		Username:      "demo-user",
		FirstName:     "Demo",
		LastName:      "User",
		IsActive:      true,
		EmailVerified: true,
	}

	if err := us.db.Create(demoUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create demo user: %w", err)
	}

	return demoUser, nil
}
