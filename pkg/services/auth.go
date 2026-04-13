// Package services provides business logic for authentication and user management.
// It handles login, user CRUD operations, token validation, and password management.
package services

import (
	"fmt"
	"time"

	"github.com/tofiquem/assingment/pkg/models"
	"gorm.io/gorm"
)

// ==================== Service Definition ====================

// AuthService handles authentication and user management operations.
// It uses JWT for token generation and bcrypt for password hashing.
type AuthService struct {
	db              *gorm.DB
	jwtService      *JWTService
	passwordService *PasswordService
}

// ==================== Constructor ====================

// NewAuthService creates a new AuthService with the given database connection.
// It initializes the required JWT and password services.
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		db:              db,
		jwtService:      NewJWTService(),
		passwordService: NewPasswordService(),
	}
}

// ==================== Authentication ====================

// Login authenticates a user with email and password.
// It returns a JWT token and user information on success.
// Errors: invalid credentials, inactive account, database errors.
func (a *AuthService) Login(email, password string) (*models.LoginResponse, error) {
	// Find user by email
	var user models.User
	if err := a.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	// Verify password
	if err := a.passwordService.VerifyPassword(user.Password, password); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	if err := a.db.Save(&user).Error; err != nil {
		// Log error but don't fail login
		fmt.Printf("Warning: failed to update last login: %v\n", err)
	}

	// Generate JWT token
	token, expiresAt, err := a.jwtService.GenerateToken(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return response with safe user data
	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user.ToSafeUser(),
	}, nil
}

// CreateUser creates a new HR or Admin user.
// It validates password strength and hashes the password before storage.
// Only admin users should have access to this function.
func (a *AuthService) CreateUser(req *models.CreateHRUserRequest) (*models.User, error) {
	// Validate password strength
	if err := a.passwordService.ValidatePasswordStrength(req.Password); err != nil {
		return nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Hash password
	hashedPassword, err := a.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      req.Role,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := a.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// ==================== User Management ====================

// GetUserByID retrieves a user by their ID.
// Returns the user with password intact (internal use only).
func (a *AuthService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := a.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email address.
// Returns the user with password intact (for authentication use).
func (a *AuthService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := a.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}

// ==================== Token Validation ====================

// ValidateToken validates a JWT token and returns the associated user.
// It checks token validity, expiration, and user account status.
func (a *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	// Validate token
	claims, err := a.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Check if token is expired
	if a.jwtService.IsTokenExpired(*claims) {
		return nil, fmt.Errorf("token has expired")
	}

	// Extract user ID
	userID, err := a.jwtService.GetUserIDFromToken(*claims)
	if err != nil {
		return nil, fmt.Errorf("failed to extract user ID: %w", err)
	}

	// Get user from database
	user, err := a.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	return user, nil
}

// UpdateUser updates an existing user's information.
// Only provided fields (non-nil) are updated, allowing partial updates.
func (a *AuthService) UpdateUser(id uint, req *models.UpdateUserRequest) (*models.User, error) {
	var user models.User
	if err := a.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Update fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := a.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

// DeleteUser permanently deletes a user from the database.
// Use with caution - this action cannot be undone.
func (a *AuthService) DeleteUser(id uint) error {
	var user models.User
	if err := a.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	if err := a.db.Delete(&user).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers returns a paginated list of users without sensitive data.
// Passwords are removed from the returned users for security.
func (a *AuthService) ListUsers(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * limit

	// Count total users
	if err := a.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users with pagination
	if err := a.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch users: %w", err)
	}

	// Remove passwords from response
	for i := range users {
		users[i] = users[i].ToSafeUser()
	}

	return users, total, nil
}

// ChangePassword changes a user's password after verifying the current password.
// It validates the new password strength before updating.
func (a *AuthService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	// Get user
	user, err := a.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify current password
	if err := a.passwordService.VerifyPassword(user.Password, currentPassword); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Validate new password
	if err := a.passwordService.ValidatePasswordStrength(newPassword); err != nil {
		return fmt.Errorf("new password validation failed: %w", err)
	}

	// Hash new password
	hashedPassword, err := a.passwordService.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	user.Password = hashedPassword
	if err := a.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
