import { timestampDate } from '@bufbuild/protobuf/wkt';
import { Code, ConnectError } from '@connectrpc/connect';
import { ArrowUpRight, Globe, LoaderCircle, PlugZap, Trash2 } from 'lucide-react';
import { useCallback, useId, useMemo, useState, type FormEvent } from 'react';

import { Button } from '../../components/ui/Button';
import { fieldClass } from '../../components/ui/field';
import { PageHeader } from '../../components/ui/PageHeader';
import { RequestError } from '../../components/ui/RequestError';
import { UserRole } from '../../gen/forgenet/v1/auth_pb';
import type { Domain } from '../../gen/forgenet/v1/domain_pb';
import { describeError } from '../../lib/api/errors';
import { useRequest } from '../../lib/useRequest';
import { useApi, useCurrentUser } from '../auth/session';

const added = new Intl.DateTimeFormat('fr-FR', { dateStyle: 'medium' });

const byHost = (a: Domain, b: Domain) => a.host.localeCompare(b.host);

/** Vrai quand le serveur répond que le reverse proxy n'est pas activé. */
function isProxyDisabled(err: unknown): boolean {
  const error = ConnectError.from(err);
  return error.code === Code.FailedPrecondition && error.rawMessage.includes('proxy.enabled');
}

export function DomainsPage() {
  const api = useApi();
  const isAdmin = useCurrentUser().role === UserRole.ADMIN;
  const load = useCallback(
    async (signal: AbortSignal) => (await api.domains.listDomains({}, { signal })).domains,
    [api]
  );
  const [request, retry] = useRequest(load);

  // Les ajouts et retraits faits sur cette page s'appliquent à la liste
  // chargée, sans la recharger ; un nouveau chargement les reprend par id.
  const [edits, setEdits] = useState<{ added: Domain[]; removed: string[] }>({
    added: [],
    removed: [],
  });
  const domains = useMemo(() => {
    if (request.status !== 'success') {
      return [];
    }
    const byId = new Map(request.data.map((d) => [d.id, d]));
    for (const d of edits.added) {
      byId.set(d.id, d);
    }
    for (const id of edits.removed) {
      byId.delete(id);
    }
    return [...byId.values()].sort(byHost);
  }, [request, edits]);

  return (
    <>
      <PageHeader
        title="Domaines"
        description="Les noms de domaine servis par le reverse proxy, chacun relié au service d’une application déployée. Le certificat HTTPS est délivré puis renouvelé automatiquement."
      />
      <div className="mt-8">
        {request.status === 'loading' && <DomainsSkeleton />}
        {request.status === 'error' &&
          (isProxyDisabled(request.error) ? (
            <ProxyDisabled message={describeError(request.error)} />
          ) : (
            <RequestError error={request.error} onRetry={retry} />
          ))}
        {request.status === 'success' && (
          <div className="grid items-start gap-8 lg:grid-cols-[minmax(0,1.6fr)_minmax(0,1fr)]">
            <section aria-labelledby="domains-list" className="animate-rise [animation-delay:60ms]">
              <h2 id="domains-list" className="text-lg font-semibold">
                {domains.length === 0
                  ? 'Aucun domaine'
                  : domains.length === 1
                    ? '1 domaine servi'
                    : `${domains.length} domaines servis`}
              </h2>
              <DomainList
                domains={domains}
                canRemove={isAdmin}
                onRemoved={(id) => setEdits((e) => ({ ...e, removed: [...e.removed, id] }))}
              />
            </section>
            <aside aria-labelledby="domains-add" className="animate-rise [animation-delay:120ms]">
              {isAdmin ? (
                <AddDomainForm
                  onAdded={(domain) => setEdits((e) => ({ ...e, added: [...e.added, domain] }))}
                />
              ) : (
                <p
                  id="domains-add"
                  className="rounded-lg border border-iron-800 bg-iron-900 px-5 py-4 text-sm text-iron-300"
                >
                  Seul un administrateur peut ajouter ou retirer un domaine.
                </p>
              )}
            </aside>
          </div>
        )}
      </div>
    </>
  );
}

function DomainList({
  domains,
  canRemove,
  onRemoved,
}: {
  domains: Domain[];
  canRemove: boolean;
  onRemoved: (id: string) => void;
}) {
  if (domains.length === 0) {
    return (
      <p className="mt-4 rounded-lg border border-dashed border-iron-700 px-5 py-10 text-center text-iron-300">
        Aucun domaine n’est encore relié à une application.
        {canRemove && ' Ajoutez le premier avec le formulaire.'}
      </p>
    );
  }
  return (
    <ul
      aria-label="Domaines servis"
      className="mt-4 divide-y divide-iron-800 overflow-hidden rounded-lg border border-iron-800"
    >
      {domains.map((d) => (
        <DomainRow key={d.id} domain={d} canRemove={canRemove} onRemoved={onRemoved} />
      ))}
    </ul>
  );
}

function DomainRow({
  domain,
  canRemove,
  onRemoved,
}: {
  domain: Domain;
  canRemove: boolean;
  onRemoved: (id: string) => void;
}) {
  const api = useApi();
  const [confirming, setConfirming] = useState(false);
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const errorId = useId();

  async function remove() {
    setPending(true);
    setError(undefined);
    try {
      await api.domains.removeDomain({ domainId: domain.id });
      onRemoved(domain.id);
    } catch (err) {
      setError(describeError(err));
      setPending(false);
    }
  }

  return (
    <li className="px-4 py-4 transition-colors hover:bg-iron-900/60">
      <div className="flex flex-wrap items-start justify-between gap-x-6 gap-y-3">
        <div className="min-w-0">
          <a
            href={`https://${domain.host}`}
            target="_blank"
            rel="noopener noreferrer"
            className="group inline-flex max-w-full items-center gap-1.5 rounded font-mono text-iron-100 hover:text-gold-300 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gold-400"
          >
            <span className="truncate">{domain.host}</span>
            <ArrowUpRight
              aria-hidden="true"
              className="size-4 shrink-0 text-iron-500 transition-colors group-hover:text-gold-300"
            />
            <span className="sr-only">(ouvre un nouvel onglet)</span>
          </a>
          <p className="mt-1 font-mono text-sm text-iron-300">
            {domain.projectId} / {domain.service}:{domain.port}
          </p>
          {domain.createdAt && (
            <p className="mt-1 text-xs text-iron-300">
              Ajouté le {added.format(timestampDate(domain.createdAt))}
            </p>
          )}
        </div>
        {canRemove && !confirming && (
          <Button
            variant="quiet"
            onClick={() => setConfirming(true)}
            aria-label={`Retirer ${domain.host}`}
            className="shrink-0"
          >
            <Trash2 aria-hidden="true" className="size-4" />
            Retirer
          </Button>
        )}
      </div>
      {confirming && (
        <div
          role="group"
          aria-label={`Retirer ${domain.host}`}
          aria-describedby={error ? errorId : undefined}
          className="mt-3 rounded-md border border-danger-400/40 bg-danger-950 px-3 py-3"
        >
          <p className="text-sm text-iron-100">
            Retirer <span className="font-mono">{domain.host}</span> ? Le site ne sera plus servi
            par ce serveur.
          </p>
          {error && (
            <p id={errorId} role="alert" className="mt-2 text-sm text-danger-400">
              {error}
            </p>
          )}
          <div className="mt-3 flex flex-wrap gap-3">
            <Button variant="danger" onClick={() => void remove()} disabled={pending}>
              {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
              {pending ? 'Retrait…' : 'Confirmer le retrait'}
            </Button>
            <Button
              variant="quiet"
              onClick={() => {
                setConfirming(false);
                setError(undefined);
              }}
              disabled={pending}
            >
              Annuler
            </Button>
          </div>
        </div>
      )}
    </li>
  );
}

interface FormValues {
  host: string;
  projectId: string;
  service: string;
  port: string;
}

const emptyForm: FormValues = { host: '', projectId: '', service: '', port: '' };

/** Vérifications locales ; le serveur valide le reste (forme du nom, service déployé). */
function checkForm(v: FormValues): Partial<Record<keyof FormValues, string>> {
  const problems: Partial<Record<keyof FormValues, string>> = {};
  if (!v.host.trim()) problems.host = 'Saisissez le nom de domaine.';
  if (!v.projectId.trim()) problems.projectId = 'Saisissez l’identifiant du projet.';
  if (!v.service.trim()) problems.service = 'Saisissez le service du fichier Compose.';
  const port = Number(v.port);
  if (!/^\d+$/.test(v.port.trim()) || port < 1 || port > 65535) {
    problems.port = 'Saisissez un port de 1 à 65535.';
  }
  return problems;
}

function AddDomainForm({ onAdded }: { onAdded: (domain: Domain) => void }) {
  const api = useApi();
  const [values, setValues] = useState<FormValues>(emptyForm);
  const [problems, setProblems] = useState<Partial<Record<keyof FormValues, string>>>({});
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const [success, setSuccess] = useState<string>();
  const ids = {
    host: useId(),
    projectId: useId(),
    service: useId(),
    port: useId(),
    error: useId(),
  };

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const found = checkForm(values);
    setProblems(found);
    setSuccess(undefined);
    if (Object.keys(found).length > 0) {
      setError(undefined);
      return;
    }
    setPending(true);
    setError(undefined);
    try {
      const res = await api.domains.addDomain({
        host: values.host.trim(),
        projectId: values.projectId.trim(),
        service: values.service.trim(),
        port: Number(values.port),
      });
      if (res.domain) {
        onAdded(res.domain);
        setSuccess(
          `${res.domain.host} est relié à ${res.domain.service}. Le certificat HTTPS est demandé : il arrive dès que le domaine pointe vers ce serveur.`
        );
      }
      setValues(emptyForm);
    } catch (err) {
      setError(describeError(err));
    } finally {
      setPending(false);
    }
  }

  const set = (key: keyof FormValues) => (e: { target: { value: string } }) =>
    setValues((v) => ({ ...v, [key]: e.target.value }));

  const input = (key: keyof FormValues) => ({
    id: ids[key],
    value: values[key],
    onChange: set(key),
    'aria-invalid': problems[key] ? true : undefined,
    'aria-describedby': problems[key] ? `${ids[key]}-hint` : undefined,
    className: `${fieldClass} font-mono ${problems[key] ? 'border-danger-400' : 'border-iron-700'}`,
  });

  const hint = (key: keyof FormValues) =>
    problems[key] && (
      <p id={`${ids[key]}-hint`} className="mt-1.5 text-sm text-danger-400">
        {problems[key]}
      </p>
    );

  return (
    <form
      noValidate
      onSubmit={submit}
      aria-labelledby="domains-add"
      aria-describedby={error ? ids.error : undefined}
      className="rounded-lg border border-iron-800 bg-iron-900 px-5 py-5"
    >
      <h2 id="domains-add" className="flex items-center gap-2 text-lg font-semibold">
        <Globe aria-hidden="true" className="size-5 text-gold-400" />
        Ajouter un domaine
      </h2>
      <p className="mt-2 text-sm text-iron-300">
        Faites d’abord pointer le domaine vers l’adresse publique de ce serveur (enregistrement A ou
        AAAA) : l’autorité de certification doit pouvoir le joindre.
      </p>

      <div className="mt-5">
        <label htmlFor={ids.host} className="text-sm font-medium">
          Nom de domaine
        </label>
        <input
          {...input('host')}
          name="host"
          inputMode="url"
          autoComplete="off"
          autoCapitalize="none"
          spellCheck={false}
          placeholder="app.exemple.com"
        />
        {hint('host')}
      </div>

      <div className="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-1 xl:grid-cols-2">
        <div>
          <label htmlFor={ids.projectId} className="text-sm font-medium">
            Projet
          </label>
          <input
            {...input('projectId')}
            name="project"
            autoComplete="off"
            autoCapitalize="none"
            spellCheck={false}
            placeholder="boutique"
          />
          {hint('projectId')}
        </div>
        <div>
          <label htmlFor={ids.service} className="text-sm font-medium">
            Service
          </label>
          <input
            {...input('service')}
            name="service"
            autoComplete="off"
            autoCapitalize="none"
            spellCheck={false}
            placeholder="web"
          />
          {hint('service')}
        </div>
      </div>

      <div className="mt-4">
        <label htmlFor={ids.port} className="text-sm font-medium">
          Port du service
        </label>
        <input
          {...input('port')}
          name="port"
          inputMode="numeric"
          autoComplete="off"
          placeholder="8080"
        />
        {problems.port ? (
          hint('port')
        ) : (
          <p className="mt-1.5 text-xs text-iron-300">
            Le port où le service écoute dans son conteneur, pas un port publié sur l’hôte.
          </p>
        )}
      </div>

      {error && (
        <p
          id={ids.error}
          role="alert"
          className="mt-5 rounded-md border border-danger-400/40 bg-danger-950 px-3 py-2.5 text-sm text-iron-100"
        >
          {error}
        </p>
      )}
      <p
        role="status"
        className="empty:hidden mt-5 rounded-md border border-ok-400/40 px-3 py-2.5 text-sm text-iron-100"
      >
        {success}
      </p>

      <Button type="submit" disabled={pending} className="mt-5 w-full">
        {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
        {pending ? 'Ajout…' : 'Ajouter le domaine'}
      </Button>
    </form>
  );
}

function ProxyDisabled({ message }: { message: string }) {
  return (
    <div className="animate-rise flex max-w-2xl items-start gap-4 rounded-lg border border-iron-700 bg-iron-900 px-5 py-6">
      <span className="flex size-10 shrink-0 items-center justify-center rounded-md bg-iron-850 text-warn-400">
        <PlugZap aria-hidden="true" className="size-5" />
      </span>
      <div>
        <p className="font-semibold">Le reverse proxy n’est pas activé sur ce serveur</p>
        <p className="mt-1 text-sm text-iron-300">{message}</p>
        <p className="mt-3 text-sm text-iron-300">
          Une fois Caddy installé, activez la section <code className="font-mono">proxy</code> de la
          configuration, puis redémarrez le serveur.
        </p>
      </div>
    </div>
  );
}

function DomainsSkeleton() {
  return (
    <div
      role="status"
      aria-live="polite"
      className="space-y-px overflow-hidden rounded-lg border border-iron-800 lg:max-w-[62%]"
    >
      <span className="sr-only">Chargement des domaines…</span>
      {[0, 1, 2].map((row) => (
        <div key={row} className="space-y-2 bg-iron-900/50 px-4 py-4">
          <span className="block h-3 w-48 animate-pulse rounded bg-iron-800" />
          <span className="block h-3 w-32 animate-pulse rounded bg-iron-800" />
        </div>
      ))}
    </div>
  );
}
