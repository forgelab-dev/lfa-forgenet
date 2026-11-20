import { create } from '@bufbuild/protobuf';
import { describe, expect, it } from 'vitest';

import { GetMailStatusResponseSchema, MailTLS } from '../../gen/forgenet/v1/mail_pb';
import { canKeepPassword, checkRelayForm, formFromStatus, requestFromForm } from './relay';

const filled = {
  host: 'smtp.exemple.org',
  port: '587',
  tls: MailTLS.MAIL_TLS_STARTTLS,
  username: 'forgenet',
  password: '',
  from: 'ForgeNet <forgenet@exemple.org>',
};

describe('relay', () => {
  it('part des réglages en vigueur, jamais du mot de passe', () => {
    const status = create(GetMailStatusResponseSchema, {
      host: 'smtp.exemple.org',
      port: 465,
      tls: MailTLS.MAIL_TLS_IMPLICIT,
      username: 'forgenet',
      from: 'forgenet@exemple.org',
      authenticated: true,
    });
    expect(formFromStatus(status)).toMatchObject({ port: '465', password: '' });
    expect(canKeepPassword(status)).toBe(true);
    expect(canKeepPassword({ ...status, passwordUnreadable: true })).toBe(false);
    expect(formFromStatus(create(GetMailStatusResponseSchema, {}))).toMatchObject({
      host: '',
      port: '587',
      tls: MailTLS.MAIL_TLS_STARTTLS,
    });
  });

  it('accepte un formulaire complet, mot de passe gardé', () => {
    expect(checkRelayForm(filled, true)).toEqual({});
    expect(requestFromForm(filled, true)).toMatchObject({ keepPassword: true, password: '' });
  });

  it('refuse ce que le serveur refuserait', () => {
    expect(checkRelayForm({ ...filled, host: 'user@smtp.exemple.org' }, true).host).toBeDefined();
    expect(checkRelayForm({ ...filled, port: '70000' }, true).port).toBeDefined();
    expect(checkRelayForm({ ...filled, port: '58a' }, true).port).toBeDefined();
    expect(checkRelayForm({ ...filled, tls: MailTLS.MAIL_TLS_NONE }, true).tls).toBeDefined();
    expect(
      checkRelayForm({ ...filled, host: 'localhost', tls: MailTLS.MAIL_TLS_NONE }, true).tls
    ).toBeUndefined();
    expect(checkRelayForm(filled, false).password).toBe('Saisissez le mot de passe du relais.');
    expect(checkRelayForm({ ...filled, password: 'a\nb' }, true).password).toBeDefined();
    expect(checkRelayForm({ ...filled, from: 'pas une adresse' }, true).from).toBeDefined();
  });

  it('sans identifiant, n’envoie aucun mot de passe', () => {
    const req = requestFromForm({ ...filled, username: ' ', password: 'oublié' }, true);
    expect(req).toMatchObject({ username: '', keepPassword: false, password: '' });
    expect(checkRelayForm({ ...filled, username: '' }, false)).toEqual({});
  });
});
