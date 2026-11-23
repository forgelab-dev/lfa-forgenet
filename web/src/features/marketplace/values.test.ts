import { create } from '@bufbuild/protobuf';
import { describe, expect, it } from 'vitest';

import { TemplateSchema, TemplateVariableSchema } from '../../gen/forgenet/v1/marketplace_pb';
import { checkValue, emptyHint, initialValues, introduction, refusal } from './values';

const password = create(TemplateVariableSchema, {
  name: 'DB_PASSWORD',
  label: 'Mot de passe',
  secret: true,
  generate: 24,
});
const user = create(TemplateVariableSchema, {
  name: 'DB_USER',
  label: 'Utilisateur',
  defaultValue: 'app',
  pattern: '^[a-z_]+$',
});
const token = create(TemplateVariableSchema, { name: 'API_TOKEN', label: 'Jeton', secret: true });
const note = create(TemplateVariableSchema, { name: 'NOTE', label: 'Note', optional: true });
const variables = [password, user, token, note];

describe('values', () => {
  it('part des valeurs par défaut du manifeste', () => {
    expect(initialValues(variables)).toEqual({
      DB_PASSWORD: '',
      DB_USER: 'app',
      API_TOKEN: '',
      NOTE: '',
    });
  });

  it('accepte un champ vide quand le serveur sait le remplir', () => {
    expect(checkValue(password, '')).toBeUndefined();
    expect(checkValue(user, '')).toBeUndefined();
    expect(checkValue(note, '')).toBeUndefined();
    expect(checkValue(token, '')).toBe('Saisissez une valeur.');
  });

  it('vérifie le motif, la taille et le caractère nul comme le serveur', () => {
    expect(checkValue(user, 'app_web')).toBeUndefined();
    expect(checkValue(user, 'App')).toBe('La valeur ne respecte pas le format attendu : ^[a-z_]+$');
    expect(checkValue(token, 'a'.repeat(4097))).toBe('4 096 octets au plus.');
    expect(checkValue(token, 'é'.repeat(2049))).toBe('4 096 octets au plus.');
    expect(checkValue(token, 'a\u0000b')).toBe('Le caractère nul est interdit.');
  });

  it('laisse au serveur un motif que le navigateur ne sait pas lire', () => {
    const go = create(TemplateVariableSchema, { name: 'X', label: 'X', pattern: '(?i)^ok$' });
    expect(checkValue(go, 'OK')).toBeUndefined();
  });

  it('dit ce que devient un champ laissé vide', () => {
    expect(emptyHint(password)).toBe(
      'Laissez vide : un mot de passe de 24 caractères sera généré.'
    );
    expect(emptyHint(user)).toBe('Laissé vide, il vaut « app ».');
    expect(emptyHint(note)).toBe('Facultatif.');
    expect(emptyHint(token)).toBeUndefined();
  });

  it('range un refus du serveur sous le champ qu’il nomme', () => {
    expect(
      refusal(
        "Valeur d'installation refusée : DB_USER : la valeur ne respecte pas le motif ^[a-z_]+$.",
        variables
      )
    ).toEqual({ field: 'DB_USER', message: 'La valeur ne respecte pas le motif ^[a-z_]+$.' });
    expect(refusal('Une application porte déjà ce nom.', variables)).toEqual({
      message: 'Une application porte déjà ce nom.',
    });
    expect(
      refusal(
        "Valeur d'installation refusée : AUTRE : variable inconnue de ce template.",
        variables
      )
    ).toEqual({
      message: "Valeur d'installation refusée : AUTRE : variable inconnue de ce template.",
    });
  });

  it('ne parle de mots de passe générés que si le template en tire', () => {
    expect(introduction(create(TemplateSchema, { variables }))).toMatch(/montré une seule fois/);
    expect(introduction(create(TemplateSchema, { variables: [user] }))).toBe(
      'Les valeurs sont chiffrées sur le serveur et ne se réaffichent pas.'
    );
    expect(introduction(create(TemplateSchema, {}))).toMatch(/aucune valeur/);
  });
});
