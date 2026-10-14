// Package config centralise la configuration de ForgeNet via Viper.
package config

import (
	"errors"
	"fmt"
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

// Config est la structure racine de la configuration ForgeNet.
type Config struct {
	Environment string       `mapstructure:"environment"` // "development", "staging", "production"
	Server      ServerConfig `mapstructure:"server"`
	Log         LogConfig    `mapstructure:"log"`
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
}
