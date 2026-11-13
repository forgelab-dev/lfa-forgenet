package docker

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"

	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
)

// ErrForeignNetwork indique un réseau du même nom que ForgeNet n'a pas créé :
// ForgeNet n'y branche aucune application (ADR-0008).
var ErrForeignNetwork = errors.New("réseau existant non créé par ForgeNet")

// EnsureNetwork crée le réseau pont name, avec ses étiquettes, s'il n'existe
// pas. Un réseau existant doit porter toutes ces étiquettes.
func (e *Engine) EnsureNetwork(ctx context.Context, name string, labels map[string]string) error {
	res, err := e.cli.NetworkInspect(ctx, name, client.NetworkInspectOptions{})
	switch {
	case err == nil:
		for k, v := range labels {
			if res.Network.Labels[k] != v {
				return fmt.Errorf("%w : %s (étiquette %s=%s absente)", ErrForeignNetwork, name, k, v)
			}
		}
		return nil
	case !cerrdefs.IsNotFound(err):
		return engineError("inspection du réseau "+name, err)
	}

	_, err = e.cli.NetworkCreate(ctx, name, client.NetworkCreateOptions{Driver: "bridge", Labels: maps.Clone(labels)})
	if err != nil && cerrdefs.IsConflict(err) {
		// Créé entre-temps par un autre appel : on vérifie qu'il est bien le nôtre.
		return e.EnsureNetwork(ctx, name, labels)
	}
	if err != nil {
		return engineError("création du réseau "+name, err)
	}
	return nil
}

// ConnectToNetwork branche le conteneur sur le réseau avec ces alias DNS. Déjà
// branché avec eux, il n'y a rien à faire ; branché sans eux, il est rebranché.
func (e *Engine) ConnectToNetwork(ctx context.Context, networkName, containerID string, aliases []string) error {
	res, err := e.inspect(ctx, containerID)
	if err != nil {
		return err
	}
	if settings := res.Container.NetworkSettings; settings != nil {
		if ep, ok := settings.Networks[networkName]; ok && ep != nil {
			if hasAll(ep.Aliases, aliases) {
				return nil
			}
			if err := e.DisconnectFromNetwork(ctx, networkName, containerID); err != nil {
				return err
			}
		}
	}

	_, err = e.cli.NetworkConnect(ctx, networkName, client.NetworkConnectOptions{
		Container:      containerID,
		EndpointConfig: &network.EndpointSettings{Aliases: slices.Clone(aliases)},
	})
	if err != nil {
		return engineError(fmt.Sprintf("branchement du conteneur %s sur %s", containerID, networkName), err)
	}
	return nil
}

// DisconnectFromNetwork débranche le conteneur du réseau. Un conteneur absent
// ou déjà débranché n'est pas une erreur : l'état est lu avant d'agir, plutôt
// que déduit du message d'erreur du démon.
func (e *Engine) DisconnectFromNetwork(ctx context.Context, networkName, containerID string) error {
	res, err := e.inspect(ctx, containerID)
	if errors.Is(err, orchestrator.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if settings := res.Container.NetworkSettings; settings == nil || settings.Networks[networkName] == nil {
		return nil
	}
	_, err = e.cli.NetworkDisconnect(ctx, networkName, client.NetworkDisconnectOptions{Container: containerID})
	if err != nil {
		return engineError(fmt.Sprintf("débranchement du conteneur %s de %s", containerID, networkName), err)
	}
	return nil
}

func hasAll(have, want []string) bool {
	for _, w := range want {
		if !slices.Contains(have, w) {
			return false
		}
	}
	return true
}
