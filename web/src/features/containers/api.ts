import type { Container } from '../../gen/forgenet/v1/container_pb';
import type { ApiClients } from '../../lib/api/client';

/** Taille de page maximale acceptée par le serveur. */
const PAGE_SIZE = 100;
/** Garde-fou : 20 pages, soit 2 000 conteneurs, bien au-delà d'un serveur réel. */
const MAX_PAGES = 20;

/**
 * Tous les conteneurs du serveur, arrêtés compris, page après page. Un serveur
 * en héberge quelques dizaines : tout charger d'un coup garde le tri et le
 * comptage simples, sans pagination à l'écran.
 */
export async function listAllContainers(
  api: ApiClients,
  signal: AbortSignal
): Promise<Container[]> {
  const containers: Container[] = [];
  let pageToken = '';
  for (let page = 0; page < MAX_PAGES; page++) {
    const res = await api.containers.listContainers(
      { includeStopped: true, pageSize: PAGE_SIZE, pageToken },
      { signal }
    );
    containers.push(...res.containers);
    if (!res.nextPageToken) {
      break;
    }
    pageToken = res.nextPageToken;
  }
  return containers;
}
