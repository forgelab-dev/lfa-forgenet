import { DeploymentState } from '../../gen/forgenet/v1/deploy_pb';

export type Tone = 'ok' | 'warn' | 'danger' | 'idle' | 'active';

interface StateInfo {
  label: string;
  tone: Tone;
}

/** Libellé et ton de chaque état, dans les mots de l'administrateur. */
export const deploymentStates: Record<DeploymentState, StateInfo> = {
  [DeploymentState.UNSPECIFIED]: { label: 'Inconnu', tone: 'idle' },
  [DeploymentState.PENDING]: { label: 'En attente', tone: 'active' },
  [DeploymentState.BUILDING]: { label: 'Vérification', tone: 'active' },
  [DeploymentState.DEPLOYING]: { label: 'Mise en service', tone: 'active' },
  [DeploymentState.HEALTH_CHECKING]: { label: 'Contrôle des services', tone: 'active' },
  [DeploymentState.SUCCESS]: { label: 'Réussi', tone: 'ok' },
  [DeploymentState.FAILED]: { label: 'Échec', tone: 'danger' },
  [DeploymentState.CANCELLED]: { label: 'Annulé', tone: 'idle' },
  [DeploymentState.ROLLING_BACK]: { label: 'Retour en arrière', tone: 'active' },
  [DeploymentState.ROLLED_BACK]: { label: 'Revenu en arrière', tone: 'warn' },
};

/** Les étapes d'un déploiement Compose, dans l'ordre de la frise. */
export const progressSteps: { state: DeploymentState; label: string }[] = [
  { state: DeploymentState.PENDING, label: 'En attente' },
  { state: DeploymentState.BUILDING, label: 'Vérification' },
  { state: DeploymentState.DEPLOYING, label: 'Mise en service' },
  { state: DeploymentState.HEALTH_CHECKING, label: 'Contrôle' },
];

const active = new Set([
  DeploymentState.PENDING,
  DeploymentState.BUILDING,
  DeploymentState.DEPLOYING,
  DeploymentState.HEALTH_CHECKING,
  DeploymentState.ROLLING_BACK,
]);

/** Vrai tant que le déploiement n'a pas atteint d'état final. */
export function isActive(state: DeploymentState): boolean {
  return active.has(state);
}

/** Le serveur n'annule qu'avant le contrôle des services. */
export function isCancellable(state: DeploymentState): boolean {
  return (
    state === DeploymentState.PENDING ||
    state === DeploymentState.BUILDING ||
    state === DeploymentState.DEPLOYING
  );
}
