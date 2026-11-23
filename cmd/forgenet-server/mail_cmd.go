package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/rs/zerolog"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/mail"
	"github.com/forgelab-dev/lfa-forgenet/internal/sealed"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
)

// errMailUsage rappelle la syntaxe de « forgenet-server mail ».
//
// « mail reset » efface les réglages du relais faits depuis le dashboard :
// ceux de .env s'appliquent de nouveau (ADR-0015). install.sh le lance quand
// il est relancé avec --smtp-host, pour que la dernière action explicite
// l'emporte.
var errMailUsage = errors.New("usage :\n" +
	"  forgenet-server mail reset [--config FICHIER]")

// runMail traite la sous-commande « mail ».
func runMail(ctx context.Context, args []string, stdout io.Writer) error {
	if len(args) == 0 || args[0] != "reset" {
		return errMailUsage
	}
	flags := flag.NewFlagSet("forgenet-server mail reset", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var configPath string
	flags.StringVar(&configPath, "config", "", "Chemin vers le fichier de configuration (YAML)")
	if err := flags.Parse(args[1:]); err != nil {
		return fmt.Errorf("%w ; %w", err, errMailUsage)
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	db, err := sqlite.Open(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("base de données : %w", err)
	}
	defer func() { _ = db.Close() }()

	gone, err := db.MailSettings().DeleteMailSettings(ctx)
	if err != nil {
		return err
	}
	if !gone {
		_, err = fmt.Fprintln(stdout, "Aucun réglage d'envoi fait depuis le dashboard : ceux du serveur s'appliquent déjà.")
		return err
	}
	// Aucun compte n'agit sur la machine du serveur.
	_ = audit.NewAuditor(db.Audit(), zerolog.Nop()).Log(ctx, "mail.settings_change", "mail:relay",
		map[string]any{"action": "rendu au serveur", "by": "server"})
	_, err = fmt.Fprintln(stdout, "Réglages d'envoi du dashboard effacés : ceux du serveur (.env) s'appliquent au prochain démarrage.")
	return err
}

// secretBox prépare le chiffrement des secrets gardés en base (ADR-0015) :
// le mot de passe du relais, les valeurs des installations (ADR-0019). nil
// sans clé réglée.
func secretBox(cfg *config.Config) (*sealed.Box, error) {
	if cfg.Secrets.Key == "" {
		return nil, nil
	}
	key, err := sealed.ParseKey(cfg.Secrets.Key)
	if err != nil {
		return nil, err
	}
	return sealed.NewBox(key)
}

// openRelay assemble le relais : .env, puis les réglages du dashboard s'ils
// existent et que leur mot de passe s'ouvre (ADR-0015).
func openRelay(ctx context.Context, cfg *config.Config, db store.Store) (*mail.Relay, error) {
	box, err := secretBox(cfg)
	if err != nil {
		return nil, err
	}
	relay := mail.NewRelay(cfg.SMTP, db.MailSettings(), box)
	if err := relay.Load(ctx); err != nil {
		return nil, fmt.Errorf("réglages d'envoi : %w", err)
	}
	return relay, nil
}

// logRelay dit au démarrage d'où viennent les réglages d'envoi, sans
// identifiants.
func logRelay(logger zerolog.Logger, relay *mail.Relay) {
	cur := relay.Current()
	if cur.Unreadable {
		logger.Warn().Msg("réglages d'envoi du dashboard illisibles (FORGENET_SECRET_KEY changée ou absente) : ceux de .env s'appliquent ; ressaisissez le mot de passe du relais dans Réglages")
	}
	if !relay.Editable() {
		logger.Info().Msg("pas de FORGENET_SECRET_KEY : le relais ne se règle pas depuis le dashboard ; la prochaine mise à jour crée la clé, ou relancez install.sh")
	}
	switch {
	case cur.Enabled:
		logger.Info().Str("relais", cur.Config.Host).Str("source", string(cur.Source)).Msg("envoi d'e-mails activé")
	case cur.Source == mail.SourceDashboard:
		logger.Info().Msg("envoi d'e-mails coupé depuis le dashboard")
	default:
		logger.Info().Msg("envoi d'e-mails désactivé : aucun relais SMTP configuré")
	}
}
