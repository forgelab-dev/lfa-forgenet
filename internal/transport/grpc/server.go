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

// NewServer instancie un serveur gRPC standardisé avec intercepteurs et health check.
func NewServer(cfg config.ServerConfig, healthSvc health.Service, logger zerolog.Logger) *Server {
	componentLogger := logger.With().Str("component", "transport.grpc").Logger()

	serverOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			RecoveryUnaryInterceptor(componentLogger),
			LoggingUnaryInterceptor(componentLogger),
		),
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
	// Si un service spécifique est demandé (ex: "readiness"), on évalue la readiness
	if req.GetService() == "readiness" {
		ready, _ := h.healthService.IsReady(ctx)
		status := grpcHealth.HealthCheckResponse_SERVING
		if !ready {
			status = grpcHealth.HealthCheckResponse_NOT_SERVING
		}
		return &grpcHealth.HealthCheckResponse{Status: status}, nil
	}

	// Par défaut, vérification de vivacité (liveness)
	live := h.healthService.IsLive(ctx)
	status := grpcHealth.HealthCheckResponse_SERVING
	if !live {
		status = grpcHealth.HealthCheckResponse_NOT_SERVING
	}
	return &grpcHealth.HealthCheckResponse{Status: status}, nil
}

func (h *healthHandler) Watch(_ *grpcHealth.HealthCheckRequest, stream grpcHealth.Health_WatchServer) error {
	status := grpcHealth.HealthCheckResponse_SERVING
	if !h.healthService.IsLive(stream.Context()) {
		status = grpcHealth.HealthCheckResponse_NOT_SERVING
	}
	return stream.Send(&grpcHealth.HealthCheckResponse{Status: status})
}
