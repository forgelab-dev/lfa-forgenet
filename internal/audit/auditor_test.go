package audit_test

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
)

func TestAuditor_LogAndList(t *testing.T) {
	ctx := context.Background()
	cfg := config.DatabaseConfig{
		Path:         ":memory:",
		BusyTimeout:  2 * time.Second,
		MaxOpenConns: 1,
		MaxIdleConns: 1,
	}

	st, err := sqlite.Open(ctx, cfg)
	require.NoError(t, err)
	defer st.Close()

	aud := audit.NewAuditor(st.Audit(), zerolog.Nop())

	// Enregistrement d'un premier événement avec contexte
	enrichedCtx := audit.WithAuditContext(ctx, "usr_admin_1", "10.0.0.1")
	err = aud.Log(enrichedCtx, "container.create", "container:web-app", map[string]any{
		"image": "nginx:alpine",
		"port":  80,
	})
	require.NoError(t, err)

	// Enregistrement d'un second événement
	err = aud.Log(enrichedCtx, "container.start", "container:web-app", nil)
	require.NoError(t, err)

	// Liste des événements
	events, err := aud.List(ctx, 10, 0)
	require.NoError(t, err)
	require.Len(t, events, 2)

	assert.Equal(t, "container.start", events[0].Action)
	assert.Equal(t, "usr_admin_1", events[0].UserID)
	assert.Equal(t, "10.0.0.1", events[0].IPAddress)

	assert.Equal(t, "container.create", events[1].Action)
	assert.Contains(t, events[1].Metadata, `"image":"nginx:alpine"`)
}
