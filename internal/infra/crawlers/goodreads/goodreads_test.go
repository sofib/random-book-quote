package goodreads

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sofib/random-book-quote/internal/infra/event"
)

func TestFetchGoodreadsQuotes(t *testing.T) {
	ctx := context.Background()
	emitter := event.NewEmitter()
	quoteCh := emitter.On("quote-found")
	quotes := make([]any, 0)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range quoteCh {
			quotes = append(quotes, data)
		}
	}()

	fetcher := NewGoodreadsFetcher(emitter, "quote-found")
	err := fetcher.Fetch(ctx)
	if err != nil {
		t.Fatalf("Failed to fetch quotes: %v", err)
	}

	emitter.Close()
	wg.Wait()
	assert.NotEmpty(t, quotes)
	for _, quote := range quotes {
		assert.NotEmpty(t, quote)
	}
}
