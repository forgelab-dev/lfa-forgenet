package domains_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/domains"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
	"github.com/forgelab-dev/lfa-forgenet/internal/proxy"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
)

// fakeProxy garde les dernières routes appliquées.
type fakeProxy struct {
	mu      sync.Mutex
	applied [][]proxy.Route
	err     error
}

func (p *fakeProxy) Apply(_ context.Context, routes []proxy.Route) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.err != nil {
		return p.err
	}
	p.applied = append(p.applied, routes)
	return nil
}

func (p *fakeProxy) hosts() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.applied) == 0 {
		return nil
	}
	var out []string
	for _, r := range p.applied[len(p.applied)-1] {
		out = append(out, r.Host)
	}
	slices.Sort(out)
	return out
}

// fakeNetwork tient les alias de chaque conteneur branché.
type fakeNetwork struct {
	mu        sync.Mutex
	ensured   int
	ensureErr error
	attached  map[string][]string
	failOn    string // conteneur dont le branchement échoue
}

func (n *fakeNetwork) EnsureNetwork(_ context.Context, name string, labels map[string]string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if name != "forgenet-edge" || labels["forgenet.managed"] != "edge" {
		return fmt.Errorf("réseau inattendu %s %v", name, labels)
	}
	n.ensured++
	return n.ensureErr
}

func (n *fakeNetwork) ConnectToNetwork(_ context.Context, _, id string, aliases []string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if id == n.failOn {
		return errors.New("branchement refusé par Docker")
	}
	if n.attached == nil {
		n.attached = map[string][]string{}
	}
	n.attached[id] = aliases
	return nil
}

func (n *fakeNetwork) DisconnectFromNetwork(_ context.Context, _, id string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.attached, id)
	return nil
}

func (n *fakeNetwork) aliases(id string) ([]string, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	a, ok := n.attached[id]
	return a, ok
}

// fakeStacks associe à chaque stack ses conteneurs.
type fakeStacks struct {
	mu       sync.Mutex
	services map[string][]compose.Service
	err      error
}

func (s *fakeStacks) Services(_ context.Context, project string) ([]compose.Service, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.services[project], s.err
}

func (s *fakeStacks) set(stack string, svc ...compose.Service) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.services[stack] = svc
}

type fakeAuditor struct {
	audit.Auditor // List n'est pas utilisée

	mu      sync.Mutex
	actions []string
}

func (a *fakeAuditor) Log(_ context.Context, action, _ string, meta map[string]any) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.actions = append(a.actions, fmt.Sprintf("%s:%v", action, meta["success"]))
	return nil
}

type fixture struct {
	svc     *domains.Service
	store   store.DomainStore
	proxy   *fakeProxy
	network *fakeNetwork
	stacks  *fakeStacks
	auditor *fakeAuditor
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	db, err := sqlite.Open(context.Background(), config.DatabaseConfig{Path: ":memory:", BusyTimeout: time.Second, MaxOpenConns: 1, MaxIdleConns: 1})
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	f := &fixture{
		store:   db.Domains(),
		proxy:   &fakeProxy{},
		network: &fakeNetwork{},
		stacks: &fakeStacks{services: map[string][]compose.Service{
			"forgenet-boutique": {
				{ID: "c-web-1", Service: "web"}, {ID: "c-web-2", Service: "web"}, {ID: "c-db-1", Service: "db"},
			},
		}},
		auditor: &fakeAuditor{},
	}
	f.svc = domains.New(f.store, f.proxy, f.network, f.stacks, f.auditor,
		domains.Config{Network: "forgenet-edge", CaddyContainer: "forgenet-caddy"}, zerolog.Nop())
	return f
}

func webRequest(host string) domains.AddRequest {
	return domains.AddRequest{Host: host, ProjectID: "boutique", Service: "web", Port: 8080, CreatedBy: "user:usr_1"}
}

const webAlias = "web.boutique.forgenet.internal"

func TestNormalizeHost(t *testing.T) {
	for in, want := range map[string]string{
		"App.Example.ORG.":  "app.example.org",
		"  app.example.org": "app.example.org",
		"bücher.example":    "xn--bcher-kva.example",
	} {
		got, err := domains.NormalizeHost(in)
		require.NoError(t, err, in)
		assert.Equal(t, want, got, in)
	}
	for _, in := range []string{"", "localhost", "a..org", "*.example.org", "1.2.3.4", "a.org {", "a.org\nb.org", "a_b.org"} {
		_, err := domains.NormalizeHost(in)
		assert.ErrorIs(t, err, proxy.ErrInvalidHost, "%q", in)
	}
}

func TestAlias(t *testing.T) {
	assert.Equal(t, webAlias, domains.Alias("boutique", "web"))
}

func TestAdd_RecordsAttachesAndServes(t *testing.T) {
	f := newFixture(t)
	d, err := f.svc.Add(context.Background(), webRequest("Shop.Example.org"))
	require.NoError(t, err)
	assert.Equal(t, "shop.example.org", d.Host)
	assert.Equal(t, "user:usr_1", d.CreatedBy)

	stored, err := f.store.GetDomain(context.Background(), d.ID)
	require.NoError(t, err)
	assert.Equal(t, "web", stored.Service)

	for _, c := range []string{"c-web-1", "c-web-2"} {
		aliases, ok := f.network.aliases(c)
		assert.True(t, ok, c)
		assert.Equal(t, []string{webAlias}, aliases)
	}
	_, dbAttached := f.network.aliases("c-db-1")
	assert.False(t, dbAttached, "la base de la stack reste hors du réseau partagé")

	require.Len(t, f.proxy.applied, 1)
	assert.Equal(t, []proxy.Route{{Host: "shop.example.org", Upstream: proxy.Upstream{Container: webAlias, Port: 8080}}}, f.proxy.applied[0])
	assert.Equal(t, []string{"domain.add:true"}, f.auditor.actions)
}

func TestAdd_RefusesInvalidRequests(t *testing.T) {
	cases := map[string]domains.AddRequest{
		"domaine":           webRequest("a.org {\n}"),
		"projet":            {Host: "a.example.org", ProjectID: "Boutique", Service: "web", Port: 80},
		"projet en tiret":   {Host: "a.example.org", ProjectID: "boutique-", Service: "web", Port: 80},
		"service majuscule": {Host: "a.example.org", ProjectID: "boutique", Service: "Web", Port: 80},
		"service pointé":    {Host: "a.example.org", ProjectID: "boutique", Service: "web.a", Port: 80},
		"port nul":          {Host: "a.example.org", ProjectID: "boutique", Service: "web", Port: 0},
		"port trop grand":   {Host: "a.example.org", ProjectID: "boutique", Service: "web", Port: 65536},
	}
	for name, req := range cases {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t)
			_, err := f.svc.Add(context.Background(), req)
			assert.ErrorIs(t, err, domains.ErrInvalid)
			assert.Empty(t, f.proxy.applied)
			assert.Empty(t, f.network.attached)
			assert.Equal(t, []string{"domain.add:false"}, f.auditor.actions)
		})
	}
}

func TestAdd_ServiceNotDeployed(t *testing.T) {
	f := newFixture(t)
	req := webRequest("a.example.org")
	req.Service = "api"
	_, err := f.svc.Add(context.Background(), req)
	assert.ErrorIs(t, err, domains.ErrNotDeployed)

	list, err := f.store.ListDomains(context.Background(), "")
	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestAdd_DuplicateHost(t *testing.T) {
	f := newFixture(t)
	_, err := f.svc.Add(context.Background(), webRequest("a.example.org"))
	require.NoError(t, err)

	_, err = f.svc.Add(context.Background(), webRequest("A.example.org."))
	assert.ErrorIs(t, err, store.ErrConflict)
	assert.Len(t, f.proxy.applied, 1, "Caddy n'est pas rechargé pour un doublon")
}

func TestAdd_CaddyRefusalLeavesNothing(t *testing.T) {
	f := newFixture(t)
	f.proxy.err = &proxy.LoadError{Message: "refusé"}

	_, err := f.svc.Add(context.Background(), webRequest("a.example.org"))
	var loadErr *proxy.LoadError
	require.ErrorAs(t, err, &loadErr)

	list, err := f.store.ListDomains(context.Background(), "")
	require.NoError(t, err)
	assert.Empty(t, list, "le domaine n'est pas gardé")
	assert.Empty(t, f.network.attached, "le service est débranché")
}

func TestAdd_RollbackKeepsServiceOfOtherDomain(t *testing.T) {
	f := newFixture(t)
	_, err := f.svc.Add(context.Background(), webRequest("a.example.org"))
	require.NoError(t, err)

	f.proxy.err = errors.New("Caddy injoignable")
	_, err = f.svc.Add(context.Background(), webRequest("b.example.org"))
	require.Error(t, err)

	_, ok := f.network.aliases("c-web-1")
	assert.True(t, ok, "a.example.org vise toujours ce service")
	list, err := f.store.ListDomains(context.Background(), "")
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "a.example.org", list[0].Host)
}

func TestAdd_ConnectFailureDetachesOthers(t *testing.T) {
	f := newFixture(t)
	f.network.failOn = "c-web-2"

	_, err := f.svc.Add(context.Background(), webRequest("a.example.org"))
	require.Error(t, err)
	assert.Empty(t, f.network.attached, "le premier conteneur branché est débranché")
	assert.Empty(t, f.proxy.applied)
}

func TestRemove(t *testing.T) {
	f := newFixture(t)
	a, err := f.svc.Add(context.Background(), webRequest("a.example.org"))
	require.NoError(t, err)
	b, err := f.svc.Add(context.Background(), webRequest("b.example.org"))
	require.NoError(t, err)

	require.NoError(t, f.svc.Remove(context.Background(), a.ID))
	assert.Equal(t, []string{"b.example.org"}, f.proxy.hosts())
	_, ok := f.network.aliases("c-web-1")
	assert.True(t, ok, "b.example.org vise encore le service")

	require.NoError(t, f.svc.Remove(context.Background(), b.ID))
	assert.Empty(t, f.proxy.hosts())
	assert.Empty(t, f.network.attached, "plus aucun domaine : le service est débranché")
	assert.Equal(t, []string{"domain.add:true", "domain.add:true", "domain.remove:true", "domain.remove:true"}, f.auditor.actions)
}

func TestRemove_CaddyRefusalChangesNothing(t *testing.T) {
	f := newFixture(t)
	a, err := f.svc.Add(context.Background(), webRequest("a.example.org"))
	require.NoError(t, err)

	f.proxy.err = errors.New("Caddy injoignable")
	require.Error(t, f.svc.Remove(context.Background(), a.ID))

	_, err = f.store.GetDomain(context.Background(), a.ID)
	require.NoError(t, err, "le domaine est toujours là")
	_, ok := f.network.aliases("c-web-1")
	assert.True(t, ok)
}

func TestRemove_Unknown(t *testing.T) {
	f := newFixture(t)
	assert.ErrorIs(t, f.svc.Remove(context.Background(), "dom_absent"), store.ErrNotFound)
	assert.Equal(t, []string{"domain.remove:false"}, f.auditor.actions)
}

func TestReattach_AfterRedeploy(t *testing.T) {
	f := newFixture(t)
	_, err := f.svc.Add(context.Background(), webRequest("a.example.org"))
	require.NoError(t, err)

	// Un redéploiement a recréé le conteneur du service.
	f.stacks.set("forgenet-boutique", compose.Service{ID: "c-web-3", Service: "web"}, compose.Service{ID: "c-db-1", Service: "db"})
	require.NoError(t, f.svc.Reattach(context.Background(), "boutique"))

	aliases, ok := f.network.aliases("c-web-3")
	assert.True(t, ok)
	assert.Equal(t, []string{webAlias}, aliases)
}

func TestReconcile(t *testing.T) {
	f := newFixture(t)
	_, err := f.svc.Add(context.Background(), webRequest("a.example.org"))
	require.NoError(t, err)
	f.network.attached = nil // redémarrage : les branchements sont à refaire

	require.NoError(t, f.svc.Reconcile(context.Background()))
	_, caddy := f.network.aliases("forgenet-caddy")
	assert.True(t, caddy, "le conteneur de Caddy est branché")
	_, web := f.network.aliases("c-web-1")
	assert.True(t, web)
	assert.Equal(t, []string{"a.example.org"}, f.proxy.hosts())
}

func TestReconcile_KeepsGoingAfterAnError(t *testing.T) {
	f := newFixture(t)
	_, err := f.svc.Add(context.Background(), webRequest("a.example.org"))
	require.NoError(t, err)
	f.stacks.err = errors.New("docker compose indisponible")
	f.proxy.applied = nil

	err = f.svc.Reconcile(context.Background())
	assert.ErrorContains(t, err, "docker compose indisponible")
	assert.Equal(t, []string{"a.example.org"}, f.proxy.hosts(), "Caddy est rechargé malgré tout")
}

func TestStart_RetriesUntilReconciled(t *testing.T) {
	f := newFixture(t)
	f.network.ensureErr = errors.New("démon Docker injoignable")
	f.svc.WithRetry(10 * time.Millisecond)

	done := make(chan error, 1)
	go func() { done <- f.svc.Start() }()
	time.Sleep(50 * time.Millisecond)
	f.network.mu.Lock()
	f.network.ensureErr = nil
	f.network.mu.Unlock()

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("la réconciliation ne s'est pas terminée")
	}
	f.network.mu.Lock()
	defer f.network.mu.Unlock()
	assert.GreaterOrEqual(t, f.network.ensured, 2, "au moins un nouvel essai")
}

func TestShutdown(t *testing.T) {
	t.Run("sans démarrage", func(t *testing.T) {
		require.NoError(t, newFixture(t).svc.Shutdown(context.Background()))
	})
	t.Run("arrête une réconciliation qui échoue", func(t *testing.T) {
		f := newFixture(t)
		f.network.ensureErr = errors.New("démon Docker injoignable")
		f.svc.WithRetry(time.Hour)
		done := make(chan error, 1)
		go func() { done <- f.svc.Start() }()
		time.Sleep(20 * time.Millisecond)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		require.NoError(t, f.svc.Shutdown(ctx))
		require.NoError(t, <-done)
	})
}

func TestAdd_Concurrent(t *testing.T) {
	f := newFixture(t)
	var wg sync.WaitGroup
	for i := range 6 {
		wg.Go(func() {
			_, err := f.svc.Add(context.Background(), webRequest(fmt.Sprintf("d%d.example.org", i)))
			assert.NoError(t, err)
		})
	}
	wg.Wait()
	assert.Len(t, f.proxy.hosts(), 6, "le dernier rechargement sert tous les domaines")
}
