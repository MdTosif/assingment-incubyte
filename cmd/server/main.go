// Package main provides the HTTP server entry point for the salary management application.
// It initializes the database, sets up routes, configures CORS, and serves static files.
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tofiquem/assingment/pkg/database"
	"github.com/tofiquem/assingment/pkg/handlers"
)

// ==================== Main Entry Point ====================

// main initializes and starts the HTTP server.
// It sets up database connection, API routes, CORS, and static file serving.
func main() {
	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	publicDir := os.Getenv("PUBLIC_DIR")
	if publicDir == "" {
		publicDir = "public"
	}
	publicDir, err := filepath.Abs(publicDir)
	if err != nil {
		log.Fatal(err)
	}

	// Create router with gorilla/mux
	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Authentication handler (registers all routes including protected ones)
	authHandler := handlers.NewAuthHandler(database.DB)
	authHandler.RegisterRoutes(r)

	// Serve static files
	r.PathPrefix("/").Handler(spaHandler(publicDir))

	// Configure CORS - allow all origins for development
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposedHeaders:   []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           86400,
		Debug:            false,
	})

	// Wrap router with CORS middleware
	handler := c.Handler(r)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + strings.TrimPrefix(p, ":")
	}
	log.Printf("listening on %s (serving static from %s)", addr, publicDir)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

// ==================== Static File Handler ====================

// spaHandler serves static assets from publicDir and falls back to index.html for client routing.
// It prevents directory traversal attacks and handles SPA client-side routing.
func spaHandler(publicDir string) http.Handler {
	fs := http.FileServer(http.Dir(publicDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		rel := strings.TrimPrefix(r.URL.Path, "/")
		if rel == "" {
			http.ServeFile(w, r, filepath.Join(publicDir, "index.html"))
			return
		}

		full := filepath.Join(publicDir, rel)
		safe, err := filepath.Rel(publicDir, full)
		if err != nil || strings.HasPrefix(safe, "..") {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		info, err := os.Stat(full)
		if err != nil || info.IsDir() {
			http.ServeFile(w, r, filepath.Join(publicDir, "index.html"))
			return
		}

		fs.ServeHTTP(w, r)
	})
}
