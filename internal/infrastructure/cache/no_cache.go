package cache

import (
	"context"
	domainCache "rss-feed/internal/domain/cache"
	"time"
)

var _ domainCache.AppCache = &NoCache{}

type NoCache struct{}

func NewDummyCache() domainCache.AppCache {
	return &NoCache{}
}

func (d *NoCache) Set(_ context.Context, _ domainCache.Key, _ interface{}, _ time.Duration) {
}

func (d *NoCache) Get(_ context.Context, _ domainCache.Key) (interface{}, bool) {
	return nil, false
}

func (d *NoCache) DoGet(_ context.Context, _ domainCache.Key, _ time.Duration, fn domainCache.CacheFn) (interface{}, error) {
	return fn()
}
