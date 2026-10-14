package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/ids"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

// ErrInvalidRole est retourné lorsqu'un compte est créé sans rôle reconnu.
var ErrInvalidRole = errors.New("rôle inconnu (valeurs acceptées : admin, member)")

// Result regroupe les données retournées lors d'une connexion réussie.
type Result struct {
	Token     string      `json:"token"`
	User      *store.User `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// Service définit les opérations d'authentification et de gestion de sessions.
type Service interface {
	Register(ctx context.Context, username, email, password string, role store.Role) (*store.User, error)
	Login(ctx context.Context, usernameOrEmail, password string) (*Result, error)
	CreateAPIToken(ctx context.Context, userID, name string, role store.Role, ttl time.Duration) (string, *store.Session, error)
	ValidateToken(ctx context.Context, tokenString string) (*Claims, error)
	RevokeSession(ctx context.Context, sessionID string) error
}

type authService struct {
	store   store.Store
	cfg     config.AuthConfig
	auditor audit.Auditor
}

// NewService instancie le service d'authentification. Sans secret configuré
// (développement, tests), un secret aléatoire de 256 bits est tiré : les
// sessions ne survivent alors pas à un redémarrage. La configuration impose un
// secret en production.
func NewService(st store.Store, cfg config.AuthConfig, auditor audit.Auditor) Service {
	if cfg.JWTSecret == "" {
		b := make([]byte, 32)
		_, _ = rand.Read(b)
		cfg.JWTSecret = hex.EncodeToString(b)
	}
	if cfg.TokenTTL <= 0 {
		cfg.TokenTTL = config.DefaultTokenTTL
	}
	if cfg.APITokenTTL <= 0 {
		cfg.APITokenTTL = config.DefaultAPITokenTTL
	}
	return &authService{store: st, cfg: cfg, auditor: auditor}
}

// Register crée un compte avec le rôle demandé. Aucun rôle n'est accordé
// implicitement : le premier administrateur est créé par la commande
// « forgenet-server user create » sur la machine, pas par la première
// personne qui atteint l'API.
func (s *authService) Register(ctx context.Context, username, email, password string, role store.Role) (*store.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(strings.ToLower(email))
	switch {
	case username == "":
		return nil, errors.New("le nom d'utilisateur est obligatoire")
	case email == "":
		return nil, errors.New("l'adresse email est obligatoire")
	case !role.IsValid():
		return nil, ErrInvalidRole
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &store.User{
		ID:           ids.New("usr"),
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.store.Users().CreateUser(ctx, user); err != nil {
		return nil, err
	}

	s.audit(ctx, "", "user.register", "user:"+user.ID, map[string]any{"username": user.Username, "role": string(user.Role)})
	return user, nil
}

// dummyHash sert à comparer un mot de passe quand l'identifiant est inconnu :
// la réponse prend alors le même temps qu'un mauvais mot de passe, et la durée
// ne révèle pas quels comptes existent.
var dummyHash = sync.OnceValue(func() string {
	h, _ := HashPassword("forgenet-compte-inexistant")
	return h
})

func (s *authService) Login(ctx context.Context, usernameOrEmail, password string) (*Result, error) {
	usernameOrEmail = strings.TrimSpace(usernameOrEmail)
	if usernameOrEmail == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.findUser(ctx, usernameOrEmail)
	if errors.Is(err, store.ErrNotFound) {
		_ = CheckPassword(dummyHash(), password)
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if err := CheckPassword(user.PasswordHash, password); err != nil {
		s.audit(ctx, user.ID, "auth.login_failed", "user:"+user.ID, map[string]any{"reason": "mot de passe incorrect"})
		return nil, ErrInvalidCredentials
	}

	// L'identifiant de session est connu avant la signature : il est porté
	// par le JWT et sert à vérifier sa révocation.
	now := time.Now().UTC()
	session := s.newSession(user, ids.New("sess"), "Session Web/CLI", user.Role, now, s.cfg.TokenTTL)
	signedToken, err := SignJWT(Claims{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      user.Role,
		SessionID: session.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(session.ExpiresAt),
			ID:        session.ID,
		},
	}, s.cfg.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("génération du JWT : %w", err)
	}

	session.TokenHash = HashToken(signedToken)
	if err := s.saveSession(ctx, session, "auth.login_success"); err != nil {
		return nil, err
	}
	return &Result{Token: signedToken, User: user, ExpiresAt: session.ExpiresAt}, nil
}

// CreateAPIToken crée un jeton API. Son rôle ne peut pas dépasser celui de
// l'utilisateur : un membre qui demande un jeton admin est refusé.
func (s *authService) CreateAPIToken(ctx context.Context, userID, name string, role store.Role, ttl time.Duration) (string, *store.Session, error) {
	user, err := s.store.Users().GetUserByID(ctx, userID)
	if err != nil {
		return "", nil, fmt.Errorf("utilisateur introuvable : %w", err)
	}

	if !role.IsValid() {
		role = user.Role
	}
	if !user.Role.Grants(role) {
		return "", nil, ErrForbidden
	}
	if ttl <= 0 {
		ttl = s.cfg.APITokenTTL
	}

	plainToken, tokenHash, err := GenerateAPIToken()
	if err != nil {
		return "", nil, err
	}

	session := s.newSession(user, ids.New("tok"), name, role, time.Now().UTC(), ttl)
	session.TokenHash = tokenHash
	if err := s.saveSession(ctx, session, "token.create"); err != nil {
		return "", nil, err
	}
	return plainToken, session, nil
}

// ValidateToken accepte un jeton API (fgn_…) ou un JWT, puis vérifie dans tous
// les cas la session en base : elle doit exister, ne pas être révoquée ni
// expirée. Toute erreur de lecture refuse le jeton (échec fermé).
func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, ErrUnauthenticated
	}

	tokenHash := HashToken(tokenString)
	var (
		session *store.Session
		err     error
	)
	if IsAPIToken(tokenString) {
		session, err = s.store.Sessions().GetSessionByTokenHash(ctx, tokenHash)
	} else {
		var claims *Claims
		if claims, err = ParseAndVerifyJWT(tokenString, s.cfg.JWTSecret); err != nil {
			return nil, err
		}
		session, err = s.store.Sessions().GetSessionByID(ctx, claims.SessionID)
		// Un JWT valide doit être celui enregistré pour sa session.
		if err == nil && session.TokenHash != tokenHash {
			return nil, ErrInvalidToken
		}
	}
	if err != nil {
		return nil, lookupError(err, "session")
	}

	return s.claimsForSession(ctx, session)
}

// claimsForSession construit l'identité d'une session active. Le rôle retenu est
// le plus bas entre celui du jeton et celui, actuel, de l'utilisateur : une
// rétrogradation prend effet sans attendre l'expiration du jeton.
func (s *authService) claimsForSession(ctx context.Context, session *store.Session) (*Claims, error) {
	if session.Revoked {
		return nil, ErrInvalidToken
	}
	if session.IsExpired() {
		return nil, ErrExpiredToken
	}

	user, err := s.store.Users().GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, lookupError(err, "utilisateur")
	}

	// last_used_at est indicatif : son écriture ne doit pas ralentir l'appel.
	bgCtx := context.WithoutCancel(ctx)
	go func() { _ = s.store.Sessions().TouchSession(bgCtx, session.ID, time.Now().UTC()) }()

	return &Claims{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      store.Lowest(session.Role, user.Role),
		SessionID: session.ID,
	}, nil
}

func (s *authService) RevokeSession(ctx context.Context, sessionID string) error {
	if err := s.store.Sessions().RevokeSession(ctx, sessionID); err != nil {
		return err
	}
	s.audit(ctx, audit.UserIDFromContext(ctx), "session.revoke", "session:"+sessionID, nil)
	return nil
}

func (s *authService) findUser(ctx context.Context, usernameOrEmail string) (*store.User, error) {
	if strings.Contains(usernameOrEmail, "@") {
		return s.store.Users().GetUserByEmail(ctx, strings.ToLower(usernameOrEmail))
	}
	return s.store.Users().GetUserByUsername(ctx, usernameOrEmail)
}

func (s *authService) newSession(user *store.User, id, name string, role store.Role, now time.Time, ttl time.Duration) *store.Session {
	return &store.Session{
		ID:        id,
		UserID:    user.ID,
		Name:      name,
		Role:      role,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}
}

func (s *authService) saveSession(ctx context.Context, session *store.Session, action string) error {
	if err := s.store.Sessions().CreateSession(ctx, session); err != nil {
		return fmt.Errorf("enregistrement de la session : %w", err)
	}
	s.audit(ctx, session.UserID, action, "session:"+session.ID, map[string]any{"name": session.Name, "role": string(session.Role)})
	return nil
}

// audit journalise une action sensible. Un échec d'écriture est déjà consigné
// par l'auditeur et ne bloque pas l'opération.
func (s *authService) audit(ctx context.Context, userID, action, resource string, metadata map[string]any) {
	if s.auditor == nil {
		return
	}
	if userID != "" {
		ctx = audit.WithAuditContext(ctx, userID, audit.IPFromContext(ctx))
	}
	_ = s.auditor.Log(ctx, action, resource, metadata)
}

// lookupError traduit une absence en jeton invalide ; toute autre erreur
// (base indisponible…) est remontée telle quelle et refuse aussi le jeton.
func lookupError(err error, what string) error {
	if errors.Is(err, store.ErrNotFound) {
		return ErrInvalidToken
	}
	return fmt.Errorf("lecture %s : %w", what, err)
}
