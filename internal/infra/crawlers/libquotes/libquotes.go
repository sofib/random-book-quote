package libquotes

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"

	"github.com/sofib/random-book-quote/internal/infra/event"
)

var urls = []string{
	"https://libquotes.com/charles-dickens",
	"https://libquotes.com/fyodor-dostoyevsky",
	"https://libquotes.com/terry-pratchett",
	"https://libquotes.com/j-r-r-tolkien",
}

type libquotesFetcher struct {
	emitter         *event.Emitter
	quoteFoundEvent event.Event
}

func NewLibquotesFetcher(
	emitter *event.Emitter,
	quoteFoundEvent event.Event,
) *libquotesFetcher { //nolint:revive // legit use
	return &libquotesFetcher{
		emitter:         emitter,
		quoteFoundEvent: quoteFoundEvent,
	}
}

func (f *libquotesFetcher) Fetch(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, url := range urls {
		g.Go(func() error {
			return f.fetchURL(ctx, url)
		})
	}

	return g.Wait()
}

func (f *libquotesFetcher) fetchURL(ctx context.Context, url string) error {
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

	doc.Find(".quote_span").Each(func(_ int, s *goquery.Selection) {
		quote := strings.TrimSpace(s.Text())

		quoteLink := s.Parent()
		author := strings.TrimSpace(quoteLink.NextAll().Find("div a").First().Text())

		f.emitter.Emit(f.quoteFoundEvent, fmt.Sprintf("%s \n - \n %s", quote, author))
	})
	return nil
}
