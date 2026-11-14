package update

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// UpdaterName est le conteneur de mise à jour : un seul à la fois.
const UpdaterName = "forgenet-updater"

// DataMount est le dossier de données du serveur dans son image (Dockerfile).
const DataMount = "/var/lib/forgenet"

// CLI lance la CLI docker, sans shell.
type CLI struct{}

// Run exécute « docker args… » et renvoie sa sortie.
func (CLI) Run(ctx context.Context, args ...string) (string, error) {
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

// LaunchOptions décrit la mise à jour à lancer.
type LaunchOptions struct {
	Container  string // conteneur du serveur : « forgenet-server » par défaut
	InstallDir string // dossier d'installation, chemin absolu sur l'hôte
	To         string // version visée (X.Y.Z) ou « latest »
	Pre        bool
	Force      bool
}

// ErrUnsupported : l'installation n'est pas une installation standard
// (FORGENET_INSTALL_DIR inconnu), la mise à jour se fait à la main.
var ErrUnsupported = errors.New("dossier d'installation inconnu : FORGENET_INSTALL_DIR (chemin absolu sur l'hôte) doit être fixé")

// Launch démarre le conteneur de mise à jour à côté du serveur et renvoie son
// identifiant. Il tourne depuis l'image en service (un code connu), en root
// pour écrire .env et la sauvegarde, avec le socket Docker, le dossier
// d'installation et le volume de données : rien de plus que ce que le serveur
// a déjà (ADR-0011).
func Launch(ctx context.Context, d Docker, o LaunchOptions) (string, error) {
	if o.InstallDir == "" || !filepath.IsAbs(o.InstallDir) {
		return "", ErrUnsupported
	}
	if o.Container == "" {
		o.Container = "forgenet-server"
	}
	if o.To == "" {
		o.To = "latest"
	}
	// Le conteneur d'une mise à jour précédente reste, arrêté, pour ses
	// journaux ; il est retiré ici, sauf s'il tourne encore.
	if status, err := d.Run(ctx, "inspect", "-f", "{{.State.Status}}", UpdaterName); err == nil {
		if s := strings.TrimSpace(status); s == "running" || s == "restarting" || s == "created" {
			return "", ErrBusy
		}
		if _, err := d.Run(ctx, "rm", UpdaterName); err != nil {
			return "", err
		}
	}
	image, err := d.Run(ctx, "inspect", "-f", "{{.Config.Image}}", o.Container)
	if err != nil {
		return "", err
	}
	volume, err := d.Run(ctx, "inspect", "-f", `{{range .Mounts}}{{if eq .Destination "`+DataMount+`"}}{{.Name}}{{end}}{{end}}`, o.Container)
	if err != nil {
		return "", err
	}
	image, volume = strings.TrimSpace(image), strings.TrimSpace(volume)
	if volume == "" {
		return "", fmt.Errorf("volume de données du conteneur %s introuvable", o.Container)
	}
	args := []string{
		"run", "-d", "--name", UpdaterName,
		"--user", "0:0",
		"-e", "FORGENET_INSTALL_DIR=" + o.InstallDir,
		"-v", "/var/run/docker.sock:/var/run/docker.sock",
		"-v", o.InstallDir + ":" + o.InstallDir,
		"-v", volume + ":" + DataMount,
		image, "update", "apply", "--to", o.To, "--install-dir", o.InstallDir,
	}
	if o.Pre {
		args = append(args, "--pre")
	}
	if o.Force {
		args = append(args, "--force")
	}
	id, err := d.Run(ctx, args...)
	if err != nil {
		return "", err
	}
	id = strings.TrimSpace(id)
	return id[:min(12, len(id))], nil
}
