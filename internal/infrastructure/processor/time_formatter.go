package processor

import (
	"rss-feed/internal/domain/rss"
	"time"
)

var _ rss.Processor = &TimeFormatter{}

var availableFormats = []string{
	time.Layout,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,
}

type TimeFormatter struct {
	format string
}

func NewTimeFormatter(format string) *TimeFormatter {
	return &TimeFormatter{format: format}
}

func NewTimeFormatters() []rss.Processor {
	formatters := make([]rss.Processor, 0, len(availableFormats))
	for _, format := range availableFormats {
		formatters = append(formatters, NewTimeFormatter(format))
	}
	return formatters
}

func (t *TimeFormatter) Name() string {
	return "time-formatter-" + t.format
}

func (t *TimeFormatter) Process(items []*rss.Item) ([]*rss.Item, error) {
	for _, item := range items {
		for _, layout := range availableFormats {
			parsedTime, err := time.Parse(layout, item.PubDate)
			if err == nil {
				item.PubDate = parsedTime.Format(t.format)
				break
			}
		}
	}

	return items, nil
}
