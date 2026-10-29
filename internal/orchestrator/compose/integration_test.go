package compose_test

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/ids"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
)

// TestEngine_RealComposeLifecycle déploie la stack de référence sur un vrai
// démon, la met à jour, puis la supprime. Il est sauté avec -short ou quand
// docker compose ne répond pas.
func TestEngine_RealComposeLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("test d'intégration Compose sauté avec -short")
	}
	// Lue au démarrage du moteur : elle ne doit pas atteindre les conteneurs.
	t.Setenv("FORGENET_TEST_SECRET", "fuite")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	engine, err := compose.New(config.DockerConfig{})
	require.NoError(t, err)
	if err := engine.Ping(ctx); err != nil {
		t.Skipf("docker compose indisponible : %v", err)
	}
	if err := exec.CommandContext(ctx, "docker", "info").Run(); err != nil {
		t.Skipf("démon Docker indisponible : %v", err)
	}

	project := ids.New("forgenet-test")
	t.Cleanup(func() {
		cleanup, done := context.WithTimeout(context.Background(), time.Minute)
		defer done()
		_ = engine.Remove(cleanup, project, true)
	})

	// Déploiement : les deux services tournent, web est sain.
	require.NoError(t, engine.Deploy(ctx, project, readFile(t, "testdata/reference.compose.yaml")))
	first := byService(ctx, t, engine, project)
	require.Len(t, first, 2)
	assert.Equal(t, orchestrator.StateRunning, first["web"].State)
	assert.Equal(t, "healthy", first["web"].Health, "--wait attend la sonde de santé")
	assert.Equal(t, orchestrator.StateRunning, first["cache"].State)
	assert.Equal(t, "absent", containerEnv(ctx, t, first["cache"].ID, "FUITE"),
		"le fichier ne lit pas l'environnement du serveur")

	// Mise à jour : web n'est pas recréé, cache est retiré, worker apparaît.
	require.NoError(t, engine.Deploy(ctx, project, readFile(t, "testdata/reference-v2.compose.yaml")))
	second := byService(ctx, t, engine, project)
	require.Len(t, second, 2)
	assert.Equal(t, first["web"].ID, second["web"].ID, "un service inchangé garde son conteneur")
	assert.Equal(t, orchestrator.StateRunning, second["worker"].State)
	assert.NotContains(t, second, "cache", "un service retiré du fichier est supprimé")

	// Suppression : plus aucun conteneur, mais le volume nommé reste.
	volume := project + "_donnees"
	require.NoError(t, engine.Remove(ctx, project, false))
	assert.Empty(t, byService(ctx, t, engine, project))
	assert.NoError(t, exec.CommandContext(ctx, "docker", "volume", "inspect", volume).Run(),
		"les données survivent à la suppression de la stack")

	require.NoError(t, engine.Remove(ctx, project, true))
	assert.Error(t, exec.CommandContext(ctx, "docker", "volume", "inspect", volume).Run(),
		"--volumes retire les données")
}

func TestEngine_RealComposeRejectsInvalidFile(t *testing.T) {
	if testing.Short() {
		t.Skip("test d'intégration Compose sauté avec -short")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	engine, err := compose.New(config.DockerConfig{})
	require.NoError(t, err)
	if err := engine.Ping(ctx); err != nil {
		t.Skipf("docker compose indisponible : %v", err)
	}

	err = engine.Validate(ctx, []byte("services:\n  web:\n    image: [pas, une, chaine]\n"))
	require.ErrorIs(t, err, compose.ErrInvalidFile)
	var cmdErr *compose.CommandError
	require.ErrorAs(t, err, &cmdErr)
	assert.NotEmpty(t, cmdErr.Output, "le message de Compose explique le refus")
}

func TestEngine_MissingDockerIsUnavailable(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	engine, err := compose.New(config.DockerConfig{})
	require.NoError(t, err)
	assert.ErrorIs(t, engine.Ping(context.Background()), compose.ErrUnavailable)
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	return content
}

func byService(ctx context.Context, t *testing.T, engine *compose.Engine, project string) map[string]compose.Service {
	t.Helper()
	services, err := engine.Services(ctx, project)
	require.NoError(t, err)
	out := make(map[string]compose.Service, len(services))
	for _, s := range services {
		out[s.Service] = s
	}
	return out
}

// containerEnv lit une variable d'environnement du conteneur.
func containerEnv(ctx context.Context, t *testing.T, id, name string) string {
	t.Helper()
	out, err := exec.CommandContext(ctx, "docker", "inspect", "--format", "{{json .Config.Env}}", id).Output()
	require.NoError(t, err)
	for _, kv := range strings.Split(strings.Trim(strings.TrimSpace(string(out)), "[]"), ",") {
		kv = strings.Trim(kv, `"`)
		if value, ok := strings.CutPrefix(kv, name+"="); ok {
			return value
		}
	}
	return ""
}
