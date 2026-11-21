package templates_test

import (
	"context"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/marketplace"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
	"github.com/forgelab-dev/lfa-forgenet/templates"
)

// TestTemplates vérifie chaque dossier de templates/ selon l'ADR-0018 :
//
//	go test ./templates/
//
// Le fichier Compose est lu par docker compose lui-même. Sans lui, le test est
// sauté sur un poste de développeur, mais échoue dans la CI : un template n'y
// passe jamais sans avoir été vérifié.
func TestTemplates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	engine, err := compose.New(config.DockerConfig{})
	require.NoError(t, err)
	if err := engine.Ping(ctx); err != nil {
		if os.Getenv("CI") != "" {
			t.Fatalf("docker compose est nécessaire pour vérifier les templates : %v", err)
		}
		t.Skipf("docker compose indisponible, vérification des templates sautée : %v", err)
	}

	// Le dossier réel, pas FS : un fichier en trop (un .env oublié) n'est pas
	// embarqué, mais doit être signalé.
	fsys := os.DirFS(".")
	dirs, err := marketplace.Dirs(fsys)
	require.NoError(t, err)
	require.Contains(t, dirs, "_exemple", "le template d'exemple sert de modèle et prouve la vérification")

	for _, dir := range dirs {
		t.Run(dir, func(t *testing.T) {
			for _, p := range marketplace.Check(ctx, fsys, dir, engine.Resolve) {
				t.Errorf("templates/%s/%v", dir, p)
			}
		})
	}
}

// TestCatalogue vérifie que le catalogue embarqué se lit, sans le dossier
// d'exemple.
func TestCatalogue(t *testing.T) {
	list, err := marketplace.Load(templates.FS)
	require.NoError(t, err)

	ids := make([]string, 0, len(list))
	for _, tpl := range list {
		ids = append(ids, tpl.ID)
	}
	assert.False(t, slices.Contains(ids, "_exemple"), "l'exemple n'entre pas au catalogue")

	dirs, err := marketplace.Dirs(os.DirFS("."))
	require.NoError(t, err)
	for _, dir := range dirs {
		if dir[0] != '_' {
			assert.Contains(t, ids, dir, "chaque dossier est embarqué")
		}
	}
}
