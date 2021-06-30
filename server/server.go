// Package server implements server logic of the application.
package server

import (
	"embed"
	"net/http"
	"time"

	"github.com/lazy-void/chatapp/chat"
	"github.com/lazy-void/chatapp/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
)

//go:embed static
var staticFiles embed.FS

//go:embed templates
var htmlTemplates embed.FS

// Application implements server logic.
type Application struct {
	Sessions sessions.Store
	Messages interface {
		Insert(text string, username string, created time.Time) (int, error)
		Latest(n int, offset int) ([]models.Message, error)
	}
	Users interface {
		Insert(username, email, password string) error
		Authenticate(email, password string) (string, error)
		Get(username string) (models.User, error)
	}
}

// NewRouter returns initialized server router.
func (app *Application) NewRouter() http.Handler {
	// start chat hub
	hub := chat.NewHub(app.Messages)
	go hub.Run()

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(csrfHandler, app.authenticate)

		r.Group(func(r chi.Router) {
			r.Use(app.requireAuthenticatedUser)

			r.Get("/", app.home)
			r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
				chat.ServeWS(hub, w, r)
			})
			r.Post("/user/logout", app.logoutUser)
		})

		r.Group(func(r chi.Router) {
			r.Use(app.requireNonAuthenticatedUser)

			r.Get("/user/signup", app.signupUserForm)
			r.Post("/user/signup", app.signupUser)
			r.Get("/user/login", app.loginUserForm)
			r.Post("/user/login", app.loginUser)
		})
	})

	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.FS(staticFiles))
		fs.ServeHTTP(w, r)
	})

	return r
}
