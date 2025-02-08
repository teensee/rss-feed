package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type consoleColor string

const (
	Reset   consoleColor = "\033[0m"
	Red     consoleColor = "\033[31m"
	Green   consoleColor = "\033[32m"
	Yellow  consoleColor = "\033[33m"
	Blue    consoleColor = "\033[34m"
	Magenta consoleColor = "\033[35m"
	Cyan    consoleColor = "\033[36m"
	Gray    consoleColor = "\033[37m"
	White   consoleColor = "\033[97m"
)

type SlogPrettyHandler struct {
	w        io.Writer
	opts     *slog.HandlerOptions
	colorize bool
	attrs    []slog.Attr
	groups   []string
}

func NewSlogPrettyHandler(w io.Writer, opts *slog.HandlerOptions, colorize bool) *SlogPrettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	return &SlogPrettyHandler{
		w:        w,
		opts:     opts,
		colorize: colorize,
	}
}

func (s *SlogPrettyHandler) Enabled(_ context.Context, l slog.Level) bool {
	minLevel := slog.LevelInfo
	if s.opts.Level != nil {
		minLevel = s.opts.Level.Level()
	}

	return l >= minLevel
}

func (s *SlogPrettyHandler) Handle(_ context.Context, r slog.Record) error { // nolint:gocritic // потому что в slog такой интерфейс
	timeString := r.Time.Format("2006-01-02 15:04:05.00000")

	logLine := fmt.Sprintf("%s %s %s", timeString, s.formatLevel(r.Level), r.Message)

	if len(s.groups) > 0 {
		logLine = strings.Join(s.groups, ".") + "." + logLine
	}

	var attrs = make([]string, 0, r.NumAttrs())

	r.Attrs(
		func(attr slog.Attr) bool {
			attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value.Any()))
			return true
		},
	)

	if len(attrs) > 0 {
		logLine += " | " + strings.Join(attrs, " | ")
	}

	_, err := fmt.Fprintln(s.w, logLine)
	if err != nil {
		return fmt.Errorf("write log failed: %w", err)
	}

	return nil
}

func (s *SlogPrettyHandler) formatLevel(l slog.Level) string {
	if !s.colorize {
		return fmt.Sprintf("[%s]", l.String())
	}

	var color consoleColor

	switch l {
	case slog.LevelDebug:
		color = Blue
	case slog.LevelInfo:
		color = Green
	case slog.LevelWarn:
		color = Yellow
	case slog.LevelError:
		color = Red
	default:
		color = Reset
	}

	return fmt.Sprintf("%s[%s]%s", color, l.String(), Reset)
}

func (s *SlogPrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SlogPrettyHandler{
		w:        s.w,
		opts:     s.opts,
		colorize: s.colorize,
		attrs:    append(s.attrs, attrs...),
		groups:   s.groups,
	}
}

func (s *SlogPrettyHandler) WithGroup(name string) slog.Handler {
	return &SlogPrettyHandler{
		w:        s.w,
		opts:     s.opts,
		colorize: s.colorize,
		attrs:    s.attrs,
		groups:   append(s.groups, name),
	}
}
