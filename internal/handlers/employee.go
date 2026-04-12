package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tofiquem/assingment/internal/database"
	"github.com/tofiquem/assingment/internal/models"
	"gorm.io/gorm"
)

type EmployeeHandler struct {
	db *gorm.DB
}

func NewEmployeeHandler() *EmployeeHandler {
	return &EmployeeHandler{
		db: database.DB,
	}
}

func (h *EmployeeHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/employees", h.GetEmployees).Methods("GET")
	r.HandleFunc("/api/employees", h.CreateEmployee).Methods("POST")
	r.HandleFunc("/api/employees/{id}", h.GetEmployee).Methods("GET")
	r.HandleFunc("/api/employees/{id}", h.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/api/employees/{id}", h.DeleteEmployee).Methods("DELETE")
}

func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 50

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

	if err := h.db.Model(&models.Employee{}).Count(&total).Error; err != nil {
		http.Error(w, "Failed to count employees", http.StatusInternalServerError)
		return
	}

	if err := h.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&employees).Error; err != nil {
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
