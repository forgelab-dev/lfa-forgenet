// Package sqlite fournit l'implémentation SQLite de store.Store.
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Enregistre le driver "sqlite"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/migrations"
)

type sqliteStore struct {
	db          *sql.DB
	users       store.UserStore
	sessions    store.SessionStore
	audit       store.AuditStore
	deployments store.DeploymentStore
}

// Open initialise une connexion à la base SQLite, applique les migrations requises et retourne un store prêt à l'emploi.
func Open(ctx context.Context, cfg config.DatabaseConfig) (store.Store, error) {
	// Le dossier parent est créé pour un chemin relatif comme absolu
	// (/var/lib/forgenet/forgenet.db à la première installation).
	if cfg.Path != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(cfg.Path), 0o750); err != nil {
			return nil, fmt.Errorf("création du dossier de la base de données : %w", err)
		}
	}

	dsn := fmt.Sprintf("%s?_pragma=busy_timeout(%d)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)&_pragma=synchronous(NORMAL)",
		cfg.Path,
		cfg.BusyTimeout.Milliseconds(),
	)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("ouverture de la base SQLite %s : %w", cfg.Path, err)
	}

	maxOpen := cfg.MaxOpenConns
	if maxOpen <= 0 {
		maxOpen = 1
	}
	maxIdle := cfg.MaxIdleConns
	if maxIdle <= 0 {
		maxIdle = 1
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)

	// Vérification de la connectivité
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("test de connexion à SQLite : %w", err)
	}

	// Application systématique des migrations SQL au démarrage
	if err := migrations.Apply(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("application des migrations SQL : %w", err)
	}

	s := &sqliteStore{
		db: db,
	}
	s.users = newUsersStore(db)
	s.sessions = newSessionsStore(db)
	s.audit = newAuditStore(db)
	s.deployments = newDeploymentsStore(db)

	return s, nil
}

func (s *sqliteStore) Users() store.UserStore {
	return s.users
}

func (s *sqliteStore) Sessions() store.SessionStore {
	return s.sessions
}

func (s *sqliteStore) Audit() store.AuditStore {
	return s.audit
}

func (s *sqliteStore) Deployments() store.DeploymentStore {
	return s.deployments
}

// Ping vérifie que la base de données est accessible.
func (s *sqliteStore) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *sqliteStore) Close() error {
	return s.db.Close()
}
