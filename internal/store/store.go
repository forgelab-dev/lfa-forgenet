// Package store définit les interfaces et modèles de persistance de ForgeNet.
// Ces interfaces abstraient le stockage des utilisateurs, sessions, RBAC, journal d'audit et déploiements,
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
	// ErrInvalidTransition est retourné quand la machine à états d'un déploiement interdit le changement demandé.
	ErrInvalidTransition = errors.New("transition d'état interdite")
)

// Role définit les permissions RBAC d'un utilisateur ou d'un token.
type Role string

const (
	// RoleAdmin donne un accès complet de gestion (orchestration, utilisateurs, proxy, secrets).
	RoleAdmin Role = "admin"
	// RoleMember donne un accès opérationnel restreint (lecture, statut, déploiement limité).
	RoleMember Role = "member"
)

// roleRank ordonne les rôles : un rôle accorde tous les droits des rangs inférieurs.
var roleRank = map[Role]int{RoleMember: 1, RoleAdmin: 2}

// IsValid vérifie que la valeur du rôle correspond à un rôle autorisé.
func (r Role) IsValid() bool {
	_, ok := roleRank[r]
	return ok
}

// Grants indique si le rôle r suffit pour une action qui exige required.
// Un rôle inconnu n'accorde rien.
func (r Role) Grants(required Role) bool {
	rank, ok := roleRank[r]
	return ok && rank >= roleRank[required] && required.IsValid()
}

// Lowest renvoie le moins élevé des deux rôles : un jeton ne peut pas
// donner plus de droits que son propriétaire n'en a aujourd'hui.
func Lowest(a, b Role) Role {
	if roleRank[a] <= roleRank[b] {
		return a
	}
	return b
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

// DeploymentSourceType indique d'où vient ce qu'un déploiement met en service.
type DeploymentSourceType string

const (
	// SourceCompose déploie une stack décrite par un fichier Docker Compose.
	SourceCompose DeploymentSourceType = "compose"
	// SourceGit construit puis déploie le contenu d'un dépôt Git.
	SourceGit DeploymentSourceType = "git"
	// SourceImage déploie une image Docker déjà construite.
	SourceImage DeploymentSourceType = "image"
)

// IsValid vérifie que le type de source est connu.
func (t DeploymentSourceType) IsValid() bool {
	switch t {
	case SourceCompose, SourceGit, SourceImage:
		return true
	}
	return false
}

// DeploymentState est l'état d'un déploiement dans son cycle de vie.
//
// Le store a ses propres valeurs plutôt que l'énumération protobuf : la
// persistance ne dépend pas du contrat réseau, que la couche transport traduit.
type DeploymentState string

// États du cycle de vie, dans l'ordre de deploy.proto.
const (
	DeploymentPending        DeploymentState = "pending"
	DeploymentBuilding       DeploymentState = "building"
	DeploymentDeploying      DeploymentState = "deploying"
	DeploymentHealthChecking DeploymentState = "health_checking"
	DeploymentSuccess        DeploymentState = "success"
	DeploymentFailed         DeploymentState = "failed"
	DeploymentCancelled      DeploymentState = "cancelled"
	DeploymentRollingBack    DeploymentState = "rolling_back"
	DeploymentRolledBack     DeploymentState = "rolled_back"
)

// deploymentTransitions reprend les transitions autorisées par deploy.proto.
// Un état absent des clés est final.
var deploymentTransitions = map[DeploymentState][]DeploymentState{
	DeploymentPending:        {DeploymentBuilding, DeploymentCancelled},
	DeploymentBuilding:       {DeploymentDeploying, DeploymentFailed, DeploymentCancelled},
	DeploymentDeploying:      {DeploymentHealthChecking, DeploymentFailed, DeploymentCancelled},
	DeploymentHealthChecking: {DeploymentSuccess, DeploymentFailed, DeploymentRollingBack},
	DeploymentRollingBack:    {DeploymentRolledBack, DeploymentFailed},
}

// IsValid vérifie que l'état est l'un des états du cycle de vie.
func (s DeploymentState) IsValid() bool {
	switch s {
	case DeploymentPending, DeploymentBuilding, DeploymentDeploying, DeploymentHealthChecking,
		DeploymentSuccess, DeploymentFailed, DeploymentCancelled, DeploymentRollingBack, DeploymentRolledBack:
		return true
	}
	return false
}

// IsTerminal indique qu'aucune transition ne part plus de cet état.
func (s DeploymentState) IsTerminal() bool {
	_, hasNext := deploymentTransitions[s]
	return s.IsValid() && !hasNext
}

// CanTransitionTo indique si le passage de s à next est autorisé.
func (s DeploymentState) CanTransitionTo(next DeploymentState) bool {
	for _, allowed := range deploymentTransitions[s] {
		if allowed == next {
			return true
		}
	}
	return false
}

// Deployment représente une opération de déploiement et son état.
type Deployment struct {
	ID         string               `json:"id"`
	ProjectID  string               `json:"project_id"`
	SourceType DeploymentSourceType `json:"source_type"`
	// SourceContent garde la source telle que reçue : le fichier Compose, ou la
	// description d'une source Git ou image.
	SourceContent        string          `json:"source_content"`
	State                DeploymentState `json:"state"`
	ErrorMessage         string          `json:"error_message,omitempty"`
	TargetEnv            string          `json:"target_env"`
	PreviousDeploymentID string          `json:"previous_deployment_id,omitempty"`
	TriggeredBy          string          `json:"triggered_by"` // ex: "user:123", "system:auto-rollback"
	// RequestID est la clé d'idempotence fournie par le client ; vide si aucune.
	RequestID  string     `json:"request_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
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

// DeploymentStore définit les opérations de persistance sur les déploiements.
type DeploymentStore interface {
	// CreateDeployment renvoie ErrConflict si la clé d'idempotence est déjà prise.
	CreateDeployment(ctx context.Context, d *Deployment) error
	GetDeployment(ctx context.Context, id string) (*Deployment, error)
	GetDeploymentByRequestID(ctx context.Context, requestID string) (*Deployment, error)
	// ListDeployments liste du plus récent au plus ancien ; projectID vide liste tout.
	ListDeployments(ctx context.Context, projectID string, limit, offset int) ([]*Deployment, error)
	// TransitionDeployment fait passer le déploiement de from à to, avec errMsg
	// pour un échec. Il renvoie ErrInvalidTransition si la machine à états
	// l'interdit, et ErrConflict si l'état n'est plus from (une annulation et
	// l'exécuteur qui avancent en même temps).
	TransitionDeployment(ctx context.Context, id string, from, to DeploymentState, errMsg string) (*Deployment, error)
}

// Store est l'interface globale du registre de persistance.
type Store interface {
	Users() UserStore
	Sessions() SessionStore
	Audit() AuditStore
	Deployments() DeploymentStore
	Ping(ctx context.Context) error
	Close() error
}
