// Package handlers provides HTTP request handlers for the salary management API.
// It handles employee CRUD operations and serves JSON responses.
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tofiquem/assingment/pkg/database"
	"github.com/tofiquem/assingment/pkg/models"
	"gorm.io/gorm"
)

// ==================== Handler Definition ====================

// EmployeeHandler handles HTTP requests for employee management.
// It provides endpoints for CRUD operations on employee records.
type EmployeeHandler struct {
	db *gorm.DB
}

// ==================== Constructor & Route Registration ====================

// NewEmployeeHandler creates a new EmployeeHandler with the default database connection.
func NewEmployeeHandler() *EmployeeHandler {
	return &EmployeeHandler{
		db: database.DB,
	}
}

// RegisterRoutes registers all employee routes with the given router.
// Routes: GET /api/employees, POST /api/employees, GET/PUT/DELETE /api/employees/{id}
func (h *EmployeeHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/employees", h.GetEmployees).Methods("GET")
	r.HandleFunc("/api/employees", h.CreateEmployee).Methods("POST")
	r.HandleFunc("/api/employees/{id}", h.GetEmployee).Methods("GET")
	r.HandleFunc("/api/employees/{id}", h.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/api/employees/{id}", h.DeleteEmployee).Methods("DELETE")
}

// ==================== CRUD Operations ====================

// GetEmployees returns a paginated list of employees with optional search.
// Query params: page (default: 1), limit (default: 50, max: 100), search (optional)
// Response: {employees, total, page, limit, pages}
func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 50
	search := r.URL.Query().Get("search")

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := (page - 1) * limit

	var employees []models.Employee
	var total int64

	query := h.db.Model(&models.Employee{})

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR job_title ILIKE ? OR country ILIKE ? OR department ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		http.Error(w, "Failed to count employees", http.StatusInternalServerError)
		return
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&employees).Error; err != nil {
		http.Error(w, "Failed to fetch employees", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"employees": employees,
		"total":     total,
		"page":      page,
		"limit":     limit,
		"pages":     (total + int64(limit) - 1) / int64(limit),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateEmployee creates a new employee from the request body.
// Expects: CreateEmployeeRequest JSON body
// Returns: 201 Created with the created employee
func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	employee := models.ToEmployee(&req)
	if err := h.db.Create(employee).Error; err != nil {
		http.Error(w, "Failed to create employee", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(employee)
}

// GetEmployee returns a single employee by ID.
// Returns: 200 OK with employee, 404 Not Found, or 400 Bad Request
func (h *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var employee models.Employee
	if err := h.db.First(&employee, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Employee not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch employee", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employee)
}

// UpdateEmployee updates an existing employee by ID.
// Expects: UpdateEmployeeRequest JSON body (partial updates supported)
// Returns: 200 OK with updated employee, 404 Not Found, or 400 Bad Request
func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var employee models.Employee
	if err := h.db.First(&employee, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Employee not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch employee", http.StatusInternalServerError)
		}
		return
	}

	employee.UpdateFromRequest(&req)
	if err := h.db.Save(&employee).Error; err != nil {
		http.Error(w, "Failed to update employee", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employee)
}

// DeleteEmployee permanently deletes an employee by ID.
// Returns: 204 No Content, 404 Not Found, or 400 Bad Request
func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var employee models.Employee
	if err := h.db.First(&employee, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Employee not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch employee", http.StatusInternalServerError)
		}
		return
	}

	if err := h.db.Delete(&employee).Error; err != nil {
		http.Error(w, "Failed to delete employee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
