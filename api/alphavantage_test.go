package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAlphaVantageClient(t *testing.T) {
	client := NewAlphaVantageClient("test_key")
	assert.NotNil(t, client)
	assert.Equal(t, "test_key", client.apiKey)
	assert.NotNil(t, client.client)
}

func TestNewAlphaVantageClientEmptyKey(t *testing.T) {
	client := NewAlphaVantageClient("")
	assert.Nil(t, client)
}

func TestGetQuoteSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GLOBAL_QUOTE", r.URL.Query().Get("function"))
		assert.Equal(t, "AAPL", r.URL.Query().Get("symbol"))

		response := `{
			"Global Quote": {
				"01. symbol": "AAPL",
				"05. price": "189.95",
				"09. change": "2.50",
				"10. change percent": "1.33%",
				"06. volume": "52345600",
				"07. latest trading day": "2024-01-15"
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, response)
	}))
	defer server.Close()

	client := &AlphaVantageClient{
		apiKey: "test_key",
		client: &http.Client{},
	}

	oldBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = oldBaseURL }()

	quote, err := client.GetQuote("AAPL")
	assert.NoError(t, err)
	assert.NotNil(t, quote)
	assert.Equal(t, "AAPL", quote.Symbol)
	assert.Equal(t, 189.95, quote.Price)
	assert.Equal(t, 2.50, quote.Change)
	assert.Equal(t, 1.33, quote.ChangePercent)
	assert.Equal(t, int64(52345600), quote.Volume)
}

func TestGetQuoteEmptySymbol(t *testing.T) {
	client := NewAlphaVantageClient("test_key")
	quote, err := client.GetQuote("")
	assert.Error(t, err)
	assert.Nil(t, quote)
	assert.Contains(t, err.Error(), "symbol cannot be empty")
}

func TestGetQuoteNilClient(t *testing.T) {
	var client *AlphaVantageClient
	quote, err := client.GetQuote("AAPL")
	assert.Error(t, err)
	assert.Nil(t, quote)
}

func TestGetQuoteEmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{"Global Quote": {}}`
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, response)
	}))
	defer server.Close()

	client := &AlphaVantageClient{
		apiKey: "test_key",
		client: &http.Client{},
	}

	oldBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = oldBaseURL }()

	quote, err := client.GetQuote("NONEXISTENT")
	assert.Error(t, err)
	assert.Nil(t, quote)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetQuoteInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "invalid json")
	}))
	defer server.Close()

	client := &AlphaVantageClient{
		apiKey: "test_key",
		client: &http.Client{},
	}

	oldBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = oldBaseURL }()

	quote, err := client.GetQuote("AAPL")
	assert.Error(t, err)
	assert.Nil(t, quote)
	assert.Contains(t, err.Error(), "parse")
}

func TestGetQuoteServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Server Error")
	}))
	defer server.Close()

	client := &AlphaVantageClient{
		apiKey: "test_key",
		client: &http.Client{},
	}

	oldBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = oldBaseURL }()

	quote, err := client.GetQuote("AAPL")
	assert.Error(t, err)
	assert.Nil(t, quote)
	assert.Contains(t, err.Error(), "status")
}

func TestGetQuoteAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"Error Message": "Invalid API call"
		}`
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, response)
	}))
	defer server.Close()

	client := &AlphaVantageClient{
		apiKey: "test_key",
		client: &http.Client{},
	}

	oldBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = oldBaseURL }()

	quote, err := client.GetQuote("AAPL")
	assert.Error(t, err)
	assert.Nil(t, quote)
}

func TestGetQuoteRateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"note": "Thank you for using Alpha Vantage! Our standard API call frequency is 5 calls per minute."
		}`
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, response)
	}))
	defer server.Close()

	client := &AlphaVantageClient{
		apiKey: "test_key",
		client: &http.Client{},
	}

	oldBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = oldBaseURL }()

	quote, err := client.GetQuote("AAPL")
	assert.Error(t, err)
	assert.Nil(t, quote)
	assert.Contains(t, err.Error(), "rate limit")
}

// Alpha Vantage's current daily-limit response: HTTP 200 with an "Information" body.
func TestGetQuoteInformationRateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"Information": "Thank you for using Alpha Vantage! Our standard API rate limit is 25 requests per day."
		}`
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, response)
	}))
	defer server.Close()

	client := &AlphaVantageClient{
		apiKey: "test_key",
		client: &http.Client{},
	}

	oldBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = oldBaseURL }()

	quote, err := client.GetQuote("AAPL")
	assert.Nil(t, quote)
	assert.ErrorIs(t, err, ErrRateLimited)
	assert.NotErrorIs(t, err, ErrSymbolNotFound)
	assert.Contains(t, err.Error(), "25 requests per day")
}

func TestGetQuoteWithChangePercent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"Global Quote": {
				"01. symbol": "MSFT",
				"05. price": "334.20",
				"09. change": "3.50",
				"10. change percent": "1.06%",
				"06. volume": "25000000",
				"07. latest trading day": "2024-01-16"
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, response)
	}))
	defer server.Close()

	client := &AlphaVantageClient{
		apiKey: "test_key",
		client: &http.Client{},
	}

	oldBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = oldBaseURL }()

	quote, err := client.GetQuote("MSFT")
	require.NoError(t, err)
	assert.Equal(t, 1.06, quote.ChangePercent)
}

func TestGetQuoteSymbolNormalization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify symbol is uppercase
		symbol := r.URL.Query().Get("symbol")
		assert.Equal(t, "AAPL", symbol)

		response := `{
			"Global Quote": {
				"01. symbol": "AAPL",
				"05. price": "189.95",
				"09. change": "0.00",
				"10. change percent": "0.00%",
				"06. volume": "0",
				"07. latest trading day": "2024-01-15"
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, response)
	}))
	defer server.Close()

	client := &AlphaVantageClient{
		apiKey: "test_key",
		client: &http.Client{},
	}

	oldBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = oldBaseURL }()

	// Pass lowercase symbol
	quote, err := client.GetQuote("aapl")
	assert.NoError(t, err)
	assert.Equal(t, "AAPL", quote.Symbol)
}

func TestGetQuoteNegativeValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"Global Quote": {
				"01. symbol": "AAPL",
				"05. price": "150.00",
				"09. change": "-10.00",
				"10. change percent": "-6.25%",
				"06. volume": "30000000",
				"07. latest trading day": "2024-01-15"
			}
		}`
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, response)
	}))
	defer server.Close()

	client := &AlphaVantageClient{
		apiKey: "test_key",
		client: &http.Client{},
	}

	oldBaseURL := BaseURL
	BaseURL = server.URL
	defer func() { BaseURL = oldBaseURL }()

	quote, err := client.GetQuote("AAPL")
	assert.NoError(t, err)
	assert.Equal(t, -10.0, quote.Change)
	assert.Equal(t, -6.25, quote.ChangePercent)
}

// Symbols with URL-special characters must arrive as one escaped value,
// not leak into extra query parameters.
func TestGetQuoteEscapesSymbol(t *testing.T) {
	tests := []struct {
		name   string
		symbol string
		want   string
	}{
		{"ampersand injection", "aapl&function=TIME_SERIES_DAILY", "AAPL&FUNCTION=TIME_SERIES_DAILY"},
		{"apikey injection", "AAPL&apikey=evil", "AAPL&APIKEY=EVIL"},
		{"space", "BRK B", "BRK B"},
		{"hash", "AAPL#x", "AAPL#X"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query()
				assert.Equal(t, []string{"GLOBAL_QUOTE"}, q["function"])
				assert.Equal(t, []string{"test_key"}, q["apikey"])
				assert.Equal(t, []string{tt.want}, q["symbol"])
				assert.Len(t, q, 3, "unexpected extra query parameters: %v", q)

				fmt.Fprint(w, `{"Global Quote": {"01. symbol": "X", "05. price": "1.00",
					"09. change": "0.00", "10. change percent": "0.00%", "06. volume": "0"}}`)
			}))
			defer server.Close()

			client := &AlphaVantageClient{
				apiKey: "test_key",
				client: &http.Client{},
			}

			oldBaseURL := BaseURL
			BaseURL = server.URL
			defer func() { BaseURL = oldBaseURL }()

			_, err := client.GetQuote(tt.symbol)
			require.NoError(t, err)
		})
	}
}
