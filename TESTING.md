# Testing Documentation

This document explains the Test-Driven Development (TDD) approach used in this Go application and how to run and maintain the tests.

## Overview

The application follows TDD principles with comprehensive test coverage including:
- Unit tests for models and business logic
- Unit tests for HTTP handlers
- Integration tests for API endpoints
- Test utilities and helpers

## Test Structure

```
internal/
  models/
    employee.go
    employee_test.go          # Model unit tests
  handlers/
    employee.go
    employee_test.go          # Handler unit tests
    analytics.go
    analytics_test.go         # Analytics handler tests
  testutils/
    testutils.go              # Test utilities and helpers
  database/
    gorm.go
cmd/
  server/
    main.go
    integration_test.go      # Integration tests
```

## Test Categories

### 1. Model Tests (`internal/models/employee_test.go`)

Tests for the Employee model including:
- `ToEmployee()` function
- `UpdateFromRequest()` method
- `BeforeCreate()` hook
- Database validation
- Unique email constraints

**Example:**
```go
func TestToEmployee(t *testing.T) {
    req := &CreateEmployeeRequest{
        FirstName: "John",
        LastName:  "Doe",
        Email:     "john@example.com",
        // ... other fields
    }
    
    employee := ToEmployee(req)
    
    assert.Equal(t, "John", employee.FirstName)
    assert.Equal(t, "Doe", employee.LastName)
    // ... more assertions
}
```

### 2. Handler Tests (`internal/handlers/`)

Tests for HTTP handlers including:
- Employee CRUD operations
- Analytics endpoints
- Error handling
- Request/response validation
- Route registration

**Example:**
```go
func TestEmployeeHandler_CreateEmployee(t *testing.T) {
    testDB := testutils.TestDB(t)
    defer testutils.CleanupTestDB(testDB)
    
    handler := NewEmployeeHandler()
    
    req := models.CreateEmployeeRequest{
        FirstName: "John",
        LastName:  "Doe",
        // ... other fields
    }
    
    rr, err := testutils.ExecuteRequest(handler.CreateEmployee, "POST", "/api/employees", req)
    
    testutils.AssertStatusCode(t, rr, http.StatusCreated)
    // ... more assertions
}
```

### 3. Integration Tests (`cmd/server/integration_test.go`)

End-to-end tests including:
- Complete API workflows
- Multiple handler interactions
- Error scenarios
- Environment variable handling
- SPA handler functionality

### 4. Test Utilities (`internal/testutils/testutils.go`)

Helper functions for testing:
- In-memory database setup
- Test data creation
- HTTP request/response helpers
- Assertion utilities
- Environment variable mocking

## Running Tests

### Using Makefile (Recommended)

```bash
# Run all tests
make test

# Run tests with verbose output and race detection
make test-verbose

# Run tests with coverage
make test-coverage

# Generate HTML coverage report
make coverage-html

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run specific test
make test-specific RUN=TestEmployeeHandler_CreateEmployee

# Run model tests only
make tdd-models

# Run handler tests only
make tdd-handlers

# Run integration tests only
make tdd-integration

# Watch for changes and run tests (requires gow)
make tdd-watch

# Run benchmark tests
make benchmark
```

### Using Go Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Run specific package tests
go test -v ./internal/models/...
go test -v ./internal/handlers/...
go test -v ./cmd/server/...

# Run specific test
go test -v -run TestEmployeeHandler_CreateEmployee ./internal/handlers/...

# Run tests with race detection
go test -race ./...

# Run benchmark tests
go test -bench=. -benchmem ./...
```

## TDD Workflow

### 1. Write the Test First

```go
func TestNewFeature(t *testing.T) {
    // Arrange
    testDB := testutils.TestDB(t)
    defer testutils.CleanupTestDB(testDB)
    
    // Act
    result := SomeFunction(testDB)
    
    // Assert
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### 2. Run the Test (It Should Fail)

```bash
make test-unit
```

### 3. Write Minimal Code to Pass

```go
func SomeFunction(db *gorm.DB) string {
    return expected
}
```

### 4. Run Test Again (Should Pass)

```bash
make test-unit
```

### 5. Refactor and Improve

Improve the implementation while keeping tests green.

## Test Best Practices

### 1. Test Structure

Use **Arrange-Act-Assert** pattern:
```go
func TestEmployeeCreation(t *testing.T) {
    // Arrange
    testDB := testutils.TestDB(t)
    defer testutils.CleanupTestDB(testDB)
    
    // Act
    employee := testutils.CreateTestEmployee(testDB, "John", "Doe", ...)
    
    // Assert
    assert.NotNil(t, employee)
    assert.Equal(t, "John", employee.FirstName)
}
```

### 2. Test Naming

Use descriptive names that explain what is being tested:
- `TestEmployeeHandler_CreateEmployee_ValidRequest`
- `TestEmployeeHandler_CreateEmployee_InvalidJSON`
- `TestEmployeeHandler_CreateEmployee_DuplicateEmail`

### 3. Test Isolation

Each test should be independent:
- Use in-memory databases
- Clean up resources in `defer`
- Don't rely on test order

### 4. Table-Driven Tests

Use table-driven tests for multiple scenarios:
```go
tests := []struct {
    name     string
    input    string
    expected string
    wantErr  bool
}{
    {"valid input", "valid", "expected", false},
    {"invalid input", "invalid", "", true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        result, err := SomeFunction(tt.input)
        if (err != nil) != tt.wantErr {
            t.Errorf("SomeFunction() error = %v, wantErr %v", err, tt.wantErr)
            return
        }
        if result != tt.expected {
            t.Errorf("SomeFunction() = %v, want %v", result, tt.expected)
        }
    })
}
```

### 5. Mock Dependencies

Use test utilities to mock external dependencies:
```go
func TestWithMockDB(t *testing.T) {
    originalDB := database.DB
    defer func() { database.DB = originalDB }()
    
    testDB := testutils.TestDB(t)
    database.DB = testDB
    
    // Test with mock database
}
```

## Coverage

The test suite aims for high coverage across:
- Model functions and validation
- HTTP handler methods
- Error handling paths
- Edge cases and boundary conditions

To check coverage:
```bash
make test-coverage
go tool cover -func=coverage.out
```

## Continuous Integration

The tests are designed to run in CI/CD environments:
- Use in-memory databases (no external dependencies)
- Fast execution (under 30 seconds)
- Clear output for debugging
- Coverage reporting

## Debugging Tests

### 1. Verbose Output
```bash
go test -v ./internal/handlers/...
```

### 2. Run Specific Test
```bash
go test -v -run TestEmployeeHandler_CreateEmployee ./internal/handlers/...
```

### 3. Test with Race Detection
```bash
go test -race ./...
```

### 4. Print Debug Information
```go
t.Logf("Debug info: %+v", someVariable)
```

## Adding New Tests

When adding new features:

1. **Write tests first** following TDD principles
2. **Use existing test utilities** (`testutils.TestDB`, `testutils.CreateTestEmployee`, etc.)
3. **Follow naming conventions** (`Test[Function]_[Scenario]`)
4. **Test both happy path and error cases**
5. **Ensure test isolation** (cleanup resources)
6. **Add integration tests** for API endpoints

## Example: Adding a New Endpoint

### 1. Write the Test First
```go
func TestEmployeeHandler_GetEmployeesByDepartment(t *testing.T) {
    testDB := testutils.TestDB(t)
    defer testutils.CleanupTestDB(testDB)
    
    // Create test data
    testutils.CreateTestEmployees(testDB)
    
    handler := NewEmployeeHandler()
    
    req, _ := http.NewRequest("GET", "/api/employees/department/Engineering", nil)
    req = mux.SetURLVars(req, map[string]string{"department": "Engineering"})
    
    rr := httptest.NewRecorder()
    handler.GetEmployeesByDepartment(rr, req)
    
    testutils.AssertStatusCode(t, rr, http.StatusOK)
    // ... more assertions
}
```

### 2. Implement the Handler
```go
func (h *EmployeeHandler) GetEmployeesByDepartment(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

### 3. Register the Route
```go
func (h *EmployeeHandler) RegisterRoutes(r *mux.Router) {
    // ... existing routes
    r.HandleFunc("/api/employees/department/{department}", h.GetEmployeesByDepartment).Methods("GET")
}
```

### 4. Run Tests
```bash
make test-unit
```

This TDD approach ensures robust, well-tested code with high confidence in the implementation.
