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
// n'accorde jamais ce rôle d'elle-même. « user set-password » rend l'accès à
// qui a perdu son mot de passe : être sur la machine du serveur en est la
// preuve, comme pour « create ». Le mot de passe est lu sur l'entrée standard
// (saisie masquée dans un terminal), jamais en argument, pour ne pas rester
// dans l'historique du shell.
var errUserUsage = errors.New("usage :\n" +
	"  forgenet-server user create --username NOM --email ADRESSE [--role admin|member] [--config FICHIER]\n" +
	"  forgenet-server user set-password --username NOM [--config FICHIER]\n" +
	"(mot de passe sur l'entrée standard)")

// runUser traite la sous-commande « user ».
func runUser(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) == 0 {
		return errUserUsage
	}
	switch args[0] {
	case "create":
		return runUserCreate(ctx, args[1:], stdin, stdout)
	case "set-password":
		return runUserSetPassword(ctx, args[1:], stdin, stdout)
	default:
		return errUserUsage
	}
}

func runUserCreate(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer) error {
	flags := flag.NewFlagSet("forgenet-server user create", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var configPath, username, email, role string
	flags.StringVar(&configPath, "config", "", "Chemin vers le fichier de configuration (YAML)")
	flags.StringVar(&username, "username", "", "Nom d'utilisateur")
	flags.StringVar(&email, "email", "", "Adresse email")
	flags.StringVar(&role, "role", string(store.RoleAdmin), "Rôle : admin ou member")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("%w ; %w", err, errUserUsage)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	password, err := readPassword(stdin, stdout, "Mot de passe : ")
	if err != nil {
		return err
	}

	return withAuth(ctx, cfg, func(authService auth.Service) error {
		user, err := authService.Register(ctx, username, email, password, store.Role(role))
		if errors.Is(err, store.ErrConflict) {
			return errors.New("ce nom d'utilisateur ou cette adresse email existe déjà")
		}
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(stdout, "Compte %s créé (rôle %s, identifiant %s).\n", user.Username, user.Role, user.ID)
		return err
	})
}

func runUserSetPassword(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer) error {
	flags := flag.NewFlagSet("forgenet-server user set-password", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var configPath, username string
	flags.StringVar(&configPath, "config", "", "Chemin vers le fichier de configuration (YAML)")
	flags.StringVar(&username, "username", "", "Nom d'utilisateur")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("%w ; %w", err, errUserUsage)
	}
	if strings.TrimSpace(username) == "" {
		return fmt.Errorf("--username est obligatoire ; %w", errUserUsage)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	password, err := readNewPassword(stdin, stdout)
	if err != nil {
		return err
	}

	return withAuth(ctx, cfg, func(authService auth.Service) error {
		revoked, err := authService.ResetPassword(ctx, username, password)
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("aucun compte ne s'appelle %q", strings.TrimSpace(username))
		}
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(stdout, "Mot de passe de %s changé ; %d session(s) et jeton(s) fermé(s).\n", strings.TrimSpace(username), revoked)
		return err
	})
}

// withAuth ouvre la base et le service d'authentification le temps de fn.
func withAuth(ctx context.Context, cfg *config.Config, fn func(auth.Service) error) error {
	db, err := sqlite.Open(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("base de données : %w", err)
	}
	defer func() { _ = db.Close() }()
	return fn(auth.NewService(db, cfg.Auth, audit.NewAuditor(db.Audit(), zerolog.Nop())))
}

// readNewPassword demande deux fois le mot de passe dans un terminal, où une
// faute de frappe ne se voit pas ; une entrée redirigée n'en donne qu'un.
func readNewPassword(stdin io.Reader, stdout io.Writer) (string, error) {
	password, err := readPassword(stdin, stdout, "Nouveau mot de passe : ")
	if err != nil || !isTerminal(stdin) {
		return password, err
	}
	again, err := readPassword(stdin, stdout, "Confirmez-le : ")
	if err != nil {
		return "", err
	}
	if again != password {
		return "", errors.New("les deux saisies diffèrent : mot de passe inchangé")
	}
	return password, nil
}

func isTerminal(stdin io.Reader) bool {
	f, ok := stdin.(*os.File)
	return ok && term.IsTerminal(int(f.Fd()))
}

// readPassword lit le mot de passe sans l'afficher dans un terminal, ou la
// première ligne de l'entrée standard quand elle est redirigée.
func readPassword(stdin io.Reader, stdout io.Writer, prompt string) (string, error) {
	if f, ok := stdin.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		_, _ = fmt.Fprint(stdout, prompt)
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
