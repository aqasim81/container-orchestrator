package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setEnv(t *testing.T, envs map[string]string) {
	t.Helper()
	for k, v := range envs {
		t.Setenv(k, v)
	}
}

func TestLoad_Defaults(t *testing.T) {
	setEnv(t, map[string]string{
		"API_KEY": "test-key",
	})

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "./data", cfg.DataDir)
	assert.Equal(t, "unix:///var/run/docker.sock", cfg.DockerHost)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "test-key", cfg.APIKey)
	assert.Equal(t, "http://localhost:3000", cfg.DashboardURL)
	assert.Equal(t, 10*time.Second, cfg.NodeHeartbeatInterval)
	assert.Equal(t, 30*time.Second, cfg.NodeHeartbeatTimeout)
	assert.Equal(t, 10*time.Second, cfg.HealthCheckInterval)
	assert.Equal(t, 10*time.Second, cfg.ReconcileInterval)
}

func TestLoad_CustomValues(t *testing.T) {
	setEnv(t, map[string]string{
		"API_KEY":                 "custom-key",
		"ORCHESTRATOR_PORT":       "9090",
		"ORCHESTRATOR_DATA_DIR":   "/tmp/data",
		"LOG_LEVEL":               "debug",
		"DASHBOARD_URL":           "http://example.com",
		"NODE_HEARTBEAT_INTERVAL": "5s",
		"NODE_HEARTBEAT_TIMEOUT":  "15s",
		"HEALTH_CHECK_INTERVAL":   "30s",
		"RECONCILE_INTERVAL":      "20s",
	})

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, 9090, cfg.Port)
	assert.Equal(t, "/tmp/data", cfg.DataDir)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "custom-key", cfg.APIKey)
	assert.Equal(t, "http://example.com", cfg.DashboardURL)
	assert.Equal(t, 5*time.Second, cfg.NodeHeartbeatInterval)
	assert.Equal(t, 15*time.Second, cfg.NodeHeartbeatTimeout)
	assert.Equal(t, 30*time.Second, cfg.HealthCheckInterval)
	assert.Equal(t, 20*time.Second, cfg.ReconcileInterval)
}

func TestLoad_MissingAPIKey(t *testing.T) {
	cfg, err := Load()

	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API_KEY")
}

func TestLoad_InvalidPort(t *testing.T) {
	setEnv(t, map[string]string{
		"API_KEY":           "test-key",
		"ORCHESTRATOR_PORT": "0",
	})

	cfg, err := Load()
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ORCHESTRATOR_PORT")
}

func TestLoad_InvalidLogLevel(t *testing.T) {
	setEnv(t, map[string]string{
		"API_KEY":   "test-key",
		"LOG_LEVEL": "verbose",
	})

	cfg, err := Load()
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "LOG_LEVEL")
}

func TestLoad_HeartbeatTimeoutMustExceedInterval(t *testing.T) {
	setEnv(t, map[string]string{
		"API_KEY":                 "test-key",
		"NODE_HEARTBEAT_INTERVAL": "30s",
		"NODE_HEARTBEAT_TIMEOUT":  "10s",
	})

	cfg, err := Load()
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NODE_HEARTBEAT_TIMEOUT")
}
