package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

var allDeploymentStates = []store.DeploymentState{
	store.DeploymentPending, store.DeploymentBuilding, store.DeploymentDeploying,
	store.DeploymentHealthChecking, store.DeploymentSuccess, store.DeploymentFailed,
	store.DeploymentCancelled, store.DeploymentRollingBack, store.DeploymentRolledBack,
}

// La machine à états doit accepter exactement les transitions écrites dans
// deploy.proto, et aucune autre.
func TestDeploymentState_TransitionsMatchContract(t *testing.T) {
	allowed := map[store.DeploymentState][]store.DeploymentState{
		store.DeploymentPending:        {store.DeploymentBuilding, store.DeploymentCancelled},
		store.DeploymentBuilding:       {store.DeploymentDeploying, store.DeploymentFailed, store.DeploymentCancelled},
		store.DeploymentDeploying:      {store.DeploymentHealthChecking, store.DeploymentFailed, store.DeploymentCancelled},
		store.DeploymentHealthChecking: {store.DeploymentSuccess, store.DeploymentFailed, store.DeploymentRollingBack},
		store.DeploymentRollingBack:    {store.DeploymentRolledBack, store.DeploymentFailed},
	}

	for _, from := range allDeploymentStates {
		for _, to := range allDeploymentStates {
			want := false
			for _, a := range allowed[from] {
				if a == to {
					want = true
				}
			}
			assert.Equal(t, want, from.CanTransitionTo(to), "%s vers %s", from, to)
		}
	}
}

func TestDeploymentState_Terminal(t *testing.T) {
	terminal := map[store.DeploymentState]bool{
		store.DeploymentSuccess: true, store.DeploymentFailed: true,
		store.DeploymentCancelled: true, store.DeploymentRolledBack: true,
	}
	for _, s := range allDeploymentStates {
		assert.True(t, s.IsValid(), s)
		assert.Equal(t, terminal[s], s.IsTerminal(), s)
	}

	unknown := store.DeploymentState("inconnu")
	assert.False(t, unknown.IsValid())
	assert.False(t, unknown.IsTerminal(), "un état inconnu n'est pas final")
	assert.False(t, unknown.CanTransitionTo(store.DeploymentPending))
}

func TestDeploymentSourceType_IsValid(t *testing.T) {
	for _, s := range []store.DeploymentSourceType{store.SourceCompose, store.SourceGit, store.SourceImage} {
		assert.True(t, s.IsValid(), s)
	}
	assert.False(t, store.DeploymentSourceType("helm").IsValid())
	assert.False(t, store.DeploymentSourceType("").IsValid())
}
