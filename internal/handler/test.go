package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
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
	rss, err := p.rss.FeedAsync(
		r.Context(),
		[]*url.URL{
			{
				Scheme:   "https",
				Host:     "habr.com",
				Path:     "/ru/rss/articles/top/daily",
				RawQuery: "limit=5",
			},
			//{
			//	Scheme:   "https",
			//	Host:     "lenta.ru",
			//	Path:     "/rss/google-newsstand/main",
			//	RawQuery: "limit=5",
			//},
			//{
			//	Scheme: "https",
			//	Host:   "lenta.ru",
			//	Path:   "/rss/last24",
			//},
			//{
			//	Scheme: "https",
			//	Host:   "dtf.ru",
			//	Path:   "/rss/all",
			//},
			{
				Scheme: "https",
				Host:   "rssexport.rbc.ru",
				Path:   "/rbcnews/news/5/full.rss",
			},
		},
	)

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
