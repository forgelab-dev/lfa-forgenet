import { ArrowLeft } from 'lucide-react';
import { Link } from 'wouter';

import { iconSource } from './catalog';

/**
 * Icône d'un template, posée sur une plaque sombre : les icônes des projets
 * sont souvent faites pour un fond clair comme pour un fond sombre.
 */
export function TemplateIcon({ svg, className = '' }: { svg: string; className?: string }) {
  return (
    <span
      className={`flex shrink-0 items-center justify-center rounded-md border border-iron-700 bg-iron-950 p-2 ${className}`}
    >
      <img src={iconSource(svg)} alt="" className="size-full object-contain" />
    </span>
  );
}

/** Lien vers le catalogue ; `className` ajuste ses marges selon la page. */
export function BackLink({ className = '-ml-2 mb-4' }: { className?: string }) {
  return (
    <Link
      href="/marketplace"
      className={`${className} inline-flex min-h-11 items-center gap-1.5 rounded-md px-2 text-sm font-medium text-iron-300 transition-colors hover:text-gold-300 focus-visible:outline-2 focus-visible:outline-gold-400`}
    >
      <ArrowLeft aria-hidden="true" className="size-4" />
      Retour au catalogue
    </Link>
  );
}
