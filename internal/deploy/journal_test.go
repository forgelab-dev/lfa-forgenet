package deploy

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
)

func TestCleanLine(t *testing.T) {
	cases := map[string]string{
		"\x1b[32m Container web Started\x1b[0m": " Container web Started",
		"tab\tici\x00\x07":                      "tab ici",
		"latin1 \xe9t\xe9":                      "latin1 �t�",
		"   ":                                   "",
		"fin   ":                                "fin",
	}
	for in, want := range cases {
		assert.Equal(t, want, cleanLine(in), "%q", in)
	}

	long := cleanLine(strings.Repeat("é", maxLogLineBytes))
	assert.True(t, strings.HasSuffix(long, "…"))
	assert.LessOrEqual(t, len(long), maxLogLineBytes+len("…"))
	assert.True(t, strings.HasPrefix(long, "éé"), "la coupe respecte les caractères multioctets")
}

func TestBroker_DropsSlowSubscriber(t *testing.T) {
	b := newBroker()
	slow, unsubscribeSlow := b.subscribe("dep_1")
	defer unsubscribeSlow()
	other, unsubscribeOther := b.subscribe("dep_2")
	defer unsubscribeOther()

	for range subscriberBuffer + 1 {
		b.publish("dep_1", event{})
	}
	n := 0
	for range slow.events {
		n++
	}
	assert.Equal(t, subscriberBuffer, n, "le tampon est livré, puis le canal est fermé")
	unsubscribeSlow() // déjà décroché : sans effet ni panique

	b.publish("dep_2", event{})
	assert.Len(t, other.events, 1, "les autres déploiements ne sont pas touchés")
}

func TestJournal_IsBoundedButKeepsTheResult(t *testing.T) {
	st, err := sqlite.Open(context.Background(), config.DatabaseConfig{Path: ":memory:", BusyTimeout: 2 * time.Second, MaxOpenConns: 1, MaxIdleConns: 1})
	require.NoError(t, err)
	t.Cleanup(func() { _ = st.Close() })
	ctx := context.Background()
	require.NoError(t, st.Deployments().CreateDeployment(ctx, &store.Deployment{
		ID: "dep_1", ProjectID: "p1", SourceType: store.SourceCompose, SourceContent: "x", TriggeredBy: "user:u1",
	}))
	svc := New(st.Deployments(), nil, nil, zerolog.Nop())

	j := svc.newJournal("dep_1")
	out := j.output(stepDeploy)
	for range maxLogLines + 50 {
		out(store.LogStderr, "ligne")
	}
	j.result(store.LogStdout, "Déploiement réussi")
	j.close()
	j.info(stepDeploy, "après la fermeture") // ignorée

	logs, err := st.Deployments().ListDeploymentLogs(ctx, "dep_1", time.Time{}, 0)
	require.NoError(t, err)
	total := len(logs)
	for len(logs) == 1000 {
		logs, err = st.Deployments().ListDeploymentLogs(ctx, "dep_1", logs[len(logs)-1].EmittedAt, 0)
		require.NoError(t, err)
		total += len(logs)
	}
	assert.Equal(t, maxLogLines+2, total, "la limite, l'avis de troncature et le bilan")
	assert.Equal(t, "Déploiement réussi", logs[len(logs)-1].Message, "le bilan passe malgré la limite")
	assert.Equal(t, truncatedNotice, logs[len(logs)-2].Message)
}
