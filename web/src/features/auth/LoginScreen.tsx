import { Eye, EyeOff, LoaderCircle } from 'lucide-react';
import { useEffect, useId, useState, type FormEvent } from 'react';
import { useLocation } from 'wouter';

import { Button } from '../../components/ui/Button';
import { fieldClass as field } from '../../components/ui/field';
import { describeError } from '../../lib/api/errors';
import { AuthFrame } from './AuthFrame';
import { useSession } from './session';

export function LoginScreen() {
  const { state, login } = useSession();
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [reveal, setReveal] = useState(false);
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const [missing, setMissing] = useState<{ identifier?: boolean; password?: boolean }>({});
  const ids = { identifier: useId(), password: useId(), error: useId() };

  const expired = state.status === 'anonymous' && state.notice === 'expired';

  useEffect(() => {
    document.title = 'Connexion · ForgeNet';
  }, []);

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const empty = { identifier: !identifier.trim(), password: !password };
    setMissing(empty);
    if (empty.identifier || empty.password) {
      setError(undefined);
      return;
    }
    setPending(true);
    setError(undefined);
    try {
      await login(identifier.trim(), password);
    } catch (err) {
      setError(describeError(err));
      setPassword('');
      setPending(false);
    }
  }

  return (
    <AuthFrame>
      <form
        noValidate
        onSubmit={submit}
        aria-describedby={error ? ids.error : undefined}
        className="animate-rise w-full max-w-sm [animation-delay:120ms]"
      >
        <h1 className="text-3xl font-semibold tracking-tight">Connexion</h1>
        <p className="mt-2 text-iron-300">Avec le compte créé sur ce serveur.</p>

        {expired && (
          <p
            role="status"
            className="mt-6 rounded-md border border-iron-700 bg-iron-850 px-3 py-2.5 text-sm text-iron-100"
          >
            Votre session a pris fin. Reconnectez-vous pour continuer.
          </p>
        )}

        <div className="mt-8">
          <label htmlFor={ids.identifier} className="text-sm font-medium">
            Identifiant
          </label>
          <input
            id={ids.identifier}
            name="login"
            autoComplete="username"
            autoCapitalize="none"
            spellCheck={false}
            placeholder="Nom d’utilisateur ou adresse e-mail"
            value={identifier}
            onChange={(e) => setIdentifier(e.target.value)}
            aria-invalid={missing.identifier || undefined}
            aria-describedby={missing.identifier ? `${ids.identifier}-hint` : undefined}
            className={`${field} ${missing.identifier ? 'border-danger-400' : 'border-iron-700'}`}
          />
          {missing.identifier && (
            <p id={`${ids.identifier}-hint`} className="mt-1.5 text-sm text-danger-400">
              Saisissez votre identifiant.
            </p>
          )}
        </div>

        <div className="mt-5">
          <label htmlFor={ids.password} className="text-sm font-medium">
            Mot de passe
          </label>
          <div className="relative">
            <input
              id={ids.password}
              name="password"
              type={reveal ? 'text' : 'password'}
              autoComplete="current-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              aria-invalid={missing.password || undefined}
              aria-describedby={missing.password ? `${ids.password}-hint` : undefined}
              className={`${field} pr-12 ${missing.password ? 'border-danger-400' : 'border-iron-700'}`}
            />
            <button
              type="button"
              onClick={() => setReveal((r) => !r)}
              aria-pressed={reveal}
              aria-label={reveal ? 'Masquer le mot de passe' : 'Afficher le mot de passe'}
              className="absolute inset-y-0 right-0 flex w-11 items-center justify-center rounded-r-md text-iron-300 transition-colors hover:text-iron-100"
            >
              {reveal ? (
                <EyeOff aria-hidden="true" className="size-4.5" />
              ) : (
                <Eye aria-hidden="true" className="size-4.5" />
              )}
            </button>
          </div>
          {missing.password && (
            <p id={`${ids.password}-hint`} className="mt-1.5 text-sm text-danger-400">
              Saisissez votre mot de passe.
            </p>
          )}
        </div>

        {error && (
          <p
            id={ids.error}
            role="alert"
            className="mt-6 rounded-md border border-danger-400/40 bg-danger-950 px-3 py-2.5 text-sm text-iron-100"
          >
            {error}
          </p>
        )}

        <Button type="submit" disabled={pending} className="mt-8 w-full">
          {pending && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
          {pending ? 'Connexion…' : 'Se connecter'}
        </Button>
        <LostPassword />
      </form>
    </AuthFrame>
  );
}

/**
 * « Mot de passe oublié ? » quand le serveur sait envoyer le lien ; sinon,
 * la commande à lancer sur le serveur. Rien tant que la réponse n'est pas
 * arrivée, ni si elle échoue : l'écran de connexion reste utilisable.
 */
function LostPassword() {
  const { api } = useSession();
  const [, navigate] = useLocation();
  const [reset, setReset] = useState<boolean>();

  useEffect(() => {
    const abort = new AbortController();
    api.auth
      .getLoginOptions({}, { signal: abort.signal })
      .then((opts) => setReset(opts.passwordReset))
      .catch(() => undefined);
    return () => abort.abort();
  }, [api]);

  if (reset === undefined) {
    return null;
  }
  if (reset) {
    return (
      <button
        type="button"
        onClick={() => navigate('/mot-de-passe-oublie')}
        className="mx-auto mt-4 flex min-h-11 items-center px-2 text-sm text-iron-300 underline-offset-4 hover:text-iron-100 hover:underline"
      >
        Mot de passe oublié ?
      </button>
    );
  }
  return (
    <p className="mt-6 text-center text-sm text-iron-300">
      Mot de passe perdu ? Sur le serveur :{' '}
      <code className="font-mono text-xs text-iron-100">
        forgenet user set-password --username NOM
      </code>
    </p>
  );
}
