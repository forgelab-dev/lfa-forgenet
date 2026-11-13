package deploy

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

const (
	// maxLogLines et maxLogBytes bornent le journal d'un déploiement : un
	// fichier bavard ne doit pas remplir le disque du serveur.
	maxLogLines = 5000
	maxLogBytes = 1 << 20
	// maxLogLineBytes borne une ligne.
	maxLogLineBytes = 4 << 10
	// journalFlushInterval espace les écritures en base : les lignes partent
	// par lots, et sont diffusées une fois écrites.
	journalFlushInterval = 200 * time.Millisecond

	stepValidate    = "validate"
	stepDeploy      = "deploy"
	stepHealthCheck = "health_check"
	stepResult      = "result"
)

var truncatedNotice = fmt.Sprintf("journal tronqué : limite de %d lignes ou de %d Kio atteinte, la suite n'est pas gardée",
	maxLogLines, maxLogBytes>>10)

// ansiEscape reconnaît les séquences de contrôle des terminaux (couleurs,
// déplacements du curseur).
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]|\x1b[@-_]`)

// journal enregistre le journal d'exécution d'un déploiement. Les lignes sont
// écrites en base par lots, puis diffusées : un abonné qui relit la base puis
// suit le flux ne perd ainsi aucune ligne (voir Service.Logs).
//
// close doit être appelée avant l'état final du déploiement : un déploiement
// terminé a tout son journal en base.
type journal struct {
	svc *Service
	id  string

	mu        sync.Mutex
	last      time.Time
	lines     int
	bytes     int
	truncated bool
	pending   []*store.DeploymentLog
	closed    bool

	stop chan struct{}
	done chan struct{}
}

func (s *Service) newJournal(id string) *journal {
	j := &journal{svc: s, id: id, stop: make(chan struct{}), done: make(chan struct{})}
	go j.loop()
	return j
}

func (j *journal) loop() {
	defer close(j.done)
	ticker := time.NewTicker(journalFlushInterval)
	defer ticker.Stop()
	for {
		select {
		case <-j.stop:
			return
		case <-ticker.C:
			j.flush()
		}
	}
}

// info ajoute une ligne de ForgeNet sur la sortie standard.
func (j *journal) info(step, message string) { j.add(store.LogStdout, step, message, false) }

// result ajoute le bilan du déploiement ; il passe même un journal tronqué.
func (j *journal) result(stream store.LogStream, message string) {
	j.add(stream, stepResult, message, true)
}

// output reçoit la sortie de docker compose pour l'étape step.
func (j *journal) output(step string) func(stream store.LogStream, line string) {
	return func(stream store.LogStream, line string) { j.add(stream, step, line, false) }
}

func (j *journal) add(stream store.LogStream, step, message string, always bool) {
	message = cleanLine(message)
	if message == "" {
		return
	}

	j.mu.Lock()
	defer j.mu.Unlock()
	if j.closed {
		return
	}
	if !always {
		if j.truncated {
			return
		}
		if j.lines >= maxLogLines || j.bytes+len(message) > maxLogBytes {
			j.truncated = true
			stream, message = store.LogStderr, truncatedNotice
		}
	}

	// Une date strictement croissante sert de curseur de reprise.
	now := time.Now().UTC().Round(0)
	if !now.After(j.last) {
		now = j.last.Add(time.Nanosecond)
	}
	j.last = now
	j.lines++
	j.bytes += len(message)
	j.pending = append(j.pending, &store.DeploymentLog{
		DeploymentID: j.id, EmittedAt: now, Stream: stream, Step: step, Message: message,
	})
}

// flush écrit les lignes en attente, puis les diffuse.
func (j *journal) flush() {
	j.mu.Lock()
	batch := j.pending
	j.pending = nil
	j.mu.Unlock()
	if len(batch) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	if err := j.svc.store.AppendDeploymentLogs(ctx, batch); err != nil {
		// Le journal aide à comprendre, il ne doit pas faire échouer le déploiement.
		j.svc.logger.Error().Err(err).Str("deployment_id", j.id).Int("lines", len(batch)).
			Msg("écriture du journal de déploiement")
	}
	for _, l := range batch {
		j.svc.events.publish(j.id, event{log: l})
	}
}

// close écrit les dernières lignes et arrête l'enregistreur ; les appels
// suivants ne font rien.
func (j *journal) close() {
	j.mu.Lock()
	if j.closed {
		j.mu.Unlock()
		return
	}
	j.closed = true
	j.mu.Unlock()

	close(j.stop)
	<-j.done
	j.flush()
}

// cleanLine rend une ligne sûre à stocker et à afficher : sans séquence de
// terminal ni caractère de contrôle, en UTF-8 valide (exigé par protobuf),
// et bornée.
func cleanLine(s string) string {
	s = ansiEscape.ReplaceAllString(s, "")
	s = strings.ToValidUTF8(s, "�")
	s = strings.Map(func(r rune) rune {
		if r == '\t' {
			return ' '
		}
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
	s = strings.TrimRight(s, " ")
	if len(s) > maxLogLineBytes {
		cut := maxLogLineBytes
		for cut > 0 && !utf8.RuneStart(s[cut]) {
			cut--
		}
		s = s[:cut] + "…"
	}
	return s
}
