package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tofiquem/assingment/internal/database"
	"github.com/tofiquem/assingment/internal/handlers"
)

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

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
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

// spaHandler serves static assets from publicDir and falls back to index.html for client routing.
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
