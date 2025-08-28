package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"cloudgate-backend/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	corsConfig.AllowWildcard = true

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
		var tokenString string

		if authHeader != "" {
			// Extract token from Bearer header
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
				c.Abort()
				return
			}
			tokenString = tokenParts[1]
		} else {
			// Fallback to cookie-based token
			if cookieToken, err := c.Cookie("access_token"); err == nil && cookieToken != "" {
				tokenString = cookieToken
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		cfg := config.LoadConfig()
		parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		if expVal, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(expVal), 0).Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				c.Abort()
				return
			}
		}

		var userID uuid.UUID
		if sub, ok := claims["sub"].(string); ok {
			if id, err := uuid.Parse(sub); err == nil {
				userID = id
			}
		}
		if userID == uuid.Nil {
			userID, _ = uuid.Parse("12345678-1234-1234-1234-123456789012")
		}

		username, _ := claims["username"].(string)
		email, _ := claims["email"].(string)

		c.Set("userID", userID)
		c.Set("username", username)
		c.Set("email", email)
		c.Next()
	}
}
