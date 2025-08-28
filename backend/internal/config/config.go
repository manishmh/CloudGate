package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Config holds the application configuration
type Config struct {
	Port                string
	AllowedOrigins      []string
	JWTSecret           string
	AccessTokenTTLMin   int
	RefreshTokenTTLHour int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Use PORT environment variable for Cloud Run compatibility
	// Default to 8081
	port := getEnv("PORT", "8081")

	// Validate required environment variables for production
	validateRequiredEnvVars()

	// Parse token lifetimes from environment with sensible defaults
	accessTTL := 15
	if v := os.Getenv("ACCESS_TOKEN_TTL_MIN"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			accessTTL = i
		}
	}
	refreshTTL := 24
	if v := os.Getenv("REFRESH_TOKEN_TTL_HOUR"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			refreshTTL = i
		}
	}

	config := &Config{
		Port:                port,
		AllowedOrigins:      strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
		JWTSecret:           getEnv("JWT_SECRET", "dev-secret-change-me"),
		AccessTokenTTLMin:   accessTTL,
		RefreshTokenTTLHour: refreshTTL,
	}

	// Log configuration (excluding sensitive values)
	log.Printf("🔧 Configuration loaded:")
	log.Printf("   Port: %s", config.Port)
	log.Printf("   Allowed Origins: %v", config.AllowedOrigins)
	log.Printf("   JWT Access TTL (min): %d", config.AccessTokenTTLMin)
	log.Printf("   JWT Refresh TTL (h): %d", config.RefreshTokenTTLHour)

	return config
}

// validateRequiredEnvVars checks if required environment variables are set for production
func validateRequiredEnvVars() {
	// Only validate in Cloud Run environment (when PORT is set by platform)
	if os.Getenv("PORT") == "" {
		return // Skip validation for local development
	}

	required := []string{"JWT_SECRET"}

	missing := []string{}
	for _, env := range required {
		if os.Getenv(env) == "" {
			missing = append(missing, env)
		}
	}

	if len(missing) > 0 {
		log.Printf("⚠️ Warning: Missing required environment variables: %v", missing)
		log.Printf("ℹ️ The application will use default values, but this may cause issues in production")
	}

	// Check database configuration
	if os.Getenv("DATABASE_URL") == "" && os.Getenv("NEON_DATABASE_URL") == "" && os.Getenv("DB_TYPE") == "" {
		log.Printf("⚠️ Warning: No database configuration found. Will use SQLite fallback.")
	}
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ValidateConfig validates the loaded configuration
func ValidateConfig(cfg *Config) error {
	if cfg.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	if cfg.JWTSecret == "" {
		return fmt.Errorf("JWT secret cannot be empty")
	}

	return nil
}
