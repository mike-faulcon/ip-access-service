package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
    "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ip-access-service/internal/config"
	"ip-access-service/internal/geoip"
	"ip-access-service/internal/httpapi"
)


func main() {
	// Initialize Config
	cfg := config.Load()

	// Initialize GeoIP reader
	geoIP, err := geoip.NewGeoIPReader(cfg.GeoIPPath)
    if err != nil {
        slog.Error("Failed to initialize GeoIP reader", "error", err)
		os.Exit(1) // TODO: abort or let the service run in a partially initialized state?
    }
    defer geoIP.Close()

	// Define API routes using the standard http.ServeMux
	mux := http.NewServeMux()

	// Setup health endpoint
	mux.HandleFunc("/health", httpapi.GetHealthHandler)
	
	// Enable the ip-check handler to get the geoIP reader via dependency injection
	h := httpapi.NewHandler(geoIP)
	mux.HandleFunc("POST /v1/check", h.PostCheckIPHandler) 

	// Wrap the mux with the logging middleware
	loggedMux := loggingMiddleware(mux)

	// Start the server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: loggedMux,
	}

	ctx, stop := signal.NotifyContext(
        context.Background(),
        syscall.SIGINT,
        syscall.SIGTERM,
    )
    defer stop()

    go func() {
        slog.Info("HTTP server starting", "addr", server.Addr)

        if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "error", err)
        }
    }()

    <-ctx.Done()

    slog.Info("shutdown signal received")

    shutdownCtx, cancel := context.WithTimeout(
        context.Background(),
        5*time.Second,
    )
    defer cancel()

    if err := server.Shutdown(shutdownCtx); err != nil {
        slog.Error("HTTP server shutdown failed", "error", err)
    }

    slog.Info("HTTP server stopped")
}

// func getReadyHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("OK"))
// }

// func getMetricsHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("OK"))
// }
