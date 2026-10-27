package handlers

import (
	"errors"
	"net/http"

	"FinTerminal/api"
)

func (h *Handlers) HandleQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	symbol := api.NormalizeSymbol(r.URL.Query().Get("symbol"))
	if symbol == "" {
		h.sendError(w, http.StatusBadRequest, "Symbol is required")
		return
	}
	if !validSymbol(symbol) {
		h.sendError(w, http.StatusBadRequest, "Invalid symbol")
		return
	}

	quote, err := h.client.GetQuote(symbol)
	if err != nil {
		status := http.StatusBadGateway
		switch {
		case errors.Is(err, api.ErrRateLimited):
			status = http.StatusTooManyRequests
		case errors.Is(err, api.ErrSymbolNotFound):
			status = http.StatusNotFound
		}
		h.sendError(w, status, "Failed to fetch quote: "+err.Error())
		return
	}

	response := QuoteResponse{
		Symbol:        quote.Symbol,
		Price:         quote.Price,
		Change:        quote.Change,
		ChangePercent: quote.ChangePercent,
		Volume:        quote.Volume,
	}

	h.sendSuccess(w, http.StatusOK, response)
}
