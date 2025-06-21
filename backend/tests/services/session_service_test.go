package services_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"cloudgate-backend/internal/models"
	"cloudgate-backend/internal/services"
)

// setupSessionTestDB initializes an in-memory SQLite database for session service tests
func setupSessionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to test database")

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Session{})
	require.NoError(t, err, "Failed to migrate database schema")

	return db
}

// setupTestSessionService sets up a test session service with database
func setupTestSessionService(t *testing.T) (*services.SessionService, *gorm.DB, *models.User) {
	db := setupSessionTestDB(t)

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

	sessionService := services.NewSessionServiceForTesting(db)
	return sessionService, db, user
}

func TestSessionService_CreateSession(t *testing.T) {
	service, db, user := setupTestSessionService(t)

	t.Run("should create session successfully", func(t *testing.T) {
		session, err := service.CreateSession(user.ID, "192.168.1.100", "Mozilla/5.0 Test Browser")

		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, user.ID, session.UserID)
		assert.Equal(t, "192.168.1.100", session.IPAddress)
		assert.Equal(t, "Mozilla/5.0 Test Browser", session.UserAgent)
		assert.True(t, session.IsActive)
		assert.NotEmpty(t, session.SessionToken)
		assert.True(t, session.ExpiresAt.After(time.Now()))

		// Verify session was stored in database
		var storedSession models.Session
		err = db.First(&storedSession, session.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, session.SessionToken, storedSession.SessionToken)
	})

	t.Run("should generate unique session tokens", func(t *testing.T) {
		session1, err := service.CreateSession(user.ID, "192.168.1.100", "Browser 1")
		assert.NoError(t, err)

		session2, err := service.CreateSession(user.ID, "192.168.1.101", "Browser 2")
		assert.NoError(t, err)

		assert.NotEqual(t, session1.SessionToken, session2.SessionToken)
	})
}

func TestSessionService_GetSessionByToken(t *testing.T) {
	service, db, user := setupTestSessionService(t)

	t.Run("should retrieve session by valid token", func(t *testing.T) {
		// Create session
		createdSession, err := service.CreateSession(user.ID, "192.168.1.100", "Test Browser")
		assert.NoError(t, err)

		// Retrieve session by token
		retrievedSession, err := service.GetSessionByToken(createdSession.SessionToken)
		assert.NoError(t, err)
		assert.Equal(t, createdSession.ID, retrievedSession.ID)
		assert.Equal(t, createdSession.SessionToken, retrievedSession.SessionToken)
		assert.Equal(t, user.ID, retrievedSession.User.ID)
	})

	t.Run("should return error for invalid token", func(t *testing.T) {
		_, err := service.GetSessionByToken("invalid-token")
		assert.Error(t, err)
	})

	t.Run("should return error for expired session", func(t *testing.T) {
		// Create session
		session, err := service.CreateSession(user.ID, "192.168.1.100", "Test Browser")
		assert.NoError(t, err)

		// Manually expire the session
		session.ExpiresAt = time.Now().Add(-1 * time.Hour)
		err = db.Save(session).Error
		assert.NoError(t, err)

		// Try to retrieve expired session
		_, err = service.GetSessionByToken(session.SessionToken)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "session expired")
	})
}

func TestSessionService_ValidateSession(t *testing.T) {
	service, _, user := setupTestSessionService(t)

	t.Run("should validate session and return user", func(t *testing.T) {
		// Create session
		session, err := service.CreateSession(user.ID, "192.168.1.100", "Test Browser")
		assert.NoError(t, err)

		// Validate session
		validatedUser, err := service.ValidateSession(session.SessionToken)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, validatedUser.ID)
		assert.Equal(t, user.Email, validatedUser.Email)
	})

	t.Run("should return error for invalid session", func(t *testing.T) {
		_, err := service.ValidateSession("invalid-token")
		assert.Error(t, err)
	})
}

func TestSessionService_RefreshSession(t *testing.T) {
	service, _, user := setupTestSessionService(t)

	t.Run("should refresh session and extend expiry", func(t *testing.T) {
		// Create session
		session, err := service.CreateSession(user.ID, "192.168.1.100", "Test Browser")
		assert.NoError(t, err)

		originalExpiry := session.ExpiresAt

		// Wait a bit to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		// Refresh session
		refreshedSession, err := service.RefreshSession(session.SessionToken)
		assert.NoError(t, err)
		assert.True(t, refreshedSession.ExpiresAt.After(originalExpiry))
	})

	t.Run("should return error for invalid session", func(t *testing.T) {
		_, err := service.RefreshSession("invalid-token")
		assert.Error(t, err)
	})
}

func TestSessionService_InvalidateSession(t *testing.T) {
	service, db, user := setupTestSessionService(t)

	t.Run("should invalidate session", func(t *testing.T) {
		// Create session
		session, err := service.CreateSession(user.ID, "192.168.1.100", "Test Browser")
		assert.NoError(t, err)

		// Invalidate session
		err = service.InvalidateSession(session.SessionToken)
		assert.NoError(t, err)

		// Verify session is inactive
		var updatedSession models.Session
		err = db.First(&updatedSession, session.ID).Error
		assert.NoError(t, err)
		assert.False(t, updatedSession.IsActive)

		// Should not be able to validate invalidated session
		_, err = service.ValidateSession(session.SessionToken)
		assert.Error(t, err)
	})

	t.Run("should handle invalidating non-existent session", func(t *testing.T) {
		err := service.InvalidateSession("non-existent-token")
		assert.NoError(t, err) // Should not error, just no effect
	})
}

func TestSessionService_InvalidateAllUserSessions(t *testing.T) {
	service, db, user := setupTestSessionService(t)

	t.Run("should invalidate all user sessions", func(t *testing.T) {
		// Create multiple sessions
		session1, err := service.CreateSession(user.ID, "192.168.1.100", "Browser 1")
		assert.NoError(t, err)

		session2, err := service.CreateSession(user.ID, "192.168.1.101", "Browser 2")
		assert.NoError(t, err)

		// Invalidate all sessions
		err = service.InvalidateAllUserSessions(user.ID)
		assert.NoError(t, err)

		// Verify all sessions are inactive
		var sessions []models.Session
		err = db.Where("user_id = ?", user.ID).Find(&sessions).Error
		assert.NoError(t, err)

		for _, session := range sessions {
			assert.False(t, session.IsActive)
		}

		// Should not be able to validate any session
		_, err = service.ValidateSession(session1.SessionToken)
		assert.Error(t, err)

		_, err = service.ValidateSession(session2.SessionToken)
		assert.Error(t, err)
	})
}

func TestSessionService_GetUserSessions(t *testing.T) {
	service, db, user := setupTestSessionService(t)

	t.Run("should get active user sessions", func(t *testing.T) {
		// Create multiple sessions
		session1, err := service.CreateSession(user.ID, "192.168.1.100", "Browser 1")
		assert.NoError(t, err)

		session2, err := service.CreateSession(user.ID, "192.168.1.101", "Browser 2")
		assert.NoError(t, err)

		// Invalidate one session
		err = service.InvalidateSession(session1.SessionToken)
		assert.NoError(t, err)

		// Get active sessions
		sessions, err := service.GetUserSessions(user.ID)
		assert.NoError(t, err)
		assert.Len(t, sessions, 1) // Only one active session
		assert.Equal(t, session2.ID, sessions[0].ID)
	})

	t.Run("should return empty list for user with no active sessions", func(t *testing.T) {
		// Create new user
		newUser := &models.User{
			ID:         uuid.New(),
			KeycloakID: "test-keycloak-id-2",
			Email:      "test2@example.com",
			Username:   "testuser2",
		}
		err := db.Create(newUser).Error
		assert.NoError(t, err)

		sessions, err := service.GetUserSessions(newUser.ID)
		assert.NoError(t, err)
		assert.Empty(t, sessions)
	})
}

func TestSessionService_CleanupExpiredSessions(t *testing.T) {
	service, db, user := setupTestSessionService(t)

	t.Run("should cleanup expired sessions", func(t *testing.T) {
		// Create session and manually expire it
		session, err := service.CreateSession(user.ID, "192.168.1.100", "Test Browser")
		assert.NoError(t, err)

		// Set expiry to 8 days ago (beyond the 7-day cleanup threshold)
		expiredTime := time.Now().Add(-8 * 24 * time.Hour)
		session.ExpiresAt = expiredTime
		err = db.Save(session).Error
		assert.NoError(t, err)

		// Run cleanup
		err = service.CleanupExpiredSessions()
		assert.NoError(t, err)

		// Verify session was deleted
		var count int64
		err = db.Model(&models.Session{}).Where("id = ?", session.ID).Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("should not cleanup recently expired sessions", func(t *testing.T) {
		// Create session and expire it recently (within 7 days)
		session, err := service.CreateSession(user.ID, "192.168.1.100", "Test Browser")
		assert.NoError(t, err)

		// Set expiry to 1 day ago (within the 7-day grace period)
		recentExpiredTime := time.Now().Add(-1 * 24 * time.Hour)
		session.ExpiresAt = recentExpiredTime
		err = db.Save(session).Error
		assert.NoError(t, err)

		// Run cleanup
		err = service.CleanupExpiredSessions()
		assert.NoError(t, err)

		// Verify session still exists
		var count int64
		err = db.Model(&models.Session{}).Where("id = ?", session.ID).Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})
}

func TestSessionService_GetSessionStats(t *testing.T) {
	service, db, user := setupTestSessionService(t)

	t.Run("should return session statistics", func(t *testing.T) {
		// Create active session
		_, err := service.CreateSession(user.ID, "192.168.1.100", "Active Browser")
		assert.NoError(t, err)

		// Create expired session
		expiredSession, err := service.CreateSession(user.ID, "192.168.1.101", "Expired Browser")
		assert.NoError(t, err)

		// Manually expire the session
		expiredSession.ExpiresAt = time.Now().Add(-1 * time.Hour)
		err = db.Save(expiredSession).Error
		assert.NoError(t, err)

		// Get stats
		stats, err := service.GetSessionStats()
		assert.NoError(t, err)

		assert.Equal(t, int64(1), stats["active_sessions"])
		assert.Equal(t, int64(1), stats["expired_sessions"])
		assert.GreaterOrEqual(t, stats["sessions_today"].(int64), int64(2))
	})
}
