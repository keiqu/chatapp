package server

import (
	"net/http"
	"time"

	"github.com/lazy-void/chatapp/pkg/models"

	"github.com/lazy-void/chatapp/pkg/chat"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Addr     string
	Messages interface {
		Insert(text string, created time.Time) (int, error)
		Get(n, offset int) ([]models.Message, error)
	}
}

func (s Server) Start() {
	// start chat hub
	hub := chat.NewHub(s.Messages)
	go hub.Run()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.home)
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static")))
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
