package marketplace

import (
	"fmt"
	"io/fs"
	"path"
	"slices"
)

// Catalog est le catalogue que présente le dashboard : les templates de fsys,
// lus une fois au démarrage du serveur. Le serveur lui donne templates.FS,
// embarqué dans le binaire ; les tests, un catalogue à eux.
type Catalog struct {
	fsys      fs.FS
	templates []Template
}

// NewCatalog lit le catalogue de fsys (Load). Un template illisible est une
// erreur : le catalogue embarqué est vérifié par la CI, il ne doit jamais
// l'être à moitié.
func NewCatalog(fsys fs.FS) (*Catalog, error) {
	list, err := Load(fsys)
	if err != nil {
		return nil, err
	}
	return &Catalog{fsys: fsys, templates: list}, nil
}

// Templates renvoie une copie du catalogue, dans l'ordre des dossiers.
func (c *Catalog) Templates() []Template {
	return slices.Clone(c.templates)
}

// Template renvoie le template id, s'il est au catalogue.
func (c *Catalog) Template(id string) (Template, bool) {
	i := slices.IndexFunc(c.templates, func(t Template) bool { return t.ID == id })
	if i < 0 {
		return Template{}, false
	}
	return c.templates[i], true
}

// Readme renvoie le README.md du template id. Un id hors du catalogue répond
// fs.ErrNotExist : seul un template du catalogue est lu, jamais un chemin
// choisi par l'appelant.
func (c *Catalog) Readme(id string) ([]byte, error) {
	return c.file(id, ReadmeFile)
}

// Compose renvoie le docker-compose.yml du template id, avec la même garde
// que Readme : c'est ce fichier, et jamais un fichier reçu, qu'une
// installation déploie.
func (c *Catalog) Compose(id string) ([]byte, error) {
	return c.file(id, ComposeFile)
}

func (c *Catalog) file(id, name string) ([]byte, error) {
	if _, ok := c.Template(id); !ok {
		return nil, fmt.Errorf("template %q : %w", id, fs.ErrNotExist)
	}
	return fs.ReadFile(c.fsys, path.Join(id, name))
}
