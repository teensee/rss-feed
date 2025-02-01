package composite

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/url"
	"rss-feed/internal/domain/rss"
	"rss-feed/pkg/cache"
	"rss-feed/pkg/http"
	"time"
)

var _ rss.Fetcher = &Fetcher{}

type Fetcher struct {
	client http.HttpClient
	cache  cache.AppCache
	l      *slog.Logger
}

func NewCompositeFetcher(client http.HttpClient, cache cache.AppCache, l *slog.Logger) rss.Fetcher {
	return &Fetcher{client: client, cache: cache, l: l}
}

func (f *Fetcher) Fetch(ctx context.Context, url url.URL) (*http.RSS, error) {
	f.l.InfoContext(ctx, fmt.Sprintf("fetching feed %s", url.String()))

	rssFeed, err := f.cache.DoGet(ctx, cache.NewMd5Key(url.String()), 1*time.Minute, func() (interface{}, error) {
		resp, err := f.client.GET(ctx, url.String(), nil)
		if err != nil {
			f.l.ErrorContext(ctx, fmt.Sprintf("fetching feed %s", url.String()), slog.Any("err", err))
			return nil, err
		}

		var rssFeed http.RSS
		err = xml.Unmarshal(resp, &rssFeed)
		if err != nil {
			return nil, err
		}

		rssFeed.XMLNSDC = "http://purl.org/dc/elements/1.1/"
		return &rssFeed, nil
	})

	if err != nil {
		return nil, err
	}

	return rssFeed.(*http.RSS), nil
}
