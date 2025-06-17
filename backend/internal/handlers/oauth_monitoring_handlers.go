package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"cloudgate-backend/internal/services"
)

// GetConnectionsHandler retrieves all OAuth connections for a user with health data
func GetConnectionsHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	monitoringService := services.NewOAuthMonitoringService(services.GetDB())
	connections, err := monitoringService.GetUserConnections(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get connections", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"connections": connections,
		"count":       len(connections),
	})
}

// GetConnectionStatsHandler retrieves aggregated connection statistics
func GetConnectionStatsHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	monitoringService := services.NewOAuthMonitoringService(services.GetDB())
	stats, err := monitoringService.GetConnectionStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get connection stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// TestConnectionHandler performs a health check on a specific connection
func TestConnectionHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	connectionID := c.Param("connectionId")
	if connectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Connection ID is required"})
		return
	}

	monitoringService := services.NewOAuthMonitoringService(services.GetDB())
	err := monitoringService.TestConnection(userID, connectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to test connection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connection test completed"})
}

// GetSecurityEventsHandler retrieves security events for a user
func GetSecurityEventsHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse limit parameter
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	monitoringService := services.NewOAuthMonitoringService(services.GetDB())
	events, err := monitoringService.GetSecurityEvents(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get security events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"count":  len(events),
	})
}

// CreateSecurityEventHandler creates a new security event
type CreateSecurityEventRequest struct {
	EventType    string  `json:"event_type" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	Severity     string  `json:"severity" binding:"required"`
	Location     string  `json:"location,omitempty"`
	RiskScore    float64 `json:"risk_score,omitempty"`
	ConnectionID string  `json:"connection_id,omitempty"`
}

func CreateSecurityEventHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request CreateSecurityEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	monitoringService := services.NewOAuthMonitoringService(services.GetDB())

	var connectionID *string
	if request.ConnectionID != "" {
		connectionID = &request.ConnectionID
	}

	err := monitoringService.CreateSecurityEvent(
		userID,
		request.EventType,
		request.Description,
		request.Severity,
		ipAddress,
		userAgent,
		request.Location,
		request.RiskScore,
		connectionID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create security event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Security event created successfully"})
}

// GetTrustedDevicesHandler retrieves trusted devices for a user
func GetTrustedDevicesHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	monitoringService := services.NewOAuthMonitoringService(services.GetDB())
	devices, err := monitoringService.GetTrustedDevices(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trusted devices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
		"count":   len(devices),
	})
}

// RegisterDeviceHandler registers or updates a device for a user
type RegisterDeviceRequest struct {
	DeviceName  string `json:"device_name" binding:"required"`
	DeviceType  string `json:"device_type" binding:"required"`
	Browser     string `json:"browser" binding:"required"`
	OS          string `json:"os" binding:"required"`
	Fingerprint string `json:"fingerprint" binding:"required"`
	Location    string `json:"location,omitempty"`
}

func RegisterDeviceHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request RegisterDeviceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ipAddress := c.ClientIP()

	monitoringService := services.NewOAuthMonitoringService(services.GetDB())
	err := monitoringService.RegisterDevice(
		userID,
		request.DeviceName,
		request.DeviceType,
		request.Browser,
		request.OS,
		request.Fingerprint,
		ipAddress,
		request.Location,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device registered successfully"})
}

// TrustDeviceHandler marks a device as trusted
func TrustDeviceHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	deviceID := c.Param("deviceId")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID is required"})
		return
	}

	monitoringService := services.NewOAuthMonitoringService(services.GetDB())
	err := monitoringService.TrustDevice(userID, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to trust device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device trusted successfully"})
}

// RevokeDeviceHandler removes a device from trusted devices
func RevokeDeviceHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	deviceID := c.Param("deviceId")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID is required"})
		return
	}

	monitoringService := services.NewOAuthMonitoringService(services.GetDB())
	err := monitoringService.RevokeDevice(userID, deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device revoked successfully"})
}

// RecordUsageHandler records usage statistics for a connection
type RecordUsageRequest struct {
	ConnectionID    string `json:"connection_id" binding:"required"`
	DataTransferred int64  `json:"data_transferred,omitempty"`
}

func RecordUsageHandler(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request RecordUsageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	monitoringService := services.NewOAuthMonitoringService(services.GetDB())
	err := monitoringService.RecordUsage(userID, request.ConnectionID, request.DataTransferred)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record usage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usage recorded successfully"})
}
