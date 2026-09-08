package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mike-faulcon/ip-access-service/internal/service"
	"github.com/mike-faulcon/ip-access-service/internal/testutil"
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

func TestPostCheckIPHandler(t *testing.T) {
	tests := []struct {
		name             string   // Clear scenario name
		ip               string   // Input IP address
		allowedCountries []string // Input list of allowed countries
		body             []byte   // Input JSON payload
		mockGeoIPCountry string
		mockGeoIPError   error
		wantStatus       int    // Expected output - http status code
		wantAllowed      bool   // Expected output - allowed or not
		wantCountry      string // Expected output - matched country
	}{
		{
			name:             "basic allowed",
			ip:               "142.251.152.119",
			allowedCountries: []string{"US"},
			mockGeoIPCountry: "US",
			wantStatus:       http.StatusOK,
			wantAllowed:      true,
			wantCountry:      "US",
		},
		{
			name:             "basic not allowed",
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
			name:             "empty country list",
			ip:               "142.251.152.119",
			allowedCountries: []string{},
			wantStatus:       http.StatusBadRequest,
		},
		{
			name:             "nil country list",
			ip:               "142.251.152.119",
			allowedCountries: nil,
			wantStatus:       http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       []byte("{invalid json]"),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lookup := testutil.MockCountryLookup{
				Country: tc.mockGeoIPCountry,
				Err:     tc.mockGeoIPError,
			}
			svc := service.NewAccessService(lookup)
			h := NewHandler(svc)

			body := tc.body
			var err error
			if body == nil {
				body, err = json.Marshal(checkRequest{
					IP:               tc.ip,
					AllowedCountries: tc.allowedCountries,
				})
				if err != nil {
					t.Fatal(err)
				}
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
