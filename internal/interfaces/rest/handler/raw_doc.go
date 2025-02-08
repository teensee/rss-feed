package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type RawDocHandler struct {
}

func NewRawDocHandler() *RawDocHandler {
	return &RawDocHandler{}
}

func (r2 RawDocHandler) Handle(w http.ResponseWriter, r *http.Request) {
	docType := chi.URLParam(r, "type")
	contentType := "application/yaml"

	if docType == "" {
		docType = "yaml"
	}

	currWd, err := os.Getwd()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	content, err := os.ReadFile(fmt.Sprintf("%s/docs/feed-api.yaml", currWd))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(content)
}
