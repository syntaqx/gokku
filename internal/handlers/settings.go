package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/syntaqx/gokku/internal/render"
)

type SettingsHandler struct {
	renderer *render.Render
}

func NewSettingsHandler(renderer *render.Render) *SettingsHandler {
	return &SettingsHandler{renderer: renderer}
}

func (h *SettingsHandler) RegisterRoutes(router chi.Router) {
	router.Get("/settings", h.Index)
}

func (h *SettingsHandler) Index(w http.ResponseWriter, r *http.Request) {
	h.renderer.HTML(w, r, http.StatusOK, "settings/index", map[string]interface{}{
		"Title": "Settings",
	})
}
