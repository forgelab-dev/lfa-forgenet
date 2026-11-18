import { CircleCheck, LoaderCircle, Mail, MailX, ShieldAlert } from 'lucide-react';
import { useCallback, useState } from 'react';

import { Button } from '../../components/ui/Button';
import { PageHeader } from '../../components/ui/PageHeader';
import { RequestError } from '../../components/ui/RequestError';
import { UserRole } from '../../gen/forgenet/v1/auth_pb';
import { MailTLS, type GetMailStatusResponse } from '../../gen/forgenet/v1/mail_pb';
import { describeError } from '../../lib/api/errors';
import { useRequest } from '../../lib/useRequest';
import { useApi, useCurrentUser } from '../auth/session';

const tlsLabel: Record<MailTLS, string> = {
  [MailTLS.MAIL_TLS_STARTTLS]: 'STARTTLS, exigé',
  [MailTLS.MAIL_TLS_IMPLICIT]: 'TLS dès la connexion',
  [MailTLS.MAIL_TLS_NONE]: 'aucune (relais local)',
  [MailTLS.MAIL_TLS_UNSPECIFIED]: 'inconnue',
};

/**
 * Réglages du serveur, réservés aux administrateurs. Ils se changent sur le
 * serveur, par install.sh : la page les montre et permet de les essayer
 * (ADR-0014).
 */
export function SettingsPage() {
  const user = useCurrentUser();
  return (
    <>
      <PageHeader
        title="Réglages"
        description="Ce que ce serveur sait faire au-delà de vos applications. Les réglages se changent sur le serveur, en relançant install.sh ; ici, on les vérifie."
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

function MailSection() {
  const api = useApi();
  const load = useCallback((signal: AbortSignal) => api.mail.getMailStatus({}, { signal }), [api]);
  const [request, retry] = useRequest(load);

  return (
    <section
      aria-labelledby="settings-mail"
      className="animate-rise rounded-lg border border-iron-800 bg-iron-900 px-5 py-5 sm:px-6"
    >
      <h2 id="settings-mail" className="flex items-center gap-2 text-lg font-semibold">
        <Mail aria-hidden="true" className="size-5 text-gold-400" />
        Envoi d’e-mails
      </h2>
      <div className="mt-4">
        {request.status === 'loading' && (
          <div role="status" className="h-24 animate-pulse rounded-md bg-iron-850">
            <span className="sr-only">Chargement des réglages d’envoi…</span>
          </div>
        )}
        {request.status === 'error' && <RequestError error={request.error} onRetry={retry} />}
        {request.status === 'success' &&
          (request.data.enabled ? <MailEnabled status={request.data} /> : <MailDisabled />)}
      </div>
    </section>
  );
}

function MailDisabled() {
  return (
    <>
      <p className="flex items-start gap-2 text-iron-300">
        <MailX aria-hidden="true" className="mt-0.5 size-4 shrink-0" />
        <span>
          Aucun relais SMTP n’est configuré : ForgeNet n’envoie aucun e-mail. Tout le reste marche,
          et un mot de passe perdu se remplace sur le serveur.
        </span>
      </p>
      <p className="mt-4 text-sm text-iron-300">
        Pour l’activer, relancez l’installation avec le relais de votre fournisseur de messagerie.
        Le mot de passe du relais est demandé, jamais passé en argument :
      </p>
      <pre className="mt-3 overflow-x-auto rounded-md bg-iron-950 px-3 py-2.5 font-mono text-xs text-iron-100">
        {
          'sudo bash install.sh \\\n  --smtp-host smtp.exemple.org \\\n  --smtp-user forgenet@exemple.org'
        }
      </pre>
    </>
  );
}

function MailEnabled({ status }: { status: GetMailStatusResponse }) {
  return (
    <>
      <dl className="grid gap-x-8 gap-y-3 text-sm sm:grid-cols-2">
        <Row label="Relais" value={`${status.host}:${status.port}`} mono />
        <Row label="Protection de la liaison" value={tlsLabel[status.tls]} />
        <Row label="Expéditeur" value={status.from} mono />
        <Row
          label="Identifiants"
          value={status.authenticated ? 'oui, gardés sur le serveur' : 'aucun'}
        />
      </dl>
      <TestMail />
    </>
  );
}

function Row({ label, value, mono = false }: { label: string; value: string; mono?: boolean }) {
  return (
    <div className="min-w-0">
      <dt className="text-xs text-iron-300">{label}</dt>
      <dd className={`mt-0.5 break-all text-iron-100 ${mono ? 'font-mono' : ''}`}>{value}</dd>
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
