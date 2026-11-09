import { ContainerState, type PortMapping } from '../../gen/forgenet/v1/container_pb';

export type Tone = 'ok' | 'warn' | 'danger' | 'idle';

/** Libellé et ton de chaque état Docker, dans les mots de l'administrateur. */
export const containerStates: Record<ContainerState, { label: string; tone: Tone }> = {
  [ContainerState.RUNNING]: { label: 'En marche', tone: 'ok' },
  [ContainerState.RESTARTING]: { label: 'Redémarre', tone: 'warn' },
  [ContainerState.PAUSED]: { label: 'En pause', tone: 'warn' },
  [ContainerState.CREATED]: { label: 'Créé', tone: 'idle' },
  [ContainerState.EXITED]: { label: 'Arrêté', tone: 'idle' },
  [ContainerState.DEAD]: { label: 'Hors service', tone: 'danger' },
  [ContainerState.UNSPECIFIED]: { label: 'Inconnu', tone: 'idle' },
};

/** « 8080 → 80/tcp » ; un port non publié s'écrit « 80/tcp ». */
export function formatPort(port: PortMapping): string {
  const target = `${port.containerPort}/${port.protocol || 'tcp'}`;
  return port.hostPort ? `${port.hostPort} → ${target}` : target;
}
