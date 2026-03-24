package messages

import (
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	SenderID  string    `json:"sender_id"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	IsOwn     bool      `json:"is_own"`
	Signature []byte    `json:"signature,omitempty"`
	FileID    string    `json:"file_id,omitempty"`
	Status    string    `json:"status"` // sent, delivered, read
}

type Chat struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // direct, group, local
	Avatar      string    `json:"avatar"`
	LastMessage string    `json:"last_message"`
	LastTime    time.Time `json:"last_time"`
	UnreadCount int       `json:"unread_count"`
}
