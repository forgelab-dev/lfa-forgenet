package deploy_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/deploy"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

// recorder garde ce qu'un flux envoie, sans course entre goroutines.
type recorder[T any] struct {
	mu    sync.Mutex
	items []T
	first chan struct{}
	once  sync.Once
}

func newRecorder[T any]() *recorder[T] { return &recorder[T]{first: make(chan struct{})} }

func (r *recorder[T]) send(item T) error {
	r.mu.Lock()
	r.items = append(r.items, item)
	r.mu.Unlock()
	r.once.Do(func() { close(r.first) })
	return nil
}

func (r *recorder[T]) all() []T {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]T(nil), r.items...)
}

func states(ds []*store.Deployment) []store.DeploymentState {
	out := make([]store.DeploymentState, 0, len(ds))
	for _, d := range ds {
		out = append(out, d.State)
	}
	return out
}

func messages(logs []*store.DeploymentLog) []string {
	out := make([]string, 0, len(logs))
	for _, l := range logs {
		out = append(out, l.Message)
	}
	return out
}

// background lance f et rend son erreur à la fin.
func background(f func() error) <-chan error {
	done := make(chan error, 1)
	go func() { done <- f() }()
	return done
}

func waitDone(t *testing.T, done <-chan error) error {
	t.Helper()
	select {
	case err := <-done:
		return err
	case <-time.After(5 * time.Second):
		t.Fatal("le flux ne s'est pas terminé")
		return nil
	}
}

func TestWatch_FollowsEveryStateToTheEnd(t *testing.T) {
	st := newStore(t)
	svc := newService(t, st, newFakeStacks())
	ctx := context.Background()

	d, err := svc.Create(ctx, request("p1"))
	require.NoError(t, err)

	rec := newRecorder[*store.Deployment]()
	done := background(func() error { return svc.Watch(ctx, d.ID, rec.send) })
	<-rec.first // abonné avant que l'exécuteur ne démarre
	run(t, svc)

	require.NoError(t, waitDone(t, done), "le flux se termine de lui-même à l'état final")
	assert.Equal(t, []store.DeploymentState{
		store.DeploymentPending, store.DeploymentBuilding, store.DeploymentDeploying,
		store.DeploymentHealthChecking, store.DeploymentSuccess,
	}, states(rec.all()))

	// Un déploiement déjà terminé n'envoie que son état.
	again := newRecorder[*store.Deployment]()
	require.NoError(t, svc.Watch(ctx, d.ID, again.send))
	assert.Equal(t, []store.DeploymentState{store.DeploymentSuccess}, states(again.all()))
}

func TestWatch_UnknownAndCancelledContext(t *testing.T) {
	st := newStore(t)
	svc := newService(t, st, newFakeStacks())

	err := svc.Watch(context.Background(), "dep_absent", func(*store.Deployment) error { return nil })
	require.ErrorIs(t, err, store.ErrNotFound)

	d, err := svc.Create(context.Background(), request("p1"))
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	rec := newRecorder[*store.Deployment]()
	done := background(func() error { return svc.Watch(ctx, d.ID, rec.send) })
	<-rec.first
	cancel()
	require.ErrorIs(t, waitDone(t, done), context.Canceled, "le flux s'arrête quand le client part")

	sendErr := errors.New("client parti")
	err = svc.Watch(context.Background(), d.ID, func(*store.Deployment) error { return sendErr })
	require.ErrorIs(t, err, sendErr, "une erreur d'envoi arrête le flux")
}

func TestLogs_FollowThenReplayAndResume(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	stacks.set(func(f *fakeStacks) {
		f.output = []string{" Container forgenet-p1-web-1 Creating", " Container forgenet-p1-web-1 Started"}
	})
	svc := newService(t, st, stacks)
	ctx := context.Background()

	d, err := svc.Create(ctx, request("p1"))
	require.NoError(t, err)

	live := newRecorder[*store.DeploymentLog]()
	follow := deploy.LogOptions{Follow: true}
	done := background(func() error { return svc.Logs(ctx, d.ID, follow, live.send) })
	// Que l'abonnement précède ou non l'exécuteur, chaque ligne arrive une fois :
	// relue en base si elle est déjà écrite, sinon par le flux.
	time.Sleep(50 * time.Millisecond)
	run(t, svc)
	require.NoError(t, waitDone(t, done), "le suivi se termine avec le déploiement")

	want := []string{
		"Vérification du fichier Compose",
		"Fichier Compose valide",
		"Mise en service de la stack forgenet-p1 (docker compose up)",
		" Container forgenet-p1-web-1 Creating",
		" Container forgenet-p1-web-1 Started",
		"Vérification des services",
		"web : running, santé healthy",
		"Déploiement réussi",
	}
	got := live.all()
	assert.Equal(t, want, messages(got), "chaque ligne une fois, dans l'ordre")
	for i := 1; i < len(got); i++ {
		assert.True(t, got[i].EmittedAt.After(got[i-1].EmittedAt), "dates strictement croissantes")
	}
	assert.Equal(t, store.LogStderr, got[3].Stream, "la progression de Compose vient de stderr")
	assert.Equal(t, "deploy", got[3].Step)
	assert.Equal(t, "result", got[len(got)-1].Step)

	// Relecture d'un déploiement terminé, avec ou sans suivi.
	replay := newRecorder[*store.DeploymentLog]()
	require.NoError(t, svc.Logs(ctx, d.ID, follow, replay.send))
	assert.Equal(t, want, messages(replay.all()))

	// Reprise après une coupure : seulement ce qui suit la dernière ligne reçue.
	resumed := newRecorder[*store.DeploymentLog]()
	require.NoError(t, svc.Logs(ctx, d.ID, deploy.LogOptions{Since: got[4].EmittedAt}, resumed.send))
	assert.Equal(t, want[5:], messages(resumed.all()))

	// Un seul flux.
	stderr := newRecorder[*store.DeploymentLog]()
	require.NoError(t, svc.Logs(ctx, d.ID, deploy.LogOptions{Stream: store.LogStderr}, stderr.send))
	assert.Equal(t, want[3:5], messages(stderr.all()))
}

func TestLogs_RecordFailureAndCancellation(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	stacks.set(func(f *fakeStacks) { f.deployErr = errors.New("image introuvable : web:latest") })
	svc := newService(t, st, stacks)
	run(t, svc)
	ctx := context.Background()

	failed, err := svc.Create(ctx, request("p1"))
	require.NoError(t, err)
	waitState(t, st, failed.ID, store.DeploymentFailed)
	rec := newRecorder[*store.DeploymentLog]()
	require.NoError(t, svc.Logs(ctx, failed.ID, deploy.LogOptions{}, rec.send))
	last := rec.all()[len(rec.all())-1]
	assert.Equal(t, "Échec : image introuvable : web:latest", last.Message)
	assert.Equal(t, store.LogStderr, last.Stream)

	stacks.set(func(f *fakeStacks) { f.deployErr, f.block = nil, true })
	blocked, err := svc.Create(ctx, request("p2"))
	require.NoError(t, err)
	<-stacks.started
	watch := newRecorder[*store.Deployment]()
	done := background(func() error { return svc.Watch(ctx, blocked.ID, watch.send) })
	<-watch.first
	_, err = svc.Cancel(ctx, blocked.ID, "mauvaise version")
	require.NoError(t, err)
	require.NoError(t, waitDone(t, done))
	ws := watch.all()
	assert.Equal(t, store.DeploymentCancelled, ws[len(ws)-1].State, "le flux d'état voit l'annulation")

	rec = newRecorder[*store.DeploymentLog]()
	require.NoError(t, svc.Logs(ctx, blocked.ID, deploy.LogOptions{}, rec.send))
	logs := rec.all()
	assert.Equal(t, "Déploiement annulé : mauvaise version", logs[len(logs)-1].Message)

	err = svc.Logs(ctx, "dep_absent", deploy.LogOptions{Follow: true}, rec.send)
	require.ErrorIs(t, err, store.ErrNotFound)
}

func TestRun_FailureMessageIsTheLastComposeLine(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	stacks.set(func(f *fakeStacks) {
		f.deployErr = &compose.CommandError{Args: []string{"up"}, ExitCode: 1,
			Output: " Image web Pulling \n Image web Error\nError response from daemon: pull access denied for web"}
	})
	svc := newService(t, st, stacks)
	run(t, svc)

	d, err := svc.Create(context.Background(), request("p1"))
	require.NoError(t, err)
	failed := waitState(t, st, d.ID, store.DeploymentFailed)
	assert.Equal(t, "docker compose up : Error response from daemon: pull access denied for web", failed.ErrorMessage,
		"une ligne lisible ; la progression complète est au journal")
}
