package update

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Latest donne la dernière version publiée ; Source le fait par l'API GitHub.
type Latest interface {
	Latest(ctx context.Context, withPre bool) (Release, error)
}

// Checker vérifie régulièrement si une version plus récente est publiée, et
// garde le résultat en mémoire pour le dashboard (ADR-0011). Il ne met rien à
// jour lui-même.
type Checker struct {
	source     Latest
	current    string
	interval   time.Duration
	firstDelay time.Duration
	minGap     time.Duration
	now        func() time.Time
	logger     zerolog.Logger

	mu        sync.Mutex
	cancel    context.CancelFunc
	stopped   bool
	done      chan struct{}
	latest    *Release
	checkedAt time.Time
	attempted time.Time
	lastErr   string
}

// Snapshot est l'état connu du Checker.
type Snapshot struct {
	// Available est la version publiée la plus récente, si elle dépasse la
	// version en service.
	Available *Release
	CheckedAt time.Time // dernière vérification réussie ; zéro si aucune
	Err       string    // erreur de la dernière vérification
}

// NewChecker prépare la vérification ; interval vaut 24 h s'il est nul.
func NewChecker(source Latest, current string, interval time.Duration, logger zerolog.Logger) *Checker {
	if interval <= 0 {
		interval = 24 * time.Hour
	}
	return &Checker{
		source: source, current: current, interval: interval, firstDelay: 30 * time.Second, minGap: time.Minute,
		now: time.Now, logger: logger.With().Str("component", "update.checker").Logger(),
	}
}

// Start vérifie peu après le démarrage, puis à chaque intervalle, jusqu'à
// Shutdown. Il bloque, comme les autres tâches de fond du serveur.
func (c *Checker) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	c.mu.Lock()
	if c.stopped { // Shutdown avant Start : rien à lancer
		c.mu.Unlock()
		cancel()
		return nil
	}
	c.cancel, c.done = cancel, make(chan struct{})
	done := c.done
	c.mu.Unlock()
	defer close(done)

	timer := time.NewTimer(c.firstDelay)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			c.Check(ctx)
			timer.Reset(c.interval)
		}
	}
}

// Shutdown arrête les vérifications ; une vérification en cours est annulée.
func (c *Checker) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	cancel, done := c.cancel, c.done
	c.stopped = true
	c.mu.Unlock()
	if cancel == nil {
		return nil
	}
	cancel()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Check interroge les releases. Appelé à moins d'une minute du précédent
// essai, il ne refait pas l'appel : le dashboard ne peut pas faire marteler
// l'API GitHub.
func (c *Checker) Check(ctx context.Context) Snapshot {
	c.mu.Lock()
	if !c.attempted.IsZero() && c.now().Sub(c.attempted) < c.minGap {
		defer c.mu.Unlock()
		return c.snapshotLocked()
	}
	c.attempted = c.now()
	c.mu.Unlock()

	rctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	rel, err := c.source.Latest(rctx, false)

	c.mu.Lock()
	defer c.mu.Unlock()
	if err != nil {
		c.lastErr = err.Error()
		c.logger.Warn().Err(err).Msg("vérification des mises à jour impossible")
		return c.snapshotLocked()
	}
	c.latest, c.checkedAt, c.lastErr = &rel, c.now(), ""
	if s := c.snapshotLocked(); s.Available != nil {
		c.logger.Info().Str("version", rel.Version).Msg("nouvelle version de ForgeNet disponible")
	}
	return c.snapshotLocked()
}

// Snapshot renvoie l'état connu, sans appel réseau.
func (c *Checker) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.snapshotLocked()
}

func (c *Checker) snapshotLocked() Snapshot {
	s := Snapshot{CheckedAt: c.checkedAt, Err: c.lastErr}
	if c.latest == nil {
		return s
	}
	latest, _ := ParseVersion(c.latest.Version)
	cur, err := ParseVersion(c.current)
	// Une version en service illisible (build de développement) n'annonce
	// rien : on ne sait pas la comparer.
	if err == nil && latest.Compare(cur) > 0 {
		rel := *c.latest
		s.Available = &rel
	}
	return s
}
