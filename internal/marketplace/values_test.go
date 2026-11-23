package marketplace_test

import (
	"crypto/rand"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/marketplace"
)

func resolveManifest() marketplace.Manifest {
	return marketplace.Manifest{Variables: []marketplace.Variable{
		{Name: "POSTGRES_USER", Label: "Utilisateur", Default: "app", Pattern: "^[a-z_][a-z0-9_]{0,62}$"},
		{Name: "POSTGRES_PASSWORD", Label: "Mot de passe", Secret: true, Generate: 24},
		{Name: "API_KEY", Label: "Clé", Secret: true},
		{Name: "NOTE", Label: "Note", Optional: true},
	}}
}

func TestResolve_DefaultsAndGenerated(t *testing.T) {
	values, err := resolveManifest().Resolve(map[string]string{"API_KEY": "cle-saisie"}, rand.Reader)
	require.NoError(t, err)
	require.Len(t, values, 4)

	// L'ordre est celui du manifeste.
	names := []string{values[0].Name, values[1].Name, values[2].Name, values[3].Name}
	assert.Equal(t, []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "API_KEY", "NOTE"}, names)

	assert.Equal(t, marketplace.Value{Name: "POSTGRES_USER", Value: "app"}, values[0])

	password := values[1]
	assert.True(t, password.Secret)
	assert.True(t, password.Generated)
	// Sans « $ » : Compose le lirait comme le début d'une variable.
	assert.Regexp(t, regexp.MustCompile(`^[A-Za-z0-9]{24}$`), password.Value)

	assert.Equal(t, marketplace.Value{Name: "API_KEY", Value: "cle-saisie", Secret: true}, values[2])
	assert.Equal(t, marketplace.Value{Name: "NOTE"}, values[3])
}

func TestResolve_GeneratedValuesDiffer(t *testing.T) {
	m := resolveManifest()
	first, err := m.Resolve(map[string]string{"API_KEY": "x"}, rand.Reader)
	require.NoError(t, err)
	second, err := m.Resolve(map[string]string{"API_KEY": "x"}, rand.Reader)
	require.NoError(t, err)
	assert.NotEqual(t, first[1].Value, second[1].Value)
}

func TestResolve_KeepsTypedSecret(t *testing.T) {
	values, err := resolveManifest().Resolve(map[string]string{"POSTGRES_PASSWORD": "choisi", "API_KEY": "x"}, rand.Reader)
	require.NoError(t, err)
	assert.Equal(t, marketplace.Value{Name: "POSTGRES_PASSWORD", Value: "choisi", Secret: true}, values[1])
}

func TestResolve_Refusals(t *testing.T) {
	cases := map[string]struct {
		input map[string]string
		want  string
	}{
		"variable inconnue":     {map[string]string{"API_KEY": "x", "AUTRE": "y"}, "AUTRE : variable inconnue"},
		"obligatoire manquante": {map[string]string{}, "API_KEY : valeur obligatoire"},
		"motif":                 {map[string]string{"API_KEY": "x", "POSTGRES_USER": "Pas Bon"}, "POSTGRES_USER : la valeur ne respecte pas le motif"},
		"caractère nul":         {map[string]string{"API_KEY": "a\x00b"}, "API_KEY : caractère nul"},
		"trop longue":           {map[string]string{"API_KEY": strings.Repeat("a", 4097)}, "API_KEY : 4096 octets au plus"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := resolveManifest().Resolve(tc.input, rand.Reader)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}

func TestResolve_ErrorNeverQuotesSecret(t *testing.T) {
	m := marketplace.Manifest{Variables: []marketplace.Variable{
		{Name: "CODE", Label: "Code", Secret: true, Pattern: "^[0-9]{4}$"},
	}}
	_, err := m.Resolve(map[string]string{"CODE": "secret-tres-prive"}, rand.Reader)
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "secret-tres-prive")
}

func TestResolve_RandomFailure(t *testing.T) {
	_, err := resolveManifest().Resolve(map[string]string{"API_KEY": "x"}, strings.NewReader(""))
	require.Error(t, err)
}
