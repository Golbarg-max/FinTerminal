package portfolio

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"FinTerminal/api"
)

type Holding struct {
	Symbol   string  `json:"symbol"`
	Quantity float64 `json:"quantity"`
	//average cost per share
	CostBasis    float64 `json:"cost_basis"`
	CurrentPrice float64 `json:"current_price"`
	LastUpdated  string  `json:"last_updated"`
}

// Validate checks if a holding has valid values
func (h *Holding) Validate() error {
	if h.Symbol == "" {
		return fmt.Errorf("symbol cannot be empty")
	}
	if h.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive, got %f", h.Quantity)
	}
	if h.CostBasis < 0 {
		return fmt.Errorf("cost basis cannot be negative, got %f", h.CostBasis)
	}
	return nil
}

type Portfolio struct {
	Holdings map[string]*Holding `json:"holdings"`
	filePath string
}

func NewPortfolio() *Portfolio {
	return &Portfolio{
		Holdings: make(map[string]*Holding),
	}
}

func (p *Portfolio) LoadFromFile(filepath string) error {
	if filepath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	p.filePath = filepath

	// Create file if it doesn't exist
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := p.SaveToFile(); err != nil {
			return fmt.Errorf("failed to create portfolio file: %w", err)
		}
		return nil
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read portfolio file: %w", err)
	}

	var temp struct {
		Holdings map[string]*Holding `json:"holdings"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to parse portfolio JSON: %w", err)
	}

	p.Holdings = temp.Holdings
	if p.Holdings == nil {
		p.Holdings = make(map[string]*Holding)
	}

	return nil
}

func (p *Portfolio) SaveToFile() error {
	if p.filePath == "" {
		return fmt.Errorf("no file path set")
	}

	dir := filepath.Dir(p.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data := struct {
		Holdings map[string]*Holding `json:"holdings"`
	}{
		Holdings: p.Holdings,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal portfolio: %w", err)
	}

	if err := os.WriteFile(p.filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write portfolio file: %w", err)
	}

	return nil
}

func (p *Portfolio) AddHolding(symbol string, quantity float64, costBasis float64) error {
	symbol = api.NormalizeSymbol(symbol)
	if symbol == "" {
		return fmt.Errorf("symbol cannot be empty")
	}
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive, got %f", quantity)
	}
	if costBasis < 0 {
		return fmt.Errorf("cost basis cannot be negative, got %f", costBasis)
	}

	holding := &Holding{
		Symbol:    symbol,
		Quantity:  quantity,
		CostBasis: costBasis,
	}

	if err := holding.Validate(); err != nil {
		return fmt.Errorf("invalid holding: %w", err)
	}

	if existing, ok := p.Holdings[symbol]; ok {
		avgCost := (existing.Quantity*existing.CostBasis + quantity*costBasis) / (existing.Quantity + quantity)
		existing.CostBasis = avgCost
		existing.Quantity += quantity

		return nil
	}

	p.Holdings[symbol] = holding
	return nil
}

func (p *Portfolio) UpdatePrice(symbol string, currentPrice float64) error {
	symbol = api.NormalizeSymbol(symbol)
	if symbol == "" {
		return fmt.Errorf("symbol cannot be empty")
	}
	if currentPrice < 0 {
		return fmt.Errorf("price cannot be negative, got %f", currentPrice)
	}

	holding, exists := p.Holdings[symbol]
	if !exists {
		return fmt.Errorf("symbol %s not found in portfolio", symbol)
	}

	holding.CurrentPrice = currentPrice
	holding.LastUpdated = time.Now().Format("2006-01-02 15:04:05")

	return nil
}

func (p *Portfolio) RemoveHolding(symbol string) error {
	symbol = api.NormalizeSymbol(symbol)
	if symbol == "" {
		return fmt.Errorf("symbol cannot be empty")
	}

	if _, exists := p.Holdings[symbol]; !exists {
		return fmt.Errorf("symbol %s not found in portfolio", symbol)
	}

	delete(p.Holdings, symbol)
	return nil
}

func (p *Portfolio) GetSymbols() []string {
	symbols := make([]string, 0, len(p.Holdings))
	for symbol := range p.Holdings {
		symbols = append(symbols, symbol)
	}
	return symbols
}

func (p *Portfolio) GetHolding(symbol string) (*Holding, error) {
	symbol = api.NormalizeSymbol(symbol)
	if symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}

	holding, exists := p.Holdings[symbol]
	if !exists {
		return nil, fmt.Errorf("symbol %s not found in portfolio", symbol)
	}

	return holding, nil
}
