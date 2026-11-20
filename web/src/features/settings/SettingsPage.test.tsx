import type { MessageInitShape } from '@bufbuild/protobuf';
import { Code, ConnectError, createRouterTransport } from '@connectrpc/connect';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { memoryLocation } from 'wouter/memory-location';

import { App } from '../../app/App';
import { AuthService, UserRole } from '../../gen/forgenet/v1/auth_pb';
import {
  GetMailStatusResponseSchema,
  MailService,
  MailSource,
  MailTLS,
  type SaveMailSettingsRequest,
} from '../../gen/forgenet/v1/mail_pb';
import { sessionInterceptors, type TransportFactory } from '../../lib/api/client';

const serverRelay = {
  enabled: true,
  host: 'smtp.exemple.org',
  port: 587,
  tls: MailTLS.MAIL_TLS_STARTTLS,
  from: 'forgenet@exemple.org',
  username: 'forgenet@exemple.org',
  authenticated: true,
  source: MailSource.SERVER,
  editable: true,
};

interface Options {
  role?: UserRole;
  status?: MessageInitShape<typeof GetMailStatusResponseSchema>;
  sendError?: ConnectError;
  saveError?: ConnectError;
}

function renderSettings({
  role = UserRole.ADMIN,
  status = serverRelay,
  sendError,
  saveError,
}: Options = {}) {
  const calls = { send: 0, saves: [] as SaveMailSettingsRequest[], disable: 0, reset: 0 };
  const factory: TransportFactory = (session) =>
    createRouterTransport(
      ({ service }) => {
        service(AuthService, {
          getCurrentUser: () => ({
            user: { id: 'u1', username: 'ana', email: 'ana@forgenet.test', role },
          }),
          logout: () => ({}),
        });
        service(MailService, {
          getMailStatus: () => status,
          sendTestMail: () => {
            calls.send++;
            if (sendError) {
              throw sendError;
            }
            return { to: 'ana@forgenet.test' };
          },
          saveMailSettings: (req) => {
            calls.saves.push(req);
            if (saveError) {
              throw saveError;
            }
            return {
              to: 'ana@forgenet.test',
              status: {
                ...serverRelay,
                host: req.host,
                port: req.port,
                tls: req.tls,
                username: req.username,
                from: req.from,
                authenticated: req.username !== '',
                source: MailSource.DASHBOARD,
                updatedBy: 'ana',
              },
            };
          },
          disableMail: () => {
            calls.disable++;
            return {
              status: {
                ...serverRelay,
                enabled: false,
                disabled: true,
                source: MailSource.DASHBOARD,
              },
            };
          },
          resetMailSettings: () => {
            calls.reset++;
            return { status: serverRelay };
          },
        });
      },
      { transport: { interceptors: sessionInterceptors(session) } }
    );
  sessionStorage.setItem(
    'forgenet.session',
    JSON.stringify({ token: 'jeton', expiresAt: Date.now() + 3_600_000 })
  );
  render(
    <App transportFactory={factory} locationHook={memoryLocation({ path: '/reglages' }).hook} />
  );
  return calls;
}

async function openForm(name = 'Modifier le relais') {
  fireEvent.click(await screen.findByRole('button', { name }));
  return screen.getByRole('form', { name: 'Réglages du relais SMTP' });
}

describe('SettingsPage', () => {
  beforeEach(() => sessionStorage.clear());

  it('montre le relais sans son mot de passe, et d’où viennent les réglages', async () => {
    renderSettings();
    expect(await screen.findByText('smtp.exemple.org:587')).toBeInTheDocument();
    expect(screen.getByText('STARTTLS, exigé')).toBeInTheDocument();
    expect(
      screen.getByText('forgenet@exemple.org, mot de passe gardé sur le serveur')
    ).toBeInTheDocument();
    expect(screen.getByText('Réglé sur le serveur (install.sh)')).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Revenir aux réglages du serveur' })).toBeNull();
  });

  it('envoie un e-mail d’essai et dit à qui', async () => {
    const calls = renderSettings();
    fireEvent.click(await screen.findByRole('button', { name: 'Envoyer un e-mail d’essai' }));
    expect(await screen.findByText(/Le relais a pris l’e-mail d’essai/)).toBeInTheDocument();
    expect(screen.getByText('ana@forgenet.test')).toBeInTheDocument();
    expect(calls.send).toBe(1);
  });

  it('montre la réponse du relais quand il refuse', async () => {
    renderSettings({
      sendError: new ConnectError(
        'le relais n’a pas pris le message : authentification refusée par le relais : 535 « identifiants refusés »',
        Code.FailedPrecondition
      ),
    });
    fireEvent.click(await screen.findByRole('button', { name: 'Envoyer un e-mail d’essai' }));
    expect(await screen.findByRole('alert')).toHaveTextContent('535');
  });

  it('enregistre et essaie un nouveau relais, en gardant le mot de passe laissé vide', async () => {
    const calls = renderSettings();
    const form = await openForm();
    const host = within(form).getByLabelText('Relais');
    expect(host).toHaveValue('smtp.exemple.org');
    expect(within(form).getByLabelText('Mot de passe du relais')).toHaveValue('');

    fireEvent.change(host, { target: { value: ' SMTP.Nouveau.org ' } });
    fireEvent.click(within(form).getByLabelText(/TLS dès la connexion/));
    expect(within(form).getByLabelText('Port')).toHaveValue('465');
    fireEvent.click(within(form).getByRole('button', { name: 'Enregistrer et essayer' }));

    expect(await screen.findByText(/Réglages enregistrés et appliqués/)).toBeInTheDocument();
    expect(calls.saves).toHaveLength(1);
    expect(calls.saves[0]).toMatchObject({
      host: 'smtp.nouveau.org',
      port: 465,
      tls: MailTLS.MAIL_TLS_IMPLICIT,
      keepPassword: true,
      password: '',
    });
    expect(screen.getByText('smtp.nouveau.org:465')).toBeInTheDocument();
    expect(screen.getByText(/Réglé ici par ana/)).toBeInTheDocument();
    await waitFor(() =>
      expect(screen.getByRole('button', { name: 'Modifier le relais' })).toHaveFocus()
    );
  });

  it('remplace le mot de passe quand il est saisi', async () => {
    const calls = renderSettings();
    const form = await openForm();
    fireEvent.change(within(form).getByLabelText('Mot de passe du relais'), {
      target: { value: 'nouveau-secret' },
    });
    fireEvent.click(within(form).getByRole('button', { name: 'Enregistrer et essayer' }));
    await screen.findByText(/Réglages enregistrés et appliqués/);
    expect(calls.saves[0]).toMatchObject({ keepPassword: false, password: 'nouveau-secret' });
  });

  it('garde le formulaire et la réponse du relais quand l’essai est refusé', async () => {
    const calls = renderSettings({
      saveError: new ConnectError(
        'le relais n’a pas pris le message : 535 « identifiants refusés » ; les réglages en place restent',
        Code.FailedPrecondition
      ),
    });
    const form = await openForm();
    fireEvent.click(within(form).getByRole('button', { name: 'Enregistrer et essayer' }));
    expect(await within(form).findByRole('alert')).toHaveTextContent(
      'les réglages en place restent'
    );
    expect(calls.saves).toHaveLength(1);
    expect(within(form).getByLabelText('Relais')).toHaveValue('smtp.exemple.org');
  });

  it('signale les saisies invalides sans rien envoyer', async () => {
    const calls = renderSettings({ status: { editable: true, source: MailSource.SERVER } });
    const form = await openForm('Régler un relais');
    fireEvent.change(within(form).getByLabelText(/Identifiant/), {
      target: { value: 'forgenet' },
    });
    fireEvent.click(within(form).getByLabelText(/Aucune/));
    fireEvent.change(within(form).getByLabelText('Relais'), {
      target: { value: 'smtp.exemple.org' },
    });
    fireEvent.click(within(form).getByRole('button', { name: 'Enregistrer et essayer' }));
    expect(within(form).getByText(/réservé à un relais sur ce serveur/)).toBeInTheDocument();
    expect(within(form).getByText('Saisissez le mot de passe du relais.')).toBeInTheDocument();
    expect(
      within(form).getByText('Indiquez l’adresse d’expéditeur des e-mails.')
    ).toBeInTheDocument();
    expect(calls.saves).toHaveLength(0);
  });

  it('coupe l’envoi après confirmation, puis le rouvre', async () => {
    const calls = renderSettings();
    fireEvent.click(await screen.findByRole('button', { name: 'Couper l’envoi' }));
    const confirm = screen.getByRole('group', { name: 'Couper l’envoi d’e-mails ?' });
    fireEvent.click(within(confirm).getByRole('button', { name: 'Couper l’envoi' }));
    expect(await screen.findByText(/L’envoi est coupé depuis cette page/)).toBeInTheDocument();
    expect(calls.disable).toBe(1);
    expect(screen.queryByRole('button', { name: 'Envoyer un e-mail d’essai' })).toBeNull();

    await openForm('Rouvrir l’envoi');
    fireEvent.click(screen.getByRole('button', { name: 'Enregistrer et essayer' }));
    expect(await screen.findByText(/Réglages enregistrés et appliqués/)).toBeInTheDocument();
    expect(calls.saves[0]).toMatchObject({ keepPassword: true });
  });

  it('revient aux réglages du serveur après confirmation', async () => {
    const calls = renderSettings({
      status: { ...serverRelay, source: MailSource.DASHBOARD, updatedBy: 'ana' },
    });
    fireEvent.click(await screen.findByRole('button', { name: 'Revenir aux réglages du serveur' }));
    const confirm = screen.getByRole('group', { name: 'Revenir aux réglages du serveur ?' });
    fireEvent.click(
      within(confirm).getByRole('button', { name: 'Revenir aux réglages du serveur' })
    );
    expect(
      await screen.findByText('Les réglages d’install.sh s’appliquent de nouveau.')
    ).toBeInTheDocument();
    expect(calls.reset).toBe(1);
    expect(screen.getByText('Réglé sur le serveur (install.sh)')).toBeInTheDocument();
  });

  it('dit quand le mot de passe enregistré ne s’ouvre plus, et le redemande', async () => {
    renderSettings({ status: { ...serverRelay, passwordUnreadable: true, authenticated: true } });
    expect(await screen.findByText(/ne s’ouvre plus/)).toBeInTheDocument();
    const form = await openForm();
    fireEvent.click(within(form).getByRole('button', { name: 'Enregistrer et essayer' }));
    expect(within(form).getByText('Saisissez le mot de passe du relais.')).toBeInTheDocument();
    fireEvent.click(within(form).getByRole('button', { name: 'Annuler' }));
    expect(
      screen.getByRole('button', { name: 'Revenir aux réglages du serveur' })
    ).toBeInTheDocument();
  });

  it('sans clé sur le serveur, montre le relais sans permettre de le changer', async () => {
    renderSettings({ status: { ...serverRelay, editable: false } });
    expect(await screen.findByText(/pas de clé FORGENET_SECRET_KEY/)).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Modifier le relais' })).toBeNull();
    expect(screen.queryByRole('button', { name: 'Couper l’envoi' })).toBeNull();
  });

  it('explique comment activer l’envoi quand aucun relais n’est réglé', async () => {
    renderSettings({ status: {} });
    expect(await screen.findByText(/Aucun relais SMTP n’est réglé/)).toBeInTheDocument();
    expect(screen.getByText(/--smtp-host smtp.exemple.org/)).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Envoyer un e-mail d’essai' })).toBeNull();
  });

  it('reste fermée aux membres, sans appel au serveur', async () => {
    const calls = renderSettings({ role: UserRole.MEMBER });
    expect(
      await screen.findByText('Les réglages du serveur sont réservés aux administrateurs.')
    ).toBeInTheDocument();
    expect(calls.send).toBe(0);
  });
});
