package cache

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/patrickmn/go-cache"
	"rss-feed/internal/domain/logging"
	"time"
)

// compile time interface checks
var _ AppCache = &GoCache{}
var _ AppCache = &DummyCache{}

type CacheFn func() (interface{}, error)

type AppCache interface {
	Set(ctx context.Context, key Key, value interface{}, expiration time.Duration)
	Get(ctx context.Context, key Key) (interface{}, bool)
	DoGet(ctx context.Context, key Key, exp time.Duration, fn CacheFn) (interface{}, error)
}

type Key string

func NewPlainKey(key string) Key {
	return Key(key)
}

func NewMd5Key(key string) Key {
	// convert to hex
	return Key(fmt.Sprintf("%x", md5.Sum([]byte(key))))
}

func (k Key) String() string {
	return string(k)
}

type DummyCache struct{}

func NewDummyCache() AppCache {
	return &DummyCache{}
}

func (d *DummyCache) Set(_ context.Context, _ Key, _ interface{}, _ time.Duration) {
}

func (d *DummyCache) Get(_ context.Context, _ Key) (interface{}, bool) {
	return nil, false
}

func (d *DummyCache) DoGet(_ context.Context, _ Key, _ time.Duration, fn CacheFn) (interface{}, error) {
	return fn()
}

type GoCache struct {
	cache *cache.Cache
	l     logging.Logger
}

func NewGoCache(defaultExpiration, cleanupInterval time.Duration, l logging.Logger) AppCache {
	return &GoCache{cache: cache.New(defaultExpiration, cleanupInterval), l: l}
}

func (c *GoCache) Set(ctx context.Context, key Key, value interface{}, expiration time.Duration) {
	c.l.Debug(ctx, fmt.Sprintf("Execute cache set with key: %s", key))
	c.cache.Set(key.String(), value, expiration)
}

func (c *GoCache) Get(ctx context.Context, key Key) (interface{}, bool) {
	c.l.Debug(ctx, fmt.Sprintf("Execute cache get with key: %s", key))
	return c.cache.Get(key.String())
}

func (c *GoCache) DoGet(ctx context.Context, key Key, exp time.Duration, fn CacheFn) (interface{}, error) {
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
