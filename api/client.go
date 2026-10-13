package api

import (
	"os"
	"strings"
)

type Client interface {
	GetQuote(symbol string) (*QuoteData, error)
}

// NewClient picks the quote client to use. An empty apiKey falls back to the
// ALPHA_VANTAGE_API_KEY environment variable (which main loads from .env).
// The bool reports whether the mock client (demo mode) was chosen.
func NewClient(apiKey string) (Client, bool) {
	if apiKey == "" {
		apiKey = os.Getenv("ALPHA_VANTAGE_API_KEY")
	}

	// Check if API key is provided and not "demo"
	if apiKey == "" || strings.ToLower(apiKey) == "demo" {
		// Use mock client
		return NewMockClient(), true
	}
	// Use real Alpha Vantage client
	return NewAlphaVantageClient(apiKey), false
}
