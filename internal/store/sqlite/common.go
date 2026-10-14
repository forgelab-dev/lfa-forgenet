package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

// timeLayout fixe la largeur des horodatages stockés en TEXT. RFC3339Nano retire
// les zéros de fin : « 05.5Z » se classerait alors avant « 05Z » et fausserait
// les ORDER BY et les comparaisons d'expiration faites par SQLite.
const timeLayout = "2006-01-02T15:04:05.000000000Z07:00"

// scanner est commun à *sql.Row et *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func formatTime(t time.Time) string {
	return t.UTC().Format(timeLayout)
}

func parseTime(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("horodatage illisible %q : %w", s, err)
	}
	return t.UTC(), nil
}

// isConstraintError reconnaît une violation d'unicité par son code SQLite, et
// non par le texte du message, qui peut changer d'une version à l'autre.
func isConstraintError(err error) bool {
	var sqliteErr *sqlite.Error
	if !errors.As(err, &sqliteErr) {
		return false
	}
	code := sqliteErr.Code()
	return code == sqlite3.SQLITE_CONSTRAINT_UNIQUE || code == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY
}

// wrapWriteError traduit une violation d'unicité en store.ErrConflict.
func wrapWriteError(err error, action string) error {
	if isConstraintError(err) {
		return store.ErrConflict
	}
	return fmt.Errorf("%s : %w", action, err)
}

// wrapReadError traduit l'absence de ligne en store.ErrNotFound.
func wrapReadError(err error, action string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return store.ErrNotFound
	}
	return fmt.Errorf("%s : %w", action, err)
}

// execOne exécute une écriture qui doit toucher exactement une ligne ;
// aucune ligne touchée signifie que la ressource n'existe pas.
func execOne(ctx context.Context, db *sql.DB, action, query string, args ...any) error {
	res, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return wrapWriteError(err, action)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s : lecture des lignes touchées : %w", action, err)
	}
	if n == 0 {
		return store.ErrNotFound
	}
	return nil
}

// queryAll exécute une requête et décode chaque ligne avec scan.
func queryAll[T any](ctx context.Context, db *sql.DB, action string, scan func(scanner) (*T, error), query string, args ...any) ([]*T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", action, err)
	}
	defer func() { _ = rows.Close() }()

	var items []*T
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s : %w", action, err)
	}
	return items, nil
}
