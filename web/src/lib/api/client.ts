import {
  Code,
  ConnectError,
  createClient,
  type Client,
  type Interceptor,
  type Transport,
} from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';

import { AuthService } from '../../gen/forgenet/v1/auth_pb';
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
 * Prévient la session quand le serveur refuse un jeton présenté : session
 * expirée ou révoquée ailleurs. Un refus sans jeton (mauvais mot de passe à la
 * connexion) ne déclenche rien, ce n'est pas une expiration.
 */
export function sessionExpiryInterceptor(onExpired: () => void): Interceptor {
  return (next) => async (req) => {
    const sentToken = req.header.has('Authorization');
    try {
      return await next(req);
    } catch (err) {
      if (sentToken && ConnectError.from(err).code === Code.Unauthenticated) {
        onExpired();
      }
      throw err;
    }
  };
}

export interface TransportOptions {
  token: TokenSource;
  onSessionExpired: () => void;
}

/** Fabrique le transport de l'application ; les tests en fournissent un faux. */
export type TransportFactory = (options: TransportOptions) => Transport;

/**
 * Intercepteurs communs à tout transport du dashboard. L'ordre compte : le
 * jeton est posé avant que la surveillance d'expiration ne lise l'en-tête.
 */
export function sessionInterceptors({ token, onSessionExpired }: TransportOptions): Interceptor[] {
  return [bearerInterceptor(token), sessionExpiryInterceptor(onSessionExpired)];
}

/**
 * Transport Connect vers le serveur ForgeNet. L'API est servie sous la même
 * origine que le dashboard (proxy Vite en développement, même domaine en
 * production) : la base est donc l'origine de la page, sans CORS.
 */
export function createApiTransport(options: TransportOptions, baseUrl = '/'): Transport {
  return createConnectTransport({ baseUrl, interceptors: sessionInterceptors(options) });
}

/** Clients typés des services du contrat api/, tirés du code généré. */
export interface ApiClients {
  auth: Client<typeof AuthService>;
  containers: Client<typeof ContainerService>;
  deployments: Client<typeof DeploymentService>;
  metrics: Client<typeof MetricsService>;
}

export function createApiClients(transport: Transport): ApiClients {
  return {
    auth: createClient(AuthService, transport),
    containers: createClient(ContainerService, transport),
    deployments: createClient(DeploymentService, transport),
    metrics: createClient(MetricsService, transport),
  };
}
