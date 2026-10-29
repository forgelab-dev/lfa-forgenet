// Package deploy enregistre les déploiements, les exécute dans l'ordre de leur
// création et les fait avancer dans leur machine à états (ADR-0004).
//
// Une stack Compose par projet, nommée « forgenet-<projet> ». Un seul
// déploiement s'exécute à la fois : deux mises à jour d'une même stack ne se
// croisent jamais, et la file reste dans la base, ce qui survit à un redémarrage.
package deploy

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/ids"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

var (
	// ErrInvalid indique une demande incomplète ou un fichier refusé.
	ErrInvalid = errors.New("demande de déploiement invalide")
	// ErrUnsupported indique une source que ForgeNet ne sait pas encore déployer.
	ErrUnsupported = errors.New("source de déploiement pas encore prise en charge")
	// ErrNotCancellable indique un déploiement qui ne peut plus être annulé.
	ErrNotCancellable = errors.New("déploiement non annulable dans son état actuel")
	// ErrBusy indique qu'un déploiement du projet est en attente ou en cours.
	ErrBusy = errors.New("un déploiement du projet est en attente ou en cours")
)

// Stacks est la part du moteur Compose dont le service a besoin.
type Stacks interface {
	Validate(ctx context.Context, content []byte) error
	Deploy(ctx context.Context, project string, content []byte) error
	Remove(ctx context.Context, project string, removeVolumes bool) error
	Services(ctx context.Context, project string) ([]compose.Service, error)
}

const (
	// MaxComposeSize borne la taille d'un fichier Compose reçu.
	MaxComposeSize = 1 << 20
	// stackPrefix distingue les stacks de ForgeNet des projets Compose de l'hôte.
	stackPrefix = "forgenet-"

	defaultTimeout      = 15 * time.Minute
	defaultPollInterval = 5 * time.Second
	// batchSize borne le nombre de déploiements lus à chaque passage de la file.
	batchSize = 20
	// writeTimeout borne l'écriture d'un état, faite même après une annulation.
	writeTimeout = 10 * time.Second
	// cancelWait borne l'attente de l'arrêt d'un déploiement annulé.
	cancelWait = 30 * time.Second

	interruptedMessage = "interrompu par l'arrêt du serveur"
)

// projectIDPattern garde des noms de stack valides pour Compose une fois préfixés.
var projectIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,62}$`)

// StackName renvoie le nom du projet Compose d'un projet ForgeNet.
func StackName(projectID string) string {
	return stackPrefix + projectID
}

// CreateRequest décrit un déploiement à créer.
type CreateRequest struct {
	ProjectID   string
	SourceType  store.DeploymentSourceType
	Content     string // fichier Compose
	TargetEnv   string
	RequestID   string // clé d'idempotence, facultative
	TriggeredBy string // ex: "user:usr_…"
}

// Service enregistre et exécute les déploiements. Start et Shutdown en font un
// serveur comme les autres pour le binaire.
type Service struct {
	store   store.DeploymentStore
	stacks  Stacks
	auditor audit.Auditor
	logger  zerolog.Logger
	timeout time.Duration
	poll    time.Duration

	wake chan struct{}

	mu      sync.Mutex
	busy    map[string]bool       // projets en cours de déploiement ou de suppression
	running map[string]*execution // déploiements en cours, par identifiant

	// stop arrête la prise de nouveaux déploiements ; kill interrompt celui en
	// cours, une fois le délai de grâce écoulé. Shutdown les ferme.
	stop     chan struct{}
	kill     chan struct{}
	stopOnce sync.Once
	killOnce sync.Once
	started  chan struct{}
	done     chan struct{}
}

type execution struct {
	cancel context.CancelFunc
	reason string // motif d'annulation, vide si personne n'a annulé
}

// Option règle le service.
type Option func(*Service)

// WithTimeout borne la durée d'un déploiement.
func WithTimeout(d time.Duration) Option { return func(s *Service) { s.timeout = d } }

// WithPollInterval règle la relecture de la file quand rien ne la réveille.
func WithPollInterval(d time.Duration) Option { return func(s *Service) { s.poll = d } }

// New crée le service ; Start lance l'exécution de la file.
func New(st store.DeploymentStore, stacks Stacks, auditor audit.Auditor, logger zerolog.Logger, opts ...Option) *Service {
	s := &Service{
		store:   st,
		stacks:  stacks,
		auditor: auditor,
		logger:  logger.With().Str("component", "deploy").Logger(),
		timeout: defaultTimeout,
		poll:    defaultPollInterval,
		wake:    make(chan struct{}, 1),
		busy:    make(map[string]bool),
		running: make(map[string]*execution),
		stop:    make(chan struct{}),
		kill:    make(chan struct{}),
		started: make(chan struct{}),
		done:    make(chan struct{}),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Create vérifie et enregistre un déploiement, en attente, puis réveille la
// file. Avec une clé d'idempotence déjà vue, il renvoie le déploiement existant.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*store.Deployment, error) {
	d, err := s.create(ctx, req)
	meta := map[string]any{"success": err == nil, "project_id": req.ProjectID, "source": string(req.SourceType)}
	resource := "project:" + req.ProjectID
	if err != nil {
		meta["error"] = err.Error()
	} else {
		resource = "deployment:" + d.ID
	}
	s.audit(ctx, "deployment.create", resource, meta)
	return d, err
}

func (s *Service) create(ctx context.Context, req CreateRequest) (*store.Deployment, error) {
	if err := validateProjectID(req.ProjectID); err != nil {
		return nil, err
	}
	if req.SourceType != store.SourceCompose {
		return nil, fmt.Errorf("%w : %s", ErrUnsupported, req.SourceType)
	}
	if strings.TrimSpace(req.Content) == "" {
		return nil, fmt.Errorf("%w : fichier Compose vide", ErrInvalid)
	}
	if len(req.Content) > MaxComposeSize {
		return nil, fmt.Errorf("%w : fichier Compose de plus de %d octets", ErrInvalid, MaxComposeSize)
	}

	if existing, err := s.byRequestID(ctx, req); existing != nil || err != nil {
		return existing, err
	}

	if err := s.stacks.Validate(ctx, []byte(req.Content)); err != nil {
		if errors.Is(err, compose.ErrInvalidFile) {
			return nil, fmt.Errorf("%w : %w", ErrInvalid, err)
		}
		return nil, err
	}

	// Le déploiement précédent du projet trace la chaîne des versions.
	previous, err := s.store.ListDeployments(ctx, req.ProjectID, 1, 0)
	if err != nil {
		return nil, err
	}
	d := &store.Deployment{
		ID:            ids.New("dep"),
		ProjectID:     req.ProjectID,
		SourceType:    req.SourceType,
		SourceContent: req.Content,
		State:         store.DeploymentPending,
		TargetEnv:     req.TargetEnv,
		TriggeredBy:   req.TriggeredBy,
		RequestID:     req.RequestID,
	}
	if len(previous) > 0 {
		d.PreviousDeploymentID = previous[0].ID
	}

	if err := s.store.CreateDeployment(ctx, d); err != nil {
		// Deux appels avec la même clé en même temps : le second rend le premier.
		if errors.Is(err, store.ErrConflict) && req.RequestID != "" {
			return s.byRequestID(ctx, req)
		}
		return nil, err
	}
	s.notify()
	return d, nil
}

// byRequestID renvoie le déploiement déjà créé avec la clé de req, ou nil.
func (s *Service) byRequestID(ctx context.Context, req CreateRequest) (*store.Deployment, error) {
	if req.RequestID == "" {
		return nil, nil
	}
	existing, err := s.store.GetDeploymentByRequestID(ctx, req.RequestID)
	if errors.Is(err, store.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if existing.ProjectID != req.ProjectID {
		return nil, fmt.Errorf("%w : clé d'idempotence déjà utilisée pour un autre projet", ErrInvalid)
	}
	return existing, nil
}

// Get renvoie un déploiement.
func (s *Service) Get(ctx context.Context, id string) (*store.Deployment, error) {
	return s.store.GetDeployment(ctx, id)
}

// List liste les déploiements d'un projet (tous si projectID est vide), du plus récent au plus ancien.
func (s *Service) List(ctx context.Context, projectID string, limit, offset int) ([]*store.Deployment, error) {
	return s.store.ListDeployments(ctx, projectID, limit, offset)
}

// Cancel annule un déploiement en attente, ou interrompt un déploiement en
// construction ou en cours de déploiement, et attend qu'il soit arrêté.
// Interrompre `docker compose up` peut laisser la stack à moitié mise à jour :
// un nouveau déploiement la remet d'aplomb.
func (s *Service) Cancel(ctx context.Context, id, reason string) (*store.Deployment, error) {
	d, err := s.cancel(ctx, id, reason)
	meta := map[string]any{"success": err == nil, "reason": reason}
	if err != nil {
		meta["error"] = err.Error()
	}
	s.audit(ctx, "deployment.cancel", "deployment:"+id, meta)
	return d, err
}

func (s *Service) cancel(ctx context.Context, id, reason string) (*store.Deployment, error) {
	message := "annulé"
	if reason = strings.TrimSpace(reason); reason != "" {
		message += " : " + reason
	}

	// Trois essais : l'exécuteur peut prendre le déploiement entre la lecture
	// de son état et l'annulation.
	for range 3 {
		d, err := s.store.GetDeployment(ctx, id)
		if err != nil {
			return nil, err
		}
		switch d.State {
		case store.DeploymentPending:
			cancelled, err := s.store.TransitionDeployment(ctx, id, store.DeploymentPending, store.DeploymentCancelled, message)
			if errors.Is(err, store.ErrConflict) {
				continue
			}
			return cancelled, err
		case store.DeploymentBuilding, store.DeploymentDeploying:
			if !s.interrupt(id, message) {
				continue
			}
			return s.waitFinished(ctx, id)
		default:
			return nil, fmt.Errorf("%w : %s", ErrNotCancellable, d.State)
		}
	}
	return nil, fmt.Errorf("%w : l'état change trop vite, réessayez", ErrNotCancellable)
}

// interrupt annule l'exécution en cours de id ; false si elle n'est pas (ou plus) en cours.
func (s *Service) interrupt(id, reason string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.running[id]
	if !ok {
		return false
	}
	run.reason = reason
	run.cancel()
	return true
}

// waitFinished attend qu'un déploiement interrompu atteigne un état final.
func (s *Service) waitFinished(ctx context.Context, id string) (*store.Deployment, error) {
	ctx, cancel := context.WithTimeout(ctx, cancelWait)
	defer cancel()
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		d, err := s.store.GetDeployment(ctx, id)
		if err != nil {
			return nil, err
		}
		if d.State.IsTerminal() {
			return d, nil
		}
		select {
		case <-ctx.Done():
			return d, nil
		case <-ticker.C:
		}
	}
}

// DeleteStack supprime la stack en service d'un projet. Il est refusé tant
// qu'un déploiement du projet attend ou s'exécute : sans cela, ce déploiement
// recréerait la stack juste après sa suppression.
func (s *Service) DeleteStack(ctx context.Context, projectID string, removeVolumes bool) error {
	err := s.deleteStack(ctx, projectID, removeVolumes)
	meta := map[string]any{"success": err == nil, "delete_volumes": removeVolumes}
	if err != nil {
		meta["error"] = err.Error()
	}
	s.audit(ctx, "stack.delete", "project:"+projectID, meta)
	return err
}

func (s *Service) deleteStack(ctx context.Context, projectID string, removeVolumes bool) error {
	if err := validateProjectID(projectID); err != nil {
		return err
	}
	if !s.acquire(projectID) {
		return ErrBusy
	}
	defer s.release(projectID)

	pending, err := s.store.ListDeploymentsByState(ctx, store.DeploymentPending, 0)
	if err != nil {
		return err
	}
	for _, d := range pending {
		if d.ProjectID == projectID {
			return fmt.Errorf("%w : %s", ErrBusy, d.ID)
		}
	}
	return s.stacks.Remove(ctx, StackName(projectID), removeVolumes)
}

func validateProjectID(id string) error {
	if !projectIDPattern.MatchString(id) {
		return fmt.Errorf("%w : identifiant de projet %q (minuscules, chiffres, « _ » et « - », 63 caractères au plus)", ErrInvalid, id)
	}
	return nil
}

func (s *Service) acquire(projectID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.busy[projectID] {
		return false
	}
	s.busy[projectID] = true
	return true
}

func (s *Service) release(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.busy, projectID)
}

func (s *Service) notify() {
	select {
	case s.wake <- struct{}{}:
	default:
	}
}

func (s *Service) audit(ctx context.Context, action, resource string, meta map[string]any) {
	if s.auditor != nil {
		_ = s.auditor.Log(ctx, action, resource, meta)
	}
}

// checkHealthy vérifie qu'après le déploiement chaque conteneur tourne, et que
// ceux qui ont une sonde de santé sont sains.
func checkHealthy(services []compose.Service) error {
	if len(services) == 0 {
		return errors.New("aucun conteneur après le déploiement")
	}
	for _, svc := range services {
		if svc.State != orchestrator.StateRunning {
			return fmt.Errorf("service %s : conteneur %s", svc.Service, svc.State)
		}
		if svc.Health != "" && svc.Health != "healthy" {
			return fmt.Errorf("service %s : santé %s", svc.Service, svc.Health)
		}
	}
	return nil
}
