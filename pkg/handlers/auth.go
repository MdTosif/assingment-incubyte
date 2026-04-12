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

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(db),
	}
}

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

// Login handles user authentication
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

// GetMe returns the current authenticated user
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(*models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToSafeUser())
}

// Logout handles user logout (client-side token removal)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a stateless JWT system, logout is handled client-side
	// We can implement token blacklisting if needed
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// ChangePassword handles password change for current user
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

// CreateUser handles user creation (admin only)
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

// ListUsers returns a paginated list of users (admin only)
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

// UpdateUser handles user updates (admin only)
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

// DeleteUser handles user deletion (admin only)
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

// authMiddleware validates JWT tokens and sets user in context
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

// hrMiddleware ensures user has HR role or higher
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

// adminMiddleware ensures user has admin role
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

// registerEmployeeRoutes registers employee routes with authentication middleware
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

// registerAnalyticsRoutes registers analytics routes with authentication middleware
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

// Helper functions
func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func parseUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	return uint(val), err
}
