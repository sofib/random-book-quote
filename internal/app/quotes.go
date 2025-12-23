package app

import (
	"math/rand" //nolint:depguard // keeping things simple
	"sync"

	"github.com/sofib/random-book-quote/internal/infra/event"
)

const (
	QuoteFoundEvent event.Event = "quote-found"
)

type QuoteCollector struct {
	mu     sync.Mutex
	quotes []any
}

func NewQuoteCollector() *QuoteCollector {
	return &QuoteCollector{
		quotes: make([]any, 0),
	}
}

func (c *QuoteCollector) Collect(quote any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.quotes = append(c.quotes, quote)
}

func (c *QuoteCollector) Quotes() []any {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.quotes
}

type QuoteRandomizer struct {
	quotes []any
}

func NewQuoteRandomizer() *QuoteRandomizer {
	return &QuoteRandomizer{}
}

func (r *QuoteRandomizer) With(quotes []any) Randomizer {
	r.quotes = quotes
	return r
}

func (r *QuoteRandomizer) Random() any {
	if len(r.quotes) == 0 {
		return nil
	}
	return r.quotes[rand.Intn(len(r.quotes))] //nolint:gosec // keeping things simple
}
