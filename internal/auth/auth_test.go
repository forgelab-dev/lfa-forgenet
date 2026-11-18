package auth_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
)

const testPassword = "motdepasse-de-test"

func setupAuthTest(t *testing.T) (auth.Service, store.Store) {
	t.Helper()
	st, err := sqlite.Open(context.Background(), config.DatabaseConfig{Path: ":memory:", BusyTimeout: 2 * time.Second})
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })

	svc := auth.NewService(st, config.AuthConfig{
		JWTSecret:   "secret-de-test-d-au-moins-32-caracteres",
		TokenTTL:    time.Hour,
		APITokenTTL: 24 * time.Hour,
	}, audit.NewAuditor(st.Audit(), zerolog.Nop()))
	return svc, st
}

func register(t *testing.T, svc auth.Service, name string, role store.Role) *store.User {
	t.Helper()
	u, err := svc.Register(context.Background(), name, name+"@forgenet.test", testPassword, role)
	require.NoError(t, err)
	return u
}

func TestPassword_HashingAndValidation(t *testing.T) {
	_, err := auth.HashPassword("court")
	require.ErrorIs(t, err, auth.ErrPasswordTooShort)
	// Huit caractères accentués font seize octets mais comptent pour huit.
	_, err = auth.HashPassword("éééééééé")
	require.NoError(t, err)
	_, err = auth.HashPassword(strings.Repeat("a", 73))
	require.ErrorIs(t, err, auth.ErrPasswordTooLong)

	hashed, err := auth.HashPassword(testPassword)
	require.NoError(t, err)
	assert.NoError(t, auth.CheckPassword(hashed, testPassword))
	assert.ErrorIs(t, auth.CheckPassword(hashed, "mauvais-mot-de-passe"), auth.ErrInvalidCredentials)
}

func TestRegister_NeverGrantsAdminImplicitly(t *testing.T) {
	svc, _ := setupAuthTest(t)
	ctx := context.Background()

	first := register(t, svc, "premier", store.RoleMember)
	assert.Equal(t, store.RoleMember, first.Role, "le premier compte garde le rôle demandé")

	_, err := svc.Register(ctx, "sansrole", "sansrole@forgenet.test", testPassword, "")
	require.ErrorIs(t, err, auth.ErrInvalidRole)

	_, err = svc.Register(ctx, "premier", "autre@forgenet.test", testPassword, store.RoleMember)
	require.ErrorIs(t, err, store.ErrConflict)
}

func TestLogin(t *testing.T) {
	svc, _ := setupAuthTest(t)
	ctx := context.Background()
	alice := register(t, svc, "alice", store.RoleAdmin)

	for name, id := range map[string]string{"mauvais mot de passe": "alice", "compte inconnu": "inconnu"} {
		t.Run(name, func(t *testing.T) {
			_, err := svc.Login(ctx, id, "mauvais-mot-de-passe")
			require.ErrorIs(t, err, auth.ErrInvalidCredentials)
		})
	}

	for _, id := range []string{"alice", "ALICE@forgenet.test"} {
		res, err := svc.Login(ctx, id, testPassword)
		require.NoError(t, err, id)
		claims, err := svc.ValidateToken(ctx, res.Token)
		require.NoError(t, err)
		assert.Equal(t, alice.ID, claims.UserID)
		assert.Equal(t, store.RoleAdmin, claims.Role)
	}
}

func TestCreateAPIToken_RoleCappedByUser(t *testing.T) {
	svc, _ := setupAuthTest(t)
	ctx := context.Background()
	member := register(t, svc, "membre", store.RoleMember)
	admin := register(t, svc, "admin", store.RoleAdmin)

	_, _, err := svc.CreateAPIToken(ctx, member.ID, "élévation", store.RoleAdmin, time.Hour)
	require.ErrorIs(t, err, auth.ErrForbidden, "un membre ne peut pas obtenir un jeton admin")

	token, session, err := svc.CreateAPIToken(ctx, member.ID, "défaut", "", time.Hour)
	require.NoError(t, err)
	assert.Equal(t, store.RoleMember, session.Role, "sans rôle demandé, le jeton prend celui de l'utilisateur")
	assert.True(t, auth.IsAPIToken(token))

	_, session, err = svc.CreateAPIToken(ctx, admin.ID, "réduit", store.RoleMember, time.Hour)
	require.NoError(t, err)
	assert.Equal(t, store.RoleMember, session.Role, "un admin peut créer un jeton aux droits réduits")
}

func TestValidateToken_RevocationAndDemotion(t *testing.T) {
	svc, st := setupAuthTest(t)
	ctx := context.Background()
	admin := register(t, svc, "admin", store.RoleAdmin)

	login, err := svc.Login(ctx, "admin", testPassword)
	require.NoError(t, err)
	apiToken, apiSession, err := svc.CreateAPIToken(ctx, admin.ID, "ci", store.RoleAdmin, time.Hour)
	require.NoError(t, err)

	t.Run("rétrogradation immédiate", func(t *testing.T) {
		admin.Role = store.RoleMember
		require.NoError(t, st.Users().UpdateUser(ctx, admin))
		for _, token := range []string{login.Token, apiToken} {
			claims, err := svc.ValidateToken(ctx, token)
			require.NoError(t, err)
			assert.Equal(t, store.RoleMember, claims.Role)
		}
	})

	t.Run("jeton API révoqué", func(t *testing.T) {
		require.NoError(t, svc.RevokeSession(ctx, apiSession.ID))
		_, err := svc.ValidateToken(ctx, apiToken)
		require.ErrorIs(t, err, auth.ErrInvalidToken)
	})

	t.Run("JWT révoqué", func(t *testing.T) {
		claims, err := svc.ValidateToken(ctx, login.Token)
		require.NoError(t, err)
		require.NoError(t, svc.RevokeSession(ctx, claims.SessionID))
		_, err = svc.ValidateToken(ctx, login.Token)
		require.ErrorIs(t, err, auth.ErrInvalidToken)
	})

	t.Run("jeton inconnu", func(t *testing.T) {
		_, err := svc.ValidateToken(ctx, "fgn_inconnu")
		require.ErrorIs(t, err, auth.ErrInvalidToken)
	})
}

// Si la base ne répond plus, un jeton jusque-là valide est refusé : la
// vérification de révocation ne doit jamais être sautée.
func TestValidateToken_FailsClosedWhenStoreUnavailable(t *testing.T) {
	svc, st := setupAuthTest(t)
	ctx := context.Background()
	register(t, svc, "admin", store.RoleAdmin)
	login, err := svc.Login(ctx, "admin", testPassword)
	require.NoError(t, err)

	require.NoError(t, st.Close())
	claims, err := svc.ValidateToken(ctx, login.Token)
	require.Error(t, err)
	assert.Nil(t, claims)
}

func TestGuard(t *testing.T) {
	svc, st := setupAuthTest(t)
	ctx := context.Background()
	member := register(t, svc, "membre", store.RoleMember)
	admin := register(t, svc, "admin", store.RoleAdmin)
	memberToken, _, err := svc.CreateAPIToken(ctx, member.ID, "membre", "", time.Hour)
	require.NoError(t, err)
	adminToken, _, err := svc.CreateAPIToken(ctx, admin.ID, "admin", "", time.Hour)
	require.NoError(t, err)

	const (
		publicMethod = "/grpc.health.v1.Health/Check"
		memberMethod = "/forgenet.v1.Exemple/Lire"
		adminMethod  = "/forgenet.v1.Exemple/Supprimer"
	)
	policy := auth.DefaultPolicy()
	policy.MinRole[memberMethod] = store.RoleMember
	unary := auth.NewGuard(svc, policy, zerolog.Nop()).Unary()

	call := func(method, token string) (any, error) {
		callCtx := ctx
		if token != "" {
			callCtx = metadata.NewIncomingContext(ctx, metadata.Pairs("authorization", "Bearer "+token))
		}
		return unary(callCtx, nil, &grpc.UnaryServerInfo{FullMethod: method}, func(ctx context.Context, _ any) (any, error) {
			if claims, ok := auth.UserFromContext(ctx); ok {
				return claims.Username, nil
			}
			return "anonyme", nil
		})
	}
	code := func(err error) codes.Code { return status.Code(err) }

	cases := []struct {
		name, method, token string
		want                codes.Code
	}{
		{"sonde publique sans jeton", publicMethod, "", codes.OK},
		{"méthode protégée sans jeton", memberMethod, "", codes.Unauthenticated},
		{"jeton invalide", memberMethod, "fgn_faux", codes.Unauthenticated},
		{"membre sur méthode membre", memberMethod, memberToken, codes.OK},
		{"membre sur méthode non déclarée (admin)", adminMethod, memberToken, codes.PermissionDenied},
		{"admin sur méthode non déclarée", adminMethod, adminToken, codes.OK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := call(tc.method, tc.token)
			assert.Equal(t, tc.want, code(err))
		})
	}

	t.Run("identité transmise au handler", func(t *testing.T) {
		resp, err := call(adminMethod, adminToken)
		require.NoError(t, err)
		assert.Equal(t, "admin", resp)
	})

	t.Run("base indisponible : erreur générique, sans détail", func(t *testing.T) {
		require.NoError(t, st.Close())
		_, err := call(memberMethod, memberToken)
		assert.Equal(t, codes.Unavailable, code(err))
		assert.NotContains(t, status.Convert(err).Message(), "sql")
	})
}

func TestGuard_Stream(t *testing.T) {
	svc, _ := setupAuthTest(t)
	stream := auth.NewGuard(svc, auth.DefaultPolicy(), zerolog.Nop()).Stream()

	err := stream(nil, &fakeStream{ctx: context.Background()}, &grpc.StreamServerInfo{FullMethod: "/forgenet.v1.Exemple/Suivre"},
		func(any, grpc.ServerStream) error { return nil })
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

type fakeStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (f *fakeStream) Context() context.Context { return f.ctx }

func TestBearerToken(t *testing.T) {
	assert.Equal(t, "abc", auth.BearerToken("Bearer abc"))
	assert.Equal(t, "abc", auth.BearerToken("  bearer   abc "))
	assert.Equal(t, "Basic abc", auth.BearerToken("Basic abc"), "autre schéma : refusé ensuite comme jeton invalide")
	assert.Empty(t, auth.BearerToken(""))
}

func TestChangePassword(t *testing.T) {
	svc, st := setupAuthTest(t)
	ctx := context.Background()
	u := register(t, svc, "alice", store.RoleMember)

	current, err := svc.Login(ctx, "alice", testPassword)
	require.NoError(t, err)
	other, err := svc.Login(ctx, "alice", testPassword)
	require.NoError(t, err)
	apiToken, _, err := svc.CreateAPIToken(ctx, u.ID, "cli", store.RoleMember, 0)
	require.NoError(t, err)
	claims, err := svc.ValidateToken(ctx, current.Token)
	require.NoError(t, err)

	_, err = svc.ChangePassword(ctx, u.ID, claims.SessionID, "mauvais-mot-de-passe", "nouveau-mot-de-passe")
	require.ErrorIs(t, err, auth.ErrWrongCurrentPassword)
	_, err = svc.ChangePassword(ctx, u.ID, claims.SessionID, testPassword, testPassword)
	require.ErrorIs(t, err, auth.ErrSamePassword)
	_, err = svc.ChangePassword(ctx, u.ID, claims.SessionID, testPassword, "court")
	require.ErrorIs(t, err, auth.ErrPasswordTooShort)

	revoked, err := svc.ChangePassword(ctx, u.ID, claims.SessionID, testPassword, "nouveau-mot-de-passe")
	require.NoError(t, err)
	assert.Equal(t, int64(2), revoked, "l'autre session et le jeton API")

	_, err = svc.ValidateToken(ctx, current.Token)
	require.NoError(t, err, "la session qui a demandé le changement reste ouverte")
	_, err = svc.ValidateToken(ctx, other.Token)
	require.ErrorIs(t, err, auth.ErrInvalidToken)
	_, err = svc.ValidateToken(ctx, apiToken)
	require.ErrorIs(t, err, auth.ErrInvalidToken)

	_, err = svc.Login(ctx, "alice", testPassword)
	require.ErrorIs(t, err, auth.ErrInvalidCredentials)
	_, err = svc.Login(ctx, "alice", "nouveau-mot-de-passe")
	require.NoError(t, err)

	actions := auditActions(t, st)
	assert.Contains(t, actions, "user.password_change_failed")
	assert.Contains(t, actions, "user.password_change")
}

func TestResetPassword(t *testing.T) {
	svc, st := setupAuthTest(t)
	ctx := context.Background()
	register(t, svc, "admin", store.RoleAdmin)
	session, err := svc.Login(ctx, "admin", testPassword)
	require.NoError(t, err)

	revoked, err := svc.ResetPassword(ctx, " admin ", "mot-de-passe-retrouve")
	require.NoError(t, err)
	assert.Equal(t, int64(1), revoked)
	_, err = svc.ValidateToken(ctx, session.Token)
	require.ErrorIs(t, err, auth.ErrInvalidToken)
	_, err = svc.Login(ctx, "admin", "mot-de-passe-retrouve")
	require.NoError(t, err)

	_, err = svc.ResetPassword(ctx, "inconnu", "mot-de-passe-retrouve")
	require.ErrorIs(t, err, store.ErrNotFound)
	_, err = svc.ResetPassword(ctx, "admin", "court")
	require.ErrorIs(t, err, auth.ErrPasswordTooShort)

	assert.Contains(t, auditActions(t, st), "user.password_reset")
}

func auditActions(t *testing.T, st store.Store) []string {
	t.Helper()
	events, err := st.Audit().List(context.Background(), 100, 0)
	require.NoError(t, err)
	actions := make([]string, 0, len(events))
	for _, e := range events {
		actions = append(actions, e.Action)
	}
	return actions
}
