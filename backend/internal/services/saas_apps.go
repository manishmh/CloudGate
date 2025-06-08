package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"cloudgate-backend/pkg/constants"
	"cloudgate-backend/pkg/types"
)

// In-memory storage for demo purposes
// In production, this would be a database
var (
	saasApps        = make(map[string]*types.SaaSApplication)
	userConnections = make(map[string]map[string]*types.UserAppConnection) // userID -> appID -> connection
)

// InitializeSaaSApps sets up the predefined SaaS applications
func InitializeSaaSApps() {
	for _, appConfig := range constants.DefaultSaaSApps {
		app := &types.SaaSApplication{
			ID:          appConfig.ID,
			Name:        appConfig.Name,
			Icon:        appConfig.Icon,
			Description: appConfig.Description,
			Category:    appConfig.Category,
			Protocol:    appConfig.Protocol,
			Status:      appConfig.Status,
			Config:      appConfig.Config,
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

// GetSaaSApp returns a specific SaaS application by ID
func GetSaaSApp(appID string) (*types.SaaSApplication, bool) {
	app, exists := saasApps[appID]
	return app, exists
}

// GetUserAppConnections returns all app connections for a user
func GetUserAppConnections(userID string) map[string]*types.UserAppConnection {
	if connections, exists := userConnections[userID]; exists {
		return connections
	}
	return make(map[string]*types.UserAppConnection)
}

// GetUserAppConnection returns a specific app connection for a user
func GetUserAppConnection(userID, appID string) (*types.UserAppConnection, bool) {
	if connections, exists := userConnections[userID]; exists {
		if connection, exists := connections[appID]; exists {
			return connection, true
		}
	}
	return nil, false
}

// CreateUserAppConnection creates a new app connection for a user
func CreateUserAppConnection(userID, appID string) *types.UserAppConnection {
	if userConnections[userID] == nil {
		userConnections[userID] = make(map[string]*types.UserAppConnection)
	}

	connection := &types.UserAppConnection{
		UserID:      userID,
		AppID:       appID,
		Status:      constants.StatusPending,
		ConnectedAt: time.Now().UTC().Format(time.RFC3339),
	}

	userConnections[userID][appID] = connection
	return connection
}

// UpdateUserAppConnection updates an existing app connection
func UpdateUserAppConnection(userID, appID string, updates map[string]interface{}) error {
	connection, exists := GetUserAppConnection(userID, appID)
	if !exists {
		return fmt.Errorf(constants.ErrConnectionNotFound)
	}

	// Update fields
	if status, ok := updates["status"].(string); ok {
		connection.Status = status
	}
	if accessToken, ok := updates["access_token"].(string); ok {
		connection.AccessToken = accessToken
	}
	if refreshToken, ok := updates["refresh_token"].(string); ok {
		connection.RefreshToken = refreshToken
	}
	if expiresAt, ok := updates["expires_at"].(string); ok {
		connection.ExpiresAt = expiresAt
	}
	if metadata, ok := updates["metadata"].(map[string]string); ok {
		connection.Metadata = metadata
	}

	connection.LastAccessAt = time.Now().UTC().Format(time.RFC3339)
	return nil
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
