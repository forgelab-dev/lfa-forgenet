package update

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
)

// launchDocker répond aux inspections de Launch et garde la commande « run ».
type launchDocker struct{ run []string }

func (d *launchDocker) Run(_ context.Context, args ...string) (string, error) {
	switch {
	case args[0] == "inspect" && args[len(args)-1] == UpdaterName:
		return "", errors.New("no such container")
	case args[0] == "inspect" && strings.Contains(args[2], ".Config.Image"):
		return "ghcr.io/forgelab-dev/forgenet-server:0.6.0\n", nil
	case args[0] == "inspect":
		return "forgenet_data\n", nil
	case args[0] == "run":
		d.run = args
		return "0123456789abcdef\n", nil
	}
	return "", errors.New("commande inattendue : " + strings.Join(args, " "))
}

func TestLaunch_MarksTheUpdaterAsSystem(t *testing.T) {
	d := &launchDocker{}
	id, err := Launch(context.Background(), d, LaunchOptions{InstallDir: "/opt/forgenet", To: "0.6.1"})
	require.NoError(t, err)
	assert.Equal(t, "0123456789ab", id)

	// Le conteneur reste après sa fin : l'étiquette le range parmi ceux de
	// ForgeNet, pas parmi les applications du tableau de bord.
	assert.Contains(t, strings.Join(d.run, " "), "--label "+orchestrator.SystemLabel+"=true")
	assert.Equal(t, "ghcr.io/forgelab-dev/forgenet-server:0.6.0", d.run[indexOf(d.run, "update")-1], "image en service")
}

func indexOf(args []string, s string) int {
	for i, a := range args {
		if a == s {
			return i
		}
	}
	return -1
}
