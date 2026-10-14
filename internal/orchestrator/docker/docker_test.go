package docker_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net/netip"
	"testing"
	"time"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/docker"
)

// fakeClient enregistre les options reçues et renvoie les réponses préparées.
type fakeClient struct {
	docker.EngineClient // méthodes non utilisées : panique si appelées

	inspect     container.InspectResponse
	err         error // renvoyée par toutes les actions
	created     client.ContainerCreateOptions
	stopOpts    client.ContainerStopOptions
	logsOpts    client.ContainerLogsOptions
	logs        []byte
	removedID   string
	forceRemove bool
}

func (f *fakeClient) ContainerInspect(context.Context, string, client.ContainerInspectOptions) (client.ContainerInspectResult, error) {
	return client.ContainerInspectResult{Container: f.inspect}, f.err
}

func (f *fakeClient) ContainerCreate(_ context.Context, opts client.ContainerCreateOptions) (client.ContainerCreateResult, error) {
	f.created = opts
	return client.ContainerCreateResult{ID: "c1"}, f.err
}

func (f *fakeClient) ContainerStart(context.Context, string, client.ContainerStartOptions) (client.ContainerStartResult, error) {
	return client.ContainerStartResult{}, f.err
}

func (f *fakeClient) ContainerStop(_ context.Context, _ string, opts client.ContainerStopOptions) (client.ContainerStopResult, error) {
	f.stopOpts = opts
	return client.ContainerStopResult{}, f.err
}

func (f *fakeClient) ContainerRestart(context.Context, string, client.ContainerRestartOptions) (client.ContainerRestartResult, error) {
	return client.ContainerRestartResult{}, f.err
}

func (f *fakeClient) ContainerRemove(_ context.Context, id string, opts client.ContainerRemoveOptions) (client.ContainerRemoveResult, error) {
	f.removedID, f.forceRemove = id, opts.Force
	return client.ContainerRemoveResult{}, f.err
}

func (f *fakeClient) ContainerLogs(_ context.Context, _ string, opts client.ContainerLogsOptions) (client.ContainerLogsResult, error) {
	f.logsOpts = opts
	return io.NopCloser(bytes.NewReader(f.logs)), f.err
}

func newEngine(t *testing.T, f *fakeClient) *docker.Engine {
	t.Helper()
	e, err := docker.NewWithClient(f, "127.0.0.1")
	require.NoError(t, err)
	return e
}

func TestNewWithClient_RejectsInvalidPublishHost(t *testing.T) {
	_, err := docker.NewWithClient(&fakeClient{}, "pas-une-ip")
	require.Error(t, err)
}

func TestCreate_PublishesPortsOnConfiguredAddress(t *testing.T) {
	f := &fakeClient{}
	e := newEngine(t, f)

	_, err := e.Create(context.Background(), orchestrator.CreateSpec{
		Name:  "web",
		Image: "nginx:alpine",
		Ports: []orchestrator.Port{
			{ContainerPort: 80, HostPort: 8080},                     // adresse par défaut
			{ContainerPort: 443, HostPort: 8443, HostIP: "0.0.0.0"}, // exposition explicite
			{ContainerPort: 9000, Protocol: "udp"},                  // exposé, non publié
		},
		MemoryBytes: 64 << 20,
		MilliCPU:    500,
	})
	require.NoError(t, err)

	cfg, host := f.created.Config, f.created.HostConfig
	assert.Equal(t, "web", f.created.Name)
	assert.Len(t, cfg.ExposedPorts, 3)
	assert.Equal(t, []network.PortBinding{{HostIP: netip.MustParseAddr("127.0.0.1"), HostPort: "8080"}},
		host.PortBindings[network.MustParsePort("80/tcp")])
	assert.Equal(t, netip.MustParseAddr("0.0.0.0"), host.PortBindings[network.MustParsePort("443/tcp")][0].HostIP)
	assert.NotContains(t, host.PortBindings, network.MustParsePort("9000/udp"))
	assert.Equal(t, int64(64<<20), host.Memory)
	assert.Equal(t, int64(500_000_000), host.NanoCPUs, "500 milli-CPU = un demi-cœur")
}

func TestCreate_RejectsInvalidSpec(t *testing.T) {
	e := newEngine(t, &fakeClient{})
	for name, spec := range map[string]orchestrator.CreateSpec{
		"sans image":        {Name: "x"},
		"protocole faux":    {Image: "nginx", Ports: []orchestrator.Port{{ContainerPort: 80, Protocol: "http"}}},
		"adresse illisible": {Image: "nginx", Ports: []orchestrator.Port{{ContainerPort: 80, HostPort: 80, HostIP: "ici"}}},
		"port nul":          {Image: "nginx", Ports: []orchestrator.Port{{ContainerPort: 0}}},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := e.Create(context.Background(), spec)
			assert.ErrorIs(t, err, orchestrator.ErrInvalidArgument)
		})
	}
}

func TestActions_TranslateEngineErrors(t *testing.T) {
	ctx := context.Background()
	actions := map[string]func(*docker.Engine, string) error{
		"Get":     func(e *docker.Engine, id string) error { _, err := e.Get(ctx, id); return err },
		"Start":   func(e *docker.Engine, id string) error { _, err := e.Start(ctx, id); return err },
		"Stop":    func(e *docker.Engine, id string) error { _, err := e.Stop(ctx, id, 0); return err },
		"Restart": func(e *docker.Engine, id string) error { _, err := e.Restart(ctx, id); return err },
		"Remove":  func(e *docker.Engine, id string) error { return e.Remove(ctx, id, false) },
	}
	for name, act := range actions {
		t.Run(name, func(t *testing.T) {
			assert.ErrorIs(t, act(newEngine(t, &fakeClient{err: cerrdefs.ErrNotFound}), "absent"), orchestrator.ErrNotFound)
			assert.ErrorIs(t, act(newEngine(t, &fakeClient{}), "  "), orchestrator.ErrInvalidArgument)

			other := errors.New("démon en panne")
			err := act(newEngine(t, &fakeClient{err: other}), "c1")
			assert.ErrorIs(t, err, other)
			assert.NotErrorIs(t, err, orchestrator.ErrNotFound)
		})
	}
}

func TestStop_Timeout(t *testing.T) {
	f := &fakeClient{}
	e := newEngine(t, f)

	_, err := e.Stop(context.Background(), "c1", 0)
	require.NoError(t, err)
	assert.Nil(t, f.stopOpts.Timeout, "0 : délai du conteneur ou du moteur")

	_, err = e.Stop(context.Background(), "c1", 5*time.Second)
	require.NoError(t, err)
	require.NotNil(t, f.stopOpts.Timeout)
	assert.Equal(t, 5, *f.stopOpts.Timeout)
}

func TestRemove_Force(t *testing.T) {
	f := &fakeClient{}
	require.NoError(t, newEngine(t, f).Remove(context.Background(), " c1 ", true))
	assert.Equal(t, "c1", f.removedID)
	assert.True(t, f.forceRemove)
}

func collectLogs(t *testing.T, e *docker.Engine, opts orchestrator.LogOptions) map[orchestrator.LogStream]string {
	t.Helper()
	got := map[orchestrator.LogStream]string{}
	err := e.Logs(context.Background(), "c1", opts, func(c orchestrator.LogChunk) error {
		got[c.Stream] += string(c.Data)
		return nil
	})
	require.NoError(t, err)
	return got
}

// frame construit une trame du flux multiplexé de Docker : un en-tête de 8
// octets (sortie, trois zéros, taille en big-endian) suivi des données.
func frame(stream stdcopy.StdType, data string) []byte {
	header := []byte{byte(stream), 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint32(header[4:], uint32(len(data)))
	return append(header, data...)
}

func TestLogs_DemultiplexesWithoutTTY(t *testing.T) {
	mux := append(frame(stdcopy.Stdout, "bonjour\n"), frame(stdcopy.Stderr, "attention\n")...)
	f := &fakeClient{inspect: container.InspectResponse{ID: "c1", Config: &container.Config{}}, logs: mux}
	got := collectLogs(t, newEngine(t, f), orchestrator.LogOptions{Stdout: true, Stderr: true, Tail: 10})

	assert.Equal(t, "bonjour\n", got[orchestrator.Stdout])
	assert.Equal(t, "attention\n", got[orchestrator.Stderr])
	assert.Equal(t, "10", f.logsOpts.Tail)
}

func TestLogs_RawStreamWithTTY(t *testing.T) {
	// Avec un TTY, Docker n'ajoute aucun en-tête : ces octets, qui ressemblent
	// à un en-tête multiplexé, doivent passer tels quels.
	raw := []byte("\x01\x00\x00\x00\x00\x00\x00\x02ok")
	f := &fakeClient{inspect: container.InspectResponse{ID: "c1", Config: &container.Config{Tty: true}}, logs: raw}
	got := collectLogs(t, newEngine(t, f), orchestrator.LogOptions{Stdout: true})

	assert.Equal(t, string(raw), got[orchestrator.Stdout])
	assert.Empty(t, got[orchestrator.Stderr])
	assert.Equal(t, "all", f.logsOpts.Tail)
}

func TestLogs_StopsWhenEmitFails(t *testing.T) {
	f := &fakeClient{inspect: container.InspectResponse{ID: "c1", Config: &container.Config{}}, logs: frame(stdcopy.Stdout, "ligne\n")}

	clientGone := errors.New("client parti")
	err := newEngine(t, f).Logs(context.Background(), "c1", orchestrator.LogOptions{Stdout: true},
		func(orchestrator.LogChunk) error { return clientGone })
	assert.ErrorIs(t, err, clientGone)
}

func TestGet_ToleratesMissingFields(t *testing.T) {
	f := &fakeClient{inspect: container.InspectResponse{ID: "c1", Name: "/web", Image: "sha256:abc"}}
	c, err := newEngine(t, f).Get(context.Background(), "c1")
	require.NoError(t, err)
	assert.Equal(t, orchestrator.Container{ID: "c1", Name: "web", Image: "sha256:abc"}, c)
}

func TestGet_MapsStateAndPorts(t *testing.T) {
	f := &fakeClient{inspect: container.InspectResponse{
		ID:      "c1",
		Name:    "/web",
		Created: "2026-10-14T12:00:00.5Z",
		State:   &container.State{Status: container.StateRunning},
		Config:  &container.Config{Image: "nginx:alpine"},
		NetworkSettings: &container.NetworkSettings{Ports: network.PortMap{
			network.MustParsePort("443/tcp"): nil,
			network.MustParsePort("80/tcp"):  {{HostIP: netip.MustParseAddr("127.0.0.1"), HostPort: "8080"}},
		}},
	}}
	c, err := newEngine(t, f).Get(context.Background(), "c1")
	require.NoError(t, err)

	assert.Equal(t, orchestrator.StateRunning, c.State)
	assert.Equal(t, "nginx:alpine", c.Image)
	assert.Equal(t, time.Date(2026, 10, 14, 12, 0, 0, 500_000_000, time.UTC), c.CreatedAt)
	assert.Equal(t, []orchestrator.Port{
		{HostIP: "127.0.0.1", HostPort: 8080, ContainerPort: 80, Protocol: "tcp"},
		{ContainerPort: 443, Protocol: "tcp"},
	}, c.Ports, "ports triés par port du conteneur")
}
