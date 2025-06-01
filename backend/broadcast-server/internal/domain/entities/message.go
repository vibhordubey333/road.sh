package entities

import "time"

// Message represents a chat message in the system
type Message struct {
    Sender    string    `json:"sender"`
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
}

// NewMessage creates a new message
func NewMessage(sender, content string) Message {
    return Message{
        Sender:    sender,
        Content:   content,
        Timestamp: time.Now(),
    }
}