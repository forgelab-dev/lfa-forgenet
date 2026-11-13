import { timestampMs } from '@bufbuild/protobuf/wkt';

import type { Deployment, Stack } from '../../gen/forgenet/v1/deploy_pb';

/**
 * Ce que l'on sait des stacks en place : `unknown` tant qu'elles ne sont pas
 * chargées, ou si le serveur n'a pas pu les lire.
 */
export type Stacks = { known: false } | { known: true; byProject: Map<string, Stack> };

/** Une application : un projet, ses déploiements et sa stack en place. */
export interface Application {
  projectId: string;
  /** Du plus récent au plus ancien ; vide pour une stack sans déploiement récent. */
  deployments: Deployment[];
  /** La stack en place ; absente si elle n'existe pas (ou plus). */
  stack?: Stack;
}

const createdMs = (d: Deployment | undefined) => (d?.createdAt ? timestampMs(d.createdAt) : 0);

/**
 * Regroupe les déploiements par projet, avec la stack en place de chacun ;
 * l'application la plus récemment déployée d'abord. Une stack sans
 * déploiement chargé (plus ancien que la page) a aussi sa ligne : elle reste
 * supprimable.
 */
export function groupByProject(deployments: Iterable<Deployment>, stacks: Stacks): Application[] {
  const byProject = new Map<string, Deployment[]>();
  for (const d of deployments) {
    const list = byProject.get(d.projectId) ?? [];
    list.push(d);
    byProject.set(d.projectId, list);
  }
  if (stacks.known) {
    for (const projectId of stacks.byProject.keys()) {
      if (!byProject.has(projectId)) byProject.set(projectId, []);
    }
  }
  return [...byProject.entries()]
    .map(([projectId, list]) => ({
      projectId,
      deployments: list.sort((a, b) => createdMs(b) - createdMs(a) || b.id.localeCompare(a.id)),
      stack: stacks.known ? stacks.byProject.get(projectId) : undefined,
    }))
    .sort(
      (a, b) =>
        createdMs(b.deployments[0]) - createdMs(a.deployments[0]) ||
        a.projectId.localeCompare(b.projectId)
    );
}
