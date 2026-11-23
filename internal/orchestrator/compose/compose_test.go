package compose_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
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

// resolved est la sortie de `config --format json` pour stack, dont le service
// a déjà une politique : aucun fichier complémentaire n'est ajouté.
const resolved = `{"name":"forgenet-validation","services":{"web":{"image":"nginx:alpine","restart":"always"}}}`

// fakeRunner enregistre les commandes et répond selon la sous-commande.
type fakeRunner struct {
	calls     []compose.Command
	outputs   map[string][]byte
	errs      map[string]error
	overrides []string // contenu des fichiers complémentaires, lu pendant la commande
}

func (f *fakeRunner) Run(_ context.Context, cmd compose.Command) ([]byte, error) {
	f.calls = append(f.calls, cmd)
	// Le premier -f est l'entrée standard ; les suivants sont des fichiers.
	for i, a := range cmd.Args {
		if a == "-f" && i+1 < len(cmd.Args) && cmd.Args[i+1] != "-" {
			content, err := os.ReadFile(cmd.Args[i+1])
			if err != nil {
				return nil, err
			}
			f.overrides = append(f.overrides, string(content))
		}
	}
	sub := subcommand(cmd.Args)
	return f.outputs[sub], f.errs[sub]
}

// deployRunner répond à la lecture du modèle par config.
func deployRunner(config string) *fakeRunner {
	return &fakeRunner{outputs: map[string][]byte{"config": []byte(config)}}
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
	r := deployRunner(resolved)
	e := newEngine(t, config.DockerConfig{}, r)

	require.NoError(t, e.Deploy(context.Background(), "prj_0a1b2c3d4e5f6a7b", []byte(stack), nil, nil))
	require.Len(t, r.calls, 2)

	assert.Equal(t, []string{"-p", "forgenet-validation", "-f", "-", "config", "--format", "json"}, r.calls[0].Args)
	assert.Equal(t, []string{
		"-p", "prj_0a1b2c3d4e5f6a7b", "--project-directory", r.calls[1].Dir, "-f", "-", "up",
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
		err := e.Deploy(context.Background(), name, []byte(stack), nil, nil)
		assert.ErrorIs(t, err, compose.ErrInvalidProject, "%q", name)
	}
	assert.Empty(t, r.calls)
}

func TestDeploy_InvalidFileDoesNotTouchTheStack(t *testing.T) {
	r := &fakeRunner{errs: map[string]error{
		"config": &compose.CommandError{Args: []string{"config"}, ExitCode: 15, Output: "services.web.image must be a string"},
	}}
	e := newEngine(t, config.DockerConfig{}, r)

	err := e.Deploy(context.Background(), "prj_1", []byte(stack), nil, nil)
	require.ErrorIs(t, err, compose.ErrInvalidFile)
	assert.Contains(t, err.Error(), "services.web.image must be a string", "le message de Compose est gardé")
	assert.Len(t, r.calls, 1, "up n'est pas lancé")
}

func TestValidate_ClassifiesErrors(t *testing.T) {
	ctx := context.Background()

	e := newEngine(t, config.DockerConfig{}, &fakeRunner{})
	assert.ErrorIs(t, e.Validate(ctx, []byte("  \n"), nil), compose.ErrInvalidFile, "un fichier vide est refusé sans commande")

	// Une annulation ou une commande absente ne disent rien du fichier.
	for _, cause := range []error{context.Canceled, compose.ErrUnavailable} {
		r := &fakeRunner{errs: map[string]error{"config": cause}}
		err := newEngine(t, config.DockerConfig{}, r).Validate(ctx, []byte(stack), nil)
		assert.ErrorIs(t, err, cause)
		assert.NotErrorIs(t, err, compose.ErrInvalidFile)
	}
}

// La vérification des templates (ADR-0018) lit le modèle sans remplacer les
// ${NOM} : ils doivent rester visibles pour être croisés avec le manifeste.
func TestResolve_KeepsVariablesUninterpolated(t *testing.T) {
	r := deployRunner(resolved)
	out, err := newEngine(t, config.DockerConfig{}, r).Resolve(context.Background(), []byte(stack))
	require.NoError(t, err)

	assert.Equal(t, resolved, string(out))
	require.Len(t, r.calls, 1)
	assert.Equal(t, []string{"-p", "forgenet-validation", "-f", "-", "config", "--no-interpolate", "--format", "json"}, r.calls[0].Args)

	failure := &compose.CommandError{Args: []string{"config"}, ExitCode: 15, Output: "services.web.image must be a string"}
	_, err = newEngine(t, config.DockerConfig{}, &fakeRunner{errs: map[string]error{"config": failure}}).Resolve(context.Background(), []byte(stack))
	assert.ErrorIs(t, err, compose.ErrInvalidFile)
}

func TestDeploy_StartFailureIsReturned(t *testing.T) {
	failure := &compose.CommandError{Args: []string{"up"}, ExitCode: 1, Output: "container web is unhealthy"}
	r := deployRunner(resolved)
	r.errs = map[string]error{"up": failure}
	err := newEngine(t, config.DockerConfig{}, r).Deploy(context.Background(), "prj_1", []byte(stack), nil, nil)

	var cmdErr *compose.CommandError
	require.ErrorAs(t, err, &cmdErr)
	assert.NotErrorIs(t, err, compose.ErrInvalidFile, "un échec au démarrage n'est pas un fichier invalide")
	assert.Equal(t, "docker compose up : code 1 : container web is unhealthy", err.Error())
}

// Un service sans politique de redémarrage reçoit unless-stopped par un
// fichier complémentaire ; une politique écrite par l'utilisateur, sous l'une
// ou l'autre forme, est gardée, et une tâche ponctuelle n'est jamais relancée.
func TestDeploy_AddsRestartPolicyOnlyWhereMissing(t *testing.T) {
	model := `{"name":"forgenet-validation","services":{
		"web":      {"image":"nginx:alpine"},
		"cache":    {"image":"redis:7-alpine","restart":"no"},
		"toujours": {"image":"nginx:alpine","restart":"always"},
		"swarm":    {"image":"nginx:alpine","deploy":{"resources":{},"restart_policy":{"condition":"on-failure"}}},
		"limites":  {"image":"nginx:alpine","deploy":{"resources":{"limits":{"memory":"64m"}}}},
		"migrate":  {"image":"alpine","command":["true"]},
		"app":      {"image":"nginx:alpine","depends_on":{
			"migrate":{"condition":"service_completed_successfully","required":true},
			"cache":{"condition":"service_started","required":true}}}
	}}`
	r := deployRunner(model)
	e := newEngine(t, config.DockerConfig{}, r)
	var journal []string
	out := func(_ compose.OutputStream, line string) { journal = append(journal, line) }

	require.NoError(t, e.Deploy(context.Background(), "prj_1", []byte(stack), nil, out))
	require.Len(t, r.calls, 2)
	up := r.calls[1].Args
	require.Len(t, up, 15)
	// Le dossier du projet reste le dossier de travail : un chemin relatif du
	// fichier ne se résout pas dans le dossier du fichier complémentaire.
	assert.Equal(t, []string{"-p", "prj_1", "--project-directory", r.calls[1].Dir, "-f", "-", "-f"}, up[:7])
	assert.Equal(t, []string{"up", "--detach", "--wait", "--wait-timeout", "90", "--remove-orphans", "--quiet-pull"}, up[8:])
	assert.Equal(t, stack, string(r.calls[1].Stdin), "le fichier de l'utilisateur passe tel quel")

	require.Len(t, r.overrides, 1)
	assert.JSONEq(t, `{"services":{
		"app":{"restart":"unless-stopped"},
		"limites":{"restart":"unless-stopped"},
		"web":{"restart":"unless-stopped"}}}`, r.overrides[0])
	assert.Equal(t, []string{"Relance après un redémarrage du serveur (restart: unless-stopped) ajoutée à : app, limites, web"}, journal)

	// Le fichier complémentaire n'est pas dans le dossier de travail, qui reste
	// vide, et il est supprimé après la commande.
	override := up[7]
	assert.NotEqual(t, r.calls[1].Dir, filepath.Dir(override))
	entries, err := os.ReadDir(r.calls[1].Dir)
	require.NoError(t, err)
	assert.Empty(t, entries)
	assert.NoFileExists(t, override)
}

func TestDeploy_RemovesOverrideWhenStartFails(t *testing.T) {
	r := deployRunner(`{"services":{"web":{"image":"nginx:alpine"},"worker":{"image":"alpine"}}}`)
	r.errs = map[string]error{"up": &compose.CommandError{Args: []string{"up"}, ExitCode: 1, Output: "boom"}}

	err := newEngine(t, config.DockerConfig{}, r).Deploy(context.Background(), "prj_1", []byte(stack), nil, nil)
	require.Error(t, err)
	require.Len(t, r.overrides, 1)
	assert.JSONEq(t, `{"services":{"web":{"restart":"unless-stopped"},"worker":{"restart":"unless-stopped"}}}`, r.overrides[0])
	assert.NoFileExists(t, r.calls[1].Args[7])
}

func TestDeploy_UnreadableModelStopsBeforeUp(t *testing.T) {
	r := deployRunner("pas du json")
	err := newEngine(t, config.DockerConfig{}, r).Deploy(context.Background(), "prj_1", []byte(stack), nil, nil)
	require.Error(t, err)
	assert.Len(t, r.calls, 1, "up n'est pas lancé")
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

// Les valeurs d'une installation (ADR-0019) s'ajoutent à l'environnement de
// la commande, pour la lecture du modèle comme pour up.
func TestDeploy_PassesVariablesToEveryCommand(t *testing.T) {
	r := deployRunner(resolved)
	e := newEngine(t, config.DockerConfig{}, r)
	vars := []string{"POSTGRES_PASSWORD=s3cret", "POSTGRES_USER=app"}

	require.NoError(t, e.Deploy(context.Background(), "prj_1", []byte(stack), vars, nil))
	require.Len(t, r.calls, 2)
	for _, c := range r.calls {
		assert.Contains(t, c.Env, "POSTGRES_PASSWORD=s3cret", subcommand(c.Args))
		assert.Contains(t, c.Env, "POSTGRES_USER=app", subcommand(c.Args))
		assert.Contains(t, c.Env, "PATH="+os.Getenv("PATH"), "l'environnement fixe reste")
	}

	r = &fakeRunner{}
	require.NoError(t, newEngine(t, config.DockerConfig{}, r).Validate(context.Background(), []byte(stack), vars))
	assert.Contains(t, r.calls[0].Env, "POSTGRES_PASSWORD=s3cret", "${VAR:?} serait refusé sans la valeur")
}

// Le manifeste interdit déjà ces noms ; le moteur ne s'y fie pas : une valeur
// ne doit jamais changer le démon visé ni la commande elle-même.
func TestDeploy_RefusesReservedVariables(t *testing.T) {
	for _, kv := range []string{
		"PATH=/tmp", "HOME=/tmp", "DOCKER_HOST=tcp://ailleurs:2375", "COMPOSE_FILE=/etc/passwd",
		"FORGENET_SECRET_KEY=x", "SANS_EGAL", "=vide", "minuscules=x",
	} {
		r := deployRunner(resolved)
		err := newEngine(t, config.DockerConfig{}, r).Deploy(context.Background(), "prj_1", []byte(stack), []string{kv}, nil)
		require.ErrorIs(t, err, compose.ErrReservedVariable, kv)
		assert.Empty(t, r.calls, "%s : aucune commande lancée", kv)
	}
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

func TestCommandError_SummaryKeepsTheCause(t *testing.T) {
	err := &compose.CommandError{Args: []string{"-p", "x", "up"}, ExitCode: 1, Output: " Image web Pulling \n Image web Error pull access denied\nError response from daemon: pull access denied for web\n"}
	assert.Equal(t, "docker compose up : Error response from daemon: pull access denied for web", err.Summary())
	assert.Contains(t, err.Error(), "Pulling", "Error garde toute la sortie")
}
