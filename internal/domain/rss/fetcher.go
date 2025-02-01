package rss

import (
	"context"
	"net/url"
	"rss-feed/pkg/http"
)

type Fetcher interface {
	Fetch(ctx context.Context, url url.URL) (*http.RSS, error)
}
