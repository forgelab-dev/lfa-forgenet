package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/forgelab-dev/lfa-forgenet/internal/store"
)

type mailSettingsStore struct {
	db *sql.DB
}

func newMailSettingsStore(db *sql.DB) store.MailSettingsStore {
	return &mailSettingsStore{db: db}
}

func (s *mailSettingsStore) GetMailSettings(ctx context.Context) (*store.MailSettings, error) {
	var m store.MailSettings
	err := s.db.QueryRowContext(ctx,
		`SELECT disabled, host, port, tls, username, password_sealed, from_address, updated_by, updated_at
		 FROM mail_settings WHERE id = 1`,
	).Scan(&m.Disabled, &m.Host, &m.Port, &m.TLS, &m.Username, &m.PasswordSealed, &m.From, &m.UpdatedBy, &m.UpdatedAt)
	if err != nil {
		return nil, wrapReadError(err, "lecture des réglages d'envoi")
	}
	return &m, nil
}

func (s *mailSettingsStore) SaveMailSettings(ctx context.Context, m *store.MailSettings) error {
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = time.Now().UTC()
	}
	var sealed any
	if len(m.PasswordSealed) > 0 {
		sealed = m.PasswordSealed
	}
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO mail_settings (id, disabled, host, port, tls, username, password_sealed, from_address, updated_by, updated_at)
		 VALUES (1, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (id) DO UPDATE SET
		   disabled = excluded.disabled, host = excluded.host, port = excluded.port, tls = excluded.tls,
		   username = excluded.username, password_sealed = excluded.password_sealed,
		   from_address = excluded.from_address, updated_by = excluded.updated_by, updated_at = excluded.updated_at`,
		m.Disabled, m.Host, m.Port, m.TLS, m.Username, sealed, m.From, m.UpdatedBy, formatTime(m.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("enregistrement des réglages d'envoi : %w", err)
	}
	return nil
}

func (s *mailSettingsStore) DeleteMailSettings(ctx context.Context) (bool, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM mail_settings WHERE id = 1`)
	if err != nil {
		return false, fmt.Errorf("effacement des réglages d'envoi : %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("effacement des réglages d'envoi : %w", err)
	}
	return n > 0, nil
}
