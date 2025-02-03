package hasher

import (
	"crypto/sha256"
	"fmt"
	"rss-feed/internal/domain/cache"
)

var _ cache.Hasher = &Sha256Hasher{}

type Sha256Hasher struct{}

func NewSha256Hasher() *Sha256Hasher {
	return &Sha256Hasher{}
}

func (p *Sha256Hasher) Hash(key string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(key)))
}
