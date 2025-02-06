package logger

import (
	"context"
	"log/slog"
	"rss-feed/internal/interfaces/rest/middleware"
)

type TraceIdSlogHandler struct {
	inner slog.Handler
}

func NewTraceIdSlogHandler(inner slog.Handler) *TraceIdSlogHandler {
	return &TraceIdSlogHandler{inner: inner}
}

func (t *TraceIdSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return t.inner.Enabled(ctx, level)
}

// nolint:gocritic // потому что в slog такой интерфейс
func (t *TraceIdSlogHandler) Handle(ctx context.Context, record slog.Record) error {
	if traceID, ok := ctx.Value(middleware.TraceIdHeader).(string); ok {
		record.Add(string(middleware.TraceIdHeader), traceID)
	}

	return t.inner.Handle(ctx, record)
}

func (t *TraceIdSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewTraceIdSlogHandler(t.inner.WithAttrs(attrs))
}

func (t *TraceIdSlogHandler) WithGroup(name string) slog.Handler {
	return NewTraceIdSlogHandler(t.inner.WithGroup(name))
}
