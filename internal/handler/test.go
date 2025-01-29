package handler

import (
	"encoding/xml"
	"net/http"
	"rss-feed/internal/domain/rss/habr"
)

type PingHandler struct {
	habr *habr.Habr
}

func NewPingHandler(habr *habr.Habr) *PingHandler {
	return &PingHandler{habr: habr}
}

func (p *PingHandler) Handle(w http.ResponseWriter, r *http.Request) {
	rss, err := p.habr.FeedCached(r.Context(), habr.NewNewFeedUrl(habr.Rate25, habr.Easy))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	body, err := xml.Marshal(rss)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write(body)

	return
}
