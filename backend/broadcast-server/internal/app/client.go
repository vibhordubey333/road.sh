package app

import (
    "bufio"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/gorilla/websocket"
    
    "vibhordubey333/road.sh/broadcast-server/internal/config"
    "vibhordubey333/road.sh/broadcast-server/internal/domain/entities"
)

// Client represents a client application
type Client struct {
    config config.ClientConfig
    conn   *websocket.Conn
    done   chan struct{}
}

// NewClient creates a new client
func NewClient(cfg config.ClientConfig) *Client {
    return &Client{
        config: cfg,
        done:   make(chan struct{}),
    }
}

// Connect connects to the server
func (c *Client) Connect() {
    // Connect to WebSocket server
    url := fmt.Sprintf("ws://%s/ws?username=%s", c.config.ServerAddr, c.config.Username)
    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        log.Fatalf("Error connecting to server: %v", err)
    }
    c.conn = conn
    
    // Handle termination signals
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
    
    // Handle incoming messages
    go c.receiveMessages()
    
    // Handle user input
    go c.sendMessages()
    
    // Wait for termination signal
    select {
    case <-sigCh:
        log.Println("Received termination signal")
    case <-c.done:
        log.Println("Connection closed")
    }
    
    // Close connection
    c.conn.Close()
}

// receiveMessages receives messages from the server
func (c *Client) receiveMessages() {
    defer close(c.done)
    
    for {
        _, data, err := c.conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            return
        }
        
        var message entities.Message
        if err := json.Unmarshal(data, &message); err != nil {
            log.Printf("Error unmarshaling message: %v", err)
            continue
        }
        
        fmt.Printf("[%s] %s: %s\n", message.Timestamp.Format("15:04:05"), message.Sender, message.Content)
    }
}

// sendMessages sends messages to the server
func (c *Client) sendMessages() {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("Connected to server. Type messages and press Enter to send.")

    for scanner.Scan() {
        text := scanner.Text()

        // Create a message
        message := entities.NewMessage(c.config.Username, text)

        // Marshal the message
        data, err := json.Marshal(message)
        if err != nil {
            log.Printf("Error marshaling message: %v", err)
            continue
        }

        // Send the message
        err = c.conn.WriteMessage(websocket.TextMessage, data)
        if err != nil {
            log.Printf("Error sending message: %v", err)
            return
        }
    }

    if err := scanner.Err(); err != nil {
        log.Printf("Error reading input: %v", err)
    }
}
