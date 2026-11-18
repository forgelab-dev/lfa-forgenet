import { createRouterTransport } from '@connectrpc/connect';
import { render, screen, within } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { memoryLocation } from 'wouter/memory-location';

import { App } from '../../app/App';
import { AuthService, UserRole } from '../../gen/forgenet/v1/auth_pb';
import { ContainerService, ContainerState } from '../../gen/forgenet/v1/container_pb';
import { sessionInterceptors, type TransportFactory } from '../../lib/api/client';

type Row = { name: string; state: ContainerState; system?: boolean };

function renderAt(path: string, containers: Row[]) {
  const factory: TransportFactory = (session) =>
    createRouterTransport(
      ({ service }) => {
        service(AuthService, {
          getCurrentUser: () => ({
            user: { id: 'u1', username: 'ana', email: 'ana@forgenet.test', role: UserRole.MEMBER },
          }),
          logout: () => ({}),
        });
        service(ContainerService, {
          listContainers: () => ({
            containers: containers.map((c) => ({ id: c.name, image: 'img', ...c })),
          }),
        });
      },
      { transport: { interceptors: sessionInterceptors(session) } }
    );
  sessionStorage.setItem(
    'forgenet.session',
    JSON.stringify({ token: 'jeton', expiresAt: Date.now() + 3_600_000 })
  );
  render(<App transportFactory={factory} locationHook={memoryLocation({ path }).hook} />);
}

// L'état du ticket #74 : deux applications, et ForgeNet après une mise à jour.
const afterUpdate: Row[] = [
  { name: 'boutique-web-1', state: ContainerState.RUNNING },
  { name: 'boutique-db-1', state: ContainerState.RUNNING },
  { name: 'forgenet-server', state: ContainerState.RUNNING, system: true },
  { name: 'forgenet-caddy', state: ContainerState.RUNNING, system: true },
  { name: 'forgenet-updater', state: ContainerState.EXITED, system: true },
];

describe('conteneurs de ForgeNet', () => {
  beforeEach(() => sessionStorage.clear());

  it('ne comptent pas comme des applications dans la vue d’ensemble', async () => {
    renderAt('/', afterUpdate);

    expect(await screen.findByText(/conteneurs d’applications/)).toHaveTextContent(
      '2 conteneurs d’applications'
    );
    expect(screen.getByText('2 en marche')).toBeInTheDocument();
    expect(screen.queryByText(/à l’arrêt/)).toBeNull();
    expect(screen.getByText(/Et 3 conteneurs de ForgeNet lui-même/)).toBeInTheDocument();
  });

  it('signalent Caddy arrêté dans la vue d’ensemble', async () => {
    renderAt('/', [
      ...afterUpdate.slice(0, 3),
      { name: 'forgenet-caddy', state: ContainerState.EXITED, system: true },
    ]);

    expect(
      await screen.findByText(/forgenet-caddy est à l’arrêt : le serveur n’est pas au complet/)
    ).toBeInTheDocument();
  });

  it('sont listés à part sur la page Conteneurs', async () => {
    renderAt('/conteneurs', afterUpdate);

    const apps = await screen.findByRole('table', { name: 'Conteneurs des applications' });
    const system = screen.getByRole('table', { name: 'Conteneurs de ForgeNet' });
    expect(within(apps).queryByText('forgenet-updater')).toBeNull();
    expect(within(system).getByText('forgenet-updater')).toBeInTheDocument();
    expect(within(system).getAllByRole('row')).toHaveLength(4);
  });
});
