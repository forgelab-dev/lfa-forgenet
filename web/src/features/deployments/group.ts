import { timestampMs } from '@bufbuild/protobuf/wkt';

import type { Deployment } from '../../gen/forgenet/v1/deploy_pb';

/** Une application : un projet et ses déploiements. */
export interface Application {
  projectId: string;
  /** Du plus récent au plus ancien. */
  deployments: Deployment[];
}

const createdMs = (d: Deployment) => (d.createdAt ? timestampMs(d.createdAt) : 0);

/** Regroupe les déploiements par projet ; l'application la plus récemment déployée d'abord. */
export function groupByProject(deployments: Iterable<Deployment>): Application[] {
  const byProject = new Map<string, Deployment[]>();
  for (const d of deployments) {
    const list = byProject.get(d.projectId) ?? [];
    list.push(d);
    byProject.set(d.projectId, list);
  }
  return [...byProject.entries()]
    .map(([projectId, list]) => ({
      projectId,
      deployments: list.sort((a, b) => createdMs(b) - createdMs(a) || b.id.localeCompare(a.id)),
    }))
    .sort((a, b) => createdMs(b.deployments[0]) - createdMs(a.deployments[0]));
}
