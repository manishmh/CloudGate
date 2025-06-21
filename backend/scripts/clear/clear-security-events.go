//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"cloudgate-backend/internal/models"
)

func main() {
	// Database connection - using SQLite like the main app
	db, err := gorm.Open(sqlite.Open("cloudgate.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Delete all security events
	result := db.Where("1 = 1").Delete(&models.SecurityEvent{})
	if result.Error != nil {
		log.Fatal("Failed to delete security events:", result.Error)
	}

	fmt.Printf("✓ Deleted %d security events\n", result.RowsAffected)
	fmt.Println("✓ Security events cleared successfully!")
}
