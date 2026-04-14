// Package services provides business logic for salary analytics.
// It handles aggregated salary statistics by country, job title, and department.
package services

import (
	"fmt"

	"github.com/tofiquem/assingment/pkg/models"
	"gorm.io/gorm"
)

// ==================== Service Definition ====================

// AnalyticsService handles salary analytics operations.
// It provides methods for calculating salary statistics and insights.
type AnalyticsService struct {
	db *gorm.DB
}

// ==================== Constructor ====================

// NewAnalyticsService creates a new AnalyticsService with the given database connection.
func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{
		db: db,
	}
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
func (s *AnalyticsService) GetSalaryByCountry() ([]CountrySalaryStats, error) {
	var stats []CountrySalaryStats

	err := s.db.Model(&models.Employee{}).
		Select("country, MIN(salary) as min, MAX(salary) as max, AVG(salary) as average, COUNT(*) as count").
		Group("country").
		Order("average DESC").
		Scan(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch salary statistics by country: %w", err)
	}

	return stats, nil
}

// GetSalaryByJobTitleInCountry returns salary statistics grouped by job title for a specific country.
// Results are ordered by average salary in descending order.
func (s *AnalyticsService) GetSalaryByJobTitleInCountry(country string) ([]JobTitleSalaryStats, error) {
	if country == "" {
		return nil, fmt.Errorf("country parameter is required")
	}

	var stats []JobTitleSalaryStats

	err := s.db.Model(&models.Employee{}).
		Select("job_title as job_title, AVG(salary) as average, COUNT(*) as count").
		Where("country = ?", country).
		Group("job_title").
		Order("average DESC").
		Scan(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch salary statistics by job title: %w", err)
	}

	return stats, nil
}

// GetDepartmentInsights returns salary statistics grouped by department.
// Results are ordered by average salary in descending order.
func (s *AnalyticsService) GetDepartmentInsights() ([]DepartmentSalaryStats, error) {
	var stats []DepartmentSalaryStats

	err := s.db.Model(&models.Employee{}).
		Select("department, MIN(salary) as min, MAX(salary) as max, AVG(salary) as average, COUNT(*) as count").
		Group("department").
		Order("average DESC").
		Scan(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch department salary insights: %w", err)
	}

	return stats, nil
}

// GetDepartmentInsightsByCountry returns salary statistics grouped by department for a specific country.
// Results are ordered by average salary in descending order.
func (s *AnalyticsService) GetDepartmentInsightsByCountry(country string) ([]DepartmentSalaryStats, error) {
	if country == "" {
		return nil, fmt.Errorf("country parameter is required")
	}

	var stats []DepartmentSalaryStats

	err := s.db.Model(&models.Employee{}).
		Select("department, MIN(salary) as min, MAX(salary) as max, AVG(salary) as average, COUNT(*) as count").
		Where("country = ?", country).
		Group("department").
		Order("average DESC").
		Scan(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch department salary insights by country: %w", err)
	}

	return stats, nil
}
