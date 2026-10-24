package api

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeClient struct {
	calls int
	err   error
	mu    sync.Mutex
}

func (f *fakeClient) GetQuote(symbol string) (*QuoteData, error) {
	f.mu.Lock()
	f.calls++
	f.mu.Unlock()
	if f.err != nil {
		return nil, f.err
	}
	return &QuoteData{Symbol: symbol, Price: 100}, nil
}

func TestCachedClientHit(t *testing.T) {
	fake := &fakeClient{}
	cached := NewCachedClient(fake, time.Minute)

	cached.GetQuote("AAPL")
	cached.GetQuote("AAPL")

	assert.Equal(t, 1, fake.calls)
}

func TestNormalization(t *testing.T) {
	normalized := NormalizeSymbol(" aapl ")
	assert.Equal(t, "AAPL", normalized)
}

func TestExpiry(t *testing.T) {
	fake := &fakeClient{}
	cached := NewCachedClient(fake, time.Millisecond)

	cached.GetQuote("AAPL")
	assert.Equal(t, 1, fake.calls)

	time.Sleep(5 * time.Millisecond)

	cached.GetQuote("AAPL")

	assert.Equal(t, 2, fake.calls)
}

func TestErrorsNotCached(t *testing.T) {
	fake := &fakeClient{err: errors.New("BOOM")}
	cached := NewCachedClient(fake, time.Minute)

	cached.GetQuote("AAPL")
	_, err := cached.GetQuote("AAPL")

	assert.Equal(t, 2, fake.calls)
	assert.Error(t, err)

}

func TestConcurrency(t *testing.T) {
	fake := &fakeClient{}
	cached := NewCachedClient(fake, time.Minute)
	symbols := []string{"AAPL", "MSFT", "TSLA"}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			cached.GetQuote(symbols[n%len(symbols)])
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	assert.GreaterOrEqual(t, fake.calls, 3)
	assert.LessOrEqual(t, fake.calls, 50)
}
