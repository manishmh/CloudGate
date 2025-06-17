package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserAppConnection represents a user's connection to a SaaS application
type UserAppConnection struct {
	ID           uuid.UUID  `gorm:"type:text;primary_key" json:"id"`
	UserID       string     `gorm:"not null;index" json:"user_id"`
	AppID        string     `gorm:"not null;index" json:"app_id"`
	Status       string     `gorm:"not null;default:'pending'" json:"status"` // "connected", "disconnected", "pending"
	AccessToken  string     `gorm:"type:text" json:"access_token,omitempty"`
	RefreshToken string     `gorm:"type:text" json:"refresh_token,omitempty"`
	TokenType    string     `json:"token_type,omitempty"`
	Scope        string     `json:"scope,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`

	// User information from OAuth provider
	UserEmail string `json:"user_email,omitempty"`
	UserName  string `json:"user_name,omitempty"`
	Username  string `json:"username,omitempty"`

	// Provider-specific metadata
	TeamName    string `json:"team_name,omitempty"`
	AccountID   string `json:"account_id,omitempty"`
	InstanceURL string `json:"instance_url,omitempty"`
	BotID       string `json:"bot_id,omitempty"`
	WorkspaceID string `json:"workspace_id,omitempty"`

	// OAuth 1.0a specific
	AccessTokenSecret string `json:"access_token_secret,omitempty"`

	ConnectedAt  time.Time  `json:"connected_at"`
	LastAccessAt *time.Time `json:"last_access_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// BeforeCreate hook to generate UUID
func (uac *UserAppConnection) BeforeCreate(tx *gorm.DB) error {
	if uac.ID == uuid.Nil {
		uac.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for UserAppConnection
func (UserAppConnection) TableName() string {
	return "user_app_connections"
}
