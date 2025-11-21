package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	RequestIDKey contextKey = "request_id"
	TraceIDKey   contextKey = "trace_id"
	SpanIDKey    contextKey = "span_id"
)

func TraceMiddleware(logger interface {
	Info(context.Context, string, ...zap.Field)
}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			traceID := r.Header.Get("X-Trace-ID")
			if traceID == "" {
				traceID = uuid.New().String()
			}

			spanID := uuid.New().String()

			ctx := r.Context()
			ctx = context.WithValue(ctx, RequestIDKey, requestID)
			ctx = context.WithValue(ctx, TraceIDKey, traceID)
			ctx = context.WithValue(ctx, SpanIDKey, spanID)

			w.Header().Set("X-Request-ID", requestID)
			w.Header().Set("X-Trace-ID", traceID)
			w.Header().Set("X-Span-ID", spanID)

			next.ServeHTTP(w, r.WithContext(ctx))

			duration := time.Since(start)
			if logger != nil {
				logger.Info(ctx, "request completed",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("request_id", requestID),
					zap.String("trace_id", traceID),
					zap.Duration("duration", duration),
				)
			}
		})
	}
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

func GetTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(TraceIDKey).(string); ok {
		return id
	}
	return ""
}
