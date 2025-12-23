package app

import (
	"context"

	"github.com/sofib/random-book-quote/internal/infra/crawlers/goodreads"
	"github.com/sofib/random-book-quote/internal/infra/crawlers/libquotes"
	"github.com/sofib/random-book-quote/internal/infra/event"
)

type Crawler interface {
	Fetch(ctx context.Context) error
}

func CrawlerFactory(emitter *event.Emitter, crawler string) Crawler {
	switch crawler {
	case "goodreads":
		return goodreads.NewGoodreadsFetcher(emitter, QuoteFoundEvent)
	case "libquotes":
		return libquotes.NewLibquotesFetcher(emitter, QuoteFoundEvent)
	default:
		return nil
	}
}
