package ws_test

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/metrics"
	"github.com/forgelab-dev/lfa-forgenet/internal/transport/ws"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// fakeSweeps joue le relevé de fond : publish enregistre un passage et
// prévient les abonnés.
type fakeSweeps struct {
	mu   sync.Mutex
	last metrics.Sweep
	subs map[chan time.Time]struct{}
}

func newFakeSweeps() *fakeSweeps { return &fakeSweeps{subs: map[chan time.Time]struct{}{}} }

func (f *fakeSweeps) LatestSweep() metrics.Sweep {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.last
}

func (f *fakeSweeps) Subscribe() (<-chan time.Time, func()) {
	ch := make(chan time.Time, 1)
	f.mu.Lock()
	f.subs[ch] = struct{}{}
	f.mu.Unlock()
	return ch, func() {
		f.mu.Lock()
		delete(f.subs, ch)
		f.mu.Unlock()
	}
}

func (f *fakeSweeps) subscribers() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.subs)
}

func (f *fakeSweeps) publish(at time.Time, ids ...string) {
	sweep := metrics.Sweep{At: at}
	for _, id := range ids {
		sweep.Containers = append(sweep.Containers, metrics.Container{ContainerID: id, CollectedAt: at, CPU: metrics.CPU{UsagePercent: 3}})
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.last = sweep
	for ch := range f.subs {
		select {
		case ch <- at:
		default:
		}
	}
}

func readAll(t *testing.T, e *env, path string) (*forgenetv1.StreamAllContainerMetricsResponse, func() *forgenetv1.StreamAllContainerMetricsResponse) {
	t.Helper()
	conn, _, err := e.dial(t, path, memberToken)
	require.NoError(t, err)
	next := func() *forgenetv1.StreamAllContainerMetricsResponse {
		return readMessage(t, conn, &forgenetv1.StreamAllContainerMetricsResponse{})
	}
	return next(), next
}

func TestAllContainersStreamFollowsTheSweeps(t *testing.T) {
	sweeps := newFakeSweeps()
	e := newEnv(t, func(c *ws.Config, _ *http.Server) { c.Sweeps = sweeps })
	sweeps.publish(collected, "a", "b")

	first, next := readAll(t, e, "/ws/v1/metrics/containers")
	assert.Equal(t, collected, first.GetCollectedAt().AsTime(), "le dernier passage part tout de suite")
	require.Len(t, first.GetContainers(), 2)
	assert.Equal(t, "a", first.GetContainers()[0].GetContainerId())

	later := collected.Add(10 * time.Second)
	sweeps.publish(later, "b")
	second := next()
	assert.Equal(t, later, second.GetCollectedAt().AsTime())
	require.Len(t, second.GetContainers(), 1, "un conteneur arrêté disparaît du relevé")
	assert.Equal(t, "b", second.GetContainers()[0].GetContainerId())
}

func TestAllContainersStreamWaitsForTheFirstSweep(t *testing.T) {
	sweeps := newFakeSweeps()
	e := newEnv(t, func(c *ws.Config, _ *http.Server) { c.Sweeps = sweeps })

	got := make(chan *forgenetv1.StreamAllContainerMetricsResponse, 1)
	go func() {
		first, _ := readAll(t, e, "/ws/v1/metrics/containers")
		got <- first
	}()
	require.Eventually(t, func() bool { return sweeps.subscribers() == 1 }, 2*time.Second, 5*time.Millisecond)
	select {
	case <-got:
		t.Fatal("aucun message avant le premier passage : il annoncerait zéro conteneur à tort")
	case <-time.After(100 * time.Millisecond):
	}
	sweeps.publish(collected, "a")
	select {
	case first := <-got:
		assert.Len(t, first.GetContainers(), 1)
	case <-time.After(5 * time.Second):
		t.Fatal("le premier passage doit partir dès qu'il a lieu")
	}
}

func TestAllContainersStreamRefusals(t *testing.T) {
	sweeps := newFakeSweeps()
	sweeps.publish(collected, "a")
	e := newEnv(t, func(c *ws.Config, _ *http.Server) { c.Sweeps = sweeps })

	_, hs, err := e.dial(t, "/ws/v1/metrics/containers", "faux")
	require.Error(t, err)
	assert.Equal(t, http.StatusUnauthorized, hs.status)

	_, hs, err = e.dial(t, "/ws/v1/metrics/containers?interval=2", memberToken)
	require.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, hs.status, "le rythme est celui du relevé de fond")

	require.Eventually(t, func() bool { return sweeps.subscribers() == 0 }, 2*time.Second, 5*time.Millisecond,
		"un flux refusé libère son abonnement")
}

func TestAllContainersRouteNeedsTheRecorder(t *testing.T) {
	e := newEnv(t, nil)
	_, hs, err := e.dial(t, "/ws/v1/metrics/containers", memberToken)
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, hs.status)
}
