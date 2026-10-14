package docker

import (
	"context"
	"io"
	"time"
)

// Orchestrator définit le contrat de gestion du cycle de vie des conteneurs pour ForgeNet.
type Orchestrator interface {
	// ListContainers retourne la liste des conteneurs (actifs ou tous).
	ListContainers(ctx context.Context, includeStopped bool) ([]*ContainerInfo, error)
	// GetContainer inspecte un conteneur par son identifiant ou son nom.
	GetContainer(ctx context.Context, containerID string) (*ContainerInfo, error)
	// CreateContainer instancie un nouveau conteneur selon la spécification fournie.
	CreateContainer(ctx context.Context, req CreateContainerRequest) (*ContainerInfo, error)
	// StartContainer démarre un conteneur arrêté.
	StartContainer(ctx context.Context, containerID string) (*ContainerInfo, error)
	// StopContainer arrête proprement un conteneur avec un délai d'attente avant SIGKILL.
	StopContainer(ctx context.Context, containerID string, timeoutSeconds uint32) (*ContainerInfo, error)
	// RestartContainer redémarre un conteneur en cours d'exécution.
	RestartContainer(ctx context.Context, containerID string) (*ContainerInfo, error)
	// DeleteContainer supprime un conteneur.
	DeleteContainer(ctx context.Context, containerID string, force bool) error
	// StreamLogs ouvre un flux de lecture des logs stdout/stderr d'un conteneur.
	StreamLogs(ctx context.Context, containerID string, follow bool, tail int, since time.Time, stdout, stderr bool) (io.ReadCloser, error)
	// Ping vérifie la disponibilité du démon Docker.
	Ping(ctx context.Context) error
}
