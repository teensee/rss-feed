package application

import (
	"context"
	"fmt"
	"log/slog"
	appDto "rss-feed/internal/application/dto"
	"rss-feed/internal/domain/logging"
	"rss-feed/internal/domain/rss"
	"sync"
)

type FeedAggregator struct {
	fetcher  rss.Fetcher
	registry rss.ProcessorRegistry
	l        logging.Logger
}

func NewFeedService(
	fetcher rss.Fetcher,
	registry rss.ProcessorRegistry,
	l logging.Logger,
) *FeedAggregator {
	return &FeedAggregator{
		fetcher:  fetcher,
		registry: registry,
		l:        l,
	}
}

func (s *FeedAggregator) AggregateFeedAsync(
	ctx context.Context,
	req appDto.AppRssFeedRequest,
) ([]*rss.Feed, error) {
	var (
		mu       sync.Mutex
		wg       sync.WaitGroup
		feedList []*rss.Feed
		errCh    = make(chan error, len(req.Items))
	)

	for _, feedItem := range req.Items {
		wg.Add(1)

		go func(feedItem *appDto.RssFeedItemProcess) {
			defer wg.Done()

			path := feedItem.Rss
			feed, err := s.doAggregate(ctx, path, feedItem.Filters)

			if err != nil {
				errCh <- fmt.Errorf("fetch rss failed: %w, path=%s", err, path)
				s.l.Error(ctx, "fetch rss failed", "path", path, "err", err)
				return
			}

			mu.Lock()
			feedList = append(feedList, feed)
			mu.Unlock()
		}(feedItem)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		// todo добавить доп логику обработки ошибок
		s.l.Warn(ctx, "RSS fetch error", slog.Any("err", err))
	}

	return feedList, nil
}

func (s *FeedAggregator) doAggregate(ctx context.Context, path string, procList []string) (*rss.Feed, error) {
	feed, err := s.fetcher.Fetch(ctx, path)

	if err != nil {
		return nil, fmt.Errorf("fetch rss failed: %w", err)
	}

	if feed == nil {
		return nil, nil
	}

	processedItems, err := s.processItems(ctx, feed.Channel.Items, procList)
	if err != nil {
		s.l.Warn(ctx, "error processing items", slog.Any("err", err))
	}

	return &rss.Feed{
		Channel: &rss.Channel{
			Title:       feed.Channel.Title,
			Link:        feed.Channel.Link,
			Description: feed.Channel.Description,
			Items:       processedItems,
		},
	}, nil
}

func (s *FeedAggregator) processItems(ctx context.Context, items []*rss.Item, procList []string) ([]*rss.Item, error) {
	var errList []error

	for _, slug := range procList {
		proc, ok := s.registry.Resolve(slug)
		if !ok {
			s.l.Debug(ctx, "processor not registered", "slug", slug)
			continue
		}

		newItems, err := proc.Process(items)
		if err != nil {
			errList = append(errList, fmt.Errorf("processor %s failed: %w", slug, err))
			continue
		}
		items = newItems
	}

	if len(errList) > 0 {
		return items, fmt.Errorf("processing errors: %v", errList)
	}

	return items, nil
}
