import { ArrowUpRight, HardDrive, Network } from 'lucide-react';
import { useCallback, useState } from 'react';

import { PageHeader } from '../../components/ui/PageHeader';
import { RequestError } from '../../components/ui/RequestError';
import { UserRole } from '../../gen/forgenet/v1/auth_pb';
import type { InstallTemplateResponse, Template } from '../../gen/forgenet/v1/marketplace_pb';
import { useRequest } from '../../lib/useRequest';
import { useApi, useCurrentUser } from '../auth/session';
import { categoryLabel } from './catalog';
import { InstallForm } from './InstallForm';
import { Installation } from './Installation';
import { MarketplacePage } from './MarketplacePage';
import { BackLink, TemplateIcon } from './parts';

/** Page d'un template : ce qu'il installe, son formulaire, puis son installation. */
export function TemplatePage({ params }: { params: { id: string } }) {
  const isAdmin = useCurrentUser().role === UserRole.ADMIN;
  // Un membre voit la même mention que sur le catalogue, sans appel refusé.
  if (!isAdmin) return <MarketplacePage />;
  return <AdminTemplate key={params.id} id={params.id} />;
}

function AdminTemplate({ id }: { id: string }) {
  const api = useApi();
  const load = useCallback(
    async (signal: AbortSignal) => (await api.marketplace.getTemplate({ id }, { signal })).template,
    [api, id]
  );
  const [request, retry] = useRequest(load);
  const [installed, setInstalled] = useState<InstallTemplateResponse>();

  return (
    <>
      <BackLink />
      {request.status === 'loading' && (
        <p role="status" className="mt-6 text-iron-300">
          Chargement du template…
        </p>
      )}
      {request.status === 'error' && (
        <div className="mt-6">
          <RequestError error={request.error} onRetry={retry} />
        </div>
      )}
      {request.status === 'success' && request.data && (
        <>
          <TemplateHeader template={request.data} />
          {installed?.deployment ? (
            <Installation
              template={request.data}
              initial={installed.deployment}
              secrets={installed.generatedSecrets}
            />
          ) : (
            <div className="mt-8 grid items-start gap-8 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.15fr)]">
              <Contents template={request.data} />
              <InstallForm template={request.data} onInstalled={setInstalled} />
            </div>
          )}
        </>
      )}
    </>
  );
}

function TemplateHeader({ template: t }: { template: Template }) {
  return (
    <PageHeader
      title={t.name}
      icon={<TemplateIcon svg={t.iconSvg} className="size-14" />}
      description={t.description}
      aside={
        <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1 text-sm">
          <dt className="text-iron-300">Version</dt>
          <dd className="font-mono">{t.version}</dd>
          <dt className="text-iron-300">Catégorie</dt>
          <dd>{categoryLabel(t.category)}</dd>
          {t.license && (
            <>
              <dt className="text-iron-300">Licence</dt>
              <dd className="font-mono">{t.license}</dd>
            </>
          )}
          {t.homepage && (
            <>
              <dt className="text-iron-300">Site</dt>
              <dd>
                <a
                  href={t.homepage}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1 text-gold-400 hover:text-gold-300"
                >
                  {new URL(t.homepage).host}
                  <ArrowUpRight aria-hidden="true" className="size-3.5" />
                  <span className="sr-only">(ouvre un nouvel onglet)</span>
                </a>
              </dd>
            </>
          )}
        </dl>
      }
    />
  );
}

/** Ce que le template met en place : ses ports internes et ses volumes. */
function Contents({ template: t }: { template: Template }) {
  return (
    <section aria-labelledby="contents" className="animate-rise [animation-delay:60ms]">
      <h2 id="contents" className="text-lg font-semibold">
        Ce qui sera installé
      </h2>
      <p className="mt-2 text-sm text-iron-300">
        Une application Docker Compose, déployée et suivie comme les autres dans la page
        Déploiements. Aucun port n’est ouvert sur Internet : un service web se relie à un nom de
        domaine depuis la page Domaines.
      </p>

      {t.ports.length > 0 && (
        <>
          <h3 className="mt-6 flex items-center gap-2 text-sm font-semibold">
            <Network aria-hidden="true" className="size-4 text-gold-400" />
            Ports
          </h3>
          <ul className="mt-2 divide-y divide-iron-800 rounded-md border border-iron-800">
            {t.ports.map((p) => (
              <li
                key={`${p.service}:${p.port}`}
                className="flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1 px-3 py-2.5"
              >
                <span className="font-mono text-sm text-iron-100">
                  {p.service}:{p.port}
                </span>
                <span className="text-sm text-iron-300">
                  {p.description}
                  <span className="ml-2 font-mono text-xs uppercase text-iron-500">
                    {p.protocol}
                  </span>
                </span>
              </li>
            ))}
          </ul>
        </>
      )}

      {t.volumes.length > 0 && (
        <>
          <h3 className="mt-6 flex items-center gap-2 text-sm font-semibold">
            <HardDrive aria-hidden="true" className="size-4 text-gold-400" />
            Données conservées
          </h3>
          <ul className="mt-2 divide-y divide-iron-800 rounded-md border border-iron-800">
            {t.volumes.map((v) => (
              <li
                key={v.name}
                className="flex flex-wrap items-baseline justify-between gap-x-4 gap-y-1 px-3 py-2.5"
              >
                <span className="font-mono text-sm text-iron-100">{v.name}</span>
                <span className="text-sm text-iron-300">{v.description}</span>
              </li>
            ))}
          </ul>
          <p className="mt-2 text-xs text-iron-300">
            Les données restent d’un redéploiement à l’autre. À la suppression de l’application,
            vous choisissez de les garder ou de les effacer.
          </p>
        </>
      )}
    </section>
  );
}
