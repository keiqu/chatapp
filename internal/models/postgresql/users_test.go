package postgresql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lazy-void/chatapp/internal/models"
)

var testUserPassword = "qwerty123"

var testUser = models.User{
	Username:       "George",
	Email:          "geor@example.com",
	HashedPassword: "$2a$12$6vzjkqafxBK8nFtvT83.ZuYKMCVAOa..lQDjySLQ6UIUo3m.2j.um",
	Created:        time.Date(2021, time.June, 12, 15, 0, 0, 0, time.UTC).Local(),
}

func TestUserModel_Insert(t *testing.T) {
	if testing.Short() {
		t.Skip("postgresql: skipping integration test")
	}

	tests := []struct {
		name                      string
		username, email, password string
		wantError                 error
	}{
		{
			name:      "Correct user",
			username:  "Alice",
			email:     "alice@google.com",
			password:  "qwerty123",
			wantError: nil,
		},
		{
			name:      "Duplicate username",
			username:  testUser.Username,
			email:     "iamgeorge@google.com",
			password:  "qwerty123",
			wantError: models.ErrDuplicateUsername,
		},
		{
			name:      "Duplicate email",
			username:  "Jim",
			email:     testUser.Email,
			password:  "qwerty123",
			wantError: models.ErrDuplicateEmail,
		},
		{
			name:      "Empty username",
			username:  "",
			email:     "jim@google.com",
			password:  "qwerty123",
			wantError: nil,
		},
		{
			name:      "Empty email",
			username:  "Casey",
			email:     "",
			password:  "qwerty123",
			wantError: nil,
		},
		{
			name:      "Empty password",
			username:  "Michael",
			email:     "michael@google.com",
			password:  "",
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserModel{DB: db}

			err := m.Insert(tt.username, tt.email, tt.password)
			assert.Equal(t, tt.wantError, err)
		})
	}
}

func TestUserModel_Authenticate(t *testing.T) {
	if testing.Short() {
		t.Skip("postgresql: skipping integration test")
	}

	tests := []struct {
		name            string
		email, password string
		wantUsername    string
		wantError       error
	}{
		{
			name:         "Correct email and password",
			email:        testUser.Email,
			password:     testUserPassword,
			wantUsername: testUser.Username,
			wantError:    nil,
		},
		{
			name:         "Incorrect email",
			email:        "somerandom@example.com",
			password:     testUserPassword,
			wantUsername: "",
			wantError:    models.ErrNoRecord,
		},
		{
			name:         "Incorrect password",
			email:        testUser.Email,
			password:     "incorrect password",
			wantUsername: "",
			wantError:    models.ErrInvalidPassword,
		},
		{
			name:         "Empty email",
			email:        "",
			password:     testUserPassword,
			wantUsername: "",
			wantError:    models.ErrNoRecord,
		},
		{
			name:         "Empty password and correct email",
			email:        testUser.Email,
			password:     "",
			wantUsername: "",
			wantError:    models.ErrInvalidPassword,
		},
		{
			name:         "Empty password and incorrect email",
			email:        "somerandom@email.com",
			password:     "",
			wantUsername: "",
			wantError:    models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserModel{DB: db}

			username, err := m.Authenticate(tt.email, tt.password)

			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, tt.wantUsername, username)
		})
	}
}

func TestUserModel_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("postgresql: skipping integration test")
	}

	tests := []struct {
		name      string
		username  string
		wantUser  models.User
		wantError error
	}{
		{
			name:      "Empty username",
			username:  "",
			wantUser:  models.User{},
			wantError: models.ErrNoRecord,
		},
		{
			name:      "Existing username",
			username:  testUser.Username,
			wantUser:  testUser,
			wantError: nil,
		},
		{
			name:      "Non-existing username",
			username:  "Emma",
			wantUser:  models.User{},
			wantError: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserModel{DB: db}

			user, err := m.Get(tt.username)

			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, tt.wantUser, user)
		})
	}
}
