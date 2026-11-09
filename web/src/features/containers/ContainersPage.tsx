import { timestampDate } from '@bufbuild/protobuf/wkt';
import { useCallback } from 'react';

import { PageHeader } from '../../components/ui/PageHeader';
import { RequestError } from '../../components/ui/RequestError';
import { ContainerState, type Container } from '../../gen/forgenet/v1/container_pb';
import { useRequest } from '../../lib/useRequest';
import { useApi } from '../auth/session';
import { listAllContainers } from './api';
import { StateBadge } from './StateBadge';
import { formatPort } from './state';

const created = new Intl.DateTimeFormat('fr-FR', { dateStyle: 'medium', timeStyle: 'short' });

/** Conteneurs en marche d'abord, puis par nom. */
function byStateThenName(a: Container, b: Container): number {
  const running =
    Number(b.state === ContainerState.RUNNING) - Number(a.state === ContainerState.RUNNING);
  return running || a.name.localeCompare(b.name, 'fr');
}

export function ContainersPage() {
  const api = useApi();
  const load = useCallback((signal: AbortSignal) => listAllContainers(api, signal), [api]);
  const [request, retry] = useRequest(load);

  return (
    <>
      <PageHeader
        title="Conteneurs"
        description="Les conteneurs Docker du serveur, arrêtés compris. Les actions (démarrer, arrêter, lire les journaux) arrivent dans une prochaine version."
      />
      <div className="mt-8">
        {request.status === 'loading' && <ContainersSkeleton />}
        {request.status === 'error' && <RequestError error={request.error} onRetry={retry} />}
        {request.status === 'success' && <ContainersTable containers={request.data} />}
      </div>
    </>
  );
}

function ContainersTable({ containers }: { containers: Container[] }) {
  if (containers.length === 0) {
    return (
      <p className="rounded-lg border border-dashed border-iron-700 px-5 py-10 text-center text-iron-300">
        Aucun conteneur sur ce serveur pour l’instant. Ils apparaîtront ici dès le premier
        déploiement.
      </p>
    );
  }

  return (
    <div className="animate-rise overflow-hidden rounded-lg border border-iron-800">
      <table className="w-full text-left text-sm">
        <caption className="sr-only">Conteneurs du serveur</caption>
        <thead className="bg-iron-900 text-iron-300">
          <tr>
            <th scope="col" className="px-4 py-3 font-medium">
              Nom
            </th>
            <th scope="col" className="hidden px-4 py-3 font-medium md:table-cell">
              Image
            </th>
            <th scope="col" className="w-px px-4 py-3 font-medium whitespace-nowrap">
              État
            </th>
            <th scope="col" className="hidden px-4 py-3 font-medium lg:table-cell">
              Ports
            </th>
            <th scope="col" className="hidden px-4 py-3 font-medium xl:table-cell">
              Créé le
            </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-iron-800">
          {[...containers].sort(byStateThenName).map((c) => (
            <tr key={c.id} className="transition-colors hover:bg-iron-900/60">
              <td className="max-w-0 px-4 py-3 md:max-w-none">
                <span className="block truncate font-mono text-iron-100">{c.name}</span>
                <span className="block truncate font-mono text-xs text-iron-300 md:hidden">
                  {c.image}
                </span>
              </td>
              <td className="hidden max-w-xs truncate px-4 py-3 font-mono text-iron-300 md:table-cell">
                {c.image}
              </td>
              <td className="w-px px-4 py-3">
                <StateBadge state={c.state} />
              </td>
              <td className="hidden px-4 py-3 font-mono text-iron-300 lg:table-cell">
                {c.ports.length ? c.ports.map(formatPort).join(', ') : 'aucun'}
              </td>
              <td className="hidden px-4 py-3 whitespace-nowrap text-iron-300 xl:table-cell">
                {c.createdAt ? created.format(timestampDate(c.createdAt)) : 'inconnue'}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function ContainersSkeleton() {
  return (
    <div
      role="status"
      aria-live="polite"
      className="space-y-px overflow-hidden rounded-lg border border-iron-800"
    >
      <span className="sr-only">Chargement des conteneurs…</span>
      {[0, 1, 2].map((row) => (
        <div key={row} className="flex items-center gap-6 bg-iron-900/50 px-4 py-4">
          <span className="h-3 w-40 animate-pulse rounded bg-iron-800" />
          <span className="hidden h-3 w-56 animate-pulse rounded bg-iron-800 md:block" />
          <span className="h-3 w-20 animate-pulse rounded bg-iron-800" />
        </div>
      ))}
    </div>
  );
}
