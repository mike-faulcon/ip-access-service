package httpapi

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"

	"ip-access-service/internal/service"
)

type Handler struct {
	accessService *service.AccessService
}

type checkRequest struct {
    IP                string   `json:"ip"`
    AllowedCountries  []string `json:"allowedCountries"`
}

type checkResponse struct {
	Allowed bool   `json:"allowed"`
	Country string `json:"country"`
}


func NewHandler(accessService *service.AccessService) *Handler {
    return &Handler{
        accessService: accessService,
    }
}

func GetHealthHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Health check request received")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) PostCheckIPHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Check IP request received")

	var req checkRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

	// slog.Info("Request", "ip", req.IP, "allow-list", req.AllowedCountries)
	slog.Info("Request", "props", req)

	// validate allowed countries list
	if len(req.AllowedCountries) == 0 {
		slog.Warn("Empty AllowedCountries List")
		http.Error(w, "Empty AllowedCountries List", http.StatusBadRequest)
		return
	}

	// validate IP request
	parsedIP := net.ParseIP(req.IP)
	if parsedIP == nil {
		slog.Warn("Invalid IP", "IP", req.IP)
		http.Error(w, "Invalid IP", http.StatusBadRequest)
		return
	}
	slog.Info("validated IP", "ip", parsedIP)

	accessResult, err := h.accessService.Check(r.Context(), parsedIP, req.AllowedCountries)
	if err != nil {
		slog.Error("Error performing access check", "error", err)
		http.Error(w, "Error performing access check", http.StatusInternalServerError)
		return
	}
	slog.Info("IP allowlist status", "allowed", accessResult.Allowed, "countryISO", accessResult.CountryISO)

	response := checkResponse{
		Allowed: accessResult.Allowed,
		Country: accessResult.CountryISO,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("Error encoding response", "error", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
