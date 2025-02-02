package rss

type Processor interface {
	Name() string
	Process(items []*Item) ([]*Item, error)
}

type ProcessorRegistry interface {
	Resolve(name string) (Processor, bool)
	Names() []string
}
