// Package proxy génère la configuration de Caddy, le reverse proxy de ForgeNet,
// et la lui fait appliquer par son API d'administration (ADR-0007).
//
// Le Caddyfile est entièrement produit à partir des routes : chaque valeur y
// entre par une liste blanche, et l'API d'administration n'est jointe que par
// un socket Unix ou une adresse de bouclage.
package proxy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog"
)

// copyMode : la copie ne contient pas de secret, mais seul le serveur l'écrit.
const copyMode = 0o640

// Loader applique un Caddyfile ; *AdminClient en est l'implémentation.
type Loader interface {
	Load(ctx context.Context, caddyfile []byte) ([]Warning, error)
}

// Manager tient la configuration de Caddy à jour : il génère le Caddyfile des
// routes, le fait appliquer, puis en garde une copie sur disque.
type Manager struct {
	mu       sync.Mutex
	settings Settings
	loader   Loader
	copyPath string
	logger   zerolog.Logger
}

// NewManager vérifie les options et prépare le client d'administration.
// copyPath reçoit la copie de la configuration appliquée ; vide, aucune copie
// n'est écrite.
func NewManager(settings Settings, copyPath string, logger zerolog.Logger) (*Manager, error) {
	if err := settings.Validate(); err != nil {
		return nil, err
	}
	client, err := NewAdminClient(settings.AdminAddress)
	if err != nil {
		return nil, err
	}
	return NewManagerWithLoader(settings, client, copyPath, logger), nil
}

// NewManagerWithLoader remplace le client d'administration, pour les tests.
// Les options sont validées à chaque Apply.
func NewManagerWithLoader(settings Settings, loader Loader, copyPath string, logger zerolog.Logger) *Manager {
	return &Manager{settings: settings, loader: loader, copyPath: copyPath, logger: logger}
}

// Apply remplace les routes servies par Caddy. Une route invalide est refusée
// avant tout appel ; une configuration refusée par Caddy (*LoadError) laisse
// l'ancienne en service. Un seul rechargement a lieu à la fois.
func (m *Manager) Apply(ctx context.Context, routes []Route) error {
	caddyfile, err := Render(m.settings, routes)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	warnings, err := m.loader.Load(ctx, caddyfile)
	for _, w := range warnings {
		m.logger.Warn().Str("avertissement", w.String()).Msg("Caddy signale un avertissement sur sa configuration")
	}
	if err != nil {
		return fmt.Errorf("rechargement de Caddy : %w", err)
	}
	m.logger.Info().Int("routes", len(routes)).Msg("configuration de Caddy appliquée")

	if m.copyPath != "" {
		// La copie sert à l'inspection : Caddy tourne déjà avec la nouvelle
		// configuration, son échec ne remet pas le rechargement en cause.
		if err := writeFileAtomic(m.copyPath, caddyfile, copyMode); err != nil {
			m.logger.Error().Err(err).Str("chemin", m.copyPath).Msg("copie du Caddyfile non écrite")
		}
	}
	return nil
}

// writeFileAtomic remplace path d'un bloc : un lecteur voit l'ancien fichier
// ou le nouveau, jamais un fichier à moitié écrit.
func writeFileAtomic(path string, data []byte, mode os.FileMode) (err error) {
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+"-*")
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tmp.Close()
			_ = os.Remove(tmp.Name())
		}
	}()
	if err = tmp.Chmod(mode); err != nil {
		return err
	}
	if _, err = tmp.Write(data); err != nil {
		return err
	}
	if err = tmp.Sync(); err != nil {
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
