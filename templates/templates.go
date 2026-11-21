// Package templates embarque le catalogue du Marketplace dans le binaire
// (ADR-0018) : un serveur propose les templates de sa version de ForgeNet,
// vérifiés par la CI, sans aller les chercher sur le réseau.
//
// Chaque dossier est un template ; son format est décrit dans README.md et
// vérifié par internal/marketplace.
package templates

import "embed"

// FS contient les fichiers de chaque template. Le dossier d'exemple, en « _ »,
// y figure aussi ; marketplace.Load l'écarte du catalogue.
//
//go:embed */manifest.yaml */docker-compose.yml */icon.svg */README.md
var FS embed.FS
