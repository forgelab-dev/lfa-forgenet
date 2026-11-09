import type { ContainerState } from '../../gen/forgenet/v1/container_pb';
import { containerStates, type Tone } from './state';

const dot: Record<Tone, string> = {
  ok: 'bg-ok-400',
  warn: 'bg-warn-400',
  danger: 'bg-danger-400',
  idle: 'border border-iron-500 bg-transparent',
};

const text: Record<Tone, string> = {
  ok: 'text-ok-400',
  warn: 'text-warn-400',
  danger: 'text-danger-400',
  idle: 'text-iron-300',
};

/** État d'un conteneur : une pastille (pleine si actif, creuse sinon) et son libellé. */
export function StateBadge({ state }: { state: ContainerState }) {
  const { label, tone } = containerStates[state];
  return (
    <span className={`inline-flex items-center gap-2 text-sm whitespace-nowrap ${text[tone]}`}>
      <span aria-hidden="true" className={`size-2 rounded-full ${dot[tone]}`} />
      {label}
    </span>
  );
}
