// Package handler provides the Vercel serverless function entry point for the API.
// It wraps the HTTP handlers and configures routing for serverless deployment.
package handler

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tofiquem/assingment/pkg/database"
	"github.com/tofiquem/assingment/pkg/handlers"
)

func init() {
	database.InitDB()
}

// ==================== Middleware ====================

// stripAPIPrefix removes the /api prefix from request URLs.
// Vercel mounts this at /api, so we need to strip the prefix for routing to work correctly.
func stripAPIPrefix(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vercel mounts this at /api, so strip /api prefix from URL
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api")
		r.URL.RawPath = strings.TrimPrefix(r.URL.RawPath, "/api")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		next.ServeHTTP(w, r)
	})
}

// ==================== Handler Entry Point ====================

// Handler is the entry point for Vercel serverless functions.
// It sets up all routes, middleware, and handles incoming HTTP requests.
func Handler(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Auth handlers
	authHandler := handlers.NewAuthHandler(database.DB)
	authHandler.RegisterRoutes(router)

	// Analytics handlers
	analyticsHandler := handlers.NewAnalyticsHandler(database.DB)
	analyticsHandler.RegisterRoutes(router)

	// Employee handlers
	employeeHandler := handlers.NewEmployeeHandler(database.DB)
	employeeHandler.RegisterRoutes(router)

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposedHeaders:   []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           86400,
	})

	handler := c.Handler(stripAPIPrefix(router))
	handler.ServeHTTP(w, r)
}
