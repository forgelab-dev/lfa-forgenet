package marketplace

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
)

// floatingTags désignent une version qui change sous le même nom.
var floatingTags = []string{"latest", "stable", "edge", "main", "master", "nightly"}

var (
	// secretLike repère une variable d'environnement dont la valeur est un
	// secret : elle doit venir d'une variable secret du manifeste.
	secretLike = regexp.MustCompile(`(?i)(PASS|SECRET|TOKEN|KEY|PRIVATE)`)
	// reference est la seule forme admise : ${NOM}.
	reference = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)
	// anyReference repère aussi les formes refusées ($NOM, ${NOM:-défaut}),
	// pour ne pas dire en plus que la variable ne sert nulle part.
	anyReference = regexp.MustCompile(`\$\{?([A-Za-z_][A-Za-z0-9_]*)`)
)

// model reprend du modèle résolu (`config --no-interpolate --format json`) ce
// que les règles de l'ADR-0018 regardent.
type model struct {
	Services map[string]service        `json:"services"`
	Volumes  map[string]external       `json:"volumes"`
	Networks map[string]external       `json:"networks"`
	Configs  map[string]map[string]any `json:"configs"`
	Secrets  map[string]map[string]any `json:"secrets"`
}

type external struct {
	External bool `json:"external"`
}

type service struct {
	Image         string            `json:"image"`
	Build         json.RawMessage   `json:"build"`
	ContainerName string            `json:"container_name"`
	Environment   json.RawMessage   `json:"environment"`
	EnvFile       json.RawMessage   `json:"env_file"`
	Ports         []json.RawMessage `json:"ports"`
	Volumes       []mount           `json:"volumes"`
	VolumesFrom   []string          `json:"volumes_from"`
	Privileged    bool              `json:"privileged"`
	CapAdd        []string          `json:"cap_add"`
	Devices       []json.RawMessage `json:"devices"`
	NetworkMode   string            `json:"network_mode"`
	Pid           string            `json:"pid"`
	Ipc           string            `json:"ipc"`
	Uts           string            `json:"uts"`
	UsernsMode    string            `json:"userns_mode"`
	SecurityOpt   []string          `json:"security_opt"`
	Healthcheck   *healthcheck      `json:"healthcheck"`
}

type mount struct {
	Type   string `json:"type"`
	Source string `json:"source"`
}

type healthcheck struct {
	Disable bool            `json:"disable"`
	Test    json.RawMessage `json:"test"`
}

// CheckCompose croise le fichier Compose d'un template, sous sa forme résolue,
// avec son manifeste.
func CheckCompose(m Manifest, resolved []byte) []error {
	var mod model
	if err := json.Unmarshal(resolved, &mod); err != nil {
		return []error{problem(ComposeFile, "modèle résolu illisible : %v", err)}
	}
	var errs []error
	add := func(format string, args ...any) { errs = append(errs, problem(ComposeFile, format, args...)) }

	if len(mod.Services) == 0 {
		add("aucun service")
	}
	used := map[string]bool{}
	for _, name := range sortedKeys(mod.Services) {
		s := mod.Services[name]
		errs = append(errs, checkService(name, s, m)...)
		for _, mt := range s.Volumes {
			if mt.Type == "volume" && mt.Source != "" {
				used[mt.Source] = true
			}
		}
	}

	errs = append(errs, checkMain(m, mod)...)
	errs = append(errs, checkVolumes(m, mod, used)...)
	for _, name := range sortedKeys(mod.Networks) {
		if mod.Networks[name].External {
			add("networks.%s : un réseau external rejoindrait un réseau qui n'est pas celui de la stack", name)
		}
	}
	for _, top := range []struct {
		kind    string
		entries map[string]map[string]any
	}{{"configs", mod.Configs}, {"secrets", mod.Secrets}} {
		for _, name := range sortedKeys(top.entries) {
			if _, ok := top.entries[name]["file"]; ok {
				add("%s.%s : file lirait un fichier du serveur ; utilisez content ou une variable", top.kind, name)
			}
		}
	}
	for _, p := range m.Ports {
		if _, ok := mod.Services[p.Service]; !ok {
			errs = append(errs, problem(ManifestFile, "ports : le service %q n'existe pas dans %s", p.Service, ComposeFile))
		}
	}

	errs = append(errs, checkVariables(m, resolved)...)
	return errs
}

func checkService(name string, s service, m Manifest) []error {
	var errs []error
	add := func(format string, args ...any) {
		errs = append(errs, problem(ComposeFile, "services."+name+" : "+format, args...))
	}

	if _, err := imageTag(s.Image); err != nil {
		add("%v", err)
	}
	if isSet(s.Build) {
		add("build est interdit : un template déploie une image publiée")
	}
	if s.ContainerName != "" {
		add("container_name est interdit : deux installations du même template se disputeraient le nom")
	}
	if isSet(s.EnvFile) {
		add("env_file est interdit : le fichier doit se suffire à lui-même (ADR-0003)")
	}
	if len(s.Ports) > 0 {
		add("ports est interdit : Docker publierait le port en contournant le pare-feu ; déclarez-le dans le manifeste (http ou tcp)")
	}
	for _, mt := range s.Volumes {
		switch {
		case mt.Type == "volume" && mt.Source == "":
			add("volume anonyme : ses données seraient perdues à la recréation du conteneur ; utilisez un volume nommé")
		case mt.Type == "bind":
			add("montage de l'hôte interdit (%s) : utilisez un volume nommé", mt.Source)
		case mt.Type != "volume" && mt.Type != "tmpfs":
			add("montage de type %q interdit", mt.Type)
		}
	}
	if len(s.VolumesFrom) > 0 {
		add("volumes_from est interdit")
	}
	if s.Privileged {
		add("privileged est interdit")
	}
	if len(s.CapAdd) > 0 {
		add("cap_add est interdit")
	}
	if len(s.Devices) > 0 {
		add("devices est interdit")
	}
	if s.NetworkMode == "host" || strings.HasPrefix(s.NetworkMode, "container:") {
		add("network_mode %q est interdit", s.NetworkMode)
	}
	for _, ns := range [][2]string{{"pid", s.Pid}, {"ipc", s.Ipc}, {"uts", s.Uts}, {"userns_mode", s.UsernsMode}} {
		if ns[1] == "host" {
			add("%s: host est interdit", ns[0])
		}
	}
	for _, opt := range s.SecurityOpt {
		if strings.Contains(opt, "unconfined") || strings.Contains(opt, "disable") {
			add("security_opt %q désactive une protection", opt)
		}
	}
	errs = append(errs, checkEnvironment(name, s.Environment, m)...)
	return errs
}

// imageTag renvoie l'étiquette d'une image épinglée, ou l'erreur qui dit
// pourquoi elle ne l'est pas. Une empreinte seule épingle aussi : l'étiquette
// rendue est alors vide.
func imageTag(image string) (string, error) {
	if image == "" {
		return "", fmt.Errorf("image obligatoire")
	}
	if strings.Contains(image, "$") {
		return "", fmt.Errorf("image %q : l'image ne dépend pas d'une variable", image)
	}
	ref, digest, pinned := strings.Cut(image, "@")
	if pinned && !strings.HasPrefix(digest, "sha256:") {
		return "", fmt.Errorf("image %q : empreinte sha256 attendue", image)
	}
	tag := ""
	if i := strings.LastIndex(ref, ":"); i > strings.LastIndex(ref, "/") {
		tag = ref[i+1:]
	}
	switch {
	case tag == "" && pinned:
		return "", nil
	case tag == "":
		return "", fmt.Errorf("image %q : ajoutez une étiquette de version (postgres:17.2-alpine), jamais latest", image)
	case slices.Contains(floatingTags, strings.ToLower(tag)) || !strings.ContainsAny(tag, "0123456789"):
		return "", fmt.Errorf("image %q : l'étiquette %q change sous le même nom ; épinglez une version", image, tag)
	}
	return tag, nil
}

func checkEnvironment(name string, raw json.RawMessage, m Manifest) []error {
	env, err := environment(raw)
	if err != nil {
		return []error{problem(ComposeFile, "services.%s.environment : %v", name, err)}
	}
	var errs []error
	for _, key := range sortedKeys(env) {
		value := env[key]
		where := fmt.Sprintf("services.%s.environment.%s", name, key)
		if value == nil {
			errs = append(errs, problem(ComposeFile,
				"%s : sans valeur, elle viendrait de l'environnement de la commande ; écrivez %s: ${NOM}", where, key))
			continue
		}
		if !secretLike.MatchString(key) {
			continue
		}
		ref := reference.FindStringSubmatch(*value)
		if ref == nil || ref[0] != *value {
			errs = append(errs, problem(ComposeFile,
				"%s : secret en dur ; la valeur doit être ${NOM}, une variable secret: true du manifeste", where))
			continue
		}
		if v, ok := variable(m, ref[1]); ok && !v.Secret {
			errs = append(errs, problem(ComposeFile,
				"%s : la variable %s doit être marquée secret: true dans le manifeste", where, ref[1]))
		}
	}
	return errs
}

// environment lit les deux formes de Compose : liste « NOM=valeur » ou
// table. Une entrée sans valeur vaut nil.
func environment(raw json.RawMessage) (map[string]*string, error) {
	if !isSet(raw) {
		return nil, nil
	}
	var table map[string]*string
	if err := json.Unmarshal(raw, &table); err == nil {
		return table, nil
	}
	var list []string
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, fmt.Errorf("forme inattendue")
	}
	table = make(map[string]*string, len(list))
	for _, entry := range list {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			table[key] = &value
		} else {
			table[key] = nil
		}
	}
	return table, nil
}

func checkMain(m Manifest, mod model) []error {
	s, ok := mod.Services[m.Service]
	if !ok {
		return []error{problem(ManifestFile, "service : %q n'existe pas dans %s", m.Service, ComposeFile)}
	}
	var errs []error
	if tag, err := imageTag(s.Image); err == nil && m.Version != "" && !strings.Contains(tag, m.Version) {
		errs = append(errs, problem(ManifestFile,
			"version : %q doit figurer dans l'étiquette de l'image du service principal (%s)", m.Version, s.Image))
	}
	h := s.Healthcheck
	if h == nil || h.Disable || !isSet(h.Test) || strings.Contains(string(h.Test), `"NONE"`) {
		errs = append(errs, problem(ComposeFile,
			"services.%s : le service principal a une sonde de santé (healthcheck), que l'installation attend", m.Service))
	}
	return errs
}

func checkVolumes(m Manifest, mod model, used map[string]bool) []error {
	var errs []error
	declared := map[string]bool{}
	for _, v := range m.Volumes {
		declared[v.Name] = true
		if _, ok := mod.Volumes[v.Name]; !ok {
			errs = append(errs, problem(ManifestFile, "volumes : %q n'est pas un volume de %s", v.Name, ComposeFile))
		}
	}
	for _, name := range sortedKeys(mod.Volumes) {
		if mod.Volumes[name].External {
			errs = append(errs, problem(ComposeFile, "volumes.%s : un volume external partagerait des données hors de la stack", name))
		}
		if !declared[name] {
			errs = append(errs, problem(ManifestFile, "volumes : %q, volume de %s, n'est pas déclaré avec sa description", name, ComposeFile))
		}
		if !used[name] {
			errs = append(errs, problem(ComposeFile, "volumes.%s : déclaré mais monté par aucun service", name))
		}
	}
	return errs
}

// checkVariables vérifie que les ${NOM} du fichier et les variables du
// manifeste se répondent, et qu'aucune autre forme de substitution n'est
// employée : la valeur par défaut vit dans le manifeste.
func checkVariables(m Manifest, resolved []byte) []error {
	var doc map[string]any
	if err := json.Unmarshal(resolved, &doc); err != nil {
		return nil // déjà signalé par CheckCompose
	}
	var errs []error
	used := map[string]bool{}
	for _, section := range []string{"services", "volumes", "networks", "configs", "secrets"} {
		walkStrings(doc[section], func(s string) {
			plain := strings.ReplaceAll(s, "$$", "")
			for _, ref := range anyReference.FindAllStringSubmatch(plain, -1) {
				used[ref[1]] = true
			}
			if strings.Contains(reference.ReplaceAllString(plain, ""), "$") {
				errs = append(errs, problem(ComposeFile,
					"%q : seule la forme ${NOM} est admise (la valeur par défaut va dans le manifeste ; $$ pour un $ littéral)", s))
			}
		})
	}
	for _, name := range sortedKeys(used) {
		if _, ok := variable(m, name); !ok {
			errs = append(errs, problem(ComposeFile, "${%s} n'est pas déclarée dans les variables du manifeste", name))
		}
	}
	for _, v := range m.Variables {
		if !used[v.Name] {
			errs = append(errs, problem(ManifestFile, "variables : %s ne sert nulle part dans %s", v.Name, ComposeFile))
		}
	}
	return errs
}

func walkStrings(node any, visit func(string)) {
	switch n := node.(type) {
	case string:
		visit(n)
	case []any:
		for _, item := range n {
			walkStrings(item, visit)
		}
	case map[string]any:
		for _, key := range sortedKeys(n) {
			walkStrings(n[key], visit)
		}
	}
}

func variable(m Manifest, name string) (Variable, bool) {
	for _, v := range m.Variables {
		if v.Name == name {
			return v, true
		}
	}
	return Variable{}, false
}

func isSet(raw json.RawMessage) bool {
	s := strings.TrimSpace(string(raw))
	return s != "" && s != "null" && s != "[]" && s != "{}"
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
