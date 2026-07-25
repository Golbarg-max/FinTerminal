package portfolio

type PerformanceMetrics struct {
	Symbol              string
	Quantity            float64
	CostPerShare        float64
	CurrentPrice        float64
	TotalCost           float64
	CurrentValue        float64
	GainLoss            float64
	GainLossPercent     float64
	ReturnOnInvestment  float64
}

func (p *Portfolio) CalculatePerformance() []PerformanceMetrics {
	var metrics []PerformanceMetrics

	for _, holding := range p.Holdings {
		totalCost := holding.Quantity * holding.CostBasis
		currentValue := holding.Quantity * holding.CurrentPrice
		gainLoss := currentValue - totalCost
		gainLossPercent := 0.0
		if totalCost > 0 {
			gainLossPercent = (gainLoss / totalCost) * 100
		}

		metric := PerformanceMetrics{
			Symbol:             holding.Symbol,
			Quantity:           holding.Quantity,
			CostPerShare:       holding.CostBasis,
			CurrentPrice:       holding.CurrentPrice,
			TotalCost:          totalCost,
			CurrentValue:       currentValue,
			GainLoss:           gainLoss,
			GainLossPercent:    gainLossPercent,
			ReturnOnInvestment: gainLossPercent,
		}

		metrics = append(metrics, metric)
	}

	return metrics
}

func (p *Portfolio) CalculatePortfolioTotals() (totalCost, totalValue, totalGainLoss float64) {
	metrics := p.CalculatePerformance()

	for _, m := range metrics {
		totalCost += m.TotalCost
		totalValue += m.CurrentValue
	}

	totalGainLoss = totalValue - totalCost

	return
}
