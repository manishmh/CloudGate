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

// SessionService handles session-related operations
type SessionService struct {
	db                *gorm.DB
	disableCleanupJob bool
}

// NewSessionService creates a new session service
func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db, disableCleanupJob: false}
}

// NewSessionServiceForTesting creates a new session service for testing (disables cleanup job)
func NewSessionServiceForTesting(db *gorm.DB) *SessionService {
	return &SessionService{db: db, disableCleanupJob: true}
}

// CreateSession creates a new session for a user
func (s *SessionService) CreateSession(userID uuid.UUID, ipAddress, userAgent string) (*models.Session, error) {
	// Generate session token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}
	sessionToken := hex.EncodeToString(tokenBytes)

	// Create session
	session := models.Session{
		UserID:       userID,
		SessionToken: sessionToken,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		ExpiresAt:    time.Now().Add(24 * time.Hour), // 24 hours default
		IsActive:     true,
	}

	if err := s.db.Create(&session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Clean up old sessions for this user (keep only last 5)
	if !s.disableCleanupJob {
		go s.cleanupOldSessions(userID)
	}

	return &session, nil
}

// GetSessionByToken retrieves a session by token
func (s *SessionService) GetSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	err := s.db.Preload("User").Where("session_token = ? AND is_active = ?", token, true).First(&session).Error
	if err != nil {
		return nil, err
	}

	// Check if session is expired
	if session.IsExpired() {
		// Mark session as inactive
		s.db.Model(&session).Update("is_active", false)
		return nil, fmt.Errorf("session expired")
	}

	return &session, nil
}

// ValidateSession validates a session and returns the user
func (s *SessionService) ValidateSession(token string) (*models.User, error) {
	session, err := s.GetSessionByToken(token)
	if err != nil {
		return nil, err
	}

	// Update last activity
	s.db.Model(session).Update("updated_at", time.Now())

	return &session.User, nil
}

// RefreshSession extends the session expiry
func (s *SessionService) RefreshSession(token string) (*models.Session, error) {
	session, err := s.GetSessionByToken(token)
	if err != nil {
		return nil, err
	}

	// Extend expiry by 24 hours
	session.ExpiresAt = time.Now().Add(24 * time.Hour)

	if err := s.db.Save(session).Error; err != nil {
		return nil, fmt.Errorf("failed to refresh session: %w", err)
	}

	return session, nil
}

// InvalidateSession invalidates a session
func (s *SessionService) InvalidateSession(token string) error {
	err := s.db.Model(&models.Session{}).Where("session_token = ?", token).Update("is_active", false).Error
	if err != nil {
		return fmt.Errorf("failed to invalidate session: %w", err)
	}
	return nil
}

// InvalidateAllUserSessions invalidates all sessions for a user
func (s *SessionService) InvalidateAllUserSessions(userID uuid.UUID) error {
	err := s.db.Model(&models.Session{}).Where("user_id = ?", userID).Update("is_active", false).Error
	if err != nil {
		return fmt.Errorf("failed to invalidate user sessions: %w", err)
	}
	return nil
}

// GetUserSessions retrieves all active sessions for a user
func (s *SessionService) GetUserSessions(userID uuid.UUID) ([]models.Session, error) {
	var sessions []models.Session
	err := s.db.Where("user_id = ? AND is_active = ?", userID, true).Order("created_at DESC").Find(&sessions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}
	return sessions, nil
}

// CleanupExpiredSessions removes expired sessions from the database
func (s *SessionService) CleanupExpiredSessions() error {
	// Delete sessions that expired more than 7 days ago
	cutoff := time.Now().Add(-7 * 24 * time.Hour)

	result := s.db.Where("expires_at < ? OR (is_active = ? AND updated_at < ?)",
		time.Now(), false, cutoff).Delete(&models.Session{})

	if result.Error != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		fmt.Printf("Cleaned up %d expired sessions\n", result.RowsAffected)
	}

	return nil
}

// cleanupOldSessions keeps only the latest 5 sessions for a user
func (s *SessionService) cleanupOldSessions(userID uuid.UUID) {
	var sessions []models.Session

	// Get all sessions for user, ordered by creation date (newest first)
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&sessions).Error
	if err != nil {
		fmt.Printf("Failed to get user sessions for cleanup: %v\n", err)
		return
	}

	// If more than 5 sessions, deactivate the older ones
	if len(sessions) > 5 {
		var oldSessionIDs []uuid.UUID
		for i := 5; i < len(sessions); i++ {
			oldSessionIDs = append(oldSessionIDs, sessions[i].ID)
		}

		err := s.db.Model(&models.Session{}).Where("id IN ?", oldSessionIDs).Update("is_active", false).Error
		if err != nil {
			fmt.Printf("Failed to cleanup old sessions: %v\n", err)
		}
	}
}

// GetSessionStats returns session statistics
func (s *SessionService) GetSessionStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total active sessions
	var activeCount int64
	if err := s.db.Model(&models.Session{}).Where("is_active = ? AND expires_at > ?", true, time.Now()).Count(&activeCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count active sessions: %w", err)
	}
	stats["active_sessions"] = activeCount

	// Total expired sessions
	var expiredCount int64
	if err := s.db.Model(&models.Session{}).Where("expires_at <= ?", time.Now()).Count(&expiredCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count expired sessions: %w", err)
	}
	stats["expired_sessions"] = expiredCount

	// Sessions created today
	today := time.Now().Truncate(24 * time.Hour)
	var todayCount int64
	if err := s.db.Model(&models.Session{}).Where("created_at >= ?", today).Count(&todayCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count today's sessions: %w", err)
	}
	stats["sessions_today"] = todayCount

	return stats, nil
}
