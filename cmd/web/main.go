package main

import (
	"database/sql"
	"flag"

	"github.com/lazy-void/chatapp/pkg/models/postgresql"

	"github.com/lazy-void/chatapp/pkg/server"
	log "github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	dsn := flag.String("dsn", "postgresql://web:pass@localhost/chatapp", "PostgreSQL connection URI.")

	db := initDB(*dsn)
	s := server.Server{
		Logger:   log.New(),
		Messages: &postgresql.MessageModel{db},
	}

	log.Println("Starting to listen...")
	s.Start()
}

func initDB(dsn string) *sql.DB {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
