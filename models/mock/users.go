// Package mock provides mocks for the database methods.
package mock

import (
	"time"

	"github.com/lazy-void/chatapp/models"
)

const (
	// ValidPassword is used for passing authentication.
	ValidPassword = "qwerty123"

	// DupeUsername fails Insert method.
	DupeUsername = "dupeUser"

	// DupeEmail fails Insert method.
	DupeEmail = "dupe@google.com"
)

// UserMock is a mock of a user.
var UserMock = models.User{
	Username:       "Fenrir",
	Email:          "fenrir@google.com",
	HashedPassword: "$2a$12$6vzjkqafxBK8nFtvT83.ZuYKMCVAOa..lQDjySLQ6UIUo3m.2j.um",
	Created:        time.Now(),
}

// UserModel implements mock methods for users table.
type UserModel struct{}

// Insert mocks insert operation into the database.
func (m *UserModel) Insert(username, email, password string) error {
	switch {
	case username == DupeUsername:
		return models.ErrDuplicateUsername
	case email == DupeEmail:
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

// Authenticate mocks check for correctness of email and password.
func (m *UserModel) Authenticate(email, password string) (string, error) {
	switch {
	case email != UserMock.Email:
		return "", models.ErrNoRecord
	case password != ValidPassword:
		return "", models.ErrInvalidPassword
	default:
		return UserMock.Username, nil
	}
}

// Get mocks operation of getting users from the database.
func (m *UserModel) Get(username string) (models.User, error) {
	switch username {
	case UserMock.Username:
		return UserMock, nil
	default:
		return models.User{}, models.ErrNoRecord
	}
}
