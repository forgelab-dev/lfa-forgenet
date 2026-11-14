import { timestampMs } from '@bufbuild/protobuf/wkt';
import {
  ArrowUpCircle,
  ArrowUpRight,
  Check,
  CircleAlert,
  CircleCheck,
  LoaderCircle,
  X,
} from 'lucide-react';
import { useCallback, useEffect, useId, useState } from 'react';

import { Button } from '../../components/ui/Button';
import {
  UpdateResult,
  UpdateStep,
  type GetUpdateStatusResponse,
  type UpdateRun,
} from '../../gen/forgenet/v1/update_pb';
import { describeError } from '../../lib/api/errors';
import { useNow } from '../../lib/live/useNow';
import { useApi } from '../auth/session';
import { dismiss, isDismissed, isFinished, progressSteps, runKey, useUpdates } from './updates';

const POLL_MS = 2000;
/** Au-delà, la mise à jour est anormalement longue : on le dit. */
const SLOW_MS = 10 * 60 * 1000;
const DAY_MS = 24 * 60 * 60 * 1000;

type Phase =
  | { kind: 'idle' }
  | { kind: 'confirm' }
  | { kind: 'starting' }
  /** ignoreKey : la mise à jour terminée d'avant, à ne pas prendre pour celle-ci. */
  | { kind: 'tracking'; ignoreKey: string }
  | { kind: 'finished'; run: UpdateRun };

/**
 * Bandeau des mises à jour, en tête de chaque page pour les administrateurs
 * (ADR-0011) : la version disponible, la confirmation, le suivi pendant la
 * coupure du serveur, puis le résultat. Il ne montre rien quand tout est à
 * jour.
 */
export function UpdateBanner() {
  const updates = useUpdates();
  const status = updates?.status;
  const [phase, setPhase] = useState<Phase>({ kind: 'idle' });
  // Stable : le suivi relance ses relectures quand ce rappel change.
  const onFinished = useCallback((run: UpdateRun) => setPhase({ kind: 'finished', run }), []);

  if (!updates || !status) {
    return null;
  }

  // Une mise à jour trouvée en cours au chargement (page rechargée pendant la
  // coupure) est suivie comme celle qu'on vient de lancer.
  const running = status.lastRun?.result === UpdateResult.RUNNING;
  const tracking: Phase | null =
    phase.kind === 'tracking'
      ? phase
      : phase.kind === 'idle' && running
        ? { kind: 'tracking', ignoreKey: '' }
        : null;

  if (tracking?.kind === 'tracking') {
    return <Progress status={status} ignoreKey={tracking.ignoreKey} onFinished={onFinished} />;
  }
  if (phase.kind === 'finished') {
    return (
      <Result
        run={phase.run}
        onClose={() => {
          dismiss(runKey(phase.run));
          setPhase({ kind: 'idle' });
          void updates.refresh();
        }}
      />
    );
  }
  return (
    <>
      <LastRunNotice status={status} />
      {status.available && (
        <Available
          status={status}
          phase={phase}
          onPhase={setPhase}
          ignoreKey={runKey(status.lastRun)}
        />
      )}
    </>
  );
}

function Available({
  status,
  phase,
  onPhase,
  ignoreKey,
}: {
  status: GetUpdateStatusResponse;
  phase: Phase;
  onPhase: (p: Phase) => void;
  ignoreKey: string;
}) {
  const api = useApi();
  const [error, setError] = useState<string>();
  const titleId = useId();
  const errorId = useId();
  const next = status.available?.version ?? '';
  const confirming = phase.kind === 'confirm' || phase.kind === 'starting';
  const starting = phase.kind === 'starting';

  async function start() {
    onPhase({ kind: 'starting' });
    setError(undefined);
    try {
      await api.updates.startUpdate({});
      onPhase({ kind: 'tracking', ignoreKey });
    } catch (err) {
      setError(describeError(err));
      onPhase({ kind: 'confirm' });
    }
  }

  return (
    <section
      aria-labelledby={titleId}
      className="mb-8 rounded-lg border border-gold-400/40 bg-gold-950 px-5 py-4"
    >
      <div className="flex flex-wrap items-start justify-between gap-x-6 gap-y-3">
        <div className="flex min-w-0 gap-3">
          <ArrowUpCircle aria-hidden="true" className="mt-0.5 size-5 shrink-0 text-gold-400" />
          <div className="min-w-0">
            <h2 id={titleId} className="font-semibold text-iron-100">
              ForgeNet <span className="font-mono">{next}</span> est disponible
            </h2>
            <p className="mt-1 text-sm text-iron-300">
              Ce serveur tourne en <span className="font-mono">{status.currentVersion}</span>.
              {status.available?.notesUrl && (
                <>
                  {' '}
                  <a
                    href={status.available.notesUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-0.5 rounded text-gold-300 underline-offset-2 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gold-400"
                  >
                    Lire les notes de version
                    <ArrowUpRight aria-hidden="true" className="size-3.5" />
                    <span className="sr-only">(ouvre un nouvel onglet)</span>
                  </a>
                </>
              )}
            </p>
          </div>
        </div>
        {status.supported && !confirming && (
          <Button onClick={() => onPhase({ kind: 'confirm' })} className="shrink-0">
            Mettre à jour
          </Button>
        )}
      </div>

      {!status.supported && (
        <p className="mt-3 text-sm text-iron-300">
          Cette installation n’a pas été faite par <span className="font-mono">install.sh</span> :
          mettez-la à jour sur le serveur, en suivant les notes de version.
        </p>
      )}

      {confirming && (
        <div
          role="group"
          aria-label={`Mettre à jour vers ${next}`}
          aria-describedby={error ? errorId : undefined}
          className="mt-4 border-t border-gold-400/25 pt-4"
        >
          <ul className="space-y-1.5 text-sm text-iron-100">
            <li className="flex gap-2">
              <Check aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-gold-400" />
              La base et la configuration sont sauvegardées avant tout changement.
            </li>
            <li className="flex gap-2">
              <Check aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-gold-400" />
              Le dashboard s’interrompt environ une minute ; les applications déployées ne sont pas
              redémarrées.
            </li>
            <li className="flex gap-2">
              <Check aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-gold-400" />
              En cas d’échec, la version {status.currentVersion} est remise en service
              automatiquement.
            </li>
          </ul>
          {error && (
            <p id={errorId} role="alert" className="mt-3 text-sm text-danger-400">
              {error}
            </p>
          )}
          <div className="mt-4 flex flex-wrap gap-3">
            <Button onClick={() => void start()} disabled={starting}>
              {starting && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
              {starting ? 'Lancement…' : `Installer la version ${next}`}
            </Button>
            <Button
              variant="quiet"
              onClick={() => {
                setError(undefined);
                onPhase({ kind: 'idle' });
              }}
              disabled={starting}
            >
              Annuler
            </Button>
          </div>
        </div>
      )}
    </section>
  );
}

/**
 * Suivi d'une mise à jour : relit l'état toutes les deux secondes. Pendant la
 * coupure, les appels échouent : c'est attendu, on l'affiche comme un
 * redémarrage, et on continue jusqu'au retour du serveur.
 */
function Progress({
  status,
  ignoreKey,
  onFinished,
}: {
  status: GetUpdateStatusResponse;
  ignoreKey: string;
  onFinished: (run: UpdateRun) => void;
}) {
  const api = useApi();
  const [run, setRun] = useState<UpdateRun | undefined>(
    runKey(status.lastRun) === ignoreKey ? undefined : status.lastRun
  );
  const [unreachable, setUnreachable] = useState(false);
  const [slow, setSlow] = useState(false);
  const target = run?.toVersion || status.available?.version || '';

  useEffect(() => {
    let active = true;
    let timer: ReturnType<typeof setTimeout>;
    const tick = async () => {
      try {
        const next = await api.updates.getUpdateStatus({});
        if (!active) {
          return;
        }
        setUnreachable(false);
        const last = next.lastRun;
        if (last && runKey(last) !== ignoreKey) {
          if (isFinished(last)) {
            onFinished(last);
            return;
          }
          setRun(last);
        }
      } catch {
        if (!active) {
          return;
        }
        setUnreachable(true);
      }
      timer = setTimeout(() => void tick(), POLL_MS);
    };
    timer = setTimeout(() => void tick(), POLL_MS);
    const slowTimer = setTimeout(() => setSlow(true), SLOW_MS);
    return () => {
      active = false;
      clearTimeout(timer);
      clearTimeout(slowTimer);
    };
  }, [api, ignoreKey, onFinished]);

  const rollingBack = run?.step === UpdateStep.ROLLBACK;
  // Serveur injoignable : il est arrêté, entre la sauvegarde et la santé.
  const current = unreachable ? 2 : stepIndex(run?.step);

  return (
    <section
      role="status"
      aria-live="polite"
      className="mb-8 rounded-lg border border-gold-400/40 bg-gold-950 px-5 py-4"
    >
      <h2 className="flex items-center gap-2 font-semibold text-iron-100">
        <LoaderCircle aria-hidden="true" className="size-5 animate-spin text-gold-400" />
        {rollingBack
          ? 'Échec : retour à la version précédente'
          : `Mise à jour vers ${target || 'la nouvelle version'}`}
      </h2>
      {rollingBack ? (
        <p className="mt-2 text-sm text-iron-300">{run?.reason}</p>
      ) : (
        <ol className="mt-3 space-y-2 text-sm">
          {progressSteps.map(({ label }, i) => (
            <li key={label} className="flex items-center gap-2">
              {i < current ? (
                <Check aria-hidden="true" className="size-4 shrink-0 text-ok-400" />
              ) : i === current ? (
                <LoaderCircle
                  aria-hidden="true"
                  className="size-4 shrink-0 animate-spin text-gold-400"
                />
              ) : (
                <span
                  aria-hidden="true"
                  className="size-4 shrink-0 rounded-full border border-iron-500"
                />
              )}
              <span className={i <= current ? 'text-iron-100' : 'text-iron-300'}>
                {label}
                <span className="sr-only">
                  {i < current ? ' : terminé' : i === current ? ' : en cours' : ' : à venir'}
                </span>
              </span>
            </li>
          ))}
        </ol>
      )}
      {unreachable && (
        <p className="mt-3 text-sm text-iron-300">
          Le serveur redémarre : le dashboard revient de lui-même, gardez cette page ouverte.
        </p>
      )}
      {slow && (
        <p className="mt-3 text-sm text-warn-400">
          C’est plus long que prévu. Sur le serveur,{' '}
          <span className="font-mono">docker logs forgenet-updater</span> montre où en est la mise à
          jour.
        </p>
      )}
    </section>
  );
}

function stepIndex(step: UpdateStep | undefined): number {
  switch (step) {
    case UpdateStep.BACKUP:
      return 1;
    case UpdateStep.SWITCH:
      return 2;
    case UpdateStep.HEALTH:
      return 3;
    default:
      return 0;
  }
}

/** Résultat d'une mise à jour qu'on vient de suivre. */
function Result({ run, onClose }: { run: UpdateRun; onClose: () => void }) {
  if (run.result === UpdateResult.SUCCEEDED) {
    return (
      <section
        role="status"
        className="mb-8 flex flex-wrap items-center justify-between gap-4 rounded-lg border border-ok-400/40 bg-iron-900 px-5 py-4"
      >
        <p className="flex items-center gap-2 text-iron-100">
          <CircleCheck aria-hidden="true" className="size-5 shrink-0 text-ok-400" />
          <span>
            ForgeNet <span className="font-mono">{run.toVersion}</span> est installé. Rechargez le
            dashboard pour passer à sa nouvelle interface.
          </span>
        </p>
        <Button onClick={() => window.location.reload()} className="shrink-0">
          Recharger le dashboard
        </Button>
      </section>
    );
  }
  return <Failure run={run} onClose={onClose} />;
}

/**
 * Avis sur la dernière mise à jour, au chargement : un échec reste affiché
 * une semaine, une réussite un jour, tant que l'administrateur ne l'a pas
 * fermé.
 */
function LastRunNotice({ status }: { status: GetUpdateStatusResponse }) {
  const now = useNow(60_000);
  const [closed, setClosed] = useState(false);
  const run = status.lastRun;
  if (closed || !run || !isFinished(run) || !run.finishedAt || isDismissed(runKey(run))) {
    return null;
  }
  const age = now - timestampMs(run.finishedAt);
  const close = () => {
    dismiss(runKey(run));
    setClosed(true);
  };
  if (run.result === UpdateResult.SUCCEEDED) {
    if (age > DAY_MS || run.toVersion !== status.currentVersion) {
      return null;
    }
    return (
      <section
        role="status"
        className="mb-8 flex items-start justify-between gap-4 rounded-lg border border-ok-400/40 bg-iron-900 px-5 py-3"
      >
        <p className="flex items-start gap-2 text-sm text-iron-100">
          <CircleCheck aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-ok-400" />
          <span>
            ForgeNet a été mis à jour de <span className="font-mono">{run.fromVersion}</span> en{' '}
            <span className="font-mono">{run.toVersion}</span>.
          </span>
        </p>
        <CloseButton onClick={close} />
      </section>
    );
  }
  if (age > 7 * DAY_MS) {
    return null;
  }
  return <Failure run={run} onClose={close} />;
}

function Failure({ run, onClose }: { run: UpdateRun; onClose: () => void }) {
  const rolledBack = run.result === UpdateResult.ROLLED_BACK;
  // Échec pendant le retour arrière : rien n'a pu être remis en place.
  const stranded = !rolledBack && run.step === UpdateStep.ROLLBACK;
  return (
    <section
      role="alert"
      className="mb-8 rounded-lg border border-danger-400/40 bg-danger-950 px-5 py-4"
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex min-w-0 gap-3">
          <CircleAlert aria-hidden="true" className="mt-0.5 size-5 shrink-0 text-danger-400" />
          <div className="min-w-0">
            <h2 className="font-semibold text-iron-100">
              {rolledBack
                ? `La mise à jour vers ${run.toVersion} a échoué`
                : `La mise à jour vers ${run.toVersion} n’a pas pu se faire`}
            </h2>
            <p className="mt-1 text-sm text-iron-100">
              {rolledBack
                ? `ForgeNet ${run.fromVersion} a été remis en service, avec ses données.`
                : stranded
                  ? 'Le retour à la version précédente a aussi échoué : une intervention sur le serveur est nécessaire.'
                  : `ForgeNet ${run.fromVersion} est toujours en service, rien n’a été modifié.`}
            </p>
            {run.reason && (
              <p className="mt-2 text-sm break-words text-iron-300">Raison : {run.reason}</p>
            )}
            {run.backup && (stranded || rolledBack) && (
              <p className="mt-1 text-sm break-all text-iron-300">
                Sauvegarde : <span className="font-mono">{run.backup}</span>
              </p>
            )}
          </div>
        </div>
        <CloseButton onClick={onClose} />
      </div>
    </section>
  );
}

function CloseButton({ onClick }: { onClick: () => void }) {
  return (
    <button
      type="button"
      onClick={onClick}
      aria-label="Fermer cet avis"
      className="-m-2 flex size-11 shrink-0 items-center justify-center rounded-md text-iron-300 hover:bg-iron-850 hover:text-iron-100 focus-visible:outline-2 focus-visible:outline-gold-400"
    >
      <X aria-hidden="true" className="size-4" />
    </button>
  );
}
