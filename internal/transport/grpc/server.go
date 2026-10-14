// Package grpc gère le serveur gRPC et ses intercepteurs.
// Ce package est une couche de transport pure, sans logique métier.
package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	grpcHealth "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/health"
)

// Server encapsule l'instance du serveur gRPC et gère son cycle de vie.
type Server struct {
	grpcServer *grpc.Server
	addr       string
	logger     zerolog.Logger
	listener   net.Listener
}

// NewServer instancie un serveur gRPC standardisé avec intercepteurs (recovery, logging, auth) et health check.
func NewServer(cfg config.ServerConfig, healthSvc health.Service, authSvc auth.Service, logger zerolog.Logger) *Server {
	componentLogger := logger.With().Str("component", "transport.grpc").Logger()

	exemptMethods := map[string]bool{
		"/grpc.health.v1.Health/Check": true,
		"/grpc.health.v1.Health/Watch": true,
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		RecoveryUnaryInterceptor(componentLogger),
		LoggingUnaryInterceptor(componentLogger),
	}
	if authSvc != nil {
		unaryInterceptors = append(unaryInterceptors, auth.UnaryServerInterceptor(authSvc, exemptMethods))
	}

	var serverOptions []grpc.ServerOption
	serverOptions = append(serverOptions, grpc.ChainUnaryInterceptor(unaryInterceptors...))

	if authSvc != nil {
		serverOptions = append(serverOptions, grpc.ChainStreamInterceptor(auth.StreamServerInterceptor(authSvc, exemptMethods)))
	}

	grpcServer := grpc.NewServer(serverOptions...)

	// Enregistrement du service de santé standard gRPC
	healthHandler := newHealthHandler(healthSvc)
	grpcHealth.RegisterHealthServer(grpcServer, healthHandler)

	return &Server{
		grpcServer: grpcServer,
		addr:       cfg.GRPCAddress(),
		logger:     componentLogger,
	}
}

// GRPCServer retourne l'instance *grpc.Server sous-jacente pour enregistrer les services de domaine.
func (s *Server) GRPCServer() *grpc.Server {
	return s.grpcServer
}

// SetListener permet d'injecter un net.Listener existant (très utile pour les tests d'intégration avec port dynamique).
func (s *Server) SetListener(ln net.Listener) {
	s.listener = ln
	s.addr = ln.Addr().String()
}

// Addr retourne l'adresse d'écoute configurée ou effective.
func (s *Server) Addr() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return s.addr
}

// Start lance l'écoute et le traitement des requêtes gRPC.
func (s *Server) Start() error {
	s.logger.Info().Str("addr", s.addr).Msg("démarrage du serveur gRPC")

	if s.listener == nil {
		ln, err := net.Listen("tcp", s.addr)
		if err != nil {
			return fmt.Errorf("écoute réseau gRPC impossible sur %s : %w", s.addr, err)
		}
		s.listener = ln
	}

	if err := s.grpcServer.Serve(s.listener); err != nil {
		return fmt.Errorf("erreur du serveur gRPC : %w", err)
	}
	return nil
}

// Shutdown arrête proprement le serveur gRPC dans la limite impartie par le contexte.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info().Msg("arrêt gracieux du serveur gRPC")

	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		return nil
	case <-ctx.Done():
		s.logger.Warn().Msg("délai de grâce dépassé, arrêt forcé du serveur gRPC")
		s.grpcServer.Stop()
		return ctx.Err()
	}
}

// healthHandler implémente l'interface grpc_health_v1.HealthServer sans logique métier,
// en déléguant intégralement au service de santé interne.
type healthHandler struct {
	grpcHealth.UnimplementedHealthServer
	healthService health.Service
}

func newHealthHandler(healthService health.Service) *healthHandler {
	return &healthHandler{
		healthService: healthService,
	}
}

func (h *healthHandler) Check(ctx context.Context, req *grpcHealth.HealthCheckRequest) (*grpcHealth.HealthCheckResponse, error) {
	// Le service « readiness » interroge les dépendances ; par défaut, seule la vivacité compte.
	if req.GetService() == "readiness" {
		ready, _ := h.healthService.IsReady(ctx)
		return servingStatus(ready), nil
	}
	return servingStatus(h.healthService.IsLive(ctx)), nil
}

func (h *healthHandler) Watch(_ *grpcHealth.HealthCheckRequest, stream grpcHealth.Health_WatchServer) error {
	return stream.Send(servingStatus(h.healthService.IsLive(stream.Context())))
}

func servingStatus(serving bool) *grpcHealth.HealthCheckResponse {
	if serving {
		return &grpcHealth.HealthCheckResponse{Status: grpcHealth.HealthCheckResponse_SERVING}
	}
	return &grpcHealth.HealthCheckResponse{Status: grpcHealth.HealthCheckResponse_NOT_SERVING}
}
