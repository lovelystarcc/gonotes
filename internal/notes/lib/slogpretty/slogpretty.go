package slogpretty

import (
	"context"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	l     *log.Logger
	attrs []slog.Attr
}

func (opts PrettyHandlerOptions) NewPrettyHandler(out io.Writer) *PrettyHandler {
	return &PrettyHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.LevelDebug
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	timeStr := color.YellowString(r.Time.Format("[15:04:05]"))
	msg := color.CyanString(r.Message)

	var kvs []string
	r.Attrs(func(a slog.Attr) bool {
		k := color.GreenString(a.Key)
		v := color.WhiteString("%v", a.Value.Any())
		kvs = append(kvs, k+"="+v)
		return true
	})
	for _, a := range h.attrs {
		k := color.GreenString(a.Key)
		v := color.WhiteString("%v", a.Value.Any())
		kvs = append(kvs, k+"="+v)
	}

	if len(kvs) > 0 {
		h.l.Println(timeStr, level, msg, kvs)
	} else {
		h.l.Println(timeStr, level, msg)
	}

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := append(append([]slog.Attr{}, h.attrs...), attrs...)
	return &PrettyHandler{
		Handler: h.Handler.WithAttrs(attrs),
		l:       h.l,
		attrs:   newAttrs,
		opts:    h.opts,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
		attrs:   h.attrs,
		opts:    h.opts,
	}
}
