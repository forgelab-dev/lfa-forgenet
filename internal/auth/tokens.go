package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// APITokenPrefix est le préfixe identifiant de manière unique les tokens API de ForgeNet.
	APITokenPrefix = "fgn_"
)

var (
	// ErrInvalidToken est retourné lorsque le jeton fourni est malformé ou que sa signature est invalide.
	ErrInvalidToken = errors.New("jeton d'authentification invalide")
	// ErrExpiredToken est retourné lorsque le jeton fourni a dépassé sa date d'expiration.
	ErrExpiredToken = errors.New("jeton d'authentification expiré")
)

// GenerateAPIToken crée un nouveau token API cryptographique aléatoire de 256 bits,
// et retourne sa valeur en clair (destinée à l'utilisateur) ainsi que son empreinte SHA-256 (stockée en base).
func GenerateAPIToken() (plainToken string, tokenHash string, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", fmt.Errorf("génération aléatoire : %w", err)
	}

	plainToken = APITokenPrefix + hex.EncodeToString(bytes)
	tokenHash = HashToken(plainToken)
	return plainToken, tokenHash, nil
}

// HashToken calcule l'empreinte SHA-256 hexadécimale d'un jeton pour sa persistance sécurisée.
func HashToken(plainToken string) string {
	hash := sha256.Sum256([]byte(plainToken))
	return hex.EncodeToString(hash[:])
}

// IsAPIToken détermine si une chaîne de jeton correspond au format d'un token API.
func IsAPIToken(token string) bool {
	return strings.HasPrefix(token, APITokenPrefix)
}

// SignJWT génère et signe un jeton JWT contenant les claims fournis avec la clé secrète HMAC-SHA256.
func SignJWT(claims Claims, secret string) (string, error) {
	if secret == "" {
		return "", errors.New("le secret de signature JWT ne peut pas être vide")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("signature du token JWT : %w", err)
	}

	return signed, nil
}

// ParseAndVerifyJWT vérifie la signature HMAC et l'expiration d'un jeton JWT, et en extrait les Claims.
func ParseAndVerifyJWT(tokenStr, secret string) (*Claims, error) {
	if secret == "" {
		return nil, errors.New("secret JWT manquant")
	}

	var claims Claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue : %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return &claims, nil
}
