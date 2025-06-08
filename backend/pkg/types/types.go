package types

// TokenIntrospectionResponse represents the response from Keycloak token introspection
type TokenIntrospectionResponse struct {
	Active            bool   `json:"active"`
	Username          string `json:"username,omitempty"`
	Email             string `json:"email,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Name              string `json:"name,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	Sub               string `json:"sub,omitempty"`
	Scope             string `json:"scope,omitempty"`
	ClientID          string `json:"client_id,omitempty"`
	TokenType         string `json:"token_type,omitempty"`
	Exp               int64  `json:"exp,omitempty"`
	Iat               int64  `json:"iat,omitempty"`
}

// UserSession represents the user session data
type UserSession struct {
	ID                string   `json:"id"`
	Email             string   `json:"email"`
	Name              string   `json:"name,omitempty"`
	PreferredUsername string   `json:"preferred_username,omitempty"`
	GivenName         string   `json:"given_name,omitempty"`
	FamilyName        string   `json:"family_name,omitempty"`
	Roles             []string `json:"roles"`
}

// TokenIntrospectionRequest represents the request for token introspection
type TokenIntrospectionRequest struct {
	Token string `json:"token" binding:"required"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Service   string `json:"service"`
}

// APIInfoResponse represents the API information response
type APIInfoResponse struct {
	Service     string   `json:"service"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Endpoints   []string `json:"endpoints"`
}

// SaaSApplication represents a SaaS application configuration
type SaaSApplication struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Icon        string            `json:"icon"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Protocol    string            `json:"protocol"` // "oauth2", "saml", "oidc"
	Status      string            `json:"status"`   // "available", "connected", "configured"
	LaunchURL   string            `json:"launch_url,omitempty"`
	Config      map[string]string `json:"config,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

// UserAppConnection represents a user's connection to a SaaS app
type UserAppConnection struct {
	UserID       string            `json:"user_id"`
	AppID        string            `json:"app_id"`
	Status       string            `json:"status"` // "connected", "disconnected", "pending"
	AccessToken  string            `json:"access_token,omitempty"`
	RefreshToken string            `json:"refresh_token,omitempty"`
	ExpiresAt    string            `json:"expires_at,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	ConnectedAt  string            `json:"connected_at"`
	LastAccessAt string            `json:"last_access_at,omitempty"`
}

// AppLaunchRequest represents a request to launch an application
type AppLaunchRequest struct {
	AppID string `json:"app_id" binding:"required"`
}

// AppLaunchResponse represents the response for launching an application
type AppLaunchResponse struct {
	LaunchURL string `json:"launch_url"`
	Method    string `json:"method"` // "redirect", "popup", "iframe"
	Token     string `json:"token,omitempty"`
	ExpiresIn int64  `json:"expires_in,omitempty"`
}

// AppConnectionRequest represents a request to connect to an application
type AppConnectionRequest struct {
	AppID string `json:"app_id" binding:"required"`
}

// AppConnectionResponse represents the response for connecting to an application
type AppConnectionResponse struct {
	AuthURL   string `json:"auth_url"`
	State     string `json:"state"`
	Challenge string `json:"challenge,omitempty"`
}

// OAuthCallbackRequest represents OAuth callback data
type OAuthCallbackRequest struct {
	AppID string `json:"app_id" binding:"required"`
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}
