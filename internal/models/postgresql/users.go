package postgresql

import (
	"database/sql"

	"github.com/lazy-void/chatapp/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username, email, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3) RETURNING id;`

	var id int
	err = m.DB.QueryRow(stmt, username, email, hashedPassword).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, models.ErrNoRecord
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}
