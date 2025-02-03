package processor

import (
	"rss-feed/internal/domain/rss"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var _ rss.Processor = &HtmlSanitizer{}

type HtmlSanitizer struct {
	policy *bluemonday.Policy
}

func NewHtmlSanitizer() *HtmlSanitizer {
	return &HtmlSanitizer{policy: bluemonday.StrictPolicy()}
}

func (h *HtmlSanitizer) Name() string {
	return "html-sanitizer"
}

func (h *HtmlSanitizer) Process(items []*rss.Item) ([]*rss.Item, error) {
	for i := range items {
		items[i].Title = h.policy.Sanitize(items[i].Title)
		items[i].Description = strings.TrimLeft(h.policy.Sanitize(items[i].Description), "\n")
	}

	return items, nil
}
