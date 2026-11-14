package metricspb

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/forgelab-dev/lfa-forgenet/internal/metrics"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// ErrUnknownKey signale une métrique que la cible ne fournit pas : une
// métrique d'extension (aucune n'existe encore), ou une métrique de l'hôte
// demandée pour un conteneur.
var ErrUnknownKey = errors.New("métrique inconnue pour cette cible")

// field décrit une métrique native : son unité, sa façon de se résumer et sa
// lecture dans une mesure.
type field[T any] struct {
	unit  forgenetv1.MetricUnit
	kind  metrics.Kind
	value func(T) float64
}

func common[T any](cpu func(T) metrics.CPU, mem func(T) metrics.Memory, nw func(T) metrics.Network, io func(T) metrics.DiskIO) map[forgenetv1.MetricName]field[T] {
	g, c := metrics.Gauge, metrics.Counter
	return map[forgenetv1.MetricName]field[T]{
		forgenetv1.MetricName_METRIC_NAME_CPU_USAGE_PERCENT:      {forgenetv1.MetricUnit_METRIC_UNIT_PERCENT, g, func(x T) float64 { return cpu(x).UsagePercent }},
		forgenetv1.MetricName_METRIC_NAME_CPU_SYSTEM_NANOSECONDS: {forgenetv1.MetricUnit_METRIC_UNIT_NANOSECONDS, c, func(x T) float64 { return float64(cpu(x).SystemNanoseconds) }},
		forgenetv1.MetricName_METRIC_NAME_MEMORY_USAGE_BYTES:     {forgenetv1.MetricUnit_METRIC_UNIT_BYTES, g, func(x T) float64 { return float64(mem(x).UsageBytes) }},
		forgenetv1.MetricName_METRIC_NAME_MEMORY_LIMIT_BYTES:     {forgenetv1.MetricUnit_METRIC_UNIT_BYTES, g, func(x T) float64 { return float64(mem(x).LimitBytes) }},
		forgenetv1.MetricName_METRIC_NAME_MEMORY_CACHE_BYTES:     {forgenetv1.MetricUnit_METRIC_UNIT_BYTES, g, func(x T) float64 { return float64(mem(x).CacheBytes) }},
		forgenetv1.MetricName_METRIC_NAME_MEMORY_USAGE_PERCENT:   {forgenetv1.MetricUnit_METRIC_UNIT_PERCENT, g, func(x T) float64 { return mem(x).UsagePercent }},
		forgenetv1.MetricName_METRIC_NAME_NETWORK_RX_BYTES_TOTAL: {forgenetv1.MetricUnit_METRIC_UNIT_BYTES, c, func(x T) float64 { return float64(nw(x).RxBytes) }},
		forgenetv1.MetricName_METRIC_NAME_NETWORK_TX_BYTES_TOTAL: {forgenetv1.MetricUnit_METRIC_UNIT_BYTES, c, func(x T) float64 { return float64(nw(x).TxBytes) }},
		forgenetv1.MetricName_METRIC_NAME_NETWORK_RX_PACKETS_TOTAL: {
			forgenetv1.MetricUnit_METRIC_UNIT_PACKETS, c, func(x T) float64 { return float64(nw(x).RxPackets) },
		},
		forgenetv1.MetricName_METRIC_NAME_NETWORK_TX_PACKETS_TOTAL: {
			forgenetv1.MetricUnit_METRIC_UNIT_PACKETS, c, func(x T) float64 { return float64(nw(x).TxPackets) },
		},
		forgenetv1.MetricName_METRIC_NAME_NETWORK_RX_ERRORS_TOTAL: {forgenetv1.MetricUnit_METRIC_UNIT_COUNT, c, func(x T) float64 { return float64(nw(x).RxErrors) }},
		forgenetv1.MetricName_METRIC_NAME_NETWORK_TX_ERRORS_TOTAL: {forgenetv1.MetricUnit_METRIC_UNIT_COUNT, c, func(x T) float64 { return float64(nw(x).TxErrors) }},
		forgenetv1.MetricName_METRIC_NAME_DISK_READ_BYTES_TOTAL:   {forgenetv1.MetricUnit_METRIC_UNIT_BYTES, c, func(x T) float64 { return float64(io(x).ReadBytes) }},
		forgenetv1.MetricName_METRIC_NAME_DISK_WRITE_BYTES_TOTAL:  {forgenetv1.MetricUnit_METRIC_UNIT_BYTES, c, func(x T) float64 { return float64(io(x).WriteBytes) }},
		forgenetv1.MetricName_METRIC_NAME_DISK_READ_OPS_TOTAL:     {forgenetv1.MetricUnit_METRIC_UNIT_OPERATIONS, c, func(x T) float64 { return float64(io(x).ReadOps) }},
		forgenetv1.MetricName_METRIC_NAME_DISK_WRITE_OPS_TOTAL:    {forgenetv1.MetricUnit_METRIC_UNIT_OPERATIONS, c, func(x T) float64 { return float64(io(x).WriteOps) }},
	}
}

var (
	hostFields = func() map[forgenetv1.MetricName]field[metrics.Host] {
		f := common(
			func(h metrics.Host) metrics.CPU { return h.CPU },
			func(h metrics.Host) metrics.Memory { return h.Memory },
			func(h metrics.Host) metrics.Network { return h.Network },
			func(h metrics.Host) metrics.DiskIO { return h.DiskIO },
		)
		g, c := metrics.Gauge, metrics.Counter
		f[forgenetv1.MetricName_METRIC_NAME_STORAGE_TOTAL_BYTES] = field[metrics.Host]{forgenetv1.MetricUnit_METRIC_UNIT_BYTES, g, func(h metrics.Host) float64 { return float64(h.Storage.TotalBytes) }}
		f[forgenetv1.MetricName_METRIC_NAME_STORAGE_USED_BYTES] = field[metrics.Host]{forgenetv1.MetricUnit_METRIC_UNIT_BYTES, g, func(h metrics.Host) float64 { return float64(h.Storage.UsedBytes) }}
		f[forgenetv1.MetricName_METRIC_NAME_STORAGE_FREE_BYTES] = field[metrics.Host]{forgenetv1.MetricUnit_METRIC_UNIT_BYTES, g, func(h metrics.Host) float64 { return float64(h.Storage.FreeBytes) }}
		f[forgenetv1.MetricName_METRIC_NAME_STORAGE_USAGE_PERCENT] = field[metrics.Host]{forgenetv1.MetricUnit_METRIC_UNIT_PERCENT, g, func(h metrics.Host) float64 { return h.Storage.UsagePercent }}
		f[forgenetv1.MetricName_METRIC_NAME_SYSTEM_UPTIME_SECONDS] = field[metrics.Host]{forgenetv1.MetricUnit_METRIC_UNIT_SECONDS, c, func(h metrics.Host) float64 { return float64(h.UptimeSeconds) }}
		f[forgenetv1.MetricName_METRIC_NAME_SYSTEM_ACTIVE_CONTAINERS_COUNT] = field[metrics.Host]{forgenetv1.MetricUnit_METRIC_UNIT_COUNT, g, func(h metrics.Host) float64 { return float64(h.ActiveContainers) }}
		return f
	}()

	containerFields = common(
		func(c metrics.Container) metrics.CPU { return c.CPU },
		func(c metrics.Container) metrics.Memory { return c.Memory },
		func(c metrics.Container) metrics.Network { return c.Network },
		func(c metrics.Container) metrics.DiskIO { return c.DiskIO },
	)
)

// HostSeries regroupe des mesures de l'hôte en séries, une par métrique
// demandée (toutes si keys est vide), dans l'ordre du contrat.
func HostSeries(samples []metrics.Host, keys []*forgenetv1.MetricKey, step time.Duration) ([]*forgenetv1.MetricSeries, error) {
	return series(samples, func(h metrics.Host) time.Time { return h.CollectedAt }, hostFields, keys, step)
}

// ContainerSeries fait de même pour un conteneur.
func ContainerSeries(samples []metrics.Container, keys []*forgenetv1.MetricKey, step time.Duration) ([]*forgenetv1.MetricSeries, error) {
	return series(samples, func(c metrics.Container) time.Time { return c.CollectedAt }, containerFields, keys, step)
}

// CheckKeys vérifie les métriques demandées avant toute lecture.
func CheckKeys(keys []*forgenetv1.MetricKey, host bool) error {
	_, err := names(keys, func(n forgenetv1.MetricName) bool {
		if host {
			_, ok := hostFields[n]
			return ok
		}
		_, ok := containerFields[n]
		return ok
	})
	return err
}

func names(keys []*forgenetv1.MetricKey, known func(forgenetv1.MetricName) bool) ([]forgenetv1.MetricName, error) {
	var out []forgenetv1.MetricName
	if len(keys) == 0 {
		// Toutes, dans l'ordre des numéros du contrat.
		for n := range forgenetv1.MetricName_name {
			name := forgenetv1.MetricName(n)
			if known(name) {
				out = append(out, name)
			}
		}
		slices.Sort(out)
		return out, nil
	}
	seen := make(map[forgenetv1.MetricName]bool, len(keys))
	for _, k := range keys {
		if custom := k.GetCustom(); custom != "" {
			return nil, fmt.Errorf("%w : %q (aucune extension installée)", ErrUnknownKey, custom)
		}
		n := k.GetKnown()
		if !known(n) {
			return nil, fmt.Errorf("%w : %s", ErrUnknownKey, n)
		}
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	return out, nil
}

func series[T any](samples []T, at func(T) time.Time, fields map[forgenetv1.MetricName]field[T], keys []*forgenetv1.MetricKey, step time.Duration) ([]*forgenetv1.MetricSeries, error) {
	wanted, err := names(keys, func(n forgenetv1.MetricName) bool { _, ok := fields[n]; return ok })
	if err != nil {
		return nil, err
	}
	out := make([]*forgenetv1.MetricSeries, 0, len(wanted))
	points := make([]metrics.Point, len(samples))
	for _, name := range wanted {
		f := fields[name]
		for i, s := range samples {
			points[i] = metrics.Point{At: at(s), Value: f.value(s)}
		}
		reduced := metrics.Downsample(points, step, f.kind)
		ms := make([]*forgenetv1.MetricSample, len(reduced))
		for i, p := range reduced {
			ms[i] = &forgenetv1.MetricSample{Timestamp: timestamppb.New(p.At), Value: p.Value}
		}
		out = append(out, &forgenetv1.MetricSeries{
			Key:     &forgenetv1.MetricKey{Key: &forgenetv1.MetricKey_Known{Known: name}},
			Unit:    f.unit,
			Samples: ms,
		})
	}
	return out, nil
}
