import { Cpu, HardDrive, MemoryStick, Network, type LucideIcon } from 'lucide-react';
import type { ReactNode } from 'react';

import type { SystemMetrics } from '../../gen/forgenet/v1/metrics_pb';
import { formatAgo, formatBytes, formatPercent, formatRate } from './format';
import type { HostCurves } from './live';
import { peak, rates, recent, toneOf, type Point, type Tone } from './series';
import { Sparkline } from './Sparkline';

interface Props {
  latest: SystemMetrics;
  curves: HostCurves;
  rangeMs: number;
  /** Heure de fin des courbes. */
  now: number;
  /** Mesures datées : liaison coupée, valeurs grisées. */
  stale: boolean;
}

const color = {
  ok: 'var(--color-ok-400)',
  warn: 'var(--color-warn-400)',
  danger: 'var(--color-danger-400)',
  inflow: 'var(--color-inflow-400)',
  outflow: 'var(--color-outflow-400)',
};

/** Les quatre instruments de l'hôte : processeur, mémoire, disque, réseau. */
export function HostGauges({ latest, curves, rangeMs, now, stale }: Props) {
  const to = now;
  const from = to - rangeMs;
  const cpu = latest.cpu?.cpuUsagePercent ?? 0;
  const cpuPeak = peak(recent(curves.cpu, rangeMs, to));
  const mem = latest.memory;
  const memTone = toneOf(mem?.usagePercent);
  const disk = latest.storage;
  const diskTone = toneOf(disk?.usagePercent);
  const rx = rates(curves.rx);
  const tx = rates(curves.tx);
  const read = rates(curves.read);
  const write = rates(curves.write);
  const errors = Number(
    (latest.network?.rxErrorsTotal ?? 0n) + (latest.network?.txErrorsTotal ?? 0n)
  );
  const windowLabel = rangeMs > 15 * 60_000 ? 'la dernière heure' : '15 minutes';

  return (
    <div
      className={`mt-4 grid grid-cols-2 gap-2 sm:gap-3 xl:grid-cols-4 ${stale ? '[&_.measure]:opacity-45' : ''}`}
    >
      <Gauge icon={Cpu} title="Processeur" tone="ok">
        <Value tone="ok" value={formatPercent(cpu)} />
        <Sub>
          {cpuPeak && cpuPeak.v >= cpu + 10
            ? `pointe à ${formatPercent(cpuPeak.v)} ${formatAgo(to - cpuPeak.t)}`
            : `${latest.cpu?.onlineCpus ?? 0} cœurs`}
        </Sub>
        <Sparkline
          className="measure mt-2.5"
          lines={[{ points: curves.cpu, color: color.ok, area: true }]}
          from={from}
          to={to}
          max={100}
          label={`Processeur sur ${windowLabel}${cpuPeak ? `, jusqu’à ${formatPercent(cpuPeak.v)}` : ''}`}
        />
      </Gauge>

      <Gauge icon={MemoryStick} title="Mémoire" tone={memTone}>
        <Value tone={memTone} value={formatPercent(mem?.usagePercent)} />
        <Sub>
          {formatBytes(Number(mem?.usageBytes ?? 0n))} sur{' '}
          {formatBytes(Number(mem?.limitBytes ?? 0n))}
          <br />+ {formatBytes(Number(mem?.cacheBytes ?? 0n))} de cache récupérable
        </Sub>
        <Sparkline
          className="measure mt-2.5"
          lines={[{ points: curves.mem, color: color[memTone], area: true }]}
          from={from}
          to={to}
          max={100}
          label={`Mémoire sur ${windowLabel}`}
        />
      </Gauge>

      <Gauge icon={HardDrive} title="Disque" tone={diskTone}>
        <Value tone={diskTone} value={formatPercent(disk?.usagePercent)} />
        {diskTone !== 'ok' && (
          <p
            className={`mt-1 text-xs ${diskTone === 'warn' ? 'text-warn-400' : 'text-danger-400'}`}
          >
            au-delà de {diskTone === 'warn' ? '85' : '95'} %
          </p>
        )}
        <Sub>
          {formatBytes(Number(disk?.freeBytes ?? 0n))} libres sur{' '}
          {formatBytes(Number(disk?.totalBytes ?? 0n))}
        </Sub>
        <div className="measure mt-3 h-1 overflow-hidden rounded-sm bg-iron-800" aria-hidden="true">
          <i
            className={`block h-full rounded-sm ${diskTone === 'ok' ? 'bg-ok-400' : diskTone === 'warn' ? 'bg-warn-400' : 'bg-danger-400'}`}
            style={{ width: `${Math.min(100, disk?.usagePercent ?? 0)}%` }}
          />
        </div>
        <Sub>
          lu {formatRate(last(read))} · écrit {formatRate(last(write))}
        </Sub>
      </Gauge>

      <Gauge icon={Network} title="Réseau" tone="ok">
        <div className="measure mt-2.5 flex flex-wrap gap-x-3.5 gap-y-1 font-mono tabular-nums">
          <Flow label="reçu" value={formatRate(last(rx))} className="text-inflow-400" />
          <Flow label="envoyé" value={formatRate(last(tx))} className="text-outflow-400" />
        </div>
        <Sub>
          {errors === 0 ? 'aucune erreur' : `${errors.toLocaleString('fr-FR')} erreurs`} depuis le
          démarrage
        </Sub>
        <Sparkline
          className="measure mt-2.5"
          lines={[
            { points: rx, color: color.inflow, area: true },
            { points: tx, color: color.outflow },
          ]}
          floor={64 * 1024}
          from={from}
          to={to}
          label={`Débit réseau sur ${windowLabel}, reçu et envoyé`}
        />
      </Gauge>
    </div>
  );
}

function last(points: Point[]): number | undefined {
  return points.length ? points[points.length - 1]!.v : undefined;
}

const border: Record<Tone, string> = {
  ok: 'border-iron-800',
  warn: 'border-warn-400/35',
  danger: 'border-danger-400/50',
};

function Gauge({
  icon: Icon,
  title,
  tone,
  children,
}: {
  icon: LucideIcon;
  title: string;
  tone: Tone;
  children: ReactNode;
}) {
  return (
    <section
      aria-label={title}
      className={`min-w-0 rounded-lg border bg-iron-900 px-3 pt-3 pb-3 sm:px-4 sm:pt-4 ${border[tone]}`}
    >
      <h3 className="flex items-center gap-2 text-sm font-normal text-iron-300">
        <Icon aria-hidden="true" className="size-4 text-iron-500" />
        {title}
      </h3>
      {children}
    </section>
  );
}

const valueTone: Record<Tone, string> = {
  ok: 'text-iron-100',
  warn: 'text-warn-400',
  danger: 'text-danger-400',
};

function Value({ value, tone }: { value: string; tone: Tone }) {
  const [number, unit] = value.split(' ');
  return (
    <p
      className={`measure mt-2.5 font-mono text-[1.375rem] leading-none font-medium tracking-[-0.02em] tabular-nums sm:text-[1.75rem] ${valueTone[tone]}`}
    >
      {number}
      {unit && <small className="ml-0.5 text-[0.9375rem] text-iron-300">{unit}</small>}
    </p>
  );
}

function Sub({ children }: { children: ReactNode }) {
  return (
    <p className="measure mt-1.5 font-mono text-xs leading-normal text-iron-300">{children}</p>
  );
}

function Flow({ label, value, className }: { label: string; value: string; className: string }) {
  return (
    <p className={`text-sm leading-none whitespace-nowrap sm:text-xl ${className}`}>
      <span className="mb-1 block font-sans text-[0.6875rem] text-iron-300">{label}</span>
      {value}
    </p>
  );
}
