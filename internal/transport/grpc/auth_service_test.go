package grpc_test

import (
	"context"
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
