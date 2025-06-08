package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID                uuid.UUID      `gorm:"type:text;primary_key" json:"id"`
	KeycloakID        string         `gorm:"uniqueIndex;not null" json:"keycloak_id"`
	Email             string         `gorm:"uniqueIndex;not null" json:"email"`
	EmailVerified     bool           `gorm:"default:false" json:"email_verified"`
	EmailVerifiedAt   *time.Time     `json:"email_verified_at,omitempty"`
	Username          string         `gorm:"uniqueIndex;not null" json:"username"`
	FirstName         string         `json:"first_name"`
	LastName          string         `json:"last_name"`
	ProfilePictureURL string         `json:"profile_picture_url,omitempty"`
	LastLoginAt       *time.Time     `json:"last_login_at,omitempty"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Sessions  []Session  `gorm:"foreignKey:UserID" json:"-"`
	AuditLogs []AuditLog `gorm:"foreignKey:UserID" json:"-"`
	AppTokens []AppToken `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// Session represents a user session
type Session struct {
	ID           uuid.UUID `gorm:"type:text;primary_key" json:"id"`
	UserID       uuid.UUID `gorm:"type:text;not null;index" json:"user_id"`
	SessionToken string    `gorm:"uniqueIndex;not null" json:"session_token"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook to generate UUID
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// IsExpired checks if the session is expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// AppToken represents OAuth tokens for SaaS applications
type AppToken struct {
	ID           uuid.UUID  `gorm:"type:text;primary_key" json:"id"`
	UserID       uuid.UUID  `gorm:"type:text;not null;index" json:"user_id"`
	AppID        string     `gorm:"not null;index" json:"app_id"`
	AccessToken  string     `gorm:"type:text" json:"-"` // Encrypted in production
	RefreshToken string     `gorm:"type:text" json:"-"` // Encrypted in production
	TokenType    string     `gorm:"default:'Bearer'" json:"token_type"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	Scope        string     `json:"scope,omitempty"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook to generate UUID
func (a *AppToken) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// IsExpired checks if the token is expired
func (a *AppToken) IsExpired() bool {
	return a.ExpiresAt != nil && time.Now().After(*a.ExpiresAt)
}

// AuditLog represents audit trail for security and compliance
type AuditLog struct {
	ID         uuid.UUID  `gorm:"type:text;primary_key" json:"id"`
	UserID     *uuid.UUID `gorm:"type:text;index" json:"user_id,omitempty"`
	Action     string     `gorm:"not null;index" json:"action"`
	Resource   string     `gorm:"index" json:"resource,omitempty"`
	ResourceID string     `gorm:"index" json:"resource_id,omitempty"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`
	Details    string     `gorm:"type:text" json:"details,omitempty"`
	Status     string     `gorm:"default:'success'" json:"status"` // success, failure, warning
	CreatedAt  time.Time  `json:"created_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook to generate UUID
func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// EmailVerification represents email verification tokens
type EmailVerification struct {
	ID        uuid.UUID  `gorm:"type:text;primary_key" json:"id"`
	UserID    uuid.UUID  `gorm:"type:text;not null;index" json:"user_id"`
	Email     string     `gorm:"not null" json:"email"`
	Token     string     `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook to generate UUID
func (e *EmailVerification) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// IsExpired checks if the verification token is expired
func (e *EmailVerification) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// IsUsed checks if the verification token has been used
func (e *EmailVerification) IsUsed() bool {
	return e.UsedAt != nil
}
