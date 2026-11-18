package sqlite_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

func seedLinkUser(t *testing.T, st store.Store, id string) {
	t.Helper()
	require.NoError(t, st.Users().CreateUser(context.Background(), &store.User{
		ID: id, Username: id, Email: id + "@forgenet.local", PasswordHash: "!", Role: store.RoleAdmin,
	}))
}

func newLink(t *testing.T, st store.Store, id, userID string, purpose store.ActionPurpose, expires time.Time) {
	t.Helper()
	require.NoError(t, st.ActionTokens().CreateActionToken(context.Background(), &store.ActionToken{
		ID: id, UserID: userID, Purpose: purpose, TokenHash: "h_" + id, ExpiresAt: expires, CreatedBy: "system:cli",
	}))
}

func TestActionTokens_ActiveOnlyBeforeExpiryAndUse(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	now := time.Now().UTC()
	seedLinkUser(t, st, "usr_a")
	newLink(t, st, "act_1", "usr_a", store.PurposeFirstLogin, now.Add(24*time.Hour))

	got, err := st.ActionTokens().GetActiveActionToken(ctx, "h_act_1", now)
	require.NoError(t, err)
	assert.Equal(t, "usr_a", got.UserID)
	assert.Equal(t, store.PurposeFirstLogin, got.Purpose)
	assert.Nil(t, got.UsedAt)

	_, err = st.ActionTokens().GetActiveActionToken(ctx, "h_act_1", now.Add(25*time.Hour))
	require.ErrorIs(t, err, store.ErrNotFound, "expiré")
	_, err = st.ActionTokens().GetActiveActionToken(ctx, "h_inconnu", now)
	require.ErrorIs(t, err, store.ErrNotFound)
}

func TestActionTokens_NewLinkReplacesThePreviousOne(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	now := time.Now().UTC()
	seedLinkUser(t, st, "usr_a")
	newLink(t, st, "act_ancien", "usr_a", store.PurposeFirstLogin, now.Add(time.Hour))
	newLink(t, st, "act_reset", "usr_a", store.PurposePasswordReset, now.Add(time.Hour))
	newLink(t, st, "act_nouveau", "usr_a", store.PurposeFirstLogin, now.Add(time.Hour))

	_, err := st.ActionTokens().GetActiveActionToken(ctx, "h_act_ancien", now)
	require.ErrorIs(t, err, store.ErrNotFound, "le nouveau lien annule l'ancien")
	_, err = st.ActionTokens().GetActiveActionToken(ctx, "h_act_nouveau", now)
	require.NoError(t, err)
	_, err = st.ActionTokens().GetActiveActionToken(ctx, "h_act_reset", now)
	require.NoError(t, err, "un lien d'un autre usage reste valable")
}

func TestActionTokens_RedeemSetsThePasswordOnce(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	now := time.Now().UTC()
	seedLinkUser(t, st, "usr_a")
	require.NoError(t, st.Sessions().CreateSession(ctx, &store.Session{
		ID: "sess_a", UserID: "usr_a", TokenHash: "h_sess_a", Name: "web", Role: store.RoleAdmin, ExpiresAt: now.Add(time.Hour),
	}))
	newLink(t, st, "act_1", "usr_a", store.PurposeFirstLogin, now.Add(time.Hour))
	newLink(t, st, "act_reset", "usr_a", store.PurposePasswordReset, now.Add(time.Hour))

	token, revoked, err := st.ActionTokens().RedeemActionToken(ctx, "h_act_1", now, "empreinte")
	require.NoError(t, err)
	assert.Equal(t, "usr_a", token.UserID)
	require.NotNil(t, token.UsedAt)
	assert.Equal(t, int64(1), revoked)

	u, err := st.Users().GetUserByID(ctx, "usr_a")
	require.NoError(t, err)
	assert.Equal(t, "empreinte", u.PasswordHash)
	sess, err := st.Sessions().GetSessionByID(ctx, "sess_a")
	require.NoError(t, err)
	assert.True(t, sess.Revoked)

	_, _, err = st.ActionTokens().RedeemActionToken(ctx, "h_act_1", now, "autre")
	require.ErrorIs(t, err, store.ErrNotFound, "un lien ne sert qu'une fois")
	_, err = st.ActionTokens().GetActiveActionToken(ctx, "h_act_reset", now)
	require.ErrorIs(t, err, store.ErrNotFound, "le mot de passe choisi annule les autres liens du compte")
	u, err = st.Users().GetUserByID(ctx, "usr_a")
	require.NoError(t, err)
	assert.Equal(t, "empreinte", u.PasswordHash)
}

func TestActionTokens_RedeemRefusesAnExpiredLink(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	now := time.Now().UTC()
	seedLinkUser(t, st, "usr_a")
	newLink(t, st, "act_1", "usr_a", store.PurposeFirstLogin, now.Add(-time.Second))

	_, _, err := st.ActionTokens().RedeemActionToken(ctx, "h_act_1", now, "empreinte")
	require.ErrorIs(t, err, store.ErrNotFound)
	u, err := st.Users().GetUserByID(ctx, "usr_a")
	require.NoError(t, err)
	assert.Equal(t, "!", u.PasswordHash, "mot de passe inchangé")
}

func TestActionTokens_ConcurrentRedeemsSucceedOnce(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	now := time.Now().UTC()
	seedLinkUser(t, st, "usr_a")
	newLink(t, st, "act_1", "usr_a", store.PurposeFirstLogin, now.Add(time.Hour))

	var (
		wg sync.WaitGroup
		mu sync.Mutex
		ok int
	)
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, _, err := st.ActionTokens().RedeemActionToken(ctx, "h_act_1", now, "empreinte"); err == nil {
				mu.Lock()
				ok++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	assert.Equal(t, 1, ok)
}

func TestSetPassword_CancelsPendingLinks(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	now := time.Now().UTC()
	seedLinkUser(t, st, "usr_a")
	newLink(t, st, "act_reset", "usr_a", store.PurposePasswordReset, now.Add(time.Hour))

	_, err := st.Users().SetPassword(ctx, "usr_a", "nouvelle", "")
	require.NoError(t, err)
	_, err = st.ActionTokens().GetActiveActionToken(ctx, "h_act_reset", now)
	require.ErrorIs(t, err, store.ErrNotFound)
}

func TestActionTokens_DeletedWithTheUser(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()
	now := time.Now().UTC()
	seedLinkUser(t, st, "usr_a")
	newLink(t, st, "act_1", "usr_a", store.PurposeInvitation, now.Add(time.Hour))

	require.NoError(t, st.Users().DeleteUser(ctx, "usr_a"))
	_, err := st.ActionTokens().GetActiveActionToken(ctx, "h_act_1", now)
	require.ErrorIs(t, err, store.ErrNotFound)
}
