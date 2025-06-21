package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditService handles comprehensive audit logging for compliance and security
type AuditService struct {
	db *gorm.DB
}

// AuditEvent represents a comprehensive audit log entry
type AuditEvent struct {
	ID              uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Timestamp       time.Time              `json:"timestamp" gorm:"not null;index"`
	EventType       AuditEventType         `json:"event_type" gorm:"not null;index"`
	Category        AuditCategory          `json:"category" gorm:"not null;index"`
	Severity        AuditSeverity          `json:"severity" gorm:"not null;index"`
	UserID          *uuid.UUID             `json:"user_id,omitempty" gorm:"type:uuid;index"`
	SessionID       *uuid.UUID             `json:"session_id,omitempty" gorm:"type:uuid;index"`
	IPAddress       string                 `json:"ip_address" gorm:"index"`
	UserAgent       string                 `json:"user_agent"`
	Resource        string                 `json:"resource" gorm:"not null;index"`
	Action          string                 `json:"action" gorm:"not null;index"`
	Outcome         AuditOutcome           `json:"outcome" gorm:"not null;index"`
	Description     string                 `json:"description" gorm:"not null"`
	Details         map[string]interface{} `json:"details" gorm:"type:jsonb"`
	RiskScore       *float64               `json:"risk_score,omitempty"`
	ComplianceFlags []string               `json:"compliance_flags" gorm:"type:text[]"`
	Tags            []string               `json:"tags" gorm:"type:text[]"`
	CorrelationID   *uuid.UUID             `json:"correlation_id,omitempty" gorm:"type:uuid;index"`
	ParentEventID   *uuid.UUID             `json:"parent_event_id,omitempty" gorm:"type:uuid;index"`
	CreatedAt       time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// AuditEventType represents the type of audit event
type AuditEventType string

const (
	// Authentication events
	EventTypeLogin           AuditEventType = "login"
	EventTypeLogout          AuditEventType = "logout"
	EventTypeLoginFailed     AuditEventType = "login_failed"
	EventTypePasswordChange  AuditEventType = "password_change"
	EventTypeAccountLocked   AuditEventType = "account_locked"
	EventTypeAccountUnlocked AuditEventType = "account_unlocked"
	EventTypeMFAEnabled      AuditEventType = "mfa_enabled"
	EventTypeMFADisabled     AuditEventType = "mfa_disabled"
	EventTypeMFAVerified     AuditEventType = "mfa_verified"
	EventTypeMFAFailed       AuditEventType = "mfa_failed"

	// Authorization events
	EventTypePermissionGranted  AuditEventType = "permission_granted"
	EventTypePermissionDenied   AuditEventType = "permission_denied"
	EventTypeRoleAssigned       AuditEventType = "role_assigned"
	EventTypeRoleRevoked        AuditEventType = "role_revoked"
	EventTypePrivilegeEscalated AuditEventType = "privilege_escalated"

	// Data access events
	EventTypeDataAccess       AuditEventType = "data_access"
	EventTypeDataModification AuditEventType = "data_modification"
	EventTypeDataDeletion     AuditEventType = "data_deletion"
	EventTypeDataExport       AuditEventType = "data_export"
	EventTypeDataImport       AuditEventType = "data_import"
	EventTypeBulkOperation    AuditEventType = "bulk_operation"

	// System events
	EventTypeSystemStartup       AuditEventType = "system_startup"
	EventTypeSystemShutdown      AuditEventType = "system_shutdown"
	EventTypeConfigurationChange AuditEventType = "configuration_change"
	EventTypeServiceFailure      AuditEventType = "service_failure"
	EventTypeBackupCreated       AuditEventType = "backup_created"
	EventTypeBackupRestored      AuditEventType = "backup_restored"

	// Security events
	EventTypeSecurityAlert           AuditEventType = "security_alert"
	EventTypeSecurityIncident        AuditEventType = "security_incident"
	EventTypeIntrusionDetected       AuditEventType = "intrusion_detected"
	EventTypeSuspiciousActivity      AuditEventType = "suspicious_activity"
	EventTypeSecurityPolicyViolation AuditEventType = "security_policy_violation"

	// OAuth and SSO events
	EventTypeOAuthAuthorization AuditEventType = "oauth_authorization"
	EventTypeOAuthTokenIssued   AuditEventType = "oauth_token_issued"
	EventTypeOAuthTokenRevoked  AuditEventType = "oauth_token_revoked"
	EventTypeSSOInitiated       AuditEventType = "sso_initiated"
	EventTypeSSOCompleted       AuditEventType = "sso_completed"
	EventTypeSSOFailed          AuditEventType = "sso_failed"

	// Administrative events
	EventTypeUserCreated     AuditEventType = "user_created"
	EventTypeUserModified    AuditEventType = "user_modified"
	EventTypeUserDeleted     AuditEventType = "user_deleted"
	EventTypeUserDeactivated AuditEventType = "user_deactivated"
	EventTypeUserReactivated AuditEventType = "user_reactivated"
	EventTypeAdminAction     AuditEventType = "admin_action"

	// API events
	EventTypeAPICall           AuditEventType = "api_call"
	EventTypeAPIError          AuditEventType = "api_error"
	EventTypeRateLimitExceeded AuditEventType = "rate_limit_exceeded"
	EventTypeAPIKeyCreated     AuditEventType = "api_key_created"
	EventTypeAPIKeyRevoked     AuditEventType = "api_key_revoked"
)

// AuditCategory represents the category of audit event
type AuditCategory string

const (
	CategoryAuthentication AuditCategory = "authentication"
	CategoryAuthorization  AuditCategory = "authorization"
	CategoryDataAccess     AuditCategory = "data_access"
	CategorySystem         AuditCategory = "system"
	CategorySecurity       AuditCategory = "security"
	CategoryCompliance     AuditCategory = "compliance"
	CategoryAdministrative AuditCategory = "administrative"
	CategoryAPI            AuditCategory = "api"
	CategoryOAuth          AuditCategory = "oauth"
	CategorySSO            AuditCategory = "sso"
)

// AuditSeverity represents the severity level of audit event
type AuditSeverity string

const (
	AuditSeverityInfo     AuditSeverity = "info"
	AuditSeverityWarning  AuditSeverity = "warning"
	AuditSeverityError    AuditSeverity = "error"
	AuditSeverityCritical AuditSeverity = "critical"
)

// AuditOutcome represents the outcome of the audited action
type AuditOutcome string

const (
	OutcomeSuccess AuditOutcome = "success"
	OutcomeFailure AuditOutcome = "failure"
	OutcomeDenied  AuditOutcome = "denied"
	OutcomeError   AuditOutcome = "error"
)

// AuditFilter represents filtering options for audit queries
type AuditFilter struct {
	StartTime     *time.Time
	EndTime       *time.Time
	EventTypes    []AuditEventType
	Categories    []AuditCategory
	Severities    []AuditSeverity
	UserID        *uuid.UUID
	IPAddress     string
	Resource      string
	Action        string
	Outcome       *AuditOutcome
	RiskScoreMin  *float64
	RiskScoreMax  *float64
	Tags          []string
	CorrelationID *uuid.UUID
	Limit         int
	Offset        int
}

// AuditStatistics represents audit statistics for reporting
type AuditStatistics struct {
	TotalEvents          int64                    `json:"total_events"`
	EventsByType         map[AuditEventType]int64 `json:"events_by_type"`
	EventsByCategory     map[AuditCategory]int64  `json:"events_by_category"`
	EventsBySeverity     map[AuditSeverity]int64  `json:"events_by_severity"`
	EventsByOutcome      map[AuditOutcome]int64   `json:"events_by_outcome"`
	TopUsers             []UserAuditSummary       `json:"top_users"`
	TopResources         []ResourceAuditSummary   `json:"top_resources"`
	SecurityEvents       int64                    `json:"security_events"`
	FailedAttempts       int64                    `json:"failed_attempts"`
	AverageRiskScore     float64                  `json:"average_risk_score"`
	ComplianceViolations int64                    `json:"compliance_violations"`
	TimeRange            AuditTimeRange           `json:"time_range"`
}

// UserAuditSummary represents audit summary for a user
type UserAuditSummary struct {
	UserID       uuid.UUID `json:"user_id"`
	EventCount   int64     `json:"event_count"`
	LastActivity time.Time `json:"last_activity"`
	RiskScore    float64   `json:"risk_score"`
}

// ResourceAuditSummary represents audit summary for a resource
type ResourceAuditSummary struct {
	Resource    string    `json:"resource"`
	EventCount  int64     `json:"event_count"`
	LastAccess  time.Time `json:"last_access"`
	UniqueUsers int64     `json:"unique_users"`
}

// AuditTimeRange represents the time range for audit statistics
type AuditTimeRange struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// ComplianceReport represents a compliance audit report
type ComplianceReport struct {
	ID              uuid.UUID                  `json:"id"`
	ReportType      ComplianceReportType       `json:"report_type"`
	TimeRange       AuditTimeRange             `json:"time_range"`
	GeneratedAt     time.Time                  `json:"generated_at"`
	GeneratedBy     uuid.UUID                  `json:"generated_by"`
	Statistics      AuditStatistics            `json:"statistics"`
	ComplianceFlags map[string]int64           `json:"compliance_flags"`
	Violations      []ComplianceViolation      `json:"violations"`
	Recommendations []ComplianceRecommendation `json:"recommendations"`
	Status          ComplianceReportStatus     `json:"status"`
}

// ComplianceReportType represents the type of compliance report
type ComplianceReportType string

const (
	ReportTypeSOX      ComplianceReportType = "sox"
	ReportTypeGDPR     ComplianceReportType = "gdpr"
	ReportTypeHIPAA    ComplianceReportType = "hipaa"
	ReportTypeSOC2     ComplianceReportType = "soc2"
	ReportTypePCI      ComplianceReportType = "pci"
	ReportTypeISO27001 ComplianceReportType = "iso27001"
	ReportTypeCustom   ComplianceReportType = "custom"
)

// ComplianceViolation represents a compliance violation
type ComplianceViolation struct {
	ID          uuid.UUID     `json:"id"`
	Type        string        `json:"type"`
	Description string        `json:"description"`
	Severity    AuditSeverity `json:"severity"`
	EventID     uuid.UUID     `json:"event_id"`
	Timestamp   time.Time     `json:"timestamp"`
	Status      string        `json:"status"`
}

// ComplianceRecommendation represents a compliance recommendation
type ComplianceRecommendation struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Category    string    `json:"category"`
}

// ComplianceReportStatus represents the status of a compliance report
type ComplianceReportStatus string

const (
	ReportStatusGenerating ComplianceReportStatus = "generating"
	ReportStatusCompleted  ComplianceReportStatus = "completed"
	ReportStatusFailed     ComplianceReportStatus = "failed"
)

// NewAuditService creates a new audit service
func NewAuditService(db *gorm.DB) *AuditService {
	service := &AuditService{
		db: db,
	}

	// Auto-migrate the audit event table
	if err := db.AutoMigrate(&AuditEvent{}); err != nil {
		log.Printf("Failed to migrate audit events table: %v", err)
	}

	return service
}

// LogEvent logs a new audit event
func (s *AuditService) LogEvent(eventType AuditEventType, category AuditCategory, severity AuditSeverity, userID *uuid.UUID, sessionID *uuid.UUID, ipAddress, userAgent, resource, action string, outcome AuditOutcome, description string, details map[string]interface{}) error {
	event := AuditEvent{
		ID:          uuid.New(),
		Timestamp:   time.Now(),
		EventType:   eventType,
		Category:    category,
		Severity:    severity,
		UserID:      userID,
		SessionID:   sessionID,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Resource:    resource,
		Action:      action,
		Outcome:     outcome,
		Description: description,
		Details:     details,
		Tags:        []string{},
	}

	// Add compliance flags based on event type and category
	event.ComplianceFlags = s.generateComplianceFlags(eventType, category, details)

	// Calculate risk score if applicable
	if riskScore := s.calculateRiskScore(eventType, category, outcome, details); riskScore > 0 {
		event.RiskScore = &riskScore
	}

	// Store the event
	if err := s.db.Create(&event).Error; err != nil {
		log.Printf("Failed to log audit event: %v", err)
		return fmt.Errorf("failed to log audit event: %w", err)
	}

	log.Printf("📋 Audit Event Logged: %s - %s - %s", eventType, resource, action)
	return nil
}

// LogAuthenticationEvent logs authentication-related events
func (s *AuditService) LogAuthenticationEvent(eventType AuditEventType, userID *uuid.UUID, sessionID *uuid.UUID, ipAddress, userAgent string, outcome AuditOutcome, details map[string]interface{}) error {
	var severity AuditSeverity
	switch outcome {
	case OutcomeSuccess:
		severity = AuditSeverityInfo
	case OutcomeFailure, OutcomeDenied:
		severity = AuditSeverityWarning
	case OutcomeError:
		severity = AuditSeverityError
	}

	description := fmt.Sprintf("Authentication event: %s", eventType)
	if userID != nil {
		description = fmt.Sprintf("Authentication event: %s for user %s", eventType, userID.String())
	}

	return s.LogEvent(eventType, CategoryAuthentication, severity, userID, sessionID, ipAddress, userAgent, "authentication", string(eventType), outcome, description, details)
}

// LogDataAccessEvent logs data access events
func (s *AuditService) LogDataAccessEvent(userID *uuid.UUID, sessionID *uuid.UUID, ipAddress, userAgent, resource, action string, outcome AuditOutcome, details map[string]interface{}) error {
	var eventType AuditEventType
	var severity AuditSeverity

	switch action {
	case "read", "view", "get":
		eventType = EventTypeDataAccess
		severity = AuditSeverityInfo
	case "create", "update", "modify":
		eventType = EventTypeDataModification
		severity = AuditSeverityInfo
	case "delete", "remove":
		eventType = EventTypeDataDeletion
		severity = AuditSeverityWarning
	case "export":
		eventType = EventTypeDataExport
		severity = AuditSeverityWarning
	case "import":
		eventType = EventTypeDataImport
		severity = AuditSeverityInfo
	default:
		eventType = EventTypeDataAccess
		severity = AuditSeverityInfo
	}

	if outcome != OutcomeSuccess {
		severity = AuditSeverityError
	}

	description := fmt.Sprintf("Data access: %s on %s", action, resource)

	return s.LogEvent(eventType, CategoryDataAccess, severity, userID, sessionID, ipAddress, userAgent, resource, action, outcome, description, details)
}

// LogSecurityEvent logs security-related events
func (s *AuditService) LogSecurityEvent(eventType AuditEventType, userID *uuid.UUID, ipAddress, userAgent, description string, details map[string]interface{}) error {
	severity := AuditSeverityCritical
	if eventType == EventTypeSuspiciousActivity {
		severity = AuditSeverityWarning
	}

	return s.LogEvent(eventType, CategorySecurity, severity, userID, nil, ipAddress, userAgent, "security", string(eventType), OutcomeSuccess, description, details)
}

// LogAdminEvent logs administrative events
func (s *AuditService) LogAdminEvent(adminUserID uuid.UUID, sessionID *uuid.UUID, ipAddress, userAgent, resource, action string, outcome AuditOutcome, description string, details map[string]interface{}) error {
	severity := AuditSeverityInfo
	if outcome != OutcomeSuccess {
		severity = AuditSeverityError
	}

	return s.LogEvent(EventTypeAdminAction, CategoryAdministrative, severity, &adminUserID, sessionID, ipAddress, userAgent, resource, action, outcome, description, details)
}

// LogAPIEvent logs API-related events
func (s *AuditService) LogAPIEvent(userID *uuid.UUID, ipAddress, userAgent, endpoint, method string, statusCode int, responseTime time.Duration, details map[string]interface{}) error {
	var eventType AuditEventType
	var outcome AuditOutcome
	var severity AuditSeverity

	if statusCode >= 200 && statusCode < 300 {
		eventType = EventTypeAPICall
		outcome = OutcomeSuccess
		severity = AuditSeverityInfo
	} else if statusCode == 429 {
		eventType = EventTypeRateLimitExceeded
		outcome = OutcomeFailure
		severity = AuditSeverityWarning
	} else if statusCode >= 400 {
		eventType = EventTypeAPIError
		outcome = OutcomeError
		severity = AuditSeverityError
	} else {
		eventType = EventTypeAPICall
		outcome = OutcomeSuccess
		severity = AuditSeverityInfo
	}

	if details == nil {
		details = make(map[string]interface{})
	}
	details["status_code"] = statusCode
	details["response_time_ms"] = responseTime.Milliseconds()
	details["method"] = method

	description := fmt.Sprintf("API call: %s %s (status: %d)", method, endpoint, statusCode)

	return s.LogEvent(eventType, CategoryAPI, severity, userID, nil, ipAddress, userAgent, endpoint, method, outcome, description, details)
}

// GetEvents retrieves audit events with filtering
func (s *AuditService) GetEvents(filter AuditFilter) ([]AuditEvent, error) {
	query := s.db.Model(&AuditEvent{})

	// Apply filters
	if filter.StartTime != nil {
		query = query.Where("timestamp >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		query = query.Where("timestamp <= ?", *filter.EndTime)
	}
	if len(filter.EventTypes) > 0 {
		query = query.Where("event_type IN ?", filter.EventTypes)
	}
	if len(filter.Categories) > 0 {
		query = query.Where("category IN ?", filter.Categories)
	}
	if len(filter.Severities) > 0 {
		query = query.Where("severity IN ?", filter.Severities)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.IPAddress != "" {
		query = query.Where("ip_address = ?", filter.IPAddress)
	}
	if filter.Resource != "" {
		query = query.Where("resource ILIKE ?", "%"+filter.Resource+"%")
	}
	if filter.Action != "" {
		query = query.Where("action ILIKE ?", "%"+filter.Action+"%")
	}
	if filter.Outcome != nil {
		query = query.Where("outcome = ?", *filter.Outcome)
	}
	if filter.RiskScoreMin != nil {
		query = query.Where("risk_score >= ?", *filter.RiskScoreMin)
	}
	if filter.RiskScoreMax != nil {
		query = query.Where("risk_score <= ?", *filter.RiskScoreMax)
	}
	if len(filter.Tags) > 0 {
		query = query.Where("tags && ?", filter.Tags)
	}
	if filter.CorrelationID != nil {
		query = query.Where("correlation_id = ?", *filter.CorrelationID)
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// Order by timestamp descending
	query = query.Order("timestamp DESC")

	var events []AuditEvent
	if err := query.Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve audit events: %w", err)
	}

	return events, nil
}

// GetStatistics generates audit statistics for a given time range
func (s *AuditService) GetStatistics(startTime, endTime time.Time) (*AuditStatistics, error) {
	stats := &AuditStatistics{
		EventsByType:     make(map[AuditEventType]int64),
		EventsByCategory: make(map[AuditCategory]int64),
		EventsBySeverity: make(map[AuditSeverity]int64),
		EventsByOutcome:  make(map[AuditOutcome]int64),
		TimeRange: AuditTimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
	}

	// Get total events count
	if err := s.db.Model(&AuditEvent{}).
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Count(&stats.TotalEvents).Error; err != nil {
		return nil, fmt.Errorf("failed to get total events count: %w", err)
	}

	// Get events by type
	var typeResults []struct {
		EventType AuditEventType `json:"event_type"`
		Count     int64          `json:"count"`
	}
	if err := s.db.Model(&AuditEvent{}).
		Select("event_type, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("event_type").
		Scan(&typeResults).Error; err != nil {
		return nil, fmt.Errorf("failed to get events by type: %w", err)
	}
	for _, result := range typeResults {
		stats.EventsByType[result.EventType] = result.Count
	}

	// Get events by category
	var categoryResults []struct {
		Category AuditCategory `json:"category"`
		Count    int64         `json:"count"`
	}
	if err := s.db.Model(&AuditEvent{}).
		Select("category, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("category").
		Scan(&categoryResults).Error; err != nil {
		return nil, fmt.Errorf("failed to get events by category: %w", err)
	}
	for _, result := range categoryResults {
		stats.EventsByCategory[result.Category] = result.Count
	}

	// Get events by severity
	var severityResults []struct {
		Severity AuditSeverity `json:"severity"`
		Count    int64         `json:"count"`
	}
	if err := s.db.Model(&AuditEvent{}).
		Select("severity, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("severity").
		Scan(&severityResults).Error; err != nil {
		return nil, fmt.Errorf("failed to get events by severity: %w", err)
	}
	for _, result := range severityResults {
		stats.EventsBySeverity[result.Severity] = result.Count
	}

	// Get events by outcome
	var outcomeResults []struct {
		Outcome AuditOutcome `json:"outcome"`
		Count   int64        `json:"count"`
	}
	if err := s.db.Model(&AuditEvent{}).
		Select("outcome, COUNT(*) as count").
		Where("timestamp BETWEEN ? AND ?", startTime, endTime).
		Group("outcome").
		Scan(&outcomeResults).Error; err != nil {
		return nil, fmt.Errorf("failed to get events by outcome: %w", err)
	}
	for _, result := range outcomeResults {
		stats.EventsByOutcome[result.Outcome] = result.Count
	}

	// Get security events count
	if err := s.db.Model(&AuditEvent{}).
		Where("timestamp BETWEEN ? AND ? AND category = ?", startTime, endTime, CategorySecurity).
		Count(&stats.SecurityEvents).Error; err != nil {
		return nil, fmt.Errorf("failed to get security events count: %w", err)
	}

	// Get failed attempts count
	if err := s.db.Model(&AuditEvent{}).
		Where("timestamp BETWEEN ? AND ? AND outcome IN ?", startTime, endTime, []AuditOutcome{OutcomeFailure, OutcomeDenied}).
		Count(&stats.FailedAttempts).Error; err != nil {
		return nil, fmt.Errorf("failed to get failed attempts count: %w", err)
	}

	// Get average risk score
	var avgRiskScore sql.NullFloat64
	if err := s.db.Model(&AuditEvent{}).
		Select("AVG(risk_score)").
		Where("timestamp BETWEEN ? AND ? AND risk_score IS NOT NULL", startTime, endTime).
		Scan(&avgRiskScore).Error; err != nil {
		return nil, fmt.Errorf("failed to get average risk score: %w", err)
	}
	if avgRiskScore.Valid {
		stats.AverageRiskScore = avgRiskScore.Float64
	}

	// Get compliance violations count
	if err := s.db.Model(&AuditEvent{}).
		Where("timestamp BETWEEN ? AND ? AND array_length(compliance_flags, 1) > 0", startTime, endTime).
		Count(&stats.ComplianceViolations).Error; err != nil {
		return nil, fmt.Errorf("failed to get compliance violations count: %w", err)
	}

	return stats, nil
}

// GenerateComplianceReport generates a comprehensive compliance report
func (s *AuditService) GenerateComplianceReport(reportType ComplianceReportType, startTime, endTime time.Time, generatedBy uuid.UUID) (*ComplianceReport, error) {
	report := &ComplianceReport{
		ID:          uuid.New(),
		ReportType:  reportType,
		TimeRange:   AuditTimeRange{StartTime: startTime, EndTime: endTime},
		GeneratedAt: time.Now(),
		GeneratedBy: generatedBy,
		Status:      ReportStatusGenerating,
	}

	// Generate statistics
	stats, err := s.GetStatistics(startTime, endTime)
	if err != nil {
		report.Status = ReportStatusFailed
		return report, fmt.Errorf("failed to generate statistics: %w", err)
	}
	report.Statistics = *stats

	// Generate compliance flags summary
	report.ComplianceFlags = make(map[string]int64)
	var flagResults []struct {
		Flag  string `json:"flag"`
		Count int64  `json:"count"`
	}
	if err := s.db.Raw(`
		SELECT unnest(compliance_flags) as flag, COUNT(*) as count
		FROM audit_events
		WHERE timestamp BETWEEN ? AND ?
		GROUP BY flag
	`, startTime, endTime).Scan(&flagResults).Error; err != nil {
		report.Status = ReportStatusFailed
		return report, fmt.Errorf("failed to get compliance flags: %w", err)
	}
	for _, result := range flagResults {
		report.ComplianceFlags[result.Flag] = result.Count
	}

	// Generate violations and recommendations based on report type
	report.Violations = s.generateComplianceViolations(reportType, startTime, endTime)
	report.Recommendations = s.generateComplianceRecommendations(reportType, report.Statistics)

	report.Status = ReportStatusCompleted
	return report, nil
}

// Helper methods

func (s *AuditService) generateComplianceFlags(eventType AuditEventType, category AuditCategory, details map[string]interface{}) []string {
	flags := make([]string, 0)

	// GDPR compliance flags
	if category == CategoryDataAccess || eventType == EventTypeDataExport {
		flags = append(flags, "gdpr-data-access")
	}
	if eventType == EventTypeDataDeletion {
		flags = append(flags, "gdpr-data-deletion")
	}

	// SOX compliance flags
	if category == CategoryAdministrative || eventType == EventTypeConfigurationChange {
		flags = append(flags, "sox-administrative-control")
	}

	// HIPAA compliance flags
	if category == CategoryDataAccess && details != nil {
		if sensitive, ok := details["sensitive_data"].(bool); ok && sensitive {
			flags = append(flags, "hipaa-phi-access")
		}
	}

	// SOC2 compliance flags
	if category == CategorySecurity || eventType == EventTypeSecurityAlert {
		flags = append(flags, "soc2-security-monitoring")
	}

	return flags
}

func (s *AuditService) calculateRiskScore(eventType AuditEventType, category AuditCategory, outcome AuditOutcome, details map[string]interface{}) float64 {
	score := 0.0

	// Base score by category
	switch category {
	case CategorySecurity:
		score += 0.8
	case CategoryAuthentication:
		score += 0.3
	case CategoryDataAccess:
		score += 0.4
	case CategoryAdministrative:
		score += 0.6
	}

	// Adjust by outcome
	switch outcome {
	case OutcomeFailure, OutcomeDenied:
		score += 0.3
	case OutcomeError:
		score += 0.5
	}

	// Adjust by event type
	switch eventType {
	case EventTypeSecurityAlert, EventTypeIntrusionDetected:
		score += 0.7
	case EventTypeLoginFailed, EventTypeMFAFailed:
		score += 0.2
	case EventTypePrivilegeEscalated:
		score += 0.9
	case EventTypeDataDeletion, EventTypeDataExport:
		score += 0.4
	}

	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}

	return score
}

func (s *AuditService) generateComplianceViolations(reportType ComplianceReportType, startTime, endTime time.Time) []ComplianceViolation {
	violations := make([]ComplianceViolation, 0)

	// This would be implemented based on specific compliance requirements
	// For now, return empty slice
	return violations
}

func (s *AuditService) generateComplianceRecommendations(reportType ComplianceReportType, stats AuditStatistics) []ComplianceRecommendation {
	recommendations := make([]ComplianceRecommendation, 0)

	// Generate recommendations based on statistics
	if stats.FailedAttempts > stats.TotalEvents/10 {
		recommendations = append(recommendations, ComplianceRecommendation{
			ID:          uuid.New(),
			Title:       "High Failed Authentication Rate",
			Description: "Consider implementing additional authentication controls or account lockout policies",
			Priority:    "high",
			Category:    "authentication",
		})
	}

	if stats.SecurityEvents > 0 {
		recommendations = append(recommendations, ComplianceRecommendation{
			ID:          uuid.New(),
			Title:       "Security Events Detected",
			Description: "Review security events and consider implementing additional monitoring",
			Priority:    "medium",
			Category:    "security",
		})
	}

	return recommendations
}
