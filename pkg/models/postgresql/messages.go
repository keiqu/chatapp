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

func (m *MessageModel) Get(id int) (models.Message, error) {
	stmt := `SELECT id, text, created
	FROM messages
	WHERE id = $1;`

	msg := models.Message{}
	err := m.DB.QueryRow(stmt, id).Scan(&msg.ID, &msg.Text, &msg.Created)
	if err == sql.ErrNoRows {
		return models.Message{}, models.ErrNoRecord
	} else if err != nil {
		return models.Message{}, err
	}

	return msg, nil
}

func (m *MessageModel) Latest(n int) ([]models.Message, error) {
	stmt := `SELECT id, text, created
	FROM messages
	ORDER BY created
 	LIMIT $1;`

	rows, err := m.DB.Query(stmt, n)
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
