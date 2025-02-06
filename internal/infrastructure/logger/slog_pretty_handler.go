package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

type SlogPrettyHandler struct {
	w    io.Writer
	opts *slog.HandlerOptions
}

func NewSlogPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *SlogPrettyHandler {
	return &SlogPrettyHandler{w: w, opts: opts}
}

func (s *SlogPrettyHandler) Enabled(_ context.Context, l slog.Level) bool {
	minLevel := slog.LevelInfo
	if s.opts.Level != nil {
		minLevel = s.opts.Level.Level()
	}

	return l >= minLevel
}

func (s *SlogPrettyHandler) Handle(ctx context.Context, r slog.Record) error { // nolint:gocritic // потому что в slog такой интерфейс
	tstr := r.Time.Format("2006-01-02 15:04:05.00000")
	lvl := r.Level.String()

	attrs := make([]string, 0)

	r.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value.Any()))
		return true
	})

	logLine := fmt.Sprintf("%s [%s]: %s", tstr, lvl, r.Message)
	if len(attrs) > 0 {
		logLine += " " + fmt.Sprint(attrs)
	}

	_, err := fmt.Fprintln(s.w, logLine)

	return err
}

func (s *SlogPrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewSlogPrettyHandler(s.w, s.opts)
}

func (s *SlogPrettyHandler) WithGroup(name string) slog.Handler {
	return NewSlogPrettyHandler(s.w, s.opts)
}
