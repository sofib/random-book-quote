package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/sofib/random-book-quote/internal/infra/event"
)

type Collector interface {
	Collect(quote any)
	Quotes() []any
}

type Randomizer interface {
	Random() any
	With(quotes []any) Randomizer
}

type App struct {
	settings   Settings
	collector  Collector
	randomizer Randomizer
}

func NewApp() *App {
	return &App{
		settings:   NewSettings(),
		collector:  NewQuoteCollector(),
		randomizer: NewQuoteRandomizer(),
	}
}

func (a *App) Run() {
	ctx := context.Background()
	emitter := event.NewEmitter()
	quoteCh := emitter.On(QuoteFoundEvent)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for quote := range quoteCh {
			a.collector.Collect(quote)
		}
	}()

	for _, crawler := range a.settings.EnabledCrawlers {
		if err := CrawlerFactory(emitter, crawler).Fetch(ctx); err != nil {
			logger.Error("Failed to fetch quotes for %s: %v", slog.String("crawler", crawler), slog.Any("error", err))
		}
	}

	emitter.Close()
	wg.Wait()
	randomQuote := a.randomizer.With(a.collector.Quotes()).Random()
	fmt.Println(randomQuote)
}
