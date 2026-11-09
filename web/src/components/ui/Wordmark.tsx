interface Props {
  className?: string;
}

/**
 * Marque ForgeNet : trois plaques empilées, la plus haute en or, comme les
 * couches d'un serveur que le produit pilote.
 */
export function Wordmark({ className = '' }: Props) {
  return (
    <span className={`inline-flex items-center gap-2.5 ${className}`}>
      <svg viewBox="0 0 24 24" aria-hidden="true" className="size-[1.1em] shrink-0">
        <rect x="5" y="3" width="14" height="5" rx="1.5" fill="var(--color-gold-400)" />
        <rect x="3" y="10" width="18" height="5" rx="1.5" fill="var(--color-iron-300)" />
        <rect x="3" y="17" width="18" height="4" rx="1.5" fill="var(--color-iron-500)" />
      </svg>
      <span className="font-bold tracking-tight">
        Forge<span className="font-normal text-iron-300">Net</span>
      </span>
    </span>
  );
}
