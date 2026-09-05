package geoip

import (
	"testing"
)

func TestInit(t *testing.T) {
	// err := Init("geoip.db")
	err := NewGeoIPReader("../../data/GeoLite2-Country.mmdb")
	if err != nil {
		t.Fatal(err)
	}
}