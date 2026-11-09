package metrics

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadStats(t *testing.T, name string) container.StatsResponse {
	t.Helper()
	raw, err := os.ReadFile("testdata/docker/" + name)
	require.NoError(t, err)
	var s container.StatsResponse
	require.NoError(t, json.Unmarshal(raw, &s))
	return s
}

// Relevé réel d'un conteneur alpine qui fait tourner une boucle (un cœur
// plein sur huit) après avoir écrit 20 Mo, démon Docker 29 en cgroup v2.
func TestFromStatsCgroupV2(t *testing.T) {
	s := loadStats(t, "stats-cgroup2.json")

	got := FromStats(s, s.PreCPUStats)

	// (4 048 693 000 - 3 041 060 000) / (139 776 270 000 000 - 139 768 340 000 000) × 8 × 100
	assert.InDelta(t, 101.65, got.CPU.UsagePercent, 0.01)
	assert.Equal(t, uint32(8), got.CPU.OnlineCPUs)
	assert.Equal(t, uint64(60905000), got.CPU.SystemNanoseconds)

	// usage - inactive_file : les 20 Mo écrits restent en cache.
	assert.Equal(t, uint64(23691264-21172224), got.Memory.UsageBytes)
	assert.Equal(t, uint64(21172224), got.Memory.CacheBytes)
	assert.Equal(t, uint64(15842967552), got.Memory.LimitBytes)
	assert.InDelta(t, 0.0159, got.Memory.UsagePercent, 0.0001)

	assert.Equal(t, Network{RxBytes: 15551, TxBytes: 126, RxPackets: 67, TxPackets: 3}, got.Network)
	assert.Equal(t, DiskIO{ReadBytes: 3 * 1699840}, got.DiskIO, "disque et deux couches device-mapper, comme docker stats")
	assert.Equal(t, s.ID, got.ContainerID)
}

func TestFromStatsCgroupV1(t *testing.T) {
	s := loadStats(t, "stats-cgroup1.json")

	got := FromStats(s, s.PreCPUStats)

	assert.Equal(t, uint32(4), got.CPU.OnlineCPUs, "déduit de percpu_usage")
	assert.InDelta(t, 40.0, got.CPU.UsagePercent, 0.001) // 1 s sur 10 s × 4 cœurs
	assert.Equal(t, uint64(200000000), got.Memory.UsageBytes, "total_inactive_file l'emporte")
	assert.InDelta(t, 20.0, got.Memory.UsagePercent, 0.001)
	assert.Equal(t, DiskIO{ReadBytes: 4096, WriteBytes: 8192, ReadOps: 1, WriteOps: 2}, got.DiskIO)
	assert.Equal(t, Network{RxBytes: 1200, TxBytes: 600, RxPackets: 12, TxPackets: 6, RxErrors: 1, TxErrors: 3}, got.Network)
	assert.Equal(t, time.Date(2026, 11, 9, 10, 0, 0, 0, time.UTC), got.CollectedAt)
}

func TestContainerCPUSkipsResetCounters(t *testing.T) {
	prev := container.CPUStats{CPUUsage: container.CPUUsage{TotalUsage: 500}, SystemUsage: 1000, OnlineCPUs: 2}
	restarted := container.CPUStats{CPUUsage: container.CPUUsage{TotalUsage: 10}, SystemUsage: 2000, OnlineCPUs: 2}

	assert.Zero(t, containerCPU(prev, restarted).UsagePercent)
	assert.Zero(t, containerCPU(prev, prev).UsagePercent, "aucun temps système écoulé")
}

func TestContainerCollectorKeepsItsPreviousSample(t *testing.T) {
	src := &stubStats{responses: []container.StatsResponse{
		{
			CPUStats:    container.CPUStats{CPUUsage: container.CPUUsage{TotalUsage: 200}, SystemUsage: 2000, OnlineCPUs: 1},
			PreCPUStats: container.CPUStats{CPUUsage: container.CPUUsage{TotalUsage: 100}, SystemUsage: 1000},
		},
		// PreCPUStats vide, comme sans includePrevious : la référence est la mesure d'avant.
		{CPUStats: container.CPUStats{CPUUsage: container.CPUUsage{TotalUsage: 500}, SystemUsage: 3000, OnlineCPUs: 1}},
	}}
	collector := NewContainerCollector(src, "abc")

	first, err := collector.Collect(context.Background())
	require.NoError(t, err)
	second, err := collector.Collect(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []bool{true, false}, src.withPrevious, "relevé précédent demandé au démon la première fois seulement")
	assert.InDelta(t, 10.0, first.CPU.UsagePercent, 0.001)
	assert.InDelta(t, 30.0, second.CPU.UsagePercent, 0.001)
	assert.Equal(t, "abc", second.ContainerID, "identifiant demandé si le démon n'en renvoie pas")
}

func TestContainerCollectorReturnsSourceError(t *testing.T) {
	boom := errors.New("conteneur introuvable")
	_, err := NewContainerCollector(&stubStats{err: boom}, "abc").Collect(context.Background())
	require.ErrorIs(t, err, boom)
}

type stubStats struct {
	responses    []container.StatsResponse
	withPrevious []bool
	err          error
}

func (s *stubStats) Stats(_ context.Context, _ string, withPrevious bool) (container.StatsResponse, error) {
	if s.err != nil {
		return container.StatsResponse{}, s.err
	}
	s.withPrevious = append(s.withPrevious, withPrevious)
	r := s.responses[0]
	s.responses = s.responses[1:]
	return r, nil
}
