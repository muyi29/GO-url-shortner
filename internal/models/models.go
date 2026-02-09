package models

import "time"

// URL represents a shortened URL in our system
// This is what we'll store in our database/storage
type URL struct {
	ID          int       `json:"id"`           // Unique identifier
	ShortCode   string    `json:"short_code"`   // The short code (e.g., "abc123")
	OriginalURL string    `json:"original_url"` // The original long URL
	CreatedAt   time.Time `json:"created_at"`   // When it was created
	Clicks      int       `json:"clicks"`       // Number of times accessed (for analytics)
}

// ShortenRequest represents the incoming request to shorten a URL
// This is what the client sends us
type ShortenRequest struct {
	URL string `json:"url"` // The long URL to shorten
}

// ShortenResponse represents the response we send back
// This is what the client receives
type ShortenResponse struct {
	ShortURL    string `json:"short_url"`    // Full shortened URL (e.g., "http://localhost:8080/abc123")
	ShortCode   string `json:"short_code"`   // Just the code (e.g., "abc123")
	OriginalURL string `json:"original_url"` // The original URL (for confirmation)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"` // Error message
}
