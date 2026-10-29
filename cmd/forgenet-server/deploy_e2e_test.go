package main_test

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/ids"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

const (
	e2eJWTSecret = "secret-de-test-de-bout-en-bout-32-caracteres"

	e2eStackV1 = `services:
  web:
    image: nginx:alpine
    healthcheck:
      test: ["CMD", "wget", "-q", "-O", "/dev/null", "http://127.0.0.1/"]
      interval: 1s
      retries: 20
`
	e2eStackV2 = e2eStackV1 + `  worker:
    image: alpine:3.21
    command: ["sleep", "3600"]
`
)

// TestServer_E2E_ComposeDeployment parcourt #7 de bout en bout sur le vrai
// binaire et un vrai démon : un administrateur déploie une stack par gRPC, la
// met à jour, puis la supprime. Sauté avec -short ou sans docker compose.
func TestServer_E2E_ComposeDeployment(t *testing.T) {
	if testing.Short() {
		t.Skip("test de bout en bout Docker sauté avec -short")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	for _, probe := range [][]string{{"compose", "version"}, {"info"}} {
		if err := exec.CommandContext(ctx, "docker", probe...).Run(); err != nil {
			t.Skipf("docker %s indisponible : %v", probe[0], err)
		}
	}

	dbPath := filepath.Join(t.TempDir(), "forgenet.db")
	grpcPort := freePort(t)
	configFile := writeConfig(t, fmt.Sprintf(`
environment: test
server:
  http_port: %d
  grpc_port: %d
  shutdown_timeout: 5s
database:
  path: %q
auth:
  jwt_secret: %q
`, freePort(t), grpcPort, dbPath, e2eJWTSecret))

	token := adminToken(ctx, t, configFile, dbPath)
	_, exited := startServer(t, nil, "--config", configFile)

	conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%d", grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()
	client := forgenetv1.NewDeploymentServiceClient(conn)
	authed := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

	project := ids.New("e2e")
	stack := "forgenet-" + project
	t.Cleanup(func() {
		_ = exec.Command("docker", "compose", "-p", stack, "down", "--volumes", "--remove-orphans").Run()
	})

	deployAndWait := func(content string) *forgenetv1.Deployment {
		t.Helper()
		var created *forgenetv1.CreateDeploymentResponse
		// Le serveur vient de démarrer : on réessaie tant que le port n'écoute pas.
		require.Eventually(t, func() bool {
			created, err = client.CreateDeployment(authed, &forgenetv1.CreateDeploymentRequest{
				ProjectId: project,
				Source:    &forgenetv1.DeploymentSource{Source: &forgenetv1.DeploymentSource_Compose{Compose: &forgenetv1.ComposeSource{ComposeContent: content}}},
			})
			return err == nil
		}, 10*time.Second, 100*time.Millisecond, "création du déploiement : %v", err)

		var d *forgenetv1.Deployment
		require.Eventually(t, func() bool {
			resp, err := client.GetDeployment(authed, &forgenetv1.GetDeploymentRequest{DeploymentId: created.GetDeployment().GetId()})
			d = resp.GetDeployment()
			return err == nil && d.GetFinishedAt() != nil
		}, 3*time.Minute, 200*time.Millisecond, "le déploiement ne se termine pas")
		require.Equal(t, forgenetv1.DeploymentState_DEPLOYMENT_STATE_SUCCESS, d.GetState(), d.GetErrorMessage())
		return d
	}

	// Déploiement : un conteneur, celui de web.
	first := deployAndWait(e2eStackV1)
	assert.Equal(t, []string{"web"}, runningServices(ctx, t, stack))

	// Mise à jour : worker s'ajoute, et la chaîne des versions est tenue.
	second := deployAndWait(e2eStackV2)
	assert.Equal(t, first.GetId(), second.GetPreviousDeploymentId())
	assert.Equal(t, []string{"web", "worker"}, runningServices(ctx, t, stack))

	// Suppression : plus aucun conteneur, l'historique reste.
	_, err = client.DeleteStack(authed, &forgenetv1.DeleteStackRequest{ProjectId: project, DeleteVolumes: true})
	require.NoError(t, err)
	assert.Empty(t, runningServices(ctx, t, stack))
	history, err := client.ListDeployments(authed, &forgenetv1.ListDeploymentsRequest{ProjectId: project})
	require.NoError(t, err)
	assert.Len(t, history.GetDeployments(), 2)

	select {
	case err := <-exited:
		t.Fatalf("le serveur s'est arrêté pendant le test : %v", err)
	default:
	}
}

// adminToken crée un administrateur par la ligne de commande, puis un jeton
// d'API directement dans la base : le contrat n'a pas encore de connexion.
func adminToken(ctx context.Context, t *testing.T, configFile, dbPath string) string {
	t.Helper()
	cmd := exec.CommandContext(ctx, binPath, "user", "create", "--config", configFile, "--username", "admin", "--email", "admin@forgenet.test")
	cmd.Stdin = strings.NewReader("motdepasse-de-test\n")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	db, err := sqlite.Open(ctx, config.DatabaseConfig{Path: dbPath, BusyTimeout: time.Second})
	require.NoError(t, err)
	defer func() { _ = db.Close() }()
	user, err := db.Users().GetUserByUsername(ctx, "admin")
	require.NoError(t, err)
	authSvc := auth.NewService(db, config.AuthConfig{JWTSecret: e2eJWTSecret}, audit.NewAuditor(db.Audit(), zerolog.Nop()))
	token, _, err := authSvc.CreateAPIToken(ctx, user.ID, "e2e", "", time.Hour)
	require.NoError(t, err)
	return token
}

// runningServices liste, triés, les services de la stack qui tournent.
func runningServices(ctx context.Context, t *testing.T, stack string) []string {
	t.Helper()
	out, err := exec.CommandContext(ctx, "docker", "ps",
		"--filter", "label=com.docker.compose.project="+stack, "--filter", "status=running",
		"--format", `{{.Label "com.docker.compose.service"}}`).Output()
	require.NoError(t, err)
	services := strings.Fields(string(out))
	if len(services) == 0 {
		return nil
	}
	slices.Sort(services)
	return services
}
