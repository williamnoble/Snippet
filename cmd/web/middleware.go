package main

import (
	"context"
	"errors"
	"github.com/justinas/nosurf"
	"net/http"
	"snippetbox/pkg/models"
)

// secureHeaders adds basic security headers to the application
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(w, r)
	})
}

//requireAuthentication checks the request context to see if the given user is authenticated. It is a wrapper to ensure
// that specified paths can only be accessed by authenticated users.
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// noSurf is used for CSRF handling, it generates a CSRF token.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := app.session.Exists(r, "authenticatedUserID")
		if !exists {
			next.ServeHTTP(w, r)
		}

		user, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return

		} else if err != nil {
			app.serverError(w, err)
		}

		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
