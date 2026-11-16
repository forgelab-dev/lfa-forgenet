package compose

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
)

// DefaultRestart est la politique donnée à un service qui n'en précise aucune :
// Docker le relance au démarrage du serveur, sauf s'il a été arrêté exprès.
const DefaultRestart = "unless-stopped"

// model reprend du modèle résolu (`config --format json`) ce qui décide de la
// politique de redémarrage. Compose y a déjà appliqué les ancres YAML, `extends`
// et les profils : un service hérité d'une ancre garde sa politique, un service
// d'un profil inactif n'apparaît pas.
type model struct {
	Services map[string]struct {
		Restart string `json:"restart"`
		Deploy  *struct {
			RestartPolicy json.RawMessage `json:"restart_policy"`
		} `json:"deploy"`
		DependsOn map[string]struct {
			Condition string `json:"condition"`
		} `json:"depends_on"`
	} `json:"services"`
}

// servicesWithoutRestart liste, triés, les services à qui ForgeNet donne
// DefaultRestart. Sont laissés tels quels :
//   - un service avec `restart:` (`no` compris) ;
//   - un service avec `deploy.restart_policy`, que Compose applique quand
//     `restart:` manque ;
//   - une tâche ponctuelle, qu'un autre service attend avec la condition
//     `service_completed_successfully` (une migration, par exemple) : relancée
//     sans fin, elle tournerait en boucle.
func servicesWithoutRestart(config []byte) ([]string, error) {
	var m model
	if err := json.Unmarshal(config, &m); err != nil {
		return nil, fmt.Errorf("lecture du modèle Compose : %w", err)
	}

	oneShot := map[string]bool{}
	for _, s := range m.Services {
		for dep, d := range s.DependsOn {
			if d.Condition == "service_completed_successfully" {
				oneShot[dep] = true
			}
		}
	}

	var names []string
	for name, s := range m.Services {
		explicit := s.Restart != "" || (s.Deploy != nil && len(s.Deploy.RestartPolicy) > 0 && string(s.Deploy.RestartPolicy) != "null")
		if !explicit && !oneShot[name] {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	return names, nil
}

// restartOverride écrit le fichier complémentaire qui donne DefaultRestart à
// services, et renvoie son chemin. Le fichier de l'utilisateur n'est jamais
// réécrit : Compose fusionne les deux, et ce qui est enregistré et affiché
// reste ce que l'utilisateur a envoyé. Le fichier est rangé hors du dossier de
// travail, qui reste vide ; l'appelant le supprime après usage.
func (e *Engine) restartOverride(services []string) (string, error) {
	override := map[string]map[string]map[string]string{"services": {}}
	for _, name := range services {
		override["services"][name] = map[string]string{"restart": DefaultRestart}
	}
	// Du JSON est du YAML valide, et l'encodeur échappe les noms de service.
	content, err := json.Marshal(override)
	if err != nil {
		return "", fmt.Errorf("politique de redémarrage : %w", err)
	}

	f, err := os.CreateTemp(e.overrideDir, "restart-*.json")
	if err != nil {
		return "", fmt.Errorf("politique de redémarrage : %w", err)
	}
	path := f.Name()
	_, err = f.Write(content)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("politique de redémarrage : %w", err)
	}
	return path, nil
}

// restartNotice est la ligne du journal qui dit à quels services la politique
// a été ajoutée.
func restartNotice(services []string) string {
	return fmt.Sprintf("Relance après un redémarrage du serveur (restart: %s) ajoutée à : %s",
		DefaultRestart, strings.Join(services, ", "))
}
