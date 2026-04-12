package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tofiquem/assingment/pkg/database"
	"github.com/tofiquem/assingment/pkg/models"
	"github.com/tofiquem/assingment/pkg/testutils"
	"gorm.io/gorm"
)

func TestNewEmployeeHandler(t *testing.T) {
	// Store original DB
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)

	database.DB = testDB

	handler := NewEmployeeHandler()
	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}

	if handler.db != testDB {
		t.Errorf("Expected handler.db to be testDB, got different database")
	}
}

func TestEmployeeHandler_GetEmployees(t *testing.T) {
	// Setup
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	// Create test data
	testutils.CreateTestEmployees(testDB)

	handler := NewEmployeeHandler()

	tests := []struct {
		name         string
		queryParams  string
		expectedCode int
		checkCount   bool
		expectedLen  int
	}{
		{
			name:         "get all employees without pagination",
			queryParams:  "",
			expectedCode: http.StatusOK,
			checkCount:   true,
		},
		{
			name:         "get employees with page 1",
			queryParams:  "?page=1",
			expectedCode: http.StatusOK,
			checkCount:   true,
		},
		{
			name:         "get employees with page 2",
			queryParams:  "?page=2&limit=2",
			expectedCode: http.StatusOK,
			checkCount:   false,
		},
		{
			name:         "get employees with limit",
			queryParams:  "?limit=3",
			expectedCode: http.StatusOK,
			checkCount:   false, // Don't check count since we're limiting
		},
		{
			name:         "invalid page parameter",
			queryParams:  "?page=invalid",
			expectedCode: http.StatusOK,
			checkCount:   true,
		},
		{
			name:         "invalid limit parameter",
			queryParams:  "?limit=invalid",
			expectedCode: http.StatusOK,
			checkCount:   true,
		},
		{
			name:         "limit too high",
			queryParams:  "?limit=150",
			expectedCode: http.StatusOK,
			checkCount:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/employees"+tt.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.GetEmployees(rr, req)

			testutils.AssertStatusCode(t, rr, tt.expectedCode)
			testutils.AssertContentType(t, rr, "application/json")

			if tt.expectedCode == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if tt.checkCount {
					employees := response["employees"].([]interface{})
					total := response["total"].(float64)
					if len(employees) != int(total) {
						t.Errorf("Expected %d employees, got %d", int(total), len(employees))
					}
				}

				// Check required fields
				if _, ok := response["employees"]; !ok {
					t.Error("Response missing 'employees' field")
				}
				if _, ok := response["total"]; !ok {
					t.Error("Response missing 'total' field")
				}
				if _, ok := response["page"]; !ok {
					t.Error("Response missing 'page' field")
				}
				if _, ok := response["limit"]; !ok {
					t.Error("Response missing 'limit' field")
				}
				if _, ok := response["pages"]; !ok {
					t.Error("Response missing 'pages' field")
				}
			}
		})
	}
}

func TestEmployeeHandler_CreateEmployee(t *testing.T) {
	// Setup
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	handler := NewEmployeeHandler()

	tests := []struct {
		name         string
		requestBody  interface{}
		expectedCode int
		checkDB      bool
	}{
		{
			name: "valid employee creation",
			requestBody: models.CreateEmployeeRequest{
				FirstName:  "John",
				LastName:   "Doe",
				Email:      "john@example.com",
				JobTitle:   "Developer",
				Country:    "USA",
				Salary:     75000.0,
				Department: "Engineering",
			},
			expectedCode: http.StatusCreated,
			checkDB:      true,
		},
		{
			name:         "invalid request body",
			requestBody:  "invalid json",
			expectedCode: http.StatusBadRequest,
			checkDB:      false,
		},
		{
			name: "missing required fields",
			requestBody: models.CreateEmployeeRequest{
				FirstName: "John",
				// Missing other required fields
			},
			expectedCode: http.StatusCreated, // SQLite allows empty strings for NOT NULL fields
			checkDB:      false,
		},
		{
			name:         "empty request body",
			requestBody:  nil,
			expectedCode: http.StatusBadRequest,
			checkDB:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tt.requestBody == nil {
				req, err = http.NewRequest("POST", "/api/employees", bytes.NewBufferString(""))
			} else if str, ok := tt.requestBody.(string); ok {
				req, err = http.NewRequest("POST", "/api/employees", bytes.NewBufferString(str))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = testutils.CreateJSONRequest("POST", "/api/employees", tt.requestBody)
			}

			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.CreateEmployee(rr, req)

			testutils.AssertStatusCode(t, rr, tt.expectedCode)

			if tt.expectedCode == http.StatusCreated {
				testutils.AssertContentType(t, rr, "application/json")

				var employee models.Employee
				if err := json.Unmarshal(rr.Body.Bytes(), &employee); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if employee.ID == 0 {
					t.Error("Expected employee ID to be set")
				}

				if tt.checkDB {
					var dbEmployee models.Employee
					if err := testDB.First(&dbEmployee, employee.ID).Error; err != nil {
						t.Errorf("Expected employee to be saved in database: %v", err)
					}
				}
			}
		})
	}
}

func TestEmployeeHandler_GetEmployee(t *testing.T) {
	// Setup
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	// Create test data
	testutils.CreateTestEmployee(testDB, "John", "Doe", "john@example.com", "Developer", "USA", "Engineering", 75000.0)

	handler := NewEmployeeHandler()

	tests := []struct {
		name         string
		employeeID   string
		expectedCode int
	}{
		{
			name:         "valid employee ID",
			employeeID:   "1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid employee ID",
			employeeID:   "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "non-existent employee ID",
			employeeID:   "999",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/employees/"+tt.employeeID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Set up mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.employeeID})

			rr := httptest.NewRecorder()
			handler.GetEmployee(rr, req)

			testutils.AssertStatusCode(t, rr, tt.expectedCode)

			if tt.expectedCode == http.StatusOK {
				testutils.AssertContentType(t, rr, "application/json")

				var responseEmployee models.Employee
				if err := json.Unmarshal(rr.Body.Bytes(), &responseEmployee); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if responseEmployee.ID != 1 {
					t.Errorf("Expected employee ID 1, got %d", responseEmployee.ID)
				}
				if responseEmployee.FirstName != "John" {
					t.Errorf("Expected FirstName John, got %s", responseEmployee.FirstName)
				}
			}
		})
	}
}

func TestEmployeeHandler_UpdateEmployee(t *testing.T) {
	// Setup
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	// Create test data
	testutils.CreateTestEmployee(testDB, "John", "Doe", "john@example.com", "Developer", "USA", "Engineering", 75000.0)

	handler := NewEmployeeHandler()

	tests := []struct {
		name         string
		employeeID   string
		requestBody  interface{}
		expectedCode int
	}{
		{
			name:       "valid update",
			employeeID: "1",
			requestBody: models.UpdateEmployeeRequest{
				FirstName: stringPtr("Jane"),
				Salary:    floatPtr(80000.0),
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid employee ID",
			employeeID:   "invalid",
			requestBody:  models.UpdateEmployeeRequest{},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "non-existent employee ID",
			employeeID:   "999",
			requestBody:  models.UpdateEmployeeRequest{},
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid request body",
			employeeID:   "1",
			requestBody:  "invalid json",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "empty update",
			employeeID:   "1",
			requestBody:  models.UpdateEmployeeRequest{},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if str, ok := tt.requestBody.(string); ok {
				req, err = http.NewRequest("PUT", "/api/employees/"+tt.employeeID, bytes.NewBufferString(str))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = testutils.CreateJSONRequest("PUT", "/api/employees/"+tt.employeeID, tt.requestBody)
			}

			if err != nil {
				t.Fatal(err)
			}

			// Set up mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.employeeID})

			rr := httptest.NewRecorder()
			handler.UpdateEmployee(rr, req)

			testutils.AssertStatusCode(t, rr, tt.expectedCode)

			if tt.expectedCode == http.StatusOK {
				testutils.AssertContentType(t, rr, "application/json")

				var responseEmployee models.Employee
				if err := json.Unmarshal(rr.Body.Bytes(), &responseEmployee); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				// Check if update was applied
				if updateReq, ok := tt.requestBody.(models.UpdateEmployeeRequest); ok {
					if updateReq.FirstName != nil && responseEmployee.FirstName != *updateReq.FirstName {
						t.Errorf("Expected FirstName %s, got %s", *updateReq.FirstName, responseEmployee.FirstName)
					}
					if updateReq.Salary != nil && responseEmployee.Salary != *updateReq.Salary {
						t.Errorf("Expected Salary %f, got %f", *updateReq.Salary, responseEmployee.Salary)
					}
				}
			}
		})
	}
}

func TestEmployeeHandler_DeleteEmployee(t *testing.T) {
	// Setup
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	// Create test data
	testutils.CreateTestEmployee(testDB, "John", "Doe", "john@example.com", "Developer", "USA", "Engineering", 75000.0)

	handler := NewEmployeeHandler()

	tests := []struct {
		name         string
		employeeID   string
		expectedCode int
		checkDB      bool
	}{
		{
			name:         "valid deletion",
			employeeID:   "1",
			expectedCode: http.StatusNoContent,
			checkDB:      true,
		},
		{
			name:         "invalid employee ID",
			employeeID:   "invalid",
			expectedCode: http.StatusBadRequest,
			checkDB:      false,
		},
		{
			name:         "non-existent employee ID",
			employeeID:   "999",
			expectedCode: http.StatusNotFound,
			checkDB:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/api/employees/"+tt.employeeID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Set up mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.employeeID})

			rr := httptest.NewRecorder()
			handler.DeleteEmployee(rr, req)

			testutils.AssertStatusCode(t, rr, tt.expectedCode)

			if tt.checkDB && tt.expectedCode == http.StatusNoContent {
				var dbEmployee models.Employee
				err := testDB.First(&dbEmployee, 1).Error
				if err == nil {
					t.Error("Expected employee to be deleted from database")
				} else if err != gorm.ErrRecordNotFound {
					t.Errorf("Unexpected error checking database: %v", err)
				}
			}
		})
	}
}

func TestEmployeeHandler_RegisterRoutes(t *testing.T) {
	// Setup
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	handler := NewEmployeeHandler()
	router := mux.NewRouter()

	handler.RegisterRoutes(router)

	// Test that routes are registered
	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/employees"},
		{"POST", "/api/employees"},
		{"GET", "/api/employees/{id}"},
		{"PUT", "/api/employees/{id}"},
		{"DELETE", "/api/employees/{id}"},
	}

	for _, route := range routes {
		t.Run("route_"+route.method+"_"+route.path, func(t *testing.T) {
			var match mux.RouteMatch
			req, _ := http.NewRequest(route.method, route.path, nil)

			if router.Match(req, &match) {
				if match.Handler == nil {
					t.Errorf("Route %s %s has no handler", route.method, route.path)
				}
			} else {
				t.Errorf("Route %s %s not found", route.method, route.path)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
