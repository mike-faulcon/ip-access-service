package testutil

import (
	"net"

	"github.com/oschwald/geoip2-golang"
)

// MockCountryLookup is a test double for GeoIP country lookups.
type MockCountryLookup struct {
	Country string
	Err     error
}

func (m MockCountryLookup) Lookup(_ net.IP) (*geoip2.Country, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return &geoip2.Country{
		Country: struct {
			Names             map[string]string `maxminddb:"names"`
			IsoCode           string            `maxminddb:"iso_code"`
			GeoNameID         uint              `maxminddb:"geoname_id"`
			IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
		}{IsoCode: m.Country},
	}, nil
}
