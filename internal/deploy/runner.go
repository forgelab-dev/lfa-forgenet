package deploy

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
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
	// Le journal est fermé avant l'état final (voir journal) ; le defer couvre
	// un déploiement annulé entre deux étapes.
	journal := s.newJournal(d.ID)
	defer journal.close()

	advance := func(to store.DeploymentState) bool {
		if !s.transition(d.ID, state, to, "") {
			return false
		}
		state = to
		return true
	}
	fail := func(err error) {
		to, message := store.DeploymentFailed, failureMessage(err)
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
		if to == store.DeploymentCancelled {
			journal.result(store.LogStderr, "Déploiement "+message)
		} else {
			journal.result(store.LogStderr, "Échec : "+message)
		}
		journal.close()
		s.transition(d.ID, state, to, message)
		log.Warn().Str("state", string(to)).Str("error", message).Msg("déploiement arrêté")
	}

	// Une source Compose n'a rien à construire : l'étape vérifie le fichier.
	if !advance(store.DeploymentBuilding) {
		return true // annulé entre-temps
	}
	journal.info(stepValidate, "Vérification du fichier Compose")
	if err := s.stacks.Validate(ctx, content); err != nil {
		fail(err)
		return true
	}
	journal.info(stepValidate, "Fichier Compose valide")
	if !advance(store.DeploymentDeploying) {
		return true
	}
	journal.info(stepDeploy, "Mise en service de la stack "+stack+" (docker compose up)")
	if err := s.stacks.Deploy(ctx, stack, content, composeOutput(journal.output(stepDeploy))); err != nil {
		fail(err)
		return true
	}
	if !advance(store.DeploymentHealthChecking) {
		return true
	}
	journal.info(stepHealthCheck, "Vérification des services")
	services, err := s.stacks.Services(ctx, stack)
	if err == nil {
		for _, svc := range services {
			journal.info(stepHealthCheck, describeService(svc))
		}
		err = checkHealthy(services)
	}
	if err != nil {
		fail(err)
		return true
	}
	journal.result(store.LogStdout, "Déploiement réussi")
	journal.close()
	if advance(store.DeploymentSuccess) {
		log.Info().Msg("déploiement réussi")
		if s.afterDeploy != nil {
			s.afterDeploy(ctx, d.ProjectID)
		}
	}
	return true
}

// transition écrit un changement d'état, même si l'exécution vient d'être
// annulée : l'état final doit toujours être enregistré. Il est ensuite diffusé.
func (s *Service) transition(id string, from, to store.DeploymentState, message string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	d, err := s.store.TransitionDeployment(ctx, id, from, to, message)
	if err != nil {
		if !errors.Is(err, store.ErrConflict) {
			s.logger.Error().Err(err).Str("deployment_id", id).Str("from", string(from)).Str("to", string(to)).
				Msg("changement d'état du déploiement")
		}
		return false
	}
	s.events.publish(id, event{deployment: d})
	return true
}

// composeOutput rattache la sortie de docker compose au journal.
func composeOutput(add func(store.LogStream, string)) compose.OutputFunc {
	return func(stream compose.OutputStream, line string) {
		if stream == compose.Stderr {
			add(store.LogStderr, line)
			return
		}
		add(store.LogStdout, line)
	}
}

// describeService résume l'état d'un service pour le journal.
func describeService(svc compose.Service) string {
	line := fmt.Sprintf("%s : %s", svc.Service, svc.State)
	if svc.Health != "" {
		line += ", santé " + svc.Health
	}
	return line
}

// failureMessage décrit l'échec en une ligne : pour une commande Compose, sa
// dernière ligne d'erreur plutôt que toute sa progression, déjà au journal.
func failureMessage(err error) string {
	message := err.Error()
	var cmdErr *compose.CommandError
	if errors.As(err, &cmdErr) {
		message = strings.Replace(message, cmdErr.Error(), cmdErr.Summary(), 1)
	}
	return message
}
