package logging_test

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/logging"
)

func TestNew_LoggerLevels(t *testing.T) {
	tests := []struct {
		levelStr      string
		expectedLevel zerolog.Level
	}{
		{"debug", zerolog.DebugLevel},
		{"info", zerolog.InfoLevel},
		{"warn", zerolog.WarnLevel},
		{"warning", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"unknown", zerolog.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.levelStr, func(t *testing.T) {
			cfg := config.LogConfig{
				Level:  tt.levelStr,
				Format: "json",
			}
			logger := logging.New(cfg)
			assert.Equal(t, tt.expectedLevel, logger.GetLevel())
		})
	}
}

func TestNew_Output(t *testing.T) {
	var buf bytes.Buffer
	cfg := config.LogConfig{
		Level:  "info",
		Format: "json",
	}
	logger := logging.New(cfg).Output(&buf)

	logger.Info().Str("test_key", "test_val").Msg("test message")

	output := buf.String()
	assert.Contains(t, output, `"level":"info"`)
	assert.Contains(t, output, `"test_key":"test_val"`)
	assert.Contains(t, output, `"message":"test message"`)
}
