package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"vibhordubey333/road.sh/broadcast-server/internal/domain/usecases"
	ws "vibhordubey333/road.sh/broadcast-server/internal/infrastructure/websocket"
)

// Server represents the HTTP server
type Server struct {
	port             int
	broadcastUseCase *usecases.BroadcastUseCase
	upgrader         websocket.Upgrader
	mux              *http.ServeMux
}

// NewServer creates a new HTTP server
func NewServer(port int, broadcastUseCase *usecases.BroadcastUseCase) *Server {
	mux := http.NewServeMux()
	
	return &Server{
		port:             port,
		broadcastUseCase: broadcastUseCase,
		mux:              mux,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections
			},
		},
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Register handlers
	s.mux.HandleFunc("/ws", s.handleWebSocket)
	
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Server started on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}

// handleWebSocket handles WebSocket connections
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	
	// Get username from query params or use default
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "anonymous"
	}
	
	// Create a new client
	hub := r.Context().Value("hub").(*ws.Hub)
	client := ws.NewWSClient(conn, hub, username)
	
	// Register client
	s.broadcastUseCase.RegisterClient(client)
	
	// Start client pumps
	go client.ReadPump()
	go client.WritePump()
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
