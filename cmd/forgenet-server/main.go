// Package main est le point d'entrée exécutable du binaire forgenet-server (control plane).
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/health"
	"github.com/forgelab-dev/lfa-forgenet/internal/logging"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	transportHTTP "github.com/forgelab-dev/lfa-forgenet/internal/transport/http"
	"github.com/forgelab-dev/lfa-forgenet/internal/version"
)

func main() {
	var (
		configPath  string
		showVersion bool
	)

	flag.StringVar(&configPath, "config", "", "Chemin vers le fichier de configuration (YAML)")
	flag.StringVar(&configPath, "c", "", "Chemin vers le fichier de configuration (raccourci)")
	flag.BoolVar(&showVersion, "version", false, "Affiche la version du binaire et quitte")
	flag.BoolVar(&showVersion, "v", false, "Affiche la version du binaire (raccourci)")
	flag.Parse()

	if showVersion {
		info := version.Get()
		fmt.Printf("forgenet-server %s (%s, %s, %s)\n", info.Version, info.GitCommit, info.BuildDate, info.GoVersion)
		os.Exit(0)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de l'initialisation de la configuration : %v\n", err)
		os.Exit(1)
	}

	logger := logging.New(cfg.Log)
	buildInfo := version.Get()

	logger.Info().
		Str("version", buildInfo.Version).
		Str("commit", buildInfo.GitCommit).
		Str("environment", cfg.Environment).
		Str("http_addr", cfg.Server.HTTPAddress()).
		Str("grpc_addr", cfg.Server.GRPCAddress()).
		Msg("démarrage de ForgeNet Server")

	// Service de santé partagé (liveness / readiness)
	healthService := health.NewService()

	// Initialisation des transports HTTP et gRPC sans logique métier
	httpServer := transportHTTP.NewServer(cfg.Server, healthService, logger)
	grpcServer := transportGRPC.NewServer(cfg.Server, healthService, logger)

	// Contexte d'écoute des signaux d'arrêt système
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	errChan := make(chan error, 2)

	// Démarrage du serveur HTTP
	go func() {
		if err := httpServer.Start(); err != nil {
			errChan <- fmt.Errorf("serveur HTTP : %w", err)
		}
	}()

	// Démarrage du serveur gRPC
	go func() {
		if err := grpcServer.Start(); err != nil {
			errChan <- fmt.Errorf("serveur gRPC : %w", err)
		}
	}()

	// Attente d'un signal d'arrêt ou d'une défaillance réseau
	var fatalErr error
	select {
	case <-ctx.Done():
		logger.Info().Msg("signal d'arrêt reçu, initialisation de la fermeture gracieuse...")
	case err := <-errChan:
		logger.Error().Err(err).Msg("défaillance inattendue d'un serveur de transport")
		fatalErr = err
	}

	// 1. Bascule de l'état de santé en NOT_SERVING pour drainer le trafic entrant
	healthService.SetServing(false)

	// 2. Contexte d'arrêt gracieux avec timeout configuré
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error().Err(err).Msg("erreur lors de l'arrêt gracieux du serveur HTTP")
		}
	}()

	go func() {
		defer wg.Done()
		if err := grpcServer.Shutdown(shutdownCtx); err != nil {
			logger.Error().Err(err).Msg("erreur lors de l'arrêt gracieux du serveur gRPC")
		}
	}()

	wg.Wait()

	if fatalErr != nil {
		logger.Error().Err(fatalErr).Msg("fermeture du serveur suite à une défaillance critique")
		os.Exit(1)
	}

	logger.Info().Msg("serveurs HTTP et gRPC arrêtés proprement. Au revoir.")
}
