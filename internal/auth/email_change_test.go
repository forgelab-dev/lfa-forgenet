package auth_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

func TestChangeOwnEmail(t *testing.T) {
	svc, st := setupAuthTest(t)
	ctx := context.Background()
	bob := register(t, svc, "bob", store.RoleMember)
	session, err := svc.Login(ctx, "bob", testPassword)
	require.NoError(t, err)
	_, _, err = svc.CreateAccessLink(ctx, bob.ID, store.PurposePasswordReset, "system:test")
	require.NoError(t, err)

	_, err = svc.ChangeOwnEmail(ctx, bob.ID, "pas-le-bon", "bob@nouveau.test")
	require.ErrorIs(t, err, auth.ErrWrongCurrentPassword, "une session volée ne suffit pas")

	change, err := svc.ChangeOwnEmail(ctx, bob.ID, testPassword, "  Bob@Nouveau.TEST ")
	require.NoError(t, err)
	assert.Equal(t, "bob@nouveau.test", change.User.Email, "adresse normalisée")
	assert.Equal(t, "bob@forgenet.test", change.OldEmail)
	assert.True(t, change.HadPassword)
	assert.Equal(t, int64(1), change.CancelledLinks)

	_, err = svc.Login(ctx, "bob@nouveau.test", testPassword)
	require.NoError(t, err, "la nouvelle adresse sert à se connecter")
	_, err = svc.ValidateToken(ctx, session.Token)
	require.NoError(t, err, "la session reste ouverte")

	actions := auditActions(t, st)
	assert.Contains(t, actions, "user.email_change_failed")
	assert.Contains(t, actions, "user.email_change")
	assertNoAddressInAudit(t, st, "bob@nouveau.test", "bob@forgenet.test")
}

func TestChangeEmail_Refusals(t *testing.T) {
	svc, _ := setupAuthTest(t)
	ctx := context.Background()
	admin := register(t, svc, "admin", store.RoleAdmin)
	bob := register(t, svc, "bob", store.RoleMember)

	for _, bad := range []string{"", "pas-une-adresse", "Bob <bob@x.test>", "bob@x.test\r\nBcc: a@b.test", strings.Repeat("a", 250) + "@x.test"} {
		_, err := svc.SetAccountEmail(ctx, admin.ID, bob.ID, bad)
		require.ErrorIs(t, err, auth.ErrInvalidEmail, "%q", bad)
	}
	_, err := svc.SetAccountEmail(ctx, admin.ID, bob.ID, "BOB@forgenet.test")
	require.ErrorIs(t, err, auth.ErrSameEmail)
	_, err = svc.SetAccountEmail(ctx, admin.ID, bob.ID, "admin@forgenet.test")
	require.ErrorIs(t, err, store.ErrConflict)
	_, err = svc.SetAccountEmail(ctx, admin.ID, "usr_absent", "libre@forgenet.test")
	require.ErrorIs(t, err, store.ErrNotFound)
}

func TestSetAccountEmail_InvitedAccount(t *testing.T) {
	svc, st := setupAuthTest(t)
	ctx := context.Background()
	admin := register(t, svc, "admin", store.RoleAdmin)
	invited, err := svc.RegisterWithoutPassword(ctx, "invite", "invite@mal-saisie.test", store.RoleMember)
	require.NoError(t, err)
	token, _, err := svc.CreateAccessLink(ctx, invited.ID, store.PurposeInvitation, "user:"+admin.ID)
	require.NoError(t, err)

	change, err := svc.SetAccountEmail(ctx, admin.ID, invited.ID, "invite@bonne.test")
	require.NoError(t, err)
	assert.False(t, change.HadPassword)
	assert.Equal(t, int64(1), change.CancelledLinks)

	// L'invitation partie à l'adresse mal saisie ne vaut plus.
	_, _, err = svc.InspectAccessLink(ctx, token)
	require.ErrorIs(t, err, auth.ErrLinkInvalid)

	events, err := st.Audit().List(ctx, 100, 0)
	require.NoError(t, err)
	var found bool
	for _, e := range events {
		if e.Action == "user.email_change" {
			found = true
			assert.Equal(t, admin.ID, e.UserID, "l'administrateur est l'auteur")
			assert.Contains(t, e.Metadata, `"by":"admin"`)
		}
	}
	assert.True(t, found)
	assertNoAddressInAudit(t, st, "invite@mal-saisie.test", "invite@bonne.test")
}

func TestSetAccountEmail_FromTheServer(t *testing.T) {
	svc, st := setupAuthTest(t)
	ctx := context.Background()
	bob := register(t, svc, "bob", store.RoleMember)

	_, err := svc.SetAccountEmail(ctx, "", bob.ID, "bob@nouveau.test")
	require.NoError(t, err)

	events, err := st.Audit().List(ctx, 100, 0)
	require.NoError(t, err)
	for _, e := range events {
		if e.Action == "user.email_change" {
			assert.Empty(t, e.UserID, "aucun compte n'agit sur la machine du serveur")
			assert.Contains(t, e.Metadata, `"by":"server"`)
		}
	}
}

// assertNoAddressInAudit vérifie que le journal ne contient aucune des
// adresses : il est lu par tous les administrateurs, et conservé.
func assertNoAddressInAudit(t *testing.T, st store.Store, addresses ...string) {
	t.Helper()
	events, err := st.Audit().List(context.Background(), 100, 0)
	require.NoError(t, err)
	for _, e := range events {
		for _, a := range addresses {
			assert.NotContains(t, e.Metadata, a, e.Action)
			assert.NotContains(t, e.Resource, a, e.Action)
		}
	}
}
