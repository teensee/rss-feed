package processor_test

import (
	"rss-feed/internal/domain/rss"
	"rss-feed/internal/infrastructure/processor"
	"slices"
	"testing"
)

type mockProcessor struct {
	name string
}

func (m *mockProcessor) Process(items []*rss.Item) ([]*rss.Item, error) {
	return items, nil
}

func (m *mockProcessor) Name() string {
	return m.name
}

func TestNewProcessorRegistry(t *testing.T) {
	processors := []rss.Processor{
		&mockProcessor{name: "processor1"},
		&mockProcessor{name: "processor2"},
		&mockProcessor{name: "processor1"}, // <-- Дубликат, не должен быть добавлен второй раз
	}

	reg := processor.NewProcessorRegistry(processors)

	if len(reg.Names()) != 2 {
		t.Errorf("Expected 2 unique processors, got %d", len(reg.Names()))
	}
}

func TestResolve(t *testing.T) {
	processors := []rss.Processor{
		&mockProcessor{name: "processor1"},
		&mockProcessor{name: "processor2"},
	}

	reg := processor.NewProcessorRegistry(processors)

	proc, found := reg.Resolve("processor1")
	if !found || proc.Name() != "processor1" {
		t.Errorf("Expected to resolve processor1, got %v", proc)
	}

	_, found = reg.Resolve("nonexistent")
	if found {
		t.Errorf("Expected nonexistent processor to not be found")
	}
}

func TestNames(t *testing.T) {
	processors := []rss.Processor{
		&mockProcessor{name: "b-processor"},
		&mockProcessor{name: "a-processor"},
		&mockProcessor{name: "c-processor"},
	}

	reg := processor.NewProcessorRegistry(processors)
	expected := []string{"a-processor", "b-processor", "c-processor"}
	actual := reg.Names()

	if !slices.Equal(expected, actual) {
		t.Errorf("Expected sorted names %v, got %v", expected, actual)
	}
}
