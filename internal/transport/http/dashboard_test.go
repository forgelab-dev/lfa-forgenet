package http_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	transportHTTP "github.com/forgelab-dev/lfa-forgenet/internal/transport/http"
)

const indexHTML = "<!doctype html><title>ForgeNet</title>"

// dashboardDir imite un build Vite, avec à côté un fichier secret hors du
// dossier et des pièges dedans : fichiers cachés et liens qui en sortent.
func dashboardDir(t *testing.T) string {
	t.Helper()
	base := t.TempDir()
	dir := filepath.Join(base, "dist")
	files := map[string]string{
		"index.html":                    indexHTML,
		"favicon.svg":                   "<svg/>",
		"assets/index-abc123.js":        "console.log(1)",
		"assets/font-abc123.woff2":      "wOF2",
		".env":                          "SECRET=1",
		"assets/.git/config":            "[core]",
		"../secret.txt":                 "hors du dossier",
		"assets/sub/index-def456.css":   "body{}",
		"assets/.hidden-in-assets.json": "{}",
	}
	for name, content := range files {
		p := filepath.Join(dir, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o755))
		require.NoError(t, os.WriteFile(p, []byte(content), 0o600))
	}
	require.NoError(t, os.Symlink(filepath.Join(base, "secret.txt"), filepath.Join(dir, "evasion.txt")))
	require.NoError(t, os.Symlink(base, filepath.Join(dir, "parent")))
	return dir
}

func serveDashboard(t *testing.T, dir string) *http.ServeMux {
	t.Helper()
	h, err := transportHTTP.NewDashboardHandler(dir)
	require.NoError(t, err)
	t.Cleanup(func() { _ = h.Close() })
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) { _, _ = io.WriteString(w, "ok") })
	h.RegisterRoutes(mux)
	return mux
}

func get(mux http.Handler, method, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(method, target, nil))
	return rec
}

func TestDashboard_ServesFilesAndPages(t *testing.T) {
	mux := serveDashboard(t, dashboardDir(t))

	cases := []struct {
		target, body, cache string
	}{
		{"/", indexHTML, "no-cache"},
		{"/domaines", indexHTML, "no-cache"},
		{"/conteneurs/abc/journal", indexHTML, "no-cache"},
		{"/assets", indexHTML, "no-cache"},
		{"/favicon.svg", "<svg/>", "no-cache"},
		{"/assets/index-abc123.js", "console.log(1)", "public, max-age=31536000, immutable"},
		{"/assets/sub/index-def456.css", "body{}", "public, max-age=31536000, immutable"},
	}
	for _, c := range cases {
		rec := get(mux, http.MethodGet, c.target)
		assert.Equal(t, http.StatusOK, rec.Code, c.target)
		assert.Equal(t, c.body, rec.Body.String(), c.target)
		assert.Equal(t, c.cache, rec.Header().Get("Cache-Control"), c.target)
	}

	assert.Equal(t, "font/woff2", get(mux, http.MethodGet, "/assets/font-abc123.woff2").Header().Get("Content-Type"))
	assert.Equal(t, "ok", get(mux, http.MethodGet, "/healthz").Body.String(), "les routes plus précises gardent la priorité")
	head := get(mux, http.MethodHead, "/")
	assert.Equal(t, http.StatusOK, head.Code)
	assert.Empty(t, head.Body.String())
}

func TestDashboard_NeverLeavesItsFolder(t *testing.T) {
	mux := serveDashboard(t, dashboardDir(t))

	for _, target := range []string{
		"/.env", "/assets/.git/config", "/assets/.hidden-in-assets.json",
		"/evasion.txt", "/parent/secret.txt", "/assets/manquant.js", "/index.html/x.js",
	} {
		rec := get(mux, http.MethodGet, target)
		assert.Equal(t, http.StatusNotFound, rec.Code, target)
		assert.NotContains(t, rec.Body.String(), "hors du dossier", target)
		assert.NotContains(t, rec.Body.String(), "SECRET", target)
	}

	// Les « .. » d'une adresse brute sont résolus avant toute lecture : au
	// pire, la page du dashboard.
	for _, raw := range []string{"/../secret.txt", "/assets/../../secret.txt", "/%2e%2e/secret.txt"} {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.URL.Path = strings.ReplaceAll(raw, "%2e", ".")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		assert.NotContains(t, rec.Body.String(), "hors du dossier", raw)
	}
}

func TestDashboard_RefusesWrites(t *testing.T) {
	mux := serveDashboard(t, dashboardDir(t))
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		rec := get(mux, method, "/")
		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, method)
		assert.Equal(t, "GET, HEAD", rec.Header().Get("Allow"), method)
	}
}

func TestDashboard_SecurityHeaders(t *testing.T) {
	mux := serveDashboard(t, dashboardDir(t))
	for _, target := range []string{"/", "/assets/index-abc123.js", "/.env"} {
		h := get(mux, http.MethodGet, target).Header()
		csp := h.Get("Content-Security-Policy")
		assert.Contains(t, csp, "default-src 'self'", target)
		assert.Contains(t, csp, "frame-ancestors 'none'", target)
		assert.NotContains(t, csp, "unsafe", target)
		assert.Equal(t, "nosniff", h.Get("X-Content-Type-Options"), target)
		assert.Equal(t, "DENY", h.Get("X-Frame-Options"), target)
		assert.Equal(t, "no-referrer", h.Get("Referrer-Policy"), target)
	}
}

func TestNewDashboardHandler_RequiresABuild(t *testing.T) {
	_, err := transportHTTP.NewDashboardHandler(t.TempDir())
	require.ErrorContains(t, err, "index.html")

	_, err = transportHTTP.NewDashboardHandler(filepath.Join(t.TempDir(), "absent"))
	require.Error(t, err)
}
