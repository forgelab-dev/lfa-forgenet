package ws

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/forgelab-dev/lfa-forgenet/internal/metrics"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// Les messages envoyés sont ceux du contrat api/ (metrics.proto) : le
// dashboard et l'application mobile les décodent avec leur code généré.

func toSystemMetrics(h metrics.Host) *forgenetv1.SystemMetrics {
	return &forgenetv1.SystemMetrics{
		CollectedAt:           timestamppb.New(h.CollectedAt),
		UptimeSeconds:         h.UptimeSeconds,
		Cpu:                   toCPU(h.CPU),
		Memory:                toMemory(h.Memory),
		Storage:               toStorage(h.Storage),
		Network:               toNetwork(h.Network),
		DiskIo:                toDiskIO(h.DiskIO),
		ActiveContainersCount: h.ActiveContainers,
	}
}

func toContainerMetrics(c metrics.Container) *forgenetv1.ContainerMetrics {
	return &forgenetv1.ContainerMetrics{
		ContainerId: c.ContainerID,
		CollectedAt: timestamppb.New(c.CollectedAt),
		Cpu:         toCPU(c.CPU),
		Memory:      toMemory(c.Memory),
		Network:     toNetwork(c.Network),
		DiskIo:      toDiskIO(c.DiskIO),
	}
}

func toCPU(c metrics.CPU) *forgenetv1.CpuMetrics {
	return &forgenetv1.CpuMetrics{
		CpuUsagePercent:      c.UsagePercent,
		OnlineCpus:           c.OnlineCPUs,
		SystemCpuNanoseconds: c.SystemNanoseconds,
	}
}

func toMemory(m metrics.Memory) *forgenetv1.MemoryMetrics {
	return &forgenetv1.MemoryMetrics{
		UsageBytes:   m.UsageBytes,
		LimitBytes:   m.LimitBytes,
		CacheBytes:   m.CacheBytes,
		UsagePercent: m.UsagePercent,
	}
}

func toStorage(s metrics.Storage) *forgenetv1.StorageMetrics {
	return &forgenetv1.StorageMetrics{
		TotalBytes:   s.TotalBytes,
		UsedBytes:    s.UsedBytes,
		FreeBytes:    s.FreeBytes,
		UsagePercent: s.UsagePercent,
	}
}

func toNetwork(n metrics.Network) *forgenetv1.NetworkMetrics {
	return &forgenetv1.NetworkMetrics{
		RxBytesTotal:   n.RxBytes,
		TxBytesTotal:   n.TxBytes,
		RxPacketsTotal: n.RxPackets,
		TxPacketsTotal: n.TxPackets,
		RxErrorsTotal:  n.RxErrors,
		TxErrorsTotal:  n.TxErrors,
	}
}

func toDiskIO(d metrics.DiskIO) *forgenetv1.DiskIoMetrics {
	return &forgenetv1.DiskIoMetrics{
		ReadBytesTotal:  d.ReadBytes,
		WriteBytesTotal: d.WriteBytes,
		ReadOpsTotal:    d.ReadOps,
		WriteOpsTotal:   d.WriteOps,
	}
}
