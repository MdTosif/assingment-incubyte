package services

import (
	"testing"

	"github.com/tofiquem/assingment/pkg/database"
	"github.com/tofiquem/assingment/pkg/models"
	"github.com/tofiquem/assingment/pkg/testutils"
)

func TestNewAnalyticsService(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAnalyticsService(testDB)
	if service == nil {
		t.Fatal("Expected analytics service to be created, got nil")
	}

	if service.db != testDB {
		t.Error("Expected service.db to be testDB")
	}
}

func TestAnalyticsService_GetSalaryByCountry(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAnalyticsService(testDB)

	// Create test employees with different countries
	employees := []models.Employee{
		{FirstName: "John", LastName: "Doe", Email: "john@example.com", JobTitle: "Developer", Country: "USA", Salary: 75000.0, Department: "Engineering"},
		{FirstName: "Jane", LastName: "Smith", Email: "jane@example.com", JobTitle: "Manager", Country: "UK", Salary: 85000.0, Department: "Management"},
		{FirstName: "Bob", LastName: "Johnson", Email: "bob@example.com", JobTitle: "Developer", Country: "USA", Salary: 80000.0, Department: "Engineering"},
		{FirstName: "Alice", LastName: "Williams", Email: "alice@example.com", JobTitle: "Designer", Country: "Canada", Salary: 70000.0, Department: "Design"},
		{FirstName: "Charlie", LastName: "Brown", Email: "charlie@example.com", JobTitle: "Developer", Country: "USA", Salary: 90000.0, Department: "Engineering"},
	}

	for _, emp := range employees {
		if err := testDB.Create(&emp).Error; err != nil {
			t.Fatalf("Failed to create test employee: %v", err)
		}
	}

	// Get salary by country
	stats, err := service.GetSalaryByCountry()
	if err != nil {
		t.Fatalf("Failed to get salary by country: %v", err)
	}

	if len(stats) != 3 {
		t.Errorf("Expected 3 countries, got %d", len(stats))
	}

	// Check USA stats (should be first due to highest average: 81500)
	usaFound := false
	for _, stat := range stats {
		if stat.Country == "USA" {
			usaFound = true
			if stat.Count != 3 {
				t.Errorf("Expected USA count 3, got %d", stat.Count)
			}
			if stat.Min != 75000.0 {
				t.Errorf("Expected USA min 75000.0, got %f", stat.Min)
			}
			if stat.Max != 90000.0 {
				t.Errorf("Expected USA max 90000.0, got %f", stat.Max)
			}
			// Average should be around 81666.67
			if stat.Average < 81666 || stat.Average > 81667 {
				t.Errorf("Expected USA average around 81666.67, got %f", stat.Average)
			}
		}
	}
	if !usaFound {
		t.Error("Expected USA to be in stats")
	}

	// Verify results are ordered by average descending
	if len(stats) >= 2 && stats[0].Average < stats[1].Average {
		t.Error("Expected results to be ordered by average salary descending")
	}
}

func TestAnalyticsService_GetSalaryByCountry_EmptyDatabase(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAnalyticsService(testDB)

	// Get salary by country with no employees
	stats, err := service.GetSalaryByCountry()
	if err != nil {
		t.Fatalf("Failed to get salary by country: %v", err)
	}

	if len(stats) != 0 {
		t.Errorf("Expected 0 countries with no employees, got %d", len(stats))
	}
}

func TestAnalyticsService_GetSalaryByJobTitleInCountry(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAnalyticsService(testDB)

	// Create test employees
	employees := []models.Employee{
		{FirstName: "John", LastName: "Doe", Email: "john@example.com", JobTitle: "Developer", Country: "USA", Salary: 75000.0, Department: "Engineering"},
		{FirstName: "Jane", LastName: "Smith", Email: "jane@example.com", JobTitle: "Manager", Country: "USA", Salary: 85000.0, Department: "Management"},
		{FirstName: "Bob", LastName: "Johnson", Email: "bob@example.com", JobTitle: "Developer", Country: "USA", Salary: 80000.0, Department: "Engineering"},
		{FirstName: "Alice", LastName: "Williams", Email: "alice@example.com", JobTitle: "Developer", Country: "UK", Salary: 70000.0, Department: "Engineering"},
	}

	for _, emp := range employees {
		if err := testDB.Create(&emp).Error; err != nil {
			t.Fatalf("Failed to create test employee: %v", err)
		}
	}

	// Get salary by job title in USA
	stats, err := service.GetSalaryByJobTitleInCountry("USA")
	if err != nil {
		t.Fatalf("Failed to get salary by job title: %v", err)
	}

	if len(stats) != 2 {
		t.Errorf("Expected 2 job titles in USA, got %d", len(stats))
	}

	// Find Developer stats
	devFound := false
	for _, stat := range stats {
		if stat.JobTitle == "Developer" {
			devFound = true
			if stat.Count != 2 {
				t.Errorf("Expected Developer count 2, got %d", stat.Count)
			}
			// Average should be 77500
			if stat.Average < 77499 || stat.Average > 77501 {
				t.Errorf("Expected Developer average 77500, got %f", stat.Average)
			}
		}
	}
	if !devFound {
		t.Error("Expected Developer to be in stats")
	}

	// Find Manager stats
	managerFound := false
	for _, stat := range stats {
		if stat.JobTitle == "Manager" {
			managerFound = true
			if stat.Count != 1 {
				t.Errorf("Expected Manager count 1, got %d", stat.Count)
			}
			if stat.Average != 85000.0 {
				t.Errorf("Expected Manager average 85000.0, got %f", stat.Average)
			}
		}
	}
	if !managerFound {
		t.Error("Expected Manager to be in stats")
	}
}

func TestAnalyticsService_GetSalaryByJobTitleInCountry_EmptyCountry(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAnalyticsService(testDB)

	// Try to get stats with empty country
	_, err := service.GetSalaryByJobTitleInCountry("")
	if err == nil {
		t.Error("Expected error for empty country parameter")
	}

	if err.Error() != "country parameter is required" {
		t.Errorf("Expected 'country parameter is required' error, got: %v", err)
	}
}

func TestAnalyticsService_GetSalaryByJobTitleInCountry_NonExistentCountry(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAnalyticsService(testDB)

	// Create test employees
	employees := []models.Employee{
		{FirstName: "John", LastName: "Doe", Email: "john@example.com", JobTitle: "Developer", Country: "USA", Salary: 75000.0, Department: "Engineering"},
	}

	for _, emp := range employees {
		if err := testDB.Create(&emp).Error; err != nil {
			t.Fatalf("Failed to create test employee: %v", err)
		}
	}

	// Get stats for non-existent country
	stats, err := service.GetSalaryByJobTitleInCountry("NonExistent")
	if err != nil {
		t.Fatalf("Failed to get salary by job title: %v", err)
	}

	if len(stats) != 0 {
		t.Errorf("Expected 0 job titles for non-existent country, got %d", len(stats))
	}
}

func TestAnalyticsService_GetDepartmentInsights(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAnalyticsService(testDB)

	// Create test employees
	employees := []models.Employee{
		{FirstName: "John", LastName: "Doe", Email: "john@example.com", JobTitle: "Developer", Country: "USA", Salary: 75000.0, Department: "Engineering"},
		{FirstName: "Jane", LastName: "Smith", Email: "jane@example.com", JobTitle: "Manager", Country: "USA", Salary: 85000.0, Department: "Management"},
		{FirstName: "Bob", LastName: "Johnson", Email: "bob@example.com", JobTitle: "Developer", Country: "UK", Salary: 80000.0, Department: "Engineering"},
		{FirstName: "Alice", LastName: "Williams", Email: "alice@example.com", JobTitle: "Designer", Country: "Canada", Salary: 70000.0, Department: "Design"},
		{FirstName: "Charlie", LastName: "Brown", Email: "charlie@example.com", JobTitle: "Developer", Country: "USA", Salary: 90000.0, Department: "Engineering"},
	}

	for _, emp := range employees {
		if err := testDB.Create(&emp).Error; err != nil {
			t.Fatalf("Failed to create test employee: %v", err)
		}
	}

	// Get department insights
	stats, err := service.GetDepartmentInsights()
	if err != nil {
		t.Fatalf("Failed to get department insights: %v", err)
	}

	if len(stats) != 3 {
		t.Errorf("Expected 3 departments, got %d", len(stats))
	}

	// Find Engineering stats
	engFound := false
	for _, stat := range stats {
		if stat.Department == "Engineering" {
			engFound = true
			if stat.Count != 3 {
				t.Errorf("Expected Engineering count 3, got %d", stat.Count)
			}
			if stat.Min != 75000.0 {
				t.Errorf("Expected Engineering min 75000.0, got %f", stat.Min)
			}
			if stat.Max != 90000.0 {
				t.Errorf("Expected Engineering max 90000.0, got %f", stat.Max)
			}
		}
	}
	if !engFound {
		t.Error("Expected Engineering to be in stats")
	}

	// Find Management stats
	mgmtFound := false
	for _, stat := range stats {
		if stat.Department == "Management" {
			mgmtFound = true
			if stat.Count != 1 {
				t.Errorf("Expected Management count 1, got %d", stat.Count)
			}
			if stat.Min != 85000.0 || stat.Max != 85000.0 {
				t.Errorf("Expected Management min and max to be 85000.0, got min=%f max=%f", stat.Min, stat.Max)
			}
		}
	}
	if !mgmtFound {
		t.Error("Expected Management to be in stats")
	}

	// Verify results are ordered by average descending
	if len(stats) >= 2 && stats[0].Average < stats[1].Average {
		t.Error("Expected results to be ordered by average salary descending")
	}
}

func TestAnalyticsService_GetDepartmentInsights_EmptyDatabase(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAnalyticsService(testDB)

	// Get department insights with no employees
	stats, err := service.GetDepartmentInsights()
	if err != nil {
		t.Fatalf("Failed to get department insights: %v", err)
	}

	if len(stats) != 0 {
		t.Errorf("Expected 0 departments with no employees, got %d", len(stats))
	}
}

func TestAnalyticsService_GetSalaryByCountry_OrderedByAverageDesc(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewAnalyticsService(testDB)

	// Create test employees with known averages
	// USA: avg = 80000, UK: avg = 70000, Canada: avg = 60000
	employees := []models.Employee{
		{FirstName: "John1", LastName: "Doe", Email: "john1@example.com", JobTitle: "Dev", Country: "USA", Salary: 80000.0, Department: "Eng"},
		{FirstName: "John2", LastName: "Smith", Email: "john2@example.com", JobTitle: "Dev", Country: "UK", Salary: 70000.0, Department: "Eng"},
		{FirstName: "John3", LastName: "Brown", Email: "john3@example.com", JobTitle: "Dev", Country: "Canada", Salary: 60000.0, Department: "Eng"},
	}

	for _, emp := range employees {
		if err := testDB.Create(&emp).Error; err != nil {
			t.Fatalf("Failed to create test employee: %v", err)
		}
	}

	// Get salary by country
	stats, err := service.GetSalaryByCountry()
	if err != nil {
		t.Fatalf("Failed to get salary by country: %v", err)
	}

	if len(stats) != 3 {
		t.Fatalf("Expected 3 countries, got %d", len(stats))
	}

	// Verify order: USA (80000), UK (70000), Canada (60000)
	expectedOrder := []string{"USA", "UK", "Canada"}
	for i, expected := range expectedOrder {
		if stats[i].Country != expected {
			t.Errorf("Expected country at position %d to be %s, got %s", i, expected, stats[i].Country)
		}
	}

	// Verify averages are in descending order
	for i := 0; i < len(stats)-1; i++ {
		if stats[i].Average < stats[i+1].Average {
			t.Errorf("Stats not ordered by average descending: %f < %f at positions %d and %d",
				stats[i].Average, stats[i+1].Average, i, i+1)
		}
	}
}
