package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcHealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

// Policy décrit qui peut appeler chaque méthode gRPC (nom complet,
// « /paquet.Service/Méthode »).
type Policy struct {
	// Public liste les méthodes appelables sans authentification.
	Public map[string]bool
	// MinRole fixe le rôle minimal d'une méthode authentifiée. Une méthode
	// absente exige le rôle admin : un service ajouté sans règle explicite
	// n'est jamais ouvert aux membres par oubli.
	MinRole map[string]store.Role
}

// DefaultPolicy est la politique du serveur : seules les sondes de santé sont
// publiques. Chaque service métier y déclare ses méthodes ouvertes aux membres.
func DefaultPolicy() Policy {
	return Policy{
		Public: map[string]bool{
			grpcHealth.Health_Check_FullMethodName: true,
			grpcHealth.Health_Watch_FullMethodName: true,
		},
		MinRole: map[string]store.Role{},
	}
}

func (p Policy) requiredRole(method string) store.Role {
	if role, ok := p.MinRole[method]; ok {
		return role
	}
	return store.RoleAdmin
}

// Guard applique l'authentification et le contrôle des rôles à chaque appel gRPC.
type Guard struct {
	auth   Service
	policy Policy
	logger zerolog.Logger
}

// NewGuard crée le garde gRPC du serveur.
func NewGuard(authSvc Service, policy Policy, logger zerolog.Logger) *Guard {
	return &Guard{auth: authSvc, policy: policy, logger: logger.With().Str("component", "auth.guard").Logger()}
}

// Unary renvoie l'intercepteur des appels simples.
func (g *Guard) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, err := g.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// Stream renvoie l'intercepteur des flux.
func (g *Guard) Stream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx, err := g.authorize(ss.Context(), info.FullMethod)
		if err != nil {
			return err
		}
		return handler(srv, &wrappedStream{ServerStream: ss, ctx: ctx})
	}
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// authorize renvoie le contexte enrichi de l'identité de l'appelant, ou une
// erreur gRPC au message générique : le détail (jeton expiré, base
// indisponible…) reste dans le journal du serveur.
func (g *Guard) authorize(ctx context.Context, method string) (context.Context, error) {
	if g.policy.Public[method] {
		return ctx, nil
	}

	claims, err := g.auth.ValidateToken(ctx, bearerToken(ctx))
	if err != nil {
		if isCredentialError(err) {
			return nil, status.Error(codes.Unauthenticated, "authentification requise ou jeton invalide")
		}
		g.logger.Error().Err(err).Str("grpc_method", method).Msg("vérification du jeton impossible")
		return nil, status.Error(codes.Unavailable, "vérification de l'authentification indisponible")
	}

	if !claims.Role.Grants(g.policy.requiredRole(method)) {
		return nil, status.Error(codes.PermissionDenied, "droits insuffisants pour cette opération")
	}

	clientIP := ""
	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		clientIP = p.Addr.String()
	}
	ctx = WithUserContext(ctx, claims)
	return audit.WithAuditContext(ctx, claims.UserID, clientIP), nil
}

func bearerToken(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	values := md.Get("authorization")
	if len(values) == 0 {
		return ""
	}
	token := strings.TrimSpace(values[0])
	if len(token) > len("bearer ") && strings.EqualFold(token[:len("bearer ")], "bearer ") {
		token = strings.TrimSpace(token[len("bearer "):])
	}
	return token
}

func isCredentialError(err error) bool {
	return errors.Is(err, ErrUnauthenticated) || errors.Is(err, ErrInvalidToken) || errors.Is(err, ErrExpiredToken)
}
