package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/health"
	transportHTTP "github.com/forgelab-dev/lfa-forgenet/internal/transport/http"
	"github.com/forgelab-dev/lfa-forgenet/internal/version"
)

func TestHealthHandler_Liveness(t *testing.T) {
	healthSvc := health.NewService()
	handler := transportHTTP.NewHealthHandler(healthSvc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Cas 1 : Serveur vivant
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp transportHTTP.HealthResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, string(health.StatusServing), resp.Status)

	// Cas 2 : Serveur non vivant (arrêt en cours)
	healthSvc.SetServing(false)
	req = httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, string(health.StatusNotServing), resp.Status)
}

func TestHealthHandler_Readiness(t *testing.T) {
	healthSvc := health.NewService()
	handler := transportHTTP.NewHealthHandler(healthSvc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Cas 1 : Prêt sans sous-système défaillant
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp transportHTTP.HealthResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, string(health.StatusServing), resp.Status)

	// Cas 2 : Un sous-système critique échoue
	healthSvc.RegisterReadinessChecker("database", health.CheckFunc(func(_ context.Context) error {
		return errors.New("connexion sqlite impossible")
	}))

	req = httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, string(health.StatusNotServing), resp.Status)
	assert.Equal(t, "connexion sqlite impossible", resp.Details["database"])
}

func TestHealthHandler_Version(t *testing.T) {
	healthSvc := health.NewService()
	handler := transportHTTP.NewHealthHandler(healthSvc)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var buildInfo version.BuildInfo
	err := json.Unmarshal(rec.Body.Bytes(), &buildInfo)
	require.NoError(t, err)
	assert.Equal(t, version.Version, buildInfo.Version)
	assert.NotEmpty(t, buildInfo.GoVersion)
}

func TestServer_StartAndShutdown(t *testing.T) {
	healthSvc := health.NewService()
	logger := zerolog.Nop()
	cfg := config.ServerConfig{
		Host:            "127.0.0.1",
		HTTPPort:        0, // Port éphémère
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    5 * time.Second,
		IdleTimeout:     5 * time.Second,
		ShutdownTimeout: 2 * time.Second,
	}

	server := transportHTTP.NewServer(cfg, healthSvc, logger)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	server.SetListener(ln)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start()
	}()

	// Vérification de la connectivité réelle
	addr := server.Addr()
	client := &http.Client{Timeout: 2 * time.Second}

	var res *http.Response
	for i := 0; i < 5; i++ {
		res, err = client.Get("http://" + addr + "/healthz")
		if err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	require.NoError(t, err)
	defer func() {
		_ = res.Body.Close()
	}()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Arrêt gracieux
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	assert.NoError(t, err)

	serverErr := <-errCh
	assert.NoError(t, serverErr)
}
