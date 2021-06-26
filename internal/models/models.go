// Package models provides abstractions from concrete data storage implementations.
package models

import (
	"errors"
	"time"
)

// Notable errors.
var (
	ErrNoRecord          = errors.New("models: no matching record found")
	ErrInvalidUsername   = errors.New("models: username doesn't exist")
	ErrInvalidPassword   = errors.New("models: invalid password")
	ErrDuplicateEmail    = errors.New("models: duplicate email")
	ErrDuplicateUsername = errors.New("models: duplicate username")
)

// Message represents row from the messages table.
type Message struct {
	ID       int64
	Username string
	Text     string
	Created  time.Time
}

// User represents row from the users table.
type User struct {
	Username       string
	Email          string
	HashedPassword string
	Created        time.Time
}
