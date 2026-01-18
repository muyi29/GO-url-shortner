package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new HTTP multiplexer (router)
	mux := http.NewServeMux()

	// Register our routes
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)

	// Server configuration
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Start the server
	fmt.Printf("🚀 Server starting on http://localhost:%s\n", port)
	fmt.Println("📝 Available endpoints:")
	fmt.Println("   GET  /        - Home page")
	fmt.Println("   GET  /health  - Health check")
	
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
	fmt.Fprintf(w, `{"message": "URL Shortener API", "version": "1.0.0"}`)
}

// healthHandler handles health check requests
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "healthy"}`)
}
