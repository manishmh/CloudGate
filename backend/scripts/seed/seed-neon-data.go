package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"cloudgate-backend/internal/models"
)

func main() {
	// Database connection - using Neon PostgreSQL
	databaseURL := "postgresql://neondb_owner:npg_AIC3QLgYf0qz@ep-cool-mud-a15d3oih-pooler.ap-southeast-1.aws.neon.tech/neondb?sslmode=require"

	// Allow override from environment
	if envURL := os.Getenv("DATABASE_URL"); envURL != "" {
		databaseURL = envURL
	}

	fmt.Println("Connecting to Neon PostgreSQL database...")
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("✓ Connected to database")

	// Run migrations
	fmt.Println("Running database migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.UserSettings{},
		&models.AppConnection{},
		&models.ConnectionHealthMetrics{},
		&models.SecurityEvent{},
		&models.TrustedDevice{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	fmt.Println("✓ Database migrations completed")

	// Demo user UUID (matches the one in middleware)
	demoUserUUID, _ := uuid.Parse("12345678-1234-1234-1234-123456789012")

	// Create demo user settings with new security fields
	userSettings := &models.UserSettings{
		UserID:             demoUserUUID,
		Language:           "en",
		Timezone:           "UTC",
		DateFormat:         "MM/DD/YYYY",
		EmailNotifications: true,
		PushNotifications:  true,
		SecurityAlerts:     true,
		AppUpdates:         true,
		WeeklyReports:      false,
		DefaultView:        "dashboard",
		ItemsPerPage:       10,
		AutoRefresh:        true,
		RefreshInterval:    30,
		AnalyticsOptIn:     true,
		ShareUsageData:     false,
		PersonalizedAds:    false,
		APIAccess:          false,
		MaxAPICalls:        1000,
		// New security fields
		TwoFactorEnabled:         true,
		LoginNotifications:       true,
		SuspiciousActivityAlerts: true,
		SessionTimeout:           720, // 12 hours in minutes
		PasswordExpiryDays:       90,
	}

	// Insert or update user settings
	var existingSettings models.UserSettings
	if err := db.Where("user_id = ?", demoUserUUID).First(&existingSettings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(userSettings).Error; err != nil {
				log.Printf("Failed to create user settings: %v", err)
			} else {
				fmt.Println("✓ Created demo user settings")
			}
		}
	} else {
		// Update existing settings with new fields
		if err := db.Model(&existingSettings).Updates(userSettings).Error; err != nil {
			log.Printf("Failed to update user settings: %v", err)
		} else {
			fmt.Println("✓ Updated demo user settings with new security fields")
		}
	}

	// Create demo OAuth connections
	connections := []models.AppConnection{
		{
			UserID:          demoUserUUID,
			AppID:           "github",
			AppName:         "GitHub",
			Provider:        "github",
			Status:          "connected",
			UserEmail:       "manishmh982@gmail.com",
			UserName:        "manishmh982",
			AccessToken:     "demo_github_token",
			RefreshToken:    "demo_github_refresh",
			TokenExpiresAt:  &[]time.Time{time.Now().Add(24 * time.Hour)}[0],
			Scopes:          "read:user,repo",
			ConnectedAt:     time.Now().Add(-7 * 24 * time.Hour),
			LastUsed:        &[]time.Time{time.Now().Add(-2 * time.Hour)}[0],
			HealthStatus:    "healthy",
			LastHealthCheck: &[]time.Time{time.Now().Add(-5 * time.Minute)}[0],
			ResponseTime:    150,
			UptimePercent:   99.8,
			ErrorCount:      0,
			UsageCount:      42,
			DataTransferred: 1024 * 1024 * 15, // 15MB
		},
		{
			UserID:          demoUserUUID,
			AppID:           "google-workspace",
			AppName:         "Google Workspace",
			Provider:        "google",
			Status:          "connected",
			UserEmail:       "manishmh982@gmail.com",
			UserName:        "Manish Kumar Saw",
			AccessToken:     "demo_google_token",
			RefreshToken:    "demo_google_refresh",
			TokenExpiresAt:  &[]time.Time{time.Now().Add(24 * time.Hour)}[0],
			Scopes:          "email,profile,drive",
			ConnectedAt:     time.Now().Add(-14 * 24 * time.Hour),
			LastUsed:        &[]time.Time{time.Now().Add(-30 * time.Minute)}[0],
			HealthStatus:    "healthy",
			LastHealthCheck: &[]time.Time{time.Now().Add(-2 * time.Minute)}[0],
			ResponseTime:    89,
			UptimePercent:   100.0,
			ErrorCount:      0,
			UsageCount:      156,
			DataTransferred: 1024 * 1024 * 45, // 45MB
		},
	}

	for _, conn := range connections {
		var existing models.AppConnection
		if err := db.Where("user_id = ? AND app_id = ?", demoUserUUID, conn.AppID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&conn).Error; err != nil {
					log.Printf("Failed to create connection %s: %v", conn.AppName, err)
				} else {
					fmt.Printf("✓ Created demo connection: %s\n", conn.AppName)
				}
			}
		} else {
			fmt.Printf("✓ Demo connection already exists: %s\n", conn.AppName)
		}
	}

	// Create demo security events
	securityEvents := []models.SecurityEvent{
		{
			UserID:      demoUserUUID,
			EventType:   "login",
			Description: "Successful login from Chrome browser",
			Severity:    "low",
			IPAddress:   "192.168.1.100",
			UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			Location:    "New York, US",
			RiskScore:   1.2,
			Resolved:    true,
		},
		{
			UserID:      demoUserUUID,
			EventType:   "device_registered",
			Description: "New device registered: Chrome on Windows",
			Severity:    "low",
			IPAddress:   "192.168.1.100",
			UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			Location:    "New York, US",
			RiskScore:   2.1,
			Resolved:    true,
		},
		{
			UserID:      demoUserUUID,
			EventType:   "failed_login",
			Description: "Failed login attempt - incorrect password",
			Severity:    "medium",
			IPAddress:   "203.0.113.1",
			UserAgent:   "Mozilla/5.0 (Unknown) Suspicious/1.0",
			Location:    "Unknown",
			RiskScore:   6.5,
			Resolved:    false,
		},
		{
			UserID:      demoUserUUID,
			EventType:   "oauth_connected",
			Description: "Connected to GitHub successfully",
			Severity:    "low",
			IPAddress:   "192.168.1.100",
			UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			Location:    "New York, US",
			RiskScore:   1.0,
			Resolved:    true,
		},
	}

	for i, event := range securityEvents {
		event.CreatedAt = time.Now().Add(-time.Duration(i*24) * time.Hour)
		var existing models.SecurityEvent
		if err := db.Where("user_id = ? AND event_type = ? AND description = ?",
			demoUserUUID, event.EventType, event.Description).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&event).Error; err != nil {
					log.Printf("Failed to create security event: %v", err)
				} else {
					fmt.Printf("✓ Created security event: %s\n", event.EventType)
				}
			}
		} else {
			fmt.Printf("✓ Security event already exists: %s\n", event.EventType)
		}
	}

	// Create demo trusted devices
	trustedDevices := []models.TrustedDevice{
		{
			UserID:      demoUserUUID,
			DeviceName:  "Windows PC",
			DeviceType:  "desktop",
			Browser:     "Chrome",
			OS:          "Windows 10",
			Fingerprint: "windows-pc-chrome-fingerprint-001",
			IPAddress:   "192.168.1.100",
			Location:    "New York, US",
			LastSeen:    time.Now().Add(-2 * time.Hour),
			Trusted:     true,
		},
		{
			UserID:      demoUserUUID,
			DeviceName:  "iPhone",
			DeviceType:  "mobile",
			Browser:     "Safari",
			OS:          "iOS 17",
			Fingerprint: "iphone-safari-fingerprint-002",
			IPAddress:   "192.168.1.101",
			Location:    "New York, US",
			LastSeen:    time.Now().Add(-6 * time.Hour),
			Trusted:     true,
		},
		{
			UserID:      demoUserUUID,
			DeviceName:  "Unknown Device",
			DeviceType:  "desktop",
			Browser:     "Unknown",
			OS:          "Unknown",
			Fingerprint: "unknown-device-fingerprint-003",
			IPAddress:   "203.0.113.1",
			Location:    "Unknown",
			LastSeen:    time.Now().Add(-3 * 24 * time.Hour),
			Trusted:     false,
		},
	}

	for _, device := range trustedDevices {
		var existing models.TrustedDevice
		if err := db.Where("user_id = ? AND device_name = ?", demoUserUUID, device.DeviceName).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&device).Error; err != nil {
					log.Printf("Failed to create trusted device: %v", err)
				} else {
					fmt.Printf("✓ Created trusted device: %s\n", device.DeviceName)
				}
			}
		} else {
			fmt.Printf("✓ Trusted device already exists: %s\n", device.DeviceName)
		}
	}

	fmt.Println("\n🎉 Demo data seeding completed successfully!")
	fmt.Println("You can now start the backend server and it will have demo data.")
}
