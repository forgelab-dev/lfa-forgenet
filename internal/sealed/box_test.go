package sealed_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/sealed"
)

func newBox(t *testing.T, fill byte) *sealed.Box {
	t.Helper()
	box, err := sealed.NewBox(bytes.Repeat([]byte{fill}, sealed.KeySize))
	require.NoError(t, err)
	return box
}

func TestBox_SealOpen(t *testing.T) {
	box := newBox(t, 7)
	ct, err := box.Seal("mail_settings.password", []byte("motdepasse-du-relais"))
	require.NoError(t, err)
	assert.NotContains(t, string(ct), "motdepasse", "jamais en clair")

	plain, err := box.Open("mail_settings.password", ct)
	require.NoError(t, err)
	assert.Equal(t, "motdepasse-du-relais", string(plain))

	again, err := box.Seal("mail_settings.password", []byte("motdepasse-du-relais"))
	require.NoError(t, err)
	assert.NotEqual(t, ct, again, "un nonce neuf à chaque chiffrement")
}

func TestBox_Refusals(t *testing.T) {
	box := newBox(t, 7)
	ct, err := box.Seal("mail_settings.password", []byte("secret"))
	require.NoError(t, err)

	_, err = newBox(t, 8).Open("mail_settings.password", ct)
	require.ErrorIs(t, err, sealed.ErrUnreadable, "autre clé")

	_, err = box.Open("autre.usage", ct)
	require.ErrorIs(t, err, sealed.ErrUnreadable, "autre usage")

	damaged := bytes.Clone(ct)
	damaged[len(damaged)-1] ^= 1
	_, err = box.Open("mail_settings.password", damaged)
	require.ErrorIs(t, err, sealed.ErrUnreadable, "chiffré abîmé")

	unknown := bytes.Clone(ct)
	unknown[0] = 9
	_, err = box.Open("mail_settings.password", unknown)
	require.ErrorIs(t, err, sealed.ErrUnreadable, "version inconnue")

	for _, short := range [][]byte{nil, {1}, ct[:10]} {
		_, err = box.Open("mail_settings.password", short)
		require.ErrorIs(t, err, sealed.ErrUnreadable)
	}
}

func TestParseKey(t *testing.T) {
	key, err := sealed.ParseKey(" " + strings.Repeat("ab", sealed.KeySize) + "\n")
	require.NoError(t, err)
	assert.Len(t, key, sealed.KeySize)

	for _, bad := range []string{"", "abc", strings.Repeat("ab", 16), strings.Repeat("zz", sealed.KeySize)} {
		_, err := sealed.ParseKey(bad)
		assert.Error(t, err, bad)
	}
	_, err = sealed.NewBox([]byte("trop court"))
	assert.Error(t, err)
}
