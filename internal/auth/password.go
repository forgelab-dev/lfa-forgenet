package auth

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrPasswordTooShort est retourné lorsque le mot de passe fait moins de 8 caractères.
	ErrPasswordTooShort = errors.New("le mot de passe doit comporter au moins 8 caractères")
	// ErrPasswordTooLong est retourné au-delà de 72 octets : bcrypt ignorerait
	// la suite, et deux mots de passe au même début seraient alors confondus.
	ErrPasswordTooLong = errors.New("le mot de passe doit comporter au plus 72 octets (environ 72 caractères sans accent)")
	// ErrInvalidCredentials est retourné lors d'un échec de correspondance identifiant/mot de passe.
	ErrInvalidCredentials = errors.New("identifiants invalides")
)

// HashPassword hache un mot de passe en clair à l'aide de l'algorithme bcrypt avec coût standard.
// La longueur minimale se compte en caractères (« é » en vaut un), la
// maximale en octets, qui est la limite réelle de bcrypt.
func HashPassword(password string) (string, error) {
	if utf8.RuneCountInString(password) < 8 {
		return "", ErrPasswordTooShort
	}
	if len(password) > 72 {
		return "", ErrPasswordTooLong
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hachage du mot de passe : %w", err)
	}

	return string(bytes), nil
}

// CheckPassword compare un mot de passe en clair avec son empreinte hachée bcrypt.
func CheckPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return ErrInvalidCredentials
	}
	return nil
}
