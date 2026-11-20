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
	"github.com/forgelab-dev/lfa-forgenet/internal/mail"
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
//
// « user create --no-password » crée un compte qui choisira son mot de passe
// par un lien, et « user login-link » fait ce lien de première connexion :
// affiché, ou envoyé à l'adresse du compte avec --send (ADR-0014). C'est ce
// qu'install.sh fait sans terminal, au lieu d'écrire un mot de passe en clair.
//
// « user suspend », « resume » et « delete » (#103) sont les gestes de la page
// Équipe, pour qui n'a plus d'administrateur pour les faire : le serveur garde
// toujours un administrateur actif. « delete » fait retaper le nom du compte.
//
// « user set-email » corrige l'adresse d'un compte (#100), comme la page
// Équipe ; avec un relais, un avis part à l'ancienne adresse.
var errUserUsage = errors.New("usage :\n" +
	"  forgenet-server user create --username NOM --email ADRESSE [--role admin|member] [--no-password] [--config FICHIER]\n" +
	"  forgenet-server user set-password --username NOM [--config FICHIER]\n" +
	"  forgenet-server user set-email --username NOM --email ADRESSE [--config FICHIER]\n" +
	"  forgenet-server user login-link --username NOM [--send] [--config FICHIER]\n" +
	"  forgenet-server user suspend|resume --username NOM [--config FICHIER]\n" +
	"  forgenet-server user delete --username NOM [--config FICHIER]\n" +
	"(mot de passe, et nom retapé pour delete, sur l'entrée standard)")

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
	case "set-email":
		return runUserSetEmail(ctx, args[1:], stdout)
	case "login-link":
		return runUserLoginLink(ctx, args[1:], stdout)
	case "suspend", "resume", "delete":
		return runUserAccountAction(ctx, args[0], args[1:], stdin, stdout)
	default:
		return errUserUsage
	}
}

func runUserCreate(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer) error {
	flags := flag.NewFlagSet("forgenet-server user create", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var configPath, username, email, role string
	var noPassword bool
	flags.StringVar(&configPath, "config", "", "Chemin vers le fichier de configuration (YAML)")
	flags.StringVar(&username, "username", "", "Nom d'utilisateur")
	flags.StringVar(&email, "email", "", "Adresse email")
	flags.StringVar(&role, "role", string(store.RoleAdmin), "Rôle : admin ou member")
	flags.BoolVar(&noPassword, "no-password", false, "Créer le compte sans mot de passe : il le choisira par un lien (user login-link)")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("%w ; %w", err, errUserUsage)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	var password string
	if !noPassword {
		if password, err = readPassword(stdin, stdout, "Mot de passe : "); err != nil {
			return err
		}
	}

	return withAuth(ctx, cfg, func(authService auth.Service, _ store.Store) error {
		var user *store.User
		if noPassword {
			user, err = authService.RegisterWithoutPassword(ctx, username, email, store.Role(role))
		} else {
			user, err = authService.Register(ctx, username, email, password, store.Role(role))
		}
		if errors.Is(err, store.ErrConflict) {
			return errors.New("ce nom d'utilisateur ou cette adresse email existe déjà")
		}
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(stdout, "Compte %s créé (rôle %s, identifiant %s).\n", user.Username, user.Role, user.ID)
		if err == nil && noPassword {
			_, err = fmt.Fprintf(stdout, "Il s'ouvrira par son lien de première connexion : forgenet user login-link --username %s\n", user.Username)
		}
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

	return withAuth(ctx, cfg, func(authService auth.Service, _ store.Store) error {
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

// runUserSetEmail change l'adresse d'un compte. Être sur la machine du
// serveur en est la preuve, comme pour « set-password ». Les adresses ne sont
// pas affichées : cette sortie peut finir dans un journal.
func runUserSetEmail(ctx context.Context, args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("forgenet-server user set-email", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var configPath, username, email string
	flags.StringVar(&configPath, "config", "", "Chemin vers le fichier de configuration (YAML)")
	flags.StringVar(&username, "username", "", "Nom d'utilisateur")
	flags.StringVar(&email, "email", "", "Nouvelle adresse e-mail")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("%w ; %w", err, errUserUsage)
	}
	username = strings.TrimSpace(username)
	switch {
	case username == "":
		return fmt.Errorf("--username est obligatoire ; %w", errUserUsage)
	case strings.TrimSpace(email) == "":
		return fmt.Errorf("--email est obligatoire ; %w", errUserUsage)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	return withAuth(ctx, cfg, func(authService auth.Service, db store.Store) error {
		user, err := db.Users().GetUserByUsername(ctx, username)
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("aucun compte ne s'appelle %q", username)
		}
		if err != nil {
			return err
		}
		change, err := authService.SetAccountEmail(ctx, "", user.ID, email)
		switch {
		case errors.Is(err, store.ErrConflict):
			return errors.New("cette adresse a déjà un compte")
		case err != nil:
			return err
		}
		if _, err := fmt.Fprintf(stdout, "Adresse du compte %s changée.\n", user.Username); err != nil {
			return err
		}

		if !change.HadPassword {
			// L'ancienne adresse n'a servi qu'à recevoir le lien, annulé avec
			// elle : il en faut un nouveau.
			_, err = fmt.Fprintf(stdout, "Son lien en attente est annulé ; refaites-le pour la nouvelle adresse : forgenet user login-link --username %s --send\n", user.Username)
			return err
		}
		if !cfg.SMTP.Enabled() {
			return nil
		}
		sendCtx, cancel := context.WithTimeout(ctx, mail.DefaultTimeout)
		defer cancel()
		if err := mail.NewSMTP(cfg.SMTP).Send(sendCtx, mail.EmailChangedMessage(change.OldEmail, user.Username, change.User.Email, "")); err != nil {
			// L'adresse est changée : l'avis manqué ne l'annule pas.
			_, err = fmt.Fprintf(stdout, "L'avis à l'ancienne adresse n'est pas parti : %v\n", err)
			return err
		}
		_, err = fmt.Fprintln(stdout, "Un avis est envoyé à l'ancienne adresse.")
		return err
	})
}

// runUserLoginLink fait un lien de première connexion pour un compte. Être sur
// la machine du serveur en est la preuve, comme pour « set-password ». Sans
// --send, le lien est écrit sur la sortie standard, qu'install.sh redirige
// vers un fichier root quand il n'a pas de terminal : il ne doit finir dans
// aucun journal.
func runUserLoginLink(ctx context.Context, args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("forgenet-server user login-link", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var configPath, username string
	var send bool
	flags.StringVar(&configPath, "config", "", "Chemin vers le fichier de configuration (YAML)")
	flags.StringVar(&username, "username", "", "Nom d'utilisateur")
	flags.BoolVar(&send, "send", false, "Envoyer le lien à l'adresse du compte au lieu de l'afficher")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("%w ; %w", err, errUserUsage)
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("--username est obligatoire ; %w", errUserUsage)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	if send && !cfg.SMTP.Enabled() {
		return fmt.Errorf("%w ; sans --send, le lien s'affiche ici", mail.ErrDisabled)
	}

	return withAuth(ctx, cfg, func(authService auth.Service, db store.Store) error {
		user, err := db.Users().GetUserByUsername(ctx, username)
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("aucun compte ne s'appelle %q", username)
		}
		if err != nil {
			return err
		}
		token, _, err := authService.CreateAccessLink(ctx, user.ID, store.PurposeFirstLogin, "system:cli")
		if errors.Is(err, auth.ErrAccountSuspended) {
			return fmt.Errorf("le compte %s est suspendu : réactivez-le d'abord (forgenet user resume --username %s)", user.Username, user.Username)
		}
		if err != nil {
			return err
		}
		local := strings.TrimSpace(cfg.Proxy.DashboardHost) == ""
		url := mail.AccessURL(mail.DashboardURL(cfg.Proxy.DashboardHost), token)

		if !send {
			_, err = fmt.Fprint(stdout, mail.FirstLoginText(user.Username, url, local, auth.HasPassword(user)))
			return err
		}
		sendCtx, cancel := context.WithTimeout(ctx, mail.DefaultTimeout)
		defer cancel()
		if err := mail.NewSMTP(cfg.SMTP).Send(sendCtx, mail.FirstLoginMessage(user.Email, user.Username, url, local, auth.HasPassword(user))); err != nil {
			return fmt.Errorf("le relais n'a pas pris le lien : %w", err)
		}
		// L'adresse n'est pas répétée : cette sortie peut finir dans un journal.
		_, err = fmt.Fprintf(stdout, "Lien de première connexion envoyé à l'adresse du compte %s (valable 24 heures).\n", user.Username)
		return err
	})
}

// runUserAccountAction suspend, réactive ou supprime un compte. Sur la
// machine du serveur, aucun compte n'agit : seul le dernier administrateur
// actif est protégé.
func runUserAccountAction(ctx context.Context, action string, args []string, stdin io.Reader, stdout io.Writer) error {
	flags := flag.NewFlagSet("forgenet-server user "+action, flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var configPath, username string
	flags.StringVar(&configPath, "config", "", "Chemin vers le fichier de configuration (YAML)")
	flags.StringVar(&username, "username", "", "Nom d'utilisateur")
	if err := flags.Parse(args); err != nil {
		return fmt.Errorf("%w ; %w", err, errUserUsage)
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("--username est obligatoire ; %w", errUserUsage)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	return withAuth(ctx, cfg, func(authService auth.Service, db store.Store) error {
		user, err := db.Users().GetUserByUsername(ctx, username)
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("aucun compte ne s'appelle %q", username)
		}
		if err != nil {
			return err
		}

		var msg string
		switch action {
		case "suspend":
			var closed int64
			if _, closed, err = authService.SuspendAccount(ctx, "", user.ID); err == nil {
				msg = fmt.Sprintf("Compte %s suspendu ; %d session(s) et jeton(s) fermé(s). Pour le rouvrir : forgenet user resume --username %s\n",
					user.Username, closed, user.Username)
			}
		case "resume":
			if _, err = authService.ResumeAccount(ctx, "", user.ID); err == nil {
				msg = fmt.Sprintf("Compte %s réactivé : il peut de nouveau se connecter.\n", user.Username)
			}
		case "delete":
			if err = confirmUsername(stdin, stdout, user.Username); err != nil {
				return err
			}
			if _, err = authService.DeleteAccount(ctx, "", user.ID); err == nil {
				msg = fmt.Sprintf("Compte %s supprimé. Le journal d'audit garde ses entrées.\n", user.Username)
			}
		}
		if errors.Is(err, store.ErrLastAdmin) {
			return fmt.Errorf("%s n'est pas touché : %w", user.Username, err)
		}
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(stdout, msg)
		return err
	})
}

// confirmUsername fait retaper le nom du compte avant une suppression. Une
// entrée redirigée en donne la première ligne : « echo NOM | … » vaut
// confirmation, en connaissance de cause.
func confirmUsername(stdin io.Reader, stdout io.Writer, username string) error {
	_, _ = fmt.Fprintf(stdout, "La suppression de %s est définitive. Tapez son nom pour confirmer : ", username)
	line, err := bufio.NewReader(stdin).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("lecture de la confirmation : %w", err)
	}
	_, _ = fmt.Fprintln(stdout)
	if strings.TrimSpace(line) != username {
		return errors.New("nom différent : rien n'est supprimé")
	}
	return nil
}

// withAuth ouvre la base et le service d'authentification le temps de fn.
func withAuth(ctx context.Context, cfg *config.Config, fn func(auth.Service, store.Store) error) error {
	db, err := sqlite.Open(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("base de données : %w", err)
	}
	defer func() { _ = db.Close() }()
	return fn(auth.NewService(db, cfg.Auth, audit.NewAuditor(db.Audit(), zerolog.Nop())), db)
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
