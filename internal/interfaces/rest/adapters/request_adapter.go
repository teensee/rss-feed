package adapters

import (
	"net/url"
	appDto "rss-feed/internal/application/dto"
	"rss-feed/internal/interfaces/rest/dto"
)

func ToAppRssFeedRequest(req dto.RssFeedRequest) (*appDto.AppRssFeedRequest, error) {
	appReq := appDto.AppRssFeedRequest{}

	for _, item := range req.Items {
		u, err := url.Parse(item.Rss)
		if err != nil {
			return nil, err
		}

		appReq.Items = append(appReq.Items, &appDto.RssFeedItemProcess{
			Rss:     u.String(),
			Filters: item.Filters,
		})
	}

	return &appReq, nil
}
