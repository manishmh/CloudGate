package services

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"cloudgate-backend/internal/models"
)

// MFAService handles MFA operations
type MFAService struct {
	db *gorm.DB
}

// NewMFAService creates a new MFA service
func NewMFAService(db *gorm.DB) *MFAService {
	return &MFAService{db: db}
}

// StoreMFASetup stores MFA setup for a user
func StoreMFASetup(userID, secret string, backupCodes []string) error {
	db := GetDB()
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete existing MFA setup if any
	if err := tx.Where("user_id = ?", userUUID).Delete(&models.MFASetup{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete existing MFA setup: %w", err)
	}

	// Create new MFA setup
	mfaSetup := models.MFASetup{
		UserID:  userUUID,
		Secret:  secret,
		Enabled: false, // Not enabled until verified
	}

	if err := tx.Create(&mfaSetup).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create MFA setup: %w", err)
	}

	// Store backup codes
	for _, code := range backupCodes {
		hashedCode := hashBackupCode(code)
		backupCode := models.BackupCode{
			MFASetupID: mfaSetup.ID,
			Code:       hashedCode,
			Used:       false,
		}
		if err := tx.Create(&backupCode).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create backup code: %w", err)
		}
	}

	return tx.Commit().Error
}

// GetMFASetup retrieves MFA setup for a user
func GetMFASetup(userID string) (*models.MFASetup, error) {
	db := GetDB()
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var mfaSetup models.MFASetup
	if err := db.Where("user_id = ?", userUUID).First(&mfaSetup).Error; err != nil {
		return nil, err
	}

	return &mfaSetup, nil
}

// EnableMFA enables MFA for a user
func EnableMFA(userID string) error {
	db := GetDB()
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	return db.Model(&models.MFASetup{}).Where("user_id = ?", userUUID).Update("enabled", true).Error
}

// DisableMFA disables MFA for a user
func DisableMFA(userID string) error {
	db := GetDB()
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	return db.Model(&models.MFASetup{}).Where("user_id = ?", userUUID).Update("enabled", false).Error
}

// UseBackupCode checks if a backup code is valid and marks it as used
func UseBackupCode(userID, code string) (bool, error) {
	db := GetDB()
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get MFA setup
	var mfaSetup models.MFASetup
	if err := db.Where("user_id = ?", userUUID).First(&mfaSetup).Error; err != nil {
		return false, err
	}

	// Hash the provided code
	hashedCode := hashBackupCode(code)

	// Check if backup code exists and is not used
	var backupCode models.BackupCode
	if err := db.Where("mfa_setup_id = ? AND code = ? AND used = ?", mfaSetup.ID, hashedCode, false).First(&backupCode).Error; err != nil {
		return false, nil // Code not found or already used
	}

	// Mark backup code as used
	now := time.Now()
	backupCode.Used = true
	backupCode.UsedAt = &now

	if err := db.Save(&backupCode).Error; err != nil {
		return false, fmt.Errorf("failed to mark backup code as used: %w", err)
	}

	return true, nil
}

// GetBackupCodesCount returns the number of unused backup codes for a user
func GetBackupCodesCount(userID string) (int, error) {
	db := GetDB()
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get MFA setup
	var mfaSetup models.MFASetup
	if err := db.Where("user_id = ?", userUUID).First(&mfaSetup).Error; err != nil {
		return 0, err
	}

	// Count unused backup codes
	var count int64
	if err := db.Model(&models.BackupCode{}).Where("mfa_setup_id = ? AND used = ?", mfaSetup.ID, false).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// ReplaceBackupCodes replaces all backup codes for a user
func ReplaceBackupCodes(userID string, newCodes []string) error {
	db := GetDB()
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Get MFA setup
	var mfaSetup models.MFASetup
	if err := db.Where("user_id = ?", userUUID).First(&mfaSetup).Error; err != nil {
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete existing backup codes
	if err := tx.Where("mfa_setup_id = ?", mfaSetup.ID).Delete(&models.BackupCode{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete existing backup codes: %w", err)
	}

	// Create new backup codes
	for _, code := range newCodes {
		hashedCode := hashBackupCode(code)
		backupCode := models.BackupCode{
			MFASetupID: mfaSetup.ID,
			Code:       hashedCode,
			Used:       false,
		}
		if err := tx.Create(&backupCode).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create backup code: %w", err)
		}
	}

	return tx.Commit().Error
}

// LogAuditEvent logs an audit event (wrapper for existing audit logging)
func LogAuditEvent(userID, action, resource, resourceID, ipAddress, userAgent, details, status string) {
	db := GetDB()

	var userUUID *uuid.UUID
	if userID != "" {
		if parsed, err := uuid.Parse(userID); err == nil {
			userUUID = &parsed
		}
	}

	auditLog := models.AuditLog{
		UserID:     userUUID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Details:    details,
		Status:     status,
	}

	if err := db.Create(&auditLog).Error; err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to create audit log: %v\n", err)
	}
}

// hashBackupCode creates a hash of the backup code for secure storage
func hashBackupCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return fmt.Sprintf("%x", hash)
}
