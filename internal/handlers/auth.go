package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"

	"github.com/syntaqx/gokku/internal/render"
	"github.com/syntaqx/gokku/internal/repository"
)

type AuthHandler struct {
	renderer       *render.Render
	userRepository *repository.UsersRepository
	store          *sessions.CookieStore
}

func NewAuthHandler(renderer *render.Render, userRepository *repository.UsersRepository, store *sessions.CookieStore) *AuthHandler {
	return &AuthHandler{renderer: renderer, userRepository: userRepository, store: store}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/login", h.Login)
	r.Post("/login", h.HandleLogin)
	r.Get("/signup", h.Signup)
	r.Post("/signup", h.HandleSignup)
	r.Get("/password_reset", h.PasswordReset)
	r.Post("/password_reset", h.HandlePasswordReset)
	r.Get("/logout", h.Logout)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.renderer.HTML(w, r, http.StatusOK, "auth/login", map[string]interface{}{
		"Title": "Login",
	}, render.HTMLOptions{Layout: "_layouts/auth"})
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	user, err := h.userRepository.ValidateUser(login, password)
	if err != nil {
		h.renderer.HTML(w, r, http.StatusOK, "auth/login", map[string]interface{}{
			"Error": "Invalid login or password",
			"Title": "Login",
		}, render.HTMLOptions{Layout: "_layouts/auth"})
		return
	}

	session, _ := h.store.Get(r, "session-name")
	session.Values["user_id"] = user.ID.String()
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	h.renderer.HTML(w, r, http.StatusOK, "auth/signup", map[string]interface{}{
		"Title": "Signup",
	}, render.HTMLOptions{Layout: "_layouts/auth"})
}

func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	email := r.FormValue("email")
	password := r.FormValue("password")

	user := &repository.User{
		Login: login,
		Email: email,
	}

	if err := user.SetPassword(password); err != nil {
		h.renderer.HTML(w, r, http.StatusOK, "auth/signup", map[string]interface{}{
			"Error": "Failed to set password",
			"Title": "Signup",
		}, render.HTMLOptions{Layout: "_layouts/auth"})
		return
	}

	if err := h.userRepository.CreateUser(user); err != nil {
		h.renderer.HTML(w, r, http.StatusOK, "auth/signup", map[string]interface{}{
			"Error": "Login or email already in use",
			"Title": "Signup",
		}, render.HTMLOptions{Layout: "_layouts/auth"})
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *AuthHandler) PasswordReset(w http.ResponseWriter, r *http.Request) {
	h.renderer.HTML(w, r, http.StatusOK, "auth/password_reset", map[string]interface{}{
		"Title": "Password Reset",
	}, render.HTMLOptions{Layout: "_layouts/auth"})
}

func (h *AuthHandler) HandlePasswordReset(w http.ResponseWriter, r *http.Request) {
	h.renderer.HTML(w, r, http.StatusOK, "auth/password_reset", map[string]interface{}{
		"Error": "Password reset functionality not implemented yet",
		"Title": "Password Reset",
	}, render.HTMLOptions{Layout: "_layouts/auth"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "session-name")
	session.Values["user_id"] = ""
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusFound)
}
