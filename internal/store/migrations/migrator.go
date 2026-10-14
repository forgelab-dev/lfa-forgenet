// Package migrations gère les migrations de schéma SQL versionnées et embarquées.
package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"
)

//go:embed *.sql
var migrationFiles embed.FS

// Migration représente une étape de migration SQL identifiée par son numéro de version.
type Migration struct {
	Version int
	Name    string
	SQL     string
}

// Load lit et trie les migrations SQL de fsys (les fichiers NNNN_nom.sql à sa racine).
func Load(fsys fs.FS) ([]Migration, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("lecture des fichiers de migration embarqués : %w", err)
	}

	var migrations []Migration
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) < 2 {
			continue
		}

		versionNum, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("nom de fichier de migration invalide %q : %w", entry.Name(), err)
		}

		content, err := fs.ReadFile(fsys, entry.Name())
		if err != nil {
			return nil, fmt.Errorf("lecture du fichier %s : %w", entry.Name(), err)
		}

		migrations = append(migrations, Migration{
			Version: versionNum,
			Name:    entry.Name(),
			SQL:     string(content),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Deux fichiers au même numéro sont refusés : sans ce contrôle, le second
	// serait considéré comme déjà appliqué et ne s'exécuterait jamais.
	for i := 1; i < len(migrations); i++ {
		if migrations[i].Version == migrations[i-1].Version {
			return nil, fmt.Errorf("numéro de migration %d en double : %s et %s",
				migrations[i].Version, migrations[i-1].Name, migrations[i].Name)
		}
	}

	return migrations, nil
}

// Apply applique, dans l'ordre et chacune dans sa transaction, les migrations
// embarquées qui ne figurent pas encore dans schema_migrations.
func Apply(ctx context.Context, db *sql.DB) error {
	migrations, err := Load(migrationFiles)
	if err != nil {
		return err
	}

	const createMigrationsTable = `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.ExecContext(ctx, createMigrationsTable); err != nil {
		return fmt.Errorf("initialisation de la table schema_migrations : %w", err)
	}

	applied, err := appliedVersions(ctx, db)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if applied[m.Version] {
			continue
		}
		if err := applyOne(ctx, db, m); err != nil {
			return err
		}
	}
	return nil
}

func appliedVersions(ctx context.Context, db *sql.DB) (map[int]bool, error) {
	rows, err := db.QueryContext(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("lecture des migrations appliquées : %w", err)
	}
	defer func() { _ = rows.Close() }()

	applied := make(map[int]bool)
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("lecture du numéro de version : %w", err)
		}
		applied[v] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("parcours de schema_migrations : %w", err)
	}
	return applied, nil
}

// applyOne exécute une migration et l'enregistre dans la même transaction :
// une migration en échec ne laisse ni schéma partiel ni version notée.
func applyOne(ctx context.Context, db *sql.DB, m Migration) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("transaction de la migration %d : %w", m.Version, err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, m.SQL); err != nil {
		return fmt.Errorf("exécution de la migration %d (%s) : %w", m.Version, m.Name, err)
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version, name) VALUES (?, ?)", m.Version, m.Name); err != nil {
		return fmt.Errorf("enregistrement de la migration %d : %w", m.Version, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("validation de la migration %d : %w", m.Version, err)
	}
	return nil
}
