// Package docker implémente orchestrator.Engine avec Docker Engine.
package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/netip"
	"strconv"
	"strings"
	"time"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
)

// Engine pilote les conteneurs d'un démon Docker.
type Engine struct {
	cli         EngineClient
	publishHost netip.Addr
}

var _ orchestrator.Engine = (*Engine)(nil)

// New se connecte au démon décrit par la configuration (variables DOCKER_* par
// défaut). La version d'API est négociée au premier appel, sauf si elle est fixée.
func New(cfg config.DockerConfig) (*Engine, error) {
	opts := []client.Opt{client.FromEnv}
	if cfg.Host != "" {
		opts = append(opts, client.WithHost(cfg.Host))
	}
	if cfg.APIVersion != "" {
		opts = append(opts, client.WithAPIVersion(cfg.APIVersion))
	}
	cli, err := client.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("client Docker : %w", err)
	}
	return NewWithClient(cli, cfg.PublishHost)
}

// NewWithClient construit l'orchestrateur autour d'un client fourni. publishHost
// est l'adresse de l'hôte sur laquelle les ports des conteneurs sont publiés.
func NewWithClient(cli EngineClient, publishHost string) (*Engine, error) {
	addr, err := netip.ParseAddr(publishHost)
	if err != nil {
		return nil, fmt.Errorf("adresse de publication des ports %q : %w", publishHost, err)
	}
	return &Engine{cli: cli, publishHost: addr}, nil
}

// Close libère la connexion au démon.
func (e *Engine) Close() error {
	return e.cli.Close()
}

// Ping vérifie que le démon répond.
func (e *Engine) Ping(ctx context.Context) error {
	if _, err := e.cli.Ping(ctx, client.PingOptions{}); err != nil {
		return fmt.Errorf("démon Docker injoignable : %w", err)
	}
	return nil
}

// List renvoie les conteneurs en cours, ou tous avec includeStopped.
func (e *Engine) List(ctx context.Context, includeStopped bool) ([]orchestrator.Container, error) {
	res, err := e.cli.ContainerList(ctx, client.ContainerListOptions{All: includeStopped})
	if err != nil {
		return nil, fmt.Errorf("liste des conteneurs : %w", err)
	}
	containers := make([]orchestrator.Container, 0, len(res.Items))
	for _, c := range res.Items {
		containers = append(containers, fromSummary(c))
	}
	return containers, nil
}

// Get renvoie l'état d'un conteneur (identifiant ou nom).
func (e *Engine) Get(ctx context.Context, id string) (orchestrator.Container, error) {
	res, err := e.inspect(ctx, id)
	if err != nil {
		return orchestrator.Container{}, err
	}
	return fromInspect(res), nil
}

// Create crée le conteneur sans le démarrer et renvoie son état.
func (e *Engine) Create(ctx context.Context, spec orchestrator.CreateSpec) (orchestrator.Container, error) {
	if strings.TrimSpace(spec.Image) == "" {
		return orchestrator.Container{}, fmt.Errorf("%w : l'image est obligatoire", orchestrator.ErrInvalidArgument)
	}
	exposed, bindings, err := e.portConfig(spec.Ports)
	if err != nil {
		return orchestrator.Container{}, err
	}

	created, err := e.cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name: spec.Name,
		Config: &container.Config{
			Image:        spec.Image,
			Env:          spec.Env,
			Labels:       spec.Labels,
			ExposedPorts: exposed,
		},
		HostConfig: &container.HostConfig{
			PortBindings: bindings,
			Resources: container.Resources{
				Memory: spec.MemoryBytes,
				// 1 cœur = 1000 milli-CPU = 1e9 nano-CPU.
				NanoCPUs: spec.MilliCPU * 1_000_000,
			},
		},
	})
	if err != nil {
		return orchestrator.Container{}, engineError("création du conteneur", err)
	}
	return e.Get(ctx, created.ID)
}

// Start démarre un conteneur et renvoie son nouvel état.
func (e *Engine) Start(ctx context.Context, id string) (orchestrator.Container, error) {
	return e.act(ctx, id, "démarrage", func(id string) error {
		_, err := e.cli.ContainerStart(ctx, id, client.ContainerStartOptions{})
		return err
	})
}

// Stop arrête un conteneur (SIGTERM, puis SIGKILL après timeout) et renvoie son nouvel état.
func (e *Engine) Stop(ctx context.Context, id string, timeout time.Duration) (orchestrator.Container, error) {
	opts := client.ContainerStopOptions{}
	if timeout > 0 {
		seconds := int(timeout.Round(time.Second) / time.Second)
		opts.Timeout = &seconds
	}
	return e.act(ctx, id, "arrêt", func(id string) error {
		_, err := e.cli.ContainerStop(ctx, id, opts)
		return err
	})
}

// Restart redémarre un conteneur et renvoie son nouvel état.
func (e *Engine) Restart(ctx context.Context, id string) (orchestrator.Container, error) {
	return e.act(ctx, id, "redémarrage", func(id string) error {
		_, err := e.cli.ContainerRestart(ctx, id, client.ContainerRestartOptions{})
		return err
	})
}

// Remove supprime un conteneur ; force arrête d'abord un conteneur en cours.
func (e *Engine) Remove(ctx context.Context, id string, force bool) error {
	id, err := requireID(id)
	if err != nil {
		return err
	}
	if _, err := e.cli.ContainerRemove(ctx, id, client.ContainerRemoveOptions{Force: force}); err != nil {
		return engineError("suppression du conteneur "+id, err)
	}
	return nil
}

// Logs sépare les sorties du conteneur. Sans TTY, Docker multiplexe stdout et
// stderr dans un même flux, décodé par stdcopy ; avec TTY, le flux est brut et
// tout provient de stdout. Le réglage est lu sur le conteneur, pas deviné à
// partir des octets reçus.
func (e *Engine) Logs(ctx context.Context, id string, opts orchestrator.LogOptions, emit func(orchestrator.LogChunk) error) error {
	info, err := e.inspect(ctx, id)
	if err != nil {
		return err
	}

	logOpts := client.ContainerLogsOptions{
		ShowStdout: opts.Stdout,
		ShowStderr: opts.Stderr,
		Follow:     opts.Follow,
		Tail:       "all",
	}
	if opts.Tail > 0 {
		logOpts.Tail = strconv.Itoa(opts.Tail)
	}
	if !opts.Since.IsZero() {
		logOpts.Since = opts.Since.UTC().Format(time.RFC3339Nano)
	}

	rc, err := e.cli.ContainerLogs(ctx, info.Container.ID, logOpts)
	if err != nil {
		return engineError("lecture des journaux", err)
	}
	defer func() { _ = rc.Close() }()

	stdout := chunkWriter{stream: orchestrator.Stdout, emit: emit}
	if info.Container.Config != nil && info.Container.Config.Tty {
		_, err = io.Copy(stdout, rc)
	} else {
		_, err = stdcopy.StdCopy(stdout, chunkWriter{stream: orchestrator.Stderr, emit: emit}, rc)
	}
	// Avec Follow, l'annulation du contexte est la fin normale de la lecture.
	if err != nil && ctx.Err() == nil {
		return err
	}
	return nil
}

// chunkWriter transmet chaque écriture à emit. Les tampons de stdcopy et
// d'io.Copy sont réutilisés : les données sont copiées avant d'être émises.
type chunkWriter struct {
	stream orchestrator.LogStream
	emit   func(orchestrator.LogChunk) error
}

func (w chunkWriter) Write(p []byte) (int, error) {
	if err := w.emit(orchestrator.LogChunk{Stream: w.stream, Data: append([]byte(nil), p...)}); err != nil {
		return 0, err
	}
	return len(p), nil
}

// act exécute une action sur un conteneur puis renvoie son nouvel état.
func (e *Engine) act(ctx context.Context, id, action string, do func(id string) error) (orchestrator.Container, error) {
	id, err := requireID(id)
	if err != nil {
		return orchestrator.Container{}, err
	}
	if err := do(id); err != nil {
		return orchestrator.Container{}, engineError(action+" du conteneur "+id, err)
	}
	return e.Get(ctx, id)
}

func (e *Engine) inspect(ctx context.Context, id string) (client.ContainerInspectResult, error) {
	id, err := requireID(id)
	if err != nil {
		return client.ContainerInspectResult{}, err
	}
	res, err := e.cli.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
	if err != nil {
		return client.ContainerInspectResult{}, engineError("inspection du conteneur "+id, err)
	}
	return res, nil
}

// portConfig traduit les ports demandés. Un port publié sans adresse l'est sur
// l'adresse configurée (127.0.0.1 par défaut) : l'exposition publique passe par
// le reverse proxy, pas par un port ouvert sur toutes les interfaces.
func (e *Engine) portConfig(ports []orchestrator.Port) (network.PortSet, network.PortMap, error) {
	exposed := network.PortSet{}
	bindings := network.PortMap{}
	for _, p := range ports {
		proto := network.IPProtocol(strings.ToLower(p.Protocol))
		if proto == "" {
			proto = network.TCP
		}
		// network.PortFrom n'écarte qu'un protocole vide : la liste est vérifiée ici.
		if !validProtocols[proto] || p.ContainerPort == 0 {
			return nil, nil, fmt.Errorf("%w : port %d/%s", orchestrator.ErrInvalidArgument, p.ContainerPort, proto)
		}
		port, _ := network.PortFrom(p.ContainerPort, proto)
		exposed[port] = struct{}{}
		if p.HostPort == 0 {
			continue
		}

		hostIP := e.publishHost
		if p.HostIP != "" {
			addr, err := netip.ParseAddr(p.HostIP)
			if err != nil {
				return nil, nil, fmt.Errorf("%w : adresse %q", orchestrator.ErrInvalidArgument, p.HostIP)
			}
			hostIP = addr
		}
		bindings[port] = append(bindings[port], network.PortBinding{HostIP: hostIP, HostPort: strconv.Itoa(int(p.HostPort))})
	}
	return exposed, bindings, nil
}

var validProtocols = map[network.IPProtocol]bool{network.TCP: true, network.UDP: true, network.SCTP: true}

func requireID(id string) (string, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return "", fmt.Errorf("%w : l'identifiant du conteneur est obligatoire", orchestrator.ErrInvalidArgument)
	}
	return id, nil
}

// engineError traduit les erreurs typées du démon en erreurs du domaine.
func engineError(action string, err error) error {
	switch {
	case cerrdefs.IsNotFound(err):
		return fmt.Errorf("%s : %w", action, errors.Join(orchestrator.ErrNotFound, err))
	case cerrdefs.IsInvalidArgument(err):
		return fmt.Errorf("%s : %w", action, errors.Join(orchestrator.ErrInvalidArgument, err))
	default:
		return fmt.Errorf("%s : %w", action, err)
	}
}
