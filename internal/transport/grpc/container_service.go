package grpc

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

const (
	defaultPageSize = 50
	maxPageSize     = 100
)

// ContainerServiceServer expose orchestrator.Engine selon le contrat
// forgenet.v1.ContainerService. L'authentification et les rôles sont
// appliqués en amont par auth.Guard ; ce service traduit et journalise.
type ContainerServiceServer struct {
	forgenetv1.UnimplementedContainerServiceServer
	engine  orchestrator.Engine
	auditor audit.Auditor
	logger  zerolog.Logger
}

var _ forgenetv1.ContainerServiceServer = (*ContainerServiceServer)(nil)

// NewContainerServiceServer instancie le service gRPC des conteneurs.
func NewContainerServiceServer(engine orchestrator.Engine, auditor audit.Auditor, logger zerolog.Logger) *ContainerServiceServer {
	return &ContainerServiceServer{
		engine:  engine,
		auditor: auditor,
		logger:  logger.With().Str("component", "grpc.container_service").Logger(),
	}
}

// ListContainers renvoie une page de conteneurs. Le jeton de page est la
// position du premier élément ; un jeton illisible est refusé.
func (s *ContainerServiceServer) ListContainers(ctx context.Context, req *forgenetv1.ListContainersRequest) (*forgenetv1.ListContainersResponse, error) {
	pageSize := int(req.GetPageSize())
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	pageSize = min(pageSize, maxPageSize)

	offset := 0
	if token := req.GetPageToken(); token != "" {
		n, err := strconv.Atoi(token)
		if err != nil || n < 0 {
			return nil, status.Error(codes.InvalidArgument, "jeton de page invalide")
		}
		offset = n
	}

	containers, err := s.engine.List(ctx, req.GetIncludeStopped())
	if err != nil {
		return nil, s.toStatus(err, "liste des conteneurs")
	}

	offset = min(offset, len(containers))
	end := min(offset+pageSize, len(containers))
	resp := &forgenetv1.ListContainersResponse{Containers: make([]*forgenetv1.Container, 0, end-offset)}
	for _, c := range containers[offset:end] {
		resp.Containers = append(resp.Containers, toProtoContainer(c))
	}
	if end < len(containers) {
		resp.NextPageToken = strconv.Itoa(end)
	}
	return resp, nil
}

// GetContainer renvoie l'état d'un conteneur.
func (s *ContainerServiceServer) GetContainer(ctx context.Context, req *forgenetv1.GetContainerRequest) (*forgenetv1.GetContainerResponse, error) {
	c, err := s.engine.Get(ctx, req.GetContainerId())
	if err != nil {
		return nil, s.toStatus(err, "inspection du conteneur")
	}
	return &forgenetv1.GetContainerResponse{Container: toProtoContainer(c)}, nil
}

// StartContainer démarre un conteneur (admin).
func (s *ContainerServiceServer) StartContainer(ctx context.Context, req *forgenetv1.StartContainerRequest) (*forgenetv1.StartContainerResponse, error) {
	c, err := s.mutate(ctx, "container.start", req.GetContainerId(), s.engine.Start)
	if err != nil {
		return nil, err
	}
	return &forgenetv1.StartContainerResponse{Container: c}, nil
}

// StopContainer arrête un conteneur dans le délai demandé (admin).
func (s *ContainerServiceServer) StopContainer(ctx context.Context, req *forgenetv1.StopContainerRequest) (*forgenetv1.StopContainerResponse, error) {
	timeout := time.Duration(req.GetTimeoutSeconds()) * time.Second
	c, err := s.mutate(ctx, "container.stop", req.GetContainerId(), func(ctx context.Context, id string) (orchestrator.Container, error) {
		return s.engine.Stop(ctx, id, timeout)
	})
	if err != nil {
		return nil, err
	}
	return &forgenetv1.StopContainerResponse{Container: c}, nil
}

// RestartContainer redémarre un conteneur (admin).
func (s *ContainerServiceServer) RestartContainer(ctx context.Context, req *forgenetv1.RestartContainerRequest) (*forgenetv1.RestartContainerResponse, error) {
	c, err := s.mutate(ctx, "container.restart", req.GetContainerId(), s.engine.Restart)
	if err != nil {
		return nil, err
	}
	return &forgenetv1.RestartContainerResponse{Container: c}, nil
}

// mutate exécute une action qui modifie un conteneur et la consigne au journal
// d'audit, réussie ou non. Le contexte porte déjà l'utilisateur et son adresse
// (posés par auth.Guard).
func (s *ContainerServiceServer) mutate(
	ctx context.Context,
	action, id string,
	do func(context.Context, string) (orchestrator.Container, error),
) (*forgenetv1.Container, error) {
	c, err := do(ctx, id)
	meta := map[string]any{"success": err == nil}
	if err != nil {
		meta["error"] = err.Error()
	}
	if s.auditor != nil {
		_ = s.auditor.Log(ctx, action, "container:"+id, meta)
	}
	if err != nil {
		return nil, s.toStatus(err, action)
	}
	return toProtoContainer(c), nil
}

// StreamContainerLogs transmet les journaux d'un conteneur, morceau par morceau.
func (s *ContainerServiceServer) StreamContainerLogs(req *forgenetv1.StreamContainerLogsRequest, stream forgenetv1.ContainerService_StreamContainerLogsServer) error {
	opts := orchestrator.LogOptions{
		Follow: req.GetFollow(),
		Tail:   int(req.GetTailLines()),
		Stdout: req.GetStream() != forgenetv1.LogStream_LOG_STREAM_STDERR,
		Stderr: req.GetStream() != forgenetv1.LogStream_LOG_STREAM_STDOUT,
	}
	if req.GetSince() != nil {
		opts.Since = req.GetSince().AsTime()
	}

	err := s.engine.Logs(stream.Context(), req.GetContainerId(), opts, func(chunk orchestrator.LogChunk) error {
		return stream.Send(&forgenetv1.StreamContainerLogsResponse{
			Chunk:   chunk.Data,
			Emitted: timestamppb.Now(),
			Source:  toProtoLogStream(chunk.Stream),
		})
	})
	if err != nil {
		return s.toStatus(err, "lecture des journaux")
	}
	return nil
}

// toStatus traduit une erreur du domaine en statut gRPC. Le détail d'une
// erreur interne (chemin du socket, message du démon…) reste dans le journal.
func (s *ContainerServiceServer) toStatus(err error, action string) error {
	switch {
	case errors.Is(err, orchestrator.ErrNotFound):
		return status.Error(codes.NotFound, "conteneur introuvable")
	case errors.Is(err, orchestrator.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, "demande invalide : identifiant ou paramètre manquant")
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "appel annulé")
	default:
		s.logger.Error().Err(err).Str("action", action).Msg("échec de l'orchestrateur")
		return status.Error(codes.Internal, "erreur du moteur de conteneurs")
	}
}

func toProtoContainer(c orchestrator.Container) *forgenetv1.Container {
	ports := make([]*forgenetv1.PortMapping, 0, len(c.Ports))
	for _, p := range c.Ports {
		ports = append(ports, &forgenetv1.PortMapping{
			HostPort:      uint32(p.HostPort),
			ContainerPort: uint32(p.ContainerPort),
			Protocol:      p.Protocol,
		})
	}
	pc := &forgenetv1.Container{
		Id:    c.ID,
		Name:  c.Name,
		Image: c.Image,
		State: toProtoState(c.State),
		Ports: ports,
	}
	if !c.CreatedAt.IsZero() {
		pc.CreatedAt = timestamppb.New(c.CreatedAt)
	}
	return pc
}

var protoStates = map[orchestrator.State]forgenetv1.ContainerState{
	orchestrator.StateCreated:    forgenetv1.ContainerState_CONTAINER_STATE_CREATED,
	orchestrator.StateRunning:    forgenetv1.ContainerState_CONTAINER_STATE_RUNNING,
	orchestrator.StatePaused:     forgenetv1.ContainerState_CONTAINER_STATE_PAUSED,
	orchestrator.StateRestarting: forgenetv1.ContainerState_CONTAINER_STATE_RESTARTING,
	orchestrator.StateExited:     forgenetv1.ContainerState_CONTAINER_STATE_EXITED,
	orchestrator.StateDead:       forgenetv1.ContainerState_CONTAINER_STATE_DEAD,
}

// toProtoState renvoie UNSPECIFIED pour un état sans équivalent dans le
// contrat (removing, inconnu).
func toProtoState(s orchestrator.State) forgenetv1.ContainerState {
	return protoStates[s]
}

func toProtoLogStream(s orchestrator.LogStream) forgenetv1.LogStream {
	if s == orchestrator.Stderr {
		return forgenetv1.LogStream_LOG_STREAM_STDERR
	}
	return forgenetv1.LogStream_LOG_STREAM_STDOUT
}
