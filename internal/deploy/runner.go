package deploy

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

// Start reprend les déploiements laissés par un arrêt, puis exécute la file
// jusqu'à Shutdown. Il bloque, comme le Start d'un serveur.
func (s *Service) Start() error {
	close(s.started)
	defer close(s.done)

	// loopCtx cesse à l'ordre d'arrêt : plus aucun déploiement n'est pris.
	// runCtx cesse au-delà du délai de grâce : celui en cours est interrompu.
	loopCtx, stopLoop := context.WithCancel(context.Background())
	defer stopLoop()
	runCtx, killRuns := context.WithCancel(context.Background())
	defer killRuns()
	go func() {
		select {
		case <-s.stop:
			stopLoop()
		case <-loopCtx.Done():
		}
	}()
	go func() {
		select {
		case <-s.kill:
			killRuns()
		case <-runCtx.Done():
		}
	}()

	select {
	case <-s.stop:
		return nil // Shutdown est arrivé avant le démarrage.
	default:
	}

	s.recoverInterrupted(loopCtx)

	ticker := time.NewTicker(s.poll)
	defer ticker.Stop()
	for {
		s.drain(loopCtx, runCtx)
		select {
		case <-loopCtx.Done():
			return nil
		case <-s.wake:
		case <-ticker.C:
		}
	}
}

// Shutdown arrête la prise de nouveaux déploiements et laisse celui en cours
// finir jusqu'à l'échéance de ctx ; au-delà, il est interrompu et noté en échec.
func (s *Service) Shutdown(ctx context.Context) error {
	s.stopOnce.Do(func() { close(s.stop) })
	select {
	case <-s.started:
	default:
		return nil // Start n'a jamais tourné : rien à attendre.
	}
	select {
	case <-s.done:
	case <-ctx.Done():
		s.killOnce.Do(func() { close(s.kill) })
		<-s.done
	}
	return nil
}

// recoverInterrupted note en échec les déploiements qu'un arrêt a surpris en
// cours d'exécution. Ceux en attente restent dans la file et seront exécutés.
func (s *Service) recoverInterrupted(ctx context.Context) {
	for _, state := range []store.DeploymentState{
		store.DeploymentBuilding, store.DeploymentDeploying,
		store.DeploymentHealthChecking, store.DeploymentRollingBack,
	} {
		for {
			list, err := s.store.ListDeploymentsByState(ctx, state, batchSize)
			if err != nil {
				s.logger.Error().Err(err).Str("state", string(state)).Msg("reprise des déploiements interrompus")
				break
			}
			for _, d := range list {
				s.transition(d.ID, state, store.DeploymentFailed, interruptedMessage)
			}
			if len(list) < batchSize {
				break
			}
		}
	}
}

// drain exécute les déploiements en attente, du plus ancien au plus récent,
// jusqu'à vider la file ou recevoir l'ordre d'arrêt.
func (s *Service) drain(loopCtx, runCtx context.Context) {
	for loopCtx.Err() == nil {
		pending, err := s.store.ListDeploymentsByState(loopCtx, store.DeploymentPending, batchSize)
		if err != nil {
			if loopCtx.Err() == nil {
				s.logger.Error().Err(err).Msg("lecture de la file des déploiements")
			}
			return
		}
		progressed := false
		for _, d := range pending {
			if loopCtx.Err() != nil {
				return
			}
			if s.execute(runCtx, d) {
				progressed = true
			}
		}
		// Rien n'a avancé (projets occupés par une suppression) : on attend le
		// prochain réveil au lieu de relire la file en boucle.
		if !progressed || len(pending) < batchSize {
			return
		}
	}
}

// execute fait passer d par toutes ses étapes. Il renvoie false si le projet
// est occupé (suppression en cours) : le déploiement reste en attente.
func (s *Service) execute(runCtx context.Context, d *store.Deployment) bool {
	if !s.acquire(d.ProjectID) {
		return false
	}
	defer s.release(d.ProjectID)

	ctx, cancel := context.WithTimeout(runCtx, s.timeout)
	defer cancel()
	run := &execution{cancel: cancel}
	s.mu.Lock()
	s.running[d.ID] = run
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		delete(s.running, d.ID)
		s.mu.Unlock()
	}()

	log := s.logger.With().Str("deployment_id", d.ID).Str("project_id", d.ProjectID).Logger()
	stack, content := StackName(d.ProjectID), []byte(d.SourceContent)
	state := store.DeploymentPending

	advance := func(to store.DeploymentState) bool {
		if !s.transition(d.ID, state, to, "") {
			return false
		}
		state = to
		return true
	}
	fail := func(err error) {
		to, message := store.DeploymentFailed, err.Error()
		s.mu.Lock()
		reason := run.reason
		s.mu.Unlock()
		switch {
		case reason != "" && state.CanTransitionTo(store.DeploymentCancelled):
			to, message = store.DeploymentCancelled, reason
		case runCtx.Err() != nil:
			message = interruptedMessage
		case errors.Is(ctx.Err(), context.DeadlineExceeded):
			message = fmt.Sprintf("délai de %s dépassé : %s", s.timeout, message)
		}
		s.transition(d.ID, state, to, message)
		log.Warn().Str("state", string(to)).Str("error", message).Msg("déploiement arrêté")
	}

	// Une source Compose n'a rien à construire : l'étape vérifie le fichier.
	if !advance(store.DeploymentBuilding) {
		return true // annulé entre-temps
	}
	if err := s.stacks.Validate(ctx, content); err != nil {
		fail(err)
		return true
	}
	if !advance(store.DeploymentDeploying) {
		return true
	}
	if err := s.stacks.Deploy(ctx, stack, content); err != nil {
		fail(err)
		return true
	}
	if !advance(store.DeploymentHealthChecking) {
		return true
	}
	services, err := s.stacks.Services(ctx, stack)
	if err == nil {
		err = checkHealthy(services)
	}
	if err != nil {
		fail(err)
		return true
	}
	if advance(store.DeploymentSuccess) {
		log.Info().Msg("déploiement réussi")
	}
	return true
}

// transition écrit un changement d'état, même si l'exécution vient d'être
// annulée : l'état final doit toujours être enregistré.
func (s *Service) transition(id string, from, to store.DeploymentState, message string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	if _, err := s.store.TransitionDeployment(ctx, id, from, to, message); err != nil {
		if !errors.Is(err, store.ErrConflict) {
			s.logger.Error().Err(err).Str("deployment_id", id).Str("from", string(from)).Str("to", string(to)).
				Msg("changement d'état du déploiement")
		}
		return false
	}
	return true
}
