import { timestampFromMs } from '@bufbuild/protobuf/wkt';
import { Code, ConnectError, createRouterTransport } from '@connectrpc/connect';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { memoryLocation } from 'wouter/memory-location';

import { App } from '../../app/App';
import { AuthService, UserRole } from '../../gen/forgenet/v1/auth_pb';
import { DomainService, type AddDomainRequest } from '../../gen/forgenet/v1/domain_pb';
import { sessionInterceptors, type TransportFactory } from '../../lib/api/client';

interface Options {
  role?: UserRole;
  disabled?: boolean;
  /** Erreur renvoyée par AddDomain, à la place d'un ajout. */
  addError?: ConnectError;
}

const created = timestampFromMs(Date.UTC(2026, 10, 13, 12));

/** Faux serveur : une session ouverte et deux domaines. */
function fakeServer({ role = UserRole.ADMIN, disabled = false, addError }: Options = {}) {
  const domains = [
    { id: 'dom_2', host: 'www.exemple.com', projectId: 'vitrine', service: 'web', port: 80 },
    { id: 'dom_1', host: 'api.exemple.com', projectId: 'boutique', service: 'api', port: 8080 },
  ].map((d) => ({ ...d, createdBy: 'user:u1', createdAt: created }));
  const calls: { added: AddDomainRequest[]; removed: string[] } = { added: [], removed: [] };

  const factory: TransportFactory = (session) =>
    createRouterTransport(
      ({ service }) => {
        service(AuthService, {
          getCurrentUser() {
            return { user: { id: 'u1', username: 'mike', email: 'mike@forgenet.test', role } };
          },
          logout() {
            return {};
          },
        });
        service(DomainService, {
          listDomains() {
            if (disabled) {
              throw new ConnectError(
                'reverse proxy désactivé : activez proxy.enabled dans la configuration du serveur',
                Code.FailedPrecondition
              );
            }
            return { domains };
          },
          addDomain(req) {
            calls.added.push(req);
            if (addError) {
              throw addError;
            }
            const domain = {
              id: 'dom_3',
              host: req.host.toLowerCase(),
              projectId: req.projectId,
              service: req.service,
              port: req.port,
              createdBy: 'user:u1',
              createdAt: created,
            };
            domains.push(domain);
            return { domain };
          },
          removeDomain(req) {
            calls.removed.push(req.domainId);
            return {};
          },
        });
      },
      { transport: { interceptors: sessionInterceptors(session) } }
    );
  return { factory, calls };
}

function renderDomains(server: ReturnType<typeof fakeServer>) {
  sessionStorage.setItem(
    'forgenet.session',
    JSON.stringify({ token: 'jeton', expiresAt: Date.now() + 3_600_000 })
  );
  const location = memoryLocation({ path: '/domaines' });
  render(<App transportFactory={server.factory} locationHook={location.hook} />);
}

async function domainRows() {
  const list = await screen.findByRole('list', { name: 'Domaines servis' });
  return within(list)
    .getAllByRole('listitem')
    .map((item) => within(item).getByRole('link').textContent);
}

function fill(label: string, value: string) {
  fireEvent.change(screen.getByLabelText(label), { target: { value } });
}

describe('DomainsPage', () => {
  beforeEach(() => sessionStorage.clear());

  it('liste les domaines par nom, avec leur cible', async () => {
    renderDomains(fakeServer());

    expect(await screen.findByRole('heading', { level: 1, name: 'Domaines' })).toBeInTheDocument();
    expect(await domainRows()).toEqual([
      'api.exemple.com(ouvre un nouvel onglet)',
      'www.exemple.com(ouvre un nouvel onglet)',
    ]);
    expect(screen.getByText('boutique / api:8080')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: '2 domaines servis' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /api\.exemple\.com/ })).toHaveAttribute(
      'href',
      'https://api.exemple.com'
    );
  });

  it('signale les champs manquants sans appeler le serveur', async () => {
    const server = fakeServer();
    renderDomains(server);

    fireEvent.click(await screen.findByRole('button', { name: 'Ajouter le domaine' }));

    expect(screen.getByText('Saisissez le nom de domaine.')).toBeInTheDocument();
    expect(screen.getByText('Saisissez un port de 1 à 65535.')).toBeInTheDocument();
    expect(screen.getByLabelText('Nom de domaine')).toHaveAttribute('aria-invalid', 'true');
    expect(server.calls.added).toEqual([]);
  });

  it('refuse un port hors bornes', async () => {
    renderDomains(fakeServer());
    await screen.findByRole('button', { name: 'Ajouter le domaine' });

    fill('Port du service', '70000');
    fireEvent.click(screen.getByRole('button', { name: 'Ajouter le domaine' }));

    expect(screen.getByText('Saisissez un port de 1 à 65535.')).toBeInTheDocument();
  });

  it('ajoute un domaine et l’insère à sa place', async () => {
    const server = fakeServer();
    renderDomains(server);
    await screen.findByRole('button', { name: 'Ajouter le domaine' });

    fill('Nom de domaine', '  Boutique.exemple.com ');
    fill('Projet', 'boutique');
    fill('Service', 'web');
    fill('Port du service', '3000');
    fireEvent.click(screen.getByRole('button', { name: 'Ajouter le domaine' }));

    expect(await screen.findByRole('status')).toHaveTextContent(
      'boutique.exemple.com est relié à web.'
    );
    expect(server.calls.added).toEqual([
      expect.objectContaining({
        host: 'Boutique.exemple.com',
        projectId: 'boutique',
        service: 'web',
        port: 3000,
      }),
    ]);
    expect(await domainRows()).toEqual([
      'api.exemple.com(ouvre un nouvel onglet)',
      'boutique.exemple.com(ouvre un nouvel onglet)',
      'www.exemple.com(ouvre un nouvel onglet)',
    ]);
    expect(screen.getByLabelText('Nom de domaine')).toHaveValue('');
  });

  it('montre le refus du serveur', async () => {
    renderDomains(
      fakeServer({
        addError: new ConnectError('ce nom de domaine est déjà utilisé', Code.AlreadyExists),
      })
    );
    await screen.findByRole('button', { name: 'Ajouter le domaine' });

    fill('Nom de domaine', 'api.exemple.com');
    fill('Projet', 'boutique');
    fill('Service', 'api');
    fill('Port du service', '8080');
    fireEvent.click(screen.getByRole('button', { name: 'Ajouter le domaine' }));

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Ce nom de domaine est déjà utilisé.'
    );
    expect(screen.getByLabelText('Nom de domaine')).toHaveValue('api.exemple.com');
  });

  it('retire un domaine après confirmation, et pas avant', async () => {
    const server = fakeServer();
    renderDomains(server);

    fireEvent.click(await screen.findByRole('button', { name: 'Retirer www.exemple.com' }));
    const confirm = screen.getByRole('group', { name: 'Retirer www.exemple.com' });
    fireEvent.click(within(confirm).getByRole('button', { name: 'Annuler' }));
    expect(server.calls.removed).toEqual([]);

    fireEvent.click(screen.getByRole('button', { name: 'Retirer www.exemple.com' }));
    fireEvent.click(screen.getByRole('button', { name: 'Confirmer le retrait' }));

    await waitFor(() => expect(server.calls.removed).toEqual(['dom_2']));
    expect(await domainRows()).toEqual(['api.exemple.com(ouvre un nouvel onglet)']);
    expect(screen.getByRole('heading', { name: '1 domaine servi' })).toBeInTheDocument();
  });

  it('laisse un membre consulter sans modifier', async () => {
    renderDomains(fakeServer({ role: UserRole.MEMBER }));

    expect(await domainRows()).toHaveLength(2);
    expect(
      screen.getByText('Seul un administrateur peut ajouter ou retirer un domaine.')
    ).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Ajouter le domaine' })).toBeNull();
    expect(screen.queryByRole('button', { name: /^Retirer/ })).toBeNull();
  });

  it('explique un reverse proxy désactivé', async () => {
    renderDomains(fakeServer({ disabled: true }));

    expect(
      await screen.findByText('Le reverse proxy n’est pas activé sur ce serveur')
    ).toBeInTheDocument();
    expect(screen.getByText(/activez proxy\.enabled/)).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Réessayer' })).toBeNull();
  });
});
