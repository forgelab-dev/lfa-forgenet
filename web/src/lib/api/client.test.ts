import { Code, ConnectError, createRouterTransport } from '@connectrpc/connect';
import { describe, expect, it, vi } from 'vitest';

import { AuthService } from '../../gen/forgenet/v1/auth_pb';
import { ContainerService } from '../../gen/forgenet/v1/container_pb';
import { createApiClients, sessionInterceptors } from './client';

/** Un faux serveur qui n'accepte que le jeton « bon-jeton ». */
function fakeServer(token: string | undefined, onSessionExpired = () => {}) {
  return createRouterTransport(
    ({ service }) => {
      service(ContainerService, {
        listContainers(_req, ctx) {
          if (ctx.requestHeader.get('Authorization') !== 'Bearer bon-jeton') {
            throw new ConnectError('authentification requise', Code.Unauthenticated);
          }
          return { containers: [{ id: 'c1', name: 'web' }] };
        },
        async *streamContainerLogs(req) {
          for (const line of ['un', 'deux']) {
            yield { chunk: new TextEncoder().encode(`${req.containerId}:${line}`) };
          }
        },
      });
      service(AuthService, {
        login() {
          throw new ConnectError('identifiant ou mot de passe incorrect', Code.Unauthenticated);
        },
      });
    },
    { transport: { interceptors: sessionInterceptors({ token: () => token, onSessionExpired }) } }
  );
}

describe('clients de l’API', () => {
  it('appelle un service du contrat avec le jeton de la session', async () => {
    const api = createApiClients(fakeServer('bon-jeton'));

    const res = await api.containers.listContainers({});

    expect(res.containers.map((c) => c.name)).toEqual(['web']);
  });

  it('n’envoie pas d’en-tête hors session, et remonte le refus', async () => {
    const api = createApiClients(fakeServer(undefined));

    await expect(api.containers.listContainers({})).rejects.toMatchObject({
      code: Code.Unauthenticated,
    });
  });

  it('lit un flux serveur message par message', async () => {
    const api = createApiClients(fakeServer('bon-jeton'));
    const lines: string[] = [];

    for await (const msg of api.containers.streamContainerLogs({ containerId: 'c1' })) {
      lines.push(new TextDecoder().decode(msg.chunk));
    }

    expect(lines).toEqual(['c1:un', 'c1:deux']);
  });

  it('signale la fin de session quand le serveur refuse le jeton présenté', async () => {
    const onExpired = vi.fn();
    const api = createApiClients(fakeServer('jeton-revoque', onExpired));

    await expect(api.containers.listContainers({})).rejects.toBeInstanceOf(ConnectError);

    expect(onExpired).toHaveBeenCalledOnce();
  });

  it('ne prend pas un mot de passe refusé pour une fin de session', async () => {
    const onExpired = vi.fn();
    const api = createApiClients(fakeServer(undefined, onExpired));

    await expect(api.auth.login({ login: 'mike', password: 'faux' })).rejects.toMatchObject({
      code: Code.Unauthenticated,
    });

    expect(onExpired).not.toHaveBeenCalled();
  });
});
