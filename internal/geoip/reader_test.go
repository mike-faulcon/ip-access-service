package geoip

import (
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