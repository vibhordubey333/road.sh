package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"

	"vibhordubey333/road.sh/broadcast-server/internal/config"
	"vibhordubey333/road.sh/broadcast-server/internal/domain/usecases"
	httpServer "vibhordubey333/road.sh/broadcast-server/internal/infrastructure/http"
	ws "vibhordubey333/road.sh/broadcast-server/internal/infrastructure/websocket"
)

// Server represents the broadcast server application
type Server struct {
	config           config.ServerConfig
	hub              *ws.Hub
	broadcastUseCase *usecases.BroadcastUseCase
	httpServer       *httpServer.Server
}

// NewServer creates a new server
func NewServer(cfg config.ServerConfig) *Server {
	// Create hub
	hub := ws.NewHub()

	// Create use case
	broadcastUseCase := usecases.NewBroadcastUseCase(hub)

	// Create HTTP server
	server := httpServer.NewServer(cfg.Port, broadcastUseCase)

	return &Server{
		config:           cfg,
		hub:              hub,
		broadcastUseCase: broadcastUseCase,
		httpServer:       server,
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

	// Create a custom HTTP server with our handler
	addr := fmt.Sprintf(":%d", s.config.Port)

	// Create a simple mux
	mux := http.NewServeMux()

	// Register the WebSocket handler directly
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the connection
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error upgrading connection: %v", err)
			return
		}

		// Get username from query params or use default
		username := r.URL.Query().Get("username")
		if username == "" {
			username = "anonymous"
		}

		log.Printf("New connection from %s", username)

		// Create a new client
		client := ws.NewWSClient(conn, s.hub, username)

		// Register client
		s.broadcastUseCase.RegisterClient(client)

		// Start client pumps
		go client.ReadPump()
		go client.WritePump()
	})

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start server
	go func() {
		log.Printf("Server started on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for termination signal
	<-ctx.Done()

	// Perform graceful shutdown
	log.Println("Shutting down server...")

	// Create a shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Shutdown the server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server stopped")
}
