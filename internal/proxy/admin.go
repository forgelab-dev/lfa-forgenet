package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	// loadTimeout borne un rechargement : Caddy répond une fois la configuration
	// appliquée, les certificats s'obtiennent ensuite en arrière-plan.
	loadTimeout = 30 * time.Second
	// maxResponse borne la réponse lue de l'API d'administration.
	maxResponse = 1 << 20
)

// ErrUnavailable indique que l'API d'administration de Caddy ne répond pas.
var ErrUnavailable = errors.New("API d'administration de Caddy injoignable")

// Warning est un avertissement de Caddy sur le Caddyfile, qui ne l'empêche pas
// d'être appliqué.
type Warning struct {
	File      string `json:"file"`
	Line      int    `json:"line"`
	Directive string `json:"directive"`
	Message   string `json:"message"`
}

func (w Warning) String() string {
	return fmt.Sprintf("%s:%d : %s", w.File, w.Line, w.Message)
}

// LoadError est une configuration refusée par Caddy. L'ancienne reste en
// service.
type LoadError struct {
	Message  string
	Warnings []Warning
}

func (e *LoadError) Error() string {
	return "configuration refusée par Caddy : " + e.Message
}

// AdminClient parle à l'API d'administration de Caddy.
type AdminClient struct {
	http *http.Client
	base string
}

// NewAdminClient prépare un client pour l'adresse d'administration de Caddy,
// limitée à un socket Unix ou à une adresse de bouclage.
func NewAdminClient(address string) (*AdminClient, error) {
	ep, err := parseAdminAddress(address)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		// Le trafic d'administration ne passe jamais par un proxy de
		// l'environnement (HTTP_PROXY…).
		Proxy: nil,
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, ep.network, ep.address)
		},
		MaxIdleConns:    1,
		IdleConnTimeout: time.Minute,
	}
	return &AdminClient{http: &http.Client{Transport: transport}, base: "http://" + ep.host}, nil
}

// Load remplace la configuration de Caddy par ce Caddyfile. Caddy l'applique
// d'un bloc ou la refuse (*LoadError) en gardant l'ancienne. Les
// avertissements renvoyés ne l'ont pas empêchée de s'appliquer.
func (c *AdminClient) Load(ctx context.Context, caddyfile []byte) ([]Warning, error) {
	ctx, cancel := context.WithTimeout(ctx, loadTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/load", bytes.NewReader(caddyfile))
	if err != nil {
		return nil, fmt.Errorf("requête vers Caddy : %w", err)
	}
	req.Header.Set("Content-Type", "text/caddyfile")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w : %w", ErrUnavailable, err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponse))
	if err != nil {
		return nil, fmt.Errorf("réponse de Caddy : %w", err)
	}

	var payload struct {
		Error    string    `json:"error"`
		Warnings []Warning `json:"warnings"`
	}
	if len(bytes.TrimSpace(body)) > 0 {
		// Un corps illisible n'empêche pas de juger du statut.
		_ = json.Unmarshal(body, &payload)
	}

	switch {
	case resp.StatusCode == http.StatusOK:
		return payload.Warnings, nil
	case resp.StatusCode == http.StatusBadRequest && payload.Error != "":
		return nil, &LoadError{Message: payload.Error, Warnings: payload.Warnings}
	default:
		return nil, fmt.Errorf("réponse inattendue de Caddy : %s : %s", resp.Status, truncate(string(body), 512))
	}
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	for n > 0 && !utf8.RuneStart(s[n]) {
		n-- // ne pas couper un caractère en deux
	}
	return s[:n] + "…"
}
