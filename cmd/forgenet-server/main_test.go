package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcHealth "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
	transportHTTP "github.com/forgelab-dev/lfa-forgenet/internal/transport/http"
	"github.com/forgelab-dev/lfa-forgenet/internal/version"
)

// binPath est le binaire compilé une seule fois pour tous les tests du paquet.
var binPath string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "forgenet-server-test")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	binPath = filepath.Join(dir, "forgenet-server")

	build := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "compilation du binaire : %v\n%s", err, out)
		_ = os.RemoveAll(dir)
		os.Exit(1)
	}

	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = ln.Close() }()
	return ln.Addr().(*net.TCPAddr).Port
}

// writeConfig écrit un fichier de configuration minimal dans un dossier temporaire.
func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "forgenet.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

// startServer lance le binaire et renvoie un canal qui reçoit son résultat de sortie.
func startServer(t *testing.T, env []string, args ...string) (*exec.Cmd, <-chan error) {
	t.Helper()
	cmd := exec.Command(binPath, args...)
	cmd.Dir = t.TempDir()
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Start())

	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()
	t.Cleanup(func() { _ = cmd.Process.Kill() })
	return cmd, exited
}

func exitCode(t *testing.T, err error) int {
	t.Helper()
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr)
	return exitErr.ExitCode()
}

func TestServer_E2E_FileAndEnvConfiguration(t *testing.T) {
	httpPort, grpcPort := freePort(t), freePort(t)

	// Le port gRPC vient du fichier, le port HTTP de l'environnement, qui doit
	// l'emporter sur la valeur du fichier.
	configFile := writeConfig(t, fmt.Sprintf(`
environment: test
server:
  http_port: 1
  grpc_port: %d
  shutdown_timeout: 3s
log:
  level: debug
  format: json
`, grpcPort))
	cmd, exited := startServer(t, []string{fmt.Sprintf("FORGENET_SERVER_HTTP_PORT=%d", httpPort)}, "--config", configFile)

	httpAddr := fmt.Sprintf("http://127.0.0.1:%d", httpPort)
	client := &http.Client{Timeout: time.Second}
	require.Eventually(t, func() bool {
		resp, err := client.Get(httpAddr + "/healthz")
		if err != nil {
			return false
		}
		_ = resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, 5*time.Second, 50*time.Millisecond, "le serveur HTTP n'a pas démarré")

	t.Run("healthz", func(t *testing.T) {
		resp, err := client.Get(httpAddr + "/healthz")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		var body transportHTTP.HealthResponse
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "SERVING", body.Status)
	})

	t.Run("readyz signale la base migrée", func(t *testing.T) {
		resp, err := client.Get(httpAddr + "/readyz")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		var body transportHTTP.HealthResponse
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "ok", body.Details["database"])
	})

	t.Run("version", func(t *testing.T) {
		resp, err := client.Get(httpAddr + "/version")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		var info version.BuildInfo
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&info))
		assert.Equal(t, version.Version, info.Version)
	})

	t.Run("grpc health", func(t *testing.T) {
		conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%d", grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)
		defer func() { _ = conn.Close() }()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		resp, err := grpcHealth.NewHealthClient(conn).Check(ctx, &grpcHealth.HealthCheckRequest{})
		require.NoError(t, err)
		assert.Equal(t, grpcHealth.HealthCheckResponse_SERVING, resp.GetStatus())
	})

	require.NoError(t, cmd.Process.Signal(syscall.SIGTERM))
	select {
	case err := <-exited:
		assert.Equal(t, 0, exitCode(t, err), "SIGTERM doit donner un arrêt propre")
	case <-time.After(5 * time.Second):
		t.Fatal("le serveur ne s'est pas arrêté dans le délai de grâce")
	}
}

func TestServer_ExitsNonZero(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() { _ = occupied.Close() }()

	cases := map[string]string{
		"port déjà occupé": fmt.Sprintf(`
server:
  http_port: %d
  grpc_port: %d
  shutdown_timeout: 1s
`, occupied.Addr().(*net.TCPAddr).Port, freePort(t)),
		"YAML invalide": "server: [cassé\n",
	}

	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			_, exited := startServer(t, nil, "--config", writeConfig(t, content))
			select {
			case err := <-exited:
				assert.Equal(t, 1, exitCode(t, err))
			case <-time.After(5 * time.Second):
				t.Fatal("le serveur aurait dû s'arrêter en erreur")
			}
		})
	}
}

func TestUserCreate_FirstAdminFromCommandLine(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "forgenet.db")
	configFile := writeConfig(t, fmt.Sprintf("database:\n  path: %q\n", dbPath))

	create := func(username string) (string, int) {
		cmd := exec.Command(binPath, "user", "create", "--config", configFile,
			"--username", username, "--email", username+"@forgenet.test")
		cmd.Stdin = strings.NewReader("motdepasse-de-test\n")
		out, err := cmd.CombinedOutput()
		return string(out), exitCode(t, err)
	}

	out, code := create("admin")
	require.Equal(t, 0, code, out)
	assert.Contains(t, out, "rôle admin")
	assert.NotContains(t, out, "motdepasse-de-test", "le mot de passe ne doit jamais être affiché")

	out, code = create("admin")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "existe déjà")

	db, err := sqlite.Open(context.Background(), config.DatabaseConfig{Path: dbPath, BusyTimeout: time.Second})
	require.NoError(t, err)
	defer func() { _ = db.Close() }()
	user, err := db.Users().GetUserByUsername(context.Background(), "admin")
	require.NoError(t, err)
	assert.Equal(t, store.RoleAdmin, user.Role)
}
