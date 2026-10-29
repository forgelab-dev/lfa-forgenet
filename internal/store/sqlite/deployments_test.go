package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

const composeWeb = "services:\n  web:\n    image: nginx:alpine\n"

func newDeployment(id, projectID string) *store.Deployment {
	return &store.Deployment{
		ID: id, ProjectID: projectID, SourceType: store.SourceCompose, SourceContent: composeWeb,
		TargetEnv: "production", TriggeredBy: "user:u1",
	}
}

func TestDeploymentStore_CreateAndGet(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	d := newDeployment("d1", "p1")
	require.NoError(t, st.Deployments().CreateDeployment(ctx, d))
	assert.Equal(t, store.DeploymentPending, d.State, "un déploiement commence en attente")

	got, err := st.Deployments().GetDeployment(ctx, "d1")
	require.NoError(t, err)
	assert.Equal(t, "p1", got.ProjectID)
	assert.Equal(t, store.SourceCompose, got.SourceType)
	assert.Equal(t, composeWeb, got.SourceContent, "le fichier Compose est relu à l'identique")
	assert.Equal(t, store.DeploymentPending, got.State)
	assert.Equal(t, "production", got.TargetEnv)
	assert.Equal(t, "user:u1", got.TriggeredBy)
	assert.Empty(t, got.ErrorMessage)
	assert.Empty(t, got.PreviousDeploymentID)
	assert.Empty(t, got.RequestID)
	assert.Nil(t, got.FinishedAt)
	assert.True(t, got.CreatedAt.Equal(d.CreatedAt))

	_, err = st.Deployments().GetDeployment(ctx, "absent")
	assert.ErrorIs(t, err, store.ErrNotFound)

	require.ErrorIs(t, st.Deployments().CreateDeployment(ctx, newDeployment("d1", "p1")), store.ErrConflict,
		"un identifiant ne sert qu'une fois")
}

func TestDeploymentStore_RejectsUnknownSourceAndState(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	d := newDeployment("d1", "p1")
	d.SourceType = "helm"
	require.Error(t, st.Deployments().CreateDeployment(ctx, d))

	d = newDeployment("d2", "p1")
	d.State = "en-cours"
	require.Error(t, st.Deployments().CreateDeployment(ctx, d))
}

func TestDeploymentStore_RequestIDIsIdempotencyKey(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	first := newDeployment("d1", "p1")
	first.RequestID = "req-42"
	require.NoError(t, st.Deployments().CreateDeployment(ctx, first))

	again := newDeployment("d2", "p1")
	again.RequestID = "req-42"
	require.ErrorIs(t, st.Deployments().CreateDeployment(ctx, again), store.ErrConflict,
		"une même clé ne crée pas un second déploiement")

	got, err := st.Deployments().GetDeploymentByRequestID(ctx, "req-42")
	require.NoError(t, err)
	assert.Equal(t, "d1", got.ID)

	// Sans clé, les déploiements ne se gênent pas : la colonne vaut NULL.
	require.NoError(t, st.Deployments().CreateDeployment(ctx, newDeployment("d3", "p1")))
	require.NoError(t, st.Deployments().CreateDeployment(ctx, newDeployment("d4", "p1")))

	_, err = st.Deployments().GetDeploymentByRequestID(ctx, "")
	assert.ErrorIs(t, err, store.ErrNotFound)
}

func TestDeploymentStore_ListByProjectNewestFirst(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	base := time.Now().UTC().Truncate(time.Second)

	for i, item := range []struct{ id, project string }{{"a1", "p1"}, {"b1", "p2"}, {"a2", "p1"}, {"a3", "p1"}} {
		d := newDeployment(item.id, item.project)
		d.CreatedAt = base.Add(time.Duration(i) * 500 * time.Millisecond)
		require.NoError(t, st.Deployments().CreateDeployment(ctx, d))
	}

	ids := func(list []*store.Deployment) []string {
		out := make([]string, 0, len(list))
		for _, d := range list {
			out = append(out, d.ID)
		}
		return out
	}

	p1, err := st.Deployments().ListDeployments(ctx, "p1", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, []string{"a3", "a2", "a1"}, ids(p1))

	all, err := st.Deployments().ListDeployments(ctx, "", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, []string{"a3", "a2", "b1", "a1"}, ids(all), "un projet vide liste tous les déploiements")

	page, err := st.Deployments().ListDeployments(ctx, "p1", 2, 1)
	require.NoError(t, err)
	assert.Equal(t, []string{"a2", "a1"}, ids(page))

	none, err := st.Deployments().ListDeployments(ctx, "inconnu", 10, 0)
	require.NoError(t, err)
	assert.Empty(t, none)
}

func TestDeploymentStore_TransitionFollowsStateMachine(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	require.NoError(t, st.Deployments().CreateDeployment(ctx, newDeployment("d1", "p1")))

	d, err := st.Deployments().TransitionDeployment(ctx, "d1", store.DeploymentPending, store.DeploymentBuilding, "")
	require.NoError(t, err)
	assert.Equal(t, store.DeploymentBuilding, d.State)
	assert.Nil(t, d.FinishedAt, "un état intermédiaire ne clôt pas le déploiement")

	d, err = st.Deployments().TransitionDeployment(ctx, "d1", store.DeploymentBuilding, store.DeploymentFailed, "image introuvable")
	require.NoError(t, err)
	assert.Equal(t, store.DeploymentFailed, d.State)
	assert.Equal(t, "image introuvable", d.ErrorMessage)
	require.NotNil(t, d.FinishedAt, "un état final date la fin du déploiement")
	assert.False(t, d.UpdatedAt.Before(d.CreatedAt))

	_, err = st.Deployments().TransitionDeployment(ctx, "d1", store.DeploymentFailed, store.DeploymentBuilding, "")
	assert.ErrorIs(t, err, store.ErrInvalidTransition, "un échec ne repart pas")

	_, err = st.Deployments().TransitionDeployment(ctx, "absent", store.DeploymentPending, store.DeploymentBuilding, "")
	assert.ErrorIs(t, err, store.ErrNotFound)
}

// Une annulation et l'exécuteur partent du même état : le second écrivain doit
// échouer au lieu d'écraser l'état choisi par le premier.
func TestDeploymentStore_TransitionRejectsStaleState(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	require.NoError(t, st.Deployments().CreateDeployment(ctx, newDeployment("d1", "p1")))

	_, err := st.Deployments().TransitionDeployment(ctx, "d1", store.DeploymentPending, store.DeploymentCancelled, "annulé par l'utilisateur")
	require.NoError(t, err)

	_, err = st.Deployments().TransitionDeployment(ctx, "d1", store.DeploymentPending, store.DeploymentBuilding, "")
	assert.ErrorIs(t, err, store.ErrConflict)

	got, err := st.Deployments().GetDeployment(ctx, "d1")
	require.NoError(t, err)
	assert.Equal(t, store.DeploymentCancelled, got.State, "l'annulation reste acquise")
}

func TestDeploymentStore_PreviousDeploymentMustExist(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	require.NoError(t, st.Deployments().CreateDeployment(ctx, newDeployment("d1", "p1")))

	orphan := newDeployment("d2", "p1")
	orphan.PreviousDeploymentID = "absent"
	require.Error(t, st.Deployments().CreateDeployment(ctx, orphan), "la chaîne de rollback ne pointe pas dans le vide")

	next := newDeployment("d3", "p1")
	next.PreviousDeploymentID = "d1"
	require.NoError(t, st.Deployments().CreateDeployment(ctx, next))

	got, err := st.Deployments().GetDeployment(ctx, "d3")
	require.NoError(t, err)
	assert.Equal(t, "d1", got.PreviousDeploymentID)
}

func TestDeploymentStore_ListByStateOldestFirst(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	base := time.Now().UTC().Truncate(time.Second)

	for i, id := range []string{"d1", "d2", "d3"} {
		d := newDeployment(id, "p1")
		d.CreatedAt = base.Add(time.Duration(i) * time.Second)
		require.NoError(t, st.Deployments().CreateDeployment(ctx, d))
	}
	_, err := st.Deployments().TransitionDeployment(ctx, "d2", store.DeploymentPending, store.DeploymentBuilding, "")
	require.NoError(t, err)

	pending, err := st.Deployments().ListDeploymentsByState(ctx, store.DeploymentPending, 10)
	require.NoError(t, err)
	require.Len(t, pending, 2)
	assert.Equal(t, "d1", pending[0].ID, "le plus ancien passe en premier")
	assert.Equal(t, "d3", pending[1].ID)

	building, err := st.Deployments().ListDeploymentsByState(ctx, store.DeploymentBuilding, 10)
	require.NoError(t, err)
	require.Len(t, building, 1)
	assert.Equal(t, "d2", building[0].ID)

	one, err := st.Deployments().ListDeploymentsByState(ctx, store.DeploymentPending, 1)
	require.NoError(t, err)
	assert.Len(t, one, 1)
}
