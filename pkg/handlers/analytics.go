// Package handlers provides HTTP request handlers for the salary management API.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tofiquem/assingment/pkg/database"
	"github.com/tofiquem/assingment/pkg/models"
	"gorm.io/gorm"
)

// ==================== Handler Definition ====================

// AnalyticsHandler handles HTTP requests for salary analytics.
// It provides aggregated salary statistics by country, job title, and department.
type AnalyticsHandler struct {
	db *gorm.DB
}

// ==================== Constructor ====================

// NewAnalyticsHandler creates a new AnalyticsHandler with the default database connection.
func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{
		db: database.DB,
	}
}

// RegisterRoutes registers all analytics routes with the given router.
// Routes: GET /api/analytics/salary/by-country, GET /api/analytics/salary/by-job-title/{country}, GET /api/analytics/salary/department-insights
func (h *AnalyticsHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/analytics/salary/by-country", h.GetSalaryByCountry).Methods("GET")
	r.HandleFunc("/api/analytics/salary/by-job-title/{country}", h.GetSalaryByJobTitleInCountry).Methods("GET")
	r.HandleFunc("/api/analytics/salary/department-insights", h.GetDepartmentInsights).Methods("GET")
}

// ==================== Response Types ====================

// CountrySalaryStats represents salary statistics aggregated by country.
type CountrySalaryStats struct {
	Country string  `json:"country"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Average float64 `json:"average"`
	Count   int64   `json:"count"`
}

// JobTitleSalaryStats represents salary statistics aggregated by job title within a country.
type JobTitleSalaryStats struct {
	JobTitle string  `json:"jobTitle"`
	Average  float64 `json:"average"`
	Count    int64   `json:"count"`
}

// DepartmentSalaryStats represents salary statistics aggregated by department.
type DepartmentSalaryStats struct {
	Department string  `json:"department"`
	Min        float64 `json:"min"`
	Max        float64 `json:"max"`
	Average    float64 `json:"average"`
	Count      int64   `json:"count"`
}

// ==================== Salary Analytics ====================

// GetSalaryByCountry returns salary statistics (min, max, average, count) grouped by country.
// Results are ordered by average salary in descending order.
func (h *AnalyticsHandler) GetSalaryByCountry(w http.ResponseWriter, r *http.Request) {
	var stats []CountrySalaryStats

	err := h.db.Model(&models.Employee{}).
		Select("country, MIN(salary) as min, MAX(salary) as max, AVG(salary) as average, COUNT(*) as count").
		Group("country").
		Order("average DESC").
		Scan(&stats).Error

	if err != nil {
		http.Error(w, "Failed to fetch salary statistics by country", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetSalaryByJobTitleInCountry returns salary statistics grouped by job title for a specific country.
// Path param: country (URL-encoded country name)
// Results are ordered by average salary in descending order.
func (h *AnalyticsHandler) GetSalaryByJobTitleInCountry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	country := vars["country"]

	if country == "" {
		http.Error(w, "Country parameter is required", http.StatusBadRequest)
		return
	}

	var stats []JobTitleSalaryStats

	err := h.db.Model(&models.Employee{}).
		Select("job_title, AVG(salary) as average, COUNT(*) as count").
		Where("country = ?", country).
		Group("job_title").
		Order("average DESC").
		Scan(&stats).Error

	if err != nil {
		http.Error(w, "Failed to fetch salary statistics by job title", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetDepartmentInsights returns salary statistics grouped by department.
// Results are ordered by average salary in descending order.
func (h *AnalyticsHandler) GetDepartmentInsights(w http.ResponseWriter, r *http.Request) {
	var stats []DepartmentSalaryStats

	err := h.db.Model(&models.Employee{}).
		Select("department, MIN(salary) as min, MAX(salary) as max, AVG(salary) as average, COUNT(*) as count").
		Group("department").
		Order("average DESC").
		Scan(&stats).Error

	if err != nil {
		http.Error(w, "Failed to fetch department salary insights", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
