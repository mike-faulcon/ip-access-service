package main

import (
	"log/slog"
	"net/http"
	"time"
)

// responseRecorder captures the status code written by the handler
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code written by the handler
func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware logs the request and response
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default if WriteHeader isn't called
		}

		start := time.Now()
		next.ServeHTTP(recorder, r)
		duration := time.Since(start)

		attrs := []any{
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.statusCode,
			"duration_ms", duration.Milliseconds(),
			"remote_addr", r.RemoteAddr,
		}

		switch {
		case recorder.statusCode >= 500:
			slog.Error("request completed", attrs...)
		case recorder.statusCode >= 400:
			slog.Warn("request completed", attrs...)
		default:
			slog.Info("request completed", attrs...)
		}
	})
}
