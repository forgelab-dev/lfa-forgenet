package domains_test

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

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/deploy"
	"github.com/forgelab-dev/lfa-forgenet/internal/domains"
	"github.com/forgelab-dev/lfa-forgenet/internal/ids"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/docker"
	"github.com/forgelab-dev/lfa-forgenet/internal/proxy"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
)

// caddyImage est épinglée par version et par empreinte (ADR-0007) ; elle sert
// aussi d'application, avec « caddy respond ».
const caddyImage = "caddy:2.11.9-alpine@sha256:003a97808e781aa9f7e75cc2f43d34d461099031faf5e20faefb9f61126836f0"

const appHost = "app.forgenet.test"

func appStack(body string) []byte {
	return fmt.Appendf(nil, "services:\n  web:\n    image: %s\n    command: [\"caddy\", \"respond\", \"--listen\", \":8080\", \"--body\", %q]\n", caddyImage, body)
}

// TestDomains_RealStack relie un domaine au service d'une vraie stack Compose
// derrière un vrai Caddy, le garde après un redéploiement qui recrée le
// conteneur, puis le retire. Sauté avec -short ou sans démon Docker.
func TestDomains_RealStack(t *testing.T) {
	if testing.Short() {
		t.Skip("test d'intégration sauté avec -short")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	if err := exec.CommandContext(ctx, "docker", "info").Run(); err != nil {
		t.Skipf("démon Docker indisponible : %v", err)
	}

	suffix := strings.TrimPrefix(ids.New("t"), "t_")
	projectID, edge, caddy := "it-"+suffix, "fn-edge-"+suffix, "fn-caddy-"+suffix
	stack := deploy.StackName(projectID)

	stacks, err := compose.New(config.DockerConfig{})
	require.NoError(t, err)
	if err := stacks.Ping(ctx); err != nil {
		t.Skipf("docker compose indisponible : %v", err)
	}
	engine, err := docker.New(config.DockerConfig{PublishHost: "127.0.0.1"})
	require.NoError(t, err)
	t.Cleanup(func() { _ = engine.Close() })

	dir, err := os.MkdirTemp("", "fnd")
	require.NoError(t, err)
	adminAddress := "unix/" + filepath.Join(dir, "admin.sock")
	t.Cleanup(func() {
		cleanCtx, done := context.WithTimeout(context.Background(), time.Minute)
		defer done()
		_ = stacks.Remove(cleanCtx, stack, true)
		_ = exec.CommandContext(cleanCtx, "docker", "rm", "-f", caddy).Run()
		_ = exec.CommandContext(cleanCtx, "docker", "network", "rm", edge).Run()
		_ = os.RemoveAll(dir)
	})

	// L'application, déployée comme le ferait le service des déploiements.
	require.NoError(t, stacks.Deploy(ctx, stack, appStack("version 1")))

	// Caddy, hors du réseau partagé : la réconciliation doit l'y brancher.
	dockerRun(ctx, t, "run", "-d", "--rm", "--name", caddy,
		"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		"-v", dir+":"+dir, "-e", "CADDY_ADMIN="+adminAddress,
		"-e", "XDG_DATA_HOME=/tmp/caddy-data", "-e", "XDG_CONFIG_HOME=/tmp/caddy-config",
		"-p", "127.0.0.1::8443", caddyImage, "caddy", "run")
	socket := filepath.Join(dir, "admin.sock")
	require.Eventually(t, func() bool { _, err := os.Stat(socket); return err == nil },
		30*time.Second, 100*time.Millisecond, "socket d'administration absent")
	httpsAddr := strings.TrimSpace(dockerRun(ctx, t, "port", caddy, "8443/tcp"))

	db, err := sqlite.Open(ctx, config.DatabaseConfig{Path: ":memory:", BusyTimeout: time.Second, MaxOpenConns: 1, MaxIdleConns: 1})
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	manager, err := proxy.NewManager(proxy.Settings{
		AdminAddress: adminAddress, LocalCerts: true, HTTPPort: 8080, HTTPSPort: 8443,
	}, "", zerolog.Nop())
	require.NoError(t, err)
	svc := domains.New(db.Domains(), manager, engine, stacks, nil,
		domains.Config{Network: edge, CaddyContainer: caddy}, zerolog.Nop())

	require.NoError(t, svc.Reconcile(ctx))
	assert.Contains(t, dockerRun(ctx, t, "inspect", "-f", "{{json .NetworkSettings.Networks}}", caddy), edge,
		"Caddy est branché sur le réseau partagé")

	d, err := svc.Add(ctx, domains.AddRequest{Host: "App.Forgenet.test", ProjectID: projectID, Service: "web", Port: 8080, CreatedBy: "user:test"})
	require.NoError(t, err)
	assert.Equal(t, appHost, d.Host)

	client := httpsClient(httpsAddr, caddyRoots(ctx, t, socket))
	assert.Equal(t, "version 1", eventuallyBody(ctx, t, client, "version 1"))

	// Un redéploiement modifié recrée le conteneur, hors du réseau partagé.
	before := serviceContainer(ctx, t, stacks, stack)
	require.NoError(t, stacks.Deploy(ctx, stack, appStack("version 2")))
	after := serviceContainer(ctx, t, stacks, stack)
	require.NotEqual(t, before, after, "le conteneur a été recréé")
	require.NoError(t, svc.Reattach(ctx, projectID))
	assert.Equal(t, "version 2", eventuallyBody(ctx, t, client, "version 2"))

	// Un redéploiement identique garde le conteneur et son branchement, sans
	// rebranchement : Compose ne touche pas au réseau ajouté par ForgeNet.
	require.NoError(t, stacks.Deploy(ctx, stack, appStack("version 2")))
	require.Equal(t, after, serviceContainer(ctx, t, stacks, stack))
	body, err := get(ctx, client)
	require.NoError(t, err)
	assert.Equal(t, "version 2", body)

	// Retiré, le domaine ne répond plus et le service quitte le réseau partagé.
	require.NoError(t, svc.Remove(ctx, d.ID))
	require.Eventually(t, func() bool { _, err := get(ctx, client); return err != nil },
		10*time.Second, 200*time.Millisecond, "le domaine retiré répond encore")
	assert.NotContains(t, dockerRun(ctx, t, "inspect", "-f", "{{json .NetworkSettings.Networks}}", after), edge)
}

func serviceContainer(ctx context.Context, t *testing.T, stacks *compose.Engine, stack string) string {
	t.Helper()
	services, err := stacks.Services(ctx, stack)
	require.NoError(t, err)
	for _, s := range services {
		if s.Service == "web" {
			return s.ID
		}
	}
	t.Fatalf("service web absent de %s", stack)
	return ""
}

func eventuallyBody(ctx context.Context, t *testing.T, client *http.Client, want string) string {
	t.Helper()
	var body string
	var err error
	require.Eventually(t, func() bool {
		body, err = get(ctx, client)
		return err == nil && body == want
	}, 30*time.Second, 200*time.Millisecond, "réponse attendue %q, reçue %q (%v)", want, body, err)
	return body
}

func dockerRun(ctx context.Context, t *testing.T, args ...string) string {
	t.Helper()
	out, err := exec.CommandContext(ctx, "docker", args...).CombinedOutput()
	require.NoError(t, err, "docker %s : %s", strings.Join(args, " "), out)
	return string(out)
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

func httpsClient(addr string, roots *x509.CertPool) *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, network, addr)
			},
			TLSClientConfig:   &tls.Config{RootCAs: roots, ServerName: appHost, MinVersion: tls.VersionTLS12},
			DisableKeepAlives: true,
		},
	}
}

func get(ctx context.Context, client *http.Client) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+appHost+"/", nil)
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
