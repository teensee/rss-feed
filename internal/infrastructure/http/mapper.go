package http

import (
	"rss-feed/internal/domain/rss"
	"rss-feed/internal/infrastructure/http/dto"
)

func toDomainModel(r *dto.RSS) *rss.Feed {
	var image *rss.Image

	if r.Channel.Image != nil {
		image = rss.NewImage(r.Channel.Image.Link, r.Channel.Image.URL, r.Channel.Image.Title)
	}

	var modelItems []*rss.Item

	if r.Channel.Items != nil && len(*r.Channel.Items) > 0 {
		for _, item := range *r.Channel.Items {
			modelItems = append(modelItems, rss.NewItem(
				item.Title,
				item.Link,
				item.Description,
				item.PubDate,
				item.Creator,
				item.Categories,
			))
		}
	}

	return rss.NewFeed(
		r.Channel.Title,
		r.Channel.Link,
		r.Channel.Description,
		r.Channel.PubDate,
		image,
		modelItems,
	)
}
