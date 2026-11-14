package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/metrics"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/transport/metricspb"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// MetricsHistory est l'historique court tenu par metrics.Recorder.
type MetricsHistory interface {
	Config() metrics.RecorderConfig
	LatestHost() (metrics.Host, bool)
	HostHistory(from, to time.Time) []metrics.Host
	LatestContainer(id string) (metrics.Container, bool)
	ContainerHistory(id string, from, to time.Time) ([]metrics.Container, bool)
}

// MetricsServiceServer sert l'état instantané et l'historique court des
// mesures (ADR-0010). Les flux en direct passent par WebSocket (ADR-0005) :
// les méthodes Stream* restent non implémentées ici.
type MetricsServiceServer struct {
	forgenetv1.UnimplementedMetricsServiceServer
	history MetricsHistory
	stats   metrics.StatsSource
	now     func() time.Time
	logger  zerolog.Logger
}

var _ forgenetv1.MetricsServiceServer = (*MetricsServiceServer)(nil)

// NewMetricsServiceServer crée le service. stats sert à mesurer sur-le-champ
// un conteneur que le relevé de fond n'a pas encore vu.
func NewMetricsServiceServer(history MetricsHistory, stats metrics.StatsSource, logger zerolog.Logger) *MetricsServiceServer {
	return &MetricsServiceServer{
		history: history,
		stats:   stats,
		now:     time.Now,
		logger:  logger.With().Str("component", "grpc.metrics_service").Logger(),
	}
}

// GetSystemMetrics renvoie la dernière mesure de l'hôte.
func (s *MetricsServiceServer) GetSystemMetrics(context.Context, *forgenetv1.GetSystemMetricsRequest) (*forgenetv1.GetSystemMetricsResponse, error) {
	h, ok := s.history.LatestHost()
	if !ok {
		return nil, status.Error(codes.Unavailable, "aucune mesure de l'hôte pour l'instant, réessayez dans quelques secondes")
	}
	return &forgenetv1.GetSystemMetricsResponse{Metrics: metricspb.ToSystemMetrics(h)}, nil
}

// GetContainerMetrics renvoie la dernière mesure d'un conteneur, ou en prend
// une si le relevé de fond ne le suit pas encore (conteneur tout juste
// démarré).
func (s *MetricsServiceServer) GetContainerMetrics(ctx context.Context, req *forgenetv1.GetContainerMetricsRequest) (*forgenetv1.GetContainerMetricsResponse, error) {
	id := req.GetContainerId()
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "container_id est obligatoire")
	}
	if c, ok := s.history.LatestContainer(id); ok && s.now().Sub(c.CollectedAt) <= 2*s.history.Config().ContainerEvery {
		return &forgenetv1.GetContainerMetricsResponse{Metrics: metricspb.ToContainerMetrics(c)}, nil
	}
	c, err := metrics.NewContainerCollector(s.stats, id).Collect(ctx)
	if err != nil {
		return nil, s.toStatus(err)
	}
	return &forgenetv1.GetContainerMetricsResponse{Metrics: metricspb.ToContainerMetrics(c)}, nil
}

// QuerySystemMetrics renvoie l'historique de l'hôte.
func (s *MetricsServiceServer) QuerySystemMetrics(_ context.Context, req *forgenetv1.QuerySystemMetricsRequest) (*forgenetv1.QuerySystemMetricsResponse, error) {
	if err := metricspb.CheckKeys(req.GetKeys(), true); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	from, to, err := s.window(req.GetRange())
	if err != nil {
		return nil, err
	}
	cfg := s.history.Config()
	step := metrics.Step(from, to, seconds(req.GetStepSeconds()), cfg.HostEvery, int(req.GetMaxPoints()))
	series, err := metricspb.HostSeries(s.history.HostHistory(from, to), req.GetKeys(), step)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &forgenetv1.QuerySystemMetricsResponse{Series: series}, nil
}

// QueryContainerMetrics renvoie l'historique d'un conteneur. Un conteneur
// arrêté garde le sien tant qu'il reste dans la fenêtre.
func (s *MetricsServiceServer) QueryContainerMetrics(_ context.Context, req *forgenetv1.QueryContainerMetricsRequest) (*forgenetv1.QueryContainerMetricsResponse, error) {
	id := req.GetContainerId()
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "container_id est obligatoire")
	}
	if err := metricspb.CheckKeys(req.GetKeys(), false); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	from, to, err := s.window(req.GetRange())
	if err != nil {
		return nil, err
	}
	samples, found := s.history.ContainerHistory(id, from, to)
	if !found {
		return nil, status.Error(codes.NotFound, "aucun historique pour ce conteneur : inconnu, ou arrêté depuis plus d'une heure")
	}
	cfg := s.history.Config()
	step := metrics.Step(from, to, seconds(req.GetStepSeconds()), cfg.ContainerEvery, int(req.GetMaxPoints()))
	series, err := metricspb.ContainerSeries(samples, req.GetKeys(), step)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// L'identifiant complet, même si la demande portait un préfixe.
	fullID := id
	if len(samples) > 0 {
		fullID = samples[len(samples)-1].ContainerID
	}
	return &forgenetv1.QueryContainerMetricsResponse{ContainerId: fullID, Series: series}, nil
}

// window lit la plage demandée : sans début, la fenêtre entière de
// l'historique ; sans fin, maintenant. Une plage plus longue que la fenêtre
// est acceptée et renvoie ce qui existe.
func (s *MetricsServiceServer) window(r *forgenetv1.TimeRange) (from, to time.Time, err error) {
	now := s.now()
	to = now
	if end := r.GetEndTime(); end != nil {
		if err := end.CheckValid(); err != nil {
			return from, to, status.Error(codes.InvalidArgument, "range.end_time invalide")
		}
		to = end.AsTime()
	}
	from = to.Add(-s.history.Config().Window)
	if start := r.GetStartTime(); start != nil {
		if err := start.CheckValid(); err != nil {
			return from, to, status.Error(codes.InvalidArgument, "range.start_time invalide")
		}
		from = start.AsTime()
	}
	if from.After(to) {
		return from, to, status.Error(codes.InvalidArgument, "range.start_time est postérieur à range.end_time")
	}
	return from, to, nil
}

func seconds(n uint32) time.Duration { return time.Duration(n) * time.Second }

func (s *MetricsServiceServer) toStatus(err error) error {
	switch {
	case errors.Is(err, orchestrator.ErrNotFound):
		return status.Error(codes.NotFound, "conteneur introuvable")
	case errors.Is(err, orchestrator.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, "identifiant de conteneur invalide")
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "appel annulé")
	default:
		s.logger.Error().Err(err).Msg("mesure d'un conteneur impossible")
		return status.Error(codes.Unavailable, "mesure indisponible")
	}
}

// StreamSystemMetrics : le flux de l'hôte est servi par WebSocket (ADR-0005).
func (s *MetricsServiceServer) StreamSystemMetrics(*forgenetv1.StreamSystemMetricsRequest, forgenetv1.MetricsService_StreamSystemMetricsServer) error {
	return status.Error(codes.Unimplemented, "flux servi par WebSocket : GET /ws/v1/metrics/system")
}

// StreamContainerMetrics : le flux d'un conteneur est servi par WebSocket (ADR-0005).
func (s *MetricsServiceServer) StreamContainerMetrics(*forgenetv1.StreamContainerMetricsRequest, forgenetv1.MetricsService_StreamContainerMetricsServer) error {
	return status.Error(codes.Unimplemented, "flux servi par WebSocket : GET /ws/v1/metrics/containers/{id}")
}
