package grpc_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/deploy"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

const composeFile = "services:\n  web:\n    image: nginx:alpine\n"

// fakeStacks répond comme un moteur Compose dont la stack démarre saine.
type fakeStacks struct {
	mu          sync.Mutex
	validateErr error
	removed     []string
}

func (f *fakeStacks) Validate(context.Context, []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.validateErr
}

func (f *fakeStacks) Deploy(context.Context, string, []byte) error { return nil }

func (f *fakeStacks) Remove(_ context.Context, project string, _ bool) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.removed = append(f.removed, project)
	return nil
}

func (f *fakeStacks) Services(context.Context, string) ([]compose.Service, error) {
	return []compose.Service{{Service: "web", State: orchestrator.StateRunning}}, nil
}

type deploymentTestServer struct {
	*testServer
	deployments forgenetv1.DeploymentServiceClient
	service     *deploy.Service
	stacks      *fakeStacks
}

// startDeploymentServer sert DeploymentService derrière le vrai garde ;
// l'exécuteur ne tourne que si withRunner est vrai.
func startDeploymentServer(t *testing.T, withRunner bool) *deploymentTestServer {
	t.Helper()
	dts := &deploymentTestServer{stacks: &fakeStacks{}}
	dts.testServer = startTestServer(t, &fakeEngine{}, func(srv *grpc.Server, st store.Store, auditor audit.Auditor) {
		dts.service = deploy.New(st.Deployments(), dts.stacks, auditor, zerolog.Nop(), deploy.WithPollInterval(20*time.Millisecond))
		forgenetv1.RegisterDeploymentServiceServer(srv, transportGRPC.NewDeploymentServiceServer(dts.service, zerolog.Nop()))
	})
	dts.deployments = forgenetv1.NewDeploymentServiceClient(dts.conn)
	if withRunner {
		go func() { _ = dts.service.Start() }()
		t.Cleanup(func() { _ = dts.service.Shutdown(context.Background()) })
	}
	return dts
}

func composeRequest(project string) *forgenetv1.CreateDeploymentRequest {
	return &forgenetv1.CreateDeploymentRequest{
		ProjectId: project,
		TargetEnv: "production",
		Source: &forgenetv1.DeploymentSource{Source: &forgenetv1.DeploymentSource_Compose{
			Compose: &forgenetv1.ComposeSource{ComposeContent: composeFile},
		}},
	}
}

// Décision du 29/10 : déployer est réservé aux administrateurs. Tout le
// service l'est, lecture comprise : le fichier Compose peut porter des secrets.
func TestDeploymentService_AdminOnly(t *testing.T) {
	ts := startDeploymentServer(t, false)

	calls := map[string]func(ctx context.Context) error{
		"Create": func(ctx context.Context) error {
			_, err := ts.deployments.CreateDeployment(ctx, composeRequest("p1"))
			return err
		},
		"Get": func(ctx context.Context) error {
			_, err := ts.deployments.GetDeployment(ctx, &forgenetv1.GetDeploymentRequest{DeploymentId: "dep_absent"})
			return err
		},
		"List": func(ctx context.Context) error {
			_, err := ts.deployments.ListDeployments(ctx, &forgenetv1.ListDeploymentsRequest{})
			return err
		},
		"Cancel": func(ctx context.Context) error {
			_, err := ts.deployments.CancelDeployment(ctx, &forgenetv1.CancelDeploymentRequest{DeploymentId: "dep_absent"})
			return err
		},
		"DeleteStack": func(ctx context.Context) error {
			_, err := ts.deployments.DeleteStack(ctx, &forgenetv1.DeleteStackRequest{ProjectId: "p2"})
			return err
		},
	}
	for name, call := range calls {
		assert.Equal(t, codes.Unauthenticated, status.Code(call(context.Background())), "%s sans jeton", name)
		assert.Equal(t, codes.PermissionDenied, status.Code(call(as(ts.memberToken))), "%s par un membre", name)
		code := status.Code(call(as(ts.adminToken)))
		assert.NotContains(t, []codes.Code{codes.Unauthenticated, codes.PermissionDenied}, code, "%s par un admin", name)
	}
}

func TestDeploymentService_CreateRunsToSuccess(t *testing.T) {
	ts := startDeploymentServer(t, true)
	ctx := as(ts.adminToken)

	created, err := ts.deployments.CreateDeployment(ctx, composeRequest("boutique"))
	require.NoError(t, err)
	d := created.GetDeployment()
	assert.Equal(t, forgenetv1.DeploymentState_DEPLOYMENT_STATE_PENDING, d.GetState(), "la création répond aussitôt, en attente")
	assert.Regexp(t, `^user:usr_[0-9a-f]{16}$`, d.GetTriggeredBy())
	assert.Equal(t, "production", d.GetTargetEnv())
	assert.Equal(t, composeFile, d.GetSource().GetCompose().GetComposeContent())
	assert.NotNil(t, d.GetStrategy().GetRecreate(), "une stack Compose est recréée")

	var got *forgenetv1.Deployment
	require.Eventually(t, func() bool {
		resp, err := ts.deployments.GetDeployment(ctx, &forgenetv1.GetDeploymentRequest{DeploymentId: d.GetId()})
		got = resp.GetDeployment()
		return err == nil && got.GetState() == forgenetv1.DeploymentState_DEPLOYMENT_STATE_SUCCESS
	}, 5*time.Second, 20*time.Millisecond)
	assert.NotNil(t, got.GetFinishedAt())
}

func TestDeploymentService_ListPagination(t *testing.T) {
	ts := startDeploymentServer(t, false)
	ctx := as(ts.adminToken)

	var ids []string
	for range 3 {
		resp, err := ts.deployments.CreateDeployment(ctx, composeRequest("p1"))
		require.NoError(t, err)
		ids = append(ids, resp.GetDeployment().GetId())
	}
	_, err := ts.deployments.CreateDeployment(ctx, composeRequest("p2"))
	require.NoError(t, err)

	page1, err := ts.deployments.ListDeployments(ctx, &forgenetv1.ListDeploymentsRequest{ProjectId: "p1", PageSize: 2})
	require.NoError(t, err)
	require.Len(t, page1.GetDeployments(), 2)
	assert.Equal(t, ids[2], page1.GetDeployments()[0].GetId(), "le plus récent d'abord")
	assert.Equal(t, ids[1], page1.GetDeployments()[0].GetPreviousDeploymentId())
	require.NotEmpty(t, page1.GetNextPageToken())

	page2, err := ts.deployments.ListDeployments(ctx, &forgenetv1.ListDeploymentsRequest{ProjectId: "p1", PageSize: 2, PageToken: page1.GetNextPageToken()})
	require.NoError(t, err)
	require.Len(t, page2.GetDeployments(), 1)
	assert.Equal(t, ids[0], page2.GetDeployments()[0].GetId())
	assert.Empty(t, page2.GetNextPageToken(), "pas de page après la dernière")

	_, err = ts.deployments.ListDeployments(ctx, &forgenetv1.ListDeploymentsRequest{PageToken: "x"})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestDeploymentService_ErrorTranslation(t *testing.T) {
	ts := startDeploymentServer(t, false)
	ctx := as(ts.adminToken)

	git := composeRequest("p1")
	git.Source = &forgenetv1.DeploymentSource{Source: &forgenetv1.DeploymentSource_Git{Git: &forgenetv1.GitSource{RepositoryUrl: "https://example.test/depot.git"}}}
	noSource := composeRequest("p1")
	noSource.Source = nil
	rolling := composeRequest("p1")
	rolling.Strategy = &forgenetv1.StrategyConfig{Strategy: &forgenetv1.StrategyConfig_Rolling{Rolling: &forgenetv1.RollingStrategy{}}}
	badProject := composeRequest("Pas_Valide")

	for name, c := range map[string]struct {
		req  *forgenetv1.CreateDeploymentRequest
		want codes.Code
	}{
		"source Git":        {git, codes.Unimplemented},
		"source absente":    {noSource, codes.InvalidArgument},
		"stratégie rolling": {rolling, codes.Unimplemented},
		"projet invalide":   {badProject, codes.InvalidArgument},
	} {
		_, err := ts.deployments.CreateDeployment(ctx, c.req)
		assert.Equal(t, c.want, status.Code(err), name)
	}

	ts.stacks.mu.Lock()
	// Comme le vrai moteur : le refus de `config` est enveloppé dans ErrInvalidFile.
	ts.stacks.validateErr = fmt.Errorf("%w : %w", compose.ErrInvalidFile,
		&compose.CommandError{Args: []string{"config"}, ExitCode: 15, Output: "services.web.image must be a string"})
	ts.stacks.mu.Unlock()
	_, err := ts.deployments.CreateDeployment(ctx, composeRequest("p1"))
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Contains(t, status.Convert(err).Message(), "services.web.image must be a string", "l'auteur du fichier sait quoi corriger")

	ts.stacks.mu.Lock()
	ts.stacks.validateErr = compose.ErrUnavailable
	ts.stacks.mu.Unlock()
	_, err = ts.deployments.CreateDeployment(ctx, composeRequest("p1"))
	assert.Equal(t, codes.Unavailable, status.Code(err))

	ts.stacks.mu.Lock()
	ts.stacks.validateErr = nil
	ts.stacks.mu.Unlock()
	created, err := ts.deployments.CreateDeployment(ctx, composeRequest("p1"))
	require.NoError(t, err)
	_, err = ts.deployments.DeleteStack(ctx, &forgenetv1.DeleteStackRequest{ProjectId: "p1"})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err), "pas de suppression avec un déploiement en attente")

	cancelled, err := ts.deployments.CancelDeployment(ctx, &forgenetv1.CancelDeploymentRequest{DeploymentId: created.GetDeployment().GetId(), Reason: "erreur de version"})
	require.NoError(t, err)
	assert.Equal(t, forgenetv1.DeploymentState_DEPLOYMENT_STATE_CANCELLED, cancelled.GetDeployment().GetState())
	_, err = ts.deployments.CancelDeployment(ctx, &forgenetv1.CancelDeploymentRequest{DeploymentId: created.GetDeployment().GetId()})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err), "un déploiement terminé ne s'annule plus")

	_, err = ts.deployments.GetDeployment(ctx, &forgenetv1.GetDeploymentRequest{DeploymentId: "dep_absent"})
	assert.Equal(t, codes.NotFound, status.Code(err))

	_, err = ts.deployments.RollbackDeployment(ctx, &forgenetv1.RollbackDeploymentRequest{TargetDeploymentId: "dep_x"})
	assert.Equal(t, codes.Unimplemented, status.Code(err), "le retour en arrière viendra avec les stratégies")
}

func TestDeploymentService_ActionsAreAudited(t *testing.T) {
	ts := startDeploymentServer(t, false)
	ctx := as(ts.adminToken)

	created, err := ts.deployments.CreateDeployment(ctx, composeRequest("p1"))
	require.NoError(t, err)
	_, err = ts.deployments.CancelDeployment(ctx, &forgenetv1.CancelDeploymentRequest{DeploymentId: created.GetDeployment().GetId()})
	require.NoError(t, err)
	_, err = ts.deployments.DeleteStack(ctx, &forgenetv1.DeleteStackRequest{ProjectId: "p1", DeleteVolumes: true})
	require.NoError(t, err)

	events, err := ts.store.Audit().List(context.Background(), 20, 0)
	require.NoError(t, err)
	actions := map[string]*store.AuditEvent{}
	for _, e := range events {
		actions[e.Action] = e
	}
	for _, action := range []string{"deployment.create", "deployment.cancel", "stack.delete"} {
		require.Contains(t, actions, action)
		assert.NotEmpty(t, actions[action].UserID, "%s porte son auteur", action)
	}
	assert.Equal(t, "deployment:"+created.GetDeployment().GetId(), actions["deployment.create"].Resource)
	assert.Equal(t, "project:p1", actions["stack.delete"].Resource)
	assert.Contains(t, actions["stack.delete"].Metadata, `"delete_volumes":true`)

	ts.stacks.mu.Lock()
	defer ts.stacks.mu.Unlock()
	assert.Equal(t, []string{"forgenet-p1"}, ts.stacks.removed)
}
