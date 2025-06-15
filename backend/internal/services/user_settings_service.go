package services

import (
	"fmt"

	"cloudgate-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSettingsService handles user settings operations
type UserSettingsService struct {
	db *gorm.DB
}

// NewUserSettingsService creates a new user settings service
func NewUserSettingsService(db *gorm.DB) *UserSettingsService {
	return &UserSettingsService{db: db}
}

// GetUserSettings retrieves user settings by user ID
func (s *UserSettingsService) GetUserSettings(userID uuid.UUID) (*models.UserSettings, error) {
	var settings models.UserSettings
	err := s.db.Where("user_id = ?", userID).First(&settings).Error

	if err == gorm.ErrRecordNotFound {
		// Create default settings if none exist
		return s.CreateDefaultSettings(userID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user settings: %w", err)
	}

	return &settings, nil
}

// CreateDefaultSettings creates default settings for a user
func (s *UserSettingsService) CreateDefaultSettings(userID uuid.UUID) (*models.UserSettings, error) {
	settings := models.GetDefaultSettings()
	settings.UserID = userID

	if err := s.db.Create(settings).Error; err != nil {
		return nil, fmt.Errorf("failed to create default settings: %w", err)
	}

	return settings, nil
}

// UpdateUserSettings updates user settings
func (s *UserSettingsService) UpdateUserSettings(userID uuid.UUID, updates map[string]interface{}) (*models.UserSettings, error) {
	var settings models.UserSettings

	// Get existing settings or create default ones
	err := s.db.Where("user_id = ?", userID).First(&settings).Error
	if err == gorm.ErrRecordNotFound {
		defaultSettings, err := s.CreateDefaultSettings(userID)
		if err != nil {
			return nil, err
		}
		settings = *defaultSettings
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user settings: %w", err)
	}

	// Update settings
	if err := s.db.Model(&settings).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update user settings: %w", err)
	}

	// Get updated settings
	if err := s.db.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated settings: %w", err)
	}

	return &settings, nil
}

// UpdateSetting updates a single setting
func (s *UserSettingsService) UpdateSetting(userID uuid.UUID, key string, value interface{}) (*models.UserSettings, error) {
	updates := map[string]interface{}{
		key: value,
	}
	return s.UpdateUserSettings(userID, updates)
}

// ResetUserSettings resets user settings to defaults
func (s *UserSettingsService) ResetUserSettings(userID uuid.UUID) (*models.UserSettings, error) {
	// Delete existing settings
	if err := s.db.Where("user_id = ?", userID).Delete(&models.UserSettings{}).Error; err != nil {
		return nil, fmt.Errorf("failed to delete existing settings: %w", err)
	}

	// Create new default settings
	return s.CreateDefaultSettings(userID)
}
