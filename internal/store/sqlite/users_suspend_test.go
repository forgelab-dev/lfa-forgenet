package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

// accounts crée des comptes et des sessions de test dans st.
type accounts struct {
	t   *testing.T
	ctx context.Context
	st  store.Store
}

func (a accounts) user(id string, role store.Role, hash string) {
	a.t.Helper()
	require.NoError(a.t, a.st.Users().CreateUser(a.ctx, &store.User{
		ID: id, Username: id, Email: id + "@forgenet.local", PasswordHash: hash, Role: role,
	}))
}

func (a accounts) session(id, userID string) {
	a.t.Helper()
	require.NoError(a.t, a.st.Sessions().CreateSession(a.ctx, &store.Session{
		ID: id, UserID: userID, TokenHash: "h_" + id, Name: id, Role: store.RoleMember,
		ExpiresAt: time.Now().Add(time.Hour),
	}))
}

func (a accounts) link(id, userID string) {
	a.t.Helper()
	require.NoError(a.t, a.st.ActionTokens().CreateActionToken(a.ctx, &store.ActionToken{
		ID: id, UserID: userID, Purpose: store.PurposePasswordReset, TokenHash: "h_" + id,
		ExpiresAt: time.Now().Add(time.Hour), CreatedBy: "system:test", CreatedAt: time.Now(),
	}))
}

func (a accounts) revoked(id string) bool {
	a.t.Helper()
	s, err := a.st.Sessions().GetSessionByID(a.ctx, id)
	require.NoError(a.t, err)
	return s.Revoked
}

func TestUserStore_SuspendAndResume(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	a := accounts{t, ctx, st}
	a.user("usr_admin", store.RoleAdmin, "h")
	a.user("usr_m", store.RoleMember, "h")
	a.session("sess_m", "usr_m")
	a.session("tok_m", "usr_m")
	a.session("sess_admin", "usr_admin")
	a.link("act_m", "usr_m")

	at := time.Date(2026, 11, 19, 10, 0, 0, 0, time.UTC)
	revoked, err := st.Users().SuspendUser(ctx, "usr_m", at)
	require.NoError(t, err)
	assert.Equal(t, int64(2), revoked)

	u, err := st.Users().GetUserByID(ctx, "usr_m")
	require.NoError(t, err)
	require.True(t, u.Suspended())
	assert.True(t, u.SuspendedAt.Equal(at))
	assert.True(t, a.revoked("sess_m"))
	assert.True(t, a.revoked("tok_m"))
	assert.False(t, a.revoked("sess_admin"), "les sessions des autres comptes restent ouvertes")

	_, err = st.ActionTokens().GetActiveActionToken(ctx, "h_act_m", time.Now())
	require.ErrorIs(t, err, store.ErrNotFound, "le lien en attente part avec la suspension")

	// Suspendre deux fois garde la première date.
	_, err = st.Users().SuspendUser(ctx, "usr_m", at.Add(time.Hour))
	require.NoError(t, err)
	u, err = st.Users().GetUserByID(ctx, "usr_m")
	require.NoError(t, err)
	assert.True(t, u.SuspendedAt.Equal(at))

	require.NoError(t, st.Users().ResumeUser(ctx, "usr_m"))
	u, err = st.Users().GetUserByID(ctx, "usr_m")
	require.NoError(t, err)
	assert.False(t, u.Suspended())
	assert.True(t, a.revoked("sess_m"), "réactiver ne rouvre pas les sessions fermées")

	_, err = st.Users().SuspendUser(ctx, "usr_absent", at)
	require.ErrorIs(t, err, store.ErrNotFound)
	require.ErrorIs(t, st.Users().ResumeUser(ctx, "usr_absent"), store.ErrNotFound)
}

func TestUserStore_LastActiveAdminIsKept(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	a := accounts{t, ctx, st}
	a.user("usr_admin", store.RoleAdmin, "h")
	// Un administrateur invité ou suspendu ne compte pas : il ne peut pas se
	// connecter.
	a.user("usr_invite", store.RoleAdmin, store.UnusablePasswordHash)
	a.user("usr_second", store.RoleAdmin, "h")

	_, err := st.Users().SuspendUser(ctx, "usr_second", time.Now())
	require.NoError(t, err)

	_, err = st.Users().SuspendUser(ctx, "usr_admin", time.Now())
	require.ErrorIs(t, err, store.ErrLastAdmin)
	require.ErrorIs(t, st.Users().DeleteUser(ctx, "usr_admin"), store.ErrLastAdmin)

	u, err := st.Users().GetUserByID(ctx, "usr_admin")
	require.NoError(t, err)
	assert.False(t, u.Suspended(), "le refus ne change rien")

	// Un administrateur déjà suspendu, ou invité, peut être retiré.
	require.NoError(t, st.Users().DeleteUser(ctx, "usr_second"))
	require.NoError(t, st.Users().DeleteUser(ctx, "usr_invite"))

	// Avec un second administrateur actif, le premier peut être suspendu.
	a.user("usr_releve", store.RoleAdmin, "h")
	_, err = st.Users().SuspendUser(ctx, "usr_admin", time.Now())
	require.NoError(t, err)
}

func TestUserStore_DeleteUserCascades(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	a := accounts{t, ctx, st}
	a.user("usr_admin", store.RoleAdmin, "h")
	a.user("usr_m", store.RoleMember, "h")
	a.session("sess_m", "usr_m")
	a.link("act_m", "usr_m")
	require.NoError(t, st.Audit().Record(ctx, &store.AuditEvent{
		ID: "evt_1", UserID: "usr_m", Action: "auth.login_success", Resource: "user:usr_m", Timestamp: time.Now(),
	}))

	require.NoError(t, st.Users().DeleteUser(ctx, "usr_m"))

	_, err := st.Users().GetUserByID(ctx, "usr_m")
	require.ErrorIs(t, err, store.ErrNotFound)
	_, err = st.Sessions().GetSessionByID(ctx, "sess_m")
	require.ErrorIs(t, err, store.ErrNotFound)
	_, err = st.ActionTokens().GetActiveActionToken(ctx, "h_act_m", time.Now())
	require.ErrorIs(t, err, store.ErrNotFound)

	events, err := st.Audit().List(ctx, 10, 0)
	require.NoError(t, err)
	require.Len(t, events, 1, "le journal d'audit garde les entrées du compte supprimé")

	require.ErrorIs(t, st.Users().DeleteUser(ctx, "usr_m"), store.ErrNotFound)
}

func TestUserStore_RevokeUserSessions(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	a := accounts{t, ctx, st}
	a.user("usr_a", store.RoleAdmin, "h")
	a.session("sess_ici", "usr_a")
	a.session("sess_perdue", "usr_a")
	a.session("tok_ci", "usr_a")

	revoked, err := st.Users().RevokeUserSessions(ctx, "usr_a", "sess_ici")
	require.NoError(t, err)
	assert.Equal(t, int64(2), revoked)
	assert.False(t, a.revoked("sess_ici"))
	assert.True(t, a.revoked("sess_perdue"))
	assert.True(t, a.revoked("tok_ci"))

	revoked, err = st.Users().RevokeUserSessions(ctx, "usr_a", "")
	require.NoError(t, err)
	assert.Equal(t, int64(1), revoked)

	u, err := st.Users().GetUserByID(ctx, "usr_a")
	require.NoError(t, err)
	assert.False(t, u.Suspended(), "fermer les sessions ne suspend pas le compte")

	_, err = st.Users().RevokeUserSessions(ctx, "usr_absent", "")
	require.ErrorIs(t, err, store.ErrNotFound)
}
