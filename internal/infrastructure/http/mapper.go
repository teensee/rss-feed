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
			modelItems = append(modelItems, &rss.Item{
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				PubDate:     item.PubDate,
				Creator:     item.Creator,
				Categories:  item.Categories,
			})
		}
	}

	return &rss.Feed{
		XMLName: r.XMLName.Local,
		Version: r.Version,
		XMLNSDC: r.XMLNSDC,
		Channel: &rss.Channel{
			Title:   r.Channel.Title,
			Link:    r.Channel.Link,
			PubDate: r.Channel.PubDate,
			Image:   image,
			Items:   modelItems,
		},
	}
}
