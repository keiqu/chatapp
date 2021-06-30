package mock

import (
	"time"

	"github.com/lazy-void/chatapp/models"
)

// MessageMock is a mock of a message.
var MessageMock = models.Message{
	ID:       1,
	Text:     "hello world",
	Username: "Fenrir",
	Created:  time.Now(),
}

// MessageModel implements mock methods for messages table.
type MessageModel struct{}

// Insert mocks insertion of message into a database.
func (m *MessageModel) Insert(text, username string, created time.Time) (int, error) {
	switch username {
	case "invalidUsername":
		return 0, models.ErrInvalidUsername
	default:
		return 2, nil
	}
}

// Latest mocks operation of getting messages from a database.
func (m *MessageModel) Latest(n, offset int) ([]models.Message, error) {
	switch {
	case n >= 1 && offset == 0:
		return []models.Message{MessageMock}, nil
	default:
		return nil, nil
	}
}
