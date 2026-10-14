// Package store définit les interfaces et modèles de persistance de ForgeNet.
// Ces interfaces abstraient le stockage des utilisateurs, sessions, RBAC et journal d'audit,
// permettant une implémentation SQLite par défaut et une extension ultérieure vers PostgreSQL sans impacter le domaine.
package store

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrNotFound est retourné lorsqu'une entité recherchée n'existe pas en base.
	ErrNotFound = errors.New("ressource non trouvée")
	// ErrConflict est retourné lorsqu'une contrainte d'unicité est violée (ex: username ou email déjà existant).
	ErrConflict = errors.New("conflit d'unicité sur la ressource")
)

// Role définit les permissions RBAC d'un utilisateur ou d'un token.
type Role string

const (
	// RoleAdmin donne un accès complet de gestion (orchestration, utilisateurs, proxy, secrets).
	RoleAdmin Role = "admin"
	// RoleMember donne un accès opérationnel restreint (lecture, statut, déploiement limité).
	RoleMember Role = "member"
)

// IsValid vérifie que la valeur du rôle correspond à un rôle autorisé.
func (r Role) IsValid() bool {
	return r == RoleAdmin || r == RoleMember
}

// User représente un compte utilisateur dans le système d'authentification.
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Session représente une session active ou un token API persistant.
type Session struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	TokenHash  string     `json:"-"`
	Name       string     `json:"name"`
	Role       Role       `json:"role"`
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	Revoked    bool       `json:"revoked"`
}

// IsExpired indique si la session a dépassé sa date d'expiration.
func (s *Session) IsExpired() bool {
	return time.Now().UTC().After(s.ExpiresAt)
}

// IsActive indique si la session est valide et non révoquée.
func (s *Session) IsActive() bool {
	return !s.Revoked && !s.IsExpired()
}

// AuditEvent représente une entrée dans le journal d'audit horodaté des actions sensibles.
type AuditEvent struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`   // ex: "container.start", "user.create", "domain.bind"
	Resource  string    `json:"resource"` // ex: "container:redis-prod", "user:admin"
	IPAddress string    `json:"ip_address,omitempty"`
	Metadata  string    `json:"metadata,omitempty"` // Détails structurés JSON
	Timestamp time.Time `json:"timestamp"`
}

// UserStore définit les opérations de persistance sur les utilisateurs.
type UserStore interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id string) error
	CountUsers(ctx context.Context) (int64, error)
}

// SessionStore définit les opérations de persistance sur les sessions et tokens API.
type SessionStore interface {
	CreateSession(ctx context.Context, session *Session) error
	GetSessionByID(ctx context.Context, id string) (*Session, error)
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (*Session, error)
	TouchSession(ctx context.Context, id string, lastUsed time.Time) error
	RevokeSession(ctx context.Context, id string) error
	DeleteExpiredSessions(ctx context.Context) error
}

// AuditStore définit les opérations de journalisation d'audit.
type AuditStore interface {
	Record(ctx context.Context, event *AuditEvent) error
	List(ctx context.Context, limit, offset int) ([]*AuditEvent, error)
}

// Store est l'interface globale du registre de persistance.
type Store interface {
	Users() UserStore
	Sessions() SessionStore
	Audit() AuditStore
	Ping(ctx context.Context) error
	Close() error
}
