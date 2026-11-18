-- Migration 0006 : liens à usage unique (ADR-0014). Seule l'empreinte SHA-256
-- du jeton est gardée : la base volée ne donne aucun lien utilisable.
-- created_by désigne l'auteur (« user:<id> », ou « system:cli » pour une
-- commande lancée sur la machine).

CREATE TABLE IF NOT EXISTS action_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    purpose TEXT NOT NULL CHECK(purpose IN ('first_login', 'password_reset', 'invitation')),
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_by TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_action_tokens_user ON action_tokens(user_id, purpose);
