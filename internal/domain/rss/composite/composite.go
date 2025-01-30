package composite

import (
	"context"
	"encoding/xml"
	"log/slog"
	"net/url"
	"rss-feed/pkg/cache"
	"rss-feed/pkg/http"
)

type ProcessedFeed struct {
	Source string     `json:"source"`
	Feed   []FeedItem `json:"feed"`
}

type FeedItem struct {
	Link  string
	Title string
}

type CompositeRss struct {
	feedList []url.URL
	client   *http.Client
	cache    cache.AppCache
	l        *slog.Logger
}

func NewCompositeRss(feedList []url.URL, client *http.Client, cache cache.AppCache, l *slog.Logger) *CompositeRss {
	return &CompositeRss{feedList: feedList, client: client, cache: cache, l: l}
}

func (r *CompositeRss) Feed(ctx context.Context) ([]ProcessedFeed, error) {
	var feeds []*http.RSS
	for _, feedUrl := range r.feedList {
		rss, err := r.fetch(ctx, feedUrl.String())

		if err != nil {
			return nil, err
		}

		feeds = append(feeds, rss)
	}

	processed := r.process(feeds)
	return processed, nil
}

func (r *CompositeRss) fetch(ctx context.Context, url string) (*http.RSS, error) {
	l := r.l.With(slog.String("url", url))
	l.InfoContext(ctx, "fetching feed")

	resp, err := r.client.GET(ctx, url, nil)
	if err != nil {
		l.ErrorContext(ctx, "fetch rss failed", slog.Any("err", err))
		return nil, err
	}

	var rss http.RSS
	err = xml.Unmarshal(resp, &rss)
	if err != nil {
		return nil, err
	}

	rss.XMLNSDC = "http://purl.org/dc/elements/1.1/"

	return &rss, nil
}

func (r *CompositeRss) process(rssList []*http.RSS) []ProcessedFeed {
	var processedFeedList []ProcessedFeed

	for _, rss := range rssList {
		items := *rss.Channel.Items
		feed := make([]FeedItem, 0, len(items))

		for _, i := range items {
			feed = append(feed, FeedItem{
				Title: i.Title,
				Link:  i.Link,
			})
		}

		processedFeedList = append(processedFeedList, ProcessedFeed{
			Source: rss.Channel.Title,
			Feed:   feed,
		})
	}

	return processedFeedList
}
