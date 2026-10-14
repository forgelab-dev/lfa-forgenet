package docker_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/docker"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

type mockDockerClient struct {
	listFn    func(ctx context.Context, options container.ListOptions) ([]types.Container, error)
	inspectFn func(ctx context.Context, containerID string) (types.ContainerJSON, error)
	createFn  func(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error)
	startFn   func(ctx context.Context, containerID string, options container.StartOptions) error
	stopFn    func(ctx context.Context, containerID string, options container.StopOptions) error
	restartFn func(ctx context.Context, containerID string, options container.StopOptions) error
	removeFn  func(ctx context.Context, containerID string, options container.RemoveOptions) error
	logsFn    func(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error)
	pingFn    func(ctx context.Context) (types.Ping, error)
	closeFn   func() error
}

func (m *mockDockerClient) ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
	if m.listFn != nil {
		return m.listFn(ctx, options)
	}
	return nil, nil
}

func (m *mockDockerClient) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	if m.inspectFn != nil {
		return m.inspectFn(ctx, containerID)
	}
	return types.ContainerJSON{}, nil
}

func (m *mockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error) {
	if m.createFn != nil {
		return m.createFn(ctx, config, hostConfig, networkingConfig, platform, containerName)
	}
	return container.CreateResponse{ID: "created-id-123"}, nil
}

func (m *mockDockerClient) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
	if m.startFn != nil {
		return m.startFn(ctx, containerID, options)
	}
	return nil
}

func (m *mockDockerClient) ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error {
	if m.stopFn != nil {
		return m.stopFn(ctx, containerID, options)
	}
	return nil
}

func (m *mockDockerClient) ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error {
	if m.restartFn != nil {
		return m.restartFn(ctx, containerID, options)
	}
	return nil
}

func (m *mockDockerClient) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error {
	if m.removeFn != nil {
		return m.removeFn(ctx, containerID, options)
	}
	return nil
}

func (m *mockDockerClient) ContainerLogs(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error) {
	if m.logsFn != nil {
		return m.logsFn(ctx, container, options)
	}
	return io.NopCloser(bytes.NewReader(nil)), nil
}

func (m *mockDockerClient) Ping(ctx context.Context) (types.Ping, error) {
	if m.pingFn != nil {
		return m.pingFn(ctx)
	}
	return types.Ping{}, nil
}

func (m *mockDockerClient) Close() error {
	if m.closeFn != nil {
		return m.closeFn()
	}
	return nil
}

func TestDockerService_Ping(t *testing.T) {
	t.Run("Succès", func(t *testing.T) {
		mock := &mockDockerClient{
			pingFn: func(_ context.Context) (types.Ping, error) {
				return types.Ping{APIVersion: "1.45"}, nil
			},
		}
		svc := docker.NewWithClient(mock)
		err := svc.Ping(context.Background())
		assert.NoError(t, err)
	})

	t.Run("Démon indisponible", func(t *testing.T) {
		mock := &mockDockerClient{
			pingFn: func(_ context.Context) (types.Ping, error) {
				return types.Ping{}, errors.New("connection refused")
			},
		}
		svc := docker.NewWithClient(mock)
		err := svc.Ping(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "démon docker injoignable")
	})
}

func TestDockerService_ListContainers(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	mock := &mockDockerClient{
		listFn: func(_ context.Context, options container.ListOptions) ([]types.Container, error) {
			assert.True(t, options.All)
			return []types.Container{
				{
					ID:      "cont-1",
					Names:   []string{"/web-service"},
					Image:   "nginx:alpine",
					State:   "running",
					Created: now.Unix(),
					Ports: []types.Port{
						{
							PrivatePort: 80,
							PublicPort:  8080,
							Type:        "tcp",
						},
					},
				},
				{
					ID:      "cont-2",
					Names:   []string{"db-service"},
					Image:   "postgres:16",
					State:   "exited",
					Created: now.Unix(),
				},
			}, nil
		},
	}

	svc := docker.NewWithClient(mock)
	containers, err := svc.ListContainers(context.Background(), true)
	require.NoError(t, err)
	require.Len(t, containers, 2)

	assert.Equal(t, "cont-1", containers[0].ID)
	assert.Equal(t, "web-service", containers[0].Name)
	assert.Equal(t, "nginx:alpine", containers[0].Image)
	assert.Equal(t, forgenetv1.ContainerState_CONTAINER_STATE_RUNNING, containers[0].State)
	require.Len(t, containers[0].Ports, 1)
	assert.Equal(t, uint32(8080), containers[0].Ports[0].HostPort)
	assert.Equal(t, uint32(80), containers[0].Ports[0].ContainerPort)
	assert.Equal(t, "tcp", containers[0].Ports[0].Protocol)

	assert.Equal(t, "cont-2", containers[1].ID)
	assert.Equal(t, "db-service", containers[1].Name)
	assert.Equal(t, forgenetv1.ContainerState_CONTAINER_STATE_EXITED, containers[1].State)

	// Test proto conversion
	proto := containers[0].ToProto()
	require.NotNil(t, proto)
	assert.Equal(t, "cont-1", proto.Id)
	assert.Equal(t, "web-service", proto.Name)
	assert.Equal(t, forgenetv1.ContainerState_CONTAINER_STATE_RUNNING, proto.State)
}

func TestDockerService_GetContainer(t *testing.T) {
	t.Run("Succès", func(t *testing.T) {
		mock := &mockDockerClient{
			inspectFn: func(_ context.Context, containerID string) (types.ContainerJSON, error) {
				assert.Equal(t, "cont-123", containerID)
				return types.ContainerJSON{
					ContainerJSONBase: &types.ContainerJSONBase{
						ID:      "cont-123",
						Name:    "/my-app",
						Created: "2026-03-01T12:00:00Z",
						State: &types.ContainerState{
							Status: "running",
						},
						Image: "sha256:abcd",
					},
					Config: &container.Config{
						Image: "my-app:v1",
					},
				}, nil
			},
		}

		svc := docker.NewWithClient(mock)
		info, err := svc.GetContainer(context.Background(), "cont-123")
		require.NoError(t, err)
		assert.Equal(t, "cont-123", info.ID)
		assert.Equal(t, "my-app", info.Name)
		assert.Equal(t, "my-app:v1", info.Image)
		assert.Equal(t, forgenetv1.ContainerState_CONTAINER_STATE_RUNNING, info.State)
	})

	t.Run("Conteneur non trouvé", func(t *testing.T) {
		mock := &mockDockerClient{
			inspectFn: func(_ context.Context, _ string) (types.ContainerJSON, error) {
				return types.ContainerJSON{}, errors.New("Error: No such container: missing-id")
			},
		}

		svc := docker.NewWithClient(mock)
		_, err := svc.GetContainer(context.Background(), "missing-id")
		assert.ErrorIs(t, err, docker.ErrContainerNotFound)
	})

	t.Run("ID vide", func(t *testing.T) {
		svc := docker.NewWithClient(&mockDockerClient{})
		_, err := svc.GetContainer(context.Background(), "   ")
		assert.Error(t, err)
	})
}

func TestDockerService_Lifecycle(t *testing.T) {
	mock := &mockDockerClient{
		inspectFn: func(_ context.Context, containerID string) (types.ContainerJSON, error) {
			return types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					ID:      containerID,
					Name:    "/test-app",
					Created: "2026-03-01T12:00:00Z",
					State: &types.ContainerState{
						Status: "running",
					},
				},
				Config: &container.Config{Image: "test:latest"},
			}, nil
		},
		startFn: func(_ context.Context, containerID string, _ container.StartOptions) error {
			assert.Equal(t, "c1", containerID)
			return nil
		},
		stopFn: func(_ context.Context, containerID string, options container.StopOptions) error {
			assert.Equal(t, "c1", containerID)
			require.NotNil(t, options.Timeout)
			assert.Equal(t, 10, *options.Timeout)
			return nil
		},
		restartFn: func(_ context.Context, containerID string, _ container.StopOptions) error {
			assert.Equal(t, "c1", containerID)
			return nil
		},
		removeFn: func(_ context.Context, containerID string, options container.RemoveOptions) error {
			assert.Equal(t, "c1", containerID)
			assert.True(t, options.Force)
			return nil
		},
	}

	svc := docker.NewWithClient(mock)

	// Start
	started, err := svc.StartContainer(context.Background(), "c1")
	require.NoError(t, err)
	assert.Equal(t, "c1", started.ID)

	// Stop
	stopped, err := svc.StopContainer(context.Background(), "c1", 10)
	require.NoError(t, err)
	assert.Equal(t, "c1", stopped.ID)

	// Restart
	restarted, err := svc.RestartContainer(context.Background(), "c1")
	require.NoError(t, err)
	assert.Equal(t, "c1", restarted.ID)

	// Delete
	err = svc.DeleteContainer(context.Background(), "c1", true)
	assert.NoError(t, err)
}

func TestDockerService_CreateContainer(t *testing.T) {
	var capturedConfig *container.Config
	var capturedHostConfig *container.HostConfig

	mock := &mockDockerClient{
		createFn: func(_ context.Context, config *container.Config, hostConfig *container.HostConfig, _ *network.NetworkingConfig, _ *ocispec.Platform, containerName string) (container.CreateResponse, error) {
			capturedConfig = config
			capturedHostConfig = hostConfig
			assert.Equal(t, "mon-conteneur", containerName)
			return container.CreateResponse{ID: "nouveau-id-789"}, nil
		},
		inspectFn: func(_ context.Context, containerID string) (types.ContainerJSON, error) {
			assert.Equal(t, "nouveau-id-789", containerID)
			return types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					ID:      "nouveau-id-789",
					Name:    "/mon-conteneur",
					Created: "2026-03-01T12:00:00Z",
					State: &types.ContainerState{
						Status: "created",
					},
				},
				Config: &container.Config{Image: "redis:7-alpine"},
			}, nil
		},
	}

	svc := docker.NewWithClient(mock)

	req := docker.CreateContainerRequest{
		Name:  "mon-conteneur",
		Image: "redis:7-alpine",
		Env:   []string{"FOO=BAR"},
		Ports: []docker.PortMapping{
			{
				HostPort:      6379,
				ContainerPort: 6379,
				Protocol:      "tcp",
			},
		},
		MemoryLimit: 512 * 1024 * 1024,
		CPULimit:    500, // 0.5 CPU
	}

	created, err := svc.CreateContainer(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "nouveau-id-789", created.ID)
	assert.Equal(t, "mon-conteneur", created.Name)
	assert.Equal(t, forgenetv1.ContainerState_CONTAINER_STATE_CREATED, created.State)

	require.NotNil(t, capturedConfig)
	assert.Equal(t, "redis:7-alpine", capturedConfig.Image)
	assert.Equal(t, []string{"FOO=BAR"}, capturedConfig.Env)

	require.NotNil(t, capturedHostConfig)
	assert.Equal(t, int64(512*1024*1024), capturedHostConfig.Memory)
	assert.Equal(t, int64(500*1_000_000), capturedHostConfig.NanoCPUs)
}

func TestDockerService_StreamLogs(t *testing.T) {
	logContent := "log line 1\nlog line 2\n"
	mock := &mockDockerClient{
		logsFn: func(_ context.Context, container string, options container.LogsOptions) (io.ReadCloser, error) {
			assert.Equal(t, "c-logs", container)
			assert.True(t, options.ShowStdout)
			assert.True(t, options.ShowStderr)
			assert.True(t, options.Follow)
			assert.Equal(t, "100", options.Tail)
			return io.NopCloser(bytes.NewBufferString(logContent)), nil
		},
	}

	svc := docker.NewWithClient(mock)
	reader, err := svc.StreamLogs(context.Background(), "c-logs", true, 100, time.Time{}, true, true)
	require.NoError(t, err)
	defer func() {
		_ = reader.Close()
	}()

	buf, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, logContent, string(buf))
}

func TestParseContainerState(t *testing.T) {
	tests := []struct {
		input    string
		expected forgenetv1.ContainerState
	}{
		{"created", forgenetv1.ContainerState_CONTAINER_STATE_CREATED},
		{"RUNNING", forgenetv1.ContainerState_CONTAINER_STATE_RUNNING},
		{"paused", forgenetv1.ContainerState_CONTAINER_STATE_PAUSED},
		{"restarting", forgenetv1.ContainerState_CONTAINER_STATE_RESTARTING},
		{"exited", forgenetv1.ContainerState_CONTAINER_STATE_EXITED},
		{"dead", forgenetv1.ContainerState_CONTAINER_STATE_DEAD},
		{"unknown", forgenetv1.ContainerState_CONTAINER_STATE_UNSPECIFIED},
		{"", forgenetv1.ContainerState_CONTAINER_STATE_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, docker.ParseContainerState(tt.input))
		})
	}
}
