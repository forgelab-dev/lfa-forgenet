package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/rs/zerolog"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/health"
)

// transport est le cycle de vie commun aux serveurs HTTP et gRPC.
type transport interface {
	Start() error
	Shutdown(ctx context.Context) error
}

type server struct {
	name      string
	transport transport
}

// serve démarre tous les serveurs, attend l'annulation de ctx ou la défaillance
// de l'un d'eux, puis les arrête tous dans le délai de grâce configuré. Il
// renvoie la défaillance éventuelle, jointe aux erreurs d'arrêt.
func serve(ctx context.Context, servers []server, healthService health.Service, cfg config.ServerConfig, logger zerolog.Logger) error {
	failures := make(chan error, len(servers))
	for _, s := range servers {
		go func() {
			if err := s.transport.Start(); err != nil {
				failures <- fmt.Errorf("serveur %s : %w", s.name, err)
			}
		}()
	}

	var failure error
	select {
	case <-ctx.Done():
		logger.Info().Msg("signal d'arrêt reçu, fermeture gracieuse")
	case failure = <-failures:
	}

	// Les sondes passent à NOT_SERVING avant l'arrêt, pour drainer le trafic.
	healthService.SetServing(false)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	shutdownErrs := make([]error, len(servers))
	var wg sync.WaitGroup
	for i, s := range servers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.transport.Shutdown(shutdownCtx); err != nil {
				shutdownErrs[i] = fmt.Errorf("arrêt du serveur %s : %w", s.name, err)
			}
		}()
	}
	wg.Wait()

	return errors.Join(append([]error{failure}, shutdownErrs...)...)
}
