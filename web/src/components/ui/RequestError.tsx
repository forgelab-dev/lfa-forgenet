import { RotateCw } from 'lucide-react';

import { describeError } from '../../lib/api/errors';
import { Button } from './Button';

interface Props {
  error: unknown;
  onRetry: () => void;
}

/** Échec d'un chargement : la raison lisible et de quoi réessayer. */
export function RequestError({ error, onRetry }: Props) {
  return (
    <div
      role="alert"
      className="flex flex-col gap-4 rounded-lg border border-danger-400/40 bg-danger-950 px-5 py-4 sm:flex-row sm:items-center sm:justify-between"
    >
      <p className="text-sm text-iron-100">{describeError(error)}</p>
      <Button variant="quiet" onClick={onRetry} className="shrink-0">
        <RotateCw aria-hidden="true" className="size-4" />
        Réessayer
      </Button>
    </div>
  );
}
