package compose_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
)

const stack = "services:\n  web:\n    image: nginx:alpine\n"

// fakeRunner enregistre les commandes et répond selon la sous-commande.
type fakeRunner struct {
	calls   []compose.Command
	outputs map[string][]byte
	errs    map[string]error
}

func (f *fakeRunner) Run(_ context.Context, cmd compose.Command) ([]byte, error) {
	f.calls = append(f.calls, cmd)
	sub := subcommand(cmd.Args)
	return f.outputs[sub], f.errs[sub]
}

func subcommand(args []string) string {
	for _, a := range args {
		switch a {
		case "config", "up", "down", "ps", "version":
			return a
		}
	}
	return ""
}

func newEngine(t *testing.T, cfg config.DockerConfig, r *fakeRunner) *compose.Engine {
	t.Helper()
	e, err := compose.New(cfg, compose.WithRunner(r), compose.WithWaitTimeout(90*time.Second))
	require.NoError(t, err)
	return e
}

func TestDeploy_ValidatesThenStartsAndWaits(t *testing.T) {
	r := &fakeRunner{}
	e := newEngine(t, config.DockerConfig{}, r)

	require.NoError(t, e.Deploy(context.Background(), "prj_0a1b2c3d4e5f6a7b", []byte(stack)))
	require.Len(t, r.calls, 2)

	assert.Equal(t, []string{"-p", "forgenet-validation", "-f", "-", "config", "--quiet"}, r.calls[0].Args)
	assert.Equal(t, []string{
		"-p", "prj_0a1b2c3d4e5f6a7b", "-f", "-", "up",
		"--detach", "--wait", "--wait-timeout", "90", "--remove-orphans", "--quiet-pull",
	}, r.calls[1].Args)
	for _, c := range r.calls {
		assert.Equal(t, stack, string(c.Stdin), "le fichier passe par l'entrée standard")
	}

	// Le dossier de travail existe et ne contient rien : un chemin relatif du
	// fichier ne lit aucun fichier du serveur.
	entries, err := os.ReadDir(r.calls[1].Dir)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

func TestDeploy_RefusesBadProjectBeforeAnyCommand(t *testing.T) {
	r := &fakeRunner{}
	e := newEngine(t, config.DockerConfig{}, r)

	for _, name := range []string{"", "Majuscules", "-tiret", "avec espace", "point.final"} {
		err := e.Deploy(context.Background(), name, []byte(stack))
		assert.ErrorIs(t, err, compose.ErrInvalidProject, "%q", name)
	}
	assert.Empty(t, r.calls)
}

func TestDeploy_InvalidFileDoesNotTouchTheStack(t *testing.T) {
	r := &fakeRunner{errs: map[string]error{
		"config": &compose.CommandError{Args: []string{"config"}, ExitCode: 15, Output: "services.web.image must be a string"},
	}}
	e := newEngine(t, config.DockerConfig{}, r)

	err := e.Deploy(context.Background(), "prj_1", []byte(stack))
	require.ErrorIs(t, err, compose.ErrInvalidFile)
	assert.Contains(t, err.Error(), "services.web.image must be a string", "le message de Compose est gardé")
	assert.Len(t, r.calls, 1, "up n'est pas lancé")
}

func TestValidate_ClassifiesErrors(t *testing.T) {
	ctx := context.Background()

	e := newEngine(t, config.DockerConfig{}, &fakeRunner{})
	assert.ErrorIs(t, e.Validate(ctx, []byte("  \n")), compose.ErrInvalidFile, "un fichier vide est refusé sans commande")

	// Une annulation ou une commande absente ne disent rien du fichier.
	for _, cause := range []error{context.Canceled, compose.ErrUnavailable} {
		r := &fakeRunner{errs: map[string]error{"config": cause}}
		err := newEngine(t, config.DockerConfig{}, r).Validate(ctx, []byte(stack))
		assert.ErrorIs(t, err, cause)
		assert.NotErrorIs(t, err, compose.ErrInvalidFile)
	}
}

func TestDeploy_StartFailureIsReturned(t *testing.T) {
	failure := &compose.CommandError{Args: []string{"up"}, ExitCode: 1, Output: "container web is unhealthy"}
	r := &fakeRunner{errs: map[string]error{"up": failure}}
	err := newEngine(t, config.DockerConfig{}, r).Deploy(context.Background(), "prj_1", []byte(stack))

	var cmdErr *compose.CommandError
	require.ErrorAs(t, err, &cmdErr)
	assert.NotErrorIs(t, err, compose.ErrInvalidFile, "un échec au démarrage n'est pas un fichier invalide")
	assert.Equal(t, "docker compose up : code 1 : container web is unhealthy", err.Error())
}

func TestRemove_KeepsVolumesUnlessAsked(t *testing.T) {
	r := &fakeRunner{}
	e := newEngine(t, config.DockerConfig{}, r)

	require.NoError(t, e.Remove(context.Background(), "prj_1", false))
	require.NoError(t, e.Remove(context.Background(), "prj_1", true))
	assert.Equal(t, []string{"-p", "prj_1", "down", "--remove-orphans", "--timeout", "10"}, r.calls[0].Args)
	assert.Equal(t, []string{"-p", "prj_1", "down", "--remove-orphans", "--timeout", "10", "--volumes"}, r.calls[1].Args)
	assert.Nil(t, r.calls[0].Stdin, "down n'a pas besoin du fichier")

	assert.ErrorIs(t, e.Remove(context.Background(), "../autre", false), compose.ErrInvalidProject)
}

func TestServices_ReadsBothOutputFormats(t *testing.T) {
	lines := `{"ID":"c1","Name":"prj_1-web-1","Service":"web","Image":"nginx:alpine","State":"running","Health":"healthy"}
{"ID":"c2","Name":"prj_1-cache-1","Service":"cache","Image":"redis:7-alpine","State":"exited","Health":""}
`
	array := `[{"ID":"c1","Name":"prj_1-web-1","Service":"web","Image":"nginx:alpine","State":"running","Health":"healthy"},` +
		`{"ID":"c2","Name":"prj_1-cache-1","Service":"cache","Image":"redis:7-alpine","State":"exited","Health":""}]`
	want := []compose.Service{
		{ID: "c1", Name: "prj_1-web-1", Service: "web", Image: "nginx:alpine", State: orchestrator.StateRunning, Health: "healthy"},
		{ID: "c2", Name: "prj_1-cache-1", Service: "cache", Image: "redis:7-alpine", State: orchestrator.StateExited},
	}

	for name, out := range map[string]string{"une ligne par conteneur": lines, "tableau": array} {
		r := &fakeRunner{outputs: map[string][]byte{"ps": []byte(out)}}
		got, err := newEngine(t, config.DockerConfig{}, r).Services(context.Background(), "prj_1")
		require.NoError(t, err, name)
		assert.Equal(t, want, got, name)
		assert.Equal(t, []string{"-p", "prj_1", "ps", "--all", "--format", "json"}, r.calls[0].Args)
	}

	r := &fakeRunner{outputs: map[string][]byte{"ps": []byte("\n")}}
	got, err := newEngine(t, config.DockerConfig{}, r).Services(context.Background(), "prj_1")
	require.NoError(t, err)
	assert.Empty(t, got, "une stack inconnue n'a aucun service")

	r = &fakeRunner{outputs: map[string][]byte{"ps": []byte("{pas du json")}}
	_, err = newEngine(t, config.DockerConfig{}, r).Services(context.Background(), "prj_1")
	assert.Error(t, err)
}

// Compose remplace ${VAR} par son environnement : un secret du serveur ne doit
// pas pouvoir être recopié dans un conteneur par le fichier.
func TestCommandEnv_OnlyPassesWhatComposeNeeds(t *testing.T) {
	t.Setenv("FORGENET_AUTH_SECRET", "ne-doit-pas-sortir")
	t.Setenv("DOCKER_HOST", "unix:///var/run/env.sock")
	t.Setenv("DOCKER_CONFIG", "/home/forgenet/.docker")

	r := &fakeRunner{}
	require.NoError(t, newEngine(t, config.DockerConfig{Host: "unix:///var/run/cfg.sock"}, r).Ping(context.Background()))
	env := r.calls[0].Env

	assert.Contains(t, env, "DOCKER_HOST=unix:///var/run/cfg.sock", "docker.host l'emporte, pour viser le même démon")
	assert.Contains(t, env, "DOCKER_CONFIG=/home/forgenet/.docker")
	assert.Contains(t, env, "PATH="+os.Getenv("PATH"))
	for _, kv := range env {
		assert.False(t, strings.HasPrefix(kv, "FORGENET_"), "variable transmise à tort : %s", kv)
	}

	r = &fakeRunner{}
	require.NoError(t, newEngine(t, config.DockerConfig{}, r).Ping(context.Background()))
	assert.Contains(t, r.calls[0].Env, "DOCKER_HOST=unix:///var/run/env.sock", "sans docker.host, l'environnement décide")
}

func TestValidProjectName(t *testing.T) {
	for _, ok := range []string{"prj_0a1b2c3d4e5f6a7b", "forgenet-test", "a", "9lives"} {
		assert.True(t, compose.ValidProjectName(ok), ok)
	}
	for _, ko := range []string{"", "_debut", "Projet", "a/b", "a.b", strings.Repeat("é", 2)} {
		assert.False(t, compose.ValidProjectName(ko), ko)
	}
}

func TestPing_ReportsUnavailable(t *testing.T) {
	r := &fakeRunner{errs: map[string]error{"version": compose.ErrUnavailable}}
	err := newEngine(t, config.DockerConfig{}, r).Ping(context.Background())
	assert.True(t, errors.Is(err, compose.ErrUnavailable))
	assert.Equal(t, []string{"version", "--short"}, r.calls[0].Args)
}
