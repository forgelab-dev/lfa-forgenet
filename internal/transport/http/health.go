// Package http gère les routes et le serveur HTTP de ForgeNet.
// Ce package est une couche de transport pure, sans logique métier.
package http

import (
	"encoding/json"
	"net/http"

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
}

// NewHealthHandler instancie le gestionnaire HTTP de santé.
func NewHealthHandler(healthService health.Service) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
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

// Readiness vérifie que le serveur et ses dépendances sont prêts à recevoir du trafic.
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	ready, details := h.healthService.IsReady(r.Context())
	if !ready {
		writeJSON(w, http.StatusServiceUnavailable, HealthResponse{
			Status:  string(health.StatusNotServing),
			Details: details,
		})
		return
	}

	writeJSON(w, http.StatusOK, HealthResponse{
		Status:  string(health.StatusServing),
		Details: details,
	})
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
