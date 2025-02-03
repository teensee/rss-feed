package hasher

import "rss-feed/internal/domain/cache"

var _ cache.Hasher = &Plain{}

type Plain struct{}

func NewPlain() *Plain {
	return &Plain{}
}

func (p *Plain) Hash(key string) string {
	return key
}
