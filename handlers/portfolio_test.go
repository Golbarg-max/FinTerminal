package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddHoldingThenGet(t *testing.T) {
	h := newTestHandlers(t)

	body := strings.NewReader(`{"symbol":"AAPL","quantity":10,"cost_basis":150}`)
	rec := do(t, h.HandleHoldings, http.MethodPost, "/api/holdings", body)
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())

	rec = do(t, h.HandlePortfolio, http.MethodGet, "/api/portfolio", nil)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var got struct {
		Success bool              `json:"success"`
		Data    PortfolioResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))

	assert.True(t, got.Success)
	require.Len(t, got.Data.Holdings, 1)
	holding := got.Data.Holdings[0]
	assert.Equal(t, "AAPL", holding.Symbol)
	assert.Equal(t, 10.0, holding.Quantity)
	assert.Equal(t, 150.0, holding.CostPerShare)
	assert.Equal(t, 1.0, holding.CurrentPrice) // countingClient always returns 1
}
