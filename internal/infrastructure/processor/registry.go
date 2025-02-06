package processor

import (
	"rss-feed/internal/domain/rss"
	"slices"
)

var _ rss.ProcessorRegistry = &Registry{}

type Registry struct {
	registry map[string]rss.Processor
	names    []string
}

func NewProcessorRegistry(processorList []rss.Processor) *Registry {
	var registry = make(map[string]rss.Processor, len(processorList))

	var names = make([]string, 0, len(processorList))

	for _, processor := range processorList {
		if _, ok := registry[processor.Name()]; !ok {
			registry[processor.Name()] = processor
			names = append(names, processor.Name())
		}
	}

	slices.Sort(names)

	return &Registry{registry: registry, names: names}
}

func (p *Registry) Resolve(name string) (rss.Processor, bool) {
	proc, ok := p.registry[name]
	return proc, ok
}

func (p *Registry) Names() []string {
	return p.names
}
