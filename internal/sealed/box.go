// Package sealed chiffre au repos les quelques secrets que ForgeNet garde en
// base, comme le mot de passe du relais SMTP (ADR-0015).
//
// AES-256-GCM de la bibliothèque standard : un secret de quelques dizaines
// d'octets n'a besoin ni d'une enveloppe de clés ni d'une dépendance.
package sealed

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// KeySize est la taille de la clé : 32 octets, soit AES-256.
const KeySize = 32

// version précède chaque chiffré : une rotation de clé pourra plus tard
// changer de format sans rendre illisibles les chiffrés en place.
const version byte = 1

// ErrUnreadable est retourné quand un chiffré ne s'ouvre pas : clé changée,
// donnée abîmée, ou chiffré prévu pour un autre usage.
var ErrUnreadable = errors.New("secret illisible avec la clé de ce serveur (FORGENET_SECRET_KEY changée ou absente)")

// Box chiffre et déchiffre avec la clé du serveur.
type Box struct {
	aead cipher.AEAD
}

// ParseKey lit une clé de 64 caractères hexadécimaux, la forme qu'install.sh
// écrit dans .env.
func ParseKey(text string) ([]byte, error) {
	key, err := hex.DecodeString(strings.TrimSpace(text))
	if err != nil || len(key) != KeySize {
		return nil, fmt.Errorf("la clé des secrets doit faire %d caractères hexadécimaux (%d octets)", 2*KeySize, KeySize)
	}
	return key, nil
}

// NewBox prépare le chiffrement avec une clé de KeySize octets.
func NewBox(key []byte) (*Box, error) {
	if len(key) != KeySize {
		return nil, fmt.Errorf("clé des secrets : %d octets attendus, %d reçus", KeySize, len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("clé des secrets : %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("clé des secrets : %w", err)
	}
	return &Box{aead: aead}, nil
}

// Seal chiffre plaintext pour l'usage purpose (ex. « mail_settings.password ») :
// version, nonce aléatoire, puis chiffré authentifié. L'usage est lié au
// chiffré comme donnée associée : copié vers un autre champ, il ne s'ouvre pas.
func (b *Box) Seal(purpose string, plaintext []byte) ([]byte, error) {
	nonce := make([]byte, b.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("tirage du nonce : %w", err)
	}
	out := make([]byte, 0, 1+len(nonce)+len(plaintext)+b.aead.Overhead())
	out = append(out, version)
	out = append(out, nonce...)
	return b.aead.Seal(out, nonce, plaintext, []byte(purpose)), nil
}

// Open déchiffre ce que Seal a produit pour le même usage ; ErrUnreadable
// pour tout autre chiffré, sans dire pourquoi.
func (b *Box) Open(purpose string, sealed []byte) ([]byte, error) {
	n := b.aead.NonceSize()
	if len(sealed) < 1+n+b.aead.Overhead() || sealed[0] != version {
		return nil, ErrUnreadable
	}
	plain, err := b.aead.Open(nil, sealed[1:1+n], sealed[1+n:], []byte(purpose))
	if err != nil {
		return nil, ErrUnreadable
	}
	return plain, nil
}
