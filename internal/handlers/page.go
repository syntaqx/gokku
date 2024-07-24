package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/syntaqx/gokku/internal/render"
)

type PageHandler struct {
	renderer *render.Render
}

func NewPageHandler(renderer *render.Render) *PageHandler {
	return &PageHandler{renderer: renderer}
}

func (h *PageHandler) RegisterRoutes(router chi.Router) {
	router.Get("/", h.Home)
	router.Get("/terms", h.Terms)
	router.Get("/privacy", h.Privacy)
}

func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	h.renderer.HTML(w, r, http.StatusOK, "dashboard", map[string]interface{}{
		"Title": "Home",
	})
}

func (h *PageHandler) Terms(w http.ResponseWriter, r *http.Request) {
	h.renderer.HTML(w, r, http.StatusOK, "page/terms", map[string]interface{}{
		"Title": "Terms of Service",
	})
}

func (h *PageHandler) Privacy(w http.ResponseWriter, r *http.Request) {
	h.renderer.HTML(w, r, http.StatusOK, "page/privacy", map[string]interface{}{
		"Title": "Privacy Policy",
	})
}
