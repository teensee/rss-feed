package composite

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"rss-feed/internal/domain/rss"
	"rss-feed/pkg/http"
	"sync"
)

type ProcessedFeed struct {
	Source string     `json:"source"`
	Feed   []FeedItem `json:"feed"`
}

type FeedItem struct {
	Link  string `json:"link,omitempty"`
	Title string `json:"title,omitempty"`
	Desc  string `json:"desc,omitempty"`
}

type CompositeRss struct {
	fetcher    rss.Fetcher
	processors []rss.Processor
	l          *slog.Logger
}

func NewCompositeRss(
	fetcher rss.Fetcher,
	processors []rss.Processor,
	l *slog.Logger,
) *CompositeRss {
	return &CompositeRss{
		fetcher:    fetcher,
		processors: processors,
		l:          l,
	}
}

func (r *CompositeRss) FeedAsync(ctx context.Context, urlList []*url.URL) ([]ProcessedFeed, error) {
	var (
		mu       sync.Mutex
		wg       sync.WaitGroup
		feedList []*http.RSS
		errCh    = make(chan error, len(urlList))
	)

	for _, feedUrl := range urlList {
		if urlList == nil {
			continue
		}

		wg.Add(1)

		go func(path url.URL) {
			defer wg.Done()

			feed, err := r.fetcher.Fetch(ctx, path)

			if err != nil {
				errCh <- fmt.Errorf("fetch rss failed: %w, path=%s", err, path.String())
				r.l.ErrorContext(ctx, "fetch rss failed", "path", path.String(), "err", err)
				return
			}

			mu.Lock()
			feedList = append(feedList, feed)
			mu.Unlock()
		}(*feedUrl)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		// todo добавить доп логику обработки ошибок
		r.l.WarnContext(ctx, "RSS fetch error", slog.Any("err", err))
	}

	processed := r.process(ctx, feedList)
	return processed, nil
}

func (r *CompositeRss) process(ctx context.Context, rssList []*http.RSS) []ProcessedFeed {
	var processedFeedList []ProcessedFeed

	// Стоит вынести вложенный цикл. Сначала обрабатываем, потом маппим
	for _, rssItem := range rssList {
		var items *[]http.Item

		var err error
		for _, processor := range r.processors {
			r.l.DebugContext(ctx, fmt.Sprintf("Run %s", processor.Name()))
			items, err = processor.Process(rssItem.Channel.Items)

			if err != nil {
				r.l.Error(fmt.Sprintf("processor: %s failed with error: %s", processor.Name(), err))
				continue
			}
		}

		if items == nil {
			// todo fix nil pointer
			continue
		}

		feed := make([]FeedItem, 0, len(*items))

		for _, i := range *items {
			feed = append(feed, FeedItem{
				Title: i.Title,
				Link:  i.Link,
				Desc:  i.Description,
			})
		}

		processedFeedList = append(processedFeedList, ProcessedFeed{
			Source: rssItem.Channel.Title,
			Feed:   feed,
		})
	}

	return processedFeedList
}
