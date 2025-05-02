package logging

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const traceIDKey contextKey = "TraceID"

// TraceIDMiddleware adds a TraceID to the context for each request
func TraceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		traceID := uuid.New().String()

		ctx := context.WithValue(r.Context(), traceIDKey, traceID)

		slog.Info("Generated TraceID", "TraceID", traceID)

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
