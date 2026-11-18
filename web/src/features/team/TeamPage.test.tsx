import { create, type MessageInitShape } from '@bufbuild/protobuf';
import { timestampFromDate } from '@bufbuild/protobuf/wkt';
import { Code, ConnectError, createRouterTransport } from '@connectrpc/connect';
import { fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { memoryLocation } from 'wouter/memory-location';

import { App } from '../../app/App';
import { AuthService, UserRole } from '../../gen/forgenet/v1/auth_pb';
import {
  MemberSchema,
  MemberStatus,
  TeamService,
  type InviteMemberRequest,
  type Member,
} from '../../gen/forgenet/v1/team_pb';
import { sessionInterceptors, type TransportFactory } from '../../lib/api/client';

const inSevenDays = new Date(Date.now() + 7 * 24 * 3_600_000);
const link = 'http://127.0.0.1:8080/acces#' + 'a'.repeat(43);

const member = (m: MessageInitShape<typeof MemberSchema>) =>
  create(MemberSchema, { role: UserRole.MEMBER, status: MemberStatus.ACTIVE, ...m });

const initialMembers = [
  member({ id: 'u1', username: 'ana', email: 'ana@forgenet.test', role: UserRole.ADMIN }),
  member({
    id: 'u2',
    username: 'bio',
    email: 'bio@exemple.org',
    status: MemberStatus.INVITED,
    invitationExpiresAt: timestampFromDate(inSevenDays),
  }),
  member({
    id: 'u3',
    username: 'eloi',
    email: 'eloi@exemple.org',
    status: MemberStatus.INVITATION_EXPIRED,
  }),
];

interface Options {
  role?: UserRole;
  mailEnabled?: boolean;
  inviteError?: ConnectError;
}

function renderTeam({ role = UserRole.ADMIN, mailEnabled = false, inviteError }: Options = {}) {
  const calls = { list: 0, invites: [] as InviteMemberRequest[], renewed: [] as string[] };
  const invited = (m: Member) => ({
    member: m,
    emailed: mailEnabled,
    invitationUrl: mailEnabled ? '' : link,
    localDashboard: true,
  });
  const factory: TransportFactory = (session) =>
    createRouterTransport(
      ({ service }) => {
        service(AuthService, {
          getCurrentUser: () => ({
            user: { id: 'u1', username: 'ana', email: 'ana@forgenet.test', role },
          }),
          logout: () => ({}),
        });
        service(TeamService, {
          listMembers: () => {
            calls.list++;
            return { members: initialMembers, mailEnabled };
          },
          inviteMember: (req) => {
            calls.invites.push(req);
            if (inviteError) {
              throw inviteError;
            }
            return {
              result: invited(
                member({
                  id: 'u9',
                  username: req.username,
                  email: req.email,
                  role: req.role,
                  status: MemberStatus.INVITED,
                  invitationExpiresAt: timestampFromDate(inSevenDays),
                })
              ),
            };
          },
          renewInvitation: (req) => {
            calls.renewed.push(req.userId);
            const before = initialMembers.find((m) => m.id === req.userId)!;
            return {
              result: invited({
                ...before,
                status: MemberStatus.INVITED,
                invitationExpiresAt: timestampFromDate(inSevenDays),
              }),
            };
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
    <App transportFactory={factory} locationHook={memoryLocation({ path: '/equipe' }).hook} />
  );
  return calls;
}

const row = (username: string) =>
  within(screen.getByRole('list', { name: 'Comptes du serveur' }))
    .getByText(username)
    .closest('li')!;

async function fillInvite(username: string, email: string) {
  fireEvent.change(await screen.findByLabelText('Nom d’utilisateur'), {
    target: { value: username },
  });
  fireEvent.change(screen.getByLabelText('Adresse e-mail'), { target: { value: email } });
}

describe('TeamPage', () => {
  beforeEach(() => sessionStorage.clear());
  afterEach(() => vi.unstubAllGlobals());

  it('réserve l’équipe aux administrateurs', async () => {
    const calls = renderTeam({ role: UserRole.MEMBER });
    expect(
      await screen.findByText('L’équipe est gérée par les administrateurs du serveur.')
    ).toBeInTheDocument();
    expect(calls.list).toBe(0);
  });

  it('montre chaque compte et où il en est', async () => {
    renderTeam();
    expect(await screen.findByRole('heading', { name: /3 comptes/ })).toHaveTextContent(
      '2 attendent leur mot de passe'
    );
    expect(within(row('ana')).getByText('(vous)')).toBeInTheDocument();
    expect(within(row('ana')).getByText('Administrateur')).toBeInTheDocument();
    expect(within(row('ana')).getByText('Actif')).toBeInTheDocument();
    expect(within(row('ana')).queryByRole('button')).not.toBeInTheDocument();
    expect(within(row('bio')).getByText(/Invité, lien valable jusqu’au/)).toBeInTheDocument();
    expect(within(row('eloi')).getByText('Invitation expirée')).toBeInTheDocument();
    expect(
      within(row('eloi')).getByRole('button', { name: 'Refaire le lien de eloi' })
    ).toBeInTheDocument();
  });

  it('sans relais, montre le lien une fois et ajoute le compte', async () => {
    const calls = renderTeam();
    expect(await screen.findByText(/le lien s’affichera ici une seule fois/)).toBeInTheDocument();
    await fillInvite(' awa.k ', 'awa@exemple.org');
    fireEvent.click(screen.getByRole('radio', { name: /^Administrateur/ }));
    fireEvent.click(screen.getByRole('button', { name: 'Créer le lien d’invitation' }));

    const reveal = await screen.findByRole('region', { name: 'Lien d’invitation de awa.k' });
    expect(within(reveal).getByLabelText('Lien d’invitation')).toHaveValue(link);
    expect(within(reveal).getByText(/ssh -N -L 8080/)).toBeInTheDocument();
    expect(calls.invites).toHaveLength(1);
    expect(calls.invites[0]).toMatchObject({
      username: 'awa.k',
      email: 'awa@exemple.org',
      role: UserRole.ADMIN,
    });
    expect(within(row('awa.k')).getByText(/Invité/)).toBeInTheDocument();
    expect(screen.getByLabelText('Nom d’utilisateur')).toHaveValue('');

    fireEvent.click(within(reveal).getByRole('button', { name: 'J’ai transmis le lien' }));
    expect(screen.queryByRole('region', { name: /Lien d’invitation/ })).not.toBeInTheDocument();
  });

  it('copie le lien dans le presse-papiers', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined);
    vi.stubGlobal('navigator', { ...navigator, clipboard: { writeText } });
    renderTeam();
    await fillInvite('awa', 'awa@exemple.org');
    fireEvent.click(await screen.findByRole('button', { name: 'Créer le lien d’invitation' }));
    fireEvent.click(await screen.findByRole('button', { name: 'Copier le lien' }));
    expect(await screen.findByText('Le lien est dans le presse-papiers.')).toBeInTheDocument();
    expect(writeText).toHaveBeenCalledWith(link);
  });

  it('avec un relais, dit où part l’invitation sans montrer le lien', async () => {
    renderTeam({ mailEnabled: true });
    await fillInvite('awa', 'awa@exemple.org');
    fireEvent.click(await screen.findByRole('button', { name: 'Envoyer l’invitation' }));
    expect(await screen.findByText(/L’invitation de awa part vers/)).toBeInTheDocument();
    expect(screen.queryByLabelText('Lien d’invitation')).not.toBeInTheDocument();
  });

  it('vérifie le formulaire avant d’appeler le serveur', async () => {
    const calls = renderTeam();
    await fillInvite('éloi k', 'pas-une-adresse');
    fireEvent.click(await screen.findByRole('button', { name: 'Créer le lien d’invitation' }));
    expect(await screen.findByText(/lettres sans accent/)).toBeInTheDocument();
    expect(screen.getByText(/Saisissez une adresse e-mail/)).toBeInTheDocument();
    expect(screen.getByLabelText('Nom d’utilisateur')).toHaveAttribute('aria-invalid', 'true');
    expect(calls.invites).toHaveLength(0);
  });

  it('montre le refus du serveur', async () => {
    renderTeam({
      inviteError: new ConnectError(
        'ce nom d’utilisateur ou cette adresse e-mail a déjà un compte',
        Code.AlreadyExists
      ),
    });
    await fillInvite('bio', 'bio@exemple.org');
    fireEvent.click(await screen.findByRole('button', { name: 'Créer le lien d’invitation' }));
    expect(await screen.findByRole('alert')).toHaveTextContent(/a déjà un compte/);
  });

  it('refait le lien d’une invitation expirée', async () => {
    const calls = renderTeam();
    fireEvent.click(await screen.findByRole('button', { name: 'Refaire le lien de eloi' }));
    expect(
      await screen.findByRole('region', { name: 'Lien d’invitation de eloi' })
    ).toBeInTheDocument();
    expect(calls.renewed).toEqual(['u3']);
    expect(within(row('eloi')).getByText(/Invité, lien valable jusqu’au/)).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: /3 comptes/ })).toHaveTextContent(
      '2 attendent leur mot de passe'
    );
  });
});
