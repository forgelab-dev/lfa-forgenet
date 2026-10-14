// Package orchestrator définit le contrat du moteur de conteneurs de ForgeNet.
//
// Le domaine ne dépend ni du SDK Docker ni du contrat protobuf : l'implémentation
// Docker vit dans orchestrator/docker, la conversion vers l'API dans
// transport/grpc. Compose (#7) et les stratégies de déploiement s'appuient sur
// Engine, jamais sur un moteur précis.
package orchestrator

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrNotFound indique qu'aucun conteneur ne correspond à l'identifiant demandé.
	ErrNotFound = errors.New("conteneur introuvable")
	// ErrInvalidArgument indique une demande incomplète ou incohérente.
	ErrInvalidArgument = errors.New("paramètre invalide")
)

// State est l'état d'un conteneur, indépendant du moteur.
type State string

// États possibles d'un conteneur.
const (
	StateUnknown    State = ""
	StateCreated    State = "created"
	StateRunning    State = "running"
	StatePaused     State = "paused"
	StateRestarting State = "restarting"
	StateRemoving   State = "removing"
	StateExited     State = "exited"
	StateDead       State = "dead"
)

// Port associe un port du conteneur à un port de l'hôte. HostPort vaut 0
// quand le port est seulement exposé, sans publication sur l'hôte.
type Port struct {
	HostIP        string
	HostPort      uint16
	ContainerPort uint16
	Protocol      string // tcp, udp ou sctp
}

// Container décrit un conteneur tel que le moteur le rapporte.
type Container struct {
	ID        string
	Name      string
	Image     string
	State     State
	Ports     []Port
	CreatedAt time.Time
}

// CreateSpec décrit le conteneur à créer.
type CreateSpec struct {
	Name        string
	Image       string
	Env         []string
	Ports       []Port // HostIP vide : adresse de publication configurée
	Labels      map[string]string
	MemoryBytes int64 // 0 : pas de limite
	MilliCPU    int64 // 1000 = un cœur ; 0 : pas de limite
}

// LogStream identifie la sortie d'origine d'un morceau de journal.
type LogStream int

// Sorties d'un conteneur.
const (
	Stdout LogStream = iota + 1
	Stderr
)

// LogOptions filtre les journaux demandés.
type LogOptions struct {
	Follow bool
	Tail   int       // 0 : tout l'historique
	Since  time.Time // zéro : depuis le début
	Stdout bool
	Stderr bool
}

// LogChunk est un morceau de journal, déjà séparé par sortie.
type LogChunk struct {
	Stream LogStream
	Data   []byte
}

// Engine gère le cycle de vie des conteneurs.
type Engine interface {
	List(ctx context.Context, includeStopped bool) ([]Container, error)
	Get(ctx context.Context, id string) (Container, error)
	Create(ctx context.Context, spec CreateSpec) (Container, error)
	Start(ctx context.Context, id string) (Container, error)
	// Stop laisse timeout au conteneur pour s'arrêter avant SIGKILL ; 0 garde
	// le délai configuré sur le conteneur ou celui du moteur.
	Stop(ctx context.Context, id string, timeout time.Duration) (Container, error)
	Restart(ctx context.Context, id string) (Container, error)
	Remove(ctx context.Context, id string, force bool) error
	// Logs appelle emit pour chaque morceau de journal, dans l'ordre, jusqu'à
	// la fin du journal (ou l'annulation de ctx avec Follow). Une erreur de
	// emit interrompt la lecture et est renvoyée telle quelle.
	Logs(ctx context.Context, id string, opts LogOptions, emit func(LogChunk) error) error
	Ping(ctx context.Context) error
}
