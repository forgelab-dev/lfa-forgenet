package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

const (
	actionTokenColumns = `id, user_id, purpose, token_hash, expires_at, used_at, created_by, created_at` //nolint:gosec // noms de colonnes, pas un secret
	selectActionToken  = `SELECT ` + actionTokenColumns + ` FROM action_tokens`
)

type actionTokensStore struct {
	db *sql.DB
}

func newActionTokensStore(db *sql.DB) store.ActionTokenStore {
	return &actionTokensStore{db: db}
}

func (s *actionTokensStore) CreateActionToken(ctx context.Context, t *store.ActionToken) error {
	if !t.Purpose.IsValid() {
		return fmt.Errorf("usage de lien inconnu : %q", t.Purpose)
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("création du lien : %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx,
		`DELETE FROM action_tokens WHERE user_id = ? AND purpose = ? AND used_at IS NULL`,
		t.UserID, string(t.Purpose)); err != nil {
		return fmt.Errorf("annulation des liens précédents : %w", err)
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO action_tokens (id, user_id, purpose, token_hash, expires_at, created_by, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.UserID, string(t.Purpose), t.TokenHash, formatTime(t.ExpiresAt), t.CreatedBy, formatTime(t.CreatedAt),
	); err != nil {
		return wrapWriteError(err, "création du lien")
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("création du lien : %w", err)
	}
	return nil
}

func (s *actionTokensStore) GetActiveActionToken(ctx context.Context, tokenHash string, now time.Time) (*store.ActionToken, error) {
	return scanActionToken(s.db.QueryRowContext(ctx,
		selectActionToken+` WHERE token_hash = ? AND used_at IS NULL AND expires_at > ?`,
		tokenHash, formatTime(now)))
}

func (s *actionTokensStore) RedeemActionToken(ctx context.Context, tokenHash string, now time.Time, passwordHash string) (*store.ActionToken, int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("utilisation du lien : %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// L'écriture vient en premier et porte toutes les conditions : c'est elle
	// qui garantit l'usage unique, même pour deux appels simultanés.
	res, err := tx.ExecContext(ctx,
		`UPDATE action_tokens SET used_at = ? WHERE token_hash = ? AND used_at IS NULL AND expires_at > ?`,
		formatTime(now), tokenHash, formatTime(now))
	if err != nil {
		return nil, 0, fmt.Errorf("utilisation du lien : %w", err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return nil, 0, fmt.Errorf("utilisation du lien : %w", err)
	} else if n == 0 {
		return nil, 0, store.ErrNotFound
	}

	token, err := scanActionToken(tx.QueryRowContext(ctx, selectActionToken+` WHERE token_hash = ?`, tokenHash))
	if err != nil {
		return nil, 0, err
	}
	revoked, err := setPasswordTx(ctx, tx, token.UserID, passwordHash, "")
	if err != nil {
		return nil, 0, err
	}
	if err := tx.Commit(); err != nil {
		return nil, 0, fmt.Errorf("utilisation du lien : %w", err)
	}
	return token, revoked, nil
}

func scanActionToken(row scanner) (*store.ActionToken, error) {
	var (
		t       store.ActionToken
		purpose string
		usedAt  sql.NullTime
	)
	if err := row.Scan(&t.ID, &t.UserID, &purpose, &t.TokenHash, &t.ExpiresAt, &usedAt, &t.CreatedBy, &t.CreatedAt); err != nil {
		return nil, wrapReadError(err, "lecture du lien")
	}
	t.Purpose = store.ActionPurpose(purpose)
	if usedAt.Valid {
		t.UsedAt = &usedAt.Time
	}
	return &t, nil
}
