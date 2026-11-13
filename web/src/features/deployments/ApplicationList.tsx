import { timestampDate } from '@bufbuild/protobuf/wkt';
import { ChevronDown, LoaderCircle, ScrollText, CircleStop, Trash2, WifiOff } from 'lucide-react';
import { useId, useState, type ReactNode } from 'react';

import { Button } from '../../components/ui/Button';
import { fieldClass } from '../../components/ui/field';
import { DeploymentState, type Deployment } from '../../gen/forgenet/v1/deploy_pb';
import { describeError } from '../../lib/api/errors';
import { useApi, useCurrentUser } from '../auth/session';
import type { Application } from './group';
import { useDeploymentWatch } from './live';
import { LogPanel } from './LogPanel';
import { DeploymentBadge, DeploymentProgress } from './Progress';
import { isActive, isCancellable } from './state';
import { ago, elapsed, fullDate, useNow } from './time';

function created(d: Deployment): Date | undefined {
  return d.createdAt ? timestampDate(d.createdAt) : undefined;
}

export function ApplicationList({
  applications,
  onUpdate,
}: {
  applications: Application[];
  onUpdate: (deployment: Deployment) => void;
}) {
  if (applications.length === 0) {
    return (
      <p className="mt-4 rounded-lg border border-dashed border-iron-700 px-5 py-10 text-center text-iron-300">
        Aucune application n’a encore été déployée sur ce serveur. Collez un fichier Compose dans le
        formulaire pour lancer la première.
      </p>
    );
  }
  return (
    <ul aria-label="Applications" className="mt-4 space-y-3">
      {applications.map((app) => (
        <ApplicationRow key={app.projectId} app={app} onUpdate={onUpdate} />
      ))}
    </ul>
  );
}

function ApplicationRow({
  app,
  onUpdate,
}: {
  app: Application;
  onUpdate: (d: Deployment) => void;
}) {
  const [latest, ...history] = app.deployments;
  const active = isActive(latest.state);
  const me = useCurrentUser();
  const now = useNow(active);
  const watch = useDeploymentWatch(latest.id, active, onUpdate);
  const [logs, setLogs] = useState<{ id: string; follow: boolean } | undefined>(() =>
    active ? { id: latest.id, follow: true } : undefined
  );
  const [showHistory, setShowHistory] = useState(false);
  const historyId = useId();
  const when = created(latest);

  // Un nouveau déploiement du projet ouvre son journal en direct. L'état est
  // ajusté pendant le rendu, comme le recommande React, plutôt que par un effet.
  const [seen, setSeen] = useState(latest.id);
  if (seen !== latest.id) {
    setSeen(latest.id);
    if (isActive(latest.state)) {
      setLogs({ id: latest.id, follow: true });
    }
  }

  const author =
    latest.triggeredBy === `user:${me.id}`
      ? 'par vous'
      : latest.triggeredBy.startsWith('user:')
        ? 'par un autre compte'
        : '';

  return (
    <li className="rounded-lg border border-iron-800 bg-iron-900/40 px-4 py-4 transition-colors hover:border-iron-700">
      <div className="flex flex-wrap items-start justify-between gap-x-6 gap-y-2">
        <div className="min-w-0">
          <h3 className="truncate font-mono text-base font-semibold text-iron-100">
            {app.projectId}
          </h3>
          <p className="mt-1 text-xs text-iron-300">
            {active && when
              ? `En cours depuis ${elapsed(now - when.getTime())}`
              : when && (
                  <time dateTime={when.toISOString()} title={fullDate(when)}>
                    {ago(when, now)}
                  </time>
                )}
            {author && ` · ${author}`}
          </p>
        </div>
        <div aria-live="polite" className="flex items-center gap-3">
          {watch.link === 'reconnecting' && (
            <span className="inline-flex items-center gap-1.5 text-xs text-warn-400">
              <WifiOff aria-hidden="true" className="size-3.5" />
              Reconnexion…
            </span>
          )}
          <DeploymentBadge state={latest.state} />
        </div>
      </div>

      {active && <DeploymentProgress state={latest.state} />}
      {!active && latest.state !== DeploymentState.SUCCESS && latest.errorMessage && (
        <p
          className={`mt-3 rounded-md border px-3 py-2 font-mono text-[0.8125rem] break-words ${
            latest.state === DeploymentState.CANCELLED
              ? 'border-iron-700 bg-iron-900 text-iron-300'
              : 'border-danger-400/30 bg-danger-950 text-iron-100'
          }`}
        >
          {latest.errorMessage}
        </p>
      )}
      {watch.link === 'error' && (
        <p role="alert" className="mt-3 text-sm text-danger-400">
          Suivi en direct interrompu : {watch.error}
        </p>
      )}

      <div className="mt-3 flex flex-wrap items-center gap-2">
        <ToggleButton
          pressed={logs?.id === latest.id}
          onClick={() =>
            setLogs((l) =>
              l?.id === latest.id ? undefined : { id: latest.id, follow: isActive(latest.state) }
            )
          }
        >
          <ScrollText aria-hidden="true" className="size-4" />
          Journal
        </ToggleButton>
        {history.length > 0 && (
          <button
            type="button"
            aria-expanded={showHistory}
            aria-controls={historyId}
            onClick={() => setShowHistory((v) => !v)}
            className="inline-flex min-h-11 items-center gap-1.5 rounded-md px-3 text-sm text-iron-300 transition-colors hover:bg-iron-850 hover:text-iron-100 focus-visible:outline-2 focus-visible:outline-gold-400"
          >
            <ChevronDown
              aria-hidden="true"
              className={`size-4 transition-transform duration-200 ${showHistory ? 'rotate-180' : ''}`}
            />
            Historique ({history.length})
          </button>
        )}
        <span className="grow" />
        {isCancellable(latest.state) && <CancelButton deployment={latest} onUpdate={onUpdate} />}
      </div>

      {logs?.id === latest.id && (
        <LogPanel
          key={latest.id}
          deploymentId={latest.id}
          follow={logs.follow}
          project={app.projectId}
        />
      )}

      {showHistory && (
        <ol
          id={historyId}
          aria-label={`Déploiements précédents de ${app.projectId}`}
          className="mt-3 divide-y divide-iron-800 border-t border-iron-800"
        >
          {history.map((d) => (
            <HistoryItem
              key={d.id}
              deployment={d}
              project={app.projectId}
              now={now}
              open={logs?.id === d.id}
              onToggle={() =>
                setLogs((l) => (l?.id === d.id ? undefined : { id: d.id, follow: false }))
              }
            />
          ))}
        </ol>
      )}

      <DeleteStack projectId={app.projectId} busy={active} />
    </li>
  );
}

function HistoryItem({
  deployment: d,
  project,
  now,
  open,
  onToggle,
}: {
  deployment: Deployment;
  project: string;
  now: number;
  open: boolean;
  onToggle: () => void;
}) {
  const when = created(d);
  return (
    <li className="py-2.5">
      <div className="flex flex-wrap items-center justify-between gap-x-4 gap-y-1">
        <div className="flex min-w-0 items-center gap-3">
          <DeploymentBadge state={d.state} />
          {when && (
            <time
              dateTime={when.toISOString()}
              title={fullDate(when)}
              className="text-xs text-iron-300"
            >
              {ago(when, now)}
            </time>
          )}
        </div>
        <ToggleButton pressed={open} onClick={onToggle} compact>
          <ScrollText aria-hidden="true" className="size-3.5" />
          Journal{when && <span className="sr-only"> du déploiement du {fullDate(when)}</span>}
        </ToggleButton>
      </div>
      {d.errorMessage && (
        <p className="mt-1 truncate font-mono text-xs text-iron-300" title={d.errorMessage}>
          {d.errorMessage}
        </p>
      )}
      {open && <LogPanel deploymentId={d.id} follow={false} project={project} />}
    </li>
  );
}

function ToggleButton({
  pressed,
  onClick,
  compact,
  children,
}: {
  pressed: boolean;
  onClick: () => void;
  compact?: boolean;
  children: ReactNode;
}) {
  return (
    <button
      type="button"
      aria-pressed={pressed}
      onClick={onClick}
      className={`inline-flex items-center gap-1.5 rounded-md border px-3 text-sm transition-colors focus-visible:outline-2 focus-visible:outline-gold-400 ${
        compact ? 'min-h-11 text-xs' : 'min-h-11'
      } ${
        pressed
          ? 'border-gold-500/60 bg-gold-950 text-gold-300'
          : 'border-iron-700 text-iron-100 hover:border-iron-500 hover:bg-iron-850'
      }`}
    >
      {children}
    </button>
  );
}

function CancelButton({
  deployment,
  onUpdate,
}: {
  deployment: Deployment;
  onUpdate: (d: Deployment) => void;
}) {
  const api = useApi();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();

  async function cancel() {
    setPending(true);
    setError(undefined);
    try {
      const res = await api.deployments.cancelDeployment({ deploymentId: deployment.id });
      if (res.deployment) onUpdate(res.deployment);
    } catch (err) {
      setError(describeError(err));
    } finally {
      setPending(false);
    }
  }

  return (
    <>
      {error && (
        <span role="alert" className="text-sm text-danger-400">
          {error}
        </span>
      )}
      <Button variant="quiet" onClick={() => void cancel()} disabled={pending}>
        {pending ? (
          <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />
        ) : (
          <CircleStop aria-hidden="true" className="size-4" />
        )}
        {pending ? 'Annulation…' : 'Annuler le déploiement'}
      </Button>
    </>
  );
}

/**
 * Suppression de la stack en service : conteneurs et réseaux, et volumes sur
 * demande. Supprimer les volumes efface les données : il faut alors retaper
 * le nom du projet.
 */
function DeleteStack({ projectId, busy }: { projectId: string; busy: boolean }) {
  const api = useApi();
  const [open, setOpen] = useState(false);
  const [withVolumes, setWithVolumes] = useState(false);
  const [typed, setTyped] = useState('');
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const [done, setDone] = useState<string>();
  const ids = { volumes: useId(), typed: useId(), error: useId() };
  const confirmed = !withVolumes || typed.trim() === projectId;

  async function remove() {
    setPending(true);
    setError(undefined);
    try {
      await api.deployments.deleteStack({ projectId, deleteVolumes: withVolumes });
      setDone(
        withVolumes
          ? 'Stack supprimée, volumes compris : les conteneurs, réseaux et données de l’application sont effacés.'
          : 'Stack supprimée : conteneurs et réseaux retirés. Les volumes, donc les données, sont gardés.'
      );
      setOpen(false);
    } catch (err) {
      setError(describeError(err));
    } finally {
      setPending(false);
    }
  }

  function close() {
    setOpen(false);
    setWithVolumes(false);
    setTyped('');
    setError(undefined);
  }

  return (
    <>
      <p
        role="status"
        className="empty:hidden mt-3 rounded-md border border-ok-400/40 px-3 py-2 text-sm text-iron-100"
      >
        {done}
      </p>
      {!open ? (
        <div className="mt-3 flex flex-wrap items-center justify-end gap-3 border-t border-iron-800 pt-3">
          {busy && (
            <span className="text-xs text-iron-300">
              Disponible une fois le déploiement terminé.
            </span>
          )}
          <Button
            variant="quiet"
            onClick={() => {
              setOpen(true);
              setDone(undefined);
            }}
            disabled={busy}
            aria-label={`Supprimer la stack ${projectId}`}
          >
            <Trash2 aria-hidden="true" className="size-4" />
            Supprimer la stack
          </Button>
        </div>
      ) : (
        <div
          role="group"
          aria-label={`Supprimer la stack ${projectId}`}
          aria-describedby={error ? ids.error : undefined}
          className="mt-3 rounded-md border border-danger-400/40 bg-danger-950 px-3 py-3"
        >
          <p className="text-sm text-iron-100">
            Arrêter et supprimer la stack <span className="font-mono">{projectId}</span> ? Ses
            conteneurs et réseaux sont retirés ; l’historique des déploiements est gardé.
          </p>
          <div className="mt-3 flex items-start gap-2.5">
            <input
              id={ids.volumes}
              type="checkbox"
              checked={withVolumes}
              onChange={(e) => {
                setWithVolumes(e.target.checked);
                setTyped('');
              }}
              className="mt-0.5 size-5 shrink-0 accent-danger-400"
            />
            <label htmlFor={ids.volumes} className="text-sm text-iron-100">
              Supprimer aussi les volumes, c’est-à-dire{' '}
              <strong>les données de l’application</strong>. C’est irréversible.
            </label>
          </div>
          {withVolumes && (
            <div className="mt-3">
              <label htmlFor={ids.typed} className="text-sm text-iron-100">
                Tapez <span className="font-mono font-semibold">{projectId}</span> pour confirmer
              </label>
              <input
                id={ids.typed}
                value={typed}
                onChange={(e) => setTyped(e.target.value)}
                autoComplete="off"
                autoCapitalize="none"
                spellCheck={false}
                className={`${fieldClass} border-iron-700 font-mono`}
              />
            </div>
          )}
          {error && (
            <p id={ids.error} role="alert" className="mt-2 text-sm text-danger-400">
              {error}
            </p>
          )}
          <div className="mt-3 flex flex-wrap gap-3">
            <Button variant="danger" onClick={() => void remove()} disabled={pending || !confirmed}>
              {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
              {pending
                ? 'Suppression…'
                : withVolumes
                  ? 'Supprimer la stack et ses données'
                  : 'Supprimer la stack'}
            </Button>
            <Button variant="quiet" onClick={close} disabled={pending}>
              Annuler
            </Button>
          </div>
        </div>
      )}
    </>
  );
}
