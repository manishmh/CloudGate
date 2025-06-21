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

	// Legal pages

	// Token introspection endpoint
	router.POST("/token/introspect", TokenIntrospectionHandler(cfg))

	// User info endpoint
	router.GET("/user/info", UserInfoHandler(cfg))

	// API info endpoint
	router.GET("/api/info", APIInfoHandler)

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

	// SaaS Applications endpoints
	router.GET("/apps", GetAppsHandler)
	router.POST("/apps/connect", ConnectAppHandler)
	router.POST("/apps/launch", LaunchAppHandler)
	router.GET("/apps/callback", OAuthCallbackHandler)

	// OAuth endpoints for real SaaS integrations
	oauthGroup := router.Group("/oauth")
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

		// Jira OAuth (OAuth 2.0)
		oauthGroup.GET("/jira/connect", JiraOAuthInitHandler)
		oauthGroup.GET("/jira/callback", JiraOAuthCallbackHandler)

		// Notion OAuth (OAuth 2.0)
		oauthGroup.GET("/notion/connect", NotionOAuthInitHandler)
		oauthGroup.GET("/notion/callback", NotionOAuthCallbackHandler)

		// Dropbox OAuth (OAuth 2.0)
		oauthGroup.GET("/dropbox/connect", DropboxOAuthInitHandler)
		oauthGroup.GET("/dropbox/callback", DropboxOAuthCallbackHandler)
	}

	// Admin endpoints (for future use)
	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/stats", AdminStatsHandler)
		adminGroup.GET("/users", AdminUsersHandler)
		adminGroup.GET("/sessions", AdminSessionsHandler)
	}

	// Dashboard routes
	dashboardHandlers := NewDashboardHandlers(userService, settingsService)
	dashboardGroup := router.Group("/dashboard")
	dashboardGroup.Use(middleware.AuthenticationMiddleware())
	{
		dashboardGroup.GET("/data", dashboardHandlers.GetDashboardData)
		dashboardGroup.GET("/metrics", dashboardHandlers.GetDashboardMetrics)
	}

	// SAML SSO endpoints
	samlGroup := router.Group("/saml")
	{
		samlGroup.GET("/:app_id/init", SAMLInitHandler)
		samlGroup.POST("/:app_id/acs", SAMLACSHandler)
		samlGroup.GET("/metadata", SAMLMetadataHandler)
	}

	// WebAuthn endpoints
	webauthnGroup := router.Group("/webauthn")
	webauthnGroup.Use(middleware.AuthenticationMiddleware())
	{
		webauthnGroup.POST("/register/begin", WebAuthnRegistrationBeginHandler)
		webauthnGroup.POST("/register/finish", WebAuthnRegistrationFinishHandler)
		webauthnGroup.POST("/authenticate/begin", WebAuthnAuthenticationBeginHandler)
		webauthnGroup.POST("/authenticate/finish", WebAuthnAuthenticationFinishHandler)
		webauthnGroup.GET("/credentials", GetWebAuthnCredentialsHandler)
		webauthnGroup.DELETE("/credentials/:credential_id", DeleteWebAuthnCredentialHandler)
	}

	// Risk assessment endpoints
	riskGroup := router.Group("/risk")
	riskGroup.Use(middleware.AuthenticationMiddleware())
	{
		riskGroup.POST("/assess", AssessRiskHandler)
		riskGroup.GET("/policy", GetPolicyDecisionHandler)
		riskGroup.GET("/history", GetRiskHistoryHandler)
		riskGroup.PUT("/thresholds", UpdateRiskThresholdsHandler)
	}

	// Adaptive Authentication endpoints
	adaptiveAuthGroup := router.Group("/api/v1/adaptive-auth")
	adaptiveAuthGroup.Use(middleware.AuthenticationMiddleware())
	{
		// Core authentication evaluation
		adaptiveAuthGroup.POST("/evaluate", adaptiveAuthHandlers.EvaluateAuthentication)

		// Risk assessment history
		adaptiveAuthGroup.GET("/history/:user_id", adaptiveAuthHandlers.GetRiskAssessmentHistory)
		adaptiveAuthGroup.GET("/latest/:user_id", adaptiveAuthHandlers.GetLatestRiskAssessment)

		// Risk threshold management
		adaptiveAuthGroup.PUT("/thresholds", adaptiveAuthHandlers.UpdateRiskThresholds)

		// Device management
		adaptiveAuthGroup.POST("/register-device", adaptiveAuthHandlers.RegisterDeviceFingerprint)
		adaptiveAuthGroup.GET("/device-status", adaptiveAuthHandlers.CheckDeviceStatus)
	}

	// Security Monitoring & Alerting endpoints
	securityGroup := router.Group("/api/v1/security")
	securityGroup.Use(middleware.AuthenticationMiddleware())
	{
		// Alert management
		securityGroup.POST("/alerts", securityMonitoringHandlers.GenerateAlert)
		securityGroup.GET("/alerts", securityMonitoringHandlers.GetAlerts)
		securityGroup.PUT("/alerts/:alert_id/status", securityMonitoringHandlers.UpdateAlertStatus)

		// Incident management
		securityGroup.POST("/incidents", securityMonitoringHandlers.CreateIncident)
		securityGroup.GET("/incidents", securityMonitoringHandlers.GetIncidents)

		// Security metrics and monitoring
		securityGroup.GET("/metrics", securityMonitoringHandlers.GetSecurityMetrics)

		// Event processing
		securityGroup.POST("/events/login", securityMonitoringHandlers.ProcessLoginEvent)
		securityGroup.POST("/events/api", securityMonitoringHandlers.ProcessAPIEvent)

		// Alert channel configuration
		securityGroup.POST("/channels", securityMonitoringHandlers.ConfigureAlertChannel)

		// Reference data
		securityGroup.GET("/alert-types", securityMonitoringHandlers.GetAlertTypes)
		securityGroup.GET("/alert-severities", securityMonitoringHandlers.GetAlertSeverities)
	}

}
