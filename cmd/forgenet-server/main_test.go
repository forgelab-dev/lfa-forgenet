package main_test

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/quotedprintable"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	grpcHealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
	transportHTTP "github.com/forgelab-dev/lfa-forgenet/internal/transport/http"
	"github.com/forgelab-dev/lfa-forgenet/internal/version"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
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
	webRoot := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(webRoot, "index.html"), []byte("<title>ForgeNet</title>"), 0o600))

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
web:
  root: %s
`, grpcPort, webRoot))
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

	t.Run("dashboard et API à la même origine", func(t *testing.T) {
		resp, err := client.Get(httpAddr + "/domaines")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Content-Security-Policy"), "default-src 'self'")

		// Les routes de l'API restent servies par l'API, jamais par la page.
		api, err := client.Post(httpAddr+"/forgenet.v1.ContainerService/ListContainers", "application/json", strings.NewReader("{}"))
		require.NoError(t, err)
		defer func() { _ = api.Body.Close() }()
		assert.Equal(t, http.StatusUnauthorized, api.StatusCode)
		assert.Empty(t, api.Header.Get("Content-Security-Policy"))
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

		// L'API des conteneurs n'est jamais servie sans authentification.
		_, err = forgenetv1.NewContainerServiceClient(conn).ListContainers(ctx, &forgenetv1.ListContainersRequest{})
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
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
		// L'API d'administration de Caddy n'a pas d'authentification : une
		// adresse exposée empêche le démarrage (ADR-0007).
		"administration de Caddy exposée": fmt.Sprintf(`
server:
  http_port: %d
  grpc_port: %d
  shutdown_timeout: 1s
proxy:
  enabled: true
  admin_address: 0.0.0.0:2019
`, freePort(t), freePort(t)),
		"dossier du dashboard sans build": fmt.Sprintf(`
server:
  http_port: %d
  grpc_port: %d
  shutdown_timeout: 1s
web:
  root: %s
`, freePort(t), freePort(t), t.TempDir()),
		"cible du dashboard invalide": fmt.Sprintf(`
server:
  http_port: %d
  grpc_port: %d
  shutdown_timeout: 1s
proxy:
  enabled: true
  admin_address: 127.0.0.1:2019
  dashboard_host: forgenet.example.org
  dashboard_upstream: "forgenet-server:8080\nrespond"
`, freePort(t), freePort(t)),
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

func TestUserSetPassword_FromCommandLine(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "forgenet.db")
	configFile := writeConfig(t, fmt.Sprintf("database:\n  path: %q\n", dbPath))
	run := func(stdin string, args ...string) (string, int) {
		cmd := exec.Command(binPath, append([]string{"user"}, append(args, "--config", configFile)...)...)
		cmd.Stdin = strings.NewReader(stdin)
		out, err := cmd.CombinedOutput()
		return string(out), exitCode(t, err)
	}

	out, code := run("motdepasse-de-test\n", "create", "--username", "admin", "--email", "admin@forgenet.test")
	require.Equal(t, 0, code, out)

	out, code = run("mot-de-passe-retrouve\n", "set-password", "--username", "admin")
	require.Equal(t, 0, code, out)
	assert.Contains(t, out, "Mot de passe de admin changé")
	assert.NotContains(t, out, "mot-de-passe-retrouve", "le mot de passe ne doit jamais être affiché")

	out, code = run("mot-de-passe-retrouve\n", "set-password", "--username", "personne")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "aucun compte")

	out, code = run("court\n", "set-password", "--username", "admin")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "8 caractères")

	out, code = run("", "set-password")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "--username est obligatoire")

	db, err := sqlite.Open(context.Background(), config.DatabaseConfig{Path: dbPath, BusyTimeout: time.Second})
	require.NoError(t, err)
	defer func() { _ = db.Close() }()
	user, err := db.Users().GetUserByUsername(context.Background(), "admin")
	require.NoError(t, err)
	assert.NoError(t, auth.CheckPassword(user.PasswordHash, "mot-de-passe-retrouve"))
}

func TestUserLoginLink_PendingAdmin(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "forgenet.db")
	configFile := writeConfig(t, fmt.Sprintf("database:\n  path: %q\n", dbPath))
	run := func(args ...string) (string, int) {
		cmd := exec.Command(binPath, append([]string{"user"}, append(args, "--config", configFile)...)...)
		cmd.Stdin = strings.NewReader("")
		out, err := cmd.CombinedOutput()
		return string(out), exitCode(t, err)
	}

	out, code := run("create", "--username", "admin", "--email", "admin@forgenet.test", "--no-password")
	require.Equal(t, 0, code, out)
	assert.Contains(t, out, "forgenet user login-link --username admin")

	out, code = run("login-link", "--username", "admin")
	require.Equal(t, 0, code, out)
	assert.Regexp(t, `  http://127\.0\.0\.1:8080/acces#[A-Za-z0-9_-]{43}\n`, out)
	assert.Contains(t, out, "ssh -N -L 8080:127.0.0.1:8080", "sans domaine, le tunnel est expliqué")
	token := out[strings.Index(out, "#")+1:]
	token = token[:strings.IndexByte(token, '\n')]

	out, code = run("login-link", "--username", "admin", "--send")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "n'est pas configuré")

	out, code = run("login-link", "--username", "personne")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "aucun compte")

	db, err := sqlite.Open(context.Background(), config.DatabaseConfig{Path: dbPath, BusyTimeout: time.Second})
	require.NoError(t, err)
	defer func() { _ = db.Close() }()
	user, err := db.Users().GetUserByUsername(context.Background(), "admin")
	require.NoError(t, err)
	assert.Equal(t, auth.UnusablePassword, user.PasswordHash)
	_, err = db.ActionTokens().GetActiveActionToken(context.Background(), auth.HashToken(token), time.Now())
	require.NoError(t, err, "le lien affiché est enregistré par son empreinte")
}

func TestUserLoginLink_SendsTheLinkByMail(t *testing.T) {
	relay := newPlainRelay(t)
	dbPath := filepath.Join(t.TempDir(), "forgenet.db")
	configFile := writeConfig(t, fmt.Sprintf(`database:
  path: %q
proxy:
  enabled: true
  dashboard_host: forgenet.exemple.org
  dashboard_upstream: forgenet-server:8080
smtp:
  host: 127.0.0.1
  port: %d
  tls: none
  from: forgenet@exemple.org
`, dbPath, relay.port))
	run := func(args ...string) (string, int) {
		cmd := exec.Command(binPath, append([]string{"user"}, append(args, "--config", configFile)...)...)
		out, err := cmd.CombinedOutput()
		return string(out), exitCode(t, err)
	}

	out, code := run("create", "--username", "admin", "--email", "admin@forgenet.test", "--no-password")
	require.Equal(t, 0, code, out)
	out, code = run("login-link", "--username", "admin", "--send")
	require.Equal(t, 0, code, out)
	assert.Contains(t, out, "envoyé à l'adresse du compte admin")
	assert.NotContains(t, out, "acces#", "le lien envoyé n'est pas affiché")
	assert.NotContains(t, out, "admin@forgenet.test", "l'adresse n'est pas affichée")

	data := relay.received()
	assert.Contains(t, data, "To: <admin@forgenet.test>")
	_, body, ok := strings.Cut(data, "\r\n\r\n")
	require.True(t, ok, data)
	decoded, err := io.ReadAll(quotedprintable.NewReader(strings.NewReader(body)))
	require.NoError(t, err)
	assert.Regexp(t, `https://forgenet\.exemple\.org/acces#[A-Za-z0-9_-]{43}`, string(decoded))
	assert.NotContains(t, string(decoded), "ssh -N", "avec un domaine, pas de tunnel")
}

// plainRelay est un relais SMTP local sans TLS (smtp.tls=none), qui garde le
// dernier message reçu.
type plainRelay struct {
	port int
	mu   sync.Mutex
	data string
}

func newPlainRelay(t *testing.T) *plainRelay {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = ln.Close() })
	r := &plainRelay{port: ln.Addr().(*net.TCPAddr).Port}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go r.handle(conn)
		}
	}()
	return r
}

func (r *plainRelay) handle(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	rd := bufio.NewReader(conn)
	say := func(s string) { _, _ = conn.Write([]byte(s + "\r\n")) }
	say("220 relais local")
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			return
		}
		switch strings.ToUpper(strings.SplitN(strings.TrimSpace(line), " ", 2)[0]) {
		case "EHLO", "HELO", "MAIL", "RCPT":
			say("250 OK")
		case "DATA":
			say("354 envoyez")
			var b strings.Builder
			for {
				l, err := rd.ReadString('\n')
				if err != nil || l == ".\r\n" {
					break
				}
				b.WriteString(l)
			}
			r.mu.Lock()
			r.data = b.String()
			r.mu.Unlock()
			say("250 reçu")
		case "QUIT":
			say("221 au revoir")
			return
		default:
			say("502 inconnu")
		}
	}
}

func (r *plainRelay) received() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.data
}

func TestUserSuspendResumeDelete_FromCommandLine(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "forgenet.db")
	configFile := writeConfig(t, fmt.Sprintf("database:\n  path: %q\n", dbPath))
	run := func(stdin string, args ...string) (string, int) {
		cmd := exec.Command(binPath, append([]string{"user"}, append(args, "--config", configFile)...)...)
		cmd.Stdin = strings.NewReader(stdin)
		out, err := cmd.CombinedOutput()
		return string(out), exitCode(t, err)
	}
	for _, name := range []string{"admin", "bob"} {
		out, code := run("motdepasse-de-test\n", "create", "--username", name, "--email", name+"@forgenet.test")
		require.Equal(t, 0, code, out)
	}

	// Le serveur garde un administrateur actif : bob est admin aussi, donc
	// admin peut être suspendu, mais pas les deux.
	out, code := run("", "suspend", "--username", "admin")
	require.Equal(t, 0, code, out)
	assert.Contains(t, out, "Compte admin suspendu")
	assert.Contains(t, out, "forgenet user resume --username admin")

	out, code = run("", "suspend", "--username", "bob")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "dernier administrateur actif")

	out, code = run("", "login-link", "--username", "admin")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "suspendu")

	out, code = run("", "resume", "--username", "admin")
	require.Equal(t, 0, code, out)
	assert.Contains(t, out, "réactivé")

	out, code = run("bobb\n", "delete", "--username", "bob")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "rien n'est supprimé")

	out, code = run("bob\n", "delete", "--username", "bob")
	require.Equal(t, 0, code, out)
	assert.Contains(t, out, "Compte bob supprimé")

	out, code = run("", "suspend", "--username", "personne")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "aucun compte")

	db, err := sqlite.Open(context.Background(), config.DatabaseConfig{Path: dbPath, BusyTimeout: time.Second})
	require.NoError(t, err)
	defer func() { _ = db.Close() }()
	_, err = db.Users().GetUserByUsername(context.Background(), "bob")
	require.ErrorIs(t, err, store.ErrNotFound)
	admin, err := db.Users().GetUserByUsername(context.Background(), "admin")
	require.NoError(t, err)
	assert.False(t, admin.Suspended())
}

func TestUserSetEmail_FromCommandLine(t *testing.T) {
	relay := newPlainRelay(t)
	dbPath := filepath.Join(t.TempDir(), "forgenet.db")
	configFile := writeConfig(t, fmt.Sprintf(`database:
  path: %q
smtp:
  host: 127.0.0.1
  port: %d
  tls: none
  from: forgenet@exemple.org
`, dbPath, relay.port))
	run := func(stdin string, args ...string) (string, int) {
		cmd := exec.Command(binPath, append([]string{"user"}, append(args, "--config", configFile)...)...)
		cmd.Stdin = strings.NewReader(stdin)
		out, err := cmd.CombinedOutput()
		return string(out), exitCode(t, err)
	}
	out, code := run("motdepasse-de-test\n", "create", "--username", "admin", "--email", "admin@forgenet.test")
	require.Equal(t, 0, code, out)
	out, code = run("", "create", "--username", "bio", "--email", "bio@mal-saisie.test", "--role", "member", "--no-password")
	require.Equal(t, 0, code, out)

	out, code = run("", "set-email", "--username", "admin", "--email", "Admin@Nouveau.test")
	require.Equal(t, 0, code, out)
	assert.Contains(t, out, "Adresse du compte admin changée")
	assert.Contains(t, out, "avis est envoyé à l'ancienne adresse")
	assert.NotContains(t, out, "admin@forgenet.test", "les adresses ne sont pas affichées : la sortie peut finir dans un journal")

	data := relay.received()
	assert.Contains(t, data, "To: <admin@forgenet.test>", "l'avis part à l'ancienne adresse")
	_, body, ok := strings.Cut(data, "\r\n\r\n")
	require.True(t, ok, data)
	decoded, err := io.ReadAll(quotedprintable.NewReader(strings.NewReader(body)))
	require.NoError(t, err)
	assert.Contains(t, string(decoded), "a été changée sur le serveur")
	assert.Contains(t, string(decoded), "a***@nouveau.test")

	// Un compte qui attend son mot de passe : pas d'avis (l'ancienne adresse
	// n'a jamais servi), mais son lien est annulé et la commande le dit.
	out, code = run("", "set-email", "--username", "bio", "--email", "bio@bonne.test")
	require.Equal(t, 0, code, out)
	assert.NotContains(t, out, "avis est envoyé")
	assert.Contains(t, out, "forgenet user login-link --username bio")

	for _, c := range []struct{ email, want string }{
		{"admin@nouveau.test", "déjà l'adresse"},
		{"bio@bonne.test", "a déjà un compte"},
		{"pas-une-adresse", "n'est pas valide"},
	} {
		out, code = run("", "set-email", "--username", "admin", "--email", c.email)
		assert.Equal(t, 1, code, c.email)
		assert.Contains(t, out, c.want, c.email)
	}
	out, code = run("", "set-email", "--username", "personne", "--email", "libre@forgenet.test")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "aucun compte")
	out, code = run("", "set-email", "--username", "admin")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "--email est obligatoire")

	// #100 : login-link ne dit plus qu'un compte qui a son mot de passe
	// l'attend encore.
	out, code = run("", "login-link", "--username", "admin")
	require.Equal(t, 0, code, out)
	assert.Contains(t, out, "a déjà un mot de passe")
	assert.NotContains(t, out, "attend son mot de passe")

	db, err := sqlite.Open(context.Background(), config.DatabaseConfig{Path: dbPath, BusyTimeout: time.Second})
	require.NoError(t, err)
	defer func() { _ = db.Close() }()
	admin, err := db.Users().GetUserByEmail(context.Background(), "admin@nouveau.test")
	require.NoError(t, err)
	assert.Equal(t, "admin", admin.Username)
}

// Les commandes du serveur lisent les réglages d'envoi du dashboard avant
// ceux de .env, et « mail reset » rend la main à .env (ADR-0015).
func TestMailReset_GivesTheServerSettingsBack(t *testing.T) {
	relay := newPlainRelay(t)
	dbPath := filepath.Join(t.TempDir(), "forgenet.db")
	configFile := writeConfig(t, fmt.Sprintf(`database:
  path: %q
smtp:
  host: 127.0.0.1
  port: %d
  tls: none
  from: forgenet@exemple.org
secrets:
  key: %s
`, dbPath, relay.port, strings.Repeat("5a", 32)))
	run := func(args ...string) (string, int) {
		cmd := exec.Command(binPath, append(args, "--config", configFile)...)
		out, err := cmd.CombinedOutput()
		return string(out), exitCode(t, err)
	}
	out, code := run("user", "create", "--username", "admin", "--email", "admin@forgenet.test", "--no-password")
	require.Equal(t, 0, code, out)

	// L'envoi a été coupé depuis le dashboard : il l'emporte sur .env.
	db, err := sqlite.Open(context.Background(), config.DatabaseConfig{Path: dbPath, BusyTimeout: time.Second})
	require.NoError(t, err)
	require.NoError(t, db.MailSettings().SaveMailSettings(context.Background(), &store.MailSettings{
		Disabled: true, Host: "127.0.0.1", Port: relay.port, TLS: "none", From: "forgenet@exemple.org",
	}))
	require.NoError(t, db.Close())

	out, code = run("user", "login-link", "--username", "admin", "--send")
	assert.Equal(t, 1, code, out)
	assert.Contains(t, out, "pas configuré")

	out, code = run("mail", "reset")
	require.Equal(t, 0, code, out)
	assert.Contains(t, out, "effacés")

	out, code = run("user", "login-link", "--username", "admin", "--send")
	require.Equal(t, 0, code, out)
	assert.Contains(t, relay.received(), "To: <admin@forgenet.test>", ".env s'applique de nouveau")

	out, code = run("mail", "reset")
	require.Equal(t, 0, code, out)
	assert.Contains(t, out, "Aucun réglage")

	out, code = run("mail")
	assert.Equal(t, 1, code)
	assert.Contains(t, out, "forgenet-server mail reset")
}
