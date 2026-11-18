import { describe, expect, it } from 'vitest';

import { checkPasswordForm, successMessage } from './password';

const valid = {
  current: 'ancien-mot-de-passe',
  next: 'nouveau-mot-de-passe',
  confirm: 'nouveau-mot-de-passe',
};

describe('checkPasswordForm', () => {
  it('accepte un formulaire complet', () => {
    expect(checkPasswordForm(valid)).toEqual({});
  });

  it('exige le mot de passe actuel', () => {
    expect(checkPasswordForm({ ...valid, current: '' }).current).toMatch(/actuel/);
  });

  it('compte le minimum en caractères, le maximum en octets', () => {
    expect(checkPasswordForm({ ...valid, next: 'court', confirm: 'court' }).next).toMatch(
      /au moins 8/
    );
    expect(checkPasswordForm({ ...valid, next: 'éééééééé', confirm: 'éééééééé' })).toEqual({});
    const long = 'é'.repeat(37);
    expect(checkPasswordForm({ ...valid, next: long, confirm: long }).next).toMatch(
      /au plus 72 octets/
    );
  });

  it('refuse le même mot de passe et une confirmation différente', () => {
    expect(
      checkPasswordForm({ ...valid, next: valid.current, confirm: valid.current }).next
    ).toMatch(/différent/);
    expect(checkPasswordForm({ ...valid, confirm: 'autre-chose' }).confirm).toMatch(/diffèrent/);
  });
});

describe('successMessage', () => {
  it('dit combien de sessions ont été fermées', () => {
    expect(successMessage(0)).toMatch(/Aucune autre session/);
    expect(successMessage(1)).toMatch(/1 autre session a été fermée/);
    expect(successMessage(3)).toMatch(/3 autres sessions ont été fermées/);
  });
});
