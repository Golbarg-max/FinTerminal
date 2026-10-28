package handlers

import "FinTerminal/api"

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
