package controllers

import (
	"context"
	"net/http"
	"strings"

	"crosshatch/internal/database/models"
	"crosshatch/internal/services"
)

const sessionCookieName = "session"

type contextKey string

const userContextKey contextKey = "user"

var publicPaths = map[string]bool{
	"/api/auth/login":  true,
	"/api/auth/logout": true,
	"/api/auth/setup":  true,
}

type AuthMiddleware struct {
	auth *services.AuthService
}

func NewAuthMiddleware(auth *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{auth: auth}
}

// Handler requires a valid session for /api routes, letting non-API paths (the
// static frontend) and the public auth endpoints through.
func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if !strings.HasPrefix(path, "/api") || publicPaths[path] {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := m.auth.Authenticate(cookie.Value)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey, user)))
	})
}

func userFromContext(ctx context.Context) *models.User {
	user, _ := ctx.Value(userContextKey).(*models.User)
	return user
}

func sessionCookie(token string, secure bool, maxAge int) http.Cookie {
	return http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}
}
