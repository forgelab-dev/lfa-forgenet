import { fromJson } from '@bufbuild/protobuf';
import { timestampFromMs, timestampMs } from '@bufbuild/protobuf/wkt';
import { useCallback, useRef, useState } from 'react';

import {
  MetricName,
  StreamAllContainerMetricsResponseSchema,
  StreamContainerMetricsResponseSchema,
  StreamSystemMetricsResponseSchema,
  type ContainerMetrics,
  type SystemMetrics,
} from '../../gen/forgenet/v1/metrics_pb';
import { readSocket } from '../../lib/live/socket';
import { useLive, type LiveStatus } from '../../lib/live/useLive';
import { useApi } from '../auth/session';
import { HISTORY_MS, byName, key, merge, rate, recent, type Point } from './series';

const json = { ignoreUnknownFields: true };

/** Les courbes de l'hôte : pourcentages, et compteurs cumulés d'où l'écran tire les débits. */
export interface HostCurves {
  cpu: Point[];
  mem: Point[];
  rx: Point[];
  tx: Point[];
  read: Point[];
  write: Point[];
}

const hostKeys: [keyof HostCurves, MetricName][] = [
  ['cpu', MetricName.CPU_USAGE_PERCENT],
  ['mem', MetricName.MEMORY_USAGE_PERCENT],
  ['rx', MetricName.NETWORK_RX_BYTES_TOTAL],
  ['tx', MetricName.NETWORK_TX_BYTES_TOTAL],
  ['read', MetricName.DISK_READ_BYTES_TOTAL],
  ['write', MetricName.DISK_WRITE_BYTES_TOTAL],
];

const emptyCurves = (): HostCurves => ({ cpu: [], mem: [], rx: [], tx: [], read: [], write: [] });

function hostPoints(m: SystemMetrics, t: number): Record<keyof HostCurves, Point> {
  return {
    cpu: { t, v: m.cpu?.cpuUsagePercent ?? 0 },
    mem: { t, v: m.memory?.usagePercent ?? 0 },
    rx: { t, v: Number(m.network?.rxBytesTotal ?? 0n) },
    tx: { t, v: Number(m.network?.txBytesTotal ?? 0n) },
    read: { t, v: Number(m.diskIo?.readBytesTotal ?? 0n) },
    write: { t, v: Number(m.diskIo?.writeBytesTotal ?? 0n) },
  };
}

/**
 * Suit l'hôte : l'heure gardée par le serveur, puis chaque mesure en direct
 * (toutes les 5 s). Après une coupure, l'historique comble le trou avant la
 * reprise du direct.
 */
export function useHostMetrics() {
  const api = useApi();
  const [latest, setLatest] = useState<SystemMetrics>();
  const [curves, setCurves] = useState<HostCurves>(emptyCurves);
  const last = useRef<number>(undefined);

  const open = useCallback(
    async (signal: AbortSignal, token: string, onLive: () => void) => {
      const since = last.current ?? Date.now() - HISTORY_MS;
      const res = await api.metrics.querySystemMetrics(
        { range: { startTime: timestampFromMs(since) }, keys: hostKeys.map(([, n]) => key(n)) },
        { signal }
      );
      const series = byName(res.series);
      setCurves((c) => {
        const next = { ...c };
        for (const [field, name] of hostKeys) next[field] = merge(c[field], series.get(name) ?? []);
        return next;
      });
      await readSocket(
        '/ws/v1/metrics/system',
        token,
        (data) => {
          const m = fromJson(StreamSystemMetricsResponseSchema, data as never, json).metrics;
          if (!m?.collectedAt) return;
          const t = timestampMs(m.collectedAt);
          last.current = t;
          const points = hostPoints(m, t);
          setLatest(m);
          setCurves((c) => {
            const next = { ...c };
            for (const [field] of hostKeys) next[field] = merge(c[field], [points[field]]);
            return next;
          });
          onLive();
        },
        signal
      );
    },
    [api]
  );

  const status = useLive(open);
  return { latest, curves, status };
}

/** Un conteneur dans le dernier passage du relevé, avec ses débits. */
export interface ContainerNow {
  id: string;
  at: number;
  cpu: number;
  memUsed: number;
  memLimit: number;
  memPercent: number;
  rxRate?: number;
  txRate?: number;
  /** Processeur sur les 15 dernières minutes. */
  trend: Point[];
}

const TREND_MS = 15 * 60_000;
const TREND_POINTS = 90;
/** Historiques demandés en même temps, pour ne pas saturer une connexion lente. */
const TREND_PARALLEL = 4;

interface Totals {
  rx: Point;
  tx: Point;
}

/**
 * Suit tous les conteneurs sur une seule liaison (ADR-0010) : un message
 * par passage du relevé, toutes les 10 s. La tendance de chaque conteneur
 * part de son historique, chargé une fois puis prolongé en direct.
 */
export function useContainersMetrics() {
  const api = useApi();
  const [sweepAt, setSweepAt] = useState<number>();
  const [containers, setContainers] = useState<ContainerNow[]>([]);
  const totals = useRef(new Map<string, Totals>());
  const trends = useRef(new Map<string, Point[]>());
  const loaded = useRef(new Set<string>());

  const open = useCallback(
    async (signal: AbortSignal, token: string, onLive: () => void) => {
      // Une reprise recharge les tendances : le trou de la coupure se comble.
      loaded.current.clear();
      const queue: string[] = [];
      let running = 0;
      const pump = () => {
        while (running < TREND_PARALLEL && queue.length > 0 && !signal.aborted) {
          const id = queue.shift()!;
          running++;
          api.metrics
            .queryContainerMetrics(
              {
                containerId: id,
                range: { startTime: timestampFromMs(Date.now() - TREND_MS) },
                maxPoints: TREND_POINTS,
                keys: [key(MetricName.CPU_USAGE_PERCENT)],
              },
              { signal }
            )
            .then((res) => {
              const points = byName(res.series).get(MetricName.CPU_USAGE_PERCENT) ?? [];
              trends.current.set(id, merge(trends.current.get(id) ?? [], points));
              setContainers((list) =>
                list.map((c) => (c.id === id ? { ...c, trend: trends.current.get(id)! } : c))
              );
            })
            .catch(() => {
              loaded.current.delete(id); // nouvel essai au passage suivant
            })
            .finally(() => {
              running--;
              pump();
            });
        }
      };

      await readSocket(
        '/ws/v1/metrics/containers',
        token,
        (data) => {
          const res = fromJson(StreamAllContainerMetricsResponseSchema, data as never, json);
          const at = res.collectedAt ? timestampMs(res.collectedAt) : Date.now();
          const seen = new Set<string>();
          const list = res.containers.map((m) => {
            seen.add(m.containerId);
            const now = toContainerNow(m, at, totals.current.get(m.containerId));
            totals.current.set(m.containerId, {
              rx: { t: now.at, v: Number(m.network?.rxBytesTotal ?? 0n) },
              tx: { t: now.at, v: Number(m.network?.txBytesTotal ?? 0n) },
            });
            const trend = recent(
              merge(trends.current.get(m.containerId) ?? [], [{ t: now.at, v: now.cpu }]),
              TREND_MS
            );
            trends.current.set(m.containerId, trend);
            if (!loaded.current.has(m.containerId)) {
              loaded.current.add(m.containerId);
              queue.push(m.containerId);
            }
            return { ...now, trend };
          });
          for (const id of [...totals.current.keys()]) {
            if (!seen.has(id)) {
              totals.current.delete(id);
              trends.current.delete(id);
            }
          }
          setSweepAt(at);
          setContainers(list);
          onLive();
          pump();
        },
        signal
      );
    },
    [api]
  );

  const status = useLive(open);
  return { sweepAt, containers, status };
}

function toContainerNow(
  m: ContainerMetrics,
  sweepAt: number,
  prev: Totals | undefined
): ContainerNow {
  const at = m.collectedAt ? timestampMs(m.collectedAt) : sweepAt;
  return {
    id: m.containerId,
    at,
    cpu: m.cpu?.cpuUsagePercent ?? 0,
    memUsed: Number(m.memory?.usageBytes ?? 0n),
    memLimit: Number(m.memory?.limitBytes ?? 0n),
    memPercent: m.memory?.usagePercent ?? 0,
    rxRate: rate(prev?.rx, { t: at, v: Number(m.network?.rxBytesTotal ?? 0n) }),
    txRate: rate(prev?.tx, { t: at, v: Number(m.network?.txBytesTotal ?? 0n) }),
    trend: [],
  };
}

/** Les courbes du détail d'un conteneur. */
export interface ContainerCurves {
  cpu: Point[];
  mem: Point[];
  read: Point[];
  write: Point[];
}

const detailKeys: [keyof ContainerCurves, MetricName][] = [
  ['cpu', MetricName.CPU_USAGE_PERCENT],
  ['mem', MetricName.MEMORY_USAGE_BYTES],
  ['read', MetricName.DISK_READ_BYTES_TOTAL],
  ['write', MetricName.DISK_WRITE_BYTES_TOTAL],
];

/**
 * Détail d'un conteneur, tant que sa ligne est ouverte : l'heure gardée,
 * puis une mesure toutes les 2 s. Un conteneur arrêté termine le suivi
 * (état `gone`).
 */
export function useContainerDetail(containerId: string): {
  latest?: ContainerMetrics;
  curves: ContainerCurves;
  status: LiveStatus;
} {
  const api = useApi();
  const [latest, setLatest] = useState<ContainerMetrics>();
  const [curves, setCurves] = useState<ContainerCurves>({ cpu: [], mem: [], read: [], write: [] });
  const last = useRef<number>(undefined);

  const open = useCallback(
    async (signal: AbortSignal, token: string, onLive: () => void) => {
      const since = last.current ?? Date.now() - HISTORY_MS;
      const res = await api.metrics.queryContainerMetrics(
        {
          containerId,
          range: { startTime: timestampFromMs(since) },
          keys: detailKeys.map(([, n]) => key(n)),
        },
        { signal }
      );
      const series = byName(res.series);
      setCurves((c) => {
        const next = { ...c };
        for (const [field, name] of detailKeys)
          next[field] = merge(c[field], series.get(name) ?? []);
        return next;
      });
      await readSocket(
        `/ws/v1/metrics/containers/${encodeURIComponent(containerId)}`,
        token,
        (data) => {
          const m = fromJson(StreamContainerMetricsResponseSchema, data as never, json).metrics;
          if (!m?.collectedAt) return;
          const t = timestampMs(m.collectedAt);
          last.current = t;
          setLatest(m);
          setCurves((c) => ({
            cpu: merge(c.cpu, [{ t, v: m.cpu?.cpuUsagePercent ?? 0 }]),
            mem: merge(c.mem, [{ t, v: Number(m.memory?.usageBytes ?? 0n) }]),
            read: merge(c.read, [{ t, v: Number(m.diskIo?.readBytesTotal ?? 0n) }]),
            write: merge(c.write, [{ t, v: Number(m.diskIo?.writeBytesTotal ?? 0n) }]),
          }));
          onLive();
        },
        signal
      );
    },
    [api, containerId]
  );

  const status = useLive(open);
  return { latest, curves, status };
}
