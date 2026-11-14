// Package update met à jour une installation de ForgeNet (ADR-0011) : il lit
// la fiche de la dernière release, sauvegarde la base et la configuration,
// bascule le serveur sur l'image exacte de la fiche, attend sa santé et revient
// en arrière en cas d'échec.
package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// OfficialImage est la seule image qu'une mise à jour accepte d'installer :
// une fiche qui en désigne une autre est refusée.
const OfficialImage = "ghcr.io/forgelab-dev/forgenet-server"

// ReleaseFile est le nom de la fiche jointe à chaque release GitHub.
const ReleaseFile = "forgenet-release.json"

var (
	digestRe = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	sha256Re = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

// InstallFiles sont les fichiers d'installation qu'une release peut faire
// évoluer. Seuls ces noms sont acceptés dans une fiche : une fiche ne peut pas
// faire écrire ailleurs dans le dossier d'installation.
var InstallFiles = []string{"docker-compose.yml", "Caddyfile.base"}

// Release est la fiche d'une version publiée.
type Release struct {
	Schema   int    `json:"schema"`
	Version  string `json:"version"`
	Image    string `json:"image"`
	Digest   string `json:"digest"`
	NotesURL string `json:"notes_url"`
	// Files donne l'empreinte SHA-256 (hexadécimale) des fichiers
	// d'installation de la release. Absent des fiches d'avant la 0.2.0 : les
	// fichiers en place sont alors gardés.
	Files map[string]string `json:"files,omitempty"`
}

// ImageRef désigne l'image par version et par empreinte : c'est l'empreinte
// qui fait foi, l'étiquette reste lisible pour l'administrateur.
func (r Release) ImageRef() string {
	return fmt.Sprintf("%s:%s@%s", r.Image, r.Version, r.Digest)
}

// ErrInvalidRelease signale une fiche inutilisable.
var ErrInvalidRelease = errors.New("fiche de release invalide")

// ParseRelease lit et valide une fiche. Les champs inconnus sont tolérés (une
// version future peut en ajouter) ; les champs connus sont vérifiés un à un.
func ParseRelease(data []byte) (Release, error) {
	var r Release
	if err := json.Unmarshal(data, &r); err != nil {
		return Release{}, fmt.Errorf("%w : %v", ErrInvalidRelease, err)
	}
	switch {
	case r.Schema != 1:
		return Release{}, fmt.Errorf("%w : schéma %d inconnu", ErrInvalidRelease, r.Schema)
	case r.Image != OfficialImage:
		return Release{}, fmt.Errorf("%w : image %q, seule %s est acceptée", ErrInvalidRelease, r.Image, OfficialImage)
	case !digestRe.MatchString(r.Digest):
		return Release{}, fmt.Errorf("%w : empreinte %q", ErrInvalidRelease, r.Digest)
	}
	if _, err := ParseVersion(r.Version); err != nil {
		return Release{}, fmt.Errorf("%w : %v", ErrInvalidRelease, err)
	}
	for name, sum := range r.Files {
		if !slices.Contains(InstallFiles, name) {
			return Release{}, fmt.Errorf("%w : fichier %q inconnu", ErrInvalidRelease, name)
		}
		if !sha256Re.MatchString(sum) {
			return Release{}, fmt.Errorf("%w : empreinte de %s %q", ErrInvalidRelease, name, sum)
		}
	}
	if !strings.HasPrefix(r.NotesURL, "https://github.com/forgelab-dev/lfa-forgenet/releases/") {
		r.NotesURL = ""
	}
	return r, nil
}

// Version est un numéro sémantique MAJEUR.MINEUR.CORRECTIF, avec une
// pré-version éventuelle (1.2.0-rc.1).
type Version struct {
	Major, Minor, Patch int
	Pre                 string
}

func (v Version) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Pre != "" {
		s += "-" + v.Pre
	}
	return s
}

// ParseVersion lit « 1.2.3 », « v1.2.3 » ou « 1.2.3-rc.1 ». Les suffixes de
// développement (« 0.1.0-alpha.dev ») sont des pré-versions comme les autres.
func ParseVersion(s string) (Version, error) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "v")
	core, pre, _ := strings.Cut(s, "-")
	parts := strings.Split(core, ".")
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("version %q : trois nombres attendus", s)
	}
	var n [3]int
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil || v < 0 || (len(p) > 1 && p[0] == '0') {
			return Version{}, fmt.Errorf("version %q : %q n'est pas un nombre", s, p)
		}
		n[i] = v
	}
	return Version{Major: n[0], Minor: n[1], Patch: n[2], Pre: pre}, nil
}

// Compare renvoie -1, 0 ou 1. Une pré-version passe avant la version finale
// (1.2.0-rc.1 < 1.2.0) ; deux pré-versions se comparent comme des textes.
func (v Version) Compare(o Version) int {
	for _, d := range [3]int{v.Major - o.Major, v.Minor - o.Minor, v.Patch - o.Patch} {
		if d != 0 {
			return sign(d)
		}
	}
	switch {
	case v.Pre == o.Pre:
		return 0
	case v.Pre == "":
		return 1
	case o.Pre == "":
		return -1
	}
	return sign(strings.Compare(v.Pre, o.Pre))
}

func sign(n int) int {
	switch {
	case n < 0:
		return -1
	case n > 0:
		return 1
	}
	return 0
}
