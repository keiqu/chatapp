package main

import (
	"database/sql"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/lazy-void/chatapp/internal/models"
	"github.com/lazy-void/chatapp/internal/models/postgresql"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type application struct {
	messages interface {
		Insert(string, time.Time) (int, error)
		Get(int, int) ([]models.Message, error)
	}
	users interface {
		Insert(username, email, password string) error
		Authenticate(email, password string) (int, error)
	}
}

func main() {
	dsn := flag.String("dsn", "postgresql://web:pass@localhost/chatapp", "PostgreSQL connection URI.")
	addr := flag.String("addr", ":4000", "Address that will be used by the server.")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006/01/02 15:04:05"})

	db := initDB(*dsn)
	defer db.Close()

	app := application{
		messages: &postgresql.MessageModel{DB: db},
		users:    &postgresql.UserModel{DB: db},
	}

	log.Info().Msg("Starting to listen...")
	err := http.ListenAndServe(*addr, app.routes())
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
