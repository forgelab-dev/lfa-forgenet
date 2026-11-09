package docker_test

import (
	"context"
	"io"
	"strings"
	"testing"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
)

func (f *fakeClient) ContainerStats(_ context.Context, _ string, opts client.ContainerStatsOptions) (client.ContainerStatsResult, error) {
	f.statsOpts = opts
	return client.ContainerStatsResult{Body: io.NopCloser(strings.NewReader(f.stats))}, f.err
}

func (f *fakeClient) ContainerList(context.Context, client.ContainerListOptions) (client.ContainerListResult, error) {
	return client.ContainerListResult{Items: make([]container.Summary, f.running)}, f.err
}

func TestEngine_StatsDecodesDaemonResponse(t *testing.T) {
	f := &fakeClient{stats: `{"id":"c1","cpu_stats":{"online_cpus":2,"cpu_usage":{"total_usage":10}},"memory_stats":{"usage":42}}`}

	stats, err := newEngine(t, f).Stats(context.Background(), " c1 ", true)

	require.NoError(t, err)
	assert.True(t, f.statsOpts.IncludePreviousSample)
	assert.False(t, f.statsOpts.Stream, "une mesure, pas un flux du démon")
	assert.Equal(t, uint32(2), stats.CPUStats.OnlineCPUs)
	assert.Equal(t, uint64(42), stats.MemoryStats.Usage)
}

func TestEngine_StatsErrors(t *testing.T) {
	_, err := newEngine(t, &fakeClient{}).Stats(context.Background(), " ", false)
	require.ErrorIs(t, err, orchestrator.ErrInvalidArgument)

	_, err = newEngine(t, &fakeClient{err: cerrdefs.ErrNotFound}).Stats(context.Background(), "absent", false)
	require.ErrorIs(t, err, orchestrator.ErrNotFound)

	_, err = newEngine(t, &fakeClient{stats: "pas du json"}).Stats(context.Background(), "c1", false)
	require.ErrorContains(t, err, "illisibles")
}

func TestEngine_CountRunning(t *testing.T) {
	n, err := newEngine(t, &fakeClient{running: 3}).CountRunning(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 3, n)
}
