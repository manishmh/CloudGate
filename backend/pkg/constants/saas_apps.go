package constants

import "time"

// SaaSAppConfig represents the configuration for a SaaS application
type SaaSAppConfig struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Icon        string            `json:"icon"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Protocol    string            `json:"protocol"`
	Status      string            `json:"status"`
	LaunchURL   string            `json:"launch_url,omitempty"`
	Config      map[string]string `json:"config,omitempty"`
}

// DefaultSaaSApps contains the predefined SaaS applications
var DefaultSaaSApps = []SaaSAppConfig{
	{
		ID:          "google-workspace",
		Name:        "Google Workspace",
		Icon:        "📧",
		Description: "Email, Drive, Calendar, and productivity tools",
		Category:    "productivity",
		Protocol:    "oauth2",
		Status:      "available",
		Config: map[string]string{
			"client_id":     "your-google-client-id",
			"client_secret": "your-google-client-secret",
			"scope":         "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile",
			"auth_url":      "https://accounts.google.com/o/oauth2/v2/auth",
			"token_url":     "https://oauth2.googleapis.com/token",
		},
	},
	{
		ID:          "microsoft-365",
		Name:        "Microsoft 365",
		Icon:        "📊",
		Description: "Office apps, Teams, SharePoint, and OneDrive",
		Category:    "productivity",
		Protocol:    "oauth2",
		Status:      "available",
		Config: map[string]string{
			"client_id":     "your-microsoft-client-id",
			"client_secret": "your-microsoft-client-secret",
			"scope":         "https://graph.microsoft.com/User.Read",
			"auth_url":      "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			"token_url":     "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		},
	},
	{
		ID:          "slack",
		Name:        "Slack",
		Icon:        "💬",
		Description: "Team communication and collaboration",
		Category:    "communication",
		Protocol:    "oauth2",
		Status:      "available",
		Config: map[string]string{
			"client_id":     "your-slack-client-id",
			"client_secret": "your-slack-client-secret",
			"scope":         "identity.basic,identity.email",
			"auth_url":      "https://slack.com/oauth/v2/authorize",
			"token_url":     "https://slack.com/api/oauth.v2.access",
		},
	},
	{
		ID:          "salesforce",
		Name:        "Salesforce",
		Icon:        "☁️",
		Description: "Customer relationship management platform",
		Category:    "crm",
		Protocol:    "oauth2",
		Status:      "available",
		Config: map[string]string{
			"client_id":     "your-salesforce-client-id",
			"client_secret": "your-salesforce-client-secret",
			"scope":         "api id",
			"auth_url":      "https://login.salesforce.com/services/oauth2/authorize",
			"token_url":     "https://login.salesforce.com/services/oauth2/token",
		},
	},
	{
		ID:          "jira",
		Name:        "Jira",
		Icon:        "🎯",
		Description: "Issue tracking and project management",
		Category:    "project-management",
		Protocol:    "oauth2",
		Status:      "available",
		Config: map[string]string{
			"client_id":     "your-jira-client-id",
			"client_secret": "your-jira-client-secret",
			"scope":         "read:jira-user read:jira-work",
			"auth_url":      "https://auth.atlassian.com/authorize",
			"token_url":     "https://auth.atlassian.com/oauth/token",
		},
	},
	{
		ID:          "trello",
		Name:        "Trello",
		Icon:        "📋",
		Description: "Project management and task organization",
		Category:    "project-management",
		Protocol:    "oauth1",
		Status:      "available",
		Config: map[string]string{
			"client_id":     "your-trello-client-id",
			"client_secret": "your-trello-client-secret",
			"scope":         "read,write",
			"auth_url":      "https://trello.com/1/OAuthAuthorizeToken",
			"token_url":     "https://trello.com/1/OAuthGetAccessToken",
		},
	},
	{
		ID:          "notion",
		Name:        "Notion",
		Icon:        "📝",
		Description: "All-in-one workspace for notes and collaboration",
		Category:    "productivity",
		Protocol:    "oauth2",
		Status:      "available",
		Config: map[string]string{
			"client_id":     "your-notion-client-id",
			"client_secret": "your-notion-client-secret",
			"scope":         "read",
			"auth_url":      "https://api.notion.com/v1/oauth/authorize",
			"token_url":     "https://api.notion.com/v1/oauth/token",
		},
	},
	{
		ID:          "github",
		Name:        "GitHub",
		Icon:        "🐙",
		Description: "Code repository and collaboration platform",
		Category:    "development",
		Protocol:    "oauth2",
		Status:      "available",
		Config: map[string]string{
			"client_id":     "your-github-client-id",
			"client_secret": "your-github-client-secret",
			"scope":         "user:email",
			"auth_url":      "https://github.com/login/oauth/authorize",
			"token_url":     "https://github.com/login/oauth/access_token",
		},
	},
	{
		ID:          "dropbox",
		Name:        "Dropbox",
		Icon:        "📦",
		Description: "Cloud storage and file synchronization",
		Category:    "storage",
		Protocol:    "oauth2",
		Status:      "available",
		Config: map[string]string{
			"client_id":     "your-dropbox-client-id",
			"client_secret": "your-dropbox-client-secret",
			"scope":         "account_info.read",
			"auth_url":      "https://www.dropbox.com/oauth2/authorize",
			"token_url":     "https://api.dropboxapi.com/oauth2/token",
		},
	},
}

// LaunchURLs contains the default launch URLs for applications
var LaunchURLs = map[string]string{
	"google-workspace": "https://workspace.google.com",
	"microsoft-365":    "https://office.com",
	"slack":            "https://slack.com/signin",
	"salesforce":       "https://login.salesforce.com",
	"jira":             "https://atlassian.net",
	"trello":           "https://trello.com",
	"notion":           "https://notion.so",
	"github":           "https://github.com",
	"dropbox":          "https://dropbox.com",
}

// Application status constants
const (
	StatusAvailable  = "available"
	StatusConnected  = "connected"
	StatusPending    = "pending"
	StatusError      = "error"
	StatusConfigured = "configured"
)

// OAuth protocol constants
const (
	ProtocolOAuth2 = "oauth2"
	ProtocolSAML   = "saml"
	ProtocolOIDC   = "oidc"
)

// Application categories
const (
	CategoryProductivity      = "productivity"
	CategoryCommunication     = "communication"
	CategoryCRM               = "crm"
	CategoryProjectManagement = "project-management"
	CategoryDocumentation     = "documentation"
	CategoryDevelopment       = "development"
	CategoryStorage           = "storage"
	CategoryAnalytics         = "analytics"
	CategorySecurity          = "security"
	CategoryFinance           = "finance"
)

// Default timeouts and limits
const (
	DefaultTokenExpiry    = 3600 * time.Second // 1 hour
	DefaultRefreshExpiry  = 24 * time.Hour     // 24 hours
	DefaultStateExpiry    = 10 * time.Minute   // 10 minutes
	MaxRetryAttempts      = 3
	DefaultRequestTimeout = 30 * time.Second
	MaxConnectionsPerUser = 50
)

// Demo configuration
const (
	DemoUserID       = "12345678-1234-1234-1234-123456789012"
	DemoAccessToken  = "demo-access-token"
	DemoRefreshToken = "demo-refresh-token"
)

// Error messages
const (
	ErrAppNotFound        = "application not found"
	ErrUserNotFound       = "user not found"
	ErrConnectionNotFound = "connection not found"
	ErrInvalidToken       = "invalid or expired token"
	ErrOAuthFailed        = "oauth authentication failed"
	ErrLaunchFailed       = "application launch failed"
)

// Success messages
const (
	MsgAppConnected     = "application connected successfully"
	MsgAppLaunched      = "application launched successfully"
	MsgConnectionUpdate = "connection updated successfully"
)
