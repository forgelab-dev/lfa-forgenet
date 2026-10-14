package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/health"
)

// Server gère le cycle de vie du serveur HTTP de ForgeNet.
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	logger     zerolog.Logger
	listener   net.Listener
}

// NewServer instancie un nouveau serveur HTTP avec ses timeouts et routes de base.
func NewServer(cfg config.ServerConfig, healthCfg config.HealthConfig, healthSvc health.Service, logger zerolog.Logger) *Server {
	mux := http.NewServeMux()
	componentLogger := logger.With().Str("component", "transport.http").Logger()

	healthHandler := NewHealthHandler(healthSvc, healthCfg.ExposeDetails, componentLogger)
	healthHandler.RegisterRoutes(mux)

	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Debug().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Dur("duration", time.Since(start)).
				Msg("requête HTTP traitée")
		})
	}

	httpServer := &http.Server{
		Addr:         cfg.HTTPAddress(),
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		httpServer: httpServer,
		mux:        mux,
		logger:     componentLogger,
	}
}

// Mux retourne le multiplexeur HTTP sous-jacent pour l'enregistrement de routes complémentaires.
func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

// Start lance l'écoute HTTP. Si un listener personnalisé a été défini (ex: pour les tests), il est utilisé.
func (s *Server) Start() error {
	s.logger.Info().Str("addr", s.httpServer.Addr).Msg("démarrage du serveur HTTP")

	var err error
	if s.listener != nil {
		err = s.httpServer.Serve(s.listener)
	} else {
		err = s.httpServer.ListenAndServe()
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("erreur du serveur HTTP : %w", err)
	}
	return nil
}

// SetListener permet d'injecter un net.Listener (notamment utile pour les tests avec port éphémère).
func (s *Server) SetListener(ln net.Listener) {
	s.listener = ln
	s.httpServer.Addr = ln.Addr().String()
}

// Addr retourne l'adresse réseau d'écoute du serveur.
func (s *Server) Addr() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return s.httpServer.Addr
}

// Shutdown arrête proprement le serveur HTTP avec un contexte d'attente (grace period).
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info().Msg("arrêt gracieux du serveur HTTP")
	return s.httpServer.Shutdown(ctx)
}
