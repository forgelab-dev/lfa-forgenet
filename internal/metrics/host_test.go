package metrics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const kB = 1024

type fixedCounter int

func (n fixedCounter) CountRunning(context.Context) (int, error) { return int(n), nil }

// gopsutil lit les fichiers de testdata/proc comme s'ils étaient /proc : le
// test couvre la lecture réelle de gopsutil et nos calculs par-dessus.
func TestHostCollectorParsesProcFixture(t *testing.T) {
	collector := NewHostCollector(GopsutilProbe{ProcRoot: "testdata/proc"}, fixedCounter(3), t.TempDir())

	got, err := collector.Collect(context.Background())
	require.NoError(t, err)

	// /proc/stat : 10 000 tops dont 8 000 inactifs et 500 en attente d'E/S,
	// 500 en mode noyau (5 s à 100 tops par seconde).
	assert.InDelta(t, 15.0, got.CPU.UsagePercent, 0.001)
	assert.Equal(t, uint32(2), got.CPU.OnlineCPUs)
	assert.Equal(t, uint64(5*time.Second), got.CPU.SystemNanoseconds)

	// /proc/meminfo : utilisé = total - disponible ; cache = Cached +
	// SReclaimable (gopsutil) + Buffers.
	assert.Equal(t, uint64(2_000_000*kB), got.Memory.UsageBytes)
	assert.Equal(t, uint64(8_000_000*kB), got.Memory.LimitBytes)
	assert.Equal(t, uint64(1_300_000*kB), got.Memory.CacheBytes)
	assert.InDelta(t, 25.0, got.Memory.UsagePercent, 0.001)

	// /proc/diskstats : nvme0n1 et sda seulement, sans leurs partitions,
	// loop0 ni dm-0 ; secteurs de 512 octets.
	assert.Equal(t, DiskIO{ReadBytes: 2400 * 512, WriteBytes: 1200 * 512, ReadOps: 110, WriteOps: 55}, got.DiskIO)

	// /proc/net/dev : eth0 et wlan0, sans lo, docker0 ni veth.
	assert.Equal(t, Network{RxBytes: 12000, TxBytes: 5000, RxPackets: 120, TxPackets: 50, RxErrors: 1, TxErrors: 2}, got.Network)

	assert.Equal(t, uint32(3), got.ActiveContainers)
	assert.NotZero(t, got.Storage.TotalBytes, "occupation lue sur le dossier temporaire")
	assert.NotZero(t, got.UptimeSeconds)
	assert.WithinDuration(t, time.Now(), got.CollectedAt, time.Minute)
}

func TestHostCollectorMeasuresCPUBetweenTwoSamples(t *testing.T) {
	probe := &stubProbe{times: []cpu.TimesStat{
		{User: 100, System: 50, Idle: 800, Iowait: 50},
		{User: 130, System: 60, Idle: 850, Iowait: 60}, // +100 dont 40 occupés
	}}
	collector := NewHostCollector(probe, nil, "/")

	first, err := collector.Collect(context.Background())
	require.NoError(t, err)
	assert.InDelta(t, 15.0, first.CPU.UsagePercent, 0.001, "depuis le démarrage")

	second, err := collector.Collect(context.Background())
	require.NoError(t, err)
	assert.InDelta(t, 40.0, second.CPU.UsagePercent, 0.001, "entre les deux relevés")
	assert.Zero(t, second.ActiveContainers, "sans compteur de conteneurs")
}

func TestHostCollectorReportsProbeFailure(t *testing.T) {
	probe := &stubProbe{times: []cpu.TimesStat{{}}, memErr: errors.New("meminfo illisible")}

	_, err := NewHostCollector(probe, nil, "/").Collect(context.Background())

	require.ErrorContains(t, err, "mémoire : meminfo illisible")
}

func TestHostCPUPercentIgnoresGuestAndCounterReset(t *testing.T) {
	prev := cpu.TimesStat{User: 100, Idle: 100, Guest: 50}
	// Le temps invité est déjà inclus dans user : il ne doit pas gonfler le total.
	assert.InDelta(t, 50.0, hostCPUPercent(&prev, cpu.TimesStat{User: 200, Idle: 200, Guest: 100}), 0.001)
	assert.Zero(t, hostCPUPercent(&prev, prev), "aucun temps écoulé")
	assert.Zero(t, hostCPUPercent(&prev, cpu.TimesStat{User: 1}), "compteurs remis à zéro")
}

func TestSumDisksKeepsWholeDisksOnly(t *testing.T) {
	disks := map[string]disk.IOCountersStat{
		"sda":       {ReadBytes: 1},
		"sda1":      {ReadBytes: 10},
		"sdaa":      {ReadBytes: 100}, // 27e disque, pas une partition de sda
		"nvme0n1":   {ReadBytes: 1000},
		"nvme0n1p2": {ReadBytes: 10000},
		"mmcblk0":   {ReadBytes: 100000},
		"mmcblk0p1": {ReadBytes: 1000000},
		"md0":       {ReadBytes: 1 << 30},
		"zram0":     {ReadBytes: 1 << 31},
	}
	assert.Equal(t, uint64(101101), sumDisks(disks).ReadBytes)
}

func TestHostMemoryToleratesAvailableAboveTotal(t *testing.T) {
	m := hostMemory(&mem.VirtualMemoryStat{Total: 100, Available: 150})
	assert.Zero(t, m.UsageBytes)
	assert.Zero(t, hostMemory(&mem.VirtualMemoryStat{}).UsagePercent, "total nul")
}

type stubProbe struct {
	times  []cpu.TimesStat
	calls  int
	memErr error
}

func (p *stubProbe) CPUTimes(context.Context) (cpu.TimesStat, error) {
	t := p.times[min(p.calls, len(p.times)-1)]
	p.calls++
	return t, nil
}

func (p *stubProbe) CPUCount(context.Context) (int, error) { return 4, nil }

func (p *stubProbe) Memory(context.Context) (*mem.VirtualMemoryStat, error) {
	if p.memErr != nil {
		return nil, p.memErr
	}
	return &mem.VirtualMemoryStat{Total: 100, Available: 50}, nil
}

func (p *stubProbe) DiskUsage(context.Context, string) (*disk.UsageStat, error) {
	return &disk.UsageStat{Total: 10, Used: 4, Free: 6, UsedPercent: 40}, nil
}

func (p *stubProbe) DiskIO(context.Context) (map[string]disk.IOCountersStat, error) {
	return nil, nil
}

func (p *stubProbe) NetIO(context.Context) ([]net.IOCountersStat, error) { return nil, nil }

func (p *stubProbe) Uptime(context.Context) (uint64, error) { return 42, nil }
