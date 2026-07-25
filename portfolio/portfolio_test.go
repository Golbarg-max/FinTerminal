package portfolio

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPortfolio(t *testing.T) {
	p := NewPortfolio()
	assert.NotNil(t, p)
	assert.Equal(t, 0, len(p.Holdings))
}

func TestAddHolding(t *testing.T) {
	p := NewPortfolio()

	err := p.AddHolding("AAPL", 10, 150.50)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(p.Holdings))

	holding, exists := p.Holdings["AAPL"]
	assert.True(t, exists)
	assert.Equal(t, "AAPL", holding.Symbol)
	assert.Equal(t, 10.0, holding.Quantity)
	assert.Equal(t, 150.50, holding.CostBasis)
}

func TestAddHoldingEmptySymbol(t *testing.T) {
	p := NewPortfolio()
	err := p.AddHolding("", 10, 150.50)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "symbol cannot be empty")
}

func TestAddHoldingNegativeQuantity(t *testing.T) {
	p := NewPortfolio()
	err := p.AddHolding("AAPL", -10, 150.50)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quantity must be positive")
}

func TestAddHoldingNegativeCostBasis(t *testing.T) {
	p := NewPortfolio()
	err := p.AddHolding("AAPL", 10, -150.50)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cost basis cannot be negative")
}

func TestAddMultipleHoldings(t *testing.T) {
	p := NewPortfolio()

	holdings := []struct {
		symbol    string
		quantity  float64
		costBasis float64
	}{
		{"AAPL", 10, 150.50},
		{"MSFT", 5, 300.00},
		{"GOOGL", 2, 2800.00},
	}

	for _, h := range holdings {
		err := p.AddHolding(h.symbol, h.quantity, h.costBasis)
		assert.NoError(t, err)
	}

	assert.Equal(t, 3, len(p.Holdings))
}

func TestUpdatePrice(t *testing.T) {
	p := NewPortfolio()
	err := p.AddHolding("AAPL", 10, 150.50)
	require.NoError(t, err)

	err = p.UpdatePrice("AAPL", 189.95)
	assert.NoError(t, err)
	assert.Equal(t, 189.95, p.Holdings["AAPL"].CurrentPrice)
	assert.NotEmpty(t, p.Holdings["AAPL"].LastUpdated)
}

func TestUpdatePriceNonexistentSymbol(t *testing.T) {
	p := NewPortfolio()
	err := p.UpdatePrice("NONEXISTENT", 100.00)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestUpdatePriceNegativePrice(t *testing.T) {
	p := NewPortfolio()
	err := p.AddHolding("AAPL", 10, 150.50)
	require.NoError(t, err)

	err = p.UpdatePrice("AAPL", -100.00)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "price cannot be negative")
}

func TestRemoveHolding(t *testing.T) {
	p := NewPortfolio()
	err := p.AddHolding("AAPL", 10, 150.50)
	require.NoError(t, err)
	err = p.AddHolding("MSFT", 5, 300.00)
	require.NoError(t, err)

	err = p.RemoveHolding("AAPL")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(p.Holdings))

	_, exists := p.Holdings["AAPL"]
	assert.False(t, exists)
}

func TestRemoveHoldingNonexistent(t *testing.T) {
	p := NewPortfolio()
	err := p.RemoveHolding("NONEXISTENT")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRemoveHoldingEmptySymbol(t *testing.T) {
	p := NewPortfolio()
	err := p.RemoveHolding("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "symbol cannot be empty")
}

func TestGetSymbols(t *testing.T) {
	p := NewPortfolio()
	err := p.AddHolding("AAPL", 10, 150.50)
	require.NoError(t, err)
	err = p.AddHolding("MSFT", 5, 300.00)
	require.NoError(t, err)
	err = p.AddHolding("GOOGL", 2, 2800.00)
	require.NoError(t, err)

	symbols := p.GetSymbols()
	assert.Equal(t, 3, len(symbols))
	assert.ElementsMatch(t, []string{"AAPL", "MSFT", "GOOGL"}, symbols)
}

func TestGetSymbolsEmpty(t *testing.T) {
	p := NewPortfolio()
	symbols := p.GetSymbols()
	assert.Empty(t, symbols)
}

func TestGetHolding(t *testing.T) {
	p := NewPortfolio()
	err := p.AddHolding("AAPL", 10, 150.50)
	require.NoError(t, err)

	holding, err := p.GetHolding("AAPL")
	assert.NoError(t, err)
	assert.NotNil(t, holding)
	assert.Equal(t, "AAPL", holding.Symbol)
	assert.Equal(t, 10.0, holding.Quantity)
}

func TestGetHoldingNonexistent(t *testing.T) {
	p := NewPortfolio()
	holding, err := p.GetHolding("NONEXISTENT")
	assert.Error(t, err)
	assert.Nil(t, holding)
	assert.Contains(t, err.Error(), "not found")
}

func TestSaveAndLoadPortfolio(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "portfolio.json")

	// Create and save
	p1 := NewPortfolio()
	err := p1.AddHolding("AAPL", 10, 150.50)
	require.NoError(t, err)
	err = p1.AddHolding("MSFT", 5, 300.00)
	require.NoError(t, err)
	p1.filePath = filePath

	err = p1.SaveToFile()
	assert.NoError(t, err)
	assert.FileExists(t, filePath)

	// Load into new portfolio
	p2 := NewPortfolio()
	err = p2.LoadFromFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(p2.Holdings))
	assert.Equal(t, 10.0, p2.Holdings["AAPL"].Quantity)
	assert.Equal(t, 5.0, p2.Holdings["MSFT"].Quantity)
}

func TestLoadNonexistentFile(t *testing.T) {
	p := NewPortfolio()
	nonexistentPath := filepath.Join(t.TempDir(), "nonexistent.json")

	err := p.LoadFromFile(nonexistentPath)
	assert.NoError(t, err)
	assert.FileExists(t, nonexistentPath)
	assert.Equal(t, 0, len(p.Holdings))
}

func TestSaveWithoutPath(t *testing.T) {
	p := NewPortfolio()
	err := p.AddHolding("AAPL", 10, 150.50)
	require.NoError(t, err)

	err = p.SaveToFile()
	assert.Error(t, err)
	assert.Equal(t, "no file path set", err.Error())
}

func TestLoadWithEmptyPath(t *testing.T) {
	p := NewPortfolio()
	err := p.LoadFromFile("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file path cannot be empty")
}

func TestPortfolioRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "roundtrip.json")

	// Create portfolio
	p1 := NewPortfolio()
	_ = p1.AddHolding("AAPL", 10, 150.50)
	_ = p1.AddHolding("MSFT", 5, 300.00)
	_ = p1.AddHolding("GOOGL", 2, 2800.00)
	p1.filePath = filePath
	_ = p1.SaveToFile()

	// Update prices
	_ = p1.UpdatePrice("AAPL", 189.95)
	_ = p1.UpdatePrice("MSFT", 334.20)
	_ = p1.UpdatePrice("GOOGL", 2950.00)
	_ = p1.SaveToFile()

	// Load fresh
	p2 := NewPortfolio()
	err := p2.LoadFromFile(filePath)
	assert.NoError(t, err)

	// Verify all data persisted
	assert.Equal(t, p1.Holdings["AAPL"].CurrentPrice, p2.Holdings["AAPL"].CurrentPrice)
	assert.Equal(t, p1.Holdings["MSFT"].CurrentPrice, p2.Holdings["MSFT"].CurrentPrice)
	assert.Equal(t, p1.Holdings["GOOGL"].CurrentPrice, p2.Holdings["GOOGL"].CurrentPrice)
}