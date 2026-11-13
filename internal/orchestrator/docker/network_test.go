package docker_test

import (
	"context"
	"errors"
	"testing"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/docker"
)

// netClient simule le démon pour les opérations réseau.
type netClient struct {
	docker.EngineClient // méthodes non utilisées : panique si appelées

	network        *network.Inspect // nil : réseau absent
	createConflict bool
	created        []client.NetworkCreateOptions
	attached       map[string]*network.EndpointSettings // réseaux du conteneur
	inspectErr     error
	connects       []client.NetworkConnectOptions
	disconnects    []string
}

func (c *netClient) NetworkInspect(context.Context, string, client.NetworkInspectOptions) (client.NetworkInspectResult, error) {
	if c.network == nil {
		return client.NetworkInspectResult{}, cerrdefs.ErrNotFound
	}
	return client.NetworkInspectResult{Network: *c.network}, nil
}

func (c *netClient) NetworkCreate(_ context.Context, name string, opts client.NetworkCreateOptions) (client.NetworkCreateResult, error) {
	c.created = append(c.created, opts)
	inspect := &network.Inspect{}
	inspect.Name = name
	inspect.Labels = opts.Labels
	if c.createConflict {
		// Un autre appel l'a créé juste avant.
		c.network = inspect
		return client.NetworkCreateResult{}, cerrdefs.ErrConflict
	}
	c.network = inspect
	return client.NetworkCreateResult{}, nil
}

func (c *netClient) ContainerInspect(context.Context, string, client.ContainerInspectOptions) (client.ContainerInspectResult, error) {
	if c.inspectErr != nil {
		return client.ContainerInspectResult{}, c.inspectErr
	}
	return client.ContainerInspectResult{Container: container.InspectResponse{
		NetworkSettings: &container.NetworkSettings{Networks: c.attached},
	}}, nil
}

func (c *netClient) NetworkConnect(_ context.Context, name string, opts client.NetworkConnectOptions) (client.NetworkConnectResult, error) {
	c.connects = append(c.connects, opts)
	if c.attached == nil {
		c.attached = map[string]*network.EndpointSettings{}
	}
	c.attached[name] = opts.EndpointConfig
	return client.NetworkConnectResult{}, nil
}

func (c *netClient) NetworkDisconnect(_ context.Context, name string, opts client.NetworkDisconnectOptions) (client.NetworkDisconnectResult, error) {
	c.disconnects = append(c.disconnects, opts.Container)
	delete(c.attached, name)
	return client.NetworkDisconnectResult{}, nil
}

func newNetEngine(t *testing.T, c *netClient) *docker.Engine {
	t.Helper()
	engine, err := docker.NewWithClient(c, "127.0.0.1")
	require.NoError(t, err)
	return engine
}

var edgeLabels = map[string]string{"forgenet.managed": "edge"}

func TestEnsureNetwork_CreatesMissingBridge(t *testing.T) {
	c := &netClient{}
	require.NoError(t, newNetEngine(t, c).EnsureNetwork(context.Background(), "forgenet-edge", edgeLabels))
	require.Len(t, c.created, 1)
	assert.Equal(t, "bridge", c.created[0].Driver)
	assert.Equal(t, edgeLabels, c.created[0].Labels)
}

func TestEnsureNetwork_KeepsOwnNetwork(t *testing.T) {
	own := &network.Inspect{}
	own.Labels = map[string]string{"forgenet.managed": "edge", "autre": "x"}
	c := &netClient{network: own}
	require.NoError(t, newNetEngine(t, c).EnsureNetwork(context.Background(), "forgenet-edge", edgeLabels))
	assert.Empty(t, c.created)
}

func TestEnsureNetwork_RefusesForeignNetwork(t *testing.T) {
	foreign := &network.Inspect{}
	foreign.Labels = map[string]string{"com.docker.compose.project": "autre"}
	c := &netClient{network: foreign}
	err := newNetEngine(t, c).EnsureNetwork(context.Background(), "forgenet-edge", edgeLabels)
	assert.ErrorIs(t, err, docker.ErrForeignNetwork)
	assert.Empty(t, c.created)
}

func TestEnsureNetwork_ConcurrentCreation(t *testing.T) {
	c := &netClient{createConflict: true}
	require.NoError(t, newNetEngine(t, c).EnsureNetwork(context.Background(), "forgenet-edge", edgeLabels))
}

func TestConnectToNetwork(t *testing.T) {
	aliases := []string{"web.boutique.forgenet.internal"}

	t.Run("branche un conteneur absent du réseau", func(t *testing.T) {
		c := &netClient{}
		require.NoError(t, newNetEngine(t, c).ConnectToNetwork(context.Background(), "forgenet-edge", "c1", aliases))
		require.Len(t, c.connects, 1)
		assert.Equal(t, "c1", c.connects[0].Container)
		assert.Equal(t, aliases, c.connects[0].EndpointConfig.Aliases)
	})

	t.Run("ne fait rien s'il est déjà branché avec l'alias", func(t *testing.T) {
		c := &netClient{attached: map[string]*network.EndpointSettings{
			"forgenet-edge": {Aliases: []string{"c1", "web.boutique.forgenet.internal"}},
		}}
		require.NoError(t, newNetEngine(t, c).ConnectToNetwork(context.Background(), "forgenet-edge", "c1", aliases))
		assert.Empty(t, c.connects)
		assert.Empty(t, c.disconnects)
	})

	t.Run("rebranche s'il manque l'alias", func(t *testing.T) {
		c := &netClient{attached: map[string]*network.EndpointSettings{"forgenet-edge": {Aliases: []string{"ancien"}}}}
		require.NoError(t, newNetEngine(t, c).ConnectToNetwork(context.Background(), "forgenet-edge", "c1", aliases))
		assert.Equal(t, []string{"c1"}, c.disconnects)
		require.Len(t, c.connects, 1)
		assert.Equal(t, aliases, c.connects[0].EndpointConfig.Aliases)
	})

	t.Run("conteneur introuvable", func(t *testing.T) {
		c := &netClient{inspectErr: cerrdefs.ErrNotFound}
		err := newNetEngine(t, c).ConnectToNetwork(context.Background(), "forgenet-edge", "c1", aliases)
		assert.ErrorIs(t, err, orchestrator.ErrNotFound)
	})
}

func TestDisconnectFromNetwork(t *testing.T) {
	t.Run("débranche un conteneur branché", func(t *testing.T) {
		c := &netClient{attached: map[string]*network.EndpointSettings{"forgenet-edge": {}}}
		require.NoError(t, newNetEngine(t, c).DisconnectFromNetwork(context.Background(), "forgenet-edge", "c1"))
		assert.Equal(t, []string{"c1"}, c.disconnects)
	})

	t.Run("rien à faire s'il n'est pas branché", func(t *testing.T) {
		c := &netClient{attached: map[string]*network.EndpointSettings{"bridge": {}}}
		require.NoError(t, newNetEngine(t, c).DisconnectFromNetwork(context.Background(), "forgenet-edge", "c1"))
		assert.Empty(t, c.disconnects)
	})

	t.Run("rien à faire si le conteneur a disparu", func(t *testing.T) {
		c := &netClient{inspectErr: cerrdefs.ErrNotFound}
		require.NoError(t, newNetEngine(t, c).DisconnectFromNetwork(context.Background(), "forgenet-edge", "c1"))
	})

	t.Run("une autre erreur remonte", func(t *testing.T) {
		c := &netClient{inspectErr: errors.New("démon injoignable")}
		assert.Error(t, newNetEngine(t, c).DisconnectFromNetwork(context.Background(), "forgenet-edge", "c1"))
	})
}
