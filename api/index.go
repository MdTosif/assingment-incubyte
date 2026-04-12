package handler

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tofiquem/assingment/internal/database"
	"github.com/tofiquem/assingment/internal/handlers"
)

func init() {
	database.InitDB()
}

// stripAPIPrefix removes /api from the request URL before passing to router
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

// Handler is the entry point for Vercel serverless functions
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
	analyticsHandler := handlers.NewAnalyticsHandler()
	analyticsHandler.RegisterRoutes(router)

	// Employee handlers
	employeeHandler := handlers.NewEmployeeHandler()
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
