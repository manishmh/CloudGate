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

// setupUserTestDB initializes an in-memory SQLite database for user service tests
func setupUserTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to test database")

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err, "Failed to migrate database schema")

	return db
}

// setupTestUserService sets up a test user service with database
func setupTestUserService(t *testing.T) (*services.UserService, *gorm.DB) {
	db := setupUserTestDB(t)
	userService := services.NewUserService(db)
	return userService, db
}

func TestUserService_CreateOrUpdateUser(t *testing.T) {
	userService, db := setupTestUserService(t)

	t.Run("should create new user successfully", func(t *testing.T) {
		keycloakID := uuid.New().String()
		email := "test@example.com"
		username := "testuser"
		firstName := "Test"
		lastName := "User"

		user, err := userService.CreateOrUpdateUser(keycloakID, email, username, firstName, lastName)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, keycloakID, user.KeycloakID)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, firstName, user.FirstName)
		assert.Equal(t, lastName, user.LastName)
		assert.True(t, user.IsActive)

		// Verify user was created in database
		var dbUser models.User
		err = db.Where("keycloak_id = ?", keycloakID).First(&dbUser).Error
		assert.NoError(t, err)
		assert.Equal(t, email, dbUser.Email)
	})

	t.Run("should update existing user", func(t *testing.T) {
		// Create initial user
		keycloakID := uuid.New().String()
		originalEmail := "original@example.com"

		user1, err := userService.CreateOrUpdateUser(keycloakID, originalEmail, "user1", "Original", "User")
		require.NoError(t, err)

		// Update user with new information
		newEmail := "updated@example.com"
		user2, err := userService.CreateOrUpdateUser(keycloakID, newEmail, "user1", "Updated", "User")

		assert.NoError(t, err)
		assert.NotNil(t, user2)
		assert.Equal(t, user1.ID, user2.ID) // Same user ID
		assert.Equal(t, newEmail, user2.Email)
		assert.Equal(t, "Updated", user2.FirstName)
		assert.NotNil(t, user2.LastLoginAt)
	})

	t.Run("should handle empty fields gracefully", func(t *testing.T) {
		keycloakID := uuid.New().String()

		user, err := userService.CreateOrUpdateUser(keycloakID, "", "", "", "")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, keycloakID, user.KeycloakID)
		assert.Equal(t, "", user.Email)
		assert.Equal(t, "", user.Username)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	userService, db := setupTestUserService(t)

	t.Run("should retrieve existing user", func(t *testing.T) {
		// Create test user directly in database
		testUser := models.User{
			ID:         uuid.New(),
			KeycloakID: uuid.New().String(),
			Email:      "test@example.com",
			Username:   "testuser",
			FirstName:  "Test",
			LastName:   "User",
			IsActive:   true,
		}
		err := db.Create(&testUser).Error
		require.NoError(t, err)

		// Retrieve user
		retrievedUser, err := userService.GetUserByID(testUser.ID)

		assert.NoError(t, err)
		assert.NotNil(t, retrievedUser)
		assert.Equal(t, testUser.ID, retrievedUser.ID)
		assert.Equal(t, testUser.Email, retrievedUser.Email)
		assert.Equal(t, testUser.Username, retrievedUser.Username)
	})

	t.Run("should return error for non-existent user", func(t *testing.T) {
		nonExistentID := uuid.New()

		user, err := userService.GetUserByID(nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("should return error for inactive user", func(t *testing.T) {
		// Create inactive user
		inactiveUser := models.User{
			ID:         uuid.New(),
			KeycloakID: uuid.New().String(),
			Email:      "inactive@example.com",
			Username:   "inactiveuser",
			IsActive:   false,
		}
		err := db.Create(&inactiveUser).Error
		require.NoError(t, err)

		// Try to retrieve inactive user
		user, err := userService.GetUserByID(inactiveUser.ID)

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserService_GetOrCreateDemoUser(t *testing.T) {
	userService, db := setupTestUserService(t)

	t.Run("should create demo user if not exists", func(t *testing.T) {
		user, err := userService.GetOrCreateDemoUser()

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "demo@cloudgate.dev", user.Email) // Fixed: should be .dev not .com
		assert.Equal(t, "demo-user", user.Username)
		assert.Equal(t, "Demo", user.FirstName)
		assert.Equal(t, "User", user.LastName)
		assert.True(t, user.IsActive)

		// Verify user was created in database
		var dbUser models.User
		err = db.Where("email = ?", "demo@cloudgate.dev").First(&dbUser).Error
		assert.NoError(t, err)
	})

	t.Run("should return existing demo user", func(t *testing.T) {
		// Create demo user first
		user1, err := userService.GetOrCreateDemoUser()
		require.NoError(t, err)

		// Get demo user again
		user2, err := userService.GetOrCreateDemoUser()

		assert.NoError(t, err)
		assert.NotNil(t, user2)
		assert.Equal(t, user1.ID, user2.ID) // Same user
		assert.Equal(t, user1.Email, user2.Email)
	})
}

// Benchmark tests
func BenchmarkUserService_CreateOrUpdateUser(b *testing.B) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		b.Fatal(err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		b.Fatal(err)
	}

	userService := services.NewUserService(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		keycloakID := uuid.New().String()
		email := "benchmark@example.com"
		username := "benchuser"

		_, err := userService.CreateOrUpdateUser(keycloakID, email, username, "Bench", "User")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUserService_GetUserByID(b *testing.B) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		b.Fatal(err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		b.Fatal(err)
	}

	userService := services.NewUserService(db)

	// Create test user
	testUser := models.User{
		ID:         uuid.New(),
		KeycloakID: uuid.New().String(),
		Email:      "benchmark@example.com",
		Username:   "benchuser",
		IsActive:   true,
	}
	err = db.Create(&testUser).Error
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := userService.GetUserByID(testUser.ID)
		if err != nil {
			b.Fatal(err)
		}
	}
}
