package services

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"cloudgate-backend/internal/models"
)

// AdaptiveAuthService handles intelligent authentication decisions
type AdaptiveAuthService struct {
	db                  *gorm.DB
	mfaService          *MFAService
	oauthMonitorService *OAuthMonitoringService
	userService         *UserService
}

// AuthContext contains all context information for authentication decision
type AuthContext struct {
	UserID            uuid.UUID              `json:"user_id"`
	Email             string                 `json:"email"`
	IPAddress         string                 `json:"ip_address"`
	UserAgent         string                 `json:"user_agent"`
	DeviceFingerprint string                 `json:"device_fingerprint"`
	Location          *GeoLocation           `json:"location,omitempty"`
	LoginTime         time.Time              `json:"login_time"`
	SessionInfo       map[string]interface{} `json:"session_info"`
	RequestHeaders    map[string]string      `json:"request_headers"`
	ApplicationID     string                 `json:"application_id,omitempty"`
}

// GeoLocation represents geographical location data
type GeoLocation struct {
	Country     string  `json:"country"`
	Region      string  `json:"region"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ISP         string  `json:"isp"`
	Timezone    string  `json:"timezone"`
	VPNDetected bool    `json:"vpn_detected"`
}

// AuthDecision represents the authentication decision
type AuthDecision struct {
	Decision        AuthDecisionType       `json:"decision"`
	RiskScore       float64                `json:"risk_score"`
	RiskLevel       string                 `json:"risk_level"`
	RequiredActions []AuthAction           `json:"required_actions"`
	Reasoning       []string               `json:"reasoning"`
	SessionDuration time.Duration          `json:"session_duration"`
	Restrictions    []AuthRestriction      `json:"restrictions"`
	Metadata        map[string]interface{} `json:"metadata"`
	ExpiresAt       time.Time              `json:"expires_at"`
}

// AuthDecisionType represents the type of authentication decision
type AuthDecisionType string

const (
	AuthDecisionAllow     AuthDecisionType = "allow"
	AuthDecisionDeny      AuthDecisionType = "deny"
	AuthDecisionChallenge AuthDecisionType = "challenge"
	AuthDecisionMonitor   AuthDecisionType = "monitor"
)

// AuthAction represents required authentication actions
type AuthAction struct {
	Type        AuthActionType         `json:"type"`
	Required    bool                   `json:"required"`
	Timeout     time.Duration          `json:"timeout"`
	Metadata    map[string]interface{} `json:"metadata"`
	Description string                 `json:"description"`
}

// AuthActionType represents the type of authentication action
type AuthActionType string

const (
	ActionMFARequired           AuthActionType = "mfa_required"
	ActionPasswordChange        AuthActionType = "password_change"
	ActionDeviceVerification    AuthActionType = "device_verification"
	ActionEmailVerification     AuthActionType = "email_verification"
	ActionCaptchaVerification   AuthActionType = "captcha_verification"
	ActionAdminApproval         AuthActionType = "admin_approval"
	ActionSecurityQuestions     AuthActionType = "security_questions"
	ActionBiometricVerification AuthActionType = "biometric_verification"
)

// AuthRestriction represents access restrictions
type AuthRestriction struct {
	Type        RestrictionType `json:"type"`
	Value       interface{}     `json:"value"`
	Description string          `json:"description"`
	ExpiresAt   *time.Time      `json:"expires_at,omitempty"`
}

// RestrictionType represents the type of restriction
type RestrictionType string

const (
	RestrictionIPWhitelist     RestrictionType = "ip_whitelist"
	RestrictionTimeWindow      RestrictionType = "time_window"
	RestrictionApplications    RestrictionType = "applications"
	RestrictionFeatures        RestrictionType = "features"
	RestrictionDataAccess      RestrictionType = "data_access"
	RestrictionSessionDuration RestrictionType = "session_duration"
)

// RiskFactors contains individual risk assessment factors
type RiskFactors struct {
	LocationRisk    float64 `json:"location_risk"`
	DeviceRisk      float64 `json:"device_risk"`
	BehavioralRisk  float64 `json:"behavioral_risk"`
	TemporalRisk    float64 `json:"temporal_risk"`
	NetworkRisk     float64 `json:"network_risk"`
	ApplicationRisk float64 `json:"application_risk"`
	HistoricalRisk  float64 `json:"historical_risk"`
	VelocityRisk    float64 `json:"velocity_risk"`
}

// NewAdaptiveAuthService creates a new adaptive authentication service
func NewAdaptiveAuthService(db *gorm.DB) *AdaptiveAuthService {
	return &AdaptiveAuthService{
		db:                  db,
		mfaService:          NewMFAService(db),
		oauthMonitorService: NewOAuthMonitoringService(db),
		userService:         NewUserService(db),
	}
}

// EvaluateAuthentication performs comprehensive authentication evaluation
func (s *AdaptiveAuthService) EvaluateAuthentication(ctx *AuthContext) (*AuthDecision, error) {
	// 1. Perform comprehensive risk assessment
	riskFactors, err := s.assessRiskFactors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to assess risk factors: %w", err)
	}

	// 2. Calculate overall risk score
	overallRisk := s.calculateOverallRisk(riskFactors)

	// 3. Determine risk level
	riskLevel := s.determineRiskLevel(overallRisk)

	// 4. Make authentication decision based on risk
	decision := s.makeAuthDecision(ctx, overallRisk, riskLevel, riskFactors)

	// 5. Store the assessment for learning
	err = s.storeAuthAssessment(ctx, decision, riskFactors)
	if err != nil {
		// Log error but don't fail the authentication
		fmt.Printf("Failed to store auth assessment: %v\n", err)
	}

	// 6. Update user behavior patterns
	go s.updateUserBehaviorPatterns(ctx, decision)

	return decision, nil
}

// assessRiskFactors evaluates all risk factors
func (s *AdaptiveAuthService) assessRiskFactors(ctx *AuthContext) (*RiskFactors, error) {
	factors := &RiskFactors{}

	// Assess location risk
	factors.LocationRisk = s.assessLocationRisk(ctx)

	// Assess device risk
	factors.DeviceRisk = s.assessDeviceRisk(ctx)

	// Assess behavioral risk
	factors.BehavioralRisk = s.assessBehavioralRisk(ctx)

	// Assess temporal risk
	factors.TemporalRisk = s.assessTemporalRisk(ctx)

	// Assess network risk
	factors.NetworkRisk = s.assessNetworkRisk(ctx)

	// Assess application risk
	factors.ApplicationRisk = s.assessApplicationRisk(ctx)

	// Assess historical risk
	factors.HistoricalRisk = s.assessHistoricalRisk(ctx)

	// Assess velocity risk
	factors.VelocityRisk = s.assessVelocityRisk(ctx)

	return factors, nil
}

// assessLocationRisk evaluates location-based risk
func (s *AdaptiveAuthService) assessLocationRisk(ctx *AuthContext) float64 {
	risk := 0.0

	// Check if location is provided
	if ctx.Location == nil {
		return 0.3 // Moderate risk for unknown location
	}

	// VPN/Proxy detection
	if ctx.Location.VPNDetected {
		risk += 0.4
	}

	// Get user's historical locations
	historicalLocations := s.getUserHistoricalLocations(ctx.UserID)

	// Check if this is a new country
	isNewCountry := true
	for _, loc := range historicalLocations {
		if loc.Country == ctx.Location.Country {
			isNewCountry = false
			break
		}
	}

	if isNewCountry {
		risk += 0.3
	}

	// Calculate distance from usual locations
	if len(historicalLocations) > 0 {
		minDistance := math.MaxFloat64
		for _, loc := range historicalLocations {
			distance := s.calculateDistance(
				ctx.Location.Latitude, ctx.Location.Longitude,
				loc.Latitude, loc.Longitude,
			)
			if distance < minDistance {
				minDistance = distance
			}
		}

		// Add risk based on distance (normalized)
		if minDistance > 1000 { // More than 1000km
			risk += math.Min(0.3, minDistance/10000)
		}
	}

	// Check for high-risk countries (simplified check)
	highRiskCountries := []string{"CN", "RU", "KP", "IR"}
	for _, country := range highRiskCountries {
		if ctx.Location.Country == country {
			risk += 0.2
			break
		}
	}

	return math.Min(risk, 1.0)
}

// assessDeviceRisk evaluates device-based risk
func (s *AdaptiveAuthService) assessDeviceRisk(ctx *AuthContext) float64 {
	risk := 0.0

	// Check if device is known
	isKnownDevice, err := IsNewDevice(ctx.UserID.String(), ctx.DeviceFingerprint)
	if err != nil {
		risk += 0.2 // Add risk for unknown device status
	} else if isKnownDevice {
		risk += 0.4 // New device adds significant risk
	}

	// Analyze user agent for suspicious patterns
	if s.isSuspiciousUserAgent(ctx.UserAgent) {
		risk += 0.3
	}

	// Check device consistency
	if s.isInconsistentDevice(ctx) {
		risk += 0.2
	}

	return math.Min(risk, 1.0)
}

// assessBehavioralRisk evaluates behavioral patterns
func (s *AdaptiveAuthService) assessBehavioralRisk(ctx *AuthContext) float64 {
	risk := 0.0

	// Get user's typical login patterns
	patterns := s.getUserBehaviorPatterns(ctx.UserID)

	// Check login time patterns
	if !s.isTypicalLoginTime(ctx.LoginTime, patterns) {
		risk += 0.2
	}

	// Check session patterns
	if !s.isTypicalSessionPattern(ctx, patterns) {
		risk += 0.2
	}

	// Check application access patterns
	if ctx.ApplicationID != "" && !s.isTypicalApplicationAccess(ctx.UserID, ctx.ApplicationID) {
		risk += 0.1
	}

	return math.Min(risk, 1.0)
}

// assessTemporalRisk evaluates time-based risk
func (s *AdaptiveAuthService) assessTemporalRisk(ctx *AuthContext) float64 {
	risk := 0.0

	// Check for unusual hours
	hour := ctx.LoginTime.Hour()
	if hour < 6 || hour > 22 {
		risk += 0.2
	}

	// Check for weekend access (if unusual for user)
	if s.isWeekendAccessUnusual(ctx.UserID, ctx.LoginTime) {
		risk += 0.1
	}

	// Check for rapid successive logins
	if s.hasRecentLogins(ctx.UserID, ctx.LoginTime) {
		risk += 0.3
	}

	return math.Min(risk, 1.0)
}

// assessNetworkRisk evaluates network-based risk
func (s *AdaptiveAuthService) assessNetworkRisk(ctx *AuthContext) float64 {
	risk := 0.0

	// Parse IP address
	ip := net.ParseIP(ctx.IPAddress)
	if ip == nil {
		return 0.5 // Invalid IP is high risk
	}

	// Check for private/local IPs in production
	if ip.IsPrivate() || ip.IsLoopback() {
		risk += 0.1
	}

	// Check IP reputation (simplified)
	if s.isHighRiskIP(ctx.IPAddress) {
		risk += 0.5
	}

	// Check for TOR exit nodes (simplified)
	if s.isTorExitNode(ctx.IPAddress) {
		risk += 0.6
	}

	return math.Min(risk, 1.0)
}

// assessApplicationRisk evaluates application-specific risk
func (s *AdaptiveAuthService) assessApplicationRisk(ctx *AuthContext) float64 {
	risk := 0.0

	if ctx.ApplicationID == "" {
		return 0.0
	}

	// Check application sensitivity level
	sensitivityLevel := s.getApplicationSensitivityLevel(ctx.ApplicationID)
	risk += sensitivityLevel * 0.3

	// Check for unusual application access
	if !s.hasUserAccessedApplication(ctx.UserID, ctx.ApplicationID) {
		risk += 0.2
	}

	return math.Min(risk, 1.0)
}

// assessHistoricalRisk evaluates historical security events
func (s *AdaptiveAuthService) assessHistoricalRisk(ctx *AuthContext) float64 {
	risk := 0.0

	// Check for recent security events
	recentEvents := s.getRecentSecurityEvents(ctx.UserID)
	if len(recentEvents) > 0 {
		risk += math.Min(float64(len(recentEvents))*0.1, 0.4)
	}

	// Check for failed login attempts
	failedAttempts := s.getRecentFailedAttempts(ctx.UserID, ctx.IPAddress)
	if failedAttempts > 0 {
		risk += math.Min(float64(failedAttempts)*0.1, 0.3)
	}

	// Check for account compromise indicators
	if s.hasCompromiseIndicators(ctx.UserID) {
		risk += 0.5
	}

	return math.Min(risk, 1.0)
}

// assessVelocityRisk evaluates login velocity
func (s *AdaptiveAuthService) assessVelocityRisk(ctx *AuthContext) float64 {
	risk := 0.0

	// Check login frequency in last hour
	recentLogins := s.getRecentLoginCount(ctx.UserID, time.Hour)
	if recentLogins > 10 {
		risk += 0.4
	} else if recentLogins > 5 {
		risk += 0.2
	}

	// Check for impossible travel
	if s.hasImpossibleTravel(ctx) {
		risk += 0.8
	}

	return math.Min(risk, 1.0)
}

// calculateOverallRisk combines all risk factors
func (s *AdaptiveAuthService) calculateOverallRisk(factors *RiskFactors) float64 {
	// Weighted combination of risk factors
	weights := map[string]float64{
		"location":    0.20,
		"device":      0.15,
		"behavioral":  0.15,
		"temporal":    0.10,
		"network":     0.15,
		"application": 0.10,
		"historical":  0.10,
		"velocity":    0.05,
	}

	totalRisk := factors.LocationRisk*weights["location"] +
		factors.DeviceRisk*weights["device"] +
		factors.BehavioralRisk*weights["behavioral"] +
		factors.TemporalRisk*weights["temporal"] +
		factors.NetworkRisk*weights["network"] +
		factors.ApplicationRisk*weights["application"] +
		factors.HistoricalRisk*weights["historical"] +
		factors.VelocityRisk*weights["velocity"]

	return math.Min(totalRisk, 1.0)
}

// determineRiskLevel categorizes risk score
func (s *AdaptiveAuthService) determineRiskLevel(riskScore float64) string {
	switch {
	case riskScore < 0.2:
		return "low"
	case riskScore < 0.4:
		return "medium"
	case riskScore < 0.7:
		return "high"
	default:
		return "critical"
	}
}

// makeAuthDecision creates authentication decision based on risk
func (s *AdaptiveAuthService) makeAuthDecision(ctx *AuthContext, riskScore float64, riskLevel string, factors *RiskFactors) *AuthDecision {
	decision := &AuthDecision{
		RiskScore:       riskScore,
		RiskLevel:       riskLevel,
		RequiredActions: []AuthAction{},
		Reasoning:       []string{},
		Restrictions:    []AuthRestriction{},
		Metadata:        make(map[string]interface{}),
		ExpiresAt:       time.Now().Add(24 * time.Hour),
	}

	// Decision logic based on risk level
	switch riskLevel {
	case "low":
		decision.Decision = AuthDecisionAllow
		decision.SessionDuration = 8 * time.Hour
		decision.Reasoning = append(decision.Reasoning, "Low risk authentication - standard access granted")

	case "medium":
		decision.Decision = AuthDecisionChallenge
		decision.SessionDuration = 4 * time.Hour
		decision.RequiredActions = append(decision.RequiredActions, AuthAction{
			Type:        ActionMFARequired,
			Required:    true,
			Timeout:     5 * time.Minute,
			Description: "Multi-factor authentication required due to elevated risk",
		})
		decision.Reasoning = append(decision.Reasoning, "Medium risk detected - MFA required")

	case "high":
		decision.Decision = AuthDecisionChallenge
		decision.SessionDuration = 2 * time.Hour

		// Multiple authentication factors required
		decision.RequiredActions = append(decision.RequiredActions,
			AuthAction{
				Type:        ActionMFARequired,
				Required:    true,
				Timeout:     5 * time.Minute,
				Description: "Multi-factor authentication required",
			},
			AuthAction{
				Type:        ActionEmailVerification,
				Required:    true,
				Timeout:     10 * time.Minute,
				Description: "Email verification required for high-risk login",
			},
		)

		// Add restrictions
		decision.Restrictions = append(decision.Restrictions, AuthRestriction{
			Type:        RestrictionSessionDuration,
			Value:       2 * time.Hour,
			Description: "Limited session duration due to high risk",
		})

		decision.Reasoning = append(decision.Reasoning, "High risk detected - enhanced verification required")

	case "critical":
		decision.Decision = AuthDecisionDeny
		decision.SessionDuration = 0
		decision.Reasoning = append(decision.Reasoning, "Critical risk detected - access denied")

		// Log security event
		go s.logSecurityEvent(ctx, "critical_risk_access_denied", riskScore)
	}

	// Add specific reasoning based on risk factors
	s.addSpecificReasoning(decision, factors)

	// Store decision metadata
	decision.Metadata["risk_factors"] = factors
	decision.Metadata["assessment_time"] = time.Now()
	decision.Metadata["user_id"] = ctx.UserID.String()

	return decision
}

// Helper methods for risk assessment (simplified implementations)

func (s *AdaptiveAuthService) getUserHistoricalLocations(userID uuid.UUID) []GeoLocation {
	// Implementation would query historical location data
	return []GeoLocation{}
}

func (s *AdaptiveAuthService) calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Haversine formula for distance calculation
	const R = 6371 // Earth's radius in kilometers

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func (s *AdaptiveAuthService) isSuspiciousUserAgent(userAgent string) bool {
	suspiciousPatterns := []string{
		"bot", "crawler", "spider", "scraper",
		"curl", "wget", "python", "automation",
	}

	userAgentLower := strings.ToLower(userAgent)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(userAgentLower, pattern) {
			return true
		}
	}
	return false
}

func (s *AdaptiveAuthService) isInconsistentDevice(ctx *AuthContext) bool {
	// Check for device fingerprint inconsistencies
	// This would involve more sophisticated device analysis
	return false
}

func (s *AdaptiveAuthService) getUserBehaviorPatterns(userID uuid.UUID) map[string]interface{} {
	// Implementation would return user behavior patterns
	return make(map[string]interface{})
}

func (s *AdaptiveAuthService) isTypicalLoginTime(loginTime time.Time, patterns map[string]interface{}) bool {
	// Implementation would check against user's typical login hours
	return true
}

func (s *AdaptiveAuthService) isTypicalSessionPattern(ctx *AuthContext, patterns map[string]interface{}) bool {
	// Implementation would analyze session patterns
	return true
}

func (s *AdaptiveAuthService) isTypicalApplicationAccess(userID uuid.UUID, appID string) bool {
	// Implementation would check user's application access history
	return true
}

func (s *AdaptiveAuthService) isWeekendAccessUnusual(userID uuid.UUID, loginTime time.Time) bool {
	// Implementation would check if weekend access is unusual for this user
	return false
}

func (s *AdaptiveAuthService) hasRecentLogins(userID uuid.UUID, loginTime time.Time) bool {
	// Implementation would check for recent login attempts
	return false
}

func (s *AdaptiveAuthService) isHighRiskIP(ipAddress string) bool {
	// Implementation would check IP against threat intelligence feeds
	return false
}

func (s *AdaptiveAuthService) isTorExitNode(ipAddress string) bool {
	// Implementation would check against TOR exit node lists
	return false
}

func (s *AdaptiveAuthService) getApplicationSensitivityLevel(appID string) float64 {
	// Implementation would return application sensitivity level (0.0-1.0)
	return 0.0
}

func (s *AdaptiveAuthService) hasUserAccessedApplication(userID uuid.UUID, appID string) bool {
	// Implementation would check user's application access history
	return true
}

func (s *AdaptiveAuthService) getRecentSecurityEvents(userID uuid.UUID) []models.SecurityEvent {
	// Implementation would query recent security events
	return []models.SecurityEvent{}
}

func (s *AdaptiveAuthService) getRecentFailedAttempts(userID uuid.UUID, ipAddress string) int {
	// Implementation would count recent failed login attempts
	return 0
}

func (s *AdaptiveAuthService) hasCompromiseIndicators(userID uuid.UUID) bool {
	// Implementation would check for account compromise indicators
	return false
}

func (s *AdaptiveAuthService) getRecentLoginCount(userID uuid.UUID, duration time.Duration) int {
	// Implementation would count recent logins within duration
	return 0
}

func (s *AdaptiveAuthService) hasImpossibleTravel(ctx *AuthContext) bool {
	// Implementation would check for impossible travel scenarios
	return false
}

func (s *AdaptiveAuthService) addSpecificReasoning(decision *AuthDecision, factors *RiskFactors) {
	if factors.LocationRisk > 0.3 {
		decision.Reasoning = append(decision.Reasoning, "Unusual or high-risk location detected")
	}
	if factors.DeviceRisk > 0.3 {
		decision.Reasoning = append(decision.Reasoning, "Unrecognized or suspicious device")
	}
	if factors.VelocityRisk > 0.3 {
		decision.Reasoning = append(decision.Reasoning, "High login velocity detected")
	}
}

func (s *AdaptiveAuthService) storeAuthAssessment(ctx *AuthContext, decision *AuthDecision, factors *RiskFactors) error {
	// Store assessment for machine learning and analysis
	assessment := map[string]interface{}{
		"user_id":    ctx.UserID.String(),
		"risk_score": decision.RiskScore,
		"risk_level": decision.RiskLevel,
		"decision":   decision.Decision,
		"factors":    factors,
		"context":    ctx,
		"timestamp":  time.Now(),
	}

	// Convert to JSON for storage
	assessmentJSON, err := json.Marshal(assessment)
	if err != nil {
		return err
	}

	// Store in risk assessment table
	return StoreRiskAssessment(map[string]interface{}{
		"user_id":       ctx.UserID.String(),
		"risk_score":    decision.RiskScore,
		"risk_level":    decision.RiskLevel,
		"factors":       string(assessmentJSON),
		"ip_address":    ctx.IPAddress,
		"user_agent":    ctx.UserAgent,
		"location_data": ctx.Location,
	})
}

func (s *AdaptiveAuthService) updateUserBehaviorPatterns(ctx *AuthContext, decision *AuthDecision) {
	// Update user behavior patterns for future assessments
	// This would involve machine learning model updates
}

func (s *AdaptiveAuthService) logSecurityEvent(ctx *AuthContext, eventType string, riskScore float64) {
	// Log security event for monitoring and alerting
	s.oauthMonitorService.CreateSecurityEvent(
		ctx.UserID.String(),
		eventType,
		fmt.Sprintf("Risk score: %.2f", riskScore),
		"critical",
		ctx.IPAddress,
		ctx.UserAgent,
		"",
		riskScore,
		nil,
	)
}
