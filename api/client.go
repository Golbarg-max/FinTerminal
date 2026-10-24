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

// NormalizeSymbol puts a ticker in its canonical form (" aapl " -> "AAPL") so
// the same stock always maps to the same key in portfolios and caches.
func NormalizeSymbol(symbol string) string {
	return strings.ToUpper(strings.TrimSpace(symbol))
}

type CachedClient struct {
	inner Client        // the real client it delegate
	ttl   time.Duration // how long a quote stays fres
	mu    sync.Mutex
	cache map[string]cacheEntry
}

type cacheEntry struct {
	quote     QuoteData
	fetchedAt time.Time
}

func NewCachedClient(inner Client, ttl time.Duration) *CachedClient {
	c := make(map[string]cacheEntry)
	return &CachedClient{inner: inner, ttl: ttl, cache: c}
}
func (c *CachedClient) GetQuote(symbol string) (*QuoteData, error) {
	symbol = NormalizeSymbol(symbol)
	c.mu.Lock()
	entry, ok := c.cache[symbol]
	c.mu.Unlock()
	if ok && time.Since(entry.fetchedAt) < c.ttl {
		return &entry.quote, nil
	}
	q, err := c.inner.GetQuote(symbol)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	c.cache[symbol] = cacheEntry{quote: *q, fetchedAt: time.Now()}
	c.mu.Unlock()
	return q, nil
}
