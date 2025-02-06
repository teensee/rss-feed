package processor

import (
	"fmt"
	"rss-feed/internal/domain/rss"
)

var _ rss.Processor = &SizeOfProcessor{}

const defaultMinimum = 8
const defaultMaximum = 64

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
	if maxSize <= 0 {
		maxSize = defaultMinimum
	}

	if maxSize > defaultMaximum {
		maxSize = defaultMaximum
	}

	return &SizeOfProcessor{maxSize: maxSize}
}
