package marketplace_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/marketplace"
)

const validIcon = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"><defs><linearGradient id="g"/></defs>` +
	`<rect width="64" height="64" fill="url(#g)"/><use href="#g"/></svg>`

func TestCheckIcon(t *testing.T) {
	assert.Empty(t, joined(marketplace.CheckIcon([]byte(validIcon))))

	cases := map[string]string{
		"vide":                  "",
		"pas un SVG":            `<html></html>`,
		"pas carré":             `<svg viewBox="0 0 64 32"></svg>`,
		"sans viewBox":          `<svg></svg>`,
		"script":                `<svg viewBox="0 0 1 1"><script>alert(1)</script></svg>`,
		"gestionnaire":          `<svg viewBox="0 0 1 1" onload="alert(1)"></svg>`,
		"lien externe":          `<svg viewBox="0 0 1 1"><image href="https://exemple.org/x.png"/></svg>`,
		"foreignObject":         `<svg viewBox="0 0 1 1"><foreignObject/></svg>`,
		"entités":               `<!DOCTYPE svg [<!ENTITY a "b">]><svg viewBox="0 0 1 1"></svg>`,
		"illisible":             `<svg viewBox="0 0 1 1">`,
		"trop lourd (> 16 Kio)": `<svg viewBox="0 0 1 1">` + strings.Repeat(" ", 16<<10) + `</svg>`,
	}
	for name, icon := range cases {
		assert.NotEmpty(t, marketplace.CheckIcon([]byte(icon)), name)
	}
}

// resolver renvoie toujours le modèle valide : la lecture par Compose est
// essayée à part (templates/templates_test.go).
func resolver(_ context.Context, _ []byte) ([]byte, error) { return []byte(validResolved), nil }

func templateFS() fstest.MapFS {
	return fstest.MapFS{
		"postgres/manifest.yaml":      {Data: []byte(validManifest)},
		"postgres/docker-compose.yml": {Data: []byte("services: {}\n")},
		"postgres/icon.svg":           {Data: []byte(validIcon)},
		"postgres/README.md":          {Data: []byte("# PostgreSQL\n")},
	}
}

func TestCheck_ValidFolder(t *testing.T) {
	assert.Empty(t, joined(marketplace.Check(context.Background(), templateFS(), "postgres", resolver)))
}

func TestCheck_FolderRules(t *testing.T) {
	ctx := context.Background()

	fsys := templateFS()
	fsys["postgres/.env"] = &fstest.MapFile{Data: []byte("POSTGRES_PASSWORD=x\n")}
	assert.Contains(t, joined(marketplace.Check(ctx, fsys, "postgres", resolver)), ".env : fichier en trop")

	fsys = templateFS()
	delete(fsys, "postgres/icon.svg")
	assert.Contains(t, joined(marketplace.Check(ctx, fsys, "postgres", resolver)), "icon.svg : fichier obligatoire absent")

	fsys = templateFS()
	fsys["postgres/README.md"] = &fstest.MapFile{Data: []byte("  \n")}
	assert.Contains(t, joined(marketplace.Check(ctx, fsys, "postgres", resolver)), "README.md : vide")

	fsys = fstest.MapFS{}
	for name, f := range templateFS() {
		fsys["Postgres"+strings.TrimPrefix(name, "postgres")] = f
	}
	assert.Contains(t, joined(marketplace.Check(ctx, fsys, "Postgres", resolver)), "nom de dossier")

	refused := func(context.Context, []byte) ([]byte, error) {
		return nil, errors.New("services.postgres.image must be a string")
	}
	assert.Contains(t, joined(marketplace.Check(ctx, templateFS(), "postgres", refused)),
		"docker-compose.yml : refusé par docker compose : services.postgres.image must be a string")
}

func TestLoad_SkipsExampleAndRejectsInvalidManifests(t *testing.T) {
	fsys := templateFS()
	for name, f := range templateFS() {
		fsys["_exemple"+strings.TrimPrefix(name, "postgres")] = f
	}
	list, err := marketplace.Load(fsys)
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, "postgres", list[0].ID)
	assert.Equal(t, "PostgreSQL", list[0].Manifest.Name)

	fsys["postgres/manifest.yaml"] = &fstest.MapFile{Data: []byte(strings.Replace(validManifest, "category: database", "category: bdd", 1))}
	_, err = marketplace.Load(fsys)
	assert.ErrorContains(t, err, "templates/postgres : manifest.yaml : category")
}
