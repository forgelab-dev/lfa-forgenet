package proxy

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateHost(t *testing.T) {
	valid := []string{
		"example.org",
		"app.example.org",
		"a-b.c-d.io",
		"forgenet.test",
		"xn--bcher-kva.example",
		"x1.example.org",
		strings.Repeat("a", 63) + ".org",
	}
	for _, h := range valid {
		assert.NoError(t, ValidateHost(h), h)
	}

	invalid := []string{
		"",
		"localhost",           // pas de domaine parent
		"Example.org",         // majuscules : le domaine doit être normalisé avant
		"example.org.",        // point final
		"*.example.org",       // joker
		"1.2.3.4",             // adresse IP
		"[::1]",               // adresse IP
		"app.123",             // dernière étiquette numérique
		"exa mple.org",        // espace
		"-app.example.org",    // tiret au bord
		"app-.example.org",    // tiret au bord
		"app..example.org",    // étiquette vide
		"app_web.example.org", // tiret bas
		"ex{ample}.org",       // accolades du Caddyfile
		"a.org\nb.org",        // retour à la ligne
		"a.org b.org",         // deux sites
		"a.org#commentaire",   // commentaire du Caddyfile
		"bücher.example",      // nom internationalisé non converti
		"a.org:443",           // port
		"http://a.org",        // schéma
		strings.Repeat("a", 64) + ".org",
		strings.Repeat("a.", 126) + "org", // 255 caractères
	}
	for _, h := range invalid {
		assert.ErrorIs(t, ValidateHost(h), ErrInvalidHost, "%q", h)
	}
}

func TestUpstream_Validate(t *testing.T) {
	for _, u := range []Upstream{
		{Container: "app-web-1", Port: 8080},
		{Container: "a", Port: 1},
		{Container: "Stack_web.1", Port: 65535},
	} {
		assert.NoError(t, u.Validate(), "%+v", u)
	}

	for _, u := range []Upstream{
		{Container: "", Port: 80},
		{Container: "-app", Port: 80},
		{Container: "app web", Port: 80},
		{Container: "app}", Port: 80},
		{Container: "app:80", Port: 80},
		{Container: "app\nfile_server", Port: 80},
		{Container: strings.Repeat("a", 129), Port: 80},
		{Container: "app", Port: 0},
		{Container: "app", Port: -1},
		{Container: "app", Port: 65536},
	} {
		assert.ErrorIs(t, u.Validate(), ErrInvalidUpstream, "%+v", u)
	}
}

func TestValidateEmail(t *testing.T) {
	for _, e := range []string{"ops@example.org", "a.b+forgenet@x.io"} {
		assert.NoError(t, validateEmail(e), e)
	}
	for _, e := range []string{
		"",
		"ops",
		"Ops <ops@example.org>",
		"ops@example.org }",
		"a{b}@x.io",
		"a`b@x.io",
		"a @x.io",
		"a@x.io\n",
		`"a b"@x.io`,
		"a#b@x.io",
		strings.Repeat("a", 250) + "@x.io",
	} {
		assert.ErrorIs(t, validateEmail(e), ErrInvalidEmail, "%q", e)
	}
}

func TestParseAdminAddress(t *testing.T) {
	cases := []struct {
		in   string
		want adminEndpoint
	}{
		{"unix//run/forgenet/caddy-admin.sock", adminEndpoint{"unix", "/run/forgenet/caddy-admin.sock", "localhost"}},
		{"unix//run/forgenet/caddy-admin.sock|0200", adminEndpoint{"unix", "/run/forgenet/caddy-admin.sock", "localhost"}},
		{"unix//tmp/a.sock|660", adminEndpoint{"unix", "/tmp/a.sock", "localhost"}},
		{"127.0.0.1:2019", adminEndpoint{"tcp", "127.0.0.1:2019", "127.0.0.1:2019"}},
		{"[::1]:2019", adminEndpoint{"tcp", "[::1]:2019", "[::1]:2019"}},
		{"localhost:2019", adminEndpoint{"tcp", "localhost:2019", "localhost:2019"}},
	}
	for _, c := range cases {
		got, err := parseAdminAddress(c.in)
		require.NoError(t, err, c.in)
		assert.Equal(t, c.want, got, c.in)
	}

	for _, in := range []string{
		"",
		":2019",              // toutes les interfaces
		"0.0.0.0:2019",       // toutes les interfaces
		"[::]:2019",          // toutes les interfaces
		"192.168.1.10:2019",  // réseau local
		"example.org:2019",   // nom résolu ailleurs
		"127.0.0.1",          // sans port
		"127.0.0.1:0",        // port nul
		"127.0.0.1:99999",    // port hors bornes
		"localhost:http",     // port nommé
		"unix//",             // chemin vide
		"unix/relative.sock", // chemin relatif
		"unix///run/../etc/x.sock",
		"unix//run/x sock",
		"unix//run/x.sock|abc",
		"unix//run/x.sock|0999",
		"unix//run/x.sock\n}",
		"unix//" + strings.Repeat("a", 110),
		"http://127.0.0.1:2019",
	} {
		_, err := parseAdminAddress(in)
		assert.ErrorIs(t, err, ErrInvalidAdminAddress, "%q", in)
	}
}
