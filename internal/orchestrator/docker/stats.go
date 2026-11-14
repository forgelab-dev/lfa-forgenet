package docker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"

	"github.com/forgelab-dev/lfa-forgenet/internal/metrics"
)

var (
	_ metrics.StatsSource    = (*Engine)(nil)
	_ metrics.RunningCounter = (*Engine)(nil)
	_ metrics.RunningLister  = (*Engine)(nil)
)

// Stats lit une mesure des ressources d'un conteneur (identifiant ou nom).
// withPrevious fait aussi remplir PreCPUStats par le démon, qui attend alors
// une seconde entre les deux relevés.
func (e *Engine) Stats(ctx context.Context, id string, withPrevious bool) (container.StatsResponse, error) {
	id, err := requireID(id)
	if err != nil {
		return container.StatsResponse{}, err
	}
	res, err := e.cli.ContainerStats(ctx, id, client.ContainerStatsOptions{IncludePreviousSample: withPrevious})
	if err != nil {
		return container.StatsResponse{}, engineError("statistiques du conteneur "+id, err)
	}
	defer func() { _ = res.Body.Close() }()

	var stats container.StatsResponse
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		return container.StatsResponse{}, fmt.Errorf("statistiques du conteneur %s illisibles : %w", id, err)
	}
	return stats, nil
}

// CountRunning compte les conteneurs en cours d'exécution.
func (e *Engine) CountRunning(ctx context.Context) (int, error) {
	res, err := e.cli.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		return 0, fmt.Errorf("liste des conteneurs : %w", err)
	}
	return len(res.Items), nil
}

// RunningIDs renvoie l'identifiant complet de chaque conteneur en cours
// d'exécution.
func (e *Engine) RunningIDs(ctx context.Context) ([]string, error) {
	res, err := e.cli.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("liste des conteneurs : %w", err)
	}
	ids := make([]string, 0, len(res.Items))
	for _, c := range res.Items {
		ids = append(ids, c.ID)
	}
	return ids, nil
}
