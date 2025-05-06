package logging

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// should not use built-in type string as key for value; define your own type to avoid collisions
type contextKey string

const traceIDKey contextKey = "TraceID"

// TraceIDMiddleware adds a TraceID to the context for each request
func TraceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if TraceID exists in the headers
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			// Generate a new TraceID if none exists in the headers
			traceID = uuid.New().String()
			slog.Info("Generated new TraceID", "TraceID", traceID)
		} else {
			slog.Info("Using existing TraceID from headers", "TraceID", traceID)
		}

		// Add the TraceID to the context
		ctx := context.WithValue(r.Context(), traceIDKey, traceID)

		// Pass the updated context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetTraceID retrieves the TraceID from the context
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}
