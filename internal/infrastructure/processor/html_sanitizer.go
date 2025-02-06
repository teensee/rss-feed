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
		var newTitle, newDescription string

		currItem := items[i]
		newTitle = h.policy.Sanitize(currItem.GetTitle())
		newDescription = strings.TrimLeft(h.policy.Sanitize(currItem.GetDescription()), "\n")

		items[i] = rss.NewItem(
			newTitle,
			currItem.GetLink(),
			newDescription,
			currItem.GetPubDate(),
			currItem.GetCreator(),
			currItem.GetCategories(),
		)
	}

	return items, nil
}
