package http

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"rss-feed/internal/domain/logging"
	"rss-feed/internal/domain/rss"
	"rss-feed/internal/infrastructure/cache"
	"rss-feed/internal/infrastructure/http/dto"
	"time"
)

var _ rss.Fetcher = &Fetcher{}

type Fetcher struct {
	client HttpClient
	cache  cache.AppCache
	l      logging.Logger
}

func NewFeedFetcher(client HttpClient, cache cache.AppCache, l logging.Logger) rss.Fetcher {
	return &Fetcher{client: client, cache: cache, l: l}
}

func (f *Fetcher) Fetch(ctx context.Context, url string) (*rss.Feed, error) {
	f.l.Info(ctx, fmt.Sprintf("fetching feed %s", url))

	rssFeed, err := f.cache.DoGet(
		ctx,
		cache.NewMd5Key(url),
		1*time.Minute,
		func() (interface{}, error) {
			return f.doFetch(ctx, url)
		},
	)

	if err != nil {
		return nil, err
	}

	return toDomainModel(rssFeed.(*dto.RSS)), nil
}

func (f *Fetcher) doFetch(ctx context.Context, url string) (interface{}, error) {
	resp, err := f.client.GET(ctx, url, nil)
	if err != nil {
		f.l.Error(ctx, fmt.Sprintf("fetching feed %s", url), slog.Any("err", err))
		return nil, err
	}

	var rssFeed dto.RSS
	err = xml.Unmarshal(resp, &rssFeed)
	if err != nil {
		return nil, err
	}

	rssFeed.XMLNSDC = "http://purl.org/dc/elements/1.1/"
	return &rssFeed, nil
}
