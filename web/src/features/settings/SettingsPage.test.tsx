import { Code, ConnectError, createRouterTransport } from '@connectrpc/connect';
import { fireEvent, render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { memoryLocation } from 'wouter/memory-location';

import { App } from '../../app/App';
import { AuthService, UserRole } from '../../gen/forgenet/v1/auth_pb';
import { MailService, MailTLS } from '../../gen/forgenet/v1/mail_pb';
import { sessionInterceptors, type TransportFactory } from '../../lib/api/client';

interface Options {
  role?: UserRole;
  enabled?: boolean;
  sendError?: ConnectError;
}

function renderSettings({ role = UserRole.ADMIN, enabled = true, sendError }: Options = {}) {
  const calls = { send: 0 };
  const factory: TransportFactory = (session) =>
    createRouterTransport(
      ({ service }) => {
        service(AuthService, {
          getCurrentUser: () => ({
            user: { id: 'u1', username: 'ana', email: 'ana@forgenet.test', role },
          }),
          logout: () => ({}),
        });
        service(MailService, {
          getMailStatus: () =>
            enabled
              ? {
                  enabled: true,
                  host: 'smtp.exemple.org',
                  port: 587,
                  tls: MailTLS.MAIL_TLS_STARTTLS,
                  from: 'forgenet@exemple.org',
                  authenticated: true,
                }
              : {},
          sendTestMail: () => {
            calls.send++;
            if (sendError) {
              throw sendError;
            }
            return { to: 'ana@forgenet.test' };
          },
        });
      },
      { transport: { interceptors: sessionInterceptors(session) } }
    );
  sessionStorage.setItem(
    'forgenet.session',
    JSON.stringify({ token: 'jeton', expiresAt: Date.now() + 3_600_000 })
  );
  render(
    <App transportFactory={factory} locationHook={memoryLocation({ path: '/reglages' }).hook} />
  );
  return calls;
}

describe('SettingsPage', () => {
  beforeEach(() => sessionStorage.clear());

  it('montre le relais sans ses identifiants', async () => {
    renderSettings();
    expect(await screen.findByText('smtp.exemple.org:587')).toBeInTheDocument();
    expect(screen.getByText('STARTTLS, exigé')).toBeInTheDocument();
    expect(screen.getByText('oui, gardés sur le serveur')).toBeInTheDocument();
  });

  it('envoie un e-mail d’essai et dit à qui', async () => {
    const calls = renderSettings();
    fireEvent.click(await screen.findByRole('button', { name: 'Envoyer un e-mail d’essai' }));
    expect(await screen.findByText(/Le relais a pris l’e-mail d’essai/)).toBeInTheDocument();
    expect(screen.getByText('ana@forgenet.test')).toBeInTheDocument();
    expect(calls.send).toBe(1);
  });

  it('montre la réponse du relais quand il refuse', async () => {
    renderSettings({
      sendError: new ConnectError(
        'le relais n’a pas pris le message : authentification refusée par le relais : 535 « identifiants refusés »',
        Code.FailedPrecondition
      ),
    });
    fireEvent.click(await screen.findByRole('button', { name: 'Envoyer un e-mail d’essai' }));
    expect(await screen.findByRole('alert')).toHaveTextContent('535');
  });

  it('explique comment activer l’envoi quand aucun relais n’est réglé', async () => {
    renderSettings({ enabled: false });
    expect(await screen.findByText(/Aucun relais SMTP n’est configuré/)).toBeInTheDocument();
    expect(screen.getByText(/--smtp-host smtp.exemple.org/)).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Envoyer un e-mail d’essai' })).toBeNull();
  });

  it('reste fermée aux membres, sans appel au serveur', async () => {
    const calls = renderSettings({ role: UserRole.MEMBER });
    expect(
      await screen.findByText('Les réglages du serveur sont réservés aux administrateurs.')
    ).toBeInTheDocument();
    expect(calls.send).toBe(0);
  });
});
