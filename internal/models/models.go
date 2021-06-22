package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Message struct {
	ID      int64
	Text    string
	Created time.Time
}

type User struct {
	ID             int64
	Username       string
	Email          string
	HashedPassword string
	Created        time.Time
}
