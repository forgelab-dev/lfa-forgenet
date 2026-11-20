package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	netmail "net/mail"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/ids"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

var (
	// ErrInvalidRole est retourné lorsqu'un compte est créé sans rôle reconnu.
	ErrInvalidRole = errors.New("rôle inconnu (valeurs acceptées : admin, member)")
	// ErrWrongCurrentPassword est retourné quand le mot de passe actuel donné
	// pour un changement ne correspond pas.
	ErrWrongCurrentPassword = errors.New("le mot de passe actuel est incorrect")
	// ErrSamePassword est retourné quand le nouveau mot de passe est l'actuel.
	ErrSamePassword = errors.New("le nouveau mot de passe doit être différent de l'actuel")
	// ErrLinkInvalid est retourné pour un lien inconnu, expiré, déjà utilisé
	// ou remplacé par un plus récent : la réponse ne dit pas lequel.
	ErrLinkInvalid = errors.New("ce lien n'est plus valable : il a déjà servi, a expiré ou a été remplacé par un plus récent")
	// ErrSelfAction est retourné quand un administrateur vise son propre
	// compte pour le suspendre ou le supprimer.
	ErrSelfAction = errors.New("vous ne pouvez pas suspendre ni supprimer votre propre compte")
	// ErrAccountSuspended est retourné quand un lien est demandé pour un
	// compte suspendu.
	ErrAccountSuspended = errors.New("ce compte est suspendu")
	// ErrInvalidEmail est retourné pour une adresse qui n'est pas une adresse
	// nue (sans nom affiché ni retour à la ligne), ou trop longue.
	ErrInvalidEmail = errors.New("l'adresse e-mail n'est pas valide")
	// ErrSameEmail est retourné quand la nouvelle adresse est déjà celle du
	// compte.
	ErrSameEmail = errors.New("c'est déjà l'adresse de ce compte")
)

// maxEmailLength est la longueur maximale d'une adresse (RFC 5321).
const maxEmailLength = 254

// NormalizeEmail met l'adresse en minuscules, sans espaces autour, et la
// refuse (ErrInvalidEmail) si ce n'est pas une adresse nue : elle part telle
// quelle dans l'en-tête To des e-mails, sans nom affiché ni retour à la ligne.
func NormalizeEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || len(email) > maxEmailLength || strings.ContainsAny(email, "\r\n") {
		return "", ErrInvalidEmail
	}
	addr, err := netmail.ParseAddress(email)
	if err != nil || addr.Address != email || addr.Name != "" {
		return "", ErrInvalidEmail
	}
	return email, nil
}

// EmailChange décrit un changement d'adresse fait : l'appelant en a besoin
// pour prévenir l'ancienne adresse et refaire une invitation annulée.
type EmailChange struct {
	// User est le compte à jour.
	User *store.User
	// OldEmail est l'adresse remplacée. Elle ne va jamais au journal.
	OldEmail string
	// HadPassword dit si le compte avait déjà choisi son mot de passe : sinon,
	// l'ancienne adresse n'a jamais servi qu'à recevoir l'invitation.
	HadPassword bool
	// CancelledLinks compte les liens à usage unique annulés (partis vers
	// l'ancienne adresse).
	CancelledLinks int64
}

// UnusablePassword est l'empreinte d'un compte qui n'a pas encore choisi son
// mot de passe (store.UnusablePasswordHash) : le store en a besoin pour
// reconnaître un administrateur actif.
const UnusablePassword = store.UnusablePasswordHash

// HasPassword dit si le compte a choisi son mot de passe. Sans lui, il ne
// s'ouvre que par un lien à usage unique.
func HasPassword(u *store.User) bool {
	return u.PasswordHash != UnusablePassword
}

// LinkTTL donne la durée de validité d'un lien selon son usage (ADR-0014).
func LinkTTL(purpose store.ActionPurpose) time.Duration {
	switch purpose {
	case store.PurposePasswordReset:
		return time.Hour
	case store.PurposeInvitation:
		return 7 * 24 * time.Hour
	default:
		return 24 * time.Hour
	}
}

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
	// ChangePassword change le mot de passe du compte après vérification de
	// l'actuel, et ferme ses sessions et jetons sauf keepSessionID. Renvoie
	// le nombre de sessions fermées.
	ChangePassword(ctx context.Context, userID, keepSessionID, current, next string) (int64, error)
	// ResetPassword fixe le mot de passe d'un compte sans l'actuel, pour la
	// commande lancée sur le serveur, et ferme toutes ses sessions.
	ResetPassword(ctx context.Context, username, next string) (int64, error)
	// RegisterWithoutPassword crée un compte qui ne s'ouvre qu'une fois son
	// mot de passe choisi par un lien à usage unique.
	RegisterWithoutPassword(ctx context.Context, username, email string, role store.Role) (*store.User, error)
	// CreateAccessLink fait un lien à usage unique pour le compte et renvoie
	// le jeton en clair, qui n'est jamais stocké. createdBy désigne l'auteur
	// (« user:<id> » ou « system:cli »).
	CreateAccessLink(ctx context.Context, userID string, purpose store.ActionPurpose, createdBy string) (string, *store.ActionToken, error)
	// InspectAccessLink décrit un lien encore valable, sans l'utiliser.
	InspectAccessLink(ctx context.Context, token string) (*store.ActionToken, *store.User, error)
	// RedeemAccessLink fixe le mot de passe choisi par le lien, ferme les
	// sessions du compte, puis en ouvre une nouvelle.
	RedeemAccessLink(ctx context.Context, token, newPassword string) (*Result, *store.ActionToken, error)
	// LookupAccount trouve un compte par son nom ou son adresse, comme Login ;
	// store.ErrNotFound s'il n'existe pas.
	LookupAccount(ctx context.Context, usernameOrEmail string) (*store.User, error)
	// ChangeOwnEmail change l'adresse du compte connecté après vérification
	// de son mot de passe : une session volée ne doit pas suffire à détourner
	// le « mot de passe oublié ». Les liens en attente sont annulés.
	ChangeOwnEmail(ctx context.Context, userID, currentPassword, newEmail string) (*EmailChange, error)
	// SetAccountEmail change l'adresse d'un compte sans son mot de passe, par
	// un administrateur (actorID) ou sur la machine du serveur (actorID vide).
	// Les liens en attente sont annulés.
	SetAccountEmail(ctx context.Context, actorID, userID, newEmail string) (*EmailChange, error)

	// Les actions suivantes sont faites par actorID, un administrateur, ou
	// par la machine du serveur quand actorID est vide (commande
	// « forgenet user »). Un administrateur ne peut ni se suspendre ni se
	// supprimer (ErrSelfAction), et le serveur garde toujours un
	// administrateur actif (store.ErrLastAdmin).

	// SuspendAccount suspend le compte et ferme ses sessions, jetons et liens.
	// Renvoie le compte à jour et le nombre de sessions fermées.
	SuspendAccount(ctx context.Context, actorID, userID string) (*store.User, int64, error)
	// ResumeAccount réactive un compte suspendu.
	ResumeAccount(ctx context.Context, actorID, userID string) (*store.User, error)
	// CloseSessions ferme les sessions et jetons du compte sans le suspendre,
	// sauf keepSessionID (la session de qui agit sur son propre compte).
	CloseSessions(ctx context.Context, actorID, userID, keepSessionID string) (int64, error)
	// DeleteAccount supprime le compte ; le journal d'audit garde ses entrées.
	// Renvoie le compte tel qu'il était.
	DeleteAccount(ctx context.Context, actorID, userID string) (*store.User, error)
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
	username, email, err := checkAccount(username, email, role)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return s.createUser(ctx, username, email, hashedPassword, role)
}

func (s *authService) RegisterWithoutPassword(ctx context.Context, username, email string, role store.Role) (*store.User, error) {
	username, email, err := checkAccount(username, email, role)
	if err != nil {
		return nil, err
	}
	return s.createUser(ctx, username, email, UnusablePassword, role)
}

// checkAccount normalise et vérifie les champs d'un nouveau compte.
func checkAccount(username, email string, role store.Role) (string, string, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(strings.ToLower(email))
	switch {
	case username == "":
		return "", "", errors.New("le nom d'utilisateur est obligatoire")
	case email == "":
		return "", "", errors.New("l'adresse email est obligatoire")
	case !role.IsValid():
		return "", "", ErrInvalidRole
	}
	return username, email, nil
}

func (s *authService) createUser(ctx context.Context, username, email, hashedPassword string, role store.Role) (*store.User, error) {
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
	if user.PasswordHash == UnusablePassword {
		// Même durée qu'un mauvais mot de passe : la réponse ne trahit pas un
		// compte qui attend encore son lien.
		_ = CheckPassword(dummyHash(), password)
		return nil, ErrInvalidCredentials
	}
	if err := CheckPassword(user.PasswordHash, password); err != nil {
		s.audit(ctx, user.ID, "auth.login_failed", "user:"+user.ID, map[string]any{"reason": "mot de passe incorrect"})
		return nil, ErrInvalidCredentials
	}
	// Le mot de passe est vérifié avant : un compte suspendu reçoit, dans le
	// même temps, le refus d'un mauvais mot de passe, sans apprendre qu'il
	// existe ni que son mot de passe est juste.
	if user.Suspended() {
		s.audit(ctx, user.ID, "auth.login_failed", "user:"+user.ID, map[string]any{"reason": "compte suspendu"})
		return nil, ErrInvalidCredentials
	}
	return s.openSession(ctx, user)
}

// openSession ouvre une session web/CLI pour un compte authentifié.
func (s *authService) openSession(ctx context.Context, user *store.User) (*Result, error) {
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
	// La suspension ferme déjà les sessions ; ce garde couvre un jeton créé
	// entre la lecture et la suspension.
	if user.Suspended() {
		return nil, ErrInvalidToken
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

func (s *authService) ChangePassword(ctx context.Context, userID, keepSessionID, current, next string) (int64, error) {
	user, err := s.store.Users().GetUserByID(ctx, userID)
	if err != nil {
		return 0, lookupError(err, "utilisateur")
	}
	if CheckPassword(user.PasswordHash, current) != nil {
		s.audit(ctx, user.ID, "user.password_change_failed", "user:"+user.ID, map[string]any{"reason": "mot de passe actuel incorrect"})
		return 0, ErrWrongCurrentPassword
	}
	// L'actuel vient d'être vérifié : comparer les textes suffit, sans
	// second hachage.
	if next == current {
		return 0, ErrSamePassword
	}
	revoked, err := s.setPassword(ctx, user, next, keepSessionID)
	if err != nil {
		return 0, err
	}
	s.audit(ctx, user.ID, "user.password_change", "user:"+user.ID, map[string]any{"revoked_sessions": revoked})
	return revoked, nil
}

func (s *authService) ResetPassword(ctx context.Context, username, next string) (int64, error) {
	user, err := s.store.Users().GetUserByUsername(ctx, strings.TrimSpace(username))
	if err != nil {
		return 0, err
	}
	revoked, err := s.setPassword(ctx, user, next, "")
	if err != nil {
		return 0, err
	}
	// Aucun compte n'agit : c'est un accès à la machine du serveur, comme
	// « user create ».
	s.audit(ctx, "", "user.password_reset", "user:"+user.ID, map[string]any{"username": user.Username, "revoked_sessions": revoked})
	return revoked, nil
}

// maxLinkTokenLength borne un jeton reçu avant tout calcul : un vrai jeton
// fait 43 caractères (32 octets en base64url).
const maxLinkTokenLength = 128

func (s *authService) CreateAccessLink(ctx context.Context, userID string, purpose store.ActionPurpose, createdBy string) (string, *store.ActionToken, error) {
	if !purpose.IsValid() {
		return "", nil, fmt.Errorf("usage de lien inconnu : %q", purpose)
	}
	user, err := s.store.Users().GetUserByID(ctx, userID)
	if err != nil {
		return "", nil, err
	}
	if user.Suspended() {
		return "", nil, ErrAccountSuspended
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", nil, fmt.Errorf("génération aléatoire : %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(b)
	now := time.Now().UTC()
	link := &store.ActionToken{
		ID:        ids.New("act"),
		UserID:    userID,
		Purpose:   purpose,
		TokenHash: HashToken(token),
		ExpiresAt: now.Add(LinkTTL(purpose)),
		CreatedBy: createdBy,
		CreatedAt: now,
	}
	if err := s.store.ActionTokens().CreateActionToken(ctx, link); err != nil {
		return "", nil, err
	}
	// Le journal garde l'usage et l'auteur, jamais le jeton.
	s.audit(ctx, "", "user.access_link_create", "user:"+userID, map[string]any{"purpose": string(purpose), "created_by": createdBy})
	return token, link, nil
}

func (s *authService) InspectAccessLink(ctx context.Context, token string) (*store.ActionToken, *store.User, error) {
	if token == "" || len(token) > maxLinkTokenLength {
		return nil, nil, ErrLinkInvalid
	}
	link, err := s.store.ActionTokens().GetActiveActionToken(ctx, HashToken(token), time.Now().UTC())
	if errors.Is(err, store.ErrNotFound) {
		return nil, nil, ErrLinkInvalid
	}
	if err != nil {
		return nil, nil, err
	}
	user, err := s.store.Users().GetUserByID(ctx, link.UserID)
	if errors.Is(err, store.ErrNotFound) {
		return nil, nil, ErrLinkInvalid
	}
	if err != nil {
		return nil, nil, err
	}
	if user.Suspended() {
		return nil, nil, ErrLinkInvalid
	}
	return link, user, nil
}

func (s *authService) RedeemAccessLink(ctx context.Context, token, newPassword string) (*Result, *store.ActionToken, error) {
	// Le lien est vérifié avant le hachage, coûteux : un appel sans lien
	// valable ne fait pas travailler bcrypt.
	if _, _, err := s.InspectAccessLink(ctx, token); err != nil {
		return nil, nil, err
	}
	// Le mot de passe est vérifié avant d'utiliser le lien : un refus le
	// laisse valable pour un nouvel essai.
	hash, err := HashPassword(newPassword)
	if err != nil {
		return nil, nil, err
	}
	link, revoked, err := s.store.ActionTokens().RedeemActionToken(ctx, HashToken(token), time.Now().UTC(), hash)
	if errors.Is(err, store.ErrNotFound) {
		return nil, nil, ErrLinkInvalid
	}
	if err != nil {
		return nil, nil, err
	}
	s.audit(ctx, link.UserID, "user.access_link_used", "user:"+link.UserID, map[string]any{"purpose": string(link.Purpose), "revoked_sessions": revoked})

	user, err := s.store.Users().GetUserByID(ctx, link.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("lecture du compte : %w", err)
	}
	res, err := s.openSession(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return res, link, nil
}

func (s *authService) setPassword(ctx context.Context, user *store.User, next, keepSessionID string) (int64, error) {
	hash, err := HashPassword(next)
	if err != nil {
		return 0, err
	}
	return s.store.Users().SetPassword(ctx, user.ID, hash, keepSessionID)
}

func (s *authService) SuspendAccount(ctx context.Context, actorID, userID string) (*store.User, int64, error) {
	if actorID != "" && actorID == userID {
		return nil, 0, ErrSelfAction
	}
	revoked, err := s.store.Users().SuspendUser(ctx, userID, time.Now().UTC())
	if err != nil {
		return nil, 0, err
	}
	user, err := s.store.Users().GetUserByID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	s.audit(ctx, actorID, "user.suspend", "user:"+userID, map[string]any{"username": user.Username, "revoked_sessions": revoked})
	return user, revoked, nil
}

func (s *authService) ResumeAccount(ctx context.Context, actorID, userID string) (*store.User, error) {
	if err := s.store.Users().ResumeUser(ctx, userID); err != nil {
		return nil, err
	}
	user, err := s.store.Users().GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	s.audit(ctx, actorID, "user.resume", "user:"+userID, map[string]any{"username": user.Username})
	return user, nil
}

func (s *authService) CloseSessions(ctx context.Context, actorID, userID, keepSessionID string) (int64, error) {
	revoked, err := s.store.Users().RevokeUserSessions(ctx, userID, keepSessionID)
	if err != nil {
		return 0, err
	}
	s.audit(ctx, actorID, "user.sessions_close", "user:"+userID, map[string]any{"revoked_sessions": revoked})
	return revoked, nil
}

func (s *authService) DeleteAccount(ctx context.Context, actorID, userID string) (*store.User, error) {
	if actorID != "" && actorID == userID {
		return nil, ErrSelfAction
	}
	user, err := s.store.Users().GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.store.Users().DeleteUser(ctx, userID); err != nil {
		return nil, err
	}
	// Le compte n'existe plus : son nom reste au journal pour lire les
	// entrées qu'il laisse.
	s.audit(ctx, actorID, "user.delete", "user:"+userID, map[string]any{"username": user.Username, "role": string(user.Role)})
	return user, nil
}

func (s *authService) ChangeOwnEmail(ctx context.Context, userID, currentPassword, newEmail string) (*EmailChange, error) {
	user, err := s.store.Users().GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	// Un compte sans mot de passe ne se connecte pas : il n'arrive pas ici.
	if CheckPassword(user.PasswordHash, currentPassword) != nil {
		s.audit(ctx, user.ID, "user.email_change_failed", "user:"+user.ID, map[string]any{"reason": "mot de passe actuel incorrect"})
		return nil, ErrWrongCurrentPassword
	}
	return s.setEmail(ctx, user.ID, user, newEmail, "self")
}

func (s *authService) SetAccountEmail(ctx context.Context, actorID, userID, newEmail string) (*EmailChange, error) {
	user, err := s.store.Users().GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	by := "admin"
	switch actorID {
	case "":
		by = "server"
	case userID:
		by = "self"
	}
	return s.setEmail(ctx, actorID, user, newEmail, by)
}

// setEmail change l'adresse de user. by dit qui agit (self, admin ou server),
// pour le journal ; les adresses n'y figurent jamais : il est lu par tous les
// administrateurs et conservé.
func (s *authService) setEmail(ctx context.Context, actorID string, user *store.User, newEmail, by string) (*EmailChange, error) {
	email, err := NormalizeEmail(newEmail)
	if err != nil {
		return nil, err
	}
	if email == user.Email {
		return nil, ErrSameEmail
	}
	cancelled, err := s.store.Users().SetEmail(ctx, user.ID, email)
	if err != nil {
		return nil, err
	}
	change := &EmailChange{OldEmail: user.Email, HadPassword: HasPassword(user), CancelledLinks: cancelled}
	updated := *user
	updated.Email = email
	change.User = &updated
	s.audit(ctx, actorID, "user.email_change", "user:"+user.ID, map[string]any{
		"username": user.Username, "by": by, "cancelled_links": cancelled,
	})
	return change, nil
}

func (s *authService) LookupAccount(ctx context.Context, usernameOrEmail string) (*store.User, error) {
	usernameOrEmail = strings.TrimSpace(usernameOrEmail)
	if usernameOrEmail == "" {
		return nil, store.ErrNotFound
	}
	return s.findUser(ctx, usernameOrEmail)
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
