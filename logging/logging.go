// logging package provides a custom slog.Handler that adds the port to each log record
package logging

import (
	"context"
	"log/slog"
	"os"
)

// PortHandler is a custom slog.Handler that adds the port to each log record
type PortHandler struct {
	h    slog.Handler
	port string
}

// NewPortHandler creates a new PortHandler with the given slog.Handler and port
// It is used to add the port information to each log record
func NewPortHandler(h slog.Handler, port string) *PortHandler {
	return &PortHandler{h: h, port: port}
}

// Handle adds the port information to the log record and passes it to the underlying handler
// It is used to ensure that all log records contain the port information
func (ph *PortHandler) Handle(ctx context.Context, r slog.Record) error {
	r.AddAttrs(slog.String("port", ph.port))
	return ph.h.Handle(ctx, r)
}

// WithAttrs adds additional attributes to the log record
// It is used to enrich the log record with more context
func (ph *PortHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PortHandler{h: ph.h.WithAttrs(attrs), port: ph.port}
}

// WithGroup adds a group to the log record
// It is used to categorize log records into groups
func (ph *PortHandler) WithGroup(name string) slog.Handler {
	return &PortHandler{h: ph.h.WithGroup(name), port: ph.port}
}

// Enabled checks if the log level is enabled for the handler
func (ph *PortHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return ph.h.Enabled(ctx, level)
}

// InitLogging initializes the logging system with a custom handler that includes the port in each log record
func InitLogging(port string) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	slogHandler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(NewPortHandler(slogHandler, port))
	slog.SetDefault(logger)
}
