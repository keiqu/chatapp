// Package postgresql implements methods for manipulating tables in postgreSQL database.
package postgresql

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"

	"github.com/jackc/pgconn"

	"github.com/lazy-void/chatapp/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// UserModel implements methods for working with users table.
type UserModel struct {
	DB *sql.DB
}

// Insert adds new user to the database.
func (m *UserModel) Insert(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3);`

	_, err = m.DB.Exec(stmt, username, email, hashedPassword)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
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

// Authenticate checks for correctness provided pair of email and password.
// In case of success user's id will be returned.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	stmt := `SELECT id, hashed_password FROM users WHERE email = $1;`

	var id int
	var hashedPassword string
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, models.ErrNoRecord
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return 0, models.ErrInvalidPassword
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

// Get gets user with provided id from the database.
func (m *UserModel) Get(id int) (models.User, error) {
	stmt := `SELECT id, username, email, hashed_password, created FROM users WHERE id = $1;`

	user := models.User{}
	err := m.DB.QueryRow(stmt, id).Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.Created)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, models.ErrNoRecord
	} else if err != nil {
		return models.User{}, err
	}

	return user, nil
}
