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

// setupOAuthTestDB initializes an in-memory SQLite database for OAuth monitoring tests
func setupOAuthTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate all required schemas
	err = db.AutoMigrate(
		&models.User{},
		&models.AppConnection{},
		&models.ConnectionHealthMetrics{},
		&models.SecurityEvent{},
		&models.TrustedDevice{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database schema: %v", err)
	}

	return db
}

// createTestUser creates a test user for OAuth monitoring tests
func createTestUser(t *testing.T, db *gorm.DB) *models.User {
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
	assert.NoError(t, err)
	return user
}

// createTestConnection creates a test app connection
func createTestConnection(t *testing.T, db *gorm.DB, userID uuid.UUID, status string) *models.AppConnection {
	connection := &models.AppConnection{
		ID:              uuid.New(),
		UserID:          userID,
		AppID:           "google-workspace",
		AppName:         "Google Workspace",
		Provider:        "google",
		Status:          status,
		AccessToken:     "test-access-token",
		RefreshToken:    "test-refresh-token",
		Scopes:          "email profile",
		UserEmail:       "test@example.com",
		UserName:        "Test User",
		ConnectedAt:     time.Now(),
		HealthStatus:    "healthy",
		ResponseTime:    150,
		ErrorCount:      0,
		UptimePercent:   99.5,
		UsageCount:      10,
		DataTransferred: 1024,
	}
	err := db.Create(connection).Error
	assert.NoError(t, err)
	return connection
}

func TestOAuthMonitoringService_GetUserConnections(t *testing.T) {
	db := setupOAuthTestDB(t)
	service := services.NewOAuthMonitoringService(db)
	user := createTestUser(t, db)

	t.Run("should return empty list for user with no connections", func(t *testing.T) {
		connections, err := service.GetUserConnections(user.ID.String())

		assert.NoError(t, err)
		assert.Empty(t, connections)
	})

	t.Run("should return user connections with health data", func(t *testing.T) {
		// Create test connections
		conn1 := createTestConnection(t, db, user.ID, "connected")
		_ = createTestConnection(t, db, user.ID, "error")

		connections, err := service.GetUserConnections(user.ID.String())

		assert.NoError(t, err)
		assert.Len(t, connections, 2)

		// Check first connection
		assert.Equal(t, conn1.ID, connections[0].ID)
		assert.Equal(t, "healthy", connections[0].Health.Status)
		assert.Equal(t, 150, connections[0].Health.ResponseTime)
		assert.Equal(t, 99.5, connections[0].Health.Uptime)
		assert.Equal(t, int64(10), connections[0].UsageCount)
		assert.Equal(t, "1.0 KB", connections[0].DataTransferred)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		_, err := service.GetUserConnections("invalid-uuid")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

func TestOAuthMonitoringService_GetConnectionStats(t *testing.T) {
	db := setupOAuthTestDB(t)
	service := services.NewOAuthMonitoringService(db)
	user := createTestUser(t, db)

	t.Run("should return zero stats for user with no connections", func(t *testing.T) {
		stats, err := service.GetConnectionStats(user.ID.String())

		assert.NoError(t, err)
		assert.Equal(t, 0, stats.TotalConnections)
		assert.Equal(t, 0, stats.ActiveConnections)
		assert.Equal(t, 0, stats.FailedConnections)
	})

	t.Run("should calculate correct statistics", func(t *testing.T) {
		// Create various connections
		createTestConnection(t, db, user.ID, "connected")
		createTestConnection(t, db, user.ID, "connected")
		createTestConnection(t, db, user.ID, "error")

		stats, err := service.GetConnectionStats(user.ID.String())

		assert.NoError(t, err)
		assert.Equal(t, 3, stats.TotalConnections)
		assert.Equal(t, 2, stats.ActiveConnections)
		assert.Equal(t, 1, stats.FailedConnections)
		assert.Equal(t, 150, stats.AverageResponseTime) // From test data
		assert.Equal(t, 99.5, stats.UptimePercentage)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		_, err := service.GetConnectionStats("invalid-uuid")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

func TestOAuthMonitoringService_TestConnection(t *testing.T) {
	db := setupOAuthTestDB(t)
	service := services.NewOAuthMonitoringService(db)
	user := createTestUser(t, db)
	connection := createTestConnection(t, db, user.ID, "connected")

	t.Run("should successfully test connection and update health", func(t *testing.T) {
		err := service.TestConnection(user.ID.String(), connection.ID.String())

		assert.NoError(t, err)

		// Verify connection was updated
		var updatedConnection models.AppConnection
		err = db.First(&updatedConnection, connection.ID).Error
		assert.NoError(t, err)
		assert.NotNil(t, updatedConnection.LastHealthCheck)
		assert.GreaterOrEqual(t, updatedConnection.ResponseTime, 0) // Response time can be 0 for very fast operations
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		err := service.TestConnection("invalid-uuid", connection.ID.String())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("should return error for invalid connection ID", func(t *testing.T) {
		err := service.TestConnection(user.ID.String(), "invalid-uuid")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid connection ID")
	})

	t.Run("should return error for non-existent connection", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := service.TestConnection(user.ID.String(), nonExistentID.String())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection not found")
	})
}

func TestOAuthMonitoringService_RecordUsage(t *testing.T) {
	db := setupOAuthTestDB(t)
	service := services.NewOAuthMonitoringService(db)
	user := createTestUser(t, db)
	connection := createTestConnection(t, db, user.ID, "connected")

	t.Run("should record usage and update statistics", func(t *testing.T) {
		err := service.RecordUsage(user.ID.String(), connection.ID.String(), 2048)

		assert.NoError(t, err)

		// Verify connection usage was updated
		var updatedConnection models.AppConnection
		err = db.First(&updatedConnection, connection.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(11), updatedConnection.UsageCount)        // Original 10 + 1
		assert.Equal(t, int64(3072), updatedConnection.DataTransferred) // Original 1024 + 2048
		assert.NotNil(t, updatedConnection.LastUsed)
	})

	t.Run("should return error for invalid user ID", func(t *testing.T) {
		err := service.RecordUsage("invalid-uuid", connection.ID.String(), 1024)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})
}

func TestOAuthMonitoringService_CreateSecurityEvent(t *testing.T) {
	db := setupOAuthTestDB(t)
	service := services.NewOAuthMonitoringService(db)
	user := createTestUser(t, db)
	connection := createTestConnection(t, db, user.ID, "connected")

	t.Run("should create security event successfully", func(t *testing.T) {
		connIDStr := connection.ID.String()
		err := service.CreateSecurityEvent(
			user.ID.String(),
			"suspicious_login",
			"Login from unusual location",
			"high",
			"192.168.1.1",
			"Mozilla/5.0",
			"Unknown Location",
			8.5,
			&connIDStr,
		)

		assert.NoError(t, err)

		// Verify event was created
		var events []models.SecurityEvent
		err = db.Where("user_id = ?", user.ID).Find(&events).Error
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, "suspicious_login", events[0].EventType)
		assert.Equal(t, "high", events[0].Severity)
		assert.Equal(t, 8.5, events[0].RiskScore)
	})

	t.Run("should create event without connection ID", func(t *testing.T) {
		err := service.CreateSecurityEvent(
			user.ID.String(),
			"failed_mfa",
			"MFA verification failed",
			"medium",
			"192.168.1.2",
			"Mozilla/5.0",
			"New York, US",
			5.0,
			nil,
		)

		assert.NoError(t, err)
	})
}

func TestOAuthMonitoringService_GetSecurityEvents(t *testing.T) {
	db := setupOAuthTestDB(t)
	service := services.NewOAuthMonitoringService(db)
	user := createTestUser(t, db)

	t.Run("should return empty list for user with no events", func(t *testing.T) {
		events, err := service.GetSecurityEvents(user.ID.String(), 10)

		assert.NoError(t, err)
		assert.Empty(t, events)
	})

	t.Run("should return security events with limit", func(t *testing.T) {
		// Create multiple events
		for i := 0; i < 5; i++ {
			err := service.CreateSecurityEvent(
				user.ID.String(),
				"login",
				"User login",
				"low",
				"192.168.1.1",
				"Mozilla/5.0",
				"Test Location",
				1.0,
				nil,
			)
			assert.NoError(t, err)
		}

		events, err := service.GetSecurityEvents(user.ID.String(), 3)

		assert.NoError(t, err)
		assert.Len(t, events, 3)
	})
}

func TestOAuthMonitoringService_TrustedDevices(t *testing.T) {
	db := setupOAuthTestDB(t)
	service := services.NewOAuthMonitoringService(db)
	user := createTestUser(t, db)

	t.Run("should register and manage trusted devices", func(t *testing.T) {
		// Register a device
		err := service.RegisterDevice(
			user.ID.String(),
			"MacBook Pro",
			"desktop",
			"Chrome",
			"macOS",
			"unique-fingerprint-123",
			"192.168.1.100",
			"San Francisco, CA",
		)
		assert.NoError(t, err)

		// Get devices
		devices, err := service.GetTrustedDevices(user.ID.String())
		assert.NoError(t, err)
		assert.Len(t, devices, 1)
		assert.Equal(t, "MacBook Pro", devices[0].DeviceName)
		assert.False(t, devices[0].Trusted) // Should not be trusted by default

		// Trust the device
		err = service.TrustDevice(user.ID.String(), devices[0].ID.String())
		assert.NoError(t, err)

		// Verify device is now trusted
		devices, err = service.GetTrustedDevices(user.ID.String())
		assert.NoError(t, err)
		assert.True(t, devices[0].Trusted)

		// Revoke the device
		err = service.RevokeDevice(user.ID.String(), devices[0].ID.String())
		assert.NoError(t, err)

		// Verify device is removed
		devices, err = service.GetTrustedDevices(user.ID.String())
		assert.NoError(t, err)
		assert.Empty(t, devices)
	})
}
