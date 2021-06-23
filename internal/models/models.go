package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord          = errors.New("models: no matching record found")
	ErrInvalidPassword   = errors.New("models: invalid password")
	ErrDuplicateEmail    = errors.New("models: duplicate email")
	ErrDuplicateUsername = errors.New("models: duplicate username")
)

type Message struct {
	ID       int64
	Username string
	Text     string
	Created  time.Time
}

type User struct {
	ID             int64
	Username       string
	Email          string
	HashedPassword string
	Created        time.Time
}
