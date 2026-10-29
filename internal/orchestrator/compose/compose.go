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

const (
	// defaultWaitTimeout borne l'attente de services démarrés et sains.
	defaultWaitTimeout = 2 * time.Minute
	// stopTimeout laisse aux conteneurs le temps de s'arrêter avant SIGKILL.
	stopTimeout = 10 * time.Second
	// validationProject nomme le projet le temps d'une validation : `config` ne crée rien.
	validationProject = "forgenet-validation"
)

// Command est une exécution de `docker compose`.
type Command struct {
	Args  []string // arguments après « docker compose »
	Stdin []byte
	Env   []string // environnement complet de la commande
	Dir   string
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

	e := &Engine{
		runner:      execRunner{},
		env:         commandEnv(os.Getenv, cfg.Host),
		dir:         dir,
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

// Validate vérifie le fichier sans rien créer.
func (e *Engine) Validate(ctx context.Context, content []byte) error {
	if len(bytes.TrimSpace(content)) == 0 {
		return fmt.Errorf("%w : fichier vide", ErrInvalidFile)
	}
	_, err := e.run(ctx, content, "-p", validationProject, "-f", "-", "config", "--quiet")
	// Seul un refus de la commande dit que le fichier est en cause ; une
	// annulation ou une commande absente remontent telles quelles.
	var cmdErr *CommandError
	if errors.As(err, &cmdErr) {
		return fmt.Errorf("%w : %w", ErrInvalidFile, err)
	}
	return err
}

// Deploy crée la stack ou la met à jour : seuls les services modifiés sont
// recréés, et ceux retirés du fichier sont supprimés. Il rend la main quand les
// services tournent et que ceux qui ont une sonde de santé sont sains.
func (e *Engine) Deploy(ctx context.Context, project string, content []byte) error {
	if err := checkProject(project); err != nil {
		return err
	}
	// La validation sépare un fichier refusé (ErrInvalidFile) d'un échec du
	// démarrage, avant de toucher à la stack en place.
	if err := e.Validate(ctx, content); err != nil {
		return err
	}
	_, err := e.run(ctx, content, "-p", project, "-f", "-", "up",
		"--detach", "--wait", "--wait-timeout", seconds(e.waitTimeout),
		"--remove-orphans", "--quiet-pull")
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
