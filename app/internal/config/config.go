// Package config loads environment-driven configuration for SUKAMOTO.
package config

import (
	"fmt"
	"os"
)

// Config holds all runtime settings sourced from environment variables.
type Config struct {
	Address      string // HTTP listen address, e.g. ":8080"
	DatabaseURL  string // postgres://user:pass@host/db
	Env          string // "development" | "production"
	SessionName  string
	SessionTTLHours int
}

// Load reads configuration from the environment with safe defaults for dev.
func Load() (*Config, error) {
	c := &Config{
		Address:         getenv("SUKAMOTO_ADDR", ":8080"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		Env:             getenv("APP_ENV", "development"),
		SessionName:     "sukamoto_session",
		SessionTTLHours: 72,
	}
	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("config: DATABASE_URL is required")
	}
	return c, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
