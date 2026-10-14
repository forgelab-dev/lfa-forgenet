package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
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

	const query = `
	INSERT INTO audit_logs (id, user_id, action, resource, ip_address, metadata, timestamp)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query,
		e.ID,
		e.UserID,
		e.Action,
		e.Resource,
		e.IPAddress,
		e.Metadata,
		e.Timestamp.Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("enregistrement audit log : %w", err)
	}

	return nil
}

func (s *auditStore) List(ctx context.Context, limit, offset int) ([]*store.AuditEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	const query = `
	SELECT id, user_id, action, resource, ip_address, metadata, timestamp
	FROM audit_logs
	ORDER BY timestamp DESC
	LIMIT ? OFFSET ?`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("requête audit logs : %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var events []*store.AuditEvent
	for rows.Next() {
		var e store.AuditEvent
		var tsStr string
		var userID, ipAddr, meta sql.NullString

		err := rows.Scan(
			&e.ID,
			&userID,
			&e.Action,
			&e.Resource,
			&ipAddr,
			&meta,
			&tsStr,
		)
		if err != nil {
			return nil, fmt.Errorf("scan audit log : %w", err)
		}

		if userID.Valid {
			e.UserID = userID.String
		}
		if ipAddr.Valid {
			e.IPAddress = ipAddr.String
		}
		if meta.Valid {
			e.Metadata = meta.String
		}
		e.Timestamp, _ = time.Parse(time.RFC3339Nano, tsStr)

		events = append(events, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("itération audit logs : %w", err)
	}

	return events, nil
}
