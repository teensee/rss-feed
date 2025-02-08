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

// NewTimeFormatter регистрирует 1 форматер
func NewTimeFormatter(format string) *TimeFormatter {
	return &TimeFormatter{format: format}
}

// NewTimeFormatters регистрирует все форматеры
func NewTimeFormatters() []rss.Processor {
	formatters := make([]rss.Processor, 0, len(availableFormats))
	for _, format := range availableFormats {
		formatters = append(formatters, NewTimeFormatter(format))
	}

	return formatters
}

// Name Возвращает текущий зарегистрированный формат
func (t *TimeFormatter) Name() string {
	return "time-formatter-" + t.format
}

// Process приводит дату к единому стилю
func (t *TimeFormatter) Process(items []*rss.Item) ([]*rss.Item, error) {
	correctLayout := ""

	for i, item := range items {
		if correctLayout == "" {
			// try check expected time format is the current time format
			if _, err := time.Parse(t.format, item.GetPubDate()); err == nil {
				continue
			}

			for _, layout := range availableFormats {
				_, err := time.Parse(layout, item.GetPubDate())
				if err == nil {
					correctLayout = layout
					break
				}
			}
		}

		parsedTime, err := time.Parse(correctLayout, item.GetPubDate())
		if err == nil {
			items[i] = rss.NewItem(
				item.GetTitle(),
				item.GetLink(),
				item.GetDescription(),
				parsedTime.Format(t.format),
				item.GetCreator(),
				item.GetCategories(),
			)
		} else {
			correctLayout = ""
		}
	}

	return items, nil
}
