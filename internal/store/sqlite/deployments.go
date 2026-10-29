package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

const (
	deploymentColumns = `id, project_id, source_type, source_content, state, error_message, target_env,
		previous_deployment_id, triggered_by, request_id, created_at, updated_at, finished_at`
	selectDeployment = `SELECT ` + deploymentColumns + ` FROM deployments`

	// defaultDeploymentPageSize borne une liste de déploiements sans limite explicite.
	defaultDeploymentPageSize = 50
)

type deploymentsStore struct {
	db *sql.DB
}

func newDeploymentsStore(db *sql.DB) store.DeploymentStore {
	return &deploymentsStore{db: db}
}

func (s *deploymentsStore) CreateDeployment(ctx context.Context, d *store.Deployment) error {
	if !d.SourceType.IsValid() {
		return fmt.Errorf("type de source de déploiement inconnu : %q", d.SourceType)
	}
	if d.State == "" {
		d.State = store.DeploymentPending
	}
	if !d.State.IsValid() {
		return fmt.Errorf("état de déploiement inconnu : %q", d.State)
	}

	now := time.Now().UTC()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	if d.UpdatedAt.IsZero() {
		d.UpdatedAt = now
	}

	var finishedAt any
	if d.FinishedAt != nil {
		finishedAt = formatTime(*d.FinishedAt)
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO deployments (`+deploymentColumns+`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.ProjectID, string(d.SourceType), d.SourceContent, string(d.State), nullIfEmpty(d.ErrorMessage),
		d.TargetEnv, nullIfEmpty(d.PreviousDeploymentID), d.TriggeredBy, nullIfEmpty(d.RequestID),
		formatTime(d.CreatedAt), formatTime(d.UpdatedAt), finishedAt,
	)
	if err != nil {
		return wrapWriteError(err, "insertion du déploiement")
	}
	return nil
}

func (s *deploymentsStore) GetDeployment(ctx context.Context, id string) (*store.Deployment, error) {
	return scanDeployment(s.db.QueryRowContext(ctx, selectDeployment+` WHERE id = ?`, id))
}

func (s *deploymentsStore) GetDeploymentByRequestID(ctx context.Context, requestID string) (*store.Deployment, error) {
	// Une clé vide n'identifie rien : sans ce garde-fou, elle ne trouverait
	// aucune ligne (NULL ≠ '') mais coûterait une requête.
	if requestID == "" {
		return nil, store.ErrNotFound
	}
	return scanDeployment(s.db.QueryRowContext(ctx, selectDeployment+` WHERE request_id = ?`, requestID))
}

func (s *deploymentsStore) ListDeployments(ctx context.Context, projectID string, limit, offset int) ([]*store.Deployment, error) {
	if limit <= 0 {
		limit = defaultDeploymentPageSize
	}
	offset = max(offset, 0)

	// L'id départage deux déploiements créés à la même nanoseconde, pour une pagination stable.
	const order = ` ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?`
	if projectID == "" {
		return queryAll(ctx, s.db, "liste des déploiements", scanDeployment,
			selectDeployment+order, limit, offset)
	}
	return queryAll(ctx, s.db, "liste des déploiements", scanDeployment,
		selectDeployment+` WHERE project_id = ?`+order, projectID, limit, offset)
}

func (s *deploymentsStore) ListDeploymentsByState(ctx context.Context, state store.DeploymentState, limit int) ([]*store.Deployment, error) {
	if limit <= 0 {
		limit = defaultDeploymentPageSize
	}
	return queryAll(ctx, s.db, "liste des déploiements par état", scanDeployment,
		selectDeployment+` WHERE state = ? ORDER BY created_at ASC, id ASC LIMIT ?`, string(state), limit)
}

func (s *deploymentsStore) TransitionDeployment(ctx context.Context, id string, from, to store.DeploymentState, errMsg string) (*store.Deployment, error) {
	if !from.CanTransitionTo(to) {
		return nil, fmt.Errorf("%w : %s vers %s", store.ErrInvalidTransition, from, to)
	}

	now := time.Now().UTC()
	var finishedAt any
	if to.IsTerminal() {
		finishedAt = formatTime(now)
	}

	// La condition sur l'état courant rend le changement atomique : entre deux
	// écrivains partis du même état, un seul l'emporte.
	err := execOne(ctx, s.db, "changement d'état du déploiement",
		`UPDATE deployments SET state = ?, error_message = ?, updated_at = ?, finished_at = COALESCE(?, finished_at)
		 WHERE id = ? AND state = ?`,
		string(to), nullIfEmpty(errMsg), formatTime(now), finishedAt, id, string(from),
	)
	if errors.Is(err, store.ErrNotFound) {
		// Aucune ligne touchée : le déploiement n'existe pas, ou il a déjà quitté l'état from.
		if _, getErr := s.GetDeployment(ctx, id); getErr != nil {
			return nil, getErr
		}
		return nil, store.ErrConflict
	}
	if err != nil {
		return nil, err
	}
	return s.GetDeployment(ctx, id)
}

func scanDeployment(row scanner) (*store.Deployment, error) {
	var (
		d                             store.Deployment
		sourceType, state             string
		errMsg, previousID, requestID sql.NullString
		finishedAt                    sql.NullTime
	)
	if err := row.Scan(&d.ID, &d.ProjectID, &sourceType, &d.SourceContent, &state, &errMsg, &d.TargetEnv,
		&previousID, &d.TriggeredBy, &requestID, &d.CreatedAt, &d.UpdatedAt, &finishedAt); err != nil {
		return nil, wrapReadError(err, "lecture du déploiement")
	}
	d.SourceType = store.DeploymentSourceType(sourceType)
	d.State = store.DeploymentState(state)
	d.ErrorMessage, d.PreviousDeploymentID, d.RequestID = errMsg.String, previousID.String, requestID.String
	if finishedAt.Valid {
		t := finishedAt.Time
		d.FinishedAt = &t
	}
	return &d, nil
}

// nullIfEmpty écrit NULL pour une chaîne vide : l'unicité de request_id et la
// clé étrangère de previous_deployment_id ignorent ainsi les valeurs absentes.
func nullIfEmpty(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
