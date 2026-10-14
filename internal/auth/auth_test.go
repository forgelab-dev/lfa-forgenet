package auth_test

import (
	"context"
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

func setupAuthTest(t *testing.T) (auth.Service, store.Store) {
	t.Helper()
	ctx := context.Background()
	dbCfg := config.DatabaseConfig{
		Path:         ":memory:",
		BusyTimeout:  2 * time.Second,
		MaxOpenConns: 1,
		MaxIdleConns: 1,
	}

	st, err := sqlite.Open(ctx, dbCfg)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = st.Close()
	})

	aud := audit.NewAuditor(st.Audit(), zerolog.Nop())
	authCfg := config.AuthConfig{
		JWTSecret:   "test-secret-key-at-least-32-characters-long",
		TokenTTL:    1 * time.Hour,
		APITokenTTL: 24 * time.Hour,
	}

	svc := auth.NewService(st, authCfg, aud)
	return svc, st
}

func TestPassword_HashingAndValidation(t *testing.T) {
	// Trop court
	_, err := auth.HashPassword("short")
	require.ErrorIs(t, err, auth.ErrPasswordTooShort)

	// Valide
	hashed, err := auth.HashPassword("supersecret123")
	require.NoError(t, err)
	assert.NotEmpty(t, hashed)

	// Vérification valide
	assert.NoError(t, auth.CheckPassword(hashed, "supersecret123"))

	// Vérification invalide
	assert.ErrorIs(t, auth.CheckPassword(hashed, "wrongpassword"), auth.ErrInvalidCredentials)
}

func TestAuthService_RegisterAndBootstrap(t *testing.T) {
	svc, _ := setupAuthTest(t)
	ctx := context.Background()

	// 1. Premier utilisateur : doit automatiquement être Admin
	u1, err := svc.Register(ctx, "firstadmin", "admin@forgenet.local", "secretpass123", store.RoleMember)
	require.NoError(t, err)
	assert.Equal(t, store.RoleAdmin, u1.Role, "Le premier utilisateur doit devenir admin")

	// 2. Deuxième utilisateur : respecte le rôle demandé (member)
	u2, err := svc.Register(ctx, "seconduser", "user@forgenet.local", "secretpass123", store.RoleMember)
	require.NoError(t, err)
	assert.Equal(t, store.RoleMember, u2.Role)

	// 3. Conflit nom d'utilisateur existant
	_, err = svc.Register(ctx, "firstadmin", "other@forgenet.local", "secretpass123", store.RoleMember)
	require.ErrorIs(t, err, store.ErrConflict)
}

func TestAuthService_LoginAndTokens(t *testing.T) {
	svc, _ := setupAuthTest(t)
	ctx := context.Background()

	_, err := svc.Register(ctx, "alice", "alice@forgenet.local", "password123", store.RoleAdmin)
	require.NoError(t, err)

	// Échec mot de passe
	_, err = svc.Login(ctx, "alice", "badpass")
	require.ErrorIs(t, err, auth.ErrInvalidCredentials)

	// Échec utilisateur inconnu
	_, err = svc.Login(ctx, "unknown", "password123")
	require.ErrorIs(t, err, auth.ErrInvalidCredentials)

	// Connexion réussie par username
	res, err := svc.Login(ctx, "alice", "password123")
	require.NoError(t, err)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, "alice", res.User.Username)

	// Connexion réussie par email
	resEmail, err := svc.Login(ctx, "alice@forgenet.local", "password123")
	require.NoError(t, err)
	assert.NotEmpty(t, resEmail.Token)

	// Validation du JWT
	claims, err := svc.ValidateToken(ctx, res.Token)
	require.NoError(t, err)
	assert.Equal(t, res.User.ID, claims.UserID)
	assert.Equal(t, store.RoleAdmin, claims.Role)

	// Création d'un Token API
	plainToken, session, err := svc.CreateAPIToken(ctx, res.User.ID, "CI Runner", store.RoleMember, 10*time.Hour)
	require.NoError(t, err)
	assert.True(t, auth.IsAPIToken(plainToken))

	// Validation du Token API
	apiClaims, err := svc.ValidateToken(ctx, plainToken)
	require.NoError(t, err)
	assert.Equal(t, res.User.ID, apiClaims.UserID)
	assert.Equal(t, store.RoleMember, apiClaims.Role)

	// Révocation de session
	err = svc.RevokeSession(ctx, session.ID)
	require.NoError(t, err)

	_, err = svc.ValidateToken(ctx, plainToken)
	require.ErrorIs(t, err, auth.ErrExpiredToken)
}

func TestRequireRole(t *testing.T) {
	ctx := context.Background()

	// Sans utilisateur dans le contexte
	assert.ErrorIs(t, auth.RequireRole(ctx, store.RoleAdmin), auth.ErrUnauthenticated)

	// Utilisateur Membre
	memberCtx := auth.WithUserContext(ctx, &auth.Claims{
		UserID: "u1",
		Role:   store.RoleMember,
	})
	assert.NoError(t, auth.RequireRole(memberCtx, store.RoleMember, store.RoleAdmin))
	assert.ErrorIs(t, auth.RequireRole(memberCtx, store.RoleAdmin), auth.ErrForbidden)

	// Utilisateur Admin
	adminCtx := auth.WithUserContext(ctx, &auth.Claims{
		UserID: "u2",
		Role:   store.RoleAdmin,
	})
	assert.NoError(t, auth.RequireRole(adminCtx, store.RoleAdmin))
	assert.NoError(t, auth.RequireRole(adminCtx, store.RoleMember, store.RoleAdmin))
}

func TestUnaryServerInterceptor(t *testing.T) {
	svc, _ := setupAuthTest(t)
	ctx := context.Background()

	admin, err := svc.Register(ctx, "bob", "bob@forgenet.local", "securepass123", store.RoleAdmin)
	require.NoError(t, err)

	authRes, err := svc.Login(ctx, "bob", "securepass123")
	require.NoError(t, err)

	exemptMethods := map[string]bool{
		"/grpc.health.v1.Health/Check": true,
	}

	interceptor := auth.UnaryServerInterceptor(svc, exemptMethods)

	dummyHandler := func(ctx context.Context, _ any) (any, error) {
		claims, ok := auth.UserFromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "non authentifié dans handler")
		}
		return claims.Username, nil
	}

	// 1. Appel méthode exemptée (pas de token)
	healthInfo := &grpc.UnaryServerInfo{FullMethod: "/grpc.health.v1.Health/Check"}
	exemptHandler := func(_ context.Context, _ any) (any, error) {
		return "health_ok", nil
	}
	resp, err := interceptor(ctx, nil, healthInfo, exemptHandler)
	require.NoError(t, err)
	assert.Equal(t, "health_ok", resp)

	// 2. Appel méthode protégée SANS token -> rejeté
	protectedInfo := &grpc.UnaryServerInfo{FullMethod: "/forgenet.v1.ContainerService/ListContainers"}
	_, err = interceptor(ctx, nil, protectedInfo, dummyHandler)
	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Unauthenticated, st.Code())

	// 3. Appel méthode protégée avec Bearer token JWT -> succès
	md := metadata.Pairs("authorization", "Bearer "+authRes.Token)
	ctxWithAuth := metadata.NewIncomingContext(ctx, md)

	resp, err = interceptor(ctxWithAuth, nil, protectedInfo, dummyHandler)
	require.NoError(t, err)
	assert.Equal(t, "bob", resp)

	// 4. Appel méthode protégée avec Token API -> succès
	plainAPIToken, _, err := svc.CreateAPIToken(ctx, admin.ID, "Test Agent", store.RoleAdmin, time.Hour)
	require.NoError(t, err)

	mdAPI := metadata.Pairs("authorization", "Bearer "+plainAPIToken)
	ctxWithAPI := metadata.NewIncomingContext(ctx, mdAPI)

	resp, err = interceptor(ctxWithAPI, nil, protectedInfo, dummyHandler)
	require.NoError(t, err)
	assert.Equal(t, "bob", resp)
}
