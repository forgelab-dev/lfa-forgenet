package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/update"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

// staleRun : au-delà, une mise à jour encore « en cours » a été interrompue
// (conteneur tué, hôte redémarré) ; le moteur, lui, finit en quelques minutes.
const staleRun = 30 * time.Minute

// UpdateDeps rassemble ce dont le service des mises à jour a besoin.
type UpdateDeps struct {
	Checker      *update.Checker
	CheckEnabled bool
	Current      string // version en service
	InstallDir   string // vide : installation non standard, mise à jour à la main
	DataDir      string // dossier de la base, où le moteur écrit maj/etat.json
	Launch       func(context.Context, update.LaunchOptions) (string, error)
	Auditor      audit.Auditor
	Now          func() time.Time
}

// UpdateServiceServer expose les mises à jour selon le contrat
// forgenet.v1.UpdateService, réservé aux administrateurs (ADR-0011).
type UpdateServiceServer struct {
	forgenetv1.UnimplementedUpdateServiceServer
	deps   UpdateDeps
	logger zerolog.Logger
}

var _ forgenetv1.UpdateServiceServer = (*UpdateServiceServer)(nil)

// NewUpdateServiceServer instancie le service gRPC des mises à jour.
func NewUpdateServiceServer(deps UpdateDeps, logger zerolog.Logger) *UpdateServiceServer {
	if deps.Now == nil {
		deps.Now = time.Now
	}
	return &UpdateServiceServer{deps: deps, logger: logger.With().Str("component", "grpc.update_service").Logger()}
}

// GetUpdateStatus donne la version en service, la version disponible et la
// dernière mise à jour lancée.
func (s *UpdateServiceServer) GetUpdateStatus(_ context.Context, _ *forgenetv1.GetUpdateStatusRequest) (*forgenetv1.GetUpdateStatusResponse, error) {
	return s.status(s.deps.Checker.Snapshot()), nil
}

// CheckForUpdates interroge les releases tout de suite.
func (s *UpdateServiceServer) CheckForUpdates(ctx context.Context, _ *forgenetv1.CheckForUpdatesRequest) (*forgenetv1.CheckForUpdatesResponse, error) {
	if !s.deps.CheckEnabled {
		return nil, status.Error(codes.FailedPrecondition, "vérification des mises à jour désactivée (updates.check: false)")
	}
	return &forgenetv1.CheckForUpdatesResponse{Status: s.status(s.deps.Checker.Check(ctx))}, nil
}

// StartUpdate lance la mise à jour vers la version disponible.
func (s *UpdateServiceServer) StartUpdate(ctx context.Context, _ *forgenetv1.StartUpdateRequest) (*forgenetv1.StartUpdateResponse, error) {
	if s.deps.InstallDir == "" {
		return nil, status.Error(codes.FailedPrecondition, "installation non standard : mettez à jour avec install.sh --update")
	}
	available := s.deps.Checker.Snapshot().Available
	if available == nil {
		return nil, status.Error(codes.FailedPrecondition, "aucune version plus récente n'est disponible")
	}
	if run := s.lastRun(); run != nil && run.GetResult() == forgenetv1.UpdateResult_UPDATE_RESULT_RUNNING {
		return nil, status.Error(codes.Aborted, update.ErrBusy.Error())
	}
	// La version est celle annoncée à l'administrateur, épinglée : le moteur
	// installe exactement ce qu'il a vu, même si une autre sort entre-temps.
	id, err := s.deps.Launch(ctx, update.LaunchOptions{InstallDir: s.deps.InstallDir, To: available.Version})
	switch {
	case errors.Is(err, update.ErrBusy):
		return nil, status.Error(codes.Aborted, err.Error())
	case err != nil:
		s.logger.Error().Err(err).Msg("lancement de la mise à jour impossible")
		return nil, status.Error(codes.Internal, "lancement de la mise à jour impossible : voir les journaux du serveur")
	}
	_ = s.deps.Auditor.Log(ctx, "update.start", "forgenet", map[string]any{
		"from": s.deps.Current, "to": available.Version, "updater": id,
	})
	s.logger.Info().Str("to", available.Version).Str("updater", id).Msg("mise à jour lancée")
	return &forgenetv1.StartUpdateResponse{ToVersion: available.Version}, nil
}

func (s *UpdateServiceServer) status(snap update.Snapshot) *forgenetv1.GetUpdateStatusResponse {
	resp := &forgenetv1.GetUpdateStatusResponse{
		CurrentVersion: s.deps.Current,
		CheckError:     snap.Err,
		CheckEnabled:   s.deps.CheckEnabled,
		Supported:      s.deps.InstallDir != "",
		LastRun:        s.lastRun(),
	}
	if snap.Available != nil {
		resp.Available = &forgenetv1.AvailableRelease{Version: snap.Available.Version, NotesUrl: snap.Available.NotesURL}
	}
	if !snap.CheckedAt.IsZero() {
		resp.CheckedAt = timestamppb.New(snap.CheckedAt)
	}
	return resp
}

// lastRun lit l'état écrit par le moteur ; illisible ou absent : rien.
func (s *UpdateServiceServer) lastRun() *forgenetv1.UpdateRun {
	st, err := update.ReadState(s.deps.DataDir)
	if err != nil {
		s.logger.Warn().Err(err).Msg("état de la dernière mise à jour illisible")
		return nil
	}
	if st == nil {
		return nil
	}
	run := &forgenetv1.UpdateRun{
		FromVersion: st.From,
		ToVersion:   st.To,
		Step:        steps[st.Step],
		Result:      results[st.Result],
		Reason:      st.Reason,
		Backup:      st.Backup,
		StartedAt:   timestamppb.New(st.Started),
	}
	if !st.Finished.IsZero() {
		run.FinishedAt = timestamppb.New(st.Finished)
	}
	if run.Result == forgenetv1.UpdateResult_UPDATE_RESULT_RUNNING && s.deps.Now().Sub(st.Started) > staleRun {
		run.Result = forgenetv1.UpdateResult_UPDATE_RESULT_FAILED
		run.Reason = "mise à jour interrompue : aucune nouvelle du conteneur forgenet-updater depuis plus de 30 minutes"
	}
	return run
}

var steps = map[string]forgenetv1.UpdateStep{
	update.StepCheck:    forgenetv1.UpdateStep_UPDATE_STEP_CHECK,
	update.StepPull:     forgenetv1.UpdateStep_UPDATE_STEP_PULL,
	update.StepBackup:   forgenetv1.UpdateStep_UPDATE_STEP_BACKUP,
	update.StepSwitch:   forgenetv1.UpdateStep_UPDATE_STEP_SWITCH,
	update.StepHealth:   forgenetv1.UpdateStep_UPDATE_STEP_HEALTH,
	update.StepRollback: forgenetv1.UpdateStep_UPDATE_STEP_ROLLBACK,
	update.StepDone:     forgenetv1.UpdateStep_UPDATE_STEP_DONE,
}

var results = map[string]forgenetv1.UpdateResult{
	update.ResultRunning:    forgenetv1.UpdateResult_UPDATE_RESULT_RUNNING,
	update.ResultSucceeded:  forgenetv1.UpdateResult_UPDATE_RESULT_SUCCEEDED,
	update.ResultRolledBack: forgenetv1.UpdateResult_UPDATE_RESULT_ROLLED_BACK,
	update.ResultFailed:     forgenetv1.UpdateResult_UPDATE_RESULT_FAILED,
}
