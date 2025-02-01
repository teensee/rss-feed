package rss

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"rss-feed/pkg/http"
	"strings"
	"unicode/utf8"
)

var _ Processor = &HtmlSanitizer{}

type Processor interface {
	Name() string
	Process(items *[]http.Item) (*[]http.Item, error)
}

type HtmlSanitizer struct {
	policy *bluemonday.Policy
}

func NewHtmlSanitizer() *HtmlSanitizer {
	return &HtmlSanitizer{policy: bluemonday.StrictPolicy()}
}

func (h *HtmlSanitizer) Name() string {
	return "html-sanitizer"
}

func (h *HtmlSanitizer) Process(items *[]http.Item) (*[]http.Item, error) {
	rssItems := *items
	for i := range rssItems {
		rssItems[i].Title = h.policy.Sanitize(rssItems[i].Title)
		rssItems[i].Description = strings.TrimLeft(h.policy.Sanitize(rssItems[i].Description), "\n")
	}

	return &rssItems, nil
}

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

func (h *MaxLengthProcessor) Process(items *[]http.Item) (*[]http.Item, error) {
	rssItems := *items
	for i := range rssItems {
		rssItems[i].Description = h.Truncate(rssItems[i].Description)
	}

	return &rssItems, nil
}

func (h *MaxLengthProcessor) Truncate(str string) string {
	if h.maxLength > 0 && utf8.RuneCountInString(str) > h.maxLength {
		runes := []rune(str)

		return string(runes[:h.maxLength]) + h.postfix
	}

	return str
}

type SizeOfProcessor struct {
	maxSize int
}

func (s *SizeOfProcessor) Name() string {
	return fmt.Sprintf("size-of-%d", s.maxSize)
}

func (s *SizeOfProcessor) Process(items *[]http.Item) (*[]http.Item, error) {
	tmp := *items
	tmp = tmp[:s.maxSize]

	return &tmp, nil
}

func NewSizeOfProcessor(maxSize int) *SizeOfProcessor {
	return &SizeOfProcessor{maxSize: maxSize}
}
