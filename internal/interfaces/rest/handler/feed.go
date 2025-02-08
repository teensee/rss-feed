package handler

import (
	"encoding/json"
	"net/http"
	"rss-feed/internal/application"
	"rss-feed/internal/domain/logging"
	"rss-feed/internal/interfaces/rest/adapters"
	"rss-feed/internal/interfaces/rest/dto"
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
	var req dto.GetFeedJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appReq, err := adapters.ToAppRssFeedRequest(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	feed, err := h.agg.AggregateFeedAsync(r.Context(), appReq)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(
		adapters.ToRssResponseList(feed),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
