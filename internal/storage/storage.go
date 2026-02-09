package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/muyi29/url-shortener/internal/models"
)

// Storage interface defines the methods our storage layer must implement
// This allows us to swap implementations later (e.g., PostgreSQL)
type Storage interface {
	Save(url *models.URL) error
	GetByShortCode(shortCode string) (*models.URL, error)
	GetByOriginalURL(originalURL string) (*models.URL, error)
	IncrementClicks(shortCode string) error
}

// InMemoryStorage implements Storage using in-memory maps
type InMemoryStorage struct {
	urls      map[string]*models.URL // shortCode -> URL
	urlsByOrig map[string]*models.URL // originalURL -> URL (for deduplication)
	mu        sync.RWMutex           // Mutex for thread-safe access
	nextID    int                    // Auto-incrementing ID
}

// NewInMemoryStorage creates a new in-memory storage instance
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		urls:      make(map[string]*models.URL),
		urlsByOrig: make(map[string]*models.URL),
		nextID:    1,
	}
}

// Save stores a new URL
func (s *InMemoryStorage) Save(url *models.URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if short code already exists (collision)
	if _, exists := s.urls[url.ShortCode]; exists {
		return errors.New("short code already exists")
	}

	// Assign ID and timestamp
	url.ID = s.nextID
	url.CreatedAt = time.Now()
	url.Clicks = 0

	// Store in both maps
	s.urls[url.ShortCode] = url
	s.urlsByOrig[url.OriginalURL] = url

	s.nextID++
	return nil
}

// GetByShortCode retrieves a URL by its short code
func (s *InMemoryStorage) GetByShortCode(shortCode string) (*models.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.urls[shortCode]
	if !exists {
		return nil, errors.New("URL not found")
	}

	return url, nil
}

// GetByOriginalURL retrieves a URL by its original URL
// This is used for deduplication - if the same URL is shortened twice,
// we return the existing short code
func (s *InMemoryStorage) GetByOriginalURL(originalURL string) (*models.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.urlsByOrig[originalURL]
	if !exists {
		return nil, errors.New("URL not found")
	}

	return url, nil
}

// IncrementClicks increments the click counter for a URL
func (s *InMemoryStorage) IncrementClicks(shortCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	url, exists := s.urls[shortCode]
	if !exists {
		return errors.New("URL not found")
	}

	url.Clicks++
	return nil
}
