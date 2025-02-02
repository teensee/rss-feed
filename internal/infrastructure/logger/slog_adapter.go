package logger

import (
	"context"
	"log/slog"
	"rss-feed/internal/domain/logging"
)

var _ logging.Logger = &SlogAdapter{}

type SlogAdapter struct {
	l *slog.Logger
}

func NewSlogAdapter(l *slog.Logger) *SlogAdapter {
	return &SlogAdapter{l: l}
}

func (s *SlogAdapter) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	s.l.Log(ctx, level, msg, args...)
}

func (s *SlogAdapter) Debug(ctx context.Context, msg string, args ...any) {
	s.log(ctx, slog.LevelDebug, msg, args...)
}

func (s *SlogAdapter) Info(ctx context.Context, msg string, args ...any) {
	s.log(ctx, slog.LevelInfo, msg, args...)
}

func (s *SlogAdapter) Warn(ctx context.Context, msg string, args ...any) {
	s.log(ctx, slog.LevelWarn, msg, args...)
}

func (s *SlogAdapter) Error(ctx context.Context, msg string, args ...any) {
	s.log(ctx, slog.LevelError, msg, args...)
}

func (s *SlogAdapter) With(args ...any) logging.Logger {
	return &SlogAdapter{l: s.l.With(args...)}
}
