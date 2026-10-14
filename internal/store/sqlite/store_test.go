package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
)

func newTestStore(t *testing.T) store.Store {
	t.Helper()
	ctx := context.Background()
	cfg := config.DatabaseConfig{
		Path:         ":memory:",
		BusyTimeout:  2 * time.Second,
		MaxOpenConns: 1,
		MaxIdleConns: 1,
	}

	st, err := sqlite.Open(ctx, cfg)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = st.Close()
	})
	return st
}

func TestStore_PingAndCheck(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	require.NoError(t, st.Ping(ctx))

	checker, ok := st.(interface{ Check(context.Context) error })
	require.True(t, ok)
	require.NoError(t, checker.Check(ctx))
}

func TestUserStore_CRUD(t *testing.T) {
	st := newTestStore(t)
	users := st.Users()
	ctx := context.Background()

	// 1. Initialement vide
	count, err := users.CountUsers(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// 2. Création utilisateur
	adminUser := &store.User{
		ID:           "usr_admin_1",
		Username:     "admin",
		Email:        "admin@forgenet.local",
		PasswordHash: "hashed_secret_admin",
		Role:         store.RoleAdmin,
	}
	err = users.CreateUser(ctx, adminUser)
	require.NoError(t, err)

	count, err = users.CountUsers(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// 3. Récupération par ID, Username, Email
	byID, err := users.GetUserByID(ctx, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, adminUser.Username, byID.Username)
	assert.Equal(t, adminUser.Email, byID.Email)
	assert.Equal(t, store.RoleAdmin, byID.Role)

	byUsername, err := users.GetUserByUsername(ctx, "admin")
	require.NoError(t, err)
	assert.Equal(t, adminUser.ID, byUsername.ID)

	byEmail, err := users.GetUserByEmail(ctx, "admin@forgenet.local")
	require.NoError(t, err)
	assert.Equal(t, adminUser.ID, byEmail.ID)

	// 4. Conflit d'unicité (même username ou email)
	duplicateUser := &store.User{
		ID:           "usr_admin_2",
		Username:     "admin",
		Email:        "other@forgenet.local",
		PasswordHash: "secret",
		Role:         store.RoleMember,
	}
	err = users.CreateUser(ctx, duplicateUser)
	require.ErrorIs(t, err, store.ErrConflict)

	// 5. Mise à jour
	byID.Role = store.RoleMember
	byID.Email = "admin_updated@forgenet.local"
	err = users.UpdateUser(ctx, byID)
	require.NoError(t, err)

	updated, err := users.GetUserByID(ctx, adminUser.ID)
	require.NoError(t, err)
	assert.Equal(t, store.RoleMember, updated.Role)
	assert.Equal(t, "admin_updated@forgenet.local", updated.Email)

	// 6. Liste
	list, err := users.ListUsers(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// 7. Suppression
	err = users.DeleteUser(ctx, adminUser.ID)
	require.NoError(t, err)

	_, err = users.GetUserByID(ctx, adminUser.ID)
	require.ErrorIs(t, err, store.ErrNotFound)
}

func TestSessionStore_Lifecycle(t *testing.T) {
	st := newTestStore(t)
	users := st.Users()
	sessions := st.Sessions()
	ctx := context.Background()

	// Création préalable d'un utilisateur pour satisfaire la clé étrangère
	u := &store.User{
		ID:           "usr_test_sess",
		Username:     "testuser",
		Email:        "test@forgenet.local",
		PasswordHash: "hash",
		Role:         store.RoleAdmin,
	}
	require.NoError(t, users.CreateUser(ctx, u))

	// 1. Création d'une session
	expiresAt := time.Now().UTC().Add(24 * time.Hour)
	sess := &store.Session{
		ID:        "sess_1",
		UserID:    u.ID,
		TokenHash: "hash_token_abc123",
		Name:      "CLI Token",
		Role:      store.RoleAdmin,
		ExpiresAt: expiresAt,
	}
	require.NoError(t, sessions.CreateSession(ctx, sess))

	// 2. Recherche par token hash
	found, err := sessions.GetSessionByTokenHash(ctx, "hash_token_abc123")
	require.NoError(t, err)
	assert.Equal(t, sess.ID, found.ID)
	assert.Equal(t, u.ID, found.UserID)
	assert.True(t, found.IsActive())
	assert.False(t, found.IsExpired())

	// 3. Touch session (mise à jour last_used_at)
	now := time.Now().UTC()
	err = sessions.TouchSession(ctx, sess.ID, now)
	require.NoError(t, err)

	refreshed, err := sessions.GetSessionByID(ctx, sess.ID)
	require.NoError(t, err)
	require.NotNil(t, refreshed.LastUsedAt)

	// 4. Révocation
	err = sessions.RevokeSession(ctx, sess.ID)
	require.NoError(t, err)

	revoked, err := sessions.GetSessionByID(ctx, sess.ID)
	require.NoError(t, err)
	assert.True(t, revoked.Revoked)
	assert.False(t, revoked.IsActive())

	// 5. Nettoyage des sessions expirées
	expiredSess := &store.Session{
		ID:        "sess_expired",
		UserID:    u.ID,
		TokenHash: "hash_expired",
		Name:      "Old Session",
		Role:      store.RoleMember,
		ExpiresAt: time.Now().UTC().Add(-1 * time.Hour), // Déjà expirée
	}
	require.NoError(t, sessions.CreateSession(ctx, expiredSess))

	err = sessions.DeleteExpiredSessions(ctx)
	require.NoError(t, err)

	_, err = sessions.GetSessionByID(ctx, "sess_expired")
	require.ErrorIs(t, err, store.ErrNotFound)
}

func TestAuditStore_RecordAndList(t *testing.T) {
	st := newTestStore(t)
	audit := st.Audit()
	ctx := context.Background()

	ev1 := &store.AuditEvent{
		ID:        "aud_1",
		UserID:    "usr_admin",
		Action:    "container.start",
		Resource:  "container:postgres-dev",
		IPAddress: "192.168.1.50",
		Metadata:  `{"port": 5432}`,
		Timestamp: time.Now().UTC().Add(-2 * time.Minute),
	}
	ev2 := &store.AuditEvent{
		ID:        "aud_2",
		UserID:    "usr_admin",
		Action:    "container.stop",
		Resource:  "container:postgres-dev",
		IPAddress: "192.168.1.50",
		Metadata:  `{"force": false}`,
		Timestamp: time.Now().UTC(),
	}

	require.NoError(t, audit.Record(ctx, ev1))
	require.NoError(t, audit.Record(ctx, ev2))

	events, err := audit.List(ctx, 10, 0)
	require.NoError(t, err)
	require.Len(t, events, 2)

	// Ordre antéchronologique (le plus récent en premier)
	assert.Equal(t, "aud_2", events[0].ID)
	assert.Equal(t, "aud_1", events[1].ID)
}
