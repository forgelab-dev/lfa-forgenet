import { Code, ConnectError } from '@connectrpc/connect';

/** Une coupure réseau ou un abonné décroché se rattrapent ; un refus, non. */
export function isTransient(err: unknown): boolean {
  const { code } = ConnectError.from(err);
  return (
    code === Code.Unavailable ||
    code === Code.Unknown ||
    code === Code.DeadlineExceeded ||
    code === Code.Internal ||
    code === Code.Aborted
  );
}

/** Attente avant la reconnexion n : 1, 2, 4, 8 puis 10 secondes. */
export function retryDelay(attempt: number): number {
  return Math.min(1000 * 2 ** attempt, 10_000);
}

/** Attend ms, ou moins si signal est annulé. */
export function sleep(ms: number, signal: AbortSignal): Promise<void> {
  return new Promise((resolve) => {
    const timer = setTimeout(resolve, ms);
    signal.addEventListener(
      'abort',
      () => {
        clearTimeout(timer);
        resolve();
      },
      { once: true }
    );
  });
}
