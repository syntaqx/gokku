package middleware

import "github.com/go-chi/chi/v5/middleware"

var (
	RequestID = middleware.RequestID
	RealIP    = middleware.RealIP
	Logger    = middleware.Logger
	Recoverer = middleware.Recoverer
)
