package connectapi_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
	"github.com/forgelab-dev/lfa-forgenet/internal/transport/connectapi"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

const password = "motdepasse-de-test"

// containers est un faux service : il renvoie l'identité vue par le handler,
// pour prouver que le garde la transmet à travers le pont vanguardgrpc.
type containers struct {
	forgenetv1.UnimplementedContainerServiceServer
}

func (containers) ListContainers(ctx context.Context, _ *forgenetv1.ListContainersRequest) (*forgenetv1.ListContainersResponse, error) {
	claims, _ := auth.UserFromContext(ctx)
	return &forgenetv1.ListContainersResponse{Containers: []*forgenetv1.Container{{Id: "c1", Name: claims.Username}}}, nil
}

func (containers) StartContainer(context.Context, *forgenetv1.StartContainerRequest) (*forgenetv1.StartContainerResponse, error) {
	return &forgenetv1.StartContainerResponse{}, nil
}

func (containers) GetContainer(context.Context, *forgenetv1.GetContainerRequest) (*forgenetv1.GetContainerResponse, error) {
	panic("défaillance simulée")
}

// StreamContainerLogs dure plus longtemps que le délai d'écriture du serveur
// de test (300 ms).
func (containers) StreamContainerLogs(_ *forgenetv1.StreamContainerLogsRequest, s grpc.ServerStreamingServer[forgenetv1.StreamContainerLogsResponse]) error {
	for i := range 3 {
		if i > 0 {
			time.Sleep(250 * time.Millisecond)
		}
		if err := s.Send(&forgenetv1.StreamContainerLogsResponse{Chunk: []byte{byte('a' + i)}}); err != nil {
			return err
		}
	}
	return nil
}

// testServer : la porte Connect devant une vraie base, un compte admin et un
// compte membre, avec leurs jetons.
type testServer struct {
	url           string
	member, admin string
}

// startServer démarre la porte. writeTimeout court (300 ms) sert au test des
// flux ; il est trop court pour Login, dont le hachage du mot de passe prend
// plus d'une seconde sous -race.
func startServer(t *testing.T, writeTimeout time.Duration) testServer {
	t.Helper()
	ctx := context.Background()
	st, err := sqlite.Open(ctx, config.DatabaseConfig{Path: ":memory:", BusyTimeout: 2 * time.Second})
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })

	authSvc := auth.NewService(st, config.AuthConfig{JWTSecret: "secret-de-test-d-au-moins-32-caracteres"}, audit.NewAuditor(st.Audit(), zerolog.Nop()))
	token := func(name string, role store.Role) string {
		u, err := authSvc.Register(ctx, name, name+"@forgenet.test", password, role)
		require.NoError(t, err)
		tok, _, err := authSvc.CreateAPIToken(ctx, u.ID, name, "", time.Hour)
		require.NoError(t, err)
		return tok
	}
	ts := testServer{admin: token("admin", store.RoleAdmin), member: token("membre", store.RoleMember)}

	api := connectapi.New(auth.NewGuard(authSvc, auth.DefaultPolicy(), zerolog.Nop()), zerolog.Nop())
	forgenetv1.RegisterAuthServiceServer(api.Registrar(), transportGRPC.NewAuthServiceServer(authSvc, st.Users(), zerolog.Nop()))
	forgenetv1.RegisterContainerServiceServer(api.Registrar(), containers{})
	mux := http.NewServeMux()
	api.Mount(mux)

	srv := httptest.NewUnstartedServer(mux)
	srv.Config.WriteTimeout = writeTimeout
	srv.Start()
	t.Cleanup(srv.Close)
	ts.url = srv.URL
	return ts
}

// call fait un appel Connect JSON, comme le dashboard.
func call(t *testing.T, base, procedure, token string, body any) (int, map[string]any) {
	t.Helper()
	payload, err := json.Marshal(body)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, base+procedure, bytes.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = res.Body.Close() }()
	var out map[string]any
	require.NoError(t, json.NewDecoder(res.Body).Decode(&out))
	return res.StatusCode, out
}

func TestConnectLoginIsPublic(t *testing.T) {
	ts := startServer(t, time.Minute)

	status, out := call(t, ts.url, "/forgenet.v1.AuthService/Login", "", map[string]string{"login": "membre", "password": password})
	require.Equal(t, http.StatusOK, status, out)
	token := out["token"].(string)

	status, out = call(t, ts.url, "/forgenet.v1.AuthService/GetCurrentUser", token, map[string]any{})
	require.Equal(t, http.StatusOK, status, out)
	assert.Equal(t, "membre", out["user"].(map[string]any)["username"])
}

func TestConnectDoorAppliesThePolicy(t *testing.T) {
	ts := startServer(t, time.Minute)
	base, member, admin := ts.url, ts.member, ts.admin

	status, out := call(t, base, "/forgenet.v1.ContainerService/ListContainers", "", map[string]any{})
	assert.Equal(t, http.StatusUnauthorized, status)
	assert.Equal(t, "unauthenticated", out["code"])

	status, out = call(t, base, "/forgenet.v1.ContainerService/ListContainers", "jeton-invalide", map[string]any{})
	assert.Equal(t, http.StatusUnauthorized, status, out)

	status, out = call(t, base, "/forgenet.v1.ContainerService/ListContainers", member, map[string]any{})
	require.Equal(t, http.StatusOK, status, out)
	first := out["containers"].([]any)[0].(map[string]any)
	assert.Equal(t, "membre", first["name"], "le handler reçoit l'identité posée par le garde")

	status, out = call(t, base, "/forgenet.v1.ContainerService/StartContainer", member, map[string]string{"containerId": "c1"})
	assert.Equal(t, http.StatusForbidden, status)
	assert.Equal(t, "permission_denied", out["code"])

	status, _ = call(t, base, "/forgenet.v1.ContainerService/StartContainer", admin, map[string]string{"containerId": "c1"})
	assert.Equal(t, http.StatusOK, status)

	status, out = call(t, base, "/forgenet.v1.AuthService/GetCurrentUser", member, map[string]any{})
	require.Equal(t, http.StatusOK, status)
	assert.Equal(t, "USER_ROLE_MEMBER", out["user"].(map[string]any)["role"])
}

func TestConnectDoorRecoversFromPanics(t *testing.T) {
	ts := startServer(t, time.Minute)
	base, member := ts.url, ts.member

	status, out := call(t, base, "/forgenet.v1.ContainerService/GetContainer", member, map[string]string{"containerId": "c1"})

	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, "erreur interne du serveur", out["message"], "le contenu de la panique reste dans le journal")
	status, _ = call(t, base, "/forgenet.v1.ContainerService/ListContainers", member, map[string]any{})
	assert.Equal(t, http.StatusOK, status, "le serveur répond toujours")
}

// Un flux Connect (application/connect+json) dure plus longtemps que le
// délai d'écriture du serveur : la porte le lève sur les flux.
func TestConnectStreamOutlivesWriteTimeout(t *testing.T) {
	ts := startServer(t, 300*time.Millisecond)
	base, member := ts.url, ts.member

	req, err := http.NewRequest(http.MethodPost, base+"/forgenet.v1.ContainerService/StreamContainerLogs",
		bytes.NewReader(envelope([]byte(`{"containerId":"c1"}`))))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/connect+json")
	req.Header.Set("Authorization", "Bearer "+member)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = res.Body.Close() }()
	require.Equal(t, http.StatusOK, res.StatusCode)

	var chunks []string
	for {
		flags, payload := readEnvelope(t, res.Body)
		if flags&0x02 != 0 { // fin de flux
			assert.JSONEq(t, `{}`, string(payload), "fin de flux sans erreur")
			break
		}
		var msg map[string]string
		require.NoError(t, json.Unmarshal(payload, &msg))
		chunks = append(chunks, msg["chunk"])
	}
	assert.Equal(t, []string{"YQ==", "Yg==", "Yw=="}, chunks, "a, b et c en base64")
}

// gRPC-Web, l'autre protocole du navigateur, passe par la même porte.
func TestGRPCWebRequest(t *testing.T) {
	ts := startServer(t, time.Minute)
	base, member := ts.url, ts.member
	msg, err := proto.Marshal(&forgenetv1.ListContainersRequest{})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, base+"/forgenet.v1.ContainerService/ListContainers", bytes.NewReader(envelope(msg)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/grpc-web+proto")
	req.Header.Set("Authorization", "Bearer "+member)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = res.Body.Close() }()
	require.Equal(t, http.StatusOK, res.StatusCode)

	_, payload := readEnvelope(t, res.Body)
	var out forgenetv1.ListContainersResponse
	require.NoError(t, proto.Unmarshal(payload, &out))
	assert.Equal(t, "membre", out.GetContainers()[0].GetName())
	_, trailer := readEnvelope(t, res.Body)
	assert.Contains(t, strings.ToLower(string(trailer)), "grpc-status: 0")
}

func envelope(payload []byte) []byte {
	out := make([]byte, 5, 5+len(payload))
	binary.BigEndian.PutUint32(out[1:], uint32(len(payload))) //nolint:gosec // message de test
	return append(out, payload...)
}

func readEnvelope(t *testing.T, r io.Reader) (byte, []byte) {
	t.Helper()
	head := make([]byte, 5)
	_, err := io.ReadFull(r, head)
	require.NoError(t, err)
	payload := make([]byte, binary.BigEndian.Uint32(head[1:]))
	_, err = io.ReadFull(r, payload)
	require.NoError(t, err)
	return head[0], payload
}
