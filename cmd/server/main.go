package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	accessv1 "github.com/mike-faulcon/ip-access-service/gen/access/v1"
	"github.com/mike-faulcon/ip-access-service/internal/config"
	"github.com/mike-faulcon/ip-access-service/internal/geoip"
	"github.com/mike-faulcon/ip-access-service/internal/grpcapi"
	"github.com/mike-faulcon/ip-access-service/internal/httpapi"
	"github.com/mike-faulcon/ip-access-service/internal/service"
)

func main() {
	// Initialize Config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	slog.Info("config loaded",
		"http_port", cfg.HTTPPort,
		"grpc_port", cfg.GRPCPort,
		"geoip_path", cfg.GeoIPPath,
	)

	// Initialize GeoIP reader
	geoIPReader, err := geoip.NewGeoIPReader(cfg.GeoIPPath)
	if err != nil {
		slog.Error("Failed to initialize GeoIP reader", "error", err)
		os.Exit(1) // TODO: abort or let the service run in a partially initialized state?
	}
	defer geoIPReader.Close()

	// Define API routes using the standard http.ServeMux
	mux := http.NewServeMux()

	// Setup health endpoint
	mux.HandleFunc("/health", httpapi.GetHealthHandler)

	accessService := service.NewAccessService(geoIPReader)

	// Enable the ip-check handler to get the access service instance via dependency injection
	httpHandler := httpapi.NewHandler(accessService)

	mux.HandleFunc("POST /v1/check", httpHandler.PostCheckIPHandler)

	// Wrap the mux with the logging middleware
	loggedMux := loggingMiddleware(mux)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: loggedMux,
	}

	grpcHandler := grpcapi.NewServer(accessService)
	grpcServer := grpc.NewServer()
	accessv1.RegisterAccessServiceServer(grpcServer, grpcHandler)
	reflection.Register(grpcServer)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		slog.Error("Failed to listen for gRPC", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("HTTP server starting", "addr", httpServer.Addr)

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "error", err)
		}
	}()

	go func() {
		slog.Info("gRPC server starting", "addr", grpcListener.Addr().String())
		if err := grpcServer.Serve(grpcListener); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	<-ctx.Done()

	slog.Info("shutdown signal received")

	grpcServer.GracefulStop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown failed", "error", err)
	}

	slog.Info("HTTP & GRPC servers stopped")
}

// func getReadyHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("OK"))
// }
