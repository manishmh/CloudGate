package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"cloudgate-backend/internal/services"
)

// SecurityMonitoringHandlers contains handlers for security monitoring
type SecurityMonitoringHandlers struct {
	securityService *services.SecurityMonitoringService
}

// NewSecurityMonitoringHandlers creates new security monitoring handlers
func NewSecurityMonitoringHandlers(service *services.SecurityMonitoringService) *SecurityMonitoringHandlers {
	return &SecurityMonitoringHandlers{
		securityService: service,
	}
}

// GenerateAlertRequest represents the request for generating a security alert
type GenerateAlertRequest struct {
	Type        string                 `json:"type" binding:"required"`
	Severity    string                 `json:"severity" binding:"required"`
	Title       string                 `json:"title" binding:"required"`
	Description string                 `json:"description" binding:"required"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// AlertResponse represents a security alert in API responses
type AlertResponse struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Source      string                 `json:"source"`
	UserID      *string                `json:"user_id,omitempty"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Timestamp   time.Time              `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
	Status      string                 `json:"status"`
	AssignedTo  *string                `json:"assigned_to,omitempty"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
	Actions     []ActionResponse       `json:"actions"`
	Tags        []string               `json:"tags"`
}

// ActionResponse represents a security action in API responses
type ActionResponse struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
	PerformedBy string                 `json:"performed_by"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// IncidentResponse represents a security incident in API responses
type IncidentResponse struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Severity    string          `json:"severity"`
	Status      string          `json:"status"`
	Alerts      []AlertResponse `json:"alerts"`
	AssignedTo  *string         `json:"assigned_to,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	ResolvedAt  *time.Time      `json:"resolved_at,omitempty"`
	Timeline    []EventResponse `json:"timeline"`
}

// EventResponse represents an incident event in API responses
type EventResponse struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
	PerformedBy string                 `json:"performed_by"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// MetricsResponse represents security metrics in API responses
type MetricsResponse struct {
	AlertsGenerated   int64         `json:"alerts_generated"`
	AlertsResolved    int64         `json:"alerts_resolved"`
	FalsePositives    int64         `json:"false_positives"`
	IncidentsCreated  int64         `json:"incidents_created"`
	IncidentsResolved int64         `json:"incidents_resolved"`
	ResponseTime      time.Duration `json:"response_time"`
}

// CreateIncidentRequest represents the request for creating a security incident
type CreateIncidentRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Severity    string   `json:"severity" binding:"required"`
	AlertIDs    []string `json:"alert_ids" binding:"required"`
}

// UpdateAlertStatusRequest represents the request for updating alert status
type UpdateAlertStatusRequest struct {
	Status     string  `json:"status" binding:"required"`
	AssignedTo *string `json:"assigned_to,omitempty"`
}

// LoginEventRequest represents a login event for monitoring
type LoginEventRequest struct {
	UserID    string  `json:"user_id" binding:"required"`
	Email     string  `json:"email" binding:"required"`
	IPAddress string  `json:"ip_address" binding:"required"`
	UserAgent string  `json:"user_agent" binding:"required"`
	Success   bool    `json:"success"`
	RiskScore float64 `json:"risk_score"`
}

// APIEventRequest represents an API event for monitoring
type APIEventRequest struct {
	Endpoint     string `json:"endpoint" binding:"required"`
	Method       string `json:"method" binding:"required"`
	IPAddress    string `json:"ip_address" binding:"required"`
	UserAgent    string `json:"user_agent" binding:"required"`
	StatusCode   int    `json:"status_code" binding:"required"`
	ResponseTime int64  `json:"response_time_ms" binding:"required"`
}

// AlertChannelRequest represents a request to configure an alert channel
type AlertChannelRequest struct {
	Type    string                 `json:"type" binding:"required"`
	Name    string                 `json:"name" binding:"required"`
	Config  map[string]interface{} `json:"config" binding:"required"`
	Enabled bool                   `json:"enabled"`
}

// GenerateAlert creates a new security alert
func (h *SecurityMonitoringHandlers) GenerateAlert(c *gin.Context) {
	var req GenerateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	// Convert string types to service types
	alertType := services.AlertType(req.Type)
	severity := services.AlertSeverity(req.Severity)

	// Generate the alert
	alert, err := h.securityService.GenerateAlert(alertType, severity, req.Title, req.Description, req.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate alert",
			"message": err.Error(),
		})
		return
	}

	// Convert to response format
	response := convertAlertToResponse(*alert)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Alert generated successfully",
		"alert":   response,
	})
}

// GetAlerts retrieves security alerts with filtering
func (h *SecurityMonitoringHandlers) GetAlerts(c *gin.Context) {
	// Parse query parameters
	filters := services.AlertFilters{}

	if alertType := c.Query("type"); alertType != "" {
		t := services.AlertType(alertType)
		filters.Type = &t
	}

	if severity := c.Query("severity"); severity != "" {
		s := services.AlertSeverity(severity)
		filters.Severity = &s
	}

	if status := c.Query("status"); status != "" {
		st := services.AlertStatus(status)
		filters.Status = &st
	}

	if userID := c.Query("user_id"); userID != "" {
		if uid, err := uuid.Parse(userID); err == nil {
			filters.UserID = &uid
		}
	}

	if ipAddress := c.Query("ip_address"); ipAddress != "" {
		filters.IPAddress = ipAddress
	}

	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			filters.StartTime = &t
		}
	}

	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			filters.EndTime = &t
		}
	}

	// Parse pagination
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 1000 {
		limit = 50
	}
	filters.Limit = limit

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	filters.Offset = offset

	// Get alerts
	alerts, err := h.securityService.GetAlerts(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve alerts",
			"message": err.Error(),
		})
		return
	}

	// Convert to response format
	response := make([]AlertResponse, len(alerts))
	for i, alert := range alerts {
		response[i] = convertAlertToResponse(alert)
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": response,
		"count":  len(response),
		"limit":  limit,
		"offset": offset,
	})
}

// UpdateAlertStatus updates the status of a security alert
func (h *SecurityMonitoringHandlers) UpdateAlertStatus(c *gin.Context) {
	alertIDStr := c.Param("alert_id")
	if alertIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Missing alert ID",
			"message": "Alert ID is required",
		})
		return
	}

	alertID, err := uuid.Parse(alertIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid alert ID",
			"message": "Alert ID must be a valid UUID",
		})
		return
	}

	var req UpdateAlertStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	// Convert status
	status := services.AlertStatus(req.Status)

	// Parse assigned to if provided
	var assignedTo *uuid.UUID
	if req.AssignedTo != nil {
		if uid, err := uuid.Parse(*req.AssignedTo); err == nil {
			assignedTo = &uid
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid assigned_to ID",
				"message": "assigned_to must be a valid UUID",
			})
			return
		}
	}

	// Update alert status
	err = h.securityService.UpdateAlertStatus(alertID, status, assignedTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update alert status",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Alert status updated successfully",
		"alert_id": alertIDStr,
		"status":   req.Status,
	})
}

// CreateIncident creates a new security incident
func (h *SecurityMonitoringHandlers) CreateIncident(c *gin.Context) {
	var req CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	// Convert severity
	severity := services.AlertSeverity(req.Severity)

	// Parse alert IDs
	alertIDs := make([]uuid.UUID, len(req.AlertIDs))
	for i, idStr := range req.AlertIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid alert ID",
				"message": "All alert IDs must be valid UUIDs",
			})
			return
		}
		alertIDs[i] = id
	}

	// Create incident
	incident, err := h.securityService.CreateIncident(req.Title, req.Description, severity, alertIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create incident",
			"message": err.Error(),
		})
		return
	}

	// Convert to response format
	response := convertIncidentToResponse(*incident)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Incident created successfully",
		"incident": response,
	})
}

// GetIncidents retrieves security incidents with filtering
func (h *SecurityMonitoringHandlers) GetIncidents(c *gin.Context) {
	// Parse query parameters
	filters := services.IncidentFilters{}

	if status := c.Query("status"); status != "" {
		s := services.IncidentStatus(status)
		filters.Status = &s
	}

	if severity := c.Query("severity"); severity != "" {
		s := services.AlertSeverity(severity)
		filters.Severity = &s
	}

	if assignedTo := c.Query("assigned_to"); assignedTo != "" {
		if uid, err := uuid.Parse(assignedTo); err == nil {
			filters.AssignedTo = &uid
		}
	}

	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			filters.StartTime = &t
		}
	}

	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			filters.EndTime = &t
		}
	}

	// Parse pagination
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}
	filters.Limit = limit

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}
	filters.Offset = offset

	// Get incidents
	incidents, err := h.securityService.GetIncidents(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve incidents",
			"message": err.Error(),
		})
		return
	}

	// Convert to response format
	response := make([]IncidentResponse, len(incidents))
	for i, incident := range incidents {
		response[i] = convertIncidentToResponse(incident)
	}

	c.JSON(http.StatusOK, gin.H{
		"incidents": response,
		"count":     len(response),
		"limit":     limit,
		"offset":    offset,
	})
}

// GetSecurityMetrics returns current security monitoring metrics
func (h *SecurityMonitoringHandlers) GetSecurityMetrics(c *gin.Context) {
	metrics := h.securityService.GetSecurityMetrics()

	response := MetricsResponse{
		AlertsGenerated:   metrics.AlertsGenerated,
		AlertsResolved:    metrics.AlertsResolved,
		FalsePositives:    metrics.FalsePositives,
		IncidentsCreated:  metrics.IncidentsCreated,
		IncidentsResolved: metrics.IncidentsResolved,
		ResponseTime:      metrics.ResponseTime,
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": response,
	})
}

// ProcessLoginEvent processes a login event for security monitoring
func (h *SecurityMonitoringHandlers) ProcessLoginEvent(c *gin.Context) {
	var req LoginEventRequest
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

	// Process login event
	err = h.securityService.ProcessLoginEvent(userID, req.Email, req.IPAddress, req.UserAgent, req.Success, req.RiskScore)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process login event",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login event processed successfully",
	})
}

// ProcessAPIEvent processes an API event for security monitoring
func (h *SecurityMonitoringHandlers) ProcessAPIEvent(c *gin.Context) {
	var req APIEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	// Convert response time to duration
	responseTime := time.Duration(req.ResponseTime) * time.Millisecond

	// Process API event
	err := h.securityService.ProcessAPIEvent(req.Endpoint, req.Method, req.IPAddress, req.UserAgent, req.StatusCode, responseTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process API event",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API event processed successfully",
	})
}

// ConfigureAlertChannel configures an alert delivery channel
func (h *SecurityMonitoringHandlers) ConfigureAlertChannel(c *gin.Context) {
	var req AlertChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	// Create alert channel based on type
	var channel services.AlertChannel
	switch req.Type {
	case "email":
		channel = &services.EmailAlertChannel{
			Enabled: req.Enabled,
		}
		// Configure email-specific settings from req.Config
	case "slack":
		channel = &services.SlackAlertChannel{
			Enabled: req.Enabled,
		}
		// Configure Slack-specific settings from req.Config
	case "webhook":
		channel = &services.WebhookAlertChannel{
			Enabled: req.Enabled,
		}
		// Configure webhook-specific settings from req.Config
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid channel type",
			"message": "Supported types: email, slack, webhook",
		})
		return
	}

	// Add channel to service
	h.securityService.AddAlertChannel(req.Name, channel)

	c.JSON(http.StatusOK, gin.H{
		"message": "Alert channel configured successfully",
		"name":    req.Name,
		"type":    req.Type,
		"enabled": req.Enabled,
	})
}

// GetAlertTypes returns available alert types
func (h *SecurityMonitoringHandlers) GetAlertTypes(c *gin.Context) {
	alertTypes := []string{
		string(services.AlertTypeLoginAnomaly),
		string(services.AlertTypeMultipleFailedLogins),
		string(services.AlertTypeSuspiciousLocation),
		string(services.AlertTypeNewDeviceAccess),
		string(services.AlertTypeBruteForceAttack),
		string(services.AlertTypeAccountLockout),
		string(services.AlertTypePrivilegeEscalation),
		string(services.AlertTypeDataExfiltration),
		string(services.AlertTypeMaliciousIP),
		string(services.AlertTypeCompromisedAccount),
		string(services.AlertTypeUnauthorizedAccess),
		string(services.AlertTypeSessionHijacking),
		string(services.AlertTypeAPIAbuse),
		string(services.AlertTypeConfigurationChange),
		string(services.AlertTypeSystemIntegrityBreach),
	}

	c.JSON(http.StatusOK, gin.H{
		"alert_types": alertTypes,
	})
}

// GetAlertSeverities returns available alert severities
func (h *SecurityMonitoringHandlers) GetAlertSeverities(c *gin.Context) {
	severities := []string{
		string(services.SeverityLow),
		string(services.SeverityMedium),
		string(services.SeverityHigh),
		string(services.SeverityCritical),
	}

	c.JSON(http.StatusOK, gin.H{
		"severities": severities,
	})
}

// Helper functions to convert service types to response types

func convertAlertToResponse(alert services.SecurityAlert) AlertResponse {
	response := AlertResponse{
		ID:          alert.ID.String(),
		Type:        string(alert.Type),
		Severity:    string(alert.Severity),
		Title:       alert.Title,
		Description: alert.Description,
		Source:      alert.Source,
		IPAddress:   alert.IPAddress,
		UserAgent:   alert.UserAgent,
		Timestamp:   alert.Timestamp,
		Metadata:    alert.Metadata,
		Status:      string(alert.Status),
		ResolvedAt:  alert.ResolvedAt,
		Actions:     make([]ActionResponse, len(alert.Actions)),
		Tags:        alert.Tags,
	}

	if alert.UserID != nil {
		userIDStr := alert.UserID.String()
		response.UserID = &userIDStr
	}

	if alert.AssignedTo != nil {
		assignedToStr := alert.AssignedTo.String()
		response.AssignedTo = &assignedToStr
	}

	for i, action := range alert.Actions {
		response.Actions[i] = ActionResponse{
			ID:          action.ID.String(),
			Type:        string(action.Type),
			Description: action.Description,
			Timestamp:   action.Timestamp,
			PerformedBy: action.PerformedBy.String(),
			Status:      string(action.Status),
			Metadata:    action.Metadata,
		}
	}

	return response
}

func convertIncidentToResponse(incident services.SecurityIncident) IncidentResponse {
	response := IncidentResponse{
		ID:          incident.ID.String(),
		Title:       incident.Title,
		Description: incident.Description,
		Severity:    string(incident.Severity),
		Status:      string(incident.Status),
		Alerts:      make([]AlertResponse, len(incident.Alerts)),
		CreatedAt:   incident.CreatedAt,
		UpdatedAt:   incident.UpdatedAt,
		ResolvedAt:  incident.ResolvedAt,
		Timeline:    make([]EventResponse, len(incident.Timeline)),
	}

	if incident.AssignedTo != nil {
		assignedToStr := incident.AssignedTo.String()
		response.AssignedTo = &assignedToStr
	}

	for i, alert := range incident.Alerts {
		response.Alerts[i] = convertAlertToResponse(alert)
	}

	for i, event := range incident.Timeline {
		response.Timeline[i] = EventResponse{
			ID:          event.ID.String(),
			Type:        event.Type,
			Description: event.Description,
			Timestamp:   event.Timestamp,
			PerformedBy: event.PerformedBy.String(),
			Metadata:    event.Metadata,
		}
	}

	return response
}
