package grpcapi

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accessv1 "github.com/mike-faulcon/ip-access-service/gen/access/v1"
	"github.com/mike-faulcon/ip-access-service/internal/service"
	"github.com/mike-faulcon/ip-access-service/internal/testutil"
)

func TestCheckAccess(t *testing.T) {
	tests := []struct {
		name             string
		ip               string
		allowedCountries []string
		mockGeoIPCountry string
		mockGeoIPError   error
		wantCode         codes.Code
		wantAllowed      bool
		wantCountry      string
	}{
		{
			name:             "allowed",
			ip:               "142.251.152.119",
			allowedCountries: []string{"US"},
			mockGeoIPCountry: "US",
			wantCode:         codes.OK,
			wantAllowed:      true,
			wantCountry:      "US",
		},
		{
			name:             "not allowed",
			ip:               "142.251.152.119",
			allowedCountries: []string{"UK", "FR"},
			mockGeoIPCountry: "US",
			wantCode:         codes.OK,
			wantAllowed:      false,
			wantCountry:      "US",
		},
		{
			name:             "invalid ip",
			ip:               "1.2.34",
			allowedCountries: []string{"US"},
			wantCode:         codes.InvalidArgument,
		},
		{
			name:             "empty allowed countries",
			ip:               "142.251.152.119",
			allowedCountries: []string{},
			wantCode:         codes.InvalidArgument,
		},
		{
			name:             "nil allowed countries",
			ip:               "142.251.152.119",
			allowedCountries: nil,
			wantCode:         codes.InvalidArgument,
		},
		{
			name:             "lookup error",
			ip:               "142.251.152.119",
			allowedCountries: []string{"US"},
			mockGeoIPError:   errors.New("db unavailable"),
			wantCode:         codes.Internal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := service.NewAccessService(testutil.MockCountryLookup{
				Country: tc.mockGeoIPCountry,
				Err:     tc.mockGeoIPError,
			})
			server := NewServer(svc)

			resp, err := server.CheckAccess(context.Background(), &accessv1.CheckAccessRequest{
				Ip:               tc.ip,
				AllowedCountries: tc.allowedCountries,
			})

			if tc.wantCode != codes.OK {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if status.Code(err) != tc.wantCode {
					t.Fatalf("code = %v, want %v, err = %v", status.Code(err), tc.wantCode, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.Allowed != tc.wantAllowed {
				t.Fatalf("allowed = %v, want %v", resp.Allowed, tc.wantAllowed)
			}
			if resp.Country != tc.wantCountry {
				t.Fatalf("country = %q, want %q", resp.Country, tc.wantCountry)
			}
		})
	}
}
