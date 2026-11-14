package update

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Docker lance la CLI docker (le conteneur de mise à jour l'embarque). Les
// arguments sont passés tels quels, jamais par un shell.
type Docker interface {
	Run(ctx context.Context, args ...string) (string, error)
}

// Étapes d'une mise à jour, telles que le dashboard les affiche.
const (
	StepCheck    = "verification"
	StepPull     = "telechargement"
	StepBackup   = "sauvegarde"
	StepSwitch   = "bascule"
	StepHealth   = "sante"
	StepRollback = "retour_arriere"
	StepDone     = "terminee"
)

// rollbackHealthTimeout est le délai minimal laissé à l'ancienne version pour
// redevenir saine après un retour arrière.
const rollbackHealthTimeout = 2 * time.Minute

// Résultats possibles.
const (
	ResultRunning    = "en_cours"
	ResultSucceeded  = "reussie"
	ResultRolledBack = "annulee"
	ResultFailed     = "echec"
)

// State est l'état d'une mise à jour, écrit dans DataDir/maj/etat.json à
// chaque étape : le serveur relancé (nouveau ou ancien) le lit pour afficher
// le résultat.
type State struct {
	From     string    `json:"de"`
	To       string    `json:"vers"`
	Step     string    `json:"etape"`
	Result   string    `json:"resultat"`
	Reason   string    `json:"raison,omitempty"`
	Backup   string    `json:"sauvegarde,omitempty"`
	Started  time.Time `json:"debut"`
	Finished time.Time `json:"fin,omitzero"`
}

// Engine applique une mise à jour sur une installation standard (ADR-0009) :
// InstallDir contient docker-compose.yml et .env, DataDir est le volume de
// données du serveur (la base SQLite).
type Engine struct {
	Docker        Docker
	InstallDir    string
	DataDir       string
	Service       string        // service compose du serveur : « server »
	Container     string        // nom du conteneur : « forgenet-server »
	DatabaseFile  string        // nom de la base dans DataDir : « forgenet.db »
	HealthTimeout time.Duration // 90 s par défaut : le premier contrôle arrive après 30 s
	Keep          int           // sauvegardes gardées : 5 par défaut
	Force         bool          // autorise une version égale ou plus ancienne
	Now           func() time.Time
	Sleep         func(context.Context, time.Duration) error
	Report        func(State) // facultatif : suivi en direct (terminal)
	// Fetch télécharge un fichier d'installation de la version visée ; requis
	// quand la fiche en liste (Release.Files).
	Fetch func(ctx context.Context, version, name string) ([]byte, error)
}

// Erreurs d'une mise à jour refusée avant toute modification.
var (
	ErrBusy      = errors.New("une mise à jour est déjà en cours")
	ErrNotNewer  = errors.New("la version visée n'est pas plus récente que la version installée")
	ErrNoInstall = errors.New("installation introuvable")
)

func (e *Engine) defaults() {
	if e.Service == "" {
		e.Service = "server"
	}
	if e.Container == "" {
		e.Container = "forgenet-server"
	}
	if e.DatabaseFile == "" {
		e.DatabaseFile = "forgenet.db"
	}
	if e.HealthTimeout <= 0 {
		e.HealthTimeout = 90 * time.Second
	}
	if e.Keep <= 0 {
		e.Keep = 5
	}
	if e.Now == nil {
		e.Now = time.Now
	}
	if e.Sleep == nil {
		e.Sleep = func(ctx context.Context, d time.Duration) error {
			t := time.NewTimer(d)
			defer t.Stop()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-t.C:
				return nil
			}
		}
	}
}

func (e *Engine) envPath() string   { return filepath.Join(e.InstallDir, ".env") }
func (e *Engine) stateDir() string  { return filepath.Join(e.DataDir, "maj") }
func (e *Engine) statePath() string { return filepath.Join(e.stateDir(), "etat.json") }

// ReadState lit l'état de la dernière mise à jour (absent : nil, nil).
func ReadState(dataDir string) (*State, error) {
	data, err := os.ReadFile(filepath.Join(dataDir, "maj", "etat.json")) //nolint:gosec // chemin de l'installation, fixé par l'administrateur
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Apply met l'installation à jour vers target. Avant l'arrêt du serveur, un
// refus ou un échec ne change rien ; après, tout échec ramène l'ancienne
// version et la base sauvegardée.
func (e *Engine) Apply(ctx context.Context, target Release) (State, error) {
	e.defaults()
	if err := os.MkdirAll(e.stateDir(), 0o750); err != nil {
		return State{}, err
	}
	_ = sameOwner(e.stateDir(), e.DataDir)
	if err := e.lock(); err != nil {
		return State{}, err
	}
	defer func() { _ = os.Remove(filepath.Join(e.stateDir(), "verrou")) }()

	env, err := readEnv(e.envPath())
	if err != nil {
		return State{}, fmt.Errorf("%w : %v", ErrNoInstall, err)
	}
	st := State{From: env["FORGENET_VERSION"], To: target.Version, Step: StepCheck, Result: ResultRunning, Started: e.Now().UTC()}
	e.save(&st)

	from, err := ParseVersion(st.From)
	if err != nil {
		return e.fail(&st, fmt.Errorf("%w : version installée %q illisible", ErrNoInstall, st.From))
	}
	to, _ := ParseVersion(target.Version)
	if to.Compare(from) <= 0 && !e.Force {
		return e.fail(&st, fmt.Errorf("%w (%s installée, %s visée) : la base est déjà migrée", ErrNotNewer, st.From, st.To))
	}

	// L'image est tirée avant d'arrêter quoi que ce soit : un registre
	// injoignable ou une empreinte inconnue ne coûte aucune coupure.
	e.step(&st, StepPull)
	if _, err := e.Docker.Run(ctx, "pull", target.ImageRef()); err != nil {
		return e.fail(&st, fmt.Errorf("téléchargement de %s : %w", target.ImageRef(), err))
	}
	files, err := e.fetchFiles(ctx, target)
	if err != nil {
		return e.fail(&st, err)
	}

	e.step(&st, StepBackup)
	if _, err := e.compose(ctx, "stop", e.Service); err != nil {
		return e.fail(&st, fmt.Errorf("arrêt du serveur : %w", err))
	}
	backup, err := e.backup(st.From)
	if err != nil {
		// Le serveur est arrêté mais rien n'a changé : on le relance tel quel.
		_, _ = e.compose(ctx, "up", "-d", "--no-deps", e.Service)
		return e.fail(&st, fmt.Errorf("sauvegarde : %w", err))
	}
	st.Backup = backup

	e.step(&st, StepSwitch)
	for _, name := range InstallFiles {
		if data, ok := files[name]; ok {
			if err := writeFile(filepath.Join(e.InstallDir, name), data, 0o644); err != nil {
				return e.rollback(ctx, &st, fmt.Errorf("écriture de %s : %w", name, err))
			}
		}
	}
	if err := writeEnv(e.envPath(), map[string]string{
		"FORGENET_VERSION":      target.Version,
		"FORGENET_IMAGE_DIGEST": target.Digest,
	}); err != nil {
		return e.rollback(ctx, &st, fmt.Errorf("écriture de .env : %w", err))
	}
	// Sans --no-deps : si docker-compose.yml a changé, compose recrée aussi
	// les autres services touchés ; sinon, seul le serveur redémarre.
	// --remove-orphans ne vise que la pile « forgenet » (serveur et Caddy) :
	// chaque projet déployé a sa propre pile « forgenet-<slug> ».
	if _, err := e.compose(ctx, "up", "-d", "--remove-orphans"); err != nil {
		return e.rollback(ctx, &st, fmt.Errorf("démarrage de %s : %w", st.To, err))
	}

	e.step(&st, StepHealth)
	if err := e.waitHealthy(ctx, e.HealthTimeout); err != nil {
		return e.rollback(ctx, &st, fmt.Errorf("la version %s ne démarre pas : %w", st.To, err))
	}

	st.Step, st.Result, st.Finished = StepDone, ResultSucceeded, e.Now().UTC()
	e.save(&st)
	return st, nil
}

// rollback remet .env et la base sauvegardés, puis relance l'ancienne version.
func (e *Engine) rollback(ctx context.Context, st *State, cause error) (State, error) {
	st.Reason = cause.Error()
	e.step(st, StepRollback)
	// Le contexte d'origine a pu expirer : le retour arrière a son propre délai.
	rctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 3*time.Minute)
	defer cancel()
	_, _ = e.compose(rctx, "stop", e.Service)
	if err := e.restore(st.Backup); err != nil {
		return e.fail(st, fmt.Errorf("%s ; retour arrière impossible (%v) : restaurer à la main depuis %s", st.Reason, err, st.Backup))
	}
	if _, err := e.compose(rctx, "up", "-d", "--remove-orphans"); err != nil {
		return e.fail(st, fmt.Errorf("%s ; redémarrage de %s impossible : %w", st.Reason, st.From, err))
	}
	// L'ancienne version a son propre délai, généreux : un délai court donné
	// pour la nouvelle (essai, réseau lent) ne doit pas faire conclure à
	// l'échec d'un retour arrière qui réussit.
	if err := e.waitHealthy(rctx, max(e.HealthTimeout, rollbackHealthTimeout)); err != nil {
		return e.fail(st, fmt.Errorf("%s ; %s ne redémarre pas : %w", st.Reason, st.From, err))
	}
	st.Result, st.Finished = ResultRolledBack, e.Now().UTC()
	e.save(st)
	return *st, cause
}

// staleLock : un verrou plus vieux vient d'une mise à jour interrompue
// (conteneur tué, hôte redémarré) ; le moteur finit en quelques minutes.
const staleLock = 30 * time.Minute

// lock prend le verrou ; un verrou abandonné est repris plutôt que de
// bloquer toutes les mises à jour suivantes.
func (e *Engine) lock() error {
	path := filepath.Join(e.stateDir(), "verrou")
	for range 2 {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600) //nolint:gosec // chemin de l'installation, fixé par l'administrateur
		if err == nil {
			return f.Close()
		}
		if !errors.Is(err, os.ErrExist) {
			return err
		}
		info, statErr := os.Stat(path)
		if statErr != nil || e.Now().Sub(info.ModTime()) < staleLock {
			return ErrBusy
		}
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return ErrBusy
}

func (e *Engine) fail(st *State, err error) (State, error) {
	st.Result, st.Reason, st.Finished = ResultFailed, err.Error(), e.Now().UTC()
	e.save(st)
	return *st, err
}

func (e *Engine) step(st *State, step string) {
	st.Step = step
	e.save(st)
}

// save écrit l'état de façon atomique ; une erreur d'écriture n'interrompt
// pas la mise à jour, l'état n'est qu'un suivi.
func (e *Engine) save(st *State) {
	if e.Report != nil {
		e.Report(*st)
	}
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return
	}
	tmp := e.statePath() + ".tmp"
	if os.WriteFile(tmp, data, 0o600) == nil {
		_ = sameOwner(tmp, e.DataDir)
		_ = os.Rename(tmp, e.statePath())
	}
}

func (e *Engine) compose(ctx context.Context, args ...string) (string, error) {
	return e.Docker.Run(ctx, append([]string{"compose", "--project-directory", e.InstallDir,
		"-f", filepath.Join(e.InstallDir, "docker-compose.yml")}, args...)...)
}

// waitHealthy attend l'état « healthy » du conteneur. Un conteneur arrêté ou
// « unhealthy » échoue tout de suite, sans attendre le délai.
func (e *Engine) waitHealthy(ctx context.Context, timeout time.Duration) error {
	deadline := e.Now().Add(timeout)
	last := ""
	for {
		out, err := e.Docker.Run(ctx, "inspect", "-f", "{{.State.Status}} {{if .State.Health}}{{.State.Health.Status}}{{end}}", e.Container)
		if err == nil {
			last = strings.TrimSpace(out)
			switch {
			case strings.HasSuffix(last, "healthy") && !strings.HasSuffix(last, "unhealthy"):
				return nil
			case strings.HasSuffix(last, "unhealthy"), strings.HasPrefix(last, "exited"), strings.HasPrefix(last, "dead"):
				return fmt.Errorf("conteneur %s", last)
			}
		}
		if !e.Now().Before(deadline) {
			if last == "" {
				last = "inconnu"
			}
			return fmt.Errorf("toujours pas sain après %s (état : %s)", timeout, last)
		}
		if err := e.Sleep(ctx, 2*time.Second); err != nil {
			return err
		}
	}
}

// dbFiles : la base et son journal WAL, qui doivent voyager ensemble.
func (e *Engine) dbFiles() []string {
	return []string{e.DatabaseFile, e.DatabaseFile + "-wal", e.DatabaseFile + "-shm"}
}

// backup copie la base (serveur arrêté) et .env dans
// InstallDir/sauvegardes/<date>-<version>/, puis ne garde que les plus récentes.
func (e *Engine) backup(version string) (string, error) {
	root := filepath.Join(e.InstallDir, "sauvegardes")
	dir := filepath.Join(root, e.Now().UTC().Format("20060102-150405")+"-"+safeName(version))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	for _, name := range e.dbFiles() {
		if err := copyFile(filepath.Join(e.DataDir, name), filepath.Join(dir, name)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
	}
	if _, err := os.Stat(filepath.Join(dir, e.DatabaseFile)); err != nil {
		return "", fmt.Errorf("base %s introuvable dans %s", e.DatabaseFile, e.DataDir)
	}
	if err := copyFile(e.envPath(), filepath.Join(dir, ".env")); err != nil {
		return "", err
	}
	for _, name := range InstallFiles {
		if err := copyFile(filepath.Join(e.InstallDir, name), filepath.Join(dir, name)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
	}
	e.prune(root)
	return dir, nil
}

// restore remet la base et .env d'une sauvegarde ; les fichiers WAL laissés
// par la nouvelle version sont retirés pour ne pas se mélanger à l'ancienne.
func (e *Engine) restore(dir string) error {
	if dir == "" {
		return errors.New("aucune sauvegarde")
	}
	for _, name := range e.dbFiles() {
		dst := filepath.Join(e.DataDir, name)
		if err := os.Remove(dst); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		err := copyFile(filepath.Join(dir, name), dst)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return err
		}
		if err := sameOwner(dst, e.DataDir); err != nil {
			return err
		}
	}
	for _, name := range InstallFiles {
		if err := copyFile(filepath.Join(dir, name), filepath.Join(e.InstallDir, name)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return copyFile(filepath.Join(dir, ".env"), e.envPath())
}

// fetchFiles télécharge les fichiers d'installation de la fiche et vérifie
// leur empreinte, avant tout arrêt : un fichier altéré annule la mise à jour
// sans coupure.
func (e *Engine) fetchFiles(ctx context.Context, target Release) (map[string][]byte, error) {
	if len(target.Files) == 0 {
		return nil, nil
	}
	if e.Fetch == nil {
		return nil, errors.New("la fiche liste des fichiers d'installation, mais aucun moyen de les télécharger")
	}
	files := map[string][]byte{}
	for name, want := range target.Files {
		data, err := e.Fetch(ctx, target.Version, name)
		if err != nil {
			return nil, fmt.Errorf("téléchargement de %s : %w", name, err)
		}
		if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != want {
			return nil, fmt.Errorf("%s de la version %s : empreinte %s, la fiche annonce %s", name, target.Version, got, want)
		}
		files[name] = data
	}
	return files, nil
}

// writeFile écrit un fichier de façon atomique.
func writeFile(path string, data []byte, mode os.FileMode) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, mode); err != nil { //nolint:gosec // fichier d'installation, lisible par compose
		return err
	}
	return os.Rename(tmp, path)
}

func (e *Engine) prune(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}
	var dirs []string
	for _, en := range entries {
		if en.IsDir() {
			dirs = append(dirs, en.Name())
		}
	}
	sort.Strings(dirs) // le nom commence par la date : ordre chronologique
	for len(dirs) > e.Keep {
		_ = os.RemoveAll(filepath.Join(root, dirs[0]))
		dirs = dirs[1:]
	}
}

func safeName(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '.' || r == '-' || (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') {
			return r
		}
		return '_'
	}, s)
}

// copyFile garde les droits du fichier d'origine (.env reste en 0600).
func copyFile(src, dst string) error {
	in, err := os.Open(src) //nolint:gosec // chemin de l'installation, fixé par l'administrateur
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	tmp := dst + ".tmp"
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode().Perm()) //nolint:gosec // chemin de l'installation, fixé par l'administrateur
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, dst)
}

// readEnv lit un fichier .env simple (CLÉ=valeur, commentaires ignorés).
func readEnv(path string) (map[string]string, error) {
	data, err := os.ReadFile(path) //nolint:gosec // chemin de l'installation, fixé par l'administrateur
	if err != nil {
		return nil, err
	}
	env := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if k, v, ok := strings.Cut(line, "="); ok {
			env[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	return env, nil
}

// writeEnv remplace ou ajoute des clés dans .env en gardant le reste du
// fichier, ses commentaires et ses droits, par écriture atomique.
func writeEnv(path string, set map[string]string) error {
	data, err := os.ReadFile(path) //nolint:gosec // chemin de l'installation, fixé par l'administrateur
	if err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	done := map[string]bool{}
	for i, line := range lines {
		k, _, ok := strings.Cut(line, "=")
		k = strings.TrimSpace(k)
		if v, want := set[k]; ok && want && !strings.HasPrefix(strings.TrimSpace(line), "#") {
			lines[i] = k + "=" + v
			done[k] = true
		}
	}
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if !done[k] {
			lines = append(lines, k+"="+set[k])
		}
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(strings.Join(lines, "\n")+"\n"), info.Mode().Perm()); err != nil { //nolint:gosec // chemin de l'installation, fixé par l'administrateur
		return err
	}
	return os.Rename(tmp, path)
}
