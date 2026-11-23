package marketplace_test

import (
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/marketplace"
	"github.com/forgelab-dev/lfa-forgenet/templates"
)

// catalogFS contient deux templates du catalogue et un brouillon écarté.
func catalogFS() fstest.MapFS {
	fsys := fstest.MapFS{}
	for _, dir := range []string{"postgres", "zz-copie", "_brouillon"} {
		for name, f := range templateFS() {
			fsys[dir+strings.TrimPrefix(name, "postgres")] = f
		}
	}
	return fsys
}

func TestCatalog(t *testing.T) {
	c, err := marketplace.NewCatalog(catalogFS())
	require.NoError(t, err)

	list := c.Templates()
	require.Len(t, list, 2)
	assert.Equal(t, "postgres", list[0].ID)
	assert.Equal(t, "zz-copie", list[1].ID)
	list[0].ID = "modifié"
	assert.Equal(t, "postgres", c.Templates()[0].ID, "la liste rendue est une copie")

	tpl, ok := c.Template("postgres")
	require.True(t, ok)
	assert.Equal(t, "PostgreSQL", tpl.Manifest.Name)
	_, ok = c.Template("_brouillon")
	assert.False(t, ok, "un dossier en « _ » n'est pas au catalogue")
	_, ok = c.Template("inconnu")
	assert.False(t, ok)

	readme, err := c.Readme("postgres")
	require.NoError(t, err)
	assert.Equal(t, "# PostgreSQL\n", string(readme))
	for _, id := range []string{"inconnu", "_brouillon", "../postgres", ""} {
		_, err := c.Readme(id)
		assert.ErrorIs(t, err, fs.ErrNotExist, id)
	}
}

func TestCatalog_RejectsInvalidTemplate(t *testing.T) {
	fsys := catalogFS()
	delete(fsys, "zz-copie/manifest.yaml")
	_, err := marketplace.NewCatalog(fsys)
	assert.ErrorContains(t, err, "templates/zz-copie")
}

// Le catalogue embarqué dans le binaire se charge : sinon le serveur ne
// démarre pas.
func TestCatalog_Embedded(t *testing.T) {
	_, err := marketplace.NewCatalog(templates.FS)
	require.NoError(t, err)
}
