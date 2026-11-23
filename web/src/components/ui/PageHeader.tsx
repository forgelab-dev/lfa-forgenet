import { useEffect, type ReactNode } from 'react';

interface Props {
  title: string;
  /** Devant le titre : l'icône de ce que la page présente. */
  icon?: ReactNode;
  description?: ReactNode;
  /** À droite du titre sur grand écran, dessous sur téléphone : un état, une action. */
  aside?: ReactNode;
}

/**
 * En-tête d'une page : son titre (le seul h1) et une phrase qui dit à quoi elle
 * sert. Le titre de l'onglet suit la page, pour qui en garde plusieurs ouverts.
 */
export function PageHeader({ title, icon, description, aside }: Props) {
  useEffect(() => {
    document.title = `${title} · ForgeNet`;
  }, [title]);

  return (
    <header className="animate-rise flex flex-wrap items-end justify-between gap-x-6 gap-y-4 border-b border-iron-800 pb-6">
      <div className="flex min-w-0 flex-[1_1_24rem] items-start gap-4">
        {icon}
        <div className="min-w-0">
          <h1 className="text-[clamp(1.75rem,3vw,2.25rem)] leading-tight font-semibold tracking-[-0.02em]">
            {title}
          </h1>
          {description && <p className="mt-2 max-w-prose text-iron-300">{description}</p>}
        </div>
      </div>
      {aside}
    </header>
  );
}
