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

// setupMFATestDB initializes an in-memory SQLite database for MFA service tests
func setupMFATestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to test database")

	// Auto-migrate all required schemas
	err = db.AutoMigrate(
		&models.User{},
		&models.MFASetup{},
		&models.BackupCode{},
	)
	require.NoError(t, err, "Failed to migrate database schema")

	return db
}

// setupTestMFAService sets up a test MFA service with database
func setupTestMFAService(t *testing.T) (*gorm.DB, *models.User) {
	db := setupMFATestDB(t)

	// Create test user
	user := &models.User{
		ID:         uuid.New(),
		KeycloakID: "test-keycloak-id",
		Email:      "test@example.com",
		Username:   "testuser",
		FirstName:  "Test",
		LastName:   "User",
		IsActive:   true,
	}
	err := db.Create(user).Error
	require.NoError(t, err)

	return db, user
}

func TestMFAService_StoreMFASetup(t *testing.T) {
	db, user := setupTestMFAService(t)

	// Mock the global DB
	originalDB := services.DB
	services.DB = db
	defer func() { services.DB = originalDB }()

	t.Run("should store MFA setup with backup codes", func(t *testing.T) {
		secret := "JBSWY3DPEHPK3PXP"
		backupCodes := []string{"123456", "789012", "345678"}

		err := services.StoreMFASetup(user.ID.String(), secret, backupCodes)
		assert.NoError(t, err)

		// Verify MFA setup was stored
		var mfaSetup models.MFASetup
		err = db.Where("user_id = ?", user.ID).First(&mfaSetup).Error
		assert.NoError(t, err)
		assert.Equal(t, secret, mfaSetup.Secret)
		assert.False(t, mfaSetup.Enabled) // Should not be enabled initially

		// Verify backup codes were stored
		var backupCodeCount int64
		err = db.Model(&models.BackupCode{}).Where("mfa_setup_id = ?", mfaSetup.ID).Count(&backupCodeCount).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(3), backupCodeCount)
	})

	t.Run("should replace existing MFA setup", func(t *testing.T) {
		// Store initial setup
		err := services.StoreMFASetup(user.ID.String(), "SECRET1", []string{"111111"})
		assert.NoError(t, err)

		// Store new setup (should replace the old one)
		err = services.StoreMFASetup(user.ID.String(), "SECRET2", []string{"222222", "333333"})
		assert.NoError(t, err)

		// Verify only one MFA setup exists
		var mfaSetupCount int64
		err = db.Model(&models.MFASetup{}).Where("user_id = ?", user.ID).Count(&mfaSetupCount).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), mfaSetupCount)

		// Verify it has the new secret
		var mfaSetup models.MFASetup
		err = db.Where("user_id = ?", user.ID).First(&mfaSetup).Error
		assert.NoError(t, err)
		assert.Equal(t, "SECRET2", mfaSetup.Secret)

		// Verify backup codes count
		var backupCodeCount int64
		err = db.Model(&models.BackupCode{}).Where("mfa_setup_id = ?", mfaSetup.ID).Count(&backupCodeCount).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(2), backupCodeCount)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		err := services.StoreMFASetup("invalid-uuid", "secret", []string{"123456"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

func TestMFAService_GetMFASetup(t *testing.T) {
	db, user := setupTestMFAService(t)

	// Mock the global DB
	originalDB := services.DB
	services.DB = db
	defer func() { services.DB = originalDB }()

	t.Run("should return MFA setup for user", func(t *testing.T) {
		// Create MFA setup
		secret := "JBSWY3DPEHPK3PXP"
		err := services.StoreMFASetup(user.ID.String(), secret, []string{"123456"})
		assert.NoError(t, err)

		// Get MFA setup
		mfaSetup, err := services.GetMFASetup(user.ID.String())
		assert.NoError(t, err)
		assert.NotNil(t, mfaSetup)
		assert.Equal(t, secret, mfaSetup.Secret)
		assert.Equal(t, user.ID, mfaSetup.UserID)
		assert.False(t, mfaSetup.Enabled)
	})

	t.Run("should return error for user without MFA setup", func(t *testing.T) {
		// Create another user without MFA setup
		newUser := models.User{
			ID:         uuid.New(),
			KeycloakID: "another-keycloak-id",
			Email:      "another@example.com",
			Username:   "anotheruser",
			IsActive:   true,
		}
		err := db.Create(&newUser).Error
		assert.NoError(t, err)

		_, err = services.GetMFASetup(newUser.ID.String())
		assert.Error(t, err)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		_, err := services.GetMFASetup("invalid-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

func TestMFAService_EnableDisableMFA(t *testing.T) {
	db, user := setupTestMFAService(t)

	// Mock the global DB
	originalDB := services.DB
	services.DB = db
	defer func() { services.DB = originalDB }()

	// Create MFA setup first
	err := services.StoreMFASetup(user.ID.String(), "JBSWY3DPEHPK3PXP", []string{"123456"})
	assert.NoError(t, err)

	t.Run("should enable MFA", func(t *testing.T) {
		err := services.EnableMFA(user.ID.String())
		assert.NoError(t, err)

		// Verify MFA is enabled
		mfaSetup, err := services.GetMFASetup(user.ID.String())
		assert.NoError(t, err)
		assert.True(t, mfaSetup.Enabled)
	})

	t.Run("should disable MFA", func(t *testing.T) {
		err := services.DisableMFA(user.ID.String())
		assert.NoError(t, err)

		// Verify MFA is disabled
		mfaSetup, err := services.GetMFASetup(user.ID.String())
		assert.NoError(t, err)
		assert.False(t, mfaSetup.Enabled)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		err := services.EnableMFA("invalid-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")

		err = services.DisableMFA("invalid-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

func TestMFAService_BackupCodes(t *testing.T) {
	db, user := setupTestMFAService(t)

	// Mock the global DB
	originalDB := services.DB
	services.DB = db
	defer func() { services.DB = originalDB }()

	// Create MFA setup with backup codes
	backupCodes := []string{"123456", "789012", "345678"}
	err := services.StoreMFASetup(user.ID.String(), "JBSWY3DPEHPK3PXP", backupCodes)
	assert.NoError(t, err)

	t.Run("should get backup codes count", func(t *testing.T) {
		count, err := services.GetBackupCodesCount(user.ID.String())
		assert.NoError(t, err)
		assert.Equal(t, 3, count)
	})

	t.Run("should use backup code successfully", func(t *testing.T) {
		valid, err := services.UseBackupCode(user.ID.String(), "123456")
		assert.NoError(t, err)
		assert.True(t, valid)

		// Verify backup codes count decreased
		count, err := services.GetBackupCodesCount(user.ID.String())
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("should not use already used backup code", func(t *testing.T) {
		// Try to use the same code again
		valid, err := services.UseBackupCode(user.ID.String(), "123456")
		assert.NoError(t, err)
		assert.False(t, valid)

		// Count should remain the same
		count, err := services.GetBackupCodesCount(user.ID.String())
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("should not use invalid backup code", func(t *testing.T) {
		valid, err := services.UseBackupCode(user.ID.String(), "invalid-code")
		assert.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("should replace backup codes", func(t *testing.T) {
		newCodes := []string{"111111", "222222", "333333", "444444"}
		err := services.ReplaceBackupCodes(user.ID.String(), newCodes)
		assert.NoError(t, err)

		// Verify new backup codes count
		count, err := services.GetBackupCodesCount(user.ID.String())
		assert.NoError(t, err)
		assert.Equal(t, 4, count)

		// Verify old codes no longer work
		valid, err := services.UseBackupCode(user.ID.String(), "789012")
		assert.NoError(t, err)
		assert.False(t, valid)

		// Verify new codes work
		valid, err = services.UseBackupCode(user.ID.String(), "111111")
		assert.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		_, err := services.GetBackupCodesCount("invalid-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")

		_, err = services.UseBackupCode("invalid-uuid", "123456")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")

		err = services.ReplaceBackupCodes("invalid-uuid", []string{"123456"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

// Benchmark tests
func BenchmarkMFAService_StoreMFASetup(b *testing.B) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		b.Fatal(err)
	}

	err = db.AutoMigrate(&models.User{}, &models.MFASetup{}, &models.BackupCode{})
	if err != nil {
		b.Fatal(err)
	}

	// Mock the global DB
	originalDB := services.DB
	services.DB = db
	defer func() { services.DB = originalDB }()

	// Create test user
	user := &models.User{
		ID:         uuid.New(),
		KeycloakID: "benchmark-user",
		Email:      "benchmark@example.com",
		Username:   "benchuser",
		IsActive:   true,
	}
	err = db.Create(user).Error
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		secret := "JBSWY3DPEHPK3PXP"
		backupCodes := []string{"123456", "789012", "345678"}

		err := services.StoreMFASetup(user.ID.String(), secret, backupCodes)
		if err != nil {
			b.Fatal(err)
		}
	}
}
