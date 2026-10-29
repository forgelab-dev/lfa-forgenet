-- Migration 0002 : déploiements (source, état du cycle de vie, chaîne de rollback).
-- request_id et previous_deployment_id valent NULL quand ils sont absents :
-- l'unicité et la clé étrangère ignorent alors la ligne.

CREATE TABLE IF NOT EXISTS deployments (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    source_type TEXT NOT NULL CHECK(source_type IN ('compose', 'git', 'image')),
    source_content TEXT NOT NULL,
    state TEXT NOT NULL CHECK(state IN ('pending', 'building', 'deploying', 'health_checking',
        'success', 'failed', 'cancelled', 'rolling_back', 'rolled_back')),
    error_message TEXT,
    target_env TEXT NOT NULL DEFAULT '',
    previous_deployment_id TEXT REFERENCES deployments(id) ON DELETE SET NULL,
    triggered_by TEXT NOT NULL,
    request_id TEXT UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_deployments_project_created ON deployments(project_id, created_at DESC);
