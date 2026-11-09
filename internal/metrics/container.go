package metrics

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/moby/moby/api/types/container"
)

// StatsSource lit les statistiques d'un conteneur auprès du démon Docker.
// withPrevious demande au démon un second relevé, une seconde plus tôt, dans
// PreCPUStats : utile pour la première mesure seulement.
type StatsSource interface {
	Stats(ctx context.Context, id string, withPrevious bool) (container.StatsResponse, error)
}

// ContainerCollector mesure un conteneur. Comme HostCollector, il garde le
// relevé processeur précédent : un collecteur par flux.
type ContainerCollector struct {
	source StatsSource
	id     string

	mu      sync.Mutex
	prevCPU *container.CPUStats
}

// NewContainerCollector crée le collecteur du conteneur id.
func NewContainerCollector(source StatsSource, id string) *ContainerCollector {
	return &ContainerCollector{source: source, id: id}
}

// Collect prend une mesure. La première demande au démon un relevé
// précédent, ce qui la retarde d'une seconde ; les suivantes repartent de la
// mesure d'avant.
func (c *ContainerCollector) Collect(ctx context.Context) (Container, error) {
	c.mu.Lock()
	first := c.prevCPU == nil
	c.mu.Unlock()

	stats, err := c.source.Stats(ctx, c.id, first)
	if err != nil {
		return Container{}, err
	}

	c.mu.Lock()
	prev := stats.PreCPUStats
	if c.prevCPU != nil {
		prev = *c.prevCPU
	}
	cur := stats.CPUStats
	c.prevCPU = &cur
	c.mu.Unlock()

	m := FromStats(stats, prev)
	if m.ContainerID == "" {
		m.ContainerID = c.id
	}
	return m, nil
}

// FromStats convertit une réponse du démon, avec prev comme relevé
// processeur de référence, selon les calculs de `docker stats`.
func FromStats(s container.StatsResponse, prev container.CPUStats) Container {
	collected := s.Read.UTC()
	if s.Read.IsZero() {
		collected = time.Now().UTC()
	}
	return Container{
		ContainerID: s.ID,
		CollectedAt: collected,
		CPU:         containerCPU(prev, s.CPUStats),
		Memory:      containerMemory(s.MemoryStats),
		Network:     containerNetwork(s.Networks),
		DiskIO:      containerDiskIO(s.BlkioStats),
	}
}

func onlineCPUs(s container.CPUStats) uint32 {
	if s.OnlineCPUs > 0 {
		return s.OnlineCPUs
	}
	return uint32(len(s.CPUUsage.PercpuUsage)) //nolint:gosec // nombre de cœurs
}

// containerCPU : part du temps de la machine consommée par le conteneur,
// multipliée par le nombre de cœurs (200 % = deux cœurs pleins).
func containerCPU(prev, cur container.CPUStats) CPU {
	cpus := onlineCPUs(cur)
	out := CPU{OnlineCPUs: cpus, SystemNanoseconds: cur.CPUUsage.UsageInKernelmode}
	// Compteurs remis à zéro (conteneur redémarré) : pas de mesure plutôt
	// qu'une valeur absurde.
	if cur.CPUUsage.TotalUsage < prev.CPUUsage.TotalUsage || cur.SystemUsage <= prev.SystemUsage {
		return out
	}
	cpuDelta := float64(cur.CPUUsage.TotalUsage - prev.CPUUsage.TotalUsage)
	systemDelta := float64(cur.SystemUsage - prev.SystemUsage)
	out.UsagePercent = cpuDelta / systemDelta * float64(cpus) * 100
	return out
}

// containerMemory retire le cache inactif de l'usage, comme `docker stats` :
// clé inactive_file en cgroup v2, total_inactive_file en cgroup v1.
func containerMemory(s container.MemoryStats) Memory {
	cache := s.Stats["inactive_file"]
	if v, ok := s.Stats["total_inactive_file"]; ok {
		cache = v
	}
	used := s.Usage
	if cache < used {
		used -= cache
	}
	return Memory{
		UsageBytes:   used,
		LimitBytes:   s.Limit,
		CacheBytes:   cache,
		UsagePercent: percent(used, s.Limit),
	}
}

func containerNetwork(nets map[string]container.NetworkStats) Network {
	var n Network
	for _, s := range nets {
		n.RxBytes += s.RxBytes
		n.TxBytes += s.TxBytes
		n.RxPackets += s.RxPackets
		n.TxPackets += s.TxPackets
		n.RxErrors += s.RxErrors
		n.TxErrors += s.TxErrors
	}
	return n
}

// containerDiskIO additionne les opérations « read » et « write » ; la casse
// de Op varie selon la version des cgroups (Read, read). Comme `docker stats`,
// toutes les entrées sont additionnées : sur un disque chiffré ou en LVM, une
// même lecture apparaît sur le disque et sur ses couches device-mapper.
// Les nombres d'opérations restent à zéro en cgroup v2, où le démon ne les
// fournit pas.
func containerDiskIO(s container.BlkioStats) DiskIO {
	var d DiskIO
	d.ReadBytes, d.WriteBytes = sumBlkio(s.IoServiceBytesRecursive)
	d.ReadOps, d.WriteOps = sumBlkio(s.IoServicedRecursive)
	return d
}

func sumBlkio(entries []container.BlkioStatEntry) (read, write uint64) {
	for _, e := range entries {
		switch strings.ToLower(e.Op) {
		case "read":
			read += e.Value
		case "write":
			write += e.Value
		}
	}
	return read, write
}
