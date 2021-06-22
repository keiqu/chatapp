package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lazy-void/chatapp/internal/chat"
)

func (app *application) routes() http.Handler {
	// start chat hub
	hub := chat.NewHub(app.messages)
	go hub.Run()

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Recoverer, csrfHandler)

	r.Get("/", app.home)
	r.Get("/user/signup", app.signupUserForm)
	r.Post("/user/signup", app.signupUser)
	r.Get("/user/login", app.loginUserForm)
	r.Post("/user/login", app.loginUser)

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWS(hub, w, r)
	})
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static")))
		fs.ServeHTTP(w, r)
	})

	return r
}
