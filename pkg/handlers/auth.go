// Package handlers provides HTTP request handlers for the salary management API.
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/tofiquem/assingment/pkg/models"
	"github.com/tofiquem/assingment/pkg/services"
	"gorm.io/gorm"
)

// ==================== Handler Definition ====================

// AuthHandler handles HTTP requests for authentication and user management.
// It provides endpoints for login, user CRUD operations, and middleware for authorization.
type AuthHandler struct {
	authService *services.AuthService
}

// ==================== Constructor ====================

// NewAuthHandler creates a new AuthHandler with the given database connection.
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(db),
	}
}

// ==================== Route Registration ====================

// RegisterRoutes registers all authentication and user routes with the given router.
// Public routes: POST /api/auth/login
// Protected routes: GET /api/auth/me, POST /api/auth/logout, POST /api/auth/change-password
// Admin routes: GET/POST /api/admin/users, PUT/DELETE /api/admin/users/{id}
func (h *AuthHandler) RegisterRoutes(r *mux.Router) {
	// Public routes
	r.HandleFunc("/api/auth/login", h.Login).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api/auth").Subrouter()
	protected.Use(h.authMiddleware)
	protected.HandleFunc("/me", h.GetMe).Methods("GET")
	protected.HandleFunc("/logout", h.Logout).Methods("POST")
	protected.HandleFunc("/change-password", h.ChangePassword).Methods("POST")

	// Admin only routes
	admin := r.PathPrefix("/api/admin").Subrouter()
	admin.Use(h.authMiddleware)
	admin.Use(h.adminMiddleware)
	admin.HandleFunc("/users", h.CreateUser).Methods("POST")
	admin.HandleFunc("/users", h.ListUsers).Methods("GET")
	admin.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	admin.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")

	// Register protected employee routes
	h.registerEmployeeRoutes(r)

	// Register protected analytics routes
	h.registerAnalyticsRoutes(r)
}

// ==================== Public Authentication ====================

// Login authenticates a user with email and password.
// Expects: LoginRequest JSON body
// Returns: 200 OK with LoginResponse (token, expiresAt, user), or 401 Unauthorized
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate user
	response, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ==================== Protected User Operations ====================

// GetMe returns the currently authenticated user from the request context.
// Returns: 200 OK with user data (safe user without password)
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(*models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToSafeUser())
}

// Logout handles user logout. In a stateless JWT system, the actual token removal
// is handled client-side. This endpoint provides a server-side confirmation.
// Returns: 200 OK with success message
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a stateless JWT system, logout is handled client-side
	// We can implement token blacklisting if needed
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// ChangePassword handles password change for the currently authenticated user.
// Expects: JSON body with currentPassword and newPassword
// Returns: 200 OK with success message, or 400 Bad Request
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(*models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	var req struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.CurrentPassword == "" || req.NewPassword == "" {
		http.Error(w, "Current password and new password are required", http.StatusBadRequest)
		return
	}

	// Change password
	if err := h.authService.ChangePassword(user.ID, req.CurrentPassword, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
}

// ==================== Admin User Management ====================

// CreateUser creates a new HR or Admin user (admin only).
// Expects: CreateHRUserRequest JSON body with all required fields
// Returns: 201 Created with created user (safe user without password), or 400 Bad Request
func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateHRUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" || req.Role == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Validate role
	if req.Role != "hr" && req.Role != "admin" {
		http.Error(w, "Invalid role. Must be 'hr' or 'admin'", http.StatusBadRequest)
		return
	}

	// Create user
	user, err := h.authService.CreateUser(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user.ToSafeUser())
}

// ListUsers returns a paginated list of all users (admin only).
// Query params: page (default: 1), limit (default: 20, max: 100)
// Response: {users, total, page, limit, pages}
func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page := 1
	limit := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := parseInt(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := parseInt(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	// Get users
	users, total, err := h.authService.ListUsers(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
		"pages": (total + int64(limit) - 1) / int64(limit),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateUser updates an existing user by ID (admin only).
// Expects: UpdateUserRequest JSON body (partial updates supported)
// Returns: 200 OK with updated user, 404 Not Found, or 400 Bad Request
func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := parseUint(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate role if provided
	if req.Role != nil && (*req.Role != "hr" && *req.Role != "admin") {
		http.Error(w, "Invalid role. Must be 'hr' or 'admin'", http.StatusBadRequest)
		return
	}

	// Update user
	user, err := h.authService.UpdateUser(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToSafeUser())
}

// DeleteUser permanently deletes a user by ID (admin only).
// Returns: 204 No Content, 404 Not Found, or 400 Bad Request
func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := parseUint(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.authService.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ==================== Middleware ====================

// authMiddleware validates JWT tokens and sets the authenticated user in the request context.
// It expects the Authorization header in the format "Bearer <token>".
// Returns: 401 Unauthorized if token is missing or invalid.
func (h *AuthHandler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid authorization header format. Expected 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		// Validate token and get user
		user, err := h.authService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Set user in context
		context.Set(r, "user", user)
		next.ServeHTTP(w, r)
	})
}

// hrMiddleware ensures the authenticated user has HR or Admin role permissions.
// Must be used after authMiddleware to ensure user is set in context.
// Returns: 403 Forbidden if user doesn't have HR access.
func (h *AuthHandler) hrMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := context.Get(r, "user").(*models.User)
		if !ok || !user.IsHR() {
			http.Error(w, "HR access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// adminMiddleware ensures the authenticated user has Admin role.
// Must be used after authMiddleware to ensure user is set in context.
// Returns: 403 Forbidden if user doesn't have admin access.
func (h *AuthHandler) adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := context.Get(r, "user").(*models.User)
		if !ok || !user.IsAdmin() {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ==================== Route Registration Helpers ====================

// registerEmployeeRoutes registers employee routes with authentication and HR middleware.
// All employee endpoints require a valid JWT token and HR/Admin role.
func (h *AuthHandler) registerEmployeeRoutes(r *mux.Router) {
	employeeHandler := NewEmployeeHandler()

	// All employee routes require authentication and HR role
	protected := r.PathPrefix("/api/employees").Subrouter()
	protected.Use(h.authMiddleware)
	protected.Use(h.hrMiddleware)

	protected.HandleFunc("", employeeHandler.GetEmployees).Methods("GET")
	protected.HandleFunc("", employeeHandler.CreateEmployee).Methods("POST")
	protected.HandleFunc("/{id}", employeeHandler.GetEmployee).Methods("GET")
	protected.HandleFunc("/{id}", employeeHandler.UpdateEmployee).Methods("PUT")
	protected.HandleFunc("/{id}", employeeHandler.DeleteEmployee).Methods("DELETE")
}

// registerAnalyticsRoutes registers analytics routes with authentication and HR middleware.
// All analytics endpoints require a valid JWT token and HR/Admin role.
func (h *AuthHandler) registerAnalyticsRoutes(r *mux.Router) {
	analyticsHandler := NewAnalyticsHandler()

	// All analytics routes require authentication and HR role
	protected := r.PathPrefix("/api/analytics").Subrouter()
	protected.Use(h.authMiddleware)
	protected.Use(h.hrMiddleware)

	protected.HandleFunc("/salary/by-country", analyticsHandler.GetSalaryByCountry).Methods("GET")
	protected.HandleFunc("/salary/by-job-title/{country}", analyticsHandler.GetSalaryByJobTitleInCountry).Methods("GET")
	protected.HandleFunc("/salary/department-insights", analyticsHandler.GetDepartmentInsights).Methods("GET")
}

// ==================== Utility Functions ====================

// parseInt parses a string to an integer.
// Returns the parsed integer or an error if parsing fails.
func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// parseUint parses a string to an unsigned integer.
// Returns the parsed uint or an error if parsing fails.
func parseUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	return uint(val), err
}
