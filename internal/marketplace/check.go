package marketplace

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strings"
)

// ResolveFunc renvoie le modèle résolu d'un fichier Compose, sans remplacer
// les ${NOM} (compose.Engine.Resolve).
type ResolveFunc func(ctx context.Context, content []byte) ([]byte, error)

// Template est un template du catalogue.
type Template struct {
	ID       string
	Manifest Manifest
}

// files sont les seuls fichiers d'un dossier de template : le fichier Compose
// ne lit jamais ses voisins (ADR-0003), un autre fichier ne servirait à rien,
// et un .env laissé là par erreur finirait dans le dépôt.
var files = []string{ManifestFile, ComposeFile, IconFile, ReadmeFile}

// Dirs liste les dossiers de templates de fsys, ceux en « _ » compris.
func Dirs(fsys fs.FS) ([]string, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			dirs = append(dirs, e.Name())
		}
	}
	return dirs, nil
}

// Load lit le catalogue : les dossiers de fsys, sauf ceux en « _ », dont le
// manifeste doit être valide. La vérification complète, qui demande Docker,
// est faite par la CI (Check).
func Load(fsys fs.FS) ([]Template, error) {
	dirs, err := Dirs(fsys)
	if err != nil {
		return nil, err
	}
	var out []Template
	for _, dir := range dirs {
		if strings.HasPrefix(dir, "_") {
			continue
		}
		if !ValidID(dir) {
			return nil, fmt.Errorf("templates/%s : nom de dossier invalide", dir)
		}
		data, err := fs.ReadFile(fsys, path.Join(dir, ManifestFile))
		if err != nil {
			return nil, fmt.Errorf("templates/%s : %w", dir, err)
		}
		m, err := ParseManifest(data)
		if err == nil {
			err = errors.Join(m.Validate()...)
		}
		if err != nil {
			return nil, fmt.Errorf("templates/%s : %w", dir, err)
		}
		out = append(out, Template{ID: dir, Manifest: m})
	}
	return out, nil
}

// Check vérifie un dossier de template entier : manifeste, fichier Compose
// croisé avec le manifeste, icône, README et absence d'autres fichiers.
// Chaque erreur est un Problem qui nomme le fichier en cause.
func Check(ctx context.Context, fsys fs.FS, dir string, resolve ResolveFunc) []error {
	var errs []error
	if !ValidID(dir) {
		errs = append(errs, problem(".", "nom de dossier %q : minuscules, chiffres et -, 31 caractères au plus", dir))
	}

	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return append(errs, problem(".", "%v", err))
	}
	for _, e := range entries {
		if !slices.Contains(files, e.Name()) {
			errs = append(errs, problem(e.Name(), "fichier en trop : un template contient %s seulement", strings.Join(files, ", ")))
		}
	}
	read := func(name string) []byte {
		data, err := fs.ReadFile(fsys, path.Join(dir, name))
		if err != nil {
			errs = append(errs, problem(name, "fichier obligatoire absent"))
			return nil
		}
		return data
	}

	manifestData, composeData := read(ManifestFile), read(ComposeFile)
	if icon := read(IconFile); icon != nil {
		errs = append(errs, CheckIcon(icon)...)
	}
	if readme := read(ReadmeFile); readme != nil && strings.TrimSpace(string(readme)) == "" {
		errs = append(errs, problem(ReadmeFile, "vide : décrivez le service, ses variables et comment s'y connecter"))
	}

	if manifestData == nil {
		return errs
	}
	m, err := ParseManifest(manifestData)
	if err != nil {
		return append(errs, err)
	}
	errs = append(errs, m.Validate()...)

	if composeData == nil {
		return errs
	}
	resolved, err := resolve(ctx, composeData)
	if err != nil {
		return append(errs, problem(ComposeFile, "refusé par docker compose : %v", err))
	}
	return append(errs, CheckCompose(m, resolved)...)
}
