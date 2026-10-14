// Package main est le point d'entrée exécutable du binaire forgenet-server (control plane).
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/health"
	"github.com/forgelab-dev/lfa-forgenet/internal/logging"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	transportHTTP "github.com/forgelab-dev/lfa-forgenet/internal/transport/http"
	"github.com/forgelab-dev/lfa-forgenet/internal/version"
)

// main se contente de traduire l'erreur de run en code de sortie : toute la
// logique vit dans run, pour que ses defer s'exécutent avant os.Exit.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, os.Args[1:], os.Stdout); err != nil {
		if !errors.Is(err, errAlreadyReported) {
			fmt.Fprintf(os.Stderr, "forgenet-server : %v\n", err)
		}
		stop()
		os.Exit(1)
	}
}

// errAlreadyReported signale une erreur déjà écrite dans le journal structuré.
var errAlreadyReported = errors.New("erreur déjà journalisée")

// run lit les options, charge la configuration et sert jusqu'à l'annulation de ctx.
func run(ctx context.Context, args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("forgenet-server", flag.ContinueOnError)
	var (
		configPath  string
		showVersion bool
	)
	flags.StringVar(&configPath, "config", "", "Chemin vers le fichier de configuration (YAML)")
	flags.StringVar(&configPath, "c", "", "Chemin vers le fichier de configuration (raccourci)")
	flags.BoolVar(&showVersion, "version", false, "Affiche la version du binaire et quitte")
	flags.BoolVar(&showVersion, "v", false, "Affiche la version du binaire (raccourci)")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if showVersion {
		info := version.Get()
		_, err := fmt.Fprintf(stdout, "forgenet-server %s (%s, %s, %s)\n", info.Version, info.GitCommit, info.BuildDate, info.GoVersion)
		return err
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
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

	healthService := health.NewService()
	servers := []server{
		{name: "HTTP", transport: transportHTTP.NewServer(cfg.Server, healthService, logger)},
		{name: "gRPC", transport: transportGRPC.NewServer(cfg.Server, healthService, logger)},
	}

	if err := serve(ctx, servers, healthService, cfg.Server, logger); err != nil {
		logger.Error().Err(err).Msg("arrêt du serveur suite à une défaillance")
		return errAlreadyReported
	}

	logger.Info().Msg("serveurs arrêtés proprement")
	return nil
}
