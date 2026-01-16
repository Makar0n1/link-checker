package handler

import (
	"net/http"

	"github.com/link-tracker/shared/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health returns service health status
// GET /health
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "backlink-service",
	})
}

// Ready returns service readiness status
// GET /ready
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}
