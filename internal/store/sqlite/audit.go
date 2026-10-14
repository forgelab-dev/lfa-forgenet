package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

const (
	auditColumns = `id, user_id, action, resource, ip_address, metadata, timestamp`

	// defaultAuditPageSize borne une lecture du journal sans limite explicite.
	defaultAuditPageSize = 50
)

type auditStore struct {
	db *sql.DB
}

func newAuditStore(db *sql.DB) store.AuditStore {
	return &auditStore{db: db}
}

func (s *auditStore) Record(ctx context.Context, e *store.AuditEvent) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO audit_logs (`+auditColumns+`) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.UserID, e.Action, e.Resource, e.IPAddress, e.Metadata, formatTime(e.Timestamp),
	)
	if err != nil {
		return wrapWriteError(err, "enregistrement de l'événement d'audit")
	}
	return nil
}

func (s *auditStore) List(ctx context.Context, limit, offset int) ([]*store.AuditEvent, error) {
	if limit <= 0 {
		limit = defaultAuditPageSize
	}
	offset = max(offset, 0)

	return queryAll(ctx, s.db, "lecture du journal d'audit", scanAuditEvent,
		`SELECT `+auditColumns+` FROM audit_logs ORDER BY timestamp DESC LIMIT ? OFFSET ?`, limit, offset)
}

func scanAuditEvent(row scanner) (*store.AuditEvent, error) {
	var (
		e                    store.AuditEvent
		userID, ip, metadata sql.NullString
	)
	if err := row.Scan(&e.ID, &userID, &e.Action, &e.Resource, &ip, &metadata, &e.Timestamp); err != nil {
		return nil, wrapReadError(err, "lecture de l'événement d'audit")
	}
	e.UserID, e.IPAddress, e.Metadata = userID.String, ip.String, metadata.String
	return &e, nil
}
