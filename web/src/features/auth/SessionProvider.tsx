import { Code, ConnectError } from '@connectrpc/connect';
import { timestampMs } from '@bufbuild/protobuf/wkt';
import { useCallback, useEffect, useMemo, useState, type ReactNode } from 'react';

import { createApiClients, type TransportFactory } from '../../lib/api/client';
import { describeError } from '../../lib/api/errors';
import { SessionContext, type OpenedSession, type Session, type SessionState } from './session';
import { createTokenStore } from './storage';

/** Durée retenue si le serveur ne donne pas la fin de validité du jeton. */
const FALLBACK_LIFETIME_MS = 60 * 60 * 1000;

interface Props {
  transportFactory: TransportFactory;
  children: ReactNode;
}

/**
 * Tient la session : le jeton, le compte connecté et les clients de l'API qui
 * présentent ce jeton. Un refus du jeton par le serveur, sur n'importe quel
 * appel, ramène à l'écran de connexion.
 */
export function SessionProvider({ transportFactory, children }: Props) {
  const [tokens] = useState(createTokenStore);
  const [state, setState] = useState<SessionState>(() =>
    tokens.get() ? { status: 'checking' } : { status: 'anonymous' }
  );

  const close = useCallback(
    (next: SessionState) => {
      tokens.clear();
      setState(next);
    },
    [tokens]
  );

  const api = useMemo(
    () =>
      createApiClients(
        transportFactory({
          token: tokens.get,
          onSessionExpired: () => close({ status: 'anonymous', notice: 'expired' }),
        })
      ),
    [transportFactory, tokens, close]
  );

  useEffect(() => {
    if (state.status !== 'checking') {
      return;
    }
    const abort = new AbortController();
    api.auth
      .getCurrentUser({}, { signal: abort.signal })
      .then(({ user }) => {
        if (user) {
          setState({ status: 'authenticated', user });
        } else {
          close({ status: 'anonymous', notice: 'expired' });
        }
      })
      .catch((err: unknown) => {
        const { code } = ConnectError.from(err);
        // Annulation : l'effet a été démonté. Refus du jeton : l'intercepteur
        // a déjà fermé la session.
        if (code !== Code.Canceled && code !== Code.Unauthenticated) {
          setState({ status: 'unreachable', message: describeError(err) });
        }
      });
    return () => abort.abort();
  }, [state.status, api, close]);

  const adopt = useCallback(
    (opened: OpenedSession) => {
      if (!opened.user || !opened.token) {
        throw new ConnectError('réponse de connexion incomplète', Code.Internal);
      }
      const expiresAt = opened.expiresAt
        ? timestampMs(opened.expiresAt)
        : Date.now() + FALLBACK_LIFETIME_MS;
      tokens.set({ token: opened.token, expiresAt });
      setState({ status: 'authenticated', user: opened.user });
    },
    [tokens]
  );

  const login = useCallback(
    async (login: string, password: string) => {
      adopt(await api.auth.login({ login, password }));
    },
    [api, adopt]
  );

  const logout = useCallback(async () => {
    try {
      await api.auth.logout({});
    } catch {
      // Jeton déjà refusé ou serveur injoignable : la session locale se
      // ferme quand même, le jeton expirera de lui-même côté serveur.
    }
    close({ status: 'anonymous' });
  }, [api, close]);

  const retry = useCallback(() => setState({ status: 'checking' }), []);
  const expire = useCallback(() => close({ status: 'anonymous', notice: 'expired' }), [close]);

  const session = useMemo<Session>(
    () => ({ state, api, login, adopt, logout, retry, token: tokens.get, expire }),
    [state, api, login, adopt, logout, retry, tokens, expire]
  );

  return <SessionContext.Provider value={session}>{children}</SessionContext.Provider>;
}
