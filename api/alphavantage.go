package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrEmptyAPIKey    = fmt.Errorf("API key cannot be empty")
	ErrSymbolNotFound = fmt.Errorf("symbol not found")
	ErrInvalidPrice   = fmt.Errorf("invalid price data")
	ErrRateLimited    = fmt.Errorf("API limit reached")
	BaseURL           = "https://www.alphavantage.co/query"
)

type AlphaVantageClient struct {
	apiKey string
	client *http.Client
}

type QuoteData struct {
	Symbol        string
	Price         float64
	Change        float64
	ChangePercent float64
	Volume        int64
	Timestamp     time.Time
}

type globalQuoteResponse struct {
	GlobalQuote struct {
		Symbol           string `json:"01. symbol"`
		Price            string `json:"05. price"`
		Change           string `json:"09. change"`
		ChangePercent    string `json:"10. change percent"`
		Volume           string `json:"06. volume"`
		LatestTradingDay string `json:"07. latest trading day"`
	} `json:"Global Quote"`
	Information string `json:"Information"`
	Note        string `json:"note"`
	Error       string `json:"Error Message"`
}

func NewAlphaVantageClient(apiKey string) *AlphaVantageClient {
	if apiKey == "" {
		return nil
	}

	return &AlphaVantageClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *AlphaVantageClient) GetQuote(symbol string) (*QuoteData, error) {
	if c == nil {
		return nil, fmt.Errorf("client is nil")
	}

	// Normalize before validating so blank input like "   " is rejected
	symbol = NormalizeSymbol(symbol)
	if symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("apikey", c.apiKey)
	params.Add("function", "GLOBAL_QUOTE")
	paramsEncoded := params.Encode()

	resp, err := c.client.Get(BaseURL + "?" + paramsEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quote for %s: %w", symbol, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d for symbol %s", resp.StatusCode, symbol)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result globalQuoteResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Check for API errors
	if result.Note != "" {
		return nil, fmt.Errorf("API rate limit or error: %s", result.Note)
	}
	if result.Information != "" {
		return nil, fmt.Errorf("%w: %s", ErrRateLimited, result.Information)
	}
	if result.Error != "" {
		return nil, fmt.Errorf("API error: %s", result.Error)
	}

	if result.GlobalQuote.Symbol == "" {
		return nil, fmt.Errorf("%w: %s", ErrSymbolNotFound, symbol)
	}

	// Parse numeric fields with error handling
	price, err := strconv.ParseFloat(result.GlobalQuote.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse price '%s': %w", result.GlobalQuote.Price, err)
	}

	change, err := strconv.ParseFloat(result.GlobalQuote.Change, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse change '%s': %w", result.GlobalQuote.Change, err)
	}

	// Remove '%' from change percent before parsing
	changePercentStr := strings.TrimSuffix(result.GlobalQuote.ChangePercent, "%")
	changePercent, err := strconv.ParseFloat(changePercentStr, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse change percent '%s': %w", result.GlobalQuote.ChangePercent, err)
	}

	volume, err := strconv.ParseInt(result.GlobalQuote.Volume, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse volume '%s': %w", result.GlobalQuote.Volume, err)
	}

	timestamp, err := time.Parse("2006-01-02", result.GlobalQuote.LatestTradingDay)
	if err != nil {
		// Use current time if parsing fails
		timestamp = time.Now()
	}

	return &QuoteData{
		Symbol:        result.GlobalQuote.Symbol,
		Price:         price,
		Change:        change,
		ChangePercent: changePercent,
		Volume:        volume,
		Timestamp:     timestamp,
	}, nil
}

// Validate checks if quote data is valid
func (q *QuoteData) Validate() error {
	if q.Symbol == "" {
		return fmt.Errorf("symbol cannot be empty")
	}
	if q.Price < 0 {
		return fmt.Errorf("price cannot be negative: %f", q.Price)
	}
	return nil
}
