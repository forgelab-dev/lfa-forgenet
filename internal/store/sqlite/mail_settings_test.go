package sqlite_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

func TestMailSettings_SaveGetDelete(t *testing.T) {
	ctx := context.Background()
	ms := newTestStore(t).MailSettings()

	_, err := ms.GetMailSettings(ctx)
	require.ErrorIs(t, err, store.ErrNotFound, "rien tant que le dashboard n'a rien enregistré")
	gone, err := ms.DeleteMailSettings(ctx)
	require.NoError(t, err)
	assert.False(t, gone)

	first := &store.MailSettings{
		Host: "smtp.exemple.org", Port: 587, TLS: "starttls", Username: "forgenet",
		PasswordSealed: []byte{1, 2, 3}, From: "forgenet@exemple.org", UpdatedBy: "usr_admin",
	}
	require.NoError(t, ms.SaveMailSettings(ctx, first))
	got, err := ms.GetMailSettings(ctx)
	require.NoError(t, err)
	assert.Equal(t, "smtp.exemple.org", got.Host)
	assert.Equal(t, []byte{1, 2, 3}, got.PasswordSealed)
	assert.False(t, got.Disabled)
	assert.False(t, got.UpdatedAt.IsZero())

	// Un second enregistrement remplace le premier d'un bloc.
	require.NoError(t, ms.SaveMailSettings(ctx, &store.MailSettings{
		Disabled: true, Host: "localhost", Port: 25, TLS: "none", From: "forgenet@exemple.org", UpdatedBy: "usr_autre",
	}))
	got, err = ms.GetMailSettings(ctx)
	require.NoError(t, err)
	assert.True(t, got.Disabled)
	assert.Equal(t, 25, got.Port)
	assert.Empty(t, got.Username)
	assert.Empty(t, got.PasswordSealed, "sans mot de passe, plus de chiffré")
	assert.Equal(t, "usr_autre", got.UpdatedBy)

	gone, err = ms.DeleteMailSettings(ctx)
	require.NoError(t, err)
	assert.True(t, gone)
	_, err = ms.GetMailSettings(ctx)
	require.ErrorIs(t, err, store.ErrNotFound)
}
