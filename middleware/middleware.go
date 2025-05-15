package middleware

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

// should not use built-in type string as key for value; define your own type to avoid collisions
type contextKey string

// define a custom type for the UserID key to avoid collisions
const UserIDKey contextKey = "UserID"

const traceIDKey contextKey = "TraceID"

// PortKey is the context key for the server port
const PortKey contextKey = "Port"

// TraceIDMiddleware adds a TraceID to the context for each request
func TraceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
			slog.Info("Generated new TraceID", "TraceID", traceID)
		} else {
			slog.Info("Using existing TraceID from headers", "TraceID", traceID)
		}

		ctx := context.WithValue(r.Context(), traceIDKey, traceID)

		slog.Info("TraceIDMiddleware executed", "TraceID", traceID)
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

// GetUserID retrieves the UserID from the context
func GetUserID(ctx context.Context) (int, error) {
	if userID, ok := ctx.Value(UserIDKey).(int); ok {
		return userID, nil
	}
	return 0, errors.New("userID not in the context")
}

// UserIDMiddleware extracts the UserID from the request, validates it as an integer, and adds it to the context
func UserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.Header.Get("X-User-ID")
		if userIDStr == "" {
			http.Error(w, "UserID is required", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "UserID must be a valid integer", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, userID)
		r = r.WithContext(ctx)

		slog.Info("UserIDMiddleware executed", "UserID", userID)
		next.ServeHTTP(w, r)
	})
}

// LoadBalancerMiddleware routes requests to one of three servers based on UserID
func LoadBalancerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetUserID(r.Context())
		if err != nil {
			http.Error(w, "UserID is missing or invalid", http.StatusBadRequest)
			return
		}

		serverIndex := userID % 3
		var serverAddr string
		switch serverIndex {
		case 0:
			serverAddr = "localhost:8081"
		case 1:
			serverAddr = "localhost:8082"
		case 2:
			serverAddr = "localhost:8083"
		}

		var requestBody string
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			} else {
				slog.Error("Failed to read request body", "error", err)
			}
		}

		slog.Info("Routing request to server",
			"ServerAddress", serverAddr,
			"UserID", userID,
			"Method", r.Method,
			"URL", r.URL.String(),
			"Body", requestBody)

		proxy := httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   serverAddr,
		})
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			slog.Error("Proxy error", "error", err)
			http.Error(w, "Bad Gateway: Unable to connect to the target server", http.StatusBadGateway)
		}

		slog.Info("LoadBalancerMiddleware executed", "ServerIndex", serverIndex)
		proxy.ServeHTTP(w, r)
	})
}

// GetPort retrieves the server port from the context
func GetPort(ctx context.Context) string {
	if port, ok := ctx.Value(PortKey).(string); ok {
		return port
	}
	return ""
}

// ChainMiddleware applies a list of middleware functions to a handler
func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
