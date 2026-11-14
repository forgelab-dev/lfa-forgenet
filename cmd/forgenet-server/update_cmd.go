package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/update"
	"github.com/forgelab-dev/lfa-forgenet/internal/version"
)

// errUpdateUsage rappelle la syntaxe de « forgenet-server update » (ADR-0011).
var errUpdateUsage = errors.New(`usage :
  forgenet-server update check [--pre]                 version publiée la plus récente
  forgenet-server update start [--to X.Y.Z] [--pre] [--force]
                                                       lance la mise à jour à côté du serveur
  forgenet-server update apply [--to X.Y.Z] ...        (dans le conteneur de mise à jour)`)

// updaterName est le conteneur de mise à jour : un seul à la fois.
const updaterName = "forgenet-updater"

// dataMount est le dossier de données du serveur dans son image (Dockerfile).
const dataMount = "/var/lib/forgenet"

func runUpdate(ctx context.Context, args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return errUpdateUsage
	}
	switch args[0] {
	case "check":
		return updateCheck(ctx, args[1:], stdout)
	case "start":
		return updateStart(ctx, args[1:], stdout)
	case "apply":
		return updateApply(ctx, args[1:], stdout)
	}
	return errUpdateUsage
}

func source() update.Source {
	return update.Source{UserAgent: "forgenet/" + version.Get().Version}
}

// target trouve la fiche visée : la dernière version, ou celle de --to.
func target(ctx context.Context, to string, pre bool) (update.Release, error) {
	if to == "" || to == "latest" {
		return source().Latest(ctx, pre)
	}
	v, err := update.ParseVersion(to)
	if err != nil {
		return update.Release{}, err
	}
	return source().ByVersion(ctx, v)
}

func updateCheck(ctx context.Context, args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("forgenet-server update check", flag.ContinueOnError)
	pre := flags.Bool("pre", false, "inclure les pré-versions")
	if err := flags.Parse(args); err != nil {
		return err
	}
	current := version.Get().Version
	rel, err := target(ctx, "", *pre)
	if err != nil {
		return err
	}
	cur, curErr := update.ParseVersion(current)
	latest, _ := update.ParseVersion(rel.Version)
	switch {
	case curErr == nil && latest.Compare(cur) <= 0:
		_, err = fmt.Fprintf(stdout, "ForgeNet %s est à jour (dernière version publiée : %s).\n", current, rel.Version)
	default:
		_, err = fmt.Fprintf(stdout, "Version installée : %s\nVersion disponible : %s (%s)\nNotes : %s\n",
			current, rel.Version, rel.Digest, rel.NotesURL)
	}
	return err
}

// dockerCLI lance la CLI docker sans shell.
type dockerCLI struct{}

func (dockerCLI) Run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", args...) //nolint:gosec // arguments construits par le moteur, sans shell
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if len(msg) > 400 {
			msg = msg[len(msg)-400:]
		}
		return string(out), fmt.Errorf("docker %s : %w : %s", args[0], err, msg)
	}
	return string(out), nil
}

var stepLabels = map[string]string{
	update.StepCheck:    "Vérification",
	update.StepPull:     "Téléchargement de l'image",
	update.StepBackup:   "Sauvegarde de la base et de la configuration",
	update.StepSwitch:   "Bascule sur la nouvelle version",
	update.StepHealth:   "Vérification de la santé du serveur",
	update.StepRollback: "Retour à la version précédente",
	update.StepDone:     "Terminée",
}

func updateApply(ctx context.Context, args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("forgenet-server update apply", flag.ContinueOnError)
	to := flags.String("to", "latest", "version visée (X.Y.Z) ou latest")
	pre := flags.Bool("pre", false, "inclure les pré-versions")
	installDir := flags.String("install-dir", envOr("FORGENET_INSTALL_DIR", "/opt/forgenet"), "dossier d'installation (docker-compose.yml, .env)")
	dataDir := flags.String("data-dir", dataMount, "volume de données du serveur")
	timeout := flags.Duration("health-timeout", 90*time.Second, "délai pour que la nouvelle version soit saine")
	force := flags.Bool("force", false, "accepter une version égale ou plus ancienne")
	if err := flags.Parse(args); err != nil {
		return err
	}
	rel, err := target(ctx, *to, *pre)
	if err != nil {
		return err
	}
	last := ""
	engine := &update.Engine{
		Docker: dockerCLI{}, InstallDir: *installDir, DataDir: *dataDir,
		HealthTimeout: *timeout, Force: *force, Fetch: source().FetchFile,
		Report: func(s update.State) {
			if s.Step != last {
				last = s.Step
				_, _ = fmt.Fprintf(stdout, "==> %s\n", stepLabels[s.Step])
			}
		},
	}
	st, err := engine.Apply(ctx, rel)
	if err != nil {
		if st.Backup != "" {
			_, _ = fmt.Fprintf(stdout, "Sauvegarde : %s\n", st.Backup)
		}
		return err
	}
	_, err = fmt.Fprintf(stdout, "ForgeNet %s est installé (sauvegarde : %s).\n", st.To, st.Backup)
	return err
}

// updateStart lance le conteneur de mise à jour à côté du serveur. Il tourne
// depuis l'image en service (un code connu), en root pour écrire .env et la
// sauvegarde, avec le socket Docker, le dossier d'installation et le volume de
// données : rien de plus que ce que le serveur a déjà (ADR-0011).
func updateStart(ctx context.Context, args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("forgenet-server update start", flag.ContinueOnError)
	to := flags.String("to", "latest", "version visée (X.Y.Z) ou latest")
	pre := flags.Bool("pre", false, "inclure les pré-versions")
	container := flags.String("container", "forgenet-server", "conteneur du serveur")
	installDir := flags.String("install-dir", os.Getenv("FORGENET_INSTALL_DIR"), "dossier d'installation sur l'hôte")
	force := flags.Bool("force", false, "accepter une version égale ou plus ancienne")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *installDir == "" || !filepath.IsAbs(*installDir) {
		return errors.New("dossier d'installation inconnu : FORGENET_INSTALL_DIR (chemin absolu sur l'hôte) doit être fixé")
	}
	// La version visée est vérifiée tout de suite : une erreur se lit ici,
	// pas dans les journaux d'un conteneur.
	rel, err := target(ctx, *to, *pre)
	if err != nil {
		return err
	}
	if !*force {
		cur, curErr := update.ParseVersion(version.Get().Version)
		next, _ := update.ParseVersion(rel.Version)
		if curErr == nil && next.Compare(cur) <= 0 {
			if *to == "" || *to == "latest" {
				_, err := fmt.Fprintf(stdout, "ForgeNet %s est à jour (dernière version publiée : %s).\n", cur, rel.Version)
				return err
			}
			return fmt.Errorf("%w (%s installée, %s visée)", update.ErrNotNewer, cur, rel.Version)
		}
	}
	d := dockerCLI{}
	// Le conteneur d'une mise à jour précédente reste, arrêté, pour ses
	// journaux ; il est retiré ici, sauf s'il tourne encore.
	if status, err := d.Run(ctx, "inspect", "-f", "{{.State.Status}}", updaterName); err == nil {
		if s := strings.TrimSpace(status); s == "running" || s == "restarting" {
			return update.ErrBusy
		}
		if _, err := d.Run(ctx, "rm", updaterName); err != nil {
			return err
		}
	}
	image, err := d.Run(ctx, "inspect", "-f", "{{.Config.Image}}", *container)
	if err != nil {
		return err
	}
	volume, err := d.Run(ctx, "inspect", "-f", `{{range .Mounts}}{{if eq .Destination "`+dataMount+`"}}{{.Name}}{{end}}{{end}}`, *container)
	if err != nil {
		return err
	}
	image, volume = strings.TrimSpace(image), strings.TrimSpace(volume)
	if volume == "" {
		return fmt.Errorf("volume de données du conteneur %s introuvable", *container)
	}
	runArgs := []string{
		"run", "-d", "--name", updaterName,
		"--user", "0:0",
		"-e", "FORGENET_INSTALL_DIR=" + *installDir,
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-v", *installDir + ":" + *installDir,
		"-v", volume + ":" + dataMount,
		image, "update", "apply", "--to", *to, "--install-dir", *installDir,
	}
	if *pre {
		runArgs = append(runArgs, "--pre")
	}
	if *force {
		runArgs = append(runArgs, "--force")
	}
	id, err := d.Run(ctx, runArgs...)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(stdout, "Mise à jour lancée dans le conteneur %s (%s).\n", updaterName, strings.TrimSpace(id)[:min(12, len(strings.TrimSpace(id)))])
	return err
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
