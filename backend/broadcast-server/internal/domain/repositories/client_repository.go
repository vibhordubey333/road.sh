package repositories

import "vibhordubey333/road.sh/broadcast-server/internal/domain/entities"

// ClientRepository defines the interface for client operations
type ClientRepository interface {
    Add(client Client) error
    Remove(client Client) error
    Broadcast(message entities.Message) error
    Count() int
}

// Client represents a connected client
type Client interface {
    ID() string
    Send(message entities.Message) error
    Close() error
}