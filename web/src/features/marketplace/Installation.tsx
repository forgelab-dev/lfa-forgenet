import { ArrowRight, Check, Copy, Eye, EyeOff, KeyRound, ShieldCheck, WifiOff } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Link } from 'wouter';

import { Button } from '../../components/ui/Button';
import { DeploymentState, type Deployment } from '../../gen/forgenet/v1/deploy_pb';
import type { GeneratedSecret, Template } from '../../gen/forgenet/v1/marketplace_pb';
import { useDeploymentWatch } from '../deployments/live';
import { LogPanel } from '../deployments/LogPanel';
import { DeploymentBadge, DeploymentProgress } from '../deployments/Progress';
import { isActive } from '../deployments/state';
import { BackLink } from './parts';

/**
 * L'installation lancée : les mots de passe générés, montrés cette seule fois,
 * puis le déploiement suivi en direct jusqu'à son issue.
 */
export function Installation({
  template,
  initial,
  secrets,
}: {
  template: Template;
  initial: Deployment;
  secrets: GeneratedSecret[];
}) {
  const [deployment, setDeployment] = useState(initial);
  const active = isActive(deployment.state);
  const watch = useDeploymentWatch(deployment.id, active, setDeployment);
  // Fixé à l'ouverture : le journal suivi se ferme seul à la fin du
  // déploiement, et le rouvrir recopierait ses lignes.
  const [follow] = useState(() => isActive(initial.state));

  return (
    <div className="mt-8 grid max-w-4xl gap-6">
      {secrets.length > 0 && <SecretsOnce secrets={secrets} />}

      <section
        aria-labelledby="installation"
        className="animate-rise rounded-lg border border-iron-800 bg-iron-900/40 px-4 py-4 [animation-delay:60ms]"
      >
        <div className="flex flex-wrap items-start justify-between gap-x-6 gap-y-2">
          <div className="min-w-0">
            <h2 id="installation" className="truncate font-mono text-base font-semibold">
              {deployment.projectId}
            </h2>
            <p className="mt-1 text-xs text-iron-300">Installation de {template.name}</p>
          </div>
          <div aria-live="polite" className="flex items-center gap-3">
            {watch.link === 'reconnecting' && (
              <span className="inline-flex items-center gap-1.5 text-xs text-warn-400">
                <WifiOff aria-hidden="true" className="size-3.5" />
                Reconnexion…
              </span>
            )}
            <DeploymentBadge state={deployment.state} />
          </div>
        </div>

        <DeploymentProgress state={deployment.state} />
        {!active && deployment.state !== DeploymentState.SUCCESS && deployment.errorMessage && (
          <p className="mt-3 rounded-md border border-danger-400/30 bg-danger-950 px-3 py-2 font-mono text-[0.8125rem] break-words text-iron-100">
            {deployment.errorMessage}
          </p>
        )}
        {watch.link === 'error' && (
          <p role="alert" className="mt-3 text-sm text-danger-400">
            Suivi en direct interrompu : {watch.error}
          </p>
        )}
        <LogPanel deploymentId={deployment.id} follow={follow} project={deployment.projectId} />

        <div className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-2">
          <Link
            href="/deploiements"
            className="inline-flex min-h-11 items-center gap-1.5 text-sm font-medium text-gold-400 hover:text-gold-300"
          >
            Voir l’application dans Déploiements
            <ArrowRight aria-hidden="true" className="size-4" />
          </Link>
          <BackLink className="" />
        </div>
      </section>
    </div>
  );
}

/**
 * Les mots de passe tirés par le serveur. Ils ne sont dans aucune réponse
 * ultérieure (ADR-0019) : cette page est la seule chance de les noter.
 * « C'est noté » les retire de la page ; avant cela, fermer ou recharger
 * l'onglet demande confirmation.
 */
function SecretsOnce({ secrets }: { secrets: GeneratedSecret[] }) {
  const [noted, setNoted] = useState(false);

  useEffect(() => {
    if (noted) return;
    const warn = (event: BeforeUnloadEvent) => event.preventDefault();
    window.addEventListener('beforeunload', warn);
    return () => window.removeEventListener('beforeunload', warn);
  }, [noted]);

  if (noted) {
    return (
      <p className="flex items-center gap-2 text-sm text-iron-300">
        <ShieldCheck aria-hidden="true" className="size-4 text-ok-400" />
        Mots de passe effacés de cette page. Pour en changer un, faites-le dans l’application.
      </p>
    );
  }

  return (
    <section
      aria-labelledby="secrets"
      className="animate-rise rounded-lg border border-gold-500/70 bg-gold-950 px-5 py-5"
    >
      <h2 id="secrets" className="flex items-center gap-2 text-lg font-semibold text-gold-300">
        <KeyRound aria-hidden="true" className="size-5" />
        Mots de passe générés
      </h2>
      <p className="mt-2 max-w-prose text-sm text-iron-100">
        Notez-les maintenant, dans votre gestionnaire de mots de passe : ForgeNet ne les réaffichera
        pas, ni ici ni ailleurs.
      </p>
      <ul className="mt-4 space-y-2">
        {secrets.map((s) => (
          <SecretRow key={s.name} secret={s} />
        ))}
      </ul>
      <Button variant="quiet" onClick={() => setNoted(true)} className="mt-5">
        <Check aria-hidden="true" className="size-4" />
        C’est noté, effacer de la page
      </Button>
    </section>
  );
}

function SecretRow({ secret: s }: { secret: GeneratedSecret }) {
  const [shown, setShown] = useState(false);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    if (!copied) return;
    const timer = setTimeout(() => setCopied(false), 2000);
    return () => clearTimeout(timer);
  }, [copied]);

  async function copy() {
    try {
      await navigator.clipboard.writeText(s.value);
      setCopied(true);
    } catch {
      // Presse-papiers refusé (page non sécurisée) : « Afficher » reste possible.
      setShown(true);
    }
  }

  const label = s.label || s.name;
  return (
    <li className="flex flex-wrap items-center gap-x-4 gap-y-2 rounded-md border border-iron-700 bg-iron-950 px-3 py-2">
      <div className="min-w-0 flex-[1_1_12rem]">
        <p className="text-sm font-medium">{label}</p>
        <p className="font-mono text-xs text-iron-500">{s.name}</p>
      </div>
      <code className="min-w-0 flex-[2_1_14rem] font-mono text-sm break-all text-iron-100 select-all">
        {shown ? s.value : '•'.repeat(Math.min(s.value.length, 24))}
      </code>
      <div className="flex gap-1">
        <button
          type="button"
          onClick={() => setShown((v) => !v)}
          aria-pressed={shown}
          aria-label={`Afficher ${label}`}
          className="flex size-11 items-center justify-center rounded-md text-iron-300 transition-colors hover:bg-iron-850 hover:text-iron-100 focus-visible:outline-2 focus-visible:outline-gold-400"
        >
          {shown ? (
            <EyeOff aria-hidden="true" className="size-4" />
          ) : (
            <Eye aria-hidden="true" className="size-4" />
          )}
        </button>
        <button
          type="button"
          onClick={() => void copy()}
          aria-label={`Copier ${label}`}
          className="flex size-11 items-center justify-center rounded-md text-iron-300 transition-colors hover:bg-iron-850 hover:text-iron-100 focus-visible:outline-2 focus-visible:outline-gold-400"
        >
          {copied ? (
            <Check aria-hidden="true" className="size-4 text-ok-400" />
          ) : (
            <Copy aria-hidden="true" className="size-4" />
          )}
        </button>
      </div>
    </li>
  );
}
