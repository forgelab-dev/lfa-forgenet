package grpc_test

import (
	"context"
	"regexp"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

var linkInURL = regexp.MustCompile(`/acces#([A-Za-z0-9_-]{43})`)

// teamFixture sert TeamService et AuthService ; outbox nil = sans relais.
// shift avance l'horloge du service, pour voir une invitation expirer.
type teamFixture struct {
	ts    *testServer
	team  forgenetv1.TeamServiceClient
	auth  forgenetv1.AuthServiceClient
	shift atomic.Int64
}

func startTeamServer(t *testing.T, outbox transportGRPC.Enqueuer, dashboardHost string) *teamFixture {
	t.Helper()
	f := &teamFixture{}
	f.ts = startTestServer(t, &fakeEngine{}, func(srv *grpc.Server, st store.Store, auditor audit.Auditor) {
		authSvc := auth.NewService(st, config.AuthConfig{JWTSecret: "secret-de-test-d-au-moins-32-caracteres"}, auditor)
		forgenetv1.RegisterAuthServiceServer(srv, transportGRPC.NewAuthServiceServer(authSvc, st.Users(), zerolog.Nop()))
		forgenetv1.RegisterTeamServiceServer(srv, transportGRPC.NewTeamServiceServer(transportGRPC.TeamDeps{
			Auth: authSvc, Users: st.Users(), Links: st.ActionTokens(), Outbox: outbox,
			DashboardHost: dashboardHost, Auditor: auditor,
			Now: func() time.Time { return time.Now().Add(time.Duration(f.shift.Load())) },
		}, zerolog.Nop()))
	})
	f.team = forgenetv1.NewTeamServiceClient(f.ts.conn)
	f.auth = forgenetv1.NewAuthServiceClient(f.ts.conn)
	return f
}

func (f *teamFixture) invite(t *testing.T, username string) *forgenetv1.InvitationResult {
	t.Helper()
	res, err := f.team.InviteMember(as(f.ts.adminToken), &forgenetv1.InviteMemberRequest{
		Username: username, Email: username + "@exemple.org", Role: forgenetv1.UserRole_USER_ROLE_MEMBER,
	})
	require.NoError(t, err)
	return res.GetResult()
}

func (f *teamFixture) member(t *testing.T, username string) *forgenetv1.Member {
	t.Helper()
	res, err := f.team.ListMembers(as(f.ts.adminToken), &forgenetv1.ListMembersRequest{})
	require.NoError(t, err)
	for _, m := range res.GetMembers() {
		if m.GetUsername() == username {
			return m
		}
	}
	t.Fatalf("membre %s absent de la liste", username)
	return nil
}

func TestTeamService_ReservedToAdmins(t *testing.T) {
	f := startTeamServer(t, nil, "")
	_, err := f.team.ListMembers(as(f.ts.memberToken), &forgenetv1.ListMembersRequest{})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = f.team.InviteMember(as(f.ts.memberToken), &forgenetv1.InviteMemberRequest{
		Username: "intrus", Email: "intrus@exemple.org", Role: forgenetv1.UserRole_USER_ROLE_ADMIN,
	})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = f.team.ListMembers(context.Background(), &forgenetv1.ListMembersRequest{})
	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	res, err := f.team.ListMembers(as(f.ts.adminToken), &forgenetv1.ListMembersRequest{})
	require.NoError(t, err)
	assert.False(t, res.GetMailEnabled())
	require.Len(t, res.GetMembers(), 2)
	for _, m := range res.GetMembers() {
		assert.Equal(t, forgenetv1.MemberStatus_MEMBER_STATUS_ACTIVE, m.GetStatus(), m.GetUsername())
		assert.Nil(t, m.GetInvitationExpiresAt())
	}
}

func TestTeamService_InviteWithoutRelayReturnsTheLinkOnce(t *testing.T) {
	f := startTeamServer(t, nil, "")
	ctx := context.Background()

	res := f.invite(t, "bio")
	assert.False(t, res.GetEmailed())
	assert.True(t, res.GetLocalDashboard(), "sans domaine, le lien passe par le tunnel")
	assert.Equal(t, forgenetv1.MemberStatus_MEMBER_STATUS_INVITED, res.GetMember().GetStatus())
	assert.Equal(t, "bio@exemple.org", res.GetMember().GetEmail())
	require.Regexp(t, `^http://127\.0\.0\.1:8080/acces#`, res.GetInvitationUrl())
	token := linkInURL.FindStringSubmatch(res.GetInvitationUrl())[1]

	info, err := f.auth.InspectAccessLink(ctx, &forgenetv1.InspectAccessLinkRequest{Token: token})
	require.NoError(t, err)
	assert.Equal(t, forgenetv1.AccessLinkPurpose_ACCESS_LINK_PURPOSE_INVITATION, info.GetPurpose())
	assert.Equal(t, "bio", info.GetUsername())
	assert.WithinDuration(t, time.Now().Add(7*24*time.Hour), info.GetExpiresAt().AsTime(), time.Minute)

	listed := f.member(t, "bio")
	assert.Equal(t, forgenetv1.MemberStatus_MEMBER_STATUS_INVITED, listed.GetStatus())
	assert.WithinDuration(t, info.GetExpiresAt().AsTime(), listed.GetInvitationExpiresAt().AsTime(), time.Second)

	// Le compte ne s'ouvre pas sans son lien.
	_, err = f.auth.Login(ctx, &forgenetv1.LoginRequest{Login: "bio", Password: "!"})
	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	_, err = f.auth.RedeemAccessLink(ctx, &forgenetv1.RedeemAccessLinkRequest{Token: token, NewPassword: "un-mot-de-passe-solide"})
	require.NoError(t, err)
	assert.Equal(t, forgenetv1.MemberStatus_MEMBER_STATUS_ACTIVE, f.member(t, "bio").GetStatus())
	assert.Equal(t, forgenetv1.UserRole_USER_ROLE_MEMBER, f.member(t, "bio").GetRole())

	events, err := f.ts.store.Audit().List(ctx, 50, 0)
	require.NoError(t, err)
	var logged bool
	for _, e := range events {
		if e.Action == "team.invite" {
			logged = true
			assert.NotContains(t, e.Metadata, token, "le lien n'entre pas au journal")
			assert.Contains(t, e.Metadata, "lien rendu")
		}
	}
	assert.True(t, logged, "invitation inscrite au journal d'audit")
}

func TestTeamService_InviteByMail(t *testing.T) {
	outbox := &transportGRPC.SyncOutboxForTests{}
	f := startTeamServer(t, outbox, "forgenet.exemple.org")

	listed, err := f.team.ListMembers(as(f.ts.adminToken), &forgenetv1.ListMembersRequest{})
	require.NoError(t, err)
	assert.True(t, listed.GetMailEnabled())

	res, err := f.team.InviteMember(as(f.ts.adminToken), &forgenetv1.InviteMemberRequest{
		Username: "chef", Email: " Chef@Exemple.org ", Role: forgenetv1.UserRole_USER_ROLE_ADMIN,
	})
	require.NoError(t, err)
	assert.True(t, res.GetResult().GetEmailed())
	assert.Empty(t, res.GetResult().GetInvitationUrl(), "le lien part par e-mail, il n'est pas rendu")
	assert.False(t, res.GetResult().GetLocalDashboard())
	assert.Equal(t, forgenetv1.UserRole_USER_ROLE_ADMIN, res.GetResult().GetMember().GetRole())

	msgs := outbox.Messages()
	require.Len(t, msgs, 1)
	assert.Equal(t, "chef@exemple.org", msgs[0].To)
	assert.Equal(t, "admin vous invite sur ForgeNet", msgs[0].Subject)
	assert.Regexp(t, `https://forgenet\.exemple\.org/acces#[A-Za-z0-9_-]{43}`, msgs[0].Body)
}

func TestTeamService_FullOutboxFallsBackToTheLink(t *testing.T) {
	outbox := &transportGRPC.SyncOutboxForTests{}
	outbox.SetFull(true)
	f := startTeamServer(t, outbox, "forgenet.exemple.org")

	res := f.invite(t, "bio")
	assert.False(t, res.GetEmailed())
	assert.Regexp(t, `^https://forgenet\.exemple\.org/acces#`, res.GetInvitationUrl())
}

func TestTeamService_InviteRejectsBadInput(t *testing.T) {
	f := startTeamServer(t, nil, "")
	ctx := as(f.ts.adminToken)
	for name, req := range map[string]*forgenetv1.InviteMemberRequest{
		"nom vide":        {Username: " ", Email: "a@exemple.org", Role: forgenetv1.UserRole_USER_ROLE_MEMBER},
		"nom à espace":    {Username: "bio sika", Email: "a@exemple.org", Role: forgenetv1.UserRole_USER_ROLE_MEMBER},
		"nom accentué":    {Username: "éloi", Email: "a@exemple.org", Role: forgenetv1.UserRole_USER_ROLE_MEMBER},
		"adresse fausse":  {Username: "bio", Email: "pas-une-adresse", Role: forgenetv1.UserRole_USER_ROLE_MEMBER},
		"adresse nommée":  {Username: "bio", Email: "Bio <bio@exemple.org>", Role: forgenetv1.UserRole_USER_ROLE_MEMBER},
		"rôle non choisi": {Username: "bio", Email: "bio@exemple.org"},
	} {
		_, err := f.team.InviteMember(ctx, req)
		assert.Equal(t, codes.InvalidArgument, status.Code(err), name)
	}

	for name, req := range map[string]*forgenetv1.InviteMemberRequest{
		"nom pris":      {Username: "membre", Email: "autre@exemple.org", Role: forgenetv1.UserRole_USER_ROLE_MEMBER},
		"adresse prise": {Username: "autre", Email: "MEMBRE@forgenet.test", Role: forgenetv1.UserRole_USER_ROLE_MEMBER},
	} {
		_, err := f.team.InviteMember(ctx, req)
		assert.Equal(t, codes.AlreadyExists, status.Code(err), name)
	}
}

func TestTeamService_InvitationExpiresThenRenews(t *testing.T) {
	f := startTeamServer(t, nil, "")
	ctx := context.Background()
	first := f.invite(t, "bio")
	oldToken := linkInURL.FindStringSubmatch(first.GetInvitationUrl())[1]

	f.shift.Store(int64(8 * 24 * time.Hour))
	expired := f.member(t, "bio")
	assert.Equal(t, forgenetv1.MemberStatus_MEMBER_STATUS_INVITATION_EXPIRED, expired.GetStatus())
	assert.Nil(t, expired.GetInvitationExpiresAt())
	f.shift.Store(0)

	renewed, err := f.team.RenewInvitation(as(f.ts.adminToken), &forgenetv1.RenewInvitationRequest{UserId: expired.GetId()})
	require.NoError(t, err)
	newURL := renewed.GetResult().GetInvitationUrl()
	require.NotEmpty(t, newURL)
	assert.NotEqual(t, first.GetInvitationUrl(), newURL)

	_, err = f.auth.InspectAccessLink(ctx, &forgenetv1.InspectAccessLinkRequest{Token: oldToken})
	assert.Equal(t, codes.NotFound, status.Code(err), "le nouveau lien annule l'ancien")
	_, err = f.auth.InspectAccessLink(ctx, &forgenetv1.InspectAccessLinkRequest{Token: linkInURL.FindStringSubmatch(newURL)[1]})
	require.NoError(t, err)
}

func TestTeamService_RenewRefusesActiveOrUnknownAccounts(t *testing.T) {
	f := startTeamServer(t, nil, "")
	admin := as(f.ts.adminToken)

	_, err := f.team.RenewInvitation(admin, &forgenetv1.RenewInvitationRequest{UserId: f.member(t, "membre").GetId()})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	_, err = f.team.RenewInvitation(admin, &forgenetv1.RenewInvitationRequest{UserId: "usr_inconnu"})
	assert.Equal(t, codes.NotFound, status.Code(err))
	_, err = f.team.RenewInvitation(as(f.ts.memberToken), &forgenetv1.RenewInvitationRequest{UserId: "usr_inconnu"})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
}
