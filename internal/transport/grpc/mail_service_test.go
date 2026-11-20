package grpc_test

import (
	"bytes"
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/mail"
	"github.com/forgelab-dev/lfa-forgenet/internal/sealed"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// fakeSender garde les messages ; err est renvoyée à la place d'une remise.
type fakeSender struct {
	mu   sync.Mutex
	sent []mail.Message
	err  error
}

func (f *fakeSender) Send(_ context.Context, m mail.Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return f.err
	}
	f.sent = append(f.sent, m)
	return nil
}

var relay = config.SMTPConfig{
	Host: "smtp.exemple.org", Port: 587, Username: "forgenet", Password: "secret-du-relais",
	From: "ForgeNet <forgenet@exemple.org>", TLS: config.SMTPStartTLS,
}

// relayLog note les réglages de chaque envoi : c'est le relais qui reçoit le
// message, pas seulement son contenu, que les tests vérifient.
type relayLog struct {
	mu      sync.Mutex
	configs []config.SMTPConfig
}

func (l *relayLog) last() config.SMTPConfig {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.configs) == 0 {
		return config.SMTPConfig{}
	}
	return l.configs[len(l.configs)-1]
}

func testSealedBox(t *testing.T) *sealed.Box {
	t.Helper()
	box, err := sealed.NewBox(bytes.Repeat([]byte{3}, sealed.KeySize))
	require.NoError(t, err)
	return box
}

func startMailServer(t *testing.T, cfg config.SMTPConfig, sender mail.Sender) (*testServer, forgenetv1.MailServiceClient) {
	t.Helper()
	ts, client, _, _ := startMailServerWith(t, cfg, sender, testSealedBox(t))
	return ts, client
}

// startMailServerWith sert le service Mail sur un relais dont chaque envoi
// passe par sender ; box nil : serveur sans FORGENET_SECRET_KEY.
func startMailServerWith(t *testing.T, cfg config.SMTPConfig, sender mail.Sender, box *sealed.Box) (*testServer, forgenetv1.MailServiceClient, *mail.Relay, *relayLog) {
	t.Helper()
	log := &relayLog{}
	var r *mail.Relay
	ts := startTestServer(t, &fakeEngine{}, func(srv *grpc.Server, st store.Store, auditor audit.Auditor) {
		r = mail.NewRelay(cfg, st.MailSettings(), box, mail.WithSenderFactory(func(c config.SMTPConfig) mail.Sender {
			log.mu.Lock()
			log.configs = append(log.configs, c)
			log.mu.Unlock()
			return sender
		}))
		require.NoError(t, r.Load(context.Background()))
		forgenetv1.RegisterMailServiceServer(srv, transportGRPC.NewMailServiceServer(r, st.Users(), auditor, zerolog.Nop()))
	})
	return ts, forgenetv1.NewMailServiceClient(ts.conn), r, log
}

func TestMailService_ReservedToAdmins(t *testing.T) {
	ts, client := startMailServer(t, relay, &fakeSender{})
	_, err := client.GetMailStatus(as(ts.memberToken), &forgenetv1.GetMailStatusRequest{})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = client.SendTestMail(as(ts.memberToken), &forgenetv1.SendTestMailRequest{})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = client.SaveMailSettings(as(ts.memberToken), &forgenetv1.SaveMailSettingsRequest{Host: "smtp.autre.org"})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = client.DisableMail(as(ts.memberToken), &forgenetv1.DisableMailRequest{})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
	_, err = client.ResetMailSettings(as(ts.memberToken), &forgenetv1.ResetMailSettingsRequest{})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
}

func TestMailService_StatusHidesCredentials(t *testing.T) {
	ts, client := startMailServer(t, relay, &fakeSender{})
	resp, err := client.GetMailStatus(as(ts.adminToken), &forgenetv1.GetMailStatusRequest{})
	require.NoError(t, err)
	assert.True(t, resp.GetEnabled())
	assert.Equal(t, "smtp.exemple.org", resp.GetHost())
	assert.Equal(t, uint32(587), resp.GetPort())
	assert.Equal(t, forgenetv1.MailTLS_MAIL_TLS_STARTTLS, resp.GetTls())
	assert.True(t, resp.GetAuthenticated())
	assert.NotContains(t, resp.String(), "secret-du-relais")

	ts, client = startMailServer(t, config.SMTPConfig{}, &fakeSender{})
	resp, err = client.GetMailStatus(as(ts.adminToken), &forgenetv1.GetMailStatusRequest{})
	require.NoError(t, err)
	assert.False(t, resp.GetEnabled())
}

func TestMailService_SendTestMailToTheAdmin(t *testing.T) {
	sender := &fakeSender{}
	ts, client := startMailServer(t, relay, sender)
	resp, err := client.SendTestMail(as(ts.adminToken), &forgenetv1.SendTestMailRequest{})
	require.NoError(t, err)
	assert.Equal(t, "admin@forgenet.test", resp.GetTo())
	require.Len(t, sender.sent, 1)
	assert.Equal(t, "admin@forgenet.test", sender.sent[0].To)
	assert.Contains(t, sender.sent[0].Body, "smtp.exemple.org")

	events, err := ts.store.Audit().List(context.Background(), 10, 0)
	require.NoError(t, err)
	var found bool
	for _, e := range events {
		found = found || e.Action == "mail.test"
	}
	assert.True(t, found, "essai inscrit au journal d'audit")
}

func TestMailService_SendTestMailErrors(t *testing.T) {
	ts, client := startMailServer(t, config.SMTPConfig{}, &fakeSender{})
	_, err := client.SendTestMail(as(ts.adminToken), &forgenetv1.SendTestMailRequest{})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	assert.Contains(t, status.Convert(err).Message(), "pas configuré")

	ts, client = startMailServer(t, relay, &fakeSender{err: errors.New(`authentification refusée par le relais : 535 "identifiants refusés"`)})
	_, err = client.SendTestMail(as(ts.adminToken), &forgenetv1.SendTestMailRequest{})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	assert.Contains(t, status.Convert(err).Message(), "535")
}

func TestMailService_SendTestMailIsThrottled(t *testing.T) {
	ts, client := startMailServer(t, relay, &fakeSender{})
	for range 10 {
		_, err := client.SendTestMail(as(ts.adminToken), &forgenetv1.SendTestMailRequest{})
		require.NoError(t, err)
	}
	_, err := client.SendTestMail(as(ts.adminToken), &forgenetv1.SendTestMailRequest{})
	assert.Equal(t, codes.ResourceExhausted, status.Code(err))
}
