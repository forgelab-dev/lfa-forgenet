// Package version fournit les métadonnées de build et de version de ForgeNet.
package version

import (
	"runtime"
	"time"
)

var (
	// Version est la version sémantique du binaire (surpassée à la compilation via -ldflags).
	Version = "0.1.0-alpha.dev"
	// GitCommit est le hash du commit git source.
	GitCommit = "unknown"
	// BuildDate est la date de compilation au format RFC3339.
	BuildDate = time.Now().UTC().Format(time.RFC3339)
)

// BuildInfo regroupe l'ensemble des informations de version et d'environnement de compilation.
type BuildInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

// Get retourne les métadonnées actuelles de l'application.
func Get() BuildInfo {
	return BuildInfo{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
	}
}
