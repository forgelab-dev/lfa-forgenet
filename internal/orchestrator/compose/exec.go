package compose

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	// waitDelay laisse à la commande le temps de se terminer après l'annulation.
	waitDelay = 5 * time.Second
	// maxErrorOutput borne le message gardé de la sortie d'erreur.
	maxErrorOutput = 4 << 10
)

// CommandError est une exécution de docker compose terminée en erreur.
type CommandError struct {
	Args     []string
	ExitCode int
	Output   string // fin de la sortie d'erreur
}

// subcommands sont les sous-commandes que lance le moteur, pour nommer celle qui échoue.
var subcommands = map[string]bool{"config": true, "up": true, "down": true, "ps": true, "version": true}

func (e *CommandError) Error() string {
	sub := "?"
	for _, a := range e.Args {
		if subcommands[a] {
			sub = a
			break
		}
	}
	return fmt.Sprintf("docker compose %s : code %d : %s", sub, e.ExitCode, e.Output)
}

type execRunner struct{}

func (execRunner) Run(ctx context.Context, c Command) ([]byte, error) {
	// Le binaire est fixe et les arguments viennent du moteur (options
	// constantes, nom de projet validé) ; ils ne passent jamais par un shell.
	cmd := exec.CommandContext(ctx, "docker", append([]string{"compose"}, c.Args...)...) //nolint:gosec // G204 : voir ci-dessus
	cmd.Env = c.Env
	cmd.Dir = c.Dir
	cmd.WaitDelay = waitDelay
	if c.Stdin != nil {
		cmd.Stdin = bytes.NewReader(c.Stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	err := cmd.Run()
	if err == nil {
		return stdout.Bytes(), nil
	}
	if errors.Is(err, exec.ErrNotFound) {
		return nil, fmt.Errorf("%w : commande docker introuvable", ErrUnavailable)
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, fmt.Errorf("docker compose interrompu : %w", ctxErr)
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return nil, fmt.Errorf("exécution de docker compose : %w", err)
	}

	output := strings.TrimSpace(stderr.String())
	// Docker sans le plugin Compose répond par une erreur d'usage.
	if strings.Contains(output, "is not a docker command") {
		return nil, fmt.Errorf("%w : plugin compose absent", ErrUnavailable)
	}
	if len(output) > maxErrorOutput {
		output = "…" + output[len(output)-maxErrorOutput:]
	}
	return nil, &CommandError{Args: c.Args, ExitCode: exitErr.ExitCode(), Output: output}
}
