package utils

import (
	"crypto/rand"
	"encoding/base64"
	"net/url"
	"strings"
)

// ValidateURL checks if a URL is valid
// Returns true if valid, false otherwise
func ValidateURL(urlStr string) bool {
	// Check if empty
	if strings.TrimSpace(urlStr) == "" {
		return false
	}

	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check if it has a scheme (http/https)
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// Check if it has a host
	if parsedURL.Host == "" {
		return false
	}

	return true
}

// GenerateShortCode generates a random short code
// Length parameter determines how many characters (default: 6)
func GenerateShortCode(length int) (string, error) {
	// Default to 6 characters if invalid length
	if length <= 0 {
		length = 6
	}

	// Generate random bytes
	// We need more bytes than the final length because base64 encoding expands the size
	numBytes := (length * 3) / 4
	if numBytes < length {
		numBytes = length
	}

	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode to base64 URL-safe format
	encoded := base64.URLEncoding.EncodeToString(randomBytes)

	// Remove special characters and truncate to desired length
	// Replace - and _ with alphanumeric characters
	encoded = strings.ReplaceAll(encoded, "-", "")
	encoded = strings.ReplaceAll(encoded, "_", "")
	encoded = strings.ReplaceAll(encoded, "=", "")

	// Truncate to desired length
	if len(encoded) > length {
		encoded = encoded[:length]
	}

	return encoded, nil
}
