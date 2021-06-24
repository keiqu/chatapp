package postgresql

import (
	"database/sql"
	"time"

	"github.com/lazy-void/chatapp/internal/models"
)

// MessageModel implements methods for working with messages table.
type MessageModel struct {
	DB *sql.DB
}

// Insert adds message to the database. In case of success id of added message is returned.
func (m *MessageModel) Insert(text string, username string, created time.Time) (int, error) {
	stmt := `INSERT INTO messages(text, user_id, created)
	VALUES($1, (SELECT id FROM users WHERE username = $2), $3)
	RETURNING id;`

	var id int
	err := m.DB.QueryRow(stmt, text, username, created).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Get sorts all messages in descending order of creation time,
// skips the first few at the given offset, and returns n objects.
// Useful for loading message history.
func (m *MessageModel) Get(n, offset int) ([]models.Message, error) {
	stmt := `SELECT m.id, username, text, m.created
	FROM messages m
	INNER JOIN users u ON m.user_id = u.id
	ORDER BY created DESC
	OFFSET $1
 	LIMIT $2;`

	rows, err := m.DB.Query(stmt, offset, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		msg := models.Message{}
		err := rows.Scan(&msg.ID, &msg.Username, &msg.Text, &msg.Created)
		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
