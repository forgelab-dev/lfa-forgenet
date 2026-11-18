import { Code, ConnectError, createRouterTransport } from '@connectrpc/connect';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { memoryLocation } from 'wouter/memory-location';

import { App } from '../../app/App';
import { AuthService, UserRole, type ChangePasswordRequest } from '../../gen/forgenet/v1/auth_pb';
import { sessionInterceptors, type TransportFactory } from '../../lib/api/client';

/** Faux serveur : une session ouverte ; ChangePassword répond `answer`. */
function fakeServer(answer: () => { revokedSessions: number } = () => ({ revokedSessions: 2 })) {
  const calls: ChangePasswordRequest[] = [];
  const factory: TransportFactory = (session) =>
    createRouterTransport(
      ({ service }) => {
        service(AuthService, {
          getCurrentUser() {
            return {
              user: {
                id: 'u1',
                username: 'ana',
                email: 'ana@forgenet.test',
                role: UserRole.MEMBER,
              },
            };
          },
          logout() {
            return {};
          },
          changePassword(req) {
            calls.push(req);
            return answer();
          },
        });
      },
      { transport: { interceptors: sessionInterceptors(session) } }
    );
  return { factory, calls };
}

function renderAccount(server: ReturnType<typeof fakeServer>, path = '/compte') {
  sessionStorage.setItem(
    'forgenet.session',
    JSON.stringify({ token: 'jeton', expiresAt: Date.now() + 3_600_000 })
  );
  const location = memoryLocation({ path });
  render(<App transportFactory={server.factory} locationHook={location.hook} />);
}

function fill(label: string, value: string) {
  fireEvent.change(screen.getByLabelText(label), { target: { value } });
}

function fillValid() {
  fill('Mot de passe actuel', 'ancien-mot-de-passe');
  fill('Nouveau mot de passe', 'nouveau-mot-de-passe');
  fill('Confirmation du nouveau mot de passe', 'nouveau-mot-de-passe');
}

const submit = () =>
  fireEvent.click(screen.getByRole('button', { name: 'Changer le mot de passe' }));

describe('AccountPage', () => {
  beforeEach(() => sessionStorage.clear());

  it('montre l’identité du compte et s’ouvre depuis la carte du compte', async () => {
    renderAccount(fakeServer(), '/');
    const nav = await screen.findAllByRole('link', { name: /Mon compte : ana, Membre/ });
    fireEvent.click(nav[0]);

    expect(
      await screen.findByRole('heading', { level: 1, name: 'Mon compte' })
    ).toBeInTheDocument();
    const identity = screen.getByRole('region', { name: 'Identité' });
    expect(within(identity).getByText('ana@forgenet.test')).toBeInTheDocument();
    expect(within(identity).getByText('Membre')).toBeInTheDocument();
    expect(nav[0]).toHaveAttribute('aria-current', 'page');
  });

  it('signale les erreurs du formulaire sans appeler le serveur', async () => {
    const server = fakeServer();
    renderAccount(server);
    await screen.findByRole('heading', { level: 1, name: 'Mon compte' });

    fill('Nouveau mot de passe', 'court');
    submit();

    expect(screen.getByText('Saisissez votre mot de passe actuel.')).toBeInTheDocument();
    expect(screen.getByLabelText('Nouveau mot de passe')).toHaveAttribute('aria-invalid', 'true');
    expect(screen.getByText(/au moins 8 caractères/)).toBeInTheDocument();
    expect(server.calls).toHaveLength(0);
  });

  it('change le mot de passe, vide le formulaire et dit combien de sessions sont fermées', async () => {
    const server = fakeServer();
    renderAccount(server);
    await screen.findByRole('heading', { level: 1, name: 'Mon compte' });

    fillValid();
    submit();

    expect(await screen.findByText(/2 autres sessions ont été fermées/)).toBeInTheDocument();
    expect(server.calls).toEqual([
      expect.objectContaining({
        currentPassword: 'ancien-mot-de-passe',
        newPassword: 'nouveau-mot-de-passe',
      }),
    ]);
    expect(screen.getByLabelText('Mot de passe actuel')).toHaveValue('');
    expect(screen.getByLabelText('Nouveau mot de passe')).toHaveValue('');
  });

  it('montre le refus du serveur et garde la session ouverte', async () => {
    renderAccount(
      fakeServer(() => {
        throw new ConnectError('le mot de passe actuel est incorrect', Code.InvalidArgument);
      })
    );
    await screen.findByRole('heading', { level: 1, name: 'Mon compte' });

    fillValid();
    submit();

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Le mot de passe actuel est incorrect.'
    );
    expect(screen.getByRole('heading', { level: 1, name: 'Mon compte' })).toBeInTheDocument();
    expect(screen.getByLabelText('Nouveau mot de passe')).toHaveValue('nouveau-mot-de-passe');
  });

  it('désactive le bouton pendant l’envoi', async () => {
    let release: (v: { revokedSessions: number }) => void = () => {};
    const pending = new Promise<{ revokedSessions: number }>((resolve) => (release = resolve));
    const server = fakeServer();
    const factory: TransportFactory = (session) =>
      createRouterTransport(
        ({ service }) => {
          service(AuthService, {
            getCurrentUser: () => ({
              user: {
                id: 'u1',
                username: 'ana',
                email: 'ana@forgenet.test',
                role: UserRole.MEMBER,
              },
            }),
            logout: () => ({}),
            changePassword: () => pending,
          });
        },
        { transport: { interceptors: sessionInterceptors(session) } }
      );
    renderAccount({ ...server, factory });
    await screen.findByRole('heading', { level: 1, name: 'Mon compte' });

    fillValid();
    submit();

    expect(await screen.findByRole('button', { name: 'Changement…' })).toBeDisabled();
    release({ revokedSessions: 0 });
    await waitFor(() => expect(screen.getByText(/Aucune autre session/)).toBeInTheDocument());
  });

  it('affiche les mots de passe à la demande', async () => {
    renderAccount(fakeServer());
    await screen.findByRole('heading', { level: 1, name: 'Mon compte' });

    expect(screen.getByLabelText('Mot de passe actuel')).toHaveAttribute('type', 'password');
    fireEvent.click(screen.getByLabelText('Afficher les mots de passe'));
    expect(screen.getByLabelText('Mot de passe actuel')).toHaveAttribute('type', 'text');
  });
});
