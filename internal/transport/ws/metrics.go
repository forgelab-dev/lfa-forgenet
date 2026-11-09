// Package ws diffuse les métriques en direct par WebSocket.
//
// Un navigateur ne sait pas ouvrir un flux gRPC : le dashboard et
// l'application mobile reçoivent donc les mesures par WebSocket, compressées
// (permessage-deflate), au format JSON des messages du contrat api/
// (StreamSystemMetricsResponse, StreamContainerMetricsResponse).
package ws

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/metrics"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

const (
	// Subprotocol est le sous-protocole négocié par le serveur.
	Subprotocol = "forgenet.v1"
	// BearerSubprotocolPrefix précède le jeton quand le client ne peut pas
	// poser d'en-tête Authorization : l'API WebSocket des navigateurs ne
	// permet que la liste des sous-protocoles. Le jeton n'apparaît ainsi ni
	// dans l'URL ni dans les journaux d'accès. Le serveur ne le renvoie pas.
	BearerSubprotocolPrefix = "forgenet.bearer."

	// DefaultSystemInterval est l'intervalle par défaut du flux de l'hôte,
	// fixé par le contrat (metrics.proto).
	DefaultSystemInterval = 5 * time.Second
	// DefaultContainerInterval est celui du flux d'un conteneur.
	DefaultContainerInterval = 2 * time.Second
	// MaxInterval borne l'intervalle demandé ; le minimum est une seconde.
	MaxInterval = time.Minute

	// DefaultMaxStreams limite les flux simultanés : chacun lit /proc ou
	// interroge le démon Docker à chaque intervalle.
	DefaultMaxStreams = 64

	// DefaultRecheckEvery : un flux dure plus longtemps qu'un jeton. Le
	// jeton est revérifié à ce rythme, pour qu'une session révoquée ou
	// expirée coupe le flux.
	DefaultRecheckEvery = time.Minute
	writeTimeout        = 10 * time.Second
)

// Codes de fermeture propres à ForgeNet (plage 4000-4999 réservée aux
// applications par la RFC 6455).
const (
	// CloseUnauthorized : le jeton a expiré ou a été révoqué pendant le flux.
	CloseUnauthorized websocket.StatusCode = 4401
	// CloseNotFound : le conteneur a disparu pendant le flux.
	CloseNotFound websocket.StatusCode = 4404
)

// TokenValidator vérifie un jeton (auth.Service le satisfait).
type TokenValidator interface {
	ValidateToken(ctx context.Context, token string) (*auth.Claims, error)
}

// HostCollector mesure l'hôte.
type HostCollector interface {
	Collect(ctx context.Context) (metrics.Host, error)
}

// ContainerCollector mesure un conteneur.
type ContainerCollector interface {
	Collect(ctx context.Context) (metrics.Container, error)
}

// Config réunit les dépendances du gestionnaire.
type Config struct {
	Auth TokenValidator
	// NewHost et NewContainer créent un collecteur par flux : chacun garde
	// son relevé processeur précédent.
	NewHost      func() HostCollector
	NewContainer func(id string) ContainerCollector
	// OriginPatterns liste les origines autorisées en plus de la même
	// origine (motifs de github.com/coder/websocket). Les clients hors
	// navigateur n'envoient pas d'Origin et ne sont pas concernés.
	OriginPatterns []string
	// MaxStreams : 0 vaut DefaultMaxStreams.
	MaxStreams int
	// RecheckEvery : 0 vaut DefaultRecheckEvery.
	RecheckEvery time.Duration
	Logger       zerolog.Logger
}

// MetricsHandler sert /ws/v1/metrics/system et /ws/v1/metrics/containers/{id}.
type MetricsHandler struct {
	cfg    Config
	slots  chan struct{}
	logger zerolog.Logger

	// ctx est annulé par Close : http.Server.Shutdown n'attend ni ne ferme
	// les connexions détournées (hijacked) par WebSocket.
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewMetricsHandler crée le gestionnaire.
func NewMetricsHandler(cfg Config) *MetricsHandler {
	if cfg.MaxStreams <= 0 {
		cfg.MaxStreams = DefaultMaxStreams
	}
	if cfg.RecheckEvery <= 0 {
		cfg.RecheckEvery = DefaultRecheckEvery
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &MetricsHandler{
		cfg:    cfg,
		slots:  make(chan struct{}, cfg.MaxStreams),
		logger: cfg.Logger.With().Str("component", "transport.ws").Logger(),
		ctx:    ctx,
		cancel: cancel,
	}
}

// RegisterRoutes ajoute les routes au multiplexeur.
func (h *MetricsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /ws/v1/metrics/system", h.serveSystem)
	mux.HandleFunc("GET /ws/v1/metrics/containers/{id}", h.serveContainer)
}

// Close ferme les flux en cours (code 1001, « arrêt du serveur ») et attend
// qu'ils soient terminés.
func (h *MetricsHandler) Close() {
	h.cancel()
	h.wg.Wait()
}

func (h *MetricsHandler) serveSystem(w http.ResponseWriter, r *http.Request) {
	collector := h.cfg.NewHost()
	h.serve(w, r, DefaultSystemInterval, func(ctx context.Context) (proto.Message, error) {
		m, err := collector.Collect(ctx)
		if err != nil {
			return nil, err
		}
		return &forgenetv1.StreamSystemMetricsResponse{Metrics: toSystemMetrics(m)}, nil
	})
}

func (h *MetricsHandler) serveContainer(w http.ResponseWriter, r *http.Request) {
	collector := h.cfg.NewContainer(r.PathValue("id"))
	h.serve(w, r, DefaultContainerInterval, func(ctx context.Context) (proto.Message, error) {
		m, err := collector.Collect(ctx)
		if err != nil {
			return nil, err
		}
		return &forgenetv1.StreamContainerMetricsResponse{Metrics: toContainerMetrics(m)}, nil
	})
}

type sample func(ctx context.Context) (proto.Message, error)

// serve fait tous les refus en HTTP, avant la poignée de main : un client
// reçoit 400, 401, 403, 404 ou 503 plutôt qu'une connexion aussitôt fermée.
// La première mesure est prise avant, pour qu'un conteneur inconnu donne 404.
func (h *MetricsHandler) serve(w http.ResponseWriter, r *http.Request, defaultInterval time.Duration, next sample) {
	h.wg.Add(1)
	defer h.wg.Done()

	token := requestToken(r)
	if !h.authorize(r.Context(), w, token) {
		return
	}
	interval, err := parseInterval(r.URL.Query().Get("interval"), defaultInterval)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	select {
	case h.slots <- struct{}{}:
		defer func() { <-h.slots }()
	default:
		http.Error(w, "trop de flux de métriques ouverts, réessayez plus tard", http.StatusServiceUnavailable)
		return
	}

	first, err := next(r.Context())
	if err != nil {
		h.refuse(w, err)
		return
	}

	// Le délai d'écriture du serveur HTTP s'appliquerait à toute la vie de
	// la connexion détournée et la couperait au bout de quelques secondes.
	rc := http.NewResponseController(w)
	_ = rc.SetReadDeadline(time.Time{})
	_ = rc.SetWriteDeadline(time.Time{})

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:    []string{Subprotocol},
		OriginPatterns:  h.cfg.OriginPatterns,
		CompressionMode: websocket.CompressionContextTakeover,
	})
	if err != nil {
		// Accept a déjà répondu (origine refusée, poignée de main invalide).
		h.logger.Debug().Err(err).Msg("poignée de main WebSocket refusée")
		return
	}
	h.stream(conn, token, interval, first, next)
}

func (h *MetricsHandler) stream(conn *websocket.Conn, token string, interval time.Duration, first proto.Message, next sample) {
	defer func() { _ = conn.CloseNow() }()

	// CloseRead lit et ignore les messages du client : ctx est annulé quand
	// il ferme la connexion. Il ne dérive pas de h.ctx : la bibliothèque
	// coupe la connexion sans message dès que le contexte d'une lecture est
	// annulé, et le client ne saurait pas que le serveur s'arrête.
	ctx := conn.CloseRead(context.Background())

	if err := send(ctx, conn, first); err != nil {
		return
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	lastCheck := time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		case <-h.ctx.Done():
			_ = conn.Close(websocket.StatusGoingAway, "arrêt du serveur")
			return
		case <-ticker.C:
		}

		if time.Since(lastCheck) >= h.cfg.RecheckEvery {
			lastCheck = time.Now()
			if !h.stillAuthorized(ctx, token) {
				_ = conn.Close(CloseUnauthorized, "jeton expiré ou révoqué")
				return
			}
		}

		msg, err := next(ctx)
		switch {
		case ctx.Err() != nil:
			continue // la fermeture est traitée en tête de boucle
		case errors.Is(err, orchestrator.ErrNotFound):
			_ = conn.Close(CloseNotFound, "conteneur introuvable")
			return
		case err != nil:
			h.logger.Error().Err(err).Msg("mesure impossible, flux fermé")
			_ = conn.Close(websocket.StatusInternalError, "mesure impossible")
			return
		}
		if err := send(ctx, conn, msg); err != nil {
			return
		}
	}
}

func send(ctx context.Context, conn *websocket.Conn, msg proto.Message) error {
	payload, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()
	return conn.Write(ctx, websocket.MessageText, payload)
}

// authorize répond 401, 403 ou 503 et renvoie false si l'appel est refusé.
// Comme la lecture des conteneurs en gRPC, les métriques sont ouvertes aux
// membres.
func (h *MetricsHandler) authorize(ctx context.Context, w http.ResponseWriter, token string) bool {
	claims, err := h.cfg.Auth.ValidateToken(ctx, token)
	switch {
	case err == nil:
	case auth.IsCredentialError(err):
		w.Header().Set("WWW-Authenticate", "Bearer")
		http.Error(w, "authentification requise ou jeton invalide", http.StatusUnauthorized)
		return false
	default:
		h.logger.Error().Err(err).Msg("vérification du jeton impossible")
		http.Error(w, "vérification de l'authentification indisponible", http.StatusServiceUnavailable)
		return false
	}
	if !claims.Role.Grants(store.RoleMember) {
		http.Error(w, "droits insuffisants pour cette opération", http.StatusForbidden)
		return false
	}
	return true
}

// stillAuthorized ne coupe le flux que si le jeton est refusé : une base
// momentanément indisponible ne ferme pas les flux en cours.
func (h *MetricsHandler) stillAuthorized(ctx context.Context, token string) bool {
	claims, err := h.cfg.Auth.ValidateToken(ctx, token)
	if err != nil {
		return !auth.IsCredentialError(err)
	}
	return claims.Role.Grants(store.RoleMember)
}

func (h *MetricsHandler) refuse(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, orchestrator.ErrNotFound):
		http.Error(w, "conteneur introuvable", http.StatusNotFound)
	case errors.Is(err, orchestrator.ErrInvalidArgument):
		http.Error(w, "identifiant de conteneur invalide", http.StatusBadRequest)
	default:
		h.logger.Error().Err(err).Msg("première mesure impossible")
		http.Error(w, "mesure indisponible", http.StatusServiceUnavailable)
	}
}

// requestToken lit le jeton dans l'en-tête Authorization, sinon dans les
// sous-protocoles proposés.
func requestToken(r *http.Request) string {
	if header := strings.TrimSpace(r.Header.Get("Authorization")); header != "" {
		if scheme, token, ok := strings.Cut(header, " "); ok && strings.EqualFold(scheme, "bearer") {
			return strings.TrimSpace(token)
		}
		return ""
	}
	for _, line := range r.Header.Values("Sec-WebSocket-Protocol") {
		for _, p := range strings.Split(line, ",") {
			if token, ok := strings.CutPrefix(strings.TrimSpace(p), BearerSubprotocolPrefix); ok {
				return token
			}
		}
	}
	return ""
}

// parseInterval lit l'intervalle en secondes, comme interval_seconds du
// contrat : absent ou 0 donne la valeur par défaut.
func parseInterval(raw string, def time.Duration) (time.Duration, error) {
	if raw == "" {
		return def, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 || time.Duration(n)*time.Second > MaxInterval {
		return 0, errors.New("interval doit être un nombre de secondes entre 1 et 60")
	}
	if n == 0 {
		return def, nil
	}
	return time.Duration(n) * time.Second, nil
}
