package handlers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (h *Handlers) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *Handlers) sendError(w http.ResponseWriter, statusCode int, message string) {
	h.sendJSON(w, statusCode, Response{
		Success: false,
		Error:   message,
	})
}

func (h *Handlers) sendSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	h.sendJSON(w, statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func (h *Handlers) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	mode := "production"
	if h.isDemoMode {
		mode = "demo"
	}

	response := map[string]string{
		"status":  "ok",
		"mode":    mode,
		"message": "FinTerminal API is running",
	}

	h.sendSuccess(w, http.StatusOK, response)
}

type PortfolioResponse struct {
	Holdings        []HoldingResponse `json:"holdings"`
	TotalCost       float64           `json:"total_cost"`
	TotalValue      float64           `json:"total_value"`
	TotalGainLoss   float64           `json:"total_gain_loss"`
	GainLossPercent float64           `json:"gain_loss_percent"`
}

type HoldingResponse struct {
	Symbol          string  `json:"symbol"`
	Quantity        float64 `json:"quantity"`
	CostPerShare    float64 `json:"cost_per_share"`
	CurrentPrice    float64 `json:"current_price"`
	TotalCost       float64 `json:"total_cost"`
	CurrentValue    float64 `json:"current_value"`
	GainLoss        float64 `json:"gain_loss"`
	GainLossPercent float64 `json:"gain_loss_percent"`
	LastUpdated     string  `json:"last_updated"`
}

type QuoteResponse struct {
	Symbol        string  `json:"symbol"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Volume        int64   `json:"volume"`
}
