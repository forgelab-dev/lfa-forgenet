package mail

import (
	"bytes"
	"context"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/sealed"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
)

// hostRecorder note, pour chaque message, le relais par lequel il est parti.
type hostRecorder struct {
	mu    sync.Mutex
	hosts []string
	fail  map[string]error
}

func (h *hostRecorder) factory(cfg config.SMTPConfig) Sender {
	return senderFunc(func(context.Context, Message) error {
		h.mu.Lock()
		defer h.mu.Unlock()
		if err := h.fail[cfg.Host]; err != nil {
			return err
		}
		h.hosts = append(h.hosts, cfg.Host+"|"+cfg.Password)
		return nil
	})
}

func (h *hostRecorder) sent() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]string(nil), h.hosts...)
}

type senderFunc func(context.Context, Message) error

func (f senderFunc) Send(ctx context.Context, m Message) error { return f(ctx, m) }

var serverSMTP = config.SMTPConfig{
	Host: "smtp.serveur.org", Port: 587, TLS: config.SMTPStartTLS,
	Username: "env", Password: "mdp-env", From: "forgenet@exemple.org",
}

var dashboardSMTP = config.SMTPConfig{
	Host: "smtp.dashboard.org", Port: 465, TLS: config.SMTPImplicit,
	Username: "dash", Password: "mdp-dashboard", From: "forgenet@exemple.org",
}

func testBox(t *testing.T, fill byte) *sealed.Box {
	t.Helper()
	box, err := sealed.NewBox(bytes.Repeat([]byte{fill}, sealed.KeySize))
	require.NoError(t, err)
	return box
}

func testStore(t *testing.T) store.Store {
	t.Helper()
	st, err := sqlite.Open(context.Background(), config.DatabaseConfig{Path: ":memory:", BusyTimeout: time.Second, MaxOpenConns: 1, MaxIdleConns: 1})
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })
	return st
}

func TestRelay_ServerSettingsUntilTheDashboardSaves(t *testing.T) {
	ctx := context.Background()
	st := testStore(t)
	rec := &hostRecorder{}
	r := NewRelay(serverSMTP, st.MailSettings(), testBox(t, 1), WithSenderFactory(rec.factory))
	require.NoError(t, r.Load(ctx))
	assert.Equal(t, SourceServer, r.Current().Source)
	assert.True(t, r.Editable())

	require.NoError(t, r.Send(ctx, essai))
	require.NoError(t, r.Save(ctx, dashboardSMTP, "usr_admin"))
	require.NoError(t, r.Send(ctx, essai))
	assert.Equal(t, []string{"smtp.serveur.org|mdp-env", "smtp.dashboard.org|mdp-dashboard"}, rec.sent(),
		"le message suivant prend le nouveau relais, sans redémarrage")

	stored, err := st.MailSettings().GetMailSettings(ctx)
	require.NoError(t, err)
	assert.NotContains(t, string(stored.PasswordSealed), "mdp-dashboard", "jamais en clair en base")

	// Au redémarrage, la base l'emporte sur .env.
	again := NewRelay(serverSMTP, st.MailSettings(), testBox(t, 1), WithSenderFactory(rec.factory))
	require.NoError(t, again.Load(ctx))
	cur := again.Current()
	assert.Equal(t, SourceDashboard, cur.Source)
	assert.Equal(t, dashboardSMTP, cur.Config)
	assert.Equal(t, "usr_admin", cur.UpdatedBy)
	assert.True(t, cur.Enabled)

	gone, err := again.Reset(ctx)
	require.NoError(t, err)
	assert.True(t, gone)
	assert.Equal(t, SourceServer, again.Current().Source)
	assert.Equal(t, serverSMTP, again.Current().Config)
}

func TestRelay_DisableKeepsTheValues(t *testing.T) {
	ctx := context.Background()
	st := testStore(t)
	rec := &hostRecorder{}
	r := NewRelay(serverSMTP, st.MailSettings(), testBox(t, 1), WithSenderFactory(rec.factory))
	require.NoError(t, r.Load(ctx))

	require.NoError(t, r.Disable(ctx, "usr_admin"))
	assert.False(t, r.Enabled())
	require.ErrorIs(t, r.Send(ctx, essai), ErrDisabled)
	assert.Empty(t, rec.sent())

	again := NewRelay(serverSMTP, st.MailSettings(), testBox(t, 1))
	require.NoError(t, again.Load(ctx))
	cur := again.Current()
	assert.False(t, cur.Enabled, "coupé depuis le dashboard, l'envoi le reste au redémarrage malgré .env")
	assert.Equal(t, SourceDashboard, cur.Source)
	assert.Equal(t, serverSMTP, cur.Config, "les valeurs restent pour rouvrir l'envoi")
}

func TestRelay_UnreadablePasswordFallsBackToTheServer(t *testing.T) {
	ctx := context.Background()
	st := testStore(t)
	r := NewRelay(serverSMTP, st.MailSettings(), testBox(t, 1))
	require.NoError(t, r.Save(ctx, dashboardSMTP, "usr_admin"))

	for name, box := range map[string]*sealed.Box{"clé changée": testBox(t, 2), "clé absente": nil} {
		t.Run(name, func(t *testing.T) {
			other := NewRelay(serverSMTP, st.MailSettings(), box)
			require.NoError(t, other.Load(ctx))
			cur := other.Current()
			assert.True(t, cur.Unreadable)
			assert.Equal(t, SourceServer, cur.Source)
			assert.Equal(t, serverSMTP, cur.Config)
		})
	}
}

func TestRelay_NotEditableWithoutAKey(t *testing.T) {
	ctx := context.Background()
	r := NewRelay(serverSMTP, testStore(t).MailSettings(), nil)
	assert.False(t, r.Editable())
	require.ErrorIs(t, r.Save(ctx, dashboardSMTP, "usr_admin"), ErrNotEditable)
	require.ErrorIs(t, r.Disable(ctx, "usr_admin"), ErrNotEditable)
	assert.Equal(t, SourceServer, r.Current().Source)
}

func TestRelay_TryChangesNothing(t *testing.T) {
	ctx := context.Background()
	rec := &hostRecorder{fail: map[string]error{"smtp.dashboard.org": assert.AnError}}
	r := NewRelay(serverSMTP, testStore(t).MailSettings(), testBox(t, 1), WithSenderFactory(rec.factory))
	require.ErrorIs(t, r.Try(ctx, dashboardSMTP, essai), assert.AnError)
	require.ErrorIs(t, r.Try(ctx, config.SMTPConfig{}, essai), ErrDisabled)
	assert.Equal(t, serverSMTP, r.Current().Config)
}

func TestOutbox_FollowsTheRelay(t *testing.T) {
	ctx := context.Background()
	r := NewRelay(config.SMTPConfig{}, testStore(t).MailSettings(), testBox(t, 1))
	o := NewOutbox(r, nil, zerolog.Nop())
	assert.False(t, o.Enabled())
	require.NoError(t, r.Save(ctx, dashboardSMTP, "usr_admin"))
	assert.True(t, o.Enabled(), "allumé sans redémarrage")
	assert.True(t, NewOutbox(&flakySender{}, nil, zerolog.Nop()).Enabled(), "un Sender sans état est réglé")
}
