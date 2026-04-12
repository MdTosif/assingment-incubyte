package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	FirstName  string    `json:"firstName" gorm:"not null"`
	LastName   string    `json:"lastName" gorm:"not null"`
	Email      string    `json:"email" gorm:"uniqueIndex;not null"`
	JobTitle   string    `json:"jobTitle" gorm:"not null;index"`
	Country    string    `json:"country" gorm:"not null;index"`
	Salary     float64   `json:"salary" gorm:"not null;index"`
	Department string    `json:"department" gorm:"not null"`
	HireDate   time.Time `json:"hireDate" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type CreateEmployeeRequest struct {
	FirstName  string  `json:"firstName" binding:"required"`
	LastName   string  `json:"lastName" binding:"required"`
	Email      string  `json:"email" binding:"required,email"`
	JobTitle   string  `json:"jobTitle" binding:"required"`
	Country    string  `json:"country" binding:"required"`
	Salary     float64 `json:"salary" binding:"required,min=0"`
	Department string  `json:"department" binding:"required"`
}

type UpdateEmployeeRequest struct {
	FirstName  *string  `json:"firstName,omitempty"`
	LastName   *string  `json:"lastName,omitempty"`
	Email      *string  `json:"email,omitempty"`
	JobTitle   *string  `json:"jobTitle,omitempty"`
	Country    *string  `json:"country,omitempty"`
	Salary     *float64 `json:"salary,omitempty"`
	Department *string  `json:"department,omitempty"`
}

func (e *Employee) BeforeCreate(tx *gorm.DB) error {
	if e.HireDate.IsZero() {
		e.HireDate = time.Now()
	}
	return nil
}

func ToEmployee(req *CreateEmployeeRequest) *Employee {
	if req == nil {
		return nil
	}

	return &Employee{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		JobTitle:   req.JobTitle,
		Country:    req.Country,
		Salary:     req.Salary,
		Department: req.Department,
	}
}

func (e *Employee) UpdateFromRequest(req *UpdateEmployeeRequest) {
	if req.FirstName != nil {
		e.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		e.LastName = *req.LastName
	}
	if req.Email != nil {
		e.Email = *req.Email
	}
	if req.JobTitle != nil {
		e.JobTitle = *req.JobTitle
	}
	if req.Country != nil {
		e.Country = *req.Country
	}
	if req.Salary != nil {
		e.Salary = *req.Salary
	}
	if req.Department != nil {
		e.Department = *req.Department
	}
}
