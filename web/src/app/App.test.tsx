import { timestampFromMs } from '@bufbuild/protobuf/wkt';
import { Code, ConnectError, createRouterTransport } from '@connectrpc/connect';
import { fireEvent, render, screen, within } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { memoryLocation } from 'wouter/memory-location';

import { AuthService, UserRole } from '../gen/forgenet/v1/auth_pb';
import { ContainerService, ContainerState } from '../gen/forgenet/v1/container_pb';
import { sessionInterceptors, type TransportFactory } from '../lib/api/client';
import { App } from './App';

const PASSWORD = 'mot-de-passe-de-test';
const mike = { id: 'u1', username: 'mike', email: 'mike@forgenet.test', role: UserRole.ADMIN };

/** Faux serveur ForgeNet : un compte, des jetons révocables, trois conteneurs. */
function fakeServer(options: { unavailable?: boolean } = {}) {
  const tokens = new Set<string>();
  const calls = { logout: 0 };
  let issued = 0;

  const factory: TransportFactory = (session) =>
    createRouterTransport(
      ({ service }) => {
        const userOf = (header: string | null) => {
          if (options.unavailable) {
            throw new ConnectError('injoignable', Code.Unavailable);
          }
          const token = header?.replace(/^Bearer /, '');
          if (!token || !tokens.has(token)) {
            throw new ConnectError('authentification requise', Code.Unauthenticated);
          }
          return token;
        };
        service(AuthService, {
          login(req) {
            if (req.login !== 'mike' || req.password !== PASSWORD) {
              throw new ConnectError('identifiant ou mot de passe incorrect', Code.Unauthenticated);
            }
            const token = `jeton-${++issued}`;
            tokens.add(token);
            return { token, expiresAt: timestampFromMs(Date.now() + 3_600_000), user: mike };
          },
          getCurrentUser(_req, ctx) {
            userOf(ctx.requestHeader.get('Authorization'));
            return { user: mike };
          },
          logout(_req, ctx) {
            calls.logout++;
            tokens.delete(userOf(ctx.requestHeader.get('Authorization')));
            return {};
          },
        });
        service(ContainerService, {
          listContainers(_req, ctx) {
            userOf(ctx.requestHeader.get('Authorization'));
            return {
              containers: [
                { id: 'c1', name: 'postgres', image: 'postgres:17', state: ContainerState.EXITED },
                {
                  id: 'c2',
                  name: 'caddy',
                  image: 'caddy:2',
                  state: ContainerState.RUNNING,
                  ports: [{ hostPort: 443, containerPort: 443, protocol: 'tcp' }],
                },
                { id: 'c3', name: 'api', image: 'forgenet/api:1', state: ContainerState.RUNNING },
              ],
            };
          },
        });
      },
      { transport: { interceptors: sessionInterceptors(session) } }
    );

  return { factory, tokens, calls };
}

function renderApp(server: ReturnType<typeof fakeServer>, path = '/') {
  const location = memoryLocation({ path, record: true });
  render(<App transportFactory={server.factory} locationHook={location.hook} />);
  return location;
}

async function signIn(password = PASSWORD) {
  fireEvent.change(await screen.findByLabelText('Identifiant'), { target: { value: 'mike' } });
  fireEvent.change(screen.getByLabelText('Mot de passe'), { target: { value: password } });
  fireEvent.click(screen.getByRole('button', { name: 'Se connecter' }));
}

function storeSession(token: string) {
  sessionStorage.setItem(
    'forgenet.session',
    JSON.stringify({ token, expiresAt: Date.now() + 3_600_000 })
  );
}

describe('App', () => {
  beforeEach(() => sessionStorage.clear());

  it('demande de se connecter quand aucune session n’est ouverte', async () => {
    renderApp(fakeServer());

    expect(await screen.findByRole('heading', { level: 1, name: 'Connexion' })).toBeInTheDocument();
    expect(screen.queryByRole('navigation', { name: 'Navigation principale' })).toBeNull();
  });

  it('signale les champs vides sans appeler le serveur', async () => {
    renderApp(fakeServer());

    fireEvent.click(await screen.findByRole('button', { name: 'Se connecter' }));

    expect(screen.getByText('Saisissez votre identifiant.')).toBeInTheDocument();
    expect(screen.getByText('Saisissez votre mot de passe.')).toBeInTheDocument();
    expect(screen.getByLabelText('Identifiant')).toHaveAttribute('aria-invalid', 'true');
  });

  it('affiche le refus du serveur et vide le mot de passe', async () => {
    renderApp(fakeServer());

    await signIn('mauvais');

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Identifiant ou mot de passe incorrect.'
    );
    expect(screen.getByLabelText('Mot de passe')).toHaveValue('');
  });

  it('ouvre la session, garde le jeton et affiche la vue d’ensemble', async () => {
    const server = fakeServer();
    renderApp(server);

    await signIn();

    expect(
      await screen.findByRole('heading', { level: 1, name: 'Vue d’ensemble' })
    ).toBeInTheDocument();
    expect(await screen.findByText('2 en marche')).toBeInTheDocument();
    expect(screen.getByText('1 à l’arrêt')).toBeInTheDocument();
    expect(screen.getAllByText('Administrateur')[0]).toBeInTheDocument();
    expect(JSON.parse(sessionStorage.getItem('forgenet.session') ?? '{}').token).toBe('jeton-1');
  });

  it('reprend une session conservée et navigue vers les conteneurs', async () => {
    const server = fakeServer();
    server.tokens.add('jeton-conserve');
    storeSession('jeton-conserve');
    const location = renderApp(server);

    const nav = await screen.findByRole('navigation', { name: 'Navigation principale' });
    expect(within(nav).getByRole('link', { name: /Vue d’ensemble/ })).toHaveAttribute(
      'aria-current',
      'page'
    );
    fireEvent.click(within(nav).getByRole('link', { name: /Conteneurs/ }));

    expect(location.history.at(-1)).toBe('/conteneurs');
    const rows = await screen.findAllByRole('row');
    // En-tête, puis les conteneurs en marche d'abord, par nom.
    expect(rows.slice(1).map((row) => within(row).getAllByRole('cell')[0].textContent)).toEqual([
      'apiforgenet/api:1',
      'caddycaddy:2',
      'postgrespostgres:17',
    ]);
    expect(screen.getByText('443 → 443/tcp')).toBeInTheDocument();
  });

  it('présente une section à venir sans faux contenu', async () => {
    const server = fakeServer();
    server.tokens.add('jeton-conserve');
    storeSession('jeton-conserve');
    renderApp(server, '/marketplace');

    expect(
      await screen.findByRole('heading', { level: 1, name: 'Marketplace' })
    ).toBeInTheDocument();
    expect(screen.getByText('En cours de construction')).toBeInTheDocument();
  });

  it('répond par une page introuvable aux adresses inconnues', async () => {
    const server = fakeServer();
    server.tokens.add('jeton-conserve');
    storeSession('jeton-conserve');
    renderApp(server, '/nulle-part');

    expect(
      await screen.findByRole('heading', { level: 1, name: 'Page introuvable' })
    ).toBeInTheDocument();
  });

  it('ferme la session côté serveur et dans le navigateur', async () => {
    const server = fakeServer();
    server.tokens.add('jeton-conserve');
    storeSession('jeton-conserve');
    renderApp(server);

    fireEvent.click(await screen.findByRole('button', { name: 'Se déconnecter' }));

    expect(await screen.findByRole('heading', { level: 1, name: 'Connexion' })).toBeInTheDocument();
    expect(server.calls.logout).toBe(1);
    expect(server.tokens.has('jeton-conserve')).toBe(false);
    expect(sessionStorage.getItem('forgenet.session')).toBeNull();
    expect(screen.queryByRole('status')).toBeNull();
  });

  it('revient à la connexion quand le serveur refuse le jeton conservé', async () => {
    storeSession('jeton-revoque');
    renderApp(fakeServer());

    expect(await screen.findByText(/Votre session a pris fin/)).toHaveAttribute('role', 'status');
    expect(screen.getByRole('heading', { level: 1, name: 'Connexion' })).toBeInTheDocument();
    expect(sessionStorage.getItem('forgenet.session')).toBeNull();
  });

  it('propose de réessayer quand le serveur ne répond pas', async () => {
    storeSession('jeton-conserve');
    renderApp(fakeServer({ unavailable: true }));

    expect(await screen.findByRole('alert')).toHaveTextContent('ne répond pas');
    expect(screen.getByRole('button', { name: 'Réessayer' })).toBeInTheDocument();
    expect(sessionStorage.getItem('forgenet.session')).not.toBeNull();
  });
});
