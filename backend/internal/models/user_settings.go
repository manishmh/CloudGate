package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSettings represents user preferences and settings
type UserSettings struct {
	ID     uuid.UUID `gorm:"type:text;primary_key" json:"id"`
	UserID uuid.UUID `gorm:"type:text;not null;uniqueIndex" json:"user_id"`

	// Appearance
	Language   string `gorm:"default:'en'" json:"language"`
	Timezone   string `gorm:"default:'America/New_York'" json:"timezone"`
	DateFormat string `gorm:"default:'MM/DD/YYYY'" json:"date_format"`

	// Notifications
	EmailNotifications bool `gorm:"default:true" json:"email_notifications"`
	PushNotifications  bool `gorm:"default:false" json:"push_notifications"`
	SecurityAlerts     bool `gorm:"default:true" json:"security_alerts"`
	AppUpdates         bool `gorm:"default:true" json:"app_updates"`
	WeeklyReports      bool `gorm:"default:false" json:"weekly_reports"`

	// Security Settings
	TwoFactorEnabled         bool `gorm:"default:false" json:"two_factor_enabled"`
	LoginNotifications       bool `gorm:"default:true" json:"login_notifications"`
	SuspiciousActivityAlerts bool `gorm:"default:true" json:"suspicious_activity_alerts"`
	SessionTimeout           int  `gorm:"default:30" json:"session_timeout"` // in minutes
	PasswordExpiryDays       int  `gorm:"default:90" json:"password_expiry_days"`

	// Dashboard
	DefaultView     string `gorm:"default:'dashboard'" json:"default_view"`
	ItemsPerPage    int    `gorm:"default:10" json:"items_per_page"`
	AutoRefresh     bool   `gorm:"default:true" json:"auto_refresh"`
	RefreshInterval int    `gorm:"default:30" json:"refresh_interval"` // in seconds

	// Privacy
	AnalyticsOptIn  bool `gorm:"default:true" json:"analytics_opt_in"`
	ShareUsageData  bool `gorm:"default:false" json:"share_usage_data"`
	PersonalizedAds bool `gorm:"default:false" json:"personalized_ads"`

	// Integration
	APIAccess   bool   `gorm:"default:false" json:"api_access"`
	WebhookURL  string `json:"webhook_url,omitempty"`
	MaxAPICalls int    `gorm:"default:1000" json:"max_api_calls"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook to generate UUID
func (us *UserSettings) BeforeCreate(tx *gorm.DB) error {
	if us.ID == uuid.Nil {
		us.ID = uuid.New()
	}
	return nil
}

// GetDefaultSettings returns default user settings
func GetDefaultSettings() *UserSettings {
	return &UserSettings{
		Language:                 "en",
		Timezone:                 "America/New_York",
		DateFormat:               "MM/DD/YYYY",
		EmailNotifications:       true,
		PushNotifications:        false,
		SecurityAlerts:           true,
		AppUpdates:               true,
		WeeklyReports:            false,
		TwoFactorEnabled:         false,
		LoginNotifications:       true,
		SuspiciousActivityAlerts: true,
		SessionTimeout:           30,
		PasswordExpiryDays:       90,
		DefaultView:              "dashboard",
		ItemsPerPage:             10,
		AutoRefresh:              true,
		RefreshInterval:          30,
		AnalyticsOptIn:           true,
		ShareUsageData:           false,
		PersonalizedAds:          false,
		APIAccess:                false,
		WebhookURL:               "",
		MaxAPICalls:              1000,
	}
}
