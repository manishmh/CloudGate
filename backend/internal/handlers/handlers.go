package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cloudgate-backend/internal/config"
	"cloudgate-backend/internal/services"
	"cloudgate-backend/pkg/constants"
	"cloudgate-backend/pkg/types"
)

// HealthCheckHandler handles health check requests
func HealthCheckHandler(c *gin.Context) {
	response := types.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Service:   "cloudgate-backend",
	}
	c.JSON(http.StatusOK, response)
}

// APIInfoHandler provides information about the API endpoints
func APIInfoHandler(c *gin.Context) {
	response := types.APIInfoResponse{
		Service:     "CloudGate SSO Backend",
		Version:     "1.0.0",
		Description: "Enterprise SSO Portal Backend API",
		Endpoints: []string{
			"GET /health - Health check",
			"POST /token/introspect - Token introspection",
			"GET /user/info - User information",
			"GET /api/info - API information",
			"GET /apps - List SaaS applications",
			"POST /apps/connect - Connect to a SaaS application",
			"POST /apps/launch - Launch a SaaS application",
			"POST /apps/callback - OAuth callback handler",
		},
	}
	c.JSON(http.StatusOK, response)
}

// TokenIntrospectionHandler handles JWT token introspection requests
func TokenIntrospectionHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.TokenIntrospectionRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Prepare introspection request to Keycloak
		introspectionURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect",
			cfg.KeycloakURL, cfg.KeycloakRealm)

		data := url.Values{}
		data.Set("token", request.Token)
		data.Set("client_id", cfg.KeycloakClientID)

		req, err := http.NewRequest("POST", introspectionURL, strings.NewReader(data.Encode()))
		if err != nil {
			log.Printf("Error creating introspection request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making introspection request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to introspect token"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading introspection response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		var introspectionResp types.TokenIntrospectionResponse
		if err := json.Unmarshal(body, &introspectionResp); err != nil {
			log.Printf("Error parsing introspection response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			return
		}

		c.JSON(http.StatusOK, introspectionResp)
	}
}

// UserInfoHandler handles user information requests
func UserInfoHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Extract token from Bearer header
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		token := tokenParts[1]

		// Get user info from Keycloak
		userInfoURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo",
			cfg.KeycloakURL, cfg.KeycloakRealm)

		req, err := http.NewRequest("GET", userInfoURL, nil)
		if err != nil {
			log.Printf("Error creating userinfo request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making userinfo request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading userinfo response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		var userInfo map[string]interface{}
		if err := json.Unmarshal(body, &userInfo); err != nil {
			log.Printf("Error parsing userinfo response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			return
		}

		c.JSON(http.StatusOK, userInfo)
	}
}

// GetAppsHandler returns all SaaS applications with user connection status
func GetAppsHandler(c *gin.Context) {
	// Get user ID from token (simplified for demo)
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	apps := services.GetAppsWithUserStatus(userID)
	c.JSON(http.StatusOK, gin.H{
		"apps":  apps,
		"count": len(apps),
	})
}

// ConnectAppHandler initiates OAuth connection to a SaaS application
func ConnectAppHandler(c *gin.Context) {
	var request types.AppConnectionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get app configuration
	app, exists := services.GetSaaSApp(request.AppID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Create or update user connection
	services.CreateUserAppConnection(userID, request.AppID)

	// Generate OAuth URL
	state := services.GenerateState()
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=%s&response_type=code&state=%s",
		app.Config["auth_url"],
		app.Config["client_id"],
		url.QueryEscape("http://localhost:8081/apps/callback"),
		url.QueryEscape(app.Config["scope"]),
		state,
	)

	// Store state for validation (in production, use Redis or database)
	// For demo, we'll skip state validation

	response := types.AppConnectionResponse{
		AuthURL: authURL,
		State:   state,
	}

	c.JSON(http.StatusOK, response)
}

// LaunchAppHandler handles application launch requests
func LaunchAppHandler(c *gin.Context) {
	var request types.AppLaunchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get app configuration
	app, exists := services.GetSaaSApp(request.AppID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Check if user is connected to the app
	connection, exists := services.GetUserAppConnection(userID, request.AppID)
	if !exists || connection.Status != constants.StatusConnected {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not connected to this application"})
		return
	}

	// Get launch URL from constants
	launchURL, exists := constants.LaunchURLs[request.AppID]
	if !exists {
		launchURL = app.LaunchURL
	}

	// If still no launch URL, use a default
	if launchURL == "" {
		launchURL = "https://example.com"
	}

	// Update last access time
	services.UpdateUserAppConnection(userID, request.AppID, map[string]interface{}{
		"last_access_at": time.Now().UTC().Format(time.RFC3339),
	})

	response := types.AppLaunchResponse{
		LaunchURL: launchURL,
		Method:    "redirect",
		ExpiresIn: 3600,
	}

	c.JSON(http.StatusOK, response)
}

// OAuthCallbackHandler handles OAuth callbacks from SaaS applications
func OAuthCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	appID := c.Query("app_id")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state parameter"})
		return
	}

	// For demo purposes, we'll simulate successful OAuth completion
	// In production, you would exchange the code for tokens

	// Simulate finding the user (in production, you'd validate the state)
	userID := constants.DemoUserID // This would come from the state parameter

	// Update connection status
	err := services.UpdateUserAppConnection(userID, appID, map[string]interface{}{
		"status":       constants.StatusConnected,
		"access_token": constants.DemoAccessToken,
		"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update connection"})
		return
	}

	// Redirect back to frontend
	c.Redirect(http.StatusFound, "http://localhost:3000/dashboard?connected="+appID)
}

// Helper function to extract user ID from context
// This gets the user ID set by the authentication middleware
func getUserIDFromContext(c *gin.Context) string {
	// Try to get userID from context (set by authentication middleware)
	userIDInterface, exists := c.Get("userID")
	if exists {
		if userID, ok := userIDInterface.(uuid.UUID); ok {
			return userID.String()
		}
	}

	// Fallback: For endpoints without authentication middleware
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Simplified: just return demo user ID if auth header exists
	return constants.DemoUserID
}

// DatabaseHealthCheckHandler checks database connectivity
func DatabaseHealthCheckHandler(c *gin.Context) {
	if err := services.DatabaseHealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "healthy",
		"database": "connected",
	})
}

// AdminStatsHandler returns system statistics (placeholder)
func AdminStatsHandler(c *gin.Context) {
	// TODO: Implement admin authentication middleware
	sessionService := services.NewSessionService(services.GetDB())
	stats, err := sessionService.GetSessionStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

// AdminUsersHandler returns user list (placeholder)
func AdminUsersHandler(c *gin.Context) {
	// TODO: Implement admin authentication middleware
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin users endpoint - not implemented yet",
	})
}

// AdminSessionsHandler returns session list (placeholder)
func AdminSessionsHandler(c *gin.Context) {
	// TODO: Implement admin authentication middleware
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin sessions endpoint - not implemented yet",
	})
}
