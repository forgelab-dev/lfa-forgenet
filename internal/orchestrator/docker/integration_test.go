package docker_test

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/ids"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/docker"
)

// integrationImage tourne en continu et écrit dans ses journaux au démarrage.
const integrationImage = "nginx:alpine"

// TestEngine_RealDockerLifecycle parcourt le cycle complet sur un vrai démon :
// création, démarrage, journaux, redémarrage, arrêt, suppression. Il est sauté
// avec -short ou quand aucun démon n'est joignable.
func TestEngine_RealDockerLifecycle(t *testing.T) {
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

	hostPort := freeTCPPort(t)
	created, err := engine.Create(ctx, orchestrator.CreateSpec{
		Name:   ids.New("forgenet-test"),
		Image:  integrationImage,
		Ports:  []orchestrator.Port{{ContainerPort: 80, HostPort: hostPort}},
		Labels: map[string]string{"forgenet.test": "true"},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = engine.Remove(context.Background(), created.ID, true) })
	assert.Equal(t, orchestrator.StateCreated, created.State)

	started, err := engine.Start(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, orchestrator.StateRunning, started.State)
	require.Contains(t, started.Ports, orchestrator.Port{HostIP: "127.0.0.1", HostPort: hostPort, ContainerPort: 80, Protocol: "tcp"},
		"le port est publié sur 127.0.0.1 seulement")

	listed, err := engine.List(ctx, false)
	require.NoError(t, err)
	assert.True(t, containsID(listed, created.ID))

	require.Eventually(t, func() bool {
		var logs strings.Builder
		_ = engine.Logs(ctx, created.ID, orchestrator.LogOptions{Stdout: true, Stderr: true},
			func(c orchestrator.LogChunk) error { logs.Write(c.Data); return nil })
		return strings.Contains(logs.String(), "start")
	}, 30*time.Second, 200*time.Millisecond, "les journaux de démarrage de nginx doivent être lisibles")

	restarted, err := engine.Restart(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, orchestrator.StateRunning, restarted.State)

	stopped, err := engine.Stop(ctx, created.ID, time.Second)
	require.NoError(t, err)
	assert.Equal(t, orchestrator.StateExited, stopped.State)

	require.NoError(t, engine.Remove(ctx, created.ID, false))
	_, err = engine.Get(ctx, created.ID)
	assert.ErrorIs(t, err, orchestrator.ErrNotFound)
}

func ensureImage(ctx context.Context, t *testing.T) {
	t.Helper()
	cli, err := client.New(client.FromEnv)
	require.NoError(t, err)
	defer func() { _ = cli.Close() }()
	if _, err := cli.ImageInspect(ctx, integrationImage); err == nil {
		return
	}
	pull, err := cli.ImagePull(ctx, integrationImage, client.ImagePullOptions{})
	require.NoError(t, err)
	require.NoError(t, pull.Wait(ctx))
}

func freeTCPPort(t *testing.T) uint16 {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = ln.Close() }()
	return uint16(ln.Addr().(*net.TCPAddr).Port)
}

func containsID(containers []orchestrator.Container, id string) bool {
	for _, c := range containers {
		if c.ID == id {
			return true
		}
	}
	return false
}
