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

var DB *gorm.DB

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
