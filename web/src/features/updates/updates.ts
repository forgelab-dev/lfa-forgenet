import { timestampMs } from '@bufbuild/protobuf/wkt';
import { createContext, useContext } from 'react';

import {
  UpdateResult,
  UpdateStep,
  type GetUpdateStatusResponse,
  type UpdateRun,
} from '../../gen/forgenet/v1/update_pb';

/** État des mises à jour, pour les administrateurs ; null pour un membre. */
export interface Updates {
  /** Absent tant que le premier appel n'a pas répondu. */
  status?: GetUpdateStatusResponse;
  /** Relit l'état (sans appel à GitHub). */
  refresh: () => Promise<void>;
  /** Interroge les releases tout de suite. */
  check: () => Promise<void>;
  checking: boolean;
  checkError?: string;
}

export const UpdatesContext = createContext<Updates | null>(null);

export function useUpdates(): Updates | null {
  return useContext(UpdatesContext);
}

export function isFinished(run: UpdateRun | undefined): boolean {
  return (
    run !== undefined &&
    run.result !== UpdateResult.RUNNING &&
    run.result !== UpdateResult.UNSPECIFIED
  );
}

/** Clé d'une mise à jour : son heure de début, unique par installation. */
export function runKey(run: UpdateRun | undefined): string {
  return run?.startedAt ? String(timestampMs(run.startedAt)) : '';
}

/**
 * Étapes montrées pendant une mise à jour, dans l'ordre. Le serveur est
 * arrêté entre la sauvegarde et la santé : ces étapes-là ne se lisent qu'après
 * coup, d'où le libellé « redémarrage » pendant la coupure.
 */
export const progressSteps: { step: UpdateStep; label: string }[] = [
  { step: UpdateStep.PULL, label: 'Téléchargement de la nouvelle version' },
  { step: UpdateStep.BACKUP, label: 'Sauvegarde de la base et de la configuration' },
  { step: UpdateStep.SWITCH, label: 'Redémarrage sur la nouvelle version' },
  { step: UpdateStep.HEALTH, label: 'Vérification du serveur' },
];

const DISMISSED_KEY = 'forgenet.updates.dismissed';

/** Vrai si l'administrateur a fermé l'avis de cette mise à jour. */
export function isDismissed(key: string): boolean {
  try {
    return key !== '' && localStorage.getItem(DISMISSED_KEY) === key;
  } catch {
    return false;
  }
}

export function dismiss(key: string): void {
  try {
    localStorage.setItem(DISMISSED_KEY, key);
  } catch {
    // Stockage indisponible (navigation privée) : l'avis reviendra.
  }
}
