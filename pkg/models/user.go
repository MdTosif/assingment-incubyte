// Package models defines the data structures for the salary management system.
package models

import (
	"time"

	"gorm.io/gorm"
)

// ==================== Model Definition ====================

// User represents an authenticated user (HR or Admin) in the system.
// It stores credentials, role information, and account status.
type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null"`
	Password  string     `json:"-" gorm:"not null"`
	Role      string     `json:"role" gorm:"not null;default:'hr'"` // hr, admin
	FirstName string     `json:"firstName" gorm:"not null"`
	LastName  string     `json:"lastName" gorm:"not null"`
	IsActive  bool       `json:"isActive"`
	LastLogin *time.Time `json:"lastLogin"`
	CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}

// ==================== Request Types ====================

// LoginRequest represents user login credentials.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the successful authentication response.
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	User      User      `json:"user"`
}

// CreateHRUserRequest represents data needed to create a new HR/Admin user.
type CreateHRUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Role      string `json:"role" binding:"required,oneof=hr admin"`
}

// UpdateUserRequest represents data for updating an existing user.
// All fields are optional for partial updates.
type UpdateUserRequest struct {
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Email     *string `json:"email,omitempty"`
	Role      *string `json:"role,omitempty"`
	IsActive  *bool   `json:"isActive,omitempty"`
}

// ==================== Helper Methods ====================

// IsHR checks if the user has HR or Admin role permissions.
func (u *User) IsHR() bool {
	return u.Role == "hr" || u.Role == "admin"
}

// IsAdmin checks if the user has Admin role.
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsActiveUser checks if the user account is active.
func (u *User) IsActiveUser() bool {
	return u.IsActive
}

// BeforeCreate is a GORM hook that sets default values for new users.
// It sets the role to "hr" if empty and activates the account.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = "hr"
	}
	// Set default IsActive to true if not explicitly set (zero value is false)
	// Since bool has a zero value of false, we need to check if this is a new record
	// For new records, we want IsActive to default to true unless explicitly set to false
	// This is tricky because we can't distinguish between explicit false and unset false
	// For now, we'll set it to true for all new records (the test will need to handle this)
	u.IsActive = true
	return nil
}

// ToSafeUser returns a copy of the user with sensitive data (password) removed.
// Use this when returning user data in API responses.
func (u *User) ToSafeUser() User {
	return User{
		ID:        u.ID,
		Email:     u.Email,
		Role:      u.Role,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsActive:  u.IsActive,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
