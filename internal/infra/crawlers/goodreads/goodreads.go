package goodreads

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"

	"github.com/sofib/random-book-quote/internal/infra/event"
)

var urls = []string{
	"https://www.goodreads.com/work/quotes/3078186-the-hitchhiker-s-guide-to-the-galaxy",
	"https://www.goodreads.com/work/quotes/2842984-mostly-harmless",
	"https://www.goodreads.com/work/quotes/1783981-foundation",
	"https://www.goodreads.com/work/quotes/4110990-good-omens-the-nice-and-accurate-prophecies-of-agnes-nutter-witch",
}

type goodreadsFetcher struct {
	emitter         *event.Emitter
	quoteFoundEvent event.Event
}

func NewGoodreadsFetcher(
	emitter *event.Emitter,
	quoteFoundEvent event.Event,
) *goodreadsFetcher { //nolint:revive // legit use
	return &goodreadsFetcher{
		emitter:         emitter,
		quoteFoundEvent: quoteFoundEvent,
	}
}

func (f *goodreadsFetcher) Fetch(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, url := range urls {
		g.Go(func() error {
			return f.fetchURL(ctx, url)
		})
	}

	return g.Wait()
}

func (f *goodreadsFetcher) fetchURL(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return err
	}

	doc.Find(".quoteText").Each(func(_ int, s *goquery.Selection) {
		f.emitter.Emit(f.quoteFoundEvent, strings.TrimSpace(s.Text()))
	})
	return nil
}
