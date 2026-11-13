package proxy

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "réécrit les Caddyfile de référence de testdata/")

func golden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if *update {
		require.NoError(t, os.WriteFile(path, got, 0o600))
	}
	want, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, string(want), string(got))
}

var productionSettings = Settings{
	AdminAddress: "unix//run/forgenet/caddy-admin.sock|0220",
	Email:        "ops@example.org",
}

func TestRender_Production(t *testing.T) {
	routes := []Route{
		{Host: "www.example.org", Upstream: Upstream{Container: "vitrine-web-1", Port: 80}},
		{Host: "api.example.org", Upstream: Upstream{Container: "boutique-api-1", Port: 8080}},
	}
	got, err := Render(productionSettings, routes)
	require.NoError(t, err)
	golden(t, "production.caddyfile", got)

	// Les routes de l'appelant ne sont pas réordonnées.
	assert.Equal(t, "www.example.org", routes[0].Host)
}

func TestRender_LocalCertificates(t *testing.T) {
	s := Settings{AdminAddress: "127.0.0.1:2019", LocalCerts: true, HTTPPort: 8080, HTTPSPort: 8443}
	got, err := Render(s, []Route{{Host: "app.forgenet.test", Upstream: Upstream{Container: "app", Port: 3000}}})
	require.NoError(t, err)
	golden(t, "local.caddyfile", got)
}

func TestRender_NoRoutes(t *testing.T) {
	got, err := Render(Settings{AdminAddress: "unix//run/forgenet/caddy-admin.sock"}, nil)
	require.NoError(t, err)
	golden(t, "empty.caddyfile", got)
}

func TestRender_SameRoutesSameFile(t *testing.T) {
	a := []Route{
		{Host: "b.example.org", Upstream: Upstream{Container: "b", Port: 80}},
		{Host: "a.example.org", Upstream: Upstream{Container: "a", Port: 80}},
	}
	b := []Route{a[1], a[0]}
	first, err := Render(productionSettings, a)
	require.NoError(t, err)
	second, err := Render(productionSettings, b)
	require.NoError(t, err)
	assert.Equal(t, first, second)
}

func TestRender_Rejects(t *testing.T) {
	ok := Upstream{Container: "app", Port: 80}
	cases := []struct {
		name     string
		settings Settings
		routes   []Route
		want     error
	}{
		{"domaine en double", productionSettings,
			[]Route{{Host: "a.example.org", Upstream: ok}, {Host: "a.example.org", Upstream: Upstream{Container: "b", Port: 80}}},
			ErrDuplicateHost},
		{"domaine invalide", productionSettings, []Route{{Host: "a.org {\n}\n:80", Upstream: ok}}, ErrInvalidHost},
		{"cible invalide", productionSettings, []Route{{Host: "a.example.org", Upstream: Upstream{Container: "app\nrespond", Port: 80}}}, ErrInvalidUpstream},
		{"administration exposée", Settings{AdminAddress: "0.0.0.0:2019"}, nil, ErrInvalidAdminAddress},
		{"administration absente", Settings{}, nil, ErrInvalidAdminAddress},
		{"courriel invalide", Settings{AdminAddress: "127.0.0.1:2019", Email: "a@b.io\n}"}, nil, ErrInvalidEmail},
		{"port invalide", Settings{AdminAddress: "127.0.0.1:2019", HTTPSPort: 70000}, nil, ErrInvalidPort},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Render(c.settings, c.routes)
			assert.ErrorIs(t, err, c.want)
			assert.Nil(t, got)
		})
	}
}

// FuzzRender vérifie qu'aucune entrée acceptée ne change la structure du
// Caddyfile : un bloc global, puis exactement un site de trois lignes par route.
func FuzzRender(f *testing.F) {
	f.Add("app.example.org", "app-web-1", 8080, "ops@example.org")
	f.Add("a.org {\n}\nb.org", "x}", 80, "a`b@c.io")
	f.Add("xn--bcher-kva.example", "a.b_c-d", 1, "")
	f.Add("*.example.org", "app", 65535, "a{$HOME}@b.io")
	f.Fuzz(func(t *testing.T, host, container string, port int, email string) {
		s := Settings{AdminAddress: "unix//run/forgenet/caddy-admin.sock", Email: email}
		out, err := Render(s, []Route{{Host: host, Upstream: Upstream{Container: container, Port: port}}})
		if err != nil {
			return
		}
		text := string(out)
		assert.Equal(t, 2, strings.Count(text, "{"), text)
		assert.Equal(t, 2, strings.Count(text, "}"), text)
		assert.NotContains(t, text, "`")
		assert.NotContains(t, text, `"`)
		assert.Contains(t, text, "\n"+host+" {\n\treverse_proxy "+container+":")
		assert.True(t, strings.HasSuffix(text, "\n}\n"))

		lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
		globals := 3 // en-tête, « { », admin
		if email != "" {
			globals++
		}
		// bloc global (+ « } »), ligne vide, site de trois lignes
		assert.Len(t, lines, globals+1+1+3, text)
	})
}
