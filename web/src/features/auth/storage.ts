/**
 * Conservation du jeton de session dans le navigateur.
 *
 * sessionStorage plutôt que localStorage : la session survit au
 * rechargement de la page mais pas à la fermeture de l'onglet, ce qui borne
 * la vie d'un jeton oublié sur un poste partagé. Un cookie HttpOnly serait
 * plus sûr encore, mais l'ADR-0006 laisse les cookies hors du sujet ; ce choix
 * pourra évoluer avec lui.
 */

const KEY = 'forgenet.session';

export interface StoredSession {
  token: string;
  /** Fin de validité, en millisecondes depuis l'époque Unix. */
  expiresAt: number;
}

/** Lit la session conservée, ou rien si elle est absente, illisible ou échue. */
export function readSession(now = Date.now()): StoredSession | undefined {
  try {
    const raw = window.sessionStorage.getItem(KEY);
    if (!raw) {
      return undefined;
    }
    const value: unknown = JSON.parse(raw);
    if (!isStoredSession(value) || value.expiresAt <= now) {
      window.sessionStorage.removeItem(KEY);
      return undefined;
    }
    return value;
  } catch {
    return undefined;
  }
}

export function writeSession(session: StoredSession): void {
  try {
    window.sessionStorage.setItem(KEY, JSON.stringify(session));
  } catch {
    // Stockage bloqué (navigation privée stricte) : la session tient en
    // mémoire jusqu'au prochain rechargement.
  }
}

export function clearSession(): void {
  try {
    window.sessionStorage.removeItem(KEY);
  } catch {
    // Rien à effacer si le stockage est inaccessible.
  }
}

function isStoredSession(value: unknown): value is StoredSession {
  if (typeof value !== 'object' || value === null) {
    return false;
  }
  const { token, expiresAt } = value as Record<string, unknown>;
  return typeof token === 'string' && token !== '' && typeof expiresAt === 'number';
}

/**
 * Jeton de la session en cours, en mémoire et dans sessionStorage. Il vit hors
 * de l'état React : l'intercepteur le lit à chaque appel sans recréer le
 * transport ni les clients.
 */
export interface TokenStore {
  get(): string | undefined;
  set(session: StoredSession): void;
  clear(): void;
}

export function createTokenStore(): TokenStore {
  let token = readSession()?.token;
  return {
    get: () => token,
    set(session) {
      token = session.token;
      writeSession(session);
    },
    clear() {
      token = undefined;
      clearSession();
    },
  };
}
