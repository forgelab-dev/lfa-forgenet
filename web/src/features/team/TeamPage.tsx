import { timestampDate } from '@bufbuild/protobuf/wkt';
import {
  Check,
  CircleCheck,
  Clock,
  Copy,
  LoaderCircle,
  MailCheck,
  ShieldAlert,
  TriangleAlert,
  UserPlus,
} from 'lucide-react';
import { useCallback, useId, useMemo, useRef, useState, type FormEvent } from 'react';
import { Link } from 'wouter';

import { Button } from '../../components/ui/Button';
import { fieldClass } from '../../components/ui/field';
import { PageHeader } from '../../components/ui/PageHeader';
import { RequestError } from '../../components/ui/RequestError';
import { UserRole } from '../../gen/forgenet/v1/auth_pb';
import { MemberStatus, type InvitationResult, type Member } from '../../gen/forgenet/v1/team_pb';
import { describeError } from '../../lib/api/errors';
import { useRequest } from '../../lib/useRequest';
import { useApi, useCurrentUser } from '../auth/session';

const until = new Intl.DateTimeFormat('fr-FR', { dateStyle: 'long', timeStyle: 'short' });

/** Même règle que le serveur : le nom sert à se connecter et figure dans les e-mails. */
const usernamePattern = /^[A-Za-z0-9._-]{2,64}$/;

/**
 * Les comptes du serveur et leurs invitations, réservés aux administrateurs
 * (ADR-0014). Un invité choisit lui-même son mot de passe, par un lien à
 * usage unique : l'administrateur ne le connaît jamais.
 */
export function TeamPage() {
  const user = useCurrentUser();
  return (
    <>
      <PageHeader
        title="Équipe"
        description="Les comptes de ce serveur. Un membre invité choisit lui-même son mot de passe, par un lien valable 7 jours et une seule fois."
      />
      <div className="mt-8">
        {user.role === UserRole.ADMIN ? (
          <Team currentUserId={user.id} />
        ) : (
          <p className="flex max-w-3xl items-start gap-3 rounded-lg border border-iron-800 px-5 py-5 text-iron-300">
            <ShieldAlert aria-hidden="true" className="mt-0.5 size-5 shrink-0 text-iron-300" />
            L’équipe est gérée par les administrateurs du serveur.
          </p>
        )}
      </div>
    </>
  );
}

function Team({ currentUserId }: { currentUserId: string }) {
  const api = useApi();
  const load = useCallback((signal: AbortSignal) => api.team.listMembers({}, { signal }), [api]);
  const [request, retry] = useRequest(load);

  // Les invitations faites ici s'appliquent à la liste chargée, sans la
  // recharger : un compte nouveau s'ajoute en bas, un compte réinvité change
  // d'état sur place.
  const [updates, setUpdates] = useState<Member[]>([]);
  const members = useMemo(() => {
    if (request.status !== 'success') {
      return [];
    }
    const byId = new Map(request.data.members.map((m) => [m.id, m]));
    for (const m of updates) {
      byId.set(m.id, m);
    }
    return [...byId.values()];
  }, [request, updates]);
  const update = useCallback((m: Member) => setUpdates((all) => [...all, m]), []);

  if (request.status === 'loading') {
    return <TeamSkeleton />;
  }
  if (request.status === 'error') {
    return <RequestError error={request.error} onRetry={retry} />;
  }
  const mailEnabled = request.data.mailEnabled;
  const waiting = members.filter((m) => m.status !== MemberStatus.ACTIVE).length;

  return (
    <div className="grid items-start gap-8 lg:grid-cols-[minmax(0,1.6fr)_minmax(0,1fr)]">
      <section aria-labelledby="team-list" className="animate-rise [animation-delay:60ms]">
        <h2 id="team-list" className="text-lg font-semibold">
          {members.length === 1 ? '1 compte' : `${members.length} comptes`}
          {waiting > 0 && (
            <span className="font-normal text-iron-300">
              {' '}
              ·{' '}
              {waiting === 1
                ? '1 attend son mot de passe'
                : `${waiting} attendent leur mot de passe`}
            </span>
          )}
        </h2>
        <ul
          aria-label="Comptes du serveur"
          className="mt-4 divide-y divide-iron-800 overflow-hidden rounded-lg border border-iron-800"
        >
          {members.map((m) => (
            <MemberRow
              key={m.id}
              member={m}
              self={m.id === currentUserId}
              mailEnabled={mailEnabled}
              onUpdated={update}
            />
          ))}
        </ul>
      </section>
      <aside className="animate-rise [animation-delay:120ms]">
        <InviteForm mailEnabled={mailEnabled} onInvited={update} />
      </aside>
    </div>
  );
}

function MemberRow({
  member,
  self,
  mailEnabled,
  onUpdated,
}: {
  member: Member;
  self: boolean;
  mailEnabled: boolean;
  onUpdated: (m: Member) => void;
}) {
  const api = useApi();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const [outcome, setOutcome] = useState<InvitationResult>();

  async function renew() {
    setPending(true);
    setError(undefined);
    setOutcome(undefined);
    try {
      const res = await api.team.renewInvitation({ userId: member.id });
      if (res.result?.member) {
        onUpdated(res.result.member);
        setOutcome(res.result);
      }
    } catch (err) {
      setError(describeError(err));
    } finally {
      setPending(false);
    }
  }

  const waiting = member.status !== MemberStatus.ACTIVE;
  return (
    <li className="px-4 py-4 transition-colors hover:bg-iron-900/60">
      <div className="flex flex-wrap items-start justify-between gap-x-6 gap-y-3">
        <div className="min-w-0">
          <p className="font-semibold text-iron-100">
            {member.username}
            {self && <span className="ml-2 text-sm font-normal text-iron-300">(vous)</span>}
          </p>
          <p className="mt-0.5 break-all font-mono text-sm text-iron-300">{member.email}</p>
          <p className="mt-2 flex flex-wrap items-center gap-2 text-xs">
            <span className="rounded border border-iron-700 px-2 py-0.5 text-iron-100">
              {member.role === UserRole.ADMIN ? 'Administrateur' : 'Membre'}
            </span>
            <StatusBadge member={member} />
          </p>
        </div>
        {waiting && (
          <Button
            variant="quiet"
            onClick={() => void renew()}
            disabled={pending}
            aria-label={`${mailEnabled ? 'Renvoyer l’invitation de' : 'Refaire le lien de'} ${member.username}`}
            className="shrink-0"
          >
            {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
            {pending ? 'Envoi…' : mailEnabled ? 'Renvoyer l’invitation' : 'Refaire le lien'}
          </Button>
        )}
      </div>
      {error && (
        <p role="alert" className="mt-3 text-sm text-danger-400">
          {error}
        </p>
      )}
      {outcome && <InvitationOutcome result={outcome} onDismiss={() => setOutcome(undefined)} />}
    </li>
  );
}

function StatusBadge({ member }: { member: Member }) {
  switch (member.status) {
    case MemberStatus.ACTIVE:
      return (
        <span className="inline-flex items-center gap-1.5 text-iron-300">
          <CircleCheck aria-hidden="true" className="size-3.5 text-ok-400" />
          Actif
        </span>
      );
    case MemberStatus.INVITED:
      return (
        <span className="inline-flex items-center gap-1.5 text-iron-100">
          <Clock aria-hidden="true" className="size-3.5 text-gold-400" />
          Invité, lien valable jusqu’au{' '}
          {member.invitationExpiresAt
            ? until.format(timestampDate(member.invitationExpiresAt))
            : '…'}
        </span>
      );
    default:
      return (
        <span className="inline-flex items-center gap-1.5 text-iron-100">
          <TriangleAlert aria-hidden="true" className="size-3.5 text-warn-400" />
          Invitation expirée
        </span>
      );
  }
}

/**
 * Ce qu'il advient d'une invitation : partie par e-mail, ou lien montré une
 * seule fois, que l'administrateur transmet lui-même.
 */
function InvitationOutcome({
  result,
  onDismiss,
}: {
  result: InvitationResult;
  onDismiss: () => void;
}) {
  const member = result.member;
  if (!member) {
    return null;
  }
  if (result.emailed) {
    return (
      <p
        role="status"
        className="mt-4 flex items-start gap-2 rounded-md border border-ok-400/40 px-3 py-2.5 text-sm text-iron-100"
      >
        <MailCheck aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-ok-400" />
        <span>
          L’invitation de {member.username} part vers{' '}
          <span className="break-all font-mono">{member.email}</span>. Le lien vaut 7 jours ; s’il
          n’arrive pas, renvoyez-la depuis la liste.
        </span>
      </p>
    );
  }
  return <LinkReveal result={result} username={member.username} onDismiss={onDismiss} />;
}

function LinkReveal({
  result,
  username,
  onDismiss,
}: {
  result: InvitationResult;
  username: string;
  onDismiss: () => void;
}) {
  const [copied, setCopied] = useState<'yes' | 'refused'>();
  const input = useRef<HTMLInputElement>(null);
  const ids = { title: useId(), field: useId() };

  async function copy() {
    try {
      await navigator.clipboard.writeText(result.invitationUrl);
      setCopied('yes');
    } catch {
      // Presse-papiers refusé : le lien est sélectionné, à copier au clavier.
      input.current?.select();
      setCopied('refused');
    }
  }

  return (
    <div
      role="region"
      aria-labelledby={ids.title}
      className="mt-4 rounded-md border border-gold-400/50 bg-gold-950 px-4 py-4"
    >
      <p id={ids.title} className="font-semibold text-iron-100">
        Lien d’invitation de {username}
      </p>
      <p className="mt-1 text-sm text-iron-100">
        Transmettez-le à {username} par un canal de confiance. Il ne sera plus affiché : s’il se
        perd, refaites-en un, ce qui annule celui-ci.
      </p>
      <label htmlFor={ids.field} className="sr-only">
        Lien d’invitation
      </label>
      <input
        ref={input}
        id={ids.field}
        readOnly
        value={result.invitationUrl}
        onFocus={(e) => e.currentTarget.select()}
        className={`${fieldClass} mt-3 border-iron-700 font-mono text-sm`}
      />
      {result.localDashboard && (
        <p className="mt-3 text-sm text-iron-300">
          Ce serveur n’a pas de domaine : {username} ouvre d’abord un tunnel SSH vers lui, puis le
          lien dans son navigateur.
          <code className="mt-1.5 block overflow-x-auto rounded bg-iron-950 px-2 py-1.5 font-mono text-xs text-iron-100">
            ssh -N -L 8080:127.0.0.1:8080 utilisateur@adresse-du-serveur
          </code>
        </p>
      )}
      <div className="mt-4 flex flex-wrap gap-3">
        <Button onClick={() => void copy()}>
          {copied === 'yes' ? (
            <Check aria-hidden="true" className="size-4" />
          ) : (
            <Copy aria-hidden="true" className="size-4" />
          )}
          {copied === 'yes' ? 'Copié' : 'Copier le lien'}
        </Button>
        <Button variant="quiet" onClick={onDismiss}>
          J’ai transmis le lien
        </Button>
      </div>
      <p role="status" className="empty:hidden mt-2 text-sm text-iron-100">
        {copied === 'yes' && 'Le lien est dans le presse-papiers.'}
        {copied === 'refused' &&
          'Le navigateur refuse la copie : le lien est sélectionné, copiez-le avec Ctrl+C.'}
      </p>
    </div>
  );
}

interface FormValues {
  username: string;
  email: string;
  role: UserRole;
}

const emptyForm: FormValues = { username: '', email: '', role: UserRole.MEMBER };

type Problems = Partial<Record<'username' | 'email', string>>;

/** Vérifications locales ; le serveur valide le reste (adresse, nom déjà pris). */
function checkForm(v: FormValues): Problems {
  const problems: Problems = {};
  if (!usernamePattern.test(v.username.trim())) {
    problems.username =
      'De 2 à 64 caractères : lettres sans accent, chiffres, point, tiret ou soulignement.';
  }
  if (!/^[^\s@]+@[^\s@]+$/.test(v.email.trim())) {
    problems.email = 'Saisissez une adresse e-mail, comme nom@exemple.org.';
  }
  return problems;
}

const roles = [
  {
    role: UserRole.MEMBER,
    label: 'Membre',
    detail: 'Consulte le serveur : conteneurs, déploiements, domaines, mesures.',
  },
  {
    role: UserRole.ADMIN,
    label: 'Administrateur',
    detail: 'Gère tout, y compris l’équipe et les réglages.',
  },
];

function InviteForm({
  mailEnabled,
  onInvited,
}: {
  mailEnabled: boolean;
  onInvited: (m: Member) => void;
}) {
  const api = useApi();
  const [values, setValues] = useState<FormValues>(emptyForm);
  const [problems, setProblems] = useState<Problems>({});
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const [outcome, setOutcome] = useState<InvitationResult>();
  const ids = { username: useId(), email: useId(), role: useId(), error: useId() };

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const found = checkForm(values);
    setProblems(found);
    if (Object.keys(found).length > 0) {
      setError(undefined);
      return;
    }
    setPending(true);
    setError(undefined);
    setOutcome(undefined);
    try {
      const res = await api.team.inviteMember({
        username: values.username.trim(),
        email: values.email.trim(),
        role: values.role,
      });
      if (res.result?.member) {
        onInvited(res.result.member);
        setOutcome(res.result);
      }
      setValues(emptyForm);
    } catch (err) {
      setError(describeError(err));
    } finally {
      setPending(false);
    }
  }

  const field = (key: 'username' | 'email') => ({
    id: ids[key],
    name: key,
    value: values[key],
    onChange: (e: { target: { value: string } }) =>
      setValues((v) => ({ ...v, [key]: e.target.value })),
    'aria-invalid': problems[key] ? true : undefined,
    'aria-describedby': problems[key] ? `${ids[key]}-hint` : undefined,
    autoComplete: 'off',
    autoCapitalize: 'none',
    spellCheck: false,
    className: `${fieldClass} font-mono ${problems[key] ? 'border-danger-400' : 'border-iron-700'}`,
  });

  const hint = (key: 'username' | 'email') =>
    problems[key] && (
      <p id={`${ids[key]}-hint`} className="mt-1.5 text-sm text-danger-400">
        {problems[key]}
      </p>
    );

  return (
    <form
      noValidate
      onSubmit={submit}
      aria-labelledby="team-invite"
      aria-describedby={error ? ids.error : undefined}
      className="rounded-lg border border-iron-800 bg-iron-900 px-5 py-5"
    >
      <h2 id="team-invite" className="flex items-center gap-2 text-lg font-semibold">
        <UserPlus aria-hidden="true" className="size-5 text-gold-400" />
        Inviter un membre
      </h2>
      <p className="mt-2 text-sm text-iron-300">
        {mailEnabled ? (
          'Le lien part par e-mail à l’adresse saisie. Il vaut 7 jours, une seule fois.'
        ) : (
          <>
            Aucun relais SMTP n’est configuré : le lien s’affichera ici une seule fois, et vous le
            transmettrez vous-même.{' '}
            <Link
              href="/reglages"
              className="rounded text-gold-300 underline underline-offset-2 hover:text-gold-400 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gold-400"
            >
              Voir l’envoi d’e-mails
            </Link>
          </>
        )}
      </p>

      <div className="mt-5">
        <label htmlFor={ids.username} className="text-sm font-medium">
          Nom d’utilisateur
        </label>
        <input {...field('username')} placeholder="prenom.nom" />
        {problems.username ? (
          hint('username')
        ) : (
          <p className="mt-1.5 text-xs text-iron-300">Il sert à se connecter.</p>
        )}
      </div>

      <div className="mt-4">
        <label htmlFor={ids.email} className="text-sm font-medium">
          Adresse e-mail
        </label>
        <input {...field('email')} type="email" inputMode="email" placeholder="nom@exemple.org" />
        {hint('email')}
      </div>

      <fieldset className="mt-5">
        <legend className="text-sm font-medium">Rôle</legend>
        <div className="mt-1.5 grid gap-2">
          {roles.map((r) => (
            <label
              key={r.role}
              className={`flex min-h-11 cursor-pointer items-start gap-3 rounded-md border px-3 py-2.5 transition-colors has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-gold-400/35 ${
                values.role === r.role
                  ? 'border-gold-400 bg-iron-950'
                  : 'border-iron-700 hover:border-iron-500'
              }`}
            >
              <input
                type="radio"
                name="role"
                value={r.role}
                checked={values.role === r.role}
                onChange={() => setValues((v) => ({ ...v, role: r.role }))}
                className="mt-1 size-4 accent-gold-400"
              />
              <span>
                <span className="block text-sm font-semibold text-iron-100">{r.label}</span>
                <span className="block text-sm text-iron-300">{r.detail}</span>
              </span>
            </label>
          ))}
        </div>
      </fieldset>

      {error && (
        <p
          id={ids.error}
          role="alert"
          className="mt-5 rounded-md border border-danger-400/40 bg-danger-950 px-3 py-2.5 text-sm text-iron-100"
        >
          {error}
        </p>
      )}

      <Button type="submit" disabled={pending} className="mt-5 w-full">
        {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
        {pending
          ? 'Invitation…'
          : mailEnabled
            ? 'Envoyer l’invitation'
            : 'Créer le lien d’invitation'}
      </Button>

      {outcome && <InvitationOutcome result={outcome} onDismiss={() => setOutcome(undefined)} />}
    </form>
  );
}

function TeamSkeleton() {
  return (
    <div
      role="status"
      aria-live="polite"
      className="space-y-px overflow-hidden rounded-lg border border-iron-800 lg:max-w-[62%]"
    >
      <span className="sr-only">Chargement de l’équipe…</span>
      {[0, 1, 2].map((row) => (
        <div key={row} className="space-y-2 bg-iron-900/50 px-4 py-4">
          <span className="block h-3 w-40 animate-pulse rounded bg-iron-800" />
          <span className="block h-3 w-56 animate-pulse rounded bg-iron-800" />
        </div>
      ))}
    </div>
  );
}
