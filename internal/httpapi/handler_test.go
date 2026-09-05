package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oschwald/geoip2-golang"
)

// ------------------------------------

func TestGetHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	if req == nil {
		t.Fatal("Request is nil")
	}

	recorder := httptest.NewRecorder()
	GetHealthHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", recorder.Code)
	}

	if recorder.Body.String() != "OK" {
		t.Errorf("Expected body OK, got %s", recorder.Body.String())
	}
}

// ------------------------------------

type mockCountryLookup struct {
    country string
    err     error
}

func (f mockCountryLookup) Lookup(net.IP) (*geoip2.Country, error) {
    if f.err != nil {
        return nil, f.err
    }
    return &geoip2.Country{
        Country: struct {
			Names             map[string]string `maxminddb:"names"`
            IsoCode           string            `maxminddb:"iso_code"`
            GeoNameID         uint              `maxminddb:"geoname_id"`
			IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
        }{IsoCode: f.country},
    }, nil
}

func TestPostCheckIPHandler(t *testing.T) {
	// Use Table Tests
	tests := []struct {
		name             string   // Clear scenario name
		ip               string   // Input IP address
		allowedCountries []string // Input list of allowed countries
		mockGeoIPCountry string
		mockGeoIPError   error
		wantStatus       int      // Expected output - http status code
		wantAllowed      bool     // Expected output - allowed or not
        wantCountry      string   // Expected output - matched country
	}{
		{
			name:             "Basic allowed case", 
			ip:               "142.251.152.119", 
			allowedCountries: []string{"US"},
			mockGeoIPCountry: "US",
			wantStatus:       http.StatusOK,
			wantAllowed:      true,
			wantCountry:      "US",
		},
		{
            name:             "Basic not allowed case",
            ip:               "142.251.152.119",
            allowedCountries: []string{"UK", "FR"},
            mockGeoIPCountry: "US",
            wantStatus:       http.StatusOK,
            wantAllowed:      false,
            wantCountry:      "US",
        },
		{
            name:             "invalid ip",
            ip:               "1.2.34",
            allowedCountries: []string{"US"},
            wantStatus:       http.StatusBadRequest,
        },
		{
            name:             "lookup error",
            ip:               "142.251.152.119",
            allowedCountries: []string{"US"},
            mockGeoIPError:   errors.New("db unavailable"),
            wantStatus:       http.StatusInternalServerError,
        },
		{
            name:             "Empty Country List",
            ip:               "142.251.152.119",
            allowedCountries: []string{},
            wantStatus:       http.StatusBadRequest,
        },
		{
            name:             "Nil Country List",
            ip:               "142.251.152.119",
            allowedCountries: nil,
            wantStatus:       http.StatusBadRequest,
        },
		// {"Case-insensitive allowed case", "142.251.152.119", ["us", "uk"], true}, // not yet implemented
		{
            name:             "Case-insensitive allowed case",
            ip:               "142.251.152.119",
            allowedCountries: []string{"us", "uk"},
            mockGeoIPCountry: "US",
            wantStatus:       http.StatusOK,
            wantAllowed:      true,
            wantCountry:      "US",
        },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHandler(mockCountryLookup{
                country: tc.mockGeoIPCountry,
                err:     tc.mockGeoIPError,
            })

			requstPayload := checkRequest{
				IP:               tc.ip,
				AllowedCountries: tc.allowedCountries,
			}

			body, err := json.Marshal(requstPayload)
            if err != nil {
                t.Fatal(err)
            }

			req := httptest.NewRequest(http.MethodPost, "/v1/check", bytes.NewReader(body))
            req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			h.PostCheckIPHandler(recorder, req)

			if recorder.Code != tc.wantStatus {
                t.Fatalf("status = %d, want %d, body = %s", recorder.Code, tc.wantStatus, recorder.Body.String())
            }

			if tc.wantStatus != http.StatusOK {
                return
            }

			var got checkResponse
            if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
                t.Fatal(err)
            }
            if got.Allowed != tc.wantAllowed {
                t.Errorf("allowed = %v, want %v", got.Allowed, tc.wantAllowed)
            }
            if got.Country != tc.wantCountry {
                t.Errorf("country = %q, want %q", got.Country, tc.wantCountry)
            }
		})
	}
}