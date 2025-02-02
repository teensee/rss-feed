package processor

import "rss-feed/internal/domain/rss"

var _ rss.ProcessorRegistry = &Registry{}

type Registry struct {
	registry map[string]rss.Processor
}

func NewProcessorRegistry(processorList []rss.Processor) *Registry {
	var registry = make(map[string]rss.Processor, len(processorList))

	for _, processor := range processorList {
		if _, ok := registry[processor.Name()]; !ok {
			registry[processor.Name()] = processor
		}
	}

	return &Registry{registry: registry}
}

func (p *Registry) Resolve(name string) (rss.Processor, bool) {
	proc, ok := p.registry[name]
	return proc, ok
}

func (p *Registry) Names() []string {
	names := make([]string, 0, len(p.registry))
	for k := range p.registry {
		names = append(names, k)
	}
	return names

}
