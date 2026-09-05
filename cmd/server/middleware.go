package main

import ( 
	"log/slog"
    "net/http"
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

		next.ServeHTTP(recorder, r)

		// Check if the request resulted in a 404 Not Found
		if recorder.statusCode == http.StatusNotFound {
			slog.Warn("not found request",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
		} else {
			slog.Info("handled request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", recorder.statusCode,
			)
		}
	})
}