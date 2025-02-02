package processor

import (
	"fmt"
	"rss-feed/internal/domain/rss"
)

var _ rss.Processor = &SizeOfProcessor{}

type SizeOfProcessor struct {
	maxSize int
}

func (s *SizeOfProcessor) Name() string {
	return fmt.Sprintf("size-of-%d", s.maxSize)
}

func (s *SizeOfProcessor) Process(items []*rss.Item) ([]*rss.Item, error) {
	items = items[:s.maxSize]

	return items, nil
}

func NewSizeOfProcessor(maxSize int) *SizeOfProcessor {
	return &SizeOfProcessor{maxSize: maxSize}
}
