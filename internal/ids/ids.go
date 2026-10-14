// Package ids génère les identifiants opaques des entités de ForgeNet.
package ids

import (
	"crypto/rand"
	"encoding/hex"
)

// New renvoie « prefix_ » suivi de 16 caractères hexadécimaux aléatoires
// (64 bits), par exemple usr_3f9a0c1d2e4b5a67. Le préfixe rend un identifiant
// lisible dans les journaux sans en révéler l'ordre de création.
func New(prefix string) string {
	b := make([]byte, 8)
	// crypto/rand.Read ne renvoie jamais d'erreur depuis Go 1.24.
	_, _ = rand.Read(b)
	return prefix + "_" + hex.EncodeToString(b)
}
