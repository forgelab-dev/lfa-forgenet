import { KeyRound, LoaderCircle, Mail, TerminalSquare } from 'lucide-react';
import { useId, useRef, useState, type FormEvent } from 'react';

import { Button } from '../../components/ui/Button';
import { PageHeader } from '../../components/ui/PageHeader';
import { fieldClass } from '../../components/ui/field';
import { UserRole } from '../../gen/forgenet/v1/auth_pb';
import { describeError } from '../../lib/api/errors';
import { useApi, useCurrentUser, useSession } from '../auth/session';
import { checkEmailForm, type EmailForm, type EmailProblems } from './email';
import {
  checkPasswordForm,
  MIN_PASSWORD_CHARS,
  successMessage,
  type PasswordForm,
  type PasswordProblems,
} from './password';

const roleLabel: Record<UserRole, string> = {
  [UserRole.ADMIN]: 'Administrateur',
  [UserRole.MEMBER]: 'Membre',
  [UserRole.UNSPECIFIED]: 'Rôle inconnu',
};

/** Compte connecté : son identité et le changement de son mot de passe. */
export function AccountPage() {
  return (
    <>
      <PageHeader
        title="Mon compte"
        description="Votre identité sur ce serveur, son adresse e-mail et le mot de passe qui la protège."
      />
      <div className="mt-8 grid items-start gap-6 lg:grid-cols-[minmax(0,1fr)_minmax(0,20rem)]">
        <ChangePasswordForm />
        <aside className="animate-rise space-y-6 [animation-delay:60ms]">
          <Identity />
          <LostPassword />
        </aside>
      </div>
    </>
  );
}

function Identity() {
  const user = useCurrentUser();
  const [editing, setEditing] = useState(false);
  const [changed, setChanged] = useState<string>();
  const changeButton = useRef<HTMLButtonElement>(null);

  const close = (message?: string) => {
    setEditing(false);
    setChanged(message);
    // Le focus revient au bouton qui a ouvert le formulaire.
    requestAnimationFrame(() => changeButton.current?.focus());
  };

  return (
    <section
      aria-labelledby="account-identity"
      className="rounded-lg border border-iron-800 px-5 py-5"
    >
      <h2 id="account-identity" className="text-base font-semibold">
        Identité
      </h2>
      <dl className="mt-3 divide-y divide-iron-800 text-sm">
        <Row label="Nom d’utilisateur" value={user.username} mono />
        <div className="py-2.5">
          <div className="flex items-start justify-between gap-3">
            <div className="min-w-0">
              <dt className="text-xs text-iron-300">Adresse e-mail</dt>
              <dd className="mt-0.5 break-all font-mono text-iron-100">{user.email || '—'}</dd>
            </div>
            {!editing && (
              <button
                ref={changeButton}
                type="button"
                aria-label="Changer d’adresse e-mail"
                onClick={() => {
                  setChanged(undefined);
                  setEditing(true);
                }}
                className="-my-2 -mr-2 inline-flex min-h-11 shrink-0 items-center rounded-md px-2.5 text-sm font-medium text-gold-400 hover:text-gold-300 focus-visible:outline-2 focus-visible:outline-gold-400"
              >
                Changer
              </button>
            )}
          </div>
          {editing && <ChangeEmailForm onDone={close} onCancel={() => close()} />}
          <p role="status" className="empty:hidden mt-2 text-xs text-ok-400">
            {changed}
          </p>
        </div>
        <Row label="Rôle" value={roleLabel[user.role]} />
      </dl>
    </section>
  );
}

const emptyEmailForm: EmailForm = { email: '', password: '' };

/**
 * Changement de sa propre adresse (#100). Le mot de passe est demandé : une
 * session volée ne doit pas suffire à détourner le « mot de passe oublié ».
 */
function ChangeEmailForm({
  onDone,
  onCancel,
}: {
  onDone: (message: string) => void;
  onCancel: () => void;
}) {
  const api = useApi();
  const { updateUser } = useSession();
  const user = useCurrentUser();
  const [values, setValues] = useState<EmailForm>(emptyEmailForm);
  const [problems, setProblems] = useState<EmailProblems>({});
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const ids = { email: useId(), password: useId(), note: useId(), error: useId() };

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const found = checkEmailForm(values, user.email);
    setProblems(found);
    setError(undefined);
    if (Object.keys(found).length > 0) {
      return;
    }
    setPending(true);
    try {
      const res = await api.auth.changeEmail({
        currentPassword: values.password,
        newEmail: values.email.trim(),
      });
      if (res.user) {
        updateUser(res.user);
      }
      onDone('Adresse changée : les prochains e-mails du compte partiront à la nouvelle.');
    } catch (err) {
      setError(describeError(err));
      setPending(false);
    }
  }

  const field = (key: keyof EmailForm) => ({
    id: ids[key],
    name: key === 'email' ? 'new-email' : 'current-password',
    value: values[key],
    onChange: (e: { target: { value: string } }) => {
      const value = e.target.value;
      setValues((v) => ({ ...v, [key]: value }));
    },
    autoCapitalize: 'none',
    spellCheck: false,
    'aria-invalid': problems[key] ? true : undefined,
    'aria-describedby': problems[key] ? `${ids[key]}-hint` : undefined,
    className: `${fieldClass} ${problems[key] ? 'border-danger-400' : 'border-iron-700'}`,
  });

  const hint = (key: keyof EmailForm) =>
    problems[key] && (
      <p id={`${ids[key]}-hint`} className="mt-1.5 text-sm text-danger-400">
        {problems[key]}
      </p>
    );

  return (
    <form
      noValidate
      onSubmit={submit}
      aria-label="Changer d’adresse e-mail"
      aria-describedby={error ? `${ids.note} ${ids.error}` : ids.note}
      className="mt-3 rounded-md border border-iron-700 bg-iron-900 px-3.5 py-3.5"
    >
      <p id={ids.note} className="flex gap-2 text-xs text-iron-300">
        <Mail aria-hidden="true" className="mt-0.5 size-3.5 shrink-0 text-gold-400" />
        <span>
          Les liens en attente envoyés à l’ancienne adresse seront annulés, et un avis y partira si
          le serveur envoie des e-mails.
        </span>
      </p>
      <input
        type="text"
        name="username"
        autoComplete="username"
        value={user.username}
        readOnly
        hidden
      />
      <div className="mt-3">
        <label htmlFor={ids.email} className="text-sm font-medium">
          Nouvelle adresse e-mail
        </label>
        {/* autoFocus : le formulaire s'ouvre sur demande, le focus y entre. */}
        <input {...field('email')} type="email" autoComplete="email" inputMode="email" autoFocus />
        {hint('email')}
      </div>
      <div className="mt-3">
        <label htmlFor={ids.password} className="text-sm font-medium">
          Votre mot de passe
        </label>
        <input {...field('password')} type="password" autoComplete="current-password" />
        {hint('password')}
      </div>
      {error && (
        <p
          id={ids.error}
          role="alert"
          className="mt-3 rounded-md border border-danger-400/40 bg-danger-950 px-3 py-2.5 text-sm text-iron-100"
        >
          {error}
        </p>
      )}
      <div className="mt-4 flex flex-wrap gap-2">
        <Button type="submit" disabled={pending}>
          {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
          {pending ? 'Enregistrement…' : 'Enregistrer l’adresse'}
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
    <div className="py-2.5 first:pt-0 last:pb-0">
      <dt className="text-xs text-iron-300">{label}</dt>
      <dd className={`mt-0.5 break-all text-iron-100 ${mono ? 'font-mono' : ''}`}>
        {value || '—'}
      </dd>
    </div>
  );
}

function LostPassword() {
  return (
    <section aria-labelledby="account-lost" className="rounded-lg border border-iron-800 px-5 py-5">
      <h2 id="account-lost" className="flex items-center gap-2 text-base font-semibold">
        <TerminalSquare aria-hidden="true" className="size-4 text-iron-300" />
        Mot de passe perdu
      </h2>
      <p className="mt-2 text-sm text-iron-300">
        Un administrateur le remplace depuis le serveur, sans l’ancien :
      </p>
      <pre className="mt-3 overflow-x-auto rounded-md bg-iron-950 px-3 py-2.5 font-mono text-xs text-iron-100">
        {'sudo forgenet user set-password \\\n  --username NOM'}
      </pre>
    </section>
  );
}

const emptyForm: PasswordForm = { current: '', next: '', confirm: '' };

function ChangePasswordForm() {
  const api = useApi();
  const { username } = useCurrentUser();
  const [values, setValues] = useState<PasswordForm>(emptyForm);
  const [problems, setProblems] = useState<PasswordProblems>({});
  const [reveal, setReveal] = useState(false);
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const [success, setSuccess] = useState<string>();
  const ids = {
    current: useId(),
    next: useId(),
    confirm: useId(),
    reveal: useId(),
    error: useId(),
  };

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const found = checkPasswordForm(values);
    setProblems(found);
    setSuccess(undefined);
    if (Object.keys(found).length > 0) {
      setError(undefined);
      return;
    }
    setPending(true);
    setError(undefined);
    try {
      const res = await api.auth.changePassword({
        currentPassword: values.current,
        newPassword: values.next,
      });
      setValues(emptyForm);
      setReveal(false);
      setSuccess(successMessage(res.revokedSessions));
    } catch (err) {
      setError(describeError(err));
    } finally {
      setPending(false);
    }
  }

  const input = (key: keyof PasswordForm) => ({
    id: ids[key],
    name: key,
    type: reveal ? 'text' : 'password',
    value: values[key],
    onChange: (e: { target: { value: string } }) => {
      const value = e.target.value;
      setValues((v) => ({ ...v, [key]: value }));
    },
    autoCapitalize: 'none',
    spellCheck: false,
    'aria-invalid': problems[key] ? true : undefined,
    // Seul le nouveau mot de passe a une aide permanente ; les autres champs
    // n'ont de texte associé qu'en erreur.
    'aria-describedby': problems[key] || key === 'next' ? `${ids[key]}-hint` : undefined,
    className: `${fieldClass} ${problems[key] ? 'border-danger-400' : 'border-iron-700'}`,
  });

  const hint = (key: keyof PasswordForm, help?: string) =>
    problems[key] ? (
      <p id={`${ids[key]}-hint`} className="mt-1.5 text-sm text-danger-400">
        {problems[key]}
      </p>
    ) : (
      help && (
        <p id={`${ids[key]}-hint`} className="mt-1.5 text-xs text-iron-300">
          {help}
        </p>
      )
    );

  return (
    <form
      noValidate
      onSubmit={submit}
      aria-labelledby="account-password"
      aria-describedby={error ? ids.error : undefined}
      className="animate-rise rounded-lg border border-iron-800 bg-iron-900 px-5 py-5 sm:px-6 sm:py-6"
    >
      <h2 id="account-password" className="flex items-center gap-2 text-lg font-semibold">
        <KeyRound aria-hidden="true" className="size-5 text-gold-400" />
        Changer de mot de passe
      </h2>
      <p className="mt-2 max-w-prose text-sm text-iron-300">
        Vos autres sessions et vos jetons d’API seront fermés : un mot de passe changé parce qu’il a
        fuité ne doit plus rien ouvrir. Cette session-ci reste ouverte.
      </p>

      {/* Champ caché pour les gestionnaires de mots de passe : ils rattachent
          ainsi le nouveau mot de passe au bon compte. */}
      <input type="text" name="username" autoComplete="username" value={username} readOnly hidden />

      <div className="mt-5 max-w-md">
        <label htmlFor={ids.current} className="text-sm font-medium">
          Mot de passe actuel
        </label>
        <input {...input('current')} autoComplete="current-password" />
        {hint('current')}
      </div>

      <div className="mt-4 max-w-md">
        <label htmlFor={ids.next} className="text-sm font-medium">
          Nouveau mot de passe
        </label>
        <input {...input('next')} autoComplete="new-password" />
        {hint(
          'next',
          `${MIN_PASSWORD_CHARS} caractères au moins. Une phrase de quelques mots se retient mieux et résiste mieux.`
        )}
      </div>

      <div className="mt-4 max-w-md">
        <label htmlFor={ids.confirm} className="text-sm font-medium">
          Confirmation du nouveau mot de passe
        </label>
        <input {...input('confirm')} autoComplete="new-password" />
        {hint('confirm')}
      </div>

      <label
        htmlFor={ids.reveal}
        className="mt-3 flex w-fit min-h-11 cursor-pointer items-center gap-2.5 text-sm text-iron-300 hover:text-iron-100"
      >
        <input
          id={ids.reveal}
          type="checkbox"
          checked={reveal}
          onChange={(e) => setReveal(e.target.checked)}
          className="size-4 accent-gold-400"
        />
        Afficher les mots de passe
      </label>

      {error && (
        <p
          id={ids.error}
          role="alert"
          className="mt-4 max-w-md rounded-md border border-danger-400/40 bg-danger-950 px-3 py-2.5 text-sm text-iron-100"
        >
          {error}
        </p>
      )}
      <p
        role="status"
        className="empty:hidden mt-4 max-w-md rounded-md border border-ok-400/40 px-3 py-2.5 text-sm text-iron-100"
      >
        {success}
      </p>

      <Button type="submit" disabled={pending} className="mt-5 w-full sm:w-auto">
        {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
        {pending ? 'Changement…' : 'Changer le mot de passe'}
      </Button>
    </form>
  );
}
