package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/mail"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// MailServiceServer montre l'envoi d'e-mails, l'essaie et le règle
// (ADR-0014, ADR-0015). Réservé aux administrateurs par la politique d'accès
// (aucune règle MinRole).
type MailServiceServer struct {
	forgenetv1.UnimplementedMailServiceServer
	relay   *mail.Relay
	users   UserLookup
	auditor audit.Auditor
	// tests freine les e-mails d'essai, par compte, qu'ils viennent du bouton
	// « Envoyer un e-mail d'essai » ou d'un enregistrement : un bouton cliqué
	// en boucle ne doit pas faire classer le relais en indésirable.
	tests  *loginLimiter
	logger zerolog.Logger
}

var _ forgenetv1.MailServiceServer = (*MailServiceServer)(nil)

// NewMailServiceServer instancie le service des e-mails.
func NewMailServiceServer(relay *mail.Relay, users UserLookup, auditor audit.Auditor, logger zerolog.Logger) *MailServiceServer {
	return &MailServiceServer{
		relay: relay, users: users, auditor: auditor,
		tests:  newLoginLimiter(time.Now),
		logger: logger.With().Str("component", "transport.grpc.mail").Logger(),
	}
}

// GetMailStatus décrit l'envoi, sans le mot de passe du relais.
func (s *MailServiceServer) GetMailStatus(ctx context.Context, _ *forgenetv1.GetMailStatusRequest) (*forgenetv1.GetMailStatusResponse, error) {
	return s.status(ctx), nil
}

func (s *MailServiceServer) status(ctx context.Context) *forgenetv1.GetMailStatusResponse {
	cur := s.relay.Current()
	res := &forgenetv1.GetMailStatusResponse{
		Enabled:            cur.Enabled,
		Editable:           s.relay.Editable(),
		Source:             mailSource(cur.Source),
		Disabled:           cur.Source == mail.SourceDashboard && !cur.Enabled && cur.Config.Enabled(),
		PasswordUnreadable: cur.Unreadable,
	}
	if cfg := cur.Config; cfg.Enabled() {
		res.Host = cfg.Host
		res.Port = uint32(cfg.Port) //nolint:gosec // port validé (1 à 65535) avant d'être retenu
		res.Tls = mailTLS(cfg.TLS)
		res.From = cfg.From
		res.Username = cfg.Username
		res.Authenticated = cfg.Password != ""
	}
	if !cur.UpdatedAt.IsZero() {
		res.UpdatedAt = timestamppb.New(cur.UpdatedAt)
		res.UpdatedBy = s.username(ctx, cur.UpdatedBy)
	}
	return res
}

// username nomme le compte qui a enregistré ; un compte supprimé depuis
// n'empêche pas de lire les réglages.
func (s *MailServiceServer) username(ctx context.Context, id string) string {
	if id == "" {
		return ""
	}
	u, err := s.users.GetUserByID(ctx, id)
	if err != nil {
		return ""
	}
	return u.Username
}

// SendTestMail envoie un e-mail d'essai à l'administrateur connecté et rend
// la réponse du relais telle quelle : c'est elle qui dit quoi corriger.
func (s *MailServiceServer) SendTestMail(ctx context.Context, _ *forgenetv1.SendTestMailRequest) (*forgenetv1.SendTestMailResponse, error) {
	cur := s.relay.Current()
	if !cur.Enabled {
		return nil, status.Error(codes.FailedPrecondition, mail.ErrDisabled.Error())
	}
	user, err := s.testRecipient(ctx)
	if err != nil {
		return nil, err
	}
	err = s.relay.Send(ctx, testMessage(user, cur.Config.Host))
	s.noteTest(ctx, user.ID, "mail.test", map[string]any{}, err)
	if err := relayRefusal(ctx, err, ""); err != nil {
		return nil, err
	}
	return &forgenetv1.SendTestMailResponse{To: user.Email}, nil
}

// SaveMailSettings essaie les nouveaux réglages, puis les enregistre : en
// cas de refus du relais, les anciens restent (ADR-0015).
func (s *MailServiceServer) SaveMailSettings(ctx context.Context, req *forgenetv1.SaveMailSettingsRequest) (*forgenetv1.SaveMailSettingsResponse, error) {
	if !s.relay.Editable() {
		return nil, status.Error(codes.FailedPrecondition, mail.ErrNotEditable.Error())
	}
	cfg, passwordChange, err := s.settingsFrom(req)
	if err != nil {
		return nil, err
	}
	user, err := s.testRecipient(ctx)
	if err != nil {
		return nil, err
	}

	err = s.relay.Try(ctx, cfg, testMessage(user, cfg.Host))
	meta := map[string]any{"host": cfg.Host}
	if err != nil {
		s.noteTest(ctx, user.ID, "mail.settings_change_failed", meta, err)
		return nil, relayRefusal(ctx, err, " ; les réglages en place restent")
	}
	if err := s.relay.Save(ctx, cfg, user.ID); err != nil {
		s.logger.Error().Err(err).Msg("enregistrement des réglages d'envoi impossible")
		return nil, status.Error(codes.Unavailable, "l'e-mail d'essai est parti, mais les réglages n'ont pas pu être enregistrés : réessayez")
	}
	s.logger.Info().Str("user_id", user.ID).Str("relais", cfg.Host).Msg("réglages d'envoi enregistrés depuis le dashboard")
	s.audit(ctx, map[string]any{
		"action": "enregistré", "host": cfg.Host, "port": cfg.Port, "tls": cfg.TLS, "password": passwordChange,
	})
	return &forgenetv1.SaveMailSettingsResponse{Status: s.status(ctx), To: user.Email}, nil
}

// settingsFrom lit et valide la demande. Le mot de passe gardé est celui en
// vigueur ; passwordChange dit au journal ce qu'il en est, sans sa valeur.
func (s *MailServiceServer) settingsFrom(req *forgenetv1.SaveMailSettingsRequest) (config.SMTPConfig, string, error) {
	tlsMode, ok := tlsFromProto(req.GetTls())
	if !ok {
		return config.SMTPConfig{}, "", status.Error(codes.InvalidArgument, "protection de la liaison inconnue")
	}
	cfg := config.SMTPConfig{
		Host:     strings.ToLower(strings.TrimSpace(req.GetHost())),
		Port:     int(req.GetPort()),
		TLS:      tlsMode,
		Username: strings.TrimSpace(req.GetUsername()),
		From:     strings.TrimSpace(req.GetFrom()),
	}
	if cfg.Host == "" {
		return cfg, "", status.Error(codes.InvalidArgument, "indiquez le relais : pour couper l'envoi, utilisez « Couper l'envoi »")
	}

	passwordChange := "aucun"
	switch {
	case cfg.Username == "":
		// Sans identifiant, pas d'authentification : un mot de passe gardé ou
		// saisi n'aurait rien à authentifier.
	case req.GetKeepPassword():
		cur := s.relay.Current()
		if cur.Unreadable || cur.Config.Password == "" {
			return cfg, "", status.Error(codes.InvalidArgument, "saisissez le mot de passe du relais : aucun n'est enregistré")
		}
		cfg.Password = cur.Config.Password
		passwordChange = "gardé"
	default:
		cfg.Password = req.GetPassword()
		passwordChange = "changé"
	}
	if err := cfg.Validate(); err != nil {
		return cfg, "", status.Error(codes.InvalidArgument, err.Error())
	}
	return cfg, passwordChange, nil
}

// DisableMail coupe l'envoi ; les valeurs restent pour le rouvrir.
func (s *MailServiceServer) DisableMail(ctx context.Context, _ *forgenetv1.DisableMailRequest) (*forgenetv1.DisableMailResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}
	err := s.relay.Disable(ctx, claims.UserID)
	switch {
	case errors.Is(err, mail.ErrNotEditable):
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	case err != nil:
		s.logger.Error().Err(err).Msg("coupure de l'envoi impossible")
		return nil, status.Error(codes.Unavailable, "coupure de l'envoi momentanément impossible")
	}
	s.logger.Info().Str("user_id", claims.UserID).Msg("envoi d'e-mails coupé depuis le dashboard")
	s.audit(ctx, map[string]any{"action": "coupé"})
	return &forgenetv1.DisableMailResponse{Status: s.status(ctx)}, nil
}

// ResetMailSettings rend la main aux réglages de .env.
func (s *MailServiceServer) ResetMailSettings(ctx context.Context, _ *forgenetv1.ResetMailSettingsRequest) (*forgenetv1.ResetMailSettingsResponse, error) {
	gone, err := s.relay.Reset(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("retour aux réglages du serveur impossible")
		return nil, status.Error(codes.Unavailable, "retour aux réglages du serveur momentanément impossible")
	}
	if gone {
		s.logger.Info().Msg("réglages d'envoi rendus au serveur (.env)")
		s.audit(ctx, map[string]any{"action": "rendu au serveur"})
	}
	return &forgenetv1.ResetMailSettingsResponse{Status: s.status(ctx)}, nil
}

// testRecipient est l'administrateur connecté, si le frein des essais le
// permet. Chaque essai compte, réussi ou non : c'est le relais qu'on ménage.
func (s *MailServiceServer) testRecipient(ctx context.Context) (*store.User, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}
	if s.tests.blocked(claims.UserID) {
		return nil, status.Error(codes.ResourceExhausted, "trop d'e-mails d'essai, réessayez dans quelques minutes")
	}
	user, err := s.users.GetUserByID(ctx, claims.UserID)
	if errors.Is(err, store.ErrNotFound) {
		return nil, status.Error(codes.Unauthenticated, "compte introuvable")
	}
	if err != nil {
		s.logger.Error().Err(err).Msg("lecture du compte impossible")
		return nil, status.Error(codes.Unavailable, "compte momentanément indisponible")
	}
	s.tests.fail(claims.UserID)
	return user, nil
}

func testMessage(user *store.User, host string) mail.Message {
	return mail.Message{
		To:      user.Email,
		Subject: "Essai d'envoi de ForgeNet",
		Body: fmt.Sprintf("Bonjour %s,\n\n"+
			"Cet e-mail confirme que votre serveur ForgeNet envoie bien ses messages par le relais %s.\n"+
			"Vous n'avez rien à faire.\n\n"+
			"— ForgeNet\n", user.Username, host),
	}
}

// noteTest inscrit l'essai au journal : le compte, le relais et l'issue,
// jamais l'adresse, le contenu ni les identifiants du relais.
func (s *MailServiceServer) noteTest(ctx context.Context, userID, action string, meta map[string]any, err error) {
	meta["result"] = "envoyé"
	if err != nil {
		meta["result"] = "refusé"
		s.logger.Warn().Err(err).Str("user_id", userID).Msg("e-mail d'essai non remis")
	} else {
		s.logger.Info().Str("user_id", userID).Msg("e-mail d'essai remis au relais")
	}
	if s.auditor != nil {
		_ = s.auditor.Log(ctx, action, "user:"+userID, meta)
	}
}

// audit inscrit un changement des réglages : jamais le mot de passe, ni
// l'identifiant du relais ou l'expéditeur, souvent des adresses (ADR-0015).
func (s *MailServiceServer) audit(ctx context.Context, meta map[string]any) {
	if s.auditor != nil {
		_ = s.auditor.Log(ctx, "mail.settings_change", "mail:relay", meta)
	}
}

// relayRefusal traduit l'issue d'un envoi d'essai ; nil s'il est parti.
func relayRefusal(ctx context.Context, err error, suffix string) error {
	switch {
	case err == nil:
		return nil
	// Seul l'abandon de l'appel lui-même (page quittée) garde son code : un
	// relais trop lent est une réponse à montrer, pas une panne du serveur.
	case ctx.Err() != nil:
		return status.FromContextError(ctx.Err()).Err()
	default:
		return status.Error(codes.FailedPrecondition, "le relais n'a pas pris le message : "+err.Error()+suffix)
	}
}

func mailSource(src mail.Source) forgenetv1.MailSource {
	switch src {
	case mail.SourceServer:
		return forgenetv1.MailSource_MAIL_SOURCE_SERVER
	case mail.SourceDashboard:
		return forgenetv1.MailSource_MAIL_SOURCE_DASHBOARD
	default:
		return forgenetv1.MailSource_MAIL_SOURCE_UNSPECIFIED
	}
}

func mailTLS(mode string) forgenetv1.MailTLS {
	switch mode {
	case config.SMTPStartTLS:
		return forgenetv1.MailTLS_MAIL_TLS_STARTTLS
	case config.SMTPImplicit:
		return forgenetv1.MailTLS_MAIL_TLS_IMPLICIT
	case config.SMTPNoTLS:
		return forgenetv1.MailTLS_MAIL_TLS_NONE
	default:
		return forgenetv1.MailTLS_MAIL_TLS_UNSPECIFIED
	}
}

func tlsFromProto(t forgenetv1.MailTLS) (string, bool) {
	switch t {
	case forgenetv1.MailTLS_MAIL_TLS_UNSPECIFIED, forgenetv1.MailTLS_MAIL_TLS_STARTTLS:
		return config.SMTPStartTLS, true
	case forgenetv1.MailTLS_MAIL_TLS_IMPLICIT:
		return config.SMTPImplicit, true
	case forgenetv1.MailTLS_MAIL_TLS_NONE:
		return config.SMTPNoTLS, true
	default:
		return "", false
	}
}
