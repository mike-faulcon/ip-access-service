package geoip

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type GeoIPReader struct {
	reader *geoip2.Reader
}

func NewGeoIPReader(path string) (*GeoIPReader, error) {
	slog.Info("Initializing GeoIP reader", "path", path)
	reader, err := geoip2.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open GeoIP database: %w", err)
	}

	return &GeoIPReader{
		reader: reader,
	}, nil
}

func (s *GeoIPReader) Lookup(ip net.IP) (*geoip2.Country, error) {
	return s.reader.Country(ip)
}

func (s *GeoIPReader) Close() error {
	return s.reader.Close()
}
