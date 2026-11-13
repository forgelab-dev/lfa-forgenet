import { Check, X } from 'lucide-react';

import { DeploymentState } from '../../gen/forgenet/v1/deploy_pb';
import { deploymentStates, progressSteps, type Tone } from './state';

const dot: Record<Tone, string> = {
  ok: 'bg-ok-400',
  warn: 'bg-warn-400',
  danger: 'bg-danger-400',
  idle: 'border border-iron-500 bg-transparent',
  active: 'bg-gold-400 animate-pulse motion-reduce:animate-none',
};

const text: Record<Tone, string> = {
  ok: 'text-ok-400',
  warn: 'text-warn-400',
  danger: 'text-danger-400',
  idle: 'text-iron-300',
  active: 'text-gold-300',
};

/** État d'un déploiement : une pastille (pulsante tant qu'il avance) et son libellé. */
export function DeploymentBadge({ state }: { state: DeploymentState }) {
  const { label, tone } = deploymentStates[state];
  return (
    <span className={`inline-flex items-center gap-2 text-sm whitespace-nowrap ${text[tone]}`}>
      <span aria-hidden="true" className={`size-2 rounded-full ${dot[tone]}`} />
      {label}
    </span>
  );
}

/**
 * Frise des étapes d'un déploiement : faites, en cours, à venir. La dernière
 * case dit l'issue une fois connue.
 */
export function DeploymentProgress({ state }: { state: DeploymentState }) {
  const current = progressSteps.findIndex((s) => s.state === state);
  // Un état final : toutes les étapes qui l'ont précédé sont faites.
  const reached = current === -1 ? progressSteps.length : current;
  const outcome =
    state === DeploymentState.SUCCESS
      ? 'ok'
      : state === DeploymentState.FAILED
        ? 'danger'
        : state === DeploymentState.CANCELLED
          ? 'idle'
          : undefined;

  return (
    <ol aria-label="Étapes du déploiement" className="mt-4 grid grid-cols-5 gap-1.5">
      {progressSteps.map((step, i) => {
        const done = i < reached;
        const now = i === current;
        return (
          <li key={step.state} aria-current={now ? 'step' : undefined} className="min-w-0">
            <span
              aria-hidden="true"
              className={`block h-1 rounded-full ${
                done
                  ? 'bg-gold-500'
                  : now
                    ? 'animate-pulse bg-gold-400 motion-reduce:animate-none'
                    : 'bg-iron-800'
              }`}
            />
            <span
              className={`mt-1.5 block truncate text-xs ${now ? 'font-semibold text-gold-300' : done ? 'text-iron-300' : 'text-iron-500'}`}
            >
              {step.label}
              <span className="sr-only">
                {done ? ' : faite' : now ? ' : en cours' : ' : à venir'}
              </span>
            </span>
          </li>
        );
      })}
      <li className="min-w-0">
        <span
          aria-hidden="true"
          className={`block h-1 rounded-full ${
            outcome === 'ok' ? 'bg-ok-400' : outcome === 'danger' ? 'bg-danger-400' : 'bg-iron-800'
          }`}
        />
        <span
          className={`mt-1.5 flex items-center gap-1 truncate text-xs ${
            outcome === 'ok'
              ? 'font-semibold text-ok-400'
              : outcome === 'danger'
                ? 'font-semibold text-danger-400'
                : 'text-iron-500'
          }`}
        >
          {outcome === 'ok' && <Check aria-hidden="true" className="size-3.5" />}
          {outcome === 'danger' && <X aria-hidden="true" className="size-3.5" />}
          {outcome ? deploymentStates[state].label : 'Résultat'}
        </span>
      </li>
    </ol>
  );
}
