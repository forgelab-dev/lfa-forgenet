interface Props {
  className?: string;
}

/**
 * Marque ForgeNet : l'enclume de la famille LucidForge, sur son socle, et une
 * étincelle d'or. Le même dessin sert aux icônes et à la vitrine
 * (docs/media/logo/) ; l'enclume prend la couleur du texte.
 */
export function Wordmark({ className = '' }: Props) {
  return (
    <span className={`inline-flex items-center gap-2.5 ${className}`}>
      <svg viewBox="0 0 48 48" aria-hidden="true" className="size-[1.25em] shrink-0">
        <path
          d="M3 12h33a7 7 0 0 1 7 7v2H31l-2 6h4l3 7H12l3-7h4l-2-6h-6c-4 0-6.5-3-8-9z"
          fill="currentColor"
        />
        <rect x="12" y="38" width="24" height="4" rx="2" fill="currentColor" opacity="0.55" />
        <path
          d="M41 2l1.7 3.8L46.5 7.5l-3.8 1.7L41 13l-1.7-3.8L35.5 7.5l3.8-1.7z"
          fill="var(--color-gold-400)"
        />
      </svg>
      <span className="font-bold tracking-tight">
        Forge<span className="font-normal text-iron-300">Net</span>
      </span>
    </span>
  );
}
