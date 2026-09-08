package geoip

import (
	"net"
	"testing"
)

func TestInit(t *testing.T) {
	reader, err := NewGeoIPReader("../../data/GeoLite2-Country.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	if reader == nil {
		t.Fatal("Reader is nil")
	}
}

func TestNewGeoIPReader_invalidPath(t *testing.T) {
	_, err := NewGeoIPReader("/nonexistent/GeoLite2-Country.mmdb")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestGeoIPReader_Lookup(t *testing.T) {
	reader, err := NewGeoIPReader("../../data/GeoLite2-Country.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	country, err := reader.Lookup(net.ParseIP("8.8.8.8"))
	if err != nil {
		t.Fatal(err)
	}
	if country.Country.IsoCode == "" {
		t.Fatal("expected country ISO code")
	}
}
