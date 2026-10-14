package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

type usersStore struct {
	db *sql.DB
}

func newUsersStore(db *sql.DB) store.UserStore {
	return &usersStore{db: db}
}

func (s *usersStore) CreateUser(ctx context.Context, u *store.User) error {
	now := time.Now().UTC()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = now
	}

	const query = `
	INSERT INTO users (id, username, email, password_hash, role, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query,
		u.ID,
		u.Username,
		u.Email,
		u.PasswordHash,
		string(u.Role),
		u.CreatedAt.Format(time.RFC3339Nano),
		u.UpdatedAt.Format(time.RFC3339Nano),
	)
	if err != nil {
		if isConstraintError(err) {
			return store.ErrConflict
		}
		return fmt.Errorf("insertion de l'utilisateur : %w", err)
	}

	return nil
}

func (s *usersStore) GetUserByID(ctx context.Context, id string) (*store.User, error) {
	const query = `
	SELECT id, username, email, password_hash, role, created_at, updated_at
	FROM users WHERE id = ?`

	return s.scanUser(s.db.QueryRowContext(ctx, query, id))
}

func (s *usersStore) GetUserByUsername(ctx context.Context, username string) (*store.User, error) {
	const query = `
	SELECT id, username, email, password_hash, role, created_at, updated_at
	FROM users WHERE username = ?`

	return s.scanUser(s.db.QueryRowContext(ctx, query, username))
}

func (s *usersStore) GetUserByEmail(ctx context.Context, email string) (*store.User, error) {
	const query = `
	SELECT id, username, email, password_hash, role, created_at, updated_at
	FROM users WHERE email = ?`

	return s.scanUser(s.db.QueryRowContext(ctx, query, email))
}

func (s *usersStore) ListUsers(ctx context.Context) ([]*store.User, error) {
	const query = `
	SELECT id, username, email, password_hash, role, created_at, updated_at
	FROM users ORDER BY created_at ASC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("requête list_users : %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var users []*store.User
	for rows.Next() {
		var u store.User
		var roleStr, createdAtStr, updatedAtStr string
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.PasswordHash,
			&roleStr,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("scan utilisateur : %w", err)
		}
		u.Role = store.Role(roleStr)
		u.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAtStr)
		u.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updatedAtStr)
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("itération utilisateurs : %w", err)
	}

	return users, nil
}

func (s *usersStore) UpdateUser(ctx context.Context, u *store.User) error {
	u.UpdatedAt = time.Now().UTC()

	const query = `
	UPDATE users
	SET username = ?, email = ?, password_hash = ?, role = ?, updated_at = ?
	WHERE id = ?`

	res, err := s.db.ExecContext(ctx, query,
		u.Username,
		u.Email,
		u.PasswordHash,
		string(u.Role),
		u.UpdatedAt.Format(time.RFC3339Nano),
		u.ID,
	)
	if err != nil {
		if isConstraintError(err) {
			return store.ErrConflict
		}
		return fmt.Errorf("mise à jour de l'utilisateur : %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("lecture des lignes affectées : %w", err)
	}
	if rows == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (s *usersStore) DeleteUser(ctx context.Context, id string) error {
	const query = `DELETE FROM users WHERE id = ?`

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("suppression de l'utilisateur : %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("lecture des lignes supprimées : %w", err)
	}
	if rows == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (s *usersStore) CountUsers(ctx context.Context) (int64, error) {
	const query = `SELECT COUNT(*) FROM users`
	var count int64
	if err := s.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("comptage des utilisateurs : %w", err)
	}
	return count, nil
}

func (s *usersStore) scanUser(row *sql.Row) (*store.User, error) {
	var u store.User
	var roleStr, createdAtStr, updatedAtStr string

	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&roleStr,
		&createdAtStr,
		&updatedAtStr,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, fmt.Errorf("lecture de l'utilisateur : %w", err)
	}

	u.Role = store.Role(roleStr)
	u.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAtStr)
	u.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updatedAtStr)

	return &u, nil
}

func isConstraintError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint") || strings.Contains(msg, "constraint failed")
}
