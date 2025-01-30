package handler

import (
	"encoding/json"
	"net/http"
	"rss-feed/internal/domain/rss/composite"
	"rss-feed/internal/domain/rss/habr"
)

type PingHandler struct {
	habr *habr.Habr
	rss  *composite.CompositeRss
}

func NewPingHandler(habr *habr.Habr, rss *composite.CompositeRss) *PingHandler {
	return &PingHandler{habr: habr, rss: rss}
}

func (p *PingHandler) Handle(w http.ResponseWriter, r *http.Request) {
	rss, err := p.rss.Feed(r.Context())
	//rss, err := p.habr.FeedCached(r.Context(), habr.NewNewFeedUrl(habr.Rate25, habr.Easy))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	body, err := json.Marshal(rss)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)

	return
}
