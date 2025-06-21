package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cloudgate-backend/internal/services"
)

// AdaptiveAuthHandlers contains handlers for adaptive authentication
type AdaptiveAuthHandlers struct {
	adaptiveAuthService *services.AdaptiveAuthService
}

// NewAdaptiveAuthHandlers creates new adaptive auth handlers
func NewAdaptiveAuthHandlers(service *services.AdaptiveAuthService) *AdaptiveAuthHandlers {
	return &AdaptiveAuthHandlers{
		adaptiveAuthService: service,
	}
}

// EvaluateAuthenticationRequest represents the request for authentication evaluation
type EvaluateAuthenticationRequest struct {
	UserID            string                 `json:"user_id" binding:"required"`
	Email             string                 `json:"email" binding:"required"`
	IPAddress         string                 `json:"ip_address" binding:"required"`
	UserAgent         string                 `json:"user_agent" binding:"required"`
	DeviceFingerprint string                 `json:"device_fingerprint" binding:"required"`
	Location          *LocationRequest       `json:"location,omitempty"`
	SessionInfo       map[string]interface{} `json:"session_info,omitempty"`
	RequestHeaders    map[string]string      `json:"request_headers,omitempty"`
	ApplicationID     string                 `json:"application_id,omitempty"`
}

// LocationRequest represents location data in the request
type LocationRequest struct {
	Country     string  `json:"country"`
	Region      string  `json:"region"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ISP         string  `json:"isp"`
	Timezone    string  `json:"timezone"`
	VPNDetected bool    `json:"vpn_detected"`
}

// EvaluateAuthenticationResponse represents the response for authentication evaluation
type EvaluateAuthenticationResponse struct {
	Decision        string                    `json:"decision"`
	RiskScore       float64                   `json:"risk_score"`
	RiskLevel       string                    `json:"risk_level"`
	RequiredActions []AuthActionResponse      `json:"required_actions"`
	Reasoning       []string                  `json:"reasoning"`
	SessionDuration int64                     `json:"session_duration_seconds"`
	Restrictions    []AuthRestrictionResponse `json:"restrictions"`
	Metadata        map[string]interface{}    `json:"metadata"`
	ExpiresAt       time.Time                 `json:"expires_at"`
}

// AuthActionResponse represents an authentication action in the response
type AuthActionResponse struct {
	Type        string                 `json:"type"`
	Required    bool                   `json:"required"`
	Timeout     int64                  `json:"timeout_seconds"`
	Metadata    map[string]interface{} `json:"metadata"`
	Description string                 `json:"description"`
}

// AuthRestrictionResponse represents an access restriction in the response
type AuthRestrictionResponse struct {
	Type        string      `json:"type"`
	Value       interface{} `json:"value"`
	Description string      `json:"description"`
	ExpiresAt   *time.Time  `json:"expires_at,omitempty"`
}

// EvaluateAuthentication evaluates an authentication request
func (h *AdaptiveAuthHandlers) EvaluateAuthentication(c *gin.Context) {
	var req EvaluateAuthenticationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid UUID",
		})
		return
	}

	// Convert request to auth context
	authContext := &services.AuthContext{
		UserID:            userID,
		Email:             req.Email,
		IPAddress:         req.IPAddress,
		UserAgent:         req.UserAgent,
		DeviceFingerprint: req.DeviceFingerprint,
		LoginTime:         time.Now(),
		SessionInfo:       req.SessionInfo,
		RequestHeaders:    req.RequestHeaders,
		ApplicationID:     req.ApplicationID,
	}

	// Convert location if provided
	if req.Location != nil {
		authContext.Location = &services.GeoLocation{
			Country:     req.Location.Country,
			Region:      req.Location.Region,
			City:        req.Location.City,
			Latitude:    req.Location.Latitude,
			Longitude:   req.Location.Longitude,
			ISP:         req.Location.ISP,
			Timezone:    req.Location.Timezone,
			VPNDetected: req.Location.VPNDetected,
		}
	}

	// Evaluate authentication
	decision, err := h.adaptiveAuthService.EvaluateAuthentication(authContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Authentication evaluation failed",
			"message": err.Error(),
		})
		return
	}

	// Convert decision to response format
	response := &EvaluateAuthenticationResponse{
		Decision:        string(decision.Decision),
		RiskScore:       decision.RiskScore,
		RiskLevel:       decision.RiskLevel,
		RequiredActions: make([]AuthActionResponse, len(decision.RequiredActions)),
		Reasoning:       decision.Reasoning,
		SessionDuration: int64(decision.SessionDuration.Seconds()),
		Restrictions:    make([]AuthRestrictionResponse, len(decision.Restrictions)),
		Metadata:        decision.Metadata,
		ExpiresAt:       decision.ExpiresAt,
	}

	// Convert actions
	for i, action := range decision.RequiredActions {
		response.RequiredActions[i] = AuthActionResponse{
			Type:        string(action.Type),
			Required:    action.Required,
			Timeout:     int64(action.Timeout.Seconds()),
			Metadata:    action.Metadata,
			Description: action.Description,
		}
	}

	// Convert restrictions
	for i, restriction := range decision.Restrictions {
		response.Restrictions[i] = AuthRestrictionResponse{
			Type:        string(restriction.Type),
			Value:       restriction.Value,
			Description: restriction.Description,
			ExpiresAt:   restriction.ExpiresAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetRiskAssessmentHistory retrieves risk assessment history for a user
func (h *AdaptiveAuthHandlers) GetRiskAssessmentHistory(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing user ID",
			"message": "User ID is required",
		})
		return
	}

	// Validate user ID format
	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid UUID",
		})
		return
	}

	// Get limit parameter
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Get risk assessment history
	history, err := services.GetRiskAssessmentHistory(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve risk assessment history",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"limit":   limit,
		"history": history,
	})
}

// GetLatestRiskAssessment retrieves the latest risk assessment for a user
func (h *AdaptiveAuthHandlers) GetLatestRiskAssessment(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing user ID",
			"message": "User ID is required",
		})
		return
	}

	// Validate user ID format
	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid UUID",
		})
		return
	}

	// Get latest risk assessment
	assessment, err := services.GetLatestRiskAssessment(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve latest risk assessment",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, assessment)
}

// UpdateRiskThresholds updates the risk scoring thresholds
func (h *AdaptiveAuthHandlers) UpdateRiskThresholds(c *gin.Context) {
	var thresholds map[string]float64
	if err := c.ShouldBindJSON(&thresholds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	// Validate threshold values
	validThresholds := map[string]bool{
		"vpn_risk":         true,
		"tor_risk":         true,
		"new_device_risk":  true,
		"off_hours_risk":   true,
		"behavior_risk":    true,
		"location_risk":    true,
		"low_threshold":    true,
		"medium_threshold": true,
		"high_threshold":   true,
	}

	for key, value := range thresholds {
		if !validThresholds[key] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid threshold key",
				"message": "Unknown threshold: " + key,
			})
			return
		}

		if value < 0.0 || value > 1.0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid threshold value",
				"message": "Threshold values must be between 0.0 and 1.0",
			})
			return
		}
	}

	// Update thresholds
	err := services.UpdateRiskThresholds(thresholds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update risk thresholds",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Risk thresholds updated successfully",
		"thresholds": thresholds,
	})
}

// RegisterDeviceRequest represents a device registration request for adaptive auth
type AdaptiveAuthRegisterDeviceRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	Fingerprint string `json:"fingerprint" binding:"required"`
	DeviceName  string `json:"device_name"`
	DeviceType  string `json:"device_type"`
	Browser     string `json:"browser"`
	OS          string `json:"os"`
}

// RegisterDeviceFingerprint registers a new device fingerprint for a user
func (h *AdaptiveAuthHandlers) RegisterDeviceFingerprint(c *gin.Context) {
	var req AdaptiveAuthRegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	// Validate user ID
	if _, err := uuid.Parse(req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid UUID",
		})
		return
	}

	// Register device fingerprint
	err := services.RegisterDeviceFingerprint(
		req.UserID,
		req.Fingerprint,
		req.DeviceName,
		req.DeviceType,
		req.Browser,
		req.OS,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to register device fingerprint",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Device fingerprint registered successfully",
		"user_id":   req.UserID,
		"device_id": req.Fingerprint,
	})
}

// CheckDeviceStatus checks if a device is known for a user
func (h *AdaptiveAuthHandlers) CheckDeviceStatus(c *gin.Context) {
	userID := c.Query("user_id")
	fingerprint := c.Query("fingerprint")

	if userID == "" || fingerprint == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing parameters",
			"message": "Both user_id and fingerprint are required",
		})
		return
	}

	// Validate user ID
	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid UUID",
		})
		return
	}

	// Check if device is new
	isNew, err := services.IsNewDevice(userID, fingerprint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check device status",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":     userID,
		"fingerprint": fingerprint,
		"is_new":      isNew,
		"status": gin.H{
			"known":   !isNew,
			"trusted": false, // Would be determined by additional logic
		},
	})
}
