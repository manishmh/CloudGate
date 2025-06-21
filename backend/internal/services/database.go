package services

import (
	"fmt"
	"log"
	"os"
	"time"

	"cloudgate-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type     string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	URL      string // For Neon DATABASE_URL format
}

// InitializeDatabase initializes the database connection
func InitializeDatabase() error {
	config := getDatabaseConfig()

	var dialector gorm.Dialector

	switch config.Type {
	case "postgres":
		var dsn string
		// Use DATABASE_URL if provided (Neon format)
		if config.URL != "" {
			dsn = config.URL
		} else {
			// Build DSN from individual components
			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
				config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)
		}
		dialector = postgres.Open(dsn)
	case "sqlite":
		// Default to SQLite for development
		dbPath := config.DBName
		if dbPath == "" {
			dbPath = "cloudgate.db"
		}
		dialector = sqlite.Open(dbPath)
	default:
		return fmt.Errorf("unsupported database type: %s", config.Type)
	}

	// Configure GORM logger
	gormLogger := logger.Default
	if os.Getenv("GIN_MODE") == "release" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Only run migrations if explicitly requested via environment variable
	runMigrationsFlag := getEnv("RUN_MIGRATIONS", "false")
	if runMigrationsFlag == "true" {
		log.Println("🔄 Running database migrations...")
		if err := runMigrations(); err != nil {
			log.Printf("❌ Failed to run migrations: %v", err)
			// Don't fail startup for migration errors in production
			if os.Getenv("PORT") != "" { // Cloud Run environment
				log.Printf("⚠️ Continuing startup without migrations in Cloud Run environment")
			} else {
				return fmt.Errorf("failed to run migrations: %w", err)
			}
		} else {
			log.Printf("✅ Database migrations completed successfully")
		}
	} else {
		log.Println("ℹ️ Skipping database migrations (set RUN_MIGRATIONS=true to enable)")
	}

	log.Println("✅ Database initialized successfully")
	return nil
}

// getDatabaseConfig reads database configuration from environment variables
func getDatabaseConfig() DatabaseConfig {
	// Check for Neon DATABASE_URL first
	neonURL := getEnv("NEON_DATABASE_URL", "")
	if neonURL != "" {
		return DatabaseConfig{
			Type: "postgres",
			URL:  neonURL,
		}
	}

	// Check for standard DATABASE_URL
	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL != "" {
		return DatabaseConfig{
			Type: "postgres",
			URL:  databaseURL,
		}
	}

	// Fall back to individual environment variables
	return DatabaseConfig{
		Type:     getEnv("DB_TYPE", "sqlite"),
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "cloudgate"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "cloudgate.db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// runMigrations runs database migrations
func runMigrations() error {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.AppToken{},
		&models.AuditLog{},
		&models.EmailVerification{},
		&models.UserSettings{},
		&models.MFASetup{},
		&models.BackupCode{},
		&models.AppConnection{},
		&models.ConnectionHealthMetrics{},
		&models.SecurityEvent{},
		&models.TrustedDevice{},
		&RiskAssessment{},
		&RiskThresholds{},
		&DeviceFingerprint{},
		&WebAuthnCredential{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// HealthCheck performs a database health check
func DatabaseHealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
