import type { ReactNode } from 'react';

import { formatBytes, formatPercent, formatRate } from './format';
import { useContainerDetail } from './live';
import { rates } from './series';
import { Sparkline } from './Sparkline';

interface Props {
  containerId: string;
  rangeMs: number;
  now: number;
}

/**
 * Détail d'un conteneur, ouvert depuis sa ligne : processeur, mémoire et
 * disque sur la fenêtre choisie, mesurés toutes les 2 s tant qu'il est ouvert.
 */
export function ContainerDetail({ containerId, rangeMs, now }: Props) {
  const { latest, curves, status } = useContainerDetail(containerId);
  const to = now;
  const from = to - rangeMs;
  const read = rates(curves.read);
  const write = rates(curves.write);
  const limit = Number(latest?.memory?.limitBytes ?? 0n);
  const used = Number(latest?.memory?.usageBytes ?? 0n);
  const memPercent = latest?.memory?.usagePercent;
  const axis = rangeMs > 15 * 60_000 ? ['−1 h', '−30 min'] : ['−15 min', '−7 min 30'];
  const stale = status.state !== 'live';

  return (
    <div className={stale ? '[&_svg]:opacity-60' : ''}>
      <div className="grid gap-3 md:grid-cols-3">
        <Chart
          title="Processeur"
          value={latest ? formatPercent(latest.cpu?.cpuUsagePercent) : '—'}
          axis={axis}
        >
          <Sparkline
            height={88}
            lines={[{ points: curves.cpu, color: 'var(--color-ok-400)', area: true }]}
            floor={10}
            from={from}
            to={to}
            label="Processeur du conteneur sur la fenêtre choisie"
          />
        </Chart>
        <Chart
          title="Mémoire"
          value={latest ? `${formatBytes(used)} / ${formatBytes(limit)}` : '—'}
          warn={(memPercent ?? 0) >= 85}
          axis={axis}
        >
          <Sparkline
            height={88}
            lines={[
              {
                points: curves.mem,
                color: (memPercent ?? 0) >= 85 ? 'var(--color-warn-400)' : 'var(--color-ok-400)',
                area: true,
              },
            ]}
            from={from}
            to={to}
            max={limit > 0 ? limit : undefined}
            label="Mémoire du conteneur, rapportée à sa limite"
          />
        </Chart>
        <Chart
          title="Disque"
          value={`lu ${formatRate(read.at(-1)?.v)} · écrit ${formatRate(write.at(-1)?.v)}`}
          axis={axis}
        >
          <Sparkline
            height={88}
            lines={[
              { points: read, color: 'var(--color-inflow-400)', area: true },
              { points: write, color: 'var(--color-outflow-400)' },
            ]}
            floor={64 * 1024}
            from={from}
            to={to}
            label="Lectures et écritures disque du conteneur"
          />
        </Chart>
      </div>
      <p className="mt-3 font-mono text-xs text-iron-300">
        {status.state === 'gone'
          ? 'Ce conteneur s’est arrêté : ses mesures restent visibles une heure.'
          : `${containerId.slice(0, 12)} · mesuré toutes les 2 s tant que la ligne est ouverte · historique gardé 1 h par le serveur`}
      </p>
    </div>
  );
}

interface ChartProps {
  title: string;
  value: string;
  axis: string[];
  warn?: boolean;
  children: ReactNode;
}

function Chart({ title, value, axis, warn = false, children }: ChartProps) {
  return (
    <figure className="m-0 min-w-0 rounded-md border border-iron-800 bg-iron-950 p-3">
      <figcaption className="flex items-baseline justify-between gap-3 text-sm text-iron-300">
        {title}
        <span
          className={`truncate font-mono text-[0.8125rem] ${warn ? 'text-warn-400' : 'text-iron-100'}`}
        >
          {value}
        </span>
      </figcaption>
      <div className="mt-2">{children}</div>
      <div
        aria-hidden="true"
        className="mt-1 flex justify-between font-mono text-[0.6875rem] text-iron-500"
      >
        <span>{axis[0]}</span>
        <span>{axis[1]}</span>
        <span>maintenant</span>
      </div>
    </figure>
  );
}
