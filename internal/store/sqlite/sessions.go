package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

type sessionsStore struct {
	db *sql.DB
}

func newSessionsStore(db *sql.DB) store.SessionStore {
	return &sessionsStore{db: db}
}

func (s *sessionsStore) CreateSession(ctx context.Context, session *store.Session) error {
	now := time.Now().UTC()
	if session.CreatedAt.IsZero() {
		session.CreatedAt = now
	}

	const query = `
	INSERT INTO sessions (id, user_id, token_hash, name, role, expires_at, created_at, last_used_at, revoked)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var lastUsedStr any
	if session.LastUsedAt != nil {
		lastUsedStr = session.LastUsedAt.Format(time.RFC3339Nano)
	}

	revokedInt := 0
	if session.Revoked {
		revokedInt = 1
	}

	_, err := s.db.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.TokenHash,
		session.Name,
		string(session.Role),
		session.ExpiresAt.Format(time.RFC3339Nano),
		session.CreatedAt.Format(time.RFC3339Nano),
		lastUsedStr,
		revokedInt,
	)
	if err != nil {
		if isConstraintError(err) {
			return store.ErrConflict
		}
		return fmt.Errorf("création de la session : %w", err)
	}

	return nil
}

func (s *sessionsStore) GetSessionByID(ctx context.Context, id string) (*store.Session, error) {
	const query = `
	SELECT id, user_id, token_hash, name, role, expires_at, created_at, last_used_at, revoked
	FROM sessions WHERE id = ?`

	return s.scanSession(s.db.QueryRowContext(ctx, query, id))
}

func (s *sessionsStore) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*store.Session, error) {
	const query = `
	SELECT id, user_id, token_hash, name, role, expires_at, created_at, last_used_at, revoked
	FROM sessions WHERE token_hash = ?`

	return s.scanSession(s.db.QueryRowContext(ctx, query, tokenHash))
}

func (s *sessionsStore) TouchSession(ctx context.Context, id string, lastUsed time.Time) error {
	const query = `UPDATE sessions SET last_used_at = ? WHERE id = ?`

	res, err := s.db.ExecContext(ctx, query, lastUsed.Format(time.RFC3339Nano), id)
	if err != nil {
		return fmt.Errorf("mise à jour last_used_at : %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("lecture rows affected : %w", err)
	}
	if rows == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (s *sessionsStore) RevokeSession(ctx context.Context, id string) error {
	const query = `UPDATE sessions SET revoked = 1 WHERE id = ?`

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("révocation de session : %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("lecture rows affected : %w", err)
	}
	if rows == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (s *sessionsStore) DeleteExpiredSessions(ctx context.Context) error {
	const query = `DELETE FROM sessions WHERE expires_at < ?`

	nowStr := time.Now().UTC().Format(time.RFC3339Nano)
	_, err := s.db.ExecContext(ctx, query, nowStr)
	if err != nil {
		return fmt.Errorf("suppression des sessions expirées : %w", err)
	}

	return nil
}

func (s *sessionsStore) scanSession(row *sql.Row) (*store.Session, error) {
	var sess store.Session
	var roleStr, expiresStr, createdStr string
	var lastUsedStr sql.NullString
	var revokedInt int

	err := row.Scan(
		&sess.ID,
		&sess.UserID,
		&sess.TokenHash,
		&sess.Name,
		&roleStr,
		&expiresStr,
		&createdStr,
		&lastUsedStr,
		&revokedInt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, fmt.Errorf("scan session : %w", err)
	}

	sess.Role = store.Role(roleStr)
	sess.ExpiresAt, _ = time.Parse(time.RFC3339Nano, expiresStr)
	sess.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdStr)
	sess.Revoked = revokedInt == 1

	if lastUsedStr.Valid {
		t, err := time.Parse(time.RFC3339Nano, lastUsedStr.String)
		if err == nil {
			sess.LastUsedAt = &t
		}
	}

	return &sess, nil
}
