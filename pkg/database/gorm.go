// Package database provides database initialization and connection management.
// It uses GORM with SQLite for data persistence and handles schema migrations.
package database

import (
	"fmt"
	"log"
	"os"

	"github.com/tofiquem/assingment/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ==================== Global Instance ====================

// DB is the global database connection pool.
// It is initialized by InitDB and should be closed with CloseDB on shutdown.
var DB *gorm.DB

// ==================== Initialization ====================

// InitDB initializes the database connection and runs schema migrations.
// It uses SQLite with the path from DATABASE_PATH env var or defaults.
// For serverless environments (Vercel), it uses /tmp/salary_management.db.
// It also creates a default admin user if no users exist.
func InitDB() {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		// Use /tmp for serverless environments (Vercel, etc.)
		if os.Getenv("VERCEL") == "1" {
			dbPath = "/tmp/salary_management.db"
		} else {
			dbPath = "./salary_management.db"
		}
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the database schema
	if err := DB.AutoMigrate(&models.User{}, &models.Employee{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create default admin user if none exists
	createDefaultAdminUser()

	log.Println("Successfully connected to database and migrated schema")
}

// ==================== Cleanup ====================

// CloseDB closes the database connection.
// It should be called on application shutdown to release resources.
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting underlying database connection: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}

// ==================== Default Data ====================

// createDefaultAdminUser creates a default admin user if no users exist in the database.
// This ensures there is always an initial user to log in with.
// Default credentials: admin@company.com / admin123
func createDefaultAdminUser() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		// Hash default password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Warning: failed to hash default admin password: %v", err)
			return
		}

		// Create default HR admin
		admin := models.User{
			Email:     "admin@company.com",
			Password:  string(hashedPassword),
			Role:      "admin",
			FirstName: "System",
			LastName:  "Administrator",
			IsActive:  true,
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Printf("Warning: failed to create default admin user: %v", err)
			return
		}

		fmt.Println("Created default admin user:")
		fmt.Println("  Email: admin@company.com")
		fmt.Println("  Password: admin123")
		fmt.Println("  Role: admin")
		fmt.Println("Please change the default password after first login!")
	}
}
