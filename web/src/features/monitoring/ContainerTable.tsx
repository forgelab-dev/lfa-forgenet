import { ChevronRight } from 'lucide-react';
import { Fragment, useMemo, useState, type ReactNode } from 'react';

import type { Container } from '../../gen/forgenet/v1/container_pb';
import { ContainerDetail } from './ContainerDetail';
import { formatBytes, formatPercent, formatRate } from './format';
import type { ContainerNow } from './live';
import { toneOf } from './series';
import { Sparkline } from './Sparkline';

type SortKey = 'cpu' | 'mem' | 'name';

interface Props {
  containers: ContainerNow[];
  /** Conteneurs connus du serveur, par identifiant : nom et image. */
  known: Map<string, Container>;
  /** Mémoire de l'hôte : une limite égale signifie « sans limite ». */
  hostMemory?: number;
  rangeMs: number;
  /** Heure de fin des courbes. */
  now: number;
  stale: boolean;
}

/**
 * Tableau des conteneurs en marche, trié par processeur. Une ligne s'ouvre
 * sur son détail ; une seule à la fois, pour ne suivre qu'un flux de plus.
 */
export function ContainerTable({ containers, known, hostMemory, rangeMs, now, stale }: Props) {
  const [sort, setSort] = useState<SortKey>('cpu');
  const [open, setOpen] = useState<string>();

  const rows = useMemo(() => {
    const name = (c: ContainerNow) => known.get(c.id)?.name ?? c.id.slice(0, 12);
    return [...containers].sort((a, b) => {
      if (sort === 'name') return name(a).localeCompare(name(b), 'fr');
      if (sort === 'mem') return b.memUsed - a.memUsed;
      return b.cpu - a.cpu;
    });
  }, [containers, known, sort]);

  const header = (key: SortKey, label: string, align = 'text-right') => (
    <th
      scope="col"
      aria-sort={sort === key ? (key === 'name' ? 'ascending' : 'descending') : 'none'}
      className={`px-4 py-2.5 text-xs font-medium whitespace-nowrap text-iron-300 ${align}`}
    >
      <button
        type="button"
        onClick={() => setSort(key)}
        className={`inline-flex min-h-8 items-center gap-1 hover:text-iron-100 ${sort === key ? 'text-iron-100' : ''}`}
      >
        {label}
        {sort === key && <span aria-hidden="true">{key === 'name' ? '↑' : '↓'}</span>}
      </button>
    </th>
  );

  return (
    <table className={`w-full border-collapse text-sm ${stale ? '[&_.measure]:opacity-45' : ''}`}>
      <thead className="bg-iron-900 max-sm:hidden">
        <tr className="border-b border-iron-800">
          {header('name', 'Conteneur', 'text-left')}
          {header('cpu', 'Processeur')}
          {header('mem', 'Mémoire')}
          <th
            scope="col"
            className="px-4 py-2.5 text-right text-xs font-medium whitespace-nowrap text-iron-300"
          >
            Réseau ↓ / ↑
          </th>
          <th scope="col" className="px-4 py-2.5 text-right text-xs font-medium text-iron-300">
            15 min
          </th>
        </tr>
      </thead>
      <tbody>
        {rows.map((c) => {
          const info = known.get(c.id);
          const expanded = open === c.id;
          const unlimited =
            c.memLimit === 0 || (hostMemory !== undefined && c.memLimit >= hostMemory * 0.98);
          const memTone = unlimited ? 'ok' : toneOf(c.memPercent);
          const detailId = `detail-${c.id.slice(0, 12)}`;
          return (
            <Fragment key={c.id}>
              <tr
                onClick={() => setOpen(expanded ? undefined : c.id)}
                className={`cursor-pointer border-t border-iron-800 first:border-t-0 hover:bg-iron-900 max-sm:grid max-sm:grid-cols-2 max-sm:gap-x-3 max-sm:gap-y-1.5 max-sm:px-4 max-sm:py-3 ${expanded ? 'bg-iron-900' : ''}`}
              >
                <td className="min-w-0 px-4 py-3 max-sm:col-span-2 max-sm:p-0">
                  <button
                    type="button"
                    aria-expanded={expanded}
                    aria-controls={detailId}
                    onClick={(e) => {
                      e.stopPropagation();
                      setOpen(expanded ? undefined : c.id);
                    }}
                    className="flex min-h-11 w-full min-w-0 items-center gap-2.5 text-left"
                  >
                    <ChevronRight
                      aria-hidden="true"
                      className={`size-4 shrink-0 transition-transform duration-200 ease-out-quart ${expanded ? 'rotate-90 text-gold-400' : 'text-iron-500'}`}
                    />
                    <span className="min-w-0">
                      <span className="block truncate font-mono text-[0.8125rem] font-medium">
                        {info?.name ?? c.id.slice(0, 12)}
                      </span>
                      <span className="block truncate text-xs text-iron-300">
                        {info?.image ?? 'conteneur'}
                      </span>
                    </span>
                  </button>
                </td>
                <Cell label="Processeur">
                  {formatPercent(c.cpu)}
                  <span
                    aria-hidden="true"
                    className="ml-2.5 inline-block h-1 w-14 overflow-hidden rounded-sm bg-iron-800 align-middle max-sm:hidden"
                  >
                    <i
                      className="block h-full bg-ok-400"
                      style={{ width: `${Math.min(100, c.cpu)}%` }}
                    />
                  </span>
                </Cell>
                <Cell
                  label="Mémoire"
                  className={
                    memTone === 'warn'
                      ? 'text-warn-400'
                      : memTone === 'danger'
                        ? 'text-danger-400'
                        : ''
                  }
                >
                  {formatBytes(c.memUsed)}{' '}
                  <span className="text-iron-300">
                    / {unlimited ? 'sans limite' : formatBytes(c.memLimit)}
                  </span>
                </Cell>
                <Cell label="Réseau ↓ / ↑" className="max-sm:col-span-2">
                  <span className="text-inflow-400">{formatRate(c.rxRate)}</span>{' '}
                  <span className="text-iron-300">/</span>{' '}
                  <span className="text-outflow-400">{formatRate(c.txRate)}</span>
                </Cell>
                <td className="px-4 py-3 max-sm:hidden">
                  <Sparkline
                    className="measure ml-auto w-24!"
                    height={24}
                    lines={[{ points: c.trend, color: 'var(--color-ok-400)', area: true }]}
                    floor={10}
                    from={now - 15 * 60_000}
                    to={now}
                    label={`Processeur de ${info?.name ?? 'ce conteneur'} sur 15 minutes`}
                  />
                </td>
              </tr>
              {expanded && (
                <tr id={detailId} className="bg-iron-900 max-sm:block">
                  <td colSpan={5} className="px-4 pt-1 pb-5 sm:pl-10 max-sm:block">
                    <ContainerDetail containerId={c.id} rangeMs={rangeMs} now={now} />
                  </td>
                </tr>
              )}
            </Fragment>
          );
        })}
      </tbody>
    </table>
  );
}

function Cell({
  label,
  className = '',
  children,
}: {
  label: string;
  className?: string;
  children: ReactNode;
}) {
  return (
    <td
      data-label={label}
      className={`measure px-4 py-3 text-right font-mono whitespace-nowrap tabular-nums max-sm:p-0 max-sm:text-left max-sm:text-[0.8125rem] max-sm:before:block max-sm:before:font-sans max-sm:before:text-[0.6875rem] max-sm:before:text-iron-300 max-sm:before:content-[attr(data-label)] ${className}`}
    >
      {children}
    </td>
  );
}
