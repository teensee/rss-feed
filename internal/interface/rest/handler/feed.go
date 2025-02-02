package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"rss-feed/internal/application"
	appDto "rss-feed/internal/application/dto"
	"rss-feed/internal/domain/logging"
	"rss-feed/internal/interface/rest/dto"
)

type FeedHandler struct {
	agg *application.FeedAggregator
	l   logging.Logger
}

func NewFeedHandler(agg *application.FeedAggregator, l logging.Logger) *FeedHandler {
	return &FeedHandler{
		agg: agg,
		l:   l,
	}
}

func (h *FeedHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req dto.RssFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appReq := appDto.AppRssFeedRequest{}

	for _, item := range req.Items {
		u, err := url.Parse(item.Rss)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		appReq.Items = append(appReq.Items, &appDto.RssFeedItemProcess{
			Rss:     u.String(),
			Filters: item.Filters,
		})

	}

	feed, err := h.agg.AggregateFeedAsync(r.Context(), appReq)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// todo: move to mapper
	var resp = make([]dto.ProcessedFeedResponse, 0, len(feed))
	for _, f := range feed {

		var newsList = make([]dto.FeedItemResponse, 0, len(f.Channel.Items))
		for _, item := range f.Channel.Items {
			newsList = append(newsList, dto.FeedItemResponse{
				Title:   item.Title,
				Link:    item.Link,
				Desc:    item.Description,
				PubDate: item.PubDate,
			})
		}

		resp = append(resp, dto.ProcessedFeedResponse{
			Source: f.Channel.Title,
			Feed:   newsList,
		})
	}

	body, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)

	return
}
