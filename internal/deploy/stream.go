package deploy

import (
	"context"
	"errors"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

// ErrLagging indique un abonné décroché parce qu'il ne suivait plus le
// rythme du déploiement : il peut se reconnecter, sans rien perdre, en
// reprenant à sa dernière ligne.
var ErrLagging = errors.New("flux en retard sur le déploiement : reconnectez-vous")

// logPage est la taille des pages de journal relues en base.
const logPage = 500

// Watch envoie l'état du déploiement id, puis chacun de ses changements,
// jusqu'à un état final ou la fin de ctx.
func (s *Service) Watch(ctx context.Context, id string, send func(*store.Deployment) error) error {
	// L'abonnement précède la lecture : un changement survenu entre les deux
	// arrive par le flux au lieu d'être perdu.
	sub, unsubscribe := s.events.subscribe(id)
	defer unsubscribe()

	d, err := s.store.GetDeployment(ctx, id)
	if err != nil {
		return err
	}
	if err := send(d); err != nil {
		return err
	}
	last := d
	for !last.State.IsTerminal() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev, ok := <-sub.events:
			if !ok {
				return ErrLagging
			}
			// Un changement antérieur à la lecture est déjà connu.
			if ev.deployment == nil || !ev.deployment.UpdatedAt.After(last.UpdatedAt) {
				continue
			}
			last = ev.deployment
			if err := send(last); err != nil {
				return err
			}
		}
	}
	return nil
}

// LogOptions règle la lecture d'un journal de déploiement.
type LogOptions struct {
	// Since reprend après cette ligne (exclue) ; zéro lit depuis le début.
	Since time.Time
	// Follow continue en direct jusqu'à la fin du déploiement.
	Follow bool
	// Stream ne garde qu'un flux ; vide garde les deux.
	Stream store.LogStream
}

// Logs envoie le journal du déploiement id : les lignes enregistrées, puis,
// avec Follow, les nouvelles en direct jusqu'à l'état final du déploiement.
func (s *Service) Logs(ctx context.Context, id string, opts LogOptions, send func(*store.DeploymentLog) error) error {
	var sub *subscription
	if opts.Follow {
		var unsubscribe func()
		sub, unsubscribe = s.events.subscribe(id)
		defer unsubscribe()
	}

	d, err := s.store.GetDeployment(ctx, id)
	if err != nil {
		return err
	}

	last := opts.Since
	deliver := func(l *store.DeploymentLog) error {
		if !l.EmittedAt.After(last) {
			return nil // déjà envoyée
		}
		last = l.EmittedAt
		if opts.Stream != "" && l.Stream != opts.Stream {
			return nil
		}
		return send(l)
	}
	replay := func() error {
		for {
			page, err := s.store.ListDeploymentLogs(ctx, id, last, logPage)
			if err != nil {
				return err
			}
			for _, l := range page {
				if err := deliver(l); err != nil {
					return err
				}
			}
			if len(page) < logPage {
				return nil
			}
		}
	}

	if err := replay(); err != nil {
		return err
	}
	if !opts.Follow || d.State.IsTerminal() {
		return nil
	}

	// Les lignes ne sont diffusées qu'une fois en base : celles écrites avant
	// l'abonnement ont été relues, les suivantes arrivent par le flux.
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev, ok := <-sub.events:
			if !ok {
				return ErrLagging
			}
			if ev.log != nil {
				if err := deliver(ev.log); err != nil {
					return err
				}
				continue
			}
			// Le journal est complet en base avant l'état final : une dernière
			// relecture rattrape ce qui aurait précédé l'abonnement.
			if ev.deployment != nil && ev.deployment.State.IsTerminal() {
				return replay()
			}
		}
	}
}
