package proxy_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/ids"
	"github.com/forgelab-dev/lfa-forgenet/internal/proxy"
)

// caddyImage est épinglée par version et par empreinte (ADR-0007).
const caddyImage = "caddy:2.11.9-alpine@sha256:003a97808e781aa9f7e75cc2f43d34d461099031faf5e20faefb9f61126836f0"

const testHost = "app.forgenet.test"

// TestManager_RealCaddy applique une route à un vrai Caddy, vérifie la réponse
// HTTPS de l'application, le refus d'une configuration cassée et le retrait de
// la route. Il est sauté avec -short ou sans démon Docker.
func TestManager_RealCaddy(t *testing.T) {
	if testing.Short() {
		t.Skip("test d'intégration Caddy sauté avec -short")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	if err := exec.CommandContext(ctx, "docker", "info").Run(); err != nil {
		t.Skipf("démon Docker indisponible : %v", err)
	}

	suffix := strings.TrimPrefix(ids.New("t"), "t_")
	network, upstream, caddy := "fn-edge-"+suffix, "fn-app-"+suffix, "fn-caddy-"+suffix

	// Chemin court : un chemin de socket Unix est limité à 108 octets.
	dir, err := os.MkdirTemp("", "fnc")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	adminAddress := "unix/" + filepath.Join(dir, "admin.sock") + "|0200"

	docker(ctx, t, "network", "create", network)
	t.Cleanup(func() { dockerCleanup("network", "rm", network) })

	// L'application : un Caddy qui répond « bonjour » sur le seul réseau partagé.
	docker(ctx, t, "run", "-d", "--rm", "--name", upstream, "--network", network,
		caddyImage, "caddy", "respond", "--listen", ":8080", "--body", "bonjour")
	t.Cleanup(func() { dockerCleanup("rm", "-f", upstream) })

	// Le proxy : son API d'administration écoute sur un socket du dossier
	// partagé, monté au même chemin, sous l'utilisateur du test. Hors de la
	// racine, il écoute sur 8080 et 8443 ; seul 8443 est publié, sur le bouclage.
	docker(ctx, t, "run", "-d", "--rm", "--name", caddy, "--network", network,
		"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		"-v", dir+":"+dir,
		"-e", "CADDY_ADMIN="+adminAddress,
		"-e", "XDG_DATA_HOME=/tmp/caddy-data", "-e", "XDG_CONFIG_HOME=/tmp/caddy-config",
		"-p", "127.0.0.1::8443",
		caddyImage, "caddy", "run")
	t.Cleanup(func() { dockerCleanup("rm", "-f", caddy) })

	socket := filepath.Join(dir, "admin.sock")
	require.Eventually(t, func() bool {
		_, err := os.Stat(socket)
		return err == nil
	}, 30*time.Second, 100*time.Millisecond, "socket d'administration absent")
	httpsAddr := strings.TrimSpace(docker(ctx, t, "port", caddy, "8443/tcp"))

	settings := proxy.Settings{AdminAddress: adminAddress, LocalCerts: true, HTTPPort: 8080, HTTPSPort: 8443}
	copyPath := filepath.Join(dir, "Caddyfile")
	manager, err := proxy.NewManager(settings, copyPath, zerolog.Nop())
	require.NoError(t, err)

	route := proxy.Route{Host: testHost, Upstream: proxy.Upstream{Container: upstream, Port: 8080}}
	require.NoError(t, manager.Apply(ctx, []proxy.Route{route}))
	applied, err := os.ReadFile(copyPath)
	require.NoError(t, err)
	assert.Contains(t, string(applied), "reverse_proxy "+upstream+":8080")

	// Le certificat est signé par l'autorité interne de Caddy : le client lui
	// fait confiance, et à elle seule.
	roots := caddyRoots(ctx, t, socket)
	client := httpsClient(httpsAddr, roots)
	var body string
	require.Eventually(t, func() bool {
		body, err = get(ctx, client)
		return err == nil
	}, 30*time.Second, 200*time.Millisecond, "la route ne répond pas : %v", err)
	assert.Equal(t, "bonjour", body)

	// Une configuration cassée est refusée, et l'ancienne reste en service.
	admin, err := proxy.NewAdminClient(adminAddress)
	require.NoError(t, err)
	_, err = admin.Load(ctx, []byte(testHost+" {\n\treverse_proxy\n"))
	var loadErr *proxy.LoadError
	require.ErrorAs(t, err, &loadErr)
	body, err = get(ctx, client)
	require.NoError(t, err)
	assert.Equal(t, "bonjour", body)

	// Sans route, le domaine ne répond plus.
	require.NoError(t, manager.Apply(ctx, nil))
	require.Eventually(t, func() bool {
		_, err := get(ctx, client)
		return err != nil
	}, 10*time.Second, 200*time.Millisecond, "la route retirée répond encore")
}

func docker(ctx context.Context, t *testing.T, args ...string) string {
	t.Helper()
	out, err := exec.CommandContext(ctx, "docker", args...).CombinedOutput()
	require.NoError(t, err, "docker %s : %s", strings.Join(args, " "), out)
	return string(out)
}

func dockerCleanup(args ...string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = exec.CommandContext(ctx, "docker", args...).Run()
}

// caddyRoots lit la racine de l'autorité interne de Caddy par son API.
func caddyRoots(ctx context.Context, t *testing.T, socket string) *x509.CertPool {
	t.Helper()
	client := &http.Client{Transport: &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", socket)
		},
	}}
	pool := x509.NewCertPool()
	require.Eventually(t, func() bool {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/pki/ca/local", nil)
		if err != nil {
			return false
		}
		resp, err := client.Do(req)
		if err != nil {
			return false
		}
		defer func() { _ = resp.Body.Close() }()
		var ca struct {
			Root string `json:"root_certificate"`
		}
		if resp.StatusCode != http.StatusOK || json.NewDecoder(resp.Body).Decode(&ca) != nil {
			return false
		}
		return pool.AppendCertsFromPEM([]byte(ca.Root))
	}, 30*time.Second, 200*time.Millisecond, "racine de l'autorité interne de Caddy introuvable")
	return pool
}

// httpsClient joint Caddy sur son port publié en présentant le domaine testé.
func httpsClient(addr string, roots *x509.CertPool) *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, network, addr)
			},
			TLSClientConfig:   &tls.Config{RootCAs: roots, ServerName: testHost, MinVersion: tls.VersionTLS12},
			DisableKeepAlives: true,
		},
	}
}

func get(ctx context.Context, client *http.Client) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+testHost+"/", nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	return string(body), nil
}
