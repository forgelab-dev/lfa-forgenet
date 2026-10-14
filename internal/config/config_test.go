package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := config.Load("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, config.DefaultEnvironment, cfg.Environment)
	assert.Equal(t, config.DefaultHost, cfg.Server.Host)
	assert.Equal(t, config.DefaultHTTPPort, cfg.Server.HTTPPort)
	assert.Equal(t, config.DefaultGRPCPort, cfg.Server.GRPCPort)
	assert.Equal(t, config.DefaultReadTimeout, cfg.Server.ReadTimeout)
	assert.Equal(t, config.DefaultWriteTimeout, cfg.Server.WriteTimeout)
	assert.Equal(t, config.DefaultShutdownTimeout, cfg.Server.ShutdownTimeout)
	assert.Equal(t, config.DefaultLogLevel, cfg.Log.Level)
	assert.Equal(t, config.DefaultLogFormat, cfg.Log.Format)

	assert.Equal(t, "127.0.0.1:8080", cfg.Server.HTTPAddress())
	assert.Equal(t, "127.0.0.1:9090", cfg.Server.GRPCAddress())
}

func TestLoad_EnvOverrides(t *testing.T) {
	t.Setenv("FORGENET_SERVER_HOST", "192.168.1.10")
	t.Setenv("FORGENET_SERVER_HTTP_PORT", "8888")
	t.Setenv("FORGENET_SERVER_GRPC_PORT", "9999")
	t.Setenv("FORGENET_LOG_LEVEL", "debug")
	t.Setenv("FORGENET_LOG_FORMAT", "json")
	t.Setenv("FORGENET_ENVIRONMENT", "production")

	cfg, err := config.Load("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "192.168.1.10", cfg.Server.Host)
	assert.Equal(t, 8888, cfg.Server.HTTPPort)
	assert.Equal(t, 9999, cfg.Server.GRPCPort)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, "production", cfg.Environment)
}

func TestLoad_ConfigFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "custom-config.yaml")

	yamlContent := `
environment: staging
server:
  host: "10.0.0.1"
  http_port: 7000
  grpc_port: 7001
  shutdown_timeout: 5s
log:
  level: "warn"
  format: "json"
`
	err := os.WriteFile(configPath, []byte(yamlContent), 0o600)
	require.NoError(t, err)

	cfg, err := config.Load(configPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "staging", cfg.Environment)
	assert.Equal(t, "10.0.0.1", cfg.Server.Host)
	assert.Equal(t, 7000, cfg.Server.HTTPPort)
	assert.Equal(t, 7001, cfg.Server.GRPCPort)
	assert.Equal(t, 5*time.Second, cfg.Server.ShutdownTimeout)
	assert.Equal(t, "warn", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
}

func TestLoad_MissingConfigFileExplicit(t *testing.T) {
	// Si un chemin est explicitement spécifié mais n'existe pas, une erreur doit être renvoyée
	_, err := config.Load("/chemin/vers/fichier/inexistant.yaml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "lecture du fichier de configuration")
}

func TestLoad_InvalidYAMLFile(t *testing.T) {
	// Fichier YAML syntaxiquement invalide passé explicitement
	tempDir := t.TempDir()
	corruptConfigPath := filepath.Join(tempDir, "corrupt.yaml")
	corruptYAML := `
server:
  http_port: 8080
    indentation_invalide: true
  invalid_yaml: [ unclosed bracket
`
	err := os.WriteFile(corruptConfigPath, []byte(corruptYAML), 0o600)
	require.NoError(t, err)

	_, err = config.Load(corruptConfigPath)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "lecture du fichier de configuration")
}

func TestLoad_ConfigFileNotFoundIgnored_DefaultFallback(t *testing.T) {
	// Vérifie que l'absence de fichier par défaut dans un répertoire propre ne bloque pas et applique les defaults
	cfg, err := config.Load("")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, config.DefaultHost, cfg.Server.Host)
	assert.Equal(t, config.DefaultHTTPPort, cfg.Server.HTTPPort)
}

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name        string
		modify      func(c *config.Config)
		expectedErr string
	}{
		{
			name: "Port HTTP invalide (0)",
			modify: func(c *config.Config) {
				c.Server.HTTPPort = 0
			},
			expectedErr: "port HTTP invalide",
		},
		{
			name: "Port gRPC invalide (> 65535)",
			modify: func(c *config.Config) {
				c.Server.GRPCPort = 70000
			},
			expectedErr: "port gRPC invalide",
		},
		{
			name: "Ports identiques",
			modify: func(c *config.Config) {
				c.Server.HTTPPort = 8080
				c.Server.GRPCPort = 8080
			},
			expectedErr: "ne peuvent pas être identiques",
		},
		{
			name: "Shutdown timeout négatif",
			modify: func(c *config.Config) {
				c.Server.ShutdownTimeout = -1 * time.Second
			},
			expectedErr: "shutdown_timeout doit être strictement positif",
		},
		{
			name: "Log level invalide",
			modify: func(c *config.Config) {
				c.Log.Level = "invalid"
			},
			expectedErr: "niveau de log inconnu",
		},
		{
			name: "Log format invalide",
			modify: func(c *config.Config) {
				c.Log.Format = "xml"
			},
			expectedErr: "format de log inconnu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Load("")
			require.NoError(t, err)

			tt.modify(cfg)
			err = cfg.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}
