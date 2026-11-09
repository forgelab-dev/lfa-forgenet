package grpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoginLimiterWindow(t *testing.T) {
	now := time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)
	l := newLoginLimiter(func() time.Time { return now })

	for range maxLoginFailures - 1 {
		l.fail("ip|membre")
	}
	assert.False(t, l.blocked("ip|membre"), "sous le seuil")
	l.fail("ip|membre")
	assert.True(t, l.blocked("ip|membre"), "au seuil")
	assert.False(t, l.blocked("ip|admin"), "clé indépendante")

	now = now.Add(loginWindow)
	assert.False(t, l.blocked("ip|membre"), "fenêtre échue")

	l.fail("ip|membre")
	l.reset("ip|membre")
	assert.Empty(t, l.failures, "une connexion réussie efface les échecs")
}
