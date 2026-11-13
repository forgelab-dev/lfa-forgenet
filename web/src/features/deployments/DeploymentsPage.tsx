import { ShieldCheck } from 'lucide-react';
import { useCallback, useEffect, useMemo, useState } from 'react';

import { PageHeader } from '../../components/ui/PageHeader';
import { RequestError } from '../../components/ui/RequestError';
import { UserRole } from '../../gen/forgenet/v1/auth_pb';
import type { Deployment } from '../../gen/forgenet/v1/deploy_pb';
import { useRequest } from '../../lib/useRequest';
import { useApi, useCurrentUser } from '../auth/session';
import { ApplicationList } from './ApplicationList';
import { DeployForm } from './DeployForm';
import { groupByProject, type Stacks } from './group';
import { isActive } from './state';

/** Les déploiements chargés : le serveur n'en renvoie pas plus de 100 par page. */
const pageSize = 100;

export function DeploymentsPage() {
  const isAdmin = useCurrentUser().role === UserRole.ADMIN;
  return (
    <>
      <PageHeader
        title="Déploiements"
        description="Les applications Docker Compose de ce serveur. Chaque déploiement se suit en direct, étape par étape, avec son journal."
      />
      <div className="mt-8">{isAdmin ? <AdminView /> : <MemberNotice />}</div>
    </>
  );
}

function AdminView() {
  const api = useApi();
  const load = useCallback(
    async (signal: AbortSignal) =>
      (await api.deployments.listDeployments({ pageSize }, { signal })).deployments,
    [api]
  );
  const [request, retry] = useRequest(load);

  // Les stacks en place, relues à la fin d'un déploiement et après une
  // suppression. Leur échec n'empêche pas d'afficher les déploiements.
  const [stacks, setStacks] = useState<Stacks>({ known: false });
  const [stacksVersion, setStacksVersion] = useState(0);
  const refreshStacks = useCallback(() => setStacksVersion((v) => v + 1), []);
  useEffect(() => {
    const abort = new AbortController();
    api.deployments.listStacks({}, { signal: abort.signal }).then(
      (res) => {
        if (!abort.signal.aborted) {
          setStacks({ known: true, byProject: new Map(res.stacks.map((s) => [s.projectId, s])) });
        }
      },
      () => {
        if (!abort.signal.aborted) setStacks({ known: false });
      }
    );
    return () => abort.abort();
  }, [api, stacksVersion]);

  // Les créations et les états reçus en direct s'appliquent à la liste
  // chargée, sans la recharger ; un nouveau chargement les reprend par id.
  const [updates, setUpdates] = useState<Map<string, Deployment>>(() => new Map());
  const applyUpdate = useCallback((d: Deployment) => {
    setUpdates((current) => new Map(current).set(d.id, d));
    // Un déploiement terminé a pu créer, changer ou laisser à moitié sa stack.
    if (!isActive(d.state)) {
      setStacksVersion((v) => v + 1);
    }
  }, []);

  const applications = useMemo(() => {
    if (request.status !== 'success') {
      return [];
    }
    const byId = new Map(request.data.map((d) => [d.id, d]));
    for (const [id, d] of updates) {
      byId.set(id, d);
    }
    return groupByProject(byId.values(), stacks);
  }, [request, updates, stacks]);

  return (
    <>
      {request.status === 'loading' && <DeploymentsSkeleton />}
      {request.status === 'error' && <RequestError error={request.error} onRetry={retry} />}
      {request.status === 'success' && (
        <div className="grid items-start gap-8 lg:grid-cols-[minmax(0,1.6fr)_minmax(0,1fr)]">
          <section
            aria-labelledby="applications"
            className="animate-rise min-w-0 [animation-delay:60ms]"
          >
            <h2 id="applications" className="text-lg font-semibold">
              {applications.length === 0
                ? 'Aucune application'
                : applications.length === 1
                  ? '1 application'
                  : `${applications.length} applications`}
            </h2>
            <ApplicationList
              applications={applications}
              stacksKnown={stacks.known}
              onUpdate={applyUpdate}
              onStackDeleted={refreshStacks}
            />
          </section>
          <aside className="animate-rise [animation-delay:120ms] lg:sticky lg:top-6">
            <DeployForm onCreated={applyUpdate} />
          </aside>
        </div>
      )}
    </>
  );
}

function MemberNotice() {
  return (
    <div className="animate-rise flex max-w-2xl items-start gap-4 rounded-lg border border-iron-700 bg-iron-900 px-5 py-6">
      <span className="flex size-10 shrink-0 items-center justify-center rounded-md bg-iron-850 text-gold-400">
        <ShieldCheck aria-hidden="true" className="size-5" />
      </span>
      <div>
        <p className="font-semibold">Réservé aux administrateurs</p>
        <p className="mt-1 text-sm text-iron-300">
          Un fichier Compose donne un accès complet au serveur et peut contenir des secrets : seul
          un administrateur peut voir et lancer les déploiements.
        </p>
      </div>
    </div>
  );
}

function DeploymentsSkeleton() {
  return (
    <div role="status" aria-live="polite" className="space-y-3 lg:max-w-[62%]">
      <span className="sr-only">Chargement des déploiements…</span>
      {[0, 1].map((row) => (
        <div
          key={row}
          className="space-y-2 rounded-lg border border-iron-800 bg-iron-900/50 px-4 py-4"
        >
          <span className="block h-3.5 w-40 animate-pulse rounded bg-iron-800" />
          <span className="block h-3 w-24 animate-pulse rounded bg-iron-800" />
        </div>
      ))}
    </div>
  );
}
