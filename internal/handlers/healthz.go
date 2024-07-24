package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/syntaqx/gokku/internal/render"
)

type HealthHandler struct {
	renderer *render.Render
}

func NewHealthHandler(renderer *render.Render) *HealthHandler {
	return &HealthHandler{renderer: renderer}
}

func (h *HealthHandler) RegisterRoutes(router chi.Router) {
	router.Get("/healthz", h.Health)
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.renderer.JSON(w, r, http.StatusOK, map[string]interface{}{
		"status": "OK",
	})
}
