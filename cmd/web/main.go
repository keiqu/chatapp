package main

import (
	"database/sql"
	"flag"
	"net/http"
	"os"

	"github.com/lazy-void/chatapp/internal/chat"
	"github.com/lazy-void/chatapp/internal/models/postgresql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type application struct{}

func main() {
	dsn := flag.String("dsn", "postgresql://web:pass@localhost/chatapp", "PostgreSQL connection URI.")
	addr := flag.String("addr", ":4000", "Address that will be used by the server.")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006/01/02 15:04:05"})

	db := initDB(*dsn)
	defer db.Close()

	app := application{}

	// start chat hub
	hub := chat.NewHub(&postgresql.MessageModel{DB: db})
	go hub.Run()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", app.home)
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static")))
		fs.ServeHTTP(w, r)
	})
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWS(hub, w, r)
	})

	log.Info().Msg("Starting to listen...")
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Error starting server.")
	}
}

func initDB(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Cannot open DB.")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Cannot connect to DB.")
	}

	return db
}
