package websocket

import (
    "encoding/json"
    "log"
    "sync"
    "time"

    "github.com/google/uuid"
    "github.com/gorilla/websocket"
    
    "vibhordubey333/road.sh/broadcast-server/internal/domain/entities"
)

// WSClient represents a WebSocket client connection
type WSClient struct {
    id       string
    conn     *websocket.Conn
    send     chan []byte
    hub      *Hub
    username string
    mu       sync.Mutex
}

// NewWSClient creates a new WebSocket client
func NewWSClient(conn *websocket.Conn, hub *Hub, username string) *WSClient {
    return &WSClient{
        id:       uuid.New().String(),
        conn:     conn,
        send:     make(chan []byte, 256),
        hub:      hub,
        username: username,
    }
}

// ID returns the client's unique ID
func (c *WSClient) ID() string {
    return c.id
}

// Send sends a message to the client
func (c *WSClient) Send(message entities.Message) error {
    data, err := json.Marshal(message)
    if err != nil {
        return err
    }
    
    c.mu.Lock()
    defer c.mu.Unlock()
    
    select {
    case c.send <- data:
        return nil
    default:
        return c.Close()
    }
}

// Close closes the client connection
func (c *WSClient) Close() error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    close(c.send)
    return c.conn.Close()
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *WSClient) ReadPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    
    c.conn.SetReadLimit(4096)
    c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })
    
    for {
        _, data, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Error reading message: %v", err)
            }
            break
        }
        
        // Create a message from the raw data
        message := entities.NewMessage(c.username, string(data))
        c.hub.broadcast <- message
    }
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *WSClient) WritePump() {
    ticker := time.NewTicker(54 * time.Second)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()
    
    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if !ok {
                // Channel closed
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            
            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)
            
            // Add queued messages
            n := len(c.send)
            for i := 0; i < n; i++ {
                w.Write([]byte{'\n'})
                w.Write(<-c.send)
            }
            
            if err := w.Close(); err != nil {
                return
            }
            
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}