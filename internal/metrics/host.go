package metrics

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/common"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

// HostProbe lit les compteurs bruts de la machine. GopsutilProbe est
// l'implémentation réelle ; les tests en substituent une autre.
type HostProbe interface {
	CPUTimes(ctx context.Context) (cpu.TimesStat, error)
	CPUCount(ctx context.Context) (int, error)
	Memory(ctx context.Context) (*mem.VirtualMemoryStat, error)
	DiskUsage(ctx context.Context, path string) (*disk.UsageStat, error)
	DiskIO(ctx context.Context) (map[string]disk.IOCountersStat, error)
	NetIO(ctx context.Context) ([]net.IOCountersStat, error)
	Uptime(ctx context.Context) (uint64, error)
}

// GopsutilProbe lit les compteurs avec gopsutil. ProcRoot remplace /proc
// (vide : variable HOST_PROC, sinon /proc) : un serveur lancé dans un
// conteneur y monte le /proc de l'hôte, sans quoi il ne verrait que le sien.
type GopsutilProbe struct {
	ProcRoot string
}

func (p GopsutilProbe) ctx(ctx context.Context) context.Context {
	if p.ProcRoot == "" {
		return ctx
	}
	return context.WithValue(ctx, common.EnvKey, common.EnvMap{common.HostProcEnvKey: p.ProcRoot})
}

// CPUTimes renvoie les temps cumulés de tous les processeurs réunis.
func (p GopsutilProbe) CPUTimes(ctx context.Context) (cpu.TimesStat, error) {
	times, err := cpu.TimesWithContext(p.ctx(ctx), false)
	if err != nil {
		return cpu.TimesStat{}, err
	}
	if len(times) == 0 {
		return cpu.TimesStat{}, fmt.Errorf("aucun temps processeur lu")
	}
	return times[0], nil
}

// CPUCount renvoie le nombre de processeurs logiques.
func (p GopsutilProbe) CPUCount(ctx context.Context) (int, error) {
	return cpu.CountsWithContext(p.ctx(ctx), true)
}

// Memory lit /proc/meminfo.
func (p GopsutilProbe) Memory(ctx context.Context) (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemoryWithContext(p.ctx(ctx))
}

// DiskUsage lit l'occupation du système de fichiers qui porte path.
func (p GopsutilProbe) DiskUsage(ctx context.Context, path string) (*disk.UsageStat, error) {
	return disk.UsageWithContext(p.ctx(ctx), path)
}

// DiskIO lit /proc/diskstats, partitions comprises.
func (p GopsutilProbe) DiskIO(ctx context.Context) (map[string]disk.IOCountersStat, error) {
	return disk.IOCountersWithContext(p.ctx(ctx))
}

// NetIO lit /proc/net/dev, interface par interface.
func (p GopsutilProbe) NetIO(ctx context.Context) ([]net.IOCountersStat, error) {
	return net.IOCountersWithContext(p.ctx(ctx), true)
}

// Uptime renvoie le temps écoulé depuis le démarrage, en secondes.
func (p GopsutilProbe) Uptime(ctx context.Context) (uint64, error) {
	return host.UptimeWithContext(p.ctx(ctx))
}

// RunningCounter compte les conteneurs en cours d'exécution.
type RunningCounter interface {
	CountRunning(ctx context.Context) (int, error)
}

// HostCollector mesure l'hôte. Il garde le relevé processeur précédent pour
// calculer l'utilisation entre deux mesures : chaque flux a le sien, pour que
// deux abonnés à des rythmes différents ne faussent pas leurs mesures.
type HostCollector struct {
	probe       HostProbe
	containers  RunningCounter
	storagePath string
	now         func() time.Time

	mu      sync.Mutex
	prevCPU *cpu.TimesStat
}

// NewHostCollector crée un collecteur. storagePath désigne le système de
// fichiers dont l'occupation est mesurée ; containers peut être nil.
func NewHostCollector(probe HostProbe, containers RunningCounter, storagePath string) *HostCollector {
	return &HostCollector{probe: probe, containers: containers, storagePath: storagePath, now: time.Now}
}

// Collect prend une mesure. La première mesure d'un collecteur donne
// l'utilisation du processeur depuis le démarrage de la machine.
func (c *HostCollector) Collect(ctx context.Context) (Host, error) {
	times, err := c.probe.CPUTimes(ctx)
	if err != nil {
		return Host{}, fmt.Errorf("temps processeur : %w", err)
	}
	cpus, err := c.probe.CPUCount(ctx)
	if err != nil {
		return Host{}, fmt.Errorf("nombre de processeurs : %w", err)
	}
	vm, err := c.probe.Memory(ctx)
	if err != nil {
		return Host{}, fmt.Errorf("mémoire : %w", err)
	}
	usage, err := c.probe.DiskUsage(ctx, c.storagePath)
	if err != nil {
		return Host{}, fmt.Errorf("occupation de %s : %w", c.storagePath, err)
	}
	disks, err := c.probe.DiskIO(ctx)
	if err != nil {
		return Host{}, fmt.Errorf("entrées-sorties disque : %w", err)
	}
	nics, err := c.probe.NetIO(ctx)
	if err != nil {
		return Host{}, fmt.Errorf("trafic réseau : %w", err)
	}
	uptime, err := c.probe.Uptime(ctx)
	if err != nil {
		return Host{}, fmt.Errorf("durée de fonctionnement : %w", err)
	}

	var running int
	if c.containers != nil {
		if running, err = c.containers.CountRunning(ctx); err != nil {
			return Host{}, fmt.Errorf("conteneurs actifs : %w", err)
		}
	}

	c.mu.Lock()
	prev := c.prevCPU
	c.prevCPU = &times
	c.mu.Unlock()

	return Host{
		CollectedAt:   c.now().UTC(),
		UptimeSeconds: uptime,
		CPU: CPU{
			UsagePercent:      hostCPUPercent(prev, times),
			OnlineCPUs:        uint32(max(cpus, 0)), //nolint:gosec // borné par max
			SystemNanoseconds: uint64(times.System * float64(time.Second)),
		},
		Memory:           hostMemory(vm),
		Storage:          hostStorage(usage),
		Network:          sumNetworks(nics),
		DiskIO:           sumDisks(disks),
		ActiveContainers: uint32(max(running, 0)), //nolint:gosec // borné par max
	}, nil
}

// hostCPUPercent applique le calcul de top : part du temps non inactif entre
// deux relevés. Les temps invités (guest) sont déjà comptés dans user et nice
// sous Linux : ils n'entrent pas dans le total.
func hostCPUPercent(prev *cpu.TimesStat, cur cpu.TimesStat) float64 {
	var before cpu.TimesStat
	if prev != nil {
		before = *prev
	}
	totalBefore, busyBefore := cpuBusy(before)
	totalNow, busyNow := cpuBusy(cur)
	total := totalNow - totalBefore
	if total <= 0 {
		return 0
	}
	busy := busyNow - busyBefore
	return clampPercent(busy / total * 100)
}

func cpuBusy(t cpu.TimesStat) (total, busy float64) {
	busy = t.User + t.System + t.Nice + t.Irq + t.Softirq + t.Steal
	total = busy + t.Idle + t.Iowait
	return total, busy
}

func clampPercent(v float64) float64 {
	return min(max(v, 0), 100)
}

// hostMemory suit free(1) : « utilisé » = total - disponible, le cache
// récupérable étant mis à part.
func hostMemory(vm *mem.VirtualMemoryStat) Memory {
	used := vm.Total - min(vm.Available, vm.Total)
	return Memory{
		UsageBytes:   used,
		LimitBytes:   vm.Total,
		CacheBytes:   vm.Cached + vm.Buffers,
		UsagePercent: percent(used, vm.Total),
	}
}

func hostStorage(u *disk.UsageStat) Storage {
	return Storage{
		TotalBytes:   u.Total,
		UsedBytes:    u.Used,
		FreeBytes:    u.Free,
		UsagePercent: u.UsedPercent,
	}
}

// Interfaces virtuelles écartées du trafic de l'hôte : le trafic d'un
// conteneur traverse son veth, le pont Docker puis la carte physique, et
// serait sinon compté deux ou trois fois.
var virtualNICPrefixes = []string{"lo", "veth", "docker", "br-", "virbr", "cni", "flannel", "cali", "vxlan", "tunl"}

func sumNetworks(nics []net.IOCountersStat) Network {
	var n Network
	for _, nic := range nics {
		if isVirtualNIC(nic.Name) {
			continue
		}
		n.RxBytes += nic.BytesRecv
		n.TxBytes += nic.BytesSent
		n.RxPackets += nic.PacketsRecv
		n.TxPackets += nic.PacketsSent
		n.RxErrors += nic.Errin
		n.TxErrors += nic.Errout
	}
	return n
}

func isVirtualNIC(name string) bool {
	for _, prefix := range virtualNICPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// Périphériques écartés des entrées-sorties de l'hôte : ceux qui ne sont pas
// des disques (loop, ram, zram, lecteurs optiques) et ceux qui reposent sur
// un autre disque (device-mapper, RAID logiciel), déjà comptés en dessous.
var virtualDiskPrefixes = []string{"loop", "ram", "zram", "sr", "fd", "dm-", "md"}

// sumDisks additionne les disques entiers. /proc/diskstats liste aussi leurs
// partitions (sda1, nvme0n1p1) : une entrée dont le nom prolonge celui d'un
// autre disque retenu est une partition, comptée dans son disque.
func sumDisks(disks map[string]disk.IOCountersStat) DiskIO {
	whole := make(map[string]bool, len(disks))
	for name := range disks {
		if !isVirtualDisk(name) {
			whole[name] = true
		}
	}
	var d DiskIO
	for name, stat := range disks {
		if !whole[name] || isPartition(name, whole) {
			continue
		}
		d.ReadBytes += stat.ReadBytes
		d.WriteBytes += stat.WriteBytes
		d.ReadOps += stat.ReadCount
		d.WriteOps += stat.WriteCount
	}
	return d
}

func isVirtualDisk(name string) bool {
	for _, prefix := range virtualDiskPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// isPartition reconnaît sda1 (disque + numéro) et nvme0n1p1 (disque + « p » +
// numéro), mais pas sdaa, le 27e disque, qui prolonge aussi « sda ».
func isPartition(name string, disks map[string]bool) bool {
	for other := range disks {
		suffix, ok := strings.CutPrefix(name, other)
		if !ok || suffix == "" {
			continue
		}
		suffix = strings.TrimPrefix(suffix, "p")
		if suffix != "" && strings.Trim(suffix, "0123456789") == "" {
			return true
		}
	}
	return false
}
