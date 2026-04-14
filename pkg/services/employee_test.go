package services

import (
	"testing"

	"github.com/tofiquem/assingment/pkg/database"
	"github.com/tofiquem/assingment/pkg/models"
	"github.com/tofiquem/assingment/pkg/testutils"
	"gorm.io/gorm"
)

func TestNewEmployeeService(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewEmployeeService(testDB)
	if service == nil {
		t.Fatal("Expected employee service to be created, got nil")
	}

	if service.db != testDB {
		t.Error("Expected service.db to be testDB")
	}
}

func TestEmployeeService_CreateEmployee(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewEmployeeService(testDB)

	// Test creating employee
	req := &models.CreateEmployeeRequest{
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@example.com",
		JobTitle:   "Developer",
		Country:    "USA",
		Salary:     75000.0,
		Department: "Engineering",
	}

	employee, err := service.CreateEmployee(req)
	if err != nil {
		t.Fatalf("Failed to create employee: %v", err)
	}

	if employee.ID == 0 {
		t.Error("Expected employee ID to be set")
	}

	if employee.FirstName != "John" {
		t.Errorf("Expected first name John, got %s", employee.FirstName)
	}

	if employee.Email != "john.doe@example.com" {
		t.Errorf("Expected email john.doe@example.com, got %s", employee.Email)
	}

	if employee.Salary != 75000.0 {
		t.Errorf("Expected salary 75000.0, got %f", employee.Salary)
	}
}

func TestEmployeeService_CreateEmployee_DuplicateEmail(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewEmployeeService(testDB)

	// Create first employee
	req1 := &models.CreateEmployeeRequest{
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "duplicate@example.com",
		JobTitle:   "Developer",
		Country:    "USA",
		Salary:     75000.0,
		Department: "Engineering",
	}

	_, err := service.CreateEmployee(req1)
	if err != nil {
		t.Fatalf("Failed to create first employee: %v", err)
	}

	// Try to create employee with same email
	req2 := &models.CreateEmployeeRequest{
		FirstName:  "Jane",
		LastName:   "Smith",
		Email:      "duplicate@example.com",
		JobTitle:   "Manager",
		Country:    "UK",
		Salary:     85000.0,
		Department: "Management",
	}

	_, err = service.CreateEmployee(req2)
	if err == nil {
		t.Error("Expected error for duplicate email")
	}
}

func TestEmployeeService_GetEmployeeByID(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewEmployeeService(testDB)

	// Create test employee
	req := &models.CreateEmployeeRequest{
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john@example.com",
		JobTitle:   "Developer",
		Country:    "USA",
		Salary:     75000.0,
		Department: "Engineering",
	}

	created, err := service.CreateEmployee(req)
	if err != nil {
		t.Fatalf("Failed to create test employee: %v", err)
	}

	// Get employee by ID
	found, err := service.GetEmployeeByID(created.ID)
	if err != nil {
		t.Fatalf("Failed to get employee by ID: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("Expected employee ID %d, got %d", created.ID, found.ID)
	}

	if found.Email != "john@example.com" {
		t.Errorf("Expected email john@example.com, got %s", found.Email)
	}
}

func TestEmployeeService_GetEmployeeByID_NonExistent(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewEmployeeService(testDB)

	// Try to get non-existent employee
	_, err := service.GetEmployeeByID(999)
	if err == nil {
		t.Error("Expected error for non-existent employee")
	}

	if err.Error() != "employee not found" {
		t.Errorf("Expected 'employee not found' error, got: %v", err)
	}
}

func TestEmployeeService_ListEmployees(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewEmployeeService(testDB)

	// Create test employees
	employees := []models.CreateEmployeeRequest{
		{FirstName: "John", LastName: "Doe", Email: "john@example.com", JobTitle: "Developer", Country: "USA", Salary: 75000.0, Department: "Engineering"},
		{FirstName: "Jane", LastName: "Smith", Email: "jane@example.com", JobTitle: "Manager", Country: "UK", Salary: 85000.0, Department: "Management"},
		{FirstName: "Bob", LastName: "Johnson", Email: "bob@example.com", JobTitle: "Developer", Country: "USA", Salary: 80000.0, Department: "Engineering"},
	}

	for _, emp := range employees {
		if _, err := service.CreateEmployee(&emp); err != nil {
			t.Fatalf("Failed to create test employee: %v", err)
		}
	}

	// List employees
	response, err := service.ListEmployees(1, 10, "")
	if err != nil {
		t.Fatalf("Failed to list employees: %v", err)
	}

	if response.Total != 3 {
		t.Errorf("Expected total 3, got %d", response.Total)
	}

	if len(response.Employees) != 3 {
		t.Errorf("Expected 3 employees, got %d", len(response.Employees))
	}

	if response.Page != 1 {
		t.Errorf("Expected page 1, got %d", response.Page)
	}

	if response.Limit != 10 {
		t.Errorf("Expected limit 10, got %d", response.Limit)
	}
}

func TestEmployeeService_ListEmployees_WithSearch(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewEmployeeService(testDB)

	// Create test employees
	employees := []models.CreateEmployeeRequest{
		{FirstName: "John", LastName: "Doe", Email: "john@example.com", JobTitle: "Developer", Country: "USA", Salary: 75000.0, Department: "Engineering"},
		{FirstName: "Jane", LastName: "Smith", Email: "jane@example.com", JobTitle: "Manager", Country: "UK", Salary: 85000.0, Department: "Management"},
		{FirstName: "Bob", LastName: "Johnson", Email: "bob@example.com", JobTitle: "Designer", Country: "USA", Salary: 70000.0, Department: "Design"},
	}

	for _, emp := range employees {
		if _, err := service.CreateEmployee(&emp); err != nil {
			t.Fatalf("Failed to create test employee: %v", err)
		}
	}

	// Search by first name (use unique term that won't match other names)
	response, err := service.ListEmployees(1, 10, "Jane")
	if err != nil {
		t.Fatalf("Failed to list employees with search: %v", err)
	}

	if response.Total != 1 {
		t.Errorf("Expected total 1 for 'Jane' search, got %d", response.Total)
	}

	// Search by job title
	response, err = service.ListEmployees(1, 10, "Developer")
	if err != nil {
		t.Fatalf("Failed to list employees with search: %v", err)
	}

	if response.Total != 1 {
		t.Errorf("Expected total 1 for 'Developer' search, got %d", response.Total)
	}

	// Search by country
	response, err = service.ListEmployees(1, 10, "USA")
	if err != nil {
		t.Fatalf("Failed to list employees with search: %v", err)
	}

	if response.Total != 2 {
		t.Errorf("Expected total 2 for 'USA' search, got %d", response.Total)
	}
}

func TestEmployeeService_ListEmployees_Pagination(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewEmployeeService(testDB)

	// Create 5 test employees
	for i := 0; i < 5; i++ {
		req := &models.CreateEmployeeRequest{
			FirstName:  "Employee",
			LastName:   string(rune('A' + i)),
			Email:      "employee" + string(rune('0'+i)) + "@example.com",
			JobTitle:   "Developer",
			Country:    "USA",
			Salary:     75000.0,
			Department: "Engineering",
		}
		if _, err := service.CreateEmployee(req); err != nil {
			t.Fatalf("Failed to create test employee: %v", err)
		}
	}

	// Test first page with limit 2
	response, err := service.ListEmployees(1, 2, "")
	if err != nil {
		t.Fatalf("Failed to list employees: %v", err)
	}

	if response.Total != 5 {
		t.Errorf("Expected total 5, got %d", response.Total)
	}

	if len(response.Employees) != 2 {
		t.Errorf("Expected 2 employees on page 1, got %d", len(response.Employees))
	}

	if response.Pages != 3 {
		t.Errorf("Expected 3 pages, got %d", response.Pages)
	}

	// Test second page
	response, err = service.ListEmployees(2, 2, "")
	if err != nil {
		t.Fatalf("Failed to list employees: %v", err)
	}

	if len(response.Employees) != 2 {
		t.Errorf("Expected 2 employees on page 2, got %d", len(response.Employees))
	}

	// Test third page (should have 1 employee)
	response, err = service.ListEmployees(3, 2, "")
	if err != nil {
		t.Fatalf("Failed to list employees: %v", err)
	}

	if len(response.Employees) != 1 {
		t.Errorf("Expected 1 employee on page 3, got %d", len(response.Employees))
	}
}

func TestEmployeeService_UpdateEmployee(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewEmployeeService(testDB)

	// Create test employee
	req := &models.CreateEmployeeRequest{
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john@example.com",
		JobTitle:   "Developer",
		Country:    "USA",
		Salary:     75000.0,
		Department: "Engineering",
	}

	created, err := service.CreateEmployee(req)
	if err != nil {
		t.Fatalf("Failed to create test employee: %v", err)
	}

	// Update employee
	newFirstName := "Jonathan"
	newSalary := 80000.0
	updateReq := &models.UpdateEmployeeRequest{
		FirstName: &newFirstName,
		Salary:    &newSalary,
	}

	updated, err := service.UpdateEmployee(created.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update employee: %v", err)
	}

	if updated.FirstName != "Jonathan" {
		t.Errorf("Expected first name Jonathan, got %s", updated.FirstName)
	}

	if updated.Salary != 80000.0 {
		t.Errorf("Expected salary 80000.0, got %f", updated.Salary)
	}

	// Verify unchanged fields
	if updated.LastName != "Doe" {
		t.Errorf("Expected last name to remain Doe, got %s", updated.LastName)
	}

	if updated.Email != "john@example.com" {
		t.Errorf("Expected email to remain john@example.com, got %s", updated.Email)
	}
}

func TestEmployeeService_UpdateEmployee_NonExistent(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewEmployeeService(testDB)

	// Try to update non-existent employee
	newFirstName := "Updated"
	updateReq := &models.UpdateEmployeeRequest{
		FirstName: &newFirstName,
	}

	_, err := service.UpdateEmployee(999, updateReq)
	if err == nil {
		t.Error("Expected error for non-existent employee")
	}

	if err.Error() != "employee not found" {
		t.Errorf("Expected 'employee not found' error, got: %v", err)
	}
}

func TestEmployeeService_DeleteEmployee(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewEmployeeService(testDB)

	// Create test employee
	req := &models.CreateEmployeeRequest{
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john@example.com",
		JobTitle:   "Developer",
		Country:    "USA",
		Salary:     75000.0,
		Department: "Engineering",
	}

	created, err := service.CreateEmployee(req)
	if err != nil {
		t.Fatalf("Failed to create test employee: %v", err)
	}

	// Delete employee
	err = service.DeleteEmployee(created.ID)
	if err != nil {
		t.Fatalf("Failed to delete employee: %v", err)
	}

	// Verify employee is deleted
	_, err = service.GetEmployeeByID(created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted employee")
	}

	if err != gorm.ErrRecordNotFound && err.Error() != "employee not found" {
		t.Errorf("Expected record not found error, got: %v", err)
	}
}

func TestEmployeeService_DeleteEmployee_NonExistent(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewEmployeeService(testDB)

	// Try to delete non-existent employee
	err := service.DeleteEmployee(999)
	if err == nil {
		t.Error("Expected error for non-existent employee")
	}

	if err.Error() != "employee not found" {
		t.Errorf("Expected 'employee not found' error, got: %v", err)
	}
}

func TestEmployeeService_ListEmployees_InvalidPagination(t *testing.T) {
	// Setup test database
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	testDB := testutils.TestDB(t)
	defer testutils.CleanupTestDB(testDB)
	database.DB = testDB

	service := NewEmployeeService(testDB)

	// Create test employee
	req := &models.CreateEmployeeRequest{
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john@example.com",
		JobTitle:   "Developer",
		Country:    "USA",
		Salary:     75000.0,
		Department: "Engineering",
	}
	if _, err := service.CreateEmployee(req); err != nil {
		t.Fatalf("Failed to create test employee: %v", err)
	}

	// Test with invalid page (should default to 1)
	response, err := service.ListEmployees(0, 10, "")
	if err != nil {
		t.Fatalf("Failed to list employees: %v", err)
	}

	if response.Page != 1 {
		t.Errorf("Expected page to default to 1, got %d", response.Page)
	}

	// Test with limit > 100 (should default to 50)
	response, err = service.ListEmployees(1, 200, "")
	if err != nil {
		t.Fatalf("Failed to list employees: %v", err)
	}

	if response.Limit != 50 {
		t.Errorf("Expected limit to default to 50, got %d", response.Limit)
	}
}
