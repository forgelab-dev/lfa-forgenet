package grpc

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/forgelab-dev/lfa-forgenet/internal/mail"
)

// SyncOutboxForTests sert aux tests du paquet grpc_test.
type SyncOutboxForTests = syncOutbox

// syncOutbox prépare chaque envoi tout de suite et garde les messages.
type syncOutbox struct {
	mu   sync.Mutex
	sent []mail.Message
	full bool
}

func (o *syncOutbox) Enqueue(job mail.Job) bool {
	if o.full {
		return false
	}
	msg, _, err := job.Build(context.Background())
	if err == nil && msg != nil {
		o.mu.Lock()
		o.sent = append(o.sent, *msg)
		o.mu.Unlock()
	}
	return true
}

// SetFull fait refuser les envois suivants, comme une file pleine.
func (o *syncOutbox) SetFull(full bool) { o.full = full }

func (o *syncOutbox) Messages() []mail.Message {
	o.mu.Lock()
	defer o.mu.Unlock()
	return append([]mail.Message(nil), o.sent...)
}

func TestIsPublicAddress(t *testing.T) {
	for host, want := range map[string]bool{
		"203.0.113.7": true, "2001:db8::1": true, "8.8.8.8": true,
		"127.0.0.1": false, "::1": false, "172.18.0.3": false, "10.0.0.2": false,
		"192.168.1.4": false, "fe80::1": false, "": false, "pas-une-ip": false,
		"::ffff:127.0.0.1": false,
	} {
		assert.Equal(t, want, isPublicAddress(host), host)
	}
}
