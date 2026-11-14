package grpc_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forgelab-dev/lfa-forgenet/internal/audit"
	"github.com/forgelab-dev/lfa-forgenet/internal/store"
	transportGRPC "github.com/forgelab-dev/lfa-forgenet/internal/transport/grpc"
	"github.com/forgelab-dev/lfa-forgenet/internal/update"
	forgenetv1 "github.com/forgelab-dev/lfa-forgenet/pkg/gen/forgenet/v1"
)

type latestSource struct {
	version string
	err     error
}

func (s latestSource) Latest(context.Context, bool) (update.Release, error) {
	if s.err != nil {
		return update.Release{}, s.err
	}
	return update.Release{Version: s.version, NotesURL: "https://github.com/forgelab-dev/lfa-forgenet/releases/tag/v" + s.version}, nil
}

type updateFixture struct {
	ts       *testServer
	client   forgenetv1.UpdateServiceClient
	launched []update.LaunchOptions
	launch   error
	dataDir  string
	auditor  audit.Auditor
}

func startUpdateServer(t *testing.T, latest string, installDir string) *updateFixture {
	t.Helper()
	f := &updateFixture{dataDir: t.TempDir()}
	checker := update.NewChecker(latestSource{version: latest}, "0.2.0", 0, zerolog.Nop())
	checker.Check(context.Background())
	f.ts = startTestServer(t, &fakeEngine{}, func(srv *grpc.Server, _ store.Store, auditor audit.Auditor) {
		f.auditor = auditor
		forgenetv1.RegisterUpdateServiceServer(srv, transportGRPC.NewUpdateServiceServer(transportGRPC.UpdateDeps{
			Checker: checker, CheckEnabled: true, Current: "0.2.0",
			InstallDir: installDir, DataDir: f.dataDir, Auditor: auditor,
			Launch: func(_ context.Context, o update.LaunchOptions) (string, error) {
				f.launched = append(f.launched, o)
				return "abc123", f.launch
			},
		}, zerolog.Nop()))
	})
	f.client = forgenetv1.NewUpdateServiceClient(f.ts.conn)
	return f
}

func (f *updateFixture) writeState(t *testing.T, st update.State) {
	t.Helper()
	data, err := json.Marshal(st)
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Join(f.dataDir, "maj"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(f.dataDir, "maj", "etat.json"), data, 0o600))
}

func TestUpdateService_AdminsOnly(t *testing.T) {
	f := startUpdateServer(t, "0.3.0", "/opt/forgenet")
	_, err := f.client.GetUpdateStatus(as(f.ts.memberToken), &forgenetv1.GetUpdateStatusRequest{})
	assert.Equal(t, codes.PermissionDenied, status.Code(err), "les versions et l'état de l'installation sont réservés aux admins")
	_, err = f.client.StartUpdate(as(f.ts.memberToken), &forgenetv1.StartUpdateRequest{})
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
	assert.Empty(t, f.launched)
}

func TestUpdateService_Status(t *testing.T) {
	f := startUpdateServer(t, "0.3.0", "/opt/forgenet")
	resp, err := f.client.GetUpdateStatus(as(f.ts.adminToken), &forgenetv1.GetUpdateStatusRequest{})
	require.NoError(t, err)
	assert.Equal(t, "0.2.0", resp.GetCurrentVersion())
	assert.Equal(t, "0.3.0", resp.GetAvailable().GetVersion())
	assert.True(t, resp.GetSupported())
	assert.NotNil(t, resp.GetCheckedAt())
	assert.Nil(t, resp.GetLastRun())

	f.writeState(t, update.State{From: "0.1.1", To: "0.2.0", Step: update.StepRollback, Result: update.ResultRolledBack,
		Reason: "la version 0.2.0 ne démarre pas", Started: time.Now().Add(-time.Hour), Finished: time.Now()})
	resp, err = f.client.GetUpdateStatus(as(f.ts.adminToken), &forgenetv1.GetUpdateStatusRequest{})
	require.NoError(t, err)
	assert.Equal(t, forgenetv1.UpdateResult_UPDATE_RESULT_ROLLED_BACK, resp.GetLastRun().GetResult())
	assert.Equal(t, forgenetv1.UpdateStep_UPDATE_STEP_ROLLBACK, resp.GetLastRun().GetStep())
	assert.Contains(t, resp.GetLastRun().GetReason(), "ne démarre pas")
}

func TestUpdateService_UpToDate(t *testing.T) {
	f := startUpdateServer(t, "0.2.0", "/opt/forgenet")
	resp, err := f.client.GetUpdateStatus(as(f.ts.adminToken), &forgenetv1.GetUpdateStatusRequest{})
	require.NoError(t, err)
	assert.Nil(t, resp.GetAvailable(), "même version : rien à annoncer")
	_, err = f.client.StartUpdate(as(f.ts.adminToken), &forgenetv1.StartUpdateRequest{})
	assert.Equal(t, codes.FailedPrecondition, status.Code(err))
	assert.Empty(t, f.launched)
}

func TestUpdateService_Start(t *testing.T) {
	f := startUpdateServer(t, "0.3.0", "/opt/forgenet")
	resp, err := f.client.StartUpdate(as(f.ts.adminToken), &forgenetv1.StartUpdateRequest{})
	require.NoError(t, err)
	assert.Equal(t, "0.3.0", resp.GetToVersion())
	require.Len(t, f.launched, 1)
	assert.Equal(t, update.LaunchOptions{InstallDir: "/opt/forgenet", To: "0.3.0"}, f.launched[0],
		"la version annoncée est épinglée, jamais « latest »")

	events, err := f.auditor.List(context.Background(), 20, 0)
	require.NoError(t, err)
	var found bool
	for _, e := range events {
		found = found || (e.Action == "update.start" && e.Resource == "forgenet")
	}
	assert.True(t, found, "la mise à jour est inscrite au journal d'audit")
}

func TestUpdateService_StartRefusals(t *testing.T) {
	t.Run("installation non standard", func(t *testing.T) {
		f := startUpdateServer(t, "0.3.0", "")
		_, err := f.client.StartUpdate(as(f.ts.adminToken), &forgenetv1.StartUpdateRequest{})
		assert.Equal(t, codes.FailedPrecondition, status.Code(err))
		assert.Empty(t, f.launched)
	})
	t.Run("déjà en cours", func(t *testing.T) {
		f := startUpdateServer(t, "0.3.0", "/opt/forgenet")
		f.writeState(t, update.State{From: "0.2.0", To: "0.3.0", Step: update.StepPull, Result: update.ResultRunning, Started: time.Now()})
		_, err := f.client.StartUpdate(as(f.ts.adminToken), &forgenetv1.StartUpdateRequest{})
		assert.Equal(t, codes.Aborted, status.Code(err))
		assert.Empty(t, f.launched)
	})
	t.Run("conteneur déjà là", func(t *testing.T) {
		f := startUpdateServer(t, "0.3.0", "/opt/forgenet")
		f.launch = update.ErrBusy
		_, err := f.client.StartUpdate(as(f.ts.adminToken), &forgenetv1.StartUpdateRequest{})
		assert.Equal(t, codes.Aborted, status.Code(err))
	})
	t.Run("échec du lancement", func(t *testing.T) {
		f := startUpdateServer(t, "0.3.0", "/opt/forgenet")
		f.launch = errors.New("docker : socket introuvable")
		_, err := f.client.StartUpdate(as(f.ts.adminToken), &forgenetv1.StartUpdateRequest{})
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.NotContains(t, err.Error(), "socket", "le détail reste dans les journaux du serveur")
	})
}

func TestUpdateService_StaleRunIsReportedAsFailed(t *testing.T) {
	f := startUpdateServer(t, "0.3.0", "/opt/forgenet")
	f.writeState(t, update.State{From: "0.2.0", To: "0.3.0", Step: update.StepHealth, Result: update.ResultRunning, Started: time.Now().Add(-2 * time.Hour)})
	resp, err := f.client.GetUpdateStatus(as(f.ts.adminToken), &forgenetv1.GetUpdateStatusRequest{})
	require.NoError(t, err)
	assert.Equal(t, forgenetv1.UpdateResult_UPDATE_RESULT_FAILED, resp.GetLastRun().GetResult())
	assert.Contains(t, resp.GetLastRun().GetReason(), "interrompue")

	_, err = f.client.StartUpdate(as(f.ts.adminToken), &forgenetv1.StartUpdateRequest{})
	require.NoError(t, err, "une mise à jour interrompue ne bloque pas la suivante")
}
