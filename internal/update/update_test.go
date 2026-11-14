package update

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const digest = "sha256:3539049c40232701e8dae4dabc62b99cc961dd39f569160365a76bc8a06491b3"

func fiche(version string) string {
	return fmt.Sprintf(`{"schema":1,"version":%q,"image":%q,"digest":%q,"notes_url":"https://github.com/forgelab-dev/lfa-forgenet/releases/tag/v%s"}`,
		version, OfficialImage, digest, version)
}

func TestParseRelease(t *testing.T) {
	r, err := ParseRelease([]byte(fiche("0.2.0")))
	require.NoError(t, err)
	assert.Equal(t, "ghcr.io/forgelab-dev/forgenet-server:0.2.0@"+digest, r.ImageRef())

	bad := map[string]string{
		"autre image":   strings.Replace(fiche("0.2.0"), OfficialImage, "docker.io/pirate/forgenet", 1),
		"empreinte":     strings.Replace(fiche("0.2.0"), digest, "sha256:abc", 1),
		"schéma":        strings.Replace(fiche("0.2.0"), `"schema":1`, `"schema":2`, 1),
		"version":       strings.Replace(fiche("0.2.0"), `"version":"0.2.0"`, `"version":"latest"`, 1),
		"JSON invalide": "{",
	}
	for name, data := range bad {
		_, err := ParseRelease([]byte(data))
		assert.ErrorIs(t, err, ErrInvalidRelease, name)
	}

	r, err = ParseRelease([]byte(strings.Replace(fiche("0.2.0"), "https://github.com/forgelab-dev", "https://ailleurs.example", 1)))
	require.NoError(t, err)
	assert.Empty(t, r.NotesURL, "un lien de notes étranger est écarté")
}

func TestVersionCompare(t *testing.T) {
	order := []string{"0.1.0-alpha.dev", "0.1.0", "0.1.1", "0.2.0-rc.1", "0.2.0", "1.0.0"}
	for i := 1; i < len(order); i++ {
		a, err := ParseVersion(order[i-1])
		require.NoError(t, err)
		b, err := ParseVersion(order[i])
		require.NoError(t, err)
		assert.Equal(t, -1, a.Compare(b), "%s < %s", a, b)
		assert.Equal(t, 1, b.Compare(a))
	}
	v, err := ParseVersion("v0.1.1")
	require.NoError(t, err)
	assert.Equal(t, "0.1.1", v.String())
	for _, s := range []string{"1.2", "1.2.x", "01.2.3", "-1.2.3"} {
		_, err := ParseVersion(s)
		assert.Error(t, err, s)
	}
}

// fakeGitHub sert une liste de releases et leurs fiches.
func fakeGitHub(t *testing.T, releases map[string]string, drafts ...string) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "forgenet/test", r.Header.Get("User-Agent"))
		if strings.HasPrefix(r.URL.Path, "/fiche/") {
			_, _ = fmt.Fprint(w, releases[strings.TrimPrefix(r.URL.Path, "/fiche/")])
			return
		}
		var items []string
		for tag := range releases {
			draft := false
			for _, d := range drafts {
				draft = draft || d == tag
			}
			items = append(items, fmt.Sprintf(`{"tag_name":%q,"draft":%t,"assets":[{"name":%q,"browser_download_url":"%s/fiche/%s"}]}`,
				tag, draft, ReleaseFile, srv.URL, tag))
		}
		items = append(items, `{"tag_name":"v9.9.9","draft":false,"assets":[]}`) // sans fiche
		_, _ = fmt.Fprint(w, "["+strings.Join(items, ",")+"]")
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestSourceLatest(t *testing.T) {
	gh := fakeGitHub(t, map[string]string{
		"v0.1.1":      fiche("0.1.1"),
		"v0.2.0":      fiche("0.2.0"),
		"v0.3.0-rc.1": fiche("0.3.0-rc.1"),
		"v0.4.0":      fiche("0.4.0"),
	}, "v0.4.0")
	src := Source{API: gh.URL, UserAgent: "forgenet/test"}

	r, err := src.Latest(context.Background(), false)
	require.NoError(t, err)
	assert.Equal(t, "0.2.0", r.Version, "ni brouillon, ni pré-version, ni release sans fiche")

	r, err = src.Latest(context.Background(), true)
	require.NoError(t, err)
	assert.Equal(t, "0.3.0-rc.1", r.Version)

	want, _ := ParseVersion("0.1.1")
	r, err = src.ByVersion(context.Background(), want)
	require.NoError(t, err)
	assert.Equal(t, "0.1.1", r.Version)
}

func TestSourceRefusesMismatchedFiche(t *testing.T) {
	gh := fakeGitHub(t, map[string]string{"v0.2.0": fiche("0.1.0")})
	_, err := Source{API: gh.URL, UserAgent: "forgenet/test"}.Latest(context.Background(), false)
	assert.ErrorIs(t, err, ErrInvalidRelease, "une fiche qui annonce une autre version que sa release")
}

// fakeDocker joue le démon : pull, compose stop/up, inspect.
type fakeDocker struct {
	mu       sync.Mutex
	calls    []string
	pullErr  error
	upErr    error
	health   []string // états successifs renvoyés par inspect
	env      string   // .env lu au moment du « up », pour vérifier la bascule
	envPath  string
	upCount  int
	failUpOn int // numéro d'appel à « up » qui échoue (1 = le premier)
}

func (d *fakeDocker) Run(_ context.Context, args ...string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	joined := strings.Join(args, " ")
	d.calls = append(d.calls, joined)
	switch {
	case args[0] == "pull":
		return "", d.pullErr
	case strings.Contains(joined, " up "):
		d.upCount++
		data, _ := os.ReadFile(d.envPath)
		d.env = string(data)
		if d.failUpOn == d.upCount {
			return "", d.upErr
		}
	case args[0] == "inspect":
		if len(d.health) == 0 {
			return "running healthy", nil
		}
		h := d.health[0]
		if len(d.health) > 1 {
			d.health = d.health[1:]
		}
		return h, nil
	}
	return "", nil
}

type install struct {
	dir, data string
	docker    *fakeDocker
	engine    *Engine
	clock     time.Time
}

func newInstall(t *testing.T) *install {
	t.Helper()
	root := t.TempDir()
	in := &install{dir: filepath.Join(root, "opt"), data: filepath.Join(root, "data"), clock: time.Date(2026, 11, 14, 14, 0, 0, 0, time.UTC)}
	require.NoError(t, os.MkdirAll(in.dir, 0o755))
	require.NoError(t, os.MkdirAll(in.data, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(in.dir, ".env"), []byte("# commentaire gardé\nFORGENET_VERSION=0.1.1\nFORGENET_JWT_SECRET=secret\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(in.data, "forgenet.db"), []byte("base 0.1.1"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(in.data, "forgenet.db-wal"), []byte("journal 0.1.1"), 0o600))
	in.docker = &fakeDocker{envPath: filepath.Join(in.dir, ".env")}
	in.engine = &Engine{
		Docker: in.docker, InstallDir: in.dir, DataDir: in.data, HealthTimeout: 20 * time.Second,
		Now:   func() time.Time { return in.clock },
		Sleep: func(context.Context, time.Duration) error { in.clock = in.clock.Add(2 * time.Second); return nil },
	}
	return in
}

func release(t *testing.T, v string) Release {
	t.Helper()
	r, err := ParseRelease([]byte(fiche(v)))
	require.NoError(t, err)
	return r
}

func TestApplySucceeds(t *testing.T) {
	in := newInstall(t)
	in.docker.health = []string{"running starting", "running starting", "running healthy"}

	st, err := in.engine.Apply(context.Background(), release(t, "0.2.0"))
	require.NoError(t, err)
	assert.Equal(t, ResultSucceeded, st.Result)
	assert.Equal(t, "0.1.1", st.From)

	assert.Equal(t, "pull ghcr.io/forgelab-dev/forgenet-server:0.2.0@"+digest, in.docker.calls[0], "l'image est tirée avant tout arrêt")
	assert.Contains(t, in.docker.env, "FORGENET_VERSION=0.2.0", "le serveur redémarre sur la nouvelle version")
	assert.Contains(t, in.docker.env, "FORGENET_IMAGE_DIGEST="+digest)
	assert.Contains(t, in.docker.env, "# commentaire gardé")
	assert.Contains(t, in.docker.env, "FORGENET_JWT_SECRET=secret")

	saved, err := os.ReadFile(filepath.Join(st.Backup, "forgenet.db"))
	require.NoError(t, err)
	assert.Equal(t, "base 0.1.1", string(saved))
	info, err := os.Stat(filepath.Join(st.Backup, ".env"))
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm(), ".env sauvegardé reste privé")

	got, err := ReadState(in.data)
	require.NoError(t, err)
	assert.Equal(t, StepDone, got.Step)
	assert.NoFileExists(t, filepath.Join(in.data, "maj", "verrou"))
}

func TestApplyRollsBackWhenNewVersionIsUnhealthy(t *testing.T) {
	in := newInstall(t)
	// La nouvelle version reste « starting » jusqu'au délai ; l'ancienne repart.
	in.docker.health = []string{"running starting", "running starting", "running starting", "running starting",
		"running starting", "running starting", "running starting", "running starting", "running starting",
		"running starting", "running starting", "running healthy"}
	// La nouvelle version a touché la base avant de tomber.
	in.docker.failUpOn = 0

	st, err := in.engine.Apply(context.Background(), release(t, "0.2.0"))
	require.Error(t, err)
	assert.Equal(t, ResultRolledBack, st.Result)
	assert.Contains(t, st.Reason, "ne démarre pas")

	env, _ := os.ReadFile(filepath.Join(in.dir, ".env"))
	assert.Contains(t, string(env), "FORGENET_VERSION=0.1.1", ".env d'avant remis en place")
	assert.NotContains(t, string(env), "FORGENET_IMAGE_DIGEST")
	db, _ := os.ReadFile(filepath.Join(in.data, "forgenet.db"))
	assert.Equal(t, "base 0.1.1", string(db))
	assert.Contains(t, in.docker.env, "FORGENET_VERSION=0.1.1", "l'ancienne version est relancée")
}

func TestApplyRollsBackWhenStartFails(t *testing.T) {
	in := newInstall(t)
	in.docker.failUpOn, in.docker.upErr = 1, errors.New("image incompatible")
	st, err := in.engine.Apply(context.Background(), release(t, "0.2.0"))
	require.Error(t, err)
	assert.Equal(t, ResultRolledBack, st.Result)
	assert.Equal(t, 2, in.docker.upCount, "un « up » raté, puis celui du retour arrière")
}

func TestApplyFailsCleanlyBeforeStoppingAnything(t *testing.T) {
	in := newInstall(t)
	in.docker.pullErr = errors.New("manifest unknown")
	st, err := in.engine.Apply(context.Background(), release(t, "0.2.0"))
	require.Error(t, err)
	assert.Equal(t, ResultFailed, st.Result)
	assert.Len(t, in.docker.calls, 1, "rien d'autre qu'un pull raté : le serveur n'a pas été arrêté")
	assert.Empty(t, st.Backup)
}

func TestApplyRefusesOlderOrSameVersion(t *testing.T) {
	for _, v := range []string{"0.1.1", "0.1.0"} {
		in := newInstall(t)
		_, err := in.engine.Apply(context.Background(), release(t, v))
		assert.ErrorIs(t, err, ErrNotNewer, v)
		assert.Empty(t, in.docker.calls, v)
	}
	in := newInstall(t)
	in.engine.Force = true
	st, err := in.engine.Apply(context.Background(), release(t, "0.1.1"))
	require.NoError(t, err, "--force réinstalle la même version")
	assert.Equal(t, ResultSucceeded, st.Result)
}

func TestApplyRefusesConcurrentRuns(t *testing.T) {
	in := newInstall(t)
	require.NoError(t, os.MkdirAll(filepath.Join(in.data, "maj"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(in.data, "maj", "verrou"), nil, 0o600))
	_, err := in.engine.Apply(context.Background(), release(t, "0.2.0"))
	assert.ErrorIs(t, err, ErrBusy)
	assert.Empty(t, in.docker.calls)
}

func TestBackupsArePruned(t *testing.T) {
	in := newInstall(t)
	in.engine.Keep = 2
	in.engine.defaults()
	for i := 0; i < 4; i++ {
		in.clock = in.clock.Add(time.Hour)
		_, err := in.engine.backup("0.1.1")
		require.NoError(t, err)
	}
	entries, err := os.ReadDir(filepath.Join(in.dir, "sauvegardes"))
	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.True(t, strings.HasPrefix(entries[1].Name(), "20261114-180000"), "les plus récentes sont gardées")
}

// Régression trouvée sur le VPS : avec un délai court pour la nouvelle version,
// le retour arrière réussissait mais était déclaré en échec, faute d'avoir
// laissé à l'ancienne version le temps de redevenir saine.
func TestRollbackHasItsOwnHealthDelay(t *testing.T) {
	in := newInstall(t)
	in.engine.HealthTimeout = 3 * time.Second
	starting := make([]string, 0, 30)
	for i := 0; i < 25; i++ { // environ 50 s de démarrage pour chaque version
		starting = append(starting, "running starting")
	}
	in.docker.health = append(starting, "running healthy")
	st, err := in.engine.Apply(context.Background(), release(t, "0.2.0"))
	require.Error(t, err)
	assert.Equal(t, ResultRolledBack, st.Result, "l'ancienne version a eu le temps de redémarrer")
}

func withFiles(t *testing.T, in *install, compose string) Release {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(in.dir, "docker-compose.yml"), []byte("compose 0.1.1"), 0o644))
	r := release(t, "0.2.0")
	r.Files = map[string]string{"docker-compose.yml": fmt.Sprintf("%x", sha256.Sum256([]byte(compose)))}
	in.engine.Fetch = func(_ context.Context, version, name string) ([]byte, error) {
		assert.Equal(t, "0.2.0", version)
		assert.Equal(t, "docker-compose.yml", name)
		return []byte("compose 0.2.0"), nil
	}
	return r
}

func TestApplyReplacesInstallFiles(t *testing.T) {
	in := newInstall(t)
	st, err := in.engine.Apply(context.Background(), withFiles(t, in, "compose 0.2.0"))
	require.NoError(t, err)
	got, _ := os.ReadFile(filepath.Join(in.dir, "docker-compose.yml"))
	assert.Equal(t, "compose 0.2.0", string(got))
	saved, _ := os.ReadFile(filepath.Join(st.Backup, "docker-compose.yml"))
	assert.Equal(t, "compose 0.1.1", string(saved), "l'ancien fichier est dans la sauvegarde")
}

func TestApplyRefusesTamperedInstallFile(t *testing.T) {
	in := newInstall(t)
	st, err := in.engine.Apply(context.Background(), withFiles(t, in, "un autre contenu"))
	require.Error(t, err)
	assert.Equal(t, ResultFailed, st.Result)
	assert.Contains(t, err.Error(), "empreinte")
	assert.Len(t, in.docker.calls, 1, "refusé après le pull, avant tout arrêt")
	got, _ := os.ReadFile(filepath.Join(in.dir, "docker-compose.yml"))
	assert.Equal(t, "compose 0.1.1", string(got))
}

func TestRollbackRestoresInstallFiles(t *testing.T) {
	in := newInstall(t)
	in.docker.failUpOn, in.docker.upErr = 1, errors.New("service inconnu")
	st, err := in.engine.Apply(context.Background(), withFiles(t, in, "compose 0.2.0"))
	require.Error(t, err)
	assert.Equal(t, ResultRolledBack, st.Result)
	got, _ := os.ReadFile(filepath.Join(in.dir, "docker-compose.yml"))
	assert.Equal(t, "compose 0.1.1", string(got))
}

func TestParseReleaseFiles(t *testing.T) {
	sum := strings.Repeat("a", 64)
	ok := strings.Replace(fiche("0.2.0"), `"schema":1`, `"schema":1,"files":{"docker-compose.yml":"`+sum+`"}`, 1)
	r, err := ParseRelease([]byte(ok))
	require.NoError(t, err)
	assert.Equal(t, sum, r.Files["docker-compose.yml"])

	for name, data := range map[string]string{
		"nom inconnu": strings.Replace(ok, "docker-compose.yml", "../../etc/passwd", 1),
		"empreinte":   strings.Replace(ok, sum, "abc", 1),
	} {
		_, err := ParseRelease([]byte(data))
		assert.ErrorIs(t, err, ErrInvalidRelease, name)
	}
}

func TestApplyTakesOverAStaleLock(t *testing.T) {
	in := newInstall(t)
	lock := filepath.Join(in.data, "maj", "verrou")
	require.NoError(t, os.MkdirAll(filepath.Dir(lock), 0o750))
	require.NoError(t, os.WriteFile(lock, nil, 0o600))
	old := in.clock.Add(-time.Hour)
	require.NoError(t, os.Chtimes(lock, old, old))

	st, err := in.engine.Apply(context.Background(), release(t, "0.2.0"))
	require.NoError(t, err, "un verrou d'une heure vient d'une mise à jour interrompue")
	assert.Equal(t, ResultSucceeded, st.Result)
}
