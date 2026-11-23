package main_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/ids"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
	"github.com/forgelab-dev/lfa-forgenet/templates"
)

const (
	// e2eManifest décrit un template dont le secret est tiré au hasard.
	e2eManifest = `schema: 1
name: Secret de bout en bout
description: Un conteneur qui n'est sain qu'avec son secret.
category: other
version: "1"
service: app
license: MIT
homepage: https://example.org
variables:
  - name: APP_SECRET
    label: Secret de l'application
    secret: true
    generate: 24
  - name: APP_NAME
    label: Nom
    default: forgenet
`
	// La sonde de santé échoue si Compose n'a pas reçu APP_SECRET ; ${…:?}
	// refuse même le fichier sans elle.
	e2eCompose = `services:
  app:
    image: alpine:3.21
    command: ["sleep", "3600"]
    environment:
      APP_SECRET: ${APP_SECRET:?}
      APP_NAME: ${APP_NAME}
    healthcheck:
      test: ["CMD-SHELL", "test -n \"$$APP_SECRET\""]
      interval: 1s
      retries: 20
`
)

// writeE2ECatalog écrit un catalogue d'un template, « secret », dans un dossier.
func writeE2ECatalog(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	icon, err := fs.ReadFile(templates.FS, "_exemple/icon.svg")
	require.NoError(t, err)
	files := map[string][]byte{
		"manifest.yaml":      []byte(e2eManifest),
		"docker-compose.yml": []byte(e2eCompose),
		"icon.svg":           icon,
		"README.md":          []byte("# Secret de bout en bout\n"),
	}
	require.NoError(t, os.Mkdir(filepath.Join(dir, "secret"), 0o700))
	for name, data := range files {
		require.NoError(t, os.WriteFile(filepath.Join(dir, "secret", name), data, 0o600))
	}
	return dir
}

// TestServer_E2E_MarketplaceInstall parcourt l'installation d'un template
// (#116, ADR-0019) sur le vrai binaire et un vrai démon : le secret tiré au
// hasard arrive dans le conteneur, n'est renvoyé qu'une fois, et ne se lit ni
// dans le journal ni en base. Sauté avec -short ou sans docker compose.
func TestServer_E2E_MarketplaceInstall(t *testing.T) {
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

	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)
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
secrets:
  key: %q
marketplace:
  dir: %q
`, freePort(t), grpcPort, dbPath, e2eJWTSecret, hex.EncodeToString(key), writeE2ECatalog(t)))

	token := adminToken(ctx, t, configFile, dbPath)
	_, exited := startServer(t, nil, "--config", configFile)

	conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%d", grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()
	marketplace := forgenetv1.NewMarketplaceServiceClient(conn)
	deployments := forgenetv1.NewDeploymentServiceClient(conn)
	authed := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

	project := ids.New("e2e")
	stack := "forgenet-" + project
	t.Cleanup(func() {
		_ = exec.Command("docker", "compose", "-p", stack, "down", "--volumes", "--remove-orphans").Run()
	})

	// Installation sans valeur : le secret est tiré, le nom prend son défaut.
	var installed *forgenetv1.InstallTemplateResponse
	require.Eventually(t, func() bool {
		installed, err = marketplace.InstallTemplate(authed, &forgenetv1.InstallTemplateRequest{TemplateId: "secret", ProjectId: project})
		return status.Code(err) != codes.Unavailable
	}, 10*time.Second, 100*time.Millisecond, "le serveur n'écoute pas")
	require.NoError(t, err)
	require.Len(t, installed.GetGeneratedSecrets(), 1)
	secret := installed.GetGeneratedSecrets()[0].GetValue()
	assert.Len(t, secret, 24)
	id := installed.GetDeployment().GetId()

	var d *forgenetv1.Deployment
	require.Eventually(t, func() bool {
		resp, err := deployments.GetDeployment(authed, &forgenetv1.GetDeploymentRequest{DeploymentId: id})
		d = resp.GetDeployment()
		return err == nil && d.GetFinishedAt() != nil
	}, 3*time.Minute, 200*time.Millisecond, "l'installation ne se termine pas")
	require.Equal(t, forgenetv1.DeploymentState_DEPLOYMENT_STATE_SUCCESS, d.GetState(), d.GetErrorMessage())
	assert.Equal(t, "secret", d.GetTemplate())

	// Le conteneur a reçu les valeurs.
	env := containerEnv(ctx, t, stack)
	assert.Contains(t, env, "APP_SECRET="+secret)
	assert.Contains(t, env, "APP_NAME=forgenet")

	// Ni le journal du déploiement ni la base ne gardent le secret en clair.
	assert.NotContains(t, deploymentLog(authed, t, deployments, id), secret)
	db, err := sqlite.Open(ctx, config.DatabaseConfig{Path: dbPath, BusyTimeout: time.Second})
	require.NoError(t, err)
	stored, err := db.Deployments().GetDeployment(ctx, id)
	require.NoError(t, err)
	_ = db.Close()
	require.NotNil(t, stored.VariablesSealed)
	assert.False(t, bytes.Contains(stored.VariablesSealed, []byte(secret)))

	// Une deuxième installation sous le même nom remplacerait l'application.
	_, err = marketplace.InstallTemplate(authed, &forgenetv1.InstallTemplateRequest{TemplateId: "secret", ProjectId: project})
	assert.Equal(t, codes.AlreadyExists, status.Code(err))

	select {
	case err := <-exited:
		t.Fatalf("le serveur s'est arrêté pendant le test : %v", err)
	default:
	}
}

// containerEnv renvoie l'environnement du conteneur de la stack.
func containerEnv(ctx context.Context, t *testing.T, stack string) []string {
	t.Helper()
	out, err := exec.CommandContext(ctx, "docker", "ps", "-q", "--filter", "label=com.docker.compose.project="+stack).Output()
	require.NoError(t, err)
	found := strings.Fields(string(out))
	require.Len(t, found, 1)
	out, err = exec.CommandContext(ctx, "docker", "inspect", "-f", "{{range .Config.Env}}{{println .}}{{end}}", found[0]).Output()
	require.NoError(t, err)
	return strings.Split(strings.TrimSpace(string(out)), "\n")
}

// deploymentLog lit tout le journal d'un déploiement terminé.
func deploymentLog(ctx context.Context, t *testing.T, client forgenetv1.DeploymentServiceClient, id string) string {
	t.Helper()
	stream, err := client.StreamDeploymentLogs(ctx, &forgenetv1.StreamDeploymentLogsRequest{DeploymentId: id})
	require.NoError(t, err)
	var all strings.Builder
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		all.WriteString(resp.GetLogEntry().GetMessage() + "\n")
	}
	require.NotEmpty(t, all.String(), "le journal a été lu")
	return all.String()
}
