package grpc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

func TestAuthService_ChangeEmail(t *testing.T) {
	ts, client, outbox := startResetServer(t, "forgenet.exemple.org")
	member := as(ts.memberToken)

	_, err := client.ChangeEmail(member, &forgenetv1.ChangeEmailRequest{CurrentPassword: "pas-le-bon", NewEmail: "membre@nouveau.org"})
	assert.Equal(t, codes.InvalidArgument, status.Code(err), "jamais UNAUTHENTICATED, qui fermerait la session")
	assert.Contains(t, status.Convert(err).Message(), "mot de passe actuel")

	for _, c := range []struct {
		email string
		code  codes.Code
	}{
		{"pas-une-adresse", codes.InvalidArgument},
		{"Membre <membre@nouveau.org>", codes.InvalidArgument},
		{"membre@forgenet.test", codes.InvalidArgument},
		{"admin@forgenet.test", codes.AlreadyExists},
	} {
		_, err = client.ChangeEmail(member, &forgenetv1.ChangeEmailRequest{CurrentPassword: "motdepasse-de-test", NewEmail: c.email})
		assert.Equal(t, c.code, status.Code(err), c.email)
	}
	assert.Empty(t, outbox.Messages(), "aucun refus n'envoie d'avis")

	res, err := client.ChangeEmail(member, &forgenetv1.ChangeEmailRequest{CurrentPassword: "motdepasse-de-test", NewEmail: " Membre@Nouveau.org "})
	require.NoError(t, err)
	assert.Equal(t, "membre@nouveau.org", res.GetUser().GetEmail())

	cur, err := client.GetCurrentUser(member, &forgenetv1.GetCurrentUserRequest{})
	require.NoError(t, err, "la session reste ouverte")
	assert.Equal(t, "membre@nouveau.org", cur.GetUser().GetEmail())

	msgs := outbox.Messages()
	require.Len(t, msgs, 1)
	assert.Equal(t, "membre@forgenet.test", msgs[0].To, "l'avis part à l'ancienne adresse")
	assert.Contains(t, msgs[0].Body, "Vous venez de changer l'adresse")
	assert.Contains(t, msgs[0].Body, "m***@nouveau.org")
}

func TestAuthService_ChangeEmailWithoutRelay(t *testing.T) {
	ts, client := startAuthServer(t)
	res, err := client.ChangeEmail(as(ts.adminToken), &forgenetv1.ChangeEmailRequest{CurrentPassword: "motdepasse-de-test", NewEmail: "chef@exemple.org"})
	require.NoError(t, err)
	assert.Equal(t, "chef@exemple.org", res.GetUser().GetEmail())
}

func TestAuthService_ChangeEmailIsThrottled(t *testing.T) {
	ts, client := startAuthServer(t)
	member := as(ts.memberToken)
	for range 10 {
		_, err := client.ChangePassword(member, &forgenetv1.ChangePasswordRequest{CurrentPassword: "pas-le-bon", NewPassword: "nouveau-mot-de-passe"})
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	}
	// Même compteur que ChangePassword : deviner le mot de passe par l'un ou
	// par l'autre revient au même.
	_, err := client.ChangeEmail(member, &forgenetv1.ChangeEmailRequest{CurrentPassword: "motdepasse-de-test", NewEmail: "membre@nouveau.org"})
	assert.Equal(t, codes.ResourceExhausted, status.Code(err))
}

func TestTeamService_UpdateMemberEmail_ActiveAccount(t *testing.T) {
	outbox := &transportGRPC.SyncOutboxForTests{}
	f := startTeamServer(t, outbox, "forgenet.exemple.org")
	membre := f.member(t, "membre")

	_, err := f.team.UpdateMemberEmail(as(f.ts.memberToken), &forgenetv1.UpdateMemberEmailRequest{UserId: membre.GetId(), Email: "x@exemple.org"})
	assert.Equal(t, codes.PermissionDenied, status.Code(err), "réservé aux administrateurs")

	res, err := f.team.UpdateMemberEmail(as(f.ts.adminToken), &forgenetv1.UpdateMemberEmailRequest{UserId: membre.GetId(), Email: "Membre@Exemple.org"})
	require.NoError(t, err)
	assert.Equal(t, "membre@exemple.org", res.GetMember().GetEmail())
	assert.Equal(t, forgenetv1.MemberStatus_MEMBER_STATUS_ACTIVE, res.GetMember().GetStatus())
	assert.Nil(t, res.GetInvitation(), "un compte actif n'a pas besoin d'invitation")

	msgs := outbox.Messages()
	require.Len(t, msgs, 1)
	assert.Equal(t, "membre@forgenet.test", msgs[0].To)
	assert.Contains(t, msgs[0].Body, "admin vient de changer l'adresse e-mail de votre compte ForgeNet membre.")
}

func TestTeamService_UpdateMemberEmail_InvitedAccountIsReinvited(t *testing.T) {
	f := startTeamServer(t, nil, "")
	first := f.invite(t, "nadia")
	oldToken := linkInURL.FindStringSubmatch(first.GetInvitationUrl())[1]

	res, err := f.team.UpdateMemberEmail(as(f.ts.adminToken), &forgenetv1.UpdateMemberEmailRequest{UserId: first.GetMember().GetId(), Email: "nadia@bonne.org"})
	require.NoError(t, err)
	assert.Equal(t, forgenetv1.MemberStatus_MEMBER_STATUS_INVITED, res.GetMember().GetStatus())
	require.NotNil(t, res.GetInvitation())
	assert.False(t, res.GetInvitation().GetEmailed())
	assert.Regexp(t, `^http://127\.0\.0\.1:8080/acces#`, res.GetInvitation().GetInvitationUrl(), "sans relais, le nouveau lien est rendu")

	// Le lien parti à l'adresse mal saisie ne vaut plus.
	_, err = f.auth.InspectAccessLink(context.Background(), &forgenetv1.InspectAccessLinkRequest{Token: oldToken})
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestTeamService_UpdateMemberEmail_InvitationByMail(t *testing.T) {
	outbox := &transportGRPC.SyncOutboxForTests{}
	f := startTeamServer(t, outbox, "forgenet.exemple.org")
	invited := f.invite(t, "nadia").GetMember()

	res, err := f.team.UpdateMemberEmail(as(f.ts.adminToken), &forgenetv1.UpdateMemberEmailRequest{UserId: invited.GetId(), Email: "nadia@bonne.org"})
	require.NoError(t, err)
	assert.True(t, res.GetInvitation().GetEmailed())

	msgs := outbox.Messages()
	require.Len(t, msgs, 2, "la première invitation, puis la nouvelle ; pas d'avis à l'adresse mal saisie")
	assert.Equal(t, "nadia@bonne.org", msgs[1].To)
	assert.Equal(t, "admin vous invite sur ForgeNet", msgs[1].Subject)
}

func TestTeamService_UpdateMemberEmail_Refusals(t *testing.T) {
	f := startTeamServer(t, nil, "")
	admin := f.member(t, "admin")
	membre := f.member(t, "membre")

	for _, c := range []struct {
		id, email string
		code      codes.Code
	}{
		{admin.GetId(), "chef@exemple.org", codes.FailedPrecondition},
		{membre.GetId(), "pas-une-adresse", codes.InvalidArgument},
		{membre.GetId(), "membre@forgenet.test", codes.InvalidArgument},
		{membre.GetId(), "admin@forgenet.test", codes.AlreadyExists},
		{"usr_absent", "libre@exemple.org", codes.NotFound},
	} {
		_, err := f.team.UpdateMemberEmail(as(f.ts.adminToken), &forgenetv1.UpdateMemberEmailRequest{UserId: c.id, Email: c.email})
		assert.Equal(t, c.code, status.Code(err), c.email)
	}
	_, err := f.team.UpdateMemberEmail(as(f.ts.adminToken), &forgenetv1.UpdateMemberEmailRequest{UserId: admin.GetId(), Email: "chef@exemple.org"})
	assert.Contains(t, status.Convert(err).Message(), "page Compte")
}

func TestTeamService_UpdateMemberEmail_SuspendedInvitationStaysSuspended(t *testing.T) {
	f := startTeamServer(t, nil, "")
	invited := f.invite(t, "nadia").GetMember()
	_, err := f.team.SuspendMember(as(f.ts.adminToken), &forgenetv1.SuspendMemberRequest{UserId: invited.GetId()})
	require.NoError(t, err)

	res, err := f.team.UpdateMemberEmail(as(f.ts.adminToken), &forgenetv1.UpdateMemberEmailRequest{UserId: invited.GetId(), Email: "nadia@bonne.org"})
	require.NoError(t, err)
	assert.Equal(t, forgenetv1.MemberStatus_MEMBER_STATUS_SUSPENDED, res.GetMember().GetStatus())
	assert.Nil(t, res.GetInvitation(), "pas d'invitation pour un compte suspendu")
}
