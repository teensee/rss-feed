package cache

import (
	"crypto/md5"
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

var _ AppCache = &GoCache{}

type AppCache interface {
	Set(key Key, value interface{}, expiration time.Duration)
	Get(key Key) (interface{}, bool)
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

type GoCache struct {
	cache *cache.Cache
}

func NewGoCache(cache *cache.Cache) AppCache {
	return &GoCache{cache: cache}
}

func (c *GoCache) Set(key Key, value interface{}, expiration time.Duration) {
	c.cache.Set(key.String(), value, expiration)
}

func (c *GoCache) Get(key Key) (interface{}, bool) {
	return c.cache.Get(key.String())
}
