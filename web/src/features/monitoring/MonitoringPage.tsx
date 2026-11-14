import { timestampMs } from '@bufbuild/protobuf/wkt';
import { Boxes } from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { Link } from 'wouter';

import { PageHeader } from '../../components/ui/PageHeader';
import { UserRole } from '../../gen/forgenet/v1/auth_pb';
import { ContainerState, type Container } from '../../gen/forgenet/v1/container_pb';
import type { LiveStatus } from '../../lib/live/useLive';
import { useNow } from '../../lib/live/useNow';
import { useApi, useCurrentUser } from '../auth/session';
import { listAllContainers } from '../containers/api';
import { ContainerTable } from './ContainerTable';
import { formatUptime } from './format';
import { HostGauges } from './HostGauges';
import { useContainersMetrics, useHostMetrics } from './live';
import { LinkBadge } from './LinkBadge';
import { RANGES } from './series';

/** Le moins bon des deux états de liaison : c'est lui que la pastille montre. */
const severity: Record<LiveStatus['state'], number> = {
  live: 0,
  paused: 1,
  connecting: 2,
  reconnecting: 3,
  gone: 4,
  error: 5,
};

export function MonitoringPage() {
  const user = useCurrentUser();
  const host = useHostMetrics();
  const fleet = useContainersMetrics();
  const known = useKnownContainers(fleet.containers.map((c) => c.id));
  const [rangeMs, setRangeMs] = useState<number>(RANGES[0].ms);
  const now = useNow();

  const status =
    severity[fleet.status.state] > severity[host.status.state] ? fleet.status : host.status;
  const hostStale = host.status.state !== 'live' && host.latest !== undefined;
  const fleetStale = fleet.status.state !== 'live' && fleet.sweepAt !== undefined;
  const lastAt = host.latest?.collectedAt ? timestampMs(host.latest.collectedAt) : undefined;
  const stopped = useMemo(
    () => [...known.values()].filter((c) => c.state !== ContainerState.RUNNING),
    [known]
  );
  const [showStopped, setShowStopped] = useState(false);

  return (
    <>
      <PageHeader
        title="Surveillance"
        description="Processeur, mémoire, disque et réseau du serveur et de chaque conteneur, en direct, avec la dernière heure."
        aside={<LinkBadge status={status} pace="toutes les 5 s" lastAt={lastAt} />}
      />

      {host.status.state === 'error' && (
        <p
          role="alert"
          className="mt-6 rounded-lg border border-danger-400/40 bg-danger-950 px-5 py-4 text-sm"
        >
          {host.status.message}
        </p>
      )}

      <section aria-labelledby="hote" className="animate-rise mt-10 [animation-delay:60ms]">
        <div className="flex flex-wrap items-center justify-between gap-x-4 gap-y-3">
          <h2 id="hote" className="text-lg font-semibold">
            Serveur
          </h2>
          <div className="flex flex-wrap items-center gap-x-4 gap-y-2">
            {host.latest && (
              <p className="font-mono text-xs text-iron-300">
                en marche depuis{' '}
                <b className="font-medium text-iron-100">
                  {formatUptime(Number(host.latest.uptimeSeconds))}
                </b>
                {' · '}
                <b className="font-medium text-iron-100">
                  {host.latest.activeContainersCount}
                </b>{' '}
                conteneurs actifs
                {' · '}
                <b className="font-medium text-iron-100">{host.latest.cpu?.onlineCpus ?? 0}</b>{' '}
                cœurs
              </p>
            )}
            <div
              role="group"
              aria-label="Fenêtre des courbes"
              className="inline-flex rounded-lg border border-iron-700 bg-iron-900 p-0.5"
            >
              {RANGES.map((r) => (
                <button
                  key={r.ms}
                  type="button"
                  aria-pressed={rangeMs === r.ms}
                  onClick={() => setRangeMs(r.ms)}
                  className={`min-h-9 rounded-md px-3 text-sm transition-colors ${
                    rangeMs === r.ms
                      ? 'bg-gold-950 text-gold-300 ring-1 ring-gold-400/50 ring-inset'
                      : 'text-iron-300 hover:text-iron-100'
                  }`}
                >
                  {r.label}
                </button>
              ))}
            </div>
          </div>
        </div>
        {host.latest ? (
          <HostGauges
            latest={host.latest}
            curves={host.curves}
            rangeMs={rangeMs}
            now={now}
            stale={hostStale}
          />
        ) : (
          <GaugesSkeleton />
        )}
      </section>

      <section aria-labelledby="conteneurs" className="animate-rise mt-10 [animation-delay:120ms]">
        <div className="flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1">
          <h2 id="conteneurs" className="text-lg font-semibold">
            Conteneurs
          </h2>
          {fleet.sweepAt !== undefined && (
            <p className="font-mono text-xs text-iron-300">
              {fleet.containers.length} {fleet.containers.length > 1 ? 'actifs' : 'actif'} · mesurés
              toutes les 10 s
            </p>
          )}
        </div>

        {fleet.sweepAt === undefined ? (
          <TableSkeleton />
        ) : fleet.containers.length === 0 ? (
          <div className="mt-4 flex max-w-2xl items-start gap-4 rounded-lg border border-iron-700 bg-iron-900 px-5 py-5">
            <span className="flex size-10 shrink-0 items-center justify-center rounded-md bg-iron-850 text-gold-400">
              <Boxes aria-hidden="true" className="size-5" />
            </span>
            <div>
              <p className="font-semibold">Aucun conteneur actif</p>
              <p className="mt-1 text-sm text-iron-300">
                Les mesures d’un conteneur apparaissent ici dès qu’il démarre.
                {user.role === UserRole.ADMIN && (
                  <>
                    {' '}
                    <Link
                      href="/deploiements"
                      className="font-medium text-gold-400 hover:text-gold-300"
                    >
                      Déployer une application
                    </Link>
                  </>
                )}
              </p>
            </div>
          </div>
        ) : (
          <div className="mt-4 overflow-hidden rounded-lg border border-iron-800">
            <ContainerTable
              containers={fleet.containers}
              known={known}
              hostMemory={host.latest?.memory ? Number(host.latest.memory.limitBytes) : undefined}
              rangeMs={rangeMs}
              now={now}
              stale={fleetStale}
            />
          </div>
        )}

        {stopped.length > 0 && fleet.sweepAt !== undefined && (
          <div className="mt-3 text-sm text-iron-300">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <p>
                {stopped.length}{' '}
                {stopped.length > 1
                  ? 'conteneurs arrêtés, sans mesure'
                  : 'conteneur arrêté, sans mesure'}
              </p>
              <button
                type="button"
                aria-expanded={showStopped}
                onClick={() => setShowStopped((v) => !v)}
                className="min-h-11 text-iron-100 underline underline-offset-4 hover:text-gold-300"
              >
                {showStopped ? 'Masquer' : 'Afficher'}
              </button>
            </div>
            {showStopped && (
              <ul className="mt-1 space-y-1 font-mono text-xs">
                {stopped.map((c) => (
                  <li key={c.id} className="truncate">
                    {c.name} <span className="text-iron-500">· {c.image}</span>
                  </li>
                ))}
              </ul>
            )}
          </div>
        )}
      </section>

      <p className="mt-10 text-xs text-iron-300">
        Le serveur garde une heure de mesures en mémoire : l’historique repart de zéro à chaque
        redémarrage.
      </p>
    </>
  );
}

/**
 * Noms et images des conteneurs, par identifiant. La liste est relue quand un
 * conteneur inconnu apparaît dans les mesures (tout juste déployé). Un échec
 * n'empêche rien : la ligne montre alors l'identifiant abrégé.
 */
function useKnownContainers(ids: string[]): Map<string, Container> {
  const api = useApi();
  const [known, setKnown] = useState<Map<string, Container>>(() => new Map());
  const loaded = useRef(false);
  // Change quand l'ensemble des conteneurs inconnus change : une relecture
  // par nouveauté, pas une par passage.
  const unknown = ids
    .filter((id) => !known.has(id))
    .sort()
    .join(',');

  useEffect(() => {
    if (loaded.current && unknown === '') return;
    const abort = new AbortController();
    listAllContainers(api, abort.signal)
      .then((list) => {
        loaded.current = true;
        setKnown(new Map(list.map((c) => [c.id, c])));
      })
      .catch(() => {});
    return () => abort.abort();
  }, [api, unknown]);

  return known;
}

function GaugesSkeleton() {
  return (
    <div
      role="status"
      aria-live="polite"
      className="mt-4 grid grid-cols-2 gap-2 sm:gap-3 xl:grid-cols-4"
    >
      <span className="sr-only">Chargement des mesures…</span>
      {[0, 1, 2, 3].map((i) => (
        <div key={i} className="space-y-3 rounded-lg border border-iron-800 bg-iron-900 p-4">
          <span className="block h-3 w-2/5 animate-pulse rounded bg-iron-800" />
          <span className="block h-6 w-1/2 animate-pulse rounded bg-iron-800" />
          <span className="block h-10 animate-pulse rounded bg-iron-800" />
        </div>
      ))}
    </div>
  );
}

function TableSkeleton() {
  return (
    <div
      role="status"
      aria-live="polite"
      className="mt-4 space-y-3 rounded-lg border border-iron-800 p-4"
    >
      <span className="sr-only">Chargement des conteneurs…</span>
      <span className="block h-3.5 w-1/3 animate-pulse rounded bg-iron-800" />
      <span className="block h-3.5 w-3/5 animate-pulse rounded bg-iron-800" />
      <span className="block h-3.5 w-2/5 animate-pulse rounded bg-iron-800" />
    </div>
  );
}
