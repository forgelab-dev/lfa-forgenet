// Package migrations gère les migrations de schéma SQL versionnées et embarquées.
package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
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

// LoadMigrations lit et trie toutes les migrations SQL embarquées dans le package.
func LoadMigrations() ([]Migration, error) {
	entries, err := migrationFiles.ReadDir(".")
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

		content, err := migrationFiles.ReadFile(entry.Name())
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

	return migrations, nil
}

// Apply applique toutes les migrations en attente de façon séquentielle et transactionnelle.
func Apply(ctx context.Context, db *sql.DB) error {
	migrations, err := LoadMigrations()
	if err != nil {
		return err
	}

	// 1. S'assurer de l'existence de la table de suivi des migrations
	const createMigrationsTable = `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.ExecContext(ctx, createMigrationsTable); err != nil {
		return fmt.Errorf("initialisation de la table schema_migrations : %w", err)
	}

	// 2. Déterminer les versions déjà appliquées
	rows, err := db.QueryContext(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("lecture des migrations appliquées : %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	applied := make(map[int]bool)
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return fmt.Errorf("lecture du numéro de version : %w", err)
		}
		applied[v] = true
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("itération sur schema_migrations : %w", err)
	}

	// 3. Appliquer chaque nouvelle migration dans une transaction dédiée
	for _, m := range migrations {
		if applied[m.Version] {
			continue
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("démarrage de la transaction pour la migration %d : %w", m.Version, err)
		}

		if _, err := tx.ExecContext(ctx, m.SQL); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("exécution de la migration %d (%s) : %w", m.Version, m.Name, err)
		}

		const insertMigration = "INSERT INTO schema_migrations (version, name) VALUES (?, ?)"
		if _, err := tx.ExecContext(ctx, insertMigration, m.Version, m.Name); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("enregistrement de la migration %d : %w", m.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("validation (commit) de la migration %d : %w", m.Version, err)
		}
	}

	return nil
}
