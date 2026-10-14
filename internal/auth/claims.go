// Package auth gère l'authentification (mots de passe, JWT, tokens API) et le contrôle d'accès RBAC.
package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

var (
	// ErrUnauthenticated est retourné lorsqu'aucun jeton valide n'est présent dans le contexte.
	ErrUnauthenticated = errors.New("authentification requise")
	// ErrForbidden est retourné lorsqu'une opération exige un rôle plus élevé que celui de l'appelant.
	ErrForbidden = errors.New("accès interdit : privilèges insuffisants")
)

type contextKey string

const (
	claimsContextKey contextKey = "forgenet_auth_claims"
)

// Claims modélise les données d'identité portées par un jeton JWT ou API.
type Claims struct {
	UserID    string     `json:"user_id"`
	Username  string     `json:"username"`
	Role      store.Role `json:"role"`
	SessionID string     `json:"session_id,omitempty"`
	jwt.RegisteredClaims
}

// WithUserContext attache les Claims d'authentification au contexte d'exécution.
func WithUserContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

// UserFromContext extrait les informations de l'utilisateur connecté du contexte.
func UserFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	return claims, ok && claims != nil
}
