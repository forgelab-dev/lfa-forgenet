package mail

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// flakySender échoue failures fois, puis réussit ; il signale chaque essai.
type flakySender struct {
	mu       sync.Mutex
	failures int
	sent     []Message
	attempts chan struct{}
}

func (f *flakySender) Send(_ context.Context, msg Message) error {
	defer func() { f.attempts <- struct{}{} }()
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failures > 0 {
		f.failures--
		return errors.New("550 boîte <ana@exemple.org> inconnue")
	}
	f.sent = append(f.sent, msg)
	return nil
}

// syncBuffer est un tampon de journal lisible pendant que la file écrit.
type syncBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

func startOutbox(t *testing.T, sender Sender, delays []time.Duration) (*Outbox, *syncBuffer) {
	t.Helper()
	logs := &syncBuffer{}
	o := NewOutbox(sender, delays, zerolog.New(logs))
	go func() { _ = o.Start() }()
	<-o.started
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		require.NoError(t, o.Shutdown(ctx))
	})
	return o, logs
}

func job(msg *Message) Job {
	return Job{Kind: "password_reset", Build: func(context.Context) (*Message, string, error) {
		return msg, "usr_ana", nil
	}}
}

func waitAttempts(t *testing.T, f *flakySender, n int) {
	t.Helper()
	for i := range n {
		select {
		case <-f.attempts:
		case <-time.After(2 * time.Second):
			t.Fatalf("essai %d attendu", i+1)
		}
	}
}

func TestOutbox_RetriesThenDelivers(t *testing.T) {
	f := &flakySender{failures: 2, attempts: make(chan struct{}, 10)}
	o, logs := startOutbox(t, f, []time.Duration{0, time.Millisecond, time.Millisecond})

	msg := &Message{To: "ana@exemple.org", Subject: "Mot de passe", Body: "lien secret"}
	require.True(t, o.Enqueue(job(msg)))
	waitAttempts(t, f, 3)

	require.Eventually(t, func() bool { return strings.Contains(logs.String(), "remis au relais") }, time.Second, 5*time.Millisecond)
	f.mu.Lock()
	assert.Len(t, f.sent, 1)
	f.mu.Unlock()
	out := logs.String()
	assert.Contains(t, out, `"user_id":"usr_ana"`)
	assert.Contains(t, out, `"kind":"password_reset"`)
	assert.Contains(t, out, "<destinataire>", "la réponse du relais est gardée, sans l'adresse")
	assert.NotContains(t, out, "ana@exemple.org")
	assert.NotContains(t, out, "lien secret")
}

func TestOutbox_GivesUpAfterAllAttempts(t *testing.T) {
	f := &flakySender{failures: 99, attempts: make(chan struct{}, 10)}
	o, logs := startOutbox(t, f, []time.Duration{0, time.Millisecond, time.Millisecond})

	require.True(t, o.Enqueue(job(&Message{To: "ana@exemple.org"})))
	waitAttempts(t, f, 3)
	require.Eventually(t, func() bool {
		return strings.Contains(logs.String(), "abandonné après tous les essais")
	}, time.Second, 5*time.Millisecond)
	assert.Empty(t, f.attempts, "trois essais, pas un de plus")
}

func TestOutbox_NothingToSend(t *testing.T) {
	f := &flakySender{attempts: make(chan struct{}, 10)}
	o, _ := startOutbox(t, f, nil)
	built := make(chan struct{})
	require.True(t, o.Enqueue(Job{Kind: "password_reset", Build: func(context.Context) (*Message, string, error) {
		close(built)
		return nil, "", nil
	}}))
	<-built
	time.Sleep(20 * time.Millisecond)
	assert.Empty(t, f.attempts)
}

func TestOutbox_ShutdownAbandonsPendingRetries(t *testing.T) {
	f := &flakySender{failures: 99, attempts: make(chan struct{}, 10)}
	logs := &syncBuffer{}
	o := NewOutbox(f, []time.Duration{0, time.Hour}, zerolog.New(logs))
	go func() { _ = o.Start() }()
	<-o.started

	require.True(t, o.Enqueue(job(&Message{To: "ana@exemple.org"})))
	waitAttempts(t, f, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	require.NoError(t, o.Shutdown(ctx), "l'attente d'un nouvel essai ne retient pas l'arrêt")
	assert.Contains(t, logs.String(), "message abandonné")
	assert.False(t, o.Enqueue(job(&Message{To: "ana@exemple.org"})), "file arrêtée")
}

func TestOutbox_ShutdownWithoutStart(t *testing.T) {
	o := NewOutbox(&flakySender{}, nil, zerolog.Nop())
	require.NoError(t, o.Shutdown(context.Background()))
}
