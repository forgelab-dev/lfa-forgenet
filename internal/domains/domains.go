// Package domains relie des noms de domaine aux services des applications
// déployées : il les enregistre, branche les conteneurs sur le réseau partagé
// avec Caddy et fait servir les routes par le reverse proxy (ADR-0007,
// ADR-0008).
//
// La base est la source de vérité. Une modification ne laisse jamais de
// demi-état, et toutes les opérations (ajout, retrait, rebranchement après un
// déploiement, réconciliation au démarrage) passent par un seul verrou.
package domains

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/net/idna"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/deploy"
	"github.com/forgelab-dev/lfa-forgenet/internal/ids"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
	"github.com/forgelab-dev/lfa-forgenet/internal/proxy"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

var (
	// ErrInvalid indique une demande incomplète ou une valeur refusée.
	ErrInvalid = errors.New("demande de domaine invalide")
	// ErrNotDeployed indique un service sans conteneur dans la stack du projet.
	ErrNotDeployed = errors.New("service non déployé")
)

// NetworkLabels marquent le réseau partagé créé par ForgeNet (ADR-0008).
var NetworkLabels = map[string]string{"forgenet.managed": "edge"}

const (
	// aliasSuffix termine l'alias DNS d'un service : « .internal » est réservé
	// aux réseaux privés et ne peut désigner aucun domaine public.
	aliasSuffix = ".forgenet.internal"
	// rollbackTimeout borne l'annulation d'un ajout, faite même si l'appel est annulé.
	rollbackTimeout = 30 * time.Second
)

// dnsLabel suit la règle des noms que résout Go (donc Caddy) : lettres
// minuscules, chiffres, « _ » et « - », sans tiret au bord, 63 caractères au plus.
var dnsLabel = regexp.MustCompile(`^[a-z0-9_]([a-z0-9_-]{0,61}[a-z0-9_])?$`)

// Proxy applique les routes ; *proxy.Manager en est l'implémentation.
type Proxy interface {
	Apply(ctx context.Context, routes []proxy.Route) error
}

// Network branche les conteneurs ; *docker.Engine en est l'implémentation.
type Network interface {
	EnsureNetwork(ctx context.Context, name string, labels map[string]string) error
	ConnectToNetwork(ctx context.Context, network, containerID string, aliases []string) error
	DisconnectFromNetwork(ctx context.Context, network, containerID string) error
}

// Stacks lit les conteneurs d'une stack ; *compose.Engine en est l'implémentation.
type Stacks interface {
	Services(ctx context.Context, project string) ([]compose.Service, error)
}

// Config règle le branchement.
type Config struct {
	// Network est le réseau partagé avec Caddy.
	Network string
	// CaddyContainer, s'il est fixé, est branché sur ce réseau à la réconciliation.
	CaddyContainer string
	// DashboardHost est le nom de domaine du dashboard (ADR-0009) : réservé,
	// il ne peut pas être relié à une application.
	DashboardHost string
}

// AddRequest décrit un domaine à ajouter.
type AddRequest struct {
	Host      string
	ProjectID string
	Service   string
	Port      int
	CreatedBy string // ex: "user:usr_…"
}

// Service gère les domaines.
type Service struct {
	mu      sync.Mutex
	store   store.DomainStore
	proxy   Proxy
	network Network
	stacks  Stacks
	auditor audit.Auditor
	cfg     Config
	logger  zerolog.Logger

	retry    time.Duration // premier délai entre deux réconciliations ratées
	stop     chan struct{}
	stopOnce sync.Once
	started  chan struct{}
	done     chan struct{}
}

// New crée le service.
func New(st store.DomainStore, p Proxy, n Network, stacks Stacks, auditor audit.Auditor, cfg Config, logger zerolog.Logger) *Service {
	return &Service{
		store: st, proxy: p, network: n, stacks: stacks, auditor: auditor, cfg: cfg,
		logger:  logger.With().Str("component", "domains").Logger(),
		retry:   5 * time.Second,
		stop:    make(chan struct{}),
		started: make(chan struct{}),
		done:    make(chan struct{}),
	}
}

// WithRetry change le premier délai entre deux réconciliations ratées, pour les tests.
func (s *Service) WithRetry(d time.Duration) *Service {
	s.retry = d
	return s
}

// NormalizeHost met un nom de domaine sous sa forme canonique (minuscules,
// sans point final, punycode) et le vérifie avec les règles du reverse proxy.
func NormalizeHost(host string) (string, error) {
	host = strings.TrimSuffix(strings.TrimSpace(host), ".")
	ascii, err := idna.Lookup.ToASCII(host)
	if err != nil {
		return "", fmt.Errorf("%w : %q n'est pas un nom de domaine", proxy.ErrInvalidHost, host)
	}
	if err := proxy.ValidateHost(ascii); err != nil {
		return "", err
	}
	return ascii, nil
}

// Alias renvoie le nom DNS du service sur le réseau partagé.
func Alias(projectID, service string) string {
	return service + "." + projectID + aliasSuffix
}

// List liste les domaines, d'un projet ou de tous, triés par nom.
func (s *Service) List(ctx context.Context, projectID string) ([]*store.Domain, error) {
	return s.store.ListDomains(ctx, projectID)
}

// Add relie un domaine au service déployé d'un projet. En cas d'échec, rien
// n'est enregistré, branché ni servi.
func (s *Service) Add(ctx context.Context, req AddRequest) (*store.Domain, error) {
	d, err := s.add(ctx, req)
	meta := map[string]any{
		"success": err == nil, "host": req.Host, "project_id": req.ProjectID,
		"service": req.Service, "port": req.Port,
	}
	resource := "project:" + req.ProjectID
	if err != nil {
		meta["error"] = err.Error()
	} else {
		resource = "domain:" + d.ID
		meta["host"] = d.Host
	}
	s.audit(ctx, "domain.add", resource, meta)
	return d, err
}

func (s *Service) add(ctx context.Context, req AddRequest) (*store.Domain, error) {
	host, err := NormalizeHost(req.Host)
	if err != nil {
		return nil, fmt.Errorf("%w : %w", ErrInvalid, err)
	}
	if s.cfg.DashboardHost != "" && host == s.cfg.DashboardHost {
		return nil, fmt.Errorf("%w : %s est réservé au dashboard de ForgeNet", ErrInvalid, host)
	}
	if err := checkTarget(req.ProjectID, req.Service, req.Port); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	containers, err := s.containers(ctx, req.ProjectID, req.Service)
	if err != nil {
		return nil, err
	}
	if len(containers) == 0 {
		return nil, fmt.Errorf("%w : %s dans le projet %s", ErrNotDeployed, req.Service, req.ProjectID)
	}

	d := &store.Domain{
		ID: ids.New("dom"), Host: host, ProjectID: req.ProjectID,
		Service: req.Service, Port: req.Port, CreatedBy: req.CreatedBy,
	}
	// L'enregistrement réserve le nom : un doublon s'arrête ici (store.ErrConflict).
	if err := s.store.CreateDomain(ctx, d); err != nil {
		return nil, err
	}

	attached, err := s.attach(ctx, req.ProjectID, req.Service, containers)
	if err == nil {
		err = s.applyAll(ctx)
	}
	if err != nil {
		s.rollbackAdd(d, attached)
		return nil, err
	}
	return d, nil
}

// rollbackAdd défait un ajout : le domaine est effacé, et le service
// débranché si ce domaine était le seul à le viser.
func (s *Service) rollbackAdd(d *store.Domain, attached []string) {
	ctx, cancel := context.WithTimeout(context.Background(), rollbackTimeout)
	defer cancel()
	if err := s.store.DeleteDomain(ctx, d.ID); err != nil {
		s.logger.Error().Err(err).Str("domain_id", d.ID).Msg("annulation de l'ajout : domaine non effacé")
		return
	}
	if s.stillTargeted(ctx, d.ProjectID, d.Service) {
		return
	}
	for _, id := range attached {
		if err := s.network.DisconnectFromNetwork(ctx, s.cfg.Network, id); err != nil {
			s.logger.Error().Err(err).Str("container_id", id).Msg("annulation de l'ajout : conteneur non débranché")
		}
	}
}

// Remove retire un domaine. Caddy cesse de le servir avant qu'il ne soit
// effacé : si Caddy refuse, rien ne change.
func (s *Service) Remove(ctx context.Context, id string) error {
	d, err := s.remove(ctx, id)
	meta := map[string]any{"success": err == nil}
	if d != nil {
		meta["host"] = d.Host
		meta["project_id"] = d.ProjectID
	}
	if err != nil {
		meta["error"] = err.Error()
	}
	s.audit(ctx, "domain.remove", "domain:"+id, meta)
	return err
}

func (s *Service) remove(ctx context.Context, id string) (*store.Domain, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	d, err := s.store.GetDomain(ctx, id)
	if err != nil {
		return nil, err
	}
	all, err := s.store.ListDomains(ctx, "")
	if err != nil {
		return d, err
	}
	remaining := make([]*store.Domain, 0, len(all))
	for _, other := range all {
		if other.ID != id {
			remaining = append(remaining, other)
		}
	}
	if err := s.proxy.Apply(ctx, routes(remaining)); err != nil {
		return d, err
	}
	if err := s.store.DeleteDomain(ctx, id); err != nil {
		// Caddy ne sert plus le domaine mais il est toujours en base : on
		// remet Caddy d'accord avec la base.
		if restoreErr := s.proxy.Apply(context.WithoutCancel(ctx), routes(all)); restoreErr != nil {
			s.logger.Error().Err(restoreErr).Msg("retrait annulé : Caddy non rechargé")
		}
		return d, err
	}

	if s.stillTargeted(ctx, d.ProjectID, d.Service) {
		return d, nil
	}
	// Le domaine est retiré : un débranchement raté laisse seulement le
	// service joignable par Caddy, sans route qui y mène.
	containers, err := s.containers(ctx, d.ProjectID, d.Service)
	if err != nil {
		s.logger.Warn().Err(err).Str("project_id", d.ProjectID).Msg("conteneurs non relus pour le débranchement")
		return d, nil
	}
	for _, c := range containers {
		if err := s.network.DisconnectFromNetwork(ctx, s.cfg.Network, c); err != nil {
			s.logger.Warn().Err(err).Str("container_id", c).Msg("conteneur non débranché")
		}
	}
	return d, nil
}

// Reattach rebranche les services d'un projet qui ont un domaine, par exemple
// après un déploiement qui a recréé leurs conteneurs.
func (s *Service) Reattach(ctx context.Context, projectID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	list, err := s.store.ListDomains(ctx, projectID)
	if err != nil {
		return err
	}
	return s.attachDomains(ctx, list)
}

// Reconcile remet tout d'accord avec la base : réseau, conteneur de Caddy,
// conteneurs des services, configuration de Caddy. Chaque étape est tentée ;
// les erreurs sont renvoyées ensemble.
func (s *Service) Reconcile(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.network.EnsureNetwork(ctx, s.cfg.Network, NetworkLabels); err != nil {
		return err
	}
	var errs []error
	if s.cfg.CaddyContainer != "" {
		if err := s.network.ConnectToNetwork(ctx, s.cfg.Network, s.cfg.CaddyContainer, nil); err != nil {
			errs = append(errs, fmt.Errorf("conteneur de Caddy : %w", err))
		}
	}
	list, err := s.store.ListDomains(ctx, "")
	if err != nil {
		return errors.Join(append(errs, err)...)
	}
	if err := s.attachDomains(ctx, list); err != nil {
		errs = append(errs, err)
	}
	if err := s.proxy.Apply(ctx, routes(list)); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// Start réconcilie au démarrage du serveur, et recommence avec un délai
// croissant tant que Caddy ou Docker ne répondent pas. Shutdown l'arrête.
func (s *Service) Start() error {
	close(s.started)
	defer close(s.done)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		select {
		case <-s.stop:
			cancel()
		case <-ctx.Done():
		}
	}()

	delay := s.retry
	for {
		err := s.Reconcile(ctx)
		if err == nil {
			s.logger.Info().Msg("domaines réconciliés avec Caddy et Docker")
			return nil
		}
		if ctx.Err() != nil {
			return nil
		}
		s.logger.Warn().Err(err).Dur("nouvel_essai", delay).Msg("réconciliation des domaines incomplète")
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(delay):
		}
		delay = min(delay*2, time.Minute)
	}
}

// Shutdown arrête la réconciliation en cours.
func (s *Service) Shutdown(ctx context.Context) error {
	s.stopOnce.Do(func() { close(s.stop) })
	select {
	case <-s.started:
	default:
		return nil // Start n'a jamais tourné : rien à attendre.
	}
	select {
	case <-s.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// attachDomains branche les services visés par ces domaines.
func (s *Service) attachDomains(ctx context.Context, list []*store.Domain) error {
	type target struct{ project, service string }
	seen := map[target]bool{}
	var errs []error
	for _, d := range list {
		t := target{d.ProjectID, d.Service}
		if seen[t] {
			continue
		}
		seen[t] = true
		containers, err := s.containers(ctx, t.project, t.service)
		if err == nil {
			_, err = s.attach(ctx, t.project, t.service, containers)
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("%s/%s : %w", t.project, t.service, err))
		}
	}
	return errors.Join(errs...)
}

// attach crée le réseau au besoin et y branche les conteneurs. Il renvoie
// ceux qui sont branchés, même en cas d'échec, pour pouvoir les débrancher.
func (s *Service) attach(ctx context.Context, projectID, service string, containers []string) ([]string, error) {
	if err := s.network.EnsureNetwork(ctx, s.cfg.Network, NetworkLabels); err != nil {
		return nil, err
	}
	alias := []string{Alias(projectID, service)}
	attached := make([]string, 0, len(containers))
	for _, c := range containers {
		if err := s.network.ConnectToNetwork(ctx, s.cfg.Network, c, alias); err != nil {
			return attached, err
		}
		attached = append(attached, c)
	}
	return attached, nil
}

// containers renvoie les identifiants des conteneurs du service dans la stack.
func (s *Service) containers(ctx context.Context, projectID, service string) ([]string, error) {
	all, err := s.stacks.Services(ctx, deploy.StackName(projectID))
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, c := range all {
		if c.Service == service {
			ids = append(ids, c.ID)
		}
	}
	return ids, nil
}

// stillTargeted indique si un domaine vise encore ce service.
func (s *Service) stillTargeted(ctx context.Context, projectID, service string) bool {
	list, err := s.store.ListDomains(ctx, projectID)
	if err != nil {
		return true // dans le doute, le service reste branché
	}
	for _, d := range list {
		if d.Service == service {
			return true
		}
	}
	return false
}

func (s *Service) applyAll(ctx context.Context) error {
	list, err := s.store.ListDomains(ctx, "")
	if err != nil {
		return err
	}
	return s.proxy.Apply(ctx, routes(list))
}

func (s *Service) audit(ctx context.Context, action, resource string, meta map[string]any) {
	if s.auditor != nil {
		_ = s.auditor.Log(ctx, action, resource, meta)
	}
}

// checkTarget vérifie le projet, le service et le port : l'alias DNS qui en
// découle doit être un nom que Caddy sait résoudre.
func checkTarget(projectID, service string, port int) error {
	if !deploy.ValidProjectID(projectID) || !dnsLabel.MatchString(projectID) {
		return fmt.Errorf("%w : identifiant de projet %q", ErrInvalid, projectID)
	}
	if !dnsLabel.MatchString(service) {
		return fmt.Errorf("%w : service %q (minuscules, chiffres, « _ » et « - », 63 caractères au plus)", ErrInvalid, service)
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("%w : port %d (1 à 65535)", ErrInvalid, port)
	}
	return nil
}

func routes(list []*store.Domain) []proxy.Route {
	out := make([]proxy.Route, 0, len(list))
	for _, d := range list {
		out = append(out, proxy.Route{
			Host:     d.Host,
			Upstream: proxy.Upstream{Container: Alias(d.ProjectID, d.Service), Port: d.Port},
		})
	}
	return out
}
