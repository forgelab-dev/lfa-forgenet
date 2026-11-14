package metrics

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
)

func TestRing_KeepsTheNewestInOrder(t *testing.T) {
	r := newRing[int](3)
	_, ok := r.last()
	assert.False(t, ok, "tampon vide")

	for i := 1; i <= 5; i++ {
		r.push(i)
	}
	var got []int
	r.each(func(v int) { got = append(got, v) })
	assert.Equal(t, []int{3, 4, 5}, got, "le plus ancien est remplacé, l'ordre est gardé")
	last, ok := r.last()
	require.True(t, ok)
	assert.Equal(t, 5, last)
}

func TestCapacity_CoversTheWindow(t *testing.T) {
	assert.Equal(t, 721, capacity(time.Hour, 5*time.Second))
	assert.Equal(t, 361, capacity(time.Hour, 10*time.Second))
}

// fakeDocker simule le démon : des conteneurs en marche, des statistiques
// dont le compteur processeur avance d'un relevé à l'autre.
type fakeDocker struct {
	mu       sync.Mutex
	running  []string
	listErr  error
	statsErr map[string]error
	delay    time.Duration // durée de chaque réponse du démon
	calls    map[string]int
	previous map[string][]bool // withPrevious reçu, appel après appel
}

func newFakeDocker(ids ...string) *fakeDocker {
	return &fakeDocker{running: ids, statsErr: map[string]error{}, calls: map[string]int{}, previous: map[string][]bool{}}
}

func (d *fakeDocker) RunningIDs(context.Context) ([]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return append([]string(nil), d.running...), d.listErr
}

func (d *fakeDocker) Stats(ctx context.Context, id string, withPrevious bool) (container.StatsResponse, error) {
	if d.delay > 0 {
		select {
		case <-time.After(d.delay):
		case <-ctx.Done():
			return container.StatsResponse{}, ctx.Err()
		}
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.statsErr[id]; err != nil {
		return container.StatsResponse{}, err
	}
	d.calls[id]++
	d.previous[id] = append(d.previous[id], withPrevious)
	n := uint64(d.calls[id])
	return container.StatsResponse{
		ID:   id,
		Read: time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC).Add(time.Duration(n) * 10 * time.Second),
		CPUStats: container.CPUStats{
			CPUUsage:    container.CPUUsage{TotalUsage: n * 1_000},
			SystemUsage: n * 10_000,
			OnlineCPUs:  2,
		},
		PreCPUStats: container.CPUStats{
			CPUUsage:    container.CPUUsage{TotalUsage: (n - 1) * 1_000},
			SystemUsage: (n - 1) * 10_000,
			OnlineCPUs:  2,
		},
		MemoryStats: container.MemoryStats{Usage: 100 << 20, Limit: 1 << 30},
	}, nil
}

type fakeHost struct {
	mu  sync.Mutex
	n   int
	err error
}

func (h *fakeHost) Collect(context.Context) (Host, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.err != nil {
		return Host{}, h.err
	}
	h.n++
	return Host{CollectedAt: time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC).Add(time.Duration(h.n) * 5 * time.Second)}, nil
}

const (
	idA = "aaaaaaaaaaaa1111111111111111111111111111111111111111111111111111"
	idB = "bbbbbbbbbbbb2222222222222222222222222222222222222222222222222222"
	// idA2 partage les 12 premiers caractères de idA.
	idA2 = "aaaaaaaaaaaa3333333333333333333333333333333333333333333333333333"
)

func newTestRecorder(docker *fakeDocker, host *fakeHost, clock *time.Time) *Recorder {
	r := NewRecorder(host, docker, docker, RecorderConfig{}, zerolog.Nop())
	r.now = func() time.Time { return *clock }
	return r
}

func TestRecorder_RecordsEveryRunningContainer(t *testing.T) {
	docker := newFakeDocker(idB, idA)
	clock := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	r := newTestRecorder(docker, &fakeHost{}, &clock)
	signals, cancel := r.Subscribe()
	defer cancel()

	assert.True(t, r.LatestSweep().At.IsZero(), "aucun passage avant le premier relevé")
	r.recordContainers(context.Background())
	clock = clock.Add(10 * time.Second)
	r.recordContainers(context.Background())

	sweep := r.LatestSweep()
	assert.Equal(t, clock, sweep.At)
	require.Len(t, sweep.Containers, 2)
	assert.Equal(t, idA, sweep.Containers[0].ContainerID, "passage trié par identifiant")
	assert.InDelta(t, 20.0, sweep.Containers[0].CPU.UsagePercent, 0.001, "1 000 / 10 000 × 2 cœurs")

	select {
	case <-signals:
	default:
		t.Fatal("un abonné est prévenu du passage")
	}

	history, found := r.ContainerHistory(idA, time.Time{}, clock.Add(time.Hour))
	require.True(t, found)
	assert.Len(t, history, 2)
	assert.Equal(t, []bool{true, false}, docker.previous[idA],
		"seule la première mesure demande un relevé précédent au démon")
}

func TestRecorder_ResolvesUniquePrefixesOnly(t *testing.T) {
	docker := newFakeDocker(idA, idB)
	clock := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	r := newTestRecorder(docker, &fakeHost{}, &clock)
	r.recordContainers(context.Background())

	_, ok := r.LatestContainer(idA[:12])
	assert.True(t, ok, "préfixe de 12 caractères")
	_, ok = r.LatestContainer(idA[:11])
	assert.False(t, ok, "préfixe trop court")

	docker.running = []string{idA, idA2, idB}
	r.recordContainers(context.Background())
	_, ok = r.LatestContainer(idA[:12])
	assert.False(t, ok, "préfixe ambigu")
	_, ok = r.LatestContainer(idA2)
	assert.True(t, ok, "identifiant complet")
}

func TestRecorder_KeepsStoppedContainersForTheWindow(t *testing.T) {
	docker := newFakeDocker(idA, idB)
	clock := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	r := newTestRecorder(docker, &fakeHost{}, &clock)
	r.recordContainers(context.Background())

	docker.running = []string{idB}
	clock = clock.Add(30 * time.Minute)
	r.recordContainers(context.Background())
	assert.Len(t, r.LatestSweep().Containers, 1, "le conteneur arrêté n'est plus dans le passage")
	_, found := r.ContainerHistory(idA, time.Time{}, clock)
	assert.True(t, found, "son historique reste lisible")

	clock = clock.Add(31 * time.Minute)
	r.recordContainers(context.Background())
	_, found = r.ContainerHistory(idA, time.Time{}, clock)
	assert.False(t, found, "au-delà de la fenêtre, il est oublié")
}

func TestRecorder_RestartedContainerGetsAFreshCollector(t *testing.T) {
	docker := newFakeDocker(idA)
	clock := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	r := newTestRecorder(docker, &fakeHost{}, &clock)
	r.recordContainers(context.Background())

	docker.running = nil
	clock = clock.Add(10 * time.Second)
	r.recordContainers(context.Background())

	docker.running = []string{idA}
	clock = clock.Add(10 * time.Second)
	r.recordContainers(context.Background())
	assert.Equal(t, []bool{true, true}, docker.previous[idA],
		"après un arrêt, le relevé processeur d'avant ne sert plus de référence")
}

func TestRecorder_SkipsContainersThatVanish(t *testing.T) {
	docker := newFakeDocker(idA, idB)
	docker.statsErr[idB] = orchestrator.ErrNotFound
	clock := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	r := newTestRecorder(docker, &fakeHost{}, &clock)
	r.recordContainers(context.Background())

	sweep := r.LatestSweep()
	require.Len(t, sweep.Containers, 1)
	assert.Equal(t, idA, sweep.Containers[0].ContainerID)
}

func TestRecorder_DockerDownKeepsTheLastSweep(t *testing.T) {
	docker := newFakeDocker(idA)
	clock := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	r := newTestRecorder(docker, &fakeHost{}, &clock)
	r.recordContainers(context.Background())
	first := r.LatestSweep().At

	docker.listErr = errors.New("démon injoignable")
	clock = clock.Add(10 * time.Second)
	r.recordContainers(context.Background())
	assert.Equal(t, first, r.LatestSweep().At, "pas de faux passage vide pendant la panne")
}

func TestRecorder_CapsTheNumberOfContainers(t *testing.T) {
	docker := newFakeDocker(idB, idA, idA2)
	clock := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	r := NewRecorder(&fakeHost{}, docker, docker, RecorderConfig{MaxContainers: 2}, zerolog.Nop())
	r.now = func() time.Time { return clock }
	r.recordContainers(context.Background())

	sweep := r.LatestSweep()
	require.Len(t, sweep.Containers, 2)
	assert.Equal(t, idA, sweep.Containers[0].ContainerID)
	assert.Equal(t, idA2, sweep.Containers[1].ContainerID)
}

func TestRecorder_HostHistoryFiltersByTime(t *testing.T) {
	host := &fakeHost{}
	clock := time.Now()
	r := newTestRecorder(newFakeDocker(), host, &clock)
	_, ok := r.LatestHost()
	assert.False(t, ok)

	for range 4 {
		r.recordHost(context.Background())
	}
	start := time.Date(2026, 11, 14, 12, 0, 10, 0, time.UTC)
	got := r.HostHistory(start, start.Add(5*time.Second))
	require.Len(t, got, 2, "bornes incluses : 12:00:10 et 12:00:15")

	host.err = errors.New("/proc illisible")
	r.recordHost(context.Background())
	latest, ok := r.LatestHost()
	require.True(t, ok)
	assert.Equal(t, start.Add(10*time.Second), latest.CollectedAt, "une mesure ratée laisse un trou, pas une valeur fausse")
}

func TestRecorder_StartAndShutdown(t *testing.T) {
	docker := newFakeDocker(idA)
	r := NewRecorder(&fakeHost{}, docker, docker,
		RecorderConfig{HostEvery: 10 * time.Millisecond, ContainerEvery: 10 * time.Millisecond}, zerolog.Nop())

	require.NoError(t, r.Shutdown(context.Background()), "arrêt sans démarrage")

	r = NewRecorder(&fakeHost{}, docker, docker,
		RecorderConfig{HostEvery: 10 * time.Millisecond, ContainerEvery: 10 * time.Millisecond}, zerolog.Nop())
	done := make(chan error, 1)
	go func() { done <- r.Start() }()
	require.Eventually(t, func() bool {
		_, host := r.LatestHost()
		return host && len(r.LatestSweep().Containers) == 1
	}, 2*time.Second, 5*time.Millisecond, "le premier relevé est immédiat")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	require.NoError(t, r.Shutdown(ctx))
	require.NoError(t, <-done)
}

// Régression trouvée par le test réel : avec un relevé toutes les secondes,
// la première mesure (une seconde d'attente côté démon) dépassait le délai
// du passage, et aucun conteneur n'était jamais mesuré.
func TestRecorder_SlowDaemonStillMeasures(t *testing.T) {
	docker := newFakeDocker(idA)
	docker.delay = 1200 * time.Millisecond
	clock := time.Date(2026, 11, 14, 12, 0, 0, 0, time.UTC)
	r := NewRecorder(&fakeHost{}, docker, docker, RecorderConfig{ContainerEvery: time.Second}, zerolog.Nop())
	r.now = func() time.Time { return clock }

	r.recordContainers(context.Background())
	assert.Len(t, r.LatestSweep().Containers, 1, "une réponse plus lente que le rythme du relevé est attendue")
}
