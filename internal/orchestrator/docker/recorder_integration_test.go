package docker_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/ids"
	"github.com/forgelab-dev/lfa-forgenet/internal/metrics"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/docker"
)

// TestRecorder_RealContainer fait tourner le relevé de fond (ADR-0010) sur un
// vrai démon : le conteneur lancé apparaît dans les passages, avec un
// historique, une mémoire et une limite lues par le démon, puis il garde son
// historique une fois arrêté.
func TestRecorder_RealContainer(t *testing.T) {
	if testing.Short() {
		t.Skip("test d'intégration Docker sauté avec -short")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	engine, err := docker.New(config.DockerConfig{PublishHost: "127.0.0.1"})
	require.NoError(t, err)
	defer func() { _ = engine.Close() }()
	if err := engine.Ping(ctx); err != nil {
		t.Skipf("démon Docker indisponible : %v", err)
	}
	ensureImage(ctx, t)

	created, err := engine.Create(ctx, orchestrator.CreateSpec{
		Name:        ids.New("forgenet-test"),
		Image:       integrationImage,
		Labels:      map[string]string{"forgenet.test": "true"},
		MemoryBytes: 64 << 20,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = engine.Remove(context.Background(), created.ID, true) })
	_, err = engine.Start(ctx, created.ID)
	require.NoError(t, err)

	running, err := engine.RunningIDs(ctx)
	require.NoError(t, err)
	assert.True(t, slices.Contains(running, created.ID), "identifiant complet dans la liste")

	recorder := metrics.NewRecorder(
		metrics.NewHostCollector(metrics.GopsutilProbe{}, engine, "/"),
		engine, engine,
		metrics.RecorderConfig{HostEvery: time.Second, ContainerEvery: time.Second},
		zerolog.Nop())
	go func() { _ = recorder.Start() }()
	t.Cleanup(func() { _ = recorder.Shutdown(context.Background()) })

	var history []metrics.Container
	require.Eventually(t, func() bool {
		history, _ = recorder.ContainerHistory(created.ID, time.Time{}, time.Now().Add(time.Minute))
		return len(history) >= 3
	}, time.Minute, 200*time.Millisecond, "trois relevés du conteneur")

	last := history[len(history)-1]
	assert.Equal(t, uint64(64<<20), last.Memory.LimitBytes, "limite posée à la création")
	assert.NotZero(t, last.Memory.UsageBytes)
	assert.NotZero(t, last.CPU.OnlineCPUs)

	sweep := recorder.LatestSweep()
	assert.True(t, slices.ContainsFunc(sweep.Containers, func(c metrics.Container) bool { return c.ContainerID == created.ID }))
	_, ok := recorder.LatestContainer(created.ID[:metrics.MinPrefix])
	assert.True(t, ok, "identifiant abrégé comme dans `docker ps`")
	host, ok := recorder.LatestHost()
	require.True(t, ok)
	assert.NotZero(t, host.Memory.LimitBytes, "mémoire de l'hôte lue")

	_, err = engine.Stop(ctx, created.ID, time.Second)
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		return !slices.ContainsFunc(recorder.LatestSweep().Containers, func(c metrics.Container) bool { return c.ContainerID == created.ID })
	}, 30*time.Second, 200*time.Millisecond, "un conteneur arrêté sort du passage")
	_, found := recorder.ContainerHistory(created.ID, time.Time{}, time.Now())
	assert.True(t, found, "son historique reste lisible après l'arrêt")
}
