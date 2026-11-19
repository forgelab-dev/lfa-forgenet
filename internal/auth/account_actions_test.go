package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

func TestSuspendAccount(t *testing.T) {
	svc, st := setupAuthTest(t)
	ctx := context.Background()
	admin := register(t, svc, "admin", store.RoleAdmin)
	bob := register(t, svc, "bob", store.RoleMember)

	session, err := svc.Login(ctx, "bob", testPassword)
	require.NoError(t, err)
	apiToken, _, err := svc.CreateAPIToken(ctx, bob.ID, "ci", store.RoleMember, time.Hour)
	require.NoError(t, err)
	_, _, err = svc.CreateAccessLink(ctx, bob.ID, store.PurposePasswordReset, "system:test")
	require.NoError(t, err)

	user, revoked, err := svc.SuspendAccount(ctx, admin.ID, bob.ID)
	require.NoError(t, err)
	assert.True(t, user.Suspended())
	assert.Equal(t, int64(2), revoked, "la session et le jeton d'API")

	// Le garde refuse au prochain appel, quel que soit le jeton.
	for name, token := range map[string]string{"session": session.Token, "jeton d'API": apiToken} {
		_, err = svc.ValidateToken(ctx, token)
		require.ErrorIs(t, err, auth.ErrInvalidToken, name)
	}

	// Même refus qu'un mauvais mot de passe, même avec le bon.
	_, err = svc.Login(ctx, "bob", testPassword)
	require.ErrorIs(t, err, auth.ErrInvalidCredentials)

	// Aucun nouveau lien pour un compte suspendu.
	_, _, err = svc.CreateAccessLink(ctx, bob.ID, store.PurposePasswordReset, "system:test")
	require.ErrorIs(t, err, auth.ErrAccountSuspended)

	_, err = svc.ResumeAccount(ctx, admin.ID, bob.ID)
	require.NoError(t, err)
	_, err = svc.Login(ctx, "bob", testPassword)
	require.NoError(t, err, "réactivé, le compte se connecte de nouveau")
	_, err = svc.ValidateToken(ctx, session.Token)
	require.ErrorIs(t, err, auth.ErrInvalidToken, "les sessions fermées le restent")

	actions := auditActions(t, st)
	assert.Contains(t, actions, "user.suspend")
	assert.Contains(t, actions, "user.resume")
}

func TestSuspendAccount_LinkStopsWorking(t *testing.T) {
	svc, _ := setupAuthTest(t)
	ctx := context.Background()
	admin := register(t, svc, "admin", store.RoleAdmin)
	invited, err := svc.RegisterWithoutPassword(ctx, "invite", "invite@forgenet.test", store.RoleMember)
	require.NoError(t, err)
	token, _, err := svc.CreateAccessLink(ctx, invited.ID, store.PurposeInvitation, "user:"+admin.ID)
	require.NoError(t, err)

	_, _, err = svc.SuspendAccount(ctx, admin.ID, invited.ID)
	require.NoError(t, err)
	_, _, err = svc.InspectAccessLink(ctx, token)
	require.ErrorIs(t, err, auth.ErrLinkInvalid)
	_, _, err = svc.RedeemAccessLink(ctx, token, "mot-de-passe-choisi")
	require.ErrorIs(t, err, auth.ErrLinkInvalid)
}

func TestAccountActions_Guards(t *testing.T) {
	svc, _ := setupAuthTest(t)
	ctx := context.Background()
	admin := register(t, svc, "admin", store.RoleAdmin)

	_, _, err := svc.SuspendAccount(ctx, admin.ID, admin.ID)
	require.ErrorIs(t, err, auth.ErrSelfAction)
	_, err = svc.DeleteAccount(ctx, admin.ID, admin.ID)
	require.ErrorIs(t, err, auth.ErrSelfAction)

	// Sur la machine du serveur (sans auteur), seul le dernier administrateur
	// actif est protégé.
	_, _, err = svc.SuspendAccount(ctx, "", admin.ID)
	require.ErrorIs(t, err, store.ErrLastAdmin)
	_, err = svc.DeleteAccount(ctx, "", admin.ID)
	require.ErrorIs(t, err, store.ErrLastAdmin)

	_, _, err = svc.SuspendAccount(ctx, admin.ID, "usr_absent")
	require.ErrorIs(t, err, store.ErrNotFound)
	_, err = svc.DeleteAccount(ctx, admin.ID, "usr_absent")
	require.ErrorIs(t, err, store.ErrNotFound)
}

func TestCloseSessions(t *testing.T) {
	svc, st := setupAuthTest(t)
	ctx := context.Background()
	admin := register(t, svc, "admin", store.RoleAdmin)

	here, err := svc.Login(ctx, "admin", testPassword)
	require.NoError(t, err)
	lost, err := svc.Login(ctx, "admin", testPassword)
	require.NoError(t, err)
	hereClaims, err := svc.ValidateToken(ctx, here.Token)
	require.NoError(t, err)

	// Sur son propre compte, la session de qui agit reste ouverte.
	revoked, err := svc.CloseSessions(ctx, admin.ID, admin.ID, hereClaims.SessionID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), revoked)
	_, err = svc.ValidateToken(ctx, here.Token)
	require.NoError(t, err)
	_, err = svc.ValidateToken(ctx, lost.Token)
	require.ErrorIs(t, err, auth.ErrInvalidToken)

	_, err = svc.Login(ctx, "admin", testPassword)
	require.NoError(t, err, "le compte n'est pas suspendu")
	assert.Contains(t, auditActions(t, st), "user.sessions_close")
}

func TestDeleteAccount(t *testing.T) {
	svc, st := setupAuthTest(t)
	ctx := context.Background()
	admin := register(t, svc, "admin", store.RoleAdmin)
	bob := register(t, svc, "bob", store.RoleMember)
	session, err := svc.Login(ctx, "bob", testPassword)
	require.NoError(t, err)

	gone, err := svc.DeleteAccount(ctx, admin.ID, bob.ID)
	require.NoError(t, err)
	assert.Equal(t, "bob", gone.Username)

	_, err = svc.ValidateToken(ctx, session.Token)
	require.ErrorIs(t, err, auth.ErrInvalidToken)
	_, err = svc.Login(ctx, "bob", testPassword)
	require.ErrorIs(t, err, auth.ErrInvalidCredentials)

	// Le nom et l'adresse redeviennent libres.
	_, err = svc.Register(ctx, "bob", "bob@forgenet.test", testPassword, store.RoleMember)
	require.NoError(t, err)

	events, err := st.Audit().List(ctx, 100, 0)
	require.NoError(t, err)
	var deleted bool
	for _, e := range events {
		if e.Action == "user.delete" {
			deleted = true
			assert.Equal(t, "user:"+bob.ID, e.Resource)
			assert.Contains(t, e.Metadata, `"username":"bob"`)
			assert.NotContains(t, e.Metadata, "@", "pas d'adresse au journal")
		}
	}
	assert.True(t, deleted)
}
