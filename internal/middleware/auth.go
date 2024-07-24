package middleware

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/sessions"

	"github.com/syntaqx/gokku/internal/repository"
)

type contextKey string

const UserContextKey contextKey = "user"

func Auth(userRepo *repository.UsersRepository, store *sessions.CookieStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "session-name")
			userID, ok := session.Values["user_id"].(string)
			if !ok || userID == "" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			uuid, err := uuid.FromString(userID)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			user, err := userRepo.GetUserByID(uuid)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) *repository.User {
	user, ok := ctx.Value(UserContextKey).(*repository.User)
	if !ok {
		return nil
	}
	return user
}
