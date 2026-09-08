package service

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/mike-faulcon/ip-access-service/internal/testutil"
)

func TestAccessService_Check(t *testing.T) {
	tests := []struct {
		name             string   // Clear scenario name
		allowedCountries []string // Input list of allowed countries
		mockGeoIPCountry string
		mockGeoIPError   error
		wantAllowed      bool   // Expected output - allowed or not
		wantCountry      string // Expected output - matched country
	}{
		{
			name:             "basic allowed",
			allowedCountries: []string{"CA"},
			mockGeoIPCountry: "CA",
			wantAllowed:      true,
			wantCountry:      "CA",
		},
		{
			name:             "not allowed",
			allowedCountries: []string{"UK", "FR"},
			mockGeoIPCountry: "US",
			wantAllowed:      false,
			wantCountry:      "US",
		},
		{
			name:             "lookup error",
			allowedCountries: []string{"US"},
			mockGeoIPError:   errors.New("db error"),
		},
		{
			name:             "empty country list",
			allowedCountries: []string{},
		},
		{
			name:             "case-insensitive",
			allowedCountries: []string{"us"},
			mockGeoIPCountry: "US",
			wantAllowed:      true,
			wantCountry:      "US",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lookup := testutil.MockCountryLookup{
				Country: tc.mockGeoIPCountry,
				Err:     tc.mockGeoIPError,
			}
			svc := NewAccessService(lookup)

			ip := net.IP{142, 251, 152, 119}

			accessRequest := AccessRequest{
				IP:               ip,
				AllowedCountries: tc.allowedCountries,
			}

			got, err := svc.Check(context.Background(), accessRequest)
			if err != nil && (tc.mockGeoIPError == nil || err.Error() != tc.mockGeoIPError.Error()) {
				t.Error("Unexpected Error", err)
			}

			if got.Allowed != tc.wantAllowed {
				t.Errorf("allowed = %v, want %v", got.Allowed, tc.wantAllowed)
			}
			if got.CountryISO != tc.wantCountry {
				t.Errorf("country = %q, want %q", got.CountryISO, tc.wantCountry)
			}
		})
	}
}
