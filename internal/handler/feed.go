package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"rss-feed/internal/domain/rss/composite"
)

type RssFeed struct {
	Rss     string   `json:"rss"`
	Filters []string `json:"filters"`
}

type RssFeedRequest struct {
	Items []RssFeed `json:"items"`
}

type FeedHandler struct {
	rss *composite.CompositeRss
	l   *slog.Logger
}

func NewFeedHandler(rss *composite.CompositeRss, l *slog.Logger) *FeedHandler {
	return &FeedHandler{
		rss: rss,
		l:   l,
	}
}

func (h *FeedHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req RssFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var urls []*url.URL
	for _, item := range req.Items {
		u, err := url.Parse(item.Rss)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		urls = append(urls, u)
	}

	rss, err := h.rss.FeedAsync(
		r.Context(),
		urls,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(rss)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)

	return
}
