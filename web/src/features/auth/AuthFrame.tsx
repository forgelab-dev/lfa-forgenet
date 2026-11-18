import type { ReactNode } from 'react';

import { ServerPlate } from '../../components/ui/ServerPlate';
import { Wordmark } from '../../components/ui/Wordmark';

/**
 * Cadre des écrans hors session (connexion, lien à usage unique) : la marque
 * et la plaque du serveur à gauche, le formulaire à droite.
 */
export function AuthFrame({ children }: { children: ReactNode }) {
  return (
    <div className="grid min-h-dvh lg:grid-cols-[minmax(0,1.15fr)_minmax(26rem,1fr)]">
      <section className="engraved relative flex flex-col justify-between gap-10 px-6 pt-10 pb-8 sm:px-10 lg:px-16 lg:py-16">
        <Wordmark className="animate-rise text-xl" />
        <div className="animate-rise max-w-xl [animation-delay:80ms]">
          <p className="text-[clamp(2.25rem,5vw,4.25rem)] leading-[1.02] font-bold tracking-[-0.03em] text-balance">
            Votre serveur, <span className="text-gold-400">sous la main.</span>
          </p>
          <p className="mt-5 max-w-md text-lg text-iron-300">
            Conteneurs, déploiements et métriques du serveur, depuis un seul tableau de bord.
          </p>
        </div>
        <ServerPlate className="animate-rise max-w-md [animation-delay:160ms]" />
      </section>

      <main className="flex items-center justify-center border-iron-800 bg-iron-900 px-6 py-12 sm:px-10 lg:border-l">
        {children}
      </main>
    </div>
  );
}
