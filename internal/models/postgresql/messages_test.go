package postgresql

import (
	"testing"
	"time"

	"github.com/lazy-void/chatapp/internal/models"

	"github.com/stretchr/testify/assert"
)

var (
	firstTestMessage = models.Message{
		ID:       1,
		Username: testUser.Username,
		Text:     "Hello World",
		Created:  time.Date(2021, time.June, 13, 15, 0, 0, 0, time.UTC).Local(),
	}
	secondTestMessage = models.Message{
		ID:       2,
		Username: testUser.Username,
		Text:     "Very loooooooooooooooooooooooooooooooooooooooooooong message!",
		Created:  time.Date(2021, time.June, 13, 15, 0, 0, 0, time.UTC).Local(),
	}
)

func TestMessageModel_Insert(t *testing.T) {
	if testing.Short() {
		t.Skip("postgresql: skipping integration test")
	}

	tests := []struct {
		name           string
		text, username string
		created        time.Time
		wantError      error
	}{
		{
			name:      "Correct message",
			text:      "Message text",
			username:  testUser.Username,
			created:   time.Now(),
			wantError: nil,
		},
		{
			name:      "Empty text",
			text:      "",
			username:  testUser.Username,
			created:   time.Now(),
			wantError: nil,
		},
		{
			name:      "Empty username",
			text:      "Message text",
			username:  "",
			created:   time.Now(),
			wantError: models.ErrInvalidUsername,
		},
		{
			name:      "Zero time",
			text:      "Message text",
			username:  testUser.Username,
			created:   time.Time{},
			wantError: nil,
		},
		{
			name:      "Non-existent username",
			text:      "Message text",
			username:  "random username",
			created:   time.Now(),
			wantError: models.ErrInvalidUsername,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := MessageModel{DB: db}

			_, err := m.Insert(tt.text, tt.username, tt.created)
			assert.Equal(t, tt.wantError, err)
		})
	}
}

func TestMessageModel_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("postgresql: skipping integration test")
	}

	tests := []struct {
		name         string
		n, offset    int
		wantMessages []models.Message
		wantError    error
	}{
		{
			name:   "n is smaller than number of elements",
			n:      1,
			offset: 0,
			wantMessages: []models.Message{
				secondTestMessage,
			},
			wantError: nil,
		},
		{
			name:   "n is equal to number of elements",
			n:      2,
			offset: 0,
			wantMessages: []models.Message{
				secondTestMessage,
				firstTestMessage,
			},
			wantError: nil,
		},
		{
			name:   "n is bigger than number of elements",
			n:      3,
			offset: 0,
			wantMessages: []models.Message{
				secondTestMessage,
				firstTestMessage,
			},
			wantError: nil,
		},
		{
			name:   "Offset skips first",
			n:      2,
			offset: 1,
			wantMessages: []models.Message{
				firstTestMessage,
			},
			wantError: nil,
		},
		{
			name:         "Offset skips all",
			n:            2,
			offset:       3,
			wantMessages: nil,
			wantError:    nil,
		},
		{
			name:         "n is zero",
			n:            0,
			offset:       1,
			wantMessages: nil,
			wantError:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := MessageModel{DB: db}

			messages, err := m.Get(tt.n, tt.offset)

			assert.Equal(t, tt.wantError, err)
			assert.Equal(t, tt.wantMessages, messages)
		})
	}
}
