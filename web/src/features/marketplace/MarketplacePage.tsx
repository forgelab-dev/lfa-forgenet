import { ArrowUpRight, PackageOpen, ShieldCheck } from 'lucide-react';
import { useCallback, useState, type ReactNode } from 'react';
import { Link } from 'wouter';

import { PageHeader } from '../../components/ui/PageHeader';
import { RequestError } from '../../components/ui/RequestError';
import { UserRole } from '../../gen/forgenet/v1/auth_pb';
import type { Template } from '../../gen/forgenet/v1/marketplace_pb';
import { useRequest } from '../../lib/useRequest';
import { useApi, useCurrentUser } from '../auth/session';
import { categoriesOf, categoryLabel, inCategory } from './catalog';
import { TemplateIcon } from './parts';

/** Le guide qui explique comment proposer un template. */
const contributorGuide =
  'https://github.com/forgelab-dev/lfa-forgenet/blob/main/templates/README.md';

export function MarketplacePage() {
  const isAdmin = useCurrentUser().role === UserRole.ADMIN;
  return (
    <>
      <PageHeader
        title="Marketplace"
        description="Des applications prêtes à installer sur ce serveur : choisissez-en une, répondez à quelques questions, ForgeNet la déploie et crée ses mots de passe."
      />
      <div className="mt-8">{isAdmin ? <Catalog /> : <AdminOnly />}</div>
    </>
  );
}

function Catalog() {
  const api = useApi();
  const load = useCallback(
    async (signal: AbortSignal) => (await api.marketplace.listTemplates({}, { signal })).templates,
    [api]
  );
  const [request, retry] = useRequest(load);
  const [category, setCategory] = useState<string>();

  if (request.status === 'loading') return <CatalogSkeleton />;
  if (request.status === 'error') return <RequestError error={request.error} onRetry={retry} />;

  const templates = request.data;
  if (templates.length === 0) return <EmptyCatalog />;

  const categories = categoriesOf(templates);
  const shown = inCategory(templates, category);
  return (
    <>
      {categories.length > 1 && (
        <div
          role="group"
          aria-label="Catégories"
          className="animate-rise flex flex-wrap gap-2 [animation-delay:60ms]"
        >
          <FilterButton pressed={!category} onClick={() => setCategory(undefined)}>
            Tout <Count n={templates.length} />
          </FilterButton>
          {categories.map((c) => (
            <FilterButton key={c} pressed={category === c} onClick={() => setCategory(c)}>
              {categoryLabel(c)} <Count n={inCategory(templates, c).length} />
            </FilterButton>
          ))}
        </div>
      )}
      <ul aria-label="Templates" className="mt-6 grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
        {shown.map((t, i) => (
          <TemplateCard key={t.id} template={t} delay={90 + i * 40} />
        ))}
      </ul>
    </>
  );
}

function TemplateCard({ template: t, delay }: { template: Template; delay: number }) {
  return (
    <li className="animate-rise" style={{ animationDelay: `${Math.min(delay, 400)}ms` }}>
      <Link
        href={`/marketplace/${t.id}`}
        className="group flex h-full gap-4 rounded-lg border border-iron-800 bg-iron-900/40 px-4 py-4 transition-colors duration-150 hover:border-gold-500/60 hover:bg-iron-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gold-400"
      >
        <TemplateIcon svg={t.iconSvg} className="size-12" />
        <div className="flex min-w-0 flex-1 flex-col">
          <div className="flex items-baseline justify-between gap-3">
            <h3 className="truncate text-base font-semibold text-iron-100">{t.name}</h3>
            <span className="shrink-0 font-mono text-xs text-iron-300">{t.version}</span>
          </div>
          <p className="mt-0.5 text-xs text-gold-400">{categoryLabel(t.category)}</p>
          <p className="mt-2 line-clamp-3 text-sm text-iron-300">{t.description}</p>
          <span className="mt-auto inline-flex items-center gap-1 pt-3 text-sm font-medium text-iron-300 transition-colors group-hover:text-gold-300">
            Installer
            <ArrowUpRight
              aria-hidden="true"
              className="size-4 transition-transform duration-150 ease-out-quart group-hover:translate-x-0.5 group-hover:-translate-y-0.5"
            />
          </span>
        </div>
      </Link>
    </li>
  );
}

function FilterButton({
  pressed,
  onClick,
  children,
}: {
  pressed: boolean;
  onClick: () => void;
  children: ReactNode;
}) {
  return (
    <button
      type="button"
      aria-pressed={pressed}
      onClick={onClick}
      className={`inline-flex min-h-11 items-center gap-2 rounded-md border px-3 text-sm font-medium transition-colors duration-150 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gold-400 ${
        pressed
          ? 'border-gold-500 bg-gold-950 text-gold-300'
          : 'border-iron-700 text-iron-300 hover:border-iron-500 hover:text-iron-100'
      }`}
    >
      {children}
    </button>
  );
}

function Count({ n }: { n: number }) {
  return <span className="font-mono text-xs opacity-80">{n}</span>;
}

function EmptyCatalog() {
  return (
    <div className="animate-rise flex max-w-2xl items-start gap-4 rounded-lg border border-dashed border-iron-700 px-5 py-6 [animation-delay:60ms]">
      <span className="flex size-10 shrink-0 items-center justify-center rounded-md bg-iron-900 text-gold-400">
        <PackageOpen aria-hidden="true" className="size-5" />
      </span>
      <div>
        <p className="font-semibold">Aucun template pour l’instant</p>
        <p className="mt-1 text-sm text-iron-300">
          Le catalogue se remplit par les contributions de la communauté : chaque template est un
          dossier de <span className="font-mono text-[0.8125rem]">templates/</span>, vérifié puis
          embarqué dans la version suivante de ForgeNet.
        </p>
        <a
          href={contributorGuide}
          target="_blank"
          rel="noreferrer"
          className="mt-3 inline-flex min-h-11 items-center gap-1.5 text-sm font-medium text-gold-400 hover:text-gold-300"
        >
          Lire le guide des contributeurs
          <ArrowUpRight aria-hidden="true" className="size-4" />
          <span className="sr-only">(ouvre un nouvel onglet)</span>
        </a>
      </div>
    </div>
  );
}

function AdminOnly() {
  return (
    <div className="animate-rise flex max-w-2xl items-start gap-4 rounded-lg border border-iron-700 bg-iron-900 px-5 py-6">
      <span className="flex size-10 shrink-0 items-center justify-center rounded-md bg-iron-850 text-gold-400">
        <ShieldCheck aria-hidden="true" className="size-5" />
      </span>
      <div>
        <p className="font-semibold">Réservé aux administrateurs</p>
        <p className="mt-1 text-sm text-iron-300">
          Installer un template déploie une application et crée ses mots de passe : seul un
          administrateur peut parcourir et installer le catalogue.
        </p>
      </div>
    </div>
  );
}

function CatalogSkeleton() {
  return (
    <div role="status" aria-live="polite" className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
      <span className="sr-only">Chargement du catalogue…</span>
      {[0, 1, 2].map((card) => (
        <div
          key={card}
          className="flex gap-4 rounded-lg border border-iron-800 bg-iron-900/50 px-4 py-4"
        >
          <span className="size-12 shrink-0 animate-pulse rounded-md bg-iron-800" />
          <div className="flex-1 space-y-2">
            <span className="block h-3.5 w-28 animate-pulse rounded bg-iron-800" />
            <span className="block h-3 w-full animate-pulse rounded bg-iron-800" />
            <span className="block h-3 w-2/3 animate-pulse rounded bg-iron-800" />
          </div>
        </div>
      ))}
    </div>
  );
}
