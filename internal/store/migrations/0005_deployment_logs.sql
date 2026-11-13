-- Migration 0005 : journal d'exécution des déploiements (sortie de docker
-- compose et étapes de ForgeNet), relu par le dashboard et diffusé en direct.
-- emitted_at est unique par déploiement (l'enregistreur l'avance d'une
-- nanoseconde au besoin) : il sert de curseur pour reprendre un flux.

CREATE TABLE IF NOT EXISTS deployment_logs (
    deployment_id TEXT NOT NULL REFERENCES deployments(id) ON DELETE CASCADE,
    emitted_at TIMESTAMP NOT NULL,
    stream TEXT NOT NULL CHECK(stream IN ('stdout', 'stderr')),
    step TEXT NOT NULL,
    message TEXT NOT NULL,
    PRIMARY KEY (deployment_id, emitted_at)
) WITHOUT ROWID;
