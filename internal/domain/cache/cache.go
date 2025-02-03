package cache

import (
	"context"
	"time"
)

type CacheFn func() (interface{}, error)

type AppCache interface {
	Set(ctx context.Context, key Key, value interface{}, expiration time.Duration)
	Get(ctx context.Context, key Key) (interface{}, bool)
	DoGet(ctx context.Context, key Key, exp time.Duration, fn CacheFn) (interface{}, error)
}

type Key string

func (k Key) String() string {
	return string(k)
}

type Hasher interface {
	Hash(key string) string
}

type KeyGenerator interface {
	FromString(plain string) Key
}

type HashGenerator struct {
	hasher Hasher
}

func NewHashGenerator(hasher Hasher) *HashGenerator {
	return &HashGenerator{hasher: hasher}
}

func (g *HashGenerator) FromString(plain string) Key {
	return Key(g.hasher.Hash(plain))
}
