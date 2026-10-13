package api

import (
	"strings"
)

type Client interface {
	GetQuote(symbol string)(*QuoteData, error)
}

func NewClient(apiKey string)(Client, bool){
	// Check if API key is provided and not "demo"
	if apiKey == "" || strings.ToLower(apiKey) == "demo" {
		// Use mock client
		return NewMockClient(), true
	} 
	// Use real Alpha Vantage client
	return NewAlphaVantageClient(apiKey), false
}