// Package config centralise la configuration de ForgeNet via Viper.
package config

import (
	"errors"
	"fmt"
	"net/netip"
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

// Config est la structure racine de la configuration ForgeNet.
type Config struct {
	Environment string         `mapstructure:"environment"` // "development", "staging", "production"
	Server      ServerConfig   `mapstructure:"server"`
	Log         LogConfig      `mapstructure:"log"`
	Database    DatabaseConfig `mapstructure:"database"`
	Auth        AuthConfig     `mapstructure:"auth"`
	Docker      DockerConfig   `mapstructure:"docker"`
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
}
