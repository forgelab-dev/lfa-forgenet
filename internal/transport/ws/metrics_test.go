package ws_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/metrics"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/transport/ws"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

const (
	memberToken = "jeton-membre"
	guestToken  = "jeton-sans-role"
	brokenToken = "base-indisponible"
)

// fakeAuth accepte memberToken tant qu'il n'est pas révoqué.
type fakeAuth struct{ revoked atomic.Bool }

func (a *fakeAuth) ValidateToken(_ context.Context, token string) (*auth.Claims, error) {
	switch {
	case token == memberToken && !a.revoked.Load():
		return &auth.Claims{UserID: "u1", Role: store.RoleMember}, nil
	case token == guestToken:
		return &auth.Claims{UserID: "u2", Role: store.Role("invite")}, nil
	case token == brokenToken:
		return nil, errors.New("base verrouillée")
	default:
		return nil, auth.ErrInvalidToken
	}
}

type hostFunc func(context.Context) (metrics.Host, error)

func (f hostFunc) Collect(ctx context.Context) (metrics.Host, error) { return f(ctx) }

type containerFunc func(context.Context) (metrics.Container, error)

func (f containerFunc) Collect(ctx context.Context) (metrics.Container, error) { return f(ctx) }

var collected = time.Date(2026, 11, 9, 12, 0, 0, 0, time.UTC)

func sampleHost(context.Context) (metrics.Host, error) {
	return metrics.Host{
		CollectedAt:      collected,
		UptimeSeconds:    3600,
		CPU:              metrics.CPU{UsagePercent: 12.5, OnlineCPUs: 4},
		Memory:           metrics.Memory{UsageBytes: 1 << 30, LimitBytes: 4 << 30, UsagePercent: 25},
		Network:          metrics.Network{RxBytes: 1000},
		ActiveContainers: 2,
	}, nil
}

type env struct {
	server  *httptest.Server
	handler *ws.MetricsHandler
	auth    *fakeAuth
	// containers simule le démon : un identifiant absent est introuvable.
	mu         sync.Mutex
	containers map[string]bool
}

func newEnv(t *testing.T, tweak func(*ws.Config, *http.Server)) *env {
	t.Helper()
	e := &env{auth: &fakeAuth{}, containers: map[string]bool{"web": true}}
	cfg := ws.Config{
		Auth:    e.auth,
		NewHost: func() ws.HostCollector { return hostFunc(sampleHost) },
		NewContainer: func(id string) ws.ContainerCollector {
			return containerFunc(func(context.Context) (metrics.Container, error) {
				e.mu.Lock()
				defer e.mu.Unlock()
				if !e.containers[id] {
					return metrics.Container{}, fmt.Errorf("statistiques : %w", orchestrator.ErrNotFound)
				}
				return metrics.Container{ContainerID: id, CollectedAt: collected, CPU: metrics.CPU{UsagePercent: 150, OnlineCPUs: 2}}, nil
			})
		},
		Logger: zerolog.Nop(),
	}
	mux := http.NewServeMux()
	srv := httptest.NewUnstartedServer(mux)
	if tweak != nil {
		tweak(&cfg, srv.Config)
	}
	e.handler = ws.NewMetricsHandler(cfg)
	e.handler.RegisterRoutes(mux)
	srv.Start()
	t.Cleanup(func() {
		e.handler.Close()
		srv.Close()
	})
	e.server = srv
	return e
}

func (e *env) url(path string) string {
	return "ws" + strings.TrimPrefix(e.server.URL, "http") + path
}

// handshake retient de la réponse HTTP ce que les tests vérifient.
type handshake struct {
	status int
	header http.Header
}

// dial ouvre un flux comme le dashboard : jeton dans les sous-protocoles.
func (e *env) dial(t *testing.T, path, token string) (*websocket.Conn, handshake, error) {
	t.Helper()
	return e.dialWith(t, path, &websocket.DialOptions{
		Subprotocols:    []string{ws.Subprotocol, ws.BearerSubprotocolPrefix + token},
		CompressionMode: websocket.CompressionContextTakeover,
	})
}

func (e *env) dialWith(t *testing.T, path string, opts *websocket.DialOptions) (*websocket.Conn, handshake, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, resp, err := websocket.Dial(ctx, e.url(path), opts)
	var hs handshake
	if resp != nil {
		hs = handshake{status: resp.StatusCode, header: resp.Header}
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}
	if conn != nil {
		t.Cleanup(func() { _ = conn.CloseNow() })
	}
	return conn, hs, err
}

func TestSystemStreamSendsContractMessagesCompressed(t *testing.T) {
	e := newEnv(t, nil)

	conn, hs, err := e.dial(t, "/ws/v1/metrics/system?interval=1", memberToken)
	require.NoError(t, err)

	assert.Equal(t, ws.Subprotocol, conn.Subprotocol(), "le jeton n'est jamais renvoyé comme sous-protocole")
	assert.Contains(t, hs.header.Get("Sec-WebSocket-Extensions"), "permessage-deflate", "flux compressé")

	first := readSystem(t, conn)
	assert.InDelta(t, 12.5, first.GetMetrics().GetCpu().GetCpuUsagePercent(), 0.001)
	assert.Equal(t, uint32(2), first.GetMetrics().GetActiveContainersCount())
	assert.Equal(t, uint64(1<<30), first.GetMetrics().GetMemory().GetUsageBytes())
	assert.Equal(t, collected, first.GetMetrics().GetCollectedAt().AsTime())

	start := time.Now()
	readSystem(t, conn)
	assert.InDelta(t, time.Second, time.Since(start), float64(500*time.Millisecond), "une mesure par intervalle demandé")
}

func TestStreamAcceptsAuthorizationHeader(t *testing.T) {
	e := newEnv(t, nil)

	conn, _, err := e.dialWith(t, "/ws/v1/metrics/containers/web", &websocket.DialOptions{
		HTTPHeader: http.Header{"Authorization": {"Bearer " + memberToken}},
	})
	require.NoError(t, err)

	msg := readContainer(t, conn)
	assert.Equal(t, "web", msg.GetMetrics().GetContainerId())
	assert.InDelta(t, 150.0, msg.GetMetrics().GetCpu().GetCpuUsagePercent(), 0.001)
}

func TestStreamRefusalsAreHTTPErrors(t *testing.T) {
	e := newEnv(t, nil)
	cases := []struct {
		name, path, token string
		status            int
	}{
		{"sans jeton", "/ws/v1/metrics/system", "", http.StatusUnauthorized},
		{"jeton invalide", "/ws/v1/metrics/system", "faux", http.StatusUnauthorized},
		{"rôle insuffisant", "/ws/v1/metrics/system", guestToken, http.StatusForbidden},
		{"vérification impossible", "/ws/v1/metrics/system", brokenToken, http.StatusServiceUnavailable},
		{"intervalle trop long", "/ws/v1/metrics/system?interval=61", memberToken, http.StatusBadRequest},
		{"intervalle négatif", "/ws/v1/metrics/system?interval=-1", memberToken, http.StatusBadRequest},
		{"conteneur inconnu", "/ws/v1/metrics/containers/absent", memberToken, http.StatusNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, hs, err := e.dial(t, tc.path, tc.token)
			require.Error(t, err)
			assert.Equal(t, tc.status, hs.status)
		})
	}
}

func TestStreamRejectsForeignOrigin(t *testing.T) {
	e := newEnv(t, nil)

	_, hs, err := e.dialWith(t, "/ws/v1/metrics/system", &websocket.DialOptions{
		HTTPHeader: http.Header{"Authorization": {"Bearer " + memberToken}, "Origin": {"https://autre-site.example"}},
	})

	require.Error(t, err)
	assert.Equal(t, http.StatusForbidden, hs.status, "un site tiers ne lit pas les métriques avec le jeton du visiteur")
}

func TestStreamLimitsConcurrentStreams(t *testing.T) {
	e := newEnv(t, func(c *ws.Config, _ *http.Server) { c.MaxStreams = 1 })

	_, _, err := e.dial(t, "/ws/v1/metrics/system", memberToken)
	require.NoError(t, err)
	_, hs, err := e.dial(t, "/ws/v1/metrics/system", memberToken)

	require.Error(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, hs.status)
}

func TestStreamClosesWhenContainerDisappears(t *testing.T) {
	e := newEnv(t, nil)
	conn, _, err := e.dial(t, "/ws/v1/metrics/containers/web?interval=1", memberToken)
	require.NoError(t, err)
	readContainer(t, conn)

	e.mu.Lock()
	delete(e.containers, "web")
	e.mu.Unlock()

	assert.Equal(t, ws.CloseNotFound, closeStatus(t, conn))
}

func TestStreamClosesWhenTokenIsRevoked(t *testing.T) {
	e := newEnv(t, func(c *ws.Config, _ *http.Server) { c.RecheckEvery = time.Millisecond })
	conn, _, err := e.dial(t, "/ws/v1/metrics/system?interval=1", memberToken)
	require.NoError(t, err)
	readSystem(t, conn)

	e.auth.revoked.Store(true)

	assert.Equal(t, ws.CloseUnauthorized, closeStatus(t, conn))
}

// Le délai d'écriture du serveur HTTP (15 s par défaut) ne doit pas couper
// un flux : réduit ici à 300 ms, le flux dure plus d'une seconde.
func TestStreamOutlivesServerWriteTimeout(t *testing.T) {
	e := newEnv(t, func(_ *ws.Config, s *http.Server) { s.WriteTimeout = 300 * time.Millisecond })
	conn, _, err := e.dial(t, "/ws/v1/metrics/system?interval=1", memberToken)
	require.NoError(t, err)

	readSystem(t, conn)
	readSystem(t, conn)
	readSystem(t, conn)
}

func TestCloseEndsStreamsWithGoingAway(t *testing.T) {
	e := newEnv(t, nil)
	conn, _, err := e.dial(t, "/ws/v1/metrics/system", memberToken)
	require.NoError(t, err)
	readSystem(t, conn)

	statuses := make(chan websocket.StatusCode, 1)
	go func() { statuses <- closeStatus(t, conn) }()
	e.handler.Close()

	assert.Equal(t, websocket.StatusGoingAway, <-statuses)
}

func readSystem(t *testing.T, conn *websocket.Conn) *forgenetv1.StreamSystemMetricsResponse {
	t.Helper()
	return readMessage(t, conn, &forgenetv1.StreamSystemMetricsResponse{})
}

func readContainer(t *testing.T, conn *websocket.Conn) *forgenetv1.StreamContainerMetricsResponse {
	t.Helper()
	return readMessage(t, conn, &forgenetv1.StreamContainerMetricsResponse{})
}

func readMessage[M proto.Message](t *testing.T, conn *websocket.Conn, msg M) M {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	typ, payload, err := conn.Read(ctx)
	require.NoError(t, err)
	require.Equal(t, websocket.MessageText, typ)
	require.NoError(t, protojson.Unmarshal(payload, msg))
	return msg
}

// closeStatus lit jusqu'à la fermeture et renvoie son code.
func closeStatus(t *testing.T, conn *websocket.Conn) websocket.StatusCode {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for {
		if _, _, err := conn.Read(ctx); err != nil {
			return websocket.CloseStatus(err)
		}
	}
}
