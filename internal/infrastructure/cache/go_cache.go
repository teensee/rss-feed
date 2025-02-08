package cache

import (
	"context"
	"fmt"
	domainCache "rss-feed/internal/domain/cache"
	"rss-feed/internal/domain/logging"
	"time"

	"github.com/patrickmn/go-cache"
)

var _ domainCache.AppCache = &GoCache{}

type GoCache struct {
	cache *cache.Cache
	l     logging.Logger
}

func NewGoCache(defaultExpiration, cleanupInterval time.Duration, l logging.Logger) domainCache.AppCache {
	return &GoCache{cache: cache.New(defaultExpiration, cleanupInterval), l: l}
}

func (c *GoCache) Set(ctx context.Context, key domainCache.Key, value interface{}, expiration time.Duration) {
	c.l.Debug(ctx, fmt.Sprintf("Execute cache set with key: %s", key))
	c.cache.Set(key.String(), value, expiration)
}

func (c *GoCache) Get(ctx context.Context, key domainCache.Key) (interface{}, bool) {
	c.l.Debug(ctx, fmt.Sprintf("Execute cache get with key: %s", key))
	return c.cache.Get(key.String())
}

func (c *GoCache) DoGet(ctx context.Context, key domainCache.Key, exp time.Duration, fn domainCache.CacheFn) (interface{}, error) {
	if res, ok := c.Get(ctx, key); ok {
		return res, nil
	}

	res, err := fn()
	if err != nil {
		c.l.Error(ctx, fmt.Sprintf("Function return error: %s, for key: %s", err, key))
		return nil, err
	}

	c.Set(ctx, key, res, exp)

	return res, nil
}
