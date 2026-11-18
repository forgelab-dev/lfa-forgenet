import { ContainerState, type Container } from '../../gen/forgenet/v1/container_pb';

/**
 * Conteneur de la dernière mise à jour (internal/update.UpdaterName). Il reste
 * arrêté après sa fin, pour ses journaux (ADR-0011) : c'est le seul conteneur
 * de ForgeNet dont l'arrêt est normal.
 */
export const UPDATER_NAME = 'forgenet-updater';

/** Sépare les applications des conteneurs de ForgeNet lui-même. */
export function splitSystem(containers: Container[]): { apps: Container[]; system: Container[] } {
  const apps: Container[] = [];
  const system: Container[] = [];
  for (const c of containers) {
    (c.system ? system : apps).push(c);
  }
  return { apps, system };
}

/**
 * Conteneurs de ForgeNet qui devraient tourner et ne tournent pas : le serveur
 * qui répond à cet écran tourne forcément, mais Caddy peut être tombé.
 */
export function stoppedSystem(system: Container[]): Container[] {
  return system.filter((c) => c.name !== UPDATER_NAME && c.state !== ContainerState.RUNNING);
}
