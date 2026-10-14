// Package docker implémente le cycle de vie et l'orchestration des conteneurs via Docker Engine.
package docker

import (
	"errors"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

var (
	// ErrContainerNotFound indique qu'aucun conteneur ne correspond à l'identifiant demandé.
	ErrContainerNotFound = errors.New("conteneur introuvable")
)

// PortMapping représente la redirection d'un port hôte vers un conteneur.
type PortMapping struct {
	HostPort      uint32 `json:"host_port"`
	ContainerPort uint32 `json:"container_port"`
	Protocol      string `json:"protocol"` // "tcp", "udp"
}

// ContainerInfo regroupe l'état d'un conteneur interrogé depuis le moteur Docker.
type ContainerInfo struct {
	ID        string                    `json:"id"`
	Name      string                    `json:"name"`
	Image     string                    `json:"image"`
	State     forgenetv1.ContainerState `json:"state"`
	Ports     []PortMapping             `json:"ports"`
	CreatedAt time.Time                 `json:"created_at"`
}

// CreateContainerRequest contient les paramètres nécessaires à l'instanciation d'un conteneur.
type CreateContainerRequest struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Env         []string          `json:"env,omitempty"`
	Ports       []PortMapping     `json:"ports,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	MemoryLimit int64             `json:"memory_limit,omitempty"` // en octets
	CPULimit    int64             `json:"cpu_limit,omitempty"`    // en milli-CPU (ex: 500 = 0.5 CPU)
}

// ToProto convertit un ContainerInfo de domaine vers le message Protobuf défini dans le contrat api/.
func (c *ContainerInfo) ToProto() *forgenetv1.Container {
	if c == nil {
		return nil
	}

	protoPorts := make([]*forgenetv1.PortMapping, len(c.Ports))
	for i, p := range c.Ports {
		protoPorts[i] = &forgenetv1.PortMapping{
			HostPort:      p.HostPort,
			ContainerPort: p.ContainerPort,
			Protocol:      p.Protocol,
		}
	}

	return &forgenetv1.Container{
		Id:        c.ID,
		Name:      c.Name,
		Image:     c.Image,
		State:     c.State,
		Ports:     protoPorts,
		CreatedAt: timestamppb.New(c.CreatedAt),
	}
}

// ParseContainerState traduit l'état textuel renvoyé par Docker Engine en ContainerState typé.
func ParseContainerState(stateStr string) forgenetv1.ContainerState {
	switch strings.ToLower(strings.TrimSpace(stateStr)) {
	case "created":
		return forgenetv1.ContainerState_CONTAINER_STATE_CREATED
	case "running":
		return forgenetv1.ContainerState_CONTAINER_STATE_RUNNING
	case "paused":
		return forgenetv1.ContainerState_CONTAINER_STATE_PAUSED
	case "restarting":
		return forgenetv1.ContainerState_CONTAINER_STATE_RESTARTING
	case "exited":
		return forgenetv1.ContainerState_CONTAINER_STATE_EXITED
	case "dead":
		return forgenetv1.ContainerState_CONTAINER_STATE_DEAD
	default:
		return forgenetv1.ContainerState_CONTAINER_STATE_UNSPECIFIED
	}
}
