package websocket

import (
    "log"
    "sync"

    "vibhordubey333/road.sh/broadcast-server/internal/domain/entities"
    "vibhordubey333/road.sh/broadcast-server/internal/domain/repositories"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
    clients    map[*WSClient]bool
    broadcast  chan entities.Message
    register   chan repositories.Client
    unregister chan repositories.Client
    mutex      sync.RWMutex
}

// NewHub creates a new hub
func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*WSClient]bool),
        broadcast:  make(chan entities.Message),
        register:   make(chan repositories.Client),
        unregister: make(chan repositories.Client),
    }
}

// Run starts the hub
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            wsClient, ok := client.(*WSClient)
            if !ok {
                log.Println("Failed to cast client to WSClient")
                continue
            }
            
            h.mutex.Lock()
            h.clients[wsClient] = true
            h.mutex.Unlock()
            log.Printf("New client connected. Total clients: %d", len(h.clients))
            
        case client := <-h.unregister:
            wsClient, ok := client.(*WSClient)
            if !ok {
                log.Println("Failed to cast client to WSClient")
                continue
            }
            
            h.mutex.Lock()
            if _, ok := h.clients[wsClient]; ok {
                delete(h.clients, wsClient)
            }
            h.mutex.Unlock()
            log.Printf("Client disconnected. Total clients: %d", len(h.clients))
            
        case message := <-h.broadcast:
            h.broadcastMessage(message)
        }
    }
}

// broadcastMessage sends a message to all connected clients
func (h *Hub) broadcastMessage(message entities.Message) {
    h.mutex.RLock()
    defer h.mutex.RUnlock()
    
    for client := range h.clients {
        err := client.Send(message)
        if err != nil {
            log.Printf("Error sending message to client: %v", err)
        }
    }
}

// Add adds a client to the hub
func (h *Hub) Add(client repositories.Client) error {
    h.register <- client
    return nil
}

// Remove removes a client from the hub
func (h *Hub) Remove(client repositories.Client) error {
    h.unregister <- client
    return nil
}

// Broadcast broadcasts a message to all clients
func (h *Hub) Broadcast(message entities.Message) error {
    h.broadcast <- message
    return nil
}

// Count returns the number of connected clients
func (h *Hub) Count() int {
    h.mutex.RLock()
    defer h.mutex.RUnlock()
    return len(h.clients)
}