package logging

import (
	"context"
	"log/slog"
	"os"
)

type PortHandler struct {
	h    slog.Handler
	port string
}

func NewPortHandler(h slog.Handler, port string) *PortHandler {
	return &PortHandler{h: h, port: port}
}

func (ph *PortHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(slog.String("port", ph.port))
	return ph.h.Handle(ctx, r)
}

func (ph *PortHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PortHandler{h: ph.h.WithAttrs(attrs), port: ph.port}
}

func (ph *PortHandler) WithGroup(name string) slog.Handler {
	return &PortHandler{h: ph.h.WithGroup(name), port: ph.port}
}

func (ph *PortHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return ph.h.Enabled(ctx, level)
}

func InitLogging(port string) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	slogHandler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(NewPortHandler(slogHandler, port))
	slog.SetDefault(logger)
}
