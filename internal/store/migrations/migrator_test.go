package migrations_test

import (
	"context"
	"database/sql"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/forgelab-dev/lfa-forgenet/internal/store/migrations"
)

func TestApply_Migrations(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	// 1. Première application : doit créer les tables et migrer avec succès
	err = migrations.Apply(ctx, db)
	require.NoError(t, err)

	// Vérification de la présence de la table schema_migrations et de la version 1
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 1").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Vérification de l'existence des tables applicatives
	tables := []string{"users", "sessions", "audit_logs"}
	for _, table := range tables {
		var name string
		err = db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		require.NoError(t, err, "la table %s doit exister", table)
		assert.Equal(t, table, name)
	}

	// 2. Deuxième application (idempotence) : aucune erreur, pas de duplication
	err = migrations.Apply(ctx, db)
	require.NoError(t, err)

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 1").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestLoad_OrdersAndRejectsDuplicateVersions(t *testing.T) {
	ordered, err := migrations.Load(fstest.MapFS{
		"0002_b.sql": {Data: []byte("SELECT 2;")},
		"0001_a.sql": {Data: []byte("SELECT 1;")},
		"README.md":  {Data: []byte("ignoré")},
	})
	require.NoError(t, err)
	require.Len(t, ordered, 2)
	assert.Equal(t, []int{1, 2}, []int{ordered[0].Version, ordered[1].Version})

	_, err = migrations.Load(fstest.MapFS{
		"0001_a.sql": {Data: []byte("SELECT 1;")},
		"0001_b.sql": {Data: []byte("SELECT 1;")},
	})
	require.ErrorContains(t, err, "en double")
}
