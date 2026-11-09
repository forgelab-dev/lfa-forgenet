import { Code, ConnectError } from '@connectrpc/connect';
import { describe, expect, it } from 'vitest';

import { describeError, isUnreachable } from './errors';

describe('describeError', () => {
  it('reprend les refus rédigés par le serveur, en phrase', () => {
    const err = new ConnectError('identifiant ou mot de passe incorrect', Code.Unauthenticated);

    expect(describeError(err)).toBe('Identifiant ou mot de passe incorrect.');
  });

  it('garde la ponctuation finale du serveur', () => {
    const err = new ConnectError('trop d’essais, réessayez plus tard.', Code.ResourceExhausted);

    expect(describeError(err)).toBe('Trop d’essais, réessayez plus tard.');
  });

  it('ne montre pas le détail des erreurs internes', () => {
    const err = new ConnectError('pq: relation "users" does not exist', Code.Internal);

    expect(describeError(err)).toBe(
      'Le serveur a rencontré une erreur (Internal). Réessayez dans un instant.'
    );
  });

  it('explique un serveur injoignable, même pour une panne réseau', () => {
    const network = new TypeError('Failed to fetch');

    expect(describeError(network)).toContain('ne répond pas');
    expect(isUnreachable(network)).toBe(true);
    expect(isUnreachable(new ConnectError('refus', Code.PermissionDenied))).toBe(false);
  });

  it('traduit un refus de droits', () => {
    expect(describeError(new ConnectError('denied', Code.PermissionDenied))).toBe(
      'Votre rôle ne permet pas cette action.'
    );
  });
});
