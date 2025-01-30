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

func (h *Habr) feed(ctx context.Context, url string) (*http.RSS, error) {
	h.l.Debug(fmt.Sprintf("execute request to %s", url))
	resp, err := h.client.GET(ctx, url, nil)

	if err != nil {
		h.l.ErrorContext(ctx, fmt.Sprintf("execute request to %s failed", url))
		return nil, err
	}

	var rss http.RSS

	err = xml.Unmarshal(resp, &rss)
	if err != nil {
		return nil, err
	}

	rss.XMLNSDC = "http://purl.org/dc/elements/1.1/"

	return &rss, nil
}

func (h *Habr) Feed(ctx context.Context, filter RssFilter) (*http.RSS, error) {
	return h.feed(ctx, filter.BuildRssUrl())
}

func (h *Habr) FeedCached(ctx context.Context, filter RssFilter) (*http.RSS, error) {
	key := cache.NewMd5Key(filter.BuildRssUrl())
	res, err := h.cache.DoGet(ctx, key, 5*time.Minute, func() (interface{}, error) {
		return h.Feed(ctx, filter)
	})

	if err != nil {
		h.l.ErrorContext(ctx, fmt.Sprintf("cache get failed: %s", err))
		return nil, err
	}

	return res.(*http.RSS), nil
}
