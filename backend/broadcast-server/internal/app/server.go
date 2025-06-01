package app

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "vibhordubey333/road.sh/broadcast-server/internal/config"
    "vibhordubey333/road.sh/broadcast-server/internal/domain/usecases"
    httpServer "vibhordubey333/road.sh/broadcast-server/internal/infrastructure/http"
    "vibhordubey333/road.sh/broadcast-server/internal/infrastructure/websocket"
)

// Server represents the broadcast server application
type Server struct {
    config          config.ServerConfig
    hub             *websocket.Hub
    broadcastUseCase *usecases.BroadcastUseCase
    httpServer      *httpServer.Server
}

// NewServer creates a new server
func NewServer(cfg config.ServerConfig) *Server {
    // Create hub
    hub := websocket.NewHub()
    
    // Create use case
    broadcastUseCase := usecases.NewBroadcastUseCase(hub)
    
    // Create HTTP server
    server := httpServer.NewServer(cfg.Port, broadcastUseCase)
    
    return &Server{
        config:          cfg,
        hub:             hub,
        broadcastUseCase: broadcastUseCase,
        httpServer:      server,
    }
}

// Start starts the server
func (s *Server) Start() {
    // Start hub
    go s.hub.Run()
    
    // Create a context that listens for termination signals
    ctx, cancel := context.WithCancel(context.Background())
    
    // Listen for termination signals
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
        <-sigCh
        log.Println("Received termination signal")
        cancel()
    }()
    
    // Add hub to context
    ctx = context.WithValue(ctx, "hub", s.hub)
    
    // Start HTTP server with the context
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        r = r.WithContext(ctx)
        s.httpServer.ServeHTTP(w, r)
    })
    
    // Start server
    go func() {
        if err := s.httpServer.Start(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Error starting server: %v", err)
        }
    }()
    
    // Wait for termination signal
    <-ctx.Done()
    
    // Perform graceful shutdown
    log.Println("Shutting down server...")
    
    // Give clients time to disconnect
    time.Sleep(time.Second)
    
    log.Println("Server stopped")
}