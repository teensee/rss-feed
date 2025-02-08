package dto

type RssFeedItemProcess struct {
	rss     string
	filters []string
}

func NewRssFeedItemProcess(rss string, filters []string) *RssFeedItemProcess {
	return &RssFeedItemProcess{rss: rss, filters: filters}
}

func (item *RssFeedItemProcess) GetRss() string {
	return item.rss
}

func (item *RssFeedItemProcess) GetFilters() []string {
	return item.filters
}

type AppRssFeedRequest struct {
	items []*RssFeedItemProcess
}

func NewAppRssFeedRequest(items []*RssFeedItemProcess) *AppRssFeedRequest {
	return &AppRssFeedRequest{items: items}
}

func (r AppRssFeedRequest) GetItems() []*RssFeedItemProcess {
	return r.items
}
