package deploy

import (
	"sync"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

// subscriberBuffer est le retard qu'un abonné peut prendre sur un déploiement
// avant d'être décroché : de quoi absorber une rafale de journal.
const subscriberBuffer = 512

// event est un changement d'état ou une ligne de journal d'un déploiement.
type event struct {
	deployment *store.Deployment
	log        *store.DeploymentLog
}

// subscription reçoit les événements d'un déploiement. Son canal est fermé
// quand l'abonné est décroché pour retard : il doit relire et se réabonner.
type subscription struct {
	events chan event
}

// broker diffuse en mémoire les événements de chaque déploiement à ses
// abonnés. Publier ne bloque jamais : un abonné trop lent est décroché, pour
// qu'un navigateur sur une liaison lente ne ralentisse pas un déploiement.
type broker struct {
	mu   sync.Mutex
	subs map[string]map[*subscription]struct{}
}

func newBroker() *broker {
	return &broker{subs: make(map[string]map[*subscription]struct{})}
}

// subscribe abonne aux événements du déploiement id ; unsubscribe libère
// l'abonnement et doit toujours être appelée.
func (b *broker) subscribe(id string) (sub *subscription, unsubscribe func()) {
	sub = &subscription{events: make(chan event, subscriberBuffer)}
	b.mu.Lock()
	if b.subs[id] == nil {
		b.subs[id] = make(map[*subscription]struct{})
	}
	b.subs[id][sub] = struct{}{}
	b.mu.Unlock()

	return sub, func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if _, ok := b.subs[id][sub]; ok {
			b.remove(id, sub)
		}
	}
}

func (b *broker) publish(id string, ev event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for sub := range b.subs[id] {
		select {
		case sub.events <- ev:
		default:
			b.remove(id, sub)
		}
	}
}

// remove retire sub et ferme son canal ; b.mu doit être tenu.
func (b *broker) remove(id string, sub *subscription) {
	delete(b.subs[id], sub)
	if len(b.subs[id]) == 0 {
		delete(b.subs, id)
	}
	close(sub.events)
}
