package grpc

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/mail"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// usernamePattern borne les noms des comptes invités : lettres et chiffres
// ASCII, point, tiret et soulignement, de 2 à 64 caractères. Le nom sert à se
// connecter et figure dans les e-mails : rien d'ambigu ni d'invisible.
var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]{2,64}$`)

// TeamUsers lit les comptes (store.UserStore le satisfait).
type TeamUsers interface {
	ListUsers(ctx context.Context) ([]*store.User, error)
	GetUserByID(ctx context.Context, id string) (*store.User, error)
}

// TeamLinks lit les liens en cours (store.ActionTokenStore le satisfait).
type TeamLinks interface {
	ListActiveActionTokens(ctx context.Context, purpose store.ActionPurpose, now time.Time) ([]*store.ActionToken, error)
}

// TeamDeps regroupe ce dont le service de l'équipe a besoin.
type TeamDeps struct {
	Auth  auth.Service
	Users TeamUsers
	Links TeamLinks
	// Outbox envoie les invitations ; nil sans relais SMTP : le lien est
	// alors rendu à l'administrateur.
	Outbox Enqueuer
	// DashboardHost est le domaine du dashboard, vide sans domaine.
	DashboardHost string
	Auditor       audit.Auditor
	Now           func() time.Time
}

// TeamServiceServer gère les membres et leurs invitations (ADR-0014).
// Réservé aux administrateurs par la politique d'accès (aucune règle MinRole).
type TeamServiceServer struct {
	forgenetv1.UnimplementedTeamServiceServer
	deps   TeamDeps
	logger zerolog.Logger
}

var _ forgenetv1.TeamServiceServer = (*TeamServiceServer)(nil)

// NewTeamServiceServer instancie le service de l'équipe.
func NewTeamServiceServer(deps TeamDeps, logger zerolog.Logger) *TeamServiceServer {
	if deps.Now == nil {
		deps.Now = time.Now
	}
	deps.DashboardHost = strings.TrimSpace(deps.DashboardHost)
	return &TeamServiceServer{deps: deps, logger: logger.With().Str("component", "transport.grpc.team").Logger()}
}

// ListMembers liste les comptes avec leur état.
func (s *TeamServiceServer) ListMembers(ctx context.Context, _ *forgenetv1.ListMembersRequest) (*forgenetv1.ListMembersResponse, error) {
	users, err := s.deps.Users.ListUsers(ctx)
	if err != nil {
		return nil, s.unavailable(err, "liste des membres impossible")
	}
	links, err := s.deps.Links.ListActiveActionTokens(ctx, store.PurposeInvitation, s.deps.Now().UTC())
	if err != nil {
		return nil, s.unavailable(err, "lecture des invitations impossible")
	}
	// Un seul lien valable par compte : un nouveau annule les précédents.
	invitations := make(map[string]time.Time, len(links))
	for _, l := range links {
		invitations[l.UserID] = l.ExpiresAt
	}
	members := make([]*forgenetv1.Member, 0, len(users))
	for _, u := range users {
		expires, invited := invitations[u.ID]
		members = append(members, toMember(u, expires, invited))
	}
	return &forgenetv1.ListMembersResponse{Members: members, MailEnabled: s.deps.Outbox != nil}, nil
}

// InviteMember crée le compte sans mot de passe, puis l'invite.
func (s *TeamServiceServer) InviteMember(ctx context.Context, req *forgenetv1.InviteMemberRequest) (*forgenetv1.InviteMemberResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}
	username := strings.TrimSpace(req.GetUsername())
	if !usernamePattern.MatchString(username) {
		return nil, status.Error(codes.InvalidArgument,
			"le nom d'utilisateur fait de 2 à 64 caractères : lettres sans accent, chiffres, point, tiret ou soulignement")
	}
	email, err := auth.NormalizeEmail(req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	role, ok := fromUserRole(req.GetRole())
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "choisissez un rôle : administrateur ou membre")
	}

	user, err := s.deps.Auth.RegisterWithoutPassword(ctx, username, email, role)
	switch {
	case errors.Is(err, store.ErrConflict):
		return nil, status.Error(codes.AlreadyExists, "ce nom d'utilisateur ou cette adresse e-mail a déjà un compte")
	case err != nil:
		return nil, s.unavailable(err, "création du compte invité impossible")
	}
	result, err := s.invite(ctx, claims, user, "team.invite")
	if err != nil {
		return nil, err
	}
	return &forgenetv1.InviteMemberResponse{Result: result}, nil
}

// RenewInvitation refait le lien d'un compte qui attend son mot de passe.
func (s *TeamServiceServer) RenewInvitation(ctx context.Context, req *forgenetv1.RenewInvitationRequest) (*forgenetv1.RenewInvitationResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}
	user, err := s.deps.Users.GetUserByID(ctx, req.GetUserId())
	switch {
	case errors.Is(err, store.ErrNotFound):
		return nil, status.Error(codes.NotFound, "ce compte n'existe plus")
	case err != nil:
		return nil, s.unavailable(err, "lecture du compte impossible")
	case user.Suspended():
		return nil, status.Error(codes.FailedPrecondition, "ce compte est suspendu : réactivez-le avant de renvoyer son invitation")
	case auth.HasPassword(user):
		return nil, status.Error(codes.FailedPrecondition, "ce compte a déjà choisi son mot de passe : il n'a plus besoin d'invitation")
	}
	result, err := s.invite(ctx, claims, user, "team.invite_renew")
	if err != nil {
		return nil, err
	}
	return &forgenetv1.RenewInvitationResponse{Result: result}, nil
}

// invite fait le lien et l'envoie par e-mail, ou le rend s'il ne peut pas
// partir. Le lien n'est jamais journalisé.
func (s *TeamServiceServer) invite(ctx context.Context, claims *auth.Claims, user *store.User, action string) (*forgenetv1.InvitationResult, error) {
	token, link, err := s.deps.Auth.CreateAccessLink(ctx, user.ID, store.PurposeInvitation, "user:"+claims.UserID)
	if err != nil {
		return nil, s.unavailable(err, "création du lien d'invitation impossible")
	}
	local := s.deps.DashboardHost == ""
	url := mail.AccessURL(mail.DashboardURL(s.deps.DashboardHost), token)
	result := &forgenetv1.InvitationResult{Member: toMember(user, link.ExpiresAt, true), LocalDashboard: local}

	if s.deps.Outbox != nil {
		msg := mail.InvitationMessage(user.Email, user.Username, claims.Username, url, local)
		result.Emailed = s.deps.Outbox.Enqueue(mail.Job{Kind: "invitation", Build: func(context.Context) (*mail.Message, string, error) {
			return &msg, user.ID, nil
		}})
		if !result.Emailed {
			// File pleine : le lien revient à l'administrateur plutôt que
			// d'être perdu.
			s.logger.Warn().Str("user_id", user.ID).Msg("invitation non mise en file d'envoi : le lien est rendu à l'administrateur")
		}
	}
	if !result.Emailed {
		result.InvitationUrl = url
	}

	sent := "lien rendu"
	if result.Emailed {
		sent = "e-mail"
	}
	if s.deps.Auditor != nil {
		actx := audit.WithAuditContext(ctx, claims.UserID, audit.IPFromContext(ctx))
		_ = s.deps.Auditor.Log(actx, action, "user:"+user.ID, map[string]any{"role": string(user.Role), "sent": sent})
	}
	return result, nil
}

// SuspendMember suspend un compte et ferme ses sessions, jetons et liens.
func (s *TeamServiceServer) SuspendMember(ctx context.Context, req *forgenetv1.SuspendMemberRequest) (*forgenetv1.SuspendMemberResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}
	user, closed, err := s.deps.Auth.SuspendAccount(ctx, claims.UserID, req.GetUserId())
	if err != nil {
		return nil, s.accountError(err, "suspension du compte impossible")
	}
	return &forgenetv1.SuspendMemberResponse{Member: toMember(user, time.Time{}, false), ClosedSessions: closed}, nil
}

// ResumeMember réactive un compte suspendu.
func (s *TeamServiceServer) ResumeMember(ctx context.Context, req *forgenetv1.ResumeMemberRequest) (*forgenetv1.ResumeMemberResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}
	user, err := s.deps.Auth.ResumeAccount(ctx, claims.UserID, req.GetUserId())
	if err != nil {
		return nil, s.accountError(err, "réactivation du compte impossible")
	}
	// Un compte invité réactivé n'a plus de lien (la suspension l'a annulé) :
	// il apparaît « invitation expirée », à renvoyer.
	return &forgenetv1.ResumeMemberResponse{Member: toMember(user, time.Time{}, false)}, nil
}

// CloseMemberSessions ferme les sessions et jetons d'un compte, sans le
// suspendre. Sur son propre compte, la session de l'appel reste ouverte.
func (s *TeamServiceServer) CloseMemberSessions(ctx context.Context, req *forgenetv1.CloseMemberSessionsRequest) (*forgenetv1.CloseMemberSessionsResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}
	keep := ""
	if req.GetUserId() == claims.UserID {
		keep = claims.SessionID
	}
	closed, err := s.deps.Auth.CloseSessions(ctx, claims.UserID, req.GetUserId(), keep)
	if err != nil {
		return nil, s.accountError(err, "fermeture des sessions impossible")
	}
	return &forgenetv1.CloseMemberSessionsResponse{ClosedSessions: closed}, nil
}

// DeleteMember supprime un compte après confirmation par son nom.
func (s *TeamServiceServer) DeleteMember(ctx context.Context, req *forgenetv1.DeleteMemberRequest) (*forgenetv1.DeleteMemberResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}
	user, err := s.deps.Users.GetUserByID(ctx, req.GetUserId())
	if err != nil {
		return nil, s.accountError(err, "lecture du compte impossible")
	}
	// La confirmation est vérifiée ici, pas seulement dans le navigateur : un
	// appel direct à l'API doit aussi retaper le nom.
	if strings.TrimSpace(req.GetConfirmUsername()) != user.Username {
		return nil, status.Error(codes.InvalidArgument, "tapez exactement le nom du compte pour confirmer la suppression")
	}
	if _, err := s.deps.Auth.DeleteAccount(ctx, claims.UserID, user.ID); err != nil {
		return nil, s.accountError(err, "suppression du compte impossible")
	}
	return &forgenetv1.DeleteMemberResponse{}, nil
}

// accountError traduit les refus d'une action sur un compte.
func (s *TeamServiceServer) accountError(err error, logMsg string) error {
	switch {
	case errors.Is(err, store.ErrNotFound):
		return status.Error(codes.NotFound, "ce compte n'existe plus")
	case errors.Is(err, auth.ErrSelfAction):
		return status.Error(codes.FailedPrecondition, "vous ne pouvez pas suspendre ni supprimer votre propre compte : un autre administrateur doit le faire")
	case errors.Is(err, store.ErrLastAdmin):
		return status.Error(codes.FailedPrecondition, "c'est le dernier administrateur actif : nommez-en un autre avant")
	}
	return s.unavailable(err, logMsg)
}

func (s *TeamServiceServer) unavailable(err error, logMsg string) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return status.FromContextError(err).Err()
	}
	s.logger.Error().Err(err).Msg(logMsg)
	return status.Error(codes.Unavailable, "équipe momentanément indisponible, réessayez")
}

func toMember(u *store.User, invitationExpires time.Time, invited bool) *forgenetv1.Member {
	m := &forgenetv1.Member{
		Id: u.ID, Username: u.Username, Email: u.Email, Role: toUserRole(u.Role),
		CreatedAt: timestamppb.New(u.CreatedAt),
	}
	switch {
	case u.Suspended():
		m.Status = forgenetv1.MemberStatus_MEMBER_STATUS_SUSPENDED
		m.SuspendedAt = timestamppb.New(*u.SuspendedAt)
	case auth.HasPassword(u):
		m.Status = forgenetv1.MemberStatus_MEMBER_STATUS_ACTIVE
	case invited:
		m.Status = forgenetv1.MemberStatus_MEMBER_STATUS_INVITED
		m.InvitationExpiresAt = timestamppb.New(invitationExpires)
	default:
		m.Status = forgenetv1.MemberStatus_MEMBER_STATUS_INVITATION_EXPIRED
	}
	return m
}

func fromUserRole(r forgenetv1.UserRole) (store.Role, bool) {
	switch r {
	case forgenetv1.UserRole_USER_ROLE_ADMIN:
		return store.RoleAdmin, true
	case forgenetv1.UserRole_USER_ROLE_MEMBER:
		return store.RoleMember, true
	default:
		return "", false
	}
}
