package api

import (
	"os"
	"strings"
	"sync"
	"time"
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

type CachedClient struct {
	inner Client          // the real client it delegate
	ttl   time.Duration   // how long a quote stays fres
	mu    sync.Mutex
	cache map[string]cacheEntry
}

type cacheEntry struct {
	quote     QuoteData
	fetchedAt time.Time
}

func NewCachedClient(inner Client, ttl time.Duration) *CachedClient
func (c *CachedClient) GetQuote(symbol string) (*QuoteData, error)