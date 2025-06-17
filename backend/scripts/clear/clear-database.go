package main

import (
	"fmt"
	"log"
	"os"

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

	// Clear all data from tables
	fmt.Println("Clearing all data from tables...")

	// Clear in order to respect foreign key constraints
	tables := []interface{}{
		&models.SecurityEvent{},
		&models.TrustedDevice{},
		&models.AppConnection{},
		&models.UserSettings{},
	}

	for _, table := range tables {
		result := db.Where("1 = 1").Delete(table)
		if result.Error != nil {
			log.Printf("Warning: Failed to clear table %T: %v", table, result.Error)
		} else {
			fmt.Printf("✓ Cleared %d records from %T\n", result.RowsAffected, table)
		}
	}

	fmt.Println("✅ Database cleared successfully!")
}
