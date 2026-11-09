import { createClient, type Client, type Interceptor, type Transport } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';

import { ContainerService } from '../../gen/forgenet/v1/container_pb';
import { DeploymentService } from '../../gen/forgenet/v1/deploy_pb';
import { MetricsService } from '../../gen/forgenet/v1/metrics_pb';

/** Renvoie le jeton de la session en cours, ou rien hors connexion. */
export type TokenSource = () => string | undefined;

/**
 * Pose `Authorization: Bearer <jeton>` sur chaque appel quand une session est
 * ouverte. Le serveur lit cet en-tête dans les deux portes, gRPC et Connect
 * (ADR-0006).
 */
export function bearerInterceptor(token: TokenSource): Interceptor {
  return (next) => (req) => {
    const value = token();
    if (value) {
      req.header.set('Authorization', `Bearer ${value}`);
    }
    return next(req);
  };
}

/**
 * Transport Connect vers le serveur ForgeNet. L'API est servie sous la même
 * origine que le dashboard (proxy Vite en développement, même domaine en
 * production) : la base est donc l'origine de la page, sans CORS.
 */
export function createApiTransport(token: TokenSource, baseUrl = '/'): Transport {
  return createConnectTransport({ baseUrl, interceptors: [bearerInterceptor(token)] });
}

/** Clients typés des services du contrat api/, tirés du code généré. */
export interface ApiClients {
  containers: Client<typeof ContainerService>;
  deployments: Client<typeof DeploymentService>;
  metrics: Client<typeof MetricsService>;
}

export function createApiClients(transport: Transport): ApiClients {
  return {
    containers: createClient(ContainerService, transport),
    deployments: createClient(DeploymentService, transport),
    metrics: createClient(MetricsService, transport),
  };
}
