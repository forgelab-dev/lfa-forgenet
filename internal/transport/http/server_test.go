package http_test

import (
	"bytes"
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
	handler := transportHTTP.NewHealthHandler(healthSvc, config.ExposeDetailsAlways, zerolog.Nop())

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
	handler := transportHTTP.NewHealthHandler(healthSvc, config.ExposeDetailsAlways, zerolog.Nop())

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
	handler := transportHTTP.NewHealthHandler(healthSvc, config.ExposeDetailsAlways, zerolog.Nop())

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

	server := transportHTTP.NewServer(cfg, config.HealthConfig{}, healthSvc, logger)

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

// Les sondes passent par le routeur réel du serveur : chaque route répond en GET
// et refuse les autres méthodes (405 et en-tête Allow, gérés par ServeMux).
func TestServer_HealthRoutes_MethodsThroughRealMux(t *testing.T) {
	srv := transportHTTP.NewServer(config.ServerConfig{}, config.HealthConfig{}, health.NewService(), zerolog.Nop())

	for _, path := range []string{"/healthz", "/readyz", "/version"} {
		t.Run(path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			srv.Mux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))

			for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
				rec := httptest.NewRecorder()
				srv.Mux().ServeHTTP(rec, httptest.NewRequest(method, path, nil))
				assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "%s %s", method, path)
				assert.Equal(t, "GET, HEAD", rec.Header().Get("Allow"))
			}
		})
	}
}

// /readyz ne détaille les dépendances qu'aux appelants autorisés par
// health.expose_details ; le statut et le code HTTP restent les mêmes.
func TestReadiness_DetailsExposure(t *testing.T) {
	healthSvc := health.NewService()
	healthSvc.RegisterReadinessChecker("docker", health.CheckFunc(func(context.Context) error {
		return errors.New("dial unix /var/run/docker.sock: connection refused")
	}))

	const (
		loopback = "127.0.0.1:41000"
		ipv6     = "[::1]:41000"
		remote   = "203.0.113.7:41000"
	)
	cases := []struct {
		name, mode, remoteAddr, proxyHeader string
		wantDetails                         bool
	}{
		{"always, requête distante", config.ExposeDetailsAlways, remote, "", true},
		{"never, requête locale", config.ExposeDetailsNever, loopback, "", false},
		{"local, boucle IPv4", config.ExposeDetailsLocal, loopback, "", true},
		{"local, boucle IPv6", config.ExposeDetailsLocal, ipv6, "", true},
		{"local, requête distante", config.ExposeDetailsLocal, remote, "", false},
		{"local, via proxy local (X-Forwarded-For)", config.ExposeDetailsLocal, loopback, "X-Forwarded-For", false},
		{"local, via proxy local (Forwarded)", config.ExposeDetailsLocal, loopback, "Forwarded", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var logs bytes.Buffer
			mux := http.NewServeMux()
			transportHTTP.NewHealthHandler(healthSvc, tc.mode, zerolog.New(&logs)).RegisterRoutes(mux)

			req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
			req.RemoteAddr = tc.remoteAddr
			if tc.proxyHeader != "" {
				req.Header.Set(tc.proxyHeader, "for=198.51.100.1")
			}
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
			var resp transportHTTP.HealthResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			assert.Equal(t, string(health.StatusNotServing), resp.Status)
			if tc.wantDetails {
				assert.Contains(t, resp.Details["docker"], "docker.sock")
			} else {
				assert.Nil(t, resp.Details)
				assert.NotContains(t, rec.Body.String(), "docker.sock")
			}
			assert.Contains(t, logs.String(), "docker.sock", "la cause est toujours journalisée côté serveur")
		})
	}
}
