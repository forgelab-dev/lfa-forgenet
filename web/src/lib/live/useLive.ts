import { Code, ConnectError } from '@connectrpc/connect';
import { useEffect, useRef, useState } from 'react';

import { useSession } from '../../features/auth/session';
import { describeError } from '../api/errors';
import { isTransient, retryDelay, sleep } from './retry';
import { SocketClosed } from './socket';

/**
 * État d'une liaison en direct :
 * - `connecting` : première ouverture, ou reprise après une pause ;
 * - `live` : le flux est ouvert ;
 * - `reconnecting` : coupure, nouvel essai à `retryAt` ;
 * - `paused` : onglet caché, flux fermé pour ménager la connexion ;
 * - `gone` : ce qui était suivi a disparu (conteneur arrêté) ;
 * - `error` : refus définitif, avec sa raison.
 */
export type LiveStatus =
  | { state: 'connecting' }
  | { state: 'live' }
  | { state: 'reconnecting'; retryAt: number }
  | { state: 'paused' }
  | { state: 'gone' }
  | { state: 'error'; message: string };

/**
 * Une ouverture du flux : rattraper ce qui manque (historique), puis suivre
 * en direct jusqu'à la coupure. `onLive` signale que le direct a repris.
 * La promesse se résout à l'annulation et se rejette à la coupure.
 */
export type LiveOpen = (signal: AbortSignal, token: string, onLive: () => void) => Promise<void>;

/**
 * Tient un flux ouvert tant que la page est visible : il est rouvert après
 * une coupure (1, 2, 4, 8 puis 10 s d'attente) et fermé quand l'onglet est
 * caché, puis repris à son retour. Un jeton refusé ferme la session.
 * `open` doit être stable (useCallback).
 */
export function useLive(open: LiveOpen, enabled = true): LiveStatus {
  const { token, expire } = useSession();
  const [status, setStatus] = useState<LiveStatus>({ state: 'connecting' });
  const [visible, setVisible] = useState(() => !document.hidden);
  const session = useRef({ token, expire });
  useEffect(() => {
    session.current = { token, expire };
  });

  useEffect(() => {
    const onChange = () => setVisible(!document.hidden);
    document.addEventListener('visibilitychange', onChange);
    return () => document.removeEventListener('visibilitychange', onChange);
  }, []);

  useEffect(() => {
    if (!enabled) {
      return;
    }
    if (!visible) {
      return;
    }
    const abort = new AbortController();
    const { signal } = abort;

    void (async () => {
      await Promise.resolve();
      if (signal.aborted) return;
      setStatus({ state: 'connecting' });
      for (let attempt = 0; !signal.aborted; attempt++) {
        const current = session.current.token();
        if (!current) {
          session.current.expire();
          return;
        }
        try {
          await open(signal, current, () => {
            attempt = 0;
            setStatus({ state: 'live' });
          });
          return; // annulé
        } catch (err) {
          if (signal.aborted) return;
          if (err instanceof SocketClosed) {
            if (err.end === 'unauthorized') {
              session.current.expire();
              return;
            }
            if (err.end === 'gone') {
              setStatus({ state: 'gone' });
              return;
            }
          } else if (ConnectError.from(err).code === Code.Unauthenticated) {
            return; // l'intercepteur de session a déjà fermé la session
          } else if (!isTransient(err)) {
            setStatus({ state: 'error', message: describeError(err) });
            return;
          }
          const delay = retryDelay(attempt);
          setStatus({ state: 'reconnecting', retryAt: Date.now() + delay });
          await sleep(delay, signal);
        }
      }
    })();
    return () => abort.abort();
  }, [open, enabled, visible]);

  return enabled && !visible ? { state: 'paused' } : status;
}
