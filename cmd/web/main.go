package main

import (
	"database/sql"
	"flag"
	"os"

	"github.com/rs/zerolog"

	"github.com/lazy-void/chatapp/pkg/models/postgresql"

	"github.com/lazy-void/chatapp/pkg/server"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	dsn := flag.String("dsn", "postgresql://web:pass@localhost/chatapp", "PostgreSQL connection URI.")
	addr := flag.String("addr", ":4000", "Address that will be used by the server.")
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006/01/02 15:04:05"})

	db := initDB(*dsn)
	defer db.Close()

	s := server.Server{
		Messages: &postgresql.MessageModel{DB: db},
		Addr:     *addr,
	}

	log.Info().Msg("Starting to listen...")
	s.Start()
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
