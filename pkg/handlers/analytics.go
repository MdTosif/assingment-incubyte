package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tofiquem/assingment/pkg/database"
	"github.com/tofiquem/assingment/pkg/models"
	"gorm.io/gorm"
)

type AnalyticsHandler struct {
	db *gorm.DB
}

func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{
		db: database.DB,
	}
}

func (h *AnalyticsHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/analytics/salary/by-country", h.GetSalaryByCountry).Methods("GET")
	r.HandleFunc("/api/analytics/salary/by-job-title/{country}", h.GetSalaryByJobTitleInCountry).Methods("GET")
	r.HandleFunc("/api/analytics/salary/department-insights", h.GetDepartmentInsights).Methods("GET")
}

type CountrySalaryStats struct {
	Country string  `json:"country"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Average float64 `json:"average"`
	Count   int64   `json:"count"`
}

type JobTitleSalaryStats struct {
	JobTitle string  `json:"jobTitle"`
	Average  float64 `json:"average"`
	Count    int64   `json:"count"`
}

type DepartmentSalaryStats struct {
	Department string  `json:"department"`
	Min        float64 `json:"min"`
	Max        float64 `json:"max"`
	Average    float64 `json:"average"`
	Count      int64   `json:"count"`
}

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
