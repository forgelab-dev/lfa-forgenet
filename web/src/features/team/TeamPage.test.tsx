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
  member({ id: 'u4', username: 'dan', email: 'dan@exemple.org' }),
  member({
    id: 'u5',
    username: 'zoe',
    email: 'zoe@exemple.org',
    status: MemberStatus.SUSPENDED,
    suspendedAt: timestampFromDate(new Date('2026-11-19T10:00:00Z')),
  }),
];

interface Options {
  role?: UserRole;
  mailEnabled?: boolean;
  inviteError?: ConnectError;
  accountError?: ConnectError;
}

function renderTeam({
  role = UserRole.ADMIN,
  mailEnabled = false,
  inviteError,
  accountError,
}: Options = {}) {
  const calls = {
    list: 0,
    invites: [] as InviteMemberRequest[],
    renewed: [] as string[],
    actions: [] as string[],
  };
  const byId = (id: string) => initialMembers.find((m) => m.id === id)!;
  const act = (name: string, id: string) => {
    calls.actions.push(`${name}:${id}`);
    if (accountError) {
      throw accountError;
    }
  };
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
          suspendMember: (req) => {
            act('suspend', req.userId);
            return {
              member: { ...byId(req.userId), status: MemberStatus.SUSPENDED },
              closedSessions: 2n,
            };
          },
          resumeMember: (req) => {
            act('resume', req.userId);
            return { member: { ...byId(req.userId), status: MemberStatus.ACTIVE } };
          },
          closeMemberSessions: (req) => {
            act('sessions', req.userId);
            return { closedSessions: 1n };
          },
          deleteMember: (req) => {
            act(`delete(${req.confirmUsername})`, req.userId);
            return {};
          },
          updateMemberEmail: (req) => {
            act(`email(${req.email})`, req.userId);
            const before = byId(req.userId);
            if (
              before.status === MemberStatus.INVITED ||
              before.status === MemberStatus.INVITATION_EXPIRED
            ) {
              const again = {
                ...before,
                email: req.email.toLowerCase(),
                status: MemberStatus.INVITED,
                invitationExpiresAt: timestampFromDate(inSevenDays),
              };
              return { member: again, invitation: invited(again) };
            }
            return { member: { ...before, email: req.email.toLowerCase() } };
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

// Le nom revient dans la confirmation d'une action : le premier est celui de
// la ligne.
const row = (username: string) =>
  within(screen.getByRole('list', { name: 'Comptes du serveur' }))
    .getAllByText(username)[0]
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
    expect(await screen.findByRole('heading', { name: /5 comptes/ })).toHaveTextContent(
      '2 attendent leur mot de passe · 1 suspendu'
    );
    expect(within(row('ana')).getByText('(vous)')).toBeInTheDocument();
    expect(within(row('ana')).getByText('Administrateur')).toBeInTheDocument();
    expect(within(row('ana')).getByText('Actif')).toBeInTheDocument();
    // Sur son propre compte : ni suspendre ni supprimer.
    expect(
      within(row('ana'))
        .getAllByRole('button')
        .map((b) => b.textContent)
    ).toEqual(['Fermer mes autres sessions']);
    expect(within(row('zoe')).getByText(/Suspendu depuis le/)).toBeInTheDocument();
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
    expect(screen.getByRole('heading', { name: /5 comptes/ })).toHaveTextContent(
      '2 attendent leur mot de passe'
    );
  });

  it('ferme ses autres sessions en gardant la sienne', async () => {
    const calls = renderTeam();
    fireEvent.click(await screen.findByRole('button', { name: 'Fermer mes autres sessions' }));
    expect(
      await within(row('ana')).findByText('Une session fermée ; la vôtre reste ouverte.')
    ).toBeInTheDocument();
    expect(calls.actions).toEqual(['sessions:u1']);
  });

  it('suspend un compte après confirmation', async () => {
    const calls = renderTeam();
    fireEvent.click(await screen.findByRole('button', { name: 'Suspendre dan' }));
    const confirm = screen.getByRole('group', { name: 'Suspendre dan' });
    expect(within(confirm).getByText(/vous pourrez le réactiver/)).toBeInTheDocument();
    expect(calls.actions).toEqual([]);

    fireEvent.click(within(confirm).getByRole('button', { name: 'Suspendre dan' }));
    expect(
      await within(row('dan')).findByText('dan est suspendu ; 2 sessions fermées.')
    ).toBeInTheDocument();
    expect(within(row('dan')).getByText('Suspendu')).toBeInTheDocument();
    expect(within(row('dan')).getByRole('button', { name: 'Réactiver dan' })).toBeInTheDocument();
    expect(calls.actions).toEqual(['suspend:u4']);
  });

  it('réactive un compte suspendu', async () => {
    const calls = renderTeam();
    fireEvent.click(await screen.findByRole('button', { name: 'Réactiver zoe' }));
    expect(
      await within(row('zoe')).findByText('zoe peut de nouveau se connecter.')
    ).toBeInTheDocument();
    expect(within(row('zoe')).getByText('Actif')).toBeInTheDocument();
    expect(calls.actions).toEqual(['resume:u5']);
  });

  it('supprime un compte une fois son nom retapé', async () => {
    const calls = renderTeam();
    fireEvent.click(await screen.findByRole('button', { name: 'Supprimer dan' }));
    const confirm = screen.getByRole('group', { name: 'Supprimer dan' });
    const submit = within(confirm).getByRole('button', { name: 'Supprimer dan' });
    expect(submit).toBeDisabled();

    fireEvent.change(within(confirm).getByLabelText(/pour confirmer/), {
      target: { value: 'da' },
    });
    expect(submit).toBeDisabled();
    fireEvent.change(within(confirm).getByLabelText(/pour confirmer/), {
      target: { value: 'dan' },
    });
    fireEvent.click(submit);

    expect(
      await screen.findByText('Le compte dan est supprimé. Le journal d’audit garde ses entrées.')
    ).toBeInTheDocument();
    expect(
      within(screen.getByRole('list', { name: 'Comptes du serveur' })).queryByText('dan')
    ).not.toBeInTheDocument();
    expect(screen.getByRole('heading', { name: /4 comptes/ })).toBeInTheDocument();
    expect(calls.actions).toEqual(['delete(dan):u4']);
  });

  it('annule une invitation en attente', async () => {
    renderTeam();
    fireEvent.click(await screen.findByRole('button', { name: 'Annuler l’invitation de bio' }));
    const confirm = screen.getByRole('group', { name: 'Supprimer bio' });
    fireEvent.change(within(confirm).getByLabelText(/pour confirmer/), {
      target: { value: 'bio' },
    });
    fireEvent.click(within(confirm).getByRole('button', { name: 'Annuler l’invitation' }));
    expect(
      await screen.findByText('L’invitation de bio est annulée : son lien ne vaut plus.')
    ).toBeInTheDocument();
  });

  it('montre le refus du serveur dans la confirmation', async () => {
    renderTeam({
      accountError: new ConnectError(
        'c’est le dernier administrateur actif : nommez-en un autre avant',
        Code.FailedPrecondition
      ),
    });
    fireEvent.click(await screen.findByRole('button', { name: 'Suspendre dan' }));
    const confirm = screen.getByRole('group', { name: 'Suspendre dan' });
    fireEvent.click(within(confirm).getByRole('button', { name: 'Suspendre dan' }));
    expect(await within(confirm).findByRole('alert')).toHaveTextContent(
      /dernier administrateur actif/
    );
    fireEvent.click(within(confirm).getByRole('button', { name: 'Garder le compte' }));
    expect(screen.queryByRole('group', { name: 'Suspendre dan' })).not.toBeInTheDocument();
  });

  it('corrige l’adresse d’un compte actif et prévient l’ancienne', async () => {
    const calls = renderTeam({ mailEnabled: true });
    await screen.findByRole('heading', { name: /5 comptes/ });

    fireEvent.click(within(row('dan')).getByRole('button', { name: 'Changer l’adresse de dan' }));
    const field = screen.getByLabelText('Nouvelle adresse de dan');
    expect(field).toHaveFocus();
    expect(
      within(row('dan')).getByText(/un avis partira à l’ancienne adresse/)
    ).toBeInTheDocument();

    fireEvent.click(within(row('dan')).getByRole('button', { name: 'Enregistrer l’adresse' }));
    expect(within(row('dan')).getByText('Saisissez la nouvelle adresse.')).toBeInTheDocument();
    fireEvent.change(field, { target: { value: 'DAN@exemple.org' } });
    fireEvent.click(within(row('dan')).getByRole('button', { name: 'Enregistrer l’adresse' }));
    expect(within(row('dan')).getByText('C’est déjà l’adresse de ce compte.')).toBeInTheDocument();
    expect(calls.actions).toEqual([]);

    fireEvent.change(field, { target: { value: ' Dan@Nouveau.org ' } });
    fireEvent.click(within(row('dan')).getByRole('button', { name: 'Enregistrer l’adresse' }));

    expect(await within(row('dan')).findByText('dan@nouveau.org')).toBeInTheDocument();
    expect(
      within(row('dan')).getByText('Adresse de dan corrigée. Un avis part à l’ancienne adresse.')
    ).toHaveAttribute('role', 'status');
    expect(calls.actions).toEqual(['email(Dan@Nouveau.org):u4']);
  });

  it('réinvite à la nouvelle adresse un compte qui attendait son mot de passe', async () => {
    renderTeam();
    await screen.findByRole('heading', { name: /5 comptes/ });

    fireEvent.click(within(row('eloi')).getByRole('button', { name: 'Changer l’adresse de eloi' }));
    expect(
      within(row('eloi')).getByText(/un nouveau lien sera fait, à lui transmettre/)
    ).toBeInTheDocument();
    fireEvent.change(screen.getByLabelText('Nouvelle adresse de eloi'), {
      target: { value: 'eloi@bonne.org' },
    });
    fireEvent.click(within(row('eloi')).getByRole('button', { name: 'Enregistrer l’adresse' }));

    expect(await within(row('eloi')).findByText('eloi@bonne.org')).toBeInTheDocument();
    expect(within(row('eloi')).getByText(/son ancien lien ne vaut plus/)).toBeInTheDocument();
    expect(within(row('eloi')).getByText(/Invité, lien valable jusqu’au/)).toBeInTheDocument();
    expect(within(row('eloi')).getByDisplayValue(link)).toBeInTheDocument();
  });

  it('montre le refus du serveur et garde la saisie', async () => {
    renderTeam({
      accountError: new ConnectError('cette adresse e-mail a déjà un compte', Code.AlreadyExists),
    });
    await screen.findByRole('heading', { name: /5 comptes/ });

    fireEvent.click(within(row('dan')).getByRole('button', { name: 'Changer l’adresse de dan' }));
    fireEvent.change(screen.getByLabelText('Nouvelle adresse de dan'), {
      target: { value: 'bio@exemple.org' },
    });
    fireEvent.click(within(row('dan')).getByRole('button', { name: 'Enregistrer l’adresse' }));

    expect(await within(row('dan')).findByRole('alert')).toHaveTextContent(/déjà un compte/);
    expect(screen.getByLabelText('Nouvelle adresse de dan')).toHaveValue('bio@exemple.org');
    fireEvent.click(within(row('dan')).getByRole('button', { name: 'Annuler' }));
    expect(screen.queryByLabelText('Nouvelle adresse de dan')).not.toBeInTheDocument();
    expect(within(row('dan')).getByText('dan@exemple.org')).toBeInTheDocument();
  });
});
