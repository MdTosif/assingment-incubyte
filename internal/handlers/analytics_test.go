package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/tofiquem/assingment/internal/database"
	"github.com/tofiquem/assingment/internal/models"
	"github.com/tofiquem/assingment/internal/testutils"
)

func TestNewAnalyticsHandler(t *testing.T) {
	// Store original DB
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)

	database.DB = testDB

	handler := NewAnalyticsHandler()
	if handler == nil {
		t.Fatal("Expected handler to be created, got nil")
	}

	if handler.db != testDB {
		t.Errorf("Expected handler.db to be testDB, got different database")
	}
}

func TestAnalyticsHandler_GetSalaryByCountry(t *testing.T) {
	// Setup
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	// Create test data with different countries and salaries
	testutils.CreateTestEmployees(testDB)

	handler := NewAnalyticsHandler()

	req, err := http.NewRequest("GET", "/api/analytics/salary/by-country", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetSalaryByCountry(rr, req)

	testutils.AssertStatusCode(t, rr, http.StatusOK)
	testutils.AssertContentType(t, rr, "application/json")

	var stats []CountrySalaryStats
	if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(stats) == 0 {
		t.Error("Expected at least one country in response")
	}

	// Verify the structure of the response
	for _, stat := range stats {
		if stat.Country == "" {
			t.Error("Country field should not be empty")
		}
		if stat.Count <= 0 {
			t.Error("Count should be greater than 0")
		}
		if stat.Average <= 0 {
			t.Error("Average salary should be greater than 0")
		}
		if stat.Min <= 0 {
			t.Error("Min salary should be greater than 0")
		}
		if stat.Max <= 0 {
			t.Error("Max salary should be greater than 0")
		}
		if stat.Min > stat.Max {
			t.Error("Min salary should not be greater than max salary")
		}
		if stat.Average < stat.Min || stat.Average > stat.Max {
			t.Error("Average salary should be between min and max")
		}
	}

	// Check that results are ordered by average descending
	for i := 1; i < len(stats); i++ {
		if stats[i-1].Average < stats[i].Average {
			t.Errorf("Results should be ordered by average descending, but found %f before %f", stats[i-1].Average, stats[i].Average)
		}
	}
}

func TestAnalyticsHandler_GetSalaryByJobTitleInCountry(t *testing.T) {
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

	handler := NewAnalyticsHandler()

	tests := []struct {
		name         string
		country      string
		expectedCode int
		expectedLen  int
	}{
		{
			name:         "valid country with employees",
			country:      "USA",
			expectedCode: http.StatusOK,
			expectedLen:  1, // Only Developer job title in USA (3 employees with same title)
		},
		{
			name:         "valid country with one employee",
			country:      "UK",
			expectedCode: http.StatusOK,
			expectedLen:  1,
		},
		{
			name:         "valid country with one employee",
			country:      "Canada",
			expectedCode: http.StatusOK,
			expectedLen:  1,
		},
		{
			name:         "country with no employees",
			country:      "Germany",
			expectedCode: http.StatusOK,
			expectedLen:  0,
		},
		{
			name:         "empty country parameter",
			country:      "",
			expectedCode: http.StatusBadRequest,
			expectedLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/analytics/salary/by-job-title/"+tt.country, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Set up mux vars
			req = mux.SetURLVars(req, map[string]string{"country": tt.country})

			rr := httptest.NewRecorder()
			handler.GetSalaryByJobTitleInCountry(rr, req)

			testutils.AssertStatusCode(t, rr, tt.expectedCode)

			if tt.expectedCode == http.StatusOK {
				testutils.AssertContentType(t, rr, "application/json")

				var stats []JobTitleSalaryStats
				if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}

				if len(stats) != tt.expectedLen {
					t.Errorf("Expected %d job titles, got %d", tt.expectedLen, len(stats))
				}

				// Verify the structure of the response
				for _, stat := range stats {
					if stat.JobTitle == "" {
						t.Error("JobTitle field should not be empty")
					}
					if stat.Count <= 0 {
						t.Error("Count should be greater than 0")
					}
					if stat.Average <= 0 {
						t.Error("Average salary should be greater than 0")
					}
				}

				// Check that results are ordered by average descending
				for i := 1; i < len(stats); i++ {
					if stats[i-1].Average < stats[i].Average {
						t.Errorf("Results should be ordered by average descending, but found %f before %f", stats[i-1].Average, stats[i].Average)
					}
				}
			}
		})
	}
}

func TestAnalyticsHandler_GetDepartmentInsights(t *testing.T) {
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

	handler := NewAnalyticsHandler()

	req, err := http.NewRequest("GET", "/api/analytics/salary/department-insights", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetDepartmentInsights(rr, req)

	testutils.AssertStatusCode(t, rr, http.StatusOK)
	testutils.AssertContentType(t, rr, "application/json")

	var stats []DepartmentSalaryStats
	if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(stats) == 0 {
		t.Error("Expected at least one department in response")
	}

	// Verify the structure of the response
	for _, stat := range stats {
		if stat.Department == "" {
			t.Error("Department field should not be empty")
		}
		if stat.Count <= 0 {
			t.Error("Count should be greater than 0")
		}
		if stat.Average <= 0 {
			t.Error("Average salary should be greater than 0")
		}
		if stat.Min <= 0 {
			t.Error("Min salary should be greater than 0")
		}
		if stat.Max <= 0 {
			t.Error("Max salary should be greater than 0")
		}
		if stat.Min > stat.Max {
			t.Error("Min salary should not be greater than max salary")
		}
		if stat.Average < stat.Min || stat.Average > stat.Max {
			t.Error("Average salary should be between min and max")
		}
	}

	// Check that results are ordered by average descending
	for i := 1; i < len(stats); i++ {
		if stats[i-1].Average < stats[i].Average {
			t.Errorf("Results should be ordered by average descending, but found %f before %f", stats[i-1].Average, stats[i].Average)
		}
	}
}

func TestAnalyticsHandler_GetSalaryByCountry_EmptyDatabase(t *testing.T) {
	// Setup
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	handler := NewAnalyticsHandler()

	req, err := http.NewRequest("GET", "/api/analytics/salary/by-country", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetSalaryByCountry(rr, req)

	testutils.AssertStatusCode(t, rr, http.StatusOK)
	testutils.AssertContentType(t, rr, "application/json")

	var stats []CountrySalaryStats
	if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(stats) != 0 {
		t.Errorf("Expected empty result for empty database, got %d countries", len(stats))
	}
}

func TestAnalyticsHandler_GetDepartmentInsights_EmptyDatabase(t *testing.T) {
	// Setup
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	handler := NewAnalyticsHandler()

	req, err := http.NewRequest("GET", "/api/analytics/salary/department-insights", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetDepartmentInsights(rr, req)

	testutils.AssertStatusCode(t, rr, http.StatusOK)
	testutils.AssertContentType(t, rr, "application/json")

	var stats []DepartmentSalaryStats
	if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(stats) != 0 {
		t.Errorf("Expected empty result for empty database, got %d departments", len(stats))
	}
}

func TestAnalyticsHandler_RegisterRoutes(t *testing.T) {
	// Setup
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	handler := NewAnalyticsHandler()
	router := mux.NewRouter()

	handler.RegisterRoutes(router)

	// Test that routes are registered
	routes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/analytics/salary/by-country"},
		{"GET", "/api/analytics/salary/by-job-title/{country}"},
		{"GET", "/api/analytics/salary/department-insights"},
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

func TestAnalyticsHandler_DataConsistency(t *testing.T) {
	// Setup
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	// Create test data with known values
	employees := []models.Employee{
		{FirstName: "John", LastName: "Doe", Email: "john@example.com", JobTitle: "Developer", Country: "USA", Salary: 75000.0, Department: "Engineering"},
		{FirstName: "Jane", LastName: "Smith", Email: "jane@example.com", JobTitle: "Developer", Country: "USA", Salary: 85000.0, Department: "Engineering"},
		{FirstName: "Bob", LastName: "Johnson", Email: "bob@example.com", JobTitle: "Manager", Country: "USA", Salary: 95000.0, Department: "Management"},
	}

	for _, emp := range employees {
		if err := testDB.Create(&emp).Error; err != nil {
			t.Fatalf("Failed to create test employee: %v", err)
		}
	}

	handler := NewAnalyticsHandler()

	// Test country statistics
	req, err := http.NewRequest("GET", "/api/analytics/salary/by-country", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.GetSalaryByCountry(rr, req)

	testutils.AssertStatusCode(t, rr, http.StatusOK)

	var countryStats []CountrySalaryStats
	if err := json.Unmarshal(rr.Body.Bytes(), &countryStats); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Should have one country (USA) with 3 employees
	if len(countryStats) != 1 {
		t.Errorf("Expected 1 country, got %d", len(countryStats))
	}

	if countryStats[0].Country != "USA" {
		t.Errorf("Expected country USA, got %s", countryStats[0].Country)
	}

	if countryStats[0].Count != 3 {
		t.Errorf("Expected 3 employees in USA, got %d", countryStats[0].Count)
	}

	expectedAvg := (75000.0 + 85000.0 + 95000.0) / 3
	if countryStats[0].Average != expectedAvg {
		t.Errorf("Expected average salary %f, got %f", expectedAvg, countryStats[0].Average)
	}

	if countryStats[0].Min != 75000.0 {
		t.Errorf("Expected min salary 75000.0, got %f", countryStats[0].Min)
	}

	if countryStats[0].Max != 95000.0 {
		t.Errorf("Expected max salary 95000.0, got %f", countryStats[0].Max)
	}

	// Test department statistics
	req2, err := http.NewRequest("GET", "/api/analytics/salary/department-insights", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr2 := httptest.NewRecorder()
	handler.GetDepartmentInsights(rr2, req2)

	testutils.AssertStatusCode(t, rr2, http.StatusOK)

	var deptStats []DepartmentSalaryStats
	if err := json.Unmarshal(rr2.Body.Bytes(), &deptStats); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Should have two departments
	if len(deptStats) != 2 {
		t.Errorf("Expected 2 departments, got %d", len(deptStats))
	}

	// Engineering department should have 2 employees with avg 80000
	var engineeringDept *DepartmentSalaryStats
	var managementDept *DepartmentSalaryStats

	for i := range deptStats {
		if deptStats[i].Department == "Engineering" {
			engineeringDept = &deptStats[i]
		} else if deptStats[i].Department == "Management" {
			managementDept = &deptStats[i]
		}
	}

	if engineeringDept == nil {
		t.Error("Engineering department not found")
	} else {
		if engineeringDept.Count != 2 {
			t.Errorf("Expected 2 employees in Engineering, got %d", engineeringDept.Count)
		}
		if engineeringDept.Average != 80000.0 {
			t.Errorf("Expected Engineering avg 80000.0, got %f", engineeringDept.Average)
		}
	}

	if managementDept == nil {
		t.Error("Management department not found")
	} else {
		if managementDept.Count != 1 {
			t.Errorf("Expected 1 employee in Management, got %d", managementDept.Count)
		}
		if managementDept.Average != 95000.0 {
			t.Errorf("Expected Management avg 95000.0, got %f", managementDept.Average)
		}
	}
}
