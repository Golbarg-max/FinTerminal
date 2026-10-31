package handlers

import (
	"encoding/json"
	"errors"
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

func TestHandleHoldings(t *testing.T) {
	tests := []struct {
		name       string
		seed       bool // add AAPL before the request
		method     string
		target     string
		body       string
		wantStatus int
	}{
		{"add", false, http.MethodPost, "/api/holdings", `{"symbol":"AAPL","quantity":10,"cost_basis":150}`, http.StatusCreated},
		{"add malformed json", false, http.MethodPost, "/api/holdings", `{"symbol":`, http.StatusBadRequest},
		{"add invalid symbol", false, http.MethodPost, "/api/holdings", `{"symbol":"AA PL","quantity":10,"cost_basis":150}`, http.StatusBadRequest},
		{"add zero quantity", false, http.MethodPost, "/api/holdings", `{"symbol":"AAPL","quantity":0,"cost_basis":150}`, http.StatusBadRequest},
		{"remove via query", true, http.MethodDelete, "/api/holdings?symbol=AAPL", "", http.StatusOK},
		{"remove via body", true, http.MethodDelete, "/api/holdings", `{"symbol":"AAPL"}`, http.StatusOK},
		{"remove missing symbol", true, http.MethodDelete, "/api/holdings", "", http.StatusBadRequest},
		{"remove not held", false, http.MethodDelete, "/api/holdings?symbol=AAPL", "", http.StatusBadRequest},
		{"wrong method", false, http.MethodGet, "/api/holdings", "", http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHandlers(t)
			if tt.seed {
				seed := strings.NewReader(`{"symbol":"AAPL","quantity":10,"cost_basis":150}`)
				rec := do(t, h.HandleHoldings, http.MethodPost, "/api/holdings", seed)
				require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
			}

			rec := do(t, h.HandleHoldings, tt.method, tt.target, strings.NewReader(tt.body))
			assert.Equal(t, tt.wantStatus, rec.Code, rec.Body.String())
		})
	}
}

func TestRemoveHoldingThenGet(t *testing.T) {
	h := newTestHandlers(t)
	seed := strings.NewReader(`{"symbol":"AAPL","quantity":10,"cost_basis":150}`)
	rec := do(t, h.HandleHoldings, http.MethodPost, "/api/holdings", seed)
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())

	rec = do(t, h.HandleHoldings, http.MethodDelete, "/api/holdings?symbol=AAPL", nil)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	rec = do(t, h.HandlePortfolio, http.MethodGet, "/api/portfolio", nil)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var got struct {
		Data PortfolioResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	assert.Empty(t, got.Data.Holdings)
}

func TestGetPortfolioQuoteErrorKeepsHolding(t *testing.T) {
	h := newTestHandlers(t)
	seed := strings.NewReader(`{"symbol":"AAPL","quantity":10,"cost_basis":150}`)
	rec := do(t, h.HandleHoldings, http.MethodPost, "/api/holdings", seed)
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())

	client := &countingClient{err: errors.New("boom")}
	h.client = client

	rec = do(t, h.HandlePortfolio, http.MethodGet, "/api/portfolio", nil)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var got struct {
		Data PortfolioResponse `json:"data"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	require.Len(t, got.Data.Holdings, 1)
	assert.Equal(t, "AAPL", got.Data.Holdings[0].Symbol)
	assert.Equal(t, 0.0, got.Data.Holdings[0].CurrentPrice) // never updated
	assert.Equal(t, 1, client.calls)
}

func TestHandlePortfolioWrongMethod(t *testing.T) {
	h := newTestHandlers(t)
	rec := do(t, h.HandlePortfolio, http.MethodPost, "/api/portfolio", nil)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestHandleRefreshPrices(t *testing.T) {
	h := newTestHandlers(t)
	for _, sym := range []string{"AAPL", "MSFT"} {
		seed := strings.NewReader(`{"symbol":"` + sym + `","quantity":1,"cost_basis":100}`)
		rec := do(t, h.HandleHoldings, http.MethodPost, "/api/holdings", seed)
		require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
	}

	client := &countingClient{}
	h.client = client

	rec := do(t, h.HandleRefreshPrices, http.MethodPost, "/api/refresh", nil)
	assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	assert.Equal(t, 2, client.calls)

	rec = do(t, h.HandleRefreshPrices, http.MethodGet, "/api/refresh", nil)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}
