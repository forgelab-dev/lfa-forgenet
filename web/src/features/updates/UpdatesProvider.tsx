import { useCallback, useEffect, useMemo, useState, type ReactNode } from 'react';

import { UserRole } from '../../gen/forgenet/v1/auth_pb';
import type { GetUpdateStatusResponse } from '../../gen/forgenet/v1/update_pb';
import { describeError } from '../../lib/api/errors';
import { useApi, useCurrentUser } from '../auth/session';
import { UpdatesContext, type Updates } from './updates';

/** Relecture de l'état : la vérification serveur, elle, est quotidienne. */
const REFRESH_MS = 30 * 60 * 1000;

/**
 * État des mises à jour, partagé par le bandeau et la barre latérale. Réservé
 * aux administrateurs (ADR-0011) : pour un membre, aucun appel n'est fait. Un
 * serveur qui ne connaît pas encore le service, ou qui ne répond pas, ne
 * montre simplement rien : ce n'est pas une alerte.
 */
export function UpdatesProvider({ children }: { children: ReactNode }) {
  const api = useApi();
  const user = useCurrentUser();
  const admin = user.role === UserRole.ADMIN;
  const [status, setStatus] = useState<GetUpdateStatusResponse>();
  const [checking, setChecking] = useState(false);
  const [checkError, setCheckError] = useState<string>();

  const refresh = useCallback(async () => {
    try {
      setStatus(await api.updates.getUpdateStatus({}));
    } catch {
      // Silencieux : voir le commentaire du composant.
    }
  }, [api]);

  useEffect(() => {
    if (!admin) {
      return;
    }
    let active = true;
    const load = async () => {
      try {
        const next = await api.updates.getUpdateStatus({});
        if (active) {
          setStatus(next);
        }
      } catch {
        // Silencieux.
      }
    };
    void load();
    const timer = setInterval(() => void load(), REFRESH_MS);
    return () => {
      active = false;
      clearInterval(timer);
    };
  }, [admin, api]);

  const check = useCallback(async () => {
    setChecking(true);
    setCheckError(undefined);
    try {
      const res = await api.updates.checkForUpdates({});
      setStatus(res.status);
    } catch (err) {
      setCheckError(describeError(err));
    } finally {
      setChecking(false);
    }
  }, [api]);

  const value = useMemo<Updates | null>(
    () => (admin ? { status, refresh, check, checking, checkError } : null),
    [admin, status, refresh, check, checking, checkError]
  );

  return <UpdatesContext.Provider value={value}>{children}</UpdatesContext.Provider>;
}
