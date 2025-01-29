package habr

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"net/url"
	"rss-feed/pkg/cache"
	"rss-feed/pkg/http"
	"time"
)

type Habr struct {
	client *http.Client
	l      *slog.Logger
	cache  cache.AppCache
}

func NewHabr(l *slog.Logger, cache cache.AppCache) *Habr {
	return &Habr{
		l: l,
		client: http.NewClient(
			url.URL{
				Scheme: "https",
				Host:   "habr.com",
				Path:   "/ru/",
			},
			l,
		),
		cache: cache,
	}
}

func (h *Habr) Feed(ctx context.Context, filter RssFilter) (*http.RSS, error) {
	endpoint := filter.BuildRssUrl()
	h.l.Debug(fmt.Sprintf("execute request to %s", endpoint))
	resp, err := h.client.GET(ctx, endpoint, nil)

	var rss http.RSS
	if err != nil {
		h.l.ErrorContext(ctx, fmt.Sprintf("execute request to %s failed", endpoint))
		return nil, err
	}

	err = xml.Unmarshal(resp, &rss)
	if err != nil {
		return nil, err
	}

	rss.XMLNSDC = "http://purl.org/dc/elements/1.1/"

	return &rss, nil
}

func (h *Habr) FeedCached(ctx context.Context, filter RssFilter) (*http.RSS, error) {
	if res, ok := h.cache.Get(cache.NewMd5Key(filter.BuildRssUrl())); ok {
		h.l.DebugContext(ctx, fmt.Sprintf("cache found for current query %s", filter.BuildRssUrl()))

		return res.(*http.RSS), nil
	}

	res, err := h.Feed(ctx, filter)

	if err != nil {
		h.l.ErrorContext(ctx, fmt.Sprintf("fetch rss failed: %s", err))
		return nil, err
	}

	h.cache.Set(cache.NewMd5Key(filter.BuildRssUrl()), res, 5*time.Minute)

	return res, nil
}
