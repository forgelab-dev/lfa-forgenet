package update

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type countingSource struct {
	calls   atomic.Int32
	version string
	err     error
}

func (s *countingSource) Latest(context.Context, bool) (Release, error) {
	s.calls.Add(1)
	if s.err != nil {
		return Release{}, s.err
	}
	return Release{Version: s.version}, nil
}

func TestCheckerAnnouncesOnlyNewerVersions(t *testing.T) {
	for current, want := range map[string]bool{"0.1.1": true, "0.2.0": false, "0.3.0": false, "dev": false} {
		c := NewChecker(&countingSource{version: "0.2.0"}, current, 0, zerolog.Nop())
		s := c.Check(context.Background())
		assert.Equal(t, want, s.Available != nil, "version en service %s", current)
		assert.False(t, s.CheckedAt.IsZero())
	}
}

func TestCheckerThrottlesAndKeepsLastResult(t *testing.T) {
	src := &countingSource{version: "0.2.0"}
	c := NewChecker(src, "0.1.1", 0, zerolog.Nop())
	clock := time.Date(2026, 11, 14, 15, 0, 0, 0, time.UTC)
	c.now = func() time.Time { return clock }

	c.Check(context.Background())
	c.Check(context.Background())
	assert.Equal(t, int32(1), src.calls.Load(), "deux clics rapprochés : un seul appel à GitHub")

	clock = clock.Add(2 * time.Minute)
	src.err = errors.New("GitHub injoignable")
	s := c.Check(context.Background())
	assert.Equal(t, int32(2), src.calls.Load())
	assert.Equal(t, "GitHub injoignable", s.Err)
	require.NotNil(t, s.Available, "une panne passagère n'efface pas la version connue")
	assert.Equal(t, "0.2.0", s.Available.Version)
}

func TestCheckerStartAndShutdown(t *testing.T) {
	src := &countingSource{version: "0.2.0"}
	c := NewChecker(src, "0.1.1", time.Hour, zerolog.Nop())
	c.firstDelay = time.Millisecond
	errc := make(chan error, 1)
	go func() { errc <- c.Start() }()
	require.Eventually(t, func() bool { return src.calls.Load() == 1 }, time.Second, 5*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	require.NoError(t, c.Shutdown(ctx))
	require.NoError(t, <-errc)
	assert.NotNil(t, c.Snapshot().Available)
}
