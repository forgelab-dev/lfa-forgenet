package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

func TestUserStore_SetEmail(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	a := accounts{t, ctx, st}
	a.user("usr_m", store.RoleMember, "h")
	a.user("usr_other", store.RoleMember, "h")
	a.session("sess_m", "usr_m")
	a.link("act_m", "usr_m")

	cancelled, err := st.Users().SetEmail(ctx, "usr_m", "nouvelle@forgenet.local")
	require.NoError(t, err)
	assert.Equal(t, int64(1), cancelled)

	u, err := st.Users().GetUserByEmail(ctx, "nouvelle@forgenet.local")
	require.NoError(t, err)
	assert.Equal(t, "usr_m", u.ID)
	_, err = st.Users().GetUserByEmail(ctx, "usr_m@forgenet.local")
	require.ErrorIs(t, err, store.ErrNotFound)

	// Un lien parti vers l'ancienne adresse ne vaut plus : elle était
	// peut-être celle d'un inconnu.
	_, err = st.ActionTokens().GetActiveActionToken(ctx, "h_act_m", time.Now())
	require.ErrorIs(t, err, store.ErrNotFound)
	// Les sessions restent : changer d'adresse ne déconnecte personne.
	assert.False(t, a.revoked("sess_m"))
}

func TestUserStore_SetEmailRefusals(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	a := accounts{t, ctx, st}
	a.user("usr_m", store.RoleMember, "h")
	a.user("usr_other", store.RoleMember, "h")
	a.link("act_m", "usr_m")

	_, err := st.Users().SetEmail(ctx, "usr_m", "usr_other@forgenet.local")
	require.ErrorIs(t, err, store.ErrConflict)
	// Le refus ne touche à rien : le lien du compte vaut toujours.
	_, err = st.ActionTokens().GetActiveActionToken(ctx, "h_act_m", time.Now())
	require.NoError(t, err)

	_, err = st.Users().SetEmail(ctx, "usr_absent", "libre@forgenet.local")
	require.ErrorIs(t, err, store.ErrNotFound)
}
