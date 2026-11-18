import { Code, ConnectError } from '@connectrpc/connect';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { LinkIcon, LoaderCircle, RotateCw } from 'lucide-react';
import { useCallback, useEffect, useId, useState, type FormEvent } from 'react';
import { useLocation } from 'wouter';

import { Button } from '../../components/ui/Button';
import { fieldClass as field } from '../../components/ui/field';
import { AccessLinkPurpose, type InspectAccessLinkResponse } from '../../gen/forgenet/v1/auth_pb';
import { describeError } from '../../lib/api/errors';
import { checkChosenPassword, MIN_PASSWORD_CHARS } from '../account/password';
import { AuthFrame } from './AuthFrame';
import { useSession } from './session';

/**
 * Lit le jeton dans le fragment de l'adresse (/acces#<jeton>), puis l'en
 * retire : il ne reste ni dans la barre d'adresse ni dans l'historique. Le
 * navigateur n'envoie jamais le fragment au serveur (ADR-0014).
 */
function takeTokenFromFragment(): string {
  const token = decodeURIComponent(window.location.hash.replace(/^#/, ''));
  if (window.location.hash) {
    window.history.replaceState(window.history.state, '', window.location.pathname);
  }
  return token;
}

const copy: Record<AccessLinkPurpose, { title: string; intro: (username: string) => string }> = {
  [AccessLinkPurpose.FIRST_LOGIN]: {
    title: 'Premier accès',
    intro: (u) => `Choisissez le mot de passe du compte ${u}. Vous serez connecté juste après.`,
  },
  [AccessLinkPurpose.PASSWORD_RESET]: {
    title: 'Nouveau mot de passe',
    intro: (u) =>
      `Choisissez un nouveau mot de passe pour ${u}. Vos sessions ouvertes ailleurs seront fermées.`,
  },
  [AccessLinkPurpose.INVITATION]: {
    title: 'Rejoindre ce serveur',
    intro: (u) =>
      `Vous êtes invité sur ce serveur ForgeNet sous le nom ${u}. Choisissez votre mot de passe pour entrer.`,
  },
  [AccessLinkPurpose.UNSPECIFIED]: {
    title: 'Choisir un mot de passe',
    intro: (u) => `Choisissez le mot de passe du compte ${u}.`,
  },
};

const until = new Intl.DateTimeFormat('fr-FR', { dateStyle: 'long', timeStyle: 'short' });

type LinkState =
  | { kind: 'checking' }
  | { kind: 'invalid' }
  | { kind: 'unreachable'; message: string }
  | { kind: 'ready'; link: InspectAccessLinkResponse };

/**
 * Page d'un lien à usage unique : premier accès, mot de passe oublié ou
 * invitation. Elle s'affiche avec ou sans session, et en ouvre une une fois
 * le mot de passe choisi.
 */
export function AccessLinkScreen() {
  const { api } = useSession();
  const [token, setToken] = useState(takeTokenFromFragment);
  const [state, setState] = useState<LinkState>(() =>
    token ? { kind: 'checking' } : { kind: 'invalid' }
  );

  // Un autre lien collé dans le même onglet ne change que le fragment : le
  // navigateur ne recharge pas la page, c'est à elle de le relire.
  useEffect(() => {
    const onHashChange = () => {
      const next = takeTokenFromFragment();
      if (next) {
        setToken(next);
        setState({ kind: 'checking' });
      }
    };
    window.addEventListener('hashchange', onHashChange);
    return () => window.removeEventListener('hashchange', onHashChange);
  }, []);

  useEffect(() => {
    document.title = 'Accès · ForgeNet';
  }, []);

  useEffect(() => {
    if (state.kind !== 'checking') {
      return;
    }
    const abort = new AbortController();
    api.auth
      .inspectAccessLink({ token }, { signal: abort.signal })
      .then((link) => setState({ kind: 'ready', link }))
      .catch((err: unknown) => {
        const { code } = ConnectError.from(err);
        if (code === Code.Canceled) {
          return;
        }
        setState(
          code === Code.NotFound
            ? { kind: 'invalid' }
            : { kind: 'unreachable', message: describeError(err) }
        );
      });
    return () => abort.abort();
  }, [api, token, state.kind]);

  const invalidate = useCallback(() => setState({ kind: 'invalid' }), []);

  return (
    <AuthFrame>
      <div className="animate-rise w-full max-w-sm [animation-delay:120ms]">
        {state.kind === 'checking' && (
          <p role="status" className="flex items-center gap-2 text-iron-300">
            <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />
            Vérification du lien…
          </p>
        )}
        {state.kind === 'invalid' && <InvalidLink />}
        {state.kind === 'unreachable' && (
          <div role="alert">
            <h1 className="text-3xl font-semibold tracking-tight">Lien non vérifié</h1>
            <p className="mt-3 text-iron-300">{state.message}</p>
            <Button className="mt-6" onClick={() => setState({ kind: 'checking' })}>
              <RotateCw aria-hidden="true" className="size-4" />
              Réessayer
            </Button>
          </div>
        )}
        {state.kind === 'ready' && (
          <ChoosePassword token={token} link={state.link} onInvalid={invalidate} />
        )}
      </div>
    </AuthFrame>
  );
}

function InvalidLink() {
  const [, navigate] = useLocation();
  return (
    <>
      <h1 className="text-3xl font-semibold tracking-tight">Ce lien ne vaut plus</h1>
      <p className="mt-3 text-iron-300">
        Un lien d’accès ne sert qu’une fois, pendant une durée limitée, et un lien plus récent
        remplace les précédents.
      </p>
      <p className="mt-4 text-sm text-iron-300">
        Pour un premier accès, refaites-en un sur le serveur :
      </p>
      <pre className="mt-2 overflow-x-auto rounded-md bg-iron-950 px-3 py-2.5 font-mono text-xs text-iron-100">
        sudo forgenet user login-link --username NOM
      </pre>
      <p className="mt-4 text-sm text-iron-300">
        Pour un mot de passe oublié, redemandez un lien depuis l’écran de connexion.
      </p>
      <Button variant="quiet" className="mt-6 w-full" onClick={() => navigate('/')}>
        Aller à la connexion
      </Button>
    </>
  );
}

function ChoosePassword({
  token,
  link,
  onInvalid,
}: {
  token: string;
  link: InspectAccessLinkResponse;
  onInvalid: () => void;
}) {
  const { api, adopt } = useSession();
  const [, navigate] = useLocation();
  const [next, setNext] = useState('');
  const [confirm, setConfirm] = useState('');
  const [reveal, setReveal] = useState(false);
  const [problems, setProblems] = useState<Partial<Record<'next' | 'confirm', string>>>({});
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const ids = { next: useId(), confirm: useId(), reveal: useId(), error: useId() };
  const text = copy[link.purpose] ?? copy[AccessLinkPurpose.UNSPECIFIED];

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const found = checkChosenPassword(next, confirm);
    setProblems(found);
    if (Object.keys(found).length > 0) {
      setError(undefined);
      return;
    }
    setPending(true);
    setError(undefined);
    try {
      const opened = await api.auth.redeemAccessLink({ token, newPassword: next });
      adopt(opened);
      navigate('/', { replace: true });
    } catch (err) {
      if (ConnectError.from(err).code === Code.NotFound) {
        onInvalid();
        return;
      }
      setError(describeError(err));
      setPending(false);
    }
  }

  const input = (key: 'next' | 'confirm', value: string, set: (v: string) => void) => ({
    id: ids[key],
    name: key === 'next' ? 'new-password' : 'confirm-password',
    type: reveal ? 'text' : 'password',
    autoComplete: 'new-password',
    autoCapitalize: 'none',
    spellCheck: false,
    value,
    onChange: (e: { target: { value: string } }) => set(e.target.value),
    'aria-invalid': problems[key] ? true : undefined,
    'aria-describedby': problems[key] || key === 'next' ? `${ids[key]}-hint` : undefined,
    className: `${field} ${problems[key] ? 'border-danger-400' : 'border-iron-700'}`,
  });

  return (
    <form noValidate onSubmit={submit} aria-describedby={error ? ids.error : undefined}>
      <h1 className="text-3xl font-semibold tracking-tight">{text.title}</h1>
      <p className="mt-2 text-iron-300">{text.intro(link.username)}</p>
      {link.expiresAt && (
        <p className="mt-3 flex items-start gap-2 text-sm text-iron-300">
          <LinkIcon aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-gold-400" />
          <span>
            Lien valable une seule fois, jusqu’au {until.format(timestampDate(link.expiresAt))}.
          </span>
        </p>
      )}

      {/* Champ caché pour les gestionnaires de mots de passe : ils rattachent
          ainsi le mot de passe choisi au bon compte. */}
      <input
        type="text"
        name="username"
        autoComplete="username"
        value={link.username}
        readOnly
        hidden
      />

      <div className="mt-8">
        <label htmlFor={ids.next} className="text-sm font-medium">
          Mot de passe
        </label>
        <input {...input('next', next, setNext)} />
        {problems.next ? (
          <p id={`${ids.next}-hint`} className="mt-1.5 text-sm text-danger-400">
            {problems.next}
          </p>
        ) : (
          <p id={`${ids.next}-hint`} className="mt-1.5 text-xs text-iron-300">
            {MIN_PASSWORD_CHARS} caractères au moins. Une phrase de quelques mots se retient mieux
            et résiste mieux.
          </p>
        )}
      </div>

      <div className="mt-5">
        <label htmlFor={ids.confirm} className="text-sm font-medium">
          Confirmation
        </label>
        <input {...input('confirm', confirm, setConfirm)} />
        {problems.confirm && (
          <p id={`${ids.confirm}-hint`} className="mt-1.5 text-sm text-danger-400">
            {problems.confirm}
          </p>
        )}
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
          className="mt-4 rounded-md border border-danger-400/40 bg-danger-950 px-3 py-2.5 text-sm text-iron-100"
        >
          {error}
        </p>
      )}

      <Button type="submit" disabled={pending} className="mt-6 w-full">
        {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
        {pending ? 'Enregistrement…' : 'Choisir ce mot de passe'}
      </Button>
    </form>
  );
}
