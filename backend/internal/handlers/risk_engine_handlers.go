package handlers

import (
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"cloudgate-backend/internal/services"
)

// Risk scoring structures
type RiskAssessment struct {
	UserID            string          `json:"user_id"`
	SessionID         string          `json:"session_id"`
	IPAddress         string          `json:"ip_address"`
	UserAgent         string          `json:"user_agent"`
	Location          LocationInfo    `json:"location"`
	DeviceFingerprint string          `json:"device_fingerprint"`
	BehaviorSignals   BehaviorSignals `json:"behavior_signals"`
	RiskScore         float64         `json:"risk_score"`
	RiskLevel         string          `json:"risk_level"`
	Factors           []RiskFactor    `json:"risk_factors"`
	Recommendations   []string        `json:"recommendations"`
	Timestamp         time.Time       `json:"timestamp"`
}

type LocationInfo struct {
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	ISP       string  `json:"isp"`
	IsVPN     bool    `json:"is_vpn"`
	IsTor     bool    `json:"is_tor"`
	IsProxy   bool    `json:"is_proxy"`
}

type BehaviorSignals struct {
	TypingPattern     TypingPattern     `json:"typing_pattern"`
	MouseMovement     MouseMovement     `json:"mouse_movement"`
	NavigationPattern NavigationPattern `json:"navigation_pattern"`
	TimePatterns      TimePatterns      `json:"time_patterns"`
}

type TypingPattern struct {
	AvgKeydownTime float64   `json:"avg_keydown_time"`
	AvgKeyupTime   float64   `json:"avg_keyup_time"`
	TypingRhythm   float64   `json:"typing_rhythm"`
	PausePatterns  []float64 `json:"pause_patterns"`
}

type MouseMovement struct {
	AvgSpeed        float64 `json:"avg_speed"`
	ClickFrequency  float64 `json:"click_frequency"`
	MovementPattern string  `json:"movement_pattern"`
	ScrollBehavior  float64 `json:"scroll_behavior"`
}

type NavigationPattern struct {
	PageSequence    []string `json:"page_sequence"`
	SessionDuration float64  `json:"session_duration"`
	ClickDepth      int      `json:"click_depth"`
	BackButtonUsage int      `json:"back_button_usage"`
}

type TimePatterns struct {
	LoginTime        string  `json:"login_time"`
	TypicalHours     []int   `json:"typical_hours"`
	WeekdayPattern   []int   `json:"weekday_pattern"`
	SessionFrequency float64 `json:"session_frequency"`
}

type RiskFactor struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Weight      float64 `json:"weight"`
	Score       float64 `json:"score"`
	Severity    string  `json:"severity"`
}

type PolicyDecision struct {
	Action        string          `json:"action"` // allow, deny, step_up, monitor
	Confidence    float64         `json:"confidence"`
	RequiredMFA   []string        `json:"required_mfa,omitempty"`
	SessionLimits SessionLimits   `json:"session_limits,omitempty"`
	Monitoring    MonitoringLevel `json:"monitoring"`
	Explanation   string          `json:"explanation"`
}

type SessionLimits struct {
	MaxDuration   int  `json:"max_duration_minutes"`
	IdleTimeout   int  `json:"idle_timeout_minutes"`
	RequireReauth bool `json:"require_reauth"`
	LimitedAccess bool `json:"limited_access"`
}

type MonitoringLevel struct {
	Level            string `json:"level"` // none, basic, enhanced, full
	LogAllActions    bool   `json:"log_all_actions"`
	RealTimeAlerts   bool   `json:"real_time_alerts"`
	BehaviorTracking bool   `json:"behavior_tracking"`
}

// AssessRiskHandler performs comprehensive risk assessment
func AssessRiskHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse request body for additional context
	var contextData map[string]interface{}
	if err := c.ShouldBindJSON(&contextData); err != nil {
		// If no body provided, continue with request headers only
		contextData = make(map[string]interface{})
	}

	// Gather risk signals
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Get device fingerprint from request
	deviceFingerprint := ""
	if fp, exists := contextData["device_fingerprint"]; exists {
		if fpStr, ok := fp.(string); ok {
			deviceFingerprint = fpStr
		}
	}

	// Perform geolocation lookup
	location := performGeolocation(ipAddress)

	// Get behavior signals
	behaviorSignals := extractBehaviorSignals(contextData)

	// Calculate risk score
	assessment := calculateRiskScore(userID, ipAddress, userAgent, deviceFingerprint, location, behaviorSignals)

	// Store assessment for future reference
	err := services.StoreRiskAssessment(assessment)
	if err != nil {
		log.Printf("Error storing risk assessment: %v", err)
		// Don't fail the request for this
	}

	c.JSON(http.StatusOK, assessment)
}

// GetPolicyDecisionHandler returns policy decision based on risk assessment
func GetPolicyDecisionHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get latest risk assessment
	assessment, err := services.GetLatestRiskAssessment(userID)
	if err != nil {
		log.Printf("Error getting risk assessment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get risk assessment"})
		return
	}

	// Convert assessment to our internal type for policy decision
	assessmentMap, ok := assessment.(map[string]interface{})
	if !ok {
		log.Printf("Error converting assessment to map")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process risk assessment"})
		return
	}

	// Create a proper RiskAssessment struct for policy decision
	riskScore, _ := assessmentMap["risk_score"].(float64)
	riskLevel, _ := assessmentMap["risk_level"].(string)

	internalAssessment := RiskAssessment{
		RiskScore: riskScore,
		RiskLevel: riskLevel,
	}

	// Make policy decision
	decision := makePolicyDecision(internalAssessment)

	// Log policy decision
	services.LogAuditEvent(userID, "policy_decision", "security", userID, c.ClientIP(), c.GetHeader("User-Agent"),
		fmt.Sprintf("Policy decision: %s (risk: %.2f)", decision.Action, riskScore), "info")

	c.JSON(http.StatusOK, decision)
}

// GetRiskHistoryHandler returns risk assessment history for a user
func GetRiskHistoryHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	assessments, err := services.GetRiskAssessmentHistory(userID, limit)
	if err != nil {
		log.Printf("Error getting risk history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get risk history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"assessments": assessments,
		"count":       len(assessments),
	})
}

// UpdateRiskThresholdsHandler updates risk scoring thresholds
func UpdateRiskThresholdsHandler(c *gin.Context) {
	// This would typically be an admin-only endpoint
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var thresholds map[string]float64
	if err := c.ShouldBindJSON(&thresholds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := services.UpdateRiskThresholds(thresholds)
	if err != nil {
		log.Printf("Error updating risk thresholds: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update thresholds"})
		return
	}

	// Log threshold update
	services.LogAuditEvent(userID, "risk_thresholds_updated", "security", "global", c.ClientIP(), c.GetHeader("User-Agent"),
		"Risk scoring thresholds updated", "info")

	c.JSON(http.StatusOK, gin.H{
		"message":    "Risk thresholds updated successfully",
		"thresholds": thresholds,
	})
}

// Helper functions

func performGeolocation(ipAddress string) LocationInfo {
	// Simplified geolocation - in production, use MaxMind GeoIP2 or similar service
	location := LocationInfo{
		Country:  "Unknown",
		Region:   "Unknown",
		City:     "Unknown",
		Timezone: "UTC",
		ISP:      "Unknown",
		IsVPN:    false,
		IsTor:    false,
		IsProxy:  false,
	}

	// Check for private/local IPs
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return location
	}

	if ip.IsPrivate() || ip.IsLoopback() {
		location.Country = "Local"
		location.Region = "Private Network"
		location.City = "Local"
		return location
	}

	// Mock geolocation data based on IP patterns
	if strings.HasPrefix(ipAddress, "192.168.") || strings.HasPrefix(ipAddress, "10.") {
		location.Country = "US"
		location.Region = "California"
		location.City = "San Francisco"
		location.Latitude = 37.7749
		location.Longitude = -122.4194
		location.Timezone = "America/Los_Angeles"
		location.ISP = "Local Network"
	}

	// In production, integrate with MaxMind GeoIP2:
	// db, err := geoip2.Open("GeoLite2-City.mmdb")
	// record, err := db.City(ip)
	// location.Country = record.Country.IsoCode
	// location.City = record.City.Names["en"]
	// etc.

	return location
}

func extractBehaviorSignals(contextData map[string]interface{}) BehaviorSignals {
	signals := BehaviorSignals{
		TypingPattern: TypingPattern{
			AvgKeydownTime: 100.0,
			AvgKeyupTime:   50.0,
			TypingRhythm:   1.0,
		},
		MouseMovement: MouseMovement{
			AvgSpeed:        150.0,
			ClickFrequency:  2.0,
			MovementPattern: "normal",
		},
		NavigationPattern: NavigationPattern{
			SessionDuration: 300.0,
			ClickDepth:      5,
		},
		TimePatterns: TimePatterns{
			LoginTime:        time.Now().Format("15:04"),
			TypicalHours:     []int{9, 10, 11, 14, 15, 16},
			SessionFrequency: 1.5,
		},
	}

	// Extract behavior data from context if provided
	if typing, exists := contextData["typing_pattern"]; exists {
		if typingMap, ok := typing.(map[string]interface{}); ok {
			if avgKeydown, exists := typingMap["avg_keydown_time"]; exists {
				if avgKeydownFloat, ok := avgKeydown.(float64); ok {
					signals.TypingPattern.AvgKeydownTime = avgKeydownFloat
				}
			}
		}
	}

	return signals
}

func calculateRiskScore(userID, ipAddress, userAgent, deviceFingerprint string, location LocationInfo, behavior BehaviorSignals) RiskAssessment {
	assessment := RiskAssessment{
		UserID:            userID,
		IPAddress:         ipAddress,
		UserAgent:         userAgent,
		DeviceFingerprint: deviceFingerprint,
		Location:          location,
		BehaviorSignals:   behavior,
		Timestamp:         time.Now(),
		Factors:           []RiskFactor{},
	}

	totalScore := 0.0

	// Location-based risk factors
	if location.IsVPN {
		factor := RiskFactor{
			Type:        "location",
			Description: "VPN usage detected",
			Weight:      0.3,
			Score:       0.6,
			Severity:    "medium",
		}
		assessment.Factors = append(assessment.Factors, factor)
		totalScore += factor.Weight * factor.Score
	}

	if location.IsTor {
		factor := RiskFactor{
			Type:        "location",
			Description: "Tor network usage detected",
			Weight:      0.5,
			Score:       0.9,
			Severity:    "high",
		}
		assessment.Factors = append(assessment.Factors, factor)
		totalScore += factor.Weight * factor.Score
	}

	// Time-based risk factors
	currentHour := time.Now().Hour()
	isOffHours := currentHour < 6 || currentHour > 22
	if isOffHours {
		factor := RiskFactor{
			Type:        "temporal",
			Description: "Login outside typical hours",
			Weight:      0.2,
			Score:       0.4,
			Severity:    "low",
		}
		assessment.Factors = append(assessment.Factors, factor)
		totalScore += factor.Weight * factor.Score
	}

	// Device fingerprint risk
	if deviceFingerprint == "" {
		factor := RiskFactor{
			Type:        "device",
			Description: "No device fingerprint available",
			Weight:      0.2,
			Score:       0.3,
			Severity:    "low",
		}
		assessment.Factors = append(assessment.Factors, factor)
		totalScore += factor.Weight * factor.Score
	}

	// Check for new device
	isNewDevice, err := services.IsNewDevice(userID, deviceFingerprint)
	if err == nil && isNewDevice {
		factor := RiskFactor{
			Type:        "device",
			Description: "New device detected",
			Weight:      0.4,
			Score:       0.7,
			Severity:    "medium",
		}
		assessment.Factors = append(assessment.Factors, factor)
		totalScore += factor.Weight * factor.Score
	}

	// Behavior analysis
	if behavior.TypingPattern.AvgKeydownTime > 200 || behavior.TypingPattern.AvgKeydownTime < 50 {
		factor := RiskFactor{
			Type:        "behavior",
			Description: "Unusual typing pattern detected",
			Weight:      0.15,
			Score:       0.5,
			Severity:    "low",
		}
		assessment.Factors = append(assessment.Factors, factor)
		totalScore += factor.Weight * factor.Score
	}

	// Normalize score to 0-1 range
	assessment.RiskScore = math.Min(totalScore, 1.0)

	// Determine risk level
	if assessment.RiskScore < 0.3 {
		assessment.RiskLevel = "low"
	} else if assessment.RiskScore < 0.6 {
		assessment.RiskLevel = "medium"
	} else if assessment.RiskScore < 0.8 {
		assessment.RiskLevel = "high"
	} else {
		assessment.RiskLevel = "critical"
	}

	// Generate recommendations
	assessment.Recommendations = generateRecommendations(assessment)

	return assessment
}

func makePolicyDecision(assessment RiskAssessment) PolicyDecision {
	decision := PolicyDecision{
		Confidence: 0.8,
		Monitoring: MonitoringLevel{
			Level:            "basic",
			LogAllActions:    false,
			RealTimeAlerts:   false,
			BehaviorTracking: false,
		},
	}

	switch assessment.RiskLevel {
	case "low":
		decision.Action = "allow"
		decision.Explanation = "Low risk - standard access granted"
		decision.SessionLimits = SessionLimits{
			MaxDuration: 480, // 8 hours
			IdleTimeout: 60,  // 1 hour
		}

	case "medium":
		decision.Action = "step_up"
		decision.RequiredMFA = []string{"totp", "sms"}
		decision.Explanation = "Medium risk - additional authentication required"
		decision.SessionLimits = SessionLimits{
			MaxDuration:   240, // 4 hours
			IdleTimeout:   30,  // 30 minutes
			RequireReauth: true,
		}
		decision.Monitoring.Level = "enhanced"
		decision.Monitoring.LogAllActions = true

	case "high":
		decision.Action = "step_up"
		decision.RequiredMFA = []string{"webauthn", "totp"}
		decision.Explanation = "High risk - strong authentication and monitoring required"
		decision.SessionLimits = SessionLimits{
			MaxDuration:   120, // 2 hours
			IdleTimeout:   15,  // 15 minutes
			RequireReauth: true,
			LimitedAccess: true,
		}
		decision.Monitoring.Level = "full"
		decision.Monitoring.LogAllActions = true
		decision.Monitoring.RealTimeAlerts = true
		decision.Monitoring.BehaviorTracking = true

	case "critical":
		decision.Action = "deny"
		decision.Explanation = "Critical risk - access denied"
		decision.Monitoring.Level = "full"
		decision.Monitoring.LogAllActions = true
		decision.Monitoring.RealTimeAlerts = true
	}

	return decision
}

func generateRecommendations(assessment RiskAssessment) []string {
	recommendations := []string{}

	if assessment.RiskScore > 0.5 {
		recommendations = append(recommendations, "Enable additional MFA methods")
		recommendations = append(recommendations, "Monitor session activity closely")
	}

	if assessment.Location.IsVPN || assessment.Location.IsTor {
		recommendations = append(recommendations, "Verify user identity through alternative means")
		recommendations = append(recommendations, "Consider blocking anonymous network access")
	}

	for _, factor := range assessment.Factors {
		if factor.Type == "device" && strings.Contains(factor.Description, "New device") {
			recommendations = append(recommendations, "Send device registration notification to user")
			recommendations = append(recommendations, "Require device verification")
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Continue with standard security monitoring")
	}

	return recommendations
}
