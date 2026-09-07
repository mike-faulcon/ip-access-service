package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	cfg := Load()

	if cfg.HTTPPort <= 0 {
		t.Fatal("Config Port is invalid")
	}

	if cfg.GeoIPPath == "" {
		t.Fatal("Config GeoIPPath is invalid")
	}
}
