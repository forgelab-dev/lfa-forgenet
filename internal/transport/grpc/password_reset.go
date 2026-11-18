package grpc

import (
	"context"
	"errors"
	"net/netip"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/mail"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

const (
	// resetsPerAccount demandes par heure au plus pour un même compte ; les
	// suivantes sont ignorées sans le dire (ADR-0014).
	resetsPerAccount = 5
	// resetsPerAddress demandes par heure au plus depuis une même adresse
	// publique, tous comptes confondus.
	resetsPerAddress = 20
	resetWindow      = time.Hour
	// maxLoginLength borne l'identifiant reçu : une adresse e-mail fait au
	// plus 254 caractères.
	maxLoginLength = 254
)

// Enqueuer met un envoi en tâche de fond (mail.Outbox le satisfait).
type Enqueuer interface {
	Enqueue(job mail.Job) bool
}

// passwordReset tient ce dont « mot de passe oublié » a besoin : la file
// d'envoi, l'adresse du dashboard pour le lien, et les freins.
type passwordReset struct {
	outbox        Enqueuer
	dashboardHost string
	perAccount    *loginLimiter
	perAddress    *loginLimiter
}

// EnablePasswordReset active « mot de passe oublié » : à n'appeler que si un
// relais SMTP est configuré. dashboardHost est le domaine du dashboard, vide
// sans domaine (le lien passe alors par le tunnel SSH).
func (s *AuthServiceServer) EnablePasswordReset(outbox Enqueuer, dashboardHost string) *AuthServiceServer {
	s.reset = &passwordReset{
		outbox:        outbox,
		dashboardHost: strings.TrimSpace(dashboardHost),
		perAccount:    newRateLimiter(time.Now, resetsPerAccount, resetWindow),
		perAddress:    newRateLimiter(time.Now, resetsPerAddress, resetWindow),
	}
	return s
}

// GetLoginOptions dit si l'écran de connexion peut proposer « mot de passe
// oublié ».
func (s *AuthServiceServer) GetLoginOptions(context.Context, *forgenetv1.GetLoginOptionsRequest) (*forgenetv1.GetLoginOptionsResponse, error) {
	return &forgenetv1.GetLoginOptionsResponse{PasswordReset: s.reset != nil}, nil
}

// RequestPasswordReset met la demande en tâche de fond et répond tout de
// suite, toujours de la même façon : ni la réponse ni son délai ne disent si
// un compte correspond.
func (s *AuthServiceServer) RequestPasswordReset(ctx context.Context, req *forgenetv1.RequestPasswordResetRequest) (*forgenetv1.RequestPasswordResetResponse, error) {
	r := s.reset
	if r == nil {
		return nil, status.Error(codes.FailedPrecondition,
			"l'envoi d'e-mails n'est pas configuré sur ce serveur : un mot de passe perdu se remplace sur le serveur, par « forgenet user set-password --username NOM »")
	}
	login := strings.TrimSpace(req.GetLogin())
	if login == "" || len(login) > maxLoginLength {
		return nil, status.Error(codes.InvalidArgument, "saisissez votre nom d'utilisateur ou votre adresse e-mail")
	}
	// Derrière Caddy ou un tunnel SSH, toutes les demandes arrivent de la même
	// adresse privée : la freiner bloquerait tout le monde. Le frein par
	// compte, lui, s'applique toujours.
	if host := clientHost(ctx); isPublicAddress(host) {
		if r.perAddress.blocked(host) {
			return nil, status.Error(codes.ResourceExhausted, "trop de demandes, réessayez dans une heure")
		}
		r.perAddress.fail(host)
	}

	job := mail.Job{Kind: "password_reset", Build: func(ctx context.Context) (*mail.Message, string, error) {
		return s.buildPasswordReset(ctx, r, login)
	}}
	if !r.outbox.Enqueue(job) {
		return nil, status.Error(codes.Unavailable, "envoi momentanément impossible, réessayez dans quelques minutes")
	}
	return &forgenetv1.RequestPasswordResetResponse{}, nil
}

// buildPasswordReset tourne en tâche de fond : il cherche le compte, fait le
// lien et rédige l'e-mail. Sans compte, ou au-delà du frein, rien ne part.
func (s *AuthServiceServer) buildPasswordReset(ctx context.Context, r *passwordReset, login string) (*mail.Message, string, error) {
	user, err := s.auth.LookupAccount(ctx, login)
	if errors.Is(err, store.ErrNotFound) {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", err
	}
	if r.perAccount.blocked(user.ID) {
		s.logger.Warn().Str("user_id", user.ID).Msg("mot de passe oublié : demandes trop nombreuses, celle-ci est ignorée")
		return nil, user.ID, nil
	}
	r.perAccount.fail(user.ID)

	token, _, err := s.auth.CreateAccessLink(ctx, user.ID, store.PurposePasswordReset, "system:password_reset")
	if err != nil {
		return nil, user.ID, err
	}
	url := mail.AccessURL(mail.DashboardURL(r.dashboardHost), token)
	msg := mail.PasswordResetMessage(user.Email, user.Username, url, r.dashboardHost == "")
	return &msg, user.ID, nil
}

// isPublicAddress dit si host est une adresse IP joignable depuis Internet,
// et non celle d'un proxy local ou d'un tunnel.
func isPublicAddress(host string) bool {
	addr, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}
	addr = addr.Unmap()
	return !addr.IsLoopback() && !addr.IsPrivate() && !addr.IsLinkLocalUnicast() && !addr.IsUnspecified()
}
