package grpc

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

const (
	// maxLoginFailures échecs dans loginWindow bloquent la paire (adresse,
	// identifiant) jusqu'à la fin de la fenêtre.
	maxLoginFailures = 10
	loginWindow      = 15 * time.Minute
)

// UserLookup lit un compte (store.UserStore le satisfait).
type UserLookup interface {
	GetUserByID(ctx context.Context, id string) (*store.User, error)
}

// AuthServiceServer ouvre et ferme les sessions selon le contrat
// forgenet.v1.AuthService. Login est public dans la politique d'accès ; les
// autres méthodes exigent un jeton, vérifié en amont par auth.Guard.
type AuthServiceServer struct {
	forgenetv1.UnimplementedAuthServiceServer
	auth    auth.Service
	users   UserLookup
	limiter *loginLimiter
	logger  zerolog.Logger
}

var _ forgenetv1.AuthServiceServer = (*AuthServiceServer)(nil)

// NewAuthServiceServer instancie le service des sessions.
func NewAuthServiceServer(authSvc auth.Service, users UserLookup, logger zerolog.Logger) *AuthServiceServer {
	return &AuthServiceServer{
		auth:    authSvc,
		users:   users,
		limiter: newLoginLimiter(time.Now),
		logger:  logger.With().Str("component", "transport.grpc.auth").Logger(),
	}
}

// Login ouvre une session. Le refus ne dit pas si le compte existe.
func (s *AuthServiceServer) Login(ctx context.Context, req *forgenetv1.LoginRequest) (*forgenetv1.LoginResponse, error) {
	key := clientHost(ctx) + "|" + strings.ToLower(strings.TrimSpace(req.GetLogin()))
	if s.limiter.blocked(key) {
		return nil, status.Error(codes.ResourceExhausted, "trop d'échecs de connexion, réessayez dans quelques minutes")
	}

	res, err := s.auth.Login(ctx, req.GetLogin(), req.GetPassword())
	switch {
	case errors.Is(err, auth.ErrInvalidCredentials):
		s.limiter.fail(key)
		return nil, status.Error(codes.Unauthenticated, "identifiant ou mot de passe incorrect")
	case err != nil:
		return nil, s.unavailable(err, "ouverture de session impossible", "connexion momentanément indisponible")
	}
	s.limiter.reset(key)

	return &forgenetv1.LoginResponse{
		Token:     res.Token,
		ExpiresAt: timestamppb.New(res.ExpiresAt),
		User:      toUser(res.User, res.User.Role),
	}, nil
}

// Logout révoque la session du jeton présenté.
func (s *AuthServiceServer) Logout(ctx context.Context, _ *forgenetv1.LogoutRequest) (*forgenetv1.LogoutResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok || claims.SessionID == "" {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}
	if err := s.auth.RevokeSession(ctx, claims.SessionID); err != nil {
		return nil, s.unavailable(err, "révocation de session impossible", "déconnexion momentanément indisponible")
	}
	return &forgenetv1.LogoutResponse{}, nil
}

// GetCurrentUser renvoie le compte du jeton, avec le rôle retenu par la
// session (le plus bas entre celui du jeton et celui du compte).
func (s *AuthServiceServer) GetCurrentUser(ctx context.Context, _ *forgenetv1.GetCurrentUserRequest) (*forgenetv1.GetCurrentUserResponse, error) {
	claims, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentification requise")
	}
	user, err := s.users.GetUserByID(ctx, claims.UserID)
	if errors.Is(err, store.ErrNotFound) {
		return nil, status.Error(codes.Unauthenticated, "compte introuvable")
	}
	if err != nil {
		return nil, s.unavailable(err, "lecture du compte impossible", "compte momentanément indisponible")
	}
	return &forgenetv1.GetCurrentUserResponse{User: toUser(user, claims.Role)}, nil
}

// unavailable traduit une panne de la base en UNAVAILABLE, journalisée. Un
// appel abandonné par le client (page quittée, onglet fermé) n'est pas une
// panne : il garde son code (CANCELED, DEADLINE_EXCEEDED) et n'encombre pas le
// journal des erreurs.
func (s *AuthServiceServer) unavailable(err error, logMsg, clientMsg string) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return status.FromContextError(err).Err()
	}
	s.logger.Error().Err(err).Msg(logMsg)
	return status.Error(codes.Unavailable, clientMsg)
}

func toUser(u *store.User, role store.Role) *forgenetv1.User {
	return &forgenetv1.User{Id: u.ID, Username: u.Username, Email: u.Email, Role: toUserRole(role)}
}

func toUserRole(r store.Role) forgenetv1.UserRole {
	switch r {
	case store.RoleAdmin:
		return forgenetv1.UserRole_USER_ROLE_ADMIN
	case store.RoleMember:
		return forgenetv1.UserRole_USER_ROLE_MEMBER
	default:
		return forgenetv1.UserRole_USER_ROLE_UNSPECIFIED
	}
}

// clientHost renvoie l'adresse du client sans le port. Derrière un reverse
// proxy, c'est celle du proxy : la limite porte alors sur l'identifiant.
func clientHost(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok || p.Addr == nil {
		return ""
	}
	addr := p.Addr.String()
	if host, _, err := net.SplitHostPort(addr); err == nil {
		return host
	}
	return addr
}

// loginLimiter compte les échecs de connexion par clé, sur une fenêtre fixe
// ouverte au premier échec. En mémoire : un redémarrage remet les compteurs à
// zéro, ce qui reste bien plus lent qu'une attaque par dictionnaire.
type loginLimiter struct {
	now func() time.Time

	mu       sync.Mutex
	failures map[string]*loginFailures
}

type loginFailures struct {
	count int
	since time.Time
}

func newLoginLimiter(now func() time.Time) *loginLimiter {
	return &loginLimiter{now: now, failures: make(map[string]*loginFailures)}
}

func (l *loginLimiter) blocked(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, ok := l.failures[key]
	if !ok {
		return false
	}
	if l.now().Sub(f.since) >= loginWindow {
		delete(l.failures, key)
		return false
	}
	return f.count >= maxLoginFailures
}

func (l *loginLimiter) fail(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	// Les fenêtres échues sont oubliées ici : la table ne grossit qu'avec
	// les clés en échec depuis moins de loginWindow.
	for k, f := range l.failures {
		if now.Sub(f.since) >= loginWindow {
			delete(l.failures, k)
		}
	}
	f, ok := l.failures[key]
	if !ok {
		f = &loginFailures{since: now}
		l.failures[key] = f
	}
	f.count++
}

func (l *loginLimiter) reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.failures, key)
}
