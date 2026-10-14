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
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcHealth "google.golang.org/grpc/health/grpc_health_v1"

	transportHTTP "github.com/forgelab-dev/lfa-forgenet/internal/transport/http"
	"github.com/forgelab-dev/lfa-forgenet/internal/version"
)

func getFreePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

func TestServer_E2E_Lifecycle(t *testing.T) {
	// Compilation du binaire de test dans un dossier temporaire
	tempDir := t.TempDir()
	binPath := filepath.Join(tempDir, "forgenet-server-test")

	buildCmd := exec.Command("go", "build", "-o", binPath, "./cmd/forgenet-server")
	buildCmd.Dir = "../.."
	out, err := buildCmd.CombinedOutput()
	require.NoError(t, err, string(out))

	httpPort := getFreePort(t)
	grpcPort := getFreePort(t)

	// Fichier de configuration de test
	configFile := filepath.Join(tempDir, "config.yaml")
	yamlContent := fmt.Sprintf(`
environment: test
server:
  host: "127.0.0.1"
  http_port: %d
  grpc_port: %d
  shutdown_timeout: 3s
log:
  level: "debug"
  format: "json"
`, httpPort, grpcPort)
	require.NoError(t, os.WriteFile(configFile, []byte(yamlContent), 0o600))

	// Lancement du binaire en sous-processus
	cmd := exec.Command(binPath, "--config", configFile)
	cmd.Dir = "../.."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	require.NoError(t, err)

	serverExited := make(chan error, 1)
	go func() {
		serverExited <- cmd.Wait()
	}()

	// Attente de démarrage du serveur HTTP
	httpAddr := fmt.Sprintf("http://127.0.0.1:%d", httpPort)
	httpClient := &http.Client{Timeout: 1 * time.Second}

	var httpReady bool
	for i := 0; i < 20; i++ {
		resp, err := httpClient.Get(httpAddr + "/healthz")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			httpReady = true
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(50 * time.Millisecond)
	}
	require.True(t, httpReady, "Le serveur HTTP n'a pas répondu dans le temps imparti")

	// 1. Validation de l'endpoint HTTP /healthz
	resp, err := httpClient.Get(httpAddr + "/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var healthResp transportHTTP.HealthResponse
	err = json.NewDecoder(resp.Body).Decode(&healthResp)
	require.NoError(t, err)
	assert.Equal(t, "SERVING", healthResp.Status)

	// 2. Validation de l'endpoint HTTP /version
	respVersion, err := httpClient.Get(httpAddr + "/version")
	require.NoError(t, err)
	defer respVersion.Body.Close()
	assert.Equal(t, http.StatusOK, respVersion.StatusCode)

	var versionInfo version.BuildInfo
	err = json.NewDecoder(respVersion.Body).Decode(&versionInfo)
	require.NoError(t, err)
	assert.Equal(t, version.Version, versionInfo.Version)

	// 3. Validation de l'endpoint gRPC Health
	grpcAddr := fmt.Sprintf("127.0.0.1:%d", grpcPort)
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	grpcClient := grpcHealth.NewHealthClient(conn)
	grpcCtx, grpcCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer grpcCancel()

	grpcResp, err := grpcClient.Check(grpcCtx, &grpcHealth.HealthCheckRequest{})
	require.NoError(t, err)
	assert.Equal(t, grpcHealth.HealthCheckResponse_SERVING, grpcResp.GetStatus())

	// 4. Test d'arrêt gracieux via SIGTERM
	err = cmd.Process.Signal(syscall.SIGTERM)
	require.NoError(t, err)

	select {
	case waitErr := <-serverExited:
		assert.NoError(t, waitErr, "Le serveur doit s'arrêter proprement avec un code de retour 0 sur SIGTERM")
	case <-time.After(5 * time.Second):
		t.Fatal("Le serveur ne s'est pas arrêté dans le délai de grâce")
	}
}

func TestServer_PortOccupied_ExitsNonZero(t *testing.T) {
	tempDir := t.TempDir()
	binPath := filepath.Join(tempDir, "forgenet-server-port-test")

	buildCmd := exec.Command("go", "build", "-o", binPath, "./cmd/forgenet-server")
	buildCmd.Dir = "../.."
	out, err := buildCmd.CombinedOutput()
	require.NoError(t, err, string(out))

	// On réserve délibérément un port pour simuler un conflit
	occupiedListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer func() {
		_ = occupiedListener.Close()
	}()
	occupiedPort := occupiedListener.Addr().(*net.TCPAddr).Port

	freeGRPCPort := getFreePort(t)

	configFile := filepath.Join(tempDir, "conflict-config.yaml")
	yamlContent := fmt.Sprintf(`
environment: test
server:
  host: "127.0.0.1"
  http_port: %d
  grpc_port: %d
  shutdown_timeout: 1s
log:
  level: "error"
  format: "json"
`, occupiedPort, freeGRPCPort)
	require.NoError(t, os.WriteFile(configFile, []byte(yamlContent), 0o600))

	cmd := exec.Command(binPath, "--config", configFile)
	cmd.Dir = "../.."

	err = cmd.Start()
	require.NoError(t, err)

	waitErr := cmd.Wait()
	require.Error(t, waitErr, "Le binaire doit échouer si le port d'écoute est déjà occupé")

	var exitErr *exec.ExitError
	require.ErrorAs(t, waitErr, &exitErr)
	assert.Equal(t, 1, exitErr.ExitCode(), "Le code de sortie doit être égal à 1 lors d'un échec de bind de port")
}
