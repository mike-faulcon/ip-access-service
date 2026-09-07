package service

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/oschwald/geoip2-golang"
)

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
	tests := []struct {
		name             string   // Clear scenario name
		allowedCountries []string // Input list of allowed countries
		mockGeoIPCountry string
		mockGeoIPError   error
		wantAllowed      bool     // Expected output - allowed or not
        wantCountry      string   // Expected output - matched country
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
            lookup := mockCountryLookup{
                country: tc.mockGeoIPCountry,
                err:     tc.mockGeoIPError,
            }
            svc := NewAccessService(lookup)

			ip := net.IP{142, 251, 152, 119}

			got, err := svc.Check(context.Background(), ip, tc.allowedCountries)
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