package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"cloudgate-backend/internal/models"
)

// OAuthMonitoringService handles OAuth connection monitoring
type OAuthMonitoringService struct {
	db *gorm.DB
}

// NewOAuthMonitoringService creates a new OAuth monitoring service
func NewOAuthMonitoringService(db *gorm.DB) *OAuthMonitoringService {
	return &OAuthMonitoringService{db: db}
}

// ConnectionStats represents aggregated connection statistics
type ConnectionStats struct {
	TotalConnections    int     `json:"total_connections"`
	ActiveConnections   int     `json:"active_connections"`
	FailedConnections   int     `json:"failed_connections"`
	AverageResponseTime int     `json:"average_response_time"`
	UptimePercentage    float64 `json:"uptime_percentage"`
}

// EnhancedConnection represents a connection with health and usage data
type EnhancedConnection struct {
	models.AppConnection
	Health          ConnectionHealth `json:"health"`
	UsageCount      int64            `json:"usage_count"`
	DataTransferred string           `json:"data_transferred"`
	LastUsed        *string          `json:"last_used,omitempty"`
}

// ConnectionHealth represents health status of a connection
type ConnectionHealth struct {
	Status       string  `json:"status"`
	LastCheck    string  `json:"last_check"`
	ResponseTime int     `json:"response_time"`
	Uptime       float64 `json:"uptime"`
	ErrorCount   int     `json:"error_count"`
}

// GetUserConnections retrieves all connections for a user with health data
func (s *OAuthMonitoringService) GetUserConnections(userID string) ([]EnhancedConnection, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var connections []models.AppConnection
	if err := s.db.Where("user_id = ?", userUUID).Find(&connections).Error; err != nil {
		return nil, fmt.Errorf("failed to get connections: %w", err)
	}

	var enhancedConnections []EnhancedConnection
	for _, conn := range connections {
		enhanced := EnhancedConnection{
			AppConnection: conn,
			Health: ConnectionHealth{
				Status:       conn.HealthStatus,
				LastCheck:    formatTime(conn.LastHealthCheck),
				ResponseTime: conn.ResponseTime,
				Uptime:       conn.UptimePercent,
				ErrorCount:   conn.ErrorCount,
			},
			UsageCount:      conn.UsageCount,
			DataTransferred: formatBytes(conn.DataTransferred),
		}

		if conn.LastUsed != nil {
			lastUsed := conn.LastUsed.Format(time.RFC3339)
			enhanced.LastUsed = &lastUsed
		}

		enhancedConnections = append(enhancedConnections, enhanced)
	}

	return enhancedConnections, nil
}

// GetConnectionStats calculates aggregated statistics for user connections
func (s *OAuthMonitoringService) GetConnectionStats(userID string) (*ConnectionStats, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var stats ConnectionStats

	// Get total connections
	var totalCount int64
	if err := s.db.Model(&models.AppConnection{}).Where("user_id = ?", userUUID).Count(&totalCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count total connections: %w", err)
	}
	stats.TotalConnections = int(totalCount)

	// Get active connections
	var activeCount int64
	if err := s.db.Model(&models.AppConnection{}).Where("user_id = ? AND status = ?", userUUID, "connected").Count(&activeCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count active connections: %w", err)
	}
	stats.ActiveConnections = int(activeCount)

	// Get failed connections
	var failedCount int64
	if err := s.db.Model(&models.AppConnection{}).Where("user_id = ? AND status = ?", userUUID, "error").Count(&failedCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count failed connections: %w", err)
	}
	stats.FailedConnections = int(failedCount)

	// Calculate average response time for active connections
	var avgResponseTime float64
	if err := s.db.Model(&models.AppConnection{}).Where("user_id = ? AND status = ?", userUUID, "connected").Select("AVG(response_time)").Scan(&avgResponseTime).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate average response time: %w", err)
	}
	stats.AverageResponseTime = int(avgResponseTime)

	// Calculate average uptime percentage
	var avgUptime float64
	if err := s.db.Model(&models.AppConnection{}).Where("user_id = ? AND status = ?", userUUID, "connected").Select("AVG(uptime_percent)").Scan(&avgUptime).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate average uptime: %w", err)
	}
	stats.UptimePercentage = avgUptime

	return &stats, nil
}

// TestConnection performs a health check on a specific connection
func (s *OAuthMonitoringService) TestConnection(userID, connectionID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	connUUID, err := uuid.Parse(connectionID)
	if err != nil {
		return fmt.Errorf("invalid connection ID: %w", err)
	}

	var connection models.AppConnection
	if err := s.db.Where("id = ? AND user_id = ?", connUUID, userUUID).First(&connection).Error; err != nil {
		return fmt.Errorf("connection not found: %w", err)
	}

	// Perform health check based on the provider
	startTime := time.Now()
	success, statusCode, errorMsg := s.performHealthCheck(&connection)
	responseTime := int(time.Since(startTime).Milliseconds())

	// Update connection health
	now := time.Now()
	updates := map[string]interface{}{
		"last_health_check": now,
		"response_time":     responseTime,
	}

	if success {
		updates["health_status"] = "healthy"
		updates["error_count"] = 0
	} else {
		updates["health_status"] = "error"
		updates["error_count"] = connection.ErrorCount + 1
		updates["last_error"] = errorMsg
		updates["last_error_at"] = now
	}

	if err := s.db.Model(&connection).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update connection health: %w", err)
	}

	// Record health metrics
	healthMetric := models.ConnectionHealthMetrics{
		ConnectionID:   connUUID,
		Timestamp:      now,
		ResponseTime:   responseTime,
		Success:        success,
		ErrorMessage:   errorMsg,
		HTTPStatusCode: statusCode,
	}

	if err := s.db.Create(&healthMetric).Error; err != nil {
		// Log error but don't fail the health check
		fmt.Printf("Failed to record health metrics: %v\n", err)
	}

	return nil
}

// performHealthCheck performs the actual health check based on provider
func (s *OAuthMonitoringService) performHealthCheck(connection *models.AppConnection) (success bool, statusCode int, errorMsg string) {
	// This is a simplified health check - in production, you'd make actual API calls
	// to each provider's health/user info endpoint using the stored access token

	switch connection.Provider {
	case "google":
		return s.checkGoogleHealth(connection)
	case "microsoft":
		return s.checkMicrosoftHealth(connection)
	case "slack":
		return s.checkSlackHealth(connection)
	case "github":
		return s.checkGitHubHealth(connection)
	default:
		// Generic health check
		return s.checkGenericHealth(connection)
	}
}

// checkGoogleHealth checks Google Workspace connection health
func (s *OAuthMonitoringService) checkGoogleHealth(connection *models.AppConnection) (bool, int, string) {
	// In production, make a call to Google's userinfo endpoint
	// For now, simulate based on token expiry and error count
	if connection.TokenExpiresAt != nil && connection.TokenExpiresAt.Before(time.Now()) {
		return false, 401, "Token expired"
	}
	if connection.ErrorCount > 5 {
		return false, 500, "Too many recent errors"
	}
	return true, 200, ""
}

// checkMicrosoftHealth checks Microsoft 365 connection health
func (s *OAuthMonitoringService) checkMicrosoftHealth(connection *models.AppConnection) (bool, int, string) {
	// Similar to Google, check Microsoft Graph API health
	if connection.TokenExpiresAt != nil && connection.TokenExpiresAt.Before(time.Now()) {
		return false, 401, "Token expired"
	}
	if connection.ErrorCount > 3 {
		return false, 429, "Rate limited due to errors"
	}
	return true, 200, ""
}

// checkSlackHealth checks Slack connection health
func (s *OAuthMonitoringService) checkSlackHealth(connection *models.AppConnection) (bool, int, string) {
	// Check Slack API health
	if connection.ErrorCount > 10 {
		return false, 503, "Service unavailable"
	}
	return true, 200, ""
}

// checkGitHubHealth checks GitHub connection health
func (s *OAuthMonitoringService) checkGitHubHealth(connection *models.AppConnection) (bool, int, string) {
	// Check GitHub API health
	if connection.TokenExpiresAt != nil && connection.TokenExpiresAt.Before(time.Now()) {
		return false, 401, "Token expired"
	}
	return true, 200, ""
}

// checkGenericHealth performs a generic health check
func (s *OAuthMonitoringService) checkGenericHealth(connection *models.AppConnection) (bool, int, string) {
	// Generic health check logic
	if connection.ErrorCount > 5 {
		return false, 500, "Too many errors"
	}
	return true, 200, ""
}

// RecordUsage records usage statistics for a connection
func (s *OAuthMonitoringService) RecordUsage(userID, connectionID string, dataTransferred int64) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	connUUID, err := uuid.Parse(connectionID)
	if err != nil {
		return fmt.Errorf("invalid connection ID: %w", err)
	}

	now := time.Now()
	updates := map[string]interface{}{
		"usage_count":      gorm.Expr("usage_count + 1"),
		"data_transferred": gorm.Expr("data_transferred + ?", dataTransferred),
		"last_used":        now,
	}

	return s.db.Model(&models.AppConnection{}).
		Where("id = ? AND user_id = ?", connUUID, userUUID).
		Updates(updates).Error
}

// GetSecurityEvents retrieves security events for a user
func (s *OAuthMonitoringService) GetSecurityEvents(userID string, limit int) ([]models.SecurityEvent, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var events []models.SecurityEvent
	query := s.db.Where("user_id = ?", userUUID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to get security events: %w", err)
	}

	return events, nil
}

// CreateSecurityEvent creates a new security event
func (s *OAuthMonitoringService) CreateSecurityEvent(userID string, eventType, description, severity, ipAddress, userAgent, location string, riskScore float64, connectionID *string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	event := models.SecurityEvent{
		UserID:      userUUID,
		EventType:   eventType,
		Description: description,
		Severity:    severity,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Location:    location,
		RiskScore:   riskScore,
	}

	if connectionID != nil {
		connUUID, err := uuid.Parse(*connectionID)
		if err == nil {
			event.ConnectionID = &connUUID
		}
	}

	return s.db.Create(&event).Error
}

// GetTrustedDevices retrieves trusted devices for a user
func (s *OAuthMonitoringService) GetTrustedDevices(userID string) ([]models.TrustedDevice, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var devices []models.TrustedDevice
	if err := s.db.Where("user_id = ?", userUUID).Order("last_seen DESC").Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("failed to get trusted devices: %w", err)
	}

	return devices, nil
}

// RegisterDevice registers or updates a device for a user
func (s *OAuthMonitoringService) RegisterDevice(userID, deviceName, deviceType, browser, os, fingerprint, ipAddress, location string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if device already exists
	var device models.TrustedDevice
	err = s.db.Where("user_id = ? AND fingerprint = ?", userUUID, fingerprint).First(&device).Error

	if err == gorm.ErrRecordNotFound {
		// Create new device
		device = models.TrustedDevice{
			UserID:      userUUID,
			DeviceName:  deviceName,
			DeviceType:  deviceType,
			Browser:     browser,
			OS:          os,
			Fingerprint: fingerprint,
			IPAddress:   ipAddress,
			Location:    location,
			Trusted:     false, // New devices are not trusted by default
			LastSeen:    time.Now(),
		}
		return s.db.Create(&device).Error
	} else if err != nil {
		return fmt.Errorf("failed to check existing device: %w", err)
	}

	// Update existing device
	updates := map[string]interface{}{
		"last_seen":   time.Now(),
		"ip_address":  ipAddress,
		"location":    location,
		"device_name": deviceName,
		"browser":     browser,
		"os":          os,
	}

	return s.db.Model(&device).Updates(updates).Error
}

// TrustDevice marks a device as trusted
func (s *OAuthMonitoringService) TrustDevice(userID, deviceID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	deviceUUID, err := uuid.Parse(deviceID)
	if err != nil {
		return fmt.Errorf("invalid device ID: %w", err)
	}

	return s.db.Model(&models.TrustedDevice{}).
		Where("id = ? AND user_id = ?", deviceUUID, userUUID).
		Update("trusted", true).Error
}

// RevokeDevice removes a device from trusted devices
func (s *OAuthMonitoringService) RevokeDevice(userID, deviceID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	deviceUUID, err := uuid.Parse(deviceID)
	if err != nil {
		return fmt.Errorf("invalid device ID: %w", err)
	}

	return s.db.Where("id = ? AND user_id = ?", deviceUUID, userUUID).Delete(&models.TrustedDevice{}).Error
}

// Helper functions

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
