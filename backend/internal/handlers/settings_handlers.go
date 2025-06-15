package handlers

import (
	"net/http"

	"cloudgate-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SettingsHandlers contains settings-related HTTP handlers
type SettingsHandlers struct {
	settingsService *services.UserSettingsService
}

// NewSettingsHandlers creates new settings handlers
func NewSettingsHandlers(settingsService *services.UserSettingsService) *SettingsHandlers {
	return &SettingsHandlers{
		settingsService: settingsService,
	}
}

// GetUserSettings retrieves the current user's settings
func (h *SettingsHandlers) GetUserSettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	settings, err := h.settingsService.GetUserSettings(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"settings": settings,
	})
}

// UpdateUserSettings updates the current user's settings
func (h *SettingsHandlers) UpdateUserSettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings, err := h.settingsService.UpdateUserSettings(userID.(uuid.UUID), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Settings updated successfully",
		"settings": settings,
	})
}

// UpdateSingleSetting updates a single setting
func (h *SettingsHandlers) UpdateSingleSetting(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		Key   string      `json:"key" binding:"required"`
		Value interface{} `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings, err := h.settingsService.UpdateSetting(userID.(uuid.UUID), req.Key, req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update setting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Setting updated successfully",
		"settings": settings,
	})
}

// ResetUserSettings resets user settings to defaults
func (h *SettingsHandlers) ResetUserSettings(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	settings, err := h.settingsService.ResetUserSettings(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset user settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Settings reset to defaults successfully",
		"settings": settings,
	})
}
