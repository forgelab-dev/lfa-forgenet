import { useEffect, useState } from 'react';

/**
 * Heure courante, rafraîchie toutes les `every` millisecondes : les courbes
 * glissent et les comptes à rebours avancent sans lire l'horloge pendant le
 * rendu.
 */
export function useNow(every = 1000): number {
  const [now, setNow] = useState(() => Date.now());
  useEffect(() => {
    const timer = setInterval(() => setNow(Date.now()), every);
    return () => clearInterval(timer);
  }, [every]);
  return now;
}
