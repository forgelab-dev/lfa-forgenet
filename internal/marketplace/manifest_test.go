package marketplace_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/marketplace"
)

// validManifest est l'exemple de l'ADR-0018.
const validManifest = `schema: 1
name: PostgreSQL
description: Base de données relationnelle, la plus répandue des bases libres.
category: database
version: "17.2"
service: postgres
license: PostgreSQL
homepage: https://www.postgresql.org

variables:
  - name: POSTGRES_USER
    label: Utilisateur créé à l'installation
    default: app
    pattern: "^[a-z_][a-z0-9_]{0,62}$"
  - name: POSTGRES_DB
    label: Base créée à l'installation
    default: app
    pattern: "^[a-z_][a-z0-9_]{0,62}$"
  - name: POSTGRES_PASSWORD
    label: Mot de passe de l'utilisateur
    secret: true
    generate: 32

ports:
  - service: postgres
    port: 5432
    protocol: tcp
    description: Connexions des autres applications

volumes:
  - name: data
    description: Fichiers de la base
`

func parse(t *testing.T, content string) marketplace.Manifest {
	t.Helper()
	m, err := marketplace.ParseManifest([]byte(content))
	require.NoError(t, err)
	return m
}

func joined(errs []error) string {
	var b strings.Builder
	for _, e := range errs {
		b.WriteString(e.Error())
		b.WriteString("\n")
	}
	return b.String()
}

func TestManifest_ADRExampleIsValid(t *testing.T) {
	m := parse(t, validManifest)
	assert.Empty(t, joined(m.Validate()))
	assert.Equal(t, "17.2", m.Version)
	assert.Equal(t, 32, m.Variables[2].Generate)
}

func TestParseManifest_RejectsUnknownFieldsAndExtraDocuments(t *testing.T) {
	_, err := marketplace.ParseManifest([]byte(strings.Replace(validManifest, "generate: 32", "generat: 32", 1)))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "generat", "une faute de frappe ne passe pas en silence")

	_, err = marketplace.ParseManifest([]byte(validManifest + "---\nschema: 1\n"))
	assert.ErrorContains(t, err, "un seul document")

	_, err = marketplace.ParseManifest(nil)
	assert.ErrorContains(t, err, "vide")
}

func TestManifest_Rules(t *testing.T) {
	cases := []struct {
		name, from, to, want string
	}{
		{"schéma inconnu", "schema: 1", "schema: 2", "schema : 1 attendu"},
		{"nom absent", "name: PostgreSQL", "name: ''", "name : obligatoire"},
		{"description trop longue", "description: Base de données relationnelle, la plus répandue des bases libres.",
			"description: " + strings.Repeat("a", 161), "160 caractères"},
		{"catégorie hors liste", "category: database", "category: bdd", "n'est pas dans la liste"},
		{"version avec espace", `version: "17.2"`, `version: "17 2"`, "version"},
		{"licence non SPDX", "license: PostgreSQL", "license: libre !", "SPDX"},
		{"page en http", "homepage: https://www.postgresql.org", "homepage: http://www.postgresql.org", "https://"},
		{"nom de variable en minuscules", "name: POSTGRES_DB", "name: postgres_db", "majuscules"},
		{"variable réservée", "name: POSTGRES_DB", "name: COMPOSE_PROJECT_NAME", "nom réservé"},
		{"variable en double", "name: POSTGRES_DB", "name: POSTGRES_USER", "déclarée deux fois"},
		{"secret avec valeur par défaut", "    generate: 32", "    default: motdepasse", "pas de valeur par défaut"},
		{"generate sans secret", "    pattern: \"^[a-z_][a-z0-9_]{0,62}$\"\n  - name: POSTGRES_PASSWORD",
			"    generate: 32\n  - name: POSTGRES_PASSWORD", "réservé aux variables secret"},
		{"generate trop court", "generate: 32", "generate: 8", "16 à 128"},
		{"motif invalide", `default: app
    pattern: "^[a-z_][a-z0-9_]{0,62}$"
  - name: POSTGRES_DB`, `default: app
    pattern: "("
  - name: POSTGRES_DB`, "expression régulière invalide"},
		{"défaut hors motif", "default: app\n    pattern: \"^[a-z_][a-z0-9_]{0,62}$\"\n  - name: POSTGRES_DB",
			"default: App\n    pattern: \"^[a-z_][a-z0-9_]{0,62}$\"\n  - name: POSTGRES_DB", "ne respecte pas pattern"},
		{"protocole inconnu", "protocol: tcp", "protocol: udp", `"http" ou "tcp"`},
		{"port hors bornes", "port: 5432", "port: 70000", "de 1 à 65535"},
		{"volume sans description", "description: Fichiers de la base", "description: ''", "ce que le volume garde"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Contains(t, validManifest, tc.from, "le cas doit modifier l'exemple")
			m := parse(t, strings.Replace(validManifest, tc.from, tc.to, 1))
			assert.Contains(t, joined(m.Validate()), tc.want)
		})
	}
}

// Un service tcp est joignable par les applications reliées : il ne s'installe
// jamais sans authentification (revue de #120).
func TestManifest_TCPPortNeedsASecret(t *testing.T) {
	m := parse(t, validManifest)
	m.Variables = m.Variables[:2]
	assert.Contains(t, joined(m.Validate()), "un port tcp exige au moins une variable secret")

	m.Ports[0].Protocol = marketplace.ProtocolHTTP
	assert.NotContains(t, joined(m.Validate()), "port tcp")
}

func TestManifest_LicenseExpressions(t *testing.T) {
	for _, ok := range []string{"MIT", "AGPL-3.0-only", "RSALv2 OR SSPL-1.0 OR AGPL-3.0-only", "GPL-2.0-only WITH Classpath-exception-2.0", "(MIT OR Apache-2.0)"} {
		m := parse(t, validManifest)
		m.License = ok
		assert.NotContains(t, joined(m.Validate()), "license", ok)
	}
	for _, bad := range []string{"", "OR MIT", "MIT OR", "MIT Apache-2.0", "MIT AND AND Apache-2.0"} {
		m := parse(t, validManifest)
		m.License = bad
		assert.Contains(t, joined(m.Validate()), "license", bad)
	}
}

func TestValidID(t *testing.T) {
	for _, ok := range []string{"postgres", "minio", "uptime-kuma", "_exemple", "n8n"} {
		assert.True(t, marketplace.ValidID(ok), ok)
	}
	for _, bad := range []string{"", "Postgres", "-pg", "pg_16", "a" + strings.Repeat("b", 31), "_"} {
		assert.False(t, marketplace.ValidID(bad), bad)
	}
}
