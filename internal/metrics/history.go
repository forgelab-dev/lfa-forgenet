package metrics

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
)

// Historique court (ADR-0010) : le serveur relève lui-même l'hôte et chaque
// conteneur en cours d'exécution, et garde la dernière heure en mémoire.
// Rien n'est écrit sur le disque ; l'historique repart de zéro à chaque
// démarrage du serveur.
const (
	DefaultHistoryWindow  = time.Hour
	DefaultHostEvery      = 5 * time.Second
	DefaultContainerEvery = 10 * time.Second
	// DefaultMaxContainers borne la mémoire : au-delà, les conteneurs en trop
	// (dans l'ordre de leurs identifiants) ne sont pas suivis.
	DefaultMaxContainers = 200
	// DefaultParallelStats limite les appels simultanés au démon Docker.
	DefaultParallelStats = 4

	// sweepSlack s'ajoute au délai d'un passage : la première mesure d'un
	// conteneur attend une seconde (relevé précédent demandé au démon), et
	// un démon chargé répond lentement. Les passages ne se chevauchent pas.
	sweepSlack = 5 * time.Second
)

// RunningLister liste les conteneurs en cours d'exécution.
type RunningLister interface {
	RunningIDs(ctx context.Context) ([]string, error)
}

// HostSampler mesure l'hôte (HostCollector le satisfait).
type HostSampler interface {
	Collect(ctx context.Context) (Host, error)
}

// RecorderConfig règle le relevé de fond ; une valeur nulle prend la valeur
// par défaut.
type RecorderConfig struct {
	Window         time.Duration
	HostEvery      time.Duration
	ContainerEvery time.Duration
	MaxContainers  int
	Parallel       int
}

func (c RecorderConfig) withDefaults() RecorderConfig {
	if c.Window <= 0 {
		c.Window = DefaultHistoryWindow
	}
	if c.HostEvery <= 0 {
		c.HostEvery = DefaultHostEvery
	}
	if c.ContainerEvery <= 0 {
		c.ContainerEvery = DefaultContainerEvery
	}
	if c.MaxContainers <= 0 {
		c.MaxContainers = DefaultMaxContainers
	}
	if c.Parallel <= 0 {
		c.Parallel = DefaultParallelStats
	}
	return c
}

// Sweep est un passage sur tous les conteneurs en cours d'exécution.
type Sweep struct {
	At         time.Time
	Containers []Container // triés par identifiant
}

type containerTrack struct {
	collector *ContainerCollector
	samples   *ring[Container]
	// lastSeen : dernière mesure réussie ; lastListed : dernier passage où le
	// démon le donnait en cours d'exécution, mesure réussie ou non.
	lastSeen, lastListed time.Time
}

// Recorder relève l'hôte et les conteneurs en tâche de fond et garde leur
// historique récent. Il suit le cycle de vie Start/Shutdown des serveurs.
type Recorder struct {
	cfg    RecorderConfig
	host   HostSampler
	stats  StatsSource
	lister RunningLister
	logger zerolog.Logger
	now    func() time.Time

	mu          sync.RWMutex
	hostSamples *ring[Host]
	containers  map[string]*containerTrack
	sweep       Sweep
	lastListing time.Time // dernière liste des conteneurs obtenue du démon
	subscribers map[chan time.Time]struct{}

	hostFailing, containerFailing failing

	stop     chan struct{}
	stopOnce sync.Once
	started  chan struct{}
	done     chan struct{}
}

// NewRecorder crée le relevé de fond.
func NewRecorder(host HostSampler, stats StatsSource, lister RunningLister, cfg RecorderConfig, logger zerolog.Logger) *Recorder {
	cfg = cfg.withDefaults()
	return &Recorder{
		cfg:         cfg,
		host:        host,
		stats:       stats,
		lister:      lister,
		logger:      logger.With().Str("component", "metrics.recorder").Logger(),
		now:         time.Now,
		hostSamples: newRing[Host](capacity(cfg.Window, cfg.HostEvery)),
		containers:  make(map[string]*containerTrack),
		subscribers: make(map[chan time.Time]struct{}),
		stop:        make(chan struct{}),
		started:     make(chan struct{}),
		done:        make(chan struct{}),
	}
}

// capacity : de quoi couvrir la fenêtre, plus un point pour ses deux bords.
func capacity(window, every time.Duration) int {
	return int(window/every) + 1
}

// Config renvoie les réglages effectifs, valeurs par défaut comprises.
func (r *Recorder) Config() RecorderConfig { return r.cfg }

// Start relève jusqu'à Shutdown. Le premier relevé est immédiat.
func (r *Recorder) Start() error {
	close(r.started)
	defer close(r.done)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		select {
		case <-r.stop:
			cancel()
		case <-ctx.Done():
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); r.every(ctx, r.cfg.HostEvery, r.recordHost) }()
	go func() { defer wg.Done(); r.every(ctx, r.cfg.ContainerEvery, r.recordContainers) }()
	wg.Wait()
	return nil
}

// Shutdown arrête le relevé et attend la fin du passage en cours.
func (r *Recorder) Shutdown(ctx context.Context) error {
	r.stopOnce.Do(func() { close(r.stop) })
	select {
	case <-r.started:
	default:
		return nil // Start n'a jamais tourné : rien à attendre.
	}
	select {
	case <-r.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *Recorder) every(ctx context.Context, interval time.Duration, record func(context.Context)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		record(ctx)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

// failing ne journalise qu'au changement d'état : une panne qui dure ne
// remplit pas le journal toutes les cinq secondes.
type failing struct {
	mu   sync.Mutex
	down bool
}

func (f *failing) set(down bool) (changed bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	changed = f.down != down
	f.down = down
	return changed
}

func (r *Recorder) recordHost(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.HostEvery)
	defer cancel()
	m, err := r.host.Collect(ctx)
	if err != nil {
		if ctx.Err() == nil && r.hostFailing.set(true) {
			r.logger.Warn().Err(err).Msg("relevé de l'hôte impossible, l'historique a un trou")
		}
		return
	}
	if r.hostFailing.set(false) {
		r.logger.Info().Msg("relevé de l'hôte rétabli")
	}
	r.mu.Lock()
	r.hostSamples.push(m)
	r.mu.Unlock()
}

func (r *Recorder) recordContainers(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, r.cfg.ContainerEvery+sweepSlack)
	defer cancel()
	ids, err := r.lister.RunningIDs(ctx)
	if err != nil {
		if ctx.Err() == nil && r.containerFailing.set(true) {
			r.logger.Warn().Err(err).Msg("liste des conteneurs impossible, relevé suspendu")
		}
		return
	}
	if r.containerFailing.set(false) {
		r.logger.Info().Msg("relevé des conteneurs rétabli")
	}
	slices.Sort(ids)
	if len(ids) > r.cfg.MaxContainers {
		r.logger.Warn().Int("en_cours", len(ids)).Int("suivis", r.cfg.MaxContainers).
			Msg("trop de conteneurs, les derniers ne sont pas suivis")
		ids = ids[:r.cfg.MaxContainers]
	}

	collectors := r.collectorsFor(ids, r.now())
	results := make([]*Container, len(ids))
	slots := make(chan struct{}, r.cfg.Parallel)
	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Add(1)
		slots <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-slots }()
			m, err := collectors[i].Collect(ctx)
			switch {
			case err == nil:
				results[i] = &m
			case errors.Is(err, orchestrator.ErrNotFound), ctx.Err() != nil:
				// Arrêté entre la liste et la mesure, ou passage trop long.
			default:
				r.logger.Debug().Err(err).Str("container_id", id).Msg("mesure d'un conteneur impossible")
			}
		}()
	}
	wg.Wait()
	if ctx.Err() != nil && errors.Is(ctx.Err(), context.Canceled) {
		return // arrêt du serveur : pas de demi-passage
	}

	now := r.now()
	sweep := Sweep{At: now, Containers: make([]Container, 0, len(ids))}
	r.mu.Lock()
	for i, id := range ids {
		if results[i] == nil {
			continue
		}
		t := r.containers[id]
		t.samples.push(*results[i])
		t.lastSeen = now
		sweep.Containers = append(sweep.Containers, *results[i])
	}
	// Un conteneur arrêté garde son historique tant qu'il reste dans la
	// fenêtre (lastListed est toujours au moins aussi récent que lastSeen).
	for id, t := range r.containers {
		if now.Sub(t.lastListed) > r.cfg.Window {
			delete(r.containers, id)
		}
	}
	r.sweep = sweep
	for ch := range r.subscribers {
		select {
		case ch <- now:
		default: // l'abonné n'a pas lu le passage précédent : il lira le dernier
		}
	}
	r.mu.Unlock()
}

// collectorsFor renvoie un collecteur par conteneur. Un conteneur absent de
// la liste précédente (arrêté puis relancé) repart d'un nouveau collecteur :
// son relevé processeur d'avant l'arrêt fausserait la première mesure. Une
// mesure ratée d'un conteneur toujours listé garde le collecteur et son
// dernier relevé : la mesure suivante n'a pas à attendre une seconde de plus.
func (r *Recorder) collectorsFor(ids []string, now time.Time) []*ContainerCollector {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*ContainerCollector, len(ids))
	for i, id := range ids {
		t, ok := r.containers[id]
		if !ok {
			t = &containerTrack{samples: newRing[Container](capacity(r.cfg.Window, r.cfg.ContainerEvery))}
			r.containers[id] = t
		}
		if t.collector == nil || t.lastListed.Before(r.lastListing) {
			t.collector = NewContainerCollector(r.stats, id)
		}
		t.lastListed = now
		out[i] = t.collector
	}
	r.lastListing = now
	return out
}

// LatestHost renvoie la dernière mesure de l'hôte.
func (r *Recorder) LatestHost() (Host, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.hostSamples.last()
}

// HostHistory renvoie les mesures de l'hôte prises entre from et to inclus,
// de la plus ancienne à la plus récente.
func (r *Recorder) HostHistory(from, to time.Time) []Host {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []Host
	r.hostSamples.each(func(h Host) {
		if !h.CollectedAt.Before(from) && !h.CollectedAt.After(to) {
			out = append(out, h)
		}
	})
	return out
}

// ContainerHistory renvoie les mesures d'un conteneur prises entre from et
// to. id est l'identifiant complet, ou un préfixe d'au moins 12 caractères
// qui n'en désigne qu'un. found est faux si le conteneur n'a pas
// d'historique (inconnu, ou arrêté depuis plus longtemps que la fenêtre).
func (r *Recorder) ContainerHistory(id string, from, to time.Time) (samples []Container, found bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t := r.track(id)
	if t == nil {
		return nil, false
	}
	t.samples.each(func(c Container) {
		if !c.CollectedAt.Before(from) && !c.CollectedAt.After(to) {
			samples = append(samples, c)
		}
	})
	return samples, true
}

// LatestContainer renvoie la dernière mesure d'un conteneur suivi.
func (r *Recorder) LatestContainer(id string) (Container, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t := r.track(id)
	if t == nil {
		return Container{}, false
	}
	return t.samples.last()
}

// MinPrefix est la longueur minimale d'un identifiant abrégé, celle qu'affiche
// `docker ps`.
const MinPrefix = 12

func (r *Recorder) track(id string) *containerTrack {
	if t, ok := r.containers[id]; ok {
		return t
	}
	if len(id) < MinPrefix {
		return nil
	}
	var match *containerTrack
	for full, t := range r.containers {
		if strings.HasPrefix(full, id) {
			if match != nil {
				return nil // préfixe ambigu
			}
			match = t
		}
	}
	return match
}

// LatestSweep renvoie le dernier passage sur les conteneurs (At est nul tant
// qu'aucun passage n'a abouti).
func (r *Recorder) LatestSweep() Sweep {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return Sweep{At: r.sweep.At, Containers: slices.Clone(r.sweep.Containers)}
}

// Subscribe prévient à chaque passage sur les conteneurs. Un abonné lent ne
// retient pas le relevé : il reçoit au plus un signal en attente, puis lit
// LatestSweep. cancel libère l'abonnement.
func (r *Recorder) Subscribe() (signals <-chan time.Time, cancel func()) {
	ch := make(chan time.Time, 1)
	r.mu.Lock()
	r.subscribers[ch] = struct{}{}
	r.mu.Unlock()
	var once sync.Once
	return ch, func() {
		once.Do(func() {
			r.mu.Lock()
			delete(r.subscribers, ch)
			r.mu.Unlock()
		})
	}
}

// ring est un tampon circulaire de taille fixe : le plus ancien élément est
// remplacé quand il est plein.
type ring[T any] struct {
	items []T
	start int
	n     int
}

func newRing[T any](size int) *ring[T] {
	return &ring[T]{items: make([]T, size)}
}

func (r *ring[T]) push(v T) {
	if r.n < len(r.items) {
		r.items[(r.start+r.n)%len(r.items)] = v
		r.n++
		return
	}
	r.items[r.start] = v
	r.start = (r.start + 1) % len(r.items)
}

// each parcourt les éléments du plus ancien au plus récent.
func (r *ring[T]) each(fn func(T)) {
	for i := range r.n {
		fn(r.items[(r.start+i)%len(r.items)])
	}
}

func (r *ring[T]) last() (T, bool) {
	if r.n == 0 {
		var zero T
		return zero, false
	}
	return r.items[(r.start+r.n-1)%len(r.items)], true
}
