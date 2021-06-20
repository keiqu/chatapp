package postgresql

import (
	"database/sql"
	"time"

	"github.com/lazy-void/chatapp/pkg/models"
)

type MessageModel struct {
	DB *sql.DB
}

func (m *MessageModel) Insert(text string, created time.Time) (int, error) {
	stmt := `INSERT INTO messages(text, created)
	VALUES($1, $2)
	RETURNING id;`

	var id int
	err := m.DB.QueryRow(stmt, text, created).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *MessageModel) Get(n, offset int) ([]models.Message, error) {
	stmt := `SELECT id, text, created
	FROM messages
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
		err := rows.Scan(&msg.ID, &msg.Text, &msg.Created)
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
