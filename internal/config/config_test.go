package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		wantHTTP int
		wantGRPC int
		wantPath string
	}{
		{
			name:     "defaults when unset",
			wantHTTP: DefaultHTTPPort,
			wantGRPC: DefaultGRPCPort,
			wantPath: DefaultGeoIPPath,
		},
		{
			name: "env overrides",
			env: map[string]string{
				"HTTP_PORT":     "3001",
				"GRPC_PORT":     "4001",
				"GEOIP_DB_PATH": "/custom/db.mmdb",
			},
			wantHTTP: 3001,
			wantGRPC: 4001,
			wantPath: "/custom/db.mmdb",
		},
		{
			name:     "invalid int falls back",
			env:      map[string]string{"HTTP_PORT": "abc"},
			wantHTTP: DefaultHTTPPort,
			wantGRPC: DefaultGRPCPort,
			wantPath: DefaultGeoIPPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HTTP_PORT", "")
			t.Setenv("GRPC_PORT", "")
			t.Setenv("GEOIP_DB_PATH", "")
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg := Load()
			if cfg.HTTPPort != tt.wantHTTP {
				t.Fatalf("HTTPPort = %d, want %d", cfg.HTTPPort, tt.wantHTTP)
			}
			if cfg.GRPCPort != tt.wantGRPC {
				t.Fatalf("GRPCPort = %d, want %d", cfg.GRPCPort, tt.wantGRPC)
			}
			if cfg.GeoIPPath != tt.wantPath {
				t.Fatalf("GeoIPPath = %q, want %q", cfg.GeoIPPath, tt.wantPath)
			}
		})
	}
}
