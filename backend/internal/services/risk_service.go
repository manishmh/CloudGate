package services

import (
	"cloudgate-backend/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RiskAssessment represents a risk assessment record
type RiskAssessment struct {
	ID                uuid.UUID `gorm:"type:text;primary_key" json:"id"`
	UserID            uuid.UUID `gorm:"type:text;not null;index" json:"user_id"`
	SessionID         string    `gorm:"type:text" json:"session_id"`
	IPAddress         string    `gorm:"type:text" json:"ip_address"`
	UserAgent         string    `gorm:"type:text" json:"user_agent"`
	Location          string    `gorm:"type:text" json:"location"` // JSON serialized LocationInfo
	DeviceFingerprint string    `gorm:"type:text" json:"device_fingerprint"`
	BehaviorSignals   string    `gorm:"type:text" json:"behavior_signals"` // JSON serialized
	RiskScore         float64   `gorm:"not null" json:"risk_score"`
	RiskLevel         string    `gorm:"type:text;not null" json:"risk_level"`
	Factors           string    `gorm:"type:text" json:"risk_factors"`    // JSON serialized
	Recommendations   string    `gorm:"type:text" json:"recommendations"` // JSON serialized
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Relationships
	User models.User `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook for RiskAssessment
func (r *RiskAssessment) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// RiskThresholds represents configurable risk scoring thresholds
type RiskThresholds struct {
	ID              uuid.UUID `gorm:"type:text;primary_key" json:"id"`
	VPNRisk         float64   `gorm:"default:0.3" json:"vpn_risk"`
	TorRisk         float64   `gorm:"default:0.9" json:"tor_risk"`
	NewDeviceRisk   float64   `gorm:"default:0.7" json:"new_device_risk"`
	OffHoursRisk    float64   `gorm:"default:0.4" json:"off_hours_risk"`
	BehaviorRisk    float64   `gorm:"default:0.5" json:"behavior_risk"`
	LocationRisk    float64   `gorm:"default:0.6" json:"location_risk"`
	LowThreshold    float64   `gorm:"default:0.3" json:"low_threshold"`
	MediumThreshold float64   `gorm:"default:0.6" json:"medium_threshold"`
	HighThreshold   float64   `gorm:"default:0.8" json:"high_threshold"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BeforeCreate hook for RiskThresholds
func (r *RiskThresholds) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// DeviceFingerprint represents a stored device fingerprint
type DeviceFingerprint struct {
	ID          uuid.UUID `gorm:"type:text;primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:text;not null;index" json:"user_id"`
	Fingerprint string    `gorm:"type:text;not null" json:"fingerprint"`
	DeviceName  string    `gorm:"type:text" json:"device_name"`
	DeviceType  string    `gorm:"type:text" json:"device_type"`
	Browser     string    `gorm:"type:text" json:"browser"`
	OS          string    `gorm:"type:text" json:"os"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
	IsTrusted   bool      `gorm:"default:false" json:"is_trusted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	User models.User `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook for DeviceFingerprint
func (d *DeviceFingerprint) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// WebAuthnCredential represents a WebAuthn credential
type WebAuthnCredential struct {
	ID                uuid.UUID  `gorm:"type:text;primary_key" json:"id"`
	UserID            uuid.UUID  `gorm:"type:text;not null;index" json:"user_id"`
	CredentialID      string     `gorm:"type:text;not null;uniqueIndex" json:"credential_id"`
	PublicKey         []byte     `gorm:"type:bytea" json:"public_key"`
	AttestationObject []byte     `gorm:"type:bytea" json:"attestation_object"`
	Counter           uint32     `gorm:"default:0" json:"counter"`
	DeviceName        string     `gorm:"type:text" json:"device_name"`
	CreatedAt         time.Time  `json:"created_at"`
	LastUsed          *time.Time `json:"last_used,omitempty"`

	// Relationships
	User models.User `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook for WebAuthnCredential
func (w *WebAuthnCredential) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

// StoreRiskAssessment stores a risk assessment in the database
func StoreRiskAssessment(assessment interface{}) error {
	db := GetDB()

	// Convert the assessment to our internal structure
	assessmentData, err := json.Marshal(assessment)
	if err != nil {
		return fmt.Errorf("failed to marshal assessment: %w", err)
	}

	var assessmentMap map[string]interface{}
	if err := json.Unmarshal(assessmentData, &assessmentMap); err != nil {
		return fmt.Errorf("failed to unmarshal assessment: %w", err)
	}

	// Extract user ID
	userIDStr, ok := assessmentMap["user_id"].(string)
	if !ok {
		return fmt.Errorf("invalid user_id in assessment")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("invalid user_id format: %w", err)
	}

	// Serialize complex fields
	locationJSON, _ := json.Marshal(assessmentMap["location"])
	behaviorJSON, _ := json.Marshal(assessmentMap["behavior_signals"])
	factorsJSON, _ := json.Marshal(assessmentMap["risk_factors"])
	recommendationsJSON, _ := json.Marshal(assessmentMap["recommendations"])

	// Create database record
	riskAssessment := RiskAssessment{
		UserID:            userID,
		SessionID:         getStringField(assessmentMap, "session_id"),
		IPAddress:         getStringField(assessmentMap, "ip_address"),
		UserAgent:         getStringField(assessmentMap, "user_agent"),
		Location:          string(locationJSON),
		DeviceFingerprint: getStringField(assessmentMap, "device_fingerprint"),
		BehaviorSignals:   string(behaviorJSON),
		RiskScore:         getFloatField(assessmentMap, "risk_score"),
		RiskLevel:         getStringField(assessmentMap, "risk_level"),
		Factors:           string(factorsJSON),
		Recommendations:   string(recommendationsJSON),
	}

	return db.Create(&riskAssessment).Error
}

// GetLatestRiskAssessment retrieves the latest risk assessment for a user
func GetLatestRiskAssessment(userID string) (interface{}, error) {
	db := GetDB()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var assessment RiskAssessment
	err = db.Where("user_id = ?", userUUID).
		Order("created_at DESC").
		First(&assessment).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get latest risk assessment: %w", err)
	}

	// Convert back to the expected format
	result := map[string]interface{}{
		"user_id":            assessment.UserID.String(),
		"session_id":         assessment.SessionID,
		"ip_address":         assessment.IPAddress,
		"user_agent":         assessment.UserAgent,
		"device_fingerprint": assessment.DeviceFingerprint,
		"risk_score":         assessment.RiskScore,
		"risk_level":         assessment.RiskLevel,
		"timestamp":          assessment.CreatedAt,
	}

	// Deserialize JSON fields
	if assessment.Location != "" {
		var location interface{}
		json.Unmarshal([]byte(assessment.Location), &location)
		result["location"] = location
	}

	if assessment.BehaviorSignals != "" {
		var behavior interface{}
		json.Unmarshal([]byte(assessment.BehaviorSignals), &behavior)
		result["behavior_signals"] = behavior
	}

	if assessment.Factors != "" {
		var factors interface{}
		json.Unmarshal([]byte(assessment.Factors), &factors)
		result["risk_factors"] = factors
	}

	if assessment.Recommendations != "" {
		var recommendations interface{}
		json.Unmarshal([]byte(assessment.Recommendations), &recommendations)
		result["recommendations"] = recommendations
	}

	return result, nil
}

// GetRiskAssessmentHistory retrieves risk assessment history for a user
func GetRiskAssessmentHistory(userID string, limit int) ([]interface{}, error) {
	db := GetDB()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var assessments []RiskAssessment
	query := db.Where("user_id = ?", userUUID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err = query.Find(&assessments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get risk assessment history: %w", err)
	}

	// Convert to expected format
	results := make([]interface{}, len(assessments))
	for i, assessment := range assessments {
		result := map[string]interface{}{
			"user_id":            assessment.UserID.String(),
			"session_id":         assessment.SessionID,
			"ip_address":         assessment.IPAddress,
			"user_agent":         assessment.UserAgent,
			"device_fingerprint": assessment.DeviceFingerprint,
			"risk_score":         assessment.RiskScore,
			"risk_level":         assessment.RiskLevel,
			"timestamp":          assessment.CreatedAt,
		}

		// Deserialize JSON fields
		if assessment.Location != "" {
			var location interface{}
			json.Unmarshal([]byte(assessment.Location), &location)
			result["location"] = location
		}

		results[i] = result
	}

	return results, nil
}

// UpdateRiskThresholds updates risk scoring thresholds
func UpdateRiskThresholds(thresholds map[string]float64) error {
	db := GetDB()

	// Get or create risk thresholds record
	var riskThresholds RiskThresholds
	err := db.First(&riskThresholds).Error
	if err == gorm.ErrRecordNotFound {
		// Create new record with defaults
		riskThresholds = RiskThresholds{}
	} else if err != nil {
		return fmt.Errorf("failed to get risk thresholds: %w", err)
	}

	// Update thresholds
	for key, value := range thresholds {
		switch key {
		case "vpn_risk":
			riskThresholds.VPNRisk = value
		case "tor_risk":
			riskThresholds.TorRisk = value
		case "new_device_risk":
			riskThresholds.NewDeviceRisk = value
		case "off_hours_risk":
			riskThresholds.OffHoursRisk = value
		case "behavior_risk":
			riskThresholds.BehaviorRisk = value
		case "location_risk":
			riskThresholds.LocationRisk = value
		case "low_threshold":
			riskThresholds.LowThreshold = value
		case "medium_threshold":
			riskThresholds.MediumThreshold = value
		case "high_threshold":
			riskThresholds.HighThreshold = value
		}
	}

	return db.Save(&riskThresholds).Error
}

// IsNewDevice checks if a device fingerprint is new for a user
func IsNewDevice(userID, deviceFingerprint string) (bool, error) {
	if deviceFingerprint == "" {
		return true, nil // No fingerprint means unknown device
	}

	db := GetDB()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return true, fmt.Errorf("invalid user ID: %w", err)
	}

	var count int64
	err = db.Model(&DeviceFingerprint{}).
		Where("user_id = ? AND fingerprint = ?", userUUID, deviceFingerprint).
		Count(&count).Error

	if err != nil {
		return true, fmt.Errorf("failed to check device fingerprint: %w", err)
	}

	return count == 0, nil
}

// RegisterDeviceFingerprint registers a new device fingerprint
func RegisterDeviceFingerprint(userID, fingerprint, deviceName, deviceType, browser, os string) error {
	db := GetDB()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if already exists
	var existing DeviceFingerprint
	err = db.Where("user_id = ? AND fingerprint = ?", userUUID, fingerprint).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new fingerprint record
		deviceFP := DeviceFingerprint{
			UserID:      userUUID,
			Fingerprint: fingerprint,
			DeviceName:  deviceName,
			DeviceType:  deviceType,
			Browser:     browser,
			OS:          os,
			FirstSeen:   time.Now(),
			LastSeen:    time.Now(),
		}
		return db.Create(&deviceFP).Error
	} else if err != nil {
		return fmt.Errorf("failed to check existing fingerprint: %w", err)
	} else {
		// Update last seen
		return db.Model(&existing).Update("last_seen", time.Now()).Error
	}
}

// WebAuthn credential management functions

// GetUserWebAuthnCredentials retrieves WebAuthn credentials for a user
func GetUserWebAuthnCredentials(userID string) ([]WebAuthnCredential, error) {
	db := GetDB()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var credentials []WebAuthnCredential
	err = db.Where("user_id = ?", userUUID).Find(&credentials).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get WebAuthn credentials: %w", err)
	}

	return credentials, nil
}

// StoreWebAuthnCredential stores a WebAuthn credential
func StoreWebAuthnCredential(userID, credentialID string, attestationObject []byte) error {
	db := GetDB()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	credential := WebAuthnCredential{
		UserID:            userUUID,
		CredentialID:      credentialID,
		AttestationObject: attestationObject,
		DeviceName:        "WebAuthn Device",
	}

	return db.Create(&credential).Error
}

// VerifyWebAuthnCredential verifies if a WebAuthn credential exists
func VerifyWebAuthnCredential(userID, credentialID string) (bool, error) {
	db := GetDB()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user ID: %w", err)
	}

	var count int64
	err = db.Model(&WebAuthnCredential{}).
		Where("user_id = ? AND credential_id = ?", userUUID, credentialID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to verify WebAuthn credential: %w", err)
	}

	return count > 0, nil
}

// UpdateWebAuthnCredentialUsage updates the last used timestamp for a credential
func UpdateWebAuthnCredentialUsage(userID, credentialID string) error {
	db := GetDB()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	now := time.Now()
	return db.Model(&WebAuthnCredential{}).
		Where("user_id = ? AND credential_id = ?", userUUID, credentialID).
		Update("last_used", &now).Error
}

// DeleteWebAuthnCredential deletes a WebAuthn credential
func DeleteWebAuthnCredential(userID, credentialID string) error {
	db := GetDB()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	return db.Where("user_id = ? AND credential_id = ?", userUUID, credentialID).
		Delete(&WebAuthnCredential{}).Error
}

// Helper functions
func getStringField(data map[string]interface{}, key string) string {
	if value, exists := data[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

func getFloatField(data map[string]interface{}, key string) float64 {
	if value, exists := data[key]; exists {
		if f, ok := value.(float64); ok {
			return f
		}
	}
	return 0.0
}
