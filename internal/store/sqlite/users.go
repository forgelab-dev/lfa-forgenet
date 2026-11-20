package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

const (
	userColumns = `id, username, email, password_hash, role, created_at, updated_at, suspended_at`
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
		`INSERT INTO users (`+userColumns+`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		u.ID, u.Username, u.Email, u.PasswordHash, string(u.Role), formatTime(u.CreatedAt), formatTime(u.UpdatedAt),
		nullableTime(u.SuspendedAt),
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
	revoked, err := revokeSessionsTx(ctx, tx, userID, keepSessionID)
	if err != nil {
		return 0, err
	}
	// Un lien encore en attente (mot de passe oublié demandé par un tiers,
	// par exemple) ne doit pas survivre au mot de passe qu'il remplacerait.
	if _, err := tx.ExecContext(ctx, `DELETE FROM action_tokens WHERE user_id = ? AND used_at IS NULL`, userID); err != nil {
		return 0, fmt.Errorf("annulation des liens : %w", err)
	}
	return revoked, nil
}

func (s *usersStore) SetEmail(ctx context.Context, userID, email string) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("changement de l'adresse : %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `UPDATE users SET email = ?, updated_at = ? WHERE id = ?`,
		email, formatTime(time.Now()), userID)
	if err != nil {
		return 0, wrapWriteError(err, "changement de l'adresse")
	}
	if n, err := res.RowsAffected(); err != nil {
		return 0, fmt.Errorf("changement de l'adresse : %w", err)
	} else if n == 0 {
		return 0, store.ErrNotFound
	}
	// Une invitation ou un « mot de passe oublié » envoyé à une adresse mal
	// saisie ouvrirait le compte à qui la lit : ces liens partent avec elle.
	res, err = tx.ExecContext(ctx, `DELETE FROM action_tokens WHERE user_id = ? AND used_at IS NULL`, userID)
	if err != nil {
		return 0, fmt.Errorf("annulation des liens : %w", err)
	}
	cancelled, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("annulation des liens : %w", err)
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("changement de l'adresse : %w", err)
	}
	return cancelled, nil
}

func (s *usersStore) SuspendUser(ctx context.Context, id string, at time.Time) (int64, error) {
	var revoked int64
	err := s.guardedWrite(ctx, "suspension du compte", id, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx,
			`UPDATE users SET suspended_at = COALESCE(suspended_at, ?), updated_at = ? WHERE id = ?`,
			formatTime(at), formatTime(time.Now()), id); err != nil {
			return err
		}
		var err error
		if revoked, err = revokeSessionsTx(ctx, tx, id, ""); err != nil {
			return err
		}
		// Un lien en attente (invitation, mot de passe oublié) rouvrirait le
		// compte : il part avec la suspension.
		_, err = tx.ExecContext(ctx, `DELETE FROM action_tokens WHERE user_id = ? AND used_at IS NULL`, id)
		return err
	})
	return revoked, err
}

func (s *usersStore) ResumeUser(ctx context.Context, id string) error {
	return execOne(ctx, s.db, "réactivation du compte",
		`UPDATE users SET suspended_at = NULL, updated_at = ? WHERE id = ?`, formatTime(time.Now()), id)
}

func (s *usersStore) RevokeUserSessions(ctx context.Context, userID, keepSessionID string) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("fermeture des sessions : %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var exists int
	if err := tx.QueryRowContext(ctx, `SELECT 1 FROM users WHERE id = ?`, userID).Scan(&exists); err != nil {
		return 0, wrapReadError(err, "fermeture des sessions")
	}
	revoked, err := revokeSessionsTx(ctx, tx, userID, keepSessionID)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("fermeture des sessions : %w", err)
	}
	return revoked, nil
}

// DeleteUser supprime le compte ; ses sessions, jetons et liens partent avec
// lui par ON DELETE CASCADE. audit_logs n'a pas de clé étrangère : le journal
// garde ses entrées.
func (s *usersStore) DeleteUser(ctx context.Context, id string) error {
	return s.guardedWrite(ctx, "suppression de l'utilisateur", id, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
		return err
	})
}

// activeAdmin désigne un administrateur qui peut se connecter : ni suspendu,
// ni en attente de son mot de passe.
const activeAdmin = `role = 'admin' AND suspended_at IS NULL AND password_hash <> '` + store.UnusablePasswordHash + `'`

// guardedWrite exécute write dans une transaction qui refuse de retirer le
// dernier administrateur actif. La lecture du compte précède l'écriture dans
// la même transaction : en mode WAL, si un autre écrivain a validé entre les
// deux, SQLite refuse l'écriture (SQLITE_BUSY_SNAPSHOT) au lieu de laisser
// deux suspensions croisées retirer les deux derniers administrateurs.
func (s *usersStore) guardedWrite(ctx context.Context, action, id string, write func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s : %w", action, err)
	}
	defer func() { _ = tx.Rollback() }()

	var isActiveAdmin bool
	err = tx.QueryRowContext(ctx, `SELECT `+activeAdmin+` FROM users WHERE id = ?`, id).Scan(&isActiveAdmin)
	if err != nil {
		return wrapReadError(err, action)
	}
	if isActiveAdmin {
		var others int
		if err := tx.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM users WHERE `+activeAdmin+` AND id <> ?`, id).Scan(&others); err != nil {
			return fmt.Errorf("%s : %w", action, err)
		}
		if others == 0 {
			return store.ErrLastAdmin
		}
	}
	if err := write(tx); err != nil {
		return fmt.Errorf("%s : %w", action, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s : %w", action, err)
	}
	return nil
}

// revokeSessionsTx ferme dans tx les sessions et jetons du compte, sauf
// keepSessionID, et renvoie leur nombre.
func revokeSessionsTx(ctx context.Context, tx *sql.Tx, userID, keepSessionID string) (int64, error) {
	res, err := tx.ExecContext(ctx, `UPDATE sessions SET revoked = 1 WHERE user_id = ? AND id <> ? AND revoked = 0`,
		userID, keepSessionID)
	if err != nil {
		return 0, fmt.Errorf("révocation des sessions : %w", err)
	}
	revoked, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("révocation des sessions : %w", err)
	}
	return revoked, nil
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
		u         store.User
		role      string
		suspended sql.NullTime
	)
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &role, &u.CreatedAt, &u.UpdatedAt, &suspended); err != nil {
		return nil, wrapReadError(err, "lecture de l'utilisateur")
	}
	u.Role = store.Role(role)
	if suspended.Valid {
		u.SuspendedAt = &suspended.Time
	}
	return &u, nil
}
