package dto

import "rss-feed/internal/application/dto"

type RssFeed struct {
	Rss     string   `json:"rss"`
	Filters []string `json:"filters"`
}

type RssFeedRequest struct {
	Items []RssFeed `json:"items"`
}

func (r *RssFeedRequest) ToApplicationModel() *dto.AppRssFeedRequest {
	appItems := make([]*dto.RssFeedItemProcess, 0, len(r.Items))

	for _, item := range r.Items {
		appItems = append(appItems, &dto.RssFeedItemProcess{
			Rss:     item.Rss,
			Filters: item.Filters,
		})
	}
	return &dto.AppRssFeedRequest{
		Items: appItems,
	}
}

type ProcessedFeedResponse struct {
	Source string             `json:"source"`
	Feed   []FeedItemResponse `json:"feed"`
}

type FeedItemResponse struct {
	Link    string `json:"link,omitempty"`
	Title   string `json:"title,omitempty"`
	Desc    string `json:"desc,omitempty"`
	PubDate string `json:"pubDate,omitempty"`
}
