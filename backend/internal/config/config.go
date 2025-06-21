package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Config holds the application configuration
type Config struct {
	KeycloakURL      string
	KeycloakRealm    string
	KeycloakClientID string
	Port             string
	AllowedOrigins   []string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Use PORT environment variable for Cloud Run compatibility
	// Cloud Run sets PORT=8080 by default
	port := getEnv("PORT", "8080")

	// Validate required environment variables for production
	validateRequiredEnvVars()

	config := &Config{
		KeycloakURL:      getEnv("KEYCLOAK_URL", "http://localhost:8080"),
		KeycloakRealm:    getEnv("KEYCLOAK_REALM", "cloudgate"),
		KeycloakClientID: getEnv("KEYCLOAK_CLIENT_ID", "cloudgate-frontend"),
		Port:             port,
		AllowedOrigins:   strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
	}

	// Log configuration (excluding sensitive values)
	log.Printf("🔧 Configuration loaded:")
	log.Printf("   Port: %s", config.Port)
	log.Printf("   Keycloak URL: %s", config.KeycloakURL)
	log.Printf("   Keycloak Realm: %s", config.KeycloakRealm)
	log.Printf("   Keycloak Client ID: %s", config.KeycloakClientID)
	log.Printf("   Allowed Origins: %v", config.AllowedOrigins)

	return config
}

// validateRequiredEnvVars checks if required environment variables are set for production
func validateRequiredEnvVars() {
	// Only validate in Cloud Run environment (when PORT is set by platform)
	if os.Getenv("PORT") == "" {
		return // Skip validation for local development
	}

	required := []string{
		"KEYCLOAK_URL",
		"KEYCLOAK_REALM",
		"KEYCLOAK_CLIENT_ID",
	}

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

	if cfg.KeycloakURL == "" {
		return fmt.Errorf("keycloak URL cannot be empty")
	}

	if cfg.KeycloakRealm == "" {
		return fmt.Errorf("keycloak realm cannot be empty")
	}

	if cfg.KeycloakClientID == "" {
		return fmt.Errorf("keycloak client ID cannot be empty")
	}

	return nil
}
