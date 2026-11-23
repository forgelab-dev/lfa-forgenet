package deploy

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/forgelab-dev/lfa-forgenet/internal/sealed"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

var (
	// ErrExists indique qu'une installation vise un nom d'application déjà
	// pris : elle remplacerait la stack en place.
	ErrExists = errors.New("une application porte déjà ce nom")
	// ErrNoKey indique des valeurs à chiffrer ou à ouvrir sur un serveur sans
	// clé des secrets.
	ErrNoKey = errors.New("valeurs d'installation impossibles à chiffrer : FORGENET_SECRET_KEY n'est pas réglée sur ce serveur")
)

// variablesPurpose lie le chiffré à son usage (sealed.Box) : un autre
// secret de la base, copié dans la colonne, ne s'ouvrirait pas.
const variablesPurpose = "deployment.variables"

// mask remplace une valeur secrète dans le journal et les messages d'échec.
const mask = "••••••"

// Variable est la valeur d'une variable d'une installation (ADR-0019).
type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	// Secret : la valeur est masquée dans le journal et les messages.
	Secret bool `json:"secret,omitempty"`
}

// WithBox donne la clé qui chiffre les valeurs des installations (ADR-0015).
// Sans elle, un déploiement avec des variables est refusé (ErrNoKey).
func WithBox(box *sealed.Box) Option { return func(s *Service) { s.box = box } }

// sealVariables chiffre toutes les valeurs ensemble : c'est plus simple qu'un
// mélange de clair et de chiffré, et rien ne les réaffiche.
func (s *Service) sealVariables(vars []Variable) ([]byte, error) {
	if s.box == nil {
		return nil, ErrNoKey
	}
	plain, err := json.Marshal(vars)
	if err != nil {
		return nil, fmt.Errorf("valeurs d'installation : %w", err)
	}
	return s.box.Seal(variablesPurpose, plain)
}

// openVariables rend les valeurs d'un déploiement ; aucune sans chiffré.
func (s *Service) openVariables(d *store.Deployment) ([]Variable, error) {
	if d.VariablesSealed == nil {
		return nil, nil
	}
	if s.box == nil {
		return nil, ErrNoKey
	}
	plain, err := s.box.Open(variablesPurpose, d.VariablesSealed)
	if err != nil {
		return nil, fmt.Errorf("valeurs d'installation : %w", err)
	}
	var vars []Variable
	if err := json.Unmarshal(plain, &vars); err != nil {
		return nil, fmt.Errorf("valeurs d'installation illisibles : %w", err)
	}
	return vars, nil
}

// environment écrit les valeurs pour la commande docker compose.
func environment(vars []Variable) []string {
	if len(vars) == 0 {
		return nil
	}
	env := make([]string, 0, len(vars))
	for _, v := range vars {
		env = append(env, v.Name+"="+v.Value)
	}
	return env
}

// variableNames liste les noms, pour l'audit : jamais les valeurs.
func variableNames(vars []Variable) []string {
	names := make([]string, 0, len(vars))
	for _, v := range vars {
		names = append(names, v.Name)
	}
	return names
}

// redactor renvoie une fonction qui masque les valeurs secrètes de vars. Les
// plus longues passent d'abord : une valeur qui en contient une autre est
// masquée entière.
func redactor(vars []Variable) func(string) string {
	var secrets []string
	for _, v := range vars {
		if v.Secret && v.Value != "" {
			secrets = append(secrets, v.Value)
		}
	}
	if len(secrets) == 0 {
		return func(s string) string { return s }
	}
	sort.Slice(secrets, func(i, j int) bool { return len(secrets[i]) > len(secrets[j]) })
	pairs := make([]string, 0, 2*len(secrets))
	for _, secret := range secrets {
		pairs = append(pairs, secret, mask)
	}
	r := strings.NewReplacer(pairs...)
	return r.Replace
}
