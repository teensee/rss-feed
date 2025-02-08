package processor

import (
	"rss-feed/internal/domain/rss"
	"testing"
	"time"
)

func TestTimeFormatter_Process_AllFormats(t *testing.T) {
	now := time.Now()

	for _, format := range availableFormats {
		t.Run("Testing format: "+format, func(t *testing.T) {
			formattedDate := now.Format(format)

			formatter := NewTimeFormatter(format)

			items := []*rss.Item{
				rss.NewItem("Test Title", "http://example.com", "Test Description", formattedDate, "Author", []string{"Tech"}),
			}

			processedItems, err := formatter.Process(items)

			if err != nil {
				t.Errorf("Unexpected error for format %s: %v", format, err)
			}

			expectedDate := now.Format(format)
			if processedItems[0].GetPubDate() != expectedDate {
				t.Errorf("Date format mismatch for format %s: expected %s, got %s",
					format, expectedDate, processedItems[0].GetPubDate())
			}
		})
	}
}

func TestTimeFormatter_Process_InvalidDate(t *testing.T) {
	formatter := NewTimeFormatter(time.RFC3339)

	invalidDate := "InvalidDateString"
	items := []*rss.Item{
		rss.NewItem("Test Title", "http://example.com", "Test Description", invalidDate, "Author", []string{"Tech"}),
	}

	processedItems, err := formatter.Process(items)

	if err != nil {
		t.Errorf("Unexpected error for invalid date: %v", err)
	}

	if processedItems[0].GetPubDate() != invalidDate {
		t.Errorf("Expected unchanged date for invalid input, got %s", processedItems[0].GetPubDate())
	}
}

// Бенчмарк-тест
func BenchmarkTimeFormatter_Process(b *testing.B) {
	now := time.Now()
	formattedDate := now.Format(time.RFC3339Nano)
	formatter := NewTimeFormatter(time.RubyDate)

	var items []*rss.Item

	for i := 0; i < 1_000; i++ {
		items = append(items, rss.NewItem("Test Title", "http://example.com", "Test Description", formattedDate, "Author", []string{"Tech"}))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		testItems := make([]*rss.Item, len(items))
		copy(testItems, items)
		_, _ = formatter.Process(testItems)
	}
}
