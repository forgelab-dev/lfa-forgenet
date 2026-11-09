import { Eye, EyeOff, LoaderCircle } from 'lucide-react';
import { useEffect, useId, useState, type FormEvent } from 'react';

import { Button } from '../../components/ui/Button';
import { ServerPlate } from '../../components/ui/ServerPlate';
import { Wordmark } from '../../components/ui/Wordmark';
import { describeError } from '../../lib/api/errors';
import { useSession } from './session';

const field =
  'mt-1.5 block min-h-11 w-full rounded-md border bg-iron-950 px-3 text-base text-iron-100 ' +
  'placeholder:text-iron-300/80 transition-colors duration-150 hover:border-iron-500 ' +
  'focus-visible:border-gold-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gold-400/35';

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
    <div className="grid min-h-dvh lg:grid-cols-[minmax(0,1.15fr)_minmax(26rem,1fr)]">
      <section className="engraved relative flex flex-col justify-between gap-10 px-6 pt-10 pb-8 sm:px-10 lg:px-16 lg:py-16">
        <Wordmark className="animate-rise text-xl" />
        <div className="animate-rise max-w-xl [animation-delay:80ms]">
          <p className="text-[clamp(2.25rem,5vw,4.25rem)] leading-[1.02] font-bold tracking-[-0.03em] text-balance">
            Votre serveur, <span className="text-gold-400">sous la main.</span>
          </p>
          <p className="mt-5 max-w-md text-lg text-iron-300">
            Conteneurs, déploiements et métriques du serveur, depuis un seul tableau de bord.
          </p>
        </div>
        <ServerPlate className="animate-rise max-w-md [animation-delay:160ms]" />
      </section>

      <main className="flex items-center justify-center border-iron-800 bg-iron-900 px-6 py-12 sm:px-10 lg:border-l">
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
        </form>
      </main>
    </div>
  );
}
