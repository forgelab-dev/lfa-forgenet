import { timestampDate } from '@bufbuild/protobuf/wkt';
import { ArrowDown, Check, Copy, LoaderCircle, WifiOff } from 'lucide-react';
import { useEffect, useLayoutEffect, useRef, useState } from 'react';

import type { DeploymentLogEntry } from '../../gen/forgenet/v1/deploy_pb';
import { LogStream } from '../../gen/forgenet/v1/container_pb';
import { useDeploymentLogs } from './live';
import { clock } from './time';

/** Distance au bas en deçà de laquelle le journal suit les nouvelles lignes. */
const stickThreshold = 24;

function lineTone(entry: DeploymentLogEntry): string {
  if (entry.step === 'result') {
    if (entry.message.startsWith('Déploiement annulé')) return 'font-semibold text-iron-100';
    return entry.stream === LogStream.STDERR
      ? 'font-semibold text-danger-400'
      : 'font-semibold text-ok-400';
  }
  if (entry.message.startsWith('journal tronqué')) return 'text-warn-400';
  if (/\berror\b|\berreur\b|\bunhealthy\b|denied/i.test(entry.message)) return 'text-danger-400';
  // Les lignes de ForgeNet (sortie standard) passent devant la progression de Compose.
  return entry.stream === LogStream.STDOUT ? 'text-iron-100' : 'text-iron-300';
}

/**
 * Journal d'exécution d'un déploiement, en direct tant qu'il avance. Il suit
 * les nouvelles lignes, sauf si l'on remonte pour lire : un bouton ramène
 * alors en bas.
 */
export function LogPanel({
  deploymentId,
  follow,
  project,
}: {
  deploymentId: string;
  follow: boolean;
  project: string;
}) {
  const { lines, link, error } = useDeploymentLogs(deploymentId, follow);
  const box = useRef<HTMLDivElement>(null);
  const [stick, setStick] = useState(true);
  const [copied, setCopied] = useState(false);

  useLayoutEffect(() => {
    const el = box.current;
    if (el && stick) {
      el.scrollTop = el.scrollHeight;
    }
  }, [lines, stick]);

  useEffect(() => {
    if (!copied) return;
    const timer = setTimeout(() => setCopied(false), 2000);
    return () => clearTimeout(timer);
  }, [copied]);

  function onScroll() {
    const el = box.current;
    if (el) {
      setStick(el.scrollHeight - el.scrollTop - el.clientHeight < stickThreshold);
    }
  }

  async function copy() {
    const text = lines
      .map((l) => `${l.emittedAt ? clock(timestampDate(l.emittedAt)) : ''}  ${l.message}`)
      .join('\n');
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
    } catch {
      // Presse-papiers refusé (page non sécurisée) : rien à signaler de plus.
    }
  }

  return (
    <div className="mt-4 overflow-hidden rounded-md border border-iron-800 bg-iron-950">
      <div className="flex flex-wrap items-center justify-between gap-2 border-b border-iron-800 px-3 py-1.5">
        <p className="flex items-center gap-2 text-xs text-iron-300">
          {link === 'connecting' && (
            <>
              <LoaderCircle aria-hidden="true" className="size-3.5 animate-spin" />
              Ouverture du journal…
            </>
          )}
          {link === 'live' && follow && (
            <>
              <span
                aria-hidden="true"
                className="size-1.5 animate-pulse rounded-full bg-gold-400 motion-reduce:animate-none"
              />
              En direct
            </>
          )}
          {link === 'reconnecting' && (
            <>
              <WifiOff aria-hidden="true" className="size-3.5 text-warn-400" />
              Liaison coupée, reconnexion…
            </>
          )}
          {(link === 'done' || (link === 'live' && !follow)) && (
            <>
              {lines.length} ligne{lines.length > 1 ? 's' : ''}
            </>
          )}
          {link === 'error' && <span className="text-danger-400">{error}</span>}
        </p>
        <button
          type="button"
          onClick={() => void copy()}
          disabled={lines.length === 0}
          className="inline-flex min-h-11 items-center gap-1.5 rounded px-2 text-xs text-iron-300 transition-colors hover:bg-iron-850 hover:text-iron-100 focus-visible:outline-2 focus-visible:outline-gold-400 disabled:opacity-50"
        >
          {copied ? (
            <Check aria-hidden="true" className="size-3.5 text-ok-400" />
          ) : (
            <Copy aria-hidden="true" className="size-3.5" />
          )}
          {copied ? 'Copié' : 'Copier'}
        </button>
      </div>
      <div className="relative">
        {/* Un lecteur d'écran n'annonce pas chaque ligne : l'état du déploiement l'est déjà. */}
        <div
          ref={box}
          onScroll={onScroll}
          role="log"
          aria-live="off"
          aria-label={`Journal du déploiement de ${project}`}
          tabIndex={0}
          className="max-h-80 overflow-y-auto px-3 py-2 font-mono text-[0.8125rem] leading-relaxed focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-gold-400"
        >
          {lines.length === 0 && link !== 'connecting' && (
            <p className="py-2 text-iron-500">
              {follow
                ? 'Le journal commence dès que le déploiement démarre.'
                : 'Aucune ligne enregistrée pour ce déploiement.'}
            </p>
          )}
          {lines.map((l) => (
            <p
              key={l.emittedAt ? `${l.emittedAt.seconds}.${l.emittedAt.nanos}` : l.message}
              className="flex gap-3 whitespace-pre-wrap [overflow-wrap:anywhere]"
            >
              <span aria-hidden="true" className="shrink-0 text-iron-500 select-none">
                {l.emittedAt ? clock(timestampDate(l.emittedAt)) : ''}
              </span>
              <span className={lineTone(l)}>{l.message}</span>
            </p>
          ))}
        </div>
        {!stick && (
          <button
            type="button"
            onClick={() => setStick(true)}
            className="absolute right-3 bottom-3 inline-flex min-h-11 items-center gap-1.5 rounded-md border border-iron-700 bg-iron-900 px-3 text-xs font-semibold text-iron-100 shadow-lg transition-colors hover:border-iron-500"
          >
            <ArrowDown aria-hidden="true" className="size-3.5" />
            Suivre les nouvelles lignes
          </button>
        )}
      </div>
    </div>
  );
}
