package rss

import (
	"context"
)

type Fetcher interface {
	Fetch(ctx context.Context, url string) (*Feed, error)
}
