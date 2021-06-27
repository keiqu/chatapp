package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/lazy-void/chatapp/chat"
	"github.com/lazy-void/chatapp/models"

	"github.com/justinas/nosurf"
)

func csrfHandler(next http.Handler) http.Handler {
	return nosurf.New(next)
}

func (app *Application) requireNonAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.authenticatedUser(r) != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := app.getUserSession(r)
		username, ok := s.Values[usernameSessionKey].(string)
		if !ok {
			app.deleteAuthCookie(w, r)
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.Users.Get(username)
		if errors.Is(err, models.ErrNoRecord) {
			app.deleteAuthCookie(w, r)
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), chat.ContextUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
