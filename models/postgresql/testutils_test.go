package postgresql

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("pgx", "postgresql://test_web:pass@localhost/test_chatapp")
	if err != nil {
		t.Fatal(err)
	}

	setup, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(setup))
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		defer db.Close()

		teardown, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(teardown))
		if err != nil {
			t.Fatal(err)
		}
	}
}
