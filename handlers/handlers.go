package handlers

import (
	"FinTerminal/api"
	"FinTerminal/portfolio"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Handlers struct {
	client        APIClient
	portfolioPath string
	isDemoMode    bool
	mu            sync.Mutex
}

type APIClient interface {
	GetQuote(symbol string) (*api.QuoteData, error)
}

func NewHandlers(apiKey string) *Handlers {

	client, isDemoMode := api.NewClient(apiKey)
	if !isDemoMode {
		client = api.NewCachedClient(client, 5 * time.Minute)
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
