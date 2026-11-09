import { useCallback, useEffect, useState } from 'react';

export type RequestState<T> =
  { status: 'loading' } | { status: 'error'; error: unknown } | { status: 'success'; data: T };

/**
 * Lance un appel au montage et garde son état (chargement, erreur, réponse).
 * `run` doit être stable (useCallback) : un nouveau `run` relance l'appel.
 * L'appel en cours est annulé au démontage.
 */
export function useRequest<T>(
  run: (signal: AbortSignal) => Promise<T>
): [RequestState<T>, () => void] {
  const [state, setState] = useState<RequestState<T>>({ status: 'loading' });
  const [attempt, setAttempt] = useState(0);

  useEffect(() => {
    const abort = new AbortController();
    run(abort.signal).then(
      (data) => {
        if (!abort.signal.aborted) {
          setState({ status: 'success', data });
        }
      },
      (error: unknown) => {
        if (!abort.signal.aborted) {
          setState({ status: 'error', error });
        }
      }
    );
    return () => abort.abort();
  }, [run, attempt]);

  const retry = useCallback(() => {
    setState({ status: 'loading' });
    setAttempt((n) => n + 1);
  }, []);

  return [state, retry];
}
