package adapters

import (
	"rss-feed/internal/domain/rss"
	"rss-feed/internal/interfaces/rest/dto"
)

func ToRssResponseList(feed []*rss.Feed) []dto.FeedListResponse {
	var resp = make([]dto.FeedListResponse, 0, len(feed))

	for _, f := range feed {
		var newsList = make([]dto.FeedListItemResponse, 0, len(f.GetItems()))
		for _, item := range f.GetItems() {
			newsList = append(newsList, dto.FeedListItemResponse{
				Title:       item.GetTitle(),
				Link:        item.GetLink(),
				Description: item.GetDescription(),
				PubDate:     item.GetPubDate(),
			})
		}

		resp = append(resp, dto.FeedListResponse{
			Source: f.GetTitle(),
			Feed:   newsList,
		})
	}

	return resp
}
