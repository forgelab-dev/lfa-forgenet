import type { Timestamp } from '@bufbuild/protobuf/wkt';
import { useEffect, useRef, useState } from 'react';

import type { Deployment, DeploymentLogEntry } from '../../gen/forgenet/v1/deploy_pb';
import { describeError } from '../../lib/api/errors';
import { isTransient, retryDelay, sleep } from '../../lib/live/retry';
import { useApi } from '../auth/session';

/**
 * État de la liaison en direct :
 * - `connecting` : premier abonnement en cours ;
 * - `live` : le flux est ouvert ;
 * - `reconnecting` : le flux est coupé, nouvel essai sous peu ;
 * - `done` : le serveur a fermé le flux, le déploiement est terminé ;
 * - `error` : refus définitif (déploiement introuvable, session…).
 */
export type Link = 'connecting' | 'live' | 'reconnecting' | 'done' | 'error';

/** Garde-fou de mémoire : le serveur borne déjà un journal à 5 000 lignes. */
const maxLines = 5000;

/**
 * Ouvre un flux et le rouvre après une coupure, jusqu'à sa fin normale ou un
 * refus. `open` reçoit le signal d'annulation ; il doit reprendre là où le
 * flux précédent s'est arrêté.
 */
async function keepAlive(
  open: (signal: AbortSignal) => Promise<void>,
  signal: AbortSignal,
  setLink: (link: Link) => void,
  setError: (message: string) => void
): Promise<void> {
  for (let attempt = 0; !signal.aborted; attempt++) {
    try {
      await open(signal);
      if (!signal.aborted) {
        setLink('done');
      }
      return;
    } catch (err) {
      if (signal.aborted) {
        return;
      }
      if (!isTransient(err)) {
        setError(describeError(err));
        setLink('error');
        return;
      }
      setLink('reconnecting');
      await sleep(retryDelay(attempt), signal);
    }
  }
}

/**
 * Suit l'état d'un déploiement en direct (StreamDeployment). `onUpdate`
 * reçoit chaque nouvel état, pour tenir la liste à jour.
 */
export function useDeploymentWatch(
  deploymentId: string,
  enabled: boolean,
  onUpdate?: (deployment: Deployment) => void
) {
  const api = useApi();
  const [link, setLink] = useState<Link>('connecting');
  const [error, setError] = useState<string>();
  const update = useRef(onUpdate);
  useEffect(() => {
    update.current = onUpdate;
  });

  useEffect(() => {
    if (!enabled) {
      return;
    }
    const abort = new AbortController();
    void keepAlive(
      async (signal) => {
        for await (const res of api.deployments.streamDeployment({ deploymentId }, { signal })) {
          setLink('live');
          if (res.deployment) {
            update.current?.(res.deployment);
          }
        }
      },
      abort.signal,
      setLink,
      setError
    );
    return () => abort.abort();
  }, [api, deploymentId, enabled]);

  return { link, error };
}

/**
 * Lit le journal d'un déploiement : les lignes enregistrées puis, avec
 * `follow`, les nouvelles en direct. Après une coupure, il reprend après la
 * dernière ligne reçue : rien n'est perdu ni doublé. Le composant qui l'appelle
 * est remonté (key) pour un autre déploiement.
 */
export function useDeploymentLogs(deploymentId: string, follow: boolean) {
  const api = useApi();
  const [lines, setLines] = useState<DeploymentLogEntry[]>([]);
  const [link, setLink] = useState<Link>('connecting');
  const [error, setError] = useState<string>();

  useEffect(() => {
    const abort = new AbortController();
    let since: Timestamp | undefined;
    void keepAlive(
      async (signal) => {
        const stream = api.deployments.streamDeploymentLogs(
          { deploymentId, follow, since },
          { signal }
        );
        setLink('live');
        for await (const res of stream) {
          const entry = res.logEntry;
          if (!entry) {
            continue;
          }
          since = entry.emittedAt;
          setLines((current) =>
            current.length >= maxLines ? [...current.slice(1), entry] : [...current, entry]
          );
        }
      },
      abort.signal,
      setLink,
      setError
    );
    return () => abort.abort();
  }, [api, deploymentId, follow]);

  return { lines, link, error };
}
