package httpapi

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

type Handler struct {
	geoIP countryLookup
}

type checkRequest struct {
    IP                string   `json:"ip"`
    AllowedCountries  []string `json:"allowedCountries"`
}

type checkResponse struct {
	Allowed bool   `json:"allowed"`
	Country string `json:"country"`
}

type countryLookup interface {
    Lookup(ip net.IP) (*geoip2.Country, error)
}

func NewHandler(geoIP countryLookup) *Handler {
    return &Handler{
        geoIP: geoIP,
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

	allowed, countryISO, err := doCheck(h.geoIP, parsedIP, req.AllowedCountries)
	if err != nil {
		slog.Error("Error looking up Country by IP", "error", err)
		http.Error(w, "Error looking up Country by IP", http.StatusInternalServerError)
		return
	}
	slog.Info("IP allowlist status", "allowed", allowed, "countryISO", countryISO)

	response := checkResponse{
		Allowed: allowed,
		Country: countryISO,
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

func doCheck(lookup countryLookup, ip net.IP, allowedCountries []string) (bool, string, error) {
	// perform lookup in the geoip database
	geoIPRecord, err := lookup.Lookup(ip)
	if err != nil {
		return false, "", err
	}

	allowed := isCountryAllowed(geoIPRecord.Country.IsoCode, allowedCountries)

	return allowed, geoIPRecord.Country.IsoCode, nil
}

func isCountryAllowed(countryISO string, allowedCountries []string) bool {
	// make this check case-insenstive
	uppercased := make([]string, len(allowedCountries))
	for i, str := range allowedCountries {
		uppercased[i] = strings.ToUpper(str)
	}

    return slices.Contains(uppercased, countryISO)
}