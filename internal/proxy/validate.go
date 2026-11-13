package proxy

import (
	"errors"
	"fmt"
	"net"
	"net/mail"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	// ErrInvalidHost indique un nom de domaine refusé.
	ErrInvalidHost = errors.New("nom de domaine invalide")
	// ErrInvalidUpstream indique une cible (conteneur ou port) refusée.
	ErrInvalidUpstream = errors.New("cible invalide")
	// ErrInvalidEmail indique une adresse ACME refusée.
	ErrInvalidEmail = errors.New("adresse électronique invalide")
	// ErrInvalidAdminAddress indique une adresse d'administration refusée :
	// seuls un socket Unix et une adresse de bouclage sont admis.
	ErrInvalidAdminAddress = errors.New("adresse d'administration de Caddy invalide")
	// ErrInvalidPort indique un port d'écoute hors de 1 à 65535.
	ErrInvalidPort = errors.New("port d'écoute invalide")
)

const (
	maxHostLength  = 253
	maxLabelLength = 63
	maxEmailLength = 254
	// maxSocketPath laisse la place du zéro final dans sun_path (108 octets
	// sous Linux, 104 sous macOS).
	maxSocketPath = 103
	// unixPrefix précède le chemin absolu du socket : « unix//run/x.sock » désigne
	// /run/x.sock, la barre après « unix » sépare le réseau de l'adresse.
	unixPrefix = "unix/"
)

// Chaque valeur écrite dans le Caddyfile passe par une de ces listes blanches :
// aucune ne laisse passer d'espace, de guillemet, d'accolade, de dièse ni de
// retour à la ligne, qui changeraient la structure du fichier.
var (
	// hostLabel : lettres minuscules, chiffres et tirets, sans tiret au bord.
	hostLabel = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)
	// containerName reprend la règle de Docker pour un nom de conteneur.
	containerName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,127}$`)
	// emailChars borne les caractères d'une adresse ACME. L'accent grave en est
	// exclu : le Caddyfile s'en sert comme guillemet.
	emailChars = regexp.MustCompile(`^[A-Za-z0-9.!$%&'*+/=?^_|~@-]+$`)
	// socketPath : chemin absolu fait de caractères ordinaires.
	socketPath = regexp.MustCompile(`^/[A-Za-z0-9._/-]+$`)
	// socketMode : droits octaux du socket, comme « 0200 » ou « 660 ».
	socketMode = regexp.MustCompile(`^0?[0-7]{3}$`)
	// numericLabel détecte une dernière étiquette tout en chiffres (adresse IP).
	numericLabel = regexp.MustCompile(`^[0-9]+$`)
)

// ValidateHost vérifie qu'un nom de domaine peut entrer dans le Caddyfile : en
// minuscules, sans point final, au moins deux étiquettes, ni joker ni adresse
// IP. Un nom internationalisé doit déjà être converti en punycode (« xn-- »).
func ValidateHost(host string) error {
	if host == "" || len(host) > maxHostLength {
		return fmt.Errorf("%w : %d caractères (1 à %d)", ErrInvalidHost, len(host), maxHostLength)
	}
	labels := strings.Split(host, ".")
	if len(labels) < 2 {
		return fmt.Errorf("%w : %q n'a pas de domaine parent", ErrInvalidHost, host)
	}
	for _, label := range labels {
		if len(label) > maxLabelLength || !hostLabel.MatchString(label) {
			return fmt.Errorf("%w : étiquette %q dans %q", ErrInvalidHost, label, host)
		}
	}
	if numericLabel.MatchString(labels[len(labels)-1]) {
		return fmt.Errorf("%w : %q ressemble à une adresse IP", ErrInvalidHost, host)
	}
	return nil
}

// Validate vérifie la cible d'une route : un nom de conteneur Docker et un port.
func (u Upstream) Validate() error {
	if !containerName.MatchString(u.Container) {
		return fmt.Errorf("%w : conteneur %q", ErrInvalidUpstream, u.Container)
	}
	if u.Port < 1 || u.Port > 65535 {
		return fmt.Errorf("%w : port %d", ErrInvalidUpstream, u.Port)
	}
	return nil
}

// ParseUpstream lit une cible écrite « conteneur:port », telle qu'elle figure
// dans la configuration du serveur.
func ParseUpstream(s string) (Upstream, error) {
	container, rawPort, ok := strings.Cut(s, ":")
	if !ok {
		return Upstream{}, fmt.Errorf("%w : %q (attendu : conteneur:port)", ErrInvalidUpstream, s)
	}
	port, err := strconv.Atoi(rawPort)
	if err != nil {
		return Upstream{}, fmt.Errorf("%w : port %q", ErrInvalidUpstream, rawPort)
	}
	u := Upstream{Container: container, Port: port}
	if err := u.Validate(); err != nil {
		return Upstream{}, err
	}
	return u, nil
}

// validateEmail n'accepte qu'une adresse nue (sans nom affiché) faite de
// caractères ordinaires.
func validateEmail(email string) error {
	if len(email) > maxEmailLength || !emailChars.MatchString(email) {
		return fmt.Errorf("%w : %q", ErrInvalidEmail, email)
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email || addr.Name != "" {
		return fmt.Errorf("%w : %q", ErrInvalidEmail, email)
	}
	return nil
}

// adminEndpoint est l'adresse d'administration vue par le client HTTP.
type adminEndpoint struct {
	network string // "unix" ou "tcp"
	address string // chemin du socket ou hôte:port
	host    string // en-tête Host attendu par Caddy
}

// parseAdminAddress n'admet qu'un socket Unix (« unix//chemin[|mode] ») ou une
// adresse de bouclage (« 127.0.0.1:2019 », « [::1]:2019 », « localhost:2019 »).
// L'API d'administration de Caddy n'a pas d'authentification : une adresse
// joignable depuis le réseau donnerait le contrôle du trafic à n'importe qui.
func parseAdminAddress(addr string) (adminEndpoint, error) {
	if rest, ok := strings.CutPrefix(addr, unixPrefix); ok {
		path, mode, hasMode := strings.Cut(rest, "|")
		if hasMode && !socketMode.MatchString(mode) {
			return adminEndpoint{}, fmt.Errorf("%w : droits %q du socket", ErrInvalidAdminAddress, mode)
		}
		if len(path) > maxSocketPath || !socketPath.MatchString(path) || filepath.Clean(path) != path {
			return adminEndpoint{}, fmt.Errorf("%w : chemin de socket %q", ErrInvalidAdminAddress, path)
		}
		// Pour un socket, Caddy accepte l'origine « localhost ».
		return adminEndpoint{network: "unix", address: path, host: "localhost"}, nil
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return adminEndpoint{}, fmt.Errorf("%w : %q", ErrInvalidAdminAddress, addr)
	}
	if n, err := strconv.Atoi(port); err != nil || n < 1 || n > 65535 {
		return adminEndpoint{}, fmt.Errorf("%w : port %q", ErrInvalidAdminAddress, port)
	}
	if host != "localhost" {
		ip := net.ParseIP(host)
		if ip == nil || !ip.IsLoopback() {
			return adminEndpoint{}, fmt.Errorf("%w : %q n'est pas une adresse de bouclage", ErrInvalidAdminAddress, host)
		}
	}
	return adminEndpoint{network: "tcp", address: addr, host: addr}, nil
}

func validatePort(name string, port int) error {
	if port != 0 && (port < 1 || port > 65535) {
		return fmt.Errorf("%w : %s %d", ErrInvalidPort, name, port)
	}
	return nil
}
