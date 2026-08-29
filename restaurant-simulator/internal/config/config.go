// Package config loads restaurant-simulator settings from the environment.
package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds every environment-driven setting for the simulator.
type Config struct {
	MainServiceURL   string
	RestaurantID     string
	ListenPort       string
	MenuSyncInterval time.Duration
	PollInterval     time.Duration
	LogLevel         string
}

// Load reads configuration from the environment, defaulting to values that
// match the fixed restaurant seeded by migrations/000002_seed_restaurant.
func Load() (Config, error) {
	menuSync, err := parseDuration("MENU_SYNC_INTERVAL", "60s")
	if err != nil {
		return Config{}, err
	}
	poll, err := parseDuration("POLL_INTERVAL", "10s")
	if err != nil {
		return Config{}, err
	}

	return Config{
		MainServiceURL:   getenv("MAIN_SERVICE_URL", "http://localhost:8080"),
		RestaurantID:     getenv("RESTAURANT_ID", "123e4567-e89b-12d3-a456-426614174000"),
		ListenPort:       getenv("LISTEN_PORT", "8081"),
		MenuSyncInterval: menuSync,
		PollInterval:     poll,
		LogLevel:         getenv("LOG_LEVEL", "info"),
	}, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseDuration(key, def string) (time.Duration, error) {
	raw := getenv(key, def)
	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("config: invalid duration for %s=%q: %w", key, raw, err)
	}
	return d, nil
}
