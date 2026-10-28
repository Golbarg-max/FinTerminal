package handlers

import (
	"FinTerminal/api"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleQuote(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		symbol     string
		clientErr  error
		wantStatus int
	}{
		{"success", http.MethodGet, "AAPL", nil, http.StatusOK},
		{"rate limited", http.MethodGet, "AAPL", fmt.Errorf("fetching: %w", api.ErrRateLimited), http.StatusTooManyRequests},
		{"symbol not found", http.MethodGet, "ZZZZ", fmt.Errorf("fetching: %w", api.ErrSymbolNotFound), http.StatusNotFound},
		{"wrong method", http.MethodPost, "AAPL", nil, http.StatusMethodNotAllowed},
		{"other error", http.MethodGet, "AAPL", fmt.Errorf("fetching: %w", errors.New("boom")), http.StatusBadGateway},
		{"missing symbol", http.MethodGet, "", nil, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &countingClient{err: tt.clientErr}
			h := &Handlers{client: client}

			req := httptest.NewRequest(tt.method, "/api/quote?symbol="+tt.symbol, nil)
			rec := httptest.NewRecorder()
			h.HandleQuote(rec, req)
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}

}
