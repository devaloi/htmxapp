package server

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
	if cfg.Addr() != ":8080" {
		t.Errorf("expected :8080, got %s", cfg.Addr())
	}
}

func TestConfig_Addr(t *testing.T) {
	cfg := Config{Host: "localhost", Port: 3000}
	if got := cfg.Addr(); got != "localhost:3000" {
		t.Errorf("expected localhost:3000, got %s", got)
	}
}

func TestFromEnv(t *testing.T) {
	t.Setenv("HTMXAPP_HOST", "0.0.0.0")
	t.Setenv("HTMXAPP_PORT", "9090")
	t.Setenv("HTMXAPP_SEED", "false")

	cfg := FromEnv()
	if cfg.Host != "0.0.0.0" {
		t.Errorf("expected 0.0.0.0, got %s", cfg.Host)
	}
	if cfg.Port != 9090 {
		t.Errorf("expected 9090, got %d", cfg.Port)
	}
	if cfg.Seed {
		t.Error("expected seed=false")
	}
}
