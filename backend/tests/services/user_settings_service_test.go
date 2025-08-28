package services_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"cloudgate-backend/internal/models"
	"cloudgate-backend/internal/services"
)

// setupUserSettingsTestDB initializes an in-memory SQLite database for user settings service tests
func setupUserSettingsTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to test database")

	// Auto-migrate all required schemas
	err = db.AutoMigrate(
		&models.User{},
		&models.UserSettings{},
	)
	require.NoError(t, err, "Failed to migrate database schema")

	return db
}

// setupTestUserSettingsService sets up a test user settings service with database
func setupTestUserSettingsService(t *testing.T) (*services.UserSettingsService, *gorm.DB, *models.User) {
	db := setupUserSettingsTestDB(t)

	// Create test user
	kc := "test-keycloak-id"
	user := &models.User{
		ID:         uuid.New(),
		KeycloakID: &kc,
		Email:      "test@example.com",
		Username:   "testuser",
		FirstName:  "Test",
		LastName:   "User",
		IsActive:   true,
	}
	err := db.Create(user).Error
	require.NoError(t, err)

	userSettingsService := services.NewUserSettingsService(db)
	return userSettingsService, db, user
}

func TestUserSettingsService_GetUserSettings(t *testing.T) {
	service, db, user := setupTestUserSettingsService(t)

	t.Run("should create default settings for new user", func(t *testing.T) {
		settings, err := service.GetUserSettings(user.ID)

		assert.NoError(t, err)
		assert.NotNil(t, settings)
		assert.Equal(t, user.ID, settings.UserID)

		// Verify default values
		assert.Equal(t, "en", settings.Language)
		assert.Equal(t, "America/New_York", settings.Timezone)
		assert.Equal(t, "MM/DD/YYYY", settings.DateFormat)
		assert.True(t, settings.EmailNotifications)
		assert.False(t, settings.PushNotifications)
		assert.True(t, settings.SecurityAlerts)
		assert.Equal(t, 30, settings.SessionTimeout)
		assert.Equal(t, 90, settings.PasswordExpiryDays)
		assert.Equal(t, "dashboard", settings.DefaultView)
		assert.Equal(t, 10, settings.ItemsPerPage)
		assert.Equal(t, 1000, settings.MaxAPICalls)

		// Verify settings were saved to database
		var dbSettings models.UserSettings
		err = db.Where("user_id = ?", user.ID).First(&dbSettings).Error
		assert.NoError(t, err)
		assert.Equal(t, settings.ID, dbSettings.ID)
	})

	t.Run("should return existing settings", func(t *testing.T) {
		// Create custom settings
		customSettings := &models.UserSettings{
			UserID:             user.ID,
			Language:           "es",
			Timezone:           "Europe/Madrid",
			EmailNotifications: false,
			SessionTimeout:     60,
		}
		err := db.Create(customSettings).Error
		assert.NoError(t, err)

		// Create new user for this test
		kc2 := "test-keycloak-id-2"
		newUser := &models.User{
			ID:         uuid.New(),
			KeycloakID: &kc2,
			Email:      "test2@example.com",
			Username:   "testuser2",
		}
		err = db.Create(newUser).Error
		assert.NoError(t, err)

		customSettings.UserID = newUser.ID
		customSettings.ID = uuid.New()
		err = db.Create(customSettings).Error
		assert.NoError(t, err)

		// Get settings
		settings, err := service.GetUserSettings(newUser.ID)

		assert.NoError(t, err)
		assert.Equal(t, "es", settings.Language)
		assert.Equal(t, "Europe/Madrid", settings.Timezone)
		assert.False(t, settings.EmailNotifications)
		assert.Equal(t, 60, settings.SessionTimeout)
	})
}

func TestUserSettingsService_CreateDefaultSettings(t *testing.T) {
	service, db, _ := setupTestUserSettingsService(t)

	t.Run("should create default settings", func(t *testing.T) {
		kc3 := "test-keycloak-id-3"
		newUser := &models.User{
			ID:         uuid.New(),
			KeycloakID: &kc3,
			Email:      "test3@example.com",
			Username:   "testuser3",
		}
		err := db.Create(newUser).Error
		assert.NoError(t, err)

		settings, err := service.CreateDefaultSettings(newUser.ID)

		assert.NoError(t, err)
		assert.NotNil(t, settings)
		assert.Equal(t, newUser.ID, settings.UserID)
		assert.NotEqual(t, uuid.Nil, settings.ID)

		// Verify default values match the model defaults
		defaultSettings := models.GetDefaultSettings()
		assert.Equal(t, defaultSettings.Language, settings.Language)
		assert.Equal(t, defaultSettings.Timezone, settings.Timezone)
		assert.Equal(t, defaultSettings.EmailNotifications, settings.EmailNotifications)
		assert.Equal(t, defaultSettings.SessionTimeout, settings.SessionTimeout)
	})

	t.Run("should save settings to database", func(t *testing.T) {
		kc4 := "test-keycloak-id-4"
		newUser := &models.User{
			ID:         uuid.New(),
			KeycloakID: &kc4,
			Email:      "test4@example.com",
			Username:   "testuser4",
		}
		err := db.Create(newUser).Error
		assert.NoError(t, err)

		settings, err := service.CreateDefaultSettings(newUser.ID)
		assert.NoError(t, err)

		// Verify in database
		var dbSettings models.UserSettings
		err = db.Where("user_id = ?", newUser.ID).First(&dbSettings).Error
		assert.NoError(t, err)
		assert.Equal(t, settings.ID, dbSettings.ID)
		assert.Equal(t, settings.Language, dbSettings.Language)
	})
}

func TestUserSettingsService_UpdateUserSettings(t *testing.T) {
	service, db, user := setupTestUserSettingsService(t)

	t.Run("should update existing settings", func(t *testing.T) {
		// First get/create settings
		_, err := service.GetUserSettings(user.ID)
		assert.NoError(t, err)

		// Update settings
		updates := map[string]interface{}{
			"language":            "fr",
			"timezone":            "Europe/Paris",
			"email_notifications": false,
			"session_timeout":     45,
			"items_per_page":      20,
		}

		updatedSettings, err := service.UpdateUserSettings(user.ID, updates)

		assert.NoError(t, err)
		assert.Equal(t, "fr", updatedSettings.Language)
		assert.Equal(t, "Europe/Paris", updatedSettings.Timezone)
		assert.False(t, updatedSettings.EmailNotifications)
		assert.Equal(t, 45, updatedSettings.SessionTimeout)
		assert.Equal(t, 20, updatedSettings.ItemsPerPage)

		// Verify changes were saved to database
		var dbSettings models.UserSettings
		err = db.Where("user_id = ?", user.ID).First(&dbSettings).Error
		assert.NoError(t, err)
		assert.Equal(t, "fr", dbSettings.Language)
		assert.Equal(t, "Europe/Paris", dbSettings.Timezone)
		assert.False(t, dbSettings.EmailNotifications)
	})

	t.Run("should create default settings if none exist and then update", func(t *testing.T) {
		kc5 := "test-keycloak-id-5"
		newUser := &models.User{
			ID:         uuid.New(),
			KeycloakID: &kc5,
			Email:      "test5@example.com",
			Username:   "testuser5",
		}
		err := db.Create(newUser).Error
		assert.NoError(t, err)

		// Update settings for user with no existing settings
		updates := map[string]interface{}{
			"language":        "de",
			"session_timeout": 120,
		}

		updatedSettings, err := service.UpdateUserSettings(newUser.ID, updates)

		assert.NoError(t, err)
		assert.Equal(t, "de", updatedSettings.Language)
		assert.Equal(t, 120, updatedSettings.SessionTimeout)

		// Should have other default values
		assert.Equal(t, "America/New_York", updatedSettings.Timezone) // Default
		assert.True(t, updatedSettings.EmailNotifications)            // Default
	})

	t.Run("should handle partial updates", func(t *testing.T) {
		// Get existing settings
		existingSettings, err := service.GetUserSettings(user.ID)
		assert.NoError(t, err)

		originalLanguage := existingSettings.Language
		originalTimeout := existingSettings.SessionTimeout

		// Update only one field
		updates := map[string]interface{}{
			"push_notifications": true,
		}

		updatedSettings, err := service.UpdateUserSettings(user.ID, updates)

		assert.NoError(t, err)
		assert.True(t, updatedSettings.PushNotifications)

		// Other fields should remain unchanged
		assert.Equal(t, originalLanguage, updatedSettings.Language)
		assert.Equal(t, originalTimeout, updatedSettings.SessionTimeout)
	})
}

func TestUserSettingsService_UpdateSetting(t *testing.T) {
	service, _, user := setupTestUserSettingsService(t)

	t.Run("should update single setting", func(t *testing.T) {
		// First get/create settings
		_, err := service.GetUserSettings(user.ID)
		assert.NoError(t, err)

		// Update single setting
		updatedSettings, err := service.UpdateSetting(user.ID, "language", "it")

		assert.NoError(t, err)
		assert.Equal(t, "it", updatedSettings.Language)
	})

	t.Run("should update boolean setting", func(t *testing.T) {
		updatedSettings, err := service.UpdateSetting(user.ID, "two_factor_enabled", true)

		assert.NoError(t, err)
		assert.True(t, updatedSettings.TwoFactorEnabled)
	})

	t.Run("should update numeric setting", func(t *testing.T) {
		updatedSettings, err := service.UpdateSetting(user.ID, "max_api_calls", 2000)

		assert.NoError(t, err)
		assert.Equal(t, 2000, updatedSettings.MaxAPICalls)
	})
}

func TestUserSettingsService_ResetUserSettings(t *testing.T) {
	service, db, user := setupTestUserSettingsService(t)

	t.Run("should reset settings to defaults", func(t *testing.T) {
		// First create and modify settings
		_, err := service.GetUserSettings(user.ID)
		assert.NoError(t, err)

		// Update some settings
		updates := map[string]interface{}{
			"language":            "zh",
			"timezone":            "Asia/Shanghai",
			"email_notifications": false,
			"session_timeout":     180,
			"items_per_page":      50,
		}
		_, err = service.UpdateUserSettings(user.ID, updates)
		assert.NoError(t, err)

		// Reset settings
		resetSettings, err := service.ResetUserSettings(user.ID)

		assert.NoError(t, err)
		assert.NotNil(t, resetSettings)

		// Verify settings are back to defaults
		defaultSettings := models.GetDefaultSettings()
		assert.Equal(t, defaultSettings.Language, resetSettings.Language)
		assert.Equal(t, defaultSettings.Timezone, resetSettings.Timezone)
		assert.Equal(t, defaultSettings.EmailNotifications, resetSettings.EmailNotifications)
		assert.Equal(t, defaultSettings.SessionTimeout, resetSettings.SessionTimeout)
		assert.Equal(t, defaultSettings.ItemsPerPage, resetSettings.ItemsPerPage)

		// Verify in database
		var dbSettings models.UserSettings
		err = db.Where("user_id = ?", user.ID).First(&dbSettings).Error
		assert.NoError(t, err)
		assert.Equal(t, defaultSettings.Language, dbSettings.Language)
		assert.Equal(t, defaultSettings.Timezone, dbSettings.Timezone)
	})

	t.Run("should handle resetting non-existent settings", func(t *testing.T) {
		kc6 := "test-keycloak-id-6"
		newUser := &models.User{
			ID:         uuid.New(),
			KeycloakID: &kc6,
			Email:      "test6@example.com",
			Username:   "testuser6",
		}
		err := db.Create(newUser).Error
		assert.NoError(t, err)

		// Reset settings for user with no existing settings
		resetSettings, err := service.ResetUserSettings(newUser.ID)

		assert.NoError(t, err)
		assert.NotNil(t, resetSettings)
		assert.Equal(t, newUser.ID, resetSettings.UserID)

		// Should have default values
		defaultSettings := models.GetDefaultSettings()
		assert.Equal(t, defaultSettings.Language, resetSettings.Language)
		assert.Equal(t, defaultSettings.EmailNotifications, resetSettings.EmailNotifications)
	})
}
