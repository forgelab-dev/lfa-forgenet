import { timestampDate } from '@bufbuild/protobuf/wkt';
import {
  CircleCheck,
  KeyRound,
  LoaderCircle,
  Mail,
  MailX,
  ServerCog,
  ShieldAlert,
  TriangleAlert,
} from 'lucide-react';
import {
  useCallback,
  useId,
  useRef,
  useState,
  type FormEvent,
  type ReactNode,
  type RefObject,
} from 'react';

import { Button } from '../../components/ui/Button';
import { PageHeader } from '../../components/ui/PageHeader';
import { RequestError } from '../../components/ui/RequestError';
import { fieldClass } from '../../components/ui/field';
import { UserRole } from '../../gen/forgenet/v1/auth_pb';
import { MailSource, MailTLS, type GetMailStatusResponse } from '../../gen/forgenet/v1/mail_pb';
import { describeError } from '../../lib/api/errors';
import { useRequest } from '../../lib/useRequest';
import { useApi, useCurrentUser } from '../auth/session';
import {
  canKeepPassword,
  checkRelayForm,
  defaultPort,
  formFromStatus,
  requestFromForm,
  type RelayForm,
  type RelayProblems,
} from './relay';

const tlsLabel: Record<MailTLS, string> = {
  [MailTLS.MAIL_TLS_STARTTLS]: 'STARTTLS, exigé',
  [MailTLS.MAIL_TLS_IMPLICIT]: 'TLS dès la connexion',
  [MailTLS.MAIL_TLS_NONE]: 'aucune (relais local)',
  [MailTLS.MAIL_TLS_UNSPECIFIED]: 'inconnue',
};

const tlsChoices: { tls: MailTLS; label: string; detail: string }[] = [
  {
    tls: MailTLS.MAIL_TLS_STARTTLS,
    label: 'STARTTLS',
    detail: 'Port 587, le plus courant : la liaison passe en TLS avant tout envoi.',
  },
  {
    tls: MailTLS.MAIL_TLS_IMPLICIT,
    label: 'TLS dès la connexion',
    detail: 'Port 465, si votre fournisseur le demande.',
  },
  {
    tls: MailTLS.MAIL_TLS_NONE,
    label: 'Aucune',
    detail: 'Seulement pour un relais installé sur ce serveur (localhost).',
  },
];

const updatedAt = new Intl.DateTimeFormat('fr-FR', { dateStyle: 'long', timeStyle: 'short' });

/**
 * Réglages du serveur, réservés aux administrateurs. Le relais SMTP se règle
 * ici ou par install.sh ; la dernière action explicite l'emporte (ADR-0015).
 */
export function SettingsPage() {
  const user = useCurrentUser();
  return (
    <>
      <PageHeader
        title="Réglages"
        description="Ce que ce serveur sait faire au-delà de vos applications, et par quel moyen."
      />
      <div className="mt-8 max-w-3xl">
        {user.role === UserRole.ADMIN ? (
          <MailSection />
        ) : (
          <p className="flex items-start gap-3 rounded-lg border border-iron-800 px-5 py-5 text-iron-300">
            <ShieldAlert aria-hidden="true" className="mt-0.5 size-5 shrink-0 text-iron-300" />
            Les réglages du serveur sont réservés aux administrateurs.
          </p>
        )}
      </div>
    </>
  );
}

/** Ce que la section montre : les réglages, le formulaire, ou une confirmation. */
type Mode = 'view' | 'edit' | 'confirm-disable' | 'confirm-reset';

function MailSection() {
  const api = useApi();
  const load = useCallback((signal: AbortSignal) => api.mail.getMailStatus({}, { signal }), [api]);
  const [request, retry] = useRequest(load);
  // Réponse d'un enregistrement, d'une coupure ou d'un retour : elle remplace
  // l'état chargé, sans recharger la page.
  const [changed, setChanged] = useState<GetMailStatusResponse>();
  const [notice, setNotice] = useState<string>();
  const [mode, setMode] = useState<Mode>('view');
  const editButton = useRef<HTMLButtonElement>(null);

  const status = changed ?? (request.status === 'success' ? request.data : undefined);

  // Le formulaire ou la confirmation refermés, le focus revient au premier
  // bouton des actions plutôt que de se perdre en haut de la page.
  function backToView() {
    setMode('view');
    requestAnimationFrame(() => editButton.current?.focus());
  }

  function done(next: GetMailStatusResponse, message: string) {
    setChanged(next);
    setNotice(message);
    backToView();
  }

  const cancel = backToView;

  return (
    <section
      aria-labelledby="settings-mail"
      className="animate-rise rounded-lg border border-iron-800 bg-iron-900 px-5 py-5 sm:px-6"
    >
      <div className="flex flex-wrap items-baseline justify-between gap-x-4 gap-y-2">
        <h2 id="settings-mail" className="flex items-center gap-2 text-lg font-semibold">
          <Mail aria-hidden="true" className="size-5 text-gold-400" />
          Envoi d’e-mails
        </h2>
        {status && <SourceTag status={status} />}
      </div>
      <div className="mt-4">
        {request.status === 'loading' && (
          <div role="status" className="h-24 animate-pulse rounded-md bg-iron-850">
            <span className="sr-only">Chargement des réglages d’envoi…</span>
          </div>
        )}
        {request.status === 'error' && <RequestError error={request.error} onRetry={retry} />}
        {status && (
          <>
            {status.passwordUnreadable && <UnreadableWarning />}
            <p
              role="status"
              className="empty:hidden mb-4 flex items-start gap-2 rounded-md border border-ok-400/40 px-3 py-2.5 text-sm text-iron-100"
            >
              {notice && mode === 'view' && (
                <>
                  <CircleCheck aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-ok-400" />
                  <span>{notice}</span>
                </>
              )}
            </p>
            {mode === 'edit' ? (
              <RelayFormPanel
                status={status}
                onSaved={(next, to) =>
                  done(
                    next,
                    `Réglages enregistrés et appliqués : le relais a pris l’e-mail d’essai pour ${to}.`
                  )
                }
                onCancel={cancel}
              />
            ) : (
              <MailOverview status={status} />
            )}
            {mode === 'confirm-disable' && (
              <Confirm
                title="Couper l’envoi d’e-mails ?"
                detail="Les invitations, « Mot de passe oublié ? » et les avis aux comptes ne partiront plus. Les réglages restent ici pour rouvrir l’envoi."
                action="Couper l’envoi"
                run={async () => (await api.mail.disableMail({})).status}
                onDone={(next) => done(next, 'Envoi coupé.')}
                onCancel={cancel}
              />
            )}
            {mode === 'confirm-reset' && (
              <Confirm
                title="Revenir aux réglages du serveur ?"
                detail="Les réglages enregistrés ici sont effacés, mot de passe compris. Ceux d’install.sh s’appliquent de nouveau, tout de suite."
                action="Revenir aux réglages du serveur"
                run={async () => (await api.mail.resetMailSettings({})).status}
                onDone={(next) =>
                  done(
                    next,
                    next.enabled
                      ? 'Les réglages d’install.sh s’appliquent de nouveau.'
                      : 'Réglages effacés : install.sh n’en avait pas, l’envoi est éteint.'
                  )
                }
                onCancel={cancel}
              />
            )}
            {mode === 'view' && (
              <Actions
                status={status}
                editButton={editButton}
                onEdit={() => {
                  setNotice(undefined);
                  setMode('edit');
                }}
                onDisable={() => {
                  setNotice(undefined);
                  setMode('confirm-disable');
                }}
                onReset={() => {
                  setNotice(undefined);
                  setMode('confirm-reset');
                }}
              />
            )}
          </>
        )}
      </div>
    </section>
  );
}

/** D'où viennent les réglages : la page, ou install.sh sur le serveur. */
function SourceTag({ status }: { status: GetMailStatusResponse }) {
  if (status.source === MailSource.DASHBOARD) {
    const when = status.updatedAt ? updatedAt.format(timestampDate(status.updatedAt)) : undefined;
    return (
      <p className="text-xs text-iron-300">
        Réglé ici{status.updatedBy && <> par {status.updatedBy}</>}
        {when && <>, le {when}</>}
      </p>
    );
  }
  return (
    <p className="flex items-center gap-1.5 text-xs text-iron-300">
      <ServerCog aria-hidden="true" className="size-3.5" />
      Réglé sur le serveur (install.sh)
    </p>
  );
}

function UnreadableWarning() {
  return (
    <p
      role="alert"
      className="mb-4 flex items-start gap-2 rounded-md border border-warn-400/50 px-3 py-2.5 text-sm text-iron-100"
    >
      <TriangleAlert aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-warn-400" />
      <span>
        Le mot de passe du relais enregistré ici ne s’ouvre plus : la clé du serveur
        (FORGENET_SECRET_KEY) a changé. Les réglages d’install.sh s’appliquent en attendant.
        Enregistrez de nouveau le relais avec son mot de passe.
      </span>
    </p>
  );
}

function MailOverview({ status }: { status: GetMailStatusResponse }) {
  if (status.enabled) {
    return (
      <>
        <RelayDetails status={status} />
        <TestMail />
      </>
    );
  }
  if (status.disabled) {
    return (
      <>
        <p className="flex items-start gap-2 text-iron-300">
          <MailX aria-hidden="true" className="mt-0.5 size-4 shrink-0" />
          <span>
            L’envoi est coupé depuis cette page : ForgeNet n’envoie aucun e-mail. Les réglages
            ci-dessous restent pour le rouvrir.
          </span>
        </p>
        <div className="mt-4 opacity-75">
          <RelayDetails status={status} />
        </div>
      </>
    );
  }
  return <NoRelay editable={status.editable} />;
}

function RelayDetails({ status }: { status: GetMailStatusResponse }) {
  return (
    <dl className="grid gap-x-8 gap-y-3 text-sm sm:grid-cols-2">
      <Row label="Relais" value={`${status.host}:${status.port}`} mono />
      <Row label="Protection de la liaison" value={tlsLabel[status.tls]} />
      <Row label="Expéditeur" value={status.from} mono />
      <Row
        label="Identifiants"
        value={
          status.authenticated ? `${status.username}, mot de passe gardé sur le serveur` : 'aucun'
        }
      />
    </dl>
  );
}

function NoRelay({ editable }: { editable: boolean }) {
  return (
    <>
      <p className="flex items-start gap-2 text-iron-300">
        <MailX aria-hidden="true" className="mt-0.5 size-4 shrink-0" />
        <span>
          Aucun relais SMTP n’est réglé : ForgeNet n’envoie aucun e-mail. Tout le reste marche, et
          un mot de passe perdu se remplace sur le serveur.
        </span>
      </p>
      <p className="mt-4 text-sm text-iron-300">
        {editable
          ? 'Réglez ici le relais de votre fournisseur de messagerie, ou relancez l’installation sur le serveur. Le mot de passe du relais y est demandé, jamais passé en argument :'
          : 'Pour l’activer, relancez l’installation avec le relais de votre fournisseur de messagerie. Le mot de passe du relais est demandé, jamais passé en argument :'}
      </p>
      <pre className="mt-3 overflow-x-auto rounded-md bg-iron-950 px-3 py-2.5 font-mono text-xs text-iron-100">
        {
          'sudo bash install.sh \\\n  --smtp-host smtp.exemple.org \\\n  --smtp-user forgenet@exemple.org'
        }
      </pre>
    </>
  );
}

function Actions({
  status,
  editButton,
  onEdit,
  onDisable,
  onReset,
}: {
  status: GetMailStatusResponse;
  editButton: RefObject<HTMLButtonElement | null>;
  onEdit: () => void;
  onDisable: () => void;
  onReset: () => void;
}) {
  if (!status.editable) {
    return (
      <p className="mt-6 flex items-start gap-2 border-t border-iron-800 pt-5 text-sm text-iron-300">
        <KeyRound aria-hidden="true" className="mt-0.5 size-4 shrink-0" />
        <span>
          Ce serveur n’a pas de clé FORGENET_SECRET_KEY pour garder le mot de passe du relais : les
          réglages se changent par install.sh, qui la crée à la prochaine mise à jour.
        </span>
      </p>
    );
  }
  const editLabel = status.disabled
    ? 'Rouvrir l’envoi'
    : status.host
      ? 'Modifier le relais'
      : 'Régler un relais';
  return (
    <div className="mt-6 flex flex-wrap gap-2 border-t border-iron-800 pt-5">
      <Button ref={editButton} variant={status.enabled ? 'quiet' : 'primary'} onClick={onEdit}>
        {editLabel}
      </Button>
      {status.enabled && (
        <Button variant="quiet" onClick={onDisable}>
          Couper l’envoi
        </Button>
      )}
      {/* Des réglages illisibles restent en base : ce bouton les efface aussi. */}
      {(status.source === MailSource.DASHBOARD || status.passwordUnreadable) && (
        <Button variant="quiet" onClick={onReset}>
          Revenir aux réglages du serveur
        </Button>
      )}
    </div>
  );
}

function Confirm({
  title,
  detail,
  action,
  run,
  onDone,
  onCancel,
}: {
  title: string;
  detail: string;
  action: string;
  run: () => Promise<GetMailStatusResponse | undefined>;
  onDone: (next: GetMailStatusResponse) => void;
  onCancel: () => void;
}) {
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const titleId = useId();

  async function confirm() {
    setPending(true);
    setError(undefined);
    try {
      const next = await run();
      if (next) {
        onDone(next);
      }
    } catch (err) {
      setError(describeError(err));
      setPending(false);
    }
  }

  return (
    <div
      role="group"
      aria-labelledby={titleId}
      className="mt-6 rounded-md border border-iron-700 bg-iron-950 px-4 py-4"
    >
      <p id={titleId} className="text-sm font-semibold text-iron-100">
        {title}
      </p>
      <p className="mt-1 text-sm text-iron-300">{detail}</p>
      {error && (
        <p
          role="alert"
          className="mt-3 rounded-md border border-danger-400/40 bg-danger-950 px-3 py-2.5 text-sm text-iron-100"
        >
          {error}
        </p>
      )}
      <div className="mt-4 flex flex-wrap gap-2">
        {/* autoFocus : la confirmation s'ouvre sur demande, le focus y entre. */}
        <Button onClick={() => void confirm()} disabled={pending} autoFocus>
          {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
          {action}
        </Button>
        <Button variant="quiet" onClick={onCancel} disabled={pending}>
          Annuler
        </Button>
      </div>
    </div>
  );
}

function RelayFormPanel({
  status,
  onSaved,
  onCancel,
}: {
  status: GetMailStatusResponse;
  onSaved: (next: GetMailStatusResponse, to: string) => void;
  onCancel: () => void;
}) {
  const api = useApi();
  const keepable = canKeepPassword(status);
  const [values, setValues] = useState<RelayForm>(() => formFromStatus(status));
  const [problems, setProblems] = useState<RelayProblems>({});
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const ids = {
    host: useId(),
    port: useId(),
    tls: useId(),
    username: useId(),
    password: useId(),
    from: useId(),
    note: useId(),
    error: useId(),
  };

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const found = checkRelayForm(values, keepable);
    setProblems(found);
    setError(undefined);
    if (Object.keys(found).length > 0) {
      return;
    }
    setPending(true);
    try {
      const res = await api.mail.saveMailSettings(requestFromForm(values, keepable));
      if (res.status) {
        onSaved(res.status, res.to);
      }
    } catch (err) {
      setError(describeError(err));
      setPending(false);
    }
  }

  function chooseTLS(tls: MailTLS) {
    setValues((v) => {
      // Le port suit la protection tant qu'il est l'habituel de l'ancienne.
      const usual = v.port === String(defaultPort[v.tls]);
      return { ...v, tls, port: usual ? String(defaultPort[tls]) : v.port };
    });
  }

  const set = (key: keyof RelayForm) => (e: { target: { value: string } }) => {
    const value = e.target.value;
    setValues((v) => ({ ...v, [key]: value }));
  };

  const field = (key: Exclude<keyof RelayForm, 'tls'>) => ({
    id: ids[key],
    name: `smtp-${key}`,
    value: values[key],
    onChange: set(key),
    autoCapitalize: 'none',
    spellCheck: false,
    'aria-invalid': problems[key] ? true : undefined,
    'aria-describedby': problems[key] ? `${ids[key]}-hint` : `${ids[key]}-help`,
    className: `${fieldClass} ${problems[key] ? 'border-danger-400' : 'border-iron-700'}`,
  });

  const hint = (key: keyof RelayForm, help?: ReactNode) =>
    problems[key] ? (
      <p id={`${ids[key]}-hint`} className="mt-1.5 text-sm text-danger-400">
        {problems[key]}
      </p>
    ) : (
      help && (
        <p id={`${ids[key]}-help`} className="mt-1.5 text-xs text-iron-300">
          {help}
        </p>
      )
    );

  return (
    <form
      noValidate
      onSubmit={submit}
      aria-label="Réglages du relais SMTP"
      aria-describedby={error ? `${ids.note} ${ids.error}` : ids.note}
      className="rounded-md border border-iron-700 bg-iron-950/40 px-4 py-4 sm:px-5"
    >
      <p id={ids.note} className="flex gap-2 text-sm text-iron-300">
        <Mail aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-gold-400" />
        <span>
          À l’enregistrement, un e-mail d’essai part à votre adresse par ce relais. Les réglages ne
          s’appliquent que s’il l’accepte ; sinon, ceux en place restent.
        </span>
      </p>

      <div className="mt-4 grid gap-4 sm:grid-cols-[minmax(0,1fr)_8rem]">
        <div>
          <label htmlFor={ids.host} className="text-sm font-medium">
            Relais
          </label>
          {/* autoFocus : le formulaire s'ouvre sur demande, le focus y entre. */}
          <input {...field('host')} placeholder="smtp.exemple.org" autoFocus />
          {hint('host', 'Le serveur SMTP de votre fournisseur de messagerie.')}
        </div>
        <div>
          <label htmlFor={ids.port} className="text-sm font-medium">
            Port
          </label>
          <input {...field('port')} inputMode="numeric" />
          {hint('port')}
        </div>
      </div>

      <fieldset className="mt-5" aria-describedby={problems.tls ? `${ids.tls}-hint` : undefined}>
        <legend className="text-sm font-medium">Protection de la liaison</legend>
        <div className="mt-1.5 grid gap-2">
          {tlsChoices.map((c) => (
            <label
              key={c.tls}
              className={`flex min-h-11 cursor-pointer items-start gap-3 rounded-md border px-3 py-2.5 transition-colors has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-gold-400/35 ${
                values.tls === c.tls
                  ? 'border-gold-400 bg-iron-950'
                  : 'border-iron-700 hover:border-iron-500'
              }`}
            >
              <input
                type="radio"
                name="smtp-tls"
                value={c.tls}
                checked={values.tls === c.tls}
                onChange={() => chooseTLS(c.tls)}
                className="mt-1 size-4 accent-gold-400"
              />
              <span>
                <span className="block text-sm font-semibold text-iron-100">{c.label}</span>
                <span className="block text-sm text-iron-300">{c.detail}</span>
              </span>
            </label>
          ))}
        </div>
        {problems.tls && (
          <p id={`${ids.tls}-hint`} className="mt-1.5 text-sm text-danger-400">
            {problems.tls}
          </p>
        )}
      </fieldset>

      <div className="mt-5 grid gap-4 sm:grid-cols-2">
        <div>
          <label htmlFor={ids.username} className="text-sm font-medium">
            Identifiant <span className="font-normal text-iron-300">(facultatif)</span>
          </label>
          <input {...field('username')} autoComplete="off" />
          {hint('username', 'Souvent l’adresse du compte d’envoi. Vide : sans authentification.')}
        </div>
        <div>
          <label htmlFor={ids.password} className="text-sm font-medium">
            Mot de passe du relais
          </label>
          <input
            {...field('password')}
            type="password"
            autoComplete="new-password"
            disabled={!values.username.trim()}
            placeholder={keepable && values.username.trim() ? '••••••••' : undefined}
            className={`${field('password').className} disabled:cursor-not-allowed disabled:opacity-55`}
          />
          {hint(
            'password',
            !values.username.trim()
              ? 'Sans identifiant, pas de mot de passe.'
              : keepable
                ? 'Vide : le mot de passe enregistré est gardé. Il n’est jamais réaffiché.'
                : 'Gardé chiffré sur le serveur, jamais réaffiché.'
          )}
        </div>
      </div>

      <div className="mt-4">
        <label htmlFor={ids.from} className="text-sm font-medium">
          Expéditeur
        </label>
        <input
          {...field('from')}
          type="email"
          inputMode="email"
          placeholder="forgenet@exemple.org"
        />
        {hint(
          'from',
          'L’adresse qui apparaît chez le destinataire ; votre relais doit l’accepter.'
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

      <div className="mt-5 flex flex-wrap gap-2">
        <Button type="submit" disabled={pending}>
          {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
          {pending ? 'Essai en cours…' : 'Enregistrer et essayer'}
        </Button>
        <Button type="button" variant="quiet" onClick={onCancel} disabled={pending}>
          Annuler
        </Button>
      </div>
    </form>
  );
}

function Row({ label, value, mono = false }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="min-w-0">
      <dt className="text-xs text-iron-300">{label}</dt>
      {/* Une adresse ou un hôte se coupe n'importe où ; une phrase, entre les mots. */}
      <dd className={`mt-0.5 text-iron-100 ${mono ? 'break-all font-mono' : 'break-words'}`}>
        {value}
      </dd>
    </div>
  );
}

type TestState =
  | { kind: 'idle' }
  | { kind: 'sending' }
  | { kind: 'sent'; to: string }
  | { kind: 'failed'; message: string };

function TestMail() {
  const api = useApi();
  const [state, setState] = useState<TestState>({ kind: 'idle' });

  async function send() {
    setState({ kind: 'sending' });
    try {
      const res = await api.mail.sendTestMail({});
      setState({ kind: 'sent', to: res.to });
    } catch (err) {
      setState({ kind: 'failed', message: describeError(err) });
    }
  }

  return (
    <div className="mt-6 border-t border-iron-800 pt-5">
      <p className="text-sm text-iron-300">
        Un e-mail d’essai part à l’adresse de votre compte. Le relais répond tout de suite : sa
        réponse dit quoi corriger en cas de refus.
      </p>
      {state.kind === 'failed' && (
        <p
          role="alert"
          className="mt-4 rounded-md border border-danger-400/40 bg-danger-950 px-3 py-2.5 text-sm text-iron-100"
        >
          {state.message}
        </p>
      )}
      <p
        role="status"
        className="empty:hidden mt-4 flex items-start gap-2 rounded-md border border-ok-400/40 px-3 py-2.5 text-sm text-iron-100"
      >
        {state.kind === 'sent' && (
          <>
            <CircleCheck aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-ok-400" />
            <span>
              Le relais a pris l’e-mail d’essai pour <span className="font-mono">{state.to}</span>.
              Vérifiez la boîte, et les indésirables s’il n’arrive pas.
            </span>
          </>
        )}
      </p>
      <Button
        variant="quiet"
        onClick={() => void send()}
        disabled={state.kind === 'sending'}
        className="mt-4"
      >
        {state.kind === 'sending' && (
          <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />
        )}
        {state.kind === 'sending' ? 'Envoi…' : 'Envoyer un e-mail d’essai'}
      </Button>
    </div>
  );
}
