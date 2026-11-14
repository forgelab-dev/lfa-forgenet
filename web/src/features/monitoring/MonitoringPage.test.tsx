import { create, toJson } from '@bufbuild/protobuf';
import { timestampFromMs, timestampMs } from '@bufbuild/protobuf/wkt';
import { createRouterTransport } from '@connectrpc/connect';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { memoryLocation } from 'wouter/memory-location';

import { App } from '../../app/App';
import { AuthService, UserRole } from '../../gen/forgenet/v1/auth_pb';
import { ContainerService, ContainerState } from '../../gen/forgenet/v1/container_pb';
import {
  MetricName,
  MetricsService,
  StreamAllContainerMetricsResponseSchema,
  StreamContainerMetricsResponseSchema,
  StreamSystemMetricsResponseSchema,
  type QueryContainerMetricsRequest,
  type QuerySystemMetricsRequest,
} from '../../gen/forgenet/v1/metrics_pb';
import { sessionInterceptors, type TransportFactory } from '../../lib/api/client';

const GiB = 1024 ** 3;
const MiB = 1024 ** 2;
const WEB = 'a1f3c9e20b7d' + '0'.repeat(52);
const DB = 'b2e4d8f10c6a' + '0'.repeat(52);

/** Faux WebSocket : le test ouvre, alimente et coupe chaque flux. */
class FakeSocket {
  static all: FakeSocket[] = [];
  onopen: (() => void) | null = null;
  onmessage: ((e: { data: string }) => void) | null = null;
  onclose: ((e: { code: number }) => void) | null = null;
  closed = false;

  constructor(
    readonly url: string,
    readonly protocols: string[]
  ) {
    FakeSocket.all.push(this);
    queueMicrotask(() => this.onopen?.());
  }

  get path() {
    return new URL(this.url).pathname;
  }

  send() {}

  close(code = 1000) {
    if (this.closed) return;
    this.closed = true;
    this.onclose?.({ code });
  }

  emit(json: unknown) {
    this.onmessage?.({ data: JSON.stringify(json) });
  }
}

function openSocket(path: string): FakeSocket {
  const found = FakeSocket.all.filter((s) => s.path === path && !s.closed).at(-1);
  if (!found) throw new Error(`aucun flux ouvert sur ${path}`);
  return found;
}

async function socketOn(path: string): Promise<FakeSocket> {
  let socket: FakeSocket | undefined;
  await waitFor(() => {
    socket = openSocket(path);
  });
  return socket!;
}

function systemMessage(at: number, extra: { cpu?: number; disk?: number } = {}) {
  return toJson(
    StreamSystemMetricsResponseSchema,
    create(StreamSystemMetricsResponseSchema, {
      metrics: {
        collectedAt: timestampFromMs(at),
        uptimeSeconds: BigInt(12 * 86400 + 4 * 3600),
        activeContainersCount: 2,
        cpu: { cpuUsagePercent: extra.cpu ?? 37.2, onlineCpus: 4 },
        memory: {
          usageBytes: BigInt(5 * GiB),
          limitBytes: BigInt(8 * GiB),
          cacheBytes: BigInt(GiB),
          usagePercent: 62.5,
        },
        storage: {
          totalBytes: BigInt(128 * GiB),
          freeBytes: BigInt(15 * GiB),
          usedBytes: BigInt(113 * GiB),
          usagePercent: extra.disk ?? 40,
        },
        network: { rxBytesTotal: BigInt(at), txBytesTotal: BigInt(Math.floor(at / 2)) },
      },
    })
  );
}

function sweep(at: number, containers: { id: string; cpu: number; mem: number; limit: number }[]) {
  return toJson(
    StreamAllContainerMetricsResponseSchema,
    create(StreamAllContainerMetricsResponseSchema, {
      collectedAt: timestampFromMs(at),
      containers: containers.map((c) => ({
        containerId: c.id,
        collectedAt: timestampFromMs(at),
        cpu: { cpuUsagePercent: c.cpu, onlineCpus: 4 },
        memory: {
          usageBytes: BigInt(Math.round(c.mem)),
          limitBytes: BigInt(c.limit),
          usagePercent: (c.mem / c.limit) * 100,
        },
      })),
    })
  );
}

function fakeServer(role = UserRole.ADMIN) {
  const calls = {
    system: [] as QuerySystemMetricsRequest[],
    container: [] as QueryContainerMetricsRequest[],
    lists: 0,
  };
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
        service(ContainerService, {
          listContainers() {
            calls.lists++;
            return {
              containers: [
                {
                  id: WEB,
                  name: 'forgenet-boutique-web-1',
                  image: 'nginx:alpine',
                  state: ContainerState.RUNNING,
                },
                {
                  id: DB,
                  name: 'forgenet-boutique-db-1',
                  image: 'postgres:16',
                  state: ContainerState.RUNNING,
                },
                {
                  id: 'c'.repeat(64),
                  name: 'forgenet-ancien-1',
                  image: 'redis:7',
                  state: ContainerState.EXITED,
                },
              ],
            };
          },
        });
        service(MetricsService, {
          querySystemMetrics(req) {
            calls.system.push(req);
            const now = Date.now();
            return {
              series: [
                {
                  key: { key: { case: 'known', value: MetricName.CPU_USAGE_PERCENT } },
                  samples: [
                    { timestamp: timestampFromMs(now - 6 * 60_000), value: 81 },
                    { timestamp: timestampFromMs(now - 5 * 60_000), value: 20 },
                  ],
                },
              ],
            };
          },
          queryContainerMetrics(req) {
            calls.container.push(req);
            return { containerId: req.containerId, series: [] };
          },
        });
      },
      { transport: { interceptors: sessionInterceptors(session) } }
    );
  return { factory, calls };
}

function renderPage(server: ReturnType<typeof fakeServer>) {
  sessionStorage.setItem(
    'forgenet.session',
    JSON.stringify({ token: 'jeton', expiresAt: Date.now() + 3_600_000 })
  );
  const location = memoryLocation({ path: '/surveillance' });
  render(<App transportFactory={server.factory} locationHook={location.hook} />);
}

let hidden = false;

beforeEach(() => {
  FakeSocket.all = [];
  hidden = false;
  vi.stubGlobal('WebSocket', FakeSocket);
  Object.defineProperty(document, 'hidden', { configurable: true, get: () => hidden });
});

afterEach(() => {
  vi.unstubAllGlobals();
  sessionStorage.clear();
});

describe('écran Surveillance', () => {
  it('montre l’hôte et les conteneurs en direct, jeton en sous-protocole', async () => {
    const server = fakeServer();
    renderPage(server);

    const host = await socketOn('/ws/v1/metrics/system');
    expect(host.protocols).toEqual(['forgenet.v1', 'forgenet.bearer.jeton']);
    expect(server.calls.system[0]?.range?.startTime).toBeDefined();

    host.emit(systemMessage(Date.now()));
    expect(await screen.findByText('12 j 4 h')).toBeInTheDocument();
    const memory = screen.getByRole('region', { name: 'Mémoire' });
    expect(within(memory).getByText(/5,0 Go sur 8,0 Go/)).toBeInTheDocument();
    expect(
      within(screen.getByRole('region', { name: 'Processeur' })).getByText(/pointe à 81 %/)
    ).toBeInTheDocument();
    const fleet = await socketOn('/ws/v1/metrics/containers');
    fleet.emit(
      sweep(Date.now(), [
        { id: DB, cpu: 18.4, mem: 1.4 * GiB, limit: 8 * GiB },
        { id: WEB, cpu: 42.1, mem: 890 * MiB, limit: 1 * GiB },
      ])
    );
    await screen.findByRole('button', { name: /forgenet-boutique-db-1/ });
    const names = screen
      .getAllByRole('button', { expanded: false })
      .map((r) => r.textContent ?? '')
      .filter((t) => t.startsWith('forgenet-'));
    expect(names[0]).toMatch(/^forgenet-boutique-web-1/); // trié par processeur
    expect(screen.getByText('sans limite', { exact: false })).toBeInTheDocument(); // limite = mémoire de l'hôte
    expect(screen.getByText('2 actifs · mesurés toutes les 10 s')).toBeInTheDocument();
    expect(screen.getByText('En direct')).toBeInTheDocument();
    expect(screen.getByText('1 conteneur arrêté, sans mesure')).toBeInTheDocument();
  });

  it('signale le disque au-delà de 85 %', async () => {
    renderPage(fakeServer());
    (await socketOn('/ws/v1/metrics/system')).emit(systemMessage(Date.now(), { disk: 88 }));
    const disk = await screen.findByRole('region', { name: 'Disque' });
    expect(within(disk).getByText('au-delà de 85 %')).toBeInTheDocument();
  });

  it('comble le trou après une coupure, puis reprend le direct', async () => {
    const server = fakeServer();
    renderPage(server);
    const host = await socketOn('/ws/v1/metrics/system');
    const last = Date.now();
    host.emit(systemMessage(last));
    (await socketOn('/ws/v1/metrics/containers')).emit(sweep(last, []));
    await screen.findByText('En direct');

    host.close(1006);
    expect(await screen.findByText('Liaison coupée')).toBeInTheDocument();
    expect(screen.getByText(/mesures de/)).toBeInTheDocument();

    await waitFor(() => expect(server.calls.system).toHaveLength(2), { timeout: 3000 });
    expect(timestampMs(server.calls.system[1]!.range!.startTime!)).toBe(last);
    const again = await socketOn('/ws/v1/metrics/system');
    expect(again).not.toBe(host);
    again.emit(systemMessage(Date.now()));
    expect(await screen.findByText('En direct')).toBeInTheDocument();
  });

  it('ferme la session quand le flux refuse le jeton', async () => {
    renderPage(fakeServer());
    (await socketOn('/ws/v1/metrics/system')).close(4401);
    expect(await screen.findByText(/Votre session a pris fin/)).toBeInTheDocument();
  });

  it('ouvre le détail d’un conteneur et suit son flux à part', async () => {
    const server = fakeServer();
    renderPage(server);
    (await socketOn('/ws/v1/metrics/containers')).emit(
      sweep(Date.now(), [{ id: WEB, cpu: 42.1, mem: 890 * MiB, limit: 1 * GiB }])
    );
    const row = await screen.findByRole('button', { name: /forgenet-boutique-web-1/ });
    fireEvent.click(row);
    expect(row).toHaveAttribute('aria-expanded', 'true');

    const detail = await socketOn(`/ws/v1/metrics/containers/${WEB}`);
    await waitFor(() =>
      expect(server.calls.container.some((c) => c.keys.length === 4 && c.containerId === WEB)).toBe(
        true
      )
    );
    detail.emit(
      toJson(
        StreamContainerMetricsResponseSchema,
        create(StreamContainerMetricsResponseSchema, {
          metrics: {
            containerId: WEB,
            collectedAt: timestampFromMs(Date.now()),
            cpu: { cpuUsagePercent: 40 },
            memory: { usageBytes: BigInt(900 * MiB), limitBytes: BigInt(GiB), usagePercent: 87.9 },
          },
        })
      )
    );
    expect(await screen.findByText('900 Mo / 1,0 Go')).toBeInTheDocument();

    detail.close(4404);
    expect(await screen.findByText(/Ce conteneur s’est arrêté/)).toBeInTheDocument();

    fireEvent.click(row);
    expect(row).toHaveAttribute('aria-expanded', 'false');
  });

  it('met les flux en pause quand l’onglet est caché', async () => {
    const server = fakeServer();
    renderPage(server);
    const host = await socketOn('/ws/v1/metrics/system');
    host.emit(systemMessage(Date.now()));
    (await socketOn('/ws/v1/metrics/containers')).emit(sweep(Date.now(), []));
    await screen.findByText('En direct');

    hidden = true;
    document.dispatchEvent(new Event('visibilitychange'));
    expect(await screen.findByText('En pause')).toBeInTheDocument();
    await waitFor(() => expect(FakeSocket.all.every((s) => s.closed)).toBe(true));

    hidden = false;
    document.dispatchEvent(new Event('visibilitychange'));
    const again = await socketOn('/ws/v1/metrics/system');
    expect(again).not.toBe(host);
    expect(server.calls.system).toHaveLength(2);
  });

  it('dit quand aucun conteneur ne tourne, sans proposer de déployer à un membre', async () => {
    renderPage(fakeServer(UserRole.MEMBER));
    (await socketOn('/ws/v1/metrics/containers')).emit(sweep(Date.now(), []));
    expect(await screen.findByText('Aucun conteneur actif')).toBeInTheDocument();
    expect(
      screen.queryByRole('link', { name: 'Déployer une application' })
    ).not.toBeInTheDocument();
  });

  it('liste les conteneurs arrêtés à la demande', async () => {
    renderPage(fakeServer());
    (await socketOn('/ws/v1/metrics/containers')).emit(
      sweep(Date.now(), [{ id: WEB, cpu: 1, mem: MiB, limit: GiB }])
    );
    fireEvent.click(await screen.findByRole('button', { name: 'Afficher' }));
    expect(screen.getByText('forgenet-ancien-1')).toBeInTheDocument();
  });
});
