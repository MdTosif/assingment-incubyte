// Package services provides business logic for employee management.
// It handles CRUD operations and search functionality for employee records.
package services

import (
	"fmt"

	"github.com/tofiquem/assingment/pkg/models"
	"gorm.io/gorm"
)

// ==================== Service Definition ====================

// EmployeeService handles employee management operations.
// It provides methods for CRUD operations and searching employees.
type EmployeeService struct {
	db *gorm.DB
}

// ==================== Constructor ====================

// NewEmployeeService creates a new EmployeeService with the given database connection.
func NewEmployeeService(db *gorm.DB) *EmployeeService {
	return &EmployeeService{
		db: db,
	}
}

// ==================== Response Types ====================

// EmployeeListResponse represents a paginated list of employees.
type EmployeeListResponse struct {
	Employees []models.Employee `json:"employees"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	Limit     int               `json:"limit"`
	Pages     int64             `json:"pages"`
}

// ==================== CRUD Operations ====================

// ListEmployees returns a paginated list of employees with optional search.
// Query params: page (default: 1), limit (default: 50, max: 100), search (optional)
// Response: EmployeeListResponse with employees, total, page, limit, and pages.
func (s *EmployeeService) ListEmployees(page, limit int, search string) (*EmployeeListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	var employees []models.Employee
	var total int64

	query := s.db.Model(&models.Employee{})

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"LOWER(first_name) LIKE LOWER(?) OR LOWER(last_name) LIKE LOWER(?) OR LOWER(email) LIKE LOWER(?) OR LOWER(job_title) LIKE LOWER(?) OR LOWER(country) LIKE LOWER(?) OR LOWER(department) LIKE LOWER(?)",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count employees: %w", err)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch employees: %w", err)
	}

	pages := (total + int64(limit) - 1) / int64(limit)

	return &EmployeeListResponse{
		Employees: employees,
		Total:     total,
		Page:      page,
		Limit:     limit,
		Pages:     pages,
	}, nil
}

// CreateEmployee creates a new employee from the request.
// Returns the created employee or an error if creation fails.
func (s *EmployeeService) CreateEmployee(req *models.CreateEmployeeRequest) (*models.Employee, error) {
	employee := models.ToEmployee(req)
	if err := s.db.Create(employee).Error; err != nil {
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}
	return employee, nil
}

// GetEmployeeByID retrieves an employee by their ID.
// Returns the employee or an error if not found.
func (s *EmployeeService) GetEmployeeByID(id uint) (*models.Employee, error) {
	var employee models.Employee
	if err := s.db.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &employee, nil
}

// UpdateEmployee updates an existing employee by ID.
// Only provided fields (non-nil) are updated, allowing partial updates.
// Returns the updated employee or an error if update fails.
func (s *EmployeeService) UpdateEmployee(id uint, req *models.UpdateEmployeeRequest) (*models.Employee, error) {
	var employee models.Employee
	if err := s.db.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	employee.UpdateFromRequest(req)
	if err := s.db.Save(&employee).Error; err != nil {
		return nil, fmt.Errorf("failed to update employee: %w", err)
	}

	return &employee, nil
}

// DeleteEmployee permanently deletes an employee by ID.
// Returns an error if the employee is not found or deletion fails.
func (s *EmployeeService) DeleteEmployee(id uint) error {
	var employee models.Employee
	if err := s.db.First(&employee, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("employee not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	if err := s.db.Delete(&employee).Error; err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	return nil
}
