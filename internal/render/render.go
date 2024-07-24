package render

import (
	"net/http"

	"github.com/syntaqx/gokku/internal/middleware"
	"github.com/unrolled/render"
)

type (
	Options     = render.Options
	HTMLOptions = render.HTMLOptions
)

type Render struct {
	renderer *render.Render
}

func New(options Options) *Render {
	return &Render{renderer: render.New(options)}
}

func (r *Render) HTML(w http.ResponseWriter, req *http.Request, status int, name string, binding map[string]interface{}, options ...HTMLOptions) {
	if binding == nil {
		binding = make(map[string]interface{})
	}

	auth := middleware.GetUserFromContext(req.Context())
	if auth != nil {
		binding["Auth"] = auth
	}

	r.renderer.HTML(w, status, name, binding, options...)
}

func (r *Render) JSON(w http.ResponseWriter, req *http.Request, status int, v interface{}) {
	r.renderer.JSON(w, status, v)
}
