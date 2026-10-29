package grpc_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/health"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// fakeEngine simule le moteur : containers pour la liste, err pour les actions.
type fakeEngine struct {
	orchestrator.Engine
	containers []orchestrator.Container
	err        error
	logs       []orchestrator.LogChunk
	stopped    time.Duration
}

func (f *fakeEngine) List(context.Context, bool) ([]orchestrator.Container, error) {
	return f.containers, f.err
}

func (f *fakeEngine) Get(_ context.Context, id string) (orchestrator.Container, error) {
	return orchestrator.Container{ID: id, State: orchestrator.StateRunning}, f.err
}

func (f *fakeEngine) Start(ctx context.Context, id string) (orchestrator.Container, error) {
	return f.Get(ctx, id)
}

func (f *fakeEngine) Stop(ctx context.Context, id string, timeout time.Duration) (orchestrator.Container, error) {
	f.stopped = timeout
	c, err := f.Get(ctx, id)
	c.State = orchestrator.StateExited
	return c, err
}

func (f *fakeEngine) Restart(ctx context.Context, id string) (orchestrator.Container, error) {
	return f.Get(ctx, id)
}

func (f *fakeEngine) Logs(_ context.Context, _ string, _ orchestrator.LogOptions, emit func(orchestrator.LogChunk) error) error {
	for _, c := range f.logs {
		if err := emit(c); err != nil {
			return err
		}
	}
	return f.err
}

// testServer démarre un vrai serveur gRPC (garde et politique compris) devant
// le faux moteur, avec un compte admin et un compte membre.
type testServer struct {
	conn        *grpc.ClientConn
	containers  forgenetv1.ContainerServiceClient
	adminToken  string
	memberToken string
	store       store.Store
}

// registerFunc ajoute un service au serveur de test, avec sa base et son audit.
type registerFunc func(srv *grpc.Server, st store.Store, auditor audit.Auditor)

func startTestServer(t *testing.T, engine orchestrator.Engine, register ...registerFunc) *testServer {
	t.Helper()
	ctx := context.Background()
	st, err := sqlite.Open(ctx, config.DatabaseConfig{Path: ":memory:", BusyTimeout: 2 * time.Second})
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })

	auditor := audit.NewAuditor(st.Audit(), zerolog.Nop())
	authSvc := auth.NewService(st, config.AuthConfig{JWTSecret: "secret-de-test-d-au-moins-32-caracteres"}, auditor)
	token := func(name string, role store.Role) string {
		u, err := authSvc.Register(ctx, name, name+"@forgenet.test", "motdepasse-de-test", role)
		require.NoError(t, err)
		tok, _, err := authSvc.CreateAPIToken(ctx, u.ID, name, "", time.Hour)
		require.NoError(t, err)
		return tok
	}

	srv := transportGRPC.NewServer(config.ServerConfig{}, health.NewService(), auth.NewGuard(authSvc, auth.DefaultPolicy(), zerolog.Nop()), zerolog.Nop())
	forgenetv1.RegisterContainerServiceServer(srv.GRPCServer(), transportGRPC.NewContainerServiceServer(engine, auditor, zerolog.Nop()))
	for _, r := range register {
		r(srv.GRPCServer(), st, auditor)
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	srv.SetListener(ln)
	go func() { _ = srv.Start() }()
	t.Cleanup(func() { _ = srv.Shutdown(context.Background()) })

	conn, err := grpc.NewClient(srv.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	return &testServer{
		conn:        conn,
		containers:  forgenetv1.NewContainerServiceClient(conn),
		adminToken:  token("admin", store.RoleAdmin),
		memberToken: token("membre", store.RoleMember),
		store:       st,
	}
}

func as(token string) context.Context {
	return metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)
}

func TestContainerService_RolePolicy(t *testing.T) {
	ts := startTestServer(t, &fakeEngine{})
	id := "c1"

	calls := map[string]func(ctx context.Context) error{
		"List": func(ctx context.Context) error {
			_, err := ts.containers.ListContainers(ctx, &forgenetv1.ListContainersRequest{})
			return err
		},
		"Get": func(ctx context.Context) error {
			_, err := ts.containers.GetContainer(ctx, &forgenetv1.GetContainerRequest{ContainerId: id})
			return err
		},
		"Logs": func(ctx context.Context) error {
			stream, err := ts.containers.StreamContainerLogs(ctx, &forgenetv1.StreamContainerLogsRequest{ContainerId: id})
			if err != nil {
				return err
			}
			_, err = stream.Recv()
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		},
		"Start": func(ctx context.Context) error {
			_, err := ts.containers.StartContainer(ctx, &forgenetv1.StartContainerRequest{ContainerId: id})
			return err
		},
		"Stop": func(ctx context.Context) error {
			_, err := ts.containers.StopContainer(ctx, &forgenetv1.StopContainerRequest{ContainerId: id})
			return err
		},
		"Restart": func(ctx context.Context) error {
			_, err := ts.containers.RestartContainer(ctx, &forgenetv1.RestartContainerRequest{ContainerId: id})
			return err
		},
	}
	// Décision du 14/10 : les membres consultent, seuls les admins agissent.
	memberAllowed := map[string]bool{"List": true, "Get": true, "Logs": true}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, codes.Unauthenticated, status.Code(call(context.Background())), "sans jeton")
			assert.Equal(t, codes.OK, status.Code(call(as(ts.adminToken))), "admin")

			want := codes.PermissionDenied
			if memberAllowed[name] {
				want = codes.OK
			}
			assert.Equal(t, want, status.Code(call(as(ts.memberToken))), "membre")
		})
	}
}

func TestContainerService_ListPagination(t *testing.T) {
	engine := &fakeEngine{}
	for i := range 5 {
		engine.containers = append(engine.containers, orchestrator.Container{ID: fmt.Sprintf("c%d", i)})
	}
	ts := startTestServer(t, engine)
	ctx := as(ts.memberToken)

	var ids []string
	token := ""
	for {
		resp, err := ts.containers.ListContainers(ctx, &forgenetv1.ListContainersRequest{PageSize: 2, PageToken: token})
		require.NoError(t, err)
		for _, c := range resp.GetContainers() {
			ids = append(ids, c.GetId())
		}
		if token = resp.GetNextPageToken(); token == "" {
			break
		}
	}
	assert.Equal(t, []string{"c0", "c1", "c2", "c3", "c4"}, ids)

	_, err := ts.containers.ListContainers(ctx, &forgenetv1.ListContainersRequest{PageToken: "abc"})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestContainerService_ErrorTranslation(t *testing.T) {
	cases := []struct {
		err  error
		want codes.Code
	}{
		{fmt.Errorf("inspection : %w", orchestrator.ErrNotFound), codes.NotFound},
		{fmt.Errorf("x : %w", orchestrator.ErrInvalidArgument), codes.InvalidArgument},
		{errors.New("dial unix /var/run/docker.sock: connection refused"), codes.Internal},
	}
	for _, tc := range cases {
		t.Run(tc.want.String(), func(t *testing.T) {
			ts := startTestServer(t, &fakeEngine{err: tc.err})
			_, err := ts.containers.GetContainer(as(ts.adminToken), &forgenetv1.GetContainerRequest{ContainerId: "c1"})
			assert.Equal(t, tc.want, status.Code(err))
			assert.NotContains(t, status.Convert(err).Message(), "docker.sock", "aucun détail interne n'est renvoyé")
		})
	}
}

func TestContainerService_ActionsAreAudited(t *testing.T) {
	engine := &fakeEngine{}
	ts := startTestServer(t, engine)

	resp, err := ts.containers.StopContainer(as(ts.adminToken), &forgenetv1.StopContainerRequest{ContainerId: "c1", TimeoutSeconds: 7})
	require.NoError(t, err)
	assert.Equal(t, forgenetv1.ContainerState_CONTAINER_STATE_EXITED, resp.GetContainer().GetState())
	assert.Equal(t, 7*time.Second, engine.stopped)

	events, err := ts.store.Audit().List(context.Background(), 50, 0)
	require.NoError(t, err)
	var stop *store.AuditEvent
	for _, e := range events {
		if e.Action == "container.stop" {
			stop = e
		}
	}
	require.NotNil(t, stop, "l'arrêt doit figurer au journal d'audit")
	assert.Equal(t, "container:c1", stop.Resource)
	assert.NotEmpty(t, stop.UserID)
	assert.Contains(t, stop.IPAddress, "127.0.0.1", "l'adresse de l'appelant est conservée")
	assert.Contains(t, stop.Metadata, `"success":true`)
}

func TestContainerService_StreamLogs(t *testing.T) {
	engine := &fakeEngine{logs: []orchestrator.LogChunk{
		{Stream: orchestrator.Stdout, Data: []byte("bonjour\n")},
		{Stream: orchestrator.Stderr, Data: []byte("attention\n")},
	}}
	ts := startTestServer(t, engine)

	stream, err := ts.containers.StreamContainerLogs(as(ts.memberToken), &forgenetv1.StreamContainerLogsRequest{ContainerId: "c1"})
	require.NoError(t, err)

	var got []*forgenetv1.StreamContainerLogsResponse
	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		got = append(got, msg)
	}
	require.Len(t, got, 2)
	assert.Equal(t, "bonjour\n", string(got[0].GetChunk()))
	assert.Equal(t, forgenetv1.LogStream_LOG_STREAM_STDOUT, got[0].GetSource())
	assert.Equal(t, forgenetv1.LogStream_LOG_STREAM_STDERR, got[1].GetSource())
}
