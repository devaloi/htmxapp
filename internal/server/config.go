package server

import (
	"os"
	"strconv"
)

// Config holds server configuration.
type Config struct {
	Host string
	Port int
	Seed bool
}

// DefaultConfig returns the default server configuration.
func DefaultConfig() Config {
	return Config{
		Host: "",
		Port: 8080,
		Seed: true,
	}
}

// FromEnv reads configuration from environment variables.
func FromEnv() Config {
	cfg := DefaultConfig()

	if host := os.Getenv("HTMXAPP_HOST"); host != "" {
		cfg.Host = host
	}
	if port := os.Getenv("HTMXAPP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}
	if seed := os.Getenv("HTMXAPP_SEED"); seed == "false" || seed == "0" {
		cfg.Seed = false
	}

	return cfg
}

// Addr returns the listen address string.
func (c Config) Addr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}
