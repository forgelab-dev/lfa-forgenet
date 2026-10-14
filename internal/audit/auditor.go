// Package audit assure la journalisation horodatée et infalsifiable de toutes les actions sensibles du système.
package audit

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

type contextKey string

const (
	// UserIDContextKey est la clé de contexte pour l'identifiant de l'utilisateur actif.
	UserIDContextKey contextKey = "forgenet_audit_user_id"
	// IPContextKey est la clé de contexte pour l'adresse IP distante du client.
	IPContextKey contextKey = "forgenet_audit_client_ip"
)

// WithAuditContext enrichit le contexte d'exécution avec les métadonnées nécessaires à la traçabilité.
func WithAuditContext(ctx context.Context, userID, ip string) context.Context {
	ctx = context.WithValue(ctx, UserIDContextKey, userID)
	return context.WithValue(ctx, IPContextKey, ip)
}

// UserIDFromContext extrait l'ID utilisateur du contexte s'il est présent.
func UserIDFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(UserIDContextKey).(string); ok {
		return val
	}
	return ""
}

// IPFromContext extrait l'adresse IP client du contexte si elle est présente.
func IPFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(IPContextKey).(string); ok {
		return val
	}
	return ""
}

// Auditor définit le contrat de traçabilité des opérations sensibles.
type Auditor interface {
	Log(ctx context.Context, action, resource string, metadata map[string]any) error
	List(ctx context.Context, limit, offset int) ([]*store.AuditEvent, error)
}

type auditor struct {
	auditStore store.AuditStore
	logger     zerolog.Logger
}

// NewAuditor instancie un nouveau service d'audit.
func NewAuditor(auditStore store.AuditStore, logger zerolog.Logger) Auditor {
	return &auditor{
		auditStore: auditStore,
		logger:     logger.With().Str("component", "audit").Logger(),
	}
}

// Log persiste un événement d'audit et le journalise dans les logs structurés.
func (a *auditor) Log(ctx context.Context, action, resource string, metadata map[string]any) error {
	userID := UserIDFromContext(ctx)
	clientIP := IPFromContext(ctx)

	var metaStr string
	if len(metadata) > 0 {
		raw, err := json.Marshal(metadata)
		if err == nil {
			metaStr = string(raw)
		}
	}

	eventID := generateAuditID()
	event := &store.AuditEvent{
		ID:        eventID,
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		IPAddress: clientIP,
		Metadata:  metaStr,
		Timestamp: time.Now().UTC(),
	}

	if err := a.auditStore.Record(ctx, event); err != nil {
		a.logger.Error().Err(err).
			Str("action", action).
			Str("resource", resource).
			Str("user_id", userID).
			Msg("échec de l'enregistrement de l'événement d'audit en base")
		return fmt.Errorf("enregistrement audit : %w", err)
	}

	a.logger.Info().
		Str("audit_id", eventID).
		Str("user_id", userID).
		Str("action", action).
		Str("resource", resource).
		Str("ip", clientIP).
		Msg("action sensible journalisée")

	return nil
}

// List retourne les entrées du journal d'audit de manière paginée et antéchronologique.
func (a *auditor) List(ctx context.Context, limit, offset int) ([]*store.AuditEvent, error) {
	return a.auditStore.List(ctx, limit, offset)
}

func generateAuditID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return "aud_" + hex.EncodeToString(b)
}
