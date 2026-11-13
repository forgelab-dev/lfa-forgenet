package http

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

// dashboardCSP n'autorise que l'origine du dashboard, sans aucune exception
// « unsafe » : le build Vite n'a ni script ni style en ligne, et les polices
// sont servies avec lui. connect-src 'self' couvre l'API Connect et les flux
// WebSocket de la même origine (ADR-0006). data: reste permis pour les images
// et polices que Vite intègre directement dans le CSS.
const dashboardCSP = "default-src 'self'; script-src 'self'; style-src 'self'; " +
	"img-src 'self' data:; font-src 'self' data:; connect-src 'self'; " +
	"object-src 'none'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"

// immutableAssets : les fichiers de assets/ portent l'empreinte de leur
// contenu dans leur nom, ils ne changent donc jamais.
const immutableAssets = "public, max-age=31536000, immutable"

// fontTypes complète la table de types du paquet mime, qui ignore les polices
// quand l'image n'a pas de /etc/mime.types.
var fontTypes = map[string]string{".woff2": "font/woff2", ".woff": "font/woff"}

// DashboardHandler sert le dashboard compilé (ADR-0009) : les fichiers du
// dossier, et index.html pour toute adresse de page, que le routeur du
// dashboard interprète.
type DashboardHandler struct {
	root *os.Root
}

// NewDashboardHandler ouvre le dossier du dashboard compilé. Les fichiers sont
// lus par os.Root : ni « .. » ni un lien symbolique ne sortent du dossier.
func NewDashboardHandler(dir string) (*DashboardHandler, error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf("dossier du dashboard : %w", err)
	}
	info, err := root.Stat("index.html")
	if err != nil || !info.Mode().IsRegular() {
		_ = root.Close()
		return nil, fmt.Errorf("dossier du dashboard %s : index.html absent", dir)
	}
	return &DashboardHandler{root: root}, nil
}

// Close libère le dossier.
func (h *DashboardHandler) Close() error {
	return h.root.Close()
}

// RegisterRoutes installe le dashboard sous « / » : les routes de l'API et
// des sondes, plus précises, gardent la priorité.
func (h *DashboardHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/", h)
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setDashboardHeaders(w.Header())
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimPrefix(path.Clean("/"+r.URL.Path), "/")
	if name == "" {
		name = "index.html"
	}
	if hidden(name) {
		http.NotFound(w, r)
		return
	}
	f, info, err := h.open(name)
	if errors.Is(err, fs.ErrNotExist) && path.Ext(name) == "" {
		// Une adresse de page (/domaines…) : le dashboard s'en charge. Une
		// ressource absente (/assets/x.js) reste une 404, jamais une page.
		name = "index.html"
		f, info, err = h.open(name)
	}
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer func() { _ = f.Close() }()

	if strings.HasPrefix(name, "assets/") {
		w.Header().Set("Cache-Control", immutableAssets)
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}
	if t, ok := fontTypes[path.Ext(name)]; ok {
		w.Header().Set("Content-Type", t)
	}
	http.ServeContent(w, r, name, info.ModTime(), f)
}

// open n'ouvre qu'un fichier ordinaire : un dossier compte comme absent, et
// aucun contenu de dossier n'est jamais listé.
func (h *DashboardHandler) open(name string) (*os.File, fs.FileInfo, error) {
	f, err := h.root.Open(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil, err
		}
		// Sortie du dossier, droits refusés… : traités comme absents, sans
		// repli sur index.html.
		return nil, nil, fs.ErrPermission
	}
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() {
		_ = f.Close()
		return nil, nil, fs.ErrNotExist
	}
	return f, info, nil
}

// hidden reconnaît un fichier ou un dossier caché (.env, .git…), jamais servi.
func hidden(name string) bool {
	for part := range strings.SplitSeq(name, "/") {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

func setDashboardHeaders(h http.Header) {
	h.Set("Content-Security-Policy", dashboardCSP)
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("X-Frame-Options", "DENY")
	h.Set("Referrer-Policy", "no-referrer")
	h.Set("Cross-Origin-Opener-Policy", "same-origin")
	h.Set("Cross-Origin-Resource-Policy", "same-origin")
	h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()")
}
