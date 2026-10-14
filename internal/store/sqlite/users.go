package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

const (
	userColumns = `id, username, email, password_hash, role, created_at, updated_at`
	selectUser  = `SELECT ` + userColumns + ` FROM users`
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

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO users (`+userColumns+`) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.Email, u.PasswordHash, string(u.Role), formatTime(u.CreatedAt), formatTime(u.UpdatedAt),
	)
	if err != nil {
		return wrapWriteError(err, "insertion de l'utilisateur")
	}
	return nil
}

func (s *usersStore) GetUserByID(ctx context.Context, id string) (*store.User, error) {
	return s.getOne(ctx, selectUser+` WHERE id = ?`, id)
}

func (s *usersStore) GetUserByUsername(ctx context.Context, username string) (*store.User, error) {
	return s.getOne(ctx, selectUser+` WHERE username = ?`, username)
}

func (s *usersStore) GetUserByEmail(ctx context.Context, email string) (*store.User, error) {
	return s.getOne(ctx, selectUser+` WHERE email = ?`, email)
}

func (s *usersStore) getOne(ctx context.Context, query, arg string) (*store.User, error) {
	return scanUser(s.db.QueryRowContext(ctx, query, arg))
}

func (s *usersStore) ListUsers(ctx context.Context) ([]*store.User, error) {
	return queryAll(ctx, s.db, "liste des utilisateurs", scanUser,
		selectUser+` ORDER BY created_at ASC`)
}

func (s *usersStore) UpdateUser(ctx context.Context, u *store.User) error {
	u.UpdatedAt = time.Now().UTC()
	return execOne(ctx, s.db, "mise à jour de l'utilisateur",
		`UPDATE users SET username = ?, email = ?, password_hash = ?, role = ?, updated_at = ? WHERE id = ?`,
		u.Username, u.Email, u.PasswordHash, string(u.Role), formatTime(u.UpdatedAt), u.ID,
	)
}

func (s *usersStore) DeleteUser(ctx context.Context, id string) error {
	return execOne(ctx, s.db, "suppression de l'utilisateur", `DELETE FROM users WHERE id = ?`, id)
}

func (s *usersStore) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return 0, fmt.Errorf("comptage des utilisateurs : %w", err)
	}
	return count, nil
}

func scanUser(row scanner) (*store.User, error) {
	var (
		u                    store.User
		role                 string
		createdAt, updatedAt string
	)
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &role, &createdAt, &updatedAt); err != nil {
		return nil, wrapReadError(err, "lecture de l'utilisateur")
	}

	var err error
	u.Role = store.Role(role)
	if u.CreatedAt, err = parseTime(createdAt); err != nil {
		return nil, err
	}
	if u.UpdatedAt, err = parseTime(updatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
