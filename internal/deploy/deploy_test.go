package deploy_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/deploy"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
)

const stackFile = "services:\n  web:\n    image: nginx:alpine\n"

// fakeStacks imite le moteur Compose. Deploy peut rester bloqué jusqu'à
// l'annulation de son contexte, pour éprouver l'annulation et l'arrêt.
type fakeStacks struct {
	mu          sync.Mutex
	validateErr error
	deployErr   error
	services    []compose.Service
	block       bool
	deployed    []string
	removed     []string
	started     chan struct{} // reçoit un signal à chaque Deploy
}

func newFakeStacks() *fakeStacks {
	return &fakeStacks{
		services: []compose.Service{{Service: "web", State: orchestrator.StateRunning, Health: "healthy"}},
		started:  make(chan struct{}, 10),
	}
}

func (f *fakeStacks) Validate(context.Context, []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.validateErr
}

func (f *fakeStacks) Deploy(ctx context.Context, project string, _ []byte) error {
	f.mu.Lock()
	f.deployed = append(f.deployed, project)
	block, err := f.block, f.deployErr
	f.mu.Unlock()
	f.started <- struct{}{}
	if block {
		<-ctx.Done()
		return ctx.Err()
	}
	return err
}

func (f *fakeStacks) Remove(_ context.Context, project string, removeVolumes bool) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	entry := project
	if removeVolumes {
		entry += "+volumes"
	}
	f.removed = append(f.removed, entry)
	return nil
}

func (f *fakeStacks) Services(context.Context, string) ([]compose.Service, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.services, nil
}

func (f *fakeStacks) set(fn func(*fakeStacks)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	fn(f)
}

func newStore(t *testing.T) store.Store {
	t.Helper()
	st, err := sqlite.Open(context.Background(), config.DatabaseConfig{Path: ":memory:", BusyTimeout: 2 * time.Second, MaxOpenConns: 1, MaxIdleConns: 1})
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })
	return st
}

func newService(t *testing.T, st store.Store, stacks deploy.Stacks, opts ...deploy.Option) *deploy.Service {
	t.Helper()
	opts = append([]deploy.Option{deploy.WithPollInterval(20 * time.Millisecond)}, opts...)
	return deploy.New(st.Deployments(), stacks, nil, zerolog.Nop(), opts...)
}

// run lance l'exécuteur et l'arrête à la fin du test.
func run(t *testing.T, svc *deploy.Service) {
	t.Helper()
	stopped := make(chan error, 1)
	go func() { stopped <- svc.Start() }()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		require.NoError(t, svc.Shutdown(ctx))
		require.NoError(t, <-stopped)
	})
}

func request(project string) deploy.CreateRequest {
	return deploy.CreateRequest{ProjectID: project, SourceType: store.SourceCompose, Content: stackFile, TriggeredBy: "user:u1"}
}

func waitState(t *testing.T, st store.Store, id string, want store.DeploymentState) *store.Deployment {
	t.Helper()
	var d *store.Deployment
	require.Eventually(t, func() bool {
		var err error
		d, err = st.Deployments().GetDeployment(context.Background(), id)
		return err == nil && d.State == want
	}, 5*time.Second, 10*time.Millisecond, "état attendu : %s", want)
	return d
}

func TestCreate_RejectsBadRequestsWithoutStoring(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	svc := newService(t, st, stacks)
	ctx := context.Background()

	cases := map[string]struct {
		req  deploy.CreateRequest
		want error
	}{
		"projet vide":          {deploy.CreateRequest{SourceType: store.SourceCompose, Content: stackFile}, deploy.ErrInvalid},
		"projet en majuscules": {deploy.CreateRequest{ProjectID: "Boutique", SourceType: store.SourceCompose, Content: stackFile}, deploy.ErrInvalid},
		"projet trop long":     {deploy.CreateRequest{ProjectID: strings.Repeat("a", 64), SourceType: store.SourceCompose, Content: stackFile}, deploy.ErrInvalid},
		"source git":           {deploy.CreateRequest{ProjectID: "p1", SourceType: store.SourceGit, Content: "x"}, deploy.ErrUnsupported},
		"fichier vide":         {deploy.CreateRequest{ProjectID: "p1", SourceType: store.SourceCompose, Content: "  "}, deploy.ErrInvalid},
		"fichier trop gros":    {deploy.CreateRequest{ProjectID: "p1", SourceType: store.SourceCompose, Content: strings.Repeat("x", deploy.MaxComposeSize+1)}, deploy.ErrInvalid},
	}
	for name, c := range cases {
		_, err := svc.Create(ctx, c.req)
		assert.ErrorIs(t, err, c.want, name)
	}

	stacks.set(func(f *fakeStacks) { f.validateErr = compose.ErrInvalidFile })
	_, err := svc.Create(ctx, request("p1"))
	assert.ErrorIs(t, err, deploy.ErrInvalid, "un fichier refusé par Compose est une demande invalide")

	stacks.set(func(f *fakeStacks) { f.validateErr = compose.ErrUnavailable })
	_, err = svc.Create(ctx, request("p1"))
	assert.ErrorIs(t, err, compose.ErrUnavailable)
	assert.NotErrorIs(t, err, deploy.ErrInvalid, "une commande absente n'est pas la faute du client")

	all, err := st.Deployments().ListDeployments(ctx, "", 10, 0)
	require.NoError(t, err)
	assert.Empty(t, all, "rien n'est enregistré")
}

func TestCreate_IdempotentAndChained(t *testing.T) {
	st := newStore(t)
	svc := newService(t, st, newFakeStacks())
	ctx := context.Background()

	first := request("p1")
	first.RequestID = "req-1"
	d1, err := svc.Create(ctx, first)
	require.NoError(t, err)
	assert.Equal(t, store.DeploymentPending, d1.State)
	assert.True(t, strings.HasPrefix(d1.ID, "dep_"))
	assert.Empty(t, d1.PreviousDeploymentID, "le premier déploiement n'a pas de prédécesseur")

	again, err := svc.Create(ctx, first)
	require.NoError(t, err)
	assert.Equal(t, d1.ID, again.ID, "la même clé rend le même déploiement")

	other := request("p2")
	other.RequestID = "req-1"
	_, err = svc.Create(ctx, other)
	assert.ErrorIs(t, err, deploy.ErrInvalid, "une clé ne sert pas à deux projets")

	d2, err := svc.Create(ctx, request("p1"))
	require.NoError(t, err)
	assert.Equal(t, d1.ID, d2.PreviousDeploymentID, "le déploiement suivant pointe vers le précédent")
}

func TestRun_DeploysToSuccessInOrder(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	svc := newService(t, st, stacks)
	ctx := context.Background()

	d1, err := svc.Create(ctx, request("p1"))
	require.NoError(t, err)
	d2, err := svc.Create(ctx, request("p2"))
	require.NoError(t, err)
	run(t, svc)

	done1 := waitState(t, st, d1.ID, store.DeploymentSuccess)
	waitState(t, st, d2.ID, store.DeploymentSuccess)
	assert.NotNil(t, done1.FinishedAt)
	assert.Empty(t, done1.ErrorMessage)

	stacks.mu.Lock()
	defer stacks.mu.Unlock()
	assert.Equal(t, []string{"forgenet-p1", "forgenet-p2"}, stacks.deployed, "une stack par projet, dans l'ordre de création")
}

func TestRun_FailuresAreRecorded(t *testing.T) {
	cases := map[string]struct {
		setup func(*fakeStacks)
		want  string
	}{
		"échec de docker compose up": {
			func(f *fakeStacks) { f.deployErr = errors.New("docker compose up : code 1 : image introuvable") },
			"image introuvable",
		},
		"service arrêté": {
			func(f *fakeStacks) {
				f.services = []compose.Service{{Service: "web", State: orchestrator.StateExited}}
			},
			"service web : conteneur exited",
		},
		"service malade": {
			func(f *fakeStacks) {
				f.services = []compose.Service{{Service: "web", State: orchestrator.StateRunning, Health: "unhealthy"}}
			},
			"service web : santé unhealthy",
		},
		"aucun conteneur": {func(f *fakeStacks) { f.services = nil }, "aucun conteneur"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			st := newStore(t)
			stacks := newFakeStacks()
			stacks.set(c.setup)
			svc := newService(t, st, stacks)
			d, err := svc.Create(context.Background(), request("p1"))
			require.NoError(t, err)
			run(t, svc)

			failed := waitState(t, st, d.ID, store.DeploymentFailed)
			assert.Contains(t, failed.ErrorMessage, c.want)
			assert.NotNil(t, failed.FinishedAt)
		})
	}
}

func TestRun_TimeoutIsReported(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	stacks.set(func(f *fakeStacks) { f.block = true })
	svc := newService(t, st, stacks, deploy.WithTimeout(100*time.Millisecond))
	d, err := svc.Create(context.Background(), request("p1"))
	require.NoError(t, err)
	run(t, svc)

	failed := waitState(t, st, d.ID, store.DeploymentFailed)
	assert.Contains(t, failed.ErrorMessage, "délai de 100ms dépassé")
}

func TestCancel_PendingRunningAndFinished(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	svc := newService(t, st, stacks)
	ctx := context.Background()

	// En attente : l'exécuteur ne tourne pas encore.
	pending, err := svc.Create(ctx, request("p1"))
	require.NoError(t, err)
	cancelled, err := svc.Cancel(ctx, pending.ID, "mauvaise version")
	require.NoError(t, err)
	assert.Equal(t, store.DeploymentCancelled, cancelled.State)
	assert.Equal(t, "annulé : mauvaise version", cancelled.ErrorMessage)

	// En cours : Deploy reste bloqué jusqu'à l'annulation.
	stacks.set(func(f *fakeStacks) { f.block = true })
	running, err := svc.Create(ctx, request("p1"))
	require.NoError(t, err)
	run(t, svc)
	<-stacks.started
	stopped, err := svc.Cancel(ctx, running.ID, "")
	require.NoError(t, err)
	assert.Equal(t, store.DeploymentCancelled, stopped.State, "Cancel attend l'arrêt effectif")
	assert.Equal(t, "annulé", stopped.ErrorMessage)

	// Terminé : plus rien à annuler.
	_, err = svc.Cancel(ctx, running.ID, "")
	assert.ErrorIs(t, err, deploy.ErrNotCancellable)

	_, err = svc.Cancel(ctx, "dep_absent", "")
	assert.ErrorIs(t, err, store.ErrNotFound)
}

func TestStart_RecoversInterruptedDeployments(t *testing.T) {
	st := newStore(t)
	ctx := context.Background()

	// Un arrêt brutal a laissé un déploiement en cours et un autre en attente.
	interrupted := &store.Deployment{ID: "dep_coupe", ProjectID: "p1", SourceType: store.SourceCompose, SourceContent: stackFile, TriggeredBy: "user:u1"}
	require.NoError(t, st.Deployments().CreateDeployment(ctx, interrupted))
	_, err := st.Deployments().TransitionDeployment(ctx, "dep_coupe", store.DeploymentPending, store.DeploymentBuilding, "")
	require.NoError(t, err)
	waiting := &store.Deployment{ID: "dep_attente", ProjectID: "p2", SourceType: store.SourceCompose, SourceContent: stackFile, TriggeredBy: "user:u1"}
	require.NoError(t, st.Deployments().CreateDeployment(ctx, waiting))

	run(t, newService(t, st, newFakeStacks()))

	failed := waitState(t, st, "dep_coupe", store.DeploymentFailed)
	assert.Equal(t, "interrompu par l'arrêt du serveur", failed.ErrorMessage)
	waitState(t, st, "dep_attente", store.DeploymentSuccess)
}

func TestShutdown_InterruptsAfterGracePeriod(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	stacks.set(func(f *fakeStacks) { f.block = true })
	svc := newService(t, st, stacks)
	d, err := svc.Create(context.Background(), request("p1"))
	require.NoError(t, err)

	stopped := make(chan error, 1)
	go func() { stopped <- svc.Start() }()
	<-stacks.started

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	require.NoError(t, svc.Shutdown(ctx))
	require.NoError(t, <-stopped)

	failed, err := st.Deployments().GetDeployment(context.Background(), d.ID)
	require.NoError(t, err)
	assert.Equal(t, store.DeploymentFailed, failed.State)
	assert.Equal(t, "interrompu par l'arrêt du serveur", failed.ErrorMessage)
}

func TestShutdown_BeforeStartReturns(t *testing.T) {
	svc := newService(t, newStore(t), newFakeStacks())
	require.NoError(t, svc.Shutdown(context.Background()))
	require.NoError(t, svc.Start(), "un Start tardif s'arrête aussitôt")
}

func TestDeleteStack_RefusedWhileProjectIsBusy(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	svc := newService(t, st, stacks)
	ctx := context.Background()

	require.NoError(t, svc.DeleteStack(ctx, "p1", false))
	require.NoError(t, svc.DeleteStack(ctx, "p2", true))
	assert.ErrorIs(t, svc.DeleteStack(ctx, "../p1", false), deploy.ErrInvalid)

	// Un déploiement en attente recréerait la stack juste après sa suppression.
	pending, err := svc.Create(ctx, request("p1"))
	require.NoError(t, err)
	assert.ErrorIs(t, svc.DeleteStack(ctx, "p1", false), deploy.ErrBusy)

	// Pendant l'exécution aussi.
	stacks.set(func(f *fakeStacks) { f.block = true })
	run(t, svc)
	<-stacks.started
	assert.ErrorIs(t, svc.DeleteStack(ctx, "p1", false), deploy.ErrBusy)
	_, err = svc.Cancel(ctx, pending.ID, "")
	require.NoError(t, err)
	require.NoError(t, svc.DeleteStack(ctx, "p1", false), "une fois le projet libre, la suppression passe")

	stacks.mu.Lock()
	defer stacks.mu.Unlock()
	assert.Equal(t, []string{"forgenet-p1", "forgenet-p2+volumes", "forgenet-p1"}, stacks.removed)
}
