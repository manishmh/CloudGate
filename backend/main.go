package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"cloudgate-backend/internal/config"
	"cloudgate-backend/internal/handlers"
	"cloudgate-backend/internal/middleware"
	"cloudgate-backend/internal/services"
)

// CloudGate Backend - v1.0.1 - CI/CD Pipeline Test
func main() {
	// Load .env file (optional for Cloud Run)
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found or error loading .env file: %v", err)
		log.Printf("Continuing with system environment variables...")
	} else {
		log.Printf("Successfully loaded .env file")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Validate configuration
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatal("❌ Configuration validation failed:", err)
	}
	log.Printf("✅ Configuration validated successfully")

	// Initialize database with retry logic for Cloud Run
	log.Printf("🔄 Initializing database connection...")
	maxRetries := 3
	var dbErr error
	for i := 0; i < maxRetries; i++ {
		if dbErr = services.InitializeDatabase(); dbErr != nil {
			log.Printf("❌ Database initialization attempt %d/%d failed: %v", i+1, maxRetries, dbErr)
			if i < maxRetries-1 {
				time.Sleep(time.Duration(i+1) * time.Second)
			}
		} else {
			log.Printf("✅ Database initialized successfully")
			break
		}
	}

	if dbErr != nil {
		log.Fatal("❌ Failed to initialize database after retries:", dbErr)
	}
	defer services.CloseDatabase()

	// Initialize SaaS applications
	log.Printf("🔄 Initializing SaaS applications...")
	services.InitializeSaaSApps()
	log.Printf("✅ SaaS applications initialized")

	// Initialize demo user for development (optional for production)
	if os.Getenv("SKIP_DEMO_USER") != "true" {
		log.Printf("🔄 Initializing demo user...")
		userService := services.NewUserService(services.GetDB())
		_, err := userService.GetOrCreateDemoUser()
		if err != nil {
			log.Printf("⚠️ Warning: Failed to create demo user: %v", err)
		} else {
			log.Printf("✅ Demo user initialized successfully")
		}
	}

	// Set Gin mode for production
	if os.Getenv("GIN_MODE") == "" {
		if os.Getenv("PORT") != "" { // Cloud Run sets PORT
			gin.SetMode(gin.ReleaseMode)
		} else {
			gin.SetMode(gin.DebugMode)
		}
	}

	// Create router
	router := gin.Default()

	// Setup middleware
	router.Use(middleware.SetupCORS(cfg))
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(handlers.DetailedRequestLogger()) // Add detailed logging

	// Setup routes
	handlers.SetupRoutes(router, cfg)

	// Start session cleanup routine
	go func() {
		sessionService := services.NewSessionService(services.GetDB())
		for {
			time.Sleep(1 * time.Hour) // Run every hour
			if err := sessionService.CleanupExpiredSessions(); err != nil {
				log.Printf("Failed to cleanup expired sessions: %v", err)
			}
		}
	}()

	// Log startup information
	log.Printf("🚀 ========================================")
	log.Printf("🚀 CloudGate Backend Starting")
	log.Printf("🚀 ========================================")
	log.Printf("📅 Timestamp: %s", time.Now().UTC().Format(time.RFC3339))
	log.Printf("🌐 Port: %s", cfg.Port)
	log.Printf("🔐 Keycloak URL: %s", cfg.KeycloakURL)
	log.Printf("🏰 Keycloak Realm: %s", cfg.KeycloakRealm)
	log.Printf("🔧 Keycloak Client ID: %s", cfg.KeycloakClientID)
	log.Printf("🌍 Allowed Origins: %v", cfg.AllowedOrigins)
	log.Printf("📦 SaaS Applications: %d", len(services.GetAllSaaSApps()))
	log.Printf("💾 Database: Initialized and migrations completed")
	log.Printf("🔄 Session cleanup: Running every hour")
	log.Printf("📝 Logging: Enhanced debugging enabled")
	log.Printf("🚀 ========================================")

	// Start server - bind to all interfaces for Cloud Run
	address := "0.0.0.0:" + cfg.Port
	log.Printf("🚀 Server starting on %s...", address)
	if err := router.Run(address); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
