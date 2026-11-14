package grpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/metrics"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

const metricsContainer = "c0ffee000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

// fakeHistory tient un historique fixe, à la place de metrics.Recorder.
type fakeHistory struct {
	host       []metrics.Host
	containers map[string][]metrics.Container
}

func (h *fakeHistory) Config() metrics.RecorderConfig {
	return metrics.RecorderConfig{Window: time.Hour, HostEvery: 5 * time.Second, ContainerEvery: 10 * time.Second}
}

func (h *fakeHistory) LatestHost() (metrics.Host, bool) {
	if len(h.host) == 0 {
		return metrics.Host{}, false
	}
	return h.host[len(h.host)-1], true
}

func (h *fakeHistory) HostHistory(from, to time.Time) []metrics.Host {
	var out []metrics.Host
	for _, s := range h.host {
		if !s.CollectedAt.Before(from) && !s.CollectedAt.After(to) {
			out = append(out, s)
		}
	}
	return out
}

func (h *fakeHistory) LatestContainer(id string) (metrics.Container, bool) {
	s := h.containers[id]
	if len(s) == 0 {
		return metrics.Container{}, false
	}
	return s[len(s)-1], true
}

func (h *fakeHistory) ContainerHistory(id string, from, to time.Time) ([]metrics.Container, bool) {
	all, ok := h.containers[id]
	if !ok {
		return nil, false
	}
	var out []metrics.Container
	for _, s := range all {
		if !s.CollectedAt.Before(from) && !s.CollectedAt.After(to) {
			out = append(out, s)
		}
	}
	return out, true
}

// fakeStats mesure sur-le-champ un seul conteneur connu.
type fakeStats struct{ calls int }

func (s *fakeStats) Stats(_ context.Context, id string, _ bool) (container.StatsResponse, error) {
	if id != "neuf" {
		return container.StatsResponse{}, orchestrator.ErrNotFound
	}
	s.calls++
	return container.StatsResponse{ID: "neuf", Read: time.Now(), MemoryStats: container.MemoryStats{Usage: 42, Limit: 100}}, nil
}

type metricsTestServer struct {
	*testServer
	metrics forgenetv1.MetricsServiceClient
	stats   *fakeStats
}

func startMetricsServer(t *testing.T, history *fakeHistory) *metricsTestServer {
	t.Helper()
	mts := &metricsTestServer{stats: &fakeStats{}}
	mts.testServer = startTestServer(t, &fakeEngine{}, func(srv *grpc.Server, _ store.Store, _ audit.Auditor) {
		forgenetv1.RegisterMetricsServiceServer(srv, transportGRPC.NewMetricsServiceServer(history, mts.stats, zerolog.Nop()))
	})
	mts.metrics = forgenetv1.NewMetricsServiceClient(mts.conn)
	return mts
}

// recentHistory : une heure de l'hôte toutes les 5 s, et un conteneur
// mesuré toutes les 10 s, jusqu'à maintenant.
func recentHistory() *fakeHistory {
	now := time.Now()
	h := &fakeHistory{containers: map[string][]metrics.Container{}}
	for i := 720; i >= 0; i-- {
		at := now.Add(-time.Duration(i) * 5 * time.Second)
		h.host = append(h.host, metrics.Host{
			CollectedAt: at,
			CPU:         metrics.CPU{UsagePercent: float64(i % 100)},
			Network:     metrics.Network{RxBytes: uint64(720-i) * 1000},
			Storage:     metrics.Storage{UsagePercent: 60},
		})
	}
	for i := 360; i >= 0; i-- {
		at := now.Add(-time.Duration(i) * 10 * time.Second)
		h.containers[metricsContainer] = append(h.containers[metricsContainer], metrics.Container{
			ContainerID: metricsContainer,
			CollectedAt: at,
			Memory:      metrics.Memory{UsageBytes: 100, LimitBytes: 1000, UsagePercent: 10},
		})
	}
	return h
}

func known(n forgenetv1.MetricName) *forgenetv1.MetricKey {
	return &forgenetv1.MetricKey{Key: &forgenetv1.MetricKey_Known{Known: n}}
}

func TestMetricsService_RolePolicy(t *testing.T) {
	ts := startMetricsServer(t, recentHistory())
	calls := map[string]func(ctx context.Context) error{
		"GetSystem": func(ctx context.Context) error {
			_, err := ts.metrics.GetSystemMetrics(ctx, &forgenetv1.GetSystemMetricsRequest{})
			return err
		},
		"QuerySystem": func(ctx context.Context) error {
			_, err := ts.metrics.QuerySystemMetrics(ctx, &forgenetv1.QuerySystemMetricsRequest{})
			return err
		},
		"GetContainer": func(ctx context.Context) error {
			_, err := ts.metrics.GetContainerMetrics(ctx, &forgenetv1.GetContainerMetricsRequest{ContainerId: metricsContainer})
			return err
		},
		"QueryContainer": func(ctx context.Context) error {
			_, err := ts.metrics.QueryContainerMetrics(ctx, &forgenetv1.QueryContainerMetricsRequest{ContainerId: metricsContainer})
			return err
		},
	}
	for name, call := range calls {
		assert.Equal(t, codes.Unauthenticated, status.Code(call(context.Background())), "%s sans jeton", name)
		assert.Equal(t, codes.OK, status.Code(call(as(ts.memberToken))), "%s par un membre", name)
		assert.Equal(t, codes.OK, status.Code(call(as(ts.adminToken))), "%s par un admin", name)
	}
}

func TestMetricsService_GetSystemMetrics(t *testing.T) {
	empty := startMetricsServer(t, &fakeHistory{})
	_, err := empty.metrics.GetSystemMetrics(as(empty.memberToken), &forgenetv1.GetSystemMetricsRequest{})
	assert.Equal(t, codes.Unavailable, status.Code(err), "aucune mesure juste après le démarrage")

	history := recentHistory()
	ts := startMetricsServer(t, history)
	res, err := ts.metrics.GetSystemMetrics(as(ts.memberToken), &forgenetv1.GetSystemMetricsRequest{})
	require.NoError(t, err)
	assert.Equal(t, history.host[len(history.host)-1].CollectedAt.UTC(), res.GetMetrics().GetCollectedAt().AsTime())
}

func TestMetricsService_GetContainerMetrics(t *testing.T) {
	ts := startMetricsServer(t, recentHistory())
	ctx := as(ts.memberToken)

	res, err := ts.metrics.GetContainerMetrics(ctx, &forgenetv1.GetContainerMetricsRequest{ContainerId: metricsContainer})
	require.NoError(t, err)
	assert.Equal(t, uint64(100), res.GetMetrics().GetMemory().GetUsageBytes())
	assert.Zero(t, ts.stats.calls, "une mesure récente du relevé suffit")

	res, err = ts.metrics.GetContainerMetrics(ctx, &forgenetv1.GetContainerMetricsRequest{ContainerId: "neuf"})
	require.NoError(t, err, "un conteneur pas encore relevé est mesuré sur-le-champ")
	assert.Equal(t, uint64(42), res.GetMetrics().GetMemory().GetUsageBytes())

	_, err = ts.metrics.GetContainerMetrics(ctx, &forgenetv1.GetContainerMetricsRequest{ContainerId: "absent"})
	assert.Equal(t, codes.NotFound, status.Code(err))
	_, err = ts.metrics.GetContainerMetrics(ctx, &forgenetv1.GetContainerMetricsRequest{})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestMetricsService_QuerySystemMetrics(t *testing.T) {
	ts := startMetricsServer(t, recentHistory())
	ctx := as(ts.memberToken)

	all, err := ts.metrics.QuerySystemMetrics(ctx, &forgenetv1.QuerySystemMetricsRequest{})
	require.NoError(t, err)
	assert.Len(t, all.GetSeries(), 22, "toutes les métriques de l'hôte, sans la valeur nulle")
	assert.Equal(t, forgenetv1.MetricName_METRIC_NAME_CPU_USAGE_PERCENT, all.GetSeries()[0].GetKey().GetKnown(), "dans l'ordre du contrat")
	assert.LessOrEqual(t, len(all.GetSeries()[0].GetSamples()), 721, "une heure au rythme du relevé")

	now := time.Now()
	res, err := ts.metrics.QuerySystemMetrics(ctx, &forgenetv1.QuerySystemMetricsRequest{
		Range:     &forgenetv1.TimeRange{StartTime: timestamppb.New(now.Add(-15 * time.Minute)), EndTime: timestamppb.New(now)},
		MaxPoints: 90,
		Keys: []*forgenetv1.MetricKey{
			known(forgenetv1.MetricName_METRIC_NAME_NETWORK_RX_BYTES_TOTAL),
			known(forgenetv1.MetricName_METRIC_NAME_CPU_USAGE_PERCENT),
		},
	})
	require.NoError(t, err)
	require.Len(t, res.GetSeries(), 2, "seulement les métriques demandées")
	rx := res.GetSeries()[0]
	assert.Equal(t, forgenetv1.MetricUnit_METRIC_UNIT_BYTES, rx.GetUnit())
	assert.LessOrEqual(t, len(rx.GetSamples()), 91, "au plus 90 tranches, plus une à cheval sur un bord")
	assert.GreaterOrEqual(t, len(rx.GetSamples()), 89)
	assert.Equal(t, forgenetv1.MetricUnit_METRIC_UNIT_PERCENT, res.GetSeries()[1].GetUnit())

	_, err = ts.metrics.QuerySystemMetrics(ctx, &forgenetv1.QuerySystemMetricsRequest{
		Range: &forgenetv1.TimeRange{StartTime: timestamppb.New(now), EndTime: timestamppb.New(now.Add(-time.Minute))},
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err), "début après la fin")

	_, err = ts.metrics.QuerySystemMetrics(ctx, &forgenetv1.QuerySystemMetricsRequest{
		Keys: []*forgenetv1.MetricKey{{Key: &forgenetv1.MetricKey_Custom{Custom: "acme.queue_depth"}}},
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err), "aucune extension installée")
}

func TestMetricsService_QueryContainerMetrics(t *testing.T) {
	ts := startMetricsServer(t, recentHistory())
	ctx := as(ts.memberToken)

	res, err := ts.metrics.QueryContainerMetrics(ctx, &forgenetv1.QueryContainerMetricsRequest{
		ContainerId: metricsContainer,
		Keys:        []*forgenetv1.MetricKey{known(forgenetv1.MetricName_METRIC_NAME_MEMORY_USAGE_PERCENT)},
		StepSeconds: 60,
	})
	require.NoError(t, err)
	assert.Equal(t, metricsContainer, res.GetContainerId())
	require.Len(t, res.GetSeries(), 1)
	samples := res.GetSeries()[0].GetSamples()
	assert.InDelta(t, 61, len(samples), 1, "une heure par tranches d'une minute")
	assert.InDelta(t, 10.0, samples[0].GetValue(), 0.001)

	all, err := ts.metrics.QueryContainerMetrics(ctx, &forgenetv1.QueryContainerMetricsRequest{ContainerId: metricsContainer})
	require.NoError(t, err)
	assert.Len(t, all.GetSeries(), 16, "les métriques d'un conteneur, sans disque hôte ni uptime")

	_, err = ts.metrics.QueryContainerMetrics(ctx, &forgenetv1.QueryContainerMetricsRequest{
		ContainerId: metricsContainer,
		Keys:        []*forgenetv1.MetricKey{known(forgenetv1.MetricName_METRIC_NAME_STORAGE_USED_BYTES)},
	})
	assert.Equal(t, codes.InvalidArgument, status.Code(err), "le disque de l'hôte n'est pas une métrique de conteneur")

	_, err = ts.metrics.QueryContainerMetrics(ctx, &forgenetv1.QueryContainerMetricsRequest{ContainerId: "absent"})
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestMetricsService_StreamsPointToWebSocket(t *testing.T) {
	ts := startMetricsServer(t, recentHistory())
	stream, err := ts.metrics.StreamSystemMetrics(as(ts.adminToken), &forgenetv1.StreamSystemMetricsRequest{})
	require.NoError(t, err)
	_, err = stream.Recv()
	assert.Equal(t, codes.Unimplemented, status.Code(err))
	assert.Contains(t, status.Convert(err).Message(), "/ws/v1/metrics/system")
}
