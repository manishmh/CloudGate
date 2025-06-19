package config

import (
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
	port := getEnv("PORT", "")
	if port == "" {
		port = "8081" // Default for local development
	}

	config := &Config{
		KeycloakURL:      getEnv("KEYCLOAK_URL", "http://localhost:8080"),
		KeycloakRealm:    getEnv("KEYCLOAK_REALM", "cloudgate"),
		KeycloakClientID: getEnv("KEYCLOAK_CLIENT_ID", "cloudgate-frontend"),
		Port:             port,
		AllowedOrigins:   strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:3000"), ","),
	}
	return config
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
