package services_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"cloudgate-backend/internal/models"
	"cloudgate-backend/internal/services"
)

// setupRiskTestDB initializes an in-memory SQLite database for risk service tests
func setupRiskTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate all required schemas
	err = db.AutoMigrate(
		&models.User{},
		&services.RiskAssessment{},
		&services.RiskThresholds{},
		&services.DeviceFingerprint{},
		&services.WebAuthnCredential{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database schema: %v", err)
	}

	return db
}

// setupTestRiskService sets up a test risk service with database
func setupTestRiskService(t *testing.T) (*gorm.DB, *models.User) {
	db := setupRiskTestDB(t)

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
	assert.NoError(t, err)

	return db, user
}

func TestRiskService_StoreRiskAssessment(t *testing.T) {
	db, user := setupTestRiskService(t)

	// Mock the global DB
	originalDB := services.DB
	services.DB = db
	defer func() { services.DB = originalDB }()

	t.Run("should store risk assessment successfully", func(t *testing.T) {
		assessment := map[string]interface{}{
			"user_id":            user.ID.String(),
			"session_id":         "test-session-123",
			"ip_address":         "192.168.1.100",
			"user_agent":         "Mozilla/5.0 (Test Browser)",
			"location":           map[string]string{"country": "US", "city": "San Francisco"},
			"device_fingerprint": "test-fingerprint-123",
			"behavior_signals":   map[string]float64{"typing_speed": 120.5, "mouse_movement": 85.2},
			"risk_score":         0.65,
			"risk_level":         "medium",
			"risk_factors":       []string{"new_device", "unusual_location"},
			"recommendations":    []string{"require_mfa", "verify_email"},
		}

		err := services.StoreRiskAssessment(assessment)
		assert.NoError(t, err)

		// Verify assessment was stored
		var storedAssessment services.RiskAssessment
		err = db.Where("user_id = ?", user.ID).First(&storedAssessment).Error
		assert.NoError(t, err)
		assert.Equal(t, "test-session-123", storedAssessment.SessionID)
		assert.Equal(t, 0.65, storedAssessment.RiskScore)
		assert.Equal(t, "medium", storedAssessment.RiskLevel)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		assessment := map[string]interface{}{
			"user_id":    "invalid-uuid",
			"risk_score": 0.5,
		}

		err := services.StoreRiskAssessment(assessment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user_id format")
	})

	t.Run("should return error for missing user ID", func(t *testing.T) {
		assessment := map[string]interface{}{
			"risk_score": 0.5,
		}

		err := services.StoreRiskAssessment(assessment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user_id in assessment")
	})
}

func TestRiskService_GetLatestRiskAssessment(t *testing.T) {
	db, user := setupTestRiskService(t)

	// Mock the global DB
	originalDB := services.DB
	services.DB = db
	defer func() { services.DB = originalDB }()

	t.Run("should return latest risk assessment", func(t *testing.T) {
		// Store multiple assessments
		assessment1 := map[string]interface{}{
			"user_id":    user.ID.String(),
			"risk_score": 0.3,
			"risk_level": "low",
		}
		assessment2 := map[string]interface{}{
			"user_id":    user.ID.String(),
			"risk_score": 0.8,
			"risk_level": "high",
		}

		err := services.StoreRiskAssessment(assessment1)
		assert.NoError(t, err)

		time.Sleep(10 * time.Millisecond) // Ensure different timestamps

		err = services.StoreRiskAssessment(assessment2)
		assert.NoError(t, err)

		// Get latest assessment
		latest, err := services.GetLatestRiskAssessment(user.ID.String())
		assert.NoError(t, err)

		latestMap := latest.(map[string]interface{})
		assert.Equal(t, 0.8, latestMap["risk_score"])
		assert.Equal(t, "high", latestMap["risk_level"])
	})

	t.Run("should return error for non-existent user", func(t *testing.T) {
		nonExistentUserID := uuid.New().String()
		_, err := services.GetLatestRiskAssessment(nonExistentUserID)
		assert.Error(t, err)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		_, err := services.GetLatestRiskAssessment("invalid-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

func TestRiskService_GetRiskAssessmentHistory(t *testing.T) {
	db, user := setupTestRiskService(t)

	// Mock the global DB
	originalDB := services.DB
	services.DB = db
	defer func() { services.DB = originalDB }()

	t.Run("should return risk assessment history with limit", func(t *testing.T) {
		// Store multiple assessments
		for i := 0; i < 5; i++ {
			assessment := map[string]interface{}{
				"user_id":    user.ID.String(),
				"risk_score": float64(i) * 0.2,
				"risk_level": "test",
			}
			err := services.StoreRiskAssessment(assessment)
			assert.NoError(t, err)
			time.Sleep(5 * time.Millisecond) // Ensure different timestamps
		}

		// Get history with limit
		history, err := services.GetRiskAssessmentHistory(user.ID.String(), 3)
		assert.NoError(t, err)
		assert.Len(t, history, 3)

		// Should be in descending order (latest first)
		firstAssessment := history[0].(map[string]interface{})
		assert.Equal(t, 0.8, firstAssessment["risk_score"]) // Latest assessment
	})

	t.Run("should return empty history for user with no assessments", func(t *testing.T) {
		newUser := &models.User{
			ID:         uuid.New(),
			KeycloakID: "test-keycloak-id-2",
			Email:      "test2@example.com",
			Username:   "testuser2",
		}
		err := db.Create(newUser).Error
		assert.NoError(t, err)

		history, err := services.GetRiskAssessmentHistory(newUser.ID.String(), 10)
		assert.NoError(t, err)
		assert.Empty(t, history)
	})
}

func TestRiskService_RiskThresholds(t *testing.T) {
	db, _ := setupTestRiskService(t)

	// Mock the global DB
	originalDB := services.DB
	services.DB = db
	defer func() { services.DB = originalDB }()

	t.Run("should update risk thresholds", func(t *testing.T) {
		thresholds := map[string]float64{
			"vpn_risk":         0.4,
			"tor_risk":         0.95,
			"new_device_risk":  0.8,
			"low_threshold":    0.25,
			"medium_threshold": 0.55,
			"high_threshold":   0.85,
		}

		err := services.UpdateRiskThresholds(thresholds)
		assert.NoError(t, err)

		// Verify thresholds were updated
		var storedThresholds services.RiskThresholds
		err = db.First(&storedThresholds).Error
		assert.NoError(t, err)
		assert.Equal(t, 0.4, storedThresholds.VPNRisk)
		assert.Equal(t, 0.95, storedThresholds.TorRisk)
		assert.Equal(t, 0.8, storedThresholds.NewDeviceRisk)
	})

	t.Run("should update existing thresholds", func(t *testing.T) {
		// Update again with different values
		newThresholds := map[string]float64{
			"vpn_risk":      0.5,
			"tor_risk":      0.99,
			"low_threshold": 0.2,
		}

		err := services.UpdateRiskThresholds(newThresholds)
		assert.NoError(t, err)

		// Verify only specified thresholds were updated
		var storedThresholds services.RiskThresholds
		err = db.First(&storedThresholds).Error
		assert.NoError(t, err)
		assert.Equal(t, 0.5, storedThresholds.VPNRisk)
		assert.Equal(t, 0.99, storedThresholds.TorRisk)
		assert.Equal(t, 0.2, storedThresholds.LowThreshold)
		// These should remain unchanged from previous test
		assert.Equal(t, 0.8, storedThresholds.NewDeviceRisk)
	})
}

func TestRiskService_DeviceFingerprinting(t *testing.T) {
	db, user := setupTestRiskService(t)

	// Mock the global DB
	originalDB := services.DB
	services.DB = db
	defer func() { services.DB = originalDB }()

	fingerprint := "test-fingerprint-123"

	t.Run("should detect new device", func(t *testing.T) {
		isNew, err := services.IsNewDevice(user.ID.String(), fingerprint)
		assert.NoError(t, err)
		assert.True(t, isNew)
	})

	t.Run("should register device fingerprint", func(t *testing.T) {
		err := services.RegisterDeviceFingerprint(
			user.ID.String(),
			fingerprint,
			"MacBook Pro",
			"desktop",
			"Chrome",
			"macOS",
		)
		assert.NoError(t, err)

		// Verify device was registered
		var deviceFP services.DeviceFingerprint
		err = db.Where("user_id = ? AND fingerprint = ?", user.ID, fingerprint).First(&deviceFP).Error
		assert.NoError(t, err)
		assert.Equal(t, "MacBook Pro", deviceFP.DeviceName)
		assert.Equal(t, "desktop", deviceFP.DeviceType)
		assert.False(t, deviceFP.IsTrusted) // Should not be trusted by default
	})

	t.Run("should detect existing device", func(t *testing.T) {
		isNew, err := services.IsNewDevice(user.ID.String(), fingerprint)
		assert.NoError(t, err)
		assert.False(t, isNew)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		_, err := services.IsNewDevice("invalid-uuid", fingerprint)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")

		err = services.RegisterDeviceFingerprint("invalid-uuid", fingerprint, "Device", "mobile", "Safari", "iOS")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

func TestRiskService_WebAuthnCredentials(t *testing.T) {
	db, user := setupTestRiskService(t)

	// Mock the global DB
	originalDB := services.DB
	services.DB = db
	defer func() { services.DB = originalDB }()

	credentialID := "test-credential-123"
	attestationObject := []byte("test-attestation-data")

	t.Run("should store WebAuthn credential", func(t *testing.T) {
		err := services.StoreWebAuthnCredential(user.ID.String(), credentialID, attestationObject)
		assert.NoError(t, err)

		// Verify credential was stored
		var credential services.WebAuthnCredential
		err = db.Where("user_id = ? AND credential_id = ?", user.ID, credentialID).First(&credential).Error
		assert.NoError(t, err)
		assert.Equal(t, credentialID, credential.CredentialID)
		assert.Equal(t, attestationObject, credential.AttestationObject)
	})

	t.Run("should get user WebAuthn credentials", func(t *testing.T) {
		credentials, err := services.GetUserWebAuthnCredentials(user.ID.String())
		assert.NoError(t, err)
		assert.Len(t, credentials, 1)
		assert.Equal(t, credentialID, credentials[0].CredentialID)
	})

	t.Run("should verify WebAuthn credential", func(t *testing.T) {
		valid, err := services.VerifyWebAuthnCredential(user.ID.String(), credentialID)
		assert.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("should update WebAuthn credential usage", func(t *testing.T) {
		err := services.UpdateWebAuthnCredentialUsage(user.ID.String(), credentialID)
		assert.NoError(t, err)

		// Verify last_used was updated
		var credential services.WebAuthnCredential
		err = db.Where("user_id = ? AND credential_id = ?", user.ID, credentialID).First(&credential).Error
		assert.NoError(t, err)
		assert.NotNil(t, credential.LastUsed)
	})

	t.Run("should delete WebAuthn credential", func(t *testing.T) {
		err := services.DeleteWebAuthnCredential(user.ID.String(), credentialID)
		assert.NoError(t, err)

		// Verify credential was deleted
		var count int64
		err = db.Model(&services.WebAuthnCredential{}).Where("user_id = ? AND credential_id = ?", user.ID, credentialID).Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("should return false for non-existent credential", func(t *testing.T) {
		valid, err := services.VerifyWebAuthnCredential(user.ID.String(), "non-existent-credential")
		assert.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		_, err := services.GetUserWebAuthnCredentials("invalid-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")

		err = services.StoreWebAuthnCredential("invalid-uuid", credentialID, attestationObject)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")

		_, err = services.VerifyWebAuthnCredential("invalid-uuid", credentialID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}
