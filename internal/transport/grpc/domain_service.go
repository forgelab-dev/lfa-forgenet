package grpc

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/domains"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/docker"
	"github.com/forgelab-dev/lfa-forgenet/internal/proxy"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// proxyDisabled répond à chaque appel tant que le reverse proxy est désactivé.
const proxyDisabled = "reverse proxy désactivé : activez proxy.enabled dans la configuration du serveur"

// DomainServiceServer expose domains.Service selon le contrat
// forgenet.v1.DomainService. Lister est ouvert aux membres ; ajouter et
// retirer sont réservés aux administrateurs (politique d'auth.Guard).
type DomainServiceServer struct {
	forgenetv1.UnimplementedDomainServiceServer
	domains *domains.Service // nil : reverse proxy désactivé
	logger  zerolog.Logger
}

var _ forgenetv1.DomainServiceServer = (*DomainServiceServer)(nil)

// NewDomainServiceServer instancie le service gRPC des domaines. Avec un
// service nil, chaque méthode répond FAILED_PRECONDITION.
func NewDomainServiceServer(d *domains.Service, logger zerolog.Logger) *DomainServiceServer {
	return &DomainServiceServer{domains: d, logger: logger.With().Str("component", "grpc.domain_service").Logger()}
}

// ListDomains liste les domaines, d'un projet ou de tous.
func (s *DomainServiceServer) ListDomains(ctx context.Context, req *forgenetv1.ListDomainsRequest) (*forgenetv1.ListDomainsResponse, error) {
	if s.domains == nil {
		return nil, status.Error(codes.FailedPrecondition, proxyDisabled)
	}
	list, err := s.domains.List(ctx, req.GetProjectId())
	if err != nil {
		return nil, s.toStatus(err, "liste des domaines")
	}
	resp := &forgenetv1.ListDomainsResponse{Domains: make([]*forgenetv1.Domain, 0, len(list))}
	for _, d := range list {
		resp.Domains = append(resp.Domains, toProtoDomain(d))
	}
	return resp, nil
}

// AddDomain relie un domaine au service déployé d'un projet.
func (s *DomainServiceServer) AddDomain(ctx context.Context, req *forgenetv1.AddDomainRequest) (*forgenetv1.AddDomainResponse, error) {
	if s.domains == nil {
		return nil, status.Error(codes.FailedPrecondition, proxyDisabled)
	}
	if req.GetPort() > 65535 {
		return nil, status.Errorf(codes.InvalidArgument, "port %d hors de 1 à 65535", req.GetPort())
	}
	d, err := s.domains.Add(ctx, domains.AddRequest{
		Host:      req.GetHost(),
		ProjectID: req.GetProjectId(),
		Service:   req.GetService(),
		Port:      int(req.GetPort()),
		CreatedBy: "user:" + audit.UserIDFromContext(ctx),
	})
	if err != nil {
		return nil, s.toStatus(err, "ajout du domaine")
	}
	return &forgenetv1.AddDomainResponse{Domain: toProtoDomain(d)}, nil
}

// RemoveDomain retire un domaine.
func (s *DomainServiceServer) RemoveDomain(ctx context.Context, req *forgenetv1.RemoveDomainRequest) (*forgenetv1.RemoveDomainResponse, error) {
	if s.domains == nil {
		return nil, status.Error(codes.FailedPrecondition, proxyDisabled)
	}
	if err := s.domains.Remove(ctx, req.GetDomainId()); err != nil {
		return nil, s.toStatus(err, "retrait du domaine")
	}
	return &forgenetv1.RemoveDomainResponse{}, nil
}

func (s *DomainServiceServer) toStatus(err error, action string) error {
	var loadErr *proxy.LoadError
	switch {
	case errors.Is(err, domains.ErrInvalid):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domains.ErrNotDeployed):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, store.ErrConflict):
		return status.Error(codes.AlreadyExists, "ce nom de domaine est déjà utilisé")
	case errors.Is(err, store.ErrNotFound):
		return status.Error(codes.NotFound, "domaine introuvable")
	case errors.As(err, &loadErr):
		return status.Error(codes.FailedPrecondition, loadErr.Error())
	case errors.Is(err, docker.ErrForeignNetwork), errors.Is(err, orchestrator.ErrInvalidArgument):
		// Réseau partagé usurpé, ou service que Docker ne peut pas brancher
		// (network_mode: host) : l'administrateur doit agir sur l'hôte.
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, proxy.ErrUnavailable):
		return status.Error(codes.Unavailable, "reverse proxy injoignable")
	case errors.Is(err, compose.ErrUnavailable):
		return status.Error(codes.Unavailable, "docker compose indisponible sur l'hôte")
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "appel annulé")
	default:
		s.logger.Error().Err(err).Str("action", action).Msg("échec du service des domaines")
		return status.Error(codes.Internal, "erreur du service des domaines")
	}
}

func toProtoDomain(d *store.Domain) *forgenetv1.Domain {
	return &forgenetv1.Domain{
		Id:        d.ID,
		Host:      d.Host,
		ProjectId: d.ProjectID,
		Service:   d.Service,
		Port:      uint32(d.Port), //nolint:gosec // G115 : port vérifié de 1 à 65535 avant l'enregistrement
		CreatedBy: d.CreatedBy,
		CreatedAt: timestamppb.New(d.CreatedAt),
	}
}
