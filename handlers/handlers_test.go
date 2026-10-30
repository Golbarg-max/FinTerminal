package handlers

import (
	"FinTerminal/api"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

// countingClient records how many times GetQuote was called.
type countingClient struct {
	calls int
	err   error
}

func (c *countingClient) GetQuote(symbol string) (*api.QuoteData, error) {
	c.calls++
	if c.err != nil {
		return nil, c.err
	}
	return &api.QuoteData{Symbol: symbol, Price: 1}, nil
}

func do(t *testing.T, handler http.HandlerFunc, method string, target string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, body)
	rec := httptest.NewRecorder()
	handler(rec, req)
	return rec
}


func newTestHandlers(t *testing.T) *Handlers {
	client := &countingClient{}
	h := &Handlers{
		client: client,
		portfolioPath: filepath.Join(t.TempDir(), "portfolio.json"),
	}
	return h
}