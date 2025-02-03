package handler

import (
	"encoding/json"
	"net/http"
	"rss-feed/internal/domain/rss"
)

type ProcessorListHandler struct {
	registry rss.ProcessorRegistry
}

func NewProcessorListHandler(registry rss.ProcessorRegistry) *ProcessorListHandler {
	return &ProcessorListHandler{registry: registry}
}

func (p *ProcessorListHandler) Handle(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(p.registry.Names())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(res)
}
