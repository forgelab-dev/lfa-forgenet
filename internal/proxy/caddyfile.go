package proxy

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// ErrDuplicateHost indique deux routes pour le même nom de domaine.
var ErrDuplicateHost = errors.New("nom de domaine déjà routé")

// header ouvre chaque Caddyfile généré.
const header = "# Généré par ForgeNet à chaque changement de domaine : ne pas modifier à la main.\n"

// Settings sont les options globales de Caddy fixées par ForgeNet.
type Settings struct {
	// AdminAddress est l'adresse de l'API d'administration, telle que Caddy
	// l'écoute : « unix//chemin[|mode] » ou une adresse de bouclage.
	AdminAddress string
	// Email reçoit les avis de l'autorité ACME (expiration, incident). Facultatif.
	Email string
	// LocalCerts fait signer les certificats par l'autorité interne de Caddy au
	// lieu de Let's Encrypt : réservé aux essais, sans domaine public.
	LocalCerts bool
	// HTTPPort et HTTPSPort changent les ports d'écoute de Caddy ; 0 garde 80 et 443.
	HTTPPort  int
	HTTPSPort int
	// Dashboard, s'il est fixé, sert le dashboard de ForgeNet sous son propre
	// nom de domaine, en HTTPS seulement (ADR-0009). Ce nom est réservé :
	// aucune route d'application ne peut le prendre.
	Dashboard *Route
}

// hsts demande au navigateur de ne plus joindre le dashboard qu'en HTTPS,
// pendant un an. Sans includeSubDomains : les sous-domaines peuvent servir
// des applications que ForgeNet ne gère pas.
const hsts = "max-age=31536000"

// Validate vérifie chaque option avant qu'elle n'entre dans le Caddyfile.
func (s Settings) Validate() error {
	if _, err := parseAdminAddress(s.AdminAddress); err != nil {
		return err
	}
	if s.Dashboard != nil {
		if err := ValidateHost(s.Dashboard.Host); err != nil {
			return fmt.Errorf("dashboard : %w", err)
		}
		if err := s.Dashboard.Upstream.Validate(); err != nil {
			return fmt.Errorf("dashboard : %w", err)
		}
	}
	if s.Email != "" {
		if err := validateEmail(s.Email); err != nil {
			return err
		}
	}
	if err := validatePort("http_port", s.HTTPPort); err != nil {
		return err
	}
	return validatePort("https_port", s.HTTPSPort)
}

// Upstream est la cible d'une route : un conteneur joint par son nom sur le
// réseau partagé avec Caddy (ADR-0007), et le port où il écoute.
type Upstream struct {
	Container string
	Port      int
}

// Route envoie le trafic d'un nom de domaine vers un conteneur.
type Route struct {
	Host     string
	Upstream Upstream
}

// Render produit le Caddyfile des routes. Toute valeur est validée avant
// d'être écrite, et les sites sont triés par domaine : les mêmes routes donnent
// toujours le même fichier.
func Render(s Settings, routes []Route) ([]byte, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	sorted := slices.Clone(routes)
	slices.SortFunc(sorted, func(a, b Route) int { return strings.Compare(a.Host, b.Host) })
	for i, r := range sorted {
		if err := ValidateHost(r.Host); err != nil {
			return nil, err
		}
		if err := r.Upstream.Validate(); err != nil {
			return nil, fmt.Errorf("%s : %w", r.Host, err)
		}
		if (i > 0 && sorted[i-1].Host == r.Host) || (s.Dashboard != nil && s.Dashboard.Host == r.Host) {
			return nil, fmt.Errorf("%w : %s", ErrDuplicateHost, r.Host)
		}
	}

	var b bytes.Buffer
	b.WriteString(header)
	b.WriteString("{\n")
	// L'adresse d'administration est toujours écrite : sans elle, Caddy
	// reviendrait à son adresse par défaut, localhost:2019.
	fmt.Fprintf(&b, "\tadmin %s\n", s.AdminAddress)
	if s.Email != "" {
		fmt.Fprintf(&b, "\temail %s\n", s.Email)
	}
	if s.LocalCerts {
		b.WriteString("\tlocal_certs\n")
		// L'autorité de test ne doit jamais être installée dans le magasin du système.
		b.WriteString("\tskip_install_trust\n")
	}
	if s.HTTPPort != 0 {
		fmt.Fprintf(&b, "\thttp_port %d\n", s.HTTPPort)
	}
	if s.HTTPSPort != 0 {
		fmt.Fprintf(&b, "\thttps_port %d\n", s.HTTPSPort)
	}
	b.WriteString("}\n")

	if d := s.Dashboard; d != nil {
		fmt.Fprintf(&b, "\n%s {\n\theader Strict-Transport-Security %s\n\treverse_proxy %s:%d\n}\n",
			d.Host, hsts, d.Upstream.Container, d.Upstream.Port)
	}
	for _, r := range sorted {
		fmt.Fprintf(&b, "\n%s {\n\treverse_proxy %s:%d\n}\n", r.Host, r.Upstream.Container, r.Upstream.Port)
	}
	return b.Bytes(), nil
}
