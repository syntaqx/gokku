package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/syntaqx/gokku/internal/render"
)

type ActionsHandler struct {
	renderer *render.Render
}

func NewActionsHandler(renderer *render.Render) *ActionsHandler {
	return &ActionsHandler{renderer: renderer}
}

func (h *ActionsHandler) RegisterRoutes(router chi.Router) {
	router.Get("/actions", h.List)
}

func (h *ActionsHandler) List(w http.ResponseWriter, r *http.Request) {
	// Stub actions data
	actions := []struct {
		ID          string
		Name        string
		Description string
		Type        string
		Hidden      bool
	}{
		{
			ID:          "action1",
			Name:        "Manual Backup",
			Description: "Perform a manual backup of the database.",
			Type:        "Manual",
			Hidden:      false,
		},
		{
			ID:          "action2",
			Name:        "Nightly Build",
			Description: "Run the nightly build pipeline.",
			Type:        "Scheduled",
			Hidden:      false,
		},
		{
			ID:          "action3",
			Name:        "Hidden Maintenance",
			Description: "Perform hidden maintenance tasks.",
			Type:        "Hidden",
			Hidden:      true,
		},
	}

	h.renderer.HTML(w, r, http.StatusOK, "actions/list", map[string]interface{}{
		"Title":   "Actions",
		"Actions": actions,
	})
}
