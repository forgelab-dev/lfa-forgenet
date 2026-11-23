import { create } from '@bufbuild/protobuf';
import { timestampFromMs } from '@bufbuild/protobuf/wkt';
import { Code, ConnectError, createRouterTransport } from '@connectrpc/connect';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { memoryLocation } from 'wouter/memory-location';

import { App } from '../../app/App';
import { AuthService, UserRole } from '../../gen/forgenet/v1/auth_pb';
import { LogStream } from '../../gen/forgenet/v1/container_pb';
import {
  DeploymentSchema,
  DeploymentService,
  DeploymentState,
} from '../../gen/forgenet/v1/deploy_pb';
import {
  MarketplaceService,
  TemplateSchema,
  type InstallTemplateRequest,
} from '../../gen/forgenet/v1/marketplace_pb';
import { sessionInterceptors, type TransportFactory } from '../../lib/api/client';

const icon =
  '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1 1"><rect width="1" height="1"/></svg>';
const generated = 'Xk2pQ8vT4mLr9sWb3nYc7hJd';

const postgres = create(TemplateSchema, {
  id: 'postgres',
  name: 'PostgreSQL',
  description: 'La base de données relationnelle.',
  category: 'database',
  version: '17.2',
  license: 'PostgreSQL',
  homepage: 'https://www.postgresql.org',
  iconSvg: icon,
  variables: [
    { name: 'POSTGRES_PASSWORD', label: 'Mot de passe', secret: true, generate: 24 },
    { name: 'POSTGRES_USER', label: 'Utilisateur', defaultValue: 'app', pattern: '^[a-z_]+$' },
  ],
  ports: [{ service: 'db', port: 5432, protocol: 'tcp', description: 'Connexions SQL' }],
  volumes: [{ name: 'data', description: 'Les données de la base' }],
});

const redis = create(TemplateSchema, {
  id: 'redis',
  name: 'Redis',
  description: 'Un cache en mémoire.',
  category: 'cache',
  version: '8.0',
  iconSvg: icon,
});

interface Options {
  role?: UserRole;
  templates?: (typeof postgres)[];
  /** Erreur renvoyée par InstallTemplate, à la place d'une installation. */
  installError?: ConnectError;
}

function fakeServer({
  role = UserRole.ADMIN,
  templates = [postgres, redis],
  installError,
}: Options = {}) {
  const calls: { installed: InstallTemplateRequest[]; listed: number } = {
    installed: [],
    listed: 0,
  };
  const created = timestampFromMs(Date.now());
  const deployment = (state: DeploymentState) =>
    create(DeploymentSchema, {
      id: 'dep_1',
      projectId: 'postgres',
      state,
      template: 'postgres',
      triggeredBy: 'user:u1',
      createdAt: created,
    });

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
        service(MarketplaceService, {
          listTemplates() {
            calls.listed++;
            return { templates };
          },
          getTemplate(req) {
            const template = templates.find((t) => t.id === req.id);
            if (!template) {
              throw new ConnectError(`template "${req.id}" absent du catalogue`, Code.NotFound);
            }
            return { template, readme: '' };
          },
          installTemplate(req) {
            calls.installed.push(req);
            if (installError) {
              throw installError;
            }
            return {
              deployment: deployment(DeploymentState.PENDING),
              generatedSecrets: [
                { name: 'POSTGRES_PASSWORD', label: 'Mot de passe', value: generated },
              ],
            };
          },
        });
        service(DeploymentService, {
          async *streamDeployment() {
            yield { deployment: deployment(DeploymentState.DEPLOYING) };
            yield { deployment: { ...deployment(DeploymentState.SUCCESS), finishedAt: created } };
          },
          async *streamDeploymentLogs() {
            yield {
              logEntry: {
                deploymentId: 'dep_1',
                step: 'deploy',
                stream: LogStream.STDOUT,
                message: 'Container postgres-db-1 Started',
                emittedAt: created,
              },
            };
          },
        });
      },
      { transport: { interceptors: sessionInterceptors(session) } }
    );
  return { factory, calls };
}

function renderAt(server: ReturnType<typeof fakeServer>, path: string) {
  sessionStorage.setItem(
    'forgenet.session',
    JSON.stringify({ token: 'jeton', expiresAt: Date.now() + 3_600_000 })
  );
  const location = memoryLocation({ path, record: true });
  render(<App transportFactory={server.factory} locationHook={location.hook} />);
  return location;
}

function fill(label: string | RegExp, value: string) {
  fireEvent.change(screen.getByLabelText(label), { target: { value } });
}

describe('MarketplacePage', () => {
  beforeEach(() => sessionStorage.clear());

  it('présente le catalogue et le filtre par catégorie', async () => {
    renderAt(fakeServer(), '/marketplace');

    const list = await screen.findByRole('list', { name: 'Templates' });
    expect(
      within(list)
        .getAllByRole('heading')
        .map((h) => h.textContent)
    ).toEqual(['PostgreSQL', 'Redis']);
    const card = within(list).getByRole('link', { name: /PostgreSQL/ });
    expect(card).toHaveAttribute('href', '/marketplace/postgres');
    // L'icône est une image : son SVG n'entre jamais dans la page.
    expect(card.querySelector('img')?.getAttribute('src')).toMatch(/^data:image\/svg\+xml/);

    fireEvent.click(screen.getByRole('button', { name: /Cache/ }));
    expect(
      within(list)
        .getAllByRole('heading')
        .map((h) => h.textContent)
    ).toEqual(['Redis']);
    expect(screen.getByRole('button', { name: /Cache/ })).toHaveAttribute('aria-pressed', 'true');
    fireEvent.click(screen.getByRole('button', { name: /Tout/ }));
    expect(within(list).getAllByRole('heading')).toHaveLength(2);
  });

  it('renvoie au guide des contributeurs quand le catalogue est vide', async () => {
    renderAt(fakeServer({ templates: [] }), '/marketplace');

    expect(await screen.findByText('Aucun template pour l’instant')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /guide des contributeurs/ })).toHaveAttribute(
      'href',
      'https://github.com/forgelab-dev/lfa-forgenet/blob/main/templates/README.md'
    );
  });

  it('réserve le catalogue aux administrateurs, sans appeler le serveur', async () => {
    const server = fakeServer({ role: UserRole.MEMBER });
    renderAt(server, '/marketplace');

    expect(await screen.findByText('Réservé aux administrateurs')).toBeInTheDocument();
    expect(server.calls.listed).toBe(0);
  });
});

describe('TemplatePage', () => {
  beforeEach(() => sessionStorage.clear());

  it('installe le template, montre le mot de passe généré une fois, puis suit le déploiement', async () => {
    const server = fakeServer();
    renderAt(server, '/marketplace/postgres');

    expect(
      await screen.findByRole('heading', { level: 1, name: 'PostgreSQL' })
    ).toBeInTheDocument();
    expect(screen.getByText('db:5432')).toBeInTheDocument();
    expect(screen.getByLabelText('Nom de l’application')).toHaveValue('postgres');
    expect(screen.getByLabelText('Utilisateur')).toHaveValue('app');
    const password = screen.getByLabelText('Mot de passe');
    expect(password).toHaveAttribute('type', 'password');
    expect(password).toHaveAttribute('autocomplete', 'new-password');
    expect(
      screen.getByText('Laissez vide : un mot de passe de 24 caractères sera généré.')
    ).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Installer PostgreSQL' }));

    const secrets = await screen.findByRole('region', { name: 'Mots de passe générés' });
    expect(server.calls.installed).toHaveLength(1);
    expect(server.calls.installed[0]).toMatchObject({
      templateId: 'postgres',
      projectId: 'postgres',
      values: { POSTGRES_USER: 'app' },
    });
    expect(server.calls.installed[0].requestId).not.toBe('');
    expect(within(secrets).getByText(/ne les réaffichera pas/)).toBeInTheDocument();

    // Masqué d'abord, lisible sur demande.
    expect(within(secrets).queryByText(generated)).not.toBeInTheDocument();
    fireEvent.click(within(secrets).getByRole('button', { name: 'Afficher Mot de passe' }));
    expect(within(secrets).getByText(generated)).toBeInTheDocument();

    expect(await screen.findAllByText('Réussi')).not.toHaveLength(0);
    expect(await screen.findByText('Container postgres-db-1 Started')).toBeInTheDocument();
    expect(
      screen.getByRole('link', { name: /Voir l’application dans Déploiements/ })
    ).toHaveAttribute('href', '/deploiements');

    // Noté : la valeur quitte la page pour de bon.
    fireEvent.click(within(secrets).getByRole('button', { name: /C’est noté/ }));
    expect(screen.queryByText(generated)).not.toBeInTheDocument();
    expect(screen.queryByRole('region', { name: 'Mots de passe générés' })).not.toBeInTheDocument();
    expect(screen.getByText(/effacés de cette page/)).toBeInTheDocument();
  });

  it('vérifie les valeurs avant l’envoi', async () => {
    const server = fakeServer();
    renderAt(server, '/marketplace/postgres');
    await screen.findByLabelText('Utilisateur');

    fill('Utilisateur', 'App');
    fill('Nom de l’application', 'Mauvais Nom');
    fireEvent.click(screen.getByRole('button', { name: 'Installer PostgreSQL' }));

    expect(
      await screen.findByText('La valeur ne respecte pas le format attendu : ^[a-z_]+$')
    ).toBeInTheDocument();
    expect(screen.getByLabelText('Utilisateur')).toHaveAttribute('aria-invalid', 'true');
    expect(screen.getByLabelText('Nom de l’application')).toHaveAttribute('aria-invalid', 'true');
    expect(server.calls.installed).toHaveLength(0);
  });

  it('range un nom déjà pris sous le champ du nom', async () => {
    const server = fakeServer({
      installError: new ConnectError('une application porte déjà ce nom', Code.AlreadyExists),
    });
    renderAt(server, '/marketplace/postgres');
    await screen.findByLabelText('Nom de l’application');

    fireEvent.click(screen.getByRole('button', { name: 'Installer PostgreSQL' }));

    await waitFor(() =>
      expect(screen.getByLabelText('Nom de l’application')).toHaveAttribute('aria-invalid', 'true')
    );
    expect(screen.getByText('Une application porte déjà ce nom.')).toBeInTheDocument();
  });

  it('range un refus du serveur sous la variable qu’il nomme', async () => {
    const server = fakeServer({
      installError: new ConnectError(
        "valeur d'installation refusée : POSTGRES_USER : valeur obligatoire",
        Code.InvalidArgument
      ),
    });
    renderAt(server, '/marketplace/postgres');
    await screen.findByLabelText('Utilisateur');

    fireEvent.click(screen.getByRole('button', { name: 'Installer PostgreSQL' }));

    expect(await screen.findByText('Valeur obligatoire.')).toBeInTheDocument();
    expect(screen.getByLabelText('Utilisateur')).toHaveAttribute('aria-invalid', 'true');
  });

  it('dit qu’un template est absent du catalogue', async () => {
    renderAt(fakeServer(), '/marketplace/inconnu');

    expect(await screen.findByText('Template "inconnu" absent du catalogue.')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /Retour au catalogue/ })).toHaveAttribute(
      'href',
      '/marketplace'
    );
  });
});
