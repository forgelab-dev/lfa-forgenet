import { Code, ConnectError } from '@connectrpc/connect';

/**
 * Traduit une erreur d'appel en phrase lisible par l'administrateur.
 *
 * Le serveur rédige déjà en français les refus qu'il veut montrer (mot de
 * passe incorrect, trop d'essais, requête invalide) : on les reprend. Les
 * autres codes reçoivent une phrase d'ici, sans détail interne.
 */
export function describeError(err: unknown): string {
  const error = ConnectError.from(err);
  switch (error.code) {
    case Code.Unauthenticated:
    case Code.ResourceExhausted:
    case Code.InvalidArgument:
    case Code.AlreadyExists:
    case Code.FailedPrecondition:
    case Code.NotFound:
      return sentence(error.rawMessage) || fallback(error.code);
    default:
      return fallback(error.code);
  }
}

/** Vrai quand l'erreur vient du réseau ou d'un serveur injoignable. */
export function isUnreachable(err: unknown): boolean {
  const { code } = ConnectError.from(err);
  // fetch qui échoue (serveur arrêté, réseau coupé) arrive en Unknown.
  return code === Code.Unavailable || code === Code.Unknown || code === Code.DeadlineExceeded;
}

function fallback(code: Code): string {
  switch (code) {
    case Code.Unavailable:
    case Code.Unknown:
    case Code.DeadlineExceeded:
      return 'Le serveur ForgeNet ne répond pas. Vérifiez qu’il est démarré et que votre connexion fonctionne.';
    case Code.PermissionDenied:
      return 'Votre rôle ne permet pas cette action.';
    case Code.Unauthenticated:
      return 'Votre session a pris fin. Reconnectez-vous.';
    default:
      return `Le serveur a rencontré une erreur (${Code[code]}). Réessayez dans un instant.`;
  }
}

function sentence(message: string): string {
  const text = message.trim();
  if (!text) {
    return '';
  }
  const capitalized = text.charAt(0).toUpperCase() + text.slice(1);
  return /[.!?]$/.test(capitalized) ? capitalized : `${capitalized}.`;
}
