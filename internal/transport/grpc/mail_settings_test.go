package grpc_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/mail"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// newSettings sont les réglages qu'un administrateur saisit dans la page.
func newSettings() *forgenetv1.SaveMailSettingsRequest {
	return &forgenetv1.SaveMailSettingsRequest{
		Host: " SMTP.Nouveau.org ", Port: 465, Tls: forgenetv1.MailTLS_MAIL_TLS_IMPLICIT,
		Username: "envois@nouveau.org", Password: "mdp-du-nouveau-relais", From: "envois@nouveau.org",
	}
}

func auditEvents(t *testing.T, st store.Store, action string) []*store.AuditEvent {
	t.Helper()
	events, err := st.Audit().List(context.Background(), 100, 0)
	require.NoError(t, err)
	var out []*store.AuditEvent
	for _, e := range events {
		if e.Action == action {
			out = append(out, e)
		}
	}
	return out
}

func TestMailService_SaveTriesThenApplies(t *testing.T) {
	sender := &fakeSender{}
	ts, client, r, log := startMailServerWith(t, relay, sender, testSealedBox(t))

	res, err := client.SaveMailSettings(as(ts.adminToken), newSettings())
	require.NoError(t, err)
	assert.Equal(t, "admin@forgenet.test", res.GetTo())

	// L'essai est parti par le nouveau relais, à l'administrateur connecté.
	require.Len(t, sender.sent, 1)
	assert.Equal(t, "admin@forgenet.test", sender.sent[0].To)
	assert.Contains(t, sender.sent[0].Body, "smtp.nouveau.org")
	assert.Equal(t, "smtp.nouveau.org", log.last().Host)
	assert.Equal(t, "mdp-du-nouveau-relais", log.last().Password)

	st := res.GetStatus()
	assert.True(t, st.GetEnabled())
	assert.Equal(t, forgenetv1.MailSource_MAIL_SOURCE_DASHBOARD, st.GetSource())
	assert.Equal(t, "smtp.nouveau.org", st.GetHost())
	assert.Equal(t, uint32(465), st.GetPort())
	assert.Equal(t, "envois@nouveau.org", st.GetUsername())
	assert.True(t, st.GetAuthenticated())
	assert.Equal(t, "admin", st.GetUpdatedBy())
	assert.NotNil(t, st.GetUpdatedAt())
	assert.NotContains(t, res.String(), "mdp-du-nouveau-relais", "le mot de passe ne revient jamais au navigateur")

	// Appliqué sans redémarrage : le message suivant prend le nouveau relais.
	require.NoError(t, r.Send(context.Background(), mail.Message{To: "x@exemple.org", Subject: "s", Body: "b"}))
	assert.Equal(t, "smtp.nouveau.org", log.last().Host)

	events := auditEvents(t, ts.store, "mail.settings_change")
	require.Len(t, events, 1)
	assert.Contains(t, events[0].Metadata, `"password":"changé"`)
	assert.Contains(t, events[0].Metadata, `"host":"smtp.nouveau.org"`)
	for _, secret := range []string{"mdp-du-nouveau-relais", "envois@nouveau.org", "@"} {
		assert.NotContains(t, events[0].Metadata, secret, "ni mot de passe, ni adresse au journal")
	}
}

func TestMailService_RefusedTryKeepsTheOldSettings(t *testing.T) {
	sender := &fakeSender{err: errors.New(`authentification refusée par le relais : 535 "identifiants refusés"`)}
	ts, client, r, _ := startMailServerWith(t, relay, sender, testSealedBox(t))

	_, err := client.SaveMailSettings(as(ts.adminToken), newSettings())
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	assert.Contains(t, status.Convert(err).Message(), "535")
	assert.Contains(t, status.Convert(err).Message(), "les réglages en place restent")

	assert.Equal(t, relay, r.Current().Config, "les anciens réglages restent")
	assert.Equal(t, mail.SourceServer, r.Current().Source)
	_, err = ts.store.MailSettings().GetMailSettings(context.Background())
	require.ErrorIs(t, err, store.ErrNotFound, "rien d'enregistré")

	failed := auditEvents(t, ts.store, "mail.settings_change_failed")
	require.Len(t, failed, 1)
	assert.Contains(t, failed[0].Metadata, `"result":"refusé"`)
	assert.Empty(t, auditEvents(t, ts.store, "mail.settings_change"))
}

func TestMailService_KeepPassword(t *testing.T) {
	ts, client, _, log := startMailServerWith(t, relay, &fakeSender{}, testSealedBox(t))
	req := newSettings()
	req.Username, req.Password, req.KeepPassword = "forgenet", "ignoré", true

	res, err := client.SaveMailSettings(as(ts.adminToken), req)
	require.NoError(t, err)
	assert.Equal(t, "secret-du-relais", log.last().Password, "le mot de passe en vigueur est gardé")
	assert.True(t, res.GetStatus().GetAuthenticated())
	assert.Contains(t, auditEvents(t, ts.store, "mail.settings_change")[0].Metadata, `"password":"gardé"`)

	// Sans identifiant, pas d'authentification : aucun mot de passe gardé.
	req = newSettings()
	req.Username, req.Password = "", ""
	res, err = client.SaveMailSettings(as(ts.adminToken), req)
	require.NoError(t, err)
	assert.False(t, res.GetStatus().GetAuthenticated())
	assert.Empty(t, log.last().Password)
}

func TestMailService_SaveRefusals(t *testing.T) {
	ts, client, _, log := startMailServerWith(t, config.SMTPConfig{}, &fakeSender{}, testSealedBox(t))
	for name, edit := range map[string]func(*forgenetv1.SaveMailSettingsRequest){
		"sans relais":                   func(r *forgenetv1.SaveMailSettingsRequest) { r.Host = " " },
		"expéditeur invalide":           func(r *forgenetv1.SaveMailSettingsRequest) { r.From = "pas une adresse" },
		"sans protection vers dehors":   func(r *forgenetv1.SaveMailSettingsRequest) { r.Tls = forgenetv1.MailTLS_MAIL_TLS_NONE },
		"port nul":                      func(r *forgenetv1.SaveMailSettingsRequest) { r.Port = 0 },
		"identifiant sans mot de passe": func(r *forgenetv1.SaveMailSettingsRequest) { r.Password = "" },
		"mot de passe gardé absent":     func(r *forgenetv1.SaveMailSettingsRequest) { r.KeepPassword = true },
		"saut de ligne":                 func(r *forgenetv1.SaveMailSettingsRequest) { r.Password = "a\r\nb" },
		"protection inconnue":           func(r *forgenetv1.SaveMailSettingsRequest) { r.Tls = 42 },
	} {
		req := newSettings()
		edit(req)
		_, err := client.SaveMailSettings(as(ts.adminToken), req)
		assert.Equal(t, codes.InvalidArgument, status.Code(err), name)
	}
	assert.Empty(t, log.configs, "aucun essai pour des réglages invalides")
}

func TestMailService_NotEditableWithoutKey(t *testing.T) {
	ts, client, _, _ := startMailServerWith(t, relay, &fakeSender{}, nil)
	st, err := client.GetMailStatus(as(ts.adminToken), &forgenetv1.GetMailStatusRequest{})
	require.NoError(t, err)
	assert.False(t, st.GetEditable())
	assert.True(t, st.GetEnabled(), "l'envoi de .env marche toujours")

	_, err = client.SaveMailSettings(as(ts.adminToken), newSettings())
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	assert.Contains(t, status.Convert(err).Message(), "FORGENET_SECRET_KEY")
	_, err = client.DisableMail(as(ts.adminToken), &forgenetv1.DisableMailRequest{})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
}

func TestMailService_DisableThenReset(t *testing.T) {
	ts, client, r, _ := startMailServerWith(t, relay, &fakeSender{}, testSealedBox(t))

	res, err := client.DisableMail(as(ts.adminToken), &forgenetv1.DisableMailRequest{})
	require.NoError(t, err)
	st := res.GetStatus()
	assert.False(t, st.GetEnabled())
	assert.True(t, st.GetDisabled())
	assert.Equal(t, "smtp.exemple.org", st.GetHost(), "les valeurs restent pour rouvrir l'envoi")
	assert.False(t, r.Enabled())

	_, err = client.SendTestMail(as(ts.adminToken), &forgenetv1.SendTestMailRequest{})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))

	// Rouvrir en gardant le mot de passe enregistré.
	req := &forgenetv1.SaveMailSettingsRequest{Host: st.GetHost(), Port: st.GetPort(), Tls: st.GetTls(), Username: st.GetUsername(), From: st.GetFrom(), KeepPassword: true}
	saved, err := client.SaveMailSettings(as(ts.adminToken), req)
	require.NoError(t, err)
	assert.True(t, saved.GetStatus().GetEnabled())

	reset, err := client.ResetMailSettings(as(ts.adminToken), &forgenetv1.ResetMailSettingsRequest{})
	require.NoError(t, err)
	assert.Equal(t, forgenetv1.MailSource_MAIL_SOURCE_SERVER, reset.GetStatus().GetSource())
	assert.True(t, reset.GetStatus().GetEnabled())

	var actions []string
	for _, e := range auditEvents(t, ts.store, "mail.settings_change") {
		actions = append(actions, e.Metadata)
	}
	joined := strings.Join(actions, " ")
	for _, want := range []string{`"action":"coupé"`, `"action":"enregistré"`, `"action":"rendu au serveur"`} {
		assert.Contains(t, joined, want)
	}
}

func TestMailService_StatusSaysWhenThePasswordIsUnreadable(t *testing.T) {
	ts, client, _, _ := startMailServerWith(t, relay, &fakeSender{}, testSealedBox(t))
	_, err := client.SaveMailSettings(as(ts.adminToken), newSettings())
	require.NoError(t, err)

	// Même base, autre clé : le serveur reprend .env et le dit.
	other := mail.NewRelay(relay, ts.store.MailSettings(), nil)
	require.NoError(t, other.Load(context.Background()))
	assert.True(t, other.Current().Unreadable)
	assert.Equal(t, relay, other.Current().Config)
}

// Couper ou rouvrir l'envoi depuis la page Réglages change aussitôt ce que
// proposent l'écran de connexion et la page Équipe (ADR-0015).
func TestMailState_FollowedWithoutRestart(t *testing.T) {
	_, client, outbox := startResetServer(t, "forgenet.exemple.org")
	outbox.SetOff(true)
	opts, err := client.GetLoginOptions(context.Background(), &forgenetv1.GetLoginOptionsRequest{})
	require.NoError(t, err)
	assert.False(t, opts.GetPasswordReset(), "envoi coupé : pas de « mot de passe oublié »")
	_, err = client.RequestPasswordReset(context.Background(), &forgenetv1.RequestPasswordResetRequest{Login: "membre"})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))

	outbox.SetOff(false)
	opts, err = client.GetLoginOptions(context.Background(), &forgenetv1.GetLoginOptionsRequest{})
	require.NoError(t, err)
	assert.True(t, opts.GetPasswordReset(), "rouvert sans redémarrage")

	teamOutbox := &transportGRPC.SyncOutboxForTests{}
	teamOutbox.SetOff(true)
	f := startTeamServer(t, teamOutbox, "")
	list, err := f.team.ListMembers(as(f.ts.adminToken), &forgenetv1.ListMembersRequest{})
	require.NoError(t, err)
	assert.False(t, list.GetMailEnabled())
	res := f.invite(t, "nadia")
	assert.False(t, res.GetEmailed())
	assert.NotEmpty(t, res.GetInvitationUrl(), "envoi coupé : le lien revient à l'administrateur")
	assert.Empty(t, teamOutbox.Messages())

	teamOutbox.SetOff(false)
	list, err = f.team.ListMembers(as(f.ts.adminToken), &forgenetv1.ListMembersRequest{})
	require.NoError(t, err)
	assert.True(t, list.GetMailEnabled())
	assert.True(t, f.invite(t, "sami").GetEmailed())
}
