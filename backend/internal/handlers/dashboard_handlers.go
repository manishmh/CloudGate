package handlers

import (
	"net/http"
	"time"

	"cloudgate-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DashboardHandlers contains dashboard-related HTTP handlers
type DashboardHandlers struct {
	userService     *services.UserService
	settingsService *services.UserSettingsService
}

// NewDashboardHandlers creates new dashboard handlers
func NewDashboardHandlers(userService *services.UserService, settingsService *services.UserSettingsService) *DashboardHandlers {
	return &DashboardHandlers{
		userService:     userService,
		settingsService: settingsService,
	}
}

// DashboardData represents the dashboard overview data
type DashboardData struct {
	User        UserProfile      `json:"user"`
	Metrics     DashboardMetrics `json:"metrics"`
	Connections []AppConnection  `json:"connections"`
	Activity    []ActivityItem   `json:"recent_activity"`
	Features    []FeatureCard    `json:"features"`
}

type UserProfile struct {
	ID                string `json:"id"`
	Email             string `json:"email"`
	Username          string `json:"username"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	ProfilePictureURL string `json:"profile_picture_url,omitempty"`
	LastLoginAt       string `json:"last_login_at,omitempty"`
}

type DashboardMetrics struct {
	TotalApps     int    `json:"total_apps"`
	ConnectedApps int    `json:"connected_apps"`
	RecentLogins  int    `json:"recent_logins"`
	SecurityScore int    `json:"security_score"`
	LastActivity  string `json:"last_activity"`
}

type AppConnection struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	ConnectURL  string `json:"connect_url"`
	LastUsed    string `json:"last_used,omitempty"`
}

type ActivityItem struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
	Icon        string `json:"icon"`
	Severity    string `json:"severity"`
}

type FeatureCard struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Stats       string   `json:"stats"`
	Color       string   `json:"color"`
	Features    []string `json:"features"`
}

// GetDashboardData retrieves comprehensive dashboard data
func (h *DashboardHandlers) GetDashboardData(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get user profile
	user, err := h.userService.GetUserByID(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	// Get user settings for personalization
	settings, err := h.settingsService.GetUserSettings(userID.(uuid.UUID))
	if err != nil {
		// If settings don't exist, create default ones
		settings, _ = h.settingsService.CreateDefaultSettings(userID.(uuid.UUID))
	}

	// Build dashboard data (settings can be used for personalization later)
	_ = settings // TODO: Use settings for dashboard personalization

	// Build dashboard data
	dashboardData := DashboardData{
		User: UserProfile{
			ID:                user.ID.String(),
			Email:             user.Email,
			Username:          user.Username,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			ProfilePictureURL: user.ProfilePictureURL,
			LastLoginAt:       formatTime(user.LastLoginAt),
		},
		Metrics: DashboardMetrics{
			TotalApps:     5,               // TODO: Get from actual app connections
			ConnectedApps: 3,               // TODO: Get from actual connected apps
			RecentLogins:  12,              // TODO: Get from audit logs
			SecurityScore: 95,              // TODO: Calculate based on security factors
			LastActivity:  "2 minutes ago", // TODO: Get from recent activity
		},
		Connections: getDefaultConnections(),
		Activity:    getRecentActivity(),
		Features:    getFeatureCards(),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dashboardData,
	})
}

// GetDashboardMetrics retrieves just the metrics for quick updates
func (h *DashboardHandlers) GetDashboardMetrics(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// TODO: Implement actual metrics calculation based on userID
	_ = userID // TODO: Use userID for personalized metrics
	metrics := DashboardMetrics{
		TotalApps:     5,
		ConnectedApps: 3,
		RecentLogins:  12,
		SecurityScore: 95,
		LastActivity:  "2 minutes ago",
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"metrics": metrics,
	})
}

// Helper functions
func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02T15:04:05Z")
}

func getDefaultConnections() []AppConnection {
	return []AppConnection{
		{
			Name:        "Google Workspace",
			Status:      "connected",
			Icon:        "🔍",
			Description: "Access Gmail, Drive, Calendar, and more",
			ConnectURL:  "/oauth/google/connect",
			LastUsed:    "2 hours ago",
		},
		{
			Name:        "Microsoft 365",
			Status:      "connected",
			Icon:        "🏢",
			Description: "Access Outlook, OneDrive, Teams, and more",
			ConnectURL:  "/oauth/microsoft/connect",
			LastUsed:    "1 day ago",
		},
		{
			Name:        "Slack",
			Status:      "disconnected",
			Icon:        "💬",
			Description: "Access your Slack workspaces",
			ConnectURL:  "/oauth/slack/connect",
		},
		{
			Name:        "GitHub",
			Status:      "connected",
			Icon:        "🐙",
			Description: "Access your repositories and organizations",
			ConnectURL:  "/oauth/github/connect",
			LastUsed:    "3 hours ago",
		},
		{
			Name:        "Trello",
			Status:      "disconnected",
			Icon:        "📋",
			Description: "Manage your boards and projects",
			ConnectURL:  "/oauth/trello/connect",
		},
	}
}

func getRecentActivity() []ActivityItem {
	return []ActivityItem{
		{
			ID:          "1",
			Type:        "login",
			Description: "Successful login from New York, NY",
			Timestamp:   "2 minutes ago",
			Icon:        "HiShieldCheck",
			Severity:    "success",
		},
		{
			ID:          "2",
			Type:        "app_launch",
			Description: "Launched Google Workspace",
			Timestamp:   "2 hours ago",
			Icon:        "HiViewGrid",
			Severity:    "info",
		},
		{
			ID:          "3",
			Type:        "connection",
			Description: "Connected to GitHub successfully",
			Timestamp:   "1 day ago",
			Icon:        "HiLink",
			Severity:    "success",
		},
		{
			ID:          "4",
			Type:        "security",
			Description: "Security scan completed - no issues found",
			Timestamp:   "2 days ago",
			Icon:        "HiShieldCheck",
			Severity:    "success",
		},
		{
			ID:          "5",
			Type:        "login",
			Description: "Failed login attempt blocked",
			Timestamp:   "3 days ago",
			Icon:        "HiExclamationCircle",
			Severity:    "warning",
		},
	}
}

func getFeatureCards() []FeatureCard {
	return []FeatureCard{
		{
			ID:          "adaptive-security",
			Title:       "Adaptive Security",
			Description: "AI-powered threat detection and response",
			Icon:        "🛡️",
			Stats:       "99.9% uptime",
			Color:       "from-blue-600 to-purple-600",
			Features:    []string{"Real-time monitoring", "Threat intelligence", "Auto-response"},
		},
		{
			ID:          "seamless-sso",
			Title:       "Seamless SSO",
			Description: "One-click access to all your applications",
			Icon:        "🔐",
			Stats:       "Sub-second login",
			Color:       "from-green-500 to-teal-600",
			Features:    []string{"SAML 2.0", "OAuth 2.0", "OpenID Connect"},
		},
		{
			ID:          "enterprise-ready",
			Title:       "Enterprise Ready",
			Description: "Scalable infrastructure for teams of any size",
			Icon:        "🏢",
			Stats:       "10K+ users",
			Color:       "from-orange-500 to-red-600",
			Features:    []string{"99.99% SLA", "24/7 support", "Global CDN"},
		},
		{
			ID:          "compliance",
			Title:       "Compliance First",
			Description: "Meet industry standards and regulations",
			Icon:        "📋",
			Stats:       "SOC 2 Type II",
			Color:       "from-purple-500 to-pink-600",
			Features:    []string{"GDPR ready", "HIPAA compliant", "ISO 27001"},
		},
	}
}
