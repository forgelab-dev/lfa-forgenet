package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/docker/go-connections/nat"
)

// Service implémente l'interface Orchestrator en s'appuyant sur Docker Engine.
type Service struct {
	cli EngineClient
}

// Assure que Service implémente bien l'interface Orchestrator.
var _ Orchestrator = (*Service)(nil)

// New instancie un orchestrateur Docker à partir des paramètres d'hôte et de version d'API.
// Si dockerHost est vide, la configuration standard de l'environnement est utilisée.
func New(dockerHost, apiVersion string) (*Service, error) {
	opts := []client.Opt{
		client.FromEnv,
	}
	if dockerHost != "" {
		opts = append(opts, client.WithHost(dockerHost))
	}
	if apiVersion != "" {
		opts = append(opts, client.WithVersion(apiVersion))
	} else {
		opts = append(opts, client.WithAPIVersionNegotiation())
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, fmt.Errorf("initialisation client docker: %w", err)
	}

	return NewWithClient(cli), nil
}

// NewWithClient instancie un orchestrateur Docker autour d'un EngineClient fourni (pratique pour les mocks et tests).
func NewWithClient(cli EngineClient) *Service {
	return &Service{
		cli: cli,
	}
}

// Ping vérifie la communication avec le démon Docker.
func (s *Service) Ping(ctx context.Context) error {
	_, err := s.cli.Ping(ctx)
	if err != nil {
		return fmt.Errorf("démon docker injoignable: %w", err)
	}
	return nil
}

// Close libère les ressources et connexions sous-jacentes du client Docker.
func (s *Service) Close() error {
	if s.cli != nil {
		return s.cli.Close()
	}
	return nil
}

// ListContainers retourne l'ensemble des conteneurs selon le filtre d'arrêt.
func (s *Service) ListContainers(ctx context.Context, includeStopped bool) ([]*ContainerInfo, error) {
	opts := container.ListOptions{
		All: includeStopped,
	}

	containers, err := s.cli.ContainerList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("liste des conteneurs: %w", err)
	}

	results := make([]*ContainerInfo, 0, len(containers))
	for _, c := range containers {
		results = append(results, summaryToContainerInfo(c))
	}

	return results, nil
}

// GetContainer inspecte un conteneur par son identifiant ou son nom.
func (s *Service) GetContainer(ctx context.Context, containerID string) (*ContainerInfo, error) {
	trimmedID := strings.TrimSpace(containerID)
	if trimmedID == "" {
		return nil, errors.New("l'identifiant du conteneur ne peut pas être vide")
	}

	cJSON, err := s.cli.ContainerInspect(ctx, trimmedID)
	if err != nil {
		if isNotFound(err) {
			return nil, ErrContainerNotFound
		}
		return nil, fmt.Errorf("inspection du conteneur %q: %w", trimmedID, err)
	}

	return inspectToContainerInfo(cJSON), nil
}

// CreateContainer instancie un nouveau conteneur.
func (s *Service) CreateContainer(ctx context.Context, req CreateContainerRequest) (*ContainerInfo, error) {
	imageName := strings.TrimSpace(req.Image)
	if imageName == "" {
		return nil, errors.New("le nom de l'image est obligatoire")
	}

	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}

	for _, p := range req.Ports {
		proto := strings.ToLower(p.Protocol)
		if proto == "" {
			proto = "tcp"
		}
		natPort, err := nat.NewPort(proto, strconv.FormatUint(uint64(p.ContainerPort), 10))
		if err != nil {
			return nil, fmt.Errorf("port invalide %d/%s: %w", p.ContainerPort, proto, err)
		}

		exposedPorts[natPort] = struct{}{}
		if p.HostPort > 0 {
			portBindings[natPort] = []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.FormatUint(uint64(p.HostPort), 10),
				},
			}
		}
	}

	cfg := &container.Config{
		Image:        imageName,
		Env:          req.Env,
		Labels:       req.Labels,
		ExposedPorts: exposedPorts,
	}

	hostCfg := &container.HostConfig{
		PortBindings: portBindings,
	}

	if req.MemoryLimit > 0 {
		hostCfg.Memory = req.MemoryLimit
	}
	if req.CPULimit > 0 {
		// Convertit milli-CPU en NanoCPUs (1 CPU = 1000m = 1e9 nanoCPUs)
		hostCfg.NanoCPUs = req.CPULimit * 1_000_000
	}

	created, err := s.cli.ContainerCreate(ctx, cfg, hostCfg, nil, nil, req.Name)
	if err != nil {
		return nil, fmt.Errorf("création du conteneur: %w", err)
	}

	return s.GetContainer(ctx, created.ID)
}

// StartContainer démarre un conteneur arrêté.
func (s *Service) StartContainer(ctx context.Context, containerID string) (*ContainerInfo, error) {
	trimmedID := strings.TrimSpace(containerID)
	if trimmedID == "" {
		return nil, errors.New("l'identifiant du conteneur ne peut pas être vide")
	}

	err := s.cli.ContainerStart(ctx, trimmedID, container.StartOptions{})
	if err != nil {
		if isNotFound(err) {
			return nil, ErrContainerNotFound
		}
		return nil, fmt.Errorf("démarrage du conteneur %q: %w", trimmedID, err)
	}

	return s.GetContainer(ctx, trimmedID)
}

// StopContainer arrête proprement un conteneur avec un délai de grâce.
func (s *Service) StopContainer(ctx context.Context, containerID string, timeoutSeconds uint32) (*ContainerInfo, error) {
	trimmedID := strings.TrimSpace(containerID)
	if trimmedID == "" {
		return nil, errors.New("l'identifiant du conteneur ne peut pas être vide")
	}

	timeoutVal := int(timeoutSeconds)
	var stopOpts container.StopOptions
	if timeoutSeconds > 0 {
		stopOpts.Timeout = &timeoutVal
	}

	err := s.cli.ContainerStop(ctx, trimmedID, stopOpts)
	if err != nil {
		if isNotFound(err) {
			return nil, ErrContainerNotFound
		}
		return nil, fmt.Errorf("arrêt du conteneur %q: %w", trimmedID, err)
	}

	return s.GetContainer(ctx, trimmedID)
}

// RestartContainer redémarre un conteneur.
func (s *Service) RestartContainer(ctx context.Context, containerID string) (*ContainerInfo, error) {
	trimmedID := strings.TrimSpace(containerID)
	if trimmedID == "" {
		return nil, errors.New("l'identifiant du conteneur ne peut pas être vide")
	}

	err := s.cli.ContainerRestart(ctx, trimmedID, container.StopOptions{})
	if err != nil {
		if isNotFound(err) {
			return nil, ErrContainerNotFound
		}
		return nil, fmt.Errorf("redémarrage du conteneur %q: %w", trimmedID, err)
	}

	return s.GetContainer(ctx, trimmedID)
}

// DeleteContainer supprime un conteneur.
func (s *Service) DeleteContainer(ctx context.Context, containerID string, force bool) error {
	trimmedID := strings.TrimSpace(containerID)
	if trimmedID == "" {
		return errors.New("l'identifiant du conteneur ne peut pas être vide")
	}

	err := s.cli.ContainerRemove(ctx, trimmedID, container.RemoveOptions{
		Force: force,
	})
	if err != nil {
		if isNotFound(err) {
			return ErrContainerNotFound
		}
		return fmt.Errorf("suppression du conteneur %q: %w", trimmedID, err)
	}

	return nil
}

// StreamLogs ouvre un flux de lecture des logs d'un conteneur.
func (s *Service) StreamLogs(
	ctx context.Context,
	containerID string,
	follow bool,
	tail int,
	since time.Time,
	stdout, stderr bool,
) (io.ReadCloser, error) {
	trimmedID := strings.TrimSpace(containerID)
	if trimmedID == "" {
		return nil, errors.New("l'identifiant du conteneur ne peut pas être vide")
	}

	tailStr := "all"
	if tail > 0 {
		tailStr = strconv.Itoa(tail)
	}

	sinceStr := ""
	if !since.IsZero() {
		sinceStr = since.Format(time.RFC3339)
	}

	opts := container.LogsOptions{
		ShowStdout: stdout,
		ShowStderr: stderr,
		Follow:     follow,
		Tail:       tailStr,
		Since:      sinceStr,
	}

	reader, err := s.cli.ContainerLogs(ctx, trimmedID, opts)
	if err != nil {
		if isNotFound(err) {
			return nil, ErrContainerNotFound
		}
		return nil, fmt.Errorf("lecture des logs du conteneur %q: %w", trimmedID, err)
	}

	return reader, nil
}

// summaryToContainerInfo convertit la structure condensée Docker en ContainerInfo du domaine.
func summaryToContainerInfo(c types.Container) *ContainerInfo {
	name := ""
	if len(c.Names) > 0 {
		name = strings.TrimPrefix(c.Names[0], "/")
	}

	ports := make([]PortMapping, 0, len(c.Ports))
	for _, p := range c.Ports {
		ports = append(ports, PortMapping{
			HostPort:      uint32(p.PublicPort),
			ContainerPort: uint32(p.PrivatePort),
			Protocol:      p.Type,
		})
	}

	return &ContainerInfo{
		ID:        c.ID,
		Name:      name,
		Image:     c.Image,
		State:     ParseContainerState(c.State),
		Ports:     ports,
		CreatedAt: time.Unix(c.Created, 0).UTC(),
	}
}

// inspectToContainerInfo convertit le rapport d'inspection Docker complet en ContainerInfo du domaine.
func inspectToContainerInfo(c types.ContainerJSON) *ContainerInfo {
	name := strings.TrimPrefix(c.Name, "/")

	var ports []PortMapping
	if c.NetworkSettings != nil && c.NetworkSettings.Ports != nil {
		for portProto, bindings := range c.NetworkSettings.Ports {
			cPort, proto := parsePortProto(string(portProto))
			if len(bindings) > 0 {
				for _, b := range bindings {
					hPort, _ := strconv.ParseUint(b.HostPort, 10, 32)
					ports = append(ports, PortMapping{
						HostPort:      uint32(hPort),
						ContainerPort: cPort,
						Protocol:      proto,
					})
				}
			} else {
				ports = append(ports, PortMapping{
					HostPort:      0,
					ContainerPort: cPort,
					Protocol:      proto,
				})
			}
		}
	}

	createdAt, err := time.Parse(time.RFC3339Nano, c.Created)
	if err != nil {
		createdAt = time.Now().UTC()
	}

	state := ParseContainerState("")
	if c.State != nil {
		state = ParseContainerState(c.State.Status)
	}

	image := c.Config.Image
	if image == "" {
		image = c.Image
	}

	return &ContainerInfo{
		ID:        c.ID,
		Name:      name,
		Image:     image,
		State:     state,
		Ports:     ports,
		CreatedAt: createdAt.UTC(),
	}
}

func parsePortProto(s string) (uint32, string) {
	parts := strings.Split(s, "/")
	p, _ := strconv.ParseUint(parts[0], 10, 32)
	proto := "tcp"
	if len(parts) > 1 && parts[1] != "" {
		proto = parts[1]
	}
	return uint32(p), proto
}

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errdefs.IsNotFound(err) || client.IsErrNotFound(err) || strings.Contains(strings.ToLower(err.Error()), "not found") || strings.Contains(strings.ToLower(err.Error()), "no such container")
}
