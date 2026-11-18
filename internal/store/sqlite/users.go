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

func (s *usersStore) SetPassword(ctx context.Context, userID, passwordHash, keepSessionID string) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("changement du mot de passe : %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	revoked, err := setPasswordTx(ctx, tx, userID, passwordHash, keepSessionID)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("changement du mot de passe : %w", err)
	}
	return revoked, nil
}

// setPasswordTx fixe le mot de passe dans tx, ferme les sessions du compte
// sauf keepSessionID et annule ses liens à usage unique encore inutilisés.
func setPasswordTx(ctx context.Context, tx *sql.Tx, userID, passwordHash, keepSessionID string) (int64, error) {
	res, err := tx.ExecContext(ctx, `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`,
		passwordHash, formatTime(time.Now()), userID)
	if err != nil {
		return 0, fmt.Errorf("changement du mot de passe : %w", err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return 0, fmt.Errorf("changement du mot de passe : %w", err)
	} else if n == 0 {
		return 0, store.ErrNotFound
	}

	// L'ancien mot de passe a pu fuiter : toute session ouverte avec lui, et
	// tout jeton créé par qui le détenait, est fermé avec lui.
	res, err = tx.ExecContext(ctx, `UPDATE sessions SET revoked = 1 WHERE user_id = ? AND id <> ? AND revoked = 0`,
		userID, keepSessionID)
	if err != nil {
		return 0, fmt.Errorf("révocation des sessions : %w", err)
	}
	revoked, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("révocation des sessions : %w", err)
	}
	// Un lien encore en attente (mot de passe oublié demandé par un tiers,
	// par exemple) ne doit pas survivre au mot de passe qu'il remplacerait.
	if _, err := tx.ExecContext(ctx, `DELETE FROM action_tokens WHERE user_id = ? AND used_at IS NULL`, userID); err != nil {
		return 0, fmt.Errorf("annulation des liens : %w", err)
	}
	return revoked, nil
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
		u    store.User
		role string
	)
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, wrapReadError(err, "lecture de l'utilisateur")
	}
	u.Role = store.Role(role)
	return &u, nil
}
