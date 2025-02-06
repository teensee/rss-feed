package processor

import (
	"rss-feed/internal/domain/rss"
	"unicode/utf8"
)

var _ rss.Processor = &MaxLengthProcessor{}

type MaxLengthProcessor struct {
	maxLength int
	postfix   string
}

type MaxLengthOption struct {
	MaxLength int
	Postfix   string
}

func NewMaxLengthProcessor(opts ...MaxLengthOption) *MaxLengthProcessor {
	var opt MaxLengthOption

	switch len(opts) {
	case 0:
		opt = MaxLengthOption{
			MaxLength: 128,
			Postfix:   "...",
		}
	default:
		opt = opts[0]
	}

	return &MaxLengthProcessor{
		maxLength: opt.MaxLength,
		postfix:   opt.Postfix,
	}
}

func (h *MaxLengthProcessor) Name() string {
	return "max-length"
}

func (h *MaxLengthProcessor) Process(items []*rss.Item) ([]*rss.Item, error) {
	for i, item := range items {
		truncatedDesc := h.Truncate(item.GetDescription())
		if truncatedDesc != item.GetDescription() {
			items[i] = rss.NewItem(
				item.GetTitle(),
				item.GetLink(),
				truncatedDesc,
				item.GetPubDate(),
				item.GetCreator(),
				item.GetCategories(),
			)
		}
	}

	return items, nil
}

func (h *MaxLengthProcessor) Truncate(str string) string {
	if h.maxLength > 0 && utf8.RuneCountInString(str) > h.maxLength {
		runes := []rune(str)

		return string(runes[:h.maxLength]) + h.postfix
	}

	return str
}
