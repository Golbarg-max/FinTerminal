package handlers

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidSymbol(t *testing.T) {
	tests := []struct {
		symbol string
		want   bool
	}{
		{"AAPL", true},
		{"aapl", true},
		{" msft ", true},
		{"BRK.B", true},
		{"BF-B", true},
		{"SHOP.TO", true},
		{"A", true},
		{"ABCDEFGHIJ", true},
		{"", false},
		{"   ", false},
		{"ABCDEFGHIJK", false},
		{".AAPL", false},
		{"A B", false},
		{"AAPL&X=1", false},
		{"AAPL#", false},
		{"!!!", false},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			assert.Equal(t, tt.want, validSymbol(tt.symbol))
		})
	}
}

func TestHandleQuoteRejectsInvalidSymbol(t *testing.T) {
	client := &countingClient{}
	h := &Handlers{client: client}

	req := httptest.NewRequest(http.MethodGet, "/api/quote?symbol=AAPL%26apikey%3Devil", nil)
	rec := httptest.NewRecorder()
	h.HandleQuote(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid symbol")
	assert.Equal(t, 0, client.calls, "invalid symbol should not reach the API")
}

func TestAddHoldingRejectsInvalidSymbol(t *testing.T) {
	path := filepath.Join(t.TempDir(), "portfolio.json")
	h := &Handlers{client: &countingClient{}, portfolioPath: path}

	body := strings.NewReader(`{"symbol": "!!!", "quantity": 1, "cost_basis": 10}`)
	req := httptest.NewRequest(http.MethodPost, "/api/holdings", body)
	rec := httptest.NewRecorder()
	h.HandleHoldings(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid symbol")
	assert.NoFileExists(t, path)
}
