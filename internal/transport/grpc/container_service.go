package grpc

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/docker"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

const (
	defaultPageSize = 50
	maxPageSize     = 100
)

// ContainerServiceServer implémente le contrat protobuf ContainerServiceServer.
type ContainerServiceServer struct {
	forgenetv1.UnimplementedContainerServiceServer
	orchestrator docker.Orchestrator
	auditor      audit.Auditor
	logger       zerolog.Logger
}

// Assure que ContainerServiceServer implémente bien forgenetv1.ContainerServiceServer.
var _ forgenetv1.ContainerServiceServer = (*ContainerServiceServer)(nil)

// NewContainerServiceServer instancie le serveur gRPC ContainerService.
func NewContainerServiceServer(orch docker.Orchestrator, auditor audit.Auditor, logger zerolog.Logger) *ContainerServiceServer {
	return &ContainerServiceServer{
		orchestrator: orch,
		auditor:      auditor,
		logger:       logger.With().Str("component", "grpc.container_service").Logger(),
	}
}

// ListContainers retourne la liste paginée des conteneurs existants.
func (s *ContainerServiceServer) ListContainers(ctx context.Context, req *forgenetv1.ListContainersRequest) (*forgenetv1.ListContainersResponse, error) {
	_, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}

	includeStopped := false
	var pageSize int32 = defaultPageSize
	pageToken := ""

	if req != nil {
		includeStopped = req.GetIncludeStopped()
		if req.GetPageSize() > 0 {
			pageSize = req.GetPageSize()
			if pageSize > maxPageSize {
				pageSize = maxPageSize
			}
		}
		pageToken = req.GetPageToken()
	}

	allContainers, err := s.orchestrator.ListContainers(ctx, includeStopped)
	if err != nil {
		s.logger.Error().Err(err).Msg("échec de la récupération des conteneurs")
		return nil, status.Errorf(codes.Internal, "erreur orchestrateur : %v", err)
	}

	// Gestion de l'offset de pagination
	offset := 0
	if pageToken != "" {
		parsed, err := strconv.Atoi(pageToken)
		if err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	if offset > len(allContainers) {
		offset = len(allContainers)
	}

	end := offset + int(pageSize)
	nextPageToken := ""
	if end < len(allContainers) {
		nextPageToken = strconv.Itoa(end)
	} else {
		end = len(allContainers)
	}

	paged := allContainers[offset:end]
	protoContainers := make([]*forgenetv1.Container, len(paged))
	for i, c := range paged {
		protoContainers[i] = c.ToProto()
	}

	return &forgenetv1.ListContainersResponse{
		Containers:    protoContainers,
		NextPageToken: nextPageToken,
	}, nil
}

// GetContainer récupère les informations détaillées d'un conteneur spécifique.
func (s *ContainerServiceServer) GetContainer(ctx context.Context, req *forgenetv1.GetContainerRequest) (*forgenetv1.GetContainerResponse, error) {
	_, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}

	if req == nil || req.GetContainerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "l'identifiant du conteneur est requis")
	}

	info, err := s.orchestrator.GetContainer(ctx, req.GetContainerId())
	if err != nil {
		if errors.Is(err, docker.ErrContainerNotFound) {
			return nil, status.Errorf(codes.NotFound, "conteneur %q introuvable", req.GetContainerId())
		}
		s.logger.Error().Err(err).Str("container_id", req.GetContainerId()).Msg("erreur inspection conteneur")
		return nil, status.Errorf(codes.Internal, "erreur orchestrateur : %v", err)
	}

	return &forgenetv1.GetContainerResponse{
		Container: info.ToProto(),
	}, nil
}

// StartContainer démarre un conteneur arrêté.
func (s *ContainerServiceServer) StartContainer(ctx context.Context, req *forgenetv1.StartContainerRequest) (*forgenetv1.StartContainerResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}

	if !isAuthorizedOperator(claims.Role) {
		return nil, status.Error(codes.PermissionDenied, "droits insuffisants pour démarrer un conteneur")
	}

	if req == nil || req.GetContainerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "l'identifiant du conteneur est requis")
	}

	info, err := s.orchestrator.StartContainer(ctx, req.GetContainerId())
	if err != nil {
		s.auditAction(ctx, claims.UserID, "container.start", req.GetContainerId(), false, map[string]string{"error": err.Error()})
		if errors.Is(err, docker.ErrContainerNotFound) {
			return nil, status.Errorf(codes.NotFound, "conteneur %q introuvable", req.GetContainerId())
		}
		s.logger.Error().Err(err).Str("container_id", req.GetContainerId()).Msg("erreur démarrage conteneur")
		return nil, status.Errorf(codes.Internal, "erreur démarrage conteneur : %v", err)
	}

	s.auditAction(ctx, claims.UserID, "container.start", req.GetContainerId(), true, nil)

	return &forgenetv1.StartContainerResponse{
		Container: info.ToProto(),
	}, nil
}

// StopContainer arrête proprement un conteneur en cours d'exécution.
func (s *ContainerServiceServer) StopContainer(ctx context.Context, req *forgenetv1.StopContainerRequest) (*forgenetv1.StopContainerResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}

	if !isAuthorizedOperator(claims.Role) {
		return nil, status.Error(codes.PermissionDenied, "droits insuffisants pour arrêter un conteneur")
	}

	if req == nil || req.GetContainerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "l'identifiant du conteneur est requis")
	}

	timeoutSec := req.GetTimeoutSeconds()
	info, err := s.orchestrator.StopContainer(ctx, req.GetContainerId(), timeoutSec)
	if err != nil {
		s.auditAction(ctx, claims.UserID, "container.stop", req.GetContainerId(), false, map[string]string{
			"error":   err.Error(),
			"timeout": strconv.FormatUint(uint64(timeoutSec), 10),
		})
		if errors.Is(err, docker.ErrContainerNotFound) {
			return nil, status.Errorf(codes.NotFound, "conteneur %q introuvable", req.GetContainerId())
		}
		s.logger.Error().Err(err).Str("container_id", req.GetContainerId()).Msg("erreur arrêt conteneur")
		return nil, status.Errorf(codes.Internal, "erreur arrêt conteneur : %v", err)
	}

	s.auditAction(ctx, claims.UserID, "container.stop", req.GetContainerId(), true, map[string]string{
		"timeout": strconv.FormatUint(uint64(timeoutSec), 10),
	})

	return &forgenetv1.StopContainerResponse{
		Container: info.ToProto(),
	}, nil
}

// RestartContainer redémarre un conteneur.
func (s *ContainerServiceServer) RestartContainer(ctx context.Context, req *forgenetv1.RestartContainerRequest) (*forgenetv1.RestartContainerResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}

	if !isAuthorizedOperator(claims.Role) {
		return nil, status.Error(codes.PermissionDenied, "droits insuffisants pour redémarrer un conteneur")
	}

	if req == nil || req.GetContainerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "l'identifiant du conteneur est requis")
	}

	info, err := s.orchestrator.RestartContainer(ctx, req.GetContainerId())
	if err != nil {
		s.auditAction(ctx, claims.UserID, "container.restart", req.GetContainerId(), false, map[string]string{"error": err.Error()})
		if errors.Is(err, docker.ErrContainerNotFound) {
			return nil, status.Errorf(codes.NotFound, "conteneur %q introuvable", req.GetContainerId())
		}
		s.logger.Error().Err(err).Str("container_id", req.GetContainerId()).Msg("erreur redémarrage conteneur")
		return nil, status.Errorf(codes.Internal, "erreur redémarrage conteneur : %v", err)
	}

	s.auditAction(ctx, claims.UserID, "container.restart", req.GetContainerId(), true, nil)

	return &forgenetv1.RestartContainerResponse{
		Container: info.ToProto(),
	}, nil
}

// StreamContainerLogs ouvre un flux gRPC pour lire les logs d'un conteneur en temps réel.
func (s *ContainerServiceServer) StreamContainerLogs(
	req *forgenetv1.StreamContainerLogsRequest,
	stream forgenetv1.ContainerService_StreamContainerLogsServer,
) error {
	ctx := stream.Context()
	_, ok := auth.UserFromContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "authentification requise")
	}

	if req == nil || req.GetContainerId() == "" {
		return status.Error(codes.InvalidArgument, "l'identifiant du conteneur est requis")
	}

	showStdout := true
	showStderr := true
	switch req.GetStream() {
	case forgenetv1.LogStream_LOG_STREAM_STDOUT:
		showStderr = false
	case forgenetv1.LogStream_LOG_STREAM_STDERR:
		showStdout = false
	case forgenetv1.LogStream_LOG_STREAM_UNSPECIFIED:
		// Garde les deux
	}

	var since time.Time
	if req.GetSince() != nil {
		since = req.GetSince().AsTime()
	}

	reader, err := s.orchestrator.StreamLogs(ctx, req.GetContainerId(), req.GetFollow(), int(req.GetTailLines()), since, showStdout, showStderr)
	if err != nil {
		if errors.Is(err, docker.ErrContainerNotFound) {
			return status.Errorf(codes.NotFound, "conteneur %q introuvable", req.GetContainerId())
		}
		s.logger.Error().Err(err).Str("container_id", req.GetContainerId()).Msg("erreur ouverture flux logs")
		return status.Errorf(codes.Internal, "erreur lecture des logs : %v", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	return forwardLogs(ctx, reader, req.GetStream(), stream)
}

func forwardLogs(
	ctx context.Context,
	reader io.Reader,
	requestedStream forgenetv1.LogStream,
	stream forgenetv1.ContainerService_StreamContainerLogsServer,
) error {
	header := make([]byte, 8)
	buf := make([]byte, 4096)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Tentative de lecture de l'en-tête multiplexé Docker (8 octets)
		n, err := io.ReadFull(reader, header)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				if n > 0 {
					// Données résiduelles non multiplexées
					_ = sendChunk(stream, header[:n], requestedStream)
				}
				return nil
			}
			return err
		}

		// Détection d'un en-tête Docker multiplexé valide :
		// header[0]: 1 (stdout) ou 2 (stderr)
		// header[1..3]: 0
		if (header[0] == 1 || header[0] == 2) && header[1] == 0 && header[2] == 0 && header[3] == 0 {
			streamSource := forgenetv1.LogStream_LOG_STREAM_STDOUT
			if header[0] == 2 {
				streamSource = forgenetv1.LogStream_LOG_STREAM_STDERR
			}

			frameSize := binary.BigEndian.Uint32(header[4:8])
			if frameSize == 0 {
				continue
			}

			// Lecture du contenu de la frame par morceaux
			var remaining = int64(frameSize)
			for remaining > 0 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				toRead := int64(len(buf))
				if remaining < toRead {
					toRead = remaining
				}

				chunkN, readErr := reader.Read(buf[:toRead])
				if chunkN > 0 {
					if err := sendChunk(stream, buf[:chunkN], streamSource); err != nil {
						return err
					}
					remaining -= int64(chunkN)
				}
				if readErr != nil {
					if errors.Is(readErr, io.EOF) {
						return nil
					}
					return readErr
				}
			}
		} else {
			// Flux brut non multiplexé (ex: TTY alloué ou flux textuel)
			defaultSource := requestedStream
			if defaultSource == forgenetv1.LogStream_LOG_STREAM_UNSPECIFIED {
				defaultSource = forgenetv1.LogStream_LOG_STREAM_STDOUT
			}

			if err := sendChunk(stream, header[:n], defaultSource); err != nil {
				return err
			}

			rawN, rawErr := reader.Read(buf)
			if rawN > 0 {
				if err := sendChunk(stream, buf[:rawN], defaultSource); err != nil {
					return err
				}
			}
			if rawErr != nil {
				if errors.Is(rawErr, io.EOF) {
					return nil
				}
				return rawErr
			}
		}
	}
}

func sendChunk(stream forgenetv1.ContainerService_StreamContainerLogsServer, data []byte, source forgenetv1.LogStream) error {
	cp := make([]byte, len(data))
	copy(cp, data)
	return stream.Send(&forgenetv1.StreamContainerLogsResponse{
		Chunk:   cp,
		Emitted: timestamppb.Now(),
		Source:  source,
	})
}

func isAuthorizedOperator(role store.Role) bool {
	return role == store.RoleAdmin || role == store.RoleMember
}

func (s *ContainerServiceServer) auditAction(
	ctx context.Context,
	userID string,
	action string,
	containerID string,
	success bool,
	meta map[string]string,
) {
	if s.auditor == nil {
		return
	}
	auditCtx := audit.WithAuditContext(ctx, userID, "")
	metadata := make(map[string]any, len(meta)+1)
	for k, v := range meta {
		metadata[k] = v
	}
	metadata["success"] = success
	_ = s.auditor.Log(auditCtx, action, fmt.Sprintf("container:%s", containerID), metadata)
}
