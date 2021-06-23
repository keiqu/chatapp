package main

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/lazy-void/chatapp/internal/models"
)

func csrfHandler(next http.Handler) http.Handler {
	return nosurf.New(next)
}

func (app *application) requireNonAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.isAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, _ := app.sessions.Get(r, "user")
		defer func() {
			err := s.Save(r, w)
			if err != nil {
				app.serverError(w, err)
				return
			}

			next.ServeHTTP(w, r)
		}()

		userID, ok := s.Values["userID"].(int)
		if !ok {
			delete(s.Values, "userID")
			return
		}

		_, err := app.users.Get(userID)
		if err == models.ErrNoRecord {
			delete(s.Values, "userID")
		}
	})
}
