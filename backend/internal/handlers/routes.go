package handlers

import (
	"cloudgate-backend/internal/config"
	"cloudgate-backend/internal/middleware"
	"cloudgate-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the API routes for the application
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// Initialize services
	db := services.GetDB()
	userService := services.NewUserService(db)
	sessionService := services.NewSessionService(db)
	settingsService := services.NewUserSettingsService(db)
	adaptiveAuthService := services.NewAdaptiveAuthService(db)
	securityMonitoringService := services.NewSecurityMonitoringService(db)

	// Initialize handlers
	userHandlers := NewUserHandlers(userService, sessionService)
	settingsHandlers := NewSettingsHandlers(settingsService)
	dashboardHandlers := NewDashboardHandlers(userService, settingsService)
	adaptiveAuthHandlers := NewAdaptiveAuthHandlers(adaptiveAuthService)
	securityMonitoringHandlers := NewSecurityMonitoringHandlers(securityMonitoringService)

	// Add global OPTIONS handler for CORS preflight
	router.OPTIONS("/*cors", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(204)
	})

	// Health check endpoint
	router.GET("/health", HealthCheckHandler)
	router.GET("/health/db", DatabaseHealthCheckHandler)

	// Auth endpoints (JWT-based)
	router.POST("/auth/register", RegisterHandler(userService))
	router.POST("/auth/login", LoginHandler(userService, sessionService, cfg))
	router.POST("/auth/refresh", RefreshHandler(sessionService, cfg))
	router.POST("/auth/logout", LogoutHandler(sessionService))

	// API info endpoint
	router.GET("/api/info", APIInfoHandler)

	// Dashboard endpoints (protected)
	dashboardGroup := router.Group("/dashboard")
	dashboardGroup.Use(middleware.AuthenticationMiddleware())
	{
		dashboardGroup.GET("/data", dashboardHandlers.GetDashboardData)
		dashboardGroup.GET("/metrics", dashboardHandlers.GetDashboardMetrics)
	}

	// User profile endpoints
	userGroup := router.Group("/user")
	userGroup.Use(middleware.AuthenticationMiddleware())
	{
		userGroup.GET("/profile", userHandlers.GetProfile)
		userGroup.PUT("/profile", userHandlers.UpdateProfile)
		userGroup.POST("/email/verify", userHandlers.SendEmailVerification)
		userGroup.GET("/email/verify", userHandlers.VerifyEmail)
		userGroup.GET("/audit-logs", userHandlers.GetAuditLogs)
		userGroup.GET("/sessions", userHandlers.GetSessions)
		userGroup.DELETE("/sessions/:token", userHandlers.InvalidateSession)
		userGroup.DELETE("/sessions", userHandlers.InvalidateAllSessions)
		userGroup.DELETE("/account", userHandlers.DeactivateAccount)
	}

	// User settings endpoints
	userSettingsGroup := router.Group("/user/settings")
	userSettingsGroup.Use(middleware.AuthenticationMiddleware())
	{
		userSettingsGroup.GET("", settingsHandlers.GetUserSettings)
		userSettingsGroup.PUT("", settingsHandlers.UpdateUserSettings)
		userSettingsGroup.PUT("/single", settingsHandlers.UpdateSingleSetting)
		userSettingsGroup.POST("/reset", settingsHandlers.ResetUserSettings)
	}

	// MFA endpoints
	mfaGroup := router.Group("/user/mfa")
	mfaGroup.Use(middleware.AuthenticationMiddleware())
	{
		mfaGroup.GET("/status", GetMFAStatusHandler)
		mfaGroup.POST("/setup", SetupMFAHandler)
		mfaGroup.POST("/verify-setup", VerifyMFASetupHandler)
		mfaGroup.POST("/verify", VerifyMFAHandler)
		mfaGroup.POST("/disable", DisableMFAHandler)
		mfaGroup.POST("/backup-codes/regenerate", RegenerateBackupCodesHandler)
	}

	// OAuth Monitoring endpoints
	monitoringGroup := router.Group("/user/monitoring")
	monitoringGroup.Use(middleware.AuthenticationMiddleware())
	{
		// Connection monitoring
		monitoringGroup.GET("/connections", GetConnectionsHandler)
		monitoringGroup.GET("/connections/stats", GetConnectionStatsHandler)
		monitoringGroup.POST("/connections/:connectionId/test", TestConnectionHandler)
		monitoringGroup.POST("/connections/usage", RecordUsageHandler)

		// Security events
		monitoringGroup.GET("/security/events", GetSecurityEventsHandler)
		monitoringGroup.POST("/security/events", CreateSecurityEventHandler)

		// Device management
		monitoringGroup.GET("/devices", GetTrustedDevicesHandler)
		monitoringGroup.POST("/devices", RegisterDeviceHandler)
		monitoringGroup.PUT("/devices/:deviceId/trust", TrustDeviceHandler)
		monitoringGroup.DELETE("/devices/:deviceId", RevokeDeviceHandler)
	}

	// SaaS Applications endpoints (protected)
	appsGroup := router.Group("/apps")
	appsGroup.Use(middleware.AuthenticationMiddleware())
	{
		appsGroup.GET("", GetAppsHandler)
		appsGroup.POST("/connect", ConnectAppHandler)
		appsGroup.POST("/launch", LaunchAppHandler)
		appsGroup.GET("/callback", OAuthCallbackHandler)
	}

	// OAuth endpoints for real SaaS integrations (protected for user context)
	oauthGroup := router.Group("/oauth")
	oauthGroup.Use(middleware.AuthenticationMiddleware())
	{
		// Google OAuth (OAuth 2.0)
		oauthGroup.GET("/google/connect", GoogleOAuthInitHandler)
		oauthGroup.GET("/google/callback", GoogleOAuthCallbackHandler)

		// Microsoft OAuth (OAuth 2.0)
		oauthGroup.GET("/microsoft/connect", MicrosoftOAuthInitHandler)
		oauthGroup.GET("/microsoft/callback", MicrosoftOAuthCallbackHandler)

		// Slack OAuth (OAuth 2.0)
		oauthGroup.GET("/slack/connect", SlackOAuthInitHandler)
		oauthGroup.GET("/slack/callback", SlackOAuthCallbackHandler)

		// GitHub OAuth (OAuth 2.0)
		oauthGroup.GET("/github/connect", GitHubOAuthInitHandler)
		oauthGroup.GET("/github/callback", GitHubOAuthCallbackHandler)

		// Trello OAuth (OAuth 1.0a)
		oauthGroup.GET("/trello/connect", TrelloOAuthInitHandler)
		oauthGroup.GET("/trello/callback", TrelloOAuthCallbackHandler)

		// Salesforce OAuth (OAuth 2.0)
		oauthGroup.GET("/salesforce/connect", SalesforceOAuthInitHandler)
		oauthGroup.GET("/salesforce/callback", SalesforceOAuthCallbackHandler)
	}

	// Adaptive Authentication endpoints
	adaptiveAuthGroup := router.Group("/api/v1/adaptive-auth")
	adaptiveAuthGroup.Use(middleware.AuthenticationMiddleware())
	{
		adaptiveAuthGroup.POST("/evaluate", adaptiveAuthHandlers.EvaluateAuthentication)
		adaptiveAuthGroup.GET("/history/:userId", adaptiveAuthHandlers.GetRiskAssessmentHistory)
		adaptiveAuthGroup.GET("/latest/:userId", adaptiveAuthHandlers.GetLatestRiskAssessment)
		adaptiveAuthGroup.PUT("/thresholds", adaptiveAuthHandlers.UpdateRiskThresholds)
		adaptiveAuthGroup.POST("/register-device", adaptiveAuthHandlers.RegisterDeviceFingerprint)
		adaptiveAuthGroup.GET("/device-status", adaptiveAuthHandlers.CheckDeviceStatus)
	}

	// WebAuthn endpoints (protected)
	webauthnGroup := router.Group("/webauthn")
	webauthnGroup.Use(middleware.AuthenticationMiddleware())
	{
		webauthnGroup.GET("/credentials", GetWebAuthnCredentialsHandler)
		webauthnGroup.DELETE("/credentials/:credential_id", DeleteWebAuthnCredentialHandler)
	}

	// Security monitoring endpoints (protected)
	securityGroup := router.Group("/api/v1/security")
	securityGroup.Use(middleware.AuthenticationMiddleware())
	{
		// Map to implemented handlers
		securityGroup.POST("/alerts/generate", securityMonitoringHandlers.GenerateAlert)
		securityGroup.GET("/alerts", securityMonitoringHandlers.GetAlerts)
		securityGroup.GET("/metrics", securityMonitoringHandlers.GetSecurityMetrics)
	}
}
