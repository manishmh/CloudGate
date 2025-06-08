package handlers

import (
	"net/http"
	"strconv"

	"cloudgate-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandlers contains user-related HTTP handlers
type UserHandlers struct {
	userService    *services.UserService
	sessionService *services.SessionService
}

// NewUserHandlers creates new user handlers
func NewUserHandlers(userService *services.UserService, sessionService *services.SessionService) *UserHandlers {
	return &UserHandlers{
		userService:    userService,
		sessionService: sessionService,
	}
}

// GetProfile retrieves the current user's profile
func (h *UserHandlers) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userService.GetUserByID(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateProfile updates the current user's profile
func (h *UserHandlers) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		FirstName         string `json:"first_name" binding:"required"`
		LastName          string `json:"last_name" binding:"required"`
		ProfilePictureURL string `json:"profile_picture_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.UpdateUserProfile(userID.(uuid.UUID), req.FirstName, req.LastName, req.ProfilePictureURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Get updated user
	user, err := h.userService.GetUserByID(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// SendEmailVerification sends an email verification token
func (h *UserHandlers) SendEmailVerification(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verification, err := h.userService.CreateEmailVerification(userID.(uuid.UUID), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create verification token"})
		return
	}

	// TODO: Send email with verification link
	// For now, we'll return the token (in production, this should be sent via email)
	c.JSON(http.StatusOK, gin.H{
		"message": "Verification email sent",
		"token":   verification.Token, // Remove this in production
	})
}

// VerifyEmail verifies an email using the verification token
func (h *UserHandlers) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	user, err := h.userService.VerifyEmail(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
		"user":    user,
	})
}

// GetAuditLogs retrieves audit logs for the current user
func (h *UserHandlers) GetAuditLogs(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	logs, err := h.userService.GetUserAuditLogs(userID.(uuid.UUID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"audit_logs": logs,
	})
}

// GetSessions retrieves active sessions for the current user
func (h *UserHandlers) GetSessions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	sessions, err := h.sessionService.GetUserSessions(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
	})
}

// InvalidateSession invalidates a specific session
func (h *UserHandlers) InvalidateSession(c *gin.Context) {
	sessionToken := c.Param("token")
	if sessionToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session token is required"})
		return
	}

	err := h.sessionService.InvalidateSession(sessionToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Session invalidated successfully",
	})
}

// InvalidateAllSessions invalidates all sessions for the current user
func (h *UserHandlers) InvalidateAllSessions(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.sessionService.InvalidateAllUserSessions(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All sessions invalidated successfully",
	})
}

// DeactivateAccount deactivates the current user's account
func (h *UserHandlers) DeactivateAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.userService.DeactivateUser(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deactivated successfully",
	})
}
