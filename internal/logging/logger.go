// Package logging configure le logger structuré Zerolog pour l'ensemble des modules.
package logging

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
)

// New initialise et retourne une nouvelle instance de zerolog.Logger configurée selon LogConfig.
func New(cfg config.LogConfig) zerolog.Logger {
	level := parseLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	var output io.Writer = os.Stdout
	if strings.ToLower(cfg.Format) == "console" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	logger := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()

	// Initialise le logger par défaut de zerolog
	log.Logger = logger

	return logger
}

func parseLevel(levelStr string) zerolog.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zerolog.DebugLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
