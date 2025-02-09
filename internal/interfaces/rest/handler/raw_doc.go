package handler

import (
	"fmt"
	"net/http"
	"os"
)

type RawDocHandler struct {
}

func NewRawDocHandler() *RawDocHandler {
	return &RawDocHandler{}
}

func (r2 RawDocHandler) Handle(w http.ResponseWriter, r *http.Request) {
	contentType := "application/yaml"

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
	_, _ = w.Write(content)
}
