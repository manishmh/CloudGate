package services

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"cloudgate-backend/internal/models"
	"cloudgate-backend/pkg/constants"
	"cloudgate-backend/pkg/types"

	"github.com/google/uuid"
)

var saasApps map[string]*types.SaaSApplication

// InitializeSaaSApps initializes the SaaS applications catalog
func InitializeSaaSApps() {
	saasApps = make(map[string]*types.SaaSApplication)

	// Google Workspace
	saasApps["google-workspace"] = &types.SaaSApplication{
		ID:          "google-workspace",
		Name:        "Google Workspace",
		Icon:        "🔍",
		Description: "Access Gmail, Drive, Calendar, and more",
		Category:    "productivity",
		Protocol:    "oauth2",
		Status:      "available",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	// Microsoft 365
	saasApps["microsoft-365"] = &types.SaaSApplication{
		ID:          "microsoft-365",
		Name:        "Microsoft 365",
		Icon:        "🏢",
		Description: "Access Outlook, OneDrive, Teams, and more",
		Category:    "productivity",
		Protocol:    "oauth2",
		Status:      "available",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	// Slack
	saasApps["slack"] = &types.SaaSApplication{
		ID:          "slack",
		Name:        "Slack",
		Icon:        "💬",
		Description: "Access your Slack workspaces",
		Category:    "communication",
		Protocol:    "oauth2",
		Status:      "available",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	// GitHub
	saasApps["github"] = &types.SaaSApplication{
		ID:          "github",
		Name:        "GitHub",
		Icon:        "🐙",
		Description: "Access your repositories and organizations",
		Category:    "development",
		Protocol:    "oauth2",
		Status:      "available",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	// Trello
	saasApps["trello"] = &types.SaaSApplication{
		ID:          "trello",
		Name:        "Trello",
		Icon:        "📋",
		Description: "Manage your boards and projects",
		Category:    "productivity",
		Protocol:    "oauth1",
		Status:      "available",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	// Salesforce
	saasApps["salesforce"] = &types.SaaSApplication{
		ID:          "salesforce",
		Name:        "Salesforce",
		Icon:        "☁️",
		Description: "Access your CRM and sales data",
		Category:    "crm",
		Protocol:    "oauth2",
		Status:      "available",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	// Jira
	saasApps["jira"] = &types.SaaSApplication{
		ID:          "jira",
		Name:        "Jira",
		Icon:        "🎯",
		Description: "Manage your projects and issues",
		Category:    "productivity",
		Protocol:    "oauth2",
		Status:      "available",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	// Notion
	saasApps["notion"] = &types.SaaSApplication{
		ID:          "notion",
		Name:        "Notion",
		Icon:        "📝",
		Description: "Access your workspace and documents",
		Category:    "productivity",
		Protocol:    "oauth2",
		Status:      "available",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	// Dropbox
	saasApps["dropbox"] = &types.SaaSApplication{
		ID:          "dropbox",
		Name:        "Dropbox",
		Icon:        "📦",
		Description: "Access your cloud storage",
		Category:    "storage",
		Protocol:    "oauth2",
		Status:      "available",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
	}
}

// GetAllSaaSApps returns all available SaaS applications
func GetAllSaaSApps() []*types.SaaSApplication {
	apps := make([]*types.SaaSApplication, 0, len(saasApps))
	for _, app := range saasApps {
		apps = append(apps, app)
	}
	return apps
}

// formatTimePtr formats a time pointer to string, returns empty if nil
func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

// buildMetadata creates a metadata map from database fields for backward compatibility
func buildMetadata(dbConn *models.AppConnection) map[string]string {
	metadata := make(map[string]string)

	if dbConn.UserEmail != "" {
		metadata["user_email"] = dbConn.UserEmail
	}
	if dbConn.UserName != "" {
		metadata["user_name"] = dbConn.UserName
	}
	if dbConn.Scopes != "" {
		metadata["scope"] = dbConn.Scopes
	}
	if dbConn.Provider != "" {
		metadata["provider"] = dbConn.Provider
	}
	if !dbConn.ConnectedAt.IsZero() {
		metadata["connected_at"] = dbConn.ConnectedAt.Format(time.RFC3339)
	}

	return metadata
}

// GetSaaSApp returns a specific SaaS application by ID
func GetSaaSApp(appID string) (*types.SaaSApplication, bool) {
	app, exists := saasApps[appID]
	return app, exists
}

// GetUserAppConnections returns all app connections for a user
func GetUserAppConnections(userID string) map[string]*types.UserAppConnection {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return make(map[string]*types.UserAppConnection)
	}

	var dbConnections []models.AppConnection
	DB.Where("user_id = ?", userUUID).Find(&dbConnections)

	connections := make(map[string]*types.UserAppConnection)
	for _, dbConn := range dbConnections {
		connections[dbConn.AppID] = &types.UserAppConnection{
			UserID:       dbConn.UserID.String(),
			AppID:        dbConn.AppID,
			Status:       dbConn.Status,
			AccessToken:  dbConn.AccessToken,
			RefreshToken: dbConn.RefreshToken,
			ExpiresAt:    formatTimePtr(dbConn.TokenExpiresAt),
			Metadata:     buildMetadata(&dbConn),
			ConnectedAt:  dbConn.ConnectedAt.Format(time.RFC3339),
			LastAccessAt: formatTimePtr(dbConn.LastUsed),
		}
	}
	return connections
}

// GetUserAppConnection returns a specific app connection for a user
func GetUserAppConnection(userID, appID string) (*types.UserAppConnection, bool) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, false
	}

	var dbConn models.AppConnection
	result := DB.Where("user_id = ? AND app_id = ?", userUUID, appID).First(&dbConn)
	if result.Error != nil {
		return nil, false
	}

	connection := &types.UserAppConnection{
		UserID:       dbConn.UserID.String(),
		AppID:        dbConn.AppID,
		Status:       dbConn.Status,
		AccessToken:  dbConn.AccessToken,
		RefreshToken: dbConn.RefreshToken,
		ExpiresAt:    formatTimePtr(dbConn.TokenExpiresAt),
		Metadata:     buildMetadata(&dbConn),
		ConnectedAt:  dbConn.ConnectedAt.Format(time.RFC3339),
		LastAccessAt: formatTimePtr(dbConn.LastUsed),
	}
	return connection, true
}

// CreateUserAppConnection creates a new app connection for a user
func CreateUserAppConnection(userID, appID string) *types.UserAppConnection {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil
	}

	now := time.Now().UTC()
	dbConn := models.AppConnection{
		UserID:      userUUID,
		AppID:       appID,
		Status:      constants.StatusPending,
		ConnectedAt: now,
	}

	DB.Create(&dbConn)

	return &types.UserAppConnection{
		UserID:      userID,
		AppID:       appID,
		Status:      constants.StatusPending,
		ConnectedAt: now.Format(time.RFC3339),
	}
}

// UpdateUserAppConnection updates an existing app connection or creates it if it doesn't exist
func UpdateUserAppConnection(userID, appID string, updates map[string]interface{}) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	var dbConn models.AppConnection
	result := DB.Where("user_id = ? AND app_id = ?", userUUID, appID).First(&dbConn)

	if result.Error != nil {
		// Create new connection if it doesn't exist
		now := time.Now().UTC()
		dbConn = models.AppConnection{
			UserID:      userUUID,
			AppID:       appID,
			Status:      constants.StatusPending,
			ConnectedAt: now,
		}
	}

	// Update fields from the updates map
	if status, ok := updates["status"].(string); ok {
		dbConn.Status = status
	}
	if accessToken, ok := updates["access_token"].(string); ok {
		dbConn.AccessToken = accessToken
	}
	if refreshToken, ok := updates["refresh_token"].(string); ok {
		dbConn.RefreshToken = refreshToken
	}
	if scopes, ok := updates["scope"].(string); ok {
		dbConn.Scopes = scopes
	}
	if expiresAtStr, ok := updates["expires_at"].(string); ok {
		if expiresAt, err := time.Parse(time.RFC3339, expiresAtStr); err == nil {
			dbConn.TokenExpiresAt = &expiresAt
		}
	}

	// Handle OAuth-specific fields
	if userEmail, ok := updates["user_email"].(string); ok {
		dbConn.UserEmail = userEmail
	}
	if userName, ok := updates["user_name"].(string); ok {
		dbConn.UserName = userName
	}
	if provider, ok := updates["provider"].(string); ok {
		dbConn.Provider = provider
	}
	if appName, ok := updates["app_name"].(string); ok {
		dbConn.AppName = appName
	}

	// Update last access time
	now := time.Now().UTC()
	dbConn.LastUsed = &now

	// Save to database
	if result.Error != nil {
		// Create new record
		return DB.Create(&dbConn).Error
	} else {
		// Update existing record
		return DB.Save(&dbConn).Error
	}
}

// GenerateState generates a random state string for OAuth
func GenerateState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetAppsWithUserStatus returns all apps with their connection status for a user
func GetAppsWithUserStatus(userID string) []*types.SaaSApplication {
	apps := GetAllSaaSApps()
	connections := GetUserAppConnections(userID)

	for _, app := range apps {
		if conn, exists := connections[app.ID]; exists {
			// Update app status based on connection
			app.Status = conn.Status
		}

	}

	return apps
}

// getFromMetadata safely gets a value from metadata map
func getFromMetadata(metadata map[string]string, key string) string {
	if metadata == nil {
		return ""
	}
	return metadata[key]
}
