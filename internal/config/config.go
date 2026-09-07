package config

import (
	"log/slog"
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

func Load() Config {
	return Config{
		HTTPPort:  getEnvInt("HTTP_PORT", DefaultHTTPPort),
		GRPCPort:  getEnvInt("GRPC_PORT", DefaultGRPCPort),
		GeoIPPath: getEnvString("GEOIP_DB_PATH", DefaultGeoIPPath),
	}
}

func getEnvString(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		slog.Warn("Environment variable (string) not found", "key", key, "fallback", fallback, "value", value)
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		slog.Warn("Environment variable (int) not found", "key", key, "fallback", fallback, "value", value)
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
