// Package health fournit les abstractions et l'état de santé du système pour les transports HTTP et gRPC.
package health

import (
	"context"
	"sync"
)

// Status représente l'état de santé d'un service ou sous-système.
type Status string

const (
	// StatusServing indique que le composant est pleinement opérationnel et prêt à traiter des requêtes.
	StatusServing Status = "SERVING"
	// StatusNotServing indique que le composant est en cours d'initialisation, dégradé ou en arrêt.
	StatusNotServing Status = "NOT_SERVING"
)

// Checker définit le contrat de vérification de santé d'un sous-système (base de données, docker, etc.).
type Checker interface {
	Check(ctx context.Context) error
}

// CheckFunc est un adaptateur permettant d'utiliser une fonction ordinaire comme Checker.
type CheckFunc func(ctx context.Context) error

// Check exécute la fonction de vérification sous-jacente.
func (f CheckFunc) Check(ctx context.Context) error {
	return f(ctx)
}

// Service définit les opérations de vérification de vivacité (liveness) et de préparation (readiness).
type Service interface {
	IsLive(ctx context.Context) bool
	IsReady(ctx context.Context) (bool, map[string]string)
	RegisterReadinessChecker(name string, checker Checker)
	SetServing(serving bool)
}

type service struct {
	mu        sync.RWMutex
	isServing bool
	checkers  map[string]Checker
}

// NewService instancie un nouveau service de santé centralisé.
func NewService() Service {
	return &service{
		isServing: true,
		checkers:  make(map[string]Checker),
	}
}

// IsLive vérifie si le processus est vivant (liveness probe).
func (s *service) IsLive(_ context.Context) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isServing
}

// IsReady vérifie si tous les sous-systèmes dépendants sont fonctionnels (readiness probe).
func (s *service) IsReady(ctx context.Context) (bool, map[string]string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isServing {
		return false, map[string]string{"system": "serveur en arrêt ou non démarré"}
	}

	details := make(map[string]string, len(s.checkers))
	allReady := true

	for name, checker := range s.checkers {
		if err := checker.Check(ctx); err != nil {
			allReady = false
			details[name] = err.Error()
		} else {
			details[name] = "ok"
		}
	}

	return allReady, details
}

// RegisterReadinessChecker enregistre un vérificateur de santé nommé pour la readiness.
func (s *service) RegisterReadinessChecker(name string, checker Checker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkers[name] = checker
}

// SetServing permet de basculer l'état global du serveur (ex: pendant le graceful shutdown).
func (s *service) SetServing(serving bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isServing = serving
}
