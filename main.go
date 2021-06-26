package main

import (
	"database/sql"
	"flag"
	"net/http"
	"os"

	"github.com/lazy-void/chatapp/models/postgresql"
	"github.com/lazy-void/chatapp/server"

	"github.com/gorilla/sessions"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	dsn := flag.String("dsn", "postgresql://web:pass@localhost/chatapp", "PostgreSQL connection URI.")
	addr := flag.String("addr", ":4000", "Address that will be used by the server.")
	secret := flag.String("secret", "946IpCV9y5Vlur8YvODJEhaOY8m9J1E4", "Secret for the session manager.")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006/01/02 15:04:05"})

	db := initDB(*dsn)
	defer db.Close()

	cs := sessions.NewCookieStore([]byte(*secret))
	cs.Options.SameSite = http.SameSiteLaxMode
	app := server.Application{
		Sessions: cs,
		Messages: &postgresql.MessageModel{DB: db},
		Users:    &postgresql.UserModel{DB: db},
	}

	log.Info().Msgf("Starting to listen on %s...", *addr)
	err := http.ListenAndServe(*addr, app.NewRouter())
	if err != nil {
		log.Err(err).
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
