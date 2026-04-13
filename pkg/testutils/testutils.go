// Package testutils provides utilities for testing the salary management application.
// It includes helpers for database setup, test data creation, HTTP request handling,
// and response assertions.
package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/tofiquem/assingment/pkg/database"
	"github.com/tofiquem/assingment/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ==================== Database Setup ====================

// TestDB creates an in-memory SQLite database for testing.
// It auto-migrates the User and Employee models.
func TestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&models.User{}, &models.Employee{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// CleanupTestDB closes the test database connection.
func CleanupTestDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// ==================== Test Data Creation ====================

// CreateTestEmployee creates a single test employee record in the database.
func CreateTestEmployee(db *gorm.DB, firstName, lastName, email, jobTitle, country, department string, salary float64) *models.Employee {
	employee := &models.Employee{
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		JobTitle:   jobTitle,
		Country:    country,
		Salary:     salary,
		Department: department,
	}

	if err := db.Create(employee).Error; err != nil {
		panic(err)
	}

	return employee
}

// CreateTestEmployees creates multiple predefined test employees for testing.
func CreateTestEmployees(db *gorm.DB) []models.Employee {
	employees := []models.Employee{
		{FirstName: "John", LastName: "Doe", Email: "john@example.com", JobTitle: "Developer", Country: "USA", Salary: 75000.0, Department: "Engineering"},
		{FirstName: "Jane", LastName: "Smith", Email: "jane@example.com", JobTitle: "Manager", Country: "UK", Salary: 85000.0, Department: "Management"},
		{FirstName: "Bob", LastName: "Johnson", Email: "bob@example.com", JobTitle: "Developer", Country: "USA", Salary: 80000.0, Department: "Engineering"},
		{FirstName: "Alice", LastName: "Williams", Email: "alice@example.com", JobTitle: "Designer", Country: "Canada", Salary: 70000.0, Department: "Design"},
		{FirstName: "Charlie", LastName: "Brown", Email: "charlie@example.com", JobTitle: "Developer", Country: "USA", Salary: 90000.0, Department: "Engineering"},
	}

	for _, emp := range employees {
		if err := db.Create(&emp).Error; err != nil {
			panic(err)
		}
	}

	return employees
}

// MockDB sets up a mock database and overrides the global database.DB variable.
// Returns the test database for cleanup.
func MockDB(t *testing.T) *gorm.DB {
	// Create test DB
	testDB := TestDB(t)

	// Override global DB for testing
	database.DB = testDB

	// Return cleanup function
	return testDB
}

// RestoreDB restores the original database connection to the global database.DB variable.
func RestoreDB(originalDB *gorm.DB) {
	database.DB = originalDB
}

// ==================== HTTP Request Helpers ====================

// CreateJSONRequest creates an HTTP request with a JSON body.
// Sets the Content-Type header to application/json.
func CreateJSONRequest(method, url string, body interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// ExecuteRequest executes an HTTP request against a handler and returns the response recorder.
func ExecuteRequest(handler http.HandlerFunc, method, url string, body interface{}) (*httptest.ResponseRecorder, error) {
	req, err := CreateJSONRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr, nil
}

// ==================== Response Assertions ====================

// ParseJSONResponse parses a JSON response body into the provided struct.
func ParseJSONResponse(rr *httptest.ResponseRecorder, v interface{}) error {
	return json.Unmarshal(rr.Body.Bytes(), v)
}

// AssertStatusCode checks if the response status code matches the expected value.
func AssertStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	if rr.Code != expected {
		t.Errorf("Expected status code %d, got %d", expected, rr.Code)
	}
}

// AssertContentType checks if the response Content-Type header matches the expected value.
func AssertContentType(t *testing.T, rr *httptest.ResponseRecorder, expected string) {
	contentType := rr.Header().Get("Content-Type")
	if contentType != expected {
		t.Errorf("Expected content type %s, got %s", expected, contentType)
	}
}

// ==================== Environment Helpers ====================

// SetTestEnv sets an environment variable and returns a cleanup function to restore the original value.
func SetTestEnv(key, value string) func() {
	original := os.Getenv(key)
	os.Setenv(key, value)
	return func() {
		if original == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, original)
		}
	}
}
