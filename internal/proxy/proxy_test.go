package proxy

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeLoader enregistre les Caddyfile reçus et compte les appels simultanés.
type fakeLoader struct {
	mu       sync.Mutex
	loaded   [][]byte
	warnings []Warning
	err      error

	inFlight, maxInFlight atomic.Int32
}

func (f *fakeLoader) Load(_ context.Context, caddyfile []byte) ([]Warning, error) {
	n := f.inFlight.Add(1)
	defer f.inFlight.Add(-1)
	for {
		m := f.maxInFlight.Load()
		if n <= m || f.maxInFlight.CompareAndSwap(m, n) {
			break
		}
	}
	time.Sleep(5 * time.Millisecond)

	f.mu.Lock()
	defer f.mu.Unlock()
	f.loaded = append(f.loaded, caddyfile)
	return f.warnings, f.err
}

func (f *fakeLoader) calls() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.loaded)
}

var testSettings = Settings{AdminAddress: "unix//run/forgenet/caddy-admin.sock"}

var appRoute = Route{Host: "app.example.org", Upstream: Upstream{Container: "app-web-1", Port: 8080}}

func TestManager_ApplyLoadsAndKeepsCopy(t *testing.T) {
	loader := &fakeLoader{}
	copyPath := filepath.Join(t.TempDir(), "Caddyfile")
	m := NewManagerWithLoader(testSettings, loader, copyPath, zerolog.Nop())

	require.NoError(t, m.Apply(context.Background(), []Route{appRoute}))

	want, err := Render(testSettings, []Route{appRoute})
	require.NoError(t, err)
	require.Equal(t, 1, loader.calls())
	assert.Equal(t, want, loader.loaded[0])

	written, err := os.ReadFile(copyPath)
	require.NoError(t, err)
	assert.Equal(t, want, written)
	info, err := os.Stat(copyPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(copyMode), info.Mode().Perm())

	entries, err := os.ReadDir(filepath.Dir(copyPath))
	require.NoError(t, err)
	assert.Len(t, entries, 1, "aucun fichier temporaire ne reste")
}

func TestManager_InvalidRouteNeverReachesCaddy(t *testing.T) {
	loader := &fakeLoader{}
	copyPath := filepath.Join(t.TempDir(), "Caddyfile")
	m := NewManagerWithLoader(testSettings, loader, copyPath, zerolog.Nop())

	err := m.Apply(context.Background(), []Route{{Host: "a.org {\n}", Upstream: appRoute.Upstream}})
	assert.ErrorIs(t, err, ErrInvalidHost)
	assert.Zero(t, loader.calls())
	assert.NoFileExists(t, copyPath)
}

func TestManager_RefusedConfigKeepsPreviousCopy(t *testing.T) {
	loader := &fakeLoader{}
	copyPath := filepath.Join(t.TempDir(), "Caddyfile")
	m := NewManagerWithLoader(testSettings, loader, copyPath, zerolog.Nop())
	require.NoError(t, m.Apply(context.Background(), []Route{appRoute}))
	before, err := os.ReadFile(copyPath)
	require.NoError(t, err)

	loader.err = &LoadError{Message: "loading config: refusé"}
	err = m.Apply(context.Background(), nil)
	var loadErr *LoadError
	require.ErrorAs(t, err, &loadErr)

	after, err := os.ReadFile(copyPath)
	require.NoError(t, err)
	assert.Equal(t, before, after, "la copie reste celle de la configuration en service")
}

func TestManager_LogsWarnings(t *testing.T) {
	var logs bytes.Buffer
	loader := &fakeLoader{warnings: []Warning{{File: "Caddyfile", Line: 4, Message: "attention"}}}
	m := NewManagerWithLoader(testSettings, loader, "", zerolog.New(&logs))

	require.NoError(t, m.Apply(context.Background(), []Route{appRoute}))
	assert.Contains(t, logs.String(), "Caddyfile:4 : attention")
}

func TestManager_CopyFailureDoesNotFailApply(t *testing.T) {
	var logs bytes.Buffer
	loader := &fakeLoader{}
	copyPath := filepath.Join(t.TempDir(), "absent", "Caddyfile")
	m := NewManagerWithLoader(testSettings, loader, copyPath, zerolog.New(&logs))

	require.NoError(t, m.Apply(context.Background(), []Route{appRoute}))
	assert.Equal(t, 1, loader.calls())
	assert.Contains(t, logs.String(), "copie du Caddyfile non écrite")
}

func TestManager_OneReloadAtATime(t *testing.T) {
	loader := &fakeLoader{}
	m := NewManagerWithLoader(testSettings, loader, "", zerolog.Nop())

	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() { assert.NoError(t, m.Apply(context.Background(), []Route{appRoute})) })
	}
	wg.Wait()
	assert.Equal(t, 8, loader.calls())
	assert.Equal(t, int32(1), loader.maxInFlight.Load())
}

func TestNewManager_RefusesInvalidSettings(t *testing.T) {
	_, err := NewManager(Settings{AdminAddress: "0.0.0.0:2019"}, "", zerolog.Nop())
	assert.ErrorIs(t, err, ErrInvalidAdminAddress)

	_, err = NewManager(Settings{AdminAddress: "127.0.0.1:2019", Email: "pas une adresse"}, "", zerolog.Nop())
	assert.True(t, errors.Is(err, ErrInvalidEmail))
}
