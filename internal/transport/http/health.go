// Package http gère les routes et le serveur HTTP de ForgeNet.
// Ce package est une couche de transport pure, sans logique métier.
package http

import (
	"encoding/json"
	"net"
	"net/http"
	"net/netip"

	"github.com/rs/zerolog"

	"github.com/forgelab-dev/lfa-forgenet/internal/config"
	"github.com/forgelab-dev/lfa-forgenet/internal/health"
	"github.com/forgelab-dev/lfa-forgenet/internal/version"
)

// HealthResponse représente la charge utile JSON standard retournée par les endpoints de santé.
type HealthResponse struct {
	Status  string            `json:"status"`
	Details map[string]string `json:"details,omitempty"`
}

// HealthHandler regroupe les handlers HTTP d'observabilité.
type HealthHandler struct {
	healthService health.Service
	exposeDetails string
	logger        zerolog.Logger
}

// NewHealthHandler instancie le gestionnaire HTTP de santé. exposeDetails
// prend une valeur de health.expose_details (always, local ou never).
func NewHealthHandler(healthService health.Service, exposeDetails string, logger zerolog.Logger) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
		exposeDetails: exposeDetails,
		logger:        logger,
	}
}

// RegisterRoutes attache les endpoints de santé au routeur HTTP.
func (h *HealthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", h.Liveness)
	mux.HandleFunc("GET /readyz", h.Readiness)
	mux.HandleFunc("GET /version", h.Version)
}

// Liveness vérifie que le serveur est vivant et traite des requêtes.
func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	if !h.healthService.IsLive(r.Context()) {
		writeJSON(w, http.StatusServiceUnavailable, HealthResponse{
			Status: string(health.StatusNotServing),
		})
		return
	}

	writeJSON(w, http.StatusOK, HealthResponse{
		Status: string(health.StatusServing),
	})
}

// Readiness vérifie que le serveur et ses dépendances sont prêts à recevoir du
// trafic. La cause d'une indisponibilité est toujours journalisée, mais n'est
// renvoyée qu'aux appelants autorisés par health.expose_details.
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	ready, details := h.healthService.IsReady(r.Context())

	code, state := http.StatusOK, health.StatusServing
	if !ready {
		code, state = http.StatusServiceUnavailable, health.StatusNotServing
		h.logger.Warn().Interface("details", details).Msg("service non prêt")
	}

	resp := HealthResponse{Status: string(state)}
	if h.showDetails(r) {
		resp.Details = details
	}
	writeJSON(w, code, resp)
}

func (h *HealthHandler) showDetails(r *http.Request) bool {
	switch h.exposeDetails {
	case config.ExposeDetailsAlways:
		return true
	case config.ExposeDetailsLocal:
		return isDirectLocalRequest(r)
	default:
		return false
	}
}

// isDirectLocalRequest reconnaît une requête émise depuis la machine même.
// Derrière un reverse proxy local (Caddy), toutes les requêtes arrivent de
// 127.0.0.1 : une requête portant un en-tête de proxy n'est donc jamais
// considérée comme locale.
func isDirectLocalRequest(r *http.Request) bool {
	for _, header := range []string{"Forwarded", "X-Forwarded-For", "X-Real-Ip"} {
		if r.Header.Get(header) != "" {
			return false
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}
	addr, err := netip.ParseAddr(host)
	return err == nil && addr.Unmap().IsLoopback()
}

// Version retourne les métadonnées de build du binaire sous forme JSON.
func (h *HealthHandler) Version(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, version.Get())
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
