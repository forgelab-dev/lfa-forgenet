package health_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/forgelab-dev/lfa-forgenet/internal/health"
)

func TestService_Liveness(t *testing.T) {
	svc := health.NewService()
	ctx := context.Background()

	assert.True(t, svc.IsLive(ctx))

	svc.SetServing(false)
	assert.False(t, svc.IsLive(ctx))
}

func TestService_Readiness(t *testing.T) {
	svc := health.NewService()
	ctx := context.Background()

	// Initialement prêt (aucun sous-système critique défaillant)
	ready, details := svc.IsReady(ctx)
	assert.True(t, ready)
	assert.Empty(t, details)

	// Ajout d'un sous-système OK
	svc.RegisterReadinessChecker("sqlite", health.CheckFunc(func(_ context.Context) error {
		return nil
	}))

	ready, details = svc.IsReady(ctx)
	assert.True(t, ready)
	assert.Equal(t, "ok", details["sqlite"])

	// Ajout d'un sous-système en erreur
	svc.RegisterReadinessChecker("docker", health.CheckFunc(func(_ context.Context) error {
		return errors.New("daemon docker inaccessible")
	}))

	ready, details = svc.IsReady(ctx)
	assert.False(t, ready)
	assert.Equal(t, "ok", details["sqlite"])
	assert.Equal(t, "daemon docker inaccessible", details["docker"])

	// Si non serving, readiness retourne faux immédiatement
	svc.SetServing(false)
	ready, _ = svc.IsReady(ctx)
	assert.False(t, ready)
}
