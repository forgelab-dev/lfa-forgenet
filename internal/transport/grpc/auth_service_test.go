package grpc_test

import (
	"context"
	"regexp"
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

// Mot de passe des comptes créés par startTestServer.
const testPassword = "motdepasse-de-test"

func startAuthServer(t *testing.T) (*testServer, forgenetv1.AuthServiceClient) {
	t.Helper()
	ts := startTestServer(t, &fakeEngine{}, func(srv *grpc.Server, st store.Store, auditor audit.Auditor) {
		authSvc := auth.NewService(st, config.AuthConfig{JWTSecret: "secret-de-test-d-au-moins-32-caracteres"}, auditor)
		forgenetv1.RegisterAuthServiceServer(srv, transportGRPC.NewAuthServiceServer(authSvc, st.Users(), zerolog.Nop()))
	})
	return ts, forgenetv1.NewAuthServiceClient(ts.conn)
}

func TestAuthService_LoginOpensAUsableSession(t *testing.T) {
	_, client := startAuthServer(t)

	for _, login := range []string{"membre", "membre@forgenet.test"} {
		res, err := client.Login(context.Background(), &forgenetv1.LoginRequest{Login: login, Password: testPassword})
		require.NoError(t, err, "connexion par %s", login)
		assert.NotEmpty(t, res.GetToken())
		assert.WithinDuration(t, time.Now().Add(config.DefaultTokenTTL), res.GetExpiresAt().AsTime(), time.Minute)
		assert.Equal(t, forgenetv1.UserRole_USER_ROLE_MEMBER, res.GetUser().GetRole())

		me, err := client.GetCurrentUser(as(res.GetToken()), &forgenetv1.GetCurrentUserRequest{})
		require.NoError(t, err)
		assert.Equal(t, "membre", me.GetUser().GetUsername())
		assert.Equal(t, "membre@forgenet.test", me.GetUser().GetEmail())
	}
}

func TestAuthService_LoginRefusalDoesNotRevealAccounts(t *testing.T) {
	_, client := startAuthServer(t)

	_, wrongPassword := client.Login(context.Background(), &forgenetv1.LoginRequest{Login: "membre", Password: "faux"})
	_, unknownUser := client.Login(context.Background(), &forgenetv1.LoginRequest{Login: "personne", Password: "faux"})

	assert.Equal(t, codes.Unauthenticated, status.Code(wrongPassword))
	assert.Equal(t, status.Convert(wrongPassword).Message(), status.Convert(unknownUser).Message())
}

func TestAuthService_RepeatedFailuresAreThrottled(t *testing.T) {
	_, client := startAuthServer(t)
	ctx := context.Background()

	for range 10 {
		_, err := client.Login(ctx, &forgenetv1.LoginRequest{Login: "membre", Password: "faux"})
		require.Equal(t, codes.Unauthenticated, status.Code(err))
	}
	_, err := client.Login(ctx, &forgenetv1.LoginRequest{Login: "MEMBRE ", Password: testPassword})

	assert.Equal(t, codes.ResourceExhausted, status.Code(err), "même le bon mot de passe est refusé pendant la fenêtre")
	_, err = client.Login(ctx, &forgenetv1.LoginRequest{Login: "admin", Password: testPassword})
	assert.NoError(t, err, "les autres comptes restent ouverts")
}

func TestAuthService_LogoutRevokesTheSession(t *testing.T) {
	_, client := startAuthServer(t)
	res, err := client.Login(context.Background(), &forgenetv1.LoginRequest{Login: "admin", Password: testPassword})
	require.NoError(t, err)

	_, err = client.Logout(as(res.GetToken()), &forgenetv1.LogoutRequest{})
	require.NoError(t, err)

	_, err = client.GetCurrentUser(as(res.GetToken()), &forgenetv1.GetCurrentUserRequest{})
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestAuthService_OnlyLoginIsPublic(t *testing.T) {
	_, client := startAuthServer(t)

	_, err := client.GetCurrentUser(context.Background(), &forgenetv1.GetCurrentUserRequest{})
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
	_, err = client.Logout(context.Background(), &forgenetv1.LogoutRequest{})
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestAuthService_ChangePassword(t *testing.T) {
	ts, client := startAuthServer(t)
	ctx := context.Background()
	login := func(password string) (*forgenetv1.LoginResponse, error) {
		return client.Login(ctx, &forgenetv1.LoginRequest{Login: "membre", Password: password})
	}
	current, err := login(testPassword)
	require.NoError(t, err)
	other, err := login(testPassword)
	require.NoError(t, err)

	for name, req := range map[string]*forgenetv1.ChangePasswordRequest{
		"actuel incorrect": {CurrentPassword: "faux", NewPassword: "nouveau-mot-de-passe"},
		"trop court":       {CurrentPassword: testPassword, NewPassword: "court"},
		"identique":        {CurrentPassword: testPassword, NewPassword: testPassword},
	} {
		_, err := client.ChangePassword(as(current.GetToken()), req)
		assert.Equal(t, codes.InvalidArgument, status.Code(err), name)
		assert.NotEmpty(t, status.Convert(err).Message(), name)
	}

	res, err := client.ChangePassword(as(current.GetToken()), &forgenetv1.ChangePasswordRequest{
		CurrentPassword: testPassword, NewPassword: "nouveau-mot-de-passe",
	})
	require.NoError(t, err)
	// L'autre session, et le jeton API que startTestServer crée pour le membre.
	assert.Equal(t, int32(2), res.GetRevokedSessions())

	_, err = client.GetCurrentUser(as(current.GetToken()), &forgenetv1.GetCurrentUserRequest{})
	require.NoError(t, err, "la session du demandeur reste ouverte")
	_, err = client.GetCurrentUser(as(other.GetToken()), &forgenetv1.GetCurrentUserRequest{})
	assert.Equal(t, codes.Unauthenticated, status.Code(err), "les autres sessions sont fermées")
	_, err = client.GetCurrentUser(as(ts.memberToken), &forgenetv1.GetCurrentUserRequest{})
	assert.Equal(t, codes.Unauthenticated, status.Code(err), "les jetons API aussi")
	_, err = login("nouveau-mot-de-passe")
	require.NoError(t, err)
}

func TestAuthService_ChangePasswordNeedsASession(t *testing.T) {
	_, client := startAuthServer(t)
	_, err := client.ChangePassword(context.Background(), &forgenetv1.ChangePasswordRequest{
		CurrentPassword: testPassword, NewPassword: "nouveau-mot-de-passe",
	})
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestAuthService_ChangePasswordIsThrottled(t *testing.T) {
	_, client := startAuthServer(t)
	res, err := client.Login(context.Background(), &forgenetv1.LoginRequest{Login: "membre", Password: testPassword})
	require.NoError(t, err)
	ctx := as(res.GetToken())

	for range 10 {
		_, err := client.ChangePassword(ctx, &forgenetv1.ChangePasswordRequest{CurrentPassword: "faux", NewPassword: "nouveau-mot-de-passe"})
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	}
	_, err = client.ChangePassword(ctx, &forgenetv1.ChangePasswordRequest{CurrentPassword: testPassword, NewPassword: "nouveau-mot-de-passe"})
	assert.Equal(t, codes.ResourceExhausted, status.Code(err), "même l'actuel exact est refusé pendant la fenêtre")
}

// pendingAccount crée un compte sans mot de passe et son lien de première
// connexion, comme install.sh sans terminal.
func pendingAccount(t *testing.T, ts *testServer, name string) string {
	t.Helper()
	ctx := context.Background()
	svc := auth.NewService(ts.store, config.AuthConfig{JWTSecret: "secret-de-test-d-au-moins-32-caracteres"}, nil)
	u, err := svc.RegisterWithoutPassword(ctx, name, name+"@forgenet.test", store.RoleAdmin)
	require.NoError(t, err)
	token, _, err := svc.CreateAccessLink(ctx, u.ID, store.PurposeFirstLogin, "system:cli")
	require.NoError(t, err)
	return token
}

func TestAuthService_AccessLinkWithoutSession(t *testing.T) {
	ts, client := startAuthServer(t)
	ctx := context.Background()
	token := pendingAccount(t, ts, "premier")

	// Les deux appels sont publics : aucun jeton de session présenté.
	info, err := client.InspectAccessLink(ctx, &forgenetv1.InspectAccessLinkRequest{Token: token})
	require.NoError(t, err)
	assert.Equal(t, "premier", info.GetUsername())
	assert.Equal(t, forgenetv1.AccessLinkPurpose_ACCESS_LINK_PURPOSE_FIRST_LOGIN, info.GetPurpose())
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), info.GetExpiresAt().AsTime(), time.Minute)

	_, err = client.RedeemAccessLink(ctx, &forgenetv1.RedeemAccessLinkRequest{Token: token, NewPassword: "court"})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	res, err := client.RedeemAccessLink(ctx, &forgenetv1.RedeemAccessLinkRequest{Token: token, NewPassword: "mot-de-passe-choisi"})
	require.NoError(t, err, "un mot de passe refusé laisse le lien valable")
	assert.Equal(t, forgenetv1.UserRole_USER_ROLE_ADMIN, res.GetUser().GetRole())
	me, err := client.GetCurrentUser(as(res.GetToken()), &forgenetv1.GetCurrentUserRequest{})
	require.NoError(t, err, "la réponse ouvre une session")
	assert.Equal(t, "premier", me.GetUser().GetUsername())

	_, err = client.Login(ctx, &forgenetv1.LoginRequest{Login: "premier", Password: "mot-de-passe-choisi"})
	require.NoError(t, err)
}

func TestAuthService_AccessLinkRefusalsAreAllTheSame(t *testing.T) {
	ts, client := startAuthServer(t)
	ctx := context.Background()
	used := pendingAccount(t, ts, "premier")
	_, err := client.RedeemAccessLink(ctx, &forgenetv1.RedeemAccessLinkRequest{Token: used, NewPassword: "mot-de-passe-choisi"})
	require.NoError(t, err)

	var messages []string
	for _, token := range []string{used, "inconnu", ""} {
		_, err := client.InspectAccessLink(ctx, &forgenetv1.InspectAccessLinkRequest{Token: token})
		assert.Equal(t, codes.NotFound, status.Code(err), token)
		_, err = client.RedeemAccessLink(ctx, &forgenetv1.RedeemAccessLinkRequest{Token: token, NewPassword: "encore-un-autre"})
		assert.Equal(t, codes.NotFound, status.Code(err), token)
		messages = append(messages, status.Convert(err).Message())
	}
	assert.Equal(t, messages[0], messages[1], "un lien utilisé ou inconnu reçoit la même réponse")
}

func TestAuthService_PendingAccountCannotLogIn(t *testing.T) {
	ts, client := startAuthServer(t)
	pendingAccount(t, ts, "premier")
	_, err := client.Login(context.Background(), &forgenetv1.LoginRequest{Login: "premier", Password: "!"})
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

// startResetServer sert AuthService avec « mot de passe oublié » activé et
// une file d'envoi synchrone.
func startResetServer(t *testing.T, dashboardHost string) (*testServer, forgenetv1.AuthServiceClient, *transportGRPC.SyncOutboxForTests) {
	t.Helper()
	outbox := &transportGRPC.SyncOutboxForTests{}
	ts := startTestServer(t, &fakeEngine{}, func(srv *grpc.Server, st store.Store, auditor audit.Auditor) {
		authSvc := auth.NewService(st, config.AuthConfig{JWTSecret: "secret-de-test-d-au-moins-32-caracteres"}, auditor)
		api := transportGRPC.NewAuthServiceServer(authSvc, st.Users(), zerolog.Nop()).EnablePasswordReset(outbox, dashboardHost)
		forgenetv1.RegisterAuthServiceServer(srv, api)
	})
	return ts, forgenetv1.NewAuthServiceClient(ts.conn), outbox
}

func TestAuthService_LoginOptions(t *testing.T) {
	_, plain := startAuthServer(t)
	opts, err := plain.GetLoginOptions(context.Background(), &forgenetv1.GetLoginOptionsRequest{})
	require.NoError(t, err, "public")
	assert.False(t, opts.GetPasswordReset())
	_, err = plain.RequestPasswordReset(context.Background(), &forgenetv1.RequestPasswordResetRequest{Login: "admin"})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err), "sans relais")
	assert.Contains(t, status.Convert(err).Message(), "forgenet user set-password")

	_, client, _ := startResetServer(t, "")
	opts, err = client.GetLoginOptions(context.Background(), &forgenetv1.GetLoginOptionsRequest{})
	require.NoError(t, err)
	assert.True(t, opts.GetPasswordReset())
}

func TestAuthService_PasswordResetByMail(t *testing.T) {
	_, client, outbox := startResetServer(t, "forgenet.exemple.org")
	ctx := context.Background()

	// Même réponse, compte ou non.
	for _, login := range []string{"personne", "MEMBRE@forgenet.test "} {
		_, err := client.RequestPasswordReset(ctx, &forgenetv1.RequestPasswordResetRequest{Login: login})
		require.NoError(t, err, login)
	}
	msgs := outbox.Messages()
	require.Len(t, msgs, 1, "seul le compte existant reçoit un e-mail")
	assert.Equal(t, "membre@forgenet.test", msgs[0].To)
	link := regexp.MustCompile(`https://forgenet\.exemple\.org/acces#([A-Za-z0-9_-]{43})`).FindStringSubmatch(msgs[0].Body)
	require.Len(t, link, 2, msgs[0].Body)

	info, err := client.InspectAccessLink(ctx, &forgenetv1.InspectAccessLinkRequest{Token: link[1]})
	require.NoError(t, err)
	assert.Equal(t, forgenetv1.AccessLinkPurpose_ACCESS_LINK_PURPOSE_PASSWORD_RESET, info.GetPurpose())
	assert.WithinDuration(t, time.Now().Add(time.Hour), info.GetExpiresAt().AsTime(), time.Minute)

	_, err = client.RedeemAccessLink(ctx, &forgenetv1.RedeemAccessLinkRequest{Token: link[1], NewPassword: "mot-de-passe-retrouve"})
	require.NoError(t, err)
	_, err = client.Login(ctx, &forgenetv1.LoginRequest{Login: "membre", Password: testPassword})
	assert.Equal(t, codes.Unauthenticated, status.Code(err), "l'ancien mot de passe ne vaut plus")
	_, err = client.Login(ctx, &forgenetv1.LoginRequest{Login: "membre", Password: "mot-de-passe-retrouve"})
	require.NoError(t, err)
}

func TestAuthService_PasswordResetIsThrottledSilently(t *testing.T) {
	_, client, outbox := startResetServer(t, "")
	ctx := context.Background()
	for range 7 {
		_, err := client.RequestPasswordReset(ctx, &forgenetv1.RequestPasswordResetRequest{Login: "membre"})
		require.NoError(t, err, "la réponse ne change pas au-delà du frein")
	}
	msgs := outbox.Messages()
	assert.Len(t, msgs, 5, "5 e-mails par heure au plus pour un compte")
	assert.Contains(t, msgs[0].Body, "http://127.0.0.1:8080/acces#")
	assert.Contains(t, msgs[0].Body, "ssh -N", "sans domaine, le tunnel est expliqué")
}

func TestAuthService_PasswordResetRejectsAnEmptyLogin(t *testing.T) {
	_, client, outbox := startResetServer(t, "")
	_, err := client.RequestPasswordReset(context.Background(), &forgenetv1.RequestPasswordResetRequest{Login: "  "})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Empty(t, outbox.Messages())
}
