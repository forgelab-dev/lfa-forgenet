package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/mail"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// MailServiceServer montre l'envoi d'e-mails et l'essaie (ADR-0014). Réservé
// aux administrateurs par la politique d'accès (aucune règle MinRole).
type MailServiceServer struct {
	forgenetv1.UnimplementedMailServiceServer
	cfg     config.SMTPConfig
	sender  mail.Sender
	users   UserLookup
	auditor audit.Auditor
	// tests freine les e-mails d'essai, par compte : un bouton cliqué en
	// boucle ne doit pas faire classer le relais en indésirable.
	tests  *loginLimiter
	logger zerolog.Logger
}

var _ forgenetv1.MailServiceServer = (*MailServiceServer)(nil)

// NewMailServiceServer instancie le service des e-mails.
func NewMailServiceServer(cfg config.SMTPConfig, sender mail.Sender, users UserLookup, auditor audit.Auditor, logger zerolog.Logger) *MailServiceServer {
	return &MailServiceServer{
		cfg: cfg, sender: sender, users: users, auditor: auditor,
		tests:  newLoginLimiter(time.Now),
		logger: logger.With().Str("component", "transport.grpc.mail").Logger(),
	}
}

// GetMailStatus décrit le relais, sans ses identifiants.
func (s *MailServiceServer) GetMailStatus(context.Context, *forgenetv1.GetMailStatusRequest) (*forgenetv1.GetMailStatusResponse, error) {
	if !s.cfg.Enabled() {
		return &forgenetv1.GetMailStatusResponse{}, nil
	}
	return &forgenetv1.GetMailStatusResponse{
		Enabled:       true,
		Host:          s.cfg.Host,
		Port:          uint32(s.cfg.Port), //nolint:gosec // port validé (1 à 65535) au chargement
		Tls:           mailTLS(s.cfg.TLS),
		From:          s.cfg.From,
		Authenticated: s.cfg.Username != "",
	}, nil
}

// SendTestMail envoie un e-mail d'essai à l'administrateur connecté et rend
// la réponse du relais telle quelle : c'est elle qui dit quoi corriger.
func (s *MailServiceServer) SendTestMail(ctx context.Context, _ *forgenetv1.SendTestMailRequest) (*forgenetv1.SendTestMailResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}
	if !s.cfg.Enabled() {
		return nil, status.Error(codes.FailedPrecondition, mail.ErrDisabled.Error())
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

	// Chaque essai compte, réussi ou non : c'est le relais qu'on ménage.
	s.tests.fail(claims.UserID)
	err = s.sender.Send(ctx, mail.Message{
		To:      user.Email,
		Subject: "Essai d'envoi de ForgeNet",
		Body: fmt.Sprintf("Bonjour %s,\n\n"+
			"Cet e-mail confirme que votre serveur ForgeNet envoie bien ses messages par le relais %s.\n"+
			"Vous n'avez rien à faire.\n\n"+
			"— ForgeNet\n", user.Username, s.cfg.Host),
	})
	// Le journal note le compte et l'issue, jamais l'adresse ni le contenu.
	meta := map[string]any{"result": "envoyé"}
	if err != nil {
		meta["result"] = "refusé"
		s.logger.Warn().Err(err).Str("user_id", user.ID).Msg("e-mail d'essai non remis")
	} else {
		s.logger.Info().Str("user_id", user.ID).Msg("e-mail d'essai remis au relais")
	}
	if s.auditor != nil {
		_ = s.auditor.Log(ctx, "mail.test", "user:"+user.ID, meta)
	}
	switch {
	// Seul l'abandon de l'appel lui-même (page quittée) garde son code : un
	// relais trop lent est une réponse à montrer, pas une panne du serveur.
	case ctx.Err() != nil:
		return nil, status.FromContextError(ctx.Err()).Err()
	case err != nil:
		return nil, status.Error(codes.FailedPrecondition, "le relais n'a pas pris le message : "+err.Error())
	}
	return &forgenetv1.SendTestMailResponse{To: user.Email}, nil
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
