package marketplace_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/forgelab-dev/lfa-forgenet/internal/marketplace"
)

// validResolved reprend la forme de `docker compose config --no-interpolate
// --format json` pour le template PostgreSQL de l'ADR-0018.
const validResolved = `{
  "name": "forgenet-validation",
  "services": {
    "postgres": {
      "image": "postgres:17.2-alpine",
      "environment": {
        "POSTGRES_USER": "${POSTGRES_USER}",
        "POSTGRES_DB": "${POSTGRES_DB}",
        "POSTGRES_PASSWORD": "${POSTGRES_PASSWORD}",
        "PGDATA": "/var/lib/postgresql/data/pgdata"
      },
      "healthcheck": {"test": ["CMD-SHELL", "pg_isready -U $${POSTGRES_USER}"], "interval": "10s"},
      "networks": {"default": null},
      "volumes": [{"type": "volume", "source": "data", "target": "/var/lib/postgresql/data", "volume": {}}]
    }
  },
  "networks": {"default": {"name": "forgenet-validation_default"}},
  "volumes": {"data": {"name": "forgenet-validation_data"}}
}`

// mutate applique change au modèle valide et renvoie le JSON obtenu.
func mutate(t *testing.T, change func(doc map[string]any, pg map[string]any)) []byte {
	t.Helper()
	var doc map[string]any
	require.NoError(t, json.Unmarshal([]byte(validResolved), &doc))
	change(doc, doc["services"].(map[string]any)["postgres"].(map[string]any))
	out, err := json.Marshal(doc)
	require.NoError(t, err)
	return out
}

func TestCheckCompose_ValidTemplatePasses(t *testing.T) {
	assert.Empty(t, joined(marketplace.CheckCompose(parse(t, validManifest), []byte(validResolved))))
}

func TestCheckCompose_Rules(t *testing.T) {
	cases := []struct {
		name   string
		change func(doc, pg map[string]any)
		want   string
	}{
		{"image sans étiquette", func(_, pg map[string]any) { pg["image"] = "postgres" }, "ajoutez une étiquette"},
		{"image latest", func(_, pg map[string]any) { pg["image"] = "postgres:latest" }, "change sous le même nom"},
		{"étiquette sans chiffre", func(_, pg map[string]any) { pg["image"] = "postgres:alpine" }, "change sous le même nom"},
		{"image variable", func(_, pg map[string]any) { pg["image"] = "postgres:${TAG}" }, "ne dépend pas d'une variable"},
		{"version absente de l'étiquette", func(_, pg map[string]any) { pg["image"] = "postgres:16.4-alpine" }, "doit figurer dans l'étiquette"},
		{"build", func(_, pg map[string]any) { pg["build"] = map[string]any{"context": "/tmp"} }, "build est interdit"},
		{"container_name", func(_, pg map[string]any) { pg["container_name"] = "pg" }, "container_name est interdit"},
		{"env_file", func(_, pg map[string]any) { pg["env_file"] = []any{map[string]any{"path": "/tmp/.env"}} }, "env_file est interdit"},
		{"port publié", func(_, pg map[string]any) {
			pg["ports"] = []any{map[string]any{"target": 5432, "published": "5432"}}
		}, "contournant le pare-feu"},
		{"montage de l'hôte", func(_, pg map[string]any) {
			pg["volumes"] = append(pg["volumes"].([]any), map[string]any{"type": "bind", "source": "/var/run/docker.sock"})
		}, "montage de l'hôte interdit"},
		{"volume anonyme", func(_, pg map[string]any) {
			pg["volumes"] = append(pg["volumes"].([]any), map[string]any{"type": "volume", "target": "/x"})
		}, "volume anonyme"},
		{"privileged", func(_, pg map[string]any) { pg["privileged"] = true }, "privileged est interdit"},
		{"cap_add", func(_, pg map[string]any) { pg["cap_add"] = []any{"SYS_ADMIN"} }, "cap_add est interdit"},
		{"devices", func(_, pg map[string]any) { pg["devices"] = []any{map[string]any{"source": "/dev/fuse"}} }, "devices est interdit"},
		{"réseau de l'hôte", func(_, pg map[string]any) { pg["network_mode"] = "host" }, "network_mode"},
		{"pid de l'hôte", func(_, pg map[string]any) { pg["pid"] = "host" }, "pid: host est interdit"},
		{"protection désactivée", func(_, pg map[string]any) { pg["security_opt"] = []any{"seccomp:unconfined"} }, "désactive une protection"},
		{"secret en dur", func(_, pg map[string]any) {
			pg["environment"].(map[string]any)["POSTGRES_PASSWORD"] = "motdepasse"
		}, "secret en dur"},
		{"secret dans une variable non secrète", func(_, pg map[string]any) {
			pg["environment"].(map[string]any)["API_KEY"] = "${POSTGRES_DB}"
		}, "doit être marquée secret: true"},
		{"variable sans valeur", func(_, pg map[string]any) { pg["environment"].(map[string]any)["TZ"] = nil }, "sans valeur"},
		{"variable non déclarée", func(_, pg map[string]any) { pg["command"] = []any{"postgres", "-c", "max_connections=${MAX}"} }, "${MAX} n'est pas déclarée"},
		{"valeur par défaut dans le fichier", func(_, pg map[string]any) {
			pg["environment"].(map[string]any)["POSTGRES_DB"] = "${POSTGRES_DB:-app}"
		}, "seule la forme ${NOM}"},
		{"variable nue", func(_, pg map[string]any) { pg["command"] = []any{"echo", "$POSTGRES_DB"} }, "seule la forme ${NOM}"},
		{"variable du manifeste inutilisée", func(_, pg map[string]any) {
			delete(pg["environment"].(map[string]any), "POSTGRES_DB")
		}, "POSTGRES_DB ne sert nulle part"},
		{"sans sonde de santé", func(_, pg map[string]any) { delete(pg, "healthcheck") }, "sonde de santé"},
		{"sonde désactivée", func(_, pg map[string]any) { pg["healthcheck"] = map[string]any{"disable": true} }, "sonde de santé"},
		{"volume non déclaré", func(doc, pg map[string]any) {
			doc["volumes"].(map[string]any)["logs"] = map[string]any{"name": "x_logs"}
			pg["volumes"] = append(pg["volumes"].([]any), map[string]any{"type": "volume", "source": "logs", "target": "/logs"})
		}, `"logs", volume de docker-compose.yml, n'est pas déclaré`},
		{"volume external", func(doc, _ map[string]any) {
			doc["volumes"].(map[string]any)["data"] = map[string]any{"name": "data", "external": true}
		}, "volume external"},
		{"réseau external", func(doc, _ map[string]any) {
			doc["networks"].(map[string]any)["edge"] = map[string]any{"name": "forgenet-edge", "external": true}
		}, "réseau external"},
		{"fichier de configuration du serveur", func(doc, _ map[string]any) {
			doc["configs"] = map[string]any{"conf": map[string]any{"file": "/etc/passwd"}}
		}, "configs.conf : file lirait un fichier du serveur"},
		{"service principal absent", func(doc, pg map[string]any) {
			doc["services"] = map[string]any{"db": pg}
		}, `service : "postgres" n'existe pas`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := marketplace.CheckCompose(parse(t, validManifest), mutate(t, tc.change))
			assert.Contains(t, joined(errs), tc.want)
		})
	}
}

func TestCheckCompose_EnvironmentAsList(t *testing.T) {
	resolved := mutate(t, func(_, pg map[string]any) {
		pg["environment"] = []any{
			"POSTGRES_USER=${POSTGRES_USER}", "POSTGRES_DB=${POSTGRES_DB}", "POSTGRES_PASSWORD=${POSTGRES_PASSWORD}",
		}
	})
	assert.Empty(t, joined(marketplace.CheckCompose(parse(t, validManifest), resolved)))

	resolved = mutate(t, func(_, pg map[string]any) {
		pg["environment"] = []any{"POSTGRES_PASSWORD=secret", "POSTGRES_USER=${POSTGRES_USER}", "POSTGRES_DB=${POSTGRES_DB}"}
	})
	assert.Contains(t, joined(marketplace.CheckCompose(parse(t, validManifest), resolved)), "secret en dur")
}

func TestCheckCompose_DigestPinsWithoutTag(t *testing.T) {
	m := parse(t, validManifest)
	sidecar := mutate(t, func(doc, _ map[string]any) {
		doc["services"].(map[string]any)["outil"] = map[string]any{
			"image": "busybox@sha256:" + "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		}
	})
	assert.Empty(t, joined(marketplace.CheckCompose(m, sidecar)))
}

func TestCheckCompose_PortOfUnknownService(t *testing.T) {
	m := parse(t, validManifest)
	m.Ports[0].Service = "web"
	assert.Contains(t, joined(marketplace.CheckCompose(m, []byte(validResolved))), `le service "web" n'existe pas`)
}
