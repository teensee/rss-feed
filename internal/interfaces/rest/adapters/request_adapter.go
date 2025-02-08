package adapters

import (
	"net/url"
	appDto "rss-feed/internal/application/dto"
	"rss-feed/internal/interfaces/rest/dto"
)

func ToAppRssFeedRequest(req dto.GetFeedJSONRequestBody) (*appDto.AppRssFeedRequest, error) {
	var items = make([]*appDto.RssFeedItemProcess, 0, len(req.Items))

	for _, item := range req.Items {
		u, err := url.Parse(item.Rss)
		if err != nil {
			return nil, err
		}

		items = append(items, appDto.NewRssFeedItemProcess(u.String(), item.Filters))
	}

	return appDto.NewAppRssFeedRequest(items), nil
}
