-- Migration 0004 : domaines servis par le reverse proxy (ADR-0008).
-- Un nom de domaine ne vise qu'un service ; host UNIQUE porte déjà son index.

CREATE TABLE IF NOT EXISTS domains (
    id TEXT PRIMARY KEY,
    host TEXT NOT NULL UNIQUE,
    project_id TEXT NOT NULL,
    service TEXT NOT NULL,
    port INTEGER NOT NULL CHECK(port BETWEEN 1 AND 65535),
    created_by TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_domains_project ON domains(project_id, host);
