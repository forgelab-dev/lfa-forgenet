// Package compose déploie, met à jour et supprime des stacks Docker Compose en
// pilotant la commande `docker compose` (ADR-0003).
//
// Chaque stack est un projet Compose nommé ; le fichier est passé sur l'entrée
// standard, la commande ne voit que les variables d'environnement dont elle a
// besoin et travaille dans un dossier vide.
package compose

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
)

var (
	// ErrInvalidProject indique un nom de projet refusé par Compose.
	ErrInvalidProject = errors.New("nom de projet Compose invalide")
	// ErrInvalidFile indique un fichier Compose que la commande refuse.
	ErrInvalidFile = errors.New("fichier Compose invalide")
	// ErrUnavailable indique que la commande docker compose est absente.
	ErrUnavailable = errors.New("docker compose indisponible")
	// ErrReservedVariable indique une variable que le moteur ne transmet pas :
	// mal formée, ou qui changerait la commande elle-même.
	ErrReservedVariable = errors.New("variable refusée par le moteur Compose")
)

// projectName reprend la règle de Compose : minuscules, chiffres, « _ » et « - »,
// en commençant par une lettre ou un chiffre.
var projectName = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

// passedEnv liste les seules variables transmises à la commande : Compose
// remplace ${VAR} par son environnement, qui ne doit pas être celui du serveur.
var passedEnv = []string{
	"PATH", "HOME",
	"DOCKER_HOST", "DOCKER_CONTEXT", "DOCKER_CONFIG", "DOCKER_CERT_PATH", "DOCKER_TLS_VERIFY",
}

// variableName est la forme des variables d'une installation (ADR-0019),
// celle des manifestes du Marketplace.
var variableName = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// reservedPrefixes changeraient le démon visé, la commande ou le serveur.
var reservedPrefixes = []string{"DOCKER_", "COMPOSE_", "FORGENET_"}

const (
	// defaultWaitTimeout borne l'attente de services démarrés et sains.
	defaultWaitTimeout = 2 * time.Minute
	// stopTimeout laisse aux conteneurs le temps de s'arrêter avant SIGKILL.
	stopTimeout = 10 * time.Second
	// validationProject nomme le projet le temps d'une validation : `config` ne crée rien.
	validationProject = "forgenet-validation"
)

// OutputStream est le flux d'une ligne de sortie de la commande.
type OutputStream int

const (
	// Stdout est la sortie standard.
	Stdout OutputStream = iota
	// Stderr est la sortie d'erreur, où Compose écrit aussi sa progression.
	Stderr
)

// OutputFunc reçoit la sortie de la commande ligne par ligne, au fil de
// l'exécution. Elle peut être appelée depuis deux goroutines à la fois (une
// par flux) et ne doit pas bloquer longtemps.
type OutputFunc func(stream OutputStream, line string)

// Command est une exécution de `docker compose`.
type Command struct {
	Args   []string // arguments après « docker compose »
	Stdin  []byte
	Env    []string // environnement complet de la commande
	Dir    string
	Output OutputFunc // facultative : reçoit la sortie en direct
}

// Runner exécute une commande et renvoie sa sortie standard. Une sortie en
// erreur porte le message de la commande.
type Runner interface {
	Run(ctx context.Context, cmd Command) ([]byte, error)
}

// Service est l'état d'un conteneur d'une stack.
type Service struct {
	ID      string // identifiant du conteneur
	Name    string // nom du conteneur
	Service string // nom du service dans le fichier
	Image   string
	State   orchestrator.State
	Health  string // healthy, unhealthy, starting, ou vide sans sonde
}

// Engine gère les stacks Compose.
type Engine struct {
	runner      Runner
	env         []string
	dir         string
	overrideDir string // fichiers complémentaires (politique de redémarrage), hors de dir
	waitTimeout time.Duration
}

// Option règle le moteur.
type Option func(*Engine)

// WithRunner remplace l'exécution réelle, pour les tests.
func WithRunner(r Runner) Option { return func(e *Engine) { e.runner = r } }

// WithWaitTimeout change l'attente maximale de services sains.
func WithWaitTimeout(d time.Duration) Option { return func(e *Engine) { e.waitTimeout = d } }

// New prépare le moteur pour le démon de cfg (vide : celui de l'environnement).
func New(cfg config.DockerConfig, opts ...Option) (*Engine, error) {
	// Le dossier est vide : un chemin relatif du fichier ne lit rien du serveur.
	dir := filepath.Join(os.TempDir(), "forgenet-compose")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("création du dossier de travail Compose : %w", err)
	}
	overrideDir := filepath.Join(os.TempDir(), "forgenet-compose-overrides")
	if err := os.MkdirAll(overrideDir, 0o700); err != nil {
		return nil, fmt.Errorf("création du dossier des fichiers complémentaires : %w", err)
	}

	e := &Engine{
		runner:      execRunner{},
		env:         commandEnv(os.Getenv, cfg.Host),
		dir:         dir,
		overrideDir: overrideDir,
		waitTimeout: defaultWaitTimeout,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e, nil
}

// commandEnv garde les variables utiles à la commande ; host, s'il est fixé,
// l'emporte sur DOCKER_HOST pour viser le même démon que le moteur de conteneurs.
func commandEnv(getenv func(string) string, host string) []string {
	env := make([]string, 0, len(passedEnv))
	for _, name := range passedEnv {
		value := getenv(name)
		if name == "DOCKER_HOST" && host != "" {
			value = host
		}
		if value != "" {
			env = append(env, name+"="+value)
		}
	}
	return env
}

// ValidProjectName indique si name peut nommer une stack.
func ValidProjectName(name string) bool {
	return projectName.MatchString(name)
}

// Validate vérifie le fichier sans rien créer. vars (« NOM=valeur ») sont les
// valeurs d'une installation : un fichier qui écrit ${NOM:?} serait refusé
// sans elles.
func (e *Engine) Validate(ctx context.Context, content []byte, vars []string) error {
	env, err := e.withVars(vars)
	if err != nil {
		return err
	}
	_, err = e.config(ctx, content, env, "--quiet")
	return err
}

// withVars ajoute vars à l'environnement fixe de la commande. Le manifeste
// interdit déjà les noms réservés ; le moteur ne s'y fie pas : une valeur ne
// doit jamais changer le démon visé ni la commande elle-même. L'erreur nomme
// la variable, jamais sa valeur.
func (e *Engine) withVars(vars []string) ([]string, error) {
	if len(vars) == 0 {
		return e.env, nil
	}
	env := make([]string, 0, len(e.env)+len(vars))
	env = append(env, e.env...)
	for _, kv := range vars {
		name, _, ok := strings.Cut(kv, "=")
		if !ok || !variableName.MatchString(name) || reserved(name) {
			return nil, fmt.Errorf("%w : %q", ErrReservedVariable, name)
		}
		env = append(env, kv)
	}
	return env, nil
}

func reserved(name string) bool {
	if slices.Contains(passedEnv, name) {
		return true
	}
	for _, p := range reservedPrefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

// Resolve renvoie le modèle résolu du fichier, en JSON, sans rien créer.
// Compose y a appliqué ancres, `extends` et syntaxes courtes, mais les ${NOM}
// restent tels quels : la vérification des templates (ADR-0018) les croise
// avec les variables du manifeste.
func (e *Engine) Resolve(ctx context.Context, content []byte) ([]byte, error) {
	return e.config(ctx, content, e.env, "--no-interpolate", "--format", "json")
}

// config lance `config` sur le fichier, dans l'environnement env, avec args
// en plus, et renvoie sa sortie.
func (e *Engine) config(ctx context.Context, content []byte, env []string, args ...string) ([]byte, error) {
	if len(bytes.TrimSpace(content)) == 0 {
		return nil, fmt.Errorf("%w : fichier vide", ErrInvalidFile)
	}
	out, err := e.runner.Run(ctx, Command{
		Args:  append([]string{"-p", validationProject, "-f", "-", "config"}, args...),
		Stdin: content, Env: env, Dir: e.dir,
	})
	// Seul un refus de la commande dit que le fichier est en cause ; une
	// annulation ou une commande absente remontent telles quelles.
	var cmdErr *CommandError
	if errors.As(err, &cmdErr) {
		return nil, fmt.Errorf("%w : %w", ErrInvalidFile, err)
	}
	return out, err
}

// Deploy crée la stack ou la met à jour : seuls les services modifiés sont
// recréés, et ceux retirés du fichier sont supprimés. Il rend la main quand les
// services tournent et que ceux qui ont une sonde de santé sont sains. La
// progression de `up` est passée à out, si elle est fournie.
//
// Un service sans politique de redémarrage reçoit DefaultRestart, pour revenir
// seul après un redémarrage du serveur (#70) : sans cela, Docker applique `no`
// et l'application reste arrêtée.
//
// vars (« NOM=valeur ») sont les valeurs d'une installation (ADR-0019) :
// Compose les substitue aux ${NOM} du fichier, qui n'est jamais réécrit.
func (e *Engine) Deploy(ctx context.Context, project string, content []byte, vars []string, out OutputFunc) error {
	if err := checkProject(project); err != nil {
		return err
	}
	env, err := e.withVars(vars)
	if err != nil {
		return err
	}
	// La lecture du modèle valide aussi le fichier : elle sépare un fichier
	// refusé (ErrInvalidFile) d'un échec du démarrage, avant de toucher à la
	// stack en place.
	resolved, err := e.config(ctx, content, env, "--format", "json")
	if err != nil {
		return err
	}
	services, err := servicesWithoutRestart(resolved)
	if err != nil {
		return err
	}

	files := []string{"-f", "-"}
	if len(services) > 0 {
		override, err := e.restartOverride(services)
		if err != nil {
			return err
		}
		defer func() { _ = os.Remove(override) }()
		files = append(files, "-f", override)
		if out != nil {
			out(Stdout, restartNotice(services))
		}
	}

	// Sans --project-directory, Compose prendrait le dossier du fichier
	// complémentaire comme dossier du projet, puisque l'entrée standard n'en a
	// pas : un chemin relatif comme « ./donnees » changerait d'emplacement.
	// --quiet-pull garde les événements de téléchargement par service, sans
	// une ligne par couche et par pourcentage.
	args := append([]string{"-p", project, "--project-directory", e.dir}, files...)
	args = append(args, "up", "--detach", "--wait", "--wait-timeout", seconds(e.waitTimeout),
		"--remove-orphans", "--quiet-pull")
	_, err = e.runner.Run(ctx, Command{Args: args, Stdin: content, Env: env, Dir: e.dir, Output: out})
	return err
}

// Remove arrête et supprime les conteneurs et réseaux de la stack. Les volumes
// nommés ne partent qu'avec removeVolumes : ce sont les données de l'application.
func (e *Engine) Remove(ctx context.Context, project string, removeVolumes bool) error {
	if err := checkProject(project); err != nil {
		return err
	}
	args := []string{"-p", project, "down", "--remove-orphans", "--timeout", seconds(stopTimeout)}
	if removeVolumes {
		args = append(args, "--volumes")
	}
	_, err := e.run(ctx, nil, args...)
	return err
}

// Services liste les conteneurs de la stack, arrêtés compris. Une stack
// inconnue n'a aucun service.
func (e *Engine) Services(ctx context.Context, project string) ([]Service, error) {
	if err := checkProject(project); err != nil {
		return nil, err
	}
	out, err := e.run(ctx, nil, "-p", project, "ps", "--all", "--format", "json")
	if err != nil {
		return nil, err
	}
	return parseServices(out)
}

// Project est un projet Compose en place sur l'hôte et l'état de ses conteneurs.
type Project struct {
	Name       string
	Containers int // conteneurs du projet, arrêtés compris
	Running    int // conteneurs en marche
}

// Projects liste les projets Compose de l'hôte, arrêtés compris.
func (e *Engine) Projects(ctx context.Context) ([]Project, error) {
	out, err := e.run(ctx, nil, "ls", "--all", "--format", "json")
	if err != nil {
		return nil, err
	}
	return parseProjects(out)
}

// statusCount reconnaît un compte de `ls` : « running(2) », « exited(1) »…
var statusCount = regexp.MustCompile(`([a-z]+)\((\d+)\)`)

// parseProjects lit la sortie de `ls --format json`, par exemple
// [{"Name":"x","Status":"exited(1), running(1)","ConfigFiles":"-"}].
func parseProjects(out []byte) ([]Project, error) {
	out = bytes.TrimSpace(out)
	if len(out) == 0 {
		return nil, nil
	}
	var entries []struct {
		Name   string `json:"Name"`
		Status string `json:"Status"`
	}
	if err := json.Unmarshal(out, &entries); err != nil {
		return nil, fmt.Errorf("lecture des projets Compose : %w", err)
	}
	projects := make([]Project, 0, len(entries))
	for _, entry := range entries {
		p := Project{Name: entry.Name}
		for _, m := range statusCount.FindAllStringSubmatch(entry.Status, -1) {
			n, err := strconv.Atoi(m[2])
			if err != nil {
				return nil, fmt.Errorf("lecture des projets Compose : état %q", entry.Status)
			}
			p.Containers += n
			if m[1] == "running" {
				p.Running += n
			}
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// Ping vérifie que la commande docker compose répond.
func (e *Engine) Ping(ctx context.Context) error {
	_, err := e.run(ctx, nil, "version", "--short")
	return err
}

func (e *Engine) run(ctx context.Context, stdin []byte, args ...string) ([]byte, error) {
	return e.runner.Run(ctx, Command{Args: args, Stdin: stdin, Env: e.env, Dir: e.dir})
}

func checkProject(project string) error {
	if !ValidProjectName(project) {
		return fmt.Errorf("%w : %q", ErrInvalidProject, project)
	}
	return nil
}

func seconds(d time.Duration) string {
	return strconv.Itoa(int(d.Round(time.Second) / time.Second))
}

// psEntry reprend les champs utiles de `docker compose ps --format json`.
type psEntry struct {
	ID      string `json:"ID"`
	Name    string `json:"Name"`
	Service string `json:"Service"`
	Image   string `json:"Image"`
	State   string `json:"State"`
	Health  string `json:"Health"`
}

// parseServices lit la sortie de `ps --format json` : un objet par ligne depuis
// Compose 2.21, un tableau avant.
func parseServices(out []byte) ([]Service, error) {
	out = bytes.TrimSpace(out)
	if len(out) == 0 {
		return nil, nil
	}

	var entries []psEntry
	if out[0] == '[' {
		if err := json.Unmarshal(out, &entries); err != nil {
			return nil, fmt.Errorf("lecture de l'état de la stack : %w", err)
		}
	} else {
		scanner := bufio.NewScanner(bytes.NewReader(out))
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := bytes.TrimSpace(scanner.Bytes())
			if len(line) == 0 {
				continue
			}
			var entry psEntry
			if err := json.Unmarshal(line, &entry); err != nil {
				return nil, fmt.Errorf("lecture de l'état de la stack : %w", err)
			}
			entries = append(entries, entry)
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("lecture de l'état de la stack : %w", err)
		}
	}

	services := make([]Service, 0, len(entries))
	for _, entry := range entries {
		services = append(services, Service{
			ID:      entry.ID,
			Name:    entry.Name,
			Service: entry.Service,
			Image:   entry.Image,
			State:   orchestrator.State(strings.ToLower(entry.State)),
			Health:  entry.Health,
		})
	}
	return services, nil
}
