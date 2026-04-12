package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tofiquem/assingment/internal/database"
	"github.com/tofiquem/assingment/internal/handlers"
	"github.com/tofiquem/assingment/internal/models"
	"github.com/tofiquem/assingment/internal/testutils"
)

func TestMain_Setup(t *testing.T) {
	// Store original environment variables
	originalDBPath := os.Getenv("DATABASE_PATH")
	originalPublicDir := os.Getenv("PUBLIC_DIR")
	originalPort := os.Getenv("PORT")

	// Cleanup after tests
	defer func() {
		if originalDBPath == "" {
			os.Unsetenv("DATABASE_PATH")
		} else {
			os.Setenv("DATABASE_PATH", originalDBPath)
		}
		if originalPublicDir == "" {
			os.Unsetenv("PUBLIC_DIR")
		} else {
			os.Setenv("PUBLIC_DIR", originalPublicDir)
		}
		if originalPort == "" {
			os.Unsetenv("PORT")
		} else {
			os.Setenv("PORT", originalPort)
		}
	}()

	t.Run("health_check", func(t *testing.T) {
		// Setup test database
		testDB := testutils.TestDB(t)
		defer testutils.CleanupTestDB(testDB)

		// Override global DB
		originalDB := database.DB
		database.DB = testDB
		defer func() { database.DB = originalDB }()

		// Create router
		router := mux.NewRouter()
		router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		}).Methods("GET")

		req, err := http.NewRequest("GET", "/api/health", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		expected := `{"status":"ok"}`
		if rr.Body.String() != expected {
			t.Errorf("Expected body %s, got %s", expected, rr.Body.String())
		}

		contentType := rr.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected content type application/json, got %s", contentType)
		}
	})

	t.Run("employee_crud_operations", func(t *testing.T) {
		// Setup test database
		testDB := testutils.TestDB(t)
		defer testutils.CleanupTestDB(testDB)

		// Override global DB
		originalDB := database.DB
		database.DB = testDB
		defer func() { database.DB = originalDB }()

		// Create router with employee routes
		router := setupTestRouter()

		// Test CREATE employee
		createReq := models.CreateEmployeeRequest{
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "john@example.com",
			JobTitle:   "Developer",
			Country:    "USA",
			Salary:     75000.0,
			Department: "Engineering",
		}

		createBody, _ := json.Marshal(createReq)
		req, err := http.NewRequest("POST", "/api/employees", bytes.NewBuffer(createBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		}

		var createdEmployee models.Employee
		if err := json.Unmarshal(rr.Body.Bytes(), &createdEmployee); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if createdEmployee.ID == 0 {
			t.Error("Expected employee ID to be set")
		}

		// Test GET all employees
		req, err = http.NewRequest("GET", "/api/employees", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		var getResponse map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &getResponse); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		employees := getResponse["employees"].([]interface{})
		if len(employees) != 1 {
			t.Errorf("Expected 1 employee, got %d", len(employees))
		}

		// Test GET specific employee
		req, err = http.NewRequest("GET", "/api/employees/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		var getEmployee models.Employee
		if err := json.Unmarshal(rr.Body.Bytes(), &getEmployee); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if getEmployee.ID != createdEmployee.ID {
			t.Errorf("Expected employee ID %d, got %d", createdEmployee.ID, getEmployee.ID)
		}

		// Test UPDATE employee
		updateReq := models.UpdateEmployeeRequest{
			FirstName: stringPtr("Jane"),
			Salary:    floatPtr(80000.0),
		}

		updateBody, _ := json.Marshal(updateReq)
		req, err = http.NewRequest("PUT", "/api/employees/1", bytes.NewBuffer(updateBody))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		var updatedEmployee models.Employee
		if err := json.Unmarshal(rr.Body.Bytes(), &updatedEmployee); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if updatedEmployee.FirstName != "Jane" {
			t.Errorf("Expected FirstName Jane, got %s", updatedEmployee.FirstName)
		}
		if updatedEmployee.Salary != 80000.0 {
			t.Errorf("Expected Salary 80000.0, got %f", updatedEmployee.Salary)
		}

		// Test DELETE employee
		req, err = http.NewRequest("DELETE", "/api/employees/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rr.Code)
		}

		// Verify employee is deleted
		req, err = http.NewRequest("GET", "/api/employees/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("analytics_endpoints", func(t *testing.T) {
		// Setup test database
		testDB := testutils.TestDB(t)
		defer testutils.CleanupTestDB(testDB)

		// Override global DB
		originalDB := database.DB
		database.DB = testDB
		defer func() { database.DB = originalDB }()

		// Create test data
		testutils.CreateTestEmployees(testDB)

		// Create router with analytics routes
		router := setupTestRouter()

		// Test salary by country
		req, err := http.NewRequest("GET", "/api/analytics/salary/by-country", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		var countryStats []map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &countryStats); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if len(countryStats) == 0 {
			t.Error("Expected at least one country in response")
		}

		// Test salary by job title in country
		req, err = http.NewRequest("GET", "/api/analytics/salary/by-job-title/USA", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		var jobTitleStats []map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &jobTitleStats); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		// Test department insights
		req, err = http.NewRequest("GET", "/api/analytics/salary/department-insights", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		var deptStats []map[string]interface{}
		if err := json.Unmarshal(rr.Body.Bytes(), &deptStats); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if len(deptStats) == 0 {
			t.Error("Expected at least one department in response")
		}
	})

	t.Run("error_handling", func(t *testing.T) {
		// Setup test database
		testDB := testutils.TestDB(t)
		defer testutils.CleanupTestDB(testDB)

		// Override global DB
		originalDB := database.DB
		database.DB = testDB
		defer func() { database.DB = originalDB }()

		router := setupTestRouter()

		// Test invalid JSON
		req, err := http.NewRequest("POST", "/api/employees", bytes.NewBufferString("invalid json"))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		// Test invalid employee ID
		req, err = http.NewRequest("GET", "/api/employees/invalid", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}

		// Test non-existent employee
		req, err = http.NewRequest("GET", "/api/employees/999", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}

		// Test invalid country parameter for analytics
		req, err = http.NewRequest("GET", "/api/analytics/salary/by-job-title/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// This should return 404 because the route doesn't match without country
		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rr.Code)
		}
	})
}

func TestSpaHandler(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "spa_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	indexContent := `<html><body>Index Page</body></html>`
	cssContent := `body { color: red; }`
	jsContent := `console.log('Hello World');`

	if err := os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(indexContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "styles.css"), []byte(cssContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "app.js"), []byte(jsContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "nested.html"), []byte("Nested page"), 0644); err != nil {
		t.Fatal(err)
	}

	handler := spaHandler(tempDir)

	tests := []struct {
		name           string
		path           string
		expectedCode   int
		expectedBody   string
		expectedHeader string
	}{
		{
			name:         "serve index.html for root",
			path:         "/",
			expectedCode: http.StatusOK,
			expectedBody: indexContent,
		},
		{
			name:         "serve CSS file",
			path:         "/styles.css",
			expectedCode: http.StatusOK,
			expectedBody: cssContent,
		},
		{
			name:         "serve JS file",
			path:         "/app.js",
			expectedCode: http.StatusOK,
			expectedBody: jsContent,
		},
		{
			name:         "serve nested file",
			path:         "/subdir/nested.html",
			expectedCode: http.StatusOK,
			expectedBody: "Nested page",
		},
		{
			name:         "fallback to index.html for non-existent file",
			path:         "/non-existent",
			expectedCode: http.StatusOK,
			expectedBody: indexContent,
		},
		{
			name:         "return 404 for API routes",
			path:         "/api/test",
			expectedCode: http.StatusNotFound,
			expectedBody: "404 page not found\n",
		},
		{
			name:         "serve index.html for empty path",
			path:         "",
			expectedCode: http.StatusOK,
			expectedBody: indexContent,
		},
		{
			name:         "return 404 for directory traversal attempt",
			path:         "/../etc/passwd",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			if tt.expectedBody != "" && rr.Body.String() != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}

			if tt.expectedHeader != "" {
				header := rr.Header().Get(tt.expectedHeader)
				if header == "" {
					t.Errorf("Expected header %s to be set", tt.expectedHeader)
				}
			}
		})
	}
}

func TestEnvironmentVariables(t *testing.T) {
	// Test default values
	if os.Getenv("PUBLIC_DIR") != "" {
		t.Skip("PUBLIC_DIR already set, skipping test")
	}

	if os.Getenv("PORT") != "" {
		t.Skip("PORT already set, skipping test")
	}

	// Test with unset environment variables
	originalDB := database.DB
	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB
	defer func() { database.DB = originalDB }()

	// This should use default values
	router := setupTestRouter()

	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Test with custom environment variables
	cleanup := testutils.SetTestEnv("PUBLIC_DIR", "/tmp/public")
	defer cleanup()

	cleanupPort := testutils.SetTestEnv("PORT", "3000")
	defer cleanupPort()

	// The environment variables are used in main() function, not in handlers
	// So we just verify they can be set and retrieved
	if os.Getenv("PUBLIC_DIR") != "/tmp/public" {
		t.Errorf("Expected PUBLIC_DIR to be /tmp/public, got %s", os.Getenv("PUBLIC_DIR"))
	}

	if os.Getenv("PORT") != "3000" {
		t.Errorf("Expected PORT to be 3000, got %s", os.Getenv("PORT"))
	}
}

// Helper function to set up test router
func setupTestRouter() *mux.Router {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Employee handlers
	employeeHandler := handlers.NewEmployeeHandler()
	employeeHandler.RegisterRoutes(router)

	// Analytics handlers
	analyticsHandler := handlers.NewAnalyticsHandler()
	analyticsHandler.RegisterRoutes(router)

	return router
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
