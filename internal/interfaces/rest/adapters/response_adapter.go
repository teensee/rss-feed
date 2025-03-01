package adapters

import (
	"rss-feed/internal/domain/rss"
	"rss-feed/internal/interfaces/rest/dto"
)

func ToRssResponseList(feed []*rss.Feed) dto.FeedListResponse {
	var items = make([]dto.FeedListObject, 0, len(feed))

	for _, f := range feed {
		var newsList = make([]dto.FeedListItemResponse, 0, len(f.GetItems()))
		for _, item := range f.GetItems() {
			newsList = append(newsList, dto.FeedListItemResponse{
				Title:       item.GetTitle(),
				Link:        item.GetLink(),
				Description: item.GetDescription(),
				PubDate:     item.GetPubDate(),
				Author:      item.GetCreator(),
			})
		}

		items = append(items, dto.FeedListObject{
			Source: f.GetLink(),
			Feed:   newsList,
		})
	}

	return dto.FeedListResponse{
		Items: items,
	}
}

func ToProcessorResponseList(processorList []string) dto.ProcessorListResponse {
	return dto.ProcessorListResponse{
		Items: processorList,
	}
}
