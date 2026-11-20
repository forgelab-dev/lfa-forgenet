package mail

import (
	"context"
	"net/mail"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// DefaultRetryDelays espacent les trois essais d'un envoi en tâche de fond
// (ADR-0014) : tout de suite, puis 30 secondes et 2 minutes plus tard.
var DefaultRetryDelays = []time.Duration{0, 30 * time.Second, 2 * time.Minute}

// outboxCapacity borne les envois en attente : au-delà, une demande est
// abandonnée (et notée) plutôt que de faire grossir la mémoire.
const outboxCapacity = 100

// Job est un envoi en tâche de fond. Build prépare le message au moment de
// l'envoi : il peut chercher le compte, créer un lien, ou renvoyer nil s'il
// n'y a rien à envoyer. Il ne tourne qu'une fois ; seul l'envoi est réessayé.
type Job struct {
	// Kind nomme le message pour le journal (ex. « password_reset »).
	Kind  string
	Build func(ctx context.Context) (msg *Message, userID string, err error)
}

// Outbox envoie des messages en tâche de fond, avec trois essais espacés. Le
// journal note le type de message, le compte et l'issue ; jamais l'adresse,
// le lien ni le contenu (ADR-0014).
type Outbox struct {
	sender Sender
	delays []time.Duration
	logger zerolog.Logger

	jobs     chan Job
	stop     chan struct{}
	stopOnce sync.Once
	started  chan struct{}
	done     chan struct{}
}

// NewOutbox prépare la file ; delays vide prend DefaultRetryDelays.
func NewOutbox(sender Sender, delays []time.Duration, logger zerolog.Logger) *Outbox {
	if len(delays) == 0 {
		delays = DefaultRetryDelays
	}
	return &Outbox{
		sender: sender, delays: delays,
		logger:  logger.With().Str("component", "mail.outbox").Logger(),
		jobs:    make(chan Job, outboxCapacity),
		stop:    make(chan struct{}),
		started: make(chan struct{}),
		done:    make(chan struct{}),
	}
}

// Enabled dit si les messages partent : la file suit l'état du relais, qui
// peut changer sans redémarrage (ADR-0015). Un Sender qui ne le dit pas est
// considéré comme réglé.
func (o *Outbox) Enabled() bool {
	if e, ok := o.sender.(interface{ Enabled() bool }); ok {
		return e.Enabled()
	}
	return true
}

// Enqueue ajoute un envoi sans attendre. Il renvoie false si la file est
// pleine ou arrêtée.
func (o *Outbox) Enqueue(job Job) bool {
	select {
	case <-o.stop:
		return false
	default:
	}
	select {
	case o.jobs <- job:
		return true
	default:
		o.logger.Warn().Str("kind", job.Kind).Msg("file d'envoi pleine : message abandonné")
		return false
	}
}

// Start traite la file jusqu'à Shutdown. Chaque envoi a sa goroutine : un
// relais injoignable, qui fait attendre les essais suivants, ne retient pas
// les autres messages.
func (o *Outbox) Start() error {
	close(o.started)
	defer close(o.done)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	defer wg.Wait()
	for {
		select {
		case <-o.stop:
			cancel()
			return nil
		case job := <-o.jobs:
			wg.Add(1)
			go func() {
				defer wg.Done()
				o.run(ctx, job)
			}()
		}
	}
}

// Shutdown arrête la file : les envois en cours s'interrompent, ceux qui
// attendaient un nouvel essai sont abandonnés.
func (o *Outbox) Shutdown(ctx context.Context) error {
	o.stopOnce.Do(func() { close(o.stop) })
	select {
	case <-o.started:
	default:
		return nil
	}
	select {
	case <-o.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (o *Outbox) run(ctx context.Context, job Job) {
	buildCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	msg, userID, err := job.Build(buildCtx)
	cancel()
	log := o.logger.With().Str("kind", job.Kind).Logger()
	if userID != "" {
		log = log.With().Str("user_id", userID).Logger()
	}
	if err != nil {
		log.Error().Err(err).Msg("message non préparé")
		return
	}
	if msg == nil {
		return // rien à envoyer (aucun compte, demande freinée…)
	}

	for attempt, delay := range o.delays {
		if delay > 0 {
			select {
			case <-ctx.Done():
				log.Warn().Int("attempt", attempt+1).Msg("arrêt du serveur : message abandonné")
				return
			case <-time.After(delay):
			}
		}
		sendCtx, cancel := context.WithTimeout(ctx, DefaultTimeout)
		err = o.sender.Send(sendCtx, *msg)
		cancel()
		if err == nil {
			log.Info().Int("attempt", attempt+1).Msg("message remis au relais")
			return
		}
		// La réponse du relais aide à corriger ; elle peut citer l'adresse,
		// retirée avant d'entrer au journal.
		log.Warn().Str("error", redact(err.Error(), msg.To)).Int("attempt", attempt+1).Msg("message non remis")
	}
	log.Error().Int("attempts", len(o.delays)).Msg("message abandonné après tous les essais")
}

// redact retire de text l'adresse du destinataire, sous toutes ses formes.
func redact(text, to string) string {
	forms := []string{to}
	if a, err := mail.ParseAddress(to); err == nil {
		forms = append(forms, a.Address)
	}
	for _, f := range forms {
		if f = strings.TrimSpace(f); f != "" {
			text = strings.ReplaceAll(text, f, "<destinataire>")
		}
	}
	return text
}
