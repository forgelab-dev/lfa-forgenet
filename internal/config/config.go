// Package config centralise la configuration de ForgeNet via Viper.
package config

import (
	"errors"
	"fmt"
	"net/netip"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	// DefaultEnvPrefix est le préfixe de toutes les variables d'environnement de ForgeNet.
	DefaultEnvPrefix = "FORGENET"

	// DefaultHost est l'adresse d'écoute hôte par défaut des serveurs HTTP et gRPC.
	DefaultHost = "127.0.0.1"
	// DefaultHTTPPort est le port par défaut du serveur HTTP.
	DefaultHTTPPort = 8080
	// DefaultGRPCPort est le port par défaut du serveur gRPC.
	DefaultGRPCPort = 9090
	// DefaultReadTimeout est le délai maximal de lecture de requête HTTP.
	DefaultReadTimeout = 15 * time.Second
	// DefaultWriteTimeout est le délai maximal d'écriture de réponse HTTP.
	DefaultWriteTimeout = 15 * time.Second
	// DefaultIdleTimeout est la durée maximale de maintien de connexion inactive.
	DefaultIdleTimeout = 60 * time.Second
	// DefaultShutdownTimeout est le délai alloué au drainage gracieux des connexions.
	DefaultShutdownTimeout = 10 * time.Second

	// DefaultLogLevel est le niveau de log par défaut.
	DefaultLogLevel = "info"
	// DefaultLogFormat est le format d'affichage par défaut des logs.
	DefaultLogFormat = "console"
	// DefaultEnvironment est l'environnement par défaut de l'application.
	DefaultEnvironment = "development"

	// DefaultDatabasePath est le chemin d'accès au fichier de base SQLite par défaut.
	DefaultDatabasePath = "data/forgenet.db"
	// DefaultDatabaseBusyTimeout est le délai de grâce maximal lors d'un verrou SQLite concurrent.
	DefaultDatabaseBusyTimeout = 5 * time.Second

	// DefaultTokenTTL est la durée de validité par défaut des tokens de session JWT (24 heures).
	DefaultTokenTTL = 24 * time.Hour
	// DefaultAPITokenTTL est la durée de validité par défaut des tokens API (30 jours).
	DefaultAPITokenTTL = 30 * 24 * time.Hour
	// MinJWTSecretLength est la longueur minimale du secret JWT en production.
	MinJWTSecretLength = 32

	// DefaultProxyAdminAddress est le socket de l'API d'administration de Caddy.
	DefaultProxyAdminAddress = "unix//run/forgenet/caddy-admin.sock"
	// DefaultProxyNetwork est le réseau Docker partagé entre Caddy et les applications.
	DefaultProxyNetwork = "forgenet-edge"
)

// ServerConfig contient les paramètres réseau et timeouts des serveurs HTTP et gRPC.
type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	HTTPPort        int           `mapstructure:"http_port"`
	GRPCPort        int           `mapstructure:"grpc_port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// HTTPAddress retourne l'adresse d'écoute HTTP sous forme "host:port".
func (s ServerConfig) HTTPAddress() string {
	return fmt.Sprintf("%s:%d", s.Host, s.HTTPPort)
}

// GRPCAddress retourne l'adresse d'écoute gRPC sous forme "host:port".
func (s ServerConfig) GRPCAddress() string {
	return fmt.Sprintf("%s:%d", s.Host, s.GRPCPort)
}

// LogConfig contient les options du système de journalisation.
type LogConfig struct {
	Level  string `mapstructure:"level"`  // "debug", "info", "warn", "error"
	Format string `mapstructure:"format"` // "console", "json"
}

// DatabaseConfig contient les paramètres de persistance SQLite.
type DatabaseConfig struct {
	Path         string        `mapstructure:"path"`
	BusyTimeout  time.Duration `mapstructure:"busy_timeout"`
	MaxOpenConns int           `mapstructure:"max_open_conns"`
	MaxIdleConns int           `mapstructure:"max_idle_conns"`
}

// AuthConfig contient les options de sécurité et de génération de jetons d'authentification.
type AuthConfig struct {
	JWTSecret   string        `mapstructure:"jwt_secret"`
	TokenTTL    time.Duration `mapstructure:"token_ttl"`
	APITokenTTL time.Duration `mapstructure:"api_token_ttl"`
}

// DockerConfig contient les options de connexion au démon Docker.
type DockerConfig struct {
	Host       string `mapstructure:"host"`        // vide : variables DOCKER_* ou socket par défaut
	APIVersion string `mapstructure:"api_version"` // vide : version négociée avec le démon
	// PublishHost est l'adresse de l'hôte où les ports des conteneurs sont
	// publiés. Par défaut 127.0.0.1 : l'accès public passe par le reverse proxy.
	PublishHost string `mapstructure:"publish_host"`
}

// MetricsConfig règle la collecte des métriques de l'hôte.
type MetricsConfig struct {
	// StoragePath désigne le système de fichiers dont l'occupation est
	// mesurée (« / » par défaut). Le /proc lu est celui de la variable
	// HOST_PROC de gopsutil, sinon /proc : dans un conteneur, y monter le
	// /proc de l'hôte, sans quoi le serveur ne mesure que son conteneur.
	StoragePath string `mapstructure:"storage_path"`
}

// Valeurs de health.expose_details.
const (
	// ExposeDetailsAlways montre toujours la cause d'une indisponibilité.
	ExposeDetailsAlways = "always"
	// ExposeDetailsLocal ne la montre qu'aux requêtes locales directes.
	ExposeDetailsLocal = "local"
	// ExposeDetailsNever ne la montre jamais (statut seul).
	ExposeDetailsNever = "never"
)

// HealthConfig règle ce que les sondes HTTP révèlent.
type HealthConfig struct {
	// ExposeDetails décide si /readyz détaille l'état de chaque dépendance
	// (base, Docker…), dont les messages peuvent révéler des chemins internes.
	// Par défaut « local » : seule une requête venue de la machine elle-même,
	// sans passer par un proxy, les reçoit. Les détails sont toujours écrits
	// dans le journal du serveur.
	ExposeDetails string `mapstructure:"expose_details"`
}

// ProxyConfig règle le reverse proxy Caddy et les domaines (ADR-0007, ADR-0008).
type ProxyConfig struct {
	// Enabled active les domaines. Désactivé par défaut : sans Caddy installé,
	// le service des domaines répond FAILED_PRECONDITION.
	Enabled bool `mapstructure:"enabled"`
	// AdminAddress est l'adresse de l'API d'administration de Caddy : un socket
	// Unix (« unix//chemin », recommandé) ou une adresse de bouclage. Toute
	// autre adresse est refusée au démarrage.
	AdminAddress string `mapstructure:"admin_address"`
	// Email reçoit les avis de l'autorité ACME (facultatif).
	Email string `mapstructure:"email"`
	// CaddyfilePath reçoit une copie de la configuration appliquée ; vide :
	// « Caddyfile » dans le dossier de la base de données.
	CaddyfilePath string `mapstructure:"caddyfile_path"`
	// Network est le réseau Docker partagé entre Caddy et les applications.
	Network string `mapstructure:"network"`
	// CaddyContainer, s'il est fixé, est branché sur ce réseau au démarrage.
	CaddyContainer string `mapstructure:"caddy_container"`
	// LocalCerts fait signer les certificats par l'autorité interne de Caddy :
	// pour les essais sans domaine public. Refusé en production.
	LocalCerts bool `mapstructure:"local_certs"`
	// HTTPPort et HTTPSPort changent les ports d'écoute de Caddy ; 0 garde 80 et 443.
	HTTPPort  int `mapstructure:"http_port"`
	HTTPSPort int `mapstructure:"https_port"`
	// DashboardHost est le nom de domaine du dashboard, servi en HTTPS par
	// Caddy (ADR-0009). Vide : le dashboard n'est pas publié, on l'ouvre par
	// un tunnel SSH vers le port HTTP du serveur.
	DashboardHost string `mapstructure:"dashboard_host"`
	// DashboardUpstream est le serveur ForgeNet tel que Caddy le joint,
	// « conteneur:port » (forgenet-server:8080 dans l'installation standard).
	DashboardUpstream string `mapstructure:"dashboard_upstream"`
}

// WebConfig règle le service du dashboard par le serveur HTTP (ADR-0009).
type WebConfig struct {
	// Root est le dossier du dashboard compilé (web/dist). Vide : le serveur
	// ne sert que l'API, comme en développement derrière Vite.
	Root string `mapstructure:"root"`
}

// Config est la structure racine de la configuration ForgeNet.
type Config struct {
	Environment string         `mapstructure:"environment"` // "development", "staging", "production"
	Server      ServerConfig   `mapstructure:"server"`
	Log         LogConfig      `mapstructure:"log"`
	Database    DatabaseConfig `mapstructure:"database"`
	Auth        AuthConfig     `mapstructure:"auth"`
	Docker      DockerConfig   `mapstructure:"docker"`
	Health      HealthConfig   `mapstructure:"health"`
	Metrics     MetricsConfig  `mapstructure:"metrics"`
	Proxy       ProxyConfig    `mapstructure:"proxy"`
	Web         WebConfig      `mapstructure:"web"`
}

// Load charge la configuration depuis un fichier optionnel, les variables d'environnement et les valeurs par défaut.
func Load(configPath string) (*Config, error) {
	v := viper.New()

	setDefaults(v)

	v.SetEnvPrefix(DefaultEnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("lecture du fichier de configuration %s : %w", configPath, err)
		}
	} else {
		v.SetConfigName("forgenet")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/forgenet")
		if err := v.ReadInConfig(); err != nil {
			var notFoundErr viper.ConfigFileNotFoundError
			if !errors.As(err, &notFoundErr) {
				return nil, fmt.Errorf("erreur de lecture du fichier de configuration : %w", err)
			}
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("désérialisation de la configuration : %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration invalide : %w", err)
	}

	return &cfg, nil
}

// Validate s'assure de la cohérence de l'ensemble des paramètres de configuration.
func (c *Config) Validate() error {
	if c.Server.HTTPPort <= 0 || c.Server.HTTPPort > 65535 {
		return fmt.Errorf("port HTTP invalide : %d (doit être compris entre 1 et 65535)", c.Server.HTTPPort)
	}
	if c.Server.GRPCPort <= 0 || c.Server.GRPCPort > 65535 {
		return fmt.Errorf("port gRPC invalide : %d (doit être compris entre 1 et 65535)", c.Server.GRPCPort)
	}
	if c.Server.HTTPPort == c.Server.GRPCPort {
		return fmt.Errorf("les ports HTTP et gRPC ne peuvent pas être identiques (%d)", c.Server.HTTPPort)
	}
	if c.Server.ShutdownTimeout <= 0 {
		return fmt.Errorf("le shutdown_timeout doit être strictement positif")
	}

	if c.Database.Path == "" {
		return fmt.Errorf("le chemin de la base de données (database.path) ne peut pas être vide")
	}
	if c.Database.BusyTimeout <= 0 {
		return fmt.Errorf("le busy_timeout de la base de données doit être strictement positif")
	}

	if c.Auth.TokenTTL <= 0 {
		return fmt.Errorf("le token_ttl d'authentification doit être strictement positif")
	}
	if c.Auth.APITokenTTL <= 0 {
		return fmt.Errorf("le api_token_ttl d'authentification doit être strictement positif")
	}
	// HS256 demande une clé d'au moins 256 bits (RFC 7518, §3.2).
	if c.Environment == "production" && len(c.Auth.JWTSecret) < MinJWTSecretLength {
		return fmt.Errorf("en production, auth.jwt_secret doit comporter au moins %d caractères", MinJWTSecretLength)
	}

	if _, err := netip.ParseAddr(c.Docker.PublishHost); err != nil {
		return fmt.Errorf("docker.publish_host doit être une adresse IP : %q", c.Docker.PublishHost)
	}

	if !filepath.IsAbs(c.Metrics.StoragePath) {
		return fmt.Errorf("metrics.storage_path doit être un chemin absolu : %q", c.Metrics.StoragePath)
	}

	switch c.Health.ExposeDetails {
	case ExposeDetailsAlways, ExposeDetailsLocal, ExposeDetailsNever:
	default:
		return fmt.Errorf("health.expose_details inconnu : %q (valeurs acceptées : always, local, never)", c.Health.ExposeDetails)
	}

	if err := c.Proxy.validate(c.Environment); err != nil {
		return err
	}
	if c.Web.Root != "" && !filepath.IsAbs(c.Web.Root) {
		return fmt.Errorf("web.root doit être un chemin absolu : %q", c.Web.Root)
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[strings.ToLower(c.Log.Level)] {
		return fmt.Errorf("niveau de log inconnu : %q (valeurs acceptées: debug, info, warn, error)", c.Log.Level)
	}

	validLogFormats := map[string]bool{"console": true, "json": true}
	if !validLogFormats[strings.ToLower(c.Log.Format)] {
		return fmt.Errorf("format de log inconnu : %q (valeurs acceptées: console, json)", c.Log.Format)
	}

	return nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("environment", DefaultEnvironment)

	v.SetDefault("server.host", DefaultHost)
	v.SetDefault("server.http_port", DefaultHTTPPort)
	v.SetDefault("server.grpc_port", DefaultGRPCPort)
	v.SetDefault("server.read_timeout", DefaultReadTimeout)
	v.SetDefault("server.write_timeout", DefaultWriteTimeout)
	v.SetDefault("server.idle_timeout", DefaultIdleTimeout)
	v.SetDefault("server.shutdown_timeout", DefaultShutdownTimeout)

	v.SetDefault("log.level", DefaultLogLevel)
	v.SetDefault("log.format", DefaultLogFormat)

	v.SetDefault("database.path", DefaultDatabasePath)
	v.SetDefault("database.busy_timeout", DefaultDatabaseBusyTimeout)
	v.SetDefault("database.max_open_conns", 1)
	v.SetDefault("database.max_idle_conns", 1)

	v.SetDefault("auth.jwt_secret", "")
	v.SetDefault("auth.token_ttl", DefaultTokenTTL)
	v.SetDefault("auth.api_token_ttl", DefaultAPITokenTTL)

	v.SetDefault("docker.host", "")
	v.SetDefault("docker.api_version", "")
	v.SetDefault("docker.publish_host", DefaultHost)

	v.SetDefault("health.expose_details", ExposeDetailsLocal)

	v.SetDefault("metrics.storage_path", "/")

	v.SetDefault("proxy.enabled", false)
	v.SetDefault("proxy.admin_address", DefaultProxyAdminAddress)
	v.SetDefault("proxy.email", "")
	v.SetDefault("proxy.caddyfile_path", "")
	v.SetDefault("proxy.network", DefaultProxyNetwork)
	v.SetDefault("proxy.caddy_container", "")
	v.SetDefault("proxy.local_certs", false)
	v.SetDefault("proxy.http_port", 0)
	v.SetDefault("proxy.https_port", 0)
	v.SetDefault("proxy.dashboard_host", "")
	v.SetDefault("proxy.dashboard_upstream", "")

	v.SetDefault("web.root", "")
}

// dockerName reprend la règle de Docker pour un nom de réseau ou de conteneur.
var dockerName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,127}$`)

// validate vérifie la section proxy quand elle est active. L'adresse
// d'administration et le courriel sont vérifiés par le paquet proxy au
// démarrage, avec ses listes blanches (ADR-0007).
func (p ProxyConfig) validate(environment string) error {
	if !p.Enabled {
		if p.DashboardHost != "" {
			return fmt.Errorf("proxy.dashboard_host demande proxy.enabled : Caddy publie le dashboard")
		}
		return nil
	}
	if (p.DashboardHost == "") != (p.DashboardUpstream == "") {
		return fmt.Errorf("proxy.dashboard_host et proxy.dashboard_upstream vont ensemble")
	}
	if !dockerName.MatchString(p.Network) {
		return fmt.Errorf("proxy.network invalide : %q", p.Network)
	}
	if p.CaddyContainer != "" && !dockerName.MatchString(p.CaddyContainer) {
		return fmt.Errorf("proxy.caddy_container invalide : %q", p.CaddyContainer)
	}
	if p.LocalCerts && environment == "production" {
		return fmt.Errorf("proxy.local_certs est réservé aux essais : refusé en production")
	}
	for name, port := range map[string]int{"proxy.http_port": p.HTTPPort, "proxy.https_port": p.HTTPSPort} {
		if port < 0 || port > 65535 {
			return fmt.Errorf("%s invalide : %d (0 pour le port par défaut, ou 1 à 65535)", name, port)
		}
	}
	return nil
}
