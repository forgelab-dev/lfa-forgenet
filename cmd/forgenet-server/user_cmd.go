package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"golang.org/x/term"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
)

// errUserUsage rappelle la syntaxe de la sous-commande.
//
// « user create » crée un compte directement dans la base, sur la machine du
// serveur. C'est la seule façon de créer le premier administrateur : l'API
// n'accorde jamais ce rôle d'elle-même. Le mot de passe est lu sur l'entrée
// standard (saisie masquée dans un terminal), jamais en argument, pour ne pas
// rester dans l'historique du shell.
var errUserUsage = errors.New("usage : forgenet-server user create --username NOM --email ADRESSE [--role admin|member] [--config FICHIER] (mot de passe sur l'entrée standard)")

// runUser traite la sous-commande « user ».
func runUser(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) == 0 || args[0] != "create" {
		return errUserUsage
	}

	flags := flag.NewFlagSet("forgenet-server user create", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var configPath, username, email, role string
	flags.StringVar(&configPath, "config", "", "Chemin vers le fichier de configuration (YAML)")
	flags.StringVar(&username, "username", "", "Nom d'utilisateur")
	flags.StringVar(&email, "email", "", "Adresse email")
	flags.StringVar(&role, "role", string(store.RoleAdmin), "Rôle : admin ou member")
	if err := flags.Parse(args[1:]); err != nil {
		return fmt.Errorf("%w ; %w", err, errUserUsage)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	password, err := readPassword(stdin, stdout)
	if err != nil {
		return err
	}

	db, err := sqlite.Open(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("base de données : %w", err)
	}
	defer func() { _ = db.Close() }()

	authService := auth.NewService(db, cfg.Auth, audit.NewAuditor(db.Audit(), zerolog.Nop()))
	user, err := authService.Register(ctx, username, email, password, store.Role(role))
	if errors.Is(err, store.ErrConflict) {
		return errors.New("ce nom d'utilisateur ou cette adresse email existe déjà")
	}
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(stdout, "Compte %s créé (rôle %s, identifiant %s).\n", user.Username, user.Role, user.ID)
	return err
}

// readPassword lit le mot de passe sans l'afficher dans un terminal, ou la
// première ligne de l'entrée standard quand elle est redirigée.
func readPassword(stdin io.Reader, stdout io.Writer) (string, error) {
	if f, ok := stdin.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		_, _ = fmt.Fprint(stdout, "Mot de passe : ")
		b, err := term.ReadPassword(int(f.Fd()))
		_, _ = fmt.Fprintln(stdout)
		if err != nil {
			return "", fmt.Errorf("lecture du mot de passe : %w", err)
		}
		return string(b), nil
	}

	line, err := bufio.NewReader(stdin).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("lecture du mot de passe : %w", err)
	}
	return strings.TrimRight(line, "\r\n"), nil
}
