-- Migration 0003 : l'exécuteur cherche les déploiements par état, du plus ancien au plus récent.

CREATE INDEX IF NOT EXISTS idx_deployments_state_created ON deployments(state, created_at ASC);
