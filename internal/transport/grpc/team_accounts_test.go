package grpc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

func TestTeamService_AccountActionsReservedToAdmins(t *testing.T) {
	f := startTeamServer(t, nil, "")
	admin := f.member(t, "admin")
	member := as(f.ts.memberToken)

	_, err := f.team.SuspendMember(member, &forgenetv1.SuspendMemberRequest{UserId: admin.GetId()})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = f.team.ResumeMember(member, &forgenetv1.ResumeMemberRequest{UserId: admin.GetId()})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = f.team.CloseMemberSessions(member, &forgenetv1.CloseMemberSessionsRequest{UserId: admin.GetId()})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = f.team.DeleteMember(member, &forgenetv1.DeleteMemberRequest{UserId: admin.GetId(), ConfirmUsername: "admin"})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
}

func TestTeamService_SuspendAndResume(t *testing.T) {
	f := startTeamServer(t, nil, "")
	ctx := context.Background()
	membre := f.member(t, "membre")

	res, err := f.team.SuspendMember(as(f.ts.adminToken), &forgenetv1.SuspendMemberRequest{UserId: membre.GetId()})
	require.NoError(t, err)
	assert.Equal(t, forgenetv1.MemberStatus_MEMBER_STATUS_SUSPENDED, res.GetMember().GetStatus())
	assert.NotNil(t, res.GetMember().GetSuspendedAt())
	assert.Equal(t, int64(1), res.GetClosedSessions(), "son jeton d'API")

	// Le jeton du membre est refusé dès l'appel suivant.
	_, err = f.team.ListMembers(as(f.ts.memberToken), &forgenetv1.ListMembersRequest{})
	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	// Même refus qu'un mauvais mot de passe.
	_, suspendedErr := f.auth.Login(ctx, &forgenetv1.LoginRequest{Login: "membre", Password: "motdepasse-de-test"})
	_, wrongErr := f.auth.Login(ctx, &forgenetv1.LoginRequest{Login: "membre", Password: "mauvais-mot-de-passe"})
	require.Error(t, suspendedErr)
	assert.Equal(t, status.Code(wrongErr), status.Code(suspendedErr))
	assert.Equal(t, status.Convert(wrongErr).Message(), status.Convert(suspendedErr).Message())

	assert.Equal(t, forgenetv1.MemberStatus_MEMBER_STATUS_SUSPENDED, f.member(t, "membre").GetStatus())

	resumed, err := f.team.ResumeMember(as(f.ts.adminToken), &forgenetv1.ResumeMemberRequest{UserId: membre.GetId()})
	require.NoError(t, err)
	assert.Equal(t, forgenetv1.MemberStatus_MEMBER_STATUS_ACTIVE, resumed.GetMember().GetStatus())
	_, err = f.auth.Login(ctx, &forgenetv1.LoginRequest{Login: "membre", Password: "motdepasse-de-test"})
	require.NoError(t, err)
}

func TestTeamService_SuspendedInvitationCannotBeRenewed(t *testing.T) {
	f := startTeamServer(t, nil, "")
	invited := f.invite(t, "nadia").GetMember()

	_, err := f.team.SuspendMember(as(f.ts.adminToken), &forgenetv1.SuspendMemberRequest{UserId: invited.GetId()})
	require.NoError(t, err)
	_, err = f.team.RenewInvitation(as(f.ts.adminToken), &forgenetv1.RenewInvitationRequest{UserId: invited.GetId()})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))

	// Réactivé, le compte n'a plus de lien : son invitation est à refaire.
	res, err := f.team.ResumeMember(as(f.ts.adminToken), &forgenetv1.ResumeMemberRequest{UserId: invited.GetId()})
	require.NoError(t, err)
	assert.Equal(t, forgenetv1.MemberStatus_MEMBER_STATUS_INVITATION_EXPIRED, res.GetMember().GetStatus())
}

func TestTeamService_GuardsOnOwnAccount(t *testing.T) {
	f := startTeamServer(t, nil, "")
	admin := f.member(t, "admin")

	_, err := f.team.SuspendMember(as(f.ts.adminToken), &forgenetv1.SuspendMemberRequest{UserId: admin.GetId()})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	_, err = f.team.DeleteMember(as(f.ts.adminToken), &forgenetv1.DeleteMemberRequest{UserId: admin.GetId(), ConfirmUsername: "admin"})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))

	_, err = f.team.SuspendMember(as(f.ts.adminToken), &forgenetv1.SuspendMemberRequest{UserId: "usr_absent"})
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestTeamService_CloseSessionsKeepsTheCallersOwn(t *testing.T) {
	f := startTeamServer(t, nil, "")
	ctx := context.Background()
	admin := f.member(t, "admin")
	other, err := f.auth.Login(ctx, &forgenetv1.LoginRequest{Login: "admin", Password: "motdepasse-de-test"})
	require.NoError(t, err)

	res, err := f.team.CloseMemberSessions(as(f.ts.adminToken), &forgenetv1.CloseMemberSessionsRequest{UserId: admin.GetId()})
	require.NoError(t, err)
	assert.Equal(t, int64(1), res.GetClosedSessions())

	_, err = f.team.ListMembers(as(f.ts.adminToken), &forgenetv1.ListMembersRequest{})
	require.NoError(t, err, "le jeton de l'appel reste ouvert")
	_, err = f.team.ListMembers(as(other.GetToken()), &forgenetv1.ListMembersRequest{})
	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	// Sur un autre compte, tout est fermé, et le compte reste actif.
	membre := f.member(t, "membre")
	res, err = f.team.CloseMemberSessions(as(f.ts.adminToken), &forgenetv1.CloseMemberSessionsRequest{UserId: membre.GetId()})
	require.NoError(t, err)
	assert.Equal(t, int64(1), res.GetClosedSessions())
	assert.Equal(t, forgenetv1.MemberStatus_MEMBER_STATUS_ACTIVE, f.member(t, "membre").GetStatus())
}

func TestTeamService_DeleteMember(t *testing.T) {
	f := startTeamServer(t, nil, "")
	invited := f.invite(t, "nadia").GetMember()

	for _, typed := range []string{"", "Nadia", "nadi"} {
		_, err := f.team.DeleteMember(as(f.ts.adminToken), &forgenetv1.DeleteMemberRequest{UserId: invited.GetId(), ConfirmUsername: typed})
		assert.Equal(t, codes.InvalidArgument, status.Code(err), typed)
	}

	// Une invitation en attente s'annule en supprimant le compte.
	_, err := f.team.DeleteMember(as(f.ts.adminToken), &forgenetv1.DeleteMemberRequest{UserId: invited.GetId(), ConfirmUsername: "nadia"})
	require.NoError(t, err)

	res, err := f.team.ListMembers(as(f.ts.adminToken), &forgenetv1.ListMembersRequest{})
	require.NoError(t, err)
	for _, m := range res.GetMembers() {
		assert.NotEqual(t, "nadia", m.GetUsername())
	}

	_, err = f.team.DeleteMember(as(f.ts.adminToken), &forgenetv1.DeleteMemberRequest{UserId: invited.GetId(), ConfirmUsername: "nadia"})
	assert.Equal(t, codes.NotFound, status.Code(err))

	// Le journal garde la suppression, sans adresse e-mail.
	events, err := f.ts.store.Audit().List(context.Background(), 100, 0)
	require.NoError(t, err)
	var found bool
	for _, e := range events {
		if e.Action == "user.delete" {
			found = true
			assert.NotContains(t, e.Metadata, "@")
		}
	}
	assert.True(t, found)
}
