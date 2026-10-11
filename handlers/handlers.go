package handlers

import (
	"FinTerminal/api"
	"FinTerminal/portfolio"
	"os"
	"path/filepath"
	"strings"
)

type Handlers struct {
	client        APIClient
	portfolioPath string
	isDemoMode    bool
}

type APIClient interface {
	GetQuote(symbol string) (*api.QuoteData, error)
}

func NewHandlers(apiKey string) *Handlers {
	var client APIClient
	isDemoMode := false

	// Check if API key is provided and not "demo"
	if apiKey == "" || strings.ToLower(apiKey) == "demo" {
		// Use mock client
		client = api.NewMockClient()
		isDemoMode = true
	} else {
		// Use real Alpha Vantage client
		client = api.NewAlphaVantageClient(apiKey)
	}

	portfolioPath := os.Getenv("FINTERMINAL_PORTFOLIO_PATH")
	if portfolioPath == "" {
		homeDir, _ := os.UserHomeDir()
		portfolioPath = filepath.Join(homeDir, ".finterminal", "portfolio.json")
	}

	return &Handlers{
		client:        client,
		portfolioPath: portfolioPath,
		isDemoMode:    isDemoMode,
	}
}

// IsDemoMode reports whether the mock quote client is in use.
func (h *Handlers) IsDemoMode() bool { return h.isDemoMode }

// PortfolioPath is the file the portfolio is read from and written to.
func (h *Handlers) PortfolioPath() string { return h.portfolioPath }

func (h *Handlers) loadPortfolio() (*portfolio.Portfolio, error) {
	p := portfolio.NewPortfolio()
	if err := p.LoadFromFile(h.portfolioPath); err != nil {
		return nil, err
	}
	return p, nil
}
