package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServiceName   string
	BindAddress   string
	AuditCapacity int
	DefaultMode   string
	PolicyFile    string
}

func Load() (Config, error) {
	cfg := Config{
		ServiceName:   envOr("GATEKEEPER_SERVICE_NAME", "gatekeeper"),
		BindAddress:   envOr("GATEKEEPER_BIND_ADDRESS", ":8080"),
		AuditCapacity: intOr("GATEKEEPER_AUDIT_CAPACITY", 256),
		DefaultMode:   envOr("GATEKEEPER_DEFAULT_MODE", "deny"),
		PolicyFile:    envOr("GATEKEEPER_POLICY_FILE", "sample_data/policies.json"),
	}

	mode := strings.ToLower(strings.TrimSpace(cfg.DefaultMode))
	if mode != "allow" && mode != "deny" {
		return Config{}, errors.New("default mode must be allow or deny")
	}
	cfg.DefaultMode = mode
	return cfg, nil
}

func envOr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}

func intOr(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.Atoi(strings.TrimSpace(value)); err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}
