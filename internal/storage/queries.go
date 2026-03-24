package storage

import (
	"database/sql"
	"time"

	"github.com/AcidicAcidity/t-mess/internal/messages"
)

// SaveMessage сохраняет сообщение в БД
func (d *Database) SaveMessage(msg *messages.Message) error {
	_, err := d.db.Exec(`
        INSERT OR REPLACE INTO messages (id, chat_id, sender_id, text, timestamp, is_own, status)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `, msg.ID, msg.ChatID, msg.SenderID, msg.Text, msg.Timestamp.Unix(), msg.IsOwn, msg.Status)
	return err
}

// GetMessages возвращает последние N сообщений для чата (от старых к новым)
func (d *Database) GetMessages(chatID string, limit int) ([]*messages.Message, error) {
	rows, err := d.db.Query(`
        SELECT id, chat_id, sender_id, text, timestamp, is_own, status
        FROM messages
        WHERE chat_id = ?
        ORDER BY timestamp ASC
        LIMIT ?
    `, chatID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*messages.Message
	for rows.Next() {
		msg := &messages.Message{}
		var ts int64
		err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Text, &ts, &msg.IsOwn, &msg.Status)
		if err != nil {
			return nil, err
		}
		msg.Timestamp = time.Unix(ts, 0)
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// GetChats возвращает список чатов (указатели)
func (d *Database) GetChats() ([]*messages.Chat, error) {
	rows, err := d.db.Query(`
        SELECT id, name, type, avatar, last_message, last_time, unread_count
        FROM chats
        ORDER BY last_time DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*messages.Chat
	for rows.Next() {
		chat := &messages.Chat{}
		var lastTime int64
		var avatar sql.NullString
		err := rows.Scan(&chat.ID, &chat.Name, &chat.Type, &avatar, &chat.LastMessage, &lastTime, &chat.UnreadCount)
		if err != nil {
			return nil, err
		}
		if avatar.Valid {
			chat.Avatar = avatar.String
		}
		chat.LastTime = time.Unix(lastTime, 0)
		chats = append(chats, chat)
	}
	return chats, nil
}

// CreateChat создаёт новый чат
func (d *Database) CreateChat(chat *messages.Chat) error {
	_, err := d.db.Exec(`
        INSERT INTO chats (id, name, type, last_message, last_time, unread_count)
        VALUES (?, ?, ?, ?, ?, ?)
    `, chat.ID, chat.Name, chat.Type, chat.LastMessage, chat.LastTime.Unix(), chat.UnreadCount)
	return err
}
