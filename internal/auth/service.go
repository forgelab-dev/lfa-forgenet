package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

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

// NewService instancie le service d'authentification et de contrôle d'accès.
func NewService(st store.Store, cfg config.AuthConfig, auditor audit.Auditor) Service {
	// Si aucun secret n'est spécifié (dev/test), génération d'un secret aléatoire de 256 bits
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

	return &authService{
		store:   st,
		cfg:     cfg,
		auditor: auditor,
	}
}

func (s *authService) Register(ctx context.Context, username, email, password string, role store.Role) (*store.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(strings.ToLower(email))

	if username == "" {
		return nil, errors.New("le nom d'utilisateur est obligatoire")
	}
	if email == "" {
		return nil, errors.New("l'adresse email est obligatoire")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	userCount, err := s.store.Users().CountUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("comptage des utilisateurs : %w", err)
	}

	// Bootstrap de la première installation : le premier utilisateur est obligatoirement administrateur
	assignedRole := role
	if userCount == 0 {
		assignedRole = store.RoleAdmin
	} else if !assignedRole.IsValid() {
		assignedRole = store.RoleMember
	}

	userID := generateID("usr")
	now := time.Now().UTC()

	user := &store.User{
		ID:           userID,
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         assignedRole,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.store.Users().CreateUser(ctx, user); err != nil {
		return nil, err
	}

	if s.auditor != nil {
		_ = s.auditor.Log(ctx, "user.register", fmt.Sprintf("user:%s", user.ID), map[string]any{
			"username": user.Username,
			"role":     string(user.Role),
		})
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, usernameOrEmail, password string) (*Result, error) {
	usernameOrEmail = strings.TrimSpace(usernameOrEmail)
	if usernameOrEmail == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	var user *store.User
	var err error

	if strings.Contains(usernameOrEmail, "@") {
		user, err = s.store.Users().GetUserByEmail(ctx, strings.ToLower(usernameOrEmail))
	} else {
		user, err = s.store.Users().GetUserByUsername(ctx, usernameOrEmail)
	}

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("recherche utilisateur : %w", err)
	}

	if err := CheckPassword(user.PasswordHash, password); err != nil {
		if s.auditor != nil {
			_ = s.auditor.Log(ctx, "auth.login_failed", fmt.Sprintf("user:%s", user.ID), map[string]any{
				"reason": "mot de passe incorrect",
			})
		}
		return nil, ErrInvalidCredentials
	}

	// Création d'une session JWT en base
	sessionID := generateID("sess")
	now := time.Now().UTC()
	expiresAt := now.Add(s.cfg.TokenTTL)

	claims := Claims{
		UserID:    user.ID,
		Username:  user.Username,
		Role:      user.Role,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        sessionID,
		},
	}

	signedToken, err := SignJWT(claims, s.cfg.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("génération du JWT : %w", err)
	}

	session := &store.Session{
		ID:        sessionID,
		UserID:    user.ID,
		TokenHash: HashToken(signedToken),
		Name:      "Web/CLI Session",
		Role:      user.Role,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}

	if err := s.store.Sessions().CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("enregistrement de session : %w", err)
	}

	if s.auditor != nil {
		enrichedCtx := audit.WithAuditContext(ctx, user.ID, audit.IPFromContext(ctx))
		_ = s.auditor.Log(enrichedCtx, "auth.login_success", fmt.Sprintf("user:%s", user.ID), map[string]any{
			"session_id": sessionID,
		})
	}

	return &Result{
		Token:     signedToken,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *authService) CreateAPIToken(ctx context.Context, userID, name string, role store.Role, ttl time.Duration) (string, *store.Session, error) {
	if ttl <= 0 {
		ttl = s.cfg.APITokenTTL
	}

	user, err := s.store.Users().GetUserByID(ctx, userID)
	if err != nil {
		return "", nil, fmt.Errorf("utilisateur introuvable : %w", err)
	}

	assignedRole := role
	if !assignedRole.IsValid() {
		assignedRole = user.Role
	}

	plainToken, tokenHash, err := GenerateAPIToken()
	if err != nil {
		return "", nil, err
	}

	now := time.Now().UTC()
	expiresAt := now.Add(ttl)
	sessionID := generateID("tok")

	session := &store.Session{
		ID:        sessionID,
		UserID:    user.ID,
		TokenHash: tokenHash,
		Name:      name,
		Role:      assignedRole,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}

	if err := s.store.Sessions().CreateSession(ctx, session); err != nil {
		return "", nil, fmt.Errorf("persistance du token API : %w", err)
	}

	if s.auditor != nil {
		enrichedCtx := audit.WithAuditContext(ctx, user.ID, audit.IPFromContext(ctx))
		_ = s.auditor.Log(enrichedCtx, "token.create", fmt.Sprintf("token:%s", sessionID), map[string]any{
			"name": name,
			"role": string(assignedRole),
		})
	}

	return plainToken, session, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, ErrUnauthenticated
	}

	// 1. Si le token commence par fgn_, c'est un token API persistant
	if IsAPIToken(tokenString) {
		tokenHash := HashToken(tokenString)
		session, err := s.store.Sessions().GetSessionByTokenHash(ctx, tokenHash)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return nil, ErrInvalidToken
			}
			return nil, fmt.Errorf("recherche de session : %w", err)
		}

		if !session.IsActive() {
			return nil, ErrExpiredToken
		}

		user, err := s.store.Users().GetUserByID(ctx, session.UserID)
		if err != nil {
			return nil, fmt.Errorf("utilisateur introuvable : %w", err)
		}

		// Mise à jour de last_used_at
		bgCtx := context.WithoutCancel(ctx)
		go func(sessID string) {
			_ = s.store.Sessions().TouchSession(bgCtx, sessID, time.Now().UTC())
		}(session.ID)

		return &Claims{
			UserID:    user.ID,
			Username:  user.Username,
			Role:      session.Role,
			SessionID: session.ID,
		}, nil
	}

	// 2. Sinon, validation du JWT
	claims, err := ParseAndVerifyJWT(tokenString, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}

	// Vérification de non-révocation si session_id est présente
	if claims.SessionID != "" {
		sess, err := s.store.Sessions().GetSessionByID(ctx, claims.SessionID)
		if err == nil && sess.Revoked {
			return nil, ErrInvalidToken
		}
	}

	return claims, nil
}

func (s *authService) RevokeSession(ctx context.Context, sessionID string) error {
	if err := s.store.Sessions().RevokeSession(ctx, sessionID); err != nil {
		return err
	}

	if s.auditor != nil {
		_ = s.auditor.Log(ctx, "session.revoke", fmt.Sprintf("session:%s", sessionID), nil)
	}

	return nil
}

func generateID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%s_%s", prefix, hex.EncodeToString(b))
}
