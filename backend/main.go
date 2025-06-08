package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"cloudgate-backend/internal/config"
	"cloudgate-backend/internal/handlers"
	"cloudgate-backend/internal/middleware"
	"cloudgate-backend/internal/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	if err := services.InitializeDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer services.CloseDatabase()

	// Initialize SaaS applications
	services.InitializeSaaSApps()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	router := gin.Default()

	// Setup middleware
	router.Use(middleware.SetupCORS(cfg))
	router.Use(middleware.SecurityHeadersMiddleware())

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
