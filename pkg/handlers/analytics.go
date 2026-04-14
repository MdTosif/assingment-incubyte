// Package handlers provides HTTP request handlers for the salary management API.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tofiquem/assingment/pkg/services"
	"gorm.io/gorm"
)

// ==================== Handler Definition ====================

// AnalyticsHandler handles HTTP requests for salary analytics.
// It provides aggregated salary statistics by country, job title, and department.
type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

// ==================== Constructor ====================

// NewAnalyticsHandler creates a new AnalyticsHandler with the given database connection.
func NewAnalyticsHandler(db *gorm.DB) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: services.NewAnalyticsService(db),
	}
}

// RegisterRoutes registers all analytics routes with the given router.
// Routes: GET /api/analytics/salary/by-country, GET /api/analytics/salary/by-job-title/{country}, GET /api/analytics/salary/department-insights, GET /api/analytics/salary/department-insights/{country}
func (h *AnalyticsHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/analytics/salary/by-country", h.GetSalaryByCountry).Methods("GET")
	r.HandleFunc("/api/analytics/salary/by-job-title/{country}", h.GetSalaryByJobTitleInCountry).Methods("GET")
	r.HandleFunc("/api/analytics/salary/department-insights", h.GetDepartmentInsights).Methods("GET")
	r.HandleFunc("/api/analytics/salary/department-insights/{country}", h.GetDepartmentInsightsByCountry).Methods("GET")
}

// ==================== Salary Analytics ====================

// GetSalaryByCountry returns salary statistics (min, max, average, count) grouped by country.
// Results are ordered by average salary in descending order.
func (h *AnalyticsHandler) GetSalaryByCountry(w http.ResponseWriter, r *http.Request) {
	stats, err := h.analyticsService.GetSalaryByCountry()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	stats, err := h.analyticsService.GetSalaryByJobTitleInCountry(country)
	if err != nil {
		if err.Error() == "country parameter is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetDepartmentInsights returns salary statistics grouped by department.
// Results are ordered by average salary in descending order.
func (h *AnalyticsHandler) GetDepartmentInsights(w http.ResponseWriter, r *http.Request) {
	stats, err := h.analyticsService.GetDepartmentInsights()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetDepartmentInsightsByCountry returns salary statistics grouped by department for a specific country.
// Path param: country (URL-encoded country name)
// Results are ordered by average salary in descending order.
func (h *AnalyticsHandler) GetDepartmentInsightsByCountry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	country := vars["country"]

	if country == "" {
		http.Error(w, "Country parameter is required", http.StatusBadRequest)
		return
	}

	stats, err := h.analyticsService.GetDepartmentInsightsByCountry(country)
	if err != nil {
		if err.Error() == "country parameter is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
