import { ArrowRight } from 'lucide-react';
import { useCallback } from 'react';
import { Link } from 'wouter';

import { PageHeader } from '../../components/ui/PageHeader';
import { RequestError } from '../../components/ui/RequestError';
import { ServerPlate } from '../../components/ui/ServerPlate';
import type { Container } from '../../gen/forgenet/v1/container_pb';
import { useRequest } from '../../lib/useRequest';
import { useApi, useCurrentUser } from '../auth/session';
import { listAllContainers } from '../containers/api';
import { containerStates, type Tone } from '../containers/state';

const tones: { tone: Tone; label: (n: number) => string; bar: string }[] = [
  { tone: 'ok', label: (n) => `${n} en marche`, bar: 'bg-ok-400' },
  { tone: 'warn', label: (n) => `${n} en pause ou en redémarrage`, bar: 'bg-warn-400' },
  { tone: 'danger', label: (n) => `${n} hors service`, bar: 'bg-danger-400' },
  { tone: 'idle', label: (n) => `${n} à l’arrêt`, bar: 'bg-iron-700' },
];

export function OverviewPage() {
  const user = useCurrentUser();
  const api = useApi();
  const load = useCallback((signal: AbortSignal) => listAllContainers(api, signal), [api]);
  const [request, retry] = useRequest(load);

  return (
    <>
      <PageHeader
        title="Vue d’ensemble"
        description={`Bonjour ${user.username}. Voici l’état du serveur au moment du chargement.`}
      />

      <div className="mt-8 grid gap-8 lg:grid-cols-[minmax(0,1.6fr)_minmax(0,1fr)]">
        <section
          aria-labelledby="overview-containers"
          className="animate-rise [animation-delay:60ms]"
        >
          <div className="flex min-h-11 items-center justify-between gap-4">
            <h2 id="overview-containers" className="text-lg font-semibold">
              Conteneurs
            </h2>
            <Link
              href="/conteneurs"
              className="inline-flex min-h-11 items-center gap-1.5 text-sm font-medium text-gold-400 hover:text-gold-300"
            >
              Voir la liste
              <ArrowRight aria-hidden="true" className="size-4" />
            </Link>
          </div>
          <div className="mt-3">
            {request.status === 'loading' && (
              <div role="status" className="h-24 animate-pulse rounded-lg bg-iron-900">
                <span className="sr-only">Chargement des conteneurs…</span>
              </div>
            )}
            {request.status === 'error' && <RequestError error={request.error} onRetry={retry} />}
            {request.status === 'success' && <ContainerSummary containers={request.data} />}
          </div>
        </section>

        <section aria-labelledby="overview-server" className="animate-rise [animation-delay:120ms]">
          <h2 id="overview-server" className="flex min-h-11 items-center text-lg font-semibold">
            Ce serveur
          </h2>
          <ServerPlate className="mt-3" />
        </section>
      </div>
    </>
  );
}

function ContainerSummary({ containers }: { containers: Container[] }) {
  if (containers.length === 0) {
    return (
      <p className="rounded-lg border border-dashed border-iron-700 px-5 py-8 text-iron-300">
        Aucun conteneur pour l’instant : le serveur est prêt à recevoir un premier déploiement.
      </p>
    );
  }

  const counts = new Map<Tone, number>();
  for (const c of containers) {
    const { tone } = containerStates[c.state];
    counts.set(tone, (counts.get(tone) ?? 0) + 1);
  }
  const present = tones.filter(({ tone }) => counts.get(tone));

  return (
    <div className="rounded-lg border border-iron-800 bg-iron-900 px-5 py-5">
      <p className="text-iron-100">
        <span className="font-mono text-2xl font-semibold">{containers.length}</span>{' '}
        {containers.length > 1 ? 'conteneurs sur ce serveur' : 'conteneur sur ce serveur'}
      </p>
      <div aria-hidden="true" className="mt-4 flex h-2 gap-0.5 overflow-hidden rounded-full">
        {present.map(({ tone, bar }) => (
          <span key={tone} className={bar} style={{ flexGrow: counts.get(tone) }} />
        ))}
      </div>
      <ul className="mt-4 flex flex-wrap gap-x-6 gap-y-2 text-sm text-iron-300">
        {present.map(({ tone, label, bar }) => (
          <li key={tone} className="inline-flex items-center gap-2">
            <span aria-hidden="true" className={`size-2 rounded-full ${bar}`} />
            {label(counts.get(tone) ?? 0)}
          </li>
        ))}
      </ul>
    </div>
  );
}
