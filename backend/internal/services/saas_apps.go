package services

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"cloudgate-backend/internal/models"
	"cloudgate-backend/pkg/constants"
	"cloudgate-backend/pkg/types"
)

// In-memory storage for SaaS apps (these are static configurations)
var (
	saasApps = make(map[string]*types.SaaSApplication)
)

// InitializeSaaSApps sets up the predefined SaaS applications with environment-based configurations
func InitializeSaaSApps() {
	for _, appConfig := range constants.DefaultSaaSApps {
		// Create a copy of the config and update with environment variables
		config := make(map[string]string)
		for k, v := range appConfig.Config {
			config[k] = v
		}

		// Update config with environment variables based on app ID
		switch appConfig.ID {
		case "google-workspace":
			if clientID := getEnv("GOOGLE_CLIENT_ID", ""); clientID != "" {
				config["client_id"] = clientID
			}
			if clientSecret := getEnv("GOOGLE_CLIENT_SECRET", ""); clientSecret != "" {
				config["client_secret"] = clientSecret
			}
		case "microsoft-365":
			if clientID := getEnv("MICROSOFT_CLIENT_ID", ""); clientID != "" {
				config["client_id"] = clientID
			}
			if clientSecret := getEnv("MICROSOFT_CLIENT_SECRET", ""); clientSecret != "" {
				config["client_secret"] = clientSecret
			}
		case "slack":
			if clientID := getEnv("SLACK_CLIENT_ID", ""); clientID != "" {
				config["client_id"] = clientID
			}
			if clientSecret := getEnv("SLACK_CLIENT_SECRET", ""); clientSecret != "" {
				config["client_secret"] = clientSecret
			}
		case "github":
			if clientID := getEnv("GITHUB_CLIENT_ID", ""); clientID != "" {
				config["client_id"] = clientID
			}
			if clientSecret := getEnv("GITHUB_CLIENT_SECRET", ""); clientSecret != "" {
				config["client_secret"] = clientSecret
			}
		case "salesforce":
			if clientID := getEnv("SALESFORCE_CLIENT_ID", ""); clientID != "" {
				config["client_id"] = clientID
			}
			if clientSecret := getEnv("SALESFORCE_CLIENT_SECRET", ""); clientSecret != "" {
				config["client_secret"] = clientSecret
			}
		case "jira":
			if clientID := getEnv("JIRA_CLIENT_ID", ""); clientID != "" {
				config["client_id"] = clientID
			}
			if clientSecret := getEnv("JIRA_CLIENT_SECRET", ""); clientSecret != "" {
				config["client_secret"] = clientSecret
			}
		case "trello":
			if clientID := getEnv("TRELLO_CLIENT_ID", ""); clientID != "" {
				config["client_id"] = clientID
			}
			if clientSecret := getEnv("TRELLO_CLIENT_SECRET", ""); clientSecret != "" {
				config["client_secret"] = clientSecret
			}
		case "notion":
			if clientID := getEnv("NOTION_CLIENT_ID", ""); clientID != "" {
				config["client_id"] = clientID
			}
			if clientSecret := getEnv("NOTION_CLIENT_SECRET", ""); clientSecret != "" {
				config["client_secret"] = clientSecret
			}
		case "dropbox":
			if clientID := getEnv("DROPBOX_CLIENT_ID", ""); clientID != "" {
				config["client_id"] = clientID
			}
			if clientSecret := getEnv("DROPBOX_CLIENT_SECRET", ""); clientSecret != "" {
				config["client_secret"] = clientSecret
			}
		}

		app := &types.SaaSApplication{
			ID:          appConfig.ID,
			Name:        appConfig.Name,
			Icon:        appConfig.Icon,
			Description: appConfig.Description,
			Category:    appConfig.Category,
			Protocol:    appConfig.Protocol,
			Status:      appConfig.Status,
			Config:      config,
			CreatedAt:   time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:   time.Now().UTC().Format(time.RFC3339),
		}
		saasApps[app.ID] = app
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

// Helper functions for database operations

// formatTimePtr formats a time pointer to RFC3339 string or empty string if nil
func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

// buildMetadata creates a metadata map from database fields for backward compatibility
func buildMetadata(dbConn *models.UserAppConnection) map[string]string {
	metadata := make(map[string]string)

	if dbConn.UserEmail != "" {
		metadata["user_email"] = dbConn.UserEmail
	}
	if dbConn.UserName != "" {
		metadata["user_name"] = dbConn.UserName
	}
	if dbConn.Username != "" {
		metadata["username"] = dbConn.Username
	}
	if dbConn.TokenType != "" {
		metadata["token_type"] = dbConn.TokenType
	}
	if dbConn.Scope != "" {
		metadata["scope"] = dbConn.Scope
	}
	if dbConn.TeamName != "" {
		metadata["team_name"] = dbConn.TeamName
	}
	if dbConn.AccountID != "" {
		metadata["account_id"] = dbConn.AccountID
	}
	if dbConn.InstanceURL != "" {
		metadata["instance_url"] = dbConn.InstanceURL
	}
	if dbConn.BotID != "" {
		metadata["bot_id"] = dbConn.BotID
	}
	if dbConn.WorkspaceID != "" {
		metadata["workspace_id"] = dbConn.WorkspaceID
	}
	if dbConn.AccessTokenSecret != "" {
		metadata["access_token_secret"] = dbConn.AccessTokenSecret
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
	var dbConnections []models.UserAppConnection
	DB.Where("user_id = ?", userID).Find(&dbConnections)

	connections := make(map[string]*types.UserAppConnection)
	for _, dbConn := range dbConnections {
		connections[dbConn.AppID] = &types.UserAppConnection{
			UserID:       dbConn.UserID,
			AppID:        dbConn.AppID,
			Status:       dbConn.Status,
			AccessToken:  dbConn.AccessToken,
			RefreshToken: dbConn.RefreshToken,
			ExpiresAt:    formatTimePtr(dbConn.ExpiresAt),
			Metadata:     buildMetadata(&dbConn),
			ConnectedAt:  dbConn.ConnectedAt.Format(time.RFC3339),
			LastAccessAt: formatTimePtr(dbConn.LastAccessAt),
		}
	}
	return connections
}

// GetUserAppConnection returns a specific app connection for a user
func GetUserAppConnection(userID, appID string) (*types.UserAppConnection, bool) {
	var dbConn models.UserAppConnection
	result := DB.Where("user_id = ? AND app_id = ?", userID, appID).First(&dbConn)
	if result.Error != nil {
		return nil, false
	}

	connection := &types.UserAppConnection{
		UserID:       dbConn.UserID,
		AppID:        dbConn.AppID,
		Status:       dbConn.Status,
		AccessToken:  dbConn.AccessToken,
		RefreshToken: dbConn.RefreshToken,
		ExpiresAt:    formatTimePtr(dbConn.ExpiresAt),
		Metadata:     buildMetadata(&dbConn),
		ConnectedAt:  dbConn.ConnectedAt.Format(time.RFC3339),
		LastAccessAt: formatTimePtr(dbConn.LastAccessAt),
	}
	return connection, true
}

// CreateUserAppConnection creates a new app connection for a user
func CreateUserAppConnection(userID, appID string) *types.UserAppConnection {
	now := time.Now().UTC()
	dbConn := models.UserAppConnection{
		UserID:      userID,
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
	var dbConn models.UserAppConnection
	result := DB.Where("user_id = ? AND app_id = ?", userID, appID).First(&dbConn)

	if result.Error != nil {
		// Create new connection if it doesn't exist
		now := time.Now().UTC()
		dbConn = models.UserAppConnection{
			UserID:      userID,
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
	if tokenType, ok := updates["token_type"].(string); ok {
		dbConn.TokenType = tokenType
	}
	if scope, ok := updates["scope"].(string); ok {
		dbConn.Scope = scope
	}
	if expiresAtStr, ok := updates["expires_at"].(string); ok {
		if expiresAt, err := time.Parse(time.RFC3339, expiresAtStr); err == nil {
			dbConn.ExpiresAt = &expiresAt
		}
	}

	// Handle OAuth-specific fields
	if userEmail, ok := updates["user_email"].(string); ok {
		dbConn.UserEmail = userEmail
	}
	if userName, ok := updates["user_name"].(string); ok {
		dbConn.UserName = userName
	}
	if username, ok := updates["username"].(string); ok {
		dbConn.Username = username
	}
	if teamName, ok := updates["team_name"].(string); ok {
		dbConn.TeamName = teamName
	}
	if accountID, ok := updates["account_id"].(string); ok {
		dbConn.AccountID = accountID
	}
	if instanceURL, ok := updates["instance_url"].(string); ok {
		dbConn.InstanceURL = instanceURL
	}
	if botID, ok := updates["bot_id"].(string); ok {
		dbConn.BotID = botID
	}
	if workspaceID, ok := updates["workspace_id"].(string); ok {
		dbConn.WorkspaceID = workspaceID
	}
	if accessTokenSecret, ok := updates["access_token_secret"].(string); ok {
		dbConn.AccessTokenSecret = accessTokenSecret
	}

	// Update last access time
	now := time.Now().UTC()
	dbConn.LastAccessAt = &now

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

// GetAppsWithUserStatus returns apps with user connection status
func GetAppsWithUserStatus(userID string) []*types.SaaSApplication {
	apps := GetAllSaaSApps()
	userConnections := GetUserAppConnections(userID)

	// Update app status based on user connections
	for _, app := range apps {
		if connection, exists := userConnections[app.ID]; exists {
			if connection.Status == constants.StatusConnected {
				app.Status = constants.StatusConnected
			} else {
				app.Status = constants.StatusPending
			}
		} else {
			app.Status = constants.StatusAvailable
		}
	}

	return apps
}
