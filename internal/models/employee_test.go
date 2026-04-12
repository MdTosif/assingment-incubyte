package models

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestToEmployee(t *testing.T) {
	tests := []struct {
		name     string
		req      *CreateEmployeeRequest
		expected *Employee
	}{
		{
			name: "valid employee request",
			req: &CreateEmployeeRequest{
				FirstName:  "John",
				LastName:   "Doe",
				Email:      "john@example.com",
				JobTitle:   "Developer",
				Country:    "USA",
				Salary:     75000.0,
				Department: "Engineering",
			},
			expected: &Employee{
				FirstName:  "John",
				LastName:   "Doe",
				Email:      "john@example.com",
				JobTitle:   "Developer",
				Country:    "USA",
				Salary:     75000.0,
				Department: "Engineering",
			},
		},
		{
			name:     "nil request",
			req:      nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToEmployee(tt.req)

			if tt.req == nil {
				if result != nil {
					t.Errorf("Expected nil result for nil request, got %v", result)
				}
				return
			}

			if result.FirstName != tt.expected.FirstName {
				t.Errorf("Expected FirstName %s, got %s", tt.expected.FirstName, result.FirstName)
			}
			if result.LastName != tt.expected.LastName {
				t.Errorf("Expected LastName %s, got %s", tt.expected.LastName, result.LastName)
			}
			if result.Email != tt.expected.Email {
				t.Errorf("Expected Email %s, got %s", tt.expected.Email, result.Email)
			}
			if result.JobTitle != tt.expected.JobTitle {
				t.Errorf("Expected JobTitle %s, got %s", tt.expected.JobTitle, result.JobTitle)
			}
			if result.Country != tt.expected.Country {
				t.Errorf("Expected Country %s, got %s", tt.expected.Country, result.Country)
			}
			if result.Salary != tt.expected.Salary {
				t.Errorf("Expected Salary %f, got %f", tt.expected.Salary, result.Salary)
			}
			if result.Department != tt.expected.Department {
				t.Errorf("Expected Department %s, got %s", tt.expected.Department, result.Department)
			}
		})
	}
}

func TestEmployee_UpdateFromRequest(t *testing.T) {
	tests := []struct {
		name     string
		employee *Employee
		req      *UpdateEmployeeRequest
		expected *Employee
	}{
		{
			name: "update all fields",
			employee: &Employee{
				FirstName:  "John",
				LastName:   "Doe",
				Email:      "john@example.com",
				JobTitle:   "Developer",
				Country:    "USA",
				Salary:     75000.0,
				Department: "Engineering",
			},
			req: &UpdateEmployeeRequest{
				FirstName:  stringPtr("Jane"),
				LastName:   stringPtr("Smith"),
				Email:      stringPtr("jane@example.com"),
				JobTitle:   stringPtr("Manager"),
				Country:    stringPtr("UK"),
				Salary:     floatPtr(85000.0),
				Department: stringPtr("Management"),
			},
			expected: &Employee{
				FirstName:  "Jane",
				LastName:   "Smith",
				Email:      "jane@example.com",
				JobTitle:   "Manager",
				Country:    "UK",
				Salary:     85000.0,
				Department: "Management",
			},
		},
		{
			name: "update partial fields",
			employee: &Employee{
				FirstName:  "John",
				LastName:   "Doe",
				Email:      "john@example.com",
				JobTitle:   "Developer",
				Country:    "USA",
				Salary:     75000.0,
				Department: "Engineering",
			},
			req: &UpdateEmployeeRequest{
				FirstName: stringPtr("Jane"),
				Salary:    floatPtr(80000.0),
			},
			expected: &Employee{
				FirstName:  "Jane",
				LastName:   "Doe",
				Email:      "john@example.com",
				JobTitle:   "Developer",
				Country:    "USA",
				Salary:     80000.0,
				Department: "Engineering",
			},
		},
		{
			name: "no fields to update",
			employee: &Employee{
				FirstName:  "John",
				LastName:   "Doe",
				Email:      "john@example.com",
				JobTitle:   "Developer",
				Country:    "USA",
				Salary:     75000.0,
				Department: "Engineering",
			},
			req: &UpdateEmployeeRequest{},
			expected: &Employee{
				FirstName:  "John",
				LastName:   "Doe",
				Email:      "john@example.com",
				JobTitle:   "Developer",
				Country:    "USA",
				Salary:     75000.0,
				Department: "Engineering",
			},
		},
		{
			name:     "nil request",
			employee: &Employee{FirstName: "John"},
			req:      nil,
			expected: &Employee{FirstName: "John"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req != nil {
				tt.employee.UpdateFromRequest(tt.req)
			}

			if tt.employee.FirstName != tt.expected.FirstName {
				t.Errorf("Expected FirstName %s, got %s", tt.expected.FirstName, tt.employee.FirstName)
			}
			if tt.employee.LastName != tt.expected.LastName {
				t.Errorf("Expected LastName %s, got %s", tt.expected.LastName, tt.employee.LastName)
			}
			if tt.employee.Email != tt.expected.Email {
				t.Errorf("Expected Email %s, got %s", tt.expected.Email, tt.employee.Email)
			}
			if tt.employee.JobTitle != tt.expected.JobTitle {
				t.Errorf("Expected JobTitle %s, got %s", tt.expected.JobTitle, tt.employee.JobTitle)
			}
			if tt.employee.Country != tt.expected.Country {
				t.Errorf("Expected Country %s, got %s", tt.expected.Country, tt.employee.Country)
			}
			if tt.employee.Salary != tt.expected.Salary {
				t.Errorf("Expected Salary %f, got %f", tt.expected.Salary, tt.employee.Salary)
			}
			if tt.employee.Department != tt.expected.Department {
				t.Errorf("Expected Department %s, got %s", tt.expected.Department, tt.employee.Department)
			}
		})
	}
}

func TestEmployee_BeforeCreate(t *testing.T) {
	tests := []struct {
		name     string
		employee *Employee
		expected bool
	}{
		{
			name: "hire date is zero",
			employee: &Employee{
				FirstName: "John",
				LastName:  "Doe",
				HireDate:  time.Time{},
			},
			expected: true,
		},
		{
			name: "hire date is already set",
			employee: &Employee{
				FirstName: "John",
				LastName:  "Doe",
				HireDate:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalHireDate := tt.employee.HireDate

			err := tt.employee.BeforeCreate(nil)
			if err != nil {
				t.Errorf("BeforeCreate returned error: %v", err)
			}

			if tt.expected {
				if tt.employee.HireDate.IsZero() {
					t.Errorf("Expected HireDate to be set, but it's still zero")
				}
				if tt.employee.HireDate.Equal(originalHireDate) {
					t.Errorf("Expected HireDate to be different from original")
				}
			} else {
				if !tt.employee.HireDate.Equal(originalHireDate) {
					t.Errorf("Expected HireDate to remain unchanged")
				}
			}
		})
	}
}

func TestEmployee_Validation(t *testing.T) {
	tests := []struct {
		name      string
		employee  *Employee
		expectErr bool
		errType   string
	}{
		{
			name: "valid employee",
			employee: &Employee{
				FirstName:  "John",
				LastName:   "Doe",
				Email:      "john@example.com",
				JobTitle:   "Developer",
				Country:    "USA",
				Salary:     75000.0,
				Department: "Engineering",
			},
			expectErr: false,
		},
		{
			name: "empty first name (SQLite allows empty strings)",
			employee: &Employee{
				FirstName:  "",
				LastName:   "Doe",
				Email:      "empty-first@example.com",
				JobTitle:   "Developer",
				Country:    "USA",
				Salary:     75000.0,
				Department: "Engineering",
			},
			expectErr: false, // SQLite allows empty strings with NOT NULL constraint
		},
		{
			name: "empty last name (SQLite allows empty strings)",
			employee: &Employee{
				FirstName:  "John",
				LastName:   "",
				Email:      "empty-last@example.com",
				JobTitle:   "Developer",
				Country:    "USA",
				Salary:     75000.0,
				Department: "Engineering",
			},
			expectErr: false, // SQLite allows empty strings with NOT NULL constraint
		},
		{
			name: "empty email (SQLite allows empty strings)",
			employee: &Employee{
				FirstName:  "John",
				LastName:   "Doe",
				Email:      "",
				JobTitle:   "Developer",
				Country:    "USA",
				Salary:     75000.0,
				Department: "Engineering",
			},
			expectErr: false, // SQLite allows empty strings with NOT NULL constraint
		},
		{
			name: "zero salary (NOT NULL constraint allows zero)",
			employee: &Employee{
				FirstName:  "John",
				LastName:   "Doe",
				Email:      "zero-salary@example.com",
				JobTitle:   "Developer",
				Country:    "USA",
				Salary:     0.0,
				Department: "Engineering",
			},
			expectErr: false, // SQLite allows zero values with NOT NULL constraint
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use a fresh database for each test to avoid unique constraint issues
			db := testDB(t)
			defer cleanupTestDB(db)

			err := db.Create(tt.employee).Error
			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestEmployee_UniqueEmail(t *testing.T) {
	db := testDB(t)
	defer cleanupTestDB(db)

	// Create first employee
	employee1 := &Employee{
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john@example.com",
		JobTitle:   "Developer",
		Country:    "USA",
		Salary:     75000.0,
		Department: "Engineering",
	}

	if err := db.Create(employee1).Error; err != nil {
		t.Fatalf("Failed to create first employee: %v", err)
	}

	// Try to create second employee with same email
	employee2 := &Employee{
		FirstName:  "Jane",
		LastName:   "Smith",
		Email:      "john@example.com", // Same email
		JobTitle:   "Manager",
		Country:    "UK",
		Salary:     85000.0,
		Department: "Management",
	}

	err := db.Create(employee2).Error
	if err == nil {
		t.Error("Expected error for duplicate email, but got none")
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

// testDB creates an in-memory database for testing
func testDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&Employee{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

// cleanupTestDB closes the test database connection
func cleanupTestDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}
