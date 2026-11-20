import { describe, expect, it } from 'vitest';

import { checkEmailForm } from './email';

describe('checkEmailForm', () => {
  it('accepte une adresse nouvelle avec le mot de passe', () => {
    expect(
      checkEmailForm({ email: ' Ana@Nouveau.org ', password: 'x' }, 'ana@forgenet.test')
    ).toEqual({});
  });

  it('demande les deux champs', () => {
    expect(checkEmailForm({ email: '', password: '' }, 'ana@forgenet.test')).toEqual({
      email: 'Saisissez la nouvelle adresse.',
      password: 'Saisissez votre mot de passe pour confirmer.',
    });
  });

  it('refuse une adresse mal formée, trop longue ou inchangée', () => {
    for (const email of ['ana', 'ana@', 'ana nouveau@x.org', `${'a'.repeat(250)}@x.org`]) {
      expect(checkEmailForm({ email, password: 'x' }, 'ana@forgenet.test').email).toMatch(
        /n’est pas valide/
      );
    }
    expect(
      checkEmailForm({ email: 'ANA@forgenet.test', password: 'x' }, 'ana@forgenet.test').email
    ).toBe('C’est déjà l’adresse de ce compte.');
  });
});
