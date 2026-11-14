import { create, type MessageInitShape } from '@bufbuild/protobuf';
import { timestampFromMs } from '@bufbuild/protobuf/wkt';
import { Code, ConnectError, createRouterTransport } from '@connectrpc/connect';
import { act, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { memoryLocation } from 'wouter/memory-location';

import { App } from '../../app/App';
import { AuthService, UserRole } from '../../gen/forgenet/v1/auth_pb';
import {
  GetUpdateStatusResponseSchema,
  UpdateResult,
  UpdateService,
  UpdateStep,
} from '../../gen/forgenet/v1/update_pb';
import { sessionInterceptors, type TransportFactory } from '../../lib/api/client';

type StatusInit = MessageInitShape<typeof GetUpdateStatusResponseSchema>;

interface Server {
  factory: TransportFactory;
  calls: { status: number; check: number; start: number };
}

/**
 * Faux serveur : `statuses` donne les réponses successives de GetUpdateStatus
 * (la dernière se répète) ; une réponse `'down'` simule le serveur arrêté.
 */
function fakeServer({
  role = UserRole.ADMIN,
  statuses,
  checked,
  startError,
}: {
  role?: UserRole;
  statuses: (StatusInit | 'down')[];
  checked?: StatusInit;
  startError?: ConnectError;
}): Server {
  const calls = { status: 0, check: 0, start: 0 };
  const status = (init: StatusInit) =>
    create(GetUpdateStatusResponseSchema, {
      currentVersion: '0.2.0',
      checkEnabled: true,
      supported: true,
      ...init,
    });
  const factory: TransportFactory = (session) =>
    createRouterTransport(
      ({ service }) => {
        service(AuthService, {
          getCurrentUser() {
            return { user: { id: 'u1', username: 'mike', email: 'mike@forgenet.test', role } };
          },
        });
        service(UpdateService, {
          getUpdateStatus() {
            const next = statuses[Math.min(calls.status, statuses.length - 1)];
            calls.status++;
            if (next === 'down') {
              throw new ConnectError('injoignable', Code.Unavailable);
            }
            return status(next);
          },
          checkForUpdates() {
            calls.check++;
            return { status: status(checked ?? {}) };
          },
          startUpdate() {
            calls.start++;
            if (startError) {
              throw startError;
            }
            return { toVersion: '0.3.0' };
          },
        });
      },
      { transport: { interceptors: sessionInterceptors(session) } }
    );
  return { factory, calls };
}

function renderApp(server: Server) {
  sessionStorage.setItem(
    'forgenet.session',
    JSON.stringify({ token: 'jeton', expiresAt: Date.now() + 3_600_000 })
  );
  // Page à venir, sans appel à l'API : seuls les avis du bandeau s'affichent.
  const location = memoryLocation({ path: '/equipe' });
  render(<App transportFactory={server.factory} locationHook={location.hook} />);
}

const available = {
  available: {
    version: '0.3.0',
    notesUrl: 'https://github.com/forgelab-dev/lfa-forgenet/releases/tag/v0.3.0',
  },
};

const step = (label: string) => screen.getByText(label).closest('li');

const minutesAgo = (n: number) => timestampFromMs(Date.now() - n * 60_000);

describe('UpdateBanner', () => {
  beforeEach(() => {
    sessionStorage.clear();
    localStorage.clear();
  });
  afterEach(() => vi.useRealTimers());

  it('annonce la version disponible avec ses notes', async () => {
    renderApp(fakeServer({ statuses: [available] }));

    expect(
      await screen.findByRole('heading', { name: 'ForgeNet 0.3.0 est disponible' })
    ).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /Lire les notes de version/ })).toHaveAttribute(
      'href',
      'https://github.com/forgelab-dev/lfa-forgenet/releases/tag/v0.3.0'
    );
    expect(screen.getByText(/Ce serveur tourne en/)).toHaveTextContent(
      'Ce serveur tourne en 0.2.0.'
    );
    expect(screen.getByText((_, el) => el?.textContent === '0.3.0 disponible')).toBeInTheDocument();
  });

  it('confirme, lance, suit la coupure, puis annonce la réussite', async () => {
    vi.useFakeTimers({ shouldAdvanceTime: true });
    const run = { fromVersion: '0.2.0', toVersion: '0.3.0', startedAt: minutesAgo(0) };
    const server = fakeServer({
      statuses: [
        available,
        { ...available, lastRun: { ...run, step: UpdateStep.PULL, result: UpdateResult.RUNNING } },
        'down',
        {
          currentVersion: '0.3.0',
          lastRun: {
            ...run,
            step: UpdateStep.DONE,
            result: UpdateResult.SUCCEEDED,
            finishedAt: minutesAgo(0),
          },
        },
      ],
    });
    renderApp(server);

    fireEvent.click(await screen.findByRole('button', { name: 'Mettre à jour' }));
    const confirm = screen.getByRole('group', { name: 'Mettre à jour vers 0.3.0' });
    expect(within(confirm).getByText(/sauvegardées avant tout changement/)).toBeInTheDocument();
    expect(server.calls.start).toBe(0);

    fireEvent.click(within(confirm).getByRole('button', { name: 'Installer la version 0.3.0' }));
    expect(
      await screen.findByRole('heading', { name: 'Mise à jour vers 0.3.0' })
    ).toBeInTheDocument();
    expect(server.calls.start).toBe(1);

    await act(() => vi.advanceTimersByTimeAsync(2100));
    expect(step('Téléchargement de la nouvelle version')).toHaveTextContent('en cours');

    await act(() => vi.advanceTimersByTimeAsync(2100));
    expect(screen.getByText(/Le serveur redémarre/)).toBeInTheDocument();
    expect(step('Redémarrage sur la nouvelle version')).toHaveTextContent('en cours');
    expect(step('Sauvegarde de la base et de la configuration')).toHaveTextContent('terminé');

    await act(() => vi.advanceTimersByTimeAsync(2100));
    expect(await screen.findByText(/est installé/)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Recharger le dashboard' })).toBeInTheDocument();
  });

  it('montre un refus du serveur dans la confirmation', async () => {
    renderApp(
      fakeServer({
        statuses: [available],
        startError: new ConnectError('une mise à jour est déjà en cours', Code.Aborted),
      })
    );

    fireEvent.click(await screen.findByRole('button', { name: 'Mettre à jour' }));
    fireEvent.click(screen.getByRole('button', { name: 'Installer la version 0.3.0' }));

    const alert = await screen.findByRole('alert');
    expect(alert).toHaveTextContent('Une mise à jour est déjà en cours.');
    expect(screen.getByRole('button', { name: 'Installer la version 0.3.0' })).toBeEnabled();
  });

  it('renvoie au serveur une installation non standard', async () => {
    renderApp(fakeServer({ statuses: [{ ...available, supported: false }] }));

    expect(await screen.findByText(/n’a pas été faite par/)).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Mettre à jour' })).not.toBeInTheDocument();
  });

  it('signale un retour arrière jusqu’à ce qu’on ferme l’avis', async () => {
    renderApp(
      fakeServer({
        statuses: [
          {
            lastRun: {
              fromVersion: '0.2.0',
              toVersion: '0.3.0',
              step: UpdateStep.ROLLBACK,
              result: UpdateResult.ROLLED_BACK,
              reason: 'la version 0.3.0 ne démarre pas',
              backup: '/opt/forgenet/sauvegardes/20261114-150000-0.2.0',
              startedAt: minutesAgo(10),
              finishedAt: minutesAgo(8),
            },
          },
        ],
      })
    );

    const alert = await screen.findByRole('alert');
    expect(alert).toHaveTextContent('La mise à jour vers 0.3.0 a échoué');
    expect(alert).toHaveTextContent('ForgeNet 0.2.0 a été remis en service, avec ses données.');
    expect(alert).toHaveTextContent('Raison : la version 0.3.0 ne démarre pas');

    fireEvent.click(within(alert).getByRole('button', { name: 'Fermer cet avis' }));
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    expect(localStorage.getItem('forgenet.updates.dismissed')).not.toBeNull();
  });

  it('vérifie sur demande depuis la barre latérale', async () => {
    const server = fakeServer({ statuses: [{}], checked: available });
    renderApp(server);

    fireEvent.click(await screen.findByRole('button', { name: 'Vérifier les mises à jour' }));

    expect(
      await screen.findByRole('heading', { name: 'ForgeNet 0.3.0 est disponible' })
    ).toBeInTheDocument();
    expect(server.calls.check).toBe(1);
  });

  it('ne montre rien et n’appelle rien pour un membre', async () => {
    const server = fakeServer({ role: UserRole.MEMBER, statuses: [available] });
    renderApp(server);

    expect(await screen.findByText('Membre')).toBeInTheDocument();
    expect(screen.queryByRole('heading', { name: /est disponible/ })).not.toBeInTheDocument();
    expect(server.calls.status).toBe(0);
  });
});
