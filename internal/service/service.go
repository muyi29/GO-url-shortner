package service

import (
	"errors"
	"fmt"

	"github.com/muyi29/url-shortener/internal/models"
	"github.com/muyi29/url-shortener/internal/storage"
	"github.com/muyi29/url-shortener/internal/utils"
)

// URLService handles the business logic for URL shortening
type URLService struct {
	storage storage.Storage
	baseURL string // Base URL for our service (e.g., "http://localhost:8080")
}

// NewURLService creates a new URL service
func NewURLService(storage storage.Storage, baseURL string) *URLService {
	return &URLService{
		storage: storage,
		baseURL: baseURL,
	}
}

// ShortenURL shortens a long URL
func (s *URLService) ShortenURL(longURL string) (*models.ShortenResponse, error) {
	// Step 1: Validate the URL
	if !utils.ValidateURL(longURL) {
		return nil, errors.New("invalid URL format")
	}

	// Step 2: Check if URL already exists (deduplication)
	existingURL, err := s.storage.GetByOriginalURL(longURL)
	if err == nil {
		// URL already exists, return existing short code
		return &models.ShortenResponse{
			ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, existingURL.ShortCode),
			ShortCode:   existingURL.ShortCode,
			OriginalURL: existingURL.OriginalURL,
		}, nil
	}

	// Step 3: Generate a short code
	// Try up to 5 times in case of collisions
	var shortCode string
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		code, err := utils.GenerateShortCode(6)
		if err != nil {
			return nil, fmt.Errorf("failed to generate short code: %w", err)
		}

		// Check if this code already exists
		_, err = s.storage.GetByShortCode(code)
		if err != nil {
			// Code doesn't exist, we can use it
			shortCode = code
			break
		}
		// Code exists, try again
	}

	if shortCode == "" {
		return nil, errors.New("failed to generate unique short code after multiple attempts")
	}

	// Step 4: Create and save the URL
	url := &models.URL{
		ShortCode:   shortCode,
		OriginalURL: longURL,
	}

	err = s.storage.Save(url)
	if err != nil {
		return nil, fmt.Errorf("failed to save URL: %w", err)
	}

	// Step 5: Return the response
	return &models.ShortenResponse{
		ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, shortCode),
		ShortCode:   shortCode,
		OriginalURL: longURL,
	}, nil
}

// GetOriginalURL retrieves the original URL by short code
func (s *URLService) GetOriginalURL(shortCode string) (string, error) {
	url, err := s.storage.GetByShortCode(shortCode)
	if err != nil {
		return "", errors.New("short URL not found")
	}

	// Increment click counter (analytics)
	_ = s.storage.IncrementClicks(shortCode)

	return url.OriginalURL, nil
}
