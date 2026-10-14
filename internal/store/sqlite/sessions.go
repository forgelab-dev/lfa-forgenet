package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

const (
	sessionColumns = `id, user_id, token_hash, name, role, expires_at, created_at, last_used_at, revoked`
	selectSession  = `SELECT ` + sessionColumns + ` FROM sessions`
)

type sessionsStore struct {
	db *sql.DB
}

func newSessionsStore(db *sql.DB) store.SessionStore {
	return &sessionsStore{db: db}
}

func (s *sessionsStore) CreateSession(ctx context.Context, session *store.Session) error {
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now().UTC()
	}

	var lastUsed any
	if session.LastUsedAt != nil {
		lastUsed = formatTime(*session.LastUsedAt)
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO sessions (`+sessionColumns+`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		session.ID, session.UserID, session.TokenHash, session.Name, string(session.Role),
		formatTime(session.ExpiresAt), formatTime(session.CreatedAt), lastUsed, session.Revoked,
	)
	if err != nil {
		return wrapWriteError(err, "création de la session")
	}
	return nil
}

func (s *sessionsStore) GetSessionByID(ctx context.Context, id string) (*store.Session, error) {
	return scanSession(s.db.QueryRowContext(ctx, selectSession+` WHERE id = ?`, id))
}

func (s *sessionsStore) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*store.Session, error) {
	return scanSession(s.db.QueryRowContext(ctx, selectSession+` WHERE token_hash = ?`, tokenHash))
}

func (s *sessionsStore) TouchSession(ctx context.Context, id string, lastUsed time.Time) error {
	return execOne(ctx, s.db, "mise à jour de last_used_at",
		`UPDATE sessions SET last_used_at = ? WHERE id = ?`, formatTime(lastUsed), id)
}

func (s *sessionsStore) RevokeSession(ctx context.Context, id string) error {
	return execOne(ctx, s.db, "révocation de la session",
		`UPDATE sessions SET revoked = 1 WHERE id = ?`, id)
}

func (s *sessionsStore) DeleteExpiredSessions(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at < ?`, formatTime(time.Now())); err != nil {
		return fmt.Errorf("suppression des sessions expirées : %w", err)
	}
	return nil
}

func scanSession(row scanner) (*store.Session, error) {
	var (
		sess                 store.Session
		role                 string
		expiresAt, createdAt string
		lastUsed             sql.NullString
	)
	err := row.Scan(&sess.ID, &sess.UserID, &sess.TokenHash, &sess.Name, &role,
		&expiresAt, &createdAt, &lastUsed, &sess.Revoked)
	if err != nil {
		return nil, wrapReadError(err, "lecture de la session")
	}

	sess.Role = store.Role(role)
	if sess.ExpiresAt, err = parseTime(expiresAt); err != nil {
		return nil, err
	}
	if sess.CreatedAt, err = parseTime(createdAt); err != nil {
		return nil, err
	}
	if lastUsed.Valid {
		t, err := parseTime(lastUsed.String)
		if err != nil {
			return nil, err
		}
		sess.LastUsedAt = &t
	}
	return &sess, nil
}
