package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	DefaultHTTPPort  = 8080
	DefaultGRPCPort  = 9090
	DefaultGeoIPPath = "data/GeoLite2-Country.mmdb"
)

type Config struct {
	HTTPPort  int
	GRPCPort  int
	GeoIPPath string
}

func Load() (Config, error) {
	httpPort, err := getEnvInt("HTTP_PORT", DefaultHTTPPort)
	if err != nil {
		return Config{}, err
	}

	grpcPort, err := getEnvInt("GRPC_PORT", DefaultGRPCPort)
	if err != nil {
		return Config{}, err
	}

	geoIPPath := getEnvString("GEOIP_DB_PATH", DefaultGeoIPPath)

	return Config{
		HTTPPort:  httpPort,
		GRPCPort:  grpcPort,
		GeoIPPath: geoIPPath,
	}, nil
}

func getEnvString(key, defaultVal string) string {
	raw, ok := os.LookupEnv(key)
	if !ok || raw == "" {
		return defaultVal
	}
	return raw
}

func getEnvInt(key string, defaultVal int) (int, error) {
	raw, ok := os.LookupEnv(key)
	if !ok || raw == "" {
		return defaultVal, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s: invalid integer %q", key, raw)
	}
	return n, nil
}
