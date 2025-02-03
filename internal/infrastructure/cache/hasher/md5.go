package hasher

import (
	"crypto/md5" // nolint:gosec // используем осознанно
	"fmt"
	"rss-feed/internal/domain/cache"
)

var _ cache.Hasher = &Md5Hasher{}

type Md5Hasher struct{}

func NewMd5Hasher() *Md5Hasher {
	return &Md5Hasher{}
}

func (p *Md5Hasher) Hash(key string) string {
	// nolint:gosec // используем осознанно
	return fmt.Sprintf("%x", md5.Sum([]byte(key)))
}
