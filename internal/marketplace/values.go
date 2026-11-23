package marketplace

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"slices"
	"sort"
	"strings"
)

// ErrInvalidValue indique une valeur d'installation refusée. Le message nomme
// la variable, jamais la valeur : elle peut être un secret.
var ErrInvalidValue = errors.New("valeur d'installation refusée")

// maxValueSize borne une valeur saisie : assez pour une clé ou un certificat,
// pas pour remplir l'environnement de la commande.
const maxValueSize = 4096

// generatedAlphabet compose les secrets tirés au hasard : lettres et chiffres
// seulement. Un « $ » serait lu par Compose comme le début d'une variable, et
// une ponctuation casserait plus d'une chaîne de connexion.
const generatedAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// Value est la valeur d'une variable pour une installation.
type Value struct {
	Name  string
	Value string
	// Secret reprend la variable du manifeste : la valeur n'est ni journalisée
	// ni réaffichée.
	Secret bool
	// Generated vaut vrai pour un secret tiré au hasard par le serveur.
	Generated bool
}

// Resolve vérifie les valeurs saisies et complète les autres : la valeur par
// défaut, ou un secret tiré au hasard dans random (crypto/rand.Reader hors des
// tests). Une valeur vide vaut une valeur absente. Le résultat suit l'ordre
// du manifeste.
func (m Manifest) Resolve(input map[string]string, random io.Reader) ([]Value, error) {
	var unknown []string
	for name := range input {
		if !slices.ContainsFunc(m.Variables, func(v Variable) bool { return v.Name == name }) {
			unknown = append(unknown, name)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return nil, fmt.Errorf("%w : %s : variable inconnue de ce template", ErrInvalidValue, strings.Join(unknown, ", "))
	}

	values := make([]Value, 0, len(m.Variables))
	for _, v := range m.Variables {
		value, err := v.resolve(input[v.Name], random)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

func (v Variable) resolve(typed string, random io.Reader) (Value, error) {
	out := Value{Name: v.Name, Value: typed, Secret: v.Secret}
	refuse := func(reason string) (Value, error) {
		return Value{}, fmt.Errorf("%w : %s : %s", ErrInvalidValue, v.Name, reason)
	}

	if typed == "" {
		switch {
		case v.Default != "":
			out.Value = v.Default
			return out, nil
		case v.Generate > 0:
			generated, err := randomString(random, v.Generate)
			if err != nil {
				return Value{}, fmt.Errorf("tirage de %s : %w", v.Name, err)
			}
			out.Value, out.Generated = generated, true
			return out, nil
		case v.Optional:
			return out, nil
		default:
			return refuse("valeur obligatoire")
		}
	}

	if len(typed) > maxValueSize {
		return refuse(fmt.Sprintf("%d octets au plus", maxValueSize))
	}
	// Une variable d'environnement ne peut pas contenir de NUL : la commande
	// ne démarrerait pas.
	if strings.ContainsRune(typed, 0) {
		return refuse("caractère nul interdit")
	}
	// Le motif est vérifié par Manifest.Validate : il compile.
	if v.Pattern != "" && !regexp.MustCompile(v.Pattern).MatchString(typed) {
		return refuse("la valeur ne respecte pas le motif " + v.Pattern)
	}
	return out, nil
}

// randomString tire n caractères de generatedAlphabet, sans biais : un octet
// au-delà du dernier multiple de la taille de l'alphabet est rejeté.
func randomString(random io.Reader, n int) (string, error) {
	const size = len(generatedAlphabet)
	limit := byte(256 - 256%size)
	out := make([]byte, 0, n)
	buf := make([]byte, n)
	for len(out) < n {
		if _, err := io.ReadFull(random, buf); err != nil {
			return "", err
		}
		for _, b := range buf {
			if b < limit && len(out) < n {
				out = append(out, generatedAlphabet[int(b)%size])
			}
		}
	}
	return string(out), nil
}
