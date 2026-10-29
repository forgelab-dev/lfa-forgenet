package grpc

import (
	"context"
	"errors"
	"strconv"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/deploy"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// DeploymentServiceServer expose deploy.Service selon le contrat
// forgenet.v1.DeploymentService. Toutes ses méthodes sont réservées aux
// administrateurs (politique par défaut d'auth.Guard) : un fichier Compose
// donne un accès complet à l'hôte, et il peut contenir des secrets.
type DeploymentServiceServer struct {
	forgenetv1.UnimplementedDeploymentServiceServer
	deployments *deploy.Service
	logger      zerolog.Logger
}

var _ forgenetv1.DeploymentServiceServer = (*DeploymentServiceServer)(nil)

// NewDeploymentServiceServer instancie le service gRPC des déploiements.
func NewDeploymentServiceServer(deployments *deploy.Service, logger zerolog.Logger) *DeploymentServiceServer {
	return &DeploymentServiceServer{
		deployments: deployments,
		logger:      logger.With().Str("component", "grpc.deployment_service").Logger(),
	}
}

// CreateDeployment enregistre un déploiement et répond aussitôt, en attente.
func (s *DeploymentServiceServer) CreateDeployment(ctx context.Context, req *forgenetv1.CreateDeploymentRequest) (*forgenetv1.CreateDeploymentResponse, error) {
	if err := checkStrategy(req.GetStrategy()); err != nil {
		return nil, err
	}
	sourceType, content, err := fromProtoSource(req.GetSource())
	if err != nil {
		return nil, err
	}

	d, err := s.deployments.Create(ctx, deploy.CreateRequest{
		ProjectID:   req.GetProjectId(),
		SourceType:  sourceType,
		Content:     content,
		TargetEnv:   req.GetTargetEnv(),
		RequestID:   req.GetRequestId(),
		TriggeredBy: "user:" + audit.UserIDFromContext(ctx),
	})
	if err != nil {
		return nil, s.toStatus(err, "création du déploiement")
	}
	return &forgenetv1.CreateDeploymentResponse{Deployment: toProtoDeployment(d)}, nil
}

// GetDeployment renvoie un déploiement.
func (s *DeploymentServiceServer) GetDeployment(ctx context.Context, req *forgenetv1.GetDeploymentRequest) (*forgenetv1.GetDeploymentResponse, error) {
	d, err := s.deployments.Get(ctx, req.GetDeploymentId())
	if err != nil {
		return nil, s.toStatus(err, "lecture du déploiement")
	}
	return &forgenetv1.GetDeploymentResponse{Deployment: toProtoDeployment(d)}, nil
}

// ListDeployments renvoie une page de l'historique, du plus récent au plus
// ancien. Le jeton de page est la position du premier élément.
func (s *DeploymentServiceServer) ListDeployments(ctx context.Context, req *forgenetv1.ListDeploymentsRequest) (*forgenetv1.ListDeploymentsResponse, error) {
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

	// Un élément de plus que la page indique s'il existe une page suivante.
	list, err := s.deployments.List(ctx, req.GetProjectId(), pageSize+1, offset)
	if err != nil {
		return nil, s.toStatus(err, "liste des déploiements")
	}
	resp := &forgenetv1.ListDeploymentsResponse{}
	if len(list) > pageSize {
		list = list[:pageSize]
		resp.NextPageToken = strconv.Itoa(offset + pageSize)
	}
	resp.Deployments = make([]*forgenetv1.Deployment, 0, len(list))
	for _, d := range list {
		resp.Deployments = append(resp.Deployments, toProtoDeployment(d))
	}
	return resp, nil
}

// CancelDeployment annule un déploiement en attente, en construction ou en
// cours de déploiement, et renvoie son état une fois arrêté.
func (s *DeploymentServiceServer) CancelDeployment(ctx context.Context, req *forgenetv1.CancelDeploymentRequest) (*forgenetv1.CancelDeploymentResponse, error) {
	d, err := s.deployments.Cancel(ctx, req.GetDeploymentId(), req.GetReason())
	if err != nil {
		return nil, s.toStatus(err, "annulation du déploiement")
	}
	return &forgenetv1.CancelDeploymentResponse{Deployment: toProtoDeployment(d)}, nil
}

// DeleteStack supprime la stack en service d'un projet.
func (s *DeploymentServiceServer) DeleteStack(ctx context.Context, req *forgenetv1.DeleteStackRequest) (*forgenetv1.DeleteStackResponse, error) {
	if err := s.deployments.DeleteStack(ctx, req.GetProjectId(), req.GetDeleteVolumes()); err != nil {
		return nil, s.toStatus(err, "suppression de la stack")
	}
	return &forgenetv1.DeleteStackResponse{}, nil
}

// checkStrategy accepte la recréation, seule stratégie d'une stack Compose
// pour l'instant ; une stratégie absente vaut recréation.
func checkStrategy(strategy *forgenetv1.StrategyConfig) error {
	if strategy == nil || strategy.GetStrategy() == nil || strategy.GetRecreate() != nil {
		return nil
	}
	return status.Error(codes.Unimplemented, "stratégie pas encore prise en charge : seule la recréation est disponible")
}

func fromProtoSource(source *forgenetv1.DeploymentSource) (store.DeploymentSourceType, string, error) {
	switch src := source.GetSource().(type) {
	case *forgenetv1.DeploymentSource_Compose:
		return store.SourceCompose, src.Compose.GetComposeContent(), nil
	case *forgenetv1.DeploymentSource_Git:
		return "", "", status.Error(codes.Unimplemented, "source Git pas encore prise en charge")
	case *forgenetv1.DeploymentSource_Image:
		return "", "", status.Error(codes.Unimplemented, "source image pas encore prise en charge")
	default:
		return "", "", status.Error(codes.InvalidArgument, "source de déploiement manquante")
	}
}

// toStatus traduit une erreur du domaine en statut gRPC. Le détail d'une
// erreur interne reste dans le journal ; celui d'une demande refusée est
// renvoyé, pour que l'auteur du fichier sache quoi corriger.
func (s *DeploymentServiceServer) toStatus(err error, action string) error {
	switch {
	case errors.Is(err, deploy.ErrInvalid):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, deploy.ErrUnsupported):
		return status.Error(codes.Unimplemented, err.Error())
	case errors.Is(err, deploy.ErrNotCancellable), errors.Is(err, deploy.ErrBusy):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, store.ErrNotFound):
		return status.Error(codes.NotFound, "déploiement introuvable")
	case errors.Is(err, compose.ErrUnavailable):
		return status.Error(codes.Unavailable, "docker compose indisponible sur l'hôte")
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "appel annulé")
	default:
		s.logger.Error().Err(err).Str("action", action).Msg("échec du service de déploiement")
		return status.Error(codes.Internal, "erreur du service de déploiement")
	}
}

var protoDeploymentStates = map[store.DeploymentState]forgenetv1.DeploymentState{
	store.DeploymentPending:        forgenetv1.DeploymentState_DEPLOYMENT_STATE_PENDING,
	store.DeploymentBuilding:       forgenetv1.DeploymentState_DEPLOYMENT_STATE_BUILDING,
	store.DeploymentDeploying:      forgenetv1.DeploymentState_DEPLOYMENT_STATE_DEPLOYING,
	store.DeploymentHealthChecking: forgenetv1.DeploymentState_DEPLOYMENT_STATE_HEALTH_CHECKING,
	store.DeploymentSuccess:        forgenetv1.DeploymentState_DEPLOYMENT_STATE_SUCCESS,
	store.DeploymentFailed:         forgenetv1.DeploymentState_DEPLOYMENT_STATE_FAILED,
	store.DeploymentCancelled:      forgenetv1.DeploymentState_DEPLOYMENT_STATE_CANCELLED,
	store.DeploymentRollingBack:    forgenetv1.DeploymentState_DEPLOYMENT_STATE_ROLLING_BACK,
	store.DeploymentRolledBack:     forgenetv1.DeploymentState_DEPLOYMENT_STATE_ROLLED_BACK,
}

func toProtoDeployment(d *store.Deployment) *forgenetv1.Deployment {
	pd := &forgenetv1.Deployment{
		Id:                   d.ID,
		ProjectId:            d.ProjectID,
		State:                protoDeploymentStates[d.State],
		Strategy:             &forgenetv1.StrategyConfig{Strategy: &forgenetv1.StrategyConfig_Recreate{Recreate: &forgenetv1.RecreateStrategy{}}},
		ErrorMessage:         d.ErrorMessage,
		TargetEnv:            d.TargetEnv,
		PreviousDeploymentId: d.PreviousDeploymentID,
		TriggeredBy:          d.TriggeredBy,
		CreatedAt:            timestamppb.New(d.CreatedAt),
		UpdatedAt:            timestamppb.New(d.UpdatedAt),
	}
	if d.SourceType == store.SourceCompose {
		pd.Source = &forgenetv1.DeploymentSource{Source: &forgenetv1.DeploymentSource_Compose{
			Compose: &forgenetv1.ComposeSource{ComposeContent: d.SourceContent},
		}}
	}
	if d.FinishedAt != nil {
		pd.FinishedAt = timestamppb.New(*d.FinishedAt)
	}
	return pd
}
