// Package config provides environment-based configuration with validation.
package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config holds all application configuration parsed from environment variables.
type Config struct {
	// DataDir is the directory for bbolt database files.
	DataDir string `env:"ORCHESTRATOR_DATA_DIR" envDefault:"./data"`

	// DockerHost is the Docker daemon socket address.
	DockerHost string `env:"DOCKER_HOST" envDefault:"unix:///var/run/docker.sock"`

	// LogLevel controls the logging verbosity (debug, info, warn, error).
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// APIKey is the required API key for authentication.
	APIKey string `env:"API_KEY,required"` //nolint:gosec // Not a hardcoded credential, populated from env.

	// DashboardURL is the dashboard origin for CORS configuration.
	DashboardURL string `env:"DASHBOARD_URL" envDefault:"http://localhost:3000"`

	// NodeHeartbeatInterval is how often nodes send heartbeats.
	NodeHeartbeatInterval time.Duration `env:"NODE_HEARTBEAT_INTERVAL" envDefault:"10s"`

	// NodeHeartbeatTimeout is the time before a node is marked NotReady.
	NodeHeartbeatTimeout time.Duration `env:"NODE_HEARTBEAT_TIMEOUT" envDefault:"30s"`

	// HealthCheckInterval is the default health check probe interval.
	HealthCheckInterval time.Duration `env:"HEALTH_CHECK_INTERVAL" envDefault:"10s"`

	// ReconcileInterval is the deployment reconciliation loop interval.
	ReconcileInterval time.Duration `env:"RECONCILE_INTERVAL" envDefault:"10s"`

	// Port is the API server listen port.
	Port int `env:"ORCHESTRATOR_PORT" envDefault:"8080"`
}

// Load parses and validates configuration from environment variables.
// It returns an error if required variables are missing or values are invalid.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	if cfg.Port < 1 || cfg.Port > 65535 {
		return fmt.Errorf("ORCHESTRATOR_PORT must be between 1 and 65535, got %d", cfg.Port)
	}

	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[cfg.LogLevel] {
		return fmt.Errorf("LOG_LEVEL must be one of debug, info, warn, error; got %q", cfg.LogLevel)
	}

	if cfg.NodeHeartbeatTimeout <= cfg.NodeHeartbeatInterval {
		return fmt.Errorf("NODE_HEARTBEAT_TIMEOUT (%s) must be greater than NODE_HEARTBEAT_INTERVAL (%s)",
			cfg.NodeHeartbeatTimeout, cfg.NodeHeartbeatInterval)
	}

	return nil
}
