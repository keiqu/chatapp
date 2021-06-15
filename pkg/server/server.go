package server

import (
	"database/sql"
	"net/http"

	"github.com/lazy-void/chatapp/pkg/chat"

	"github.com/lazy-void/chatapp/pkg/models/postgresql"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Addr     string
	DB       *sql.DB
	Messages *postgresql.MessageModel
}

func (s Server) Start() {
	staticDir := http.Dir("./ui/static")
	hub := chat.NewHub()
	ch := make(chan string)
	go hub.Run(ch)
	go s.createMessages(ch)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.home)
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/static", http.FileServer(staticDir))
		fs.ServeHTTP(w, r)
	})
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWS(hub, w, r)
	})

	err := http.ListenAndServe(s.Addr, r)
	if err != nil {
		log.Fatal().Err(err)
	}
}
