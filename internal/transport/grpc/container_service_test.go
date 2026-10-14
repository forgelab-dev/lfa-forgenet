package grpc_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/docker"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	transportgrpc "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

type mockOrchestrator struct {
	listFn    func(ctx context.Context, includeStopped bool) ([]*docker.ContainerInfo, error)
	getFn     func(ctx context.Context, containerID string) (*docker.ContainerInfo, error)
	createFn  func(ctx context.Context, req docker.CreateContainerRequest) (*docker.ContainerInfo, error)
	startFn   func(ctx context.Context, containerID string) (*docker.ContainerInfo, error)
	stopFn    func(ctx context.Context, containerID string, timeoutSeconds uint32) (*docker.ContainerInfo, error)
	restartFn func(ctx context.Context, containerID string) (*docker.ContainerInfo, error)
	deleteFn  func(ctx context.Context, containerID string, force bool) error
	streamFn  func(ctx context.Context, containerID string, follow bool, tail int, since time.Time, stdout, stderr bool) (io.ReadCloser, error)
	pingFn    func(ctx context.Context) error
}

func (m *mockOrchestrator) ListContainers(ctx context.Context, includeStopped bool) ([]*docker.ContainerInfo, error) {
	if m.listFn != nil {
		return m.listFn(ctx, includeStopped)
	}
	return nil, nil
}

func (m *mockOrchestrator) GetContainer(ctx context.Context, containerID string) (*docker.ContainerInfo, error) {
	if m.getFn != nil {
		return m.getFn(ctx, containerID)
	}
	return nil, docker.ErrContainerNotFound
}

func (m *mockOrchestrator) CreateContainer(ctx context.Context, req docker.CreateContainerRequest) (*docker.ContainerInfo, error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	return nil, nil
}

func (m *mockOrchestrator) StartContainer(ctx context.Context, containerID string) (*docker.ContainerInfo, error) {
	if m.startFn != nil {
		return m.startFn(ctx, containerID)
	}
	return nil, nil
}

func (m *mockOrchestrator) StopContainer(ctx context.Context, containerID string, timeoutSeconds uint32) (*docker.ContainerInfo, error) {
	if m.stopFn != nil {
		return m.stopFn(ctx, containerID, timeoutSeconds)
	}
	return nil, nil
}

func (m *mockOrchestrator) RestartContainer(ctx context.Context, containerID string) (*docker.ContainerInfo, error) {
	if m.restartFn != nil {
		return m.restartFn(ctx, containerID)
	}
	return nil, nil
}

func (m *mockOrchestrator) DeleteContainer(ctx context.Context, containerID string, force bool) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, containerID, force)
	}
	return nil
}

func (m *mockOrchestrator) StreamLogs(ctx context.Context, containerID string, follow bool, tail int, since time.Time, stdout, stderr bool) (io.ReadCloser, error) {
	if m.streamFn != nil {
		return m.streamFn(ctx, containerID, follow, tail, since, stdout, stderr)
	}
	return io.NopCloser(bytes.NewReader(nil)), nil
}

func (m *mockOrchestrator) Ping(ctx context.Context) error {
	if m.pingFn != nil {
		return m.pingFn(ctx)
	}
	return nil
}

type mockAuditor struct {
	loggedEvents []string
}

func (m *mockAuditor) Log(_ context.Context, action, resource string, _ map[string]any) error {
	m.loggedEvents = append(m.loggedEvents, action+":"+resource)
	return nil
}

func (m *mockAuditor) List(_ context.Context, _, _ int) ([]*store.AuditEvent, error) {
	return nil, nil
}

type mockStreamLogsServer struct {
	grpc.ServerStream
	ctx      context.Context
	sentMsgs []*forgenetv1.StreamContainerLogsResponse
}

func (m *mockStreamLogsServer) Context() context.Context {
	return m.ctx
}

func (m *mockStreamLogsServer) Send(resp *forgenetv1.StreamContainerLogsResponse) error {
	m.sentMsgs = append(m.sentMsgs, resp)
	return nil
}

func contextWithUser(role store.Role) context.Context {
	return auth.WithUserContext(context.Background(), &auth.Claims{
		UserID:   "usr_test",
		Username: "operator",
		Role:     role,
	})
}

func TestContainerService_ListContainers(t *testing.T) {
	logger := zerolog.Nop()

	t.Run("Refus si non authentifié", func(t *testing.T) {
		svc := transportgrpc.NewContainerServiceServer(&mockOrchestrator{}, &mockAuditor{}, logger)
		_, err := svc.ListContainers(context.Background(), &forgenetv1.ListContainersRequest{})
		require.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("Succès paginé", func(t *testing.T) {
		orch := &mockOrchestrator{
			listFn: func(_ context.Context, _ bool) ([]*docker.ContainerInfo, error) {
				var list []*docker.ContainerInfo
				for i := 1; i <= 3; i++ {
					list = append(list, &docker.ContainerInfo{
						ID:        "c" + string(rune('0'+i)),
						Name:      "app-" + string(rune('0'+i)),
						Image:     "nginx",
						State:     forgenetv1.ContainerState_CONTAINER_STATE_RUNNING,
						CreatedAt: time.Now().UTC(),
					})
				}
				return list, nil
			},
		}

		svc := transportgrpc.NewContainerServiceServer(orch, &mockAuditor{}, logger)
		ctx := contextWithUser(store.RoleMember)

		resp, err := svc.ListContainers(ctx, &forgenetv1.ListContainersRequest{
			PageSize: 2,
		})
		require.NoError(t, err)
		require.Len(t, resp.Containers, 2)
		assert.Equal(t, "c1", resp.Containers[0].Id)
		assert.Equal(t, "c2", resp.Containers[1].Id)
		assert.Equal(t, "2", resp.NextPageToken)

		// Deuxième page
		resp2, err := svc.ListContainers(ctx, &forgenetv1.ListContainersRequest{
			PageSize:  2,
			PageToken: resp.NextPageToken,
		})
		require.NoError(t, err)
		require.Len(t, resp2.Containers, 1)
		assert.Equal(t, "c3", resp2.Containers[0].Id)
		assert.Empty(t, resp2.NextPageToken)
	})
}

func TestContainerService_GetContainer(t *testing.T) {
	logger := zerolog.Nop()

	t.Run("Refus si non authentifié", func(t *testing.T) {
		svc := transportgrpc.NewContainerServiceServer(&mockOrchestrator{}, &mockAuditor{}, logger)
		_, err := svc.GetContainer(context.Background(), &forgenetv1.GetContainerRequest{ContainerId: "c1"})
		require.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("ID manquant", func(t *testing.T) {
		svc := transportgrpc.NewContainerServiceServer(&mockOrchestrator{}, &mockAuditor{}, logger)
		ctx := contextWithUser(store.RoleMember)
		_, err := svc.GetContainer(ctx, &forgenetv1.GetContainerRequest{ContainerId: ""})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("Non trouvé", func(t *testing.T) {
		orch := &mockOrchestrator{
			getFn: func(_ context.Context, _ string) (*docker.ContainerInfo, error) {
				return nil, docker.ErrContainerNotFound
			},
		}
		svc := transportgrpc.NewContainerServiceServer(orch, &mockAuditor{}, logger)
		ctx := contextWithUser(store.RoleMember)
		_, err := svc.GetContainer(ctx, &forgenetv1.GetContainerRequest{ContainerId: "unknown"})
		require.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("Succès", func(t *testing.T) {
		orch := &mockOrchestrator{
			getFn: func(_ context.Context, containerID string) (*docker.ContainerInfo, error) {
				return &docker.ContainerInfo{
					ID:        containerID,
					Name:      "web",
					Image:     "alpine",
					State:     forgenetv1.ContainerState_CONTAINER_STATE_RUNNING,
					CreatedAt: time.Now().UTC(),
				}, nil
			},
		}
		svc := transportgrpc.NewContainerServiceServer(orch, &mockAuditor{}, logger)
		ctx := contextWithUser(store.RoleMember)
		resp, err := svc.GetContainer(ctx, &forgenetv1.GetContainerRequest{ContainerId: "c1"})
		require.NoError(t, err)
		assert.Equal(t, "c1", resp.Container.Id)
		assert.Equal(t, "web", resp.Container.Name)
	})
}

func TestContainerService_Operations(t *testing.T) {
	logger := zerolog.Nop()

	t.Run("StartContainer avec RBAC et Audit", func(t *testing.T) {
		auditor := &mockAuditor{}
		orch := &mockOrchestrator{
			startFn: func(_ context.Context, containerID string) (*docker.ContainerInfo, error) {
				return &docker.ContainerInfo{
					ID:        containerID,
					State:     forgenetv1.ContainerState_CONTAINER_STATE_RUNNING,
					CreatedAt: time.Now().UTC(),
				}, nil
			},
		}
		svc := transportgrpc.NewContainerServiceServer(orch, auditor, logger)

		// Test rôle non autorisé (rôle vide ou guest)
		unauthCtx := contextWithUser(store.Role("guest"))
		_, err := svc.StartContainer(unauthCtx, &forgenetv1.StartContainerRequest{ContainerId: "c1"})
		require.Error(t, err)
		assert.Equal(t, codes.PermissionDenied, status.Code(err))

		// Test succès avec RoleAdmin
		adminCtx := contextWithUser(store.RoleAdmin)
		resp, err := svc.StartContainer(adminCtx, &forgenetv1.StartContainerRequest{ContainerId: "c1"})
		require.NoError(t, err)
		assert.Equal(t, "c1", resp.Container.Id)
		require.Len(t, auditor.loggedEvents, 1)
		assert.Equal(t, "container.start:container:c1", auditor.loggedEvents[0])
	})

	t.Run("StopContainer avec Audit", func(t *testing.T) {
		auditor := &mockAuditor{}
		orch := &mockOrchestrator{
			stopFn: func(_ context.Context, containerID string, timeoutSeconds uint32) (*docker.ContainerInfo, error) {
				assert.Equal(t, uint32(15), timeoutSeconds)
				return &docker.ContainerInfo{
					ID:        containerID,
					State:     forgenetv1.ContainerState_CONTAINER_STATE_EXITED,
					CreatedAt: time.Now().UTC(),
				}, nil
			},
		}
		svc := transportgrpc.NewContainerServiceServer(orch, auditor, logger)
		ctx := contextWithUser(store.RoleMember)

		resp, err := svc.StopContainer(ctx, &forgenetv1.StopContainerRequest{
			ContainerId:    "c2",
			TimeoutSeconds: 15,
		})
		require.NoError(t, err)
		assert.Equal(t, "c2", resp.Container.Id)
		require.Len(t, auditor.loggedEvents, 1)
		assert.Equal(t, "container.stop:container:c2", auditor.loggedEvents[0])
	})

	t.Run("RestartContainer avec Audit", func(t *testing.T) {
		auditor := &mockAuditor{}
		orch := &mockOrchestrator{
			restartFn: func(_ context.Context, containerID string) (*docker.ContainerInfo, error) {
				return &docker.ContainerInfo{
					ID:        containerID,
					State:     forgenetv1.ContainerState_CONTAINER_STATE_RUNNING,
					CreatedAt: time.Now().UTC(),
				}, nil
			},
		}
		svc := transportgrpc.NewContainerServiceServer(orch, auditor, logger)
		ctx := contextWithUser(store.RoleMember)

		resp, err := svc.RestartContainer(ctx, &forgenetv1.RestartContainerRequest{ContainerId: "c3"})
		require.NoError(t, err)
		assert.Equal(t, "c3", resp.Container.Id)
		require.Len(t, auditor.loggedEvents, 1)
		assert.Equal(t, "container.restart:container:c3", auditor.loggedEvents[0])
	})
}

func TestContainerService_StreamContainerLogs(t *testing.T) {
	logger := zerolog.Nop()

	t.Run("Refus si non authentifié", func(t *testing.T) {
		svc := transportgrpc.NewContainerServiceServer(&mockOrchestrator{}, &mockAuditor{}, logger)
		stream := &mockStreamLogsServer{ctx: context.Background()}
		err := svc.StreamContainerLogs(&forgenetv1.StreamContainerLogsRequest{ContainerId: "c1"}, stream)
		require.Error(t, err)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("Streaming multiplexé Docker (stdout & stderr)", func(t *testing.T) {
		var logBuf bytes.Buffer

		// Frame 1: stdout, "Hello World\n" (12 bytes)
		msg1 := []byte("Hello World\n")
		header1 := make([]byte, 8)
		header1[0] = 1 // stdout
		binary.BigEndian.PutUint32(header1[4:8], uint32(len(msg1)))
		logBuf.Write(header1)
		logBuf.Write(msg1)

		// Frame 2: stderr, "Error occurred\n" (15 bytes)
		msg2 := []byte("Error occurred\n")
		header2 := make([]byte, 8)
		header2[0] = 2 // stderr
		binary.BigEndian.PutUint32(header2[4:8], uint32(len(msg2)))
		logBuf.Write(header2)
		logBuf.Write(msg2)

		orch := &mockOrchestrator{
			streamFn: func(_ context.Context, containerID string, _ bool, _ int, _ time.Time, _, _ bool) (io.ReadCloser, error) {
				assert.Equal(t, "c-logs", containerID)
				return io.NopCloser(&logBuf), nil
			},
		}

		svc := transportgrpc.NewContainerServiceServer(orch, &mockAuditor{}, logger)
		ctx := contextWithUser(store.RoleMember)
		stream := &mockStreamLogsServer{ctx: ctx}

		err := svc.StreamContainerLogs(&forgenetv1.StreamContainerLogsRequest{
			ContainerId: "c-logs",
			Stream:      forgenetv1.LogStream_LOG_STREAM_UNSPECIFIED,
		}, stream)
		require.NoError(t, err)
		require.Len(t, stream.sentMsgs, 2)

		assert.Equal(t, "Hello World\n", string(stream.sentMsgs[0].Chunk))
		assert.Equal(t, forgenetv1.LogStream_LOG_STREAM_STDOUT, stream.sentMsgs[0].Source)

		assert.Equal(t, "Error occurred\n", string(stream.sentMsgs[1].Chunk))
		assert.Equal(t, forgenetv1.LogStream_LOG_STREAM_STDERR, stream.sentMsgs[1].Source)
	})

	t.Run("Streaming brut textuel", func(t *testing.T) {
		rawContent := "plain text log line"
		orch := &mockOrchestrator{
			streamFn: func(_ context.Context, _ string, _ bool, _ int, _ time.Time, _, _ bool) (io.ReadCloser, error) {
				return io.NopCloser(bytes.NewBufferString(rawContent)), nil
			},
		}

		svc := transportgrpc.NewContainerServiceServer(orch, &mockAuditor{}, logger)
		ctx := contextWithUser(store.RoleMember)
		stream := &mockStreamLogsServer{ctx: ctx}

		err := svc.StreamContainerLogs(&forgenetv1.StreamContainerLogsRequest{
			ContainerId: "c-raw",
		}, stream)
		require.NoError(t, err)
		require.NotEmpty(t, stream.sentMsgs)

		totalReceived := ""
		for _, m := range stream.sentMsgs {
			totalReceived += string(m.Chunk)
		}
		assert.Equal(t, rawContent, totalReceived)
	})
}
