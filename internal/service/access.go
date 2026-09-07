package service

import (
	"context"
	"net"
	"slices"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

type AccessService struct {
	geoIP countryLookup
}

type AccessRequest struct {
    IP               net.IP
    AllowedCountries []string
}

type AccessResult struct {
    Allowed    bool
    CountryISO string
}

type countryLookup interface {
	Lookup(ip net.IP) (*geoip2.Country, error)
}

func NewAccessService(geoIP countryLookup) *AccessService {
	return &AccessService{
        geoIP: geoIP,
    }
}

func (s *AccessService) Check(
    ctx context.Context,
	req AccessRequest,
) (AccessResult, error) {
    geoIPRecord, err := s.geoIP.Lookup(req.IP)
    if err != nil {
        return AccessResult{}, err
    }

	countryISO := geoIPRecord.Country.IsoCode

	return AccessResult{
        Allowed:    isCountryAllowed(countryISO, req.AllowedCountries),
        CountryISO: countryISO,
    }, nil
}

func isCountryAllowed(countryISO string, allowedCountries []string) bool {
	// make this check case-insenstive
	uppercased := make([]string, len(allowedCountries))
	for i, str := range allowedCountries {
		uppercased[i] = strings.ToUpper(str)
	}

    return slices.Contains(uppercased, countryISO)
}