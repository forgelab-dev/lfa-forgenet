import { timestampFromMs, type Timestamp } from '@bufbuild/protobuf/wkt';
import { create } from '@bufbuild/protobuf';
import { Code, ConnectError, createRouterTransport } from '@connectrpc/connect';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { memoryLocation } from 'wouter/memory-location';

import { App } from '../../app/App';
import { AuthService, UserRole } from '../../gen/forgenet/v1/auth_pb';
import { LogStream } from '../../gen/forgenet/v1/container_pb';
import {
  DeploymentLogEntrySchema,
  DeploymentSchema,
  DeploymentService,
  DeploymentState,
  type CreateDeploymentRequest,
  type Deployment,
  type DeploymentLogEntry,
  type StreamDeploymentLogsRequest,
} from '../../gen/forgenet/v1/deploy_pb';
import { sessionInterceptors, type TransportFactory } from '../../lib/api/client';

const minute = 60_000;
const t0 = Date.now() - 60 * minute;

function deployment(
  id: string,
  projectId: string,
  state: DeploymentState,
  minutesAgo: number,
  extra: { errorMessage?: string } = {}
): Deployment {
  return create(DeploymentSchema, {
    id,
    projectId,
    state,
    triggeredBy: 'user:u1',
    createdAt: timestampFromMs(Date.now() - minutesAgo * minute),
    ...extra,
  });
}

function line(
  ms: number,
  message: string,
  step = 'deploy',
  stream = LogStream.STDERR
): DeploymentLogEntry {
  return create(DeploymentLogEntrySchema, {
    emittedAt: timestampFromMs(t0 + ms),
    message,
    step,
    stream,
  });
}

const later = (ts: Timestamp | undefined, than: Timestamp | undefined) =>
  !than ||
  (ts && (ts.seconds > than.seconds || (ts.seconds === than.seconds && ts.nanos > than.nanos)));

interface Options {
  role?: UserRole;
  initial?: Deployment[];
  /** États envoyés par StreamDeployment pour un déploiement créé ici. */
  states?: DeploymentState[];
  /** Journal du déploiement créé ici. */
  logs?: DeploymentLogEntry[];
  /** Le premier flux de journal se coupe après n lignes. */
  dropLogsAfter?: number;
  deleteError?: ConnectError;
}

function fakeServer(options: Options = {}) {
  const {
    role = UserRole.ADMIN,
    initial = [],
    states = [
      DeploymentState.BUILDING,
      DeploymentState.DEPLOYING,
      DeploymentState.HEALTH_CHECKING,
      DeploymentState.SUCCESS,
    ],
    logs = [],
  } = options;
  const calls = {
    list: 0,
    created: [] as CreateDeploymentRequest[],
    cancelled: [] as string[],
    deleted: [] as { projectId: string; deleteVolumes: boolean }[],
    logRequests: [] as StreamDeploymentLogsRequest[],
  };
  let dropped = false;
  // Libère le flux d'état ; le test décide quand le déploiement avance.
  let release: () => void = () => {};
  const advance = new Promise<void>((resolve) => (release = resolve));

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
        service(DeploymentService, {
          listDeployments() {
            calls.list++;
            return { deployments: initial };
          },
          createDeployment(req) {
            calls.created.push(req);
            if (req.projectId === 'refuse') {
              throw new ConnectError(
                'fichier Compose invalide : services.web.image must be a string',
                Code.InvalidArgument
              );
            }
            return { deployment: deployment('dep_new', req.projectId, DeploymentState.PENDING, 0) };
          },
          cancelDeployment(req) {
            calls.cancelled.push(req.deploymentId);
            const d = initial.find((x) => x.id === req.deploymentId)!;
            return {
              deployment: { ...d, state: DeploymentState.CANCELLED, errorMessage: 'annulé' },
            };
          },
          deleteStack(req) {
            calls.deleted.push({ projectId: req.projectId, deleteVolumes: req.deleteVolumes });
            if (options.deleteError) throw options.deleteError;
            return {};
          },
          async *streamDeployment(req) {
            const d = initial.find((x) => x.id === req.deploymentId);
            if (d) {
              yield { deployment: d };
              return new Promise<void>(() => {}); // reste ouvert : déploiement en cours
            }
            yield {
              deployment: deployment(req.deploymentId, 'boutique', DeploymentState.PENDING, 0),
            };
            await advance;
            for (const state of states) {
              yield { deployment: deployment(req.deploymentId, 'boutique', state, 0) };
            }
          },
          async *streamDeploymentLogs(req) {
            calls.logRequests.push(req);
            const rest = logs.filter((l) => later(l.emittedAt, req.since));
            for (const [i, entry] of rest.entries()) {
              if (options.dropLogsAfter !== undefined && !dropped && i === options.dropLogsAfter) {
                dropped = true;
                throw new ConnectError(
                  'flux en retard sur le déploiement : reconnectez-vous',
                  Code.Unavailable
                );
              }
              yield { logEntry: entry };
            }
          },
        });
      },
      { transport: { interceptors: sessionInterceptors(session) } }
    );
  return { factory, calls, release: () => release() };
}

function renderPage(server: ReturnType<typeof fakeServer>) {
  sessionStorage.setItem(
    'forgenet.session',
    JSON.stringify({ token: 'jeton', expiresAt: Date.now() + 3_600_000 })
  );
  const location = memoryLocation({ path: '/deploiements' });
  render(<App transportFactory={server.factory} locationHook={location.hook} />);
}

function fill(label: string, value: string) {
  fireEvent.change(screen.getByLabelText(label), { target: { value } });
}

async function appRow(project: string) {
  const list = await screen.findByRole('list', { name: 'Applications' });
  const heading = await within(list).findByRole('heading', { name: project });
  return heading.closest('li')!;
}

describe('DeploymentsPage', () => {
  beforeEach(() => sessionStorage.clear());

  it('montre aux membres que la section est réservée, sans rien charger', async () => {
    const server = fakeServer({ role: UserRole.MEMBER });
    renderPage(server);

    expect(await screen.findByText('Réservé aux administrateurs')).toBeInTheDocument();
    expect(server.calls.list).toBe(0);
    expect(screen.queryByRole('button', { name: 'Déployer' })).toBeNull();
  });

  it('regroupe les déploiements par application, la plus récente d’abord', async () => {
    renderPage(
      fakeServer({
        initial: [
          deployment('dep_1', 'boutique', DeploymentState.SUCCESS, 30),
          deployment('dep_2', 'api', DeploymentState.FAILED, 10, {
            errorMessage: 'docker compose up : pull access denied for api',
          }),
          deployment('dep_0', 'boutique', DeploymentState.FAILED, 120),
        ],
      })
    );

    const list = await screen.findByRole('list', { name: 'Applications' });
    expect(
      within(list)
        .getAllByRole('heading')
        .map((h) => h.textContent)
    ).toEqual(['api', 'boutique']);
    expect(screen.getByRole('heading', { name: '2 applications' })).toBeInTheDocument();

    const api = await appRow('api');
    expect(within(api).getByText('Échec')).toBeInTheDocument();
    expect(
      within(api).getByText('docker compose up : pull access denied for api')
    ).toBeInTheDocument();

    const boutique = await appRow('boutique');
    expect(within(boutique).getByText('Réussi')).toBeInTheDocument();
    fireEvent.click(within(boutique).getByRole('button', { name: 'Historique (1)' }));
    const history = within(boutique).getByRole('list', {
      name: 'Déploiements précédents de boutique',
    });
    expect(within(history).getByText('Échec')).toBeInTheDocument();
  });

  it('vérifie le formulaire avant d’appeler le serveur', async () => {
    const server = fakeServer();
    renderPage(server);

    fireEvent.click(await screen.findByRole('button', { name: 'Déployer' }));
    expect(screen.getByText('Collez le fichier Compose ou importez-le.')).toBeInTheDocument();
    expect(screen.getByText('Saisissez l’identifiant du projet.')).toBeInTheDocument();

    fill('Fichier Compose', 'services:\n  web:\n    image: nginx\n');
    fill('Projet', 'Ma Boutique');
    fireEvent.click(screen.getByRole('button', { name: 'Déployer' }));
    expect(screen.getByLabelText('Projet')).toHaveAttribute('aria-invalid', 'true');
    expect(screen.getByText(/Minuscules, chiffres/)).toBeInTheDocument();
    expect(server.calls.created).toEqual([]);
  });

  it('propose l’identifiant depuis le fichier importé, sans écraser une saisie', async () => {
    renderPage(fakeServer());
    await screen.findByRole('button', { name: 'Déployer' });

    const input = document.querySelector<HTMLInputElement>('input[type="file"]')!;
    const file = new File(['services:\n  web:\n    image: nginx\n'], 'Blog Perso.yml', {
      type: 'application/yaml',
    });
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => expect(screen.getByLabelText('Projet')).toHaveValue('blog-perso'));
    expect(screen.getByLabelText('Fichier Compose')).toHaveValue(
      'services:\n  web:\n    image: nginx\n'
    );
    expect(screen.getByText(/Importé de Blog Perso\.yml/)).toBeInTheDocument();

    fill('Projet', 'mon-blog');
    fill('Fichier Compose', 'name: autre\nservices: {}\n');
    expect(screen.getByLabelText('Projet')).toHaveValue('mon-blog');
  });

  it('déploie, puis suit les étapes et le journal en direct jusqu’au succès', async () => {
    const server = fakeServer({
      logs: [
        line(0, 'Vérification du fichier Compose', 'validate', LogStream.STDOUT),
        line(1, ' Container forgenet-boutique-web-1 Started'),
        line(2, 'Déploiement réussi', 'result', LogStream.STDOUT),
      ],
    });
    renderPage(server);
    await screen.findByRole('button', { name: 'Déployer' });

    fill('Fichier Compose', 'name: boutique\nservices:\n  web:\n    image: nginx\n');
    expect(screen.getByLabelText('Projet')).toHaveValue('boutique');
    fireEvent.click(screen.getByRole('button', { name: 'Déployer' }));

    const row = await appRow('boutique');
    expect(server.calls.created).toHaveLength(1);
    expect(server.calls.created[0].source?.source).toEqual({
      case: 'compose',
      value: expect.objectContaining({
        composeContent: 'name: boutique\nservices:\n  web:\n    image: nginx\n',
      }),
    });
    expect(server.calls.created[0].requestId).not.toBe('');
    expect(screen.getByLabelText('Projet')).toHaveValue('');

    expect(within(row).getByRole('list', { name: 'Étapes du déploiement' })).toBeInTheDocument();
    expect(within(row).getByRole('button', { name: 'Annuler le déploiement' })).toBeInTheDocument();
    const log = within(row).getByRole('log', { name: 'Journal du déploiement de boutique' });
    expect(await within(log).findByText('Déploiement réussi')).toBeInTheDocument();
    expect(server.calls.logRequests[0].follow).toBe(true);

    server.release();
    await waitFor(() => expect(within(row).getByText('Réussi')).toBeInTheDocument());
    expect(within(row).queryByRole('list', { name: 'Étapes du déploiement' })).toBeNull();
    expect(within(row).queryByRole('button', { name: 'Annuler le déploiement' })).toBeNull();
  });

  it('reprend le journal après une coupure, sans perte ni doublon', async () => {
    const lines = [line(0, 'un'), line(1, 'deux'), line(2, 'trois')];
    const server = fakeServer({
      initial: [deployment('dep_1', 'boutique', DeploymentState.SUCCESS, 5)],
      logs: lines,
      dropLogsAfter: 2,
    });
    renderPage(server);

    const row = await appRow('boutique');
    fireEvent.click(within(row).getByRole('button', { name: 'Journal' }));
    const log = within(row).getByRole('log');
    expect(await within(log).findByText('deux')).toBeInTheDocument();
    expect(await within(row).findByText('Liaison coupée, reconnexion…')).toBeInTheDocument();

    expect(await within(log).findByText('trois', {}, { timeout: 3000 })).toBeInTheDocument();
    expect(within(log).getAllByText('deux')).toHaveLength(1);
    expect(server.calls.logRequests).toHaveLength(2);
    expect(server.calls.logRequests[1].since).toEqual(lines[1].emittedAt);
    expect(server.calls.logRequests[1].follow).toBe(false);
  });

  it('annule un déploiement en cours', async () => {
    const server = fakeServer({
      initial: [deployment('dep_1', 'boutique', DeploymentState.DEPLOYING, 1)],
    });
    renderPage(server);

    const row = await appRow('boutique');
    expect(within(row).getByRole('button', { name: 'Supprimer la stack boutique' })).toBeDisabled();
    fireEvent.click(within(row).getByRole('button', { name: 'Annuler le déploiement' }));

    await waitFor(() => expect(within(row).getByText('Annulé')).toBeInTheDocument());
    expect(server.calls.cancelled).toEqual(['dep_1']);
  });

  it('montre le refus du serveur à la création', async () => {
    renderPage(fakeServer());
    await screen.findByRole('button', { name: 'Déployer' });

    fill('Fichier Compose', 'services:\n  web:\n    image: [1]\n');
    fill('Projet', 'refuse');
    fireEvent.click(screen.getByRole('button', { name: 'Déployer' }));

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Fichier Compose invalide : services.web.image must be a string.'
    );
    expect(screen.getByLabelText('Projet')).toHaveValue('refuse');
  });

  it('supprime une stack, et exige le nom du projet pour effacer ses données', async () => {
    const server = fakeServer({
      initial: [deployment('dep_1', 'boutique', DeploymentState.SUCCESS, 5)],
    });
    renderPage(server);
    const row = await appRow('boutique');

    fireEvent.click(within(row).getByRole('button', { name: 'Supprimer la stack boutique' }));
    let confirm = within(row).getByRole('group', { name: 'Supprimer la stack boutique' });
    fireEvent.click(within(confirm).getByRole('button', { name: 'Annuler' }));
    expect(server.calls.deleted).toEqual([]);

    fireEvent.click(within(row).getByRole('button', { name: 'Supprimer la stack boutique' }));
    confirm = within(row).getByRole('group', { name: 'Supprimer la stack boutique' });
    fireEvent.click(within(confirm).getByRole('checkbox'));
    const erase = within(confirm).getByRole('button', {
      name: 'Supprimer la stack et ses données',
    });
    expect(erase).toBeDisabled();
    fireEvent.change(within(confirm).getByLabelText(/pour confirmer/), {
      target: { value: 'boutique' },
    });
    expect(erase).toBeEnabled();
    fireEvent.click(erase);

    expect(await within(row).findByRole('status')).toHaveTextContent(
      'Stack supprimée, volumes compris'
    );
    expect(server.calls.deleted).toEqual([{ projectId: 'boutique', deleteVolumes: true }]);
  });

  it('explique un refus de suppression', async () => {
    renderPage(
      fakeServer({
        initial: [deployment('dep_1', 'boutique', DeploymentState.SUCCESS, 5)],
        deleteError: new ConnectError(
          'un déploiement du projet est en attente ou en cours',
          Code.FailedPrecondition
        ),
      })
    );
    const row = await appRow('boutique');

    fireEvent.click(within(row).getByRole('button', { name: 'Supprimer la stack boutique' }));
    fireEvent.click(within(row).getByRole('button', { name: 'Supprimer la stack' }));
    expect(await within(row).findByRole('alert')).toHaveTextContent(
      'Un déploiement du projet est en attente ou en cours.'
    );
  });
});
