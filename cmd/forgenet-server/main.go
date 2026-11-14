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
	"path/filepath"
	"syscall"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/auth"
	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/deploy"
	"github.com/forgelab-dev/lfa-forgenet/internal/domains"
	"github.com/forgelab-dev/lfa-forgenet/internal/health"
	"github.com/forgelab-dev/lfa-forgenet/internal/logging"
	"github.com/forgelab-dev/lfa-forgenet/internal/metrics"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/compose"
	"github.com/forgelab-dev/lfa-forgenet/internal/orchestrator/docker"
	"github.com/forgelab-dev/lfa-forgenet/internal/proxy"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	"github.com/forgelab-dev/lfa-forgenet/internal/store/sqlite"
	"github.com/forgelab-dev/lfa-forgenet/internal/transport/connectapi"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	transportHTTP "github.com/forgelab-dev/lfa-forgenet/internal/transport/http"
	"github.com/forgelab-dev/lfa-forgenet/internal/transport/ws"
	"github.com/forgelab-dev/lfa-forgenet/internal/version"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// main se contente de traduire l'erreur de run en code de sortie : toute la
// logique vit dans run, pour que ses defer s'exécutent avant os.Exit.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, os.Args[1:], os.Stdin, os.Stdout); err != nil {
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
// « forgenet-server user … » administre les comptes au lieu de servir.
func run(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) > 0 && args[0] == "user" {
		return runUser(ctx, args[1:], stdin, stdout)
	}

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
	return runServer(ctx, cfg)
}

// runServer assemble les composants du serveur et sert jusqu'à l'annulation de ctx.
func runServer(ctx context.Context, cfg *config.Config) error {
	logger := logging.New(cfg.Log)
	buildInfo := version.Get()
	logger.Info().
		Str("version", buildInfo.Version).
		Str("commit", buildInfo.GitCommit).
		Str("environment", cfg.Environment).
		Str("http_addr", cfg.Server.HTTPAddress()).
		Str("grpc_addr", cfg.Server.GRPCAddress()).
		Msg("démarrage de ForgeNet Server")

	// La base est ouverte (et migrée) avant d'écouter : un serveur qui accepte
	// du trafic sans schéma à jour répondrait par des erreurs.
	db, err := sqlite.Open(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("base de données : %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error().Err(err).Msg("fermeture de la base de données")
		}
	}()

	healthService := health.NewService()
	healthService.RegisterReadinessChecker("database", health.CheckFunc(db.Ping))

	auditor := audit.NewAuditor(db.Audit(), logger)
	authService := auth.NewService(db, cfg.Auth, auditor)
	guard := auth.NewGuard(authService, auth.DefaultPolicy(), logger)
	grpcServer := transportGRPC.NewServer(cfg.Server, healthService, guard, logger)

	// Le client ne se connecte pas encore ici : une erreur signale une
	// configuration invalide, et l'indisponibilité du démon se lit dans /readyz.
	engine, err := docker.New(cfg.Docker)
	if err != nil {
		return err
	}
	defer func() { _ = engine.Close() }()
	healthService.RegisterReadinessChecker("docker", health.CheckFunc(engine.Ping))

	// Les stacks Compose visent le même démon que les conteneurs (ADR-0003).
	stacks, err := compose.New(cfg.Docker)
	if err != nil {
		return err
	}
	healthService.RegisterReadinessChecker("compose", health.CheckFunc(stacks.Ping))

	// Les domaines n'existent qu'avec le reverse proxy activé. Après chaque
	// déploiement réussi, les services qui ont un domaine sont rebranchés sur
	// le réseau partagé avec Caddy, leurs conteneurs ayant pu être recréés (ADR-0008).
	domainService, err := newDomainService(cfg, db, engine, stacks, auditor, logger)
	if err != nil {
		return err
	}
	var deployOpts []deploy.Option
	if domainService != nil {
		deployOpts = append(deployOpts, deploy.WithAfterDeploy(func(ctx context.Context, projectID string) {
			if err := domainService.Reattach(ctx, projectID); err != nil {
				logger.Warn().Err(err).Str("project_id", projectID).Msg("domaines non rebranchés après le déploiement")
			}
		}))
	}
	deployments := deploy.New(db.Deployments(), stacks, auditor, logger, deployOpts...)

	// Historique court des mesures, gardé en mémoire (ADR-0010) : le serveur
	// relève lui-même l'hôte et les conteneurs, sans rien écrire sur le disque.
	recorder := metrics.NewRecorder(
		metrics.NewHostCollector(metrics.GopsutilProbe{}, engine, cfg.Metrics.StoragePath),
		engine, engine, metrics.RecorderConfig{}, logger)

	// Les mêmes services sont servis par deux portes : gRPC (CLI, mobile) et
	// Connect sur le serveur HTTP (dashboard), avec le même garde (ADR-0006).
	httpServer := transportHTTP.NewServer(cfg.Server, cfg.Health, healthService, logger)
	connectAPI := connectapi.New(guard, logger)
	authAPI := transportGRPC.NewAuthServiceServer(authService, db.Users(), logger)
	containerAPI := transportGRPC.NewContainerServiceServer(engine, auditor, logger)
	deploymentAPI := transportGRPC.NewDeploymentServiceServer(deployments, logger)
	domainAPI := transportGRPC.NewDomainServiceServer(domainService, logger)
	metricsAPI := transportGRPC.NewMetricsServiceServer(recorder, engine, logger)
	for _, registrar := range []grpc.ServiceRegistrar{grpcServer.GRPCServer(), connectAPI.Registrar()} {
		forgenetv1.RegisterAuthServiceServer(registrar, authAPI)
		forgenetv1.RegisterContainerServiceServer(registrar, containerAPI)
		forgenetv1.RegisterDeploymentServiceServer(registrar, deploymentAPI)
		forgenetv1.RegisterDomainServiceServer(registrar, domainAPI)
		forgenetv1.RegisterMetricsServiceServer(registrar, metricsAPI)
	}
	connectAPI.Mount(httpServer.Mux())

	// Les métriques en direct partent par WebSocket : un navigateur ne lit
	// pas un flux gRPC. Chaque flux a ses collecteurs (relevé précédent).
	metricsStreams := ws.NewMetricsHandler(ws.Config{
		Auth: authService,
		NewHost: func() ws.HostCollector {
			return metrics.NewHostCollector(metrics.GopsutilProbe{}, engine, cfg.Metrics.StoragePath)
		},
		NewContainer: func(id string) ws.ContainerCollector { return metrics.NewContainerCollector(engine, id) },
		Sweeps:       recorder,
		Logger:       logger,
	})
	metricsStreams.RegisterRoutes(httpServer.Mux())
	httpServer.OnShutdown(metricsStreams.Close)

	// Le dashboard compilé est servi à la même origine que l'API (ADR-0006,
	// ADR-0009) ; sans web.root, comme en développement, l'API seule.
	if cfg.Web.Root != "" {
		dashboard, err := transportHTTP.NewDashboardHandler(cfg.Web.Root)
		if err != nil {
			return err
		}
		defer func() { _ = dashboard.Close() }()
		dashboard.RegisterRoutes(httpServer.Mux())
	}

	// L'exécuteur des déploiements s'arrête avec les serveurs : celui en cours
	// a le délai de grâce pour finir (ADR-0004).
	servers := []server{
		{name: "HTTP", transport: httpServer},
		{name: "gRPC", transport: grpcServer},
		{name: "déploiements", transport: deployments},
		{name: "relevé des mesures", transport: recorder},
	}
	if domainService != nil {
		// Réconcilie Caddy, le réseau et les conteneurs au démarrage, et
		// recommence tant que Caddy ou Docker ne répondent pas.
		servers = append(servers, server{name: "domaines", transport: domainService})
	}

	if err := serve(ctx, servers, healthService, cfg.Server, logger); err != nil {
		logger.Error().Err(err).Msg("arrêt du serveur suite à une défaillance")
		return errAlreadyReported
	}

	logger.Info().Msg("serveurs arrêtés proprement")
	return nil
}

// newDomainService prépare le service des domaines, ou renvoie nil quand le
// reverse proxy est désactivé. Une adresse d'administration exposée ou un
// courriel invalide empêchent le démarrage (ADR-0007).
func newDomainService(cfg *config.Config, db store.Store, engine *docker.Engine, stacks *compose.Engine, auditor audit.Auditor, logger zerolog.Logger) (*domains.Service, error) {
	if !cfg.Proxy.Enabled {
		return nil, nil
	}
	copyPath := cfg.Proxy.CaddyfilePath
	if copyPath == "" {
		copyPath = filepath.Join(filepath.Dir(cfg.Database.Path), "Caddyfile")
	}
	settings := proxy.Settings{
		AdminAddress: cfg.Proxy.AdminAddress,
		Email:        cfg.Proxy.Email,
		LocalCerts:   cfg.Proxy.LocalCerts,
		HTTPPort:     cfg.Proxy.HTTPPort,
		HTTPSPort:    cfg.Proxy.HTTPSPort,
	}
	if cfg.Proxy.DashboardHost != "" {
		upstream, err := proxy.ParseUpstream(cfg.Proxy.DashboardUpstream)
		if err != nil {
			return nil, fmt.Errorf("configuration du reverse proxy : proxy.dashboard_upstream : %w", err)
		}
		settings.Dashboard = &proxy.Route{Host: cfg.Proxy.DashboardHost, Upstream: upstream}
	}
	manager, err := proxy.NewManager(settings, copyPath, logger)
	if err != nil {
		return nil, fmt.Errorf("configuration du reverse proxy : %w", err)
	}
	return domains.New(db.Domains(), manager, engine, stacks, auditor, domains.Config{
		Network:        cfg.Proxy.Network,
		CaddyContainer: cfg.Proxy.CaddyContainer,
		DashboardHost:  cfg.Proxy.DashboardHost,
	}, logger), nil
}
