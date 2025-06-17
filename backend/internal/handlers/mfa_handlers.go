package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"

	"cloudgate-backend/internal/services"
)

// MFA Setup Response
type MFASetupResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	QRCodeData  string   `json:"qr_code_data_url"`
	BackupCodes []string `json:"backup_codes"`
}

// MFA Verification Request
type MFAVerifyRequest struct {
	Code string `json:"code" binding:"required"`
}

// MFA Status Response
type MFAStatusResponse struct {
	Enabled     bool    `json:"enabled"`
	SetupDate   *string `json:"setup_date,omitempty"`
	BackupCodes int     `json:"backup_codes_remaining"`
}

// SetupMFAHandler initiates MFA setup for a user
func SetupMFAHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if MFA is already enabled
	userService := services.NewUserService(services.GetDB())
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := userService.GetUserByID(userUUID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Generate TOTP secret
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "CloudGate SSO",
		AccountName: user.Email,
		SecretSize:  32,
	})
	if err != nil {
		log.Printf("Error generating TOTP secret: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate MFA secret"})
		return
	}

	// Generate QR code
	qrCodePNG, err := qrcode.Encode(secret.URL(), qrcode.Medium, 256)
	if err != nil {
		log.Printf("Error generating QR code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Convert QR code to data URL
	qrCodeDataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrCodePNG)

	// Generate backup codes
	backupCodes, err := generateBackupCodes(10)
	if err != nil {
		log.Printf("Error generating backup codes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate backup codes"})
		return
	}

	// Store MFA setup in database (but don't enable yet - requires verification)
	err = services.StoreMFASetup(userID, secret.Secret(), backupCodes)
	if err != nil {
		log.Printf("Error storing MFA setup: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store MFA setup"})
		return
	}

	response := MFASetupResponse{
		Secret:      secret.Secret(),
		QRCodeURL:   secret.URL(),
		QRCodeData:  qrCodeDataURL,
		BackupCodes: backupCodes,
	}

	c.JSON(http.StatusOK, response)
}

// VerifyMFASetupHandler verifies and enables MFA for a user
func VerifyMFASetupHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request MFAVerifyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get MFA setup from database
	mfaSetup, err := services.GetMFASetup(userID)
	if err != nil {
		log.Printf("Error getting MFA setup: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "MFA setup not found"})
		return
	}

	// Verify TOTP code
	valid := totp.Validate(request.Code, mfaSetup.Secret)
	if !valid {
		// Check if it's a backup code
		valid, err = services.UseBackupCode(userID, request.Code)
		if err != nil {
			log.Printf("Error checking backup code: %v", err)
		}

		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
			return
		}
	}

	// Enable MFA for the user
	err = services.EnableMFA(userID)
	if err != nil {
		log.Printf("Error enabling MFA: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable MFA"})
		return
	}

	// Log MFA enablement
	services.LogAuditEvent(userID, "mfa_enabled", "user", userID, c.ClientIP(), c.GetHeader("User-Agent"), "MFA successfully enabled", "success")

	c.JSON(http.StatusOK, gin.H{
		"message": "MFA enabled successfully",
		"enabled": true,
	})
}

// VerifyMFAHandler verifies MFA code during login
func VerifyMFAHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request MFAVerifyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get MFA setup from database
	mfaSetup, err := services.GetMFASetup(userID)
	if err != nil {
		log.Printf("Error getting MFA setup: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "MFA not configured"})
		return
	}

	if !mfaSetup.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MFA not enabled"})
		return
	}

	// Verify TOTP code
	valid := totp.Validate(request.Code, mfaSetup.Secret)
	if !valid {
		// Check if it's a backup code
		valid, err = services.UseBackupCode(userID, request.Code)
		if err != nil {
			log.Printf("Error checking backup code: %v", err)
		}

		if !valid {
			// Log failed MFA attempt
			services.LogAuditEvent(userID, "mfa_verification_failed", "user", userID, c.ClientIP(), c.GetHeader("User-Agent"), "MFA verification failed", "failure")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
			return
		}
	}

	// Log successful MFA verification
	services.LogAuditEvent(userID, "mfa_verification_success", "user", userID, c.ClientIP(), c.GetHeader("User-Agent"), "MFA verification successful", "success")

	c.JSON(http.StatusOK, gin.H{
		"message":  "MFA verification successful",
		"verified": true,
	})
}

// GetMFAStatusHandler returns MFA status for a user
func GetMFAStatusHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get MFA setup from database
	mfaSetup, err := services.GetMFASetup(userID)
	if err != nil {
		// MFA not configured
		c.JSON(http.StatusOK, MFAStatusResponse{
			Enabled:     false,
			BackupCodes: 0,
		})
		return
	}

	// Get backup codes count
	backupCodesCount, err := services.GetBackupCodesCount(userID)
	if err != nil {
		log.Printf("Error getting backup codes count: %v", err)
		backupCodesCount = 0
	}

	var setupDate *string
	if !mfaSetup.CreatedAt.IsZero() {
		dateStr := mfaSetup.CreatedAt.Format(time.RFC3339)
		setupDate = &dateStr
	}

	response := MFAStatusResponse{
		Enabled:     mfaSetup.Enabled,
		SetupDate:   setupDate,
		BackupCodes: backupCodesCount,
	}

	c.JSON(http.StatusOK, response)
}

// DisableMFAHandler disables MFA for a user
func DisableMFAHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request MFAVerifyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get MFA setup from database
	mfaSetup, err := services.GetMFASetup(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MFA not configured"})
		return
	}

	if !mfaSetup.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MFA not enabled"})
		return
	}

	// Verify TOTP code before disabling
	valid := totp.Validate(request.Code, mfaSetup.Secret)
	if !valid {
		// Check if it's a backup code
		valid, err = services.UseBackupCode(userID, request.Code)
		if err != nil {
			log.Printf("Error checking backup code: %v", err)
		}

		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
			return
		}
	}

	// Disable MFA
	err = services.DisableMFA(userID)
	if err != nil {
		log.Printf("Error disabling MFA: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable MFA"})
		return
	}

	// Log MFA disablement
	services.LogAuditEvent(userID, "mfa_disabled", "user", userID, c.ClientIP(), c.GetHeader("User-Agent"), "MFA disabled", "warning")

	c.JSON(http.StatusOK, gin.H{
		"message": "MFA disabled successfully",
		"enabled": false,
	})
}

// RegenerateBackupCodesHandler generates new backup codes
func RegenerateBackupCodesHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request MFAVerifyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get MFA setup from database
	mfaSetup, err := services.GetMFASetup(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MFA not configured"})
		return
	}

	if !mfaSetup.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MFA not enabled"})
		return
	}

	// Verify TOTP code before regenerating
	valid := totp.Validate(request.Code, mfaSetup.Secret)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	// Generate new backup codes
	backupCodes, err := generateBackupCodes(10)
	if err != nil {
		log.Printf("Error generating backup codes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate backup codes"})
		return
	}

	// Store new backup codes
	err = services.ReplaceBackupCodes(userID, backupCodes)
	if err != nil {
		log.Printf("Error storing backup codes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store backup codes"})
		return
	}

	// Log backup codes regeneration
	services.LogAuditEvent(userID, "backup_codes_regenerated", "user", userID, c.ClientIP(), c.GetHeader("User-Agent"), "Backup codes regenerated", "success")

	c.JSON(http.StatusOK, gin.H{
		"message":      "Backup codes regenerated successfully",
		"backup_codes": backupCodes,
	})
}

// Helper function to generate backup codes
func generateBackupCodes(count int) ([]string, error) {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		bytes := make([]byte, 5) // 5 bytes = 10 character hex string
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}
		codes[i] = fmt.Sprintf("%X", bytes)
	}
	return codes, nil
}
