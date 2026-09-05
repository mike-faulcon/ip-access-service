package main

import ( 
	"fmt"
	"log/slog"
    "net/http"
	"os"

	"ip-access-service/internal/geoip"
	"ip-access-service/internal/httpapi"
 )

 const (
	Port = 8080
	GeoIPPath = "../../data/GeoLite2-Country.mmdb"
 )

func main() {
	// Initialize GeoIP reader
	geoIP, err := geoip.NewGeoIPReader(GeoIPPath)
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
		Addr:    fmt.Sprintf(":%d", Port),
		Handler: loggedMux,
	}

	fmt.Printf("Server is running on port %d\n", Port)
	server.ListenAndServe()
}

// func getReadyHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("OK"))
// }

// func getMetricsHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("OK"))
// }
