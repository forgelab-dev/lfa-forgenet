package mail

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/sealed"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

// Source dit d'où viennent les réglages en vigueur (ADR-0015).
type Source string

const (
	// SourceServer : les réglages de .env, écrits par install.sh.
	SourceServer Source = "server"
	// SourceDashboard : les réglages enregistrés depuis la page Réglages.
	SourceDashboard Source = "dashboard"
)

// passwordPurpose lie le chiffré du mot de passe à son champ.
const passwordPurpose = "mail_settings.password"

// ErrNotEditable est retourné quand le serveur n'a pas de clé pour chiffrer
// le mot de passe du relais : les réglages se changent alors par install.sh.
var ErrNotEditable = errors.New("ce serveur n'a pas de clé FORGENET_SECRET_KEY : le relais se règle par install.sh, qui la crée à la prochaine mise à jour")

// Settings décrit l'envoi en vigueur.
type Settings struct {
	// Config porte les valeurs, mot de passe compris : il ne quitte jamais
	// le serveur.
	Config config.SMTPConfig
	// Enabled : un relais est réglé et l'envoi n'est pas coupé.
	Enabled bool
	Source  Source
	// UpdatedBy et UpdatedAt datent l'enregistrement depuis le dashboard.
	UpdatedBy string
	UpdatedAt time.Time
	// Unreadable : des réglages du dashboard existent, mais leur mot de passe
	// ne s'ouvre plus avec la clé du serveur ; ceux de .env s'appliquent.
	Unreadable bool
}

// RelayOption règle un Relay.
type RelayOption func(*Relay)

// WithSenderFactory remplace l'envoi SMTP, pour les tests.
func WithSenderFactory(f func(config.SMTPConfig) Sender) RelayOption {
	return func(r *Relay) { r.newSender = f }
}

// Relay tient les réglages d'envoi en vigueur et les remplace sans
// redémarrage (ADR-0015) : ceux du dashboard s'ils existent, sinon ceux de
// .env. Chaque message lit ceux du moment où il part.
type Relay struct {
	server    config.SMTPConfig
	store     store.MailSettingsStore
	box       *sealed.Box
	newSender func(config.SMTPConfig) Sender

	mu  sync.RWMutex
	cur Settings
}

// NewRelay part des réglages de .env ; Load lit ensuite ceux de la base. st
// ou box nil : les réglages ne se changent pas depuis le dashboard.
func NewRelay(server config.SMTPConfig, st store.MailSettingsStore, box *sealed.Box, opts ...RelayOption) *Relay {
	r := &Relay{
		server: server, store: st, box: box,
		newSender: func(cfg config.SMTPConfig) Sender { return NewSMTP(cfg) },
	}
	for _, opt := range opts {
		opt(r)
	}
	r.cur = r.fromServer(false)
	return r
}

func (r *Relay) fromServer(unreadable bool) Settings {
	return Settings{Config: r.server, Enabled: r.server.Enabled(), Source: SourceServer, Unreadable: unreadable}
}

// Load lit les réglages enregistrés depuis le dashboard. Un mot de passe qui
// ne s'ouvre plus ne bloque pas le serveur : .env s'applique, et Current le
// signale. L'erreur ne concerne que la lecture de la base.
func (r *Relay) Load(ctx context.Context) error {
	if r.store == nil {
		return nil
	}
	m, err := r.store.GetMailSettings(ctx)
	if errors.Is(err, store.ErrNotFound) {
		r.set(r.fromServer(false))
		return nil
	}
	if err != nil {
		return err
	}
	r.set(r.fromStored(m))
	return nil
}

func (r *Relay) fromStored(m *store.MailSettings) Settings {
	cfg := config.SMTPConfig{Host: m.Host, Port: m.Port, TLS: m.TLS, Username: m.Username, From: m.From}
	if len(m.PasswordSealed) > 0 {
		if r.box == nil {
			return r.fromServer(true)
		}
		plain, err := r.box.Open(passwordPurpose, m.PasswordSealed)
		if err != nil {
			return r.fromServer(true)
		}
		cfg.Password = string(plain)
	}
	return Settings{
		Config: cfg, Enabled: !m.Disabled && cfg.Enabled(), Source: SourceDashboard,
		UpdatedBy: m.UpdatedBy, UpdatedAt: m.UpdatedAt,
	}
}

func (r *Relay) set(s Settings) {
	r.mu.Lock()
	r.cur = s
	r.mu.Unlock()
}

// Current renvoie les réglages en vigueur.
func (r *Relay) Current() Settings {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.cur
}

// Enabled dit si les messages partent.
func (r *Relay) Enabled() bool { return r.Current().Enabled }

// Editable dit si les réglages se changent depuis le dashboard.
func (r *Relay) Editable() bool { return r.store != nil && r.box != nil }

// Send remet msg au relais en vigueur à cet instant.
func (r *Relay) Send(ctx context.Context, msg Message) error {
	cur := r.Current()
	if !cur.Enabled {
		return ErrDisabled
	}
	return r.newSender(cur.Config).Send(ctx, msg)
}

// Try envoie msg par cfg sans rien changer : c'est l'essai qui précède un
// enregistrement.
func (r *Relay) Try(ctx context.Context, cfg config.SMTPConfig, msg Message) error {
	if !cfg.Enabled() {
		return ErrDisabled
	}
	return r.newSender(cfg).Send(ctx, msg)
}

// Save enregistre cfg, déjà validé et essayé, et l'applique aussitôt. by est
// l'identifiant du compte.
func (r *Relay) Save(ctx context.Context, cfg config.SMTPConfig, by string) error {
	return r.persist(ctx, cfg, false, by)
}

// Disable coupe l'envoi. Les valeurs en vigueur restent enregistrées, pour
// rouvrir l'envoi sans tout ressaisir.
func (r *Relay) Disable(ctx context.Context, by string) error {
	return r.persist(ctx, r.Current().Config, true, by)
}

func (r *Relay) persist(ctx context.Context, cfg config.SMTPConfig, disabled bool, by string) error {
	if !r.Editable() {
		return ErrNotEditable
	}
	m := &store.MailSettings{
		Disabled: disabled, Host: cfg.Host, Port: cfg.Port, TLS: cfg.TLS,
		Username: cfg.Username, From: cfg.From, UpdatedBy: by, UpdatedAt: time.Now().UTC(),
	}
	if cfg.Password != "" {
		ct, err := r.box.Seal(passwordPurpose, []byte(cfg.Password))
		if err != nil {
			return err
		}
		m.PasswordSealed = ct
	}
	// Le verrou couvre l'écriture : deux enregistrements simultanés ne
	// laissent pas la mémoire et la base en désaccord.
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := r.store.SaveMailSettings(ctx, m); err != nil {
		return err
	}
	r.cur = Settings{
		Config: cfg, Enabled: !disabled && cfg.Enabled(), Source: SourceDashboard,
		UpdatedBy: by, UpdatedAt: m.UpdatedAt,
	}
	return nil
}

// Reset efface les réglages du dashboard : ceux de .env s'appliquent de
// nouveau. Renvoie true s'il y en avait.
func (r *Relay) Reset(ctx context.Context) (bool, error) {
	if r.store == nil {
		return false, nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	gone, err := r.store.DeleteMailSettings(ctx)
	if err != nil {
		return false, fmt.Errorf("retour aux réglages du serveur : %w", err)
	}
	r.cur = r.fromServer(false)
	return gone, nil
}
