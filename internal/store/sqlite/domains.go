package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

const (
	domainColumns = `id, host, project_id, service, port, created_by, created_at`
	selectDomain  = `SELECT ` + domainColumns + ` FROM domains`
)

type domainsStore struct {
	db *sql.DB
}

func newDomainsStore(db *sql.DB) store.DomainStore {
	return &domainsStore{db: db}
}

func (s *domainsStore) CreateDomain(ctx context.Context, d *store.Domain) error {
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now().UTC()
	}
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO domains (`+domainColumns+`) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.Host, d.ProjectID, d.Service, d.Port, d.CreatedBy, formatTime(d.CreatedAt),
	)
	if err != nil {
		return wrapWriteError(err, "insertion du domaine")
	}
	return nil
}

func (s *domainsStore) GetDomain(ctx context.Context, id string) (*store.Domain, error) {
	return scanDomain(s.db.QueryRowContext(ctx, selectDomain+` WHERE id = ?`, id))
}

func (s *domainsStore) ListDomains(ctx context.Context, projectID string) ([]*store.Domain, error) {
	if projectID == "" {
		return queryAll(ctx, s.db, "liste des domaines", scanDomain, selectDomain+` ORDER BY host`)
	}
	return queryAll(ctx, s.db, "liste des domaines", scanDomain,
		selectDomain+` WHERE project_id = ? ORDER BY host`, projectID)
}

func (s *domainsStore) DeleteDomain(ctx context.Context, id string) error {
	return execOne(ctx, s.db, "suppression du domaine", `DELETE FROM domains WHERE id = ?`, id)
}

func scanDomain(row scanner) (*store.Domain, error) {
	var d store.Domain
	if err := row.Scan(&d.ID, &d.Host, &d.ProjectID, &d.Service, &d.Port, &d.CreatedBy, &d.CreatedAt); err != nil {
		return nil, wrapReadError(err, "lecture du domaine")
	}
	return &d, nil
}
