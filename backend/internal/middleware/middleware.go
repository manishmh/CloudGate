package middleware

import (
	"log"
	"net/http"
	"strings"

	"cloudgate-backend/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SetupCORS configures CORS middleware for the application
func SetupCORS(cfg *config.Config) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.AllowedOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Credentials",
	}
	corsConfig.AllowCredentials = true
	corsConfig.ExposeHeaders = []string{"*"}

	// Log CORS configuration for debugging
	log.Printf("🌐 CORS Configuration:")
	log.Printf("  📍 Allowed Origins: %v", cfg.AllowedOrigins)
	log.Printf("  🔧 Allowed Methods: %v", corsConfig.AllowMethods)
	log.Printf("  📋 Allowed Headers: %v", corsConfig.AllowHeaders)
	log.Printf("  🔐 Allow Credentials: %v", corsConfig.AllowCredentials)

	return cors.New(corsConfig)
}

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

// AuthenticationMiddleware validates the JWT token and sets user context
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Extract token from Bearer header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// For demo purposes, accept demo-user-token and create a valid user context
		// In production, you would validate the JWT token against Keycloak
		if token == "demo-user-token" {
			// Create a demo user UUID for the demo user
			demoUserUUID, err := uuid.Parse("12345678-1234-1234-1234-123456789012")
			if err != nil {
				// Fallback: generate a new UUID
				demoUserUUID = uuid.New()
			}

			// Set user context
			c.Set("userID", demoUserUUID)
			c.Set("username", "demo-user")
			c.Set("email", "demo@cloudgate.com")
			c.Next()
			return
		}

		// In production, validate JWT token here
		// For now, reject any other tokens
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
	}
}
