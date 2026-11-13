package proxy

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// received est une requête reçue par le faux Caddy.
type received struct {
	method, path, contentType, host, body string
}

// fakeCaddy sert l'API d'administration sur un socket Unix et répond par respond.
func fakeCaddy(t *testing.T, respond func(w http.ResponseWriter)) (address string, got <-chan received) {
	t.Helper()
	dir, err := os.MkdirTemp("", "fnadm")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	sock := filepath.Join(dir, "admin.sock")
	ln, err := net.Listen("unix", sock)
	require.NoError(t, err)

	ch := make(chan received, 1)
	srv := &http.Server{
		ReadHeaderTimeout: time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			ch <- received{r.Method, r.URL.Path, r.Header.Get("Content-Type"), r.Host, string(body)}
			respond(w)
		}),
	}
	go func() { _ = srv.Serve(ln) }()
	t.Cleanup(func() { _ = srv.Close() })
	return "unix/" + sock, ch
}

func TestAdminClient_LoadSendsCaddyfile(t *testing.T) {
	addr, got := fakeCaddy(t, func(w http.ResponseWriter) { w.WriteHeader(http.StatusOK) })
	client, err := NewAdminClient(addr)
	require.NoError(t, err)

	warnings, err := client.Load(context.Background(), []byte("a.example.org {\n}\n"))
	require.NoError(t, err)
	assert.Empty(t, warnings)

	r := <-got
	assert.Equal(t, http.MethodPost, r.method)
	assert.Equal(t, "/load", r.path)
	assert.Equal(t, "text/caddyfile", r.contentType)
	assert.Equal(t, "localhost", r.host, "Caddy n'accepte que l'origine localhost sur un socket")
	assert.Equal(t, "a.example.org {\n}\n", r.body)
}

func TestAdminClient_LoadReturnsWarnings(t *testing.T) {
	addr, _ := fakeCaddy(t, func(w http.ResponseWriter) {
		_, _ = io.WriteString(w, `{"warnings":[{"file":"Caddyfile","line":3,"directive":"reverse_proxy","message":"attention"}]}`)
	})
	client, err := NewAdminClient(addr)
	require.NoError(t, err)

	warnings, err := client.Load(context.Background(), []byte("x"))
	require.NoError(t, err)
	require.Len(t, warnings, 1)
	assert.Equal(t, "Caddyfile:3 : attention", warnings[0].String())
}

func TestAdminClient_LoadRefused(t *testing.T) {
	addr, _ := fakeCaddy(t, func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `{"error":"loading config: Caddyfile:2 - Error during parsing","warnings":[{"file":"Caddyfile","line":2,"message":"m"}]}`)
	})
	client, err := NewAdminClient(addr)
	require.NoError(t, err)

	_, err = client.Load(context.Background(), []byte("x"))
	var loadErr *LoadError
	require.ErrorAs(t, err, &loadErr)
	assert.Equal(t, "loading config: Caddyfile:2 - Error during parsing", loadErr.Message)
	assert.Len(t, loadErr.Warnings, 1)
}

func TestAdminClient_UnexpectedResponseIsBounded(t *testing.T) {
	addr, _ := fakeCaddy(t, func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, strings.Repeat("é", 2<<20))
	})
	client, err := NewAdminClient(addr)
	require.NoError(t, err)

	_, err = client.Load(context.Background(), []byte("x"))
	require.Error(t, err)
	var loadErr *LoadError
	assert.NotErrorAs(t, err, &loadErr)
	assert.Contains(t, err.Error(), "500")
	assert.Less(t, len(err.Error()), 700)
	assert.True(t, strings.HasSuffix(err.Error(), "é…"), "coupé entre deux caractères")
}

func TestAdminClient_LoopbackIgnoresEnvironmentProxy(t *testing.T) {
	// Un proxy de l'environnement ne doit jamais voir le trafic d'administration.
	t.Setenv("HTTP_PROXY", "http://127.0.0.1:1")
	t.Setenv("NO_PROXY", "")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/load", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client, err := NewAdminClient(strings.TrimPrefix(srv.URL, "http://"))
	require.NoError(t, err)
	_, err = client.Load(context.Background(), []byte("x"))
	require.NoError(t, err)
}

func TestAdminClient_Unreachable(t *testing.T) {
	client, err := NewAdminClient("unix//tmp/forgenet-absent-" + time.Now().Format("150405.000000") + ".sock")
	require.NoError(t, err)
	_, err = client.Load(context.Background(), []byte("x"))
	assert.ErrorIs(t, err, ErrUnavailable)
}

func TestAdminClient_RespectsContext(t *testing.T) {
	addr, _ := fakeCaddy(t, func(http.ResponseWriter) { time.Sleep(2 * time.Second) })
	client, err := NewAdminClient(addr)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err = client.Load(ctx, []byte("x"))
	require.Error(t, err)
	assert.Less(t, time.Since(start), time.Second)
}

func TestNewAdminClient_RefusesExposedAddress(t *testing.T) {
	for _, addr := range []string{"0.0.0.0:2019", ":2019", "10.0.0.5:2019"} {
		_, err := NewAdminClient(addr)
		assert.ErrorIs(t, err, ErrInvalidAdminAddress, addr)
	}
}
