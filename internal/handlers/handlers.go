package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/muyi29/url-shortener/internal/models"
	"github.com/muyi29/url-shortener/internal/service"
)

// URLHandler handles HTTP requests for URL operations
type URLHandler struct {
	service *service.URLService
}

// NewURLHandler creates a new URL handler
func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

// ShortenURL handles POST /api/shorten
func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req models.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate that URL is provided
	if strings.TrimSpace(req.URL) == "" {
		sendErrorResponse(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Call service to shorten URL
	response, err := h.service.ShortenURL(req.URL)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send success response
	sendJSONResponse(w, response, http.StatusCreated)
}

// RedirectURL handles GET /:code
func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract short code from URL path
	// Path will be like "/r/abc123"
	shortCode := strings.TrimPrefix(r.URL.Path, "/r/")
	
	// Validate short code
	if shortCode == "" || shortCode == "/" || shortCode == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	// Get original URL
	originalURL, err := h.service.GetOriginalURL(shortCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Redirect to original URL
	// HTTP 302 = temporary redirect (allows us to track clicks)
	// HTTP 301 = permanent redirect (browsers cache it, can't track clicks)
	http.Redirect(w, r, originalURL, http.StatusFound)
}

// Helper function to send JSON response
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Helper function to send error response
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: message,
	})
}
