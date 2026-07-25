package api

import (
	"math/rand"
	"time"
)

type MockClient struct {
	lastUpdate time.Time
}

func NewMockClient() *MockClient {
	return &MockClient{
		lastUpdate: time.Now(),
	}
}

func (m *MockClient) GetQuote(symbol string) (*QuoteData, error) {
	mockData := map[string]*QuoteData{
		"AAPL": {
			Symbol:        "AAPL",
			Price:         189.95,
			Change:        2.50,
			ChangePercent: 1.33,
			Volume:        52345600,
			Timestamp:     time.Now(),
		},
		"MSFT": {
			Symbol:        "MSFT",
			Price:         334.20,
			Change:        3.50,
			ChangePercent: 1.06,
			Volume:        25000000,
			Timestamp:     time.Now(),
		},
		"GOOGL": {
			Symbol:        "GOOGL",
			Price:         2950.00,
			Change:        -15.00,
			ChangePercent: -0.50,
			Volume:        15000000,
			Timestamp:     time.Now(),
		},
		"TSLA": {
			Symbol:        "TSLA",
			Price:         265.30,
			Change:        5.00,
			ChangePercent: 1.92,
			Volume:        120000000,
			Timestamp:     time.Now(),
		},
		"NVDA": {
			Symbol:        "NVDA",
			Price:         875.00,
			Change:        12.50,
			ChangePercent: 1.45,
			Volume:        40000000,
			Timestamp:     time.Now(),
		},
		"AMZN": {
			Symbol:        "AMZN",
			Price:         180.50,
			Change:        -2.30,
			ChangePercent: -1.26,
			Volume:        55000000,
			Timestamp:     time.Now(),
		},
	}

	// If symbol exists in mock data, return it
	if data, exists := mockData[symbol]; exists {
		// Add slight random variation for realism
		return &QuoteData{
			Symbol:        data.Symbol,
			Price:         data.Price + (rand.Float64()-0.5)*2, // ±$1 variation
			Change:        data.Change + (rand.Float64()-0.5)*1,
			ChangePercent: data.ChangePercent + (rand.Float64()-0.5)*0.5,
			Volume:        data.Volume + rand.Int63n(5000000),
			Timestamp:     time.Now(),
		}, nil
	}

	// For unknown symbols, generate random data
	return &QuoteData{
		Symbol:        symbol,
		Price:         100.00 + rand.Float64()*200,
		Change:        (rand.Float64() - 0.5) * 10,
		ChangePercent: (rand.Float64() - 0.5) * 5,
		Volume:        rand.Int63n(100000000),
		Timestamp:     time.Now(),
	}, nil
}