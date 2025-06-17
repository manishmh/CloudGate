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

	// Initialize handlers
	userHandlers := NewUserHandlers(userService, sessionService)
	settingsHandlers := NewSettingsHandlers(settingsService)

	// Health check endpoint
	router.GET("/health", HealthCheckHandler)
	router.GET("/health/db", DatabaseHealthCheckHandler)

	// Legal pages
	router.GET("/privacy-policy", PrivacyPolicyHandler)
	router.GET("/terms", TermsHandler)

	// Token introspection endpoint
	router.POST("/token/introspect", TokenIntrospectionHandler(cfg))

	// User info endpoint
	router.GET("/user/info", UserInfoHandler(cfg))

	// API info endpoint
	router.GET("/api/info", APIInfoHandler)

	// User profile endpoints
	userGroup := router.Group("/user")
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
	{
		userSettingsGroup.GET("", settingsHandlers.GetUserSettings)
		userSettingsGroup.PUT("", settingsHandlers.UpdateUserSettings)
		userSettingsGroup.PUT("/single", settingsHandlers.UpdateSingleSetting)
		userSettingsGroup.POST("/reset", settingsHandlers.ResetUserSettings)
	}

	// MFA endpoints
	mfaGroup := router.Group("/user/mfa")
	{
		mfaGroup.GET("/status", GetMFAStatusHandler)
		mfaGroup.POST("/setup", SetupMFAHandler)
		mfaGroup.POST("/verify-setup", VerifyMFASetupHandler)
		mfaGroup.POST("/verify", VerifyMFAHandler)
		mfaGroup.POST("/disable", DisableMFAHandler)
		mfaGroup.POST("/backup-codes/regenerate", RegenerateBackupCodesHandler)
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

}
