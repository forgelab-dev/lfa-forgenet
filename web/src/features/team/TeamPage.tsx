import { timestampDate } from '@bufbuild/protobuf/wkt';
import {
  AtSign,
  Ban,
  Check,
  CircleCheck,
  Clock,
  Copy,
  LoaderCircle,
  LogOut,
  MailCheck,
  RotateCcw,
  ShieldAlert,
  Trash2,
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
import {
  MemberStatus,
  type InvitationResult,
  type Member,
  type UpdateMemberEmailResponse,
} from '../../gen/forgenet/v1/team_pb';
import { describeError } from '../../lib/api/errors';
import { useRequest } from '../../lib/useRequest';
import { checkNewEmail } from '../account/email';
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
  const [removed, setRemoved] = useState<ReadonlySet<string>>(new Set());
  const members = useMemo(() => {
    if (request.status !== 'success') {
      return [];
    }
    const byId = new Map(request.data.members.map((m) => [m.id, m]));
    for (const m of updates) {
      byId.set(m.id, m);
    }
    return [...byId.values()].filter((m) => !removed.has(m.id));
  }, [request, updates, removed]);
  const update = useCallback((m: Member) => setUpdates((all) => [...all, m]), []);
  const [notice, setNotice] = useState<string>();
  const remove = useCallback((m: Member) => {
    setRemoved((all) => new Set(all).add(m.id));
    setNotice(
      m.status === MemberStatus.ACTIVE || m.status === MemberStatus.SUSPENDED
        ? `Le compte ${m.username} est supprimé. Le journal d’audit garde ses entrées.`
        : `L’invitation de ${m.username} est annulée : son lien ne vaut plus.`
    );
  }, []);

  if (request.status === 'loading') {
    return <TeamSkeleton />;
  }
  if (request.status === 'error') {
    return <RequestError error={request.error} onRetry={retry} />;
  }
  const mailEnabled = request.data.mailEnabled;
  const waiting = members.filter(
    (m) => m.status === MemberStatus.INVITED || m.status === MemberStatus.INVITATION_EXPIRED
  ).length;
  const suspended = members.filter((m) => m.status === MemberStatus.SUSPENDED).length;

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
          {suspended > 0 && (
            <span className="font-normal text-iron-300">
              {' '}
              · {suspended === 1 ? '1 suspendu' : `${suspended} suspendus`}
            </span>
          )}
        </h2>
        <p
          role="status"
          className="empty:hidden mt-3 rounded-md border border-ok-400/40 px-3 py-2 text-sm text-iron-100"
        >
          {notice}
        </p>
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
              onRemoved={remove}
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
  onRemoved,
}: {
  member: Member;
  self: boolean;
  mailEnabled: boolean;
  onUpdated: (m: Member) => void;
  onRemoved: (m: Member) => void;
}) {
  const api = useApi();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const [outcome, setOutcome] = useState<InvitationResult>();
  const [editingEmail, setEditingEmail] = useState(false);
  const [emailDone, setEmailDone] = useState<string>();
  const emailButton = useRef<HTMLButtonElement>(null);

  const closeEmail = (message?: string) => {
    setEditingEmail(false);
    setEmailDone(message);
    // Le focus revient au bouton qui a ouvert le formulaire.
    requestAnimationFrame(() => emailButton.current?.focus());
  };

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

  const waiting =
    member.status === MemberStatus.INVITED || member.status === MemberStatus.INVITATION_EXPIRED;
  return (
    <li className="px-4 py-4 transition-colors hover:bg-iron-900/60">
      <div className="flex flex-wrap items-start justify-between gap-x-6 gap-y-3">
        <div className="min-w-0">
          <p className="font-semibold text-iron-100">
            {member.username}
            {self && <span className="ml-2 text-sm font-normal text-iron-300">(vous)</span>}
          </p>
          <p className="mt-0.5 flex flex-wrap items-center gap-x-1">
            <span className="break-all font-mono text-sm text-iron-300">{member.email}</span>
            {/* Sa propre adresse se change depuis la page Compte, avec le mot de passe (#100). */}
            {!self && !editingEmail && (
              <button
                ref={emailButton}
                type="button"
                onClick={() => {
                  setEmailDone(undefined);
                  setOutcome(undefined);
                  setEditingEmail(true);
                }}
                aria-label={`Changer l’adresse de ${member.username}`}
                className="-my-3 inline-flex min-h-11 items-center gap-1 rounded-md px-2 text-sm font-medium text-gold-400 hover:text-gold-300 focus-visible:outline-2 focus-visible:outline-gold-400"
              >
                <AtSign aria-hidden="true" className="size-3.5" />
                Changer
              </button>
            )}
          </p>
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
      {editingEmail && (
        <MemberEmailForm
          member={member}
          mailEnabled={mailEnabled}
          onCancel={() => closeEmail()}
          onChanged={(res) => {
            if (res.member) {
              onUpdated(res.member);
            }
            if (res.invitation) {
              setOutcome(res.invitation);
            }
            closeEmail(
              res.invitation
                ? `Adresse de ${member.username} corrigée ; son ancien lien ne vaut plus.`
                : `Adresse de ${member.username} corrigée.${
                    mailEnabled && member.status === MemberStatus.ACTIVE
                      ? ' Un avis part à l’ancienne adresse.'
                      : ''
                  }`
            );
          }}
        />
      )}
      <p
        role="status"
        className="empty:hidden mt-3 rounded-md border border-ok-400/40 px-3 py-2 text-sm text-iron-100"
      >
        {emailDone}
      </p>
      {outcome && <InvitationOutcome result={outcome} onDismiss={() => setOutcome(undefined)} />}
      <MemberActions member={member} self={self} onUpdated={onUpdated} onRemoved={onRemoved} />
    </li>
  );
}

type Confirming = 'suspend' | 'delete';

function plural(n: bigint, one: string, many: string) {
  return n === 1n ? one : `${n} ${many}`;
}

/**
 * Les gestes d'un administrateur sur un compte (#103) : fermer ses sessions,
 * le suspendre ou le réactiver, le supprimer. Sur son propre compte, seules
 * ses autres sessions se ferment : le serveur refuse le reste.
 */
function MemberActions({
  member,
  self,
  onUpdated,
  onRemoved,
}: {
  member: Member;
  self: boolean;
  onUpdated: (m: Member) => void;
  onRemoved: (m: Member) => void;
}) {
  const api = useApi();
  const [confirming, setConfirming] = useState<Confirming>();
  const [typed, setTyped] = useState('');
  const [pending, setPending] = useState<'sessions' | 'resume' | Confirming>();
  const [error, setError] = useState<string>();
  const [done, setDone] = useState<string>();
  const ids = { typed: useId(), error: useId() };

  const suspended = member.status === MemberStatus.SUSPENDED;
  const active = member.status === MemberStatus.ACTIVE;
  const invitation = !active && !suspended;

  async function run(kind: NonNullable<typeof pending>, action: () => Promise<string>) {
    setPending(kind);
    setError(undefined);
    setDone(undefined);
    try {
      setDone(await action());
      setConfirming(undefined);
      setTyped('');
    } catch (err) {
      setError(describeError(err));
    } finally {
      setPending(undefined);
    }
  }

  const closeSessions = () =>
    run('sessions', async () => {
      const res = await api.team.closeMemberSessions({ userId: member.id });
      const n = res.closedSessions;
      if (n === 0n) {
        return self ? 'Aucune autre session n’était ouverte.' : 'Aucune session n’était ouverte.';
      }
      return `${plural(n, 'Une session fermée', 'sessions fermées')}${self ? ' ; la vôtre reste ouverte.' : ' : il devra se reconnecter.'}`;
    });

  const suspend = () =>
    run('suspend', async () => {
      const res = await api.team.suspendMember({ userId: member.id });
      if (res.member) {
        onUpdated(res.member);
      }
      const n = res.closedSessions;
      return n === 0n
        ? `${member.username} est suspendu.`
        : `${member.username} est suspendu ; ${plural(n, 'une session fermée', 'sessions fermées')}.`;
    });

  const resume = () =>
    run('resume', async () => {
      const res = await api.team.resumeMember({ userId: member.id });
      if (res.member) {
        onUpdated(res.member);
      }
      return res.member?.status === MemberStatus.ACTIVE
        ? `${member.username} peut de nouveau se connecter.`
        : `${member.username} est réactivé ; son invitation est à renvoyer.`;
    });

  const remove = () =>
    run('delete', async () => {
      await api.team.deleteMember({ userId: member.id, confirmUsername: typed.trim() });
      onRemoved(member);
      return '';
    });

  function cancel() {
    setConfirming(undefined);
    setTyped('');
    setError(undefined);
  }

  const busy = pending !== undefined;
  return (
    <>
      <p
        role="status"
        className="empty:hidden mt-3 rounded-md border border-ok-400/40 px-3 py-2 text-sm text-iron-100"
      >
        {done}
      </p>
      {error && !confirming && (
        <p role="alert" className="mt-3 text-sm text-danger-400">
          {error}
        </p>
      )}
      {!confirming ? (
        <div className="mt-3 flex flex-wrap items-center justify-end gap-2 border-t border-iron-800 pt-3">
          {active && (
            <Button
              variant="quiet"
              onClick={() => void closeSessions()}
              disabled={busy}
              aria-label={
                self ? 'Fermer mes autres sessions' : `Fermer les sessions de ${member.username}`
              }
            >
              {pending === 'sessions' ? (
                <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />
              ) : (
                <LogOut aria-hidden="true" className="size-4" />
              )}
              {self ? 'Fermer mes autres sessions' : 'Fermer les sessions'}
            </Button>
          )}
          {!self && suspended && (
            <Button
              variant="quiet"
              onClick={() => void resume()}
              disabled={busy}
              aria-label={`Réactiver ${member.username}`}
            >
              {pending === 'resume' ? (
                <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />
              ) : (
                <RotateCcw aria-hidden="true" className="size-4" />
              )}
              Réactiver
            </Button>
          )}
          {!self && !suspended && (
            <Button
              variant="quiet"
              onClick={() => {
                setConfirming('suspend');
                setDone(undefined);
                setError(undefined);
              }}
              disabled={busy}
              aria-label={`Suspendre ${member.username}`}
            >
              <Ban aria-hidden="true" className="size-4" />
              Suspendre
            </Button>
          )}
          {!self && (
            <Button
              variant="quiet"
              onClick={() => {
                setConfirming('delete');
                setDone(undefined);
                setError(undefined);
              }}
              disabled={busy}
              aria-label={`${invitation ? 'Annuler l’invitation de' : 'Supprimer'} ${member.username}`}
            >
              <Trash2 aria-hidden="true" className="size-4" />
              {invitation ? 'Annuler l’invitation' : 'Supprimer'}
            </Button>
          )}
        </div>
      ) : (
        <div
          role="group"
          aria-label={
            confirming === 'suspend'
              ? `Suspendre ${member.username}`
              : `Supprimer ${member.username}`
          }
          aria-describedby={error ? ids.error : undefined}
          className="mt-3 rounded-md border border-danger-400/40 bg-danger-950 px-3 py-3"
        >
          {confirming === 'suspend' ? (
            <p className="text-sm text-iron-100">
              Suspendre <span className="font-semibold">{member.username}</span> ? Il ne peut plus
              se connecter, et ses sessions, jetons d’API et liens en cours sont fermés tout de
              suite. Son compte et son historique restent : vous pourrez le réactiver.
            </p>
          ) : (
            <>
              <p className="text-sm text-iron-100">
                {invitation ? (
                  <>
                    Annuler l’invitation de <span className="font-semibold">{member.username}</span>{' '}
                    ? Le compte, qui n’a pas encore servi, est supprimé et son lien ne vaut plus.
                  </>
                ) : (
                  <>
                    Supprimer <span className="font-semibold">{member.username}</span> ? Le compte,
                    ses sessions et ses jetons d’API sont effacés. C’est <strong>définitif</strong>{' '}
                    ; le journal d’audit garde ses entrées.
                  </>
                )}
              </p>
              <div className="mt-3">
                <label htmlFor={ids.typed} className="text-sm text-iron-100">
                  Tapez <span className="font-mono font-semibold">{member.username}</span> pour
                  confirmer
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
            </>
          )}
          {error && (
            <p id={ids.error} role="alert" className="mt-2 text-sm text-danger-400">
              {error}
            </p>
          )}
          <div className="mt-3 flex flex-wrap gap-3">
            {confirming === 'suspend' ? (
              <Button variant="danger" onClick={() => void suspend()} disabled={busy}>
                {busy && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
                {busy ? 'Suspension…' : `Suspendre ${member.username}`}
              </Button>
            ) : (
              <Button
                variant="danger"
                onClick={() => void remove()}
                disabled={busy || typed.trim() !== member.username}
              >
                {busy && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
                {busy
                  ? 'Suppression…'
                  : invitation
                    ? 'Annuler l’invitation'
                    : `Supprimer ${member.username}`}
              </Button>
            )}
            <Button variant="quiet" onClick={cancel} disabled={busy}>
              Garder le compte
            </Button>
          </div>
        </div>
      )}
    </>
  );
}

/**
 * Correction de l'adresse d'un compte par un administrateur (#100). Le
 * serveur annule les liens partis vers l'ancienne adresse ; un compte qui
 * attendait son mot de passe est réinvité à la nouvelle.
 */
function MemberEmailForm({
  member,
  mailEnabled,
  onChanged,
  onCancel,
}: {
  member: Member;
  mailEnabled: boolean;
  onChanged: (res: UpdateMemberEmailResponse) => void;
  onCancel: () => void;
}) {
  const api = useApi();
  const [email, setEmail] = useState('');
  const [problem, setProblem] = useState<string>();
  const [error, setError] = useState<string>();
  const [pending, setPending] = useState(false);
  const ids = { email: useId(), note: useId(), error: useId(), hint: useId() };

  const waiting =
    member.status === MemberStatus.INVITED || member.status === MemberStatus.INVITATION_EXPIRED;
  const note = waiting
    ? `Son lien actuel ne vaudra plus, et ${mailEnabled ? 'une nouvelle invitation partira à cette adresse' : 'un nouveau lien sera fait, à lui transmettre'}.`
    : member.status === MemberStatus.SUSPENDED
      ? 'Ses liens en attente seront annulés. Le compte reste suspendu.'
      : `Ses liens en attente seront annulés${mailEnabled ? ', et un avis partira à l’ancienne adresse' : ''}. Ses sessions restent ouvertes.`;

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const found = checkNewEmail(email, member.email);
    setProblem(found);
    setError(undefined);
    if (found) {
      return;
    }
    setPending(true);
    try {
      onChanged(await api.team.updateMemberEmail({ userId: member.id, email: email.trim() }));
    } catch (err) {
      setError(describeError(err));
      setPending(false);
    }
  }

  return (
    <form
      noValidate
      onSubmit={submit}
      aria-label={`Changer l’adresse de ${member.username}`}
      aria-describedby={error ? `${ids.note} ${ids.error}` : ids.note}
      className="mt-3 rounded-md border border-iron-700 bg-iron-900 px-3 py-3"
    >
      <label htmlFor={ids.email} className="text-sm font-medium text-iron-100">
        Nouvelle adresse de <span className="font-mono">{member.username}</span>
      </label>
      {/* autoFocus : le formulaire s'ouvre sur demande, le focus y entre. */}
      <input
        id={ids.email}
        type="email"
        inputMode="email"
        autoComplete="off"
        autoCapitalize="none"
        spellCheck={false}
        autoFocus
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        aria-invalid={problem ? true : undefined}
        aria-describedby={problem ? ids.hint : undefined}
        className={`${fieldClass} font-mono ${problem ? 'border-danger-400' : 'border-iron-700'}`}
      />
      {problem && (
        <p id={ids.hint} className="mt-1.5 text-sm text-danger-400">
          {problem}
        </p>
      )}
      <p id={ids.note} className="mt-2 text-xs text-iron-300">
        {note}
      </p>
      {error && (
        <p id={ids.error} role="alert" className="mt-2 text-sm text-danger-400">
          {error}
        </p>
      )}
      <div className="mt-3 flex flex-wrap gap-3">
        <Button type="submit" disabled={pending}>
          {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
          {pending ? 'Enregistrement…' : 'Enregistrer l’adresse'}
        </Button>
        <Button variant="quiet" onClick={onCancel} disabled={pending}>
          Annuler
        </Button>
      </div>
    </form>
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
    case MemberStatus.SUSPENDED:
      return (
        <span className="inline-flex items-center gap-1.5 text-iron-100">
          <Ban aria-hidden="true" className="size-3.5 text-danger-400" />
          Suspendu
          {member.suspendedAt && <> depuis le {until.format(timestampDate(member.suspendedAt))}</>}
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
