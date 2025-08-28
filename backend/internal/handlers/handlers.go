package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cloudgate-backend/internal/services"
	"cloudgate-backend/pkg/types"
)

// HealthCheckHandler handles health check requests
func HealthCheckHandler(c *gin.Context) {
	log.Printf("🏥 Health Check Request from %s", c.ClientIP())

	response := types.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Service:   "cloudgate-backend",
	}

	log.Printf("✅ Health Check Response: %+v", response)
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
			"POST /auth/login - Login with email/password",
			"POST /auth/register - Register new user",
			"POST /auth/refresh - Refresh access token",
			"POST /auth/logout - Logout and revoke refresh token",
			"GET /api/info - API information",
			"GET /apps - List SaaS applications",
			"POST /apps/connect - Connect to a SaaS application",
			"POST /apps/launch - Launch a SaaS application",
			"POST /apps/callback - OAuth callback handler",
		},
	}
	c.JSON(http.StatusOK, response)
}

// Legacy Keycloak proxy endpoints removed during JWT migration.

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

	// Simulate OAuth connection initiation
	connectionURL := fmt.Sprintf("https://auth.%s.com/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s",
		app.ID, "your_client_id", "https://yourapp.com/oauth/callback", "read write")

	response := types.AppConnectionResponse{
		AuthURL: connectionURL,
		State:   "mock_state_value",
	}

	c.JSON(http.StatusOK, response)
}

// OAuthCallbackHandler handles OAuth callback from SaaS providers
func OAuthCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OAuth callback received",
		"code":    code,
		"state":   state,
	})
}

// LaunchAppHandler simulates launching a connected SaaS application
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

	// Simulate generating a temporary access token for app launch
	launchToken := uuid.New().String()

	response := types.AppLaunchResponse{
		LaunchURL: fmt.Sprintf("https://app.%s.com/dashboard?token=%s", request.AppID, launchToken),
		Method:    "redirect",
		Token:     launchToken,
		ExpiresIn: 300,
	}

	c.JSON(http.StatusOK, response)
}

// Helper function to extract user ID from request context
func getUserIDFromContext(c *gin.Context) string {
	userID, exists := c.Get("userID")
	if !exists {
		return ""
	}
	return userID.(uuid.UUID).String()
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
