import { createContext, useContext } from 'react';

import type { User } from '../../gen/forgenet/v1/auth_pb';
import type { ApiClients } from '../../lib/api/client';

/**
 * État de la session du dashboard.
 *
 * - `checking` : un jeton conservé est en cours de vérification ;
 * - `unreachable` : le jeton n'a pas pu être vérifié (serveur injoignable) ;
 * - `anonymous` : pas de session, l'écran de connexion s'affiche ;
 * - `authenticated` : la session est ouverte.
 */
export type SessionState =
  | { status: 'checking' }
  | { status: 'unreachable'; message: string }
  | { status: 'anonymous'; notice?: 'expired' }
  | { status: 'authenticated'; user: User };

export interface Session {
  state: SessionState;
  api: ApiClients;
  /** Ouvre une session ; rejette avec l'erreur du serveur en cas de refus. */
  login(login: string, password: string): Promise<void>;
  /** Ferme la session, même si le serveur ne répond pas. */
  logout(): Promise<void>;
  /** Relance la vérification du jeton conservé. */
  retry(): void;
}

export const SessionContext = createContext<Session | null>(null);

export function useSession(): Session {
  const session = useContext(SessionContext);
  if (!session) {
    throw new Error('useSession doit être appelé sous <SessionProvider>');
  }
  return session;
}

/** Clients de l'API de la session en cours. */
export function useApi(): ApiClients {
  return useSession().api;
}

/** Compte connecté ; à n'appeler que dans l'espace authentifié. */
export function useCurrentUser(): User {
  const { state } = useSession();
  if (state.status !== 'authenticated') {
    throw new Error('useCurrentUser appelé hors d’une session ouverte');
  }
  return state.user;
}
