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

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found or error loading .env file: %v", err)
		log.Printf("Continuing with system environment variables...")
	} else {
		log.Printf("Successfully loaded .env file")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	if err := services.InitializeDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer services.CloseDatabase()

	// Initialize SaaS applications
	services.InitializeSaaSApps()

	// Initialize demo user for development
	userService := services.NewUserService(services.GetDB())
	_, err := userService.GetOrCreateDemoUser()
	if err != nil {
		log.Printf("Warning: Failed to create demo user: %v", err)
	} else {
		log.Printf("Demo user initialized successfully")
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	router := gin.Default()

	// Setup middleware
	router.Use(middleware.SetupCORS(cfg))
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(handlers.RequestResponseLogger()) // Add detailed logging

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
	log.Printf("Starting CloudGate Backend on port %s", cfg.Port)
	log.Printf("Keycloak URL: %s", cfg.KeycloakURL)
	log.Printf("Keycloak Realm: %s", cfg.KeycloakRealm)
	log.Printf("Allowed Origins: %v", cfg.AllowedOrigins)
	log.Printf("Initialized %d SaaS applications", len(services.GetAllSaaSApps()))
	log.Printf("Database initialized and migrations completed")

	// Start server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
