import { ArrowLeft, Hammer } from 'lucide-react';
import { Link } from 'wouter';

import { PageHeader } from '../components/ui/PageHeader';
import type { Section } from './navigation';

/** Section prévue mais pas encore construite : ce qu'elle fera, sans faux contenu. */
export function UpcomingPage({ section }: { section: Section }) {
  const Icon = section.icon;
  return (
    <>
      <PageHeader title={section.label} description={section.summary} />
      <div className="animate-rise mt-8 flex max-w-2xl items-start gap-4 rounded-lg border border-dashed border-iron-700 px-5 py-6 [animation-delay:60ms]">
        <span className="flex size-10 shrink-0 items-center justify-center rounded-md bg-iron-900 text-iron-300">
          <Icon aria-hidden="true" className="size-5" />
        </span>
        <div>
          <p className="flex items-center gap-2 font-semibold">
            <Hammer aria-hidden="true" className="size-4 text-gold-400" />
            En cours de construction
          </p>
          <p className="mt-1 text-sm text-iron-300">
            Cette section arrive dans une prochaine version de ForgeNet.
          </p>
        </div>
      </div>
    </>
  );
}

export function NotFoundPage() {
  return (
    <>
      <PageHeader
        title="Page introuvable"
        description="Cette adresse ne correspond à aucune page du tableau de bord."
      />
      <Link
        href="/"
        className="mt-6 inline-flex min-h-11 items-center gap-2 text-sm font-medium text-gold-400 hover:text-gold-300"
      >
        <ArrowLeft aria-hidden="true" className="size-4" />
        Revenir à la vue d’ensemble
      </Link>
    </>
  );
}
