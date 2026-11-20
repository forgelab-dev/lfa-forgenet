package grpc

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/mail"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// ChangeEmail change l'adresse du compte connecté (#100). Comme pour
// ChangePassword, les refus répondent INVALID_ARGUMENT, jamais
// UNAUTHENTICATED, et les essais du mot de passe partagent le même frein.
func (s *AuthServiceServer) ChangeEmail(ctx context.Context, req *forgenetv1.ChangeEmailRequest) (*forgenetv1.ChangeEmailResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}
	key := claims.UserID
	if s.passwords.blocked(key) {
		return nil, status.Error(codes.ResourceExhausted, "trop d'essais du mot de passe actuel, réessayez dans quelques minutes")
	}

	change, err := s.auth.ChangeOwnEmail(ctx, claims.UserID, req.GetCurrentPassword(), req.GetNewEmail())
	switch {
	case errors.Is(err, auth.ErrWrongCurrentPassword):
		s.passwords.fail(key)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	case err != nil:
		if refusal := emailChangeError(err); refusal != nil {
			return nil, refusal
		}
		if errors.Is(err, store.ErrNotFound) {
			return nil, status.Error(codes.Unauthenticated, "compte introuvable")
		}
		return nil, s.unavailable(err, "changement d'adresse impossible", "changement d'adresse momentanément indisponible")
	}
	s.passwords.reset(key)

	if s.reset != nil {
		sendEmailChangeNotice(s.reset.outbox, change, change.User.Username)
	}
	return &forgenetv1.ChangeEmailResponse{User: toUser(change.User, claims.Role)}, nil
}

// UpdateMemberEmail corrige l'adresse d'un compte (#100). Un compte qui
// attendait son mot de passe est réinvité à la nouvelle adresse ; un compte
// actif reçoit un avis à l'ancienne.
func (s *TeamServiceServer) UpdateMemberEmail(ctx context.Context, req *forgenetv1.UpdateMemberEmailRequest) (*forgenetv1.UpdateMemberEmailResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}
	// Sa propre adresse se change depuis la page Compte, avec le mot de
	// passe : une session d'administrateur volée ne doit pas y suffire.
	if req.GetUserId() == claims.UserID {
		return nil, status.Error(codes.FailedPrecondition, "changez votre propre adresse depuis la page Compte : votre mot de passe y est demandé")
	}

	change, err := s.deps.Auth.SetAccountEmail(ctx, claims.UserID, req.GetUserId(), req.GetEmail())
	if err != nil {
		if refusal := emailChangeError(err); refusal != nil {
			return nil, refusal
		}
		return nil, s.accountError(err, "changement d'adresse impossible")
	}
	user := change.User
	res := &forgenetv1.UpdateMemberEmailResponse{Member: toMember(user, time.Time{}, false)}

	switch {
	case !change.HadPassword && !user.Suspended():
		// Son lien est parti avec l'ancienne adresse : l'invitation repart à
		// la nouvelle.
		invitation, err := s.invite(ctx, claims, user, "team.invite_renew")
		if err != nil {
			return nil, err
		}
		res.Invitation = invitation
		res.Member = invitation.GetMember()
	case change.HadPassword && s.deps.Outbox != nil:
		sendEmailChangeNotice(s.deps.Outbox, change, claims.Username)
	}
	return res, nil
}

// emailChangeError traduit les refus propres à une adresse ; nil pour toute
// autre erreur.
func emailChangeError(err error) error {
	switch {
	case errors.Is(err, auth.ErrInvalidEmail), errors.Is(err, auth.ErrSameEmail):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, store.ErrConflict):
		return status.Error(codes.AlreadyExists, "cette adresse e-mail a déjà un compte")
	}
	return nil
}

// sendEmailChangeNotice prévient l'ancienne adresse d'un compte actif. Un
// compte qui attendait son mot de passe n'en reçoit pas : son ancienne
// adresse, peut-être mal saisie, n'a jamais servi qu'à l'invitation. actor
// est l'auteur du changement (le nom du compte s'il l'a fait lui-même).
func sendEmailChangeNotice(outbox Enqueuer, change *auth.EmailChange, actor string) {
	if !change.HadPassword {
		return
	}
	msg := mail.EmailChangedMessage(change.OldEmail, change.User.Username, change.User.Email, actor)
	userID := change.User.ID
	outbox.Enqueue(mail.Job{Kind: "email_change", Build: func(context.Context) (*mail.Message, string, error) {
		return &msg, userID, nil
	}})
}
