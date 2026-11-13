package docker_test

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/ids"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/docker"
)

// TestEngine_RealNetwork crée un réseau, y branche un conteneur avec un alias,
// vérifie que chaque opération peut être relancée, puis refuse un réseau du
// même nom créé hors de ForgeNet. Sauté avec -short ou sans démon Docker.
func TestEngine_RealNetwork(t *testing.T) {
	if testing.Short() {
		t.Skip("test d'intégration Docker sauté avec -short")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := exec.CommandContext(ctx, "docker", "info").Run(); err != nil {
		t.Skipf("démon Docker indisponible : %v", err)
	}

	engine, err := docker.New(config.DockerConfig{PublishHost: "127.0.0.1"})
	require.NoError(t, err)
	defer func() { _ = engine.Close() }()

	suffix := strings.TrimPrefix(ids.New("t"), "t_")
	netName, foreign, ctr := "fn-net-"+suffix, "fn-foreign-"+suffix, "fn-ctr-"+suffix
	labels := map[string]string{"forgenet.managed": "edge"}
	t.Cleanup(func() {
		cleanup("rm", "-f", ctr)
		cleanup("network", "rm", netName, foreign)
	})

	require.NoError(t, engine.EnsureNetwork(ctx, netName, labels))
	require.NoError(t, engine.EnsureNetwork(ctx, netName, labels), "un second appel ne recrée rien")
	assert.Equal(t, "edge", strings.TrimSpace(dockerOut(ctx, t, "network", "inspect", "-f", `{{index .Labels "forgenet.managed"}}`, netName)))

	run(ctx, t, "run", "-d", "--rm", "--name", ctr, integrationImage)
	alias := "web.boutique.forgenet.internal"
	require.NoError(t, engine.ConnectToNetwork(ctx, netName, ctr, []string{alias}))
	require.NoError(t, engine.ConnectToNetwork(ctx, netName, ctr, []string{alias}), "déjà branché : rien à faire")
	aliases := dockerOut(ctx, t, "inspect", "-f", `{{json (index .NetworkSettings.Networks "`+netName+`").Aliases}}`, ctr)
	assert.Contains(t, aliases, alias)

	require.NoError(t, engine.DisconnectFromNetwork(ctx, netName, ctr))
	require.NoError(t, engine.DisconnectFromNetwork(ctx, netName, ctr), "déjà débranché : rien à faire")
	networks := dockerOut(ctx, t, "inspect", "-f", `{{json .NetworkSettings.Networks}}`, ctr)
	assert.NotContains(t, networks, netName)

	run(ctx, t, "network", "create", foreign)
	assert.ErrorIs(t, engine.EnsureNetwork(ctx, foreign, labels), docker.ErrForeignNetwork)
}

func run(ctx context.Context, t *testing.T, args ...string) {
	t.Helper()
	dockerOut(ctx, t, args...)
}

func dockerOut(ctx context.Context, t *testing.T, args ...string) string {
	t.Helper()
	out, err := exec.CommandContext(ctx, "docker", args...).CombinedOutput()
	require.NoError(t, err, "docker %s : %s", strings.Join(args, " "), out)
	return string(out)
}

func cleanup(args ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = exec.CommandContext(ctx, "docker", args...).Run()
}
