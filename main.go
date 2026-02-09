package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/muyi29/url-shortener/internal/handlers"
	"github.com/muyi29/url-shortener/internal/service"
	"github.com/muyi29/url-shortener/internal/storage"
)

func main() {
	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get base URL from environment variable
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%s", port)
	}

	// Initialize storage layer
	store := storage.NewInMemoryStorage()

	// Initialize service layer
	urlService := service.NewURLService(store, baseURL)

	// Initialize handlers
	urlHandler := handlers.NewURLHandler(urlService)

	// Create a new HTTP multiplexer (router)
	mux := http.NewServeMux()

	// Register our routes
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/shorten", urlHandler.ShortenURL)
	
	// Note: This is a catch-all for redirects
	// We'll handle the routing logic inside the handler
	// In the next phase, we'll use a proper router like Gin
	mux.HandleFunc("/r/", urlHandler.RedirectURL)

	// Server configuration
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Start the server
	fmt.Printf("🚀 Server starting on %s\n", baseURL)
	fmt.Println("📝 Available endpoints:")
	fmt.Println("   GET  /             - Home page")
	fmt.Println("   GET  /health       - Health check")
	fmt.Println("   POST /api/shorten  - Shorten a URL")
	fmt.Println("   GET  /r/:code      - Redirect to original URL")
	
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// homeHandler handles requests to the root path
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Only handle exact "/" path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "URL Shortener API", "version": "1.0.0", "endpoints": {"/api/shorten": "POST - Shorten a URL", "/r/:code": "GET - Redirect to original URL"}}`)
}

// healthHandler handles health check requests
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "healthy"}`)
}
