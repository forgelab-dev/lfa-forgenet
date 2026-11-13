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
	t.Setenv("FORGENET_AUTH_JWT_SECRET", "super-secret-production-key-32bytes")

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

func TestLoad_DockerPublishHost(t *testing.T) {
	t.Chdir(t.TempDir())

	cfg, err := config.Load("")
	require.NoError(t, err)
	assert.Equal(t, "127.0.0.1", cfg.Docker.PublishHost, "par défaut, les ports des conteneurs restent locaux")

	t.Setenv("FORGENET_DOCKER_PUBLISH_HOST", "partout")
	_, err = config.Load("")
	require.ErrorContains(t, err, "docker.publish_host")
}

func TestLoad_HealthExposeDetails(t *testing.T) {
	t.Chdir(t.TempDir())

	cfg, err := config.Load("")
	require.NoError(t, err)
	assert.Equal(t, config.ExposeDetailsLocal, cfg.Health.ExposeDetails)

	t.Setenv("FORGENET_HEALTH_EXPOSE_DETAILS", "parfois")
	_, err = config.Load("")
	require.ErrorContains(t, err, "health.expose_details")
}

func TestLoad_MetricsStoragePath(t *testing.T) {
	t.Chdir(t.TempDir())

	cfg, err := config.Load("")
	require.NoError(t, err)
	assert.Equal(t, "/", cfg.Metrics.StoragePath)

	t.Setenv("FORGENET_METRICS_STORAGE_PATH", "/var/lib/forgenet")
	cfg, err = config.Load("")
	require.NoError(t, err)
	assert.Equal(t, "/var/lib/forgenet", cfg.Metrics.StoragePath)

	t.Setenv("FORGENET_METRICS_STORAGE_PATH", "relatif")
	_, err = config.Load("")
	require.ErrorContains(t, err, "metrics.storage_path")
}

func TestLoad_ProxyDefaults(t *testing.T) {
	t.Chdir(t.TempDir())

	cfg, err := config.Load("")
	require.NoError(t, err)
	assert.False(t, cfg.Proxy.Enabled, "sans Caddy installé, les domaines restent désactivés")
	assert.Equal(t, config.DefaultProxyAdminAddress, cfg.Proxy.AdminAddress)
	assert.Equal(t, config.DefaultProxyNetwork, cfg.Proxy.Network)
	assert.Empty(t, cfg.Proxy.CaddyfilePath)
	assert.False(t, cfg.Proxy.LocalCerts)
}

func TestLoad_ProxyFromEnvironment(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("FORGENET_PROXY_ENABLED", "true")
	t.Setenv("FORGENET_PROXY_ADMIN_ADDRESS", "127.0.0.1:2019")
	t.Setenv("FORGENET_PROXY_EMAIL", "ops@example.org")
	t.Setenv("FORGENET_PROXY_CADDY_CONTAINER", "forgenet-caddy")
	t.Setenv("FORGENET_PROXY_HTTPS_PORT", "8443")

	cfg, err := config.Load("")
	require.NoError(t, err)
	assert.True(t, cfg.Proxy.Enabled)
	assert.Equal(t, "127.0.0.1:2019", cfg.Proxy.AdminAddress)
	assert.Equal(t, "ops@example.org", cfg.Proxy.Email)
	assert.Equal(t, "forgenet-caddy", cfg.Proxy.CaddyContainer)
	assert.Equal(t, 8443, cfg.Proxy.HTTPSPort)
}

func TestLoad_ProxyValidation(t *testing.T) {
	cases := map[string]struct {
		env  map[string]string
		want string
	}{
		"réseau":               {map[string]string{"FORGENET_PROXY_NETWORK": "mon réseau"}, "proxy.network"},
		"conteneur":            {map[string]string{"FORGENET_PROXY_CADDY_CONTAINER": "caddy;rm"}, "proxy.caddy_container"},
		"port":                 {map[string]string{"FORGENET_PROXY_HTTP_PORT": "70000"}, "proxy.http_port"},
		"dashboard sans cible": {map[string]string{"FORGENET_PROXY_DASHBOARD_HOST": "forgenet.example.org"}, "proxy.dashboard_upstream"},
		"cible sans dashboard": {map[string]string{"FORGENET_PROXY_DASHBOARD_UPSTREAM": "forgenet-server:8080"}, "proxy.dashboard_host"},
		"local en production": {map[string]string{
			"FORGENET_PROXY_LOCAL_CERTS": "true", "FORGENET_ENVIRONMENT": "production",
			"FORGENET_AUTH_JWT_SECRET": "super-secret-production-key-32bytes",
		}, "proxy.local_certs"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			t.Setenv("FORGENET_PROXY_ENABLED", "true")
			for k, v := range c.env {
				t.Setenv(k, v)
			}
			_, err := config.Load("")
			require.ErrorContains(t, err, c.want)
		})
	}

	t.Run("dashboard sans proxy", func(t *testing.T) {
		t.Chdir(t.TempDir())
		t.Setenv("FORGENET_PROXY_DASHBOARD_HOST", "forgenet.example.org")
		_, err := config.Load("")
		require.ErrorContains(t, err, "proxy.enabled")
	})

	// Désactivé, le proxy n'est pas vérifié : une valeur inutilisée ne bloque pas le démarrage.
	t.Run("désactivé", func(t *testing.T) {
		t.Chdir(t.TempDir())
		t.Setenv("FORGENET_PROXY_NETWORK", "mon réseau")
		_, err := config.Load("")
		require.NoError(t, err)
	})
}

func TestLoad_WebRoot(t *testing.T) {
	t.Chdir(t.TempDir())

	cfg, err := config.Load("")
	require.NoError(t, err)
	assert.Empty(t, cfg.Web.Root, "sans dossier, le serveur ne sert que l'API")

	t.Setenv("FORGENET_WEB_ROOT", "/usr/share/forgenet/web")
	cfg, err = config.Load("")
	require.NoError(t, err)
	assert.Equal(t, "/usr/share/forgenet/web", cfg.Web.Root)

	t.Setenv("FORGENET_WEB_ROOT", "web/dist")
	_, err = config.Load("")
	require.ErrorContains(t, err, "web.root")
}
