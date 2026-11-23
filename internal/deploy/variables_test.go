package deploy_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/deploy"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
	"github.com/forgelab-dev/lfa-forgenet/internal/sealed"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

const secretValue = "s3cret-tres-long-a-ne-jamais-montrer"

func newBox(t *testing.T) *sealed.Box {
	t.Helper()
	key := make([]byte, sealed.KeySize)
	_, err := rand.Read(key)
	require.NoError(t, err)
	box, err := sealed.NewBox(key)
	require.NoError(t, err)
	return box
}

func install(project string) deploy.CreateRequest {
	req := request(project)
	req.Template = "postgres"
	req.NewProject = true
	req.Variables = []deploy.Variable{
		{Name: "POSTGRES_PASSWORD", Value: secretValue, Secret: true},
		{Name: "POSTGRES_USER", Value: "app"},
	}
	return req
}

var installEnv = []string{"POSTGRES_PASSWORD=" + secretValue, "POSTGRES_USER=app"}

// recordingAuditor garde les événements d'audit.
type recordingAuditor struct {
	mu     sync.Mutex
	events []map[string]any
}

func (a *recordingAuditor) Log(_ context.Context, action, resource string, meta map[string]any) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = append(a.events, map[string]any{"action": action, "resource": resource, "meta": meta})
	return nil
}

func (a *recordingAuditor) List(context.Context, int, int) ([]*store.AuditEvent, error) {
	return nil, nil
}

func TestInstall_SealsVariablesAndPassesThemToCompose(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	box := newBox(t)
	svc := newService(t, st, stacks, deploy.WithBox(box))
	ctx := context.Background()

	d, created, err := svc.CreateOnce(ctx, install("base"))
	require.NoError(t, err)
	assert.True(t, created)
	assert.Equal(t, "postgres", d.Template)

	stored, err := st.Deployments().GetDeployment(ctx, d.ID)
	require.NoError(t, err)
	require.NotNil(t, stored.VariablesSealed)
	assert.False(t, bytes.Contains(stored.VariablesSealed, []byte(secretValue)), "la base ne voit pas la valeur en clair")
	plain, err := box.Open("deployment.variables", stored.VariablesSealed)
	require.NoError(t, err)
	assert.Contains(t, string(plain), secretValue)

	run(t, svc)
	waitState(t, st, d.ID, store.DeploymentSuccess)
	stacks.mu.Lock()
	defer stacks.mu.Unlock()
	require.Len(t, stacks.deployVars, 1)
	assert.Equal(t, installEnv, stacks.deployVars[0])
	for _, vars := range stacks.validateVars {
		assert.Equal(t, installEnv, vars, "la vérification voit les mêmes valeurs")
	}
}

// Redéployer depuis la page Déploiements envoie un fichier sans valeurs : sans
// la reprise, ${POSTGRES_PASSWORD} deviendrait vide.
func TestCreate_ReusesPreviousVariables(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	svc := newService(t, st, stacks, deploy.WithBox(newBox(t)))
	ctx := context.Background()
	run(t, svc)

	first, err := svc.Create(ctx, install("base"))
	require.NoError(t, err)
	waitState(t, st, first.ID, store.DeploymentSuccess)

	second, err := svc.Create(ctx, request("base"))
	require.NoError(t, err)
	assert.Equal(t, "postgres", second.Template)
	waitState(t, st, second.ID, store.DeploymentSuccess)

	a, err := st.Deployments().GetDeployment(ctx, first.ID)
	require.NoError(t, err)
	b, err := st.Deployments().GetDeployment(ctx, second.ID)
	require.NoError(t, err)
	assert.Equal(t, a.VariablesSealed, b.VariablesSealed, "le chiffré est repris tel quel")

	stacks.mu.Lock()
	defer stacks.mu.Unlock()
	require.Len(t, stacks.deployVars, 2)
	assert.Equal(t, installEnv, stacks.deployVars[1])
}

func TestCreate_NewProjectRefusedWhenNameIsTaken(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	stacks.projects = []compose.Project{{Name: "forgenet-sur-hote", Containers: 1}}
	svc := newService(t, st, stacks, deploy.WithBox(newBox(t)))
	ctx := context.Background()

	_, err := svc.Create(ctx, request("deja"))
	require.NoError(t, err)
	_, err = svc.Create(ctx, install("deja"))
	require.ErrorIs(t, err, deploy.ErrExists, "un déploiement porte ce nom")

	_, err = svc.Create(ctx, install("sur-hote"))
	require.ErrorIs(t, err, deploy.ErrExists, "une stack porte ce nom sans déploiement enregistré")

	// Rejouer la même demande rend l'installation, sans la refuser.
	req := install("neuve")
	req.RequestID = "req-installation"
	d, created, err := svc.CreateOnce(ctx, req)
	require.NoError(t, err)
	require.True(t, created)
	again, created, err := svc.CreateOnce(ctx, req)
	require.NoError(t, err)
	assert.False(t, created, "la clé d'idempotence a rendu un déploiement existant")
	assert.Equal(t, d.ID, again.ID)
}

func TestCreate_VariablesNeedTheSecretKey(t *testing.T) {
	st := newStore(t)
	svc := newService(t, st, newFakeStacks())
	ctx := context.Background()

	_, err := svc.Create(ctx, install("base"))
	require.ErrorIs(t, err, deploy.ErrNoKey)

	// Une installation faite avec une clé, redéployée par un serveur sans clé.
	sealedBy := newService(t, st, newFakeStacks(), deploy.WithBox(newBox(t)))
	_, err = sealedBy.Create(ctx, install("base"))
	require.NoError(t, err)
	_, err = svc.Create(ctx, request("base"))
	require.ErrorIs(t, err, deploy.ErrNoKey, "les valeurs à reprendre ne s'ouvrent pas")
}

// Compose peut citer une valeur dans une erreur (« invalid port "…" ») : le
// journal et le message d'échec la masquent.
func TestRun_MasksSecretValues(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	stacks.output = []string{"service db : valeur " + secretValue + " refusée", "utilisateur app"}
	stacks.deployErr = fmt.Errorf("docker compose up : %s invalide", secretValue)
	svc := newService(t, st, stacks, deploy.WithBox(newBox(t)))
	ctx := context.Background()

	d, err := svc.Create(ctx, install("base"))
	require.NoError(t, err)
	run(t, svc)
	failed := waitState(t, st, d.ID, store.DeploymentFailed)
	assert.NotContains(t, failed.ErrorMessage, secretValue)
	assert.Contains(t, failed.ErrorMessage, "••••••")

	logs, err := st.Deployments().ListDeploymentLogs(ctx, d.ID, time.Time{}, 0)
	require.NoError(t, err)
	require.NotEmpty(t, logs)
	var all string
	for _, l := range logs {
		all += l.Message + "\n"
	}
	assert.NotContains(t, all, secretValue)
	assert.Contains(t, all, "service db : valeur •••••• refusée")
	assert.Contains(t, all, "utilisateur app", "une valeur qui n'est pas secrète reste lisible")
}

func TestRun_UnreadableVariablesFailTheDeployment(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	ctx := context.Background()

	// Chiffré par une autre clé : FORGENET_SECRET_KEY a changé entre-temps.
	other := newService(t, st, stacks, deploy.WithBox(newBox(t)))
	d, err := other.Create(ctx, install("base"))
	require.NoError(t, err)

	svc := newService(t, st, stacks, deploy.WithBox(newBox(t)))
	run(t, svc)
	failed := waitState(t, st, d.ID, store.DeploymentFailed)
	assert.Contains(t, failed.ErrorMessage, "FORGENET_SECRET_KEY")
	stacks.mu.Lock()
	defer stacks.mu.Unlock()
	assert.Empty(t, stacks.deployed, "rien n'est lancé sans les valeurs")
}

func TestCreate_AuditNamesVariablesNeverValues(t *testing.T) {
	st := newStore(t)
	stacks := newFakeStacks()
	// Le refus de Compose arrive enveloppé, comme celui du vrai moteur.
	stacks.validateErr = fmt.Errorf("%w : valeur %s refusée", compose.ErrInvalidFile, secretValue)
	auditor := &recordingAuditor{}
	svc := deploy.New(st.Deployments(), stacks, auditor, zerolog.Nop(), deploy.WithBox(newBox(t)))
	ctx := context.Background()

	_, err := svc.Create(ctx, install("base"))
	require.Error(t, err)
	assert.NotContains(t, err.Error(), secretValue, "l'erreur rendue au client non plus")

	stacks.set(func(f *fakeStacks) { f.validateErr = nil })
	_, err = svc.Create(ctx, install("base"))
	require.NoError(t, err)

	auditor.mu.Lock()
	defer auditor.mu.Unlock()
	require.Len(t, auditor.events, 2)
	for _, e := range auditor.events {
		assert.NotContains(t, fmt.Sprint(e), secretValue)
		meta := e["meta"].(map[string]any)
		assert.Equal(t, "postgres", meta["template"])
		assert.Equal(t, []string{"POSTGRES_PASSWORD", "POSTGRES_USER"}, meta["variables"])
	}
}
