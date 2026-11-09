import { useEffect, type ReactNode } from 'react';

interface Props {
  title: string;
  description?: ReactNode;
}

/**
 * En-tête d'une page : son titre (le seul h1) et une phrase qui dit à quoi elle
 * sert. Le titre de l'onglet suit la page, pour qui en garde plusieurs ouverts.
 */
export function PageHeader({ title, description }: Props) {
  useEffect(() => {
    document.title = `${title} · ForgeNet`;
  }, [title]);

  return (
    <header className="animate-rise border-b border-iron-800 pb-6">
      <h1 className="text-[clamp(1.75rem,3vw,2.25rem)] leading-tight font-semibold tracking-[-0.02em]">
        {title}
      </h1>
      {description && <p className="mt-2 max-w-prose text-iron-300">{description}</p>}
    </header>
  );
}
