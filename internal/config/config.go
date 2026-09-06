package config

import (
    "os"
    "strconv"
)

const (
	DefaultPort = 8080
	DefaultGeoIPPath = "data/GeoLite2-Country.mmdb"
)

type Config struct {
    Port      int
    GeoIPPath string
}

func Load() Config {
    return Config{
        Port:      getEnvInt("PORT", DefaultPort),
        GeoIPPath: getEnvString("GEOIP_DB_PATH", DefaultGeoIPPath),
    }
}

func getEnvString(key, fallback string) string {
    value, exists := os.LookupEnv(key)
    if !exists || value == "" {
        return fallback
    }

    return value
}

func getEnvInt(key string, fallback int) int {
    value, exists := os.LookupEnv(key)
    if !exists || value == "" {
        return fallback
    }

    parsed, err := strconv.Atoi(value)
    if err != nil {
        return fallback
    }

    return parsed
}