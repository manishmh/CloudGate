package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppConnection represents a user's connection to a SaaS application
type AppConnection struct {
	ID       uuid.UUID `gorm:"type:text;primary_key" json:"id"`
	UserID   uuid.UUID `gorm:"type:text;not null;index" json:"user_id"`
	AppID    string    `gorm:"type:text;not null" json:"app_id"`
	AppName  string    `gorm:"type:text;not null" json:"app_name"`
	Provider string    `gorm:"type:text;not null" json:"provider"`
	Status   string    `gorm:"type:text;not null;default:'pending'" json:"status"` // pending, connected, error, revoked

	// OAuth specific fields
	AccessToken    string     `gorm:"type:text" json:"-"`
	RefreshToken   string     `gorm:"type:text" json:"-"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
	Scopes         string     `gorm:"type:text" json:"scopes"`

	// Connection details
	UserEmail   string     `gorm:"type:text" json:"user_email,omitempty"`
	UserName    string     `gorm:"type:text" json:"user_name,omitempty"`
	ConnectedAt time.Time  `json:"connected_at"`
	LastUsed    *time.Time `json:"last_used,omitempty"`

	// Health monitoring
	LastHealthCheck *time.Time `json:"last_health_check,omitempty"`
	HealthStatus    string     `gorm:"type:text;default:'unknown'" json:"health_status"` // healthy, warning, error, unknown
	ResponseTime    int        `gorm:"default:0" json:"response_time_ms"`
	ErrorCount      int        `gorm:"default:0" json:"error_count"`
	UptimePercent   float64    `gorm:"default:100.0" json:"uptime_percent"`

	// Usage statistics
	UsageCount      int64      `gorm:"default:0" json:"usage_count"`
	DataTransferred int64      `gorm:"default:0" json:"data_transferred_bytes"`
	LastError       string     `gorm:"type:text" json:"last_error,omitempty"`
	LastErrorAt     *time.Time `json:"last_error_at,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook to generate UUID
func (a *AppConnection) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.ConnectedAt.IsZero() {
		a.ConnectedAt = time.Now()
	}
	return nil
}

// ConnectionHealthMetrics represents health metrics for a connection
type ConnectionHealthMetrics struct {
	ConnectionID   uuid.UUID `gorm:"type:text;primary_key" json:"connection_id"`
	Timestamp      time.Time `gorm:"not null;index" json:"timestamp"`
	ResponseTime   int       `json:"response_time_ms"`
	Success        bool      `json:"success"`
	ErrorMessage   string    `gorm:"type:text" json:"error_message,omitempty"`
	HTTPStatusCode int       `json:"http_status_code,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate hook for ConnectionHealthMetrics
func (c *ConnectionHealthMetrics) BeforeCreate(tx *gorm.DB) error {
	if c.Timestamp.IsZero() {
		c.Timestamp = time.Now()
	}
	return nil
}

// SecurityEvent represents security-related events for connections
type SecurityEvent struct {
	ID           uuid.UUID  `gorm:"type:text;primary_key" json:"id"`
	UserID       uuid.UUID  `gorm:"type:text;not null;index" json:"user_id"`
	ConnectionID *uuid.UUID `gorm:"type:text;index" json:"connection_id,omitempty"`
	EventType    string     `gorm:"type:text;not null" json:"event_type"` // login, suspicious_location, new_device, failed_mfa, token_refresh
	Description  string     `gorm:"type:text;not null" json:"description"`
	Severity     string     `gorm:"type:text;not null" json:"severity"` // low, medium, high, critical
	IPAddress    string     `gorm:"type:text" json:"ip_address"`
	UserAgent    string     `gorm:"type:text" json:"user_agent"`
	Location     string     `gorm:"type:text" json:"location,omitempty"`
	RiskScore    float64    `gorm:"default:0.0" json:"risk_score"`
	Resolved     bool       `gorm:"default:false" json:"resolved"`
	ResolvedAt   *time.Time `json:"resolved_at,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User       User           `gorm:"foreignKey:UserID" json:"-"`
	Connection *AppConnection `gorm:"foreignKey:ConnectionID" json:"-"`
}

// BeforeCreate hook for SecurityEvent
func (s *SecurityEvent) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// TrustedDevice represents a user's trusted device
type TrustedDevice struct {
	ID          uuid.UUID `gorm:"type:text;primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:text;not null;index" json:"user_id"`
	DeviceName  string    `gorm:"type:text;not null" json:"device_name"`
	DeviceType  string    `gorm:"type:text;not null" json:"device_type"` // desktop, mobile, tablet
	Browser     string    `gorm:"type:text" json:"browser"`
	OS          string    `gorm:"type:text" json:"os"`
	Fingerprint string    `gorm:"type:text;not null;uniqueIndex" json:"fingerprint"`
	IPAddress   string    `gorm:"type:text" json:"ip_address"`
	Location    string    `gorm:"type:text" json:"location,omitempty"`
	Trusted     bool      `gorm:"default:false" json:"trusted"`
	LastSeen    time.Time `json:"last_seen"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook for TrustedDevice
func (t *TrustedDevice) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.LastSeen.IsZero() {
		t.LastSeen = time.Now()
	}
	return nil
}
