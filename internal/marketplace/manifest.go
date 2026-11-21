// Package marketplace lit et vérifie les templates du Marketplace : un dossier
// de templates/ avec un manifest.yaml, un docker-compose.yml, une icône et un
// README (ADR-0018).
//
// Le type Manifest est le format qui fait foi : templates/README.md le décrit
// pour les contributeurs, et la vérification le lit en mode strict.
package marketplace

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"
)

// SchemaVersion est la version du format que ce code lit.
const SchemaVersion = 1

// Fichiers d'un template.
const (
	ManifestFile = "manifest.yaml"
	ComposeFile  = "docker-compose.yml"
	IconFile     = "icon.svg"
	ReadmeFile   = "README.md"
)

// Categories est la liste fermée des catégories du catalogue. Une catégorie
// s'ajoute par une PR qui modifie cette liste.
var Categories = []string{
	"database", "cache", "storage", "cms", "devtools", "monitoring", "communication", "other",
}

// Protocoles d'un port.
const (
	ProtocolHTTP = "http" // peut recevoir un domaine servi par Caddy (ADR-0008)
	ProtocolTCP  = "tcp"  // interne : sert d'autres applications du serveur
)

// Longueurs bornées des valeurs générées et des textes affichés.
const (
	minGenerate    = 16
	maxGenerate    = 128
	maxName        = 40
	maxDescription = 160
	maxVersion     = 64
)

var (
	// templateID reprend la règle des noms de projet Compose, en plus court :
	// l'identifiant nommera la stack installée.
	templateID   = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,30}$`)
	variableName = regexp.MustCompile(`^[A-Z][A-Z0-9_]{0,63}$`)
	volumeName   = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)
	serviceName  = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)
	spdxID       = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9.+-]*$`)
)

// reservedPrefixes changeraient le comportement de la commande docker compose
// elle-même (projet, démon visé) ou de ForgeNet : les valeurs des variables
// sont ajoutées à son environnement.
var reservedPrefixes = []string{"FORGENET_", "COMPOSE_", "DOCKER_"}

// reservedNames sont transmis par ForgeNet à la commande (ADR-0003).
var reservedNames = []string{"PATH", "HOME"}

// Manifest décrit un template : ce que le dashboard affiche et demande.
type Manifest struct {
	Schema      int        `yaml:"schema"`
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Category    string     `yaml:"category"`
	Version     string     `yaml:"version"`
	Service     string     `yaml:"service"`
	License     string     `yaml:"license"`
	Homepage    string     `yaml:"homepage"`
	Variables   []Variable `yaml:"variables"`
	Ports       []Port     `yaml:"ports"`
	Volumes     []Volume   `yaml:"volumes"`
}

// Variable est une valeur demandée à l'installation.
type Variable struct {
	Name     string `yaml:"name"`
	Label    string `yaml:"label"`
	Help     string `yaml:"help"`
	Default  string `yaml:"default"`
	Pattern  string `yaml:"pattern"`
	Optional bool   `yaml:"optional"`
	Secret   bool   `yaml:"secret"`
	Generate int    `yaml:"generate"`
}

// Port est un port où un service écoute dans son conteneur.
type Port struct {
	Service     string `yaml:"service"`
	Port        int    `yaml:"port"`
	Protocol    string `yaml:"protocol"`
	Description string `yaml:"description"`
}

// Volume est un volume nommé de la stack.
type Volume struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// Problem est un écart à une règle, rattaché au fichier qui le porte.
type Problem struct {
	File    string
	Message string
}

func (p Problem) Error() string { return p.File + " : " + p.Message }

func problem(file, format string, args ...any) error {
	return Problem{File: file, Message: fmt.Sprintf(format, args...)}
}

// ParseManifest lit un manifeste. Un champ inconnu est une erreur : une faute
// de frappe (`generat: 32`) ne passe pas en silence.
func ParseManifest(data []byte) (Manifest, error) {
	var m Manifest
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&m); err != nil {
		if errors.Is(err, io.EOF) {
			return Manifest{}, problem(ManifestFile, "fichier vide")
		}
		return Manifest{}, problem(ManifestFile, "%v", err)
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return Manifest{}, problem(ManifestFile, "un seul document YAML attendu")
	}
	return m, nil
}

// Validate vérifie les règles propres au manifeste, sans le fichier Compose.
func (m Manifest) Validate() []error {
	var errs []error
	add := func(format string, args ...any) { errs = append(errs, problem(ManifestFile, format, args...)) }

	if m.Schema != SchemaVersion {
		add("schema : %d attendu, %d trouvé", SchemaVersion, m.Schema)
	}
	checkText(add, "name", m.Name, maxName)
	checkText(add, "description", m.Description, maxDescription)
	if !slices.Contains(Categories, m.Category) {
		add("category : %q n'est pas dans la liste (%s)", m.Category, strings.Join(Categories, ", "))
	}
	switch {
	case m.Version == "":
		add("version : obligatoire")
	case len(m.Version) > maxVersion || strings.ContainsAny(m.Version, " \t\n"):
		add("version : %q doit être un mot de %d caractères au plus", m.Version, maxVersion)
	}
	if !serviceName.MatchString(m.Service) {
		add("service : le nom du service principal est obligatoire")
	}
	if !validLicense(m.License) {
		add("license : %q n'est pas un identifiant SPDX (MIT, AGPL-3.0-only, PostgreSQL…)", m.License)
	}
	if u, err := url.Parse(m.Homepage); err != nil || u.Scheme != "https" || u.Host == "" {
		add("homepage : %q doit être une adresse en https://", m.Homepage)
	}

	errs = append(errs, m.validateVariables()...)
	errs = append(errs, m.validatePorts()...)
	errs = append(errs, m.validateVolumes()...)
	return errs
}

func checkText(add func(string, ...any), field, value string, limit int) {
	switch {
	case strings.TrimSpace(value) == "":
		add("%s : obligatoire", field)
	case utf8.RuneCountInString(value) > limit:
		add("%s : %d caractères au plus", field, limit)
	case strings.ContainsAny(value, "\n\r"):
		add("%s : une seule ligne", field)
	}
}

// validLicense accepte un identifiant SPDX ou une expression simple
// (`RSALv2 OR AGPL-3.0-only`).
func validLicense(expr string) bool {
	tokens := strings.Fields(strings.NewReplacer("(", " ", ")", " ").Replace(expr))
	if len(tokens) == 0 {
		return false
	}
	for i, tok := range tokens {
		operator := tok == "AND" || tok == "OR" || tok == "WITH"
		if operator == (i%2 == 0) || (!operator && !spdxID.MatchString(tok)) {
			return false
		}
	}
	return len(tokens)%2 == 1
}

func (m Manifest) validateVariables() []error {
	var errs []error
	seen := map[string]bool{}
	for i, v := range m.Variables {
		where := fmt.Sprintf("variables[%d] (%s)", i, v.Name)
		add := func(format string, args ...any) {
			errs = append(errs, problem(ManifestFile, where+" : "+format, args...))
		}
		switch {
		case !variableName.MatchString(v.Name):
			add("le nom s'écrit en majuscules, chiffres et _, 64 caractères au plus")
		case reservedName(v.Name):
			add("nom réservé : les préfixes %s et les noms %s changeraient la commande docker compose elle-même",
				strings.Join(reservedPrefixes, ", "), strings.Join(reservedNames, ", "))
		case seen[v.Name]:
			add("déclarée deux fois")
		}
		seen[v.Name] = true

		if strings.TrimSpace(v.Label) == "" {
			add("label : obligatoire, c'est le libellé du formulaire d'installation")
		}
		if v.Secret && v.Default != "" {
			add("un secret n'a pas de valeur par défaut : elle serait écrite en dur dans le dépôt ; utilisez generate")
		}
		if v.Generate != 0 {
			if !v.Secret {
				add("generate est réservé aux variables secret: true")
			} else if v.Generate < minGenerate || v.Generate > maxGenerate {
				add("generate : une longueur de %d à %d caractères", minGenerate, maxGenerate)
			}
		}
		if v.Pattern != "" {
			re, err := regexp.Compile(v.Pattern)
			if err != nil {
				add("pattern : expression régulière invalide : %v", err)
			} else if v.Default != "" && !re.MatchString(v.Default) {
				add("default : %q ne respecte pas pattern", v.Default)
			}
		}
	}
	return errs
}

func reservedName(name string) bool {
	if slices.Contains(reservedNames, name) {
		return true
	}
	for _, p := range reservedPrefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

func (m Manifest) validatePorts() []error {
	var errs []error
	seen := map[string]bool{}
	hasTCP := false
	for i, p := range m.Ports {
		where := fmt.Sprintf("ports[%d] (%s:%d)", i, p.Service, p.Port)
		add := func(format string, args ...any) {
			errs = append(errs, problem(ManifestFile, where+" : "+format, args...))
		}
		if !serviceName.MatchString(p.Service) {
			add("service : obligatoire")
		}
		if p.Port < 1 || p.Port > 65535 {
			add("port : de 1 à 65535")
		}
		switch p.Protocol {
		case ProtocolHTTP:
		case ProtocolTCP:
			hasTCP = true
		default:
			add("protocol : %q ou %q", ProtocolHTTP, ProtocolTCP)
		}
		if strings.TrimSpace(p.Description) == "" {
			add("description : obligatoire")
		}
		key := fmt.Sprintf("%s:%d", p.Service, p.Port)
		if seen[key] {
			add("déclaré deux fois")
		}
		seen[key] = true
	}

	// Un service tcp est joignable par les applications reliées (ADR-0018,
	// revue de #120) : il ne s'installe jamais sans authentification.
	if hasTCP && !slices.ContainsFunc(m.Variables, func(v Variable) bool { return v.Secret }) {
		errs = append(errs, problem(ManifestFile,
			"ports : un port tcp exige au moins une variable secret: true, pour que le service demande une authentification"))
	}
	return errs
}

func (m Manifest) validateVolumes() []error {
	var errs []error
	seen := map[string]bool{}
	for i, v := range m.Volumes {
		where := fmt.Sprintf("volumes[%d] (%s)", i, v.Name)
		add := func(format string, args ...any) {
			errs = append(errs, problem(ManifestFile, where+" : "+format, args...))
		}
		if !volumeName.MatchString(v.Name) {
			add("name : lettres, chiffres, _, . et -")
		}
		if strings.TrimSpace(v.Description) == "" {
			add("description : obligatoire, elle dit ce que le volume garde")
		}
		if seen[v.Name] {
			add("déclaré deux fois")
		}
		seen[v.Name] = true
	}
	return errs
}

// ValidID indique si le nom d'un dossier peut identifier un template. Un
// dossier en « _ » (le template d'exemple) suit la même règle après le « _ ».
func ValidID(dir string) bool {
	return templateID.MatchString(strings.TrimPrefix(dir, "_"))
}
