import type { LiveStatus } from '../../lib/live/useLive';
import { useNow } from '../../lib/live/useNow';
import { formatClock } from './format';

interface Props {
  status: LiveStatus;
  /** Rythme des mesures en direct, ex. « toutes les 5 s ». */
  pace: string;
  /** Date de la dernière mesure reçue. */
  lastAt?: number;
}

const styles = {
  ok: 'border-iron-700 text-iron-100',
  warn: 'border-warn-400/50 text-warn-400',
  idle: 'border-iron-700 text-iron-300',
  danger: 'border-danger-400/50 bg-danger-950 text-danger-400',
};

const dots = {
  ok: 'bg-ok-400 motion-safe:animate-pulse',
  warn: 'bg-warn-400',
  idle: 'bg-iron-500',
  danger: 'bg-danger-400',
};

/** État de la liaison en direct : une seule pastille, toujours au même endroit. */
export function LinkBadge({ status, pace, lastAt }: Props) {
  const now = useNow();
  const seconds =
    status.state === 'reconnecting' ? Math.max(0, Math.ceil((status.retryAt - now) / 1000)) : 0;

  let tone: keyof typeof styles = 'idle';
  let label = 'Connexion…';
  let detail: string | undefined;
  switch (status.state) {
    case 'live':
      tone = 'ok';
      label = 'En direct';
      detail = pace;
      break;
    case 'reconnecting':
      tone = 'warn';
      label = 'Liaison coupée';
      detail = seconds > 0 ? `nouvel essai dans ${seconds} s` : 'nouvel essai…';
      break;
    case 'paused':
      label = 'En pause';
      detail = 'onglet caché';
      break;
    case 'error':
    case 'gone':
      tone = 'danger';
      label = 'Flux fermé';
      break;
  }
  const stale = status.state !== 'live' && lastAt !== undefined;

  return (
    <p
      role="status"
      aria-live="polite"
      className={`inline-flex min-h-9 items-center gap-2.5 rounded-full border bg-iron-900 px-3 text-sm whitespace-nowrap ${styles[tone]}`}
    >
      <span aria-hidden="true" className={`size-2 rounded-full ${dots[tone]}`} />
      {label}
      {detail && (
        <>
          <span aria-hidden="true" className="text-iron-500">
            ·
          </span>
          <span className="font-mono text-[0.8125rem]">{detail}</span>
        </>
      )}
      {stale && (
        <span className="font-mono text-[0.8125rem] text-iron-300">
          mesures de {formatClock(lastAt)}
        </span>
      )}
    </p>
  );
}
