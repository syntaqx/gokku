package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/syntaqx/gokku/internal/render"
)

type ApplicationsHandler struct {
	renderer *render.Render
}

func NewApplicationsHandler(renderer *render.Render) *ApplicationsHandler {
	return &ApplicationsHandler{renderer: renderer}
}

func (h *ApplicationsHandler) RegisterRoutes(router chi.Router) {
	router.Get("/applications", h.List)
}

func (h *ApplicationsHandler) List(w http.ResponseWriter, r *http.Request) {
	// Stub applications data
	applications := []struct {
		ID    string
		Name  string
		Stack string
		Tags  []string
	}{
		{
			ID:    "app1",
			Name:  "Application 1",
			Stack: "Container",
			Tags:  []string{"Tag1", "Tag2"},
		},
		{
			ID:    "app2",
			Name:  "Application 2",
			Stack: "Container",
			Tags:  []string{"Tag3"},
		},
		{
			ID:    "app3",
			Name:  "Application 3",
			Stack: "Container",
			Tags:  []string{"Tag4", "Tag5", "Tag6"},
		},
	}

	h.renderer.HTML(w, r, http.StatusOK, "applications/list", map[string]interface{}{
		"Title":        "Applications",
		"Applications": applications,
	})
}
