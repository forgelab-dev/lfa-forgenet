/**
 * Flux WebSocket des métriques (ADR-0005, ADR-0010) : le jeton passe en
 * sous-protocole, chaque message est un JSON du contrat api/.
 */

export const SUBPROTOCOL = 'forgenet.v1';
export const BEARER_PREFIX = 'forgenet.bearer.';

/** Codes de fermeture propres à ForgeNet (ADR-0005). */
export const CLOSE_UNAUTHORIZED = 4401;
export const CLOSE_NOT_FOUND = 4404;

/**
 * Fin d'un flux :
 * - `transient` : coupure réseau, serveur arrêté ou poignée de main refusée ;
 *   on réessaie. Le navigateur ne donne pas le statut HTTP d'un refus : un
 *   jeton périmé se voit à l'appel d'historique qui précède chaque reprise ;
 * - `unauthorized` : le serveur a refusé le jeton pendant le flux ;
 * - `gone` : le conteneur suivi a disparu.
 */
export type SocketEnd = 'transient' | 'unauthorized' | 'gone';

export class SocketClosed extends Error {
  constructor(
    readonly end: SocketEnd,
    readonly code: number
  ) {
    super(`flux fermé (${code})`);
    this.name = 'SocketClosed';
  }
}

export function endOf(code: number): SocketEnd {
  if (code === CLOSE_UNAUTHORIZED) return 'unauthorized';
  if (code === CLOSE_NOT_FOUND) return 'gone';
  return 'transient';
}

/** Adresse ws:// ou wss:// d'une route, à la même origine que la page (ADR-0006). */
export function socketUrl(
  path: string,
  location: Pick<Location, 'protocol' | 'host'> = window.location
) {
  const scheme = location.protocol === 'https:' ? 'wss:' : 'ws:';
  return `${scheme}//${location.host}${path}`;
}

/**
 * Ouvre un flux et transmet chaque message décodé. La promesse se rejette
 * toujours par un SocketClosed à la fermeture, sauf annulation par `signal`,
 * où elle se résout.
 */
export function readSocket(
  path: string,
  token: string,
  onMessage: (data: unknown) => void,
  signal: AbortSignal,
  onOpen?: () => void
): Promise<void> {
  return new Promise((resolve, reject) => {
    if (signal.aborted) {
      resolve();
      return;
    }
    const socket = new WebSocket(socketUrl(path), [SUBPROTOCOL, BEARER_PREFIX + token]);
    const abort = () => {
      socket.close(1000);
      resolve();
    };
    signal.addEventListener('abort', abort, { once: true });
    socket.onopen = () => onOpen?.();
    socket.onmessage = (event: MessageEvent) => {
      if (typeof event.data !== 'string') return;
      try {
        onMessage(JSON.parse(event.data));
      } catch {
        // Message illisible : ignoré, le suivant le remplacera.
      }
    };
    socket.onclose = (event: CloseEvent) => {
      signal.removeEventListener('abort', abort);
      if (signal.aborted) {
        resolve();
        return;
      }
      reject(new SocketClosed(endOf(event.code), event.code));
    };
  });
}
