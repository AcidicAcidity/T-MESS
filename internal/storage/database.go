package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to init schema: %w", err)
	}

	d := &Database{db: db}

	// Создаём дефолтный чат Notes, если его нет
	if err := d.ensureNotesChat(); err != nil {
		return nil, fmt.Errorf("failed to create notes chat: %w", err)
	}

	return d, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) DB() *sql.DB {
	return d.db
}

func initSchema(db *sql.DB) error {
	schema := `
    CREATE TABLE IF NOT EXISTS chats (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        type TEXT NOT NULL,
        avatar TEXT,
        last_message TEXT,
        last_time INTEGER,
        unread_count INTEGER DEFAULT 0,
        is_active INTEGER DEFAULT 0
    );

    CREATE TABLE IF NOT EXISTS messages (
        id TEXT PRIMARY KEY,
        chat_id TEXT NOT NULL,
        sender_id TEXT NOT NULL,
        text TEXT,
        timestamp INTEGER NOT NULL,
        is_own INTEGER NOT NULL,
        signature BLOB,
        file_id TEXT,
        status TEXT DEFAULT 'sent'
    );

    CREATE TABLE IF NOT EXISTS contacts (
        peer_id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        avatar TEXT,
        public_key BLOB,
        fingerprint TEXT,
        last_seen INTEGER,
        is_online INTEGER DEFAULT 0
    );

    CREATE TABLE IF NOT EXISTS reactions (
        message_id TEXT NOT NULL,
        user_id TEXT NOT NULL,
        emoji TEXT NOT NULL,
        timestamp INTEGER NOT NULL,
        signature BLOB,
        PRIMARY KEY(message_id, user_id)
    );
	
	CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
	);

    CREATE INDEX IF NOT EXISTS idx_messages_chat ON messages(chat_id);
    CREATE INDEX IF NOT EXISTS idx_messages_time ON messages(timestamp);
    `

	_, err := db.Exec(schema)
	return err
}

// ensureNotesChat создаёт чат Notes, если его нет
func (d *Database) ensureNotesChat() error {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM chats WHERE id = 'notes'").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		now := time.Now().Unix()
		_, err := d.db.Exec(`
            INSERT INTO chats (id, name, type, avatar, last_message, last_time, unread_count, is_active)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
        `, "notes", "📝 Notes", "local", "📝", "Welcome to T-MESS! Send a message here to test.", now, 0, 1)
		if err != nil {
			return err
		}

		// Добавляем приветственное сообщение
		msgID := uuid.New().String()
		_, err = d.db.Exec(`
            INSERT INTO messages (id, chat_id, sender_id, text, timestamp, is_own, status)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        `, msgID, "notes", "system", "Welcome to T-MESS! This is your Notes chat. Use it to test messages.", now, 0, "sent")
		return err
	}
	return nil
}
