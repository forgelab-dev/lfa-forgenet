import { Code, ConnectError, createRouterTransport } from '@connectrpc/connect';
import { describe, expect, it } from 'vitest';

import { ContainerService } from '../../gen/forgenet/v1/container_pb';
import { bearerInterceptor, createApiClients } from './client';

/** Un faux serveur qui n'accepte que le jeton « bon-jeton ». */
function fakeServer(token: string | undefined) {
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
    },
    { transport: { interceptors: [bearerInterceptor(() => token)] } }
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
});
