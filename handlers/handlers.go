package handlers

import (
	"FinTerminal/api"
	"FinTerminal/portfolio"
	"os"
	"path/filepath"
)

type Handlers struct {
	apiKey        string
	client        *api.AlphaVantageClient
	portfolioPath string
}

func NewHandlers(apiKey string) *Handlers {
	portfolioPath := os.Getenv("FINTERMINAL_PORTFOLIO_PATH")
	if portfolioPath == "" {
		homeDir, _ := os.UserHomeDir()
		portfolioPath = filepath.Join(homeDir, ".finterminal", "portfolio.json")
	}

	return &Handlers{
		apiKey:        apiKey,
		client:        api.NewAlphaVantageClient(apiKey),
		portfolioPath: portfolioPath,
	}
}

func (h *Handlers) loadPortfolio() (*portfolio.Portfolio, error) {
	p := portfolio.NewPortfolio()
	if err := p.LoadFromFile(h.portfolioPath); err != nil {
		return nil, err
	}
	return p, nil
}