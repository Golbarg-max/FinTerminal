package portfolio

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculatePerformance(t *testing.T) {
	p := NewPortfolio()
	err := p.AddHolding("AAPL", 10, 150.00)
	require.NoError(t, err)
	p.Holdings["AAPL"].CurrentPrice = 200.00

	metrics := p.CalculatePerformance()
	assert.Equal(t, 1, len(metrics))

	m := metrics[0]
	assert.Equal(t, "AAPL", m.Symbol)
	assert.Equal(t, 1500.00, m.TotalCost)
	assert.Equal(t, 2000.00, m.CurrentValue)
	assert.Equal(t, 500.00, m.GainLoss)
	assert.InDelta(t, 33.33, m.GainLossPercent, 0.5)
}

func TestCalculatePerformanceNegativeReturn(t *testing.T) {
	p := NewPortfolio()
	err := p.AddHolding("AAPL", 10, 200.00)
	require.NoError(t, err)
	p.Holdings["AAPL"].CurrentPrice = 150.00

	metrics := p.CalculatePerformance()
	m := metrics[0]

	assert.Negative(t, m.GainLoss)
	assert.Negative(t, m.GainLossPercent)
	assert.Equal(t, -500.00, m.GainLoss)
	assert.InDelta(t, -25.00, m.GainLossPercent, 0.5)
}

func TestCalculatePerformanceZeroPrice(t *testing.T) {
	p := NewPortfolio()
	err := p.AddHolding("AAPL", 10, 150.00)
	require.NoError(t, err)
	// CurrentPrice is 0 by default

	metrics := p.CalculatePerformance()
	m := metrics[0]

	assert.Equal(t, -1500.00, m.GainLoss)
	assert.Equal(t, -100.00, m.GainLossPercent)
}

func TestCalculatePerformanceBreakEven(t *testing.T) {
	p := NewPortfolio()
	err := p.AddHolding("AAPL", 10, 150.00)
	require.NoError(t, err)
	p.Holdings["AAPL"].CurrentPrice = 150.00

	metrics := p.CalculatePerformance()
	m := metrics[0]

	assert.Equal(t, 0.0, m.GainLoss)
	assert.Equal(t, 0.0, m.GainLossPercent)
}

func TestCalculatePortfolioTotals(t *testing.T) {
	p := NewPortfolio()
	_ = p.AddHolding("AAPL", 10, 150.00)
	p.Holdings["AAPL"].CurrentPrice = 200.00

	_ = p.AddHolding("MSFT", 5, 300.00)
	p.Holdings["MSFT"].CurrentPrice = 350.00

	totalCost, totalValue, totalGainLoss := p.CalculatePortfolioTotals()

	assert.Equal(t, 3000.00, totalCost)
	assert.Equal(t, 3750.00, totalValue)
	assert.Equal(t, 750.00, totalGainLoss)
}

func TestCalculatePortfolioTotalsEmpty(t *testing.T) {
	p := NewPortfolio()

	totalCost, totalValue, totalGainLoss := p.CalculatePortfolioTotals()

	assert.Zero(t, totalCost)
	assert.Zero(t, totalValue)
	assert.Zero(t, totalGainLoss)
}

func TestCalculatePortfolioTotalsSingleHolding(t *testing.T) {
	p := NewPortfolio()
	_ = p.AddHolding("AAPL", 10, 150.00)
	p.Holdings["AAPL"].CurrentPrice = 200.00

	totalCost, totalValue, totalGainLoss := p.CalculatePortfolioTotals()

	assert.Equal(t, 1500.00, totalCost)
	assert.Equal(t, 2000.00, totalValue)
	assert.Equal(t, 500.00, totalGainLoss)
}

func TestPerformanceMetricsPositiveReturn(t *testing.T) {
	p := NewPortfolio()
	_ = p.AddHolding("TEST", 10, 100)
	p.Holdings["TEST"].CurrentPrice = 150

	metrics := p.CalculatePerformance()
	m := metrics[0]

	assert.Equal(t, 10.0, m.Quantity)
	assert.Equal(t, 100.0, m.CostPerShare)
	assert.Equal(t, 150.0, m.CurrentPrice)
	assert.Equal(t, 1000.0, m.TotalCost)
	assert.Equal(t, 1500.0, m.CurrentValue)
	assert.Equal(t, 500.0, m.GainLoss)
	assert.InDelta(t, 50.0, m.GainLossPercent, 0.1)
}

func TestPerformanceMetricsMultipleHoldings(t *testing.T) {
	p := NewPortfolio()
	_ = p.AddHolding("AAPL", 10, 100)
	p.Holdings["AAPL"].CurrentPrice = 120

	_ = p.AddHolding("MSFT", 20, 50)
	p.Holdings["MSFT"].CurrentPrice = 55

	metrics := p.CalculatePerformance()
	assert.Equal(t, 2, len(metrics))

	// Verify both metrics are calculated
	gainLossTotal := 0.0
	for _, m := range metrics {
		assert.NotZero(t, m.TotalCost)
		assert.NotZero(t, m.CurrentValue)
		gainLossTotal += m.GainLoss
	}

	assert.Positive(t, gainLossTotal)
}

func TestPerformanceWithFractionalShares(t *testing.T) {
	p := NewPortfolio()
	_ = p.AddHolding("AAPL", 2.5, 150.00)
	p.Holdings["AAPL"].CurrentPrice = 200.00

	metrics := p.CalculatePerformance()
	m := metrics[0]

	assert.Equal(t, 2.5, m.Quantity)
	assert.Equal(t, 375.00, m.TotalCost)
	assert.Equal(t, 500.00, m.CurrentValue)
	assert.Equal(t, 125.00, m.GainLoss)
}

func TestPerformanceMetricsFieldAssignment(t *testing.T) {
	p := NewPortfolio()
	_ = p.AddHolding("AAPL", 5, 200.00)
	p.Holdings["AAPL"].CurrentPrice = 250.00

	metrics := p.CalculatePerformance()
	m := metrics[0]

	// Verify all fields are properly assigned
	assert.Equal(t, "AAPL", m.Symbol)
	assert.Equal(t, 5.0, m.Quantity)
	assert.Equal(t, 200.0, m.CostPerShare)
	assert.Equal(t, 250.0, m.CurrentPrice)
	assert.Equal(t, 1000.0, m.TotalCost)
	assert.Equal(t, 1250.0, m.CurrentValue)
	assert.Equal(t, 250.0, m.GainLoss)
	assert.Equal(t, m.GainLossPercent, m.ReturnOnInvestment)
}

// Holdings live in a map, whose iteration order Go randomizes, so the result
// must be sorted explicitly. Run with -count=50 to make a regression show up.
func TestCalculatePerformanceSortedBySymbol(t *testing.T) {
	p := NewPortfolio()
	for _, symbol := range []string{"TSLA", "AAPL", "NVDA", "GOOGL", "MSFT"} {
		require.NoError(t, p.AddHolding(symbol, 1, 100.00))
	}

	metrics := p.CalculatePerformance()

	var symbols []string
	for _, m := range metrics {
		symbols = append(symbols, m.Symbol)
	}
	assert.Equal(t, []string{"AAPL", "GOOGL", "MSFT", "NVDA", "TSLA"}, symbols)
}
