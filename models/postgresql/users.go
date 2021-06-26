// Package postgresql implements methods for manipulating tables in postgreSQL database.
package postgresql

import (
	"database/sql"
	"errors"

	"github.com/lazy-void/chatapp/models"

	"github.com/jackc/pgerrcode"

	"github.com/jackc/pgconn"

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
			} else if pgErr.ConstraintName == "users_pkey" {
				return models.ErrDuplicateUsername
			}
		}
	}

	return err
}

// Authenticate checks for correctness provided pair of email and password.
// In case of success username will be returned.
func (m *UserModel) Authenticate(email, password string) (string, error) {
	stmt := `SELECT username, hashed_password FROM users WHERE email = $1;`

	var username string
	var hashedPassword string
	err := m.DB.QueryRow(stmt, email).Scan(&username, &hashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return "", models.ErrNoRecord
	} else if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", models.ErrInvalidPassword
	} else if err != nil {
		return "", err
	}

	return username, nil
}

// Get gets user with provided id from the database.
func (m *UserModel) Get(username string) (models.User, error) {
	stmt := `SELECT username, email, hashed_password, created FROM users WHERE username = $1;`

	user := models.User{}
	err := m.DB.QueryRow(stmt, username).Scan(&user.Username, &user.Email, &user.HashedPassword, &user.Created)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, models.ErrNoRecord
	} else if err != nil {
		return models.User{}, err
	}

	return user, nil
}
