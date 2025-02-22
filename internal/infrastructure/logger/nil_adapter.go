package logger

import (
	"context"
	"rss-feed/internal/domain/logging"
)

var _ logging.Logger = &NilAdapter{}

type NilAdapter struct{}

func NewNilAdapter() *NilAdapter {
	return &NilAdapter{}
}

func (s *NilAdapter) Debug(_ context.Context, _ string, _ ...any) {
}

func (s *NilAdapter) Info(_ context.Context, _ string, _ ...any) {
}

func (s *NilAdapter) Warn(_ context.Context, _ string, _ ...any) {
}

func (s *NilAdapter) Error(_ context.Context, _ string, _ ...any) {
}

func (s *NilAdapter) With(_ ...any) logging.Logger {
	return s
}
