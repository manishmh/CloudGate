package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MFASetup represents MFA configuration for a user
type MFASetup struct {
	ID        uuid.UUID      `gorm:"type:text;primary_key" json:"id"`
	UserID    uuid.UUID      `gorm:"type:text;not null;index" json:"user_id"`
	Secret    string         `gorm:"type:text;not null" json:"-"` // TOTP secret, encrypted in production
	Enabled   bool           `gorm:"default:false" json:"enabled"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User        User         `gorm:"foreignKey:UserID" json:"-"`
	BackupCodes []BackupCode `gorm:"foreignKey:MFASetupID" json:"-"`
}

// BeforeCreate hook to generate UUID
func (m *MFASetup) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// BackupCode represents backup codes for MFA
type BackupCode struct {
	ID         uuid.UUID      `gorm:"type:text;primary_key" json:"id"`
	MFASetupID uuid.UUID      `gorm:"type:text;not null;index" json:"mfa_setup_id"`
	Code       string         `gorm:"type:text;not null;uniqueIndex" json:"-"` // Hashed backup code
	Used       bool           `gorm:"default:false" json:"used"`
	UsedAt     *time.Time     `json:"used_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	MFASetup MFASetup `gorm:"foreignKey:MFASetupID" json:"-"`
}

// BeforeCreate hook to generate UUID
func (b *BackupCode) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
