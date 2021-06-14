package server

import (
	"database/sql"
	"net/http"

	"github.com/lazy-void/chatapp/pkg/models/postgresql"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	*log.Logger
	DB       *sql.DB
	Messages *postgresql.MessageModel
}

func (s Server) Start() {
	dir := http.Dir("./ui/static")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", s.home)
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/static", http.FileServer(dir))
		fs.ServeHTTP(w, r)
	})
	r.Post("/message/send", s.createMessage)

	err := http.ListenAndServe(":4000", r)
	if err != nil {
		log.Fatal(err)
	}
}
