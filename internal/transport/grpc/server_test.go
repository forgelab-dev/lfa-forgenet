package grpc_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcHealth "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/health"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
)

func TestGRPCServer_HealthCheck(t *testing.T) {
	healthSvc := health.NewService()
	logger := zerolog.Nop()
	cfg := config.ServerConfig{
		Host:            "127.0.0.1",
		GRPCPort:        0,
		ShutdownTimeout: 2 * time.Second,
	}

	// Les sondes sont publiques : le garde ne consulte pas le service d'authentification.
	guard := auth.NewGuard(nil, auth.DefaultPolicy(), logger)
	server := transportGRPC.NewServer(cfg, healthSvc, guard, logger)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	server.SetListener(ln)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start()
	}()

	// Connexion client gRPC vers le serveur de test
	conn, err := grpc.NewClient(
		server.Addr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer func() {
		_ = conn.Close()
	}()

	client := grpcHealth.NewHealthClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 1. Liveness standard
	resp, err := client.Check(ctx, &grpcHealth.HealthCheckRequest{})
	require.NoError(t, err)
	assert.Equal(t, grpcHealth.HealthCheckResponse_SERVING, resp.GetStatus())

	// 2. Readiness sans dépendance défaillante
	resp, err = client.Check(ctx, &grpcHealth.HealthCheckRequest{Service: "readiness"})
	require.NoError(t, err)
	assert.Equal(t, grpcHealth.HealthCheckResponse_SERVING, resp.GetStatus())

	// 3. Défaillance d'une dépendance
	healthSvc.RegisterReadinessChecker("vault", health.CheckFunc(func(_ context.Context) error {
		return errors.New("vault scellé")
	}))

	resp, err = client.Check(ctx, &grpcHealth.HealthCheckRequest{Service: "readiness"})
	require.NoError(t, err)
	assert.Equal(t, grpcHealth.HealthCheckResponse_NOT_SERVING, resp.GetStatus())

	// 4. Liveness toujours OK même si readiness est dégradée
	resp, err = client.Check(ctx, &grpcHealth.HealthCheckRequest{})
	require.NoError(t, err)
	assert.Equal(t, grpcHealth.HealthCheckResponse_SERVING, resp.GetStatus())

	// 5. Arrêt du serveur (not serving)
	healthSvc.SetServing(false)
	resp, err = client.Check(ctx, &grpcHealth.HealthCheckRequest{})
	require.NoError(t, err)
	assert.Equal(t, grpcHealth.HealthCheckResponse_NOT_SERVING, resp.GetStatus())

	// 6. Arrêt gracieux
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	err = server.Shutdown(shutdownCtx)
	assert.NoError(t, err)

	serverErr := <-errCh
	assert.NoError(t, serverErr)
}
