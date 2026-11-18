import { LoaderCircle, MailCheck } from 'lucide-react';
import { useEffect, useId, useState, type FormEvent } from 'react';
import { useLocation } from 'wouter';

import { Button } from '../../components/ui/Button';
import { fieldClass as field } from '../../components/ui/field';
import { describeError } from '../../lib/api/errors';
import { AuthFrame } from './AuthFrame';
import { useSession } from './session';

/**
 * Demande d'un lien « mot de passe oublié » (ADR-0014). La réponse est la
 * même qu'un compte corresponde ou non : la page ne peut donc dire que « si
 * un compte correspond ».
 */
export function ForgotPasswordScreen() {
  const { api } = useSession();
  const [, navigate] = useLocation();
  const [login, setLogin] = useState('');
  const [missing, setMissing] = useState(false);
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const [sent, setSent] = useState(false);
  const ids = { login: useId(), error: useId() };

  useEffect(() => {
    document.title = 'Mot de passe oublié · ForgeNet';
  }, []);

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!login.trim()) {
      setMissing(true);
      setError(undefined);
      return;
    }
    setMissing(false);
    setPending(true);
    setError(undefined);
    try {
      await api.auth.requestPasswordReset({ login: login.trim() });
      setSent(true);
    } catch (err) {
      setError(describeError(err));
    } finally {
      setPending(false);
    }
  }

  const back = (
    <Button variant="quiet" className="mt-6 w-full" onClick={() => navigate('/')}>
      Retour à la connexion
    </Button>
  );

  return (
    <AuthFrame>
      <div className="animate-rise w-full max-w-sm [animation-delay:120ms]">
        {sent ? (
          <div role="status">
            <h1 className="text-3xl font-semibold tracking-tight">Vérifiez vos e-mails</h1>
            <p className="mt-3 flex items-start gap-2 text-iron-300">
              <MailCheck aria-hidden="true" className="mt-1 size-4 shrink-0 text-gold-400" />
              <span>
                Si un compte correspond, un e-mail vient de partir avec un lien pour choisir un
                nouveau mot de passe. Il vaut une heure, une seule fois.
              </span>
            </p>
            <p className="mt-3 text-sm text-iron-300">
              Rien après quelques minutes ? Regardez les indésirables, puis vérifiez l’identifiant
              saisi.
            </p>
            {back}
          </div>
        ) : (
          <form noValidate onSubmit={submit} aria-describedby={error ? ids.error : undefined}>
            <h1 className="text-3xl font-semibold tracking-tight">Mot de passe oublié</h1>
            <p className="mt-2 text-iron-300">
              Un lien pour en choisir un nouveau part à l’adresse du compte.
            </p>

            <div className="mt-8">
              <label htmlFor={ids.login} className="text-sm font-medium">
                Identifiant
              </label>
              <input
                id={ids.login}
                name="login"
                autoComplete="username"
                autoCapitalize="none"
                spellCheck={false}
                placeholder="Nom d’utilisateur ou adresse e-mail"
                value={login}
                onChange={(e) => setLogin(e.target.value)}
                aria-invalid={missing || undefined}
                aria-describedby={missing ? `${ids.login}-hint` : undefined}
                className={`${field} ${missing ? 'border-danger-400' : 'border-iron-700'}`}
              />
              {missing && (
                <p id={`${ids.login}-hint`} className="mt-1.5 text-sm text-danger-400">
                  Saisissez votre identifiant.
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
              {pending ? 'Envoi…' : 'Recevoir un lien'}
            </Button>
            {back}
          </form>
        )}
      </div>
    </AuthFrame>
  );
}
