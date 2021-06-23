package postgresql

import (
	"database/sql"

	"github.com/jackc/pgerrcode"

	"github.com/jackc/pgconn"

	"github.com/lazy-void/chatapp/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3);`

	_, err = m.DB.Exec(stmt, username, email, hashedPassword)
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == pgerrcode.UniqueViolation {
			if pgErr.ConstraintName == "users_email_key" {
				return models.ErrDuplicateEmail
			} else if pgErr.ConstraintName == "users_username_key" {
				return models.ErrDuplicateUsername
			}
		}
	}

	return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	stmt := `SELECT id, hashed_password FROM users WHERE email = $1;`

	var id int
	var hashedPassword string
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrNoRecord
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidPassword
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *UserModel) Get(id int) (models.User, error) {
	stmt := `SELECT id, username, email, hashed_password, created FROM users WHERE id = $1;`

	user := models.User{}
	err := m.DB.QueryRow(stmt, id).Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.Created)
	if err == sql.ErrNoRows {
		return models.User{}, models.ErrNoRecord
	} else if err != nil {
		return models.User{}, err
	}

	return user, nil
}
