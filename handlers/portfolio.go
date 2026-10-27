package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handlers) HandlePortfolio(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getPortfolio(w, r)
	default:
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *Handlers) getPortfolio(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	p, err := h.loadPortfolio()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to load portfolio: "+err.Error())
		return
	}

	for _, symbol := range p.GetSymbols() {
		quote, err := h.client.GetQuote(symbol)
		if err != nil {
			continue
		}
		_ = p.UpdatePrice(symbol, quote.Price)
	}

	metrics := p.CalculatePerformance()
	totalCost, totalValue, totalGainLoss := p.CalculatePortfolioTotals()

	holdings := make([]HoldingResponse, 0, len(metrics))
	for _, m := range metrics {
		holdings = append(holdings, HoldingResponse{
			Symbol:          m.Symbol,
			Quantity:        m.Quantity,
			CostPerShare:    m.CostPerShare,
			CurrentPrice:    m.CurrentPrice,
			TotalCost:       m.TotalCost,
			CurrentValue:    m.CurrentValue,
			GainLoss:        m.GainLoss,
			GainLossPercent: m.GainLossPercent,
			LastUpdated:     p.Holdings[m.Symbol].LastUpdated,
		})
	}

	gainLossPercent := 0.0
	if totalCost > 0 {
		gainLossPercent = (totalGainLoss / totalCost) * 100
	}

	response := PortfolioResponse{
		Holdings:        holdings,
		TotalCost:       totalCost,
		TotalValue:      totalValue,
		TotalGainLoss:   totalGainLoss,
		GainLossPercent: gainLossPercent,
	}

	h.sendSuccess(w, http.StatusOK, response)
}

func (h *Handlers) HandleHoldings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.addHolding(w, r)
	case http.MethodDelete:
		h.removeHolding(w, r)
	default:
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *Handlers) addHolding(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Symbol    string  `json:"symbol"`
		Quantity  float64 `json:"quantity"`
		CostBasis float64 `json:"cost_basis"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if !validSymbol(req.Symbol) {
		h.sendError(w, http.StatusBadRequest, "Invalid symbol")
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()

	p, err := h.loadPortfolio()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to load portfolio")
		return
	}

	if err := p.AddHolding(req.Symbol, req.Quantity, req.CostBasis); err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := p.SaveToFile(); err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to save portfolio")
		return
	}

	h.sendSuccess(w, http.StatusCreated, map[string]string{"message": "Holding added successfully"})
}

func (h *Handlers) removeHolding(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")

	// The web UI sends the symbol as a JSON body on DELETE, so accept that too.
	if symbol == "" {
		var req struct {
			Symbol string `json:"symbol"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			symbol = req.Symbol
		}
	}
	if symbol == "" {
		h.sendError(w, http.StatusBadRequest, "Symbol is required")
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()

	p, err := h.loadPortfolio()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to load portfolio")
		return
	}

	if err := p.RemoveHolding(symbol); err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := p.SaveToFile(); err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to save portfolio")
		return
	}

	h.sendSuccess(w, http.StatusOK, map[string]string{"message": "Holding removed successfully"})
}

func (h *Handlers) HandleRefreshPrices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()

	p, err := h.loadPortfolio()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to load portfolio")
		return
	}

	for _, symbol := range p.GetSymbols() {
		quote, err := h.client.GetQuote(symbol)
		if err != nil {
			continue
		}
		_ = p.UpdatePrice(symbol, quote.Price)
	}

	if err := p.SaveToFile(); err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to save portfolio")
		return
	}

	h.sendSuccess(w, http.StatusOK, map[string]string{"message": "Prices refreshed successfully"})
}
