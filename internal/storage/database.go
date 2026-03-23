package storage

import (
	"database/sql"
	"fmt"

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

	return &Database{db: db}, nil
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

    CREATE INDEX IF NOT EXISTS idx_messages_chat ON messages(chat_id);
    CREATE INDEX IF NOT EXISTS idx_messages_time ON messages(timestamp);
    `

	_, err := db.Exec(schema)
	return err
}
